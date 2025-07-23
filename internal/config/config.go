package config

import (
	"os"
)

// Config holds the application configuration
type Config struct {
	PanelURL   string
	NodeID     string
	Secret     string
	HealthPort string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	panelURL := os.Getenv("PANEL_URL")
	if panelURL == "" {
		panelURL = "ws://localhost:8080" // Default for development
	}

	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		nodeID = "node-1" // Default for development
	}

	secret := os.Getenv("AGENT_SECRET")
	if secret == "" {
		secret = "agent-secret" // Default for development
	}

	healthPort := os.Getenv("HEALTH_PORT")
	if healthPort == "" {
		healthPort = "8081" // Default for development
	}

	return &Config{
		PanelURL:   panelURL,
		NodeID:     nodeID,
		Secret:     secret,
		HealthPort: healthPort,
	}, nil
}
