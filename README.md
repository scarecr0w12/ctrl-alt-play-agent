# Ctrl-Alt-Play Agent

[![CI/CD Pipeline](https://github.com/scarecr0w12/ctrl-alt-play-agent/actions/workflows/ci.yml/badge.svg)](https://github.com/scarecr0w12/ctrl-alt-play-agent/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/scarecr0w12/ctrl-alt-play-agent)](https://goreportcard.com/report/github.com/scarecr0w12/ctrl-alt-play-agent)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A lightweight, high-performance remote game server management agent designed to work seamlessly with the [Ctrl-Alt-Play Panel](https://github.com/scarecr0w12/ctrl-alt-play-panel). Built with Go for maximum efficiency and Docker integration for reliable container management.

## 🚀 Features

### Panel Issue #27 Compatible

- **New Unified Command Format**: Full compatibility with Panel's latest command protocol
- **Backwards Compatibility**: Supports legacy message formats for seamless migration
- **Real-time Communication**: WebSocket-based bidirectional communication with Panel
- **Standardized Responses**: Structured response format with error handling

### Container Management

- **Docker Integration**: Complete Docker API integration for game server containers
- **Resource Management**: CPU, memory, and disk usage monitoring and limiting
- **Container Lifecycle**: Create, start, stop, restart, and delete server containers
- **Signal Support**: Graceful shutdowns with SIGTERM/SIGKILL and timeout handling

### Server Operations

- **Multi-Server Support**: Manage multiple game servers on a single node
- **Real-time Status**: Live server status monitoring and reporting
- **Console Access**: Execute commands and stream console output
- **File Management**: Read and write server configuration files

### Security & Reliability

- **JWT Authentication**: Secure Bearer token authentication with Panel
- **Health Monitoring**: Built-in health check endpoint for system monitoring
- **Error Handling**: Comprehensive error reporting with structured error codes
- **Connection Recovery**: Automatic reconnection and heartbeat monitoring

## 🏗️ Architecture

The Agent follows a distributed architecture pattern inspired by Pelican Panel/Wings:

```text
┌─────────────────┐    WebSocket/HTTPS    ┌─────────────────┐
│  Panel (Main)   │ ←----------------→   │  Agent (Node)   │
│                 │                      │                 │
│ - Web UI        │                      │ - Docker Mgmt   │
│ - User Auth     │                      │ - Container     │
│ - Server Config │                      │   Lifecycle     │
│ - Database      │                      │ - Log Streaming │
└─────────────────┘                      └─────────────────┘
```

### Protocol Compatibility

**Panel Issue #27 Command Format (NEW)**:
```json
{
  "id": "cmd_12345_abc",
  "type": "command",
  "timestamp": "2025-01-23T10:00:00Z",
  "agentId": "agent_uuid",
  "action": "start_server",
  "serverId": "server_123"
}
```

**Agent Response Format**:
```json
{
  "id": "cmd_12345_abc",
  "type": "response",
  "timestamp": "2025-01-23T10:00:00Z",
  "success": true,
  "message": "Server started successfully",
  "data": {
    "serverId": "server_123",
    "status": "running"
  }
}
```

## 📦 Installation

### Prerequisites

- **Go 1.24.5+** for building from source
- **Docker Engine** for container management
- **Linux/macOS/Windows** (cross-platform support)

### Quick Start

1. **Download Pre-built Binary**:
   ```bash
   # Linux (amd64)
   wget https://github.com/scarecr0w12/ctrl-alt-play-agent/releases/latest/download/agent-linux-amd64
   chmod +x agent-linux-amd64
   sudo mv agent-linux-amd64 /usr/local/bin/agent
   
   # Or build from source
   git clone https://github.com/scarecr0w12/ctrl-alt-play-agent.git
   cd ctrl-alt-play-agent
   go build -o bin/agent cmd/agent/main.go
   ```

2. **Configure Environment**:
   ```bash
   export PANEL_URL="ws://your-panel-host:8080"
   export NODE_ID="your-unique-node-id"
   export AGENT_SECRET="your-secure-agent-token"
   export HEALTH_PORT="8081"
   ```

3. **Run the Agent**:
   ```bash
   ./bin/agent
   ```

### Docker Deployment

```bash
# Pull the latest image
docker pull scarecr0w12/ctrl-alt-play-agent:latest

# Run with environment variables
docker run -d \
  --name ctrl-alt-play-agent \
  --restart unless-stopped \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e PANEL_URL="ws://panel:8080" \
  -e NODE_ID="docker-node-1" \
  -e AGENT_SECRET="your-agent-secret" \
  -p 8081:8081 \
  scarecr0w12/ctrl-alt-play-agent:latest
```

## ⚙️ Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PANEL_URL` | Panel WebSocket endpoint | `ws://localhost:8080` | ✅ |
| `NODE_ID` | Unique node identifier | `node-1` | ✅ |
| `AGENT_SECRET` | Authentication token | `agent-secret` | ✅ |
| `HEALTH_PORT` | Health check server port | `8081` | ❌ |

### Advanced Configuration

```bash
# Production example
export PANEL_URL="wss://panel.yourdomain.com:8080"
export NODE_ID="prod-node-us-east-1"
export AGENT_SECRET="$(openssl rand -hex 32)"
export HEALTH_PORT="8081"
```

## 🔌 Panel Integration

### Supported Commands

| **Command** | **Description** | **Payload Support** |
|-------------|-----------------|-------------------|
| `start_server` | Start a game server container | ❌ |
| `stop_server` | Stop server with signal/timeout | ✅ Signal, timeout |
| `restart_server` | Graceful restart sequence | ❌ |
| `create_server` | Create new server container | ✅ Full server config |
| `delete_server` | Remove server and cleanup | ❌ |
| `get_status` | Get detailed server status | ❌ |

### Stop Server Example

```json
{
  "id": "cmd_stop_123",
  "type": "command",
  "action": "stop_server",
  "serverId": "minecraft_server_1",
  "payload": {
    "signal": "SIGTERM",
    "timeout": 30
  }
}
```

### Real-time Events

The Agent broadcasts events to keep the Panel synchronized:

```json
{
  "type": "event",
  "timestamp": "2025-01-23T10:05:00Z",
  "event": "server_status_changed",
  "data": {
    "serverId": "minecraft_server_1",
    "previousStatus": "starting",
    "currentStatus": "running",
    "pid": 1234
  }
}
```

## 🧪 Testing

### Run Tests

```bash
# Unit tests
go test ./...

# Integration tests
go test ./cmd/test -v

# Coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Protocol Compatibility Test

```bash
# Test Panel Issue #27 compatibility
./scripts/test-panel-protocol.sh
```

## 🚀 Development

### Building

```bash
# Development build
go build -o bin/agent cmd/agent/main.go

# Production build with optimizations
CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/agent cmd/agent/main.go

# Cross-platform builds
make build-all
```

### Project Structure

```
ctrl-alt-play-agent/
├── cmd/
│   ├── agent/          # Main application entry point
│   └── test/           # Integration tests
├── internal/
│   ├── client/         # Panel WebSocket client
│   ├── config/         # Configuration management
│   ├── docker/         # Docker API integration
│   ├── health/         # Health check server
│   └── messages/       # Protocol message types
├── scripts/            # Build and test scripts
├── .github/workflows/  # CI/CD pipeline
└── docs/              # Additional documentation
```

### Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Commit changes: `git commit -m 'Add amazing feature'`
4. Push to branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

## 📊 Monitoring

### Health Check

The Agent provides a health check endpoint:

```bash
# Check agent health
curl http://localhost:8081/health

# Response
{
  "status": "healthy",
  "timestamp": "2025-01-23T10:00:00Z",
  "uptime": "2h30m45s",
  "panel_connected": true,
  "docker_available": true
}
```

### Logging

```bash
# View logs in real-time
journalctl -u ctrl-alt-play-agent -f

# Or with Docker
docker logs -f ctrl-alt-play-agent
```

## 🔗 Related Projects

- **[Ctrl-Alt-Play Panel](https://github.com/scarecr0w12/ctrl-alt-play-panel)** - Web management interface
- **[Pelican Panel](https://pelican.dev/)** - Inspiration for architecture design

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🤝 Support

- **Documentation**: [Wiki](https://github.com/scarecr0w12/ctrl-alt-play-agent/wiki)
- **Issues**: [GitHub Issues](https://github.com/scarecr0w12/ctrl-alt-play-agent/issues)
- **Discussions**: [GitHub Discussions](https://github.com/scarecr0w12/ctrl-alt-play-agent/discussions)

---

**Made with ❤️ for the gaming community**
