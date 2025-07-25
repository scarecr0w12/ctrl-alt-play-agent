# Contributing to Ctrl-Alt-Play Agent

Thank you for your interest in contributing to the Ctrl-Alt-Play Agent! This document provides guidelines and information for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Code Style](#code-style)
- [Documentation](#documentation)

## Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code.

## Getting Started

### Prerequisites

- Go 1.23 or later
- Docker and Docker Compose
- Git
- Make (optional, for using Makefile commands)

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork locally:

```bash
git clone https://github.com/your-username/ctrl-alt-play-agent.git
cd ctrl-alt-play-agent
```

3. Add the upstream repository:

```bash
git remote add upstream https://github.com/scarecr0w12/ctrl-alt-play-agent.git
```

## Development Setup

### Local Development

1. Install dependencies:

```bash
go mod download
```

2. Build the project:

```bash
make build
# or
go build -o bin/agent cmd/agent/main.go
```

3. Run tests:

```bash
make test
# or
go test ./...
```

4. Run the agent locally:

```bash
make dev
# or
go run cmd/agent/main.go
```

### Environment Configuration

Create a `.env` file for local development:

```env
PANEL_URL=ws://localhost:8080
NODE_ID=dev-agent
AGENT_SECRET=dev-secret
HEALTH_PORT=8081
```

### Docker Development

Build and run with Docker:

```bash
# Build image
make docker-build

# Run container
make docker-run
```

## Making Changes

### Branching Strategy

- Use descriptive branch names: `feature/add-metrics`, `fix/docker-connection`, `docs/update-api`
- Branch from `main` for new features
- Keep branches focused on a single change

### Commit Messages

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
feat: add new Docker container metrics endpoint
fix: resolve WebSocket connection timeout issues
docs: update API documentation for v1.1.0
refactor: simplify authentication middleware
test: add integration tests for Docker operations
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or modifying tests
- `chore`: Maintenance tasks

### Code Organization

```
ctrl-alt-play-agent/
├── cmd/
│   ├── agent/          # Main application entry point
│   └── test/           # Test utilities
├── internal/
│   ├── api/            # HTTP API server
│   ├── client/         # WebSocket client
│   ├── config/         # Configuration management
│   ├── docker/         # Docker operations
│   ├── health/         # Health monitoring
│   └── messages/       # Message types
├── docs/               # Documentation
├── scripts/            # Build and utility scripts
└── tests/              # Integration tests
```

## Testing

### Unit Tests

Run unit tests for all packages:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Integration Tests

Integration tests require Docker:

```bash
go test -tags=integration ./tests/...
```

### Test Guidelines

- Write unit tests for all new functionality
- Maintain test coverage above 80%
- Use table-driven tests where appropriate
- Mock external dependencies (Docker API, WebSocket connections)
- Include both positive and negative test cases

### Example Test

```go
func TestDockerManager_ListContainers(t *testing.T) {
    tests := []struct {
        name    string
        setup   func(*testing.T) *docker.Manager
        want    int
        wantErr bool
    }{
        {
            name: "successful list",
            setup: func(t *testing.T) *docker.Manager {
                // Mock setup
                return mockDockerManager()
            },
            want:    3,
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            dm := tt.setup(t)
            containers, err := dm.ListContainers(context.Background())
            
            if (err != nil) != tt.wantErr {
                t.Errorf("ListContainers() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            
            if len(containers) != tt.want {
                t.Errorf("ListContainers() = %v containers, want %v", len(containers), tt.want)
            }
        })
    }
}
```

## Submitting Changes

### Pull Request Process

1. Update your branch with the latest main:

```bash
git fetch upstream
git rebase upstream/main
```

2. Push your changes:

```bash
git push origin your-branch-name
```

3. Create a pull request on GitHub

4. Fill out the pull request template completely

5. Wait for review and address feedback

### Pull Request Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] No new warnings or errors
```

### Review Process

- All PRs require at least one approval
- CI/CD pipeline must pass
- Documentation must be updated for new features
- Breaking changes require discussion in issues first

## Code Style

### Go Style Guide

Follow the official [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) and [Effective Go](https://golang.org/doc/effective_go.html).

### Formatting

- Use `gofmt` to format code
- Use `golint` for linting
- Use `go vet` for static analysis

### Naming Conventions

- Use descriptive names for variables and functions
- Use PascalCase for exported functions and types
- Use camelCase for unexported functions and variables
- Use ALL_CAPS for constants

### Comments

- Document all exported functions and types
- Use complete sentences in comments
- Explain complex logic with inline comments

Example:

```go
// ListContainers retrieves all Docker containers from the daemon.
// It returns a slice of container summaries and any error encountered.
func (m *Manager) ListContainers(ctx context.Context) ([]container.Summary, error) {
    // Use ListOptions with All=true to include stopped containers
    return m.client.ContainerList(ctx, container.ListOptions{All: true})
}
```

### Error Handling

- Always handle errors explicitly
- Use meaningful error messages
- Wrap errors with context when appropriate

```go
containers, err := m.client.ContainerList(ctx, options)
if err != nil {
    return nil, fmt.Errorf("failed to list containers: %w", err)
}
```

## Documentation

### API Documentation

- Update `/docs/API.md` for API changes
- Include request/response examples
- Document error codes and responses

### Architecture Documentation

- Update `/docs/ARCHITECTURE.md` for structural changes
- Include diagrams for complex flows
- Explain design decisions

### Deployment Documentation

- Update `/docs/DEPLOYMENT.md` for configuration changes
- Include new environment variables
- Update deployment examples

### README Updates

- Keep feature list current
- Update quick start examples
- Maintain accurate links

## Release Process

### Version Numbering

Follow [Semantic Versioning](https://semver.org/):

- `MAJOR.MINOR.PATCH`
- Major: Breaking changes
- Minor: New features (backward compatible)
- Patch: Bug fixes (backward compatible)

### Creating a Release

1. Update version in relevant files:
   - `VERSION`
   - `package.json`
   - `cmd/agent/main.go` (health server)

2. Update `CHANGELOG.md`

3. Create release PR

4. After merge, create release tag:

```bash
git tag -a v1.2.0 -m "Release v1.2.0"
git push upstream v1.2.0
```

## Getting Help

### Communication Channels

- GitHub Issues: Bug reports and feature requests
- GitHub Discussions: Questions and general discussion
- Pull Request Reviews: Code-specific discussions

### Issue Templates

Use the provided issue templates:

- Bug Report
- Feature Request
- Documentation Update
- Performance Issue

### Asking Questions

When asking questions:

1. Search existing issues first
2. Provide relevant context
3. Include code examples when applicable
4. Specify your environment (OS, Go version, Docker version)

## Recognition

Contributors will be recognized in:

- `CONTRIBUTORS.md` file
- Release notes for significant contributions
- GitHub's contributor statistics

Thank you for contributing to Ctrl-Alt-Play Agent!
