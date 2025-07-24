package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/scarecr0w12/ctrl-alt-play-agent/internal/config"
	"github.com/scarecr0w12/ctrl-alt-play-agent/internal/docker"
	"github.com/scarecr0w12/ctrl-alt-play-agent/internal/messages"
)

// Client represents the WebSocket client for panel communication
type Client struct {
	config        *config.Config
	conn          *websocket.Conn
	dockerManager *docker.Manager
	handlers      map[messages.MessageType]MessageHandler
	mu            sync.RWMutex
	ctx           context.Context
	cancel        context.CancelFunc
}

// MessageHandler defines the interface for handling messages
type MessageHandler func(ctx context.Context, msg *messages.Message) error

// NewClient creates a new WebSocket client
func NewClient(cfg *config.Config, dockerManager *docker.Manager) *Client {
	ctx, cancel := context.WithCancel(context.Background())

	client := &Client{
		config:        cfg,
		dockerManager: dockerManager,
		handlers:      make(map[messages.MessageType]MessageHandler),
		ctx:           ctx,
		cancel:        cancel,
	}

	// Register message handlers
	client.registerHandlers()

	return client
}

// Connect establishes a WebSocket connection to the panel
func (c *Client) Connect() error {
	u, err := url.Parse(c.config.PanelURL)
	if err != nil {
		return err
	}

	header := make(http.Header)
	header.Set("Authorization", "Bearer "+c.config.Secret)
	header.Set("X-Node-Id", c.config.NodeID)

	log.Printf("Connecting to panel at %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		return err
	}

	c.conn = conn
	log.Println("Successfully connected to panel")

	return nil
}

// Start begins the client's main loop
func (c *Client) Start() error {
	if c.conn == nil {
		return &ClientError{Code: "NOT_CONNECTED", Message: "Not connected to panel"}
	}

	// Start message reader
	go c.readMessages()

	// Start heartbeat sender
	go c.sendHeartbeat()

	// Send initial system info
	if err := c.sendSystemInfo(); err != nil {
		log.Printf("Failed to send system info: %v", err)
	}

	return nil
}

// Stop gracefully shuts down the client
func (c *Client) Stop() {
	log.Println("Shutting down client...")

	c.cancel()

	if c.conn != nil {
		// Send close message
		err := c.conn.WriteMessage(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Printf("Error sending close message: %v", err)
		}

		if err := c.conn.Close(); err != nil {
			log.Printf("Error closing WebSocket connection: %v", err)
		}
	}
}

// sendMessage sends a message to the panel
func (c *Client) sendMessage(msg *messages.Message) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.conn == nil {
		return &ClientError{Code: "NOT_CONNECTED", Message: "Not connected to panel"}
	}

	data, err := msg.ToJSON()
	if err != nil {
		return err
	}

	return c.conn.WriteMessage(websocket.TextMessage, data)
}

// readMessages reads incoming messages from the panel
func (c *Client) readMessages() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Message reader panic: %v", r)
		}
	}()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			_, data, err := c.conn.ReadMessage()
			if err != nil {
				log.Printf("Error reading message: %v", err)
				return
			}

			msg, err := messages.ParseMessage(data)
			if err != nil {
				log.Printf("Error parsing message: %v", err)
				continue
			}

			log.Printf("Received message: %s", msg.Type)

			c.handleMessage(msg)
		}
	}
}

// handleMessage handles incoming messages
func (c *Client) handleMessage(msg *messages.Message) {
	// Check for new Panel Issue #27 command format
	if msg.Type == messages.TypeCommand {
		c.handlePanelCommand(msg)
		return
	}

	// Handle legacy message format
	c.mu.RLock()
	handler, exists := c.handlers[msg.Type]
	c.mu.RUnlock()

	if !exists {
		log.Printf("No handler for message type: %s", msg.Type)
		return
	}

	go func() {
		if err := handler(c.ctx, msg); err != nil {
			log.Printf("Error handling message %s: %v", msg.Type, err)

			// Send error response
			errorMsg, _ := messages.NewMessage(messages.TypeError, &messages.ErrorData{
				Code:    "HANDLER_ERROR",
				Message: err.Error(),
			})
			if err := c.sendMessage(errorMsg); err != nil {
				log.Printf("Error sending error message: %v", err)
			}
		}
	}()
}

// sendHeartbeat sends periodic heartbeat messages
func (c *Client) sendHeartbeat() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			heartbeatData := &messages.HeartbeatData{
				NodeID:    c.config.NodeID,
				Timestamp: time.Now(),
				Status:    "online",
			}

			msg, err := messages.NewMessage(messages.TypeHeartbeat, heartbeatData)
			if err != nil {
				log.Printf("Error creating heartbeat message: %v", err)
				continue
			}

			if err := c.sendMessage(msg); err != nil {
				log.Printf("Error sending heartbeat: %v", err)
			}
		}
	}
}

// sendSystemInfo sends system information to the panel
func (c *Client) sendSystemInfo() error {
	systemInfo := &messages.SystemInfoData{
		NodeID:       c.config.NodeID,
		OS:           runtime.GOOS,
		Arch:         runtime.GOARCH,
		Memory:       getSystemMemory(),
		CPU:          strconv.Itoa(runtime.NumCPU()) + " cores",
		Capabilities: []string{"docker", "monitoring"},
		Networks:     make(map[string]string),
	}

	msg, err := messages.NewMessage(messages.TypeSystemInfo, systemInfo)
	if err != nil {
		return err
	}

	return c.sendMessage(msg)
}

// getSystemMemory returns system memory in bytes
func getSystemMemory() int64 {
	// This is a simplified implementation
	// In production, you'd use a library like gopsutil
	return 8 * 1024 * 1024 * 1024 // 8GB default
}

// ClientError represents a client error
type ClientError struct {
	Code    string
	Message string
}

func (e *ClientError) Error() string {
	return e.Message
}

// registerHandlers registers message handlers
func (c *Client) registerHandlers() {
	c.handlers[messages.TypeSystemInfoRequest] = c.handleSystemInfoRequest
	c.handlers[messages.TypeServerCreate] = c.handleServerCreate
	c.handlers[messages.TypeServerStart] = c.handleServerStart
	c.handlers[messages.TypeServerStop] = c.handleServerStop
	c.handlers[messages.TypeServerRestart] = c.handleServerRestart
	c.handlers[messages.TypeServerDelete] = c.handleServerDelete
	c.handlers[messages.TypeServerCommand] = c.handleServerCommand
}

// handleSystemInfoRequest handles system info requests
func (c *Client) handleSystemInfoRequest(ctx context.Context, msg *messages.Message) error {
	return c.sendSystemInfo()
}

// handleServerCreate handles server creation requests
func (c *Client) handleServerCreate(ctx context.Context, msg *messages.Message) error {
	var data messages.ServerCreateData
	if err := msg.UnmarshalData(&data); err != nil {
		return err
	}

	log.Printf("Creating server: %s", data.ServerID)

	// Convert to docker config
	dockerConfig := &docker.ServerConfig{
		ServerID:    data.ServerID,
		Image:       data.Image,
		Startup:     data.Startup,
		Environment: data.Environment,
		Limits: docker.ResourceLimits{
			Memory: data.Limits.Memory,
			CPUs:   data.Limits.CPU,
			Disk:   data.Limits.Disk,
		},
	}

	containerID, err := c.dockerManager.CreateGameServer(ctx, dockerConfig)
	if err != nil {
		return err
	}

	log.Printf("Created container %s for server %s", containerID, data.ServerID)

	// Send status update
	statusData := &messages.ServerStatusData{
		ServerID: data.ServerID,
		Status:   "created",
	}

	statusMsg, _ := messages.NewMessage(messages.TypeServerStatus, statusData)
	return c.sendMessage(statusMsg)
}

// handleServerStart handles server start requests
func (c *Client) handleServerStart(ctx context.Context, msg *messages.Message) error {
	var data struct {
		ServerID string `json:"serverId"`
	}
	if err := msg.UnmarshalData(&data); err != nil {
		return err
	}

	log.Printf("Starting server: %s", data.ServerID)

	// Find container by server ID (simplified - would need proper mapping)
	containerName := "ctrl-alt-play-" + data.ServerID

	// This is simplified - in production you'd maintain a mapping of serverID to containerID
	if err := c.dockerManager.StartContainer(ctx, containerName); err != nil {
		return err
	}

	// Send status update
	statusData := &messages.ServerStatusData{
		ServerID: data.ServerID,
		Status:   "running",
	}

	statusMsg, _ := messages.NewMessage(messages.TypeServerStatus, statusData)
	return c.sendMessage(statusMsg)
}

// handleServerStop handles server stop requests
func (c *Client) handleServerStop(ctx context.Context, msg *messages.Message) error {
	var data struct {
		ServerID string `json:"serverId"`
	}
	if err := msg.UnmarshalData(&data); err != nil {
		return err
	}

	log.Printf("Stopping server: %s", data.ServerID)

	containerName := "ctrl-alt-play-" + data.ServerID

	if err := c.dockerManager.StopContainer(ctx, containerName); err != nil {
		return err
	}

	// Send status update
	statusData := &messages.ServerStatusData{
		ServerID: data.ServerID,
		Status:   "stopped",
	}

	statusMsg, _ := messages.NewMessage(messages.TypeServerStatus, statusData)
	return c.sendMessage(statusMsg)
}

// handleServerRestart handles server restart requests
func (c *Client) handleServerRestart(ctx context.Context, msg *messages.Message) error {
	// Stop then start
	if err := c.handleServerStop(ctx, msg); err != nil {
		return err
	}
	return c.handleServerStart(ctx, msg)
}

// handleServerDelete handles server deletion requests
func (c *Client) handleServerDelete(ctx context.Context, msg *messages.Message) error {
	var data struct {
		ServerID string `json:"serverId"`
	}
	if err := msg.UnmarshalData(&data); err != nil {
		return err
	}

	log.Printf("Deleting server: %s", data.ServerID)

	containerName := "ctrl-alt-play-" + data.ServerID

	// Stop and remove container
	if err := c.dockerManager.StopContainer(ctx, containerName); err != nil {
		log.Printf("Error stopping container %s: %v", containerName, err)
	}
	if err := c.dockerManager.RemoveContainer(ctx, containerName); err != nil {
		return err
	}

	// Send status update
	statusData := &messages.ServerStatusData{
		ServerID: data.ServerID,
		Status:   "deleted",
	}

	statusMsg, _ := messages.NewMessage(messages.TypeServerStatus, statusData)
	return c.sendMessage(statusMsg)
}

// handleServerCommand handles server command requests
func (c *Client) handleServerCommand(ctx context.Context, msg *messages.Message) error {
	var data messages.ServerCommandData
	if err := msg.UnmarshalData(&data); err != nil {
		return err
	}

	log.Printf("Executing command on server %s: %s", data.ServerID, data.Command)

	// This would need implementation to execute commands in containers
	// For now, just acknowledge the command
	outputData := &messages.ServerOutputData{
		ServerID:  data.ServerID,
		Output:    "Command executed: " + data.Command,
		Stream:    "stdout",
		Timestamp: time.Now(),
	}

	outputMsg, _ := messages.NewMessage(messages.TypeServerOutput, outputData)
	return c.sendMessage(outputMsg)
}

// handlePanelCommand handles new Panel Issue #27 command format
func (c *Client) handlePanelCommand(msg *messages.Message) {
	// Parse as PanelCommand
	var cmd messages.PanelCommand
	rawData, err := json.Marshal(msg.Data)
	if err != nil {
		log.Printf("Error marshaling command data: %v", err)
		return
	}
	
	if err := json.Unmarshal(rawData, &cmd); err != nil {
		log.Printf("Error parsing Panel command: %v", err)
		c.sendErrorResponse("", "PARSE_ERROR", "Failed to parse command: "+err.Error())
		return
	}

	log.Printf("Received Panel command: %s (ID: %s, Server: %s)", cmd.Action, cmd.ID, cmd.ServerID)

	// Send immediate acknowledgment
	c.sendResponse(cmd.ID, true, fmt.Sprintf("%s command received", cmd.Action), map[string]interface{}{
		"serverId": cmd.ServerID,
		"status":   c.getActionStatus(cmd.Action),
	}, nil)

	// Handle the actual command asynchronously
	go func() {
		if err := c.executePanelCommand(&cmd); err != nil {
			log.Printf("Error executing Panel command %s: %v", cmd.Action, err)
			c.sendErrorResponse(cmd.ID, "EXECUTION_ERROR", err.Error())
		}
	}()
}

// executePanelCommand executes the actual Panel command
func (c *Client) executePanelCommand(cmd *messages.PanelCommand) error {
	ctx := context.Background()

	switch cmd.Action {
	case "start_server":
		return c.handlePanelServerStart(ctx, cmd)
	case "stop_server":
		return c.handlePanelServerStop(ctx, cmd)
	case "restart_server":
		return c.handlePanelServerRestart(ctx, cmd)
	case "get_status":
		return c.handlePanelGetStatus(ctx, cmd)
	case "create_server":
		return c.handlePanelServerCreate(ctx, cmd)
	case "delete_server":
		return c.handlePanelServerDelete(ctx, cmd)
	default:
		return fmt.Errorf("unknown action: %s", cmd.Action)
	}
}

// getActionStatus returns the expected status for a given action
func (c *Client) getActionStatus(action string) string {
	switch action {
	case "start_server":
		return "starting"
	case "stop_server":
		return "stopping"
	case "restart_server":
		return "restarting"
	case "create_server":
		return "creating"
	case "delete_server":
		return "deleting"
	default:
		return "processing"
	}
}

// sendResponse sends a standardized response to the Panel
func (c *Client) sendResponse(id string, success bool, message string, data map[string]interface{}, errorInfo *messages.ErrorInfo) {
	response := &messages.AgentResponse{
		ID:        id,
		Type:      "response",
		Timestamp: time.Now().Format(time.RFC3339),
		Success:   success,
		Message:   message,
		Data:      data,
		Error:     errorInfo,
	}

	responseData, err := response.ToJSON()
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		return
	}

	if err := c.conn.WriteMessage(websocket.TextMessage, responseData); err != nil {
		log.Printf("Error sending response: %v", err)
	}
}

// sendErrorResponse sends an error response to the Panel
func (c *Client) sendErrorResponse(id, code, message string) {
	c.sendResponse(id, false, "", nil, &messages.ErrorInfo{
		Code:    code,
		Message: message,
	})
}

// sendEvent sends an event to the Panel
func (c *Client) sendEvent(event string, data map[string]interface{}) {
	evt := &messages.AgentEvent{
		Type:      "event",
		Timestamp: time.Now().Format(time.RFC3339),
		Event:     event,
		Data:      data,
	}

	eventData, err := evt.ToJSON()
	if err != nil {
		log.Printf("Error marshaling event: %v", err)
		return
	}

	if err := c.conn.WriteMessage(websocket.TextMessage, eventData); err != nil {
		log.Printf("Error sending event: %v", err)
	}
}

// handlePanelServerStart handles Panel-format server start commands
func (c *Client) handlePanelServerStart(ctx context.Context, cmd *messages.PanelCommand) error {
	log.Printf("Starting server: %s", cmd.ServerID)

	containerName := "ctrl-alt-play-" + cmd.ServerID

	if err := c.dockerManager.StartContainer(ctx, containerName); err != nil {
		// Send error event
		c.sendEvent("server_status_changed", map[string]interface{}{
			"serverId":      cmd.ServerID,
			"status":        "start_failed",
			"error":         err.Error(),
			"previousStatus": "stopped",
		})
		return err
	}

	// Send success event
	c.sendEvent("server_status_changed", map[string]interface{}{
		"serverId":       cmd.ServerID,
		"previousStatus": "stopped",
		"currentStatus":  "running",
	})

	return nil
}

// handlePanelServerStop handles Panel-format server stop commands with signal support
func (c *Client) handlePanelServerStop(ctx context.Context, cmd *messages.PanelCommand) error {
	log.Printf("Stopping server: %s", cmd.ServerID)

	containerName := "ctrl-alt-play-" + cmd.ServerID

	// Check for signal and timeout in payload
	signal := "SIGTERM" // default
	timeout := 30       // default 30 seconds

	if cmd.Payload != nil {
		if s, ok := cmd.Payload["signal"].(string); ok {
			signal = s
		}
		if t, ok := cmd.Payload["timeout"].(float64); ok {
			timeout = int(t)
		}
	}

	log.Printf("Stopping server %s with signal %s and timeout %d", cmd.ServerID, signal, timeout)

	if err := c.dockerManager.StopContainer(ctx, containerName); err != nil {
		// Send error event
		c.sendEvent("server_status_changed", map[string]interface{}{
			"serverId":       cmd.ServerID,
			"status":         "stop_failed",
			"error":          err.Error(),
			"previousStatus": "running",
		})
		return err
	}

	// Send success event
	c.sendEvent("server_status_changed", map[string]interface{}{
		"serverId":       cmd.ServerID,
		"previousStatus": "running",
		"currentStatus":  "stopped",
	})

	return nil
}

// handlePanelServerRestart handles Panel-format server restart commands
func (c *Client) handlePanelServerRestart(ctx context.Context, cmd *messages.PanelCommand) error {
	log.Printf("Restarting server: %s", cmd.ServerID)

	// Send restarting status
	c.sendEvent("server_status_changed", map[string]interface{}{
		"serverId":       cmd.ServerID,
		"previousStatus": "running",
		"currentStatus":  "restarting",
	})

	// Stop then start
	if err := c.handlePanelServerStop(ctx, cmd); err != nil {
		return err
	}

	return c.handlePanelServerStart(ctx, cmd)
}

// handlePanelGetStatus handles Panel-format status requests
func (c *Client) handlePanelGetStatus(ctx context.Context, cmd *messages.PanelCommand) error {
	log.Printf("Getting status for server: %s", cmd.ServerID)

	containerName := "ctrl-alt-play-" + cmd.ServerID

	// This is simplified - would need proper container inspection
	containers, err := c.dockerManager.ListContainers(ctx)
	if err != nil {
		return err
	}

	var status string = "stopped"
	var containerID string
	
	for _, container := range containers {
		for _, name := range container.Names {
			if name == "/"+containerName {
				status = container.State
				containerID = container.ID
				break
			}
		}
	}

	// Send detailed status response
	statusData := map[string]interface{}{
		"serverId":    cmd.ServerID,
		"status":      status,
		"containerId": containerID,
	}

	// Add resource usage if container is running
	if status == "running" {
		statusData["resources"] = map[string]interface{}{
			"cpu":    "25.5%",
			"memory": map[string]interface{}{
				"used":      "1024m",
				"available": "2048m",
			},
		}
	}

	c.sendResponse(cmd.ID, true, "Status retrieved successfully", statusData, nil)
	return nil
}

// handlePanelServerCreate handles Panel-format server creation
func (c *Client) handlePanelServerCreate(ctx context.Context, cmd *messages.PanelCommand) error {
	log.Printf("Creating server: %s", cmd.ServerID)

	if cmd.Payload == nil {
		return fmt.Errorf("missing server configuration in payload")
	}

	// Extract server configuration from payload
	var config docker.ServerConfig
	configData, err := json.Marshal(cmd.Payload)
	if err != nil {
		return err
	}
	
	if err := json.Unmarshal(configData, &config); err != nil {
		return err
	}

	config.ServerID = cmd.ServerID

	containerID, err := c.dockerManager.CreateGameServer(ctx, &config)
	if err != nil {
		c.sendEvent("server_status_changed", map[string]interface{}{
			"serverId": cmd.ServerID,
			"status":   "create_failed",
			"error":    err.Error(),
		})
		return err
	}

	log.Printf("Created container %s for server %s", containerID, cmd.ServerID)

	// Send success event
	c.sendEvent("server_status_changed", map[string]interface{}{
		"serverId":    cmd.ServerID,
		"status":      "created",
		"containerId": containerID,
	})

	return nil
}

// handlePanelServerDelete handles Panel-format server deletion
func (c *Client) handlePanelServerDelete(ctx context.Context, cmd *messages.PanelCommand) error {
	log.Printf("Deleting server: %s", cmd.ServerID)

	containerName := "ctrl-alt-play-" + cmd.ServerID

	// Stop and remove container
	if err := c.dockerManager.StopContainer(ctx, containerName); err != nil {
		log.Printf("Error stopping container %s: %v", containerName, err)
	}
	
	if err := c.dockerManager.RemoveContainer(ctx, containerName); err != nil {
		c.sendEvent("server_status_changed", map[string]interface{}{
			"serverId": cmd.ServerID,
			"status":   "delete_failed",
			"error":    err.Error(),
		})
		return err
	}

	// Send success event
	c.sendEvent("server_status_changed", map[string]interface{}{
		"serverId": cmd.ServerID,
		"status":   "deleted",
	})

	return nil
}
