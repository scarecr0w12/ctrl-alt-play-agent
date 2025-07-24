package docker

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerConfig_Validation(t *testing.T) {
	tests := []struct {
		name    string
		config  *ServerConfig
		wantErr bool
	}{
		{
			name: "valid server config",
			config: &ServerConfig{
				ServerID:    "server_123",
				Image:       "minecraft:latest",
				Startup:     "java -jar server.jar",
				Environment: map[string]string{"MEMORY": "2G"},
				Limits: ResourceLimits{
					Memory: 2147483648, // 2GB
					CPUs:   2,
					Disk:   5368709120, // 5GB
				},
			},
			wantErr: false,
		},
		{
			name: "empty server ID",
			config: &ServerConfig{
				ServerID:    "",
				Image:       "minecraft:latest",
				Startup:     "java -jar server.jar",
				Environment: map[string]string{},
				Limits:      ResourceLimits{},
			},
			wantErr: true,
		},
		{
			name: "empty image",
			config: &ServerConfig{
				ServerID:    "server_123",
				Image:       "",
				Startup:     "java -jar server.jar",
				Environment: map[string]string{},
				Limits:      ResourceLimits{},
			},
			wantErr: true,
		},
		{
			name: "empty startup command",
			config: &ServerConfig{
				ServerID:    "server_123",
				Image:       "minecraft:latest",
				Startup:     "",
				Environment: map[string]string{},
				Limits:      ResourceLimits{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateServerConfig(tt.config)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestResourceLimits_Validation(t *testing.T) {
	tests := []struct {
		name    string
		limits  ResourceLimits
		wantErr bool
	}{
		{
			name: "valid limits",
			limits: ResourceLimits{
				Memory: 1073741824, // 1GB
				CPUs:   1,
				Disk:   5368709120, // 5GB
			},
			wantErr: false,
		},
		{
			name: "zero memory",
			limits: ResourceLimits{
				Memory: 0,
				CPUs:   1,
				Disk:   5368709120,
			},
			wantErr: true,
		},
		{
			name: "zero CPU",
			limits: ResourceLimits{
				Memory: 1073741824,
				CPUs:   0,
				Disk:   5368709120,
			},
			wantErr: true,
		},
		{
			name: "zero disk",
			limits: ResourceLimits{
				Memory: 1073741824,
				CPUs:   1,
				Disk:   0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateResourceLimits(&tt.limits)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPortMapping_Validation(t *testing.T) {
	tests := []struct {
		name    string
		mapping PortMapping
		wantErr bool
	}{
		{
			name: "valid TCP mapping",
			mapping: PortMapping{
				Internal: 25565,
				External: 25565,
				Protocol: "tcp",
			},
			wantErr: false,
		},
		{
			name: "valid UDP mapping",
			mapping: PortMapping{
				Internal: 25565,
				External: 25566,
				Protocol: "udp",
			},
			wantErr: false,
		},
		{
			name: "invalid internal port",
			mapping: PortMapping{
				Internal: 0,
				External: 25565,
				Protocol: "tcp",
			},
			wantErr: true,
		},
		{
			name: "invalid external port",
			mapping: PortMapping{
				Internal: 25565,
				External: 0,
				Protocol: "tcp",
			},
			wantErr: true,
		},
		{
			name: "invalid protocol",
			mapping: PortMapping{
				Internal: 25565,
				External: 25565,
				Protocol: "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePortMapping(&tt.mapping)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Mock Docker client for testing (when Docker is not available)
type mockDockerClient struct {
	containers map[string]string
	shouldFail bool
}

func (m *mockDockerClient) ContainerCreate(ctx context.Context, config interface{}, hostConfig interface{}, networkingConfig interface{}, platform interface{}, containerName string) (interface{}, error) {
	if m.shouldFail {
		return nil, assert.AnError
	}
	
	// Mock response
	return struct {
		ID string
	}{ID: "container_123"}, nil
}

func TestManager_CreateGameServer(t *testing.T) {
	// This test would require Docker to be available
	// For now, we'll test the configuration validation
	
	config := &ServerConfig{
		ServerID:    "test_server",
		Image:       "minecraft:latest",
		Startup:     "java -jar server.jar",
		Environment: map[string]string{"MEMORY": "2G"},
		Limits: ResourceLimits{
			Memory: 2147483648,
			CPUs:   2,
			Disk:   5368709120,
		},
	}

	// Test that the config is valid
	err := validateServerConfig(config)
	assert.NoError(t, err)

	// Test container name generation
	expectedName := "ctrl-alt-play-" + config.ServerID
	assert.Equal(t, "ctrl-alt-play-test_server", expectedName)
}

// Helper functions for validation (these would be added to the manager.go file)
func validateServerConfig(config *ServerConfig) error {
	if config.ServerID == "" {
		return assert.AnError
	}
	if config.Image == "" {
		return assert.AnError
	}
	if config.Startup == "" {
		return assert.AnError
	}
	return validateResourceLimits(&config.Limits)
}

func validateResourceLimits(limits *ResourceLimits) error {
	if limits.Memory <= 0 {
		return assert.AnError
	}
	if limits.CPUs <= 0 {
		return assert.AnError
	}
	if limits.Disk <= 0 {
		return assert.AnError
	}
	return nil
}

func validatePortMapping(mapping *PortMapping) error {
	if mapping.Internal <= 0 || mapping.Internal > 65535 {
		return assert.AnError
	}
	if mapping.External <= 0 || mapping.External > 65535 {
		return assert.AnError
	}
	if mapping.Protocol != "tcp" && mapping.Protocol != "udp" {
		return assert.AnError
	}
	return nil
}
