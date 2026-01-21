package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/techbuzzz/agent-shaker/internal/database"
	"github.com/techbuzzz/agent-shaker/internal/handlers"
	"github.com/techbuzzz/agent-shaker/internal/middleware"
	"github.com/techbuzzz/agent-shaker/internal/websocket"
)

func main() {
	// Get database URL from environment
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://mcp:secret@localhost:5432/mcp_tracker?sslmode=disable"
	}

	// Connect to database
	db, err := database.NewDB(databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("Connected to database")

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Create WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Create handlers
	projectHandler := handlers.NewProjectHandler(db, hub)
	agentHandler := handlers.NewAgentHandler(db, hub)
	taskHandler := handlers.NewTaskHandler(db, hub)
	contextHandler := handlers.NewContextHandler(db, hub)
	wsHandler := handlers.NewWebSocketHandler(hub)
	docsHandler := handlers.NewDocsHandler()

	// Setup router
	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api").Subrouter()

	// Projects
	api.HandleFunc("/projects", projectHandler.CreateProject).Methods("POST")
	api.HandleFunc("/projects", projectHandler.ListProjects).Methods("GET")
	api.HandleFunc("/projects/{id}", projectHandler.GetProject).Methods("GET")

	// Agents
	api.HandleFunc("/agents", agentHandler.CreateAgent).Methods("POST")
	api.HandleFunc("/agents", agentHandler.ListAgents).Methods("GET")
	api.HandleFunc("/agents/{id}/status", agentHandler.UpdateAgentStatus).Methods("PUT")

	// Tasks
	api.HandleFunc("/tasks", taskHandler.CreateTask).Methods("POST")
	api.HandleFunc("/tasks", taskHandler.ListTasks).Methods("GET")
	api.HandleFunc("/tasks/{id}", taskHandler.GetTask).Methods("GET")
	api.HandleFunc("/tasks/{id}", taskHandler.UpdateTask).Methods("PUT")

	// Contexts
	api.HandleFunc("/contexts", contextHandler.CreateContext).Methods("POST")
	api.HandleFunc("/contexts", contextHandler.ListContexts).Methods("GET")
	api.HandleFunc("/contexts/{id}", contextHandler.GetContext).Methods("GET")
	api.HandleFunc("/contexts/{id}", contextHandler.UpdateContext).Methods("PUT")
	api.HandleFunc("/contexts/{id}", contextHandler.DeleteContext).Methods("DELETE")

	// Documentation
	api.HandleFunc("/docs", docsHandler.ListDocs).Methods("GET")
	api.HandleFunc("/docs/{path:.*}", docsHandler.GetDoc).Methods("GET")

	// WebSocket
	r.HandleFunc("/ws", wsHandler.HandleWebSocket)

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Serve Vue.js app (try dist first for production, fallback to static for old version)
	// Check if Vue.js build exists
	if _, err := os.Stat("./web/dist"); err == nil {
		// Serve Vue.js SPA
		spa := spaHandler{staticPath: "./web/dist", indexPath: "index.html"}
		r.PathPrefix("/").Handler(spa)
	} else {
		// Fallback to old static files
		r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/static")))
	}

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// Apply middleware
	handler := middleware.Recovery(
		middleware.Logger(
			middleware.RequestSizeLimit(10 * 1024 * 1024)( // 10MB limit
				c.Handler(r),
			),
		),
	)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func runMigrations(db *database.DB) error {
	migrationSQL, err := os.ReadFile("migrations/001_init.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(migrationSQL))
	return err
}

// spaHandler implements the http.Handler interface for serving a SPA
type spaHandler struct {
	staticPath string
	indexPath  string
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get the absolute path to prevent directory traversal
	path := r.URL.Path

	// Check if file exists
	fullPath := h.staticPath + path
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		// File does not exist, serve index.html
		http.ServeFile(w, r, h.staticPath+"/"+h.indexPath)
		return
	}

	// Otherwise, serve the file
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}
