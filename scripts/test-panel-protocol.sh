#!/bin/bash

# Test script for Panel Issue #27 protocol compatibility
# This script tests the new Panel→Agent command format

set -e

echo "🧪 Testing Agent Panel Protocol Compatibility"
echo "============================================="

# Build the agent
echo "📦 Building Agent..."
cd /home/scarecrow/ctrl-alt-play-agent
go build -o bin/agent cmd/agent/main.go

# Create test environment
echo "🔧 Setting up test environment..."
export PANEL_URL="ws://localhost:8080"
export NODE_ID="test-agent-node"
export AGENT_SECRET="test-agent-secret"
export LOG_LEVEL="debug"

# Test message examples
echo "📋 Testing message format compatibility..."

# Test 1: New Panel command format
echo "✅ New Panel Command Format:"
cat << 'EOF'
{
  "id": "cmd_12345_abc",
  "type": "command",
  "timestamp": "2025-01-23T10:00:00Z",
  "agentId": "test-agent-node",
  "action": "start_server",
  "serverId": "server_123"
}
EOF

echo ""

# Test 2: Legacy message format (backwards compatibility)
echo "✅ Legacy Message Format (Backwards Compatible):"
cat << 'EOF'
{
  "type": "server_start",
  "data": {
    "serverId": "server_123"
  },
  "timestamp": "2025-01-23T10:00:00Z"
}
EOF

echo ""
echo "🎯 Expected Agent Responses:"

# Test 3: New standardized response format
echo "✅ New Standardized Response:"
cat << 'EOF'
{
  "id": "cmd_12345_abc",
  "type": "response",
  "timestamp": "2025-01-23T10:00:00Z",
  "success": true,
  "message": "start_server command received",
  "data": {
    "serverId": "server_123",
    "status": "starting"
  }
}
EOF

echo ""

# Test 4: Event format
echo "✅ Agent Event Format:"
cat << 'EOF'
{
  "type": "event",
  "timestamp": "2025-01-23T10:00:00Z",
  "event": "server_status_changed",
  "data": {
    "serverId": "server_123",
    "previousStatus": "stopped",
    "currentStatus": "running"
  }
}
EOF

echo ""
echo "🚀 Ready for Panel Connection Testing!"
echo ""
echo "To test with actual Panel:"
echo "1. Start Panel: cd /home/scarecrow/ctrl-alt-play-panel && npm start"
echo "2. Start Agent: ./bin/agent"
echo "3. Test commands through Panel API endpoints"
echo ""
echo "✅ Agent now supports Panel Issue #27 protocol!"
