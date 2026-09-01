package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ShowAPIDocPortal serves the styled interactive developer portal documentation
func ShowAPIDocPortal(c *gin.Context) {
	htmlContent := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>InfraPilot Enterprise - Developer API Reference Portal</title>
    <link href="https://fonts.googleapis.com/css2?family=Outfit:wght@300;400;600;800&family=JetBrains+Mono&display=swap" rel="stylesheet">
    <style>
        :root {
            --bg-color: #020617;
            --card-bg: #0f172a;
            --accent-glow: #38bdf8;
            --accent-purple: #a855f7;
            --text-primary: #f1f5f9;
            --text-secondary: #94a3b8;
            --border-color: #1e293b;
            --btn-hover: #0ea5e9;
        }

        * {
            box-sizing: border-box;
            margin: 0;
            padding: 0;
        }

        body {
            font-family: 'Outfit', sans-serif;
            background-color: var(--bg-color);
            color: var(--text-primary);
            line-height: 1.6;
            padding: 2rem;
        }

        .container {
            max-width: 1100px;
            margin: 0 auto;
        }

        header {
            text-align: center;
            margin-bottom: 3rem;
            padding-bottom: 2rem;
            border-bottom: 1px solid var(--border-color);
            position: relative;
        }

        header h1 {
            font-size: 2.8rem;
            font-weight: 800;
            background: linear-gradient(135deg, var(--accent-glow) 0%, var(--accent-purple) 100%);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            margin-bottom: 0.5rem;
        }

        header p {
            font-size: 1.2rem;
            color: var(--text-secondary);
        }

        .badge {
            display: inline-block;
            padding: 0.3rem 0.8rem;
            font-size: 0.75rem;
            font-weight: 600;
            border-radius: 9999px;
            text-transform: uppercase;
            letter-spacing: 0.05em;
        }

        .badge-get { background-color: rgba(14, 165, 233, 0.15); color: var(--accent-glow); border: 1px solid rgba(14, 165, 233, 0.3); }
        .badge-post { background-color: rgba(168, 85, 247, 0.15); color: var(--accent-purple); border: 1px solid rgba(168, 85, 247, 0.3); }
        .badge-delete { background-color: rgba(239, 68, 68, 0.15); color: #f87171; border: 1px solid rgba(239, 68, 68, 0.3); }

        .api-card {
            background-color: var(--card-bg);
            border: 1px solid var(--border-color);
            border-radius: 16px;
            margin-bottom: 2rem;
            overflow: hidden;
            box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.3);
            transition: transform 0.2s, box-shadow 0.2s;
        }

        .api-card:hover {
            transform: translateY(-2px);
            box-shadow: 0 20px 35px -10px rgba(56, 189, 248, 0.15);
        }

        .api-header {
            padding: 1.5rem;
            background-color: rgba(30, 41, 59, 0.7);
            border-bottom: 1px solid var(--border-color);
            display: flex;
            align-items: center;
            justify-content: space-between;
        }

        .api-endpoint {
            display: flex;
            align-items: center;
            gap: 1rem;
            font-family: 'JetBrains Mono', monospace;
            font-size: 1.1rem;
            font-weight: 600;
        }

        .api-body {
            padding: 1.5rem;
        }

        .api-title {
            font-size: 1.3rem;
            font-weight: 600;
            margin-bottom: 0.5rem;
            color: #fff;
        }

        .api-desc {
            color: var(--text-secondary);
            margin-bottom: 1rem;
        }

        .code-block {
            font-family: 'JetBrains Mono', monospace;
            font-size: 0.9rem;
            background-color: #0b0f19;
            padding: 1rem;
            border-radius: 8px;
            border: 1px solid var(--border-color);
            overflow-x: auto;
            color: #38bdf8;
            margin-top: 0.5rem;
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>InfraPilot Enterprise Portal</h1>
            <p>API Integration and Headless Monitoring Reference v1.0</p>
        </header>

        <!-- Authentication API -->
        <div class="api-card">
            <div class="api-header">
                <div class="api-endpoint">
                    <span class="badge badge-post">POST</span>
                    <span>/api/v1/auth/login</span>
                </div>
            </div>
            <div class="api-body">
                <h2 class="api-title">User Authentication & JWT Tokens</h2>
                <p class="api-desc">Authenticates user credentials and issues stateless JWT Access Tokens and long-lived Refresh Tokens.</p>
                <div class="code-block">
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@infrapilot.com", "password": "password"}'
                </div>
            </div>
        </div>

        <!-- Password Reset Request -->
        <div class="api-card">
            <div class="api-header">
                <div class="api-endpoint">
                    <span class="badge badge-post">POST</span>
                    <span>/api/v1/auth/password-reset/request</span>
                </div>
            </div>
            <div class="api-body">
                <h2 class="api-title">Password Reset Request</h2>
                <p class="api-desc">Initiates a secure password reset token mapped in Redis cache with 15-minute expiration.</p>
                <div class="code-block">
curl -X POST http://localhost:8080/api/v1/auth/password-reset/request \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@infrapilot.com"}'
                </div>
            </div>
        </div>

        <!-- Servers API -->
        <div class="api-card">
            <div class="api-header">
                <div class="api-endpoint">
                    <span class="badge badge-get">GET</span>
                    <span>/api/v1/machines</span>
                </div>
            </div>
            <div class="api-body">
                <h2 class="api-title">Get Monitored Host Infrastructure</h2>
                <p class="api-desc">Retrieves a comprehensive list of all active registered servers and agent host specifications.</p>
                <div class="code-block">
curl -H "Authorization: Bearer &lt;JWT_TOKEN&gt;" http://localhost:8080/api/v1/machines
                </div>
            </div>
        </div>

        <!-- Docker Containers API -->
        <div class="api-card">
            <div class="api-header">
                <div class="api-endpoint">
                    <span class="badge badge-get">GET</span>
                    <span>/api/v1/docker/containers/:machineId</span>
                </div>
            </div>
            <div class="api-body">
                <h2 class="api-title">Get Running Docker Containers</h2>
                <p class="api-desc">List discovered Docker containers with CPU percent, Memory usages, and restart count states.</p>
                <div class="code-block">
curl -H "Authorization: Bearer &lt;JWT_TOKEN&gt;" http://localhost:8080/api/v1/docker/containers/101bfd6f-f4f5-492f-b5d4-67d614741b17
                </div>
            </div>
        </div>

        <!-- Kubernetes API -->
        <div class="api-card">
            <div class="api-header">
                <div class="api-endpoint">
                    <span class="badge badge-get">GET</span>
                    <span>/api/v1/kubernetes/clusters</span>
                </div>
            </div>
            <div class="api-body">
                <h2 class="api-title">Get Discovered Kubernetes Clusters</h2>
                <p class="api-desc">Retrieves a list of all discovered Kubernetes clusters, versions, and providers.</p>
                <div class="code-block">
curl -H "Authorization: Bearer &lt;JWT_TOKEN&gt;" http://localhost:8080/api/v1/kubernetes/clusters
                </div>
            </div>
        </div>

        <!-- AI Incidents API -->
        <div class="api-card">
            <div class="api-header">
                <div class="api-endpoint">
                    <span class="badge badge-get">GET</span>
                    <span>/api/v1/ai/incidents</span>
                </div>
            </div>
            <div class="api-body">
                <h2 class="api-title">Get AI Correlated Incidents</h2>
                <p class="api-desc">Retrieves multi-signal correlated incident reports and chronological timelines generated by the AI Ops engine.</p>
                <div class="code-block">
curl -H "Authorization: Bearer &lt;JWT_TOKEN&gt;" http://localhost:8080/api/v1/ai/incidents
                </div>
            </div>
        </div>

        <!-- Reports API -->
        <div class="api-card">
            <div class="api-header">
                <div class="api-endpoint">
                    <span class="badge badge-get">GET</span>
                    <span>/api/v1/reports</span>
                </div>
            </div>
            <div class="api-body">
                <h2 class="api-title">Get Executive SLA and Compliance Reports</h2>
                <p class="api-desc">Retrieves generated PDF/JSON operational audit and infrastructure health reports.</p>
                <div class="code-block">
curl -H "Authorization: Bearer &lt;JWT_TOKEN&gt;" http://localhost:8080/api/v1/reports
                </div>
            </div>
        </div>
    </div>
</body>
</html>
`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlContent))
}
