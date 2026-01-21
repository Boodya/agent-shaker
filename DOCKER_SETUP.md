# Docker Setup - Vue.js Application

## Architecture

The application now runs with **3 Docker containers**:

1. **postgres** - PostgreSQL 16 database
2. **mcp-server** - Go backend API server (port 8080)
3. **web** - Vue.js frontend with Nginx (port 80)

## Container Details

### 1. PostgreSQL Database
- **Image**: `postgres:16-alpine`
- **Port**: 5433:5432
- **Database**: mcp_tracker
- **User**: mcp
- **Password**: secret
- **Volume**: postgres_data (persistent storage)

### 2. Go Backend (mcp-server)
- **Build**: From root Dockerfile
- **Port**: 8080
- **API Endpoints**: `/api/*`
- **WebSocket**: `/ws`
- **Dependencies**: Waits for PostgreSQL to be healthy

### 3. Vue.js Frontend (web)
- **Build**: From `web/Dockerfile`
- **Port**: 80 (accessible via http://localhost)
- **Build Process**: 
  - Stage 1: Node.js 18 Alpine - builds Vue app with Vite
  - Stage 2: Nginx Alpine - serves static files
- **Nginx Configuration**: 
  - Proxies `/api/*` requests to mcp-server:8080
  - Proxies `/ws` WebSocket connections to mcp-server:8080
  - Serves Vue.js SPA with proper routing

## How to Run

### Start all services:
```powershell
docker-compose up -d
```

### Build and start (if changes were made):
```powershell
docker-compose up --build -d
```

### Stop all services:
```powershell
docker-compose down
```

### View logs:
```powershell
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f web
docker-compose logs -f mcp-server
docker-compose logs -f postgres
```

### Check status:
```powershell
docker-compose ps
```

## Access Points

- **Frontend**: http://localhost
- **Backend API**: http://localhost:8080/api
- **WebSocket**: ws://localhost:8080/ws
- **PostgreSQL**: localhost:5433

## Vue.js Application

### Structure
```
web/
├── Dockerfile          # Multi-stage build: Node.js + Nginx
├── nginx.conf          # Nginx configuration with API proxy
├── package.json        # NPM dependencies
├── vite.config.js      # Vite configuration
├── index.html          # HTML entry point
└── src/
    ├── main.js         # Vue app entry
    ├── App.vue         # Root component
    ├── router/         # Vue Router
    ├── stores/         # Pinia stores
    ├── services/       # API services
    ├── views/          # Page components
    └── composables/    # Vue composables
```

### Key Features
- **Vue 3** with Composition API
- **Vite** for fast development and optimized builds
- **Vue Router** for SPA routing
- **Pinia** for state management
- **Axios** for HTTP requests
- **Marked** for Markdown parsing
- **DOMPurify** for XSS protection

### Development
To run locally without Docker:
```powershell
cd web
npm install
npm run dev
```
Access at http://localhost:3000

### Production Build
The Docker container automatically builds for production:
```powershell
npm run build
```
Outputs to `web/dist/` and served by Nginx.

## Network Flow

```
Browser (http://localhost)
    ↓
Nginx (web container, port 80)
    ├── Static files → Vue.js SPA
    ├── /api/* → Go Backend (mcp-server:8080)
    └── /ws → WebSocket (mcp-server:8080)
        ↓
Go Backend (mcp-server container, port 8080)
    ↓
PostgreSQL (postgres container, port 5432)
```

## Benefits of This Setup

1. **Separation of Concerns**: Frontend, backend, and database are isolated
2. **Scalability**: Each service can be scaled independently
3. **Development Friendly**: Can run services individually for debugging
4. **Production Ready**: Optimized builds with Nginx serving static files
5. **Proper Routing**: Nginx handles SPA routing and API proxying
6. **WebSocket Support**: Full duplex communication for real-time updates

## Troubleshooting

### Check if containers are running:
```powershell
docker-compose ps
```

### View container logs:
```powershell
docker-compose logs -f web
```

### Restart a specific service:
```powershell
docker-compose restart web
```

### Rebuild a service:
```powershell
docker-compose up --build web
```

### Access container shell:
```powershell
docker-compose exec web sh
docker-compose exec mcp-server sh
docker-compose exec postgres sh
```

### Clean up everything:
```powershell
docker-compose down -v  # Also removes volumes
```
