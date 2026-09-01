package cache

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type MetricsQueue struct {
	mu       sync.Mutex
	filePath string
	items    []interface{}
}

var GlobalQueue *MetricsQueue

func Init(filePath string) {
	GlobalQueue = &MetricsQueue{
		filePath: filePath,
		items:    make([]interface{}, 0),
	}
	if err := GlobalQueue.load(); err != nil {
		log.Printf("[Cache] Failed to load offline queue: %v. Initializing empty queue.", err)
	}
}

func (q *MetricsQueue) load() error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if _, err := os.Stat(q.filePath); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(q.filePath)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}

	return json.Unmarshal(data, &q.items)
}

func (q *MetricsQueue) save() error {
	data, err := json.MarshalIndent(q.items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(q.filePath, data, 0600)
}

func QueueMetric(item interface{}) {
	if GlobalQueue == nil {
		return
	}
	GlobalQueue.mu.Lock()
	defer GlobalQueue.mu.Unlock()

	GlobalQueue.items = append(GlobalQueue.items, item)
	if err := GlobalQueue.save(); err != nil {
		log.Printf("[Cache] Failed to save queued metric: %v", err)
	} else {
		log.Printf("[Cache] Metric queued offline. Total queued: %d", len(GlobalQueue.items))
	}
}

func PopMetric() (interface{}, bool) {
	if GlobalQueue == nil {
		return nil, false
	}
	GlobalQueue.mu.Lock()
	defer GlobalQueue.mu.Unlock()

	if len(GlobalQueue.items) == 0 {
		return nil, false
	}

	item := GlobalQueue.items[0]
	GlobalQueue.items = GlobalQueue.items[1:]

	if err := GlobalQueue.save(); err != nil {
		log.Printf("[Cache] Failed to save queue after pop: %v", err)
	}
	return item, true
}

func QueueSize() int {
	if GlobalQueue == nil {
		return 0
	}
	GlobalQueue.mu.Lock()
	defer GlobalQueue.mu.Unlock()
	return len(GlobalQueue.items)
}
