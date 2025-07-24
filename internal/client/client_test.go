package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/scarecr0w12/ctrl-alt-play-agent/internal/config"
	"github.com/scarecr0w12/ctrl-alt-play-agent/internal/docker"
	"github.com/scarecr0w12/ctrl-alt-play-agent/internal/messages"
)

// mockDockerManager implements a mock Docker manager for testing
type mockDockerManager struct {
	containers map[string]string // containerName -> status
	errors     map[string]error  // operation -> error
}

func newMockDockerManager() *mockDockerManager {
	return &mockDockerManager{
		containers: make(map[string]string),
		errors:     make(map[string]error),
	}
}

func (m *mockDockerManager) StartContainer(ctx context.Context, containerID string) error {
	if err, exists := m.errors["start"]; exists {
		return err
	}
	m.containers[containerID] = "running"
	return nil
}

func (m *mockDockerManager) StopContainer(ctx context.Context, containerID string) error {
	if err, exists := m.errors["stop"]; exists {
		return err
	}
	m.containers[containerID] = "stopped"
	return nil
}

func (m *mockDockerManager) RestartContainer(ctx context.Context, containerID string) error {
	if err, exists := m.errors["restart"]; exists {
		return err
	}
	m.containers[containerID] = "running"
	return nil
}

func (m *mockDockerManager) DeleteContainer(ctx context.Context, containerID string) error {
	if err, exists := m.errors["delete"]; exists {
		return err
	}
	delete(m.containers, containerID)
	return nil
}

func (m *mockDockerManager) CreateGameServer(ctx context.Context, config *docker.ServerConfig) (string, error) {
	if err, exists := m.errors["create"]; exists {
		return "", err
	}
	containerID := "container_" + time.Now().Format("20060102150405")
	m.containers[config.ServerID] = "created"
	return containerID, nil
}

func (m *mockDockerManager) GetContainerStatus(ctx context.Context, containerID string) (string, error) {
	if err, exists := m.errors["status"]; exists {
		return "", err
	}
	status, exists := m.containers[containerID]
	if !exists {
		return "not_found", nil
	}
	return status, nil
}

func (m *mockDockerManager) ExecuteCommand(ctx context.Context, containerID string, command []string) (string, error) {
	if err, exists := m.errors["execute"]; exists {
		return "", err
	}
	return "command executed", nil
}

func (m *mockDockerManager) RemoveContainer(ctx context.Context, containerID string) error {
	if err, exists := m.errors["remove"]; exists {
		return err
	}
	delete(m.containers, containerID)
	return nil
}

func (m *mockDockerManager) ListContainers(ctx context.Context) ([]interface{}, error) {
	return nil, nil // Simplified for testing
}

func (m *mockDockerManager) CreateContainer(ctx context.Context, config interface{}, hostConfig interface{}, containerName string) (string, error) {
	if err, exists := m.errors["create"]; exists {
		return "", err
	}
	containerID := "container_" + time.Now().Format("20060102150405")
	m.containers[containerName] = "created"
	return containerID, nil
}

func (m *mockDockerManager) GetContainerStats(ctx context.Context, containerID string) (interface{}, error) {
	return nil, nil // Simplified for testing
}

func (m *mockDockerManager) GetContainerLogs(ctx context.Context, containerID string) (interface{}, error) {
	return nil, nil // Simplified for testing
}

func (m *mockDockerManager) Close() error {
	return nil
}

// createTestServer creates a WebSocket test server
func createTestServer(t *testing.T, handler func(*websocket.Conn)) *httptest.Server {
	upgrader := websocket.Upgrader{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		require.NoError(t, err)
		defer conn.Close()

		handler(conn)
	}))

	return server
}

func TestClient_PanelCommandHandling(t *testing.T) {
	tests := []struct {
		name        string
		command     messages.PanelCommand
		expectError bool
		mockError   string
	}{
		{
			name: "start server command",
			command: messages.PanelCommand{
				ID:        "cmd_123",
				Type:      "command",
				Timestamp: time.Now().Format(time.RFC3339),
				AgentID:   "test-agent",
				Action:    "start_server",
				ServerID:  "server_123",
			},
			expectError: false,
		},
		{
			name: "stop server command with payload",
			command: messages.PanelCommand{
				ID:        "cmd_456",
				Type:      "command",
				Timestamp: time.Now().Format(time.RFC3339),
				AgentID:   "test-agent",
				Action:    "stop_server",
				ServerID:  "server_456",
				Payload: map[string]interface{}{
					"signal":  "SIGTERM",
					"timeout": 30,
				},
			},
			expectError: false,
		},
		{
			name: "get status command",
			command: messages.PanelCommand{
				ID:        "cmd_789",
				Type:      "command",
				Timestamp: time.Now().Format(time.RFC3339),
				AgentID:   "test-agent",
				Action:    "get_status",
				ServerID:  "server_789",
			},
			expectError: false,
		},
		{
			name: "unknown action",
			command: messages.PanelCommand{
				ID:        "cmd_unknown",
				Type:      "command",
				Timestamp: time.Now().Format(time.RFC3339),
				AgentID:   "test-agent",
				Action:    "unknown_action",
				ServerID:  "server_unknown",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock docker manager
			mockDocker := newMockDockerManager()
			if tt.mockError != "" {
				mockDocker.errors[strings.Split(tt.mockError, ":")[0]] =
					fmt.Errorf(strings.Split(tt.mockError, ":")[1])
			}

			// Create test server
			var receivedMessages []map[string]interface{}
			server := createTestServer(t, func(conn *websocket.Conn) {
				for {
					_, data, err := conn.ReadMessage()
					if err != nil {
						break
					}

					var msg map[string]interface{}
					json.Unmarshal(data, &msg)
					receivedMessages = append(receivedMessages, msg)
				}
			})
			defer server.Close()

			// Convert HTTP URL to WebSocket URL
			wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

			// Create client
			cfg := &config.Config{
				PanelURL: wsURL,
				NodeID:   "test-node",
				Secret:   "test-secret",
			}

			// Create interface wrapper for mock
			var dockerManager docker.ManagerInterface = mockDocker
			client := NewClient(cfg, dockerManager)

			// Connect to test server
			err = client.Connect()
			require.NoError(t, err)
			defer client.Stop()

			// Create command message
			cmdData, err := json.Marshal(tt.command)
			require.NoError(t, err)

			var rawData json.RawMessage = cmdData
			msg := &messages.Message{
				Type:      messages.TypeCommand,
				Data:      rawData,
				Timestamp: time.Now(),
			}

			// Handle the command
			client.handleMessage(msg)

			// Wait for processing
			time.Sleep(100 * time.Millisecond)

			// Verify response was sent
			assert.GreaterOrEqual(t, len(receivedMessages), 1, "Should receive at least one response")

			if len(receivedMessages) > 0 {
				response := receivedMessages[0]
				assert.Equal(t, "response", response["type"])
				assert.Equal(t, tt.command.ID, response["id"])

				if tt.expectError {
					assert.False(t, response["success"].(bool))
					assert.NotNil(t, response["error"])
				} else {
					assert.True(t, response["success"].(bool))
				}
			}
		})
	}
}

func TestClient_LegacyMessageCompatibility(t *testing.T) {
	// Test that legacy messages still work
	mockDocker := newMockDockerManager()

	cfg := &config.Config{
		PanelURL:    "ws://localhost:8080",
		NodeID:      "test-node",
		AgentSecret: "test-secret",
	}

	client := NewClient(cfg, mockDocker)

	// Create legacy server start message
	legacyData := map[string]interface{}{
		"serverId": "server_123",
	}

	msg, err := messages.NewMessage(messages.TypeServerStart, legacyData)
	require.NoError(t, err)

	// Should not panic and should be handled
	assert.NotPanics(t, func() {
		client.handleMessage(msg)
	})
}

func TestClient_MessageRouting(t *testing.T) {
	mockDocker := newMockDockerManager()

	cfg := &config.Config{
		PanelURL:    "ws://localhost:8080",
		NodeID:      "test-node",
		AgentSecret: "test-secret",
	}

	client := NewClient(cfg, mockDocker)
	client.registerHandlers()

	tests := []struct {
		name        string
		messageType messages.MessageType
		expectRoute bool
	}{
		{
			name:        "command message routed to panel handler",
			messageType: messages.TypeCommand,
			expectRoute: true,
		},
		{
			name:        "legacy server start routed to legacy handler",
			messageType: messages.TypeServerStart,
			expectRoute: true,
		},
		{
			name:        "system info request routed to legacy handler",
			messageType: messages.TypeSystemInfoRequest,
			expectRoute: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &messages.Message{
				Type:      tt.messageType,
				Data:      json.RawMessage(`{}`),
				Timestamp: time.Now(),
			}

			if tt.messageType == messages.TypeCommand {
				// For command messages, we expect handlePanelCommand to be called
				assert.NotPanics(t, func() {
					client.handleMessage(msg)
				})
			} else {
				// For legacy messages, we expect handlers to exist
				client.mu.RLock()
				_, exists := client.handlers[tt.messageType]
				client.mu.RUnlock()

				assert.Equal(t, tt.expectRoute, exists)
			}
		})
	}
}

func TestAgentResponse_Serialization(t *testing.T) {
	response := &messages.AgentResponse{
		ID:        "cmd_123",
		Type:      "response",
		Timestamp: "2025-01-23T10:00:00Z",
		Success:   true,
		Message:   "Test message",
		Data: map[string]interface{}{
			"serverId": "server_123",
			"status":   "running",
		},
	}

	data, err := response.ToJSON()
	require.NoError(t, err)

	// Verify it can be parsed back
	var parsed map[string]interface{}
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, "cmd_123", parsed["id"])
	assert.Equal(t, "response", parsed["type"])
	assert.Equal(t, true, parsed["success"])
	assert.Equal(t, "Test message", parsed["message"])
}

func TestGetActionStatus(t *testing.T) {
	mockDocker := newMockDockerManager()
	cfg := &config.Config{}
	client := NewClient(cfg, mockDocker)

	tests := []struct {
		action   string
		expected string
	}{
		{"start_server", "starting"},
		{"stop_server", "stopping"},
		{"restart_server", "restarting"},
		{"create_server", "creating"},
		{"delete_server", "deleting"},
		{"unknown_action", "processing"},
	}

	for _, tt := range tests {
		t.Run(tt.action, func(t *testing.T) {
			status := client.getActionStatus(tt.action)
			assert.Equal(t, tt.expected, status)
		})
	}
}
