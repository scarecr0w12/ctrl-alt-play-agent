package main

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/scarecr0w12/ctrl-alt-play-agent/internal/config"
	"github.com/scarecr0w12/ctrl-alt-play-agent/internal/messages"
)

func TestPanelAgentProtocolCompatibility(t *testing.T) {
	// Test that we can handle Panel Issue #27 command format
	panelCommandJSON := `{
		"id": "cmd_12345_abc",
		"type": "command",
		"timestamp": "2025-01-23T10:00:00Z",
		"agentId": "test-agent",
		"action": "start_server",
		"serverId": "server_123"
	}`

	// Parse Panel command
	cmd, err := messages.ParsePanelCommand([]byte(panelCommandJSON))
	require.NoError(t, err)

	assert.Equal(t, "cmd_12345_abc", cmd.ID)
	assert.Equal(t, "command", cmd.Type)
	assert.Equal(t, "start_server", cmd.Action)
	assert.Equal(t, "server_123", cmd.ServerID)

	// Test Agent response format
	response := &messages.AgentResponse{
		ID:        cmd.ID,
		Type:      "response",
		Timestamp: time.Now().Format(time.RFC3339),
		Success:   true,
		Message:   "start_server command received",
		Data: map[string]interface{}{
			"serverId": cmd.ServerID,
			"status":   "starting",
		},
	}

	responseData, err := response.ToJSON()
	require.NoError(t, err)

	// Verify response can be parsed by Panel
	var parsed map[string]interface{}
	err = json.Unmarshal(responseData, &parsed)
	require.NoError(t, err)

	assert.Equal(t, cmd.ID, parsed["id"])
	assert.Equal(t, "response", parsed["type"])
	assert.Equal(t, true, parsed["success"])
}

func TestLegacyBackwardsCompatibility(t *testing.T) {
	// Test that legacy messages still work
	legacyData := map[string]interface{}{
		"serverId": "server_456",
	}

	msg, err := messages.NewMessage(messages.TypeServerStart, legacyData)
	require.NoError(t, err)

	assert.Equal(t, messages.TypeServerStart, msg.Type)

	var parsedData map[string]interface{}
	err = msg.UnmarshalData(&parsedData)
	require.NoError(t, err)

	assert.Equal(t, "server_456", parsedData["serverId"])
}

func TestConfigurationLoading(t *testing.T) {
	// Test configuration loading works
	cfg, err := config.LoadConfig()
	require.NoError(t, err)

	// Should have defaults
	assert.NotEmpty(t, cfg.PanelURL)
	assert.NotEmpty(t, cfg.NodeID)
	assert.NotEmpty(t, cfg.Secret)
	assert.NotEmpty(t, cfg.HealthPort)
}

func TestMessageTypeDefinitions(t *testing.T) {
	// Test all message types are properly defined
	messageTypes := []messages.MessageType{
		// Legacy types
		messages.TypeHeartbeat,
		messages.TypeSystemInfo,
		messages.TypeSystemInfoRequest,
		messages.TypeServerStart,
		messages.TypeServerStop,
		// New Panel Issue #27 types
		messages.TypeCommand,
		messages.TypeResponse,
		messages.TypeEvent,
	}

	for _, msgType := range messageTypes {
		assert.NotEmpty(t, string(msgType), "Message type should not be empty")
	}
}

func TestEventGeneration(t *testing.T) {
	// Test Agent event generation
	event := &messages.AgentEvent{
		Type:      "event",
		Timestamp: time.Now().Format(time.RFC3339),
		Event:     "server_status_changed",
		Data: map[string]interface{}{
			"serverId":       "server_789",
			"previousStatus": "stopped",
			"currentStatus":  "running",
		},
	}

	eventData, err := event.ToJSON()
	require.NoError(t, err)

	var parsed map[string]interface{}
	err = json.Unmarshal(eventData, &parsed)
	require.NoError(t, err)

	assert.Equal(t, "event", parsed["type"])
	assert.Equal(t, "server_status_changed", parsed["event"])

	data := parsed["data"].(map[string]interface{})
	assert.Equal(t, "server_789", data["serverId"])
	assert.Equal(t, "running", data["currentStatus"])
}
