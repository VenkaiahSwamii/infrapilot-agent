package aiops

import (
	"sync"
	"time"

	"infrapilot/backend/internal/logger"
)

// AIOpsScheduler coordinates periodic AI analysis loops
type AIOpsScheduler struct {
	anomalyEngine     *AnomalyEngine
	predictionEngine  *PredictionEngine
	forecastEngine    *ForecastEngine
	recommendEngine   *RecommendationEngine
	healthScoreEngine *HealthScoreEngine
	logger            *logger.Logger
	done              chan bool
	mu                sync.Mutex
	running           bool
}

// NewAIOpsScheduler creates a new AIOpsScheduler
func NewAIOpsScheduler(
	anomaly *AnomalyEngine,
	prediction *PredictionEngine,
	forecast *ForecastEngine,
	recommend *RecommendationEngine,
	health *HealthScoreEngine,
) *AIOpsScheduler {
	return &AIOpsScheduler{
		anomalyEngine:     anomaly,
		predictionEngine:  prediction,
		forecastEngine:    forecast,
		recommendEngine:   recommend,
		healthScoreEngine: health,
		logger:            logger.Get(),
		done:              make(chan bool),
	}
}

// Start begins periodic AIOps background monitoring
func (s *AIOpsScheduler) Start() {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return
	}
	s.running = true
	s.mu.Unlock()

	s.logger.Info("Starting AIOps Background Analysis Engine")
	go s.runLoop()
}

// Stop halts the scheduler
func (s *AIOpsScheduler) Stop() {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}
	s.running = false
	s.mu.Unlock()

	close(s.done)
	s.logger.Info("AIOps Background Engine stopped")
}

func (s *AIOpsScheduler) runLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.logger.Info("Executing AIOps predictive analysis tick")
			_, _ = s.anomalyEngine.DetectAnomaliesForMachine("m-srv-01", "default")
			_, _ = s.predictionEngine.PredictFailures("default")
		case <-s.done:
			return
		}
	}
}
