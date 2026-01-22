# MCP Endpoints Bug Fixes

## Issues Fixed

### 1. `create_task` - Missing `created_by` Field

**Problem:**
```
{"error": "pq: null value in column "created_by" of relation "tasks" violates not-null constraint"}
```

**Root Cause:**
The `tasks` table has a NOT NULL constraint on the `created_by` column, but the MCP handler's `executeCreateTask` function wasn't including this field in the INSERT statement.

**Solution:**
- Added `created_by` parameter support to the `create_task` tool
- If `created_by` is not provided, the function now automatically queries the first agent from the project
- Updated the INSERT statement to include the `created_by` field
- Updated the tool schema to document the `created_by` parameter

**Changes Made:**
```go
// Now accepts created_by parameter
createdBy, _ := args["created_by"].(string)

// Auto-assign first agent if not provided
if createdBy == "" {
    err := h.db.QueryRow(`SELECT id FROM agents WHERE project_id = $1 LIMIT 1`, projectID).Scan(&createdBy)
    if err != nil {
        return `{"error": "created_by is required or no agents found in project"}`, true
    }
}

// Include created_by in INSERT
query := `INSERT INTO tasks (id, project_id, title, description, status, priority, created_by, assigned_to) 
          VALUES ($1, $2, $3, $4, 'pending', $5, $6, $7) RETURNING id, created_at`
```

### 2. `add_context` - Malformed Array Literal for Tags

**Problem:**
```
{"error": "pq: malformed array literal: "["documentation","backend","onboarding"]""}
```

**Root Cause:**
The `executeAddContext` function was converting tags to JSON format using `json.Marshal()`, but PostgreSQL's `TEXT[]` column type expects PostgreSQL array literals, not JSON arrays. The difference:
- JSON format: `["tag1","tag2"]`
- PostgreSQL array format: `{tag1,tag2}`

**Solution:**
- Added missing `agent_id` requirement (contexts table requires this field)
- If `agent_id` is not provided, automatically query the first agent from the project
- Use `pq.Array()` from the `github.com/lib/pq` package for proper PostgreSQL array handling
- Added import for `"github.com/lib/pq"` package

**Changes Made:**
```go
// Import added
import "github.com/lib/pq"

// Handle agent_id requirement
agentID, _ := args["agent_id"].(string)

if agentID == "" {
    err := h.db.QueryRow(`SELECT id FROM agents WHERE project_id = $1 LIMIT 1`, projectID).Scan(&agentID)
    if err != nil {
        return `{"error": "agent_id is required or no agents found in project"}`, true
    }
}

// Convert tags properly for PostgreSQL
var tags []string
if tagsInterface, ok := args["tags"].([]interface{}); ok {
    for _, tag := range tagsInterface {
        if tagStr, ok := tag.(string); ok {
            tags = append(tags, tagStr)
        }
    }
}

// Use pq.Array for proper PostgreSQL array handling
err := h.db.QueryRow(query, id, projectID, agentID, title, content, pq.Array(tags)).Scan(&createdAt)
```

## Testing the Fixes

### Test create_task:

**With created_by (explicit):**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "create_task",
    "arguments": {
      "project_id": "your-project-id",
      "title": "Test Task",
      "description": "Testing create_task fix",
      "priority": "medium",
      "created_by": "agent-id"
    }
  }
}
```

**Without created_by (auto-assigned):**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "create_task",
    "arguments": {
      "project_id": "your-project-id",
      "title": "Test Task",
      "description": "Testing auto-assign",
      "priority": "high"
    }
  }
}
```

### Test add_context:

**With agent_id (explicit):**
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "add_context",
    "arguments": {
      "project_id": "your-project-id",
      "agent_id": "agent-id",
      "title": "Testing Documentation",
      "content": "This is a test context with proper array handling",
      "tags": ["documentation", "backend", "testing"]
    }
  }
}
```

**Without agent_id (auto-assigned):**
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "add_context",
    "arguments": {
      "project_id": "your-project-id",
      "title": "Auto-assigned Context",
      "content": "Testing automatic agent assignment",
      "tags": ["test", "fix"]
    }
  }
}
```

## Files Modified

1. **internal/mcp/handler.go**
   - Added `pq` package import
   - Fixed `executeCreateTask()` to handle `created_by` field
   - Fixed `executeAddContext()` to use proper PostgreSQL array format and handle `agent_id` field
   - Updated tool schemas to document the new optional parameters

2. **cmd/server/main.go** (previous fix)
   - Added CORS headers to health endpoint

## Summary

Both issues have been resolved:
- ✅ `create_task` now properly handles the `created_by` constraint
- ✅ `add_context` now properly handles PostgreSQL array format for tags
- ✅ Both functions automatically assign the first project agent if not explicitly provided
- ✅ Tool schemas updated to reflect the changes

The server has been successfully rebuilt and is ready for testing.
