# Quick Start Guide

Get MCP Task Tracker up and running in 5 minutes!

## Prerequisites

- Docker and Docker Compose (recommended)
  
OR

- Go 1.21+
- PostgreSQL 15+

## Option 1: Docker (Recommended)

### 1. Clone and Start

```bash
git clone https://github.com/techbuzzz/agent-shaker.git
cd agent-shaker
docker-compose up -d
```

### 2. Verify

```bash
curl http://localhost:8080/health
```

Expected response: `OK`

### 3. Open Web UI

Open your browser: http://localhost:8080

### 4. Run Demo

```bash
./demo.sh
```

This creates:
- 1 project
- 2 agents (backend, frontend)
- 2 tasks
- 1 documentation entry

### 5. Stop Services

```bash
docker-compose down
```

## Option 2: Local Development

### 1. Clone Repository

```bash
git clone https://github.com/techbuzzz/agent-shaker.git
cd agent-shaker
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Start PostgreSQL

```bash
# Using Docker
docker run -d \
  --name mcp-postgres \
  -e POSTGRES_DB=mcp_tracker \
  -e POSTGRES_USER=mcp \
  -e POSTGRES_PASSWORD=secret \
  -p 5432:5432 \
  postgres:15-alpine

# OR install PostgreSQL locally and create database
createdb mcp_tracker
```

### 4. Set Environment

```bash
export DATABASE_URL="postgres://mcp:secret@localhost:5432/mcp_tracker?sslmode=disable"
```

### 5. Build and Run

```bash
make build
./mcp-server
```

Or:

```bash
make run
```

### 6. Verify

```bash
curl http://localhost:8080/health
```

## First Steps

### Create a Project

```bash
curl -X POST http://localhost:8080/api/projects \
  -H "Content-Type: application/json" \
  -d '{
    "name": "My First Project",
    "description": "Testing MCP Task Tracker"
  }'
```

Save the `id` from the response (you'll need it for next steps).

### Register an Agent

```bash
curl -X POST http://localhost:8080/api/agents \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": "YOUR-PROJECT-ID",
    "name": "My Agent",
    "role": "backend",
    "team": "Development Team"
  }'
```

Save the agent `id`.

### Create a Task

```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "project_id": "YOUR-PROJECT-ID",
    "title": "My First Task",
    "description": "Testing task creation",
    "priority": "high",
    "created_by": "YOUR-AGENT-ID",
    "assigned_to": "YOUR-AGENT-ID"
  }'
```

### View in Web UI

1. Open http://localhost:8080
2. You should see your project
3. Click on it to select
4. Switch to "Agents" tab to see your agent
5. Switch to "Tasks" tab to see your task

## Using Makefile

```bash
# Build
make build

# Run locally
make run

# Run tests
make test

# Start with Docker
make docker-up

# Stop Docker services
make docker-down

# View logs
make docker-logs

# Run demo
make demo

# Clean build artifacts
make clean

# Show all commands
make help
```

## Next Steps

1. **Read the Documentation**
   - [README.md](README.md) - Overview and features
   - [docs/API.md](docs/API.md) - Complete API reference
   - [docs/COPILOT_INTEGRATION.md](docs/COPILOT_INTEGRATION.md) - Copilot setup

2. **Try the Demo**
   ```bash
   ./demo.sh
   ```

3. **Integrate with GitHub Copilot**
   - See [docs/COPILOT_INTEGRATION.md](docs/COPILOT_INTEGRATION.md)

4. **Explore the Web UI**
   - Create projects
   - Register agents
   - Create and manage tasks
   - Add documentation

5. **Test WebSocket Updates**
   - Open Web UI in two browser tabs
   - Create a task in one tab
   - Watch it appear in the other tab in real-time

## Common Issues

### Port Already in Use

If port 8080 is already in use:

```bash
# Change port in docker-compose.yml or set environment variable
export PORT=8081
./mcp-server
```

### Database Connection Failed

Check that PostgreSQL is running:

```bash
# Docker
docker ps | grep postgres

# Local
pg_isready
```

### WebSocket Not Connecting

Ensure you're using the correct project ID in the WebSocket URL:

```javascript
ws://localhost:8080/ws?project_id=YOUR-PROJECT-ID
```

## Getting Help

- Check [docs/API.md](docs/API.md) for API details
- See [CONTRIBUTING.md](CONTRIBUTING.md) for development info
- Open an issue on GitHub for bugs or questions

## What's Next?

- **Multi-team Coordination**: Set up multiple agents for different teams
- **Documentation Hub**: Add markdown documentation with tags
- **Real-time Collaboration**: Use WebSocket for live updates
- **GitHub Copilot Integration**: Connect your AI agents

Happy tracking! 🚀
