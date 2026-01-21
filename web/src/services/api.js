import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor
api.interceptors.request.use(
  (config) => {
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor
api.interceptors.response.use(
  (response) => {
    return response.data
  },
  (error) => {
    console.error('API Error:', error)
    return Promise.reject(error)
  }
)

export default {
  // Projects
  getProjects() {
    return api.get('/projects')
  },
  getProject(id) {
    return api.get(`/projects/${id}`)
  },
  createProject(data) {
    return api.post('/projects', data)
  },

  // Agents
  getAgents() {
    return api.get('/agents')
  },
  getAgent(id) {
    return api.get(`/agents/${id}`)
  },
  getProjectAgents(projectId) {
    return api.get(`/projects/${projectId}/agents`)
  },
  createAgent(data) {
    return api.post('/agents', data)
  },
  updateAgentStatus(id, status) {
    return api.put(`/agents/${id}/status`, { status })
  },

  // Tasks
  getTasks() {
    return api.get('/tasks')
  },
  getTask(id) {
    return api.get(`/tasks/${id}`)
  },
  getProjectTasks(projectId) {
    return api.get(`/projects/${projectId}/tasks`)
  },
  getAgentTasks(agentId) {
    return api.get(`/agents/${agentId}/tasks`)
  },
  createTask(data) {
    return api.post('/tasks', data)
  },
  updateTaskStatus(id, status) {
    return api.put(`/tasks/${id}/status`, { status })
  },

  // Context/Documentation
  createDocumentation(data) {
    return api.post('/documentation', data)
  },
  getTaskDocumentation(taskId) {
    return api.get(`/tasks/${taskId}/documentation`)
  },

  // Contexts
  getContexts() {
    return api.get('/contexts')
  },
  getContext(id) {
    return api.get(`/contexts/${id}`)
  },
  getProjectContexts(projectId) {
    return api.get('/contexts', { params: { project_id: projectId } })
  },
  createContext(data) {
    return api.post('/contexts', data)
  },
  updateContext(id, data) {
    return api.put(`/contexts/${id}`, data)
  },
  deleteContext(id) {
    return api.delete(`/contexts/${id}`)
  },

  // Health
  checkHealth() {
    return axios.get('/health')
  },
}
