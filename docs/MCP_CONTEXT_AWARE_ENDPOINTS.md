# MCP Context-Aware Endpoints

## Overview

The MCP server now automatically uses `project_id` and `agent_id` from the connection URL for `create_task` and `add_context` endpoints. This makes the API much more user-friendly when connected with context.

## Connection URL Format

When connecting to the MCP server, include your project and agent IDs in the URL:

```
http://localhost:8080?project_id=68488bf3-8d73-498f-b871-69d63641d6e3&agent_id=1a9c32e7-f0b0-4cc4-b1ed-6ee92f7ef184
```

## Enhanced Endpoints

### 1. `create_task` - Context-Aware Task Creation

**Before (required all parameters):**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "create_task",
    "arguments": {
      "project_id": "68488bf3-8d73-498f-b871-69d63641d6e3",
      "created_by": "1a9c32e7-f0b0-4cc4-b1ed-6ee92f7ef184",
      "title": "Implement user authentication",
      "description": "Add JWT-based authentication",
      "priority": "high"
    }
  }
}
```

**After (with context from URL - simplified!):**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "create_task",
    "arguments": {
      "title": "Implement user authentication",
      "description": "Add JWT-based authentication",
      "priority": "high"
    }
  }
}
```

**How it works:**
1. If `project_id` is not in arguments, uses `project_id` from URL context
2. If `created_by` is not in arguments, uses `agent_id` from URL context
3. If still no `created_by`, falls back to first agent in project
4. Only `title` is required when connected with context!

### 2. `add_context` - Context-Aware Documentation

**Before (required all parameters):**
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "add_context",
    "arguments": {
      "project_id": "68488bf3-8d73-498f-b871-69d63641d6e3",
      "agent_id": "1a9c32e7-f0b0-4cc4-b1ed-6ee92f7ef184",
      "title": "Authentication Implementation Notes",
      "content": "## JWT Implementation\n\nUsing RS256 algorithm...",
      "tags": ["auth", "security", "backend"]
    }
  }
}
```

**After (with context from URL - simplified!):**
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "add_context",
    "arguments": {
      "title": "Authentication Implementation Notes",
      "content": "## JWT Implementation\n\nUsing RS256 algorithm...",
      "tags": ["auth", "security", "backend"]
    }
  }
}
```

**How it works:**
1. If `project_id` is not in arguments, uses `project_id` from URL context
2. If `agent_id` is not in arguments, uses `agent_id` from URL context
3. If still no `agent_id`, falls back to first agent in project
4. Only `title` and `content` are required when connected with context!

## Priority of Values

For both endpoints, the priority order is:

1. **Explicit argument** - If you provide `project_id`, `agent_id`, or `created_by` in the arguments, it will be used
2. **URL context** - If not in arguments, uses values from the MCP connection URL
3. **Fallback** - If still not available, queries the first agent from the project

This allows maximum flexibility while providing sensible defaults.

## Benefits

### 1. **Simplified API Calls**
When connected with context, you only need to provide the essential data (title, content, etc.) without repeating the project and agent IDs.

### 2. **Reduced Errors**
Less manual copying and pasting of UUIDs means fewer mistakes.

### 3. **Better Developer Experience**
AI agents can focus on the actual content rather than managing IDs.

### 4. **Backward Compatible**
You can still explicitly provide `project_id`, `agent_id`, or `created_by` if needed - they will override the context.

## Context-Aware Tools

The following tools automatically use the URL context:

| Tool | Uses project_id from URL | Uses agent_id from URL |
|------|---------------------------|------------------------|
| `get_my_identity` | ✅ | ✅ |
| `get_my_project` | ✅ | ❌ |
| `get_my_tasks` | ❌ | ✅ |
| `update_my_status` | ❌ | ✅ |
| `claim_task` | ❌ | ✅ |
| `complete_task` | ❌ | ✅ |
| `create_task` | ✅ (optional) | ✅ (optional as created_by) |
| `add_context` | ✅ (optional) | ✅ (optional) |

## Example Workflow

### Agent Setup with Context

1. **Connect with context:**
   ```
   http://localhost:8080?project_id=PROJECT_UUID&agent_id=AGENT_UUID
   ```

2. **Check your identity:**
   ```json
   {"method": "tools/call", "params": {"name": "get_my_identity"}}
   ```
   Returns your project and agent details automatically.

3. **Create a task (simplified):**
   ```json
   {
     "method": "tools/call",
     "params": {
       "name": "create_task",
       "arguments": {
         "title": "Fix bug in login form",
         "priority": "high"
       }
     }
   }
   ```
   Automatically uses your project_id and agent_id from the connection!

4. **Add documentation (simplified):**
   ```json
   {
     "method": "tools/call",
     "params": {
       "name": "add_context",
       "arguments": {
         "title": "Bug Fix Notes",
         "content": "Fixed validation issue in email field",
         "tags": ["bugfix", "frontend"]
       }
     }
   }
   ```
   Automatically uses your project_id and agent_id from the connection!

## VS Code Integration

When you download the MCP setup files from the web UI, the `.vscode/mcp.json` file includes:

```json
{
  "mcpServers": {
    "agent-shaker": {
      "url": "http://localhost:8080?project_id=YOUR_PROJECT&agent_id=YOUR_AGENT",
      ...
    }
  }
}
```

This means when you use GitHub Copilot with the MCP server, it automatically knows:
- Which project you're working on
- Which agent you are
- And can create tasks and documentation without you specifying those IDs every time!

## Testing

Test the context-aware functionality:

```bash
# Connect with context
curl "http://localhost:8080?project_id=YOUR_PROJECT_ID&agent_id=YOUR_AGENT_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/call",
    "params": {
      "name": "create_task",
      "arguments": {
        "title": "Test task"
      }
    }
  }'
```

Should create a task automatically using the project_id and agent_id from the URL!

## Summary

The MCP server is now truly context-aware. When you connect with `project_id` and `agent_id` in the URL:

✅ **create_task** - Only requires `title` (project_id and created_by auto-filled)
✅ **add_context** - Only requires `title` and `content` (project_id and agent_id auto-filled)
✅ All context-aware tools use your connection context automatically
✅ You can still override with explicit arguments if needed
✅ Fallback to first project agent if no context available

This makes working with the MCP server much more natural and reduces the cognitive load on both humans and AI agents!
