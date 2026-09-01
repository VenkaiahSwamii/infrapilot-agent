package handlers

import (
	"infrapilot/backend/internal/services"

	"github.com/gin-gonic/gin"
)

type MachineHandler struct {
	serverHandler *ServerHandler
}

func NewMachineHandler(machineService *services.MachineService) *MachineHandler {
	return &MachineHandler{
		serverHandler: &ServerHandler{serverService: machineService.GetServerService()},
	}
}

func (h *MachineHandler) EnrollMachine(c *gin.Context) {
	h.serverHandler.EnrollServer(c)
}

func (h *MachineHandler) RegisterMachine(c *gin.Context) {
	h.serverHandler.RegisterServer(c)
}

func (h *MachineHandler) GetMachines(c *gin.Context) {
	h.serverHandler.GetServers(c)
}

func (h *MachineHandler) GetMachineByID(c *gin.Context) {
	h.serverHandler.GetServerByID(c)
}

func (h *MachineHandler) UpdateMachine(c *gin.Context) {
	h.serverHandler.UpdateServer(c)
}

func (h *MachineHandler) RotateMachineKey(c *gin.Context) {
	h.serverHandler.RotateServerKey(c)
}

func (h *MachineHandler) GetMachineKeyRotationStatus(c *gin.Context) {
	h.serverHandler.GetServerKeyRotationStatus(c)
}
