package organization

import (
	"time"

	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
)

type Repository interface {
	Create(org *models.Organization) error
	GetByID(id uuid.UUID) (*models.Organization, error)
	GetBySlug(slug string) (*models.Organization, error)
	List(userID uuid.UUID, role string) ([]models.Organization, error)
	Update(org *models.Organization) error
	Delete(id uuid.UUID) error

	GetSettings(orgID uuid.UUID) (*models.OrganizationSettings, error)
	UpdateSettings(settings *models.OrganizationSettings) error

	GetQuotas(orgID uuid.UUID) (*models.OrganizationQuota, error)
	UpdateQuotas(quotas *models.OrganizationQuota) error

	GetAuditLogs(orgID uuid.UUID, limit, offset int) ([]models.OrganizationAuditLog, int64, error)
	LogAudit(entry *models.OrganizationAuditLog) error

	GetEnrollmentTokens(orgID uuid.UUID) ([]models.EnrollmentToken, error)
	CreateEnrollmentToken(token *models.EnrollmentToken) error
}

type postgresRepository struct{}

func NewRepository() Repository {
	return &postgresRepository{}
}

func (r *postgresRepository) Create(org *models.Organization) error {
	if database.DB == nil {
		return nil
	}
	return database.DB.Create(org).Error
}

func (r *postgresRepository) GetByID(id uuid.UUID) (*models.Organization, error) {
	if database.DB == nil {
		return &models.Organization{ID: id, Name: "Default Organization", Slug: "default", Status: "active"}, nil
	}
	var org models.Organization
	if err := database.DB.First(&org, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *postgresRepository) GetBySlug(slug string) (*models.Organization, error) {
	if database.DB == nil {
		return &models.Organization{Name: "Default Organization", Slug: slug, Status: "active"}, nil
	}
	var org models.Organization
	if err := database.DB.First(&org, "slug = ?", slug).Error; err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *postgresRepository) List(userID uuid.UUID, role string) ([]models.Organization, error) {
	if database.DB == nil {
		return []models.Organization{{ID: uuid.New(), Name: "Default Organization", Slug: "default", Status: "active"}}, nil
	}
	var orgs []models.Organization
	if role == "SuperAdmin" || role == "Admin" {
		if err := database.DB.Order("created_at desc").Find(&orgs).Error; err != nil {
			return nil, err
		}
	} else {
		if err := database.DB.Table("organizations").
			Joins("join organization_users on organization_users.organization_id = organizations.id").
			Where("organization_users.user_id = ? AND organizations.status != 'archived'", userID).
			Scan(&orgs).Error; err != nil {
			return nil, err
		}
	}
	return orgs, nil
}

func (r *postgresRepository) Update(org *models.Organization) error {
	if database.DB == nil {
		return nil
	}
	return database.DB.Save(org).Error
}

func (r *postgresRepository) Delete(id uuid.UUID) error {
	if database.DB == nil {
		return nil
	}
	return database.DB.Model(&models.Organization{}).Where("id = ?", id).Update("status", "archived").Error
}

func (r *postgresRepository) GetSettings(orgID uuid.UUID) (*models.OrganizationSettings, error) {
	if database.DB == nil {
		return &models.OrganizationSettings{OrganizationID: orgID, CompanyName: "InfraPilot Enterprise", Timezone: "UTC"}, nil
	}
	var settings models.OrganizationSettings
	err := database.DB.First(&settings, "organization_id = ?", orgID).Error
	if err != nil {
		settings = models.OrganizationSettings{
			OrganizationID: orgID,
			CompanyName:    "InfraPilot Enterprise",
			Timezone:       "UTC",
			Theme:          "dark",
			PrimaryColor:   "#58a6ff",
			SidebarColor:   "#0d1117",
			LicenseTier:    "enterprise",
		}
		database.DB.Create(&settings)
	}
	return &settings, nil
}

func (r *postgresRepository) UpdateSettings(settings *models.OrganizationSettings) error {
	if database.DB == nil {
		return nil
	}
	return database.DB.Where("organization_id = ?", settings.OrganizationID).Assign(settings).FirstOrCreate(settings).Error
}

func (r *postgresRepository) GetQuotas(orgID uuid.UUID) (*models.OrganizationQuota, error) {
	if database.DB == nil {
		return &models.OrganizationQuota{OrganizationID: orgID, MaxMachines: 100, MaxUsers: 25, MaxStorageGB: 500}, nil
	}
	var quota models.OrganizationQuota
	err := database.DB.First(&quota, "organization_id = ?", orgID).Error
	if err != nil {
		quota = models.OrganizationQuota{
			OrganizationID: orgID,
			MaxMachines:    100,
			MaxUsers:       25,
			MaxStorageGB:   500,
			MaxReports:     100,
			MaxAlerts:      1000,
			RetentionDays:  90,
		}
		database.DB.Create(&quota)
	}
	return &quota, nil
}

func (r *postgresRepository) UpdateQuotas(quotas *models.OrganizationQuota) error {
	if database.DB == nil {
		return nil
	}
	return database.DB.Where("organization_id = ?", quotas.OrganizationID).Assign(quotas).FirstOrCreate(quotas).Error
}

func (r *postgresRepository) GetAuditLogs(orgID uuid.UUID, limit, offset int) ([]models.OrganizationAuditLog, int64, error) {
	if database.DB == nil {
		return []models.OrganizationAuditLog{}, 0, nil
	}
	var logs []models.OrganizationAuditLog
	var total int64
	db := database.DB.Model(&models.OrganizationAuditLog{}).Where("organization_id = ?", orgID)
	db.Count(&total)
	err := db.Order("created_at desc").Limit(limit).Offset(offset).Find(&logs).Error
	return logs, total, err
}

func (r *postgresRepository) LogAudit(entry *models.OrganizationAuditLog) error {
	if database.DB == nil {
		return nil
	}
	if entry.ID == uuid.Nil {
		entry.ID = uuid.New()
	}
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now()
	}
	return database.DB.Create(entry).Error
}

func (r *postgresRepository) GetEnrollmentTokens(orgID uuid.UUID) ([]models.EnrollmentToken, error) {
	if database.DB == nil {
		return []models.EnrollmentToken{}, nil
	}
	var tokens []models.EnrollmentToken
	err := database.DB.Where("organization_id = ?", orgID.String()).Order("created_at desc").Find(&tokens).Error
	return tokens, err
}

func (r *postgresRepository) CreateEnrollmentToken(token *models.EnrollmentToken) error {
	if database.DB == nil {
		return nil
	}
	if token.ID == uuid.Nil {
		token.ID = uuid.New()
	}
	return database.DB.Create(token).Error
}
