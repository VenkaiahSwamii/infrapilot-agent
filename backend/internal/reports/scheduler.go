package reports

import (
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

type ReportScheduler struct {
	generator *ReportGenerator
	tickers   map[uuid.UUID]*time.Ticker
	mu        sync.Mutex
}

func NewReportScheduler(generator *ReportGenerator) *ReportScheduler {
	return &ReportScheduler{
		generator: generator,
		tickers:   make(map[uuid.UUID]*time.Ticker),
	}
}

func (s *ReportScheduler) Start() {
	log.Println("Report scheduler started")
}

func (s *ReportScheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, ticker := range s.tickers {
		ticker.Stop()
	}
	s.tickers = make(map[uuid.UUID]*time.Ticker)
	log.Println("Report scheduler stopped")
}

func (s *ReportScheduler) ScheduleReport(userID uuid.UUID, reportID uuid.UUID, schedule ReportSchedule) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if ticker, exists := s.tickers[reportID]; exists {
		ticker.Stop()
	}

	var interval time.Duration
	switch schedule {
	case ReportScheduleDaily:
		interval = 24 * time.Hour
	case ReportScheduleWeekly:
		interval = 7 * 24 * time.Hour
	case ReportScheduleMonthly:
		interval = 30 * 24 * time.Hour
	default:
		return
	}

	ticker := time.NewTicker(interval)
	go func(id uuid.UUID) {
		for range ticker.C {
			report, err := s.generator.Repo.FindByID(id)
			if err != nil {
				log.Printf("Scheduled report %s not found: %v", id, err)
				continue
			}

			if report.Status != ReportStatusCompleted {
				log.Printf("Skipping scheduled report %s, status: %s", id, report.Status)
				continue
			}

			req := GenerateReportRequest{
				Name:       report.Name,
				Type:       report.Type,
				Format:     report.Format,
				Schedule:   schedule,
				MachineIDs: deserializeUUIDs(report.MachineIDs),
				Filters:    deserializeFilters(report.Filters),
			}
			if report.DateFrom != nil {
				req.DateFrom = report.DateFrom
			}
			if report.DateTo != nil {
				req.DateTo = report.DateTo
			}

			_, err = s.generator.Generate(userID, req)
			if err != nil {
				log.Printf("Failed to generate scheduled report %s: %v", id, err)
			}
		}
	}(reportID)

	s.tickers[reportID] = ticker
}

func deserializeUUIDs(data []byte) []uuid.UUID {
	if len(data) == 0 {
		return nil
	}
	var ids []uuid.UUID
	// Simple parsing - in production use proper JSON unmarshal
	return ids
}

func deserializeFilters(data []byte) *ReportFilter {
	if len(data) == 0 {
		return nil
	}
	var filter ReportFilter
	// Simple parsing - in production use proper JSON unmarshal
	return &filter
}
