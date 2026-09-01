package websocket

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"infrapilot/backend/internal/auth"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

func init() {
	ServerOrgLookup = func(serverID string) (string, error) {
		return "org-123", nil
	}
}

func TestWebSocketHub_ClientLifecycleAndBroadcast(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	// 1. Initial State
	if hub.ActiveClientsCount() != 0 {
		t.Errorf("Expected 0 active clients, got %d", hub.ActiveClientsCount())
	}

	// 2. HTTP Server with WS handler
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hub.Handle(w, r)
	}))
	defer server.Close()

	token, err := auth.GenerateToken("test-user", "test-user", "test@test.com", "admin", "org-123")
	if err != nil {
		t.Fatalf("Failed to generate test auth token: %v", err)
	}

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?token=" + token

	// 3. Connect Client 1
	conn1, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to dial WebSocket server: %v", err)
	}
	defer conn1.Close()

	// Give hub time to register client
	time.Sleep(50 * time.Millisecond)

	if hub.ActiveClientsCount() != 1 {
		t.Errorf("Expected 1 active client, got %d", hub.ActiveClientsCount())
	}

	// 4. Send Ping Event
	pingMsg := ClientMessage{
		Event:     "ping",
		Timestamp: time.Now().Unix(),
	}
	pingBytes, _ := json.Marshal(pingMsg)
	if err := conn1.WriteMessage(websocket.TextMessage, pingBytes); err != nil {
		t.Fatalf("Failed to write ping message: %v", err)
	}

	// Read Pong Response
	_, respBytes, err := conn1.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to read pong message: %v", err)
	}

	var pongResp map[string]interface{}
	_ = json.Unmarshal(respBytes, &pongResp)
	if pongResp["event"] != "pong" {
		t.Errorf("Expected event pong, got %v", pongResp["event"])
	}

	// 5. Broadcast Message Test
	broadcastMsg := map[string]string{
		"event":  "machine_status_updated",
		"status": "ONLINE",
	}
	hub.BroadcastJSON(broadcastMsg)

	_, broadcastRespBytes, err := conn1.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to receive broadcasted message: %v", err)
	}

	var rxMsg map[string]string
	_ = json.Unmarshal(broadcastRespBytes, &rxMsg)
	if rxMsg["status"] != "ONLINE" {
		t.Errorf("Expected status ONLINE, got %s", rxMsg["status"])
	}
}

func TestWebSocketHub_RoomSubscription(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hub.Handle(w, r)
	}))
	defer server.Close()

	token, err := auth.GenerateToken("test-user", "test-user", "test@test.com", "admin", "org-123")
	if err != nil {
		t.Fatalf("Failed to generate test auth token: %v", err)
	}

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?token=" + token

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to dial WebSocket server: %v", err)
	}
	defer conn.Close()

	time.Sleep(50 * time.Millisecond)

	// Subscribe to "alerts" room
	subMsg := ClientMessage{
		Event: "subscribe",
		Room:  "alerts",
	}
	subBytes, _ := json.Marshal(subMsg)
	_ = conn.WriteMessage(websocket.TextMessage, subBytes)

	// Wait for "subscribed" confirmation
	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("Failed to read subscription confirmation: %v", err)
		}
		var confirm map[string]interface{}
		_ = json.Unmarshal(msgBytes, &confirm)
		if confirm["event"] == "subscribed" && confirm["room"] == "alerts" {
			break
		}
	}

	// Broadcast event to "alerts" room
	alertEvent := Event{
		Event:   EventAlertCreated,
		Room:    "alerts",
		Payload: map[string]string{"alert_id": "alert-999", "severity": "CRITICAL"},
	}
	hub.Broadcast(alertEvent)

	_, rxBytes, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("Failed to receive room event: %v", err)
	}

	var rxEvt Event
	_ = json.Unmarshal(rxBytes, &rxEvt)
	if rxEvt.Event != EventAlertCreated {
		t.Errorf("Expected event alert.created, got %s", rxEvt.Event)
	}
}

func TestWebSocketHub_RBACSubscription(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hub.Handle(w, r)
	}))
	defer server.Close()

	// 1. Connect Client with role "viewer"
	viewerToken, err := auth.GenerateToken("viewer-user", "viewer-user", "viewer@test.com", "viewer", "org-123")
	if err != nil {
		t.Fatalf("Failed to generate viewer token: %v", err)
	}
	wsURLViewer := "ws" + strings.TrimPrefix(server.URL, "http") + "?token=" + viewerToken
	connViewer, _, err := websocket.DefaultDialer.Dial(wsURLViewer, nil)
	if err != nil {
		t.Fatalf("Failed to dial as viewer: %v", err)
	}
	defer connViewer.Close()

	// Subscribe to a restricted room: "terminal:session-123"
	subMsg := ClientMessage{
		Event: "subscribe",
		Room:  "terminal:session-123",
	}
	subBytes, _ := json.Marshal(subMsg)
	_ = connViewer.WriteMessage(websocket.TextMessage, subBytes)
	time.Sleep(50 * time.Millisecond)

	// Broadcast an event to "terminal:session-123"
	termEvent := Event{
		Event:   "terminal.data",
		Room:    "terminal:session-123",
		Payload: "secret-data",
	}
	hub.Broadcast(termEvent)

	// Verify viewer does NOT receive the message (read deadline times out)
	_ = connViewer.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
	_, _, err = connViewer.ReadMessage()
	if err == nil {
		t.Error("Expected read timeout for restricted subscription (viewer role), but received message")
	}

	// 2. Connect Client with role "admin"
	adminToken, err := auth.GenerateToken("admin-user", "admin-user", "admin@test.com", "admin", "org-123")
	if err != nil {
		t.Fatalf("Failed to generate admin token: %v", err)
	}
	wsURLAdmin := "ws" + strings.TrimPrefix(server.URL, "http") + "?token=" + adminToken
	connAdmin, _, err := websocket.DefaultDialer.Dial(wsURLAdmin, nil)
	if err != nil {
		t.Fatalf("Failed to dial as admin: %v", err)
	}
	defer connAdmin.Close()

	// Subscribe to "terminal:session-123"
	_ = connAdmin.WriteMessage(websocket.TextMessage, subBytes)

	// Wait for confirmation
	for {
		_, msgBytes, err := connAdmin.ReadMessage()
		if err != nil {
			t.Fatalf("Failed to read subscription confirmation: %v", err)
		}
		var confirm map[string]interface{}
		_ = json.Unmarshal(msgBytes, &confirm)
		if confirm["event"] == "subscribed" && confirm["room"] == "terminal:session-123" {
			break
		}
	}

	// Broadcast another event to "terminal:session-123"
	hub.Broadcast(termEvent)

	// Verify admin DOES receive the message
	_ = connAdmin.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
	_, rxBytes, err := connAdmin.ReadMessage()
	if err != nil {
		t.Fatalf("Admin failed to receive authorized terminal event: %v", err)
	}

	var rxEvt2 Event
	_ = json.Unmarshal(rxBytes, &rxEvt2)
	if rxEvt2.Event != "terminal.data" {
		t.Errorf("Expected event terminal.data, got %s", rxEvt2.Event)
	}
}

func TestWebSocketHub_AuthorizeSubscriptionMatrix(t *testing.T) {
	hub := NewHub()

	tests := []struct {
		name        string
		role        string
		orgID       string
		room        string
		shouldAllow bool
	}{
		// Viewer Role Matrix
		{"viewer-metrics-allowed", RoleViewer, "org-123", "org:org-123/server:machine01", true},
		{"viewer-alerts-allowed", RoleViewer, "org-123", "alerts", true},
		{"viewer-terminal-denied", RoleViewer, "org-123", "terminal:session-123", false},
		{"viewer-commands-denied", RoleViewer, "org-123", "commands:cmd-123", false},
		{"viewer-orchestration-denied", RoleViewer, "org-123", "orchestration:job-123", false},

		// Operator Role Matrix
		{"operator-metrics-allowed", RoleOperator, "org-123", "org:org-123/server:machine01", true},
		{"operator-terminal-allowed", RoleOperator, "org-123", "terminal:session-123", true},
		{"operator-commands-allowed", RoleOperator, "org-123", "commands:cmd-123", true},
		{"operator-orchestration-allowed", RoleOperator, "org-123", "orchestration:job-123", true},

		// Admin Role Matrix (Full Access)
		{"admin-metrics-allowed", RoleAdmin, "org-123", "org:org-123/server:machine01", true},
		{"admin-terminal-allowed", RoleAdmin, "org-123", "terminal:session-123", true},
		{"admin-commands-allowed", RoleAdmin, "org-123", "commands:cmd-123", true},

		// Guest / Unknown Role Matrix
		{"guest-metrics-allowed", RoleGuest, "org-123", "org:org-123/server:machine01", true},
		{"guest-terminal-denied", RoleGuest, "org-123", "terminal:session-123", false},
		{"guest-commands-denied", RoleGuest, "org-123", "commands:cmd-123", false},

		// Tenant Isolation Matrix (Admin role tries to cross org boundaries)
		{"tenant-isolation-cross-org-denied", RoleAdmin, "org-123", "org:org-999/server:machine01", false},
		{"tenant-isolation-server-cross-org-denied", RoleAdmin, "org-123", "server:machine-of-another-org", false},
	}

	// Update mock to return different org for cross org server test
	ServerOrgLookup = func(serverID string) (string, error) {
		if serverID == "machine-of-another-org" {
			return "org-999", nil
		}
		return "org-123", nil
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &Client{
				Role:           tc.role,
				OrganizationID: tc.orgID,
			}
			allowed := hub.AuthorizeSubscription(client, tc.room)
			if allowed != tc.shouldAllow {
				t.Errorf("Role %s subscribing to %s: expected allowed=%v, got %v", tc.role, tc.room, tc.shouldAllow, allowed)
			}
		})
	}
}

func generateExpiredToken(userID, username, email, role, organizationID string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "ChangeThisToASecretKey"
	}
	claims := auth.Claims{
		UserID:         userID,
		Username:       username,
		Email:          email,
		Role:           role,
		OrganizationID: organizationID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-24 * time.Hour)), // Expired 24 hours ago
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-48 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func TestWebSocketHub_AuthErrors(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hub.Handle(w, r)
	}))
	defer server.Close()

	// 1. Missing Token
	wsURLMissing := "ws" + strings.TrimPrefix(server.URL, "http")
	_, respMissing, err := websocket.DefaultDialer.Dial(wsURLMissing, nil)
	if err == nil {
		t.Fatal("Expected dial to fail with missing token")
	}
	if respMissing != nil && respMissing.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 401 Unauthorized for missing token, got %d", respMissing.StatusCode)
	}

	// 2. Invalid Token (garbage string)
	wsURLInvalid := "ws" + strings.TrimPrefix(server.URL, "http") + "?token=this-is-not-a-valid-jwt"
	_, respInvalid, err := websocket.DefaultDialer.Dial(wsURLInvalid, nil)
	if err == nil {
		t.Fatal("Expected dial to fail with invalid token")
	}
	if respInvalid != nil && respInvalid.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 401 Unauthorized for invalid token, got %d", respInvalid.StatusCode)
	}

	// 3. Expired Token
	expiredToken, err := generateExpiredToken("test-user", "test-user", "test@test.com", "admin", "org-123")
	if err != nil {
		t.Fatalf("Failed to generate expired token: %v", err)
	}
	wsURLExpired := "ws" + strings.TrimPrefix(server.URL, "http") + "?token=" + expiredToken
	_, respExpired, err := websocket.DefaultDialer.Dial(wsURLExpired, nil)
	if err == nil {
		t.Fatal("Expected dial to fail with expired token")
	}
	if respExpired != nil && respExpired.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 401 Unauthorized for expired token, got %d", respExpired.StatusCode)
	}
}

func TestWebSocketHub_TenantIsolationBroadcast(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hub.Handle(w, r)
	}))
	defer server.Close()

	// Connect Client 1 (Belongs to org-123)
	token1, err := auth.GenerateToken("user-1", "user-1", "user1@test.com", "admin", "org-123")
	if err != nil {
		t.Fatalf("Failed to generate token 1: %v", err)
	}
	wsURL1 := "ws" + strings.TrimPrefix(server.URL, "http") + "?token=" + token1
	conn1, _, err := websocket.DefaultDialer.Dial(wsURL1, nil)
	if err != nil {
		t.Fatalf("Failed to dial client 1: %v", err)
	}
	defer conn1.Close()

	// Start asynchronous reader to prevent read deadline errors from breaking the socket
	msgChan := make(chan []byte, 10)
	go func() {
		for {
			_, msgBytes, err := conn1.ReadMessage()
			if err != nil {
				return
			}
			msgChan <- msgBytes
		}
	}()

	// Client 1 subscribes to global
	subMsg := ClientMessage{
		Event: "subscribe",
		Room:  "global",
	}
	subBytes, _ := json.Marshal(subMsg)
	_ = conn1.WriteMessage(websocket.TextMessage, subBytes)

	// Wait for confirmation from channel
	select {
	case rxConfirm := <-msgChan:
		var confirm map[string]interface{}
		_ = json.Unmarshal(rxConfirm, &confirm)
		if confirm["event"] != "subscribed" || confirm["room"] != "global" {
			t.Fatalf("Expected subscribed confirmation, got: %v", confirm)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Timeout waiting for subscription confirmation")
	}

	// 1. Broadcast an event targeting org-456 (different organization)
	crossOrgEvent := Event{
		Event:    EventMetricUpdated,
		ServerID: "machine01",
		Room:     "org:org-456/server:machine01", // room belongs to org-456
		Payload:  "metrics-data",
	}
	hub.Broadcast(crossOrgEvent)

	// Sleep briefly and verify no messages were received in msgChan
	time.Sleep(50 * time.Millisecond)
	if len(msgChan) > 0 {
		msg := <-msgChan
		t.Errorf("Expected client 1 (org-123) to NOT receive cross-org event, but received: %s", string(msg))
	}

	// 2. Broadcast an event targeting org-123 (same organization)
	sameOrgEvent := Event{
		Event:    EventMetricUpdated,
		ServerID: "machine01",
		Room:     "org:org-123/server:machine01", // room belongs to org-123
		Payload:  "metrics-data",
	}
	hub.Broadcast(sameOrgEvent)

	// Verify client 1 (org-123) DOES receive the metrics
	select {
	case rxBytes := <-msgChan:
		var rxEvt Event
		_ = json.Unmarshal(rxBytes, &rxEvt)
		if rxEvt.Event != EventMetricUpdated {
			t.Errorf("Expected event %s, got %s", EventMetricUpdated, rxEvt.Event)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Timeout waiting for same-org metrics broadcast")
	}
}

func TestWebSocketHub_EventLevelAuthorization(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hub.Handle(w, r)
	}))
	defer server.Close()

	// 1. Connect Client 1 with role "viewer" (org-123)
	viewerToken, err := auth.GenerateToken("viewer-user", "viewer-user", "viewer@test.com", "viewer", "org-123")
	if err != nil {
		t.Fatalf("Failed to generate viewer token: %v", err)
	}
	wsURLViewer := "ws" + strings.TrimPrefix(server.URL, "http") + "?token=" + viewerToken
	connViewer, _, err := websocket.DefaultDialer.Dial(wsURLViewer, nil)
	if err != nil {
		t.Fatalf("Failed to dial viewer: %v", err)
	}
	defer connViewer.Close()

	// Async reader for viewer
	viewerChan := make(chan []byte, 10)
	go func() {
		for {
			_, msgBytes, err := connViewer.ReadMessage()
			if err != nil {
				return
			}
			viewerChan <- msgBytes
		}
	}()

	// Viewer subscribes to global
	subMsg := ClientMessage{Event: "subscribe", Room: "global"}
	subBytes, _ := json.Marshal(subMsg)
	_ = connViewer.WriteMessage(websocket.TextMessage, subBytes)
	select {
	case <-viewerChan: // wait for confirm
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Viewer subscription timeout")
	}

	// 2. Connect Client 2 with role "admin" (org-123)
	adminToken, err := auth.GenerateToken("admin-user", "admin-user", "admin@test.com", "admin", "org-123")
	if err != nil {
		t.Fatalf("Failed to generate admin token: %v", err)
	}
	wsURLAdmin := "ws" + strings.TrimPrefix(server.URL, "http") + "?token=" + adminToken
	connAdmin, _, err := websocket.DefaultDialer.Dial(wsURLAdmin, nil)
	if err != nil {
		t.Fatalf("Failed to dial admin: %v", err)
	}
	defer connAdmin.Close()

	// Async reader for admin
	adminChan := make(chan []byte, 10)
	go func() {
		for {
			_, msgBytes, err := connAdmin.ReadMessage()
			if err != nil {
				return
			}
			adminChan <- msgBytes
		}
	}()

	// Admin subscribes to global
	_ = connAdmin.WriteMessage(websocket.TextMessage, subBytes)
	select {
	case <-adminChan: // wait for confirm
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Admin subscription timeout")
	}

	// 3. Broadcast terminal event (sensitive)
	termEvent := Event{
		Event:    "terminal.data",
		ServerID: "machine01",
		Room:     "org:org-123/terminal:session-999",
		Payload:  "keystrokes",
	}
	hub.Broadcast(termEvent)

	// Verify viewer does NOT receive the terminal event
	time.Sleep(50 * time.Millisecond)
	if len(viewerChan) > 0 {
		msg := <-viewerChan
		t.Errorf("Viewer incorrectly received terminal event: %s", string(msg))
	}

	// Verify admin DOES receive the terminal event
	select {
	case rxBytes := <-adminChan:
		var rxEvt Event
		_ = json.Unmarshal(rxBytes, &rxEvt)
		if rxEvt.Event != "terminal.data" {
			t.Errorf("Expected terminal.data event, got %s", rxEvt.Event)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Admin timeout waiting for terminal event")
	}

	// 4. Broadcast metrics event (view only)
	metricsEvent := Event{
		Event:    "metric.updated",
		ServerID: "machine01",
		Room:     "org:org-123/server:machine01",
		Payload:  "metrics-payload",
	}
	hub.Broadcast(metricsEvent)

	// Verify viewer DOES receive the metrics event
	select {
	case rxBytes := <-viewerChan:
		var rxEvt Event
		_ = json.Unmarshal(rxBytes, &rxEvt)
		if rxEvt.Event != "metric.updated" {
			t.Errorf("Expected metric.updated event for viewer, got %s", rxEvt.Event)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Viewer timeout waiting for metrics event")
	}

	// Verify admin DOES receive the metrics event too
	select {
	case rxBytes := <-adminChan:
		var rxEvt Event
		_ = json.Unmarshal(rxBytes, &rxEvt)
		if rxEvt.Event != "metric.updated" {
			t.Errorf("Expected metric.updated event for admin, got %s", rxEvt.Event)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Admin timeout waiting for metrics event")
	}

	// 5. Broadcast an unknown/unclassified event (default-deny policy check)
	unknownEvent := Event{
		Event:    "secret.exposed",
		ServerID: "machine01",
		Room:     "org:org-123/server:machine01",
		Payload:  "credentials",
	}
	hub.Broadcast(unknownEvent)

	// Verify viewer does NOT receive the unknown event (blocked by default-deny)
	time.Sleep(50 * time.Millisecond)
	if len(viewerChan) > 0 {
		msg := <-viewerChan
		t.Errorf("Viewer incorrectly received unclassified event: %s", string(msg))
	}

	// Verify admin DOES receive the unknown event
	select {
	case rxBytes := <-adminChan:
		var rxEvt Event
		_ = json.Unmarshal(rxBytes, &rxEvt)
		if rxEvt.Event != "secret.exposed" {
			t.Errorf("Expected secret.exposed event for admin, got %s", rxEvt.Event)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Admin timeout waiting for unknown event")
	}
}
