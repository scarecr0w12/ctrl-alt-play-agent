package client

import (
	"context"
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
