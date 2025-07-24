# System Patterns

## Panel+Agent Distributed Architecture

Comprehensive distributed system where the Panel acts as the central control plane and Agents execute commands on remote nodes. The architecture enables centralized management of game servers across multiple Linux systems while maintaining real-time communication and monitoring.

### Key Components

- **Panel**: Central web interface and API server
- **Agent**: Lightweight daemon on remote nodes
- **WebSocket**: Real-time bidirectional communication
- **Docker API**: Container lifecycle management
- **Authentication**: Bearer token security

## Protocol Evolution (Issue #27 Breaking Changes)

The Panel implemented critical breaking changes in Issue #27, introducing a new unified command format. This requires Agent updates to handle both new and legacy protocols during migration.

### Protocol Patterns

**New Unified Format** (Post Issue #27):
```json
{
  "id": "cmd_123",
  "type": "command", 
  "action": "start_server",
  "serverId": "server_456",
  "timestamp": "2025-07-24T10:00:00Z"
}
```

**Legacy Format** (Pre Issue #27):
```json
{
  "type": "server_start",
  "data": {"serverId": "server_456"}
}
```

**Agent Response Format**:
```json
{
  "id": "cmd_123",
  "type": "response",
  "success": true,
  "message": "Server started successfully",
  "data": {"status": "running"}
}
```

**Error Response Format**:
```json
{
  "id": "cmd_123", 
  "type": "response",
  "success": false,
  "error": {
    "code": "CONTAINER_NOT_FOUND",
    "message": "Container not found"
  }
}
```

## Communication Patterns

### WebSocket Connection
- **Authentication**: `Authorization: Bearer <agent_secret>`
- **Heartbeat**: Every 30 seconds for connection health
- **Reconnection**: Automatic with exponential backoff
- **Message ID Tracking**: Request/response correlation

### Docker Integration
- **API Client**: `docker.NewClientWithOpts(docker.FromEnv)`
- **Container Lifecycle**: Create, start, stop, restart, delete
- **Resource Monitoring**: CPU, memory, network stats
- **Log Streaming**: Real-time container output

### Health Monitoring
- **HTTP Endpoint**: `GET /health` returns JSON status
- **Status Reporting**: System resources and connection state
- **Graceful Shutdown**: `context.WithCancel()` for cleanup

## Implementation Requirements

### Dual Protocol Support
- New `handlePanelCommand()` method for Issue #27 format
- Legacy message handlers for backwards compatibility
- Structured responses with success/error fields
- Proper message ID tracking for request/response matching

### Security Patterns
- Bearer token validation for all Panel communications
- Secure Docker socket access with minimal privileges
- Input validation and sanitization for all commands
- Error message sanitization to prevent information leakage


## Comprehensive Integration Testing Strategy

Multi-layered testing approach that validates both component isolation and system integration. Combines unit tests for core logic with live system testing for real-world validation. Pattern prioritizes functional verification over test coverage metrics when interface compatibility issues arise.

### Examples

- Unit tests for messages, config, and docker modules
- Integration tests for protocol compatibility
- Live panel connection testing with actual WebSocket communication
- Health endpoint validation with HTTP requests
- Error handling verification through timeout scenarios
