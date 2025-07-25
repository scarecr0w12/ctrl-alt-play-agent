# Deployment Guide

## Overview

This guide covers deploying the Ctrl-Alt-Play Agent in various environments.

## Prerequisites

- Docker installed and running
- Go 1.23+ (for building from source)
- Network access to the Ctrl-Alt-Play Panel
- Proper firewall configuration

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PANEL_URL` | `ws://localhost:8080` | WebSocket URL of the panel |
| `NODE_ID` | `node-1` | Unique identifier for this agent |
| `AGENT_SECRET` | `agent-secret` | Authentication secret |
| `HEALTH_PORT` | `8081` | Port for health and API endpoints |

## Docker Deployment (Recommended)

### 1. Using Docker Compose

Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  ctrl-alt-play-agent:
    image: ctrl-alt-play-agent:latest
    container_name: ctrl-alt-play-agent
    restart: unless-stopped
    ports:
      - "8081:8081"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    environment:
      - PANEL_URL=ws://your-panel-host:8080
      - NODE_ID=agent-node-1
      - AGENT_SECRET=your-secure-secret
      - HEALTH_PORT=8081
    networks:
      - ctrl-alt-play-network

networks:
  ctrl-alt-play-network:
    external: true
```

Deploy:
```bash
docker-compose up -d
```

### 2. Using Docker Run

```bash
docker run -d \
  --name ctrl-alt-play-agent \
  --restart unless-stopped \
  -p 8081:8081 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e PANEL_URL=ws://your-panel-host:8080 \
  -e NODE_ID=agent-node-1 \
  -e AGENT_SECRET=your-secure-secret \
  -e HEALTH_PORT=8081 \
  ctrl-alt-play-agent:latest
```

## Binary Deployment

### 1. Download Release

```bash
# Download latest release
wget https://github.com/scarecr0w12/ctrl-alt-play-agent/releases/latest/download/ctrl-alt-play-agent-linux-amd64

# Make executable
chmod +x ctrl-alt-play-agent-linux-amd64

# Move to system location
sudo mv ctrl-alt-play-agent-linux-amd64 /usr/local/bin/ctrl-alt-play-agent
```

### 2. Create Systemd Service

Create `/etc/systemd/system/ctrl-alt-play-agent.service`:

```ini
[Unit]
Description=Ctrl-Alt-Play Agent
After=network.target docker.service
Requires=docker.service

[Service]
Type=simple
User=ctrl-alt-play
Group=docker
ExecStart=/usr/local/bin/ctrl-alt-play-agent
Restart=always
RestartSec=5
Environment="PANEL_URL=ws://your-panel-host:8080"
Environment="NODE_ID=agent-node-1"
Environment="AGENT_SECRET=your-secure-secret"
Environment="HEALTH_PORT=8081"

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl daemon-reload
sudo systemctl enable ctrl-alt-play-agent
sudo systemctl start ctrl-alt-play-agent
```

## Building from Source

### 1. Clone Repository

```bash
git clone https://github.com/scarecr0w12/ctrl-alt-play-agent.git
cd ctrl-alt-play-agent
```

### 2. Build Binary

```bash
# Build for current platform
go build -o bin/agent cmd/agent/main.go

# Cross-compile for Linux
GOOS=linux GOARCH=amd64 go build -o bin/agent-linux-amd64 cmd/agent/main.go

# Cross-compile for Windows
GOOS=windows GOARCH=amd64 go build -o bin/agent-windows-amd64.exe cmd/agent/main.go
```

### 3. Build Docker Image

```bash
docker build -t ctrl-alt-play-agent:latest .
```

## Kubernetes Deployment

Create `k8s-deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ctrl-alt-play-agent
  namespace: game-servers
spec:
  replicas: 1
  selector:
    matchLabels:
      app: ctrl-alt-play-agent
  template:
    metadata:
      labels:
        app: ctrl-alt-play-agent
    spec:
      containers:
      - name: agent
        image: ctrl-alt-play-agent:latest
        ports:
        - containerPort: 8081
        env:
        - name: PANEL_URL
          value: "ws://ctrl-alt-play-panel:8080"
        - name: NODE_ID
          value: "k8s-agent-1"
        - name: AGENT_SECRET
          valueFrom:
            secretKeyRef:
              name: agent-secret
              key: secret
        - name: HEALTH_PORT
          value: "8081"
        volumeMounts:
        - name: docker-sock
          mountPath: /var/run/docker.sock
        livenessProbe:
          httpGet:
            path: /health
            port: 8081
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8081
          initialDelaySeconds: 5
          periodSeconds: 5
      volumes:
      - name: docker-sock
        hostPath:
          path: /var/run/docker.sock
---
apiVersion: v1
kind: Service
metadata:
  name: ctrl-alt-play-agent
  namespace: game-servers
spec:
  selector:
    app: ctrl-alt-play-agent
  ports:
  - port: 8081
    targetPort: 8081
  type: ClusterIP
---
apiVersion: v1
kind: Secret
metadata:
  name: agent-secret
  namespace: game-servers
type: Opaque
data:
  secret: <base64-encoded-secret>
```

Deploy:
```bash
kubectl apply -f k8s-deployment.yaml
```

## Security Considerations

### 1. Network Security

- **Firewall Rules**: Only allow necessary ports (8081 for API, panel connection)
- **VPN/Private Network**: Deploy agents in secure network segments
- **TLS/SSL**: Use HTTPS for API endpoints in production

### 2. Docker Security

- **Socket Access**: Agent requires Docker socket access for container management
- **User Permissions**: Run agent with minimal required permissions
- **Container Isolation**: Ensure proper container security policies

### 3. Authentication

- **Strong Secrets**: Use cryptographically strong secrets (32+ characters)
- **Secret Rotation**: Regularly rotate agent secrets
- **Environment Variables**: Use secure secret management (not plain text)

### 4. Monitoring

- **Health Checks**: Monitor agent health endpoint
- **Log Monitoring**: Monitor agent logs for security events
- **Resource Monitoring**: Monitor CPU, memory, and network usage

## Troubleshooting

### Common Issues

#### 1. Connection Refused to Panel

**Symptoms:** `dial tcp: connect: connection refused`

**Solutions:**
- Verify panel is running and accessible
- Check PANEL_URL configuration
- Verify network connectivity
- Check firewall rules

#### 2. Docker Permission Denied

**Symptoms:** `permission denied while trying to connect to Docker daemon`

**Solutions:**
- Add user to docker group: `sudo usermod -aG docker $USER`
- Verify Docker socket permissions
- Check Docker daemon is running

#### 3. Health Check Failures

**Symptoms:** Health endpoint returns 503 or timeouts

**Solutions:**
- Check agent process is running
- Verify port 8081 is not blocked
- Check agent logs for errors
- Verify sufficient system resources

### Log Analysis

View agent logs:
```bash
# Docker deployment
docker logs ctrl-alt-play-agent

# Systemd deployment
sudo journalctl -u ctrl-alt-play-agent -f

# Kubernetes deployment
kubectl logs -n game-servers deployment/ctrl-alt-play-agent
```

### Performance Tuning

#### System Requirements

**Minimum:**
- CPU: 1 core
- RAM: 512MB
- Disk: 1GB
- Network: 10Mbps

**Recommended:**
- CPU: 2+ cores
- RAM: 2GB+
- Disk: 10GB+
- Network: 100Mbps+

#### Configuration Tuning

- **Docker Resources**: Adjust container memory/CPU limits
- **Network Buffers**: Tune network buffer sizes for high throughput
- **File Descriptors**: Increase ulimits for handling many containers

## Monitoring and Alerting

### Health Monitoring

Set up monitoring for:
- Health endpoint availability
- WebSocket connection status
- Docker daemon connectivity
- System resource usage

### Example Monitoring Setup

#### Prometheus Configuration

```yaml
scrape_configs:
  - job_name: 'ctrl-alt-play-agent'
    static_configs:
      - targets: ['agent-host:8081']
    metrics_path: '/health'
    scrape_interval: 30s
```

#### Alert Rules

```yaml
groups:
  - name: ctrl-alt-play-agent
    rules:
      - alert: AgentDown
        expr: up{job="ctrl-alt-play-agent"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Agent is down"
          description: "Ctrl-Alt-Play Agent has been down for more than 1 minute"
      
      - alert: AgentDisconnected
        expr: agent_connected{job="ctrl-alt-play-agent"} == 0
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Agent disconnected from panel"
          description: "Agent has been disconnected from panel for more than 5 minutes"
```

## Backup and Recovery

### Configuration Backup

Important files to backup:
- Environment variables/configuration files
- SSL certificates (if using HTTPS)
- Agent logs (for troubleshooting)

### Disaster Recovery

1. **Agent Failure**: Restart agent service or container
2. **Host Failure**: Deploy agent on new host with same configuration
3. **Network Failure**: Verify network connectivity and DNS resolution
4. **Panel Failure**: Agent will automatically reconnect when panel is restored

Recovery steps:
1. Restore configuration from backup
2. Deploy agent with same NODE_ID
3. Verify connectivity to panel
4. Check container management functionality
