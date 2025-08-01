package messages

import (
	"encoding/json"
	"time"
)

// MessageType represents the type of message
type MessageType string

const (
	// Legacy incoming message types (from panel) - maintain backwards compatibility
	TypeSystemInfoRequest MessageType = "system_info_request"
	TypeServerCreate      MessageType = "server_create"
	TypeServerStart       MessageType = "server_start"
	TypeServerStop        MessageType = "server_stop"
	TypeServerRestart     MessageType = "server_restart"
	TypeServerDelete      MessageType = "server_delete"
	TypeServerCommand     MessageType = "server_command"
	TypeFileRead          MessageType = "file_read"
	TypeFileWrite         MessageType = "file_write"

	// Legacy outgoing message types (to panel)
	TypeHeartbeat    MessageType = "heartbeat"
	TypeSystemInfo   MessageType = "system_info"
	TypeServerStatus MessageType = "server_status"
	TypeServerOutput MessageType = "server_output"
	TypeFileContent  MessageType = "file_content"
	TypeError        MessageType = "error"

	// New Panel Issue #27 message types
	TypeCommand  MessageType = "command"
	TypeResponse MessageType = "response"
	TypeEvent    MessageType = "event"
)

// Message represents a WebSocket message
type Message struct {
	Type      MessageType     `json:"type"`
	Data      json.RawMessage `json:"data,omitempty"`
	MessageID string          `json:"messageId,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
}

// HeartbeatData represents heartbeat message data
type HeartbeatData struct {
	NodeID    string    `json:"nodeId"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
}

// SystemInfoData represents system information
type SystemInfoData struct {
	NodeID        string            `json:"nodeId"`
	OS            string            `json:"os"`
	Arch          string            `json:"arch"`
	Memory        int64             `json:"memory"`
	CPU           string            `json:"cpu"`
	DockerVersion string            `json:"dockerVersion"`
	Capabilities  []string          `json:"capabilities"`
	Networks      map[string]string `json:"networks"`
}

// ServerCreateData represents server creation data
type ServerCreateData struct {
	ServerID    string            `json:"serverId"`
	Image       string            `json:"image"`
	Startup     string            `json:"startup"`
	Environment map[string]string `json:"environment"`
	Limits      ResourceLimits    `json:"limits"`
	Ports       []PortMapping     `json:"ports"`
}

// ResourceLimits defines resource constraints
type ResourceLimits struct {
	Memory int64 `json:"memory"` // in bytes
	Swap   int64 `json:"swap"`   // in bytes
	Disk   int64 `json:"disk"`   // in bytes
	IO     int64 `json:"io"`     // IO weight
	CPU    int64 `json:"cpu"`    // CPU shares
}

// PortMapping defines port forwarding
type PortMapping struct {
	Internal int    `json:"internal"`
	External int    `json:"external"`
	Protocol string `json:"protocol"` // tcp/udp
}

// ServerStatusData represents server status information
type ServerStatusData struct {
	ServerID string           `json:"serverId"`
	Status   string           `json:"status"` // running, stopped, error, etc.
	Stats    *ServerStatsData `json:"stats,omitempty"`
	Error    string           `json:"error,omitempty"`
}

// ServerStatsData represents server performance statistics
type ServerStatsData struct {
	CPU     float64 `json:"cpu"`    // CPU usage percentage
	Memory  int64   `json:"memory"` // Memory usage in bytes
	Disk    int64   `json:"disk"`   // Disk usage in bytes
	Network struct {
		In  int64 `json:"in"`  // Network bytes in
		Out int64 `json:"out"` // Network bytes out
	} `json:"network"`
	Players   int       `json:"players"` // Number of players (if applicable)
	Timestamp time.Time `json:"timestamp"`
}

// ServerCommandData represents a command to execute in a server
type ServerCommandData struct {
	ServerID string `json:"serverId"`
	Command  string `json:"command"`
}

// ServerOutputData represents output from a server
type ServerOutputData struct {
	ServerID  string    `json:"serverId"`
	Output    string    `json:"output"`
	Stream    string    `json:"stream"` // stdout/stderr
	Timestamp time.Time `json:"timestamp"`
}

// FileOperationData represents file operation data
type FileOperationData struct {
	ServerID string `json:"serverId"`
	Path     string `json:"path"`
	Content  string `json:"content,omitempty"`
}

// ErrorData represents error information
type ErrorData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// NewMessage creates a new message with the current timestamp
func NewMessage(msgType MessageType, data interface{}) (*Message, error) {
	var rawData json.RawMessage
	var err error

	if data != nil {
		rawData, err = json.Marshal(data)
		if err != nil {
			return nil, err
		}
	}

	return &Message{
		Type:      msgType,
		Data:      rawData,
		Timestamp: time.Now(),
	}, nil
}

// ParseMessage parses a message from JSON bytes
func ParseMessage(data []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	return &msg, err
}

// ToJSON converts a message to JSON bytes
func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalData unmarshals the message data into the provided interface
func (m *Message) UnmarshalData(v interface{}) error {
	if m.Data == nil {
		return nil
	}
	return json.Unmarshal(m.Data, v)
}

// PanelCommand represents the new Panel Issue #27 command format
type PanelCommand struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`      // Always "command"
	Timestamp string                 `json:"timestamp"` // ISO 8601 timestamp
	AgentID   string                 `json:"agentId"`
	Action    string                 `json:"action"`    // start_server, stop_server, etc.
	ServerID  string                 `json:"serverId,omitempty"`
	Payload   map[string]interface{} `json:"payload,omitempty"`
}

// AgentResponse represents the standardized response format
type AgentResponse struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`    // Always "response"
	Timestamp string                 `json:"timestamp"`
	Success   bool                   `json:"success"`
	Message   string                 `json:"message,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Error     *ErrorInfo             `json:"error,omitempty"`
}

// ErrorInfo represents structured error information
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// AgentEvent represents events sent from Agent to Panel
type AgentEvent struct {
	Type      string                 `json:"type"`      // Always "event"
	Timestamp string                 `json:"timestamp"`
	Event     string                 `json:"event"`     // server_status_changed, server_log, etc.
	Data      map[string]interface{} `json:"data"`
}

// ParsePanelCommand parses a PanelCommand from JSON bytes
func ParsePanelCommand(data []byte) (*PanelCommand, error) {
	var cmd PanelCommand
	err := json.Unmarshal(data, &cmd)
	return &cmd, err
}

// ToJSON converts an AgentResponse to JSON bytes
func (r *AgentResponse) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}

// ToJSON converts an AgentEvent to JSON bytes
func (e *AgentEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}
