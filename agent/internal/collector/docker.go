package collector

import (
	"context"
	"encoding/json"
	"io"
	"strings"
	"sync"
	"time"

	"infrapilot/agent/internal/metrics"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

var (
	eventsBuffer []metrics.DockerEventMetric
	eventsMu     sync.Mutex
	eventsOnce   sync.Once
)

func startEventsListener(cli *client.Client) {
	go func() {
		ctx := context.Background()
		msgChan, errChan := cli.Events(ctx, events.ListOptions{})
		for {
			select {
			case msg := <-msgChan:
				eventsMu.Lock()
				actorName := msg.Actor.Attributes["name"]
				if actorName == "" {
					actorName = msg.Actor.ID
					if len(actorName) > 12 {
						actorName = actorName[:12]
					}
				}
				eventsBuffer = append(eventsBuffer, metrics.DockerEventMetric{
					Time:      time.Unix(msg.Time, msg.TimeNano),
					Type:      string(msg.Type),
					Action:    string(msg.Action),
					ActorID:   msg.Actor.ID,
					ActorName: actorName,
					Message:   strings.ToUpper(string(msg.Type)) + " " + string(msg.Action) + " (ID: " + msg.Actor.ID[:12] + ")",
				})
				// Limit buffer size to avoid memory leaks
				if len(eventsBuffer) > 1000 {
					eventsBuffer = eventsBuffer[len(eventsBuffer)-1000:]
				}
				eventsMu.Unlock()
			case err := <-errChan:
				if err != nil {
					time.Sleep(5 * time.Second)
					return
				}
			}
		}
	}()
}

func collectDockerData() (bool, metrics.DockerVersionInfo, []metrics.DockerContainerMetric, []metrics.DockerImageMetric, []metrics.DockerVolumeMetric, []metrics.DockerNetworkMetric, []metrics.DockerEventMetric) {
	var versionInfo metrics.DockerVersionInfo
	var containerList []metrics.DockerContainerMetric
	var imageList []metrics.DockerImageMetric
	var volumeList []metrics.DockerVolumeMetric
	var networkList []metrics.DockerNetworkMetric
	var eventsList []metrics.DockerEventMetric

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return false, versionInfo, containerList, imageList, volumeList, networkList, eventsList
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = cli.Ping(ctx)
	if err != nil {
		return false, versionInfo, containerList, imageList, volumeList, networkList, eventsList
	}

	// Start events listener
	eventsOnce.Do(func() {
		// Create a separate client that persists for the listener
		listenerCli, listenerErr := client.NewClientWithOpts(client.FromEnv)
		if listenerErr == nil {
			startEventsListener(listenerCli)
		}
	})

	// Get Version Info
	version, err := cli.ServerVersion(ctx)
	if err == nil {
		versionInfo = metrics.DockerVersionInfo{
			DockerVersion: version.Version,
			EngineVersion: version.Version,
			APIVersion:    version.APIVersion,
			HostOS:        version.Os,
			DockerRootDir: "/var/lib/docker", // Default root dir
		}
	}

	// 1. Containers Inventory and Live Metrics
	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err == nil {
		for _, c := range containers {
			name := ""
			if len(c.Names) > 0 {
				name = c.Names[0]
				name = strings.TrimPrefix(name, "/")
			}

			cpuPercent := 0.0
			memoryUsed := uint64(0)
			memoryLimit := uint64(0)
			memoryPct := 0.0
			networkIn := uint64(0)
			networkOut := uint64(0)
			diskRead := uint64(0)
			diskWrite := uint64(0)
			pids := 0

			stats, err := cli.ContainerStats(ctx, c.ID, false)
			if err == nil {
				body, readErr := io.ReadAll(stats.Body)
				stats.Body.Close()

				if readErr == nil {
					var response struct {
						MemoryStats struct {
							Usage uint64 `json:"usage"`
							Limit uint64 `json:"limit"`
						} `json:"memory_stats"`
						CPUStats struct {
							CPUUsage struct {
								TotalUsage uint64 `json:"total_usage"`
							} `json:"cpu_usage"`
							SystemUsage uint64 `json:"system_cpu_usage"`
						} `json:"cpu_stats"`
						PidsStats struct {
							Current int `json:"current"`
						} `json:"pids_stats"`
						Networks map[string]struct {
							RxBytes uint64 `json:"rx_bytes"`
							TxBytes uint64 `json:"tx_bytes"`
						} `json:"networks"`
						BlkioStats struct {
							IoServiceBytesRecursive []struct {
								Op    string `json:"op"`
								Value uint64 `json:"value"`
							} `json:"io_service_bytes_recursive"`
						} `json:"blkio_stats"`
					}

					json.Unmarshal(body, &response)

					memoryUsed = response.MemoryStats.Usage
					memoryLimit = response.MemoryStats.Limit
					if memoryLimit > 0 {
						memoryPct = (float64(memoryUsed) / float64(memoryLimit)) * 100.0
						cpuPercent = float64(response.CPUStats.CPUUsage.TotalUsage) / float64(response.CPUStats.SystemUsage+1) * 100.0
					}

					pids = response.PidsStats.Current
					for _, net := range response.Networks {
						networkIn += net.RxBytes
						networkOut += net.TxBytes
					}
					for _, bio := range response.BlkioStats.IoServiceBytesRecursive {
						if strings.EqualFold(bio.Op, "Read") {
							diskRead += bio.Value
						} else if strings.EqualFold(bio.Op, "Write") {
							diskWrite += bio.Value
						}
					}
				}
			}

			// Inspect to get created details and restarts
			inspect, err := cli.ContainerInspect(ctx, c.ID)
			createdTime := time.Now()
			restartCount := 0
			if err == nil {
				createdTime, _ = time.Parse(time.RFC3339Nano, inspect.Created)
				restartCount = inspect.RestartCount
			}

			containerList = append(containerList, metrics.DockerContainerMetric{
				ID:           c.ID[:12],
				Name:         name,
				Image:        c.Image,
				Status:       c.Status,
				State:        c.State,
				CPUPercent:   cpuPercent,
				MemoryUsed:   memoryUsed,
				MemoryLimit:  memoryLimit,
				MemoryPct:    memoryPct,
				NetworkIn:    networkIn,
				NetworkOut:   networkOut,
				DiskRead:     diskRead,
				DiskWrite:    diskWrite,
				PIDs:         pids,
				RestartCount: restartCount,
				Uptime:       c.Status,
				CreatedTime:  createdTime,
			})
		}
	}

	// 2. Images Collection
	images, err := cli.ImageList(ctx, image.ListOptions{All: false})
	if err == nil {
		for _, img := range images {
			name := "<none>"
			tag := "<none>"
			if len(img.RepoTags) > 0 {
				parts := strings.Split(img.RepoTags[0], ":")
				name = parts[0]
				if len(parts) > 1 {
					tag = parts[1]
				}
			}

			isUnused := img.Containers == -1 || img.Containers == 0

			imageList = append(imageList, metrics.DockerImageMetric{
				ID:          img.ID,
				Name:        name,
				Tag:         tag,
				Size:        img.Size,
				IsUnused:    isUnused,
				CreatedTime: time.Unix(img.Created, 0),
			})
		}
	}

	// 3. Volumes Collection
	volumes, err := cli.VolumeList(ctx, volume.ListOptions{})
	if err == nil {
		for _, vol := range volumes.Volumes {
			volumeList = append(volumeList, metrics.DockerVolumeMetric{
				Name:       vol.Name,
				Driver:     vol.Driver,
				MountPoint: vol.Mountpoint,
				UsageBytes: 0,
			})
		}
	}

	// 4. Networks Collection
	networks, err := cli.NetworkList(ctx, network.ListOptions{})
	if err == nil {
		for _, net := range networks {
			var connected []string
			inspect, err := cli.NetworkInspect(ctx, net.ID, network.InspectOptions{})
			if err == nil {
				for _, container := range inspect.Containers {
					connected = append(connected, container.Name)
				}
			}

			networkList = append(networkList, metrics.DockerNetworkMetric{
				ID:                  net.ID,
				Name:                net.Name,
				Driver:              net.Driver,
				Scope:               net.Scope,
				ConnectedContainers: connected,
			})
		}
	}

	// 5. Events Collection (Read and Flush)
	eventsMu.Lock()
	eventsList = make([]metrics.DockerEventMetric, len(eventsBuffer))
	copy(eventsList, eventsBuffer)
	eventsBuffer = nil
	eventsMu.Unlock()

	return true, versionInfo, containerList, imageList, volumeList, networkList, eventsList
}
