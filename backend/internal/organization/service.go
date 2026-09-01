package organization

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
)

type Service struct {
	repo        Repository
	quotaEngine *QuotaEngine
	audit       *AuditService
}

func NewService(repo Repository) *Service {
	return &Service{
		repo:        repo,
		quotaEngine: NewQuotaEngine(repo),
		audit:       NewAuditService(repo),
	}
}

func (s *Service) GetQuotaEngine() *QuotaEngine {
	return s.quotaEngine
}

func (s *Service) CreateOrganization(name, slug string, ownerID uuid.UUID) (*models.Organization, error) {
	if slug == "" {
		slug = strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	}

	org := &models.Organization{
		ID:        uuid.New(),
		Name:      name,
		Slug:      slug,
		Status:    "active",
		OwnerID:   ownerID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(org); err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	// Add owner to OrganizationUser
	if ownerID != uuid.Nil && database.DB != nil {
		orgUser := &models.OrganizationUser{
			ID:             uuid.New(),
			OrganizationID: org.ID,
			UserID:         ownerID,
			Role:           "owner",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		database.DB.Create(orgUser)
	}

	// Initialize default settings and quotas
	_, _ = s.repo.GetSettings(org.ID)
	_, _ = s.repo.GetQuotas(org.ID)

	// Create initial enrollment token
	tokenStr := s.GenerateEnrollmentTokenString(org.Slug)
	token := &models.EnrollmentToken{
		ID:             uuid.New(),
		OrganizationID: org.ID.String(),
		Name:           "Default Enrollment Token",
		TokenPrefix:    "INFRA-" + strings.ToUpper(org.Slug),
		Token:          tokenStr,
		MaxUses:        0, // unlimited
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	_ = s.repo.CreateEnrollmentToken(token)

	_ = s.audit.Log(org.ID, "SYSTEM", "127.0.0.1", "Create Organization", "Organization", "Success", "Created organization "+name)

	return org, nil
}

func (s *Service) GenerateEnrollmentTokenString(slug string) string {
	bytes := make([]byte, 4)
	_, _ = rand.Read(bytes)
	hash := strings.ToUpper(hex.EncodeToString(bytes))
	return fmt.Sprintf("INFRA-%s-%s", strings.ToUpper(slug), hash)
}

func (s *Service) GetOrganization(idOrSlug string) (*models.Organization, error) {
	if orgID, err := uuid.Parse(idOrSlug); err == nil {
		return s.repo.GetByID(orgID)
	}
	return s.repo.GetBySlug(idOrSlug)
}

func (s *Service) ListOrganizations(userID string, role string) ([]models.Organization, error) {
	uID, _ := uuid.Parse(userID)
	return s.repo.List(uID, role)
}

func (s *Service) UpdateOrganization(orgID uuid.UUID, name, slug string) (*models.Organization, error) {
	org, err := s.repo.GetByID(orgID)
	if err != nil {
		return nil, err
	}

	if name != "" {
		org.Name = name
	}
	if slug != "" {
		org.Slug = slug
	}
	org.UpdatedAt = time.Now()

	if err := s.repo.Update(org); err != nil {
		return nil, err
	}

	return org, nil
}

func (s *Service) DeleteOrganization(orgID uuid.UUID) error {
	return s.repo.Delete(orgID)
}

func (s *Service) GetDashboardMetrics(orgID uuid.UUID) (map[string]interface{}, error) {
	res := map[string]interface{}{
		"organization_id":   orgID,
		"total_servers":     0,
		"online_servers":    0,
		"offline_servers":   0,
		"docker_containers": 0,
		"k8s_clusters":      0,
		"critical_alerts":   0,
		"health_score":      98,
	}

	if database.DB == nil {
		return res, nil
	}

	var mTotal, mOnline, mOffline, dCount, cCount, aCount int64

	database.DB.Model(&models.Machine{}).Where("organization_id = ?", orgID.String()).Count(&mTotal)
	database.DB.Model(&models.Machine{}).Where("organization_id = ? AND status = ?", orgID.String(), "online").Count(&mOnline)
	database.DB.Model(&models.Machine{}).Where("organization_id = ? AND status = ?", orgID.String(), "offline").Count(&mOffline)
	database.DB.Table("docker_containers").Where("organization_id = ?", orgID.String()).Count(&dCount)
	database.DB.Table("kubernetes_clusters").Where("organization_id = ?", orgID.String()).Count(&cCount)
	database.DB.Table("linux_alerts").Where("organization_id = ? AND severity = ?", orgID.String(), "critical").Count(&aCount)

	res["total_servers"] = mTotal
	res["online_servers"] = mOnline
	res["offline_servers"] = mOffline
	res["docker_containers"] = dCount
	res["k8s_clusters"] = cCount
	res["critical_alerts"] = aCount

	healthScore := 100
	if mTotal > 0 {
		healthScore = int((float64(mOnline) / float64(mTotal)) * 100.0)
	}
	res["health_score"] = healthScore

	return res, nil
}

func (s *Service) CreateInvitation(orgID uuid.UUID, email, role, invitedBy string) (*models.OrganizationInvitation, error) {
	tokenBytes := make([]byte, 16)
	_, _ = rand.Read(tokenBytes)
	tokenStr := hex.EncodeToString(tokenBytes)

	inv := &models.OrganizationInvitation{
		ID:             uuid.New(),
		OrganizationID: orgID,
		Email:          email,
		Role:           role,
		InviteToken:    tokenStr,
		InvitedBy:      invitedBy,
		ExpiresAt:      time.Now().Add(7 * 24 * time.Hour),
		Accepted:       false,
		CreatedAt:      time.Now(),
	}

	if database.DB != nil {
		database.DB.Create(inv)
	}

	_ = s.audit.Log(orgID, invitedBy, "127.0.0.1", "Invite Team Member", "User", "Success", "Invited "+email+" as "+role)

	return inv, nil
}
