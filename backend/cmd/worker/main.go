package main

import (
	"encoding/json"
	"log"
	"time"

	"infrapilot/backend/internal/config"
	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/events"
	"infrapilot/backend/internal/pipeline"
	"infrapilot/backend/internal/repository"
	"infrapilot/backend/internal/services"
)

func main() {
	log.Println("Starting InfraPilot Ingestion and Alerting Worker...")

	config.LoadEnv()
	database.Connect()

	// Repositories
	machineRepo := repository.NewMachineRepository()
	metricRepo := repository.NewMetricRepository()

	// Services
	eventBus := events.NewEventBus()
	metricService := services.NewMetricService(metricRepo, machineRepo)
	alertEngine := services.NewAlertEngine(eventBus)

	log.Println("Worker initialized. Starting processing loop...")

	for {
		// Dequeue a batch of up to 50 metrics
		batch, tx, err := pipeline.DequeueBatch(database.DB, 50)
		if err != nil {
			log.Printf("Worker Dequeue Error: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		if len(batch) == 0 {
			_ = tx.Rollback()
			time.Sleep(1 * time.Second)
			continue
		}

		log.Printf("Processing batch of %d queued metrics...", len(batch))

		for _, item := range batch {
			var input services.SaveMetricInput
			if err := json.Unmarshal([]byte(item.Payload), &input); err != nil {
				log.Printf("Worker payload unmarshal error: %v", err)
				_ = tx.Exec("DELETE FROM queued_metrics WHERE id = ?", item.ID)
				continue
			}

			_, machine, err := metricService.SaveMetric(input)
			if err != nil {
				log.Printf("Worker SaveMetric error: %v", err)
				_ = tx.Exec("DELETE FROM queued_metrics WHERE id = ?", item.ID)
				continue
			}

			// Evaluate thresholds (CPU, RAM, Disk)
			alertEngine.Evaluate(machine.ID.String(), input.CPUUsage, input.MemoryPercent, input.DiskPercent, input.LatencyMs, input.PacketLoss)

			// Delete processed item from the queue
			if err := tx.Exec("DELETE FROM queued_metrics WHERE id = ?", item.ID); err != nil {
				log.Printf("Worker queue delete error: %v", err)
			}
		}

		if err := tx.Commit().Error; err != nil {
			log.Printf("Worker commit transaction error: %v", err)
		}

		time.Sleep(100 * time.Millisecond)
	}
}
