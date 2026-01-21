#!/bin/bash
# MCP Task Tracker Demo Script
# This script demonstrates the API usage

set -e

API_BASE="http://localhost:8080/api"

echo "🚀 MCP Task Tracker API Demo"
echo "=============================="
echo ""

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if server is running
echo "Checking server health..."
if ! curl -s -f "${API_BASE%/api}/health" > /dev/null; then
    echo "❌ Server is not running. Please start it with: docker-compose up -d"
    exit 1
fi
echo -e "${GREEN}✓${NC} Server is running"
echo ""

# 1. Create a project
echo -e "${BLUE}1. Creating a project...${NC}"
PROJECT_RESPONSE=$(curl -s -X POST "$API_BASE/projects" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "InvoiceAI Demo",
    "description": "AI-powered invoice processing system"
  }')
PROJECT_ID=$(echo $PROJECT_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
echo "   Project created: $PROJECT_ID"
echo ""

# 2. Register backend agent
echo -e "${BLUE}2. Registering backend agent...${NC}"
BACKEND_RESPONSE=$(curl -s -X POST "$API_BASE/agents" \
  -H "Content-Type: application/json" \
  -d "{
    \"project_id\": \"$PROJECT_ID\",
    \"name\": \"Backend-Copilot\",
    \"role\": \"backend\",
    \"team\": \"Backend Team\"
  }")
BACKEND_AGENT_ID=$(echo $BACKEND_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
echo "   Backend agent registered: $BACKEND_AGENT_ID"
echo ""

# 3. Register frontend agent
echo -e "${BLUE}3. Registering frontend agent...${NC}"
FRONTEND_RESPONSE=$(curl -s -X POST "$API_BASE/agents" \
  -H "Content-Type: application/json" \
  -d "{
    \"project_id\": \"$PROJECT_ID\",
    \"name\": \"Frontend-Copilot\",
    \"role\": \"frontend\",
    \"team\": \"Frontend Team\"
  }")
FRONTEND_AGENT_ID=$(echo $FRONTEND_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
echo "   Frontend agent registered: $FRONTEND_AGENT_ID"
echo ""

# 4. Create a task for backend
echo -e "${BLUE}4. Creating task for backend...${NC}"
TASK1_RESPONSE=$(curl -s -X POST "$API_BASE/tasks" \
  -H "Content-Type: application/json" \
  -d "{
    \"project_id\": \"$PROJECT_ID\",
    \"title\": \"Implement Invoice REST API\",
    \"description\": \"Create CRUD endpoints for invoice management\",
    \"priority\": \"high\",
    \"created_by\": \"$BACKEND_AGENT_ID\",
    \"assigned_to\": \"$BACKEND_AGENT_ID\"
  }")
TASK1_ID=$(echo $TASK1_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
echo "   Task created: $TASK1_ID"
echo ""

# 5. Update task status to in_progress
echo -e "${BLUE}5. Backend starts working on task...${NC}"
curl -s -X PUT "$API_BASE/tasks/$TASK1_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "in_progress",
    "output": ""
  }' > /dev/null
echo "   Task status updated to 'in_progress'"
echo ""

# 6. Complete the task
echo -e "${BLUE}6. Backend completes the task...${NC}"
curl -s -X PUT "$API_BASE/tasks/$TASK1_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "done",
    "output": "API implemented at /api/invoices with GET, POST, PUT, DELETE endpoints"
  }' > /dev/null
echo "   Task marked as 'done'"
echo ""

# 7. Add API documentation
echo -e "${BLUE}7. Backend adds API documentation...${NC}"
CONTEXT_RESPONSE=$(curl -s -X POST "$API_BASE/contexts" \
  -H "Content-Type: application/json" \
  -d "{
    \"project_id\": \"$PROJECT_ID\",
    \"agent_id\": \"$BACKEND_AGENT_ID\",
    \"task_id\": \"$TASK1_ID\",
    \"title\": \"Invoice API Documentation\",
    \"content\": \"# Invoice API\\n\\n## GET /api/invoices\\nReturns list of invoices\\n\\n## POST /api/invoices\\nCreate new invoice\\n\\n## PUT /api/invoices/{id}\\nUpdate invoice\\n\\n## DELETE /api/invoices/{id}\\nDelete invoice\",
    \"tags\": [\"api\", \"documentation\", \"backend\"]
  }")
CONTEXT_ID=$(echo $CONTEXT_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
echo "   Documentation added: $CONTEXT_ID"
echo ""

# 8. Create task for frontend
echo -e "${BLUE}8. Backend creates task for frontend...${NC}"
TASK2_RESPONSE=$(curl -s -X POST "$API_BASE/tasks" \
  -H "Content-Type: application/json" \
  -d "{
    \"project_id\": \"$PROJECT_ID\",
    \"title\": \"Implement Invoice List UI\",
    \"description\": \"Create UI component to display invoices. API is ready at GET /api/invoices\",
    \"priority\": \"high\",
    \"created_by\": \"$BACKEND_AGENT_ID\",
    \"assigned_to\": \"$FRONTEND_AGENT_ID\"
  }")
TASK2_ID=$(echo $TASK2_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
echo "   Task created for frontend: $TASK2_ID"
echo ""

# 9. List all tasks
echo -e "${BLUE}9. Listing all tasks in project...${NC}"
TASKS=$(curl -s "$API_BASE/tasks?project_id=$PROJECT_ID")
echo "   Found $(echo $TASKS | grep -o '"id"' | wc -l) tasks"
echo ""

# 10. List all agents
echo -e "${BLUE}10. Listing all agents in project...${NC}"
AGENTS=$(curl -s "$API_BASE/agents?project_id=$PROJECT_ID")
echo "   Found $(echo $AGENTS | grep -o '"id"' | wc -l) agents"
echo ""

# 11. Search documentation
echo -e "${BLUE}11. Searching for API documentation...${NC}"
DOCS=$(curl -s "$API_BASE/contexts?project_id=$PROJECT_ID&tags=api")
echo "   Found $(echo $DOCS | grep -o '"id"' | wc -l) documentation entries"
echo ""

echo -e "${GREEN}✓${NC} Demo completed successfully!"
echo ""
echo "📊 Summary:"
echo "   - Project ID: $PROJECT_ID"
echo "   - Backend Agent ID: $BACKEND_AGENT_ID"
echo "   - Frontend Agent ID: $FRONTEND_AGENT_ID"
echo "   - Tasks Created: 2"
echo "   - Documentation Added: 1"
echo ""
echo "🌐 Open Web UI: http://localhost:8080"
echo "📚 API Documentation: docs/API.md"
echo "🤖 Copilot Integration: docs/COPILOT_INTEGRATION.md"
