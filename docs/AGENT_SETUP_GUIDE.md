# Agent Setup Guide - MCP Application Integration

## Overview

This guide explains how to connect VS Code AI agents to the MCP (Multi-agent Coordination Protocol) application, enabling collaborative development where agents can coordinate their work, track tasks, and maintain context across a project.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Starting the MCP Server](#starting-the-mcp-server)
3. [Creating a Project](#creating-a-project)
4. [Registering an Agent](#registering-an-agent)
5. [Agent Usage Examples](#agent-usage-examples)
6. [API Reference](#api-reference)
7. [Troubleshooting](#troubleshooting)

---

## Prerequisites

- Docker and Docker Compose installed
- VS Code with GitHub Copilot or similar AI assistant
- Basic understanding of REST APIs
- curl, Postman, or similar tool for API testing

---

## Starting the MCP Server

### 1. Start the Application

```bash
# Navigate to the project directory
cd c:\Sources\GitHub\agent-shaker

# Start all services (PostgreSQL + MCP Server)
docker-compose up -d

# Verify services are running
docker-compose ps
```

**Expected Output:**
```
NAME                        STATUS                   PORTS
agent-shaker-mcp-server-1   Up                      0.0.0.0:8080->8080/tcp
agent-shaker-postgres-1     Up (healthy)            0.0.0.0:5433->5432/tcp
```

### 2. Verify Health

```bash
curl http://localhost:8080/health
# Response: OK
```

---

## Creating a Project

Before registering agents, create a project they will work on.

### Using curl (PowerShell)

```powershell
$project = @{
    name = "InvoiceAI"
    description = "AI-powered invoice processing system with modern frontend"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/projects" `
    -Method POST `
    -ContentType "application/json" `
    -Body $project
```

### Using curl (bash/Linux/Mac)

```bash
curl -X POST http://localhost:8080/api/projects \
  -H "Content-Type: application/json" \
  -d '{
    "name": "InvoiceAI",
    "description": "AI-powered invoice processing system with modern frontend"
  }'
```

**Response:**
```json
{
  "id": "2da7e7e6-0f8c-4f33-9ae8-f3dfb6f0ec8a",
  "name": "InvoiceAI",
  "description": "AI-powered invoice processing system with modern frontend",
  "status": "active",
  "created_at": "2026-01-21T10:00:00Z",
  "updated_at": "2026-01-21T10:00:00Z"
}
```

**📝 Save the `id` value - you'll need it to register agents!**

---

## Registering an Agent

### Example: Registering a Frontend Agent

Let's register an agent named "InvoiceAI-Frontend" who will work on the frontend components.

#### PowerShell

```powershell
# Replace with your actual project_id from the previous step
$projectId = "2da7e7e6-0f8c-4f33-9ae8-f3dfb6f0ec8a"

$agent = @{
    project_id = $projectId
    name = "InvoiceAI-Frontend"
    role = "frontend"
    team = "UI Development"
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "http://localhost:8080/api/agents" `
    -Method POST `
    -ContentType "application/json" `
    -Body $agent

# Display the agent details
$response | ConvertTo-Json
```

#### Bash/Linux/Mac

```bash
PROJECT_ID="2da7e7e6-0f8c-4f33-9ae8-f3dfb6f0ec8a"

curl -X POST http://localhost:8080/api/agents \
  -H "Content-Type: application/json" \
  -d "{
    \"project_id\": \"$PROJECT_ID\",
    \"name\": \"InvoiceAI-Frontend\",
    \"role\": \"frontend\",
    \"team\": \"UI Development\"
  }"
```

**Response:**
```json
{
  "id": "8f3e7b1c-4d2a-4c8f-9e7a-b5c3d1e9f8a2",
  "project_id": "2da7e7e6-0f8c-4f33-9ae8-f3dfb6f0ec8a",
  "name": "InvoiceAI-Frontend",
  "role": "frontend",
  "team": "UI Development",
  "status": "active",
  "last_seen": "2026-01-21T10:05:00Z",
  "created_at": "2026-01-21T10:05:00Z"
}
```

**📝 Save the agent `id` for task management!**

### Example: Registering a Backend Agent

```powershell
# PowerShell
$agent = @{
    project_id = $projectId
    name = "InvoiceAI-Backend"
    role = "backend"
    team = "API Development"
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/agents" `
    -Method POST `
    -ContentType "application/json" `
    -Body $agent
```

---

## Agent Usage Examples

### Scenario: Frontend Agent Working on Invoice UI

#### 1. Create a Task for the Frontend Agent

```powershell
$task = @{
    project_id = $projectId
    agent_id = "8f3e7b1c-4d2a-4c8f-9e7a-b5c3d1e9f8a2"  # InvoiceAI-Frontend
    title = "Create Invoice List Component"
    description = "Build a React component to display list of invoices with filtering and sorting"
    priority = "high"
    dependencies = @()
} | ConvertTo-Json

$taskResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/tasks" `
    -Method POST `
    -ContentType "application/json" `
    -Body $task

$taskId = $taskResponse.id
```

#### 2. Update Task Status as Agent Works

```powershell
# Mark as in progress
$status = @{ status = "in_progress" } | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/tasks/$taskId/status" `
    -Method PUT `
    -ContentType "application/json" `
    -Body $status

# Later, mark as done
$status = @{ status = "done" } | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/tasks/$taskId/status" `
    -Method PUT `
    -ContentType "application/json" `
    -Body $status
```

#### 3. Add Context/Documentation

```powershell
$context = @{
    task_id = $taskId
    context = @"
# Invoice List Component

## Implementation Details
- Used React functional components with hooks
- Implemented useInvoiceStore custom hook for state management
- Added sorting by date, amount, and status
- Implemented filtering by status (paid, pending, overdue)
- Added responsive design for mobile devices

## Files Created
- src/components/InvoiceList.tsx
- src/components/InvoiceListItem.tsx
- src/hooks/useInvoiceStore.ts
- src/styles/InvoiceList.css

## Dependencies Added
- date-fns for date formatting
- react-table for advanced table features

## Testing
- Unit tests added in __tests__/InvoiceList.test.tsx
- All tests passing (15/15)

## Next Steps
- Add pagination for large invoice lists
- Implement export to CSV functionality
"@
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8080/api/documentation" `
    -Method POST `
    -ContentType "application/json" `
    -Body $context
```

#### 4. Check All Tasks for the Frontend Agent

```powershell
# Get all tasks assigned to InvoiceAI-Frontend
$agentId = "8f3e7b1c-4d2a-4c8f-9e7a-b5c3d1e9f8a2"

$tasks = Invoke-RestMethod -Uri "http://localhost:8080/api/agents/$agentId/tasks" `
    -Method GET

$tasks | Format-Table -Property title, status, priority
```

---

## Complete Workflow Example

### Multi-Agent Collaboration: InvoiceAI Project

Let's set up a complete workflow with multiple agents working together:

```powershell
# 1. Create Project
$project = @{
    name = "InvoiceAI"
    description = "AI-powered invoice processing system"
} | ConvertTo-Json

$proj = Invoke-RestMethod -Uri "http://localhost:8080/api/projects" `
    -Method POST -ContentType "application/json" -Body $project

$projectId = $proj.id

# 2. Register Frontend Agent
$frontendAgent = @{
    project_id = $projectId
    name = "InvoiceAI-Frontend"
    role = "frontend"
    team = "UI Team"
} | ConvertTo-Json

$frontend = Invoke-RestMethod -Uri "http://localhost:8080/api/agents" `
    -Method POST -ContentType "application/json" -Body $frontendAgent

# 3. Register Backend Agent
$backendAgent = @{
    project_id = $projectId
    name = "InvoiceAI-Backend"
    role = "backend"
    team = "API Team"
} | ConvertTo-Json

$backend = Invoke-RestMethod -Uri "http://localhost:8080/api/agents" `
    -Method POST -ContentType "application/json" -Body $backendAgent

# 4. Create Backend Task (API Development)
$apiTask = @{
    project_id = $projectId
    agent_id = $backend.id
    title = "Create Invoice REST API"
    description = "Implement CRUD endpoints for invoice management"
    priority = "high"
    dependencies = @()
} | ConvertTo-Json

$apiTaskResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/tasks" `
    -Method POST -ContentType "application/json" -Body $apiTask

# 5. Create Frontend Task (depends on backend)
$uiTask = @{
    project_id = $projectId
    agent_id = $frontend.id
    title = "Build Invoice Dashboard"
    description = "Create React dashboard consuming Invoice API"
    priority = "high"
    dependencies = @($apiTaskResponse.id)
} | ConvertTo-Json

$uiTaskResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/tasks" `
    -Method POST -ContentType "application/json" -Body $uiTask

Write-Host "✅ Project Setup Complete!" -ForegroundColor Green
Write-Host "Project ID: $projectId"
Write-Host "Frontend Agent: $($frontend.name) (ID: $($frontend.id))"
Write-Host "Backend Agent: $($backend.name) (ID: $($backend.id))"
Write-Host "Backend Task: $($apiTaskResponse.title)"
Write-Host "Frontend Task: $($uiTaskResponse.title)"
```

---

## Integration with VS Code AI Agents

### Method 1: Using GitHub Copilot Chat Commands

In your VS Code, you can instruct Copilot to interact with the MCP server:

```
@workspace I'm working as InvoiceAI-Frontend agent (ID: 8f3e7b1c-4d2a-4c8f-9e7a-b5c3d1e9f8a2). 
Can you help me:
1. Check my current tasks from http://localhost:8080/api/agents/8f3e7b1c-4d2a-4c8f-9e7a-b5c3d1e9f8a2/tasks
2. Start working on the highest priority task
3. Update the task status to "in_progress" when I start
```

### Method 2: Create Helper Scripts

Create a PowerShell script to manage agent identity:

**agent-helper.ps1:**

```powershell
# Save this as agent-helper.ps1 in your project root

param(
    [Parameter(Mandatory=$true)]
    [string]$AgentName,
    
    [Parameter(Mandatory=$true)]
    [string]$AgentId,
    
    [Parameter(Mandatory=$false)]
    [string]$Action = "status"
)

$baseUrl = "http://localhost:8080/api"

function Get-AgentTasks {
    Write-Host "Fetching tasks for $AgentName..." -ForegroundColor Cyan
    $tasks = Invoke-RestMethod -Uri "$baseUrl/agents/$AgentId/tasks"
    $tasks | Format-Table -Property title, status, priority, created_at
}

function Update-TaskStatus {
    param([string]$TaskId, [string]$Status)
    $body = @{ status = $Status } | ConvertTo-Json
    Invoke-RestMethod -Uri "$baseUrl/tasks/$TaskId/status" `
        -Method PUT -ContentType "application/json" -Body $body
    Write-Host "✅ Task status updated to: $Status" -ForegroundColor Green
}

function Show-AgentInfo {
    Write-Host "=== Agent Information ===" -ForegroundColor Yellow
    Write-Host "Name: $AgentName"
    Write-Host "ID: $AgentId"
    Write-Host "API Base: $baseUrl"
    Write-Host "=========================" -ForegroundColor Yellow
}

switch ($Action) {
    "status" { 
        Show-AgentInfo
        Get-AgentTasks 
    }
    "tasks" { Get-AgentTasks }
    "info" { Show-AgentInfo }
    default { Write-Host "Unknown action: $Action" }
}
```

**Usage:**

```powershell
# Check status and tasks
.\agent-helper.ps1 -AgentName "InvoiceAI-Frontend" `
    -AgentId "8f3e7b1c-4d2a-4c8f-9e7a-b5c3d1e9f8a2" `
    -Action status
```

### Method 3: Environment Variables

Set up your VS Code workspace with agent identity:

**.vscode/settings.json:**

```json
{
  "terminal.integrated.env.windows": {
    "MCP_AGENT_NAME": "InvoiceAI-Frontend",
    "MCP_AGENT_ID": "8f3e7b1c-4d2a-4c8f-9e7a-b5c3d1e9f8a2",
    "MCP_PROJECT_ID": "2da7e7e6-0f8c-4f33-9ae8-f3dfb6f0ec8a",
    "MCP_API_URL": "http://localhost:8080/api"
  }
}
```

Then create a script that uses these variables:

```powershell
# quick-status.ps1
$agentId = $env:MCP_AGENT_ID
$apiUrl = $env:MCP_API_URL

Invoke-RestMethod -Uri "$apiUrl/agents/$agentId/tasks"
```

---

## API Reference

### Projects

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/projects` | POST | Create a new project |
| `/api/projects` | GET | List all projects |
| `/api/projects/{id}` | GET | Get project details |

### Agents

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/agents` | POST | Register a new agent |
| `/api/agents/{id}` | GET | Get agent details |
| `/api/projects/{project_id}/agents` | GET | List agents in a project |
| `/api/agents/{id}/status` | PUT | Update agent status |

### Tasks

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/tasks` | POST | Create a new task |
| `/api/tasks/{id}` | GET | Get task details |
| `/api/tasks/{id}/status` | PUT | Update task status |
| `/api/agents/{agent_id}/tasks` | GET | Get tasks for an agent |
| `/api/projects/{project_id}/tasks` | GET | Get tasks for a project |

### Documentation/Context

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/documentation` | POST | Add context/documentation |
| `/api/tasks/{task_id}/documentation` | GET | Get task documentation |

### Agent Roles

Available roles:
- `frontend` - UI/UX development
- `backend` - API/Server development

### Task Status

Available statuses:
- `pending` - Not started
- `in_progress` - Currently being worked on
- `done` - Completed
- `blocked` - Waiting on dependencies

### Task Priority

Available priorities:
- `low`
- `medium`
- `high`

---

## Troubleshooting

### Issue: Cannot Connect to MCP Server

**Solution:**
```powershell
# Check if services are running
docker-compose ps

# If not running, start them
docker-compose up -d

# Check logs for errors
docker-compose logs mcp-server
```

### Issue: Agent Not Found

**Solution:**
```powershell
# List all agents in your project
$projectId = "your-project-id"
Invoke-RestMethod -Uri "http://localhost:8080/api/projects/$projectId/agents"
```

### Issue: Task Creation Fails

**Solution:**
- Verify project_id exists
- Verify agent_id exists and belongs to the project
- Check that status is one of: pending, in_progress, done, blocked
- Check that priority is one of: low, medium, high

### View Server Logs

```powershell
# Real-time logs
docker-compose logs -f mcp-server

# Last 50 lines
docker-compose logs --tail=50 mcp-server
```

### Reset Database

```powershell
# Stop services
docker-compose down

# Remove volumes (WARNING: Deletes all data)
docker-compose down -v

# Start fresh
docker-compose up -d
```

---

## Best Practices

### 1. Agent Identity Management

- Use descriptive names: `ProjectName-Role` (e.g., "InvoiceAI-Frontend")
- Store agent IDs in your project's `.env` file (add to `.gitignore`)
- Document which agent is responsible for which components

### 2. Task Organization

- Create tasks with clear, actionable titles
- Use dependencies to link related tasks
- Update status regularly to track progress
- Add detailed context/documentation when completing tasks

### 3. Context Sharing

- Document implementation decisions
- List files created/modified
- Note dependencies added
- Share code patterns and conventions

### 4. Multi-Agent Coordination

- Frontend depends on Backend API tasks
- Use task dependencies to enforce workflow order
- Regular status checks to avoid conflicts
- Document interfaces between components

---

## Example Project Structure

```
InvoiceAI/
├── frontend/                 (InvoiceAI-Frontend)
│   ├── src/
│   │   ├── components/
│   │   ├── hooks/
│   │   └── pages/
│   └── agent.json           # Agent identity
│
├── backend/                  (InvoiceAI-Backend)
│   ├── api/
│   ├── services/
│   └── agent.json           # Agent identity
│
└── .vscode/
    └── settings.json        # MCP configuration
```

**agent.json example:**

```json
{
  "mcp": {
    "agent_name": "InvoiceAI-Frontend",
    "agent_id": "8f3e7b1c-4d2a-4c8f-9e7a-b5c3d1e9f8a2",
    "project_id": "2da7e7e6-0f8c-4f33-9ae8-f3dfb6f0ec8a",
    "role": "frontend",
    "team": "UI Development"
  }
}
```

---

## Additional Resources

- **API Documentation**: See `docs/API.md` for detailed API specifications
- **WebSocket Demo**: Open `http://localhost:8080` in browser for real-time updates
- **Health Check**: `http://localhost:8080/health`

---

## Quick Reference Card

```powershell
# Start services
docker-compose up -d

# Create project
$proj = Invoke-RestMethod -Uri "http://localhost:8080/api/projects" -Method POST -ContentType "application/json" -Body '{"name":"MyProject","description":"..."}'

# Register agent
$agent = Invoke-RestMethod -Uri "http://localhost:8080/api/agents" -Method POST -ContentType "application/json" -Body "{`"project_id`":`"$($proj.id)`",`"name`":`"MyAgent-Frontend`",`"role`":`"frontend`"}"

# Create task
$task = Invoke-RestMethod -Uri "http://localhost:8080/api/tasks" -Method POST -ContentType "application/json" -Body "{`"project_id`":`"$($proj.id)`",`"agent_id`":`"$($agent.id)`",`"title`":`"Task Title`",`"priority`":`"high`"}"

# Get agent tasks
Invoke-RestMethod -Uri "http://localhost:8080/api/agents/$($agent.id)/tasks"

# Update task status
Invoke-RestMethod -Uri "http://localhost:8080/api/tasks/$($task.id)/status" -Method PUT -ContentType "application/json" -Body '{"status":"in_progress"}'
```

---

**Happy Multi-Agent Development! 🚀**
