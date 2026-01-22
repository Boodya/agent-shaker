package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/techbuzzz/agent-shaker/internal/database"
)

// DashboardHandler handles dashboard statistics requests
type DashboardHandler struct {
	db *database.DB
}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler(db *database.DB) *DashboardHandler {
	return &DashboardHandler{db: db}
}

// DashboardStats represents the dashboard statistics
type DashboardStats struct {
	Projects ProjectStats `json:"projects"`
	Agents   AgentStats   `json:"agents"`
	Tasks    TaskStats    `json:"tasks"`
	Contexts ContextStats `json:"contexts"`
}

// ProjectStats represents project statistics
type ProjectStats struct {
	Total    int `json:"total"`
	Active   int `json:"active"`
	Archived int `json:"archived"`
}

// AgentStats represents agent statistics
type AgentStats struct {
	Total   int `json:"total"`
	Active  int `json:"active"`
	Idle    int `json:"idle"`
	Offline int `json:"offline"`
}

// TaskStats represents task statistics
type TaskStats struct {
	Total      int `json:"total"`
	Pending    int `json:"pending"`
	InProgress int `json:"in_progress"`
	Done       int `json:"done"`
	Blocked    int `json:"blocked"`
}

// ContextStats represents context statistics
type ContextStats struct {
	Total int `json:"total"`
}

// GetDashboardStats returns comprehensive dashboard statistics
func (h *DashboardHandler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		http.Error(w, "Database connection not available", http.StatusServiceUnavailable)
		return
	}

	stats := DashboardStats{}

	// Get project statistics
	var projectStats ProjectStats
	err := h.db.QueryRow(`
		SELECT 
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE status = 'active') as active,
			COUNT(*) FILTER (WHERE status = 'archived') as archived
		FROM projects
	`).Scan(&projectStats.Total, &projectStats.Active, &projectStats.Archived)

	if err != nil {
		log.Printf("Error fetching project stats: %v", err)
		projectStats = ProjectStats{Total: 0, Active: 0, Archived: 0}
	}
	stats.Projects = projectStats

	// Get agent statistics
	var agentStats AgentStats
	err = h.db.QueryRow(`
		SELECT 
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE status = 'active') as active,
			COUNT(*) FILTER (WHERE status = 'idle') as idle,
			COUNT(*) FILTER (WHERE status = 'offline') as offline
		FROM agents
	`).Scan(&agentStats.Total, &agentStats.Active, &agentStats.Idle, &agentStats.Offline)

	if err != nil {
		log.Printf("Error fetching agent stats: %v", err)
		agentStats = AgentStats{Total: 0, Active: 0, Idle: 0, Offline: 0}
	}
	stats.Agents = agentStats

	// Get task statistics
	var taskStats TaskStats
	err = h.db.QueryRow(`
		SELECT 
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE status = 'pending') as pending,
			COUNT(*) FILTER (WHERE status = 'in_progress') as in_progress,
			COUNT(*) FILTER (WHERE status = 'done') as done,
			COUNT(*) FILTER (WHERE status = 'blocked') as blocked
		FROM tasks
	`).Scan(&taskStats.Total, &taskStats.Pending, &taskStats.InProgress, &taskStats.Done, &taskStats.Blocked)

	if err != nil {
		log.Printf("Error fetching task stats: %v", err)
		taskStats = TaskStats{Total: 0, Pending: 0, InProgress: 0, Done: 0, Blocked: 0}
	}
	stats.Tasks = taskStats

	// Get context statistics
	var contextStats ContextStats
	err = h.db.QueryRow(`
		SELECT COUNT(*) as total
		FROM contexts
	`).Scan(&contextStats.Total)

	if err != nil {
		log.Printf("Error fetching context stats: %v", err)
		contextStats = ContextStats{Total: 0}
	}
	stats.Contexts = contextStats

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
