/**
 * Get agent name by ID
 * @param {Array} agents - Array of agent objects
 * @param {string} agentId - Agent ID to find
 * @returns {string} Agent name or 'Unknown'
 */
export const getAgentName = (agents, agentId) => {
  const agent = agents.find(a => a.id === agentId)
  return agent ? agent.name : 'Unknown'
}

/**
 * Get task title by ID
 * @param {Array} tasks - Array of task objects
 * @param {string} taskId - Task ID to find
 * @returns {string} Task title or 'Unknown'
 */
export const getTaskTitle = (tasks, taskId) => {
  const task = tasks.find(t => t.id === taskId)
  return task ? task.title : 'Unknown'
}

/**
 * Filter contexts by search and tag
 * @param {Array} contexts - Array of context objects
 * @param {string} searchQuery - Search query
 * @param {string} tagFilter - Tag to filter by
 * @returns {Array} Filtered contexts
 */
export const filterContexts = (contexts, searchQuery, tagFilter) => {
  let filtered = [...contexts]

  // Filter by search
  if (searchQuery) {
    const search = searchQuery.toLowerCase()
    filtered = filtered.filter(context =>
      context.title.toLowerCase().includes(search) ||
      context.content.toLowerCase().includes(search)
    )
  }

  // Filter by tag
  if (tagFilter) {
    filtered = filtered.filter(context =>
      context.tags && context.tags.includes(tagFilter)
    )
  }

  return filtered
}
