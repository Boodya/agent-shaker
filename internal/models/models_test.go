package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestProjectModel(t *testing.T) {
	project := Project{
		ID:          uuid.New(),
		Name:        "Test Project",
		Description: "Test Description",
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if project.Name != "Test Project" {
		t.Errorf("Expected project name to be 'Test Project', got '%s'", project.Name)
	}

	if project.Status != "active" {
		t.Errorf("Expected status to be 'active', got '%s'", project.Status)
	}
}

func TestAgentModel(t *testing.T) {
	agent := Agent{
		ID:        uuid.New(),
		ProjectID: uuid.New(),
		Name:      "Test Agent",
		Role:      "backend",
		Team:      "Backend Team",
		Status:    "active",
		LastSeen:  time.Now(),
		CreatedAt: time.Now(),
	}

	if agent.Name != "Test Agent" {
		t.Errorf("Expected agent name to be 'Test Agent', got '%s'", agent.Name)
	}

	if agent.Role != "backend" {
		t.Errorf("Expected role to be 'backend', got '%s'", agent.Role)
	}
}

func TestTaskModel(t *testing.T) {
	assignedTo := uuid.New()
	task := Task{
		ID:          uuid.New(),
		ProjectID:   uuid.New(),
		Title:       "Test Task",
		Description: "Test Description",
		Status:      "pending",
		Priority:    "high",
		CreatedBy:   uuid.New(),
		AssignedTo:  &assignedTo,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if task.Title != "Test Task" {
		t.Errorf("Expected task title to be 'Test Task', got '%s'", task.Title)
	}

	if task.Status != "pending" {
		t.Errorf("Expected status to be 'pending', got '%s'", task.Status)
	}

	if task.Priority != "high" {
		t.Errorf("Expected priority to be 'high', got '%s'", task.Priority)
	}
}

func TestContextModel(t *testing.T) {
	taskID := uuid.New()
	ctx := Context{
		ID:        uuid.New(),
		ProjectID: uuid.New(),
		AgentID:   uuid.New(),
		TaskID:    &taskID,
		Title:     "Test Context",
		Content:   "Test Content",
		Tags:      []string{"api", "documentation"},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if ctx.Title != "Test Context" {
		t.Errorf("Expected context title to be 'Test Context', got '%s'", ctx.Title)
	}

	if len(ctx.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(ctx.Tags))
	}

	if ctx.Tags[0] != "api" {
		t.Errorf("Expected first tag to be 'api', got '%s'", ctx.Tags[0])
	}
}
