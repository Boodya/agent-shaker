# DELETE Endpoints Fix

## Issue
The frontend was trying to delete projects, agents, and tasks but getting **404 Page Not Found** errors because the backend DELETE endpoints were not implemented.

## Root Cause
- Frontend API service had `deleteProject()`, `deleteAgent()`, and `deleteTask()` methods
- Backend handlers had no corresponding DELETE methods
- Backend routes had no DELETE endpoint registrations

## Solution

### 1. Added DeleteProject Handler
**File**: `internal/handlers/projects.go`

```go
func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request)
```

**Features**:
- Transaction-based deletion to ensure data consistency
- Cascade delete in correct order:
  1. Contexts (reference tasks)
  2. Tasks (reference project and agents)
  3. Agents (reference project)
  4. Project (parent entity)
- WebSocket broadcast for real-time UI updates
- Returns 204 No Content on success
- Returns 404 if project not found

### 2. Added DeleteAgent Handler
**File**: `internal/handlers/agents.go`

```go
func (h *AgentHandler) DeleteAgent(w http.ResponseWriter, r *http.Request)
```

**Features**:
- Transaction-based deletion
- Cascade delete order:
  1. Contexts (reference tasks which reference agents)
  2. Tasks (reference agent)
  3. Agent (parent entity)
- WebSocket broadcast with project_id
- Returns 204 No Content on success

### 3. Added DeleteTask Handler
**File**: `internal/handlers/tasks.go`

```go
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request)
```

**Features**:
- Transaction-based deletion
- Cascade delete order:
  1. Contexts (reference task)
  2. Task (parent entity)
- WebSocket broadcast with project_id
- Returns 204 No Content on success

### 4. Registered DELETE Routes
**File**: `cmd/server/main.go`

```go
// Projects
api.HandleFunc("/projects/{id}", projectHandler.DeleteProject).Methods("DELETE")

// Agents
api.HandleFunc("/agents/{id}", agentHandler.DeleteAgent).Methods("DELETE")

// Tasks
api.HandleFunc("/tasks/{id}", taskHandler.DeleteTask).Methods("DELETE")
```

## API Endpoints

| Method | Endpoint | Handler | Status Code |
|--------|----------|---------|-------------|
| DELETE | `/api/projects/{id}` | DeleteProject | 204 No Content |
| DELETE | `/api/agents/{id}` | DeleteAgent | 204 No Content |
| DELETE | `/api/tasks/{id}` | DeleteTask | 204 No Content |

## Database Cascade Rules

### Project Deletion
```
Project → Agents → Tasks → Contexts
```

### Agent Deletion
```
Agent → Tasks → Contexts
```

### Task Deletion
```
Task → Contexts
```

## WebSocket Events

All DELETE operations broadcast events to keep frontend in sync:

- `project_deleted` - Sent to project channel
- `agent_deleted` - Sent to project channel with agent_id
- `task_deleted` - Sent to project channel with task_id

## Testing

To test the endpoints:

```bash
# Delete a project
curl -X DELETE http://localhost:8080/api/projects/{project-id}

# Delete an agent
curl -X DELETE http://localhost:8080/api/agents/{agent-id}

# Delete a task
curl -X DELETE http://localhost:8080/api/tasks/{task-id}
```

Expected responses:
- **204 No Content** - Successfully deleted
- **404 Not Found** - Entity not found
- **400 Bad Request** - Invalid UUID format
- **500 Internal Server Error** - Database error

## Notes

1. All DELETE operations use database transactions to ensure atomicity
2. Foreign key constraints are respected through manual cascade deletes
3. WebSocket broadcasts ensure real-time UI updates
4. Frontend stores automatically remove deleted items from their state
5. No orphaned records are left in the database

## Build Verification

```bash
go build -o bin/mcp-server.exe ./cmd/server
```

✅ Build successful - all changes compile without errors.
