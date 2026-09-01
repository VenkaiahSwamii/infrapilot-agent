package repository

import (
	"infrapilot/backend/internal/database"
	"infrapilot/backend/internal/models"

	"github.com/google/uuid"
)

type MachineRepository struct {
	serverRepo *ServerRepository
}

func NewMachineRepository() *MachineRepository {
	return &MachineRepository{serverRepo: NewServerRepository()}
}

func (r *MachineRepository) Create(machine *models.Machine) error {
	return r.serverRepo.CreateServer(machine)
}

func (r *MachineRepository) FindAll() ([]models.Machine, error) {
	return r.serverRepo.ListServers()
}

func (r *MachineRepository) FindByID(id uuid.UUID) (*models.Machine, error) {
	return r.serverRepo.GetServer(id)
}

func (r *MachineRepository) FindByAPIKey(apiKey string) (*models.Machine, error) {
	var machine models.Machine
	err := database.DB.Where("api_key = ?", apiKey).First(&machine).Error
	if err != nil {
		return nil, err
	}
	return &machine, nil
}

func (r *MachineRepository) Update(machine *models.Machine) error {
	return r.serverRepo.UpdateServer(machine)
}
