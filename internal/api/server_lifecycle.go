package api

import (
	"context"
	"fmt"
	"time"

	"github.com/scarecr0w12/ctrl-alt-play-agent/internal/config"
	"github.com/scarecr0w12/ctrl-alt-play-agent/internal/docker"
)

// ServerLifecycleManager handles server-specific operations that the panel expects
type ServerLifecycleManager struct {
	dockerManager *docker.Manager
	config        *config.Config
}

// ServerInfo represents a game server instance
type ServerInfo struct {
	ServerID    string                 `json:"serverId"`
	ContainerID string                 `json:"containerId"`
	Name        string                 `json:"name"`
	Status      string                 `json:"status"`
	Config      map[string]interface{} `json:"config"`
	Metrics     map[string]interface{} `json:"metrics,omitempty"`
}

// NewServerLifecycleManager creates a new server lifecycle manager
func NewServerLifecycleManager(dockerManager *docker.Manager, config *config.Config) *ServerLifecycleManager {
	return &ServerLifecycleManager{
		dockerManager: dockerManager,
		config:        config,
	}
}

// Server lifecycle commands that the panel expects
func (s *Server) handleStartServer(data map[string]interface{}) CommandResponse {
	serverID, ok := data["serverId"].(string)
	if !ok {
		return CommandResponse{
			Success: false,
			Error:   "Missing or invalid serverId",
		}
	}

	// For now, map serverId to containerID (this could be enhanced with a mapping table)
	containerID := serverID

	ctx := context.Background()
	err := s.dockerManager.StartContainer(ctx, containerID)
	if err != nil {
		return CommandResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to start server %s: %v", serverID, err),
		}
	}

	return CommandResponse{
		Success: true,
		Data: map[string]interface{}{
			"serverId": serverID,
			"status":   "starting",
			"message":  "Server start command issued successfully",
		},
	}
}

func (s *Server) handleStopServer(data map[string]interface{}) CommandResponse {
	serverID, ok := data["serverId"].(string)
	if !ok {
		return CommandResponse{
			Success: false,
			Error:   "Missing or invalid serverId",
		}
	}

	containerID := serverID

	ctx := context.Background()
	err := s.dockerManager.StopContainer(ctx, containerID)
	if err != nil {
		return CommandResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to stop server %s: %v", serverID, err),
		}
	}

	return CommandResponse{
		Success: true,
		Data: map[string]interface{}{
			"serverId": serverID,
			"status":   "stopping",
			"message":  "Server stop command issued successfully",
		},
	}
}

func (s *Server) handleRestartServer(data map[string]interface{}) CommandResponse {
	serverID, ok := data["serverId"].(string)
	if !ok {
		return CommandResponse{
			Success: false,
			Error:   "Missing or invalid serverId",
		}
	}

	containerID := serverID

	ctx := context.Background()

	// Stop the container first
	if err := s.dockerManager.StopContainer(ctx, containerID); err != nil {
		return CommandResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to stop server %s during restart: %v", serverID, err),
		}
	}

	// Wait a moment for graceful shutdown
	time.Sleep(2 * time.Second)

	// Start the container again
	if err := s.dockerManager.StartContainer(ctx, containerID); err != nil {
		return CommandResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to start server %s during restart: %v", serverID, err),
		}
	}

	return CommandResponse{
		Success: true,
		Data: map[string]interface{}{
			"serverId": serverID,
			"status":   "restarting",
			"message":  "Server restart command issued successfully",
		},
	}
}

func (s *Server) handleKillServer(data map[string]interface{}) CommandResponse {
	serverID, ok := data["serverId"].(string)
	if !ok {
		return CommandResponse{
			Success: false,
			Error:   "Missing or invalid serverId",
		}
	}

	containerID := serverID

	ctx := context.Background()

	// Force remove the container (equivalent to kill)
	err := s.dockerManager.RemoveContainer(ctx, containerID)
	if err != nil {
		return CommandResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to kill server %s: %v", serverID, err),
		}
	}

	return CommandResponse{
		Success: true,
		Data: map[string]interface{}{
			"serverId": serverID,
			"status":   "killed",
			"message":  "Server killed successfully",
		},
	}
}

func (s *Server) handleGetServerStatus(data map[string]interface{}) CommandResponse {
	serverID, ok := data["serverId"].(string)
	if !ok {
		return CommandResponse{
			Success: false,
			Error:   "Missing or invalid serverId",
		}
	}

	containerID := serverID

	ctx := context.Background()
	containers, err := s.dockerManager.ListContainers(ctx)
	if err != nil {
		return CommandResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to get server status: %v", err),
		}
	}

	// Find the specific container/server
	for _, container := range containers {
		if container.ID == containerID || container.Names[0] == "/"+serverID {
			serverInfo := ServerInfo{
				ServerID:    serverID,
				ContainerID: container.ID,
				Name:        container.Names[0],
				Status:      container.State,
				Config: map[string]interface{}{
					"image":   container.Image,
					"ports":   container.Ports,
					"created": container.Created,
				},
			}

			return CommandResponse{
				Success: true,
				Data: map[string]interface{}{
					"server": serverInfo,
				},
			}
		}
	}

	return CommandResponse{
		Success: false,
		Error:   fmt.Sprintf("Server %s not found", serverID),
	}
}

func (s *Server) handleGetServerMetrics(data map[string]interface{}) CommandResponse {
	serverID, ok := data["serverId"].(string)
	if !ok {
		return CommandResponse{
			Success: false,
			Error:   "Missing or invalid serverId",
		}
	}

	// For now, return basic metrics - this could be enhanced with real Docker stats
	metrics := map[string]interface{}{
		"serverId":     serverID,
		"timestamp":    time.Now(),
		"cpu_usage":    "0%", // Placeholder - would need Docker stats API
		"memory_usage": "0MB",
		"network_in":   "0B",
		"network_out":  "0B",
		"uptime":       "unknown",
		"player_count": 0, // Game-specific metric
	}

	return CommandResponse{
		Success: true,
		Data: map[string]interface{}{
			"metrics": metrics,
		},
	}
}

func (s *Server) handleListServers() CommandResponse {
	ctx := context.Background()
	containers, err := s.dockerManager.ListContainers(ctx)
	if err != nil {
		return CommandResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to list servers: %v", err),
		}
	}

	var servers []ServerInfo
	for _, container := range containers {
		serverInfo := ServerInfo{
			ServerID:    container.Names[0], // Use container name as server ID
			ContainerID: container.ID,
			Name:        container.Names[0],
			Status:      container.State,
			Config: map[string]interface{}{
				"image":   container.Image,
				"ports":   container.Ports,
				"created": container.Created,
			},
		}
		servers = append(servers, serverInfo)
	}

	return CommandResponse{
		Success: true,
		Data: map[string]interface{}{
			"servers": servers,
			"count":   len(servers),
		},
	}
}
