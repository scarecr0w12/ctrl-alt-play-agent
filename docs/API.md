# Ctrl-Alt-Play Agent API Documentation

## Overview

The Ctrl-Alt-Play Agent provides a comprehensive dual communication architecture with both HTTP REST API and WebSocket capabilities for game server management, file operations, and mod management. The agent now supports all commands expected by the Ctrl-Alt-Play Panel.

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

Execute commands on the agent. All commands follow the same request/response format.

**Request Format:**

```json
{
  "action": "command_name",
  "data": {
    "param1": "value1",
    "param2": "value2"
  }
}
```

**Response Format:**

```json
{
  "success": true,
  "data": {
    "result": "data"
  },
  "error": "error message if success is false"
}
```

## Server Lifecycle Commands

These commands provide server-specific management capabilities that the panel expects.

### start_server

Start a game server.

**Parameters:**
- `serverId` (string): The ID of the server to start

**Example Request:**

```json
{
  "action": "start_server",
  "data": {
    "serverId": "minecraft-001"
  }
}
```

**Example Response:**

```json
{
  "success": true,
  "data": {
    "serverId": "minecraft-001",
    "status": "starting",
    "message": "Server start command issued successfully"
  }
}
```

### stop_server

Stop a game server gracefully.

**Parameters:**
- `serverId` (string): The ID of the server to stop

### restart_server

Restart a game server (stop then start).

**Parameters:**
- `serverId` (string): The ID of the server to restart

### kill_server

Forcefully terminate a game server.

**Parameters:**
- `serverId` (string): The ID of the server to kill

### get_server_status

Get the current status of a specific server.

**Parameters:**
- `serverId` (string): The ID of the server to check

**Example Response:**

```json
{
  "success": true,
  "data": {
    "server": {
      "serverId": "minecraft-001",
      "containerId": "abc123",
      "name": "/minecraft-001",
      "status": "running",
      "config": {
        "image": "minecraft:latest",
        "ports": ["25565:25565"],
        "created": "2024-01-20T08:00:00Z"
      }
    }
  }
}
```

### get_server_metrics

Get performance metrics for a specific server.

**Parameters:**
- `serverId` (string): The ID of the server to get metrics for

**Example Response:**

```json
{
  "success": true,
  "data": {
    "metrics": {
      "serverId": "minecraft-001",
      "timestamp": "2024-01-20T10:30:00Z",
      "cpu_usage": "15%",
      "memory_usage": "2.5GB",
      "network_in": "1.2MB",
      "network_out": "850KB",
      "uptime": "2h30m15s",
      "player_count": 0
    }
  }
}
```

### list_servers

List all available servers.

**Parameters:** None

**Example Response:**

```json
{
  "success": true,
  "data": {
    "servers": [
      {
        "serverId": "minecraft-001",
        "containerId": "abc123",
        "name": "/minecraft-001",
        "status": "running",
        "config": {
          "image": "minecraft:latest",
          "ports": ["25565:25565"],
          "created": "2024-01-20T08:00:00Z"
        }
      }
    ],
    "count": 1
  }
}
```

## File Management Commands

These commands provide file system operations within server directories.

### list_files

List files and directories within a server's directory.

**Parameters:**
- `serverId` (string): The ID of the server
- `path` (string, optional): Path within the server directory (default: "/")

**Example Response:**

```json
{
  "success": true,
  "data": {
    "serverId": "minecraft-001",
    "path": "/",
    "files": [
      {
        "name": "server.properties",
        "type": "file",
        "size": 1024,
        "modified": "2024-01-20T09:00:00Z"
      },
      {
        "name": "logs",
        "type": "directory",
        "size": 0,
        "modified": "2024-01-20T08:30:00Z"
      }
    ]
  }
}
```

### read_file

Read the contents of a file within a server's directory.

**Parameters:**
- `serverId` (string): The ID of the server
- `path` (string): Path to the file within the server directory

**Example Response:**

```json
{
  "success": true,
  "data": {
    "serverId": "minecraft-001",
    "path": "server.properties",
    "content": "server-port=25565\nmax-players=20\n...",
    "size": 1024,
    "modified": "2024-01-20T09:00:00Z"
  }
}
```

### write_file

Write content to a file within a server's directory.

**Parameters:**
- `serverId` (string): The ID of the server
- `path` (string): Path to the file within the server directory
- `content` (string): Content to write to the file

### upload_file

Upload a file to a server's directory (binary files supported via base64).

**Parameters:**
- `serverId` (string): The ID of the server
- `path` (string): Path to the file within the server directory
- `content` (string): Base64-encoded file content

### download_file

Download a file from a server's directory.

**Parameters:**
- `serverId` (string): The ID of the server
- `path` (string): Path to the file within the server directory

**Example Response:**

```json
{
  "success": true,
  "data": {
    "serverId": "minecraft-001",
    "path": "server.properties",
    "content": "c2VydmVyLXBvcnQ9MjU1NjU=...",
    "size": 1024,
    "modified": "2024-01-20T09:00:00Z",
    "encoding": "base64"
  }
}
```

## Mod Management Commands

These commands provide mod installation and management capabilities.

### install_mod

Install a mod for a specific server.

**Parameters:**
- `serverId` (string): The ID of the server
- `modId` (string): The ID of the mod to install
- `modUrl` (string, optional): URL to download the mod from
- `version` (string, optional): Version of the mod to install

### uninstall_mod

Uninstall a mod from a specific server.

**Parameters:**
- `serverId` (string): The ID of the server
- `modId` (string): The ID of the mod to uninstall

### list_mods

List all installed mods for a specific server.

**Parameters:**
- `serverId` (string): The ID of the server

**Example Response:**

```json
{
  "success": true,
  "data": {
    "serverId": "minecraft-001",
    "mods": [
      {
        "id": "worldedit",
        "name": "WorldEdit",
        "version": "7.2.12",
        "description": "World editing plugin",
        "enabled": true
      }
    ],
    "count": 1
  }
}
```

## Legacy Docker Commands

For backward compatibility, these Docker commands are still supported:

### docker.list

List all Docker containers.

**Parameters:** None

### docker.start

Start a Docker container.

**Parameters:**
- `containerId` (string): The ID or name of the container to start

### docker.stop

Stop a Docker container.

**Parameters:**
- `containerId` (string): The ID or name of the container to stop

### docker.remove

Remove a Docker container.

**Parameters:**
- `containerId` (string): The ID or name of the container to remove

### docker.inspect

Inspect a Docker container.

**Parameters:**
- `containerId` (string): The ID or name of the container to inspect

## System Commands

### system.ping

Simple ping command for connectivity testing.

**Parameters:** None

**Example Response:**

```json
{
  "success": true,
  "data": {
    "message": "pong",
    "timestamp": "2024-01-20T10:30:00Z"
  }
}
```

### system.status

Get system information including uptime, memory, and disk usage.

**Parameters:** None

**Example Response:**

```json
{
  "success": true,
  "data": {
    "uptime": "up 2 days, 5 hours, 23 minutes",
    "memory": "total: 16GB, used: 8.2GB, free: 6.5GB",
    "disk": "total: 100GB, used: 45GB, available: 50GB"
  }
}
```

## Error Handling

All API responses include a `success` field. When `success` is `false`, an `error` field provides details about what went wrong.

Common error scenarios:

- Missing or invalid authentication (HTTP 401)
- Missing required parameters (HTTP 400)
- Server/container not found
- File system errors
- Docker API errors
- Permission denied for file operations

Example error response:

```json
{
  "success": false,
  "error": "Server minecraft-001 not found"
}
```

## Security Considerations

- All file operations are restricted to the `/opt/gameservers/{serverId}` directory
- Path traversal attacks are prevented by validating paths
- Authentication is required for all API calls
- File uploads are limited and validated

## Rate Limiting

Currently, no rate limiting is implemented, but it may be added in future versions for production deployments.

## WebSocket Integration

The agent maintains a WebSocket connection to the Ctrl-Alt-Play Panel for real-time communication. The WebSocket uses the same authentication mechanism and command format as the HTTP API.
