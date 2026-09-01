package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ServeOpenAPISpec returns the complete OpenAPI 3.0.0 JSON specification
func ServeOpenAPISpec(c *gin.Context) {
	spec := gin.H{
		"openapi": "3.0.0",
		"info": gin.H{
			"title":       "InfraPilot Enterprise API",
			"version":     "11.0.0",
			"description": "Production Readiness & Cloud Platform API Specification for InfraPilot Infrastructure Monitoring and Control Plane.",
		},
		"servers": []gin.H{
			{"url": "/api/v1", "description": "Production API Cluster"},
		},
		"paths": gin.H{
			"/auth/login": gin.H{
				"post": gin.H{
					"summary":     "Authenticate User",
					"description": "Issues JWT access and refresh tokens. Enforces account lockout after 5 failed attempts.",
					"responses": gin.H{
						"200": gin.H{"description": "Success with Token Payload"},
						"401": gin.H{"description": "Invalid credentials"},
						"429": gin.H{"description": "Account locked"},
					},
				},
			},
			"/auth/logout": gin.H{
				"post": gin.H{
					"summary":     "Revoke Active Session",
					"description": "Blacklists JWT access token in Redis.",
					"responses": gin.H{
						"200": gin.H{"description": "Logged out successfully"},
					},
				},
			},
			"/health": gin.H{
				"get": gin.H{
					"summary":     "Overall System Health",
					"description": "Returns comprehensive status of PostgreSQL, Redis, AI Service, and WebSocket Hub.",
					"responses": gin.H{
						"200": gin.H{"description": "Healthy"},
						"533": gin.H{"description": "Unhealthy Component"},
					},
				},
			},
			"/ready": gin.H{
				"get": gin.H{
					"summary":     "Kubernetes Readiness Probe",
					"description": "Returns 200 OK if critical dependencies (DB, Redis) are accessible.",
					"responses": gin.H{
						"200": gin.H{"description": "Ready"},
						"503": gin.H{"description": "Not Ready"},
					},
				},
			},
			"/live": gin.H{
				"get": gin.H{
					"summary":     "Kubernetes Liveness Probe",
					"description": "Returns 200 OK if backend instance process is alive.",
					"responses": gin.H{
						"200": gin.H{"description": "Alive"},
					},
				},
			},
			"/machines": gin.H{
				"get": gin.H{
					"summary":     "List Monitored Hosts",
					"description": "Retrieves active registered servers and agent specs.",
					"responses": gin.H{
						"200": gin.H{"description": "List of Machines"},
					},
				},
			},
			"/docker/overview/{serverId}": gin.H{
				"get": gin.H{
					"summary":     "Docker Infrastructure Summary",
					"description": "Overview of Docker daemon, container counts, memory, and images.",
					"responses": gin.H{
						"200": gin.H{"description": "Docker Overview"},
					},
				},
			},
			"/kubernetes/clusters": gin.H{
				"get": gin.H{
					"summary":     "Discovered K8s Clusters",
					"description": "List of active Kubernetes cluster environments.",
					"responses": gin.H{
						"200": gin.H{"description": "List of Clusters"},
					},
				},
			},
			"/aiops/dashboard": gin.H{
				"get": gin.H{
					"summary":     "AIOps Incident & Health Score",
					"description": "Multi-signal AI anomaly detection and auto-remediation recommendations.",
					"responses": gin.H{
						"200": gin.H{"description": "AIOps Dashboard Payload"},
					},
				},
			},
			"/admin/platform/status": gin.H{
				"get": gin.H{
					"summary":     "Platform Component Health & Metrics",
					"description": "Returns backend node, DB, Redis, WS client count, and request throughput metrics.",
					"responses": gin.H{
						"200": gin.H{"description": "Platform Status Payload"},
					},
				},
			},
			"/admin/backup/run": gin.H{
				"post": gin.H{
					"summary":     "Trigger Manual Database Backup",
					"description": "Generates a SQL database snapshot and stores it in backups directory.",
					"responses": gin.H{
						"200": gin.H{"description": "Backup created"},
					},
				},
			},
		},
	}
	c.JSON(http.StatusOK, spec)
}
