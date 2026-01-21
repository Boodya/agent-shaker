# Multi-stage Dockerfile for MCP Task Tracker with Vue.js Frontend

# Stage 1: Build Vue.js frontend
FROM node:18-alpine AS frontend-builder

WORKDIR /app/web

# Copy package files
COPY web/package*.json ./

# Install dependencies
RUN npm install

# Copy source files
COPY web/ ./

# Build for production
RUN npm run build

# Stage 2: Build Go backend
FROM golang:1.24-alpine AS backend-builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /mcp-server ./cmd/server

# Stage 3: Final image
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy backend binary
COPY --from=backend-builder /mcp-server .

# Copy migrations
COPY migrations ./migrations

# Copy Vue.js built files
COPY --from=frontend-builder /app/web/dist ./web/dist

# Expose port
EXPOSE 8080

# Run the application
CMD ["./mcp-server"]
