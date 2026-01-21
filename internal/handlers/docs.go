package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
)

// DocsHandler handles documentation requests
type DocsHandler struct{}

// NewDocsHandler creates a new documentation handler
func NewDocsHandler() *DocsHandler {
	return &DocsHandler{}
}

// GetDoc serves markdown documentation files
func (h *DocsHandler) GetDoc(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	docPath := vars["path"]

	// Security: prevent directory traversal
	if strings.Contains(docPath, "..") {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	// Construct full path
	fullPath := filepath.Join(".", docPath)

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	// Read file
	content, err := os.ReadFile(fullPath)
	if err != nil {
		http.Error(w, "Failed to read document", http.StatusInternalServerError)
		return
	}

	// Return markdown content as plain text
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}

// ListDocs returns a list of available documentation files
func (h *DocsHandler) ListDocs(w http.ResponseWriter, r *http.Request) {
	docs := []map[string]string{
		{"name": "README", "path": "README.md"},
		{"name": "Quick Start", "path": "QUICKSTART.md"},
		{"name": "Agent Setup Guide", "path": "AGENT_SETUP_GUIDE.md"},
		{"name": "Quick Start - Agent", "path": "QUICKSTART_AGENT.md"},
		{"name": "Architecture", "path": "ARCHITECTURE.md"},
		{"name": "API Documentation", "path": "docs/API.md"},
		{"name": "Copilot Integration", "path": "docs/COPILOT_INTEGRATION.md"},
		{"name": "Implementation Summary", "path": "IMPLEMENTATION_SUMMARY.md"},
		{"name": "Docker Deployment", "path": "DOCKER_DEPLOYMENT.md"},
		{"name": "Contributing", "path": "CONTRIBUTING.md"},
		{"name": "Improvements", "path": "IMPROVEMENTS.md"},
		{"name": "Vue Setup", "path": "VUE_SETUP_COMPLETE.md"},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(docs)
}
