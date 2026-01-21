# Quick Start - MCP Agent Setup

## TL;DR - Get Started in 2 Minutes

### 1. Start the MCP Server

```powershell
docker-compose up -d
```

### 2. Run the Setup Script (Interactive Mode)

```powershell
.\setup-agent.ps1 -Interactive
```

Follow the prompts to:
- Enter your project name (e.g., "InvoiceAI")
- Enter your agent name (e.g., "InvoiceAI-Frontend")
- Select your role (Frontend or Backend)

### 3. Load Agent Helper Functions

```powershell
# The script will tell you the exact command, something like:
. .\agent-helper-InvoiceAI-Frontend.ps1
```

### 4. Start Working!

```powershell
# View your tasks
Get-MyTasks

# Create a new task
New-Task -Title "Build invoice list component" -Priority "high"

# Update task status
Update-TaskStatus -TaskId "task-id-here" -Status "in_progress"

# Show your agent info
Show-AgentInfo
```

---

## Example: Setting up InvoiceAI Frontend Agent

```powershell
# 1. Start services
docker-compose up -d

# 2. Setup agent (non-interactive)
.\setup-agent.ps1 -ProjectName "InvoiceAI" -AgentName "InvoiceAI-Frontend" -Role "frontend"

# 3. Load helper
. .\agent-helper-InvoiceAI-Frontend.ps1

# 4. Create your first task
New-Task -Title "Create Invoice List Component" -Description "Build React component for invoice listing" -Priority "high"

# 5. Check your tasks
Get-MyTasks
```

---

## Full Documentation

📚 **[Complete Agent Setup Guide](./AGENT_SETUP_GUIDE.md)** - Comprehensive manual with all API details, examples, and troubleshooting

---

## What Gets Created?

After running the setup script, you'll have:

1. **agent-config.json** - Your agent's identity configuration
2. **agent-helper-[YourAgentName].ps1** - PowerShell functions to interact with MCP
3. **.env.agent** - Environment variables for your agent

---

## Common Commands

```powershell
# View all your tasks
Get-MyTasks

# Create a task
New-Task -Title "Task name" -Description "Details" -Priority "high"

# Update task status (pending, in_progress, done, blocked)
Update-TaskStatus -TaskId "xxx" -Status "in_progress"

# Add documentation/context to a task
Add-TaskContext -TaskId "xxx" -Context "Implementation notes here..."

# Show agent information
Show-AgentInfo
```

---

## Multiple Agents Setup

```powershell
# Setup frontend agent
.\setup-agent.ps1 -ProjectName "InvoiceAI" -AgentName "InvoiceAI-Frontend" -Role "frontend"

# Setup backend agent
.\setup-agent.ps1 -ProjectName "InvoiceAI" -AgentName "InvoiceAI-Backend" -Role "backend"

# Now you have two agents collaborating on the same project!
```

---

## Troubleshooting

**Server not running?**
```powershell
docker-compose up -d
```

**Want to see server logs?**
```powershell
docker-compose logs -f mcp-server
```

**Reset everything?**
```powershell
docker-compose down -v
docker-compose up -d
```

---

## API Endpoints Quick Reference

- **Health**: http://localhost:8080/health
- **Projects**: http://localhost:8080/api/projects
- **Agents**: http://localhost:8080/api/agents
- **Tasks**: http://localhost:8080/api/tasks

---

## Support

- 📖 Full Guide: [AGENT_SETUP_GUIDE.md](./AGENT_SETUP_GUIDE.md)
- 🔧 API Docs: [docs/API.md](./docs/API.md)
- 🌐 WebSocket Demo: http://localhost:8080

---

**Happy Multi-Agent Development! 🚀**
