<template>
  <div class="project-detail">
    <div class="container">
      <div v-if="loading" class="loading">Loading project...</div>
      <div v-else-if="error" class="error">{{ error }}</div>
      
      <div v-else-if="project">
        <div class="page-header">
          <div>
            <h2>{{ project.name }}</h2>
            <p class="subtitle">{{ project.description }}</p>
          </div>
          <span class="badge" :class="project.status">{{ project.status }}</span>
        </div>

        <div class="tabs">
          <button 
            :class="{ active: activeTab === 'agents' }"
            @click="activeTab = 'agents'"
          >
            Agents ({{ agents.length }})
          </button>
          <button 
            :class="{ active: activeTab === 'tasks' }"
            @click="activeTab = 'tasks'"
          >
            Tasks ({{ tasks.length }})
          </button>
          <button 
            :class="{ active: activeTab === 'contexts' }"
            @click="activeTab = 'contexts'"
          >
            Contexts ({{ contexts.length }})
          </button>
        </div>

        <!-- Agents Tab -->
        <div v-if="activeTab === 'agents'" class="tab-content">
          <div class="section-header">
            <h3>Project Agents</h3>
            <button @click="showAddAgentModal = true" class="btn btn-primary">
              + Add Agent
            </button>
          </div>

          <div class="agents-grid">
            <div v-for="agent in agents" :key="agent.id" class="agent-card">
              <div class="agent-header">
                <h4>{{ agent.name }}</h4>
                <span class="status-badge" :class="agent.status">{{ agent.status }}</span>
              </div>
              <div class="agent-details">
                <p><strong>Role:</strong> <span class="badge" :class="agent.role">{{ agent.role }}</span></p>
                <p><strong>Team:</strong> {{ agent.team }}</p>
                <p><strong>Last Seen:</strong> {{ formatDate(agent.last_seen) }}</p>
              </div>
            </div>
          </div>

          <div v-if="agents.length === 0" class="empty-state">
            <p>No agents assigned to this project yet</p>
          </div>
        </div>

        <!-- Tasks Tab -->
        <div v-if="activeTab === 'tasks'" class="tab-content">
          <div class="section-header">
            <h3>Project Tasks</h3>
            <button @click="showAddTaskModal = true" class="btn btn-primary">
              + Create Task
            </button>
          </div>

          <div class="tasks-list">
            <div v-for="task in tasks" :key="task.id" class="task-card">
              <div class="task-header">
                <h4>{{ task.title }}</h4>
                <div class="task-badges">
                  <span class="priority" :class="task.priority">{{ task.priority }}</span>
                  <span class="status" :class="task.status">{{ task.status }}</span>
                </div>
              </div>
              <p>{{ task.description }}</p>
              <div class="task-footer">
                <span>Agent: {{ getAgentName(task.agent_id) }}</span>
                <span>Created {{ formatDate(task.created_at) }}</span>
              </div>
            </div>
          </div>

          <div v-if="tasks.length === 0" class="empty-state">
            <p>No tasks in this project yet</p>
          </div>
        </div>

        <!-- Contexts Tab -->
        <div v-if="activeTab === 'contexts'" class="tab-content">
          <div class="section-header">
            <h3>Project Documentation / Context</h3>
            <button @click="showAddContextModal = true" class="btn btn-primary">
              + Add Context
            </button>
          </div>

          <div class="search-filter">
            <input 
              v-model="contextSearch" 
              type="text" 
              placeholder="Search contexts..." 
              class="search-input"
            />
            <select v-model="contextTagFilter" class="filter-select">
              <option value="">All Tags</option>
              <option v-for="tag in uniqueTags" :key="tag" :value="tag">
                {{ tag }}
              </option>
            </select>
          </div>

          <div class="contexts-grid">
            <div v-for="context in filteredContexts" :key="context.id" class="context-card">
              <div class="context-header">
                <h4>{{ context.title }}</h4>
                <div class="context-actions">
                  <button @click="viewContext(context)" class="btn-icon" title="View">
                    👁️
                  </button>
                  <button @click="editContext(context)" class="btn-icon" title="Edit">
                    ✏️
                  </button>
                  <button @click="confirmDeleteContext(context)" class="btn-icon btn-danger" title="Delete">
                    🗑️
                  </button>
                </div>
              </div>
              <div class="context-tags">
                <span v-for="tag in context.tags" :key="tag" class="tag">{{ tag }}</span>
              </div>
              <div class="context-meta">
                <span>Agent: {{ getAgentName(context.agent_id) }}</span>
                <span v-if="context.task_id">Task: {{ getTaskTitle(context.task_id) }}</span>
                <span>{{ formatDate(context.created_at) }}</span>
              </div>
            </div>
          </div>

          <div v-if="filteredContexts.length === 0" class="empty-state">
            <p>{{ contexts.length === 0 ? 'No contexts in this project yet' : 'No contexts match your search' }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Add Agent Modal -->
    <div v-if="showAddAgentModal" class="modal-overlay" @click="showAddAgentModal = false">
      <div class="modal" @click.stop>
        <h3>Add Agent to Project</h3>
        <form @submit.prevent="handleAddAgent">
          <div class="form-group">
            <label>Agent Name</label>
            <input v-model="newAgent.name" type="text" required />
          </div>
          <div class="form-group">
            <label>Role</label>
            <select v-model="newAgent.role" required>
              <option value="frontend">Frontend</option>
              <option value="backend">Backend</option>
            </select>
          </div>
          <div class="form-group">
            <label>Team</label>
            <input v-model="newAgent.team" type="text" />
          </div>
          <div class="modal-actions">
            <button type="button" @click="showAddAgentModal = false" class="btn btn-secondary">
              Cancel
            </button>
            <button type="submit" class="btn btn-primary">Add Agent</button>
          </div>
        </form>
      </div>
    </div>

    <!-- Add Task Modal -->
    <div v-if="showAddTaskModal" class="modal-overlay" @click="showAddTaskModal = false">
      <div class="modal" @click.stop>
        <h3>Create New Task</h3>
        <form @submit.prevent="handleAddTask">
          <div class="form-group">
            <label>Task Title</label>
            <input v-model="newTask.title" type="text" required />
          </div>
          <div class="form-group">
            <label>Description</label>
            <textarea v-model="newTask.description" rows="4"></textarea>
          </div>
          <div class="form-group">
            <label>Assign to Agent</label>
            <select v-model="newTask.agent_id" required>
              <option value="">Select an agent</option>
              <option v-for="agent in agents" :key="agent.id" :value="agent.id">
                {{ agent.name }}
              </option>
            </select>
          </div>
          <div class="form-group">
            <label>Priority</label>
            <select v-model="newTask.priority">
              <option value="low">Low</option>
              <option value="medium">Medium</option>
              <option value="high">High</option>
            </select>
          </div>
          <div class="modal-actions">
            <button type="button" @click="showAddTaskModal = false" class="btn btn-secondary">
              Cancel
            </button>
            <button type="submit" class="btn btn-primary">Create Task</button>
          </div>
        </form>
      </div>
    </div>

    <!-- Add Context Modal -->
    <div v-if="showAddContextModal" class="modal-overlay" @click="closeContextModal">
      <div class="modal modal-large" @click.stop>
        <h3>{{ editingContext ? 'Edit Context' : 'Add Context / Documentation' }}</h3>
        <form @submit.prevent="handleSaveContext">
          <div class="form-group">
            <label>Title</label>
            <input v-model="contextForm.title" type="text" required />
          </div>
          <div class="form-group">
            <label>Agent</label>
            <select v-model="contextForm.agent_id" required>
              <option value="">Select an agent</option>
              <option v-for="agent in agents" :key="agent.id" :value="agent.id">
                {{ agent.name }}
              </option>
            </select>
          </div>
          <div class="form-group">
            <label>Related Task (Optional)</label>
            <select v-model="contextForm.task_id">
              <option value="">None</option>
              <option v-for="task in tasks" :key="task.id" :value="task.id">
                {{ task.title }}
              </option>
            </select>
          </div>
          <div class="form-group">
            <label>Content (Markdown)</label>
            <textarea 
              v-model="contextForm.content" 
              rows="12" 
              placeholder="Write your documentation in Markdown format..."
              required
            ></textarea>
            <small class="help-text">Supports Markdown: **bold**, *italic*, [link](url), etc.</small>
          </div>
          <div class="form-group">
            <label>Tags (comma-separated)</label>
            <input 
              v-model="contextForm.tagsString" 
              type="text" 
              placeholder="api, documentation, backend"
            />
          </div>
          <div class="modal-actions">
            <button type="button" @click="closeContextModal" class="btn btn-secondary">
              Cancel
            </button>
            <button type="submit" class="btn btn-primary">
              {{ editingContext ? 'Update' : 'Create' }} Context
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- View Context Modal -->
    <div v-if="showViewContextModal" class="modal-overlay" @click="showViewContextModal = false">
      <div class="modal modal-large" @click.stop>
        <div class="modal-header-flex">
          <h3>{{ viewingContext?.title }}</h3>
          <button @click="showViewContextModal = false" class="btn-close">×</button>
        </div>
        <div class="context-view">
          <div class="context-meta-bar">
            <div class="context-tags">
              <span v-for="tag in viewingContext?.tags" :key="tag" class="tag">{{ tag }}</span>
            </div>
            <div class="context-info">
              <span>Agent: {{ getAgentName(viewingContext?.agent_id) }}</span>
              <span v-if="viewingContext?.task_id">Task: {{ getTaskTitle(viewingContext?.task_id) }}</span>
              <span>{{ formatDate(viewingContext?.created_at) }}</span>
            </div>
          </div>
          <div class="markdown-content" v-html="renderMarkdown(viewingContext?.content)"></div>
        </div>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <div v-if="showDeleteConfirm" class="modal-overlay" @click="showDeleteConfirm = false">
      <div class="modal modal-small" @click.stop>
        <h3>⚠️ Confirm Delete</h3>
        <p>Are you sure you want to delete the context "{{ deletingContext?.title }}"?</p>
        <div class="modal-actions">
          <button @click="showDeleteConfirm = false" class="btn btn-secondary">
            Cancel
          </button>
          <button @click="handleDeleteContext" class="btn btn-danger">
            Delete
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useProjectStore } from '../stores/projectStore'
import { useAgentStore } from '../stores/agentStore'
import { useTaskStore } from '../stores/taskStore'
import { useContextStore } from '../stores/contextStore'
import { marked } from 'marked'
import DOMPurify from 'dompurify'

export default {
  name: 'ProjectDetail',
  setup() {
    const route = useRoute()
    const projectStore = useProjectStore()
    const agentStore = useAgentStore()
    const taskStore = useTaskStore()
    const contextStore = useContextStore()

    const activeTab = ref('agents')
    const showAddAgentModal = ref(false)
    const showAddTaskModal = ref(false)
    const showAddContextModal = ref(false)
    const showViewContextModal = ref(false)
    const showDeleteConfirm = ref(false)

    const newAgent = ref({ name: '', role: 'frontend', team: '' })
    const newTask = ref({ title: '', description: '', agent_id: '', priority: 'medium' })
    
    const contextForm = ref({
      title: '',
      agent_id: '',
      task_id: '',
      content: '',
      tagsString: ''
    })
    const editingContext = ref(null)
    const viewingContext = ref(null)
    const deletingContext = ref(null)
    
    const contextSearch = ref('')
    const contextTagFilter = ref('')

    const project = computed(() => projectStore.currentProject)
    const agents = computed(() => agentStore.agents)
    const tasks = computed(() => taskStore.tasks)
    const contexts = computed(() => contextStore.contexts)

    const uniqueTags = computed(() => {
      const tags = new Set()
      contexts.value.forEach(context => {
        if (context.tags && Array.isArray(context.tags)) {
          context.tags.forEach(tag => tags.add(tag))
        }
      })
      return Array.from(tags).sort()
    })

    const filteredContexts = computed(() => {
      let filtered = contexts.value

      // Filter by search
      if (contextSearch.value) {
        const search = contextSearch.value.toLowerCase()
        filtered = filtered.filter(context =>
          context.title.toLowerCase().includes(search) ||
          context.content.toLowerCase().includes(search)
        )
      }

      // Filter by tag
      if (contextTagFilter.value) {
        filtered = filtered.filter(context =>
          context.tags && context.tags.includes(contextTagFilter.value)
        )
      }

      return filtered
    })

    onMounted(() => {
      const projectId = route.params.id
      projectStore.fetchProject(projectId)
      agentStore.fetchProjectAgents(projectId)
      taskStore.fetchProjectTasks(projectId)
      contextStore.fetchProjectContexts(projectId)
    })

    const handleAddAgent = async () => {
      try {
        await agentStore.createAgent({
          ...newAgent.value,
          project_id: route.params.id
        })
        showAddAgentModal.value = false
        newAgent.value = { name: '', role: 'frontend', team: '' }
      } catch (error) {
        console.error('Failed to add agent:', error)
      }
    }

    const handleAddTask = async () => {
      try {
        await taskStore.createTask({
          ...newTask.value,
          project_id: route.params.id,
          dependencies: []
        })
        showAddTaskModal.value = false
        newTask.value = { title: '', description: '', agent_id: '', priority: 'medium' }
      } catch (error) {
        console.error('Failed to create task:', error)
      }
    }

    const handleSaveContext = async () => {
      try {
        const tags = contextForm.value.tagsString
          .split(',')
          .map(tag => tag.trim())
          .filter(tag => tag.length > 0)

        const contextData = {
          project_id: route.params.id,
          agent_id: contextForm.value.agent_id,
          task_id: contextForm.value.task_id || null,
          title: contextForm.value.title,
          content: contextForm.value.content,
          tags: tags
        }

        if (editingContext.value) {
          await contextStore.updateContext(editingContext.value.id, contextData)
        } else {
          await contextStore.createContext(contextData)
        }

        closeContextModal()
      } catch (error) {
        console.error('Failed to save context:', error)
      }
    }

    const viewContext = (context) => {
      viewingContext.value = context
      showViewContextModal.value = true
    }

    const editContext = (context) => {
      editingContext.value = context
      contextForm.value = {
        title: context.title,
        agent_id: context.agent_id,
        task_id: context.task_id || '',
        content: context.content,
        tagsString: context.tags ? context.tags.join(', ') : ''
      }
      showAddContextModal.value = true
    }

    const confirmDeleteContext = (context) => {
      deletingContext.value = context
      showDeleteConfirm.value = true
    }

    const handleDeleteContext = async () => {
      try {
        await contextStore.deleteContext(deletingContext.value.id)
        showDeleteConfirm.value = false
        deletingContext.value = null
      } catch (error) {
        console.error('Failed to delete context:', error)
      }
    }

    const closeContextModal = () => {
      showAddContextModal.value = false
      editingContext.value = null
      contextForm.value = {
        title: '',
        agent_id: '',
        task_id: '',
        content: '',
        tagsString: ''
      }
    }

    const renderMarkdown = (content) => {
      if (!content) return ''
      const html = marked(content)
      return DOMPurify.sanitize(html)
    }

    const getAgentName = (agentId) => {
      const agent = agents.value.find(a => a.id === agentId)
      return agent ? agent.name : 'Unknown'
    }

    const getTaskTitle = (taskId) => {
      const task = tasks.value.find(t => t.id === taskId)
      return task ? task.title : 'Unknown'
    }

    const formatDate = (dateString) => {
      return new Date(dateString).toLocaleString()
    }

    return {
      project,
      agents,
      tasks,
      contexts,
      loading: projectStore.loading,
      error: projectStore.error,
      activeTab,
      showAddAgentModal,
      showAddTaskModal,
      showAddContextModal,
      showViewContextModal,
      showDeleteConfirm,
      newAgent,
      newTask,
      contextForm,
      editingContext,
      viewingContext,
      deletingContext,
      contextSearch,
      contextTagFilter,
      uniqueTags,
      filteredContexts,
      handleAddAgent,
      handleAddTask,
      handleSaveContext,
      viewContext,
      editContext,
      confirmDeleteContext,
      handleDeleteContext,
      closeContextModal,
      renderMarkdown,
      getAgentName,
      getTaskTitle,
      formatDate
    }
  }
}
</script>

<style scoped>
.project-detail {
  min-height: 100vh;
  background: #f5f7fa;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
}

.loading, .error {
  text-align: center;
  padding: 2rem;
  font-size: 1.1rem;
}

.error {
  color: #e74c3c;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 2rem;
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.page-header h2 {
  margin: 0 0 0.5rem 0;
  color: #2c3e50;
}

.subtitle {
  color: #7f8c8d;
  margin: 0;
}

.badge {
  padding: 0.5rem 1rem;
  border-radius: 20px;
  font-size: 0.85rem;
  font-weight: 600;
  text-transform: uppercase;
}

.badge.active { background: #d4edda; color: #155724; }
.badge.completed { background: #cce5ff; color: #004085; }
.badge.archived { background: #d6d8db; color: #383d41; }

.tabs {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;
  background: white;
  padding: 1rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.tabs button {
  padding: 0.75rem 1.5rem;
  border: none;
  background: transparent;
  color: #7f8c8d;
  font-weight: 500;
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.3s;
}

.tabs button:hover {
  background: #ecf0f1;
  color: #2c3e50;
}

.tabs button.active {
  background: #3498db;
  color: white;
}

.tab-content {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.section-header h3 {
  margin: 0;
  color: #2c3e50;
}

.btn {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.3s;
}

.btn-primary {
  background: #3498db;
  color: white;
}

.btn-primary:hover {
  background: #2980b9;
}

.btn-secondary {
  background: #95a5a6;
  color: white;
}

.btn-secondary:hover {
  background: #7f8c8d;
}

.btn-danger {
  background: #e74c3c;
  color: white;
}

.btn-danger:hover {
  background: #c0392b;
}

.btn-icon {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 1.2rem;
  padding: 0.25rem 0.5rem;
  opacity: 0.7;
  transition: opacity 0.3s;
}

.btn-icon:hover {
  opacity: 1;
}

.btn-close {
  background: none;
  border: none;
  font-size: 2rem;
  cursor: pointer;
  color: #95a5a6;
  line-height: 1;
  padding: 0;
}

.btn-close:hover {
  color: #7f8c8d;
}

.agents-grid, .contexts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1.5rem;
}

.agent-card, .context-card {
  background: #f8f9fa;
  padding: 1.5rem;
  border-radius: 8px;
  transition: transform 0.2s, box-shadow 0.2s;
}

.agent-card:hover, .context-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0,0,0,0.1);
}

.agent-header, .context-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1rem;
}

.agent-header h4, .context-header h4 {
  margin: 0;
  color: #2c3e50;
}

.context-actions {
  display: flex;
  gap: 0.25rem;
}

.status-badge {
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 600;
}

.status-badge.active { background: #d4edda; color: #155724; }
.status-badge.idle { background: #fff3cd; color: #856404; }
.status-badge.offline { background: #f8d7da; color: #721c24; }

.agent-details p {
  margin: 0.5rem 0;
  color: #7f8c8d;
}

.badge.frontend { background: #e3f2fd; color: #1976d2; }
.badge.backend { background: #f3e5f5; color: #7b1fa2; }
.badge.devops { background: #e8f5e9; color: #388e3c; }

.tasks-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.task-card {
  background: #f8f9fa;
  padding: 1.5rem;
  border-radius: 8px;
}

.task-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1rem;
}

.task-header h4 {
  margin: 0;
  color: #2c3e50;
}

.task-badges {
  display: flex;
  gap: 0.5rem;
}

.priority, .status {
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 600;
}

.priority.high { background: #f8d7da; color: #721c24; }
.priority.medium { background: #fff3cd; color: #856404; }
.priority.low { background: #d1e7dd; color: #0f5132; }

.status.pending { background: #cfe2ff; color: #084298; }
.status.in_progress { background: #fff3cd; color: #856404; }
.status.done { background: #d1e7dd; color: #0f5132; }
.status.blocked { background: #f8d7da; color: #721c24; }

.task-footer {
  display: flex;
  justify-content: space-between;
  color: #7f8c8d;
  font-size: 0.85rem;
  margin-top: 1rem;
  padding-top: 1rem;
  border-top: 1px solid #dee2e6;
}

.context-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.tag {
  background: #e3f2fd;
  color: #1976d2;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 500;
}

.context-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  color: #7f8c8d;
  font-size: 0.85rem;
}

.search-filter {
  display: flex;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.search-input, .filter-select {
  padding: 0.75rem 1rem;
  border: 1px solid #dee2e6;
  border-radius: 6px;
  font-size: 0.95rem;
}

.search-input {
  flex: 1;
}

.filter-select {
  min-width: 200px;
}

.empty-state {
  text-align: center;
  padding: 3rem;
  color: #7f8c8d;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  max-width: 500px;
  width: 90%;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-large {
  max-width: 800px;
}

.modal-small {
  max-width: 400px;
}

.modal h3 {
  margin: 0 0 1.5rem 0;
  color: #2c3e50;
}

.modal-header-flex {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.modal-header-flex h3 {
  margin: 0;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  color: #2c3e50;
  font-weight: 500;
}

.form-group input,
.form-group textarea,
.form-group select {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #dee2e6;
  border-radius: 6px;
  font-size: 0.95rem;
  font-family: inherit;
}

.form-group textarea {
  resize: vertical;
  font-family: 'Courier New', monospace;
}

.help-text {
  display: block;
  margin-top: 0.5rem;
  color: #7f8c8d;
  font-size: 0.85rem;
}

.modal-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  margin-top: 2rem;
}

.context-view {
  margin-top: 1rem;
}

.context-meta-bar {
  padding: 1rem;
  background: #f8f9fa;
  border-radius: 6px;
  margin-bottom: 1.5rem;
}

.context-info {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  margin-top: 0.75rem;
  color: #7f8c8d;
  font-size: 0.85rem;
}

.markdown-content {
  line-height: 1.7;
  color: #2c3e50;
}

.markdown-content h1,
.markdown-content h2,
.markdown-content h3 {
  margin-top: 1.5rem;
  margin-bottom: 0.75rem;
  color: #2c3e50;
}

.markdown-content h1 { font-size: 2rem; }
.markdown-content h2 { font-size: 1.5rem; }
.markdown-content h3 { font-size: 1.25rem; }

.markdown-content p {
  margin-bottom: 1rem;
}

.markdown-content ul,
.markdown-content ol {
  margin-bottom: 1rem;
  padding-left: 2rem;
}

.markdown-content code {
  background: #f8f9fa;
  padding: 0.2rem 0.4rem;
  border-radius: 3px;
  font-family: 'Courier New', monospace;
  font-size: 0.9em;
}

.markdown-content pre {
  background: #2c3e50;
  color: #ecf0f1;
  padding: 1rem;
  border-radius: 6px;
  overflow-x: auto;
  margin-bottom: 1rem;
}

.markdown-content pre code {
  background: none;
  padding: 0;
  color: inherit;
}

.markdown-content a {
  color: #3498db;
  text-decoration: none;
}

.markdown-content a:hover {
  text-decoration: underline;
}

.markdown-content blockquote {
  border-left: 4px solid #3498db;
  padding-left: 1rem;
  margin: 1rem 0;
  color: #7f8c8d;
}

.markdown-content table {
  width: 100%;
  border-collapse: collapse;
  margin-bottom: 1rem;
}

.markdown-content th,
.markdown-content td {
  border: 1px solid #dee2e6;
  padding: 0.75rem;
  text-align: left;
}

.markdown-content th {
  background: #f8f9fa;
  font-weight: 600;
}
</style>
