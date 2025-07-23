# Ctrl-Alt-Play Agent

A lightweight game server management agent that communicates with the Ctrl-Alt-Play panel to manage Docker-based game servers on remote Linux systems.

## Features

- **WebSocket Communication**: Secure real-time communication with the Ctrl-Alt-Play panel
- **Docker Integration**: Full lifecycle management of game server containers
- **System Monitoring**: Real-time resource monitoring and reporting
- **Automatic Reconnection**: Robust connection handling with automatic recovery
- **Security**: Bearer token authentication and secure message handling

## Quick Start

### Prerequisites

- Go 1.21 or later
- Docker Engine installed and running
- Access to a Ctrl-Alt-Play panel instance

### Installation

1. Clone the repository:
```bash
git clone https://github.com/scarecr0w12/ctrl-alt-play-agent.git
cd ctrl-alt-play-agent
```

2. Build the agent:
```bash
make build
```

3. Run the agent:
```bash
# Set required environment variables
export PANEL_URL="ws://your-panel-host:8080"
export NODE_ID="your-node-id"
export AGENT_SECRET="your-secret-token"

# Run the agent
./build/ctrl-alt-play-agent
```

### Configuration

The agent is configured via environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `PANEL_URL` | WebSocket URL of the panel | `ws://localhost:8080` |
| `NODE_ID` | Unique identifier for this node | `node-1` |
| `AGENT_SECRET` | Authentication secret | `agent-secret` |
| `HEALTH_PORT` | Port for health check server | `8081` |

### Health Checks

The agent provides a health check endpoint for monitoring:

```bash
# Check agent health
curl http://localhost:8081/health

# Example response
{
  "status": "healthy",
  "timestamp": "2024-01-15T10:30:45Z",
  "version": "1.0.0",
  "nodeId": "node-1",
  "uptime": "2h15m30s",
  "connected": true
}
```

Status values:

- `healthy`: Agent is running and connected to panel
- `degraded`: Agent is running but not connected to panel

### Docker Deployment

Build and run using Docker:

```bash
# Build Docker image
make docker-build

# Run with environment variables
docker run --rm -it \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e PANEL_URL="ws://your-panel-host:8080" \
  -e NODE_ID="your-node-id" \
  -e AGENT_SECRET="your-secret-token" \
  ctrl-alt-play-agent
```

## Development

### Building

```bash
# Install dependencies
make install

# Build the binary
make build

# Run tests
make test

# Run linter
make lint

# Build for all platforms
make build-all
```

### Project Structure

```
.
├── cmd/agent/          # Main application entry point
├── internal/
│   ├── client/         # WebSocket client and message handling
│   ├── config/         # Configuration management
│   ├── docker/         # Docker container management
│   └── messages/       # Message types and serialization
├── build/              # Build artifacts (created by make)
├── Dockerfile          # Container build configuration
├── Makefile           # Build automation
└── README.md          # This file
```

## Communication Protocol

The agent communicates with the panel using WebSocket messages with the following structure:

```json
{
  "type": "message_type",
  "data": { "...": "..." },
  "timestamp": "2025-01-23T10:30:00Z"
}
```

### Supported Message Types

#### Incoming (from panel):
- `system_info_request` - Request for system information
- `server_create` - Create a new game server container
- `server_start` - Start a game server
- `server_stop` - Stop a game server
- `server_restart` - Restart a game server
- `server_delete` - Delete a game server and its container
- `server_command` - Execute a command in a server container

#### Outgoing (to panel):
- `heartbeat` - Periodic status update
- `system_info` - System specifications and capabilities
- `server_status` - Server status updates
- `server_output` - Console output from servers
- `error` - Error messages

## Security

- All communication uses bearer token authentication
- Docker socket access is required for container management
- The agent runs with minimal required privileges
- Supports non-root container execution

## License

This project is part of the Ctrl-Alt-Play ecosystem. See the main panel repository for licensing information.

## Contributing

Please refer to the main Ctrl-Alt-Play panel repository for contribution guidelines.

## Support

For issues and support, please use the GitHub issues in the main panel repository.
Agent for the Ctrl-Alt-Play panel

## Overview
The `ctrl-alt-play-agent` is designed to manage and run game servers via Docker, providing a seamless integration with the Ctrl-Alt-Play panel system. This agent functions similarly to the "Wings" system found in Pelican Panel and Pterodactyl panel, allowing for efficient server management and real-time communication.

## Project Structure
The project is organized into the following directories and files:

- **cmd/agent/main.go**: Entry point of the agent application. Initializes the application and starts the main event loop.
- **internal/api/router.go**: Defines API routes for managing game servers and handles incoming requests.
- **internal/config/config.go**: Manages configuration settings, loading from files or environment variables.
- **internal/docker/manager.go**: Contains logic for managing Docker containers, including starting and stopping game servers.
- **internal/server/manager.go**: Handles the lifecycle of game servers and interacts with the Docker manager.
- **internal/websocket/handler.go**: Manages WebSocket connections for real-time communication with the Ctrl-Alt-Play panel.
- **go.mod**: Module definition for the Go project, specifying dependencies and their versions.

## Setup Instructions
1. Clone the repository:
   ```
   git clone <repository-url>
   ```
2. Navigate to the project directory:
   ```
   cd ctrl-alt-play-agent
   ```
3. Install dependencies:
   ```
   go mod tidy
   ```
4. Configure the agent by editing the configuration file or setting environment variables as needed.

## Usage
To run the agent, execute the following command:
```
go run cmd/agent/main.go
```
This will start the agent and begin managing game servers as configured.

## Contributing
Contributions are welcome! Please submit a pull request or open an issue for any enhancements or bug fixes.