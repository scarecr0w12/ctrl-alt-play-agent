# Ctrl-Alt-Play Agent API Documentation

## Overview

The Ctrl-Alt-Play Agent provides a dual communication architecture with both HTTP REST API and WebSocket capabilities for comprehensive game server management.

## Base Information

- **Version**: 1.1.0
- **Base URL**: `http://localhost:8081`
- **WebSocket URL**: `ws://localhost:8080` (connects to panel)
- **Authentication**: X-API-Key header or Bearer token

## Authentication

### HTTP API
```http
X-API-Key: your-agent-secret-key
```

### WebSocket
```http
Authorization: Bearer your-agent-secret-key
```

## Endpoints

### Health Check

#### GET /health

Returns the current health status of the agent.

**Response:**
```json
{
  "status": "healthy|degraded",
  "timestamp": "2025-07-25T03:48:33.940027974Z",
  "version": "1.1.0",
  "nodeId": "node-1",
  "uptime": "10.248752897s",
  "connected": false
}
```

**Status Codes:**
- `200 OK` - Agent is healthy and connected
- `503 Service Unavailable` - Agent is degraded (not connected to panel)

### Command Execution

#### POST /api/command

Executes commands on the agent.

**Request:**
```json
{
  "action": "command_name",
  "data": {
    "key": "value"
  }
}
```

**Response:**
```json
{
  "success": true|false,
  "data": {
    "result": "data"
  },
  "error": "error message (if success=false)"
}
```

## Available Commands

### System Commands

#### system.ping
Test connectivity to the agent.

**Request:**
```json
{
  "action": "system.ping",
  "data": {}
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "message": "pong",
    "timestamp": "2025-07-25T03:49:03.032081114Z"
  }
}
```

#### system.status
Get comprehensive system information.

**Request:**
```json
{
  "action": "system.status",
  "data": {}
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "uptime": "03:49:59 up 1 day, 26 min, 4 users, load average: 0.07, 0.04, 0.06",
    "memory": "total used free shared buff/cache available...",
    "disk": "Filesystem Size Used Avail Use% Mounted on..."
  }
}
```

### Docker Commands

#### docker.list
List all Docker containers.

**Request:**
```json
{
  "action": "docker.list",
  "data": {}
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "containers": [
      {
        "Id": "container_id",
        "Names": ["/container_name"],
        "Image": "image_name",
        "State": "running|exited",
        "Status": "Up 2 hours",
        "Ports": [...],
        "Labels": {...},
        "Mounts": [...]
      }
    ]
  }
}
```

#### docker.start
Start a Docker container.

**Request:**
```json
{
  "action": "docker.start",
  "data": {
    "containerId": "container_id_or_name"
  }
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "containerId": "container_id",
    "status": "started"
  }
}
```

#### docker.stop
Stop a Docker container.

**Request:**
```json
{
  "action": "docker.stop",
  "data": {
    "containerId": "container_id_or_name"
  }
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "containerId": "container_id",
    "status": "stopped"
  }
}
```

#### docker.remove
Remove a Docker container.

**Request:**
```json
{
  "action": "docker.remove",
  "data": {
    "containerId": "container_id_or_name"
  }
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "containerId": "container_id",
    "status": "removed"
  }
}
```

#### docker.inspect
Get detailed information about a Docker container.

**Request:**
```json
{
  "action": "docker.inspect",
  "data": {
    "containerId": "container_id_or_name"
  }
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "container": {
      "Id": "full_container_id",
      "Names": ["/container_name"],
      "Image": "image_name",
      "State": "running|exited",
      "Status": "detailed_status",
      "Ports": [...],
      "Labels": {...},
      "Mounts": [...],
      "NetworkSettings": {...}
    }
  }
}
```

## Error Responses

All endpoints may return error responses in this format:

```json
{
  "success": false,
  "error": "Error description"
}
```

**Common HTTP Status Codes:**
- `400 Bad Request` - Invalid request format
- `401 Unauthorized` - Missing or invalid authentication
- `404 Not Found` - Endpoint not found
- `405 Method Not Allowed` - Invalid HTTP method
- `500 Internal Server Error` - Server error

## WebSocket Communication

The agent maintains a WebSocket connection to the panel for real-time communication. This is used for:

- Live status updates
- Real-time command execution
- Event streaming
- Connection monitoring

**Connection URL:** `ws://panel-host:8080`
**Authentication:** Bearer token in connection headers

## CORS Support

The API includes CORS headers for browser-based requests:

- `Access-Control-Allow-Origin: *`
- `Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS`
- `Access-Control-Allow-Headers: Content-Type, Authorization, X-API-Key`

## Rate Limiting

Currently, no rate limiting is implemented. In production deployments, consider implementing rate limiting at the reverse proxy level.

## Security Considerations

1. **Authentication Required**: All API endpoints require valid authentication
2. **HTTPS Recommended**: Use HTTPS in production deployments
3. **Secret Management**: Store agent secrets securely
4. **Network Security**: Restrict access to agent ports
5. **Container Security**: Ensure Docker daemon security

## Examples

### cURL Examples

```bash
# Health check
curl -s http://localhost:8081/health

# Ping test
curl -s -H "X-API-Key: agent-secret" \
     -H "Content-Type: application/json" \
     -X POST http://localhost:8081/api/command \
     -d '{"action":"system.ping","data":{}}'

# List containers
curl -s -H "X-API-Key: agent-secret" \
     -H "Content-Type: application/json" \
     -X POST http://localhost:8081/api/command \
     -d '{"action":"docker.list","data":{}}'

# Start container
curl -s -H "X-API-Key: agent-secret" \
     -H "Content-Type: application/json" \
     -X POST http://localhost:8081/api/command \
     -d '{"action":"docker.start","data":{"containerId":"container_name"}}'
```

### JavaScript Examples

```javascript
// Using fetch API
const response = await fetch('http://localhost:8081/api/command', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'X-API-Key': 'agent-secret'
  },
  body: JSON.stringify({
    action: 'system.ping',
    data: {}
  })
});

const result = await response.json();
console.log(result);
```
