package repository

import (
	"errors"
	"strings"
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
)

type ServerRepository struct{}

func NewServerRepository() *ServerRepository {
	return &ServerRepository{}
}

func (r *ServerRepository) CreateServer(server *models.Server) error {
	if database.DB == nil {
		return nil
	}
	return database.DB.Create(server).Error
}

func (r *ServerRepository) GetServerByHostname(hostname string) (*models.Server, error) {
	if database.DB == nil {
		return nil, errors.New("record not found")
	}
	var server models.Server
	err := database.DB.First(&server, "LOWER(hostname) = LOWER(?)", strings.TrimSpace(hostname)).Error
	if err != nil {
		return nil, err
	}
	return &server, nil
}

func (r *ServerRepository) GetServer(id uuid.UUID) (*models.Server, error) {
	if database.DB == nil {
		return nil, errors.New("record not found")
	}
	var server models.Server
	err := database.DB.First(&server, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &server, nil
}

func (r *ServerRepository) GetServerByIDOrHostname(identifier string) (*models.Server, error) {
	if database.DB == nil {
		if id, err := uuid.Parse(identifier); err == nil {
			return &models.Server{ID: id, Hostname: identifier}, nil
		}
		return &models.Server{ID: uuid.New(), Hostname: identifier}, nil
	}
	identifier = strings.TrimSpace(identifier)
	if identifier == "" {
		return nil, errors.New("empty identifier")
	}

	var server models.Server
	if id, err := uuid.Parse(identifier); err == nil {
		if err := database.DB.First(&server, "id = ?", id).Error; err == nil {
			return &server, nil
		}
	}

	err := database.DB.First(&server, "LOWER(hostname) = LOWER(?) OR LOWER(name) = LOWER(?) OR id::text LIKE ? OR ip_address = ?", identifier, identifier, identifier+"%", identifier).Error
	if err != nil {
		return nil, err
	}
	return &server, nil
}

func (r *ServerRepository) ListServers() ([]models.Server, error) {
	if database.DB == nil {
		return []models.Server{}, nil
	}
	var servers []models.Server
	err := database.DB.Order("created_at asc, id asc").Find(&servers).Error
	return servers, err
}

func (r *ServerRepository) UpdateServer(server *models.Server) error {
	if database.DB == nil {
		return nil
	}
	return database.DB.Save(server).Error
}

func (r *ServerRepository) DeleteServer(id uuid.UUID) error {
	if database.DB == nil {
		return nil
	}
	return database.DB.Delete(&models.Server{}, "id = ?", id).Error
}

func (r *ServerRepository) UpdateHeartbeat(id uuid.UUID, lastSeen time.Time) error {
	if database.DB == nil {
		return nil
	}
	return database.DB.Model(&models.Server{}).Where("id = ?", id).Updates(map[string]interface{}{
		"last_seen": lastSeen,
		"status":    "ONLINE",
		"online":    true,
	}).Error
}

func (r *ServerRepository) UpdateStatus(id uuid.UUID, status string, online bool) error {
	if database.DB == nil {
		return nil
	}
	return database.DB.Model(&models.Server{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status": status,
		"online": online,
	}).Error
}

func (r *ServerRepository) RotateAPIKey(id uuid.UUID, apiKey string, version int) error {
	if database.DB == nil {
		return nil
	}
	return database.DB.Model(&models.Server{}).Where("id = ?", id).Updates(map[string]interface{}{
		"api_key":         apiKey,
		"key_version":     version,
		"last_key_rotate": time.Now(),
	}).Error
}
