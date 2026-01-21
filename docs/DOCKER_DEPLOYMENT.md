# 🐳 Docker Deployment Guide

This guide covers multiple ways to deploy the MCP Task Tracker with the Vue.js frontend in Docker.

## 📦 Deployment Options

### Option 1: Single Container (Backend serves Frontend)
**Best for:** Simple deployments, minimal infrastructure

The Go backend serves the built Vue.js files directly.

```bash
# Build and run
docker-compose -f docker-compose.vue.yml up -d

# Access the app
# http://localhost:8080
```

**Architecture:**
```
┌─────────────────────────────────────┐
│  Container: mcp-server-vue          │
│  ┌───────────────────────────────┐  │
│  │  Go Backend (Port 8080)       │  │
│  │  - REST API (/api/*)          │  │
│  │  - WebSocket (/ws)            │  │
│  │  - Static Files (Vue.js)      │  │
│  └───────────────────────────────┘  │
└─────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────┐
│  Container: postgres                │
│  PostgreSQL 16 (Port 5433)          │
└─────────────────────────────────────┘
```

### Option 2: Separate Containers (Nginx + Backend)
**Best for:** Production, scaling, better separation

The Vue.js app runs in an Nginx container and proxies API requests to the backend.

```bash
# Build and run
docker-compose -f docker-compose.web.yml up -d

# Access the app
# http://localhost (port 80)
# Backend API: http://localhost:8080
```

**Architecture:**
```
┌─────────────────────────────────────┐
│  Container: mcp-web (Nginx)         │
│  - Vue.js SPA (Port 80)             │
│  - Proxies /api/* → Backend         │
│  - Proxies /ws → Backend            │
└─────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────┐
│  Container: mcp-server              │
│  Go Backend (Port 8080)             │
│  - REST API                         │
│  - WebSocket                        │
└─────────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────┐
│  Container: postgres                │
│  PostgreSQL 16 (Port 5433)          │
└─────────────────────────────────────┘
```

### Option 3: Development Mode
**Best for:** Active development with Hot Module Replacement

```bash
# Start backend
docker-compose up -d

# Start Vue.js dev server (separate terminal)
cd web
npm run dev

# Access the app
# http://localhost:3000 (Vue.js with HMR)
# Backend API: http://localhost:8080
```

## 🚀 Quick Start

### Option 1 (Recommended for Simple Setup)

```powershell
# Stop any running containers
docker-compose down

# Build and start (single container)
docker-compose -f docker-compose.vue.yml up -d --build

# Wait for services to be healthy
Start-Sleep -Seconds 10

# Check health
curl http://localhost:8080/health

# Open browser
start http://localhost:8080
```

### Option 2 (Recommended for Production)

```powershell
# Stop any running containers
docker-compose down

# Build and start (separate containers)
docker-compose -f docker-compose.web.yml up -d --build

# Wait for services to be healthy
Start-Sleep -Seconds 15

# Check health
curl http://localhost:8080/health
curl http://localhost/

# Open browser
start http://localhost
```

## 📋 Detailed Instructions

### Building the Images

#### Option 1: Single Container
```bash
# Build the multi-stage image
docker build -f Dockerfile.vue -t mcp-tracker:vue .

# Check the image
docker images | grep mcp-tracker
```

#### Option 2: Separate Containers
```bash
# Build backend
docker build -t mcp-tracker:backend .

# Build frontend
cd web
docker build -t mcp-tracker:web .
cd ..
```

### Running the Containers

#### Option 1: Single Container
```bash
# Start services
docker-compose -f docker-compose.vue.yml up -d

# View logs
docker-compose -f docker-compose.vue.yml logs -f

# Stop services
docker-compose -f docker-compose.vue.yml down
```

#### Option 2: Separate Containers
```bash
# Start services
docker-compose -f docker-compose.web.yml up -d

# View logs
docker-compose -f docker-compose.web.yml logs -f mcp-web
docker-compose -f docker-compose.web.yml logs -f mcp-server

# Stop services
docker-compose -f docker-compose.web.yml down
```

## 🔧 Configuration

### Environment Variables

Create a `.env` file in the root directory:

```env
# Database
POSTGRES_USER=mcp
POSTGRES_PASSWORD=secret
POSTGRES_DB=mcp_tracker

# Backend
PORT=8080
DATABASE_URL=postgres://mcp:secret@postgres:5432/mcp_tracker?sslmode=disable

# Frontend (for nginx)
BACKEND_URL=http://mcp-server:8080
```

### Custom Ports

#### Option 1 (Single Container)
Edit `docker-compose.vue.yml`:
```yaml
services:
  mcp-server:
    ports:
      - "3000:8080"  # Change 8080 to 3000
```

#### Option 2 (Separate Containers)
Edit `docker-compose.web.yml`:
```yaml
services:
  web:
    ports:
      - "3000:80"  # Change 80 to 3000
```

## 🧪 Testing the Deployment

### Health Checks

```bash
# Check backend health
curl http://localhost:8080/health
# Expected: "OK"

# Check frontend (Option 1)
curl http://localhost:8080/
# Expected: HTML content

# Check frontend (Option 2)
curl http://localhost/
# Expected: HTML content

# Check API
curl http://localhost:8080/api/projects
# Expected: JSON array

# Check PostgreSQL
docker exec -it mcp-postgres psql -U mcp -d mcp_tracker -c "SELECT COUNT(*) FROM projects;"
```

### Monitoring

```bash
# View all running containers
docker-compose -f docker-compose.web.yml ps

# View resource usage
docker stats

# Follow logs
docker-compose -f docker-compose.web.yml logs -f

# View specific service logs
docker logs mcp-web -f
docker logs mcp-server -f
docker logs mcp-postgres -f
```

## 🐛 Troubleshooting

### Port Conflicts

```powershell
# Check what's using port 8080
Get-NetTCPConnection -LocalPort 8080

# Kill process
Stop-Process -Id <PID> -Force

# Or change the port in docker-compose
```

### Database Connection Issues

```bash
# Check if PostgreSQL is running
docker-compose ps postgres

# Check PostgreSQL logs
docker logs mcp-postgres

# Connect to PostgreSQL
docker exec -it mcp-postgres psql -U mcp -d mcp_tracker

# Test connection from backend
docker exec -it mcp-server sh
# Inside container:
nc -zv postgres 5432
```

### Frontend Not Loading

```bash
# Check if frontend files exist (Option 1)
docker exec -it mcp-server-vue ls -la /app/web/dist

# Check nginx logs (Option 2)
docker logs mcp-web

# Test nginx config (Option 2)
docker exec -it mcp-web nginx -t

# Rebuild frontend
docker-compose -f docker-compose.web.yml up -d --build web
```

### WebSocket Not Connecting

Check browser console for errors. Common issues:

1. **CORS**: Check backend CORS configuration
2. **Proxy**: Verify nginx proxy settings (Option 2)
3. **Connection**: Ensure backend WebSocket endpoint is accessible

```bash
# Test WebSocket from command line
npm install -g wscat
wscat -c ws://localhost:8080/ws
```

## 🔐 Security Best Practices

### Production Recommendations

1. **Use secrets management**
```yaml
secrets:
  db_password:
    external: true
```

2. **Don't expose PostgreSQL port**
```yaml
# Remove this in production:
# ports:
#   - "5433:5432"
```

3. **Use HTTPS with reverse proxy**
```bash
# Use nginx or traefik as reverse proxy
# Add SSL certificates
```

4. **Set resource limits**
```yaml
services:
  mcp-server:
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
```

5. **Use read-only root filesystem**
```yaml
services:
  mcp-server:
    read_only: true
    tmpfs:
      - /tmp
```

## 📊 Performance Optimization

### Backend Optimization

```yaml
services:
  mcp-server:
    environment:
      GOMAXPROCS: 2
      GOMEMLIMIT: 512MiB
```

### Database Optimization

```yaml
services:
  postgres:
    command: >
      postgres
      -c shared_buffers=256MB
      -c max_connections=100
      -c work_mem=4MB
```

### Frontend Optimization (Nginx)

```nginx
# In nginx.conf
gzip on;
gzip_comp_level 6;
gzip_types text/plain text/css application/json application/javascript;

# Browser caching
location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
    expires 1y;
    add_header Cache-Control "public, immutable";
}
```

## 🔄 Updates and Maintenance

### Updating the Application

```bash
# Pull latest changes
git pull

# Rebuild and restart
docker-compose -f docker-compose.web.yml up -d --build

# Or force recreate
docker-compose -f docker-compose.web.yml up -d --force-recreate --build
```

### Database Backup

```bash
# Backup
docker exec mcp-postgres pg_dump -U mcp mcp_tracker > backup.sql

# Restore
cat backup.sql | docker exec -i mcp-postgres psql -U mcp -d mcp_tracker
```

### Cleaning Up

```bash
# Stop and remove containers
docker-compose -f docker-compose.web.yml down

# Remove volumes (WARNING: deletes data!)
docker-compose -f docker-compose.web.yml down -v

# Remove images
docker rmi mcp-tracker:backend mcp-tracker:web

# Clean up everything
docker system prune -a
```

## 📈 Scaling

### Horizontal Scaling

```yaml
services:
  mcp-server:
    deploy:
      replicas: 3
      
  web:
    deploy:
      replicas: 2
```

### Load Balancer

Add nginx load balancer:
```yaml
services:
  nginx-lb:
    image: nginx:alpine
    ports:
      - "80:80"
    volumes:
      - ./nginx-lb.conf:/etc/nginx/nginx.conf
```

## 🎯 Summary

**For Development:**
```bash
docker-compose up -d && cd web && npm run dev
```

**For Simple Production:**
```bash
docker-compose -f docker-compose.vue.yml up -d
```

**For Advanced Production:**
```bash
docker-compose -f docker-compose.web.yml up -d
```

Choose the option that best fits your needs! 🚀
