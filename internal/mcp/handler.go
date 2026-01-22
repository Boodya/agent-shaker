package mcp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/techbuzzz/agent-shaker/internal/database"
)

// JSON-RPC 2.0 structures
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type JSONRPCResponse struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      interface{}   `json:"id,omitempty"`
	Result  interface{}   `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
}

type JSONRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// MCP Protocol structures
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ServerCapabilities struct {
	Tools     *ToolsCapability     `json:"tools,omitempty"`
	Resources *ResourcesCapability `json:"resources,omitempty"`
	Prompts   *PromptsCapability   `json:"prompts,omitempty"`
}

type ToolsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

type ResourcesCapability struct {
	Subscribe   bool `json:"subscribe,omitempty"`
	ListChanged bool `json:"listChanged,omitempty"`
}

type PromptsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

type InitializeResult struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ServerCapabilities `json:"capabilities"`
	ServerInfo      ServerInfo         `json:"serverInfo"`
}

type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	InputSchema InputSchema `json:"inputSchema"`
}

type InputSchema struct {
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	Required   []string               `json:"required,omitempty"`
}

type ToolsListResult struct {
	Tools []Tool `json:"tools"`
}

type Resource struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	MimeType    string `json:"mimeType,omitempty"`
}

type ResourcesListResult struct {
	Resources []Resource `json:"resources"`
}

type ResourceContent struct {
	URI      string `json:"uri"`
	MimeType string `json:"mimeType,omitempty"`
	Text     string `json:"text,omitempty"`
}

type ResourcesReadResult struct {
	Contents []ResourceContent `json:"contents"`
}

type ToolCallParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

type ToolResult struct {
	Content []ToolResultContent `json:"content"`
	IsError bool                `json:"isError,omitempty"`
}

type ToolResultContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// MCPHandler handles MCP protocol requests
type MCPHandler struct {
	db       *database.DB
	sessions sync.Map
}

type Session struct {
	ID         string
	CreatedAt  time.Time
	ClientInfo map[string]interface{}
}

func NewMCPHandler(db *database.DB) *MCPHandler {
	return &MCPHandler{
		db: db,
	}
}

// HandleMCP handles the main MCP endpoint with SSE support
func (h *MCPHandler) HandleMCP(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Check for SSE request (GET with Accept: text/event-stream)
	if r.Method == "GET" {
		accept := r.Header.Get("Accept")
		if accept == "text/event-stream" {
			h.handleSSE(w, r)
			return
		}
		// Return server info for plain GET
		h.handleServerInfo(w, r)
		return
	}

	// Handle POST requests (JSON-RPC)
	if r.Method == "POST" {
		h.handleJSONRPC(w, r)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (h *MCPHandler) handleServerInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"name":            "agent-shaker",
		"version":         "1.0.0",
		"protocolVersion": "2024-11-05",
		"capabilities": map[string]interface{}{
			"tools":     map[string]bool{"listChanged": false},
			"resources": map[string]bool{"subscribe": false, "listChanged": false},
		},
	})
}

func (h *MCPHandler) handleSSE(w http.ResponseWriter, r *http.Request) {
	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	// Create session
	sessionID := uuid.New().String()
	session := &Session{
		ID:        sessionID,
		CreatedAt: time.Now(),
	}
	h.sessions.Store(sessionID, session)
	defer h.sessions.Delete(sessionID)

	log.Printf("MCP SSE connection established: %s", sessionID)

	// Send initial endpoint message
	endpointMsg := fmt.Sprintf("event: endpoint\ndata: /mcp/message?sessionId=%s\n\n", sessionID)
	w.Write([]byte(endpointMsg))
	flusher.Flush()

	// Keep connection alive with periodic pings
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	done := r.Context().Done()
	for {
		select {
		case <-done:
			log.Printf("MCP SSE connection closed: %s", sessionID)
			return
		case <-ticker.C:
			// Send ping to keep connection alive
			w.Write([]byte(": ping\n\n"))
			flusher.Flush()
		}
	}
}

func (h *MCPHandler) handleJSONRPC(w http.ResponseWriter, r *http.Request) {
	var req JSONRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, nil, -32700, "Parse error", err.Error())
		return
	}

	log.Printf("MCP Request: method=%s, id=%v", req.Method, req.ID)

	var result interface{}
	var rpcErr *JSONRPCError

	switch req.Method {
	case "initialize":
		result, rpcErr = h.handleInitialize(req.Params)
	case "initialized":
		// Client notification that initialization is complete
		result = map[string]interface{}{}
	case "tools/list":
		result, rpcErr = h.handleToolsList()
	case "tools/call":
		result, rpcErr = h.handleToolsCall(req.Params)
	case "resources/list":
		result, rpcErr = h.handleResourcesList()
	case "resources/read":
		result, rpcErr = h.handleResourcesRead(req.Params)
	case "ping":
		result = map[string]interface{}{}
	default:
		rpcErr = &JSONRPCError{
			Code:    -32601,
			Message: "Method not found",
			Data:    fmt.Sprintf("Unknown method: %s", req.Method),
		}
	}

	h.sendResponse(w, req.ID, result, rpcErr)
}

func (h *MCPHandler) handleInitialize(params json.RawMessage) (interface{}, *JSONRPCError) {
	// Parse client info if provided
	var clientParams struct {
		ProtocolVersion string                 `json:"protocolVersion"`
		Capabilities    map[string]interface{} `json:"capabilities"`
		ClientInfo      map[string]interface{} `json:"clientInfo"`
	}
	if params != nil {
		json.Unmarshal(params, &clientParams)
	}

	log.Printf("MCP Initialize - Client: %v, Protocol: %s", clientParams.ClientInfo, clientParams.ProtocolVersion)

	return InitializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities: ServerCapabilities{
			Tools: &ToolsCapability{
				ListChanged: false,
			},
			Resources: &ResourcesCapability{
				Subscribe:   false,
				ListChanged: false,
			},
		},
		ServerInfo: ServerInfo{
			Name:    "agent-shaker",
			Version: "1.0.0",
		},
	}, nil
}

func (h *MCPHandler) handleToolsList() (interface{}, *JSONRPCError) {
	tools := []Tool{
		{
			Name:        "list_projects",
			Description: "List all projects in the system",
			InputSchema: InputSchema{
				Type:       "object",
				Properties: map[string]interface{}{},
			},
		},
		{
			Name:        "get_project",
			Description: "Get details of a specific project",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"project_id": map[string]interface{}{
						"type":        "string",
						"description": "The project ID (UUID)",
					},
				},
				Required: []string{"project_id"},
			},
		},
		{
			Name:        "list_agents",
			Description: "List all agents, optionally filtered by project",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"project_id": map[string]interface{}{
						"type":        "string",
						"description": "Optional project ID to filter agents",
					},
				},
			},
		},
		{
			Name:        "get_agent",
			Description: "Get details of a specific agent",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"agent_id": map[string]interface{}{
						"type":        "string",
						"description": "The agent ID (UUID)",
					},
				},
				Required: []string{"agent_id"},
			},
		},
		{
			Name:        "list_tasks",
			Description: "List tasks, optionally filtered by project or agent",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"project_id": map[string]interface{}{
						"type":        "string",
						"description": "Optional project ID to filter tasks",
					},
					"agent_id": map[string]interface{}{
						"type":        "string",
						"description": "Optional agent ID to filter tasks",
					},
					"status": map[string]interface{}{
						"type":        "string",
						"description": "Optional status filter (pending, in_progress, done, blocked)",
					},
				},
			},
		},
		{
			Name:        "create_task",
			Description: "Create a new task in a project",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"project_id": map[string]interface{}{
						"type":        "string",
						"description": "The project ID",
					},
					"title": map[string]interface{}{
						"type":        "string",
						"description": "Task title",
					},
					"description": map[string]interface{}{
						"type":        "string",
						"description": "Task description",
					},
					"priority": map[string]interface{}{
						"type":        "string",
						"description": "Priority: low, medium, high",
						"enum":        []string{"low", "medium", "high"},
					},
					"assigned_to": map[string]interface{}{
						"type":        "string",
						"description": "Agent ID to assign the task to",
					},
				},
				Required: []string{"project_id", "title"},
			},
		},
		{
			Name:        "update_task_status",
			Description: "Update the status of a task",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"task_id": map[string]interface{}{
						"type":        "string",
						"description": "The task ID",
					},
					"status": map[string]interface{}{
						"type":        "string",
						"description": "New status: pending, in_progress, done, blocked",
						"enum":        []string{"pending", "in_progress", "done", "blocked"},
					},
				},
				Required: []string{"task_id", "status"},
			},
		},
		{
			Name:        "list_contexts",
			Description: "List documentation/contexts for a project",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"project_id": map[string]interface{}{
						"type":        "string",
						"description": "Optional project ID to filter contexts",
					},
				},
			},
		},
		{
			Name:        "add_context",
			Description: "Add documentation or context to a project",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]interface{}{
					"project_id": map[string]interface{}{
						"type":        "string",
						"description": "The project ID",
					},
					"title": map[string]interface{}{
						"type":        "string",
						"description": "Context title",
					},
					"content": map[string]interface{}{
						"type":        "string",
						"description": "Context content (markdown supported)",
					},
					"tags": map[string]interface{}{
						"type":        "array",
						"description": "Tags for categorization",
						"items":       map[string]string{"type": "string"},
					},
				},
				Required: []string{"project_id", "title", "content"},
			},
		},
		{
			Name:        "get_dashboard",
			Description: "Get dashboard statistics and overview",
			InputSchema: InputSchema{
				Type:       "object",
				Properties: map[string]interface{}{},
			},
		},
	}

	return ToolsListResult{Tools: tools}, nil
}

func (h *MCPHandler) handleToolsCall(params json.RawMessage) (interface{}, *JSONRPCError) {
	var callParams ToolCallParams
	if err := json.Unmarshal(params, &callParams); err != nil {
		return nil, &JSONRPCError{
			Code:    -32602,
			Message: "Invalid params",
			Data:    err.Error(),
		}
	}

	log.Printf("MCP Tool Call: %s with args %v", callParams.Name, callParams.Arguments)

	var resultText string
	var isError bool

	switch callParams.Name {
	case "list_projects":
		resultText, isError = h.executeListProjects()
	case "get_project":
		resultText, isError = h.executeGetProject(callParams.Arguments)
	case "list_agents":
		resultText, isError = h.executeListAgents(callParams.Arguments)
	case "get_agent":
		resultText, isError = h.executeGetAgent(callParams.Arguments)
	case "list_tasks":
		resultText, isError = h.executeListTasks(callParams.Arguments)
	case "create_task":
		resultText, isError = h.executeCreateTask(callParams.Arguments)
	case "update_task_status":
		resultText, isError = h.executeUpdateTaskStatus(callParams.Arguments)
	case "list_contexts":
		resultText, isError = h.executeListContexts(callParams.Arguments)
	case "add_context":
		resultText, isError = h.executeAddContext(callParams.Arguments)
	case "get_dashboard":
		resultText, isError = h.executeGetDashboard()
	default:
		return nil, &JSONRPCError{
			Code:    -32601,
			Message: "Unknown tool",
			Data:    fmt.Sprintf("Tool not found: %s", callParams.Name),
		}
	}

	return ToolResult{
		Content: []ToolResultContent{
			{Type: "text", Text: resultText},
		},
		IsError: isError,
	}, nil
}

func (h *MCPHandler) handleResourcesList() (interface{}, *JSONRPCError) {
	resources := []Resource{
		{
			URI:         "agent-shaker://projects",
			Name:        "Projects",
			Description: "List of all projects",
			MimeType:    "application/json",
		},
		{
			URI:         "agent-shaker://agents",
			Name:        "Agents",
			Description: "List of all agents",
			MimeType:    "application/json",
		},
		{
			URI:         "agent-shaker://tasks",
			Name:        "Tasks",
			Description: "List of all tasks",
			MimeType:    "application/json",
		},
		{
			URI:         "agent-shaker://dashboard",
			Name:        "Dashboard",
			Description: "Dashboard statistics",
			MimeType:    "application/json",
		},
	}

	return ResourcesListResult{Resources: resources}, nil
}

func (h *MCPHandler) handleResourcesRead(params json.RawMessage) (interface{}, *JSONRPCError) {
	var readParams struct {
		URI string `json:"uri"`
	}
	if err := json.Unmarshal(params, &readParams); err != nil {
		return nil, &JSONRPCError{
			Code:    -32602,
			Message: "Invalid params",
			Data:    err.Error(),
		}
	}

	var content string
	var isError bool

	switch readParams.URI {
	case "agent-shaker://projects":
		content, isError = h.executeListProjects()
	case "agent-shaker://agents":
		content, isError = h.executeListAgents(nil)
	case "agent-shaker://tasks":
		content, isError = h.executeListTasks(nil)
	case "agent-shaker://dashboard":
		content, isError = h.executeGetDashboard()
	default:
		return nil, &JSONRPCError{
			Code:    -32602,
			Message: "Unknown resource",
			Data:    fmt.Sprintf("Resource not found: %s", readParams.URI),
		}
	}

	if isError {
		return nil, &JSONRPCError{
			Code:    -32000,
			Message: "Resource read failed",
			Data:    content,
		}
	}

	return ResourcesReadResult{
		Contents: []ResourceContent{
			{
				URI:      readParams.URI,
				MimeType: "application/json",
				Text:     content,
			},
		},
	}, nil
}

// Tool execution methods
func (h *MCPHandler) executeListProjects() (string, bool) {
	if h.db == nil {
		return `{"error": "Database not connected"}`, true
	}

	rows, err := h.db.Query(`
		SELECT id, name, description, status, created_at, updated_at 
		FROM projects ORDER BY created_at DESC
	`)
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error()), true
	}
	defer rows.Close()

	var projects []map[string]interface{}
	for rows.Next() {
		var id, name, description, status string
		var createdAt, updatedAt interface{}
		if err := rows.Scan(&id, &name, &description, &status, &createdAt, &updatedAt); err != nil {
			continue
		}
		projects = append(projects, map[string]interface{}{
			"id":          id,
			"name":        name,
			"description": description,
			"status":      status,
			"created_at":  createdAt,
			"updated_at":  updatedAt,
		})
	}

	result, _ := json.MarshalIndent(projects, "", "  ")
	return string(result), false
}

func (h *MCPHandler) executeGetProject(args map[string]interface{}) (string, bool) {
	if h.db == nil {
		return `{"error": "Database not connected"}`, true
	}

	projectID, ok := args["project_id"].(string)
	if !ok {
		return `{"error": "project_id is required"}`, true
	}

	var id, name, description, status string
	var createdAt, updatedAt interface{}
	err := h.db.QueryRow(`
		SELECT id, name, description, status, created_at, updated_at 
		FROM projects WHERE id = $1
	`, projectID).Scan(&id, &name, &description, &status, &createdAt, &updatedAt)
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error()), true
	}

	result, _ := json.MarshalIndent(map[string]interface{}{
		"id":          id,
		"name":        name,
		"description": description,
		"status":      status,
		"created_at":  createdAt,
		"updated_at":  updatedAt,
	}, "", "  ")
	return string(result), false
}

func (h *MCPHandler) executeListAgents(args map[string]interface{}) (string, bool) {
	if h.db == nil {
		return `{"error": "Database not connected"}`, true
	}

	query := `SELECT id, project_id, name, role, status, team, created_at FROM agents`
	var queryArgs []interface{}

	if args != nil {
		if projectID, ok := args["project_id"].(string); ok && projectID != "" {
			query += " WHERE project_id = $1"
			queryArgs = append(queryArgs, projectID)
		}
	}
	query += " ORDER BY created_at DESC"

	rows, err := h.db.Query(query, queryArgs...)
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error()), true
	}
	defer rows.Close()

	var agents []map[string]interface{}
	for rows.Next() {
		var id, projectID, name, role, status string
		var team *string
		var createdAt interface{}
		if err := rows.Scan(&id, &projectID, &name, &role, &status, &team, &createdAt); err != nil {
			continue
		}
		agent := map[string]interface{}{
			"id":         id,
			"project_id": projectID,
			"name":       name,
			"role":       role,
			"status":     status,
			"created_at": createdAt,
		}
		if team != nil {
			agent["team"] = *team
		}
		agents = append(agents, agent)
	}

	result, _ := json.MarshalIndent(agents, "", "  ")
	return string(result), false
}

func (h *MCPHandler) executeGetAgent(args map[string]interface{}) (string, bool) {
	if h.db == nil {
		return `{"error": "Database not connected"}`, true
	}

	agentID, ok := args["agent_id"].(string)
	if !ok {
		return `{"error": "agent_id is required"}`, true
	}

	var id, projectID, name, role, status string
	var team *string
	var createdAt interface{}
	err := h.db.QueryRow(`
		SELECT id, project_id, name, role, status, team, created_at 
		FROM agents WHERE id = $1
	`, agentID).Scan(&id, &projectID, &name, &role, &status, &team, &createdAt)
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error()), true
	}

	agent := map[string]interface{}{
		"id":         id,
		"project_id": projectID,
		"name":       name,
		"role":       role,
		"status":     status,
		"created_at": createdAt,
	}
	if team != nil {
		agent["team"] = *team
	}

	result, _ := json.MarshalIndent(agent, "", "  ")
	return string(result), false
}

func (h *MCPHandler) executeListTasks(args map[string]interface{}) (string, bool) {
	if h.db == nil {
		return `{"error": "Database not connected"}`, true
	}

	query := `SELECT id, project_id, title, description, status, priority, assigned_to, created_at FROM tasks WHERE 1=1`
	var queryArgs []interface{}
	argNum := 1

	if args != nil {
		if projectID, ok := args["project_id"].(string); ok && projectID != "" {
			query += fmt.Sprintf(" AND project_id = $%d", argNum)
			queryArgs = append(queryArgs, projectID)
			argNum++
		}
		if agentID, ok := args["agent_id"].(string); ok && agentID != "" {
			query += fmt.Sprintf(" AND assigned_to = $%d", argNum)
			queryArgs = append(queryArgs, agentID)
			argNum++
		}
		if status, ok := args["status"].(string); ok && status != "" {
			query += fmt.Sprintf(" AND status = $%d", argNum)
			queryArgs = append(queryArgs, status)
			argNum++
		}
	}
	query += " ORDER BY created_at DESC"

	rows, err := h.db.Query(query, queryArgs...)
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error()), true
	}
	defer rows.Close()

	var tasks []map[string]interface{}
	for rows.Next() {
		var id, projectID, title, status, priority string
		var description, assignedTo *string
		var createdAt interface{}
		if err := rows.Scan(&id, &projectID, &title, &description, &status, &priority, &assignedTo, &createdAt); err != nil {
			continue
		}
		task := map[string]interface{}{
			"id":         id,
			"project_id": projectID,
			"title":      title,
			"status":     status,
			"priority":   priority,
			"created_at": createdAt,
		}
		if description != nil {
			task["description"] = *description
		}
		if assignedTo != nil {
			task["assigned_to"] = *assignedTo
		}
		tasks = append(tasks, task)
	}

	result, _ := json.MarshalIndent(tasks, "", "  ")
	return string(result), false
}

func (h *MCPHandler) executeCreateTask(args map[string]interface{}) (string, bool) {
	if h.db == nil {
		return `{"error": "Database not connected"}`, true
	}

	projectID, ok := args["project_id"].(string)
	if !ok {
		return `{"error": "project_id is required"}`, true
	}
	title, ok := args["title"].(string)
	if !ok {
		return `{"error": "title is required"}`, true
	}

	description, _ := args["description"].(string)
	priority, _ := args["priority"].(string)
	if priority == "" {
		priority = "medium"
	}
	assignedTo, _ := args["assigned_to"].(string)

	id := uuid.New().String()
	query := `INSERT INTO tasks (id, project_id, title, description, status, priority, assigned_to) 
	          VALUES ($1, $2, $3, $4, 'pending', $5, $6) RETURNING id, created_at`

	var createdID string
	var createdAt interface{}
	var assignedToPtr *string
	if assignedTo != "" {
		assignedToPtr = &assignedTo
	}

	err := h.db.QueryRow(query, id, projectID, title, description, priority, assignedToPtr).Scan(&createdID, &createdAt)
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error()), true
	}

	result, _ := json.MarshalIndent(map[string]interface{}{
		"success":    true,
		"id":         createdID,
		"title":      title,
		"status":     "pending",
		"priority":   priority,
		"created_at": createdAt,
	}, "", "  ")
	return string(result), false
}

func (h *MCPHandler) executeUpdateTaskStatus(args map[string]interface{}) (string, bool) {
	if h.db == nil {
		return `{"error": "Database not connected"}`, true
	}

	taskID, ok := args["task_id"].(string)
	if !ok {
		return `{"error": "task_id is required"}`, true
	}
	status, ok := args["status"].(string)
	if !ok {
		return `{"error": "status is required"}`, true
	}

	// Validate status
	validStatuses := map[string]bool{"pending": true, "in_progress": true, "done": true, "blocked": true}
	if !validStatuses[status] {
		return `{"error": "invalid status, must be one of: pending, in_progress, done, blocked"}`, true
	}

	_, err := h.db.Exec(`UPDATE tasks SET status = $1, updated_at = NOW() WHERE id = $2`, status, taskID)
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error()), true
	}

	result, _ := json.MarshalIndent(map[string]interface{}{
		"success": true,
		"task_id": taskID,
		"status":  status,
	}, "", "  ")
	return string(result), false
}

func (h *MCPHandler) executeListContexts(args map[string]interface{}) (string, bool) {
	if h.db == nil {
		return `{"error": "Database not connected"}`, true
	}

	query := `SELECT id, project_id, title, content, tags, created_at FROM contexts`
	var queryArgs []interface{}

	if args != nil {
		if projectID, ok := args["project_id"].(string); ok && projectID != "" {
			query += " WHERE project_id = $1"
			queryArgs = append(queryArgs, projectID)
		}
	}
	query += " ORDER BY created_at DESC"

	rows, err := h.db.Query(query, queryArgs...)
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error()), true
	}
	defer rows.Close()

	var contexts []map[string]interface{}
	for rows.Next() {
		var id, projectID, title, content string
		var tags interface{}
		var createdAt interface{}
		if err := rows.Scan(&id, &projectID, &title, &content, &tags, &createdAt); err != nil {
			continue
		}
		contexts = append(contexts, map[string]interface{}{
			"id":         id,
			"project_id": projectID,
			"title":      title,
			"content":    content,
			"tags":       tags,
			"created_at": createdAt,
		})
	}

	result, _ := json.MarshalIndent(contexts, "", "  ")
	return string(result), false
}

func (h *MCPHandler) executeAddContext(args map[string]interface{}) (string, bool) {
	if h.db == nil {
		return `{"error": "Database not connected"}`, true
	}

	projectID, ok := args["project_id"].(string)
	if !ok {
		return `{"error": "project_id is required"}`, true
	}
	title, ok := args["title"].(string)
	if !ok {
		return `{"error": "title is required"}`, true
	}
	content, ok := args["content"].(string)
	if !ok {
		return `{"error": "content is required"}`, true
	}

	var tagsJSON []byte
	if tags, ok := args["tags"].([]interface{}); ok {
		tagsJSON, _ = json.Marshal(tags)
	} else {
		tagsJSON = []byte("[]")
	}

	id := uuid.New().String()
	query := `INSERT INTO contexts (id, project_id, title, content, tags) VALUES ($1, $2, $3, $4, $5) RETURNING created_at`

	var createdAt interface{}
	err := h.db.QueryRow(query, id, projectID, title, content, tagsJSON).Scan(&createdAt)
	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error()), true
	}

	result, _ := json.MarshalIndent(map[string]interface{}{
		"success":    true,
		"id":         id,
		"title":      title,
		"created_at": createdAt,
	}, "", "  ")
	return string(result), false
}

func (h *MCPHandler) executeGetDashboard() (string, bool) {
	if h.db == nil {
		return `{"error": "Database not connected"}`, true
	}

	var projectCount, agentCount, taskCount, contextCount int
	var pendingTasks, inProgressTasks, doneTasks, blockedTasks int

	h.db.QueryRow("SELECT COUNT(*) FROM projects").Scan(&projectCount)
	h.db.QueryRow("SELECT COUNT(*) FROM agents").Scan(&agentCount)
	h.db.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&taskCount)
	h.db.QueryRow("SELECT COUNT(*) FROM contexts").Scan(&contextCount)

	h.db.QueryRow("SELECT COUNT(*) FROM tasks WHERE status = 'pending'").Scan(&pendingTasks)
	h.db.QueryRow("SELECT COUNT(*) FROM tasks WHERE status = 'in_progress'").Scan(&inProgressTasks)
	h.db.QueryRow("SELECT COUNT(*) FROM tasks WHERE status = 'done'").Scan(&doneTasks)
	h.db.QueryRow("SELECT COUNT(*) FROM tasks WHERE status = 'blocked'").Scan(&blockedTasks)

	result, _ := json.MarshalIndent(map[string]interface{}{
		"projects":          projectCount,
		"agents":            agentCount,
		"tasks":             taskCount,
		"contexts":          contextCount,
		"pending_tasks":     pendingTasks,
		"in_progress_tasks": inProgressTasks,
		"done_tasks":        doneTasks,
		"blocked_tasks":     blockedTasks,
	}, "", "  ")
	return string(result), false
}

func (h *MCPHandler) sendResponse(w http.ResponseWriter, id interface{}, result interface{}, rpcErr *JSONRPCError) {
	w.Header().Set("Content-Type", "application/json")

	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
	}

	if rpcErr != nil {
		resp.Error = rpcErr
	} else {
		resp.Result = result
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *MCPHandler) sendError(w http.ResponseWriter, id interface{}, code int, message string, data interface{}) {
	h.sendResponse(w, id, nil, &JSONRPCError{
		Code:    code,
		Message: message,
		Data:    data,
	})
}
