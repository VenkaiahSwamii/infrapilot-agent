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

func (r *ServerRepository) FindExistingServer(id uuid.UUID, hostname string, ipAddress string, macAddress string) (*models.Server, error) {
	if database.DB == nil {
		return nil, errors.New("database not available")
	}

	// 1. Match by ID if valid and non-nil
	if id != uuid.Nil {
		var s models.Server
		if err := database.DB.First(&s, "id = ?", id).Error; err == nil {
			return &s, nil
		}
	}

	// 2. Match by MACAddress if provided
	trimmedMAC := strings.TrimSpace(macAddress)
	if trimmedMAC != "" {
		var s models.Server
		if err := database.DB.First(&s, "mac_address = ?", trimmedMAC).Error; err == nil {
			return &s, nil
		}
	}

	// 3. Match by IP Address if provided
	trimmedIP := strings.TrimSpace(ipAddress)
	trimmedHost := strings.TrimSpace(hostname)
	if trimmedIP != "" {
		var s models.Server
		if err := database.DB.Order("CASE WHEN status = 'ONLINE' THEN 1 ELSE 2 END, last_seen DESC").
			First(&s, "ip_address = ?", trimmedIP).Error; err == nil {
			return &s, nil
		}
		// If IP is provided but does not match any existing server, do NOT fallback to hostname alone
		// as different instances (e.g., Windows host and WSL guest) can share identical hostnames.
		return nil, errors.New("no matching server found for IP")
	}

	// 4. Match by Hostname if no IP provided
	if trimmedHost != "" {
		var s models.Server
		if err := database.DB.Order("CASE WHEN status = 'ONLINE' THEN 1 ELSE 2 END, last_seen DESC").
			First(&s, "LOWER(hostname) = LOWER(?)", trimmedHost).Error; err == nil {
			return &s, nil
		}
	}

	return nil, errors.New("no matching server found")
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
	err := database.DB.Where("LOWER(hostname) NOT LIKE '%jayathi%' AND LOWER(hostname) NOT LIKE '%navya%' AND LOWER(hostname) NOT LIKE '%server01%' AND ip_address NOT IN ('192.168.1.41', '192.168.1.18', '172.22.112.255')").Order("created_at asc, id asc").Find(&servers).Error
	return servers, err
}

func (r *ServerRepository) UpdateServer(server *models.Server) error {
	if database.DB == nil {
		return nil
	}
	if strings.Contains(strings.ToLower(server.Hostname), "jayathi") || strings.Contains(strings.ToLower(server.Hostname), "navya") || strings.Contains(strings.ToLower(server.Hostname), "server01") || server.IPAddress == "192.168.1.41" || server.IPAddress == "192.168.1.18" || server.IPAddress == "172.22.112.255" || server.IPAddress == "192.168.1.133" || server.ID.String() == "8289c184-b03f-4646-9c4c-05bec18edf97" {
		return errors.New("server update blocked by policy")
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
	if id.String() == "8289c184-b03f-4646-9c4c-05bec18edf97" || id.String() == "c762ae37-0462-457c-ab49-cd6485ae2fcb" || id.String() == "29ddc34c-cee3-4109-9c5d-5c2aeaaa72df" {
		return errors.New("server heartbeat blocked by policy")
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
