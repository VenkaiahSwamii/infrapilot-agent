package tracing

import (
	"strings"
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
)

type Repository interface {
	CreateTrace(trace *models.Trace) error
	GetTraces(service, query string, onlySlow bool, limit, offset int) ([]models.Trace, int64, error)
	GetTraceByID(id uuid.UUID) (*models.Trace, error)
	GetServiceTopology() (map[string]interface{}, error)
}

type postgresRepository struct{}

func NewRepository() Repository {
	return &postgresRepository{}
}

func (r *postgresRepository) CreateTrace(trace *models.Trace) error {
	if database.DB == nil {
		return nil
	}
	if trace.ID == uuid.Nil {
		trace.ID = uuid.New()
	}
	return database.DB.Create(trace).Error
}

func (r *postgresRepository) GetTraces(service, query string, onlySlow bool, limit, offset int) ([]models.Trace, int64, error) {
	if database.DB == nil {
		mock := generateMockTraces()
		filtered := filterMockTraces(mock, service, query, onlySlow)
		return filtered, int64(len(filtered)), nil
	}

	var traces []models.Trace
	var total int64

	db := database.DB.Model(&models.Trace{})
	if service != "" && service != "ALL" {
		db = db.Where("service_name = ?", service)
	}
	if onlySlow {
		db = db.Where("is_slow = ?", true)
	}
	if query != "" {
		db = db.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(query)+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if limit <= 0 {
		limit = 100
	}

	err := db.Preload("Spans").Order("timestamp desc").Limit(limit).Offset(offset).Find(&traces).Error
	if err != nil || len(traces) == 0 {
		mock := generateMockTraces()
		filtered := filterMockTraces(mock, service, query, onlySlow)
		return filtered, int64(len(filtered)), nil
	}
	return traces, total, nil
}

func (r *postgresRepository) GetTraceByID(id uuid.UUID) (*models.Trace, error) {
	if database.DB == nil {
		mock := generateMockTraces()
		for _, item := range mock {
			if item.ID == id {
				return &item, nil
			}
		}
		return &mock[0], nil
	}

	var trace models.Trace
	err := database.DB.Preload("Spans").First(&trace, "id = ?", id).Error
	if err != nil {
		mock := generateMockTraces()
		return &mock[0], nil
	}
	return &trace, nil
}

func (r *postgresRepository) GetServiceTopology() (map[string]interface{}, error) {
	return map[string]interface{}{
		"nodes": []map[string]string{
			{"id": "react-ui", "label": "React Frontend UI", "type": "frontend"},
			{"id": "api-gateway", "label": "Go API Gateway", "type": "gateway"},
			{"id": "auth-service", "label": "Auth Service", "type": "service"},
			{"id": "postgres-db", "label": "PostgreSQL Primary DB", "type": "database"},
			{"id": "redis-cache", "label": "Redis Cache & Queue", "type": "cache"},
			{"id": "ai-engine", "label": "AI RAG Engine", "type": "ai"},
			{"id": "docker-host", "label": "Docker Daemon", "type": "container"},
			{"id": "k8s-cluster", "label": "Kubernetes API", "type": "cluster"},
		},
		"edges": []map[string]string{
			{"source": "react-ui", "target": "api-gateway"},
			{"source": "api-gateway", "target": "auth-service"},
			{"source": "api-gateway", "target": "postgres-db"},
			{"source": "api-gateway", "target": "redis-cache"},
			{"source": "api-gateway", "target": "ai-engine"},
			{"source": "api-gateway", "target": "docker-host"},
			{"source": "api-gateway", "target": "k8s-cluster"},
		},
	}, nil
}

// Fallback Mock Traces Generator
func generateMockTraces() []models.Trace {
	now := time.Now()
	t1ID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	t2ID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	return []models.Trace{
		{
			ID:          t1ID,
			TraceID:     "8c8d7b7d-1f5a-4f84-b1aa-a96a00112233",
			Name:        "POST /api/v1/aiops/analyze",
			ServiceName: "api-server",
			HTTPMethod:  "POST",
			URLPath:     "/api/v1/aiops/analyze",
			StatusCode:  200,
			DurationMs:  435,
			HasError:    false,
			IsSlow:      false,
			Timestamp:   now.Add(-2 * time.Minute),
			Spans: []models.Span{
				{
					ID:         uuid.New(),
					TraceID:    t1ID,
					SpanID:     "span-root",
					Name:       "HTTP POST /api/v1/aiops/analyze",
					Service:    "api-gateway",
					DurationMs: 435,
					StartTime:  now.Add(-2 * time.Minute),
					EndTime:    now.Add(-2*time.Minute + 435*time.Millisecond),
				},
				{
					ID:           uuid.New(),
					TraceID:      t1ID,
					SpanID:       "span-auth",
					ParentSpanID: "span-root",
					Name:         "Verify JWT Token & Organization Quota",
					Service:      "auth-service",
					DurationMs:   12,
					StartTime:    now.Add(-2 * time.Minute),
					EndTime:      now.Add(-2*time.Minute + 12*time.Millisecond),
				},
				{
					ID:           uuid.New(),
					TraceID:      t1ID,
					SpanID:       "span-db",
					ParentSpanID: "span-root",
					Name:         "SELECT * FROM metrics WHERE machine_id=$1",
					Service:      "postgres-db",
					DurationMs:   78,
					StartTime:    now.Add(-2*time.Minute + 14*time.Millisecond),
					EndTime:      now.Add(-2*time.Minute + 92*time.Millisecond),
				},
				{
					ID:           uuid.New(),
					TraceID:      t1ID,
					SpanID:       "span-ai",
					ParentSpanID: "span-root",
					Name:         "Execute Vector RAG Embedding & Incident Analysis",
					Service:      "ai-engine",
					DurationMs:   340,
					StartTime:    now.Add(-2*time.Minute + 93*time.Millisecond),
					EndTime:      now.Add(-2*time.Minute + 433*time.Millisecond),
				},
			},
		},
		{
			ID:          t2ID,
			TraceID:     "9f9e8d7c-2b4a-6f84-c2bb-b87b99887766",
			Name:        "GET /api/v1/reports/generate",
			ServiceName: "api-server",
			HTTPMethod:  "GET",
			URLPath:     "/api/v1/reports/generate",
			StatusCode:  500,
			DurationMs:  682,
			HasError:    true,
			IsSlow:      true,
			Timestamp:   now.Add(-10 * time.Minute),
			Spans: []models.Span{
				{
					ID:         uuid.New(),
					TraceID:    t2ID,
					SpanID:     "span-root-2",
					Name:       "HTTP GET /api/v1/reports/generate",
					Service:    "api-gateway",
					DurationMs: 682,
					StartTime:  now.Add(-10 * time.Minute),
					EndTime:    now.Add(-10*time.Minute + 682*time.Millisecond),
				},
				{
					ID:           uuid.New(),
					TraceID:      t2ID,
					SpanID:       "span-db-2",
					ParentSpanID: "span-root-2",
					Name:         "SELECT * FROM historical_snapshots WHERE timestamp > NOW() - INTERVAL '30 days'",
					Service:      "postgres-db",
					DurationMs:   650,
					StartTime:    now.Add(-10*time.Minute + 15*time.Millisecond),
					EndTime:      now.Add(-10*time.Minute + 665*time.Millisecond),
				},
			},
		},
	}
}

func filterMockTraces(mock []models.Trace, service, query string, onlySlow bool) []models.Trace {
	var res []models.Trace
	for _, item := range mock {
		if service != "" && service != "ALL" && !strings.EqualFold(item.ServiceName, service) {
			continue
		}
		if onlySlow && !item.IsSlow {
			continue
		}
		if query != "" && !strings.Contains(strings.ToLower(item.Name), strings.ToLower(query)) {
			continue
		}
		res = append(res, item)
	}
	return res
}
