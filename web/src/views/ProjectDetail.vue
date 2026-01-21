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
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useProjectStore } from '../stores/projectStore'
import { useAgentStore } from '../stores/agentStore'
import { useTaskStore } from '../stores/taskStore'

export default {
  name: 'ProjectDetail',
  setup() {
    const route = useRoute()
    const projectStore = useProjectStore()
    const agentStore = useAgentStore()
    const taskStore = useTaskStore()

    const activeTab = ref('agents')
    const showAddAgentModal = ref(false)
    const showAddTaskModal = ref(false)

    const newAgent = ref({ name: '', role: 'frontend', team: '' })
    const newTask = ref({ title: '', description: '', agent_id: '', priority: 'medium' })

    const project = computed(() => projectStore.currentProject)
    const agents = computed(() => agentStore.agents)
    const tasks = computed(() => taskStore.tasks)

    onMounted(() => {
      const projectId = route.params.id
      projectStore.fetchProject(projectId)
      agentStore.fetchProjectAgents(projectId)
      taskStore.fetchProjectTasks(projectId)
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

    const getAgentName = (agentId) => {
      const agent = agents.value.find(a => a.id === agentId)
      return agent ? agent.name : 'Unknown'
    }

    const formatDate = (dateString) => {
      return new Date(dateString).toLocaleString()
    }

    return {
      project,
      agents,
      tasks,
      loading: projectStore.loading,
      error: projectStore.error,
      activeTab,
      showAddAgentModal,
      showAddTaskModal,
      newAgent,
      newTask,
      handleAddAgent,
      handleAddTask,
      getAgentName,
      formatDate
    }
  }
}
</script>
