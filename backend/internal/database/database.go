package database

import (
	"time"

	"infrapilot/backend/internal/config"
	"infrapilot/backend/internal/logger"
	"infrapilot/backend/internal/models"
	"infrapilot/backend/internal/search"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect initializes the database connection
func Connect() {
	cfg := config.Get()
	dsn := cfg.GetDSN()

	logger.Info("Connecting to PostgreSQL",
		"host", cfg.DBHost,
		"port", cfg.DBPort,
		"database", cfg.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Info),
	})

	if err != nil {
		logger.Fatal("Database connection failed",
			"error", err,
		)
	}

	// Pre-create servers table if it doesn't exist
	if !db.Migrator().HasTable("servers") {
		logger.Info("Creating servers table...")
		err := db.Exec(`
			CREATE TABLE servers (
				id uuid PRIMARY KEY,
				name varchar(255),
				hostname varchar(255),
				ip_address varchar(255),
				os varchar(255),
				operating_system varchar(255),
				platform varchar(255),
				agent_version varchar(255),
				resource_type varchar(255) DEFAULT 'windows',
				organization varchar(255) DEFAULT 'Default Organization',
				kernel varchar(255),
				architecture varchar(255),
				mac_address varchar(255),
				cpu_model varchar(255),
				total_memory_gb bigint,
				total_disk_gb bigint,
				gpu varchar(255),
				virtualization varchar(255),
				cloud_provider varchar(255),
				status varchar(255) DEFAULT 'ONLINE',
				online boolean,
				health_score numeric DEFAULT 100,
				last_seen timestamp with time zone,
				retry_count bigint DEFAULT 0,
				api_key varchar(255) UNIQUE,
				key_version bigint DEFAULT 1,
				last_key_rotate timestamp with time zone,
				created_at timestamp with time zone,
				updated_at timestamp with time zone
			);
		`).Error
		if err != nil {
			logger.Error("Failed to create servers table", "error", err)
		}
	}

	// Always pre-populate empty servers table from machines table if machines has records
	if db.Migrator().HasTable("servers") && db.Migrator().HasTable("machines") {
		var serverCount int64
		db.Table("servers").Count(&serverCount)
		if serverCount == 0 {
			var machineCount int64
			db.Table("machines").Count(&machineCount)
			if machineCount > 0 {
				logger.Info("Pre-populating empty servers table from machines table...")
				err := db.Exec(`
					INSERT INTO servers (
						id, name, hostname, ip_address, os, operating_system, platform, agent_version,
						resource_type, organization, kernel, architecture, mac_address, cpu_model,
						total_memory_gb, total_disk_gb, gpu, virtualization, cloud_provider,
						status, online, health_score, last_seen, retry_count, api_key, key_version,
						last_key_rotate, created_at, updated_at
					)
					SELECT 
						id, name, hostname, ip_address, os, operating_system, platform, agent_version,
						resource_type, organization, kernel, architecture, mac_address, cpu_model,
						total_memory_gb, total_disk_gb, gpu, virtualization, cloud_provider,
						status, online, health_score, last_seen, retry_count, api_key, key_version,
						last_key_rotate, created_at, updated_at
					FROM machines
					ON CONFLICT (id) DO NOTHING;
				`).Error
				if err != nil {
					logger.Error("Failed to copy data from machines to servers", "error", err)
				} else {
					logger.Info("Successfully pre-populated servers table.")
				}
			}
		}
	}

	// Ensure incidents table organization_id column exists and handles NULL values before AutoMigrate adds NOT NULL constraint
	if db.Migrator().HasTable("incidents") {
		if !db.Migrator().HasColumn("incidents", "organization_id") {
			logger.Info("Adding nullable organization_id column to incidents table...")
			_ = db.Exec("ALTER TABLE incidents ADD COLUMN organization_id uuid;").Error
		}
		_ = db.Exec(`
			UPDATE incidents
			SET organization_id = COALESCE(
				(SELECT id FROM organizations ORDER BY created_at ASC LIMIT 1),
				'00000000-0000-0000-0000-000000000000'::uuid
			)
			WHERE organization_id IS NULL;
		`).Error
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.EnrollmentToken{},
		&models.Server{},
		&models.Metric{},
		&models.LinuxServer{},
		&models.LinuxMetric{},
		&models.LinuxProcess{},
		&models.LinuxService{},
		&models.LinuxNetwork{},
		&models.LinuxStorage{},
		&models.LinuxAlert{},
		&models.HistoricalMetric{},
		&models.QueuedMetric{},
		&models.Organization{},
		&models.OrganizationUser{},
		&models.OrganizationSettings{},
		&models.OrganizationInvitation{},
		&models.OrganizationBilling{},
		&models.OrganizationQuota{},
		&models.OrganizationSSO{},
		&models.AlertRule{},
		&models.LinuxDocker{},
		&models.LinuxKubernetes{},
		&models.KubernetesCluster{},
		&models.KubernetesNode{},
		&models.KubernetesPod{},
		&models.KubernetesDeployment{},
		&models.KubernetesStatefulSet{},
		&models.KubernetesDaemonSet{},
		&models.KubernetesService{},
		&models.KubernetesNamespace{},
		&models.KubernetesStorage{},
		&models.KubernetesEvent{},
		&models.LinuxLog{},
		&models.AuditLog{},
		&models.Report{},
		&models.DockerHost{},
		&models.AIIncident{},
		&models.AIRecommendation{},
		&models.AIPrediction{},
		&models.AIHealthScore{},
		&models.AIChatMessage{},
		&models.DockerContainer{},
		&models.DockerImage{},
		&models.DockerVolume{},
		&models.DockerNetwork{},
		&models.DockerEvent{},
		&models.Notification{},
		&models.NotificationPolicy{},
		&models.Incident{},
		&models.IncidentTimeline{},
		&models.IncidentAlert{},
		&models.RemediationPolicy{},
		&models.RemediationJob{},
		&models.Workflow{},
		&models.WorkflowStep{},
		&models.WorkflowExecution{},
		&models.WorkflowLog{},
		&models.BackupRecord{},
		&models.RestoreRecord{},
		&models.DRTestRecord{},
		&models.AnomalyRecord{},
		&models.PredictionRecord{},
		&models.CapacityForecastRecord{},
		&models.RootCauseRecord{},
		&models.AIRecommendationRecord{},
		&models.SLARecord{},
		&models.SecurityAnalyticRecord{},
		&models.CostOptimizationRecord{},
		&models.PolicyRuleRecord{},
		&models.MFASettingRecord{},
		&models.ImmutableAuditRecord{},
		&models.CertificateRecord{},
		&models.ComplianceFrameworkRecord{},
		&models.SecretVaultRecord{},
		&search.SearchIndex{},
		&search.SavedSearch{},
		&search.RecentSearch{},
		&search.SearchAnalytics{},
		&models.APIRequest{},
		&models.ServiceMetric{},
		&models.Command{},
		&models.BlockedCommand{},
		&models.FileOperation{},
		&models.TerminalSession{},
		&models.TerminalCommand{},
	)

	if err != nil {
		logger.Fatal("Migration failed",
			"error", err,
		)
	}

	DB = db

	if sqlDB, err := db.DB(); err == nil {
		sqlDB.SetMaxIdleConns(25)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(15 * time.Minute)
	}

	seedDefaultAdmin(db)

	logger.Info("PostgreSQL connected successfully",
		"host", cfg.DBHost,
		"database", cfg.DBName,
	)
}

func seedDefaultAdmin(db *gorm.DB) {
	var user models.User
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	if err := db.Where("email = ?", "admin@infrapilot.com").First(&user).Error; err != nil {
		adminUser := models.User{
			ID:        uuid.New(),
			Username:  "Admin",
			Email:     "admin@infrapilot.com",
			Password:  string(hashedPassword),
			Role:      models.RoleSuperAdmin,
			IsActive:  true,
			CreatedAt: time.Now(),
		}
		_ = db.Create(&adminUser).Error
	} else {
		_ = db.Model(&models.User{}).Where("email = ?", "admin@infrapilot.com").Update("password", string(hashedPassword)).Error
	}
}
