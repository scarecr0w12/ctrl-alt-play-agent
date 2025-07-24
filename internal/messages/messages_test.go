package messages

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPanelCommand(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected PanelCommand
		wantErr  bool
	}{
		{
			name: "valid start server command",
			input: `{
				"id": "cmd_12345_abc",
				"type": "command",
				"timestamp": "2025-01-23T10:00:00Z",
				"agentId": "test-agent",
				"action": "start_server",
				"serverId": "server_123"
			}`,
			expected: PanelCommand{
				ID:        "cmd_12345_abc",
				Type:      "command",
				Timestamp: "2025-01-23T10:00:00Z",
				AgentID:   "test-agent",
				Action:    "start_server",
				ServerID:  "server_123",
			},
			wantErr: false,
		},
		{
			name: "stop server command with payload",
			input: `{
				"id": "cmd_67890_def",
				"type": "command",
				"timestamp": "2025-01-23T10:01:00Z",
				"agentId": "test-agent",
				"action": "stop_server",
				"serverId": "server_456",
				"payload": {
					"signal": "SIGTERM",
					"timeout": 30
				}
			}`,
			expected: PanelCommand{
				ID:        "cmd_67890_def",
				Type:      "command",
				Timestamp: "2025-01-23T10:01:00Z",
				AgentID:   "test-agent",
				Action:    "stop_server",
				ServerID:  "server_456",
				Payload: map[string]interface{}{
					"signal":  "SIGTERM",
					"timeout": float64(30),
				},
			},
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			input:   `{"invalid": json}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := ParsePanelCommand([]byte(tt.input))

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected.ID, cmd.ID)
			assert.Equal(t, tt.expected.Type, cmd.Type)
			assert.Equal(t, tt.expected.AgentID, cmd.AgentID)
			assert.Equal(t, tt.expected.Action, cmd.Action)
			assert.Equal(t, tt.expected.ServerID, cmd.ServerID)

			if tt.expected.Payload != nil {
				assert.Equal(t, tt.expected.Payload, cmd.Payload)
			}
		})
	}
}

func TestAgentResponse(t *testing.T) {
	tests := []struct {
		name     string
		response AgentResponse
		expected string
	}{
		{
			name: "success response",
			response: AgentResponse{
				ID:        "cmd_12345_abc",
				Type:      "response",
				Timestamp: "2025-01-23T10:00:00Z",
				Success:   true,
				Message:   "Server started successfully",
				Data: map[string]interface{}{
					"serverId": "server_123",
					"status":   "running",
				},
			},
			expected: `{"id":"cmd_12345_abc","type":"response","timestamp":"2025-01-23T10:00:00Z","success":true,"message":"Server started successfully","data":{"serverId":"server_123","status":"running"}}`,
		},
		{
			name: "error response",
			response: AgentResponse{
				ID:        "cmd_67890_def",
				Type:      "response",
				Timestamp: "2025-01-23T10:01:00Z",
				Success:   false,
				Error: &ErrorInfo{
					Code:    "CONTAINER_NOT_FOUND",
					Message: "Docker container not found",
				},
			},
			expected: `{"id":"cmd_67890_def","type":"response","timestamp":"2025-01-23T10:01:00Z","success":false,"error":{"code":"CONTAINER_NOT_FOUND","message":"Docker container not found"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := tt.response.ToJSON()
			require.NoError(t, err)

			// Parse both as JSON to compare structure
			var expected, actual map[string]interface{}
			err = json.Unmarshal([]byte(tt.expected), &expected)
			require.NoError(t, err)
			err = json.Unmarshal(data, &actual)
			require.NoError(t, err)

			assert.Equal(t, expected, actual)
		})
	}
}

func TestAgentEvent(t *testing.T) {
	event := AgentEvent{
		Type:      "event",
		Timestamp: "2025-01-23T10:00:00Z",
		Event:     "server_status_changed",
		Data: map[string]interface{}{
			"serverId":       "server_123",
			"previousStatus": "stopped",
			"currentStatus":  "running",
		},
	}

	data, err := event.ToJSON()
	require.NoError(t, err)

	var parsed map[string]interface{}
	err = json.Unmarshal(data, &parsed)
	require.NoError(t, err)

	assert.Equal(t, "event", parsed["type"])
	assert.Equal(t, "server_status_changed", parsed["event"])
	assert.Equal(t, "2025-01-23T10:00:00Z", parsed["timestamp"])

	eventData := parsed["data"].(map[string]interface{})
	assert.Equal(t, "server_123", eventData["serverId"])
	assert.Equal(t, "stopped", eventData["previousStatus"])
	assert.Equal(t, "running", eventData["currentStatus"])
}

func TestLegacyMessageCompatibility(t *testing.T) {
	// Test that legacy message format still works
	legacyData := map[string]interface{}{
		"serverId": "server_123",
	}

	msg, err := NewMessage(TypeServerStart, legacyData)
	require.NoError(t, err)

	assert.Equal(t, TypeServerStart, msg.Type)
	assert.NotNil(t, msg.Data)
	assert.NotZero(t, msg.Timestamp)

	// Test parsing legacy message
	jsonData, err := msg.ToJSON()
	require.NoError(t, err)

	parsed, err := ParseMessage(jsonData)
	require.NoError(t, err)

	assert.Equal(t, TypeServerStart, parsed.Type)

	var data map[string]interface{}
	err = parsed.UnmarshalData(&data)
	require.NoError(t, err)
	assert.Equal(t, "server_123", data["serverId"])
}

func TestMessageTypeConstants(t *testing.T) {
	// Test that all message type constants are defined
	expectedTypes := []MessageType{
		// Legacy types
		TypeHeartbeat,
		TypeSystemInfo,
		TypeSystemInfoRequest,
		TypeServerCreate,
		TypeServerStart,
		TypeServerStop,
		TypeServerRestart,
		TypeServerDelete,
		TypeServerCommand,
		TypeServerStatus,
		TypeServerOutput,
		TypeFileRead,
		TypeFileWrite,
		TypeFileContent,
		TypeError,
		// New Panel Issue #27 types
		TypeCommand,
		TypeResponse,
		TypeEvent,
	}

	for _, msgType := range expectedTypes {
		assert.NotEmpty(t, string(msgType), "Message type should not be empty")
	}
}

func TestHeartbeatData(t *testing.T) {
	heartbeat := HeartbeatData{
		NodeID:    "test-node",
		Timestamp: time.Now(),
		Status:    "online",
	}

	msg, err := NewMessage(TypeHeartbeat, heartbeat)
	require.NoError(t, err)

	var parsed HeartbeatData
	err = msg.UnmarshalData(&parsed)
	require.NoError(t, err)

	assert.Equal(t, heartbeat.NodeID, parsed.NodeID)
	assert.Equal(t, heartbeat.Status, parsed.Status)
}

func TestSystemInfoData(t *testing.T) {
	sysInfo := SystemInfoData{
		NodeID:       "test-node",
		OS:           "linux",
		Arch:         "amd64",
		Memory:       8589934592, // 8GB
		CPU:          "4 cores",
		Capabilities: []string{"docker", "monitoring"},
		Networks:     map[string]string{"eth0": "192.168.1.100"},
	}

	msg, err := NewMessage(TypeSystemInfo, sysInfo)
	require.NoError(t, err)

	var parsed SystemInfoData
	err = msg.UnmarshalData(&parsed)
	require.NoError(t, err)

	assert.Equal(t, sysInfo.NodeID, parsed.NodeID)
	assert.Equal(t, sysInfo.OS, parsed.OS)
	assert.Equal(t, sysInfo.Memory, parsed.Memory)
	assert.Equal(t, sysInfo.Capabilities, parsed.Capabilities)
}
