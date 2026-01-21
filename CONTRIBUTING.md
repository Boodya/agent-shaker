# Contributing to MCP Task Tracker

Thank you for considering contributing to MCP Task Tracker! This document provides guidelines for contributing to the project.

## Getting Started

1. **Fork the repository**
   ```bash
   git clone https://github.com/techbuzzz/agent-shaker.git
   cd agent-shaker
   ```

2. **Set up development environment**
   ```bash
   # Install Go 1.21 or later
   # Install PostgreSQL 15 or later
   # Install Docker and Docker Compose (optional)
   ```

3. **Install dependencies**
   ```bash
   make deps
   ```

4. **Start development services**
   ```bash
   docker-compose up -d postgres
   export DATABASE_URL="postgres://mcp:secret@localhost:5432/mcp_tracker?sslmode=disable"
   ```

## Development Workflow

### Building

```bash
make build
```

### Running locally

```bash
make run
```

Or with Docker:

```bash
make docker-up
```

### Testing

```bash
make test
```

### Code formatting

```bash
make fmt
```

### Linting

```bash
make lint
```

## Project Structure

```
agent-shaker/
├── cmd/server/          # Main application entry point
├── internal/
│   ├── database/        # Database connection
│   ├── handlers/        # HTTP handlers for API endpoints
│   ├── models/          # Data models
│   └── websocket/       # WebSocket hub
├── migrations/          # Database migrations
├── web/
│   └── static/          # Web UI files
├── docs/                # Documentation
└── docker-compose.yml   # Docker setup
```

## Adding New Features

### Adding a new API endpoint

1. **Add model** (if needed) in `internal/models/`
2. **Add handler** in `internal/handlers/`
3. **Register route** in `cmd/server/main.go`
4. **Update documentation** in `docs/API.md`
5. **Add tests**

### Example: Adding a new entity

```go
// internal/models/comment.go
type Comment struct {
    ID        uuid.UUID `json:"id" db:"id"`
    TaskID    uuid.UUID `json:"task_id" db:"task_id"`
    Content   string    `json:"content" db:"content"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// internal/handlers/comments.go
type CommentHandler struct {
    db  *database.DB
    hub *websocket.Hub
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
    // Implementation
}

// cmd/server/main.go
commentHandler := handlers.NewCommentHandler(db, hub)
api.HandleFunc("/comments", commentHandler.CreateComment).Methods("POST")
```

## Database Migrations

1. Create a new migration file in `migrations/` with format `XXX_description.sql`
2. Update the migration runner in `cmd/server/main.go`

Example migration:

```sql
-- migrations/002_add_comments.sql
CREATE TABLE IF NOT EXISTS comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    task_id UUID NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    content TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_comments_task ON comments(task_id);
```

## WebSocket Events

When adding new features that should trigger real-time updates:

```go
// Broadcast to project
h.hub.BroadcastToProject(projectID, "event_type", payload)
```

Event types should be documented in `docs/API.md`.

## Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Add comments for exported functions
- Keep functions small and focused
- Use meaningful variable names

## Testing Guidelines

- Write unit tests for business logic
- Test error cases
- Use table-driven tests where appropriate
- Mock external dependencies

Example test:

```go
func TestCreateProject(t *testing.T) {
    tests := []struct {
        name    string
        input   CreateProjectRequest
        wantErr bool
    }{
        {
            name: "valid project",
            input: CreateProjectRequest{
                Name: "Test Project",
                Description: "Test Description",
            },
            wantErr: false,
        },
        // Add more test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## Documentation

- Update README.md for user-facing changes
- Update API.md for API changes
- Add code comments for complex logic
- Include examples in documentation

## Pull Request Process

1. **Create a feature branch**
   ```bash
   git checkout -b feature/my-new-feature
   ```

2. **Make your changes**
   - Write code
   - Add tests
   - Update documentation

3. **Test your changes**
   ```bash
   make test
   make build
   ```

4. **Commit your changes**
   ```bash
   git add .
   git commit -m "feat: add my new feature"
   ```

   Follow [Conventional Commits](https://www.conventionalcommits.org/):
   - `feat:` - New feature
   - `fix:` - Bug fix
   - `docs:` - Documentation changes
   - `style:` - Code style changes
   - `refactor:` - Code refactoring
   - `test:` - Adding tests
   - `chore:` - Maintenance tasks

5. **Push to your fork**
   ```bash
   git push origin feature/my-new-feature
   ```

6. **Create a Pull Request**
   - Provide a clear description
   - Reference any related issues
   - Include screenshots for UI changes

## Bug Reports

When reporting bugs, please include:

1. **Description** - Clear description of the bug
2. **Steps to reproduce** - Detailed steps to reproduce the issue
3. **Expected behavior** - What you expected to happen
4. **Actual behavior** - What actually happened
5. **Environment** - OS, Go version, etc.
6. **Logs** - Relevant error messages or logs

## Feature Requests

When requesting features, please include:

1. **Use case** - Why is this feature needed?
2. **Proposed solution** - How should it work?
3. **Alternatives** - Other approaches you've considered
4. **Additional context** - Any other relevant information

## Code of Conduct

- Be respectful and inclusive
- Focus on constructive feedback
- Help others learn and grow
- Maintain a positive community

## Questions?

- Open an issue for questions
- Check existing documentation
- Review closed issues and PRs

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

Thank you for contributing to MCP Task Tracker! 🚀
