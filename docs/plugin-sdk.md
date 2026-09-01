# InfraPilot Plugin SDK

Build custom monitoring plugins to extend InfraPilot Enterprise.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Creating a Plugin](#creating-a-plugin)
- [Plugin Interface](#plugin-interface)
- [Built-in Plugins](#built-in-plugins)
- [Plugin Configuration](#plugin-configuration)
- [Examples](#examples)
- [Best Practices](#best-practices)

## Overview

InfraPilot's plugin system allows you to extend monitoring capabilities without modifying core code. Plugins run inside the agent and collect custom metrics or perform actions.

### Plugin Types

- **Metric Plugins**: Collect custom metrics (CPU, memory, disk, network)
- **Action Plugins**: Perform remote actions (service restart, package install)
- **Notification Plugins**: Send alerts to external systems
- **Cloud Plugins**: Integrate with cloud providers (AWS, Azure, GCP)

### Example Plugins

- VMware vSphere monitoring
- Hyper-V metrics collection
- AWS CloudWatch integration
- Azure Monitor metrics
- GCP Cloud Monitoring
- Slack notifications
- Microsoft Teams notifications
- Discord webhooks
- Custom scripts and APIs

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    InfraPilot Agent                          │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ Plugin       │  │ Plugin       │  │ Plugin       │      │
│  │ Manager      │─▶│ (Docker)    │─▶│ (K8s)       │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│         │                  │                  │              │
│         └──────────────────┼──────────────────┘              │
│                            │                                 │
│                    ┌───────▼────────┐                        │
│                    │ Event Bus      │                        │
│                    │ (Redis Streams)│                       │
│                    └───────┬────────┘                        │
│                            │                                 │
│                    ┌───────▼────────┐                        │
│                    │ Metrics        │                        │
│                    │ Collector     │                        │
│                    └───────────────┘                        │
└─────────────────────────────────────────────────────────────┘
```

## Creating a Plugin

### Step 1: Define Plugin Structure

Create a new Go file in `agent/internal/plugins/`:

```go
package plugins

import "time"

// CloudWatchPlugin collects AWS CloudWatch metrics
type CloudWatchPlugin struct {
    enabled bool
    region  string
    // Add your fields
}

// NewCloudWatchPlugin creates a new CloudWatch plugin
func NewCloudWatchPlugin(region string) *CloudWatchPlugin {
    return &CloudWatchPlugin{
        enabled: true,
        region:  region,
    }
}

// Name returns the plugin name
func (p *CloudWatchPlugin) Name() string {
    return "cloudwatch"
}

// Collect gathers metrics from AWS CloudWatch
func (p *CloudWatchPlugin) Collect() (interface{}, error) {
    // Implement your metric collection logic
    return map[string]interface{}{
        "region": p.region,
        "metrics": map[string]float64{
            "cpu_utilization": 45.2,
            "memory_usage": 62.1,
        },
    }, nil
}

// HealthCheck returns plugin health status
func (p *CloudWatchPlugin) HealthCheck() PluginHealth {
    health := NewPluginHealth(p.Name())
    health.Status = "healthy"
    health.LastRun = time.Now()
    health.Metadata = map[string]string{
        "region": p.region,
    }
    return health
}

// IsEnabled returns if plugin is enabled
func (p *CloudWatchPlugin) IsEnabled() bool {
    return p.enabled
}

// SetEnabled enables/disables the plugin
func (p *CloudWatchPlugin) SetEnabled(enabled bool) {
    p.enabled = enabled
}
```

### Step 2: Register Plugin

Add your plugin to the manager in `agent/cmd/agent/main.go`:

```go
func setupPlugins() *plugins.Manager {
    manager := plugins.NewManager()

    // Register built-in plugins
    manager.Register(plugins.NewLinuxPlugin())
    manager.Register(plugins.NewDockerPlugin())
    manager.Register(plugins.NewKubernetesPlugin())

    // Register your custom plugin
    manager.Register(plugins.NewCloudWatchPlugin("us-east-1"))

    return manager
}
```

### Step 3: Add Configuration

Add plugin configuration to `agent/config.example.yaml`:

```yaml
# AWS CloudWatch Plugin
plugins:
  cloudwatch:
    enabled: true
    region: us-east-1
    access_key: ${AWS_ACCESS_KEY}
    secret_key: ${AWS_SECRET_KEY}
    metrics:
      - CPUUtilization
      - MemoryUtilization
      - DiskReadBytes
    period: 60s
```

## Plugin Interface

All plugins must implement the `Plugin` interface:

```go
type Plugin interface {
    // Name returns the unique identifier
    Name() string

    // Collect gathers metrics and returns data
    Collect() (interface{}, error)

    // HealthCheck returns current health status
    HealthCheck() PluginHealth

    // IsEnabled returns whether plugin is active
    IsEnabled() bool

    // SetEnabled toggles plugin activation
    SetEnabled(enabled bool)
}
```

### Plugin Health

```go
type PluginHealth struct {
    Plugin   string            `json:"plugin"`
    Status   string            `json:"status"` // "healthy", "degraded", "unhealthy"
    LastRun  time.Time         `json:"last_run"`
    Error    string            `json:"error,omitempty"`
    Metadata map[string]string `json:"metadata,omitempty"`
}
```

## Built-in Plugins

### Linux Plugin (`linux`)

Collects system metrics from Linux hosts:

```go
plugin := plugins.NewLinuxPlugin()
metrics := plugin.Collect()
// Returns:
// {
//   "cpu": {"usage": 45.2, "cores": 8},
//   "memory": {"total": 16384, "used": 10240, "percent": 62.5},
//   "disk": [{"mount": "/", "total": 500, "used": 200}],
//   "network": {"interfaces": ["eth0"], "bytes_in": 1024, "bytes_out": 512}
// }
```

### Docker Plugin (`docker`)

Monitors Docker containers:

```go
plugin := plugins.NewDockerPlugin()
metrics := plugin.Collect()
// Returns:
// {
//   "containers": [
//     {"id": "abc123", "name": "web", "status": "running", "cpu": 0.5, "memory": 256}
//   ],
//   "images": 5,
//   "volumes": 3
// }
```

### Kubernetes Plugin (`kubernetes`)

Tracks Kubernetes resources:

```go
plugin := plugins.NewKubernetesPlugin()
metrics := plugin.Collect()
// Returns:
// {
//   "clusters": [{"name": "prod", "context": "k8s"}],
//   "nodes": [{"name": "node-1", "status": "Ready", "cpu": 45.2}],
//   "pods": [{"name": "web-123", "status": "Running", "namespace": "default"}]
// }
```

## Plugin Configuration

### Configuration File

Plugins are configured in `agent/config.yaml`:

```yaml
plugins:
  linux:
    enabled: true
    interval: 10s

  docker:
    enabled: true
    interval: 15s
    socket: /var/run/docker.sock

  kubernetes:
    enabled: true
    interval: 30s
    kubeconfig: ~/.kube/config

  cloudwatch:
    enabled: true
    region: us-east-1
    access_key: ${AWS_ACCESS_KEY}
    secret_key: ${AWS_SECRET_KEY}
    interval: 60s
    metrics:
      - CPUUtilization
      - MemoryUtilization

  slack:
    enabled: true
    webhook_url: ${SLACK_WEBHOOK}
    channel: "#alerts"
```

### Environment Variables

Use environment variables for secrets:

```bash
export AWS_ACCESS_KEY=AKIAIOSFODNN7EXAMPLE
export AWS_SECRET_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
export SLACK_WEBHOOK=https://example.com/api/v1/slack-webhook-placeholder
```

## Examples

### Example 1: AWS CloudWatch Plugin

```go
package plugins

import (
    "fmt"
    "time"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/cloudwatch"
)

type CloudWatchPlugin struct {
    enabled bool
    region  string
    svc     *cloudwatch.CloudWatch
}

func NewCloudWatchPlugin(region string) *CloudWatchPlugin {
    sess := session.Must(session.NewSession(&aws.Config{
        Region: aws.String(region),
    }))

    return &CloudWatchPlugin{
        enabled: true,
        region:  region,
        svc:     cloudwatch.New(sess),
    }
}

func (p *CloudWatchPlugin) Name() string {
    return "cloudwatch"
}

func (p *CloudWatchPlugin) Collect() (interface{}, error) {
    input := &cloudwatch.GetMetricDataInput{
        MetricDataQueries: []*cloudwatch.MetricDataQuery{
            {
                Id: aws.String("cpu"),
                MetricStat: &cloudwatch.MetricStat{
                    Metric: &cloudwatch.Metric{
                        MetricName: aws.String("CPUUtilization"),
                        Namespace:  aws.String("AWS/EC2"),
                    },
                    Period: aws.Int64(60),
                    Stat:   aws.String("Average"),
                },
            },
        },
        StartTime: aws.Time(time.Now().Add(-5 * time.Minute)),
        EndTime:   aws.Time(time.Now()),
    }

    result, err := p.svc.GetMetricData(input)
    if err != nil {
        return nil, fmt.Errorf("failed to get CloudWatch metrics: %w", err)
    }

    metrics := make(map[string]float64)
    for _, metric := range result.MetricDataResults {
        if len(metric.Values) > 0 {
            metrics[*metric.Id] = *metric.Values[0]
        }
    }

    return map[string]interface{}{
        "region":   p.region,
        "metrics":  metrics,
        "timestamp": time.Now(),
    }, nil
}

func (p *CloudWatchPlugin) HealthCheck() PluginHealth {
    return PluginHealth{
        Plugin:   p.Name(),
        Status:   "healthy",
        LastRun:  time.Now(),
        Metadata: map[string]string{"region": p.region},
    }
}

func (p *CloudWatchPlugin) IsEnabled() bool {
    return p.enabled
}

func (p *CloudWatchPlugin) SetEnabled(enabled bool) {
    p.enabled = enabled
}
```

### Example 2: Slack Notification Plugin

```go
package plugins

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "github.com/venkaiswami/infrapilot-enterprise/agent/internal/models"
)

type SlackPlugin struct {
    enabled    bool
    webhookURL string
    channel    string
    username   string
}

type SlackMessage struct {
    Channel string `json:"channel"`
    Text    string `json:"text"`
    Username string `json:"username"`
    IconEmoji string `json:"icon_emoji"`
}

func NewSlackPlugin(webhookURL, channel string) *SlackPlugin {
    return &SlackPlugin{
        enabled:    true,
        webhookURL: webhookURL,
        channel:    channel,
        username:   "InfraPilot Bot",
    }
}

func (p *SlackPlugin) Name() string {
    return "slack"
}

func (p *SlackPlugin) Collect() (interface{}, error) {
    return map[string]interface{}{
        "status":  "connected",
        "channel": p.channel,
    }, nil
}

func (p *SlackPlugin) SendAlert(alert models.Alert) error {
    emoji := ":warning:"
    if alert.Severity == "critical" {
        emoji = ":rotating_light:"
    }

    message := SlackMessage{
        Channel:    p.channel,
        Text:       fmt.Sprintf("%s *%s*\n%s", emoji, alert.Name, alert.Message),
        Username:   p.username,
        IconEmoji:  emoji,
    }

    jsonData, _ := json.Marshal(message)
    resp, err := http.Post(p.webhookURL, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("slack webhook returned status %d", resp.StatusCode)
    }

    return nil
}

func (p *SlackPlugin) HealthCheck() PluginHealth {
    return PluginHealth{
        Plugin:   p.Name(),
        Status:   "healthy",
        LastRun:  time.Now(),
        Metadata: map[string]string{"channel": p.channel},
    }
}

func (p *SlackPlugin) IsEnabled() bool {
    return p.enabled
}

func (p *SlackPlugin) SetEnabled(enabled bool) {
    p.enabled = enabled
}
```

### Example 3: Custom Script Plugin

```go
package plugins

import (
    "bufio"
    "os/exec"
    "strings"
    "time"
)

type ScriptPlugin struct {
    enabled bool
    script  string
    args    []string
}

func NewScriptPlugin(scriptPath string, args []string) *ScriptPlugin {
    return &ScriptPlugin{
        enabled: true,
        script:  scriptPath,
        args:    args,
    }
}

func (p *ScriptPlugin) Name() string {
    return "custom_script"
}

func (p *ScriptPlugin) Collect() (interface{}, error) {
    cmd := exec.Command(p.script, p.args...)

    output, err := cmd.CombinedOutput()
    if err != nil {
        return nil, fmt.Errorf("script execution failed: %w, output: %s", err, string(output))
    }

    // Parse script output (JSON expected)
    var result map[string]interface{}
    if err := json.Unmarshal(output, &result); err != nil {
        return nil, fmt.Errorf("failed to parse script output: %w", err)
    }

    return result, nil
}

func (p *ScriptPlugin) HealthCheck() PluginHealth {
    // Check if script exists and is executable
    _, err := exec.LookPath(p.script)
    status := "healthy"
    if err != nil {
        status = "unhealthy"
    }

    return PluginHealth{
        Plugin:   p.Name(),
        Status:   status,
        LastRun:  time.Now(),
        Metadata: map[string]string{"script": p.script},
    }
}

func (p *ScriptPlugin) IsEnabled() bool {
    return p.enabled
}

func (p *ScriptPlugin) SetEnabled(enabled bool) {
    p.enabled = enabled
}
```

Usage in config:

```yaml
plugins:
  custom_script:
    enabled: true
    script: /opt/infrapilot/scripts/custom-metrics.sh
    args: ["--json", "--collect"]
    interval: 30s
```

## Best Practices

### 1. Error Handling

Always handle errors gracefully:

```go
func (p *MyPlugin) Collect() (interface{}, error) {
    data, err := fetchMetrics()
    if err != nil {
        // Log error but don't crash
        logger.Errorf("Plugin %s collection failed: %v", p.Name(), err)
        return nil, err
    }

    // Validate data
    if data == nil {
        return map[string]interface{}{}, nil
    }

    return data, nil
}
```

### 2. Timeout Configuration

Set reasonable timeouts:

```go
func (p *MyPlugin) Collect() (interface{}, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Use context for network calls
    result, err := fetchWithContext(ctx, p.endpoint)
    return result, err
}
```

### 3. Resource Management

Don't leak resources:

```go
func (p *MyPlugin) Collect() (interface{}, error) {
    client, err := createClient()
    if err != nil {
        return nil, err
    }
    defer client.Close()

    return client.Fetch()
}
```

### 4. Logging

Add structured logging:

```go
func (p *MyPlugin) Collect() (interface{}, error) {
    start := time.Now()

    data, err := p.fetch()
    duration := time.Since(start)

    logger.WithFields(logrus.Fields{
        "plugin": p.Name(),
        "duration_ms": duration.Milliseconds(),
        "error": err,
    }).Info("Plugin collection completed")

    return data, err
}
```

### 5. Configuration Validation

Validate config on startup:

```go
func NewMyPlugin(config PluginConfig) (*MyPlugin, error) {
    if config.Endpoint == "" {
        return nil, fmt.Errorf("endpoint is required")
    }

    if config.Timeout <= 0 {
        config.Timeout = 10 * time.Second
    }

    return &MyPlugin{
        enabled:  true,
        endpoint: config.Endpoint,
        timeout:  config.Timeout,
    }, nil
}
```

### 6. Test Your Plugin

Write tests:

```go
func TestMyPlugin(t *testing.T) {
    plugin := NewMyPlugin(PluginConfig{
        Endpoint: "http://example.com",
    })

    // Test Name
    assert.Equal(t, "my_plugin", plugin.Name())

    // Test Collect
    data, err := plugin.Collect()
    assert.NoError(t, err)
    assert.NotNil(t, data)

    // Test HealthCheck
    health := plugin.HealthCheck()
    assert.Equal(t, "healthy", health.Status)
}
```

## Testing Plugins

Run plugin tests:

```bash
cd agent
go test -v -run TestMyPlugin ./internal/plugins/
```

## Distributing Plugins

### Option 1: Include in Agent

Contribute plugins to the InfraPilot repository via pull requests.

### Option 2: External Plugin (Future)

Future versions will support loading external plugins:

```yaml
plugins:
  external:
    - path: /opt/plugins/cloudwatch.so
      config:
        region: us-east-1
```

## Support

- **Documentation**: https://github.com/venkaiswami/infrapilot-enterprise/tree/main/docs
- **Issues**: https://github.com/venkaiswami/infrapilot-enterprise/issues
- **Discussions**: https://github.com/venkaiswami/infrapilot-enterprise/discussions
- **Examples**: https://github.com/venkaiswami/infrapilot-enterprise/tree/main/agent/internal/plugins

Happy plugin building! 🚀