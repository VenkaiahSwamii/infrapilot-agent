package pipeline

import (
	"time"

	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Enqueue inserts a serialized metrics payload into the PostgreSQL queue.
func Enqueue(db *gorm.DB, payload string) error {
	queued := models.QueuedMetric{
		ID:        uuid.New(),
		Payload:   payload,
		CreatedAt: time.Now(),
	}
	return db.Create(&queued).Error
}

// DequeueBatch retrieves and locks a batch of queued metrics using SKIP LOCKED.
// It returns the metrics batch and the database transaction so the caller can
// commit or rollback after processing.
func DequeueBatch(db *gorm.DB, batchSize int) ([]models.QueuedMetric, *gorm.DB, error) {
	tx := db.Begin()
	if tx.Error != nil {
		return nil, nil, tx.Error
	}

	var batch []models.QueuedMetric
	err := tx.Raw("SELECT * FROM queued_metrics ORDER BY created_at ASC LIMIT ? FOR UPDATE SKIP LOCKED", batchSize).Scan(&batch).Error
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	return batch, tx, nil
}
