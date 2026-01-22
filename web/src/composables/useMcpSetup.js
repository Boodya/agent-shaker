import { computed } from 'vue'

/**
 * Composable for MCP setup configuration generation
 * @param {Object} agent - Agent object
 * @param {Object} project - Project object
 * @param {string} apiUrl - API base URL
 * @returns {Object} MCP configuration objects
 */
export const useMcpSetup = (agent, project, apiUrl) => {
  const mcpSettingsJson = computed(() => {
    if (!agent.value || !project.value) return ''
    return JSON.stringify({
      "terminal.integrated.env.windows": {
        "MCP_AGENT_NAME": agent.value.name,
        "MCP_AGENT_ID": agent.value.id,
        "MCP_PROJECT_ID": project.value.id,
        "MCP_PROJECT_NAME": project.value.name,
        "MCP_API_URL": apiUrl
      },
      "terminal.integrated.env.linux": {
        "MCP_AGENT_NAME": agent.value.name,
        "MCP_AGENT_ID": agent.value.id,
        "MCP_PROJECT_ID": project.value.id,
        "MCP_PROJECT_NAME": project.value.name,
        "MCP_API_URL": apiUrl
      },
      "terminal.integrated.env.osx": {
        "MCP_AGENT_NAME": agent.value.name,
        "MCP_AGENT_ID": agent.value.id,
        "MCP_PROJECT_ID": project.value.id,
        "MCP_PROJECT_NAME": project.value.name,
        "MCP_API_URL": apiUrl
      }
    }, null, 2)
  })

  const mcpCopilotInstructions = computed(() => {
    if (!agent.value || !project.value) return ''
    return `# Agent Identity and MCP Integration

## Your Identity
- **Agent Name**: ${agent.value.name}
- **Agent ID**: ${agent.value.id}
- **Role**: ${agent.value.role}
- **Team**: ${agent.value.team || 'Not specified'}
- **Project**: ${project.value.name}
- **Project ID**: ${project.value.id}

## MCP API Configuration
- **API URL**: ${apiUrl}

## Your Responsibilities
As the **${agent.value.role}** agent, you should:
${agent.value.role === 'frontend' ? `
- Focus on UI/UX implementation
- Work with Vue.js, React, or other frontend frameworks
- Implement responsive designs and accessibility
- Handle client-side state management
` : `
- Focus on API development and backend logic
- Work with databases and data models
- Implement business logic and validations
- Handle server-side security and authentication
`}

## Task Management
When working on tasks, use these API endpoints:

### Get Your Tasks
\`\`\`bash
curl "${apiUrl}/agents/${agent.value.id}/tasks"
\`\`\`

### Update Task Status
\`\`\`bash
curl -X PUT "${apiUrl}/tasks/{task_id}/status" \\
  -H "Content-Type: application/json" \\
  -d '{"status": "in_progress"}'
\`\`\`

Status options: \`pending\`, \`in_progress\`, \`done\`, \`blocked\`

### Add Context/Documentation
\`\`\`bash
curl -X POST "${apiUrl}/contexts" \\
  -H "Content-Type: application/json" \\
  -d '{
    "project_id": "${project.value.id}",
    "agent_id": "${agent.value.id}",
    "title": "Implementation Notes",
    "content": "Your documentation here...",
    "tags": ["documentation", "${agent.value.role}"]
  }'
\`\`\`

## Collaboration Guidelines
1. Always check for existing tasks before starting new work
2. Update task status when you begin and complete work
3. Document important decisions and implementation details
4. Check other agents' contexts to avoid conflicts
`
  })

  const mcpPowerShellScript = computed(() => {
    if (!agent.value || !project.value) return ''
    return `# MCP Agent Helper Script for PowerShell
# Agent: ${agent.value.name}
# Project: ${project.value.name}

$MCP_API_URL = "${apiUrl}"
$MCP_AGENT_ID = "${agent.value.id}"
$MCP_PROJECT_ID = "${project.value.id}"

function Get-MyTasks {
    Invoke-RestMethod -Uri "$MCP_API_URL/agents/$MCP_AGENT_ID/tasks" -Method GET
}

function Update-TaskStatus {
    param(
        [Parameter(Mandatory=$true)]
        [string]$TaskId,
        [Parameter(Mandatory=$true)]
        [ValidateSet("pending", "in_progress", "done", "blocked")]
        [string]$Status
    )
    
    $body = @{ status = $Status } | ConvertTo-Json
    Invoke-RestMethod -Uri "$MCP_API_URL/tasks/$TaskId/status" -Method PUT -Body $body -ContentType "application/json"
}

function Add-Context {
    param(
        [Parameter(Mandatory=$true)]
        [string]$Title,
        [Parameter(Mandatory=$true)]
        [string]$Content,
        [string[]]$Tags = @()
    )
    
    $body = @{
        project_id = $MCP_PROJECT_ID
        agent_id = $MCP_AGENT_ID
        title = $Title
        content = $Content
        tags = $Tags
    } | ConvertTo-Json
    
    Invoke-RestMethod -Uri "$MCP_API_URL/contexts" -Method POST -Body $body -ContentType "application/json"
}

function Get-ProjectContexts {
    Invoke-RestMethod -Uri "$MCP_API_URL/projects/$MCP_PROJECT_ID/contexts" -Method GET
}

# Usage examples:
# Get-MyTasks
# Update-TaskStatus -TaskId "task-uuid" -Status "in_progress"
# Add-Context -Title "API Design" -Content "Documentation content..." -Tags @("api", "design")
`
  })

  const mcpBashScript = computed(() => {
    if (!agent.value || !project.value) return ''
    return `#!/bin/bash
# MCP Agent Helper Script for Bash
# Agent: ${agent.value.name}
# Project: ${project.value.name}

MCP_API_URL="${apiUrl}"
MCP_AGENT_ID="${agent.value.id}"
MCP_PROJECT_ID="${project.value.id}"

# Get tasks assigned to this agent
get_my_tasks() {
    curl -s "$MCP_API_URL/agents/$MCP_AGENT_ID/tasks" | jq .
}

# Update task status
# Usage: update_task_status <task_id> <status>
# Status: pending, in_progress, done, blocked
update_task_status() {
    local task_id=$1
    local status=$2
    curl -s -X PUT "$MCP_API_URL/tasks/$task_id/status" \\
        -H "Content-Type: application/json" \\
        -d "{\\"status\\": \\"$status\\"}" | jq .
}

# Add context/documentation
# Usage: add_context "Title" "Content" "tag1,tag2"
add_context() {
    local title=$1
    local content=$2
    local tags=$3
    
    curl -s -X POST "$MCP_API_URL/contexts" \\
        -H "Content-Type: application/json" \\
        -d "{
            \\"project_id\\": \\"$MCP_PROJECT_ID\\",
            \\"agent_id\\": \\"$MCP_AGENT_ID\\",
            \\"title\\": \\"$title\\",
            \\"content\\": \\"$content\\",
            \\"tags\\": [\\"$tags\\"]
        }" | jq .
}

# Get project contexts
get_project_contexts() {
    curl -s "$MCP_API_URL/projects/$MCP_PROJECT_ID/contexts" | jq .
}

# Usage examples:
# get_my_tasks
# update_task_status "task-uuid" "in_progress"
# add_context "API Design" "Documentation content..." "api,design"
`
  })

  const mcpVSCodeJson = computed(() => {
    if (!agent.value || !project.value) return ''
    
    const baseUrl = apiUrl.replace('/api', '')
    const mcpUrl = `${baseUrl}?project_id=${project.value.id}&agent_id=${agent.value.id}`
    
    const config = {
      "mcpServers": {
        "agent-shaker": {
          "url": mcpUrl,
          "type": "http",
          "metadata": {
            "name": "Agent Shaker MCP Server",
            "version": "1.0.0",
            "description": "Multi-agent coordination platform for collaborative development",
            "capabilities": [
              "resources",
              "tools",
              "prompts",
              "context-sharing"
            ]
          },
          "project": {
            "id": project.value.id,
            "name": project.value.name,
            "description": project.value.description || "",
            "status": project.value.status,
            "type": "multi-agent"
          },
          "agent": {
            "id": agent.value.id,
            "name": agent.value.name,
            "role": agent.value.role,
            "team": agent.value.team || "default",
            "status": agent.value.status
          }
        }
      }
    }
    
    return JSON.stringify(config, null, 2)
  })

  return {
    mcpSettingsJson,
    mcpCopilotInstructions,
    mcpPowerShellScript,
    mcpBashScript,
    mcpVSCodeJson
  }
}

/**
 * Download file helper
 * @param {string} filename - Name of file to download
 * @param {string} content - File content
 * @param {string} mimeType - MIME type
 */
export const downloadFile = (filename, content, mimeType = 'text/plain') => {
  const blob = new Blob([content], { type: mimeType })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = filename
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}

/**
 * Download all MCP files as zip
 * @param {Object} mcpConfig - MCP configuration object
 * @param {string} agentName - Agent name for filename
 */
export const downloadAllMcpFiles = async (mcpConfig, agentName) => {
  const { default: JSZip } = await import('jszip')
  const zip = new JSZip()
  
  zip.file('.vscode/settings.json', mcpConfig.mcpSettingsJson)
  zip.file('.vscode/mcp.json', mcpConfig.mcpVSCodeJson)
  zip.file('.github/copilot-instructions.md', mcpConfig.mcpCopilotInstructions)
  zip.file('scripts/mcp-agent.ps1', mcpConfig.mcpPowerShellScript)
  zip.file('scripts/mcp-agent.sh', mcpConfig.mcpBashScript)
  
  const readmeContent = `# MCP Setup Files for ${agentName}

## Contents
- .vscode/settings.json - VS Code environment variables
- .vscode/mcp.json - Enhanced MCP server configuration
- .github/copilot-instructions.md - GitHub Copilot agent instructions
- scripts/mcp-agent.ps1 - PowerShell helper script
- scripts/mcp-agent.sh - Bash helper script

## Setup Instructions
1. Extract this zip to your project's root directory
2. Restart VS Code to apply environment variables
3. Start using Copilot with your agent identity!
`
  zip.file('MCP_SETUP_README.md', readmeContent)
  
  const content = await zip.generateAsync({ type: 'blob' })
  const agentSlug = agentName.toLowerCase().replace(/[^a-z0-9]+/g, '-')
  downloadFile(`mcp-setup-${agentSlug}.zip`, content, 'application/zip')
}
