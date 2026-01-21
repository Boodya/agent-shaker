<template>
  <div class="dashboard">
    <div class="container">
      <div class="dashboard-header">
        <h2>Dashboard</h2>
        <p class="subtitle">Overview of your MCP Task Tracker</p>
      </div>

      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-icon">📁</div>
          <div class="stat-content">
            <h3>{{ projects.length }}</h3>
            <p>Projects</p>
          </div>
        </div>

        <div class="stat-card">
          <div class="stat-icon">🤖</div>
          <div class="stat-content">
            <h3>{{ agents.length }}</h3>
            <p>Agents</p>
          </div>
        </div>

        <div class="stat-card">
          <div class="stat-icon">📋</div>
          <div class="stat-content">
            <h3>{{ tasks.length }}</h3>
            <p>Total Tasks</p>
          </div>
        </div>

        <div class="stat-card">
          <div class="stat-icon">✅</div>
          <div class="stat-content">
            <h3>{{ completedTasks }}</h3>
            <p>Completed</p>
          </div>
        </div>
      </div>

      <div class="dashboard-grid">
        <div class="card">
          <h3>Recent Projects</h3>
          <div class="project-list">
            <div v-for="project in recentProjects" :key="project.id" class="project-item">
              <router-link :to="`/projects/${project.id}`" class="project-link">
                <h4>{{ project.name }}</h4>
                <p>{{ project.description }}</p>
                <span class="badge" :class="project.status">{{ project.status }}</span>
              </router-link>
            </div>
            <p v-if="projects.length === 0" class="empty-state">No projects yet</p>
          </div>
        </div>

        <div class="card">
          <h3>Active Agents</h3>
          <div class="agent-list">
            <div v-for="agent in activeAgents" :key="agent.id" class="agent-item">
              <div class="agent-info">
                <span class="agent-name">{{ agent.name }}</span>
                <span class="agent-role" :class="agent.role">{{ agent.role }}</span>
              </div>
              <span class="status-badge" :class="agent.status">{{ agent.status }}</span>
            </div>
            <p v-if="activeAgents.length === 0" class="empty-state">No active agents</p>
          </div>
        </div>

        <div class="card full-width">
          <h3>Recent Tasks</h3>
          <div class="task-list">
            <div v-for="task in recentTasks" :key="task.id" class="task-item">
              <div class="task-content">
                <h4>{{ task.title }}</h4>
                <p>{{ task.description }}</p>
              </div>
              <div class="task-meta">
                <span class="priority" :class="task.priority">{{ task.priority }}</span>
                <span class="status" :class="task.status">{{ task.status }}</span>
              </div>
            </div>
            <p v-if="tasks.length === 0" class="empty-state">No tasks yet</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { computed, onMounted } from 'vue'
import { useProjectStore } from '../stores/projectStore'
import { useAgentStore } from '../stores/agentStore'
import { useTaskStore } from '../stores/taskStore'

export default {
  name: 'Dashboard',
  setup() {
    const projectStore = useProjectStore()
    const agentStore = useAgentStore()
    const taskStore = useTaskStore()

    onMounted(() => {
      projectStore.fetchProjects()
      agentStore.fetchAgents()
      taskStore.fetchTasks()
    })

    const recentProjects = computed(() => projectStore.projects.slice(0, 5))
    const activeAgents = computed(() => 
      agentStore.agents.filter(a => a.status === 'active').slice(0, 5)
    )
    const recentTasks = computed(() => taskStore.tasks.slice(0, 10))
    const completedTasks = computed(() => 
      taskStore.tasks.filter(t => t.status === 'done').length
    )

    return {
      projects: projectStore.projects,
      agents: agentStore.agents,
      tasks: taskStore.tasks,
      recentProjects,
      activeAgents,
      recentTasks,
      completedTasks
    }
  }
}
</script>
