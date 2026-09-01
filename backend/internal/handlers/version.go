package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Version(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"project":  "InfraPilot Enterprise",
		"version":  "1.0.0",
		"edition":  "Enterprise",
		"build":    "2026.07.05",
		"language": "Go 1.24",
		"status":   "Stable",
	})
}
