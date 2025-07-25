package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/scarecr0w12/ctrl-alt-play-agent/internal/config"
	"github.com/scarecr0w12/ctrl-alt-play-agent/internal/docker"
	"github.com/scarecr0w12/ctrl-alt-play-agent/internal/health"
)

// Server provides REST API endpoints for the panel
type Server struct {
	config        *config.Config
	dockerManager *docker.Manager
}

// NewServer creates a new API server
func NewServer(cfg *config.Config, dockerManager *docker.Manager) *Server {
	return &Server{
		config:        cfg,
		dockerManager: dockerManager,
	}
}

// CommandRequest represents a command request from the panel
type CommandRequest struct {
	Action string                 `json:"action"`
	Data   map[string]interface{} `json:"data"`
}

// CommandResponse represents a command response to the panel
type CommandResponse struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data,omitempty"`
	Error   string                 `json:"error,omitempty"`
}

// authenticateRequest validates the request using the agent secret
func (s *Server) authenticateRequest(r *http.Request) bool {
	// Check for API key in header (panel uses this format)
	apiKey := r.Header.Get("X-API-Key")
	if apiKey == s.config.Secret {
		return true
	}

	// Check for Bearer token (WebSocket auth format)
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		return token == s.config.Secret
	}

	return false
}

// CommandHandler handles command requests from the panel
func (s *Server) CommandHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Authenticate the request
		if !s.authenticateRequest(r) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Parse the command request
		var cmdReq CommandRequest
		if err := json.NewDecoder(r.Body).Decode(&cmdReq); err != nil {
			response := CommandResponse{
				Success: false,
				Error:   "Invalid request format",
			}
			s.sendResponse(w, response, http.StatusBadRequest)
			return
		}

		// Execute the command
		response := s.executeCommand(cmdReq)
		s.sendResponse(w, response, http.StatusOK)
	}
}

// executeCommand processes the command request
func (s *Server) executeCommand(req CommandRequest) CommandResponse {
	log.Printf("Executing command: %s", req.Action)

	switch req.Action {
	case "docker.list":
		return s.handleDockerList()
	case "docker.start":
		return s.handleDockerStart(req.Data)
	case "docker.stop":
		return s.handleDockerStop(req.Data)
	case "docker.remove":
		return s.handleDockerRemove(req.Data)
	case "docker.inspect":
		return s.handleDockerInspect(req.Data)
	case "system.status":
		return s.handleSystemStatus()
	case "system.ping":
		return CommandResponse{
			Success: true,
			Data: map[string]interface{}{
				"message":   "pong",
				"timestamp": time.Now(),
			},
		}
	default:
		return CommandResponse{
			Success: false,
			Error:   fmt.Sprintf("Unknown action: %s", req.Action),
		}
	}
}

// Docker command handlers
func (s *Server) handleDockerList() CommandResponse {
	ctx := context.Background()
	containers, err := s.dockerManager.ListContainers(ctx)
	if err != nil {
		return CommandResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	return CommandResponse{
		Success: true,
		Data: map[string]interface{}{
			"containers": containers,
		},
	}
}

func (s *Server) handleDockerStart(data map[string]interface{}) CommandResponse {
	containerID, ok := data["containerId"].(string)
	if !ok {
		return CommandResponse{
			Success: false,
			Error:   "Missing or invalid containerId",
		}
	}

	ctx := context.Background()
	err := s.dockerManager.StartContainer(ctx, containerID)
	if err != nil {
		return CommandResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	return CommandResponse{
		Success: true,
		Data: map[string]interface{}{
			"containerId": containerID,
			"status":      "started",
		},
	}
}

func (s *Server) handleDockerStop(data map[string]interface{}) CommandResponse {
	containerID, ok := data["containerId"].(string)
	if !ok {
		return CommandResponse{
			Success: false,
			Error:   "Missing or invalid containerId",
		}
	}

	ctx := context.Background()
	err := s.dockerManager.StopContainer(ctx, containerID)
	if err != nil {
		return CommandResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	return CommandResponse{
		Success: true,
		Data: map[string]interface{}{
			"containerId": containerID,
			"status":      "stopped",
		},
	}
}

func (s *Server) handleDockerRemove(data map[string]interface{}) CommandResponse {
	containerID, ok := data["containerId"].(string)
	if !ok {
		return CommandResponse{
			Success: false,
			Error:   "Missing or invalid containerId",
		}
	}

	ctx := context.Background()
	err := s.dockerManager.RemoveContainer(ctx, containerID)
	if err != nil {
		return CommandResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	return CommandResponse{
		Success: true,
		Data: map[string]interface{}{
			"containerId": containerID,
			"status":      "removed",
		},
	}
}

func (s *Server) handleDockerInspect(data map[string]interface{}) CommandResponse {
	containerID, ok := data["containerId"].(string)
	if !ok {
		return CommandResponse{
			Success: false,
			Error:   "Missing or invalid containerId",
		}
	}

	// For inspection, we'll get container info from the list
	ctx := context.Background()
	containers, err := s.dockerManager.ListContainers(ctx)
	if err != nil {
		return CommandResponse{
			Success: false,
			Error:   err.Error(),
		}
	}

	// Find the specific container
	for _, container := range containers {
		if container.ID == containerID || strings.HasPrefix(container.ID, containerID) {
			return CommandResponse{
				Success: true,
				Data: map[string]interface{}{
					"container": container,
				},
			}
		}
	}

	return CommandResponse{
		Success: false,
		Error:   "Container not found",
	}
}

func (s *Server) handleSystemStatus() CommandResponse {
	// Get system information
	systemInfo := make(map[string]interface{})

	// Get uptime
	if output, err := exec.Command("uptime").Output(); err == nil {
		systemInfo["uptime"] = strings.TrimSpace(string(output))
	}

	// Get memory info
	if output, err := exec.Command("free", "-h").Output(); err == nil {
		systemInfo["memory"] = strings.TrimSpace(string(output))
	}

	// Get disk info
	if output, err := exec.Command("df", "-h").Output(); err == nil {
		systemInfo["disk"] = strings.TrimSpace(string(output))
	}

	return CommandResponse{
		Success: true,
		Data:    systemInfo,
	}
}

// sendResponse sends a JSON response
func (s *Server) sendResponse(w http.ResponseWriter, response CommandResponse, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// StartServer starts the API server with health endpoint
func (s *Server) StartServer(port string, healthServer *health.Server) error {
	// Add health endpoint
	http.HandleFunc("/health", healthServer.Handler())

	// Add API command endpoint
	http.HandleFunc("/api/command", s.CommandHandler())

	// Add CORS headers for browser requests
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-Key")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.URL.Path == "/api/command" {
			s.CommandHandler()(w, r)
		} else {
			http.NotFound(w, r)
		}
	})

	// Root redirects to health
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/health", http.StatusMovedPermanently)
		} else {
			http.NotFound(w, r)
		}
	})

	log.Printf("Combined API/Health server starting on port %s", port)
	return http.ListenAndServe(":"+port, nil)
}
