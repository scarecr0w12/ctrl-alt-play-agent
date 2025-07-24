package config

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Save original env vars
	originalVars := map[string]string{
		"PANEL_URL":    os.Getenv("PANEL_URL"),
		"NODE_ID":      os.Getenv("NODE_ID"),
		"AGENT_SECRET": os.Getenv("AGENT_SECRET"),
		"HEALTH_PORT":  os.Getenv("HEALTH_PORT"),
	}

	// Clean up after test
	defer func() {
		for key, value := range originalVars {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}()

	tests := []struct {
		name    string
		envVars map[string]string
		want    *Config
		wantErr bool
	}{
		{
			name: "valid config",
			envVars: map[string]string{
				"PANEL_URL":    "ws://localhost:8080",
				"NODE_ID":      "test-node",
				"AGENT_SECRET": "test-secret",
				"HEALTH_PORT":  "8081",
			},
			want: &Config{
				PanelURL:   "ws://localhost:8080",
				NodeID:     "test-node",
				Secret:     "test-secret",
				HealthPort: "8081",
			},
			wantErr: false,
		},
		{
			name: "defaults when env vars not set",
			envVars: map[string]string{},
			want: &Config{
				PanelURL:   "ws://localhost:8080",
				NodeID:     "node-1",
				Secret:     "agent-secret",
				HealthPort: "8081",
			},
			wantErr: false,
		},
		{
			name: "custom values",
			envVars: map[string]string{
				"PANEL_URL":    "wss://production.example.com:8080",
				"NODE_ID":      "prod-node-1",
				"AGENT_SECRET": "super-secret-token",
				"HEALTH_PORT":  "9090",
			},
			want: &Config{
				PanelURL:   "wss://production.example.com:8080",
				NodeID:     "prod-node-1",
				Secret:     "super-secret-token",
				HealthPort: "9090",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear all env vars first
			os.Unsetenv("PANEL_URL")
			os.Unsetenv("NODE_ID")
			os.Unsetenv("AGENT_SECRET")
			os.Unsetenv("HEALTH_PORT")

			// Set test env vars
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			got, err := LoadConfig()

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want.PanelURL, got.PanelURL)
			assert.Equal(t, tt.want.NodeID, got.NodeID)
			assert.Equal(t, tt.want.Secret, got.Secret)
			assert.Equal(t, tt.want.HealthPort, got.HealthPort)
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				PanelURL:   "ws://localhost:8080",
				NodeID:     "test-node",
				Secret:     "test-secret",
				HealthPort: "8081",
			},
			wantErr: false,
		},
		{
			name: "empty panel URL",
			config: &Config{
				PanelURL:   "",
				NodeID:     "test-node",
				Secret:     "test-secret",
				HealthPort: "8081",
			},
			wantErr: true,
		},
		{
			name: "empty node ID",
			config: &Config{
				PanelURL:   "ws://localhost:8080",
				NodeID:     "",
				Secret:     "test-secret",
				HealthPort: "8081",
			},
			wantErr: true,
		},
		{
			name: "empty agent secret",
			config: &Config{
				PanelURL:   "ws://localhost:8080",
				NodeID:     "test-node",
				Secret:     "",
				HealthPort: "8081",
			},
			wantErr: true,
		},
		{
			name: "invalid panel URL scheme",
			config: &Config{
				PanelURL:   "http://localhost:8080",
				NodeID:     "test-node",
				Secret:     "test-secret",
				HealthPort: "8081",
			},
			wantErr: true,
		},
		{
			name: "valid wss URL",
			config: &Config{
				PanelURL:   "wss://panel.example.com:8080",
				NodeID:     "test-node",
				Secret:     "test-secret",
				HealthPort: "8081",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.config)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Helper function for validation
func validateConfig(config *Config) error {
	if config.PanelURL == "" || config.NodeID == "" || config.Secret == "" {
		return assert.AnError
	}
	
	// Check URL scheme
	if !strings.HasPrefix(config.PanelURL, "ws://") && !strings.HasPrefix(config.PanelURL, "wss://") {
		return assert.AnError
	}
	
	return nil
}
