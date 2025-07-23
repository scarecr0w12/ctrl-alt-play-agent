package docker

import (
	"context"
	"io"
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// Manager handles Docker operations for game servers
type Manager struct {
	client *client.Client
}

// NewManager creates a new Docker manager
func NewManager() (*Manager, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &Manager{
		client: cli,
	}, nil
}

// CreateContainer creates a new container for a game server
func (m *Manager) CreateContainer(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, containerName string) (string, error) {
	resp, err := m.client.ContainerCreate(ctx, config, hostConfig, nil, nil, containerName)
	if err != nil {
		return "", err
	}
	return resp.ID, nil
}

// StartContainer starts a container
func (m *Manager) StartContainer(ctx context.Context, containerID string) error {
	return m.client.ContainerStart(ctx, containerID, container.StartOptions{})
}

// StopContainer stops a container
func (m *Manager) StopContainer(ctx context.Context, containerID string) error {
	timeout := 30 // 30 seconds timeout
	return m.client.ContainerStop(ctx, containerID, container.StopOptions{Timeout: &timeout})
}

// RemoveContainer removes a container
func (m *Manager) RemoveContainer(ctx context.Context, containerID string) error {
	return m.client.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true})
}

// GetContainerStats gets real-time stats for a container
func (m *Manager) GetContainerStats(ctx context.Context, containerID string) (io.ReadCloser, error) {
	stats, err := m.client.ContainerStats(ctx, containerID, false)
	if err != nil {
		return nil, err
	}
	return stats.Body, nil
}

// GetContainerLogs gets logs from a container
func (m *Manager) GetContainerLogs(ctx context.Context, containerID string) (io.ReadCloser, error) {
	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       "100",
	}
	return m.client.ContainerLogs(ctx, containerID, options)
}

// ListContainers lists all containers
func (m *Manager) ListContainers(ctx context.Context) ([]container.Summary, error) {
	return m.client.ContainerList(ctx, container.ListOptions{All: true})
}

// Close closes the Docker client connection
func (m *Manager) Close() error {
	if m.client != nil {
		return m.client.Close()
	}
	return nil
}

// ServerConfig represents configuration for a game server container
type ServerConfig struct {
	ServerID    string            `json:"serverId"`
	Image       string            `json:"image"`
	Startup     string            `json:"startup"`
	Environment map[string]string `json:"environment"`
	Limits      ResourceLimits    `json:"limits"`
	Ports       []PortMapping     `json:"ports"`
}

// ResourceLimits defines resource constraints for containers
type ResourceLimits struct {
	Memory int64 `json:"memory"` // in bytes
	CPUs   int64 `json:"cpu"`    // CPU shares
	Disk   int64 `json:"disk"`   // disk space in bytes
}

// PortMapping defines port forwarding for containers
type PortMapping struct {
	Internal int    `json:"internal"`
	External int    `json:"external"`
	Protocol string `json:"protocol"` // tcp/udp
}

// CreateGameServer creates a game server container based on configuration
func (m *Manager) CreateGameServer(ctx context.Context, config *ServerConfig) (string, error) {
	// Prepare container configuration
	containerConfig := &container.Config{
		Image: config.Image,
		Env:   make([]string, 0, len(config.Environment)),
		Cmd:   []string{"/bin/sh", "-c", config.Startup},
		Labels: map[string]string{
			"ctrl-alt-play.server-id": config.ServerID,
			"ctrl-alt-play.managed":   "true",
		},
	}

	// Add environment variables
	for key, value := range config.Environment {
		containerConfig.Env = append(containerConfig.Env, key+"="+value)
	}

	// Prepare host configuration with resource limits
	hostConfig := &container.HostConfig{
		RestartPolicy: container.RestartPolicy{Name: "unless-stopped"},
		Resources: container.Resources{
			Memory: config.Limits.Memory,
		},
	}

	// Create container name
	containerName := "ctrl-alt-play-" + config.ServerID

	log.Printf("Creating container %s with image %s", containerName, config.Image)

	return m.CreateContainer(ctx, containerConfig, hostConfig, containerName)
}
