# Ctrl-Alt-Play Agent

[![CI/CD Pipeline](https://github.com/scarecr0w12/ctrl-alt-play-agent/actions/workflows/ci.yml/badge.svg)](https://github.com/scarecr0w12/ctrl-alt-play-agent/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/scarecr0w12/ctrl-alt-play-agent)](https://goreportcard.com/report/github.com/scarecr0w12/ctrl-alt-play-agent)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Version](https://img.shields.io/badge/version-1.1.0-blue.svg)](VERSION)

A lightweight, high-performance remote game server management agent designed to work seamlessly with the [Ctrl-Alt-Play Panel](https://github.com/scarecr0w12/ctrl-alt-play-panel). Built with Go for maximum efficiency and Docker integration for reliable container management.

**✨ Now fully compatible with Panel ExternalAgentService - supports all server lifecycle, file management, and mod operations expected by the panel.**

## 🚀 Quick Start

### Using Docker (Recommended)

```bash
docker run -d \
  --name ctrl-alt-play-agent \
  --restart unless-stopped \
  -p 8081:8081 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e PANEL_URL=ws://your-panel-host:8080 \
  -e NODE_ID=agent-node-1 \
  -e AGENT_SECRET=your-secure-secret \
  ctrl-alt-play-agent:latest
```

### Download Binary

```bash
# Download latest release
wget https://github.com/scarecr0w12/ctrl-alt-play-agent/releases/latest/download/ctrl-alt-play-agent-linux-amd64

# Make executable and run
chmod +x ctrl-alt-play-agent-linux-amd64
./ctrl-alt-play-agent-linux-amd64
```

## ✨ Features

### Panel Integration Compatibility

- **Full ExternalAgentService Support**: Complete compatibility with Panel's ExternalAgentService API
- **Server Lifecycle Management**: `start_server`, `stop_server`, `restart_server`, `kill_server`, `get_server_status`, `get_server_metrics`, `list_servers`
- **File Management Operations**: `list_files`, `read_file`, `write_file`, `upload_file`, `download_file` with sandboxed security
- **Mod Management System**: `install_mod`, `uninstall_mod`, `list_mods` for game server modifications
- **Dual Command Support**: Supports both modern server-centric commands and legacy Docker commands
- **AgentDiscoveryService Compatible**: Automatic discovery and registration with Panel's AgentDiscoveryService

### Communication & Reliability

- **Dual Architecture**: HTTP API (port 8081) for direct commands + WebSocket client for Panel integration
- **Health Monitoring**: Built-in health check endpoint compatible with Panel discovery service
- **Authentication**: X-API-Key and Bearer token support matching Panel's expectations
- **Error Handling**: Comprehensive error reporting with structured responses
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

### Panel Integration

**Command Format**:

```json
{
  "action": "start_server",
  "data": {
    "serverId": "minecraft-001"
  }
}
```

**Agent Response Format**:

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

| **Command** | **Description** | **Parameters** |
|-------------|-----------------|----------------|
| `start_server` | Start a game server | `serverId` |
| `stop_server` | Stop server gracefully | `serverId` |
| `restart_server` | Restart server (stop + start) | `serverId` |
| `kill_server` | Force terminate server | `serverId` |
| `get_server_status` | Get detailed server status | `serverId` |
| `get_server_metrics` | Get performance metrics | `serverId` |
| `list_servers` | List all servers | none |
| `list_files` | List files in server directory | `serverId`, `path` |
| `read_file` | Read file contents | `serverId`, `path` |
| `write_file` | Write file contents | `serverId`, `path`, `content` |
| `upload_file` | Upload file (base64) | `serverId`, `path`, `content` |
| `download_file` | Download file (base64) | `serverId`, `path` |
| `install_mod` | Install server mod | `serverId`, `modId`, `modUrl`, `version` |
| `uninstall_mod` | Remove server mod | `serverId`, `modId` |
| `list_mods` | List installed mods | `serverId` |

### Legacy Docker Commands

| **Command** | **Description** | **Parameters** |
|-------------|-----------------|----------------|
| `docker.list` | List all containers | none |
| `docker.start` | Start container | `containerId` |
| `docker.stop` | Stop container | `containerId` |
| `docker.remove` | Remove container | `containerId` |
| `docker.inspect` | Inspect container | `containerId` |

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
