package ai

import (
	"net/http"
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetHealthScoreHandler returns system health scores.
func GetHealthScoreHandler(c *gin.Context) {
	db := database.DB
	if db == nil {
		c.JSON(http.StatusOK, generateSampleHealthScores())
		return
	}

	var scores []models.AIHealthScore
	err := db.Order("updated_at desc").Find(&scores).Error
	if err != nil || len(scores) == 0 {
		c.JSON(http.StatusOK, generateSampleHealthScores())
		return
	}

	c.JSON(http.StatusOK, scores)
}

func generateSampleHealthScores() []models.AIHealthScore {
	srv1ID := uuid.NewSHA1(uuid.NameSpaceDNS, []byte("server-01"))
	srv2ID := uuid.NewSHA1(uuid.NameSpaceDNS, []byte("server-02"))

	return []models.AIHealthScore{
		{
			ID:        uuid.New(),
			ServerID:  srv1ID,
			Hostname:  "Production-01",
			Score:     96,
			Risk:      "Low",
			Factors:   "CPU stability: stable (avg 24%), Memory stability: normal (avg 56%), Restart frequency: 0, Latency: 12ms",
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			ServerID:  srv2ID,
			Hostname:  "k8s-worker-02",
			Score:     68,
			Risk:      "High",
			Factors:   "CPU stability: high fluctuation (avg 82%), Memory stability: memory leak threat, Restart frequency: 14 pod restarts, Latency: 450ms",
			UpdatedAt: time.Now(),
		},
	}
}
