package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// DocsHandler handles documentation requests
type DocsHandler struct {
	githubBaseURL string
}

// NewDocsHandler creates a new documentation handler
func NewDocsHandler() *DocsHandler {
	return &DocsHandler{
		githubBaseURL: "https://raw.githubusercontent.com/techbuzzz/agent-shaker/main",
	}
}

// GetDoc serves markdown documentation files from GitHub
func (h *DocsHandler) GetDoc(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	docPath := vars["path"]

	if docPath == "" {
		http.Error(w, "Missing path", http.StatusBadRequest)
		return
	}

	// Security: prevent directory traversal
	if strings.Contains(docPath, "..") {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	// Normalize path
	docPath = strings.TrimPrefix(docPath, "/")

	// Try organized location first (docs/<...>), then the legacy path
	tryPaths := []string{
		fmt.Sprintf("docs/%s", docPath),
		docPath,
	}

	var lastErr error
	for _, p := range tryPaths {
		githubURL := fmt.Sprintf("%s/%s", h.githubBaseURL, p)
		resp, err := http.Get(githubURL)
		if err != nil {
			lastErr = err
			continue
		}

		func() {
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusNotFound {
				// try next path
				return
			}

			if resp.StatusCode != http.StatusOK {
				lastErr = fmt.Errorf("unexpected status %d", resp.StatusCode)
				return
			}

			content, err := io.ReadAll(resp.Body)
			if err != nil {
				lastErr = err
				return
			}

			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write(content)
			lastErr = nil
		}()

		if lastErr == nil {
			return
		}
	}

	if lastErr != nil {
		http.Error(w, "Failed to fetch document from GitHub", http.StatusInternalServerError)
		return
	}

	http.Error(w, "Document not found", http.StatusNotFound)
}

// ListDocs returns a list of available documentation files with GitHub URLs
func (h *DocsHandler) ListDocs(w http.ResponseWriter, r *http.Request) {
	// The API advertises organized locations under `docs/<category>/` while
	// keeping `legacy_path` fields so consumers can migrate gradually.
	docs := []map[string]string{
		{
			"name":           "README",
			"path":           "docs/getting-started/README.md",
			"legacy_path":    "README.md",
			"url":            "https://github.com/techbuzzz/agent-shaker/blob/main/README.md",
			"raw_url":        fmt.Sprintf("%s/docs/getting-started/README.md", h.githubBaseURL),
			"legacy_raw_url": fmt.Sprintf("%s/README.md", h.githubBaseURL),
			"category":       "Getting Started",
		},
		{
			"name":           "Quick Start Guide",
			"path":           "docs/getting-started/QUICKSTART_GUIDE.md",
			"legacy_path":    "docs/QUICKSTART_GUIDE.md",
			"url":            "https://github.com/techbuzzz/agent-shaker/blob/main/docs/QUICKSTART_GUIDE.md",
			"raw_url":        fmt.Sprintf("%s/docs/getting-started/QUICKSTART_GUIDE.md", h.githubBaseURL),
			"legacy_raw_url": fmt.Sprintf("%s/docs/QUICKSTART_GUIDE.md", h.githubBaseURL),
			"category":       "Getting Started",
		},
		{
			"name":           "Quick Start - Agent",
			"path":           "docs/getting-started/QUICKSTART_AGENT.md",
			"legacy_path":    "docs/QUICKSTART_AGENT.md",
			"url":            "https://github.com/techbuzzz/agent-shaker/blob/main/docs/QUICKSTART_AGENT.md",
			"raw_url":        fmt.Sprintf("%s/docs/getting-started/QUICKSTART_AGENT.md", h.githubBaseURL),
			"legacy_raw_url": fmt.Sprintf("%s/docs/QUICKSTART_AGENT.md", h.githubBaseURL),
			"category":       "Getting Started",
		},
		{
			"name":           "Agent Setup Guide",
			"path":           "docs/configuration/AGENT_SETUP_GUIDE.md",
			"legacy_path":    "docs/AGENT_SETUP_GUIDE.md",
			"url":            "https://github.com/techbuzzz/agent-shaker/blob/main/docs/AGENT_SETUP_GUIDE.md",
			"raw_url":        fmt.Sprintf("%s/docs/configuration/AGENT_SETUP_GUIDE.md", h.githubBaseURL),
			"legacy_raw_url": fmt.Sprintf("%s/docs/AGENT_SETUP_GUIDE.md", h.githubBaseURL),
			"category":       "Configuration",
		},
		{
			"name":           "Architecture",
			"path":           "docs/architecture/ARCHITECTURE.md",
			"legacy_path":    "docs/ARCHITECTURE.md",
			"url":            "https://github.com/techbuzzz/agent-shaker/blob/main/docs/ARCHITECTURE.md",
			"raw_url":        fmt.Sprintf("%s/docs/architecture/ARCHITECTURE.md", h.githubBaseURL),
			"legacy_raw_url": fmt.Sprintf("%s/docs/ARCHITECTURE.md", h.githubBaseURL),
			"category":       "Architecture",
		},
		{
			"name":           "Architecture Summary",
			"path":           "docs/architecture/ARCHITECTURE_SUMMARY.md",
			"legacy_path":    "docs/ARCHITECTURE_SUMMARY.md",
			"url":            "https://github.com/techbuzzz/agent-shaker/blob/main/docs/ARCHITECTURE_SUMMARY.md",
			"raw_url":        fmt.Sprintf("%s/docs/architecture/ARCHITECTURE_SUMMARY.md", h.githubBaseURL),
			"legacy_raw_url": fmt.Sprintf("%s/docs/ARCHITECTURE_SUMMARY.md", h.githubBaseURL),
			"category":       "Architecture",
		},
		{
			"name":           "Clean Architecture",
			"path":           "docs/architecture/CLEAN_ARCHITECTURE.md",
			"legacy_path":    "docs/CLEAN_ARCHITECTURE.md",
			"url":            "https://github.com/techbuzzz/agent-shaker/blob/main/docs/CLEAN_ARCHITECTURE.md",
			"raw_url":        fmt.Sprintf("%s/docs/architecture/CLEAN_ARCHITECTURE.md", h.githubBaseURL),
			"legacy_raw_url": fmt.Sprintf("%s/docs/CLEAN_ARCHITECTURE.md", h.githubBaseURL),
			"category":       "Architecture",
		},
		{
			"name":           "API Documentation",
			"path":           "docs/api/API.md",
			"legacy_path":    "docs/API.md",
			"url":            "https://github.com/techbuzzz/agent-shaker/blob/main/docs/API.md",
			"raw_url":        fmt.Sprintf("%s/docs/api/API.md", h.githubBaseURL),
			"legacy_raw_url": fmt.Sprintf("%s/docs/API.md", h.githubBaseURL),
			"category":       "API Reference",
		},
		{
			"name":           "MCP Quick Start",
			"path":           "docs/mcp/MCP_QUICKSTART.md",
			"legacy_path":    "docs/MCP_QUICKSTART.md",
			"url":            "https://github.com/techbuzzz/agent-shaker/blob/main/docs/MCP_QUICKSTART.md",
			"raw_url":        fmt.Sprintf("%s/docs/mcp/MCP_QUICKSTART.md", h.githubBaseURL),
			"legacy_raw_url": fmt.Sprintf("%s/docs/MCP_QUICKSTART.md", h.githubBaseURL),
			"category":       "MCP Integration",
		},
		{
			"name":           "Copilot Integration",
			"path":           "docs/mcp/COPILOT_INTEGRATION.md",
			"legacy_path":    "docs/COPILOT_INTEGRATION.md",
			"url":            "https://github.com/techbuzzz/agent-shaker/blob/main/docs/COPILOT_INTEGRATION.md",
			"raw_url":        fmt.Sprintf("%s/docs/mcp/COPILOT_INTEGRATION.md", h.githubBaseURL),
			"legacy_raw_url": fmt.Sprintf("%s/docs/COPILOT_INTEGRATION.md", h.githubBaseURL),
			"category":       "MCP Integration",
		},
		{
			"name":           "Copilot MCP Integration",
			"path":           "docs/mcp/COPILOT_MCP_INTEGRATION.md",
			"legacy_path":    "docs/COPILOT_MCP_INTEGRATION.md",
			"url":            "https://github.com/techbuzzz/agent-shaker/blob/main/docs/COPILOT_MCP_INTEGRATION.md",
			"raw_url":        fmt.Sprintf("%s/docs/mcp/COPILOT_MCP_INTEGRATION.md", h.githubBaseURL),
			"legacy_raw_url": fmt.Sprintf("%s/docs/COPILOT_MCP_INTEGRATION.md", h.githubBaseURL),
			"category":       "MCP Integration",
		},
		{
			"name":           "Docker Deployment",
			"path":           "docs/deployment/DOCKER_DEPLOYMENT.md",
			"legacy_path":    "docs/DOCKER_DEPLOYMENT.md",
			"url":            "https://github.com/techbuzzz/agent-shaker/blob/main/docs/DOCKER_DEPLOYMENT.md",
			"raw_url":        fmt.Sprintf("%s/docs/deployment/DOCKER_DEPLOYMENT.md", h.githubBaseURL),
			"legacy_raw_url": fmt.Sprintf("%s/docs/DOCKER_DEPLOYMENT.md", h.githubBaseURL),
			"category":       "Deployment",
		},
		{
			"name":           "Docker Setup",
			"path":           "docs/deployment/DOCKER_SETUP.md",
			"legacy_path":    "docs/DOCKER_SETUP.md",
			"url":            "https://github.com/techbuzzz/agent-shaker/blob/main/docs/DOCKER_SETUP.md",
			"raw_url":        fmt.Sprintf("%s/docs/deployment/DOCKER_SETUP.md", h.githubBaseURL),
			"legacy_raw_url": fmt.Sprintf("%s/docs/DOCKER_SETUP.md", h.githubBaseURL),
			"category":       "Deployment",
		},
		{
			"name":           "Contributing",
			"path":           "docs/development/CONTRIBUTING.md",
			"legacy_path":    "docs/CONTRIBUTING.md",
			"url":            "https://github.com/techbuzzz/agent-shaker/blob/main/docs/CONTRIBUTING.md",
			"raw_url":        fmt.Sprintf("%s/docs/development/CONTRIBUTING.md", h.githubBaseURL),
			"legacy_raw_url": fmt.Sprintf("%s/docs/CONTRIBUTING.md", h.githubBaseURL),
			"category":       "Development",
		},
		{
			"name":           "Implementation Summary",
			"path":           "docs/development/IMPLEMENTATION_SUMMARY.md",
			"legacy_path":    "docs/IMPLEMENTATION_SUMMARY.md",
			"url":            "https://github.com/techbuzzz/agent-shaker/blob/main/docs/IMPLEMENTATION_SUMMARY.md",
			"raw_url":        fmt.Sprintf("%s/docs/development/IMPLEMENTATION_SUMMARY.md", h.githubBaseURL),
			"legacy_raw_url": fmt.Sprintf("%s/docs/IMPLEMENTATION_SUMMARY.md", h.githubBaseURL),
			"category":       "Development",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"docs":       docs,
		"repository": "https://github.com/techbuzzz/agent-shaker",
		"branch":     "main",
	})
}
