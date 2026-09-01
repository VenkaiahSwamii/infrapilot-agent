package websocket

// PublishEvent sends a generic payload of type T to a room derived from serverID
func PublishEvent[T any](eventName string, serverID string, payload T) {
	room := ""
	if serverID != "" {
		var orgID string
		if ServerOrgLookup != nil {
			orgID, _ = ServerOrgLookup(serverID)
		}
		if orgID != "" {
			room = "org:" + orgID + "/server:" + serverID
		} else {
			room = "server:" + serverID
		}
	}
	evt := Event{
		Event:    eventName,
		ServerID: serverID,
		Room:     room,
		Payload:  payload,
	}
	WS.Broadcast(evt)
}

// PublishMetricUpdated sends a type-safe metric.updated event
func PublishMetricUpdated(serverID string, payload MetricUpdatedPayload) {
	PublishEvent(EventMetricUpdated, serverID, payload)
}

// PublishServerOnline sends a type-safe server.online event
func PublishServerOnline(serverID string, payload ServerStatusPayload) {
	PublishEvent(EventServerOnline, serverID, payload)
}

// PublishServerOffline sends a type-safe server.offline event
func PublishServerOffline(serverID string, payload ServerStatusPayload) {
	PublishEvent(EventServerOffline, serverID, payload)
}

// PublishDockerEvent sends a type-safe docker.event event
func PublishDockerEvent(serverID string, payload DockerEventPayload) {
	PublishEvent("docker.event", serverID, payload)
}
