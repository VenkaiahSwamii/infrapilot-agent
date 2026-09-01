package workers

import (
	"context"
	"fmt"
	"sync"
	"time"

	"infrapilot/backend/internal/queue"
)

// Task represents a worker task
type Task func(ctx context.Context, msg *queue.Message) error

// WorkerPool manages a pool of workers
type WorkerPool struct {
	name          string
	task          Task
	numWorkers    int
	taskQueue     chan *queue.Message
	wg            sync.WaitGroup
	stopChan      chan struct{}
	queue         *queue.Queue
	streamName    string
	consumerGroup string
	running       bool
	mu            sync.RWMutex
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(name string, task Task, numWorkers, queueSize int, q *queue.Queue, streamName, consumerGroup string) *WorkerPool {
	return &WorkerPool{
		name:          name,
		task:          task,
		numWorkers:    numWorkers,
		taskQueue:     make(chan *queue.Message, queueSize),
		stopChan:      make(chan struct{}),
		queue:         q,
		streamName:    streamName,
		consumerGroup: consumerGroup,
	}
}

// Start starts the worker pool
func (wp *WorkerPool) Start() error {
	wp.mu.Lock()
	if wp.running {
		wp.mu.Unlock()
		return fmt.Errorf("worker pool %s is already running", wp.name)
	}
	wp.running = true
	wp.mu.Unlock()

	// Create consumer group
	if err := wp.queue.CreateConsumerGroup(wp.streamName, wp.consumerGroup); err != nil {
		return fmt.Errorf("failed to create consumer group: %w", err)
	}

	// Start workers
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}

	// Start monitor goroutine
	go wp.monitor()

	fmt.Printf("Worker pool %s started with %d workers\n", wp.name, wp.numWorkers)
	return nil
}

// Stop stops the worker pool gracefully
func (wp *WorkerPool) Stop() {
	wp.mu.Lock()
	if !wp.running {
		wp.mu.Unlock()
		return
	}
	wp.running = false
	wp.mu.Unlock()

	close(wp.stopChan)
	wp.wg.Wait()

	fmt.Printf("Worker pool %s stopped\n", wp.name)
}

// worker processes tasks from the queue
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	for {
		select {
		case <-wp.stopChan:
			return
		default:
			// Consume messages from stream
			messages, err := wp.queue.Consume(wp.streamName, wp.consumerGroup, fmt.Sprintf("worker-%d", id), 10)
			if err != nil {
				fmt.Printf("Worker %s-%d error consuming: %v\n", wp.name, id, err)
				time.Sleep(time.Second)
				continue
			}

			for _, msg := range messages {
				if err := wp.task(context.Background(), msg); err != nil {
					fmt.Printf("Worker %s-%d error processing message %s: %v\n", wp.name, id, msg.ID, err)

					// Increment retry count
					msg.Retries++
					if msg.Retries < msg.MaxRetries {
						// Re-publish for retry
						if pubErr := wp.queue.Publish(wp.streamName, msg); pubErr != nil {
							fmt.Printf("Worker %s-%d error re-publishing message %s: %v\n", wp.name, id, msg.ID, pubErr)
						}
					} else {
						// Send to DLQ
						if dlqErr := wp.queue.PublishToDLQ(wp.streamName, msg); dlqErr != nil {
							fmt.Printf("Worker %s-%d error sending message %s to DLQ: %v\n", wp.name, id, msg.ID, dlqErr)
						}
					}
				}

				// Acknowledge message
				if ackErr := wp.queue.Ack(wp.streamName, wp.consumerGroup, msg.ID); ackErr != nil {
					fmt.Printf("Worker %s-%d error acknowledging message %s: %v\n", wp.name, id, msg.ID, ackErr)
				}
			}
		}
	}
}

// monitor monitors the worker pool health
func (wp *WorkerPool) monitor() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Check stream length
			length, err := wp.queue.GetStreamLength(wp.streamName)
			if err != nil {
				fmt.Printf("Worker pool %s monitor error: %v\n", wp.name, err)
				continue
			}

			if length > 1000 {
				fmt.Printf("Worker pool %s has high queue length: %d\n", wp.name, length)
			}
		case <-wp.stopChan:
			return
		}
	}
}

// IsRunning returns whether the worker pool is running
func (wp *WorkerPool) IsRunning() bool {
	wp.mu.RLock()
	defer wp.mu.RUnlock()
	return wp.running
}

// Submit submits a task to the worker pool (in-memory queue)
func (wp *WorkerPool) Submit(msg *queue.Message) error {
	select {
	case wp.taskQueue <- msg:
		return nil
	default:
		return fmt.Errorf("worker pool %s task queue is full", wp.name)
	}
}

// Metrics returns worker pool metrics
type WorkerMetrics struct {
	Name        string
	NumWorkers  int
	Running     bool
	QueueLength int64
}

// GetMetrics returns metrics for the worker pool
func (wp *WorkerPool) GetMetrics() WorkerMetrics {
	wp.mu.RLock()
	defer wp.mu.RUnlock()

	length, _ := wp.queue.GetStreamLength(wp.streamName)

	return WorkerMetrics{
		Name:        wp.name,
		NumWorkers:  wp.numWorkers,
		Running:     wp.running,
		QueueLength: length,
	}
}
