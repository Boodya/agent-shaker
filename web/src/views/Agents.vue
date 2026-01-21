<template>
  <div class="agents-page">
    <div class="container">
      <div class="page-header">
        <h2>All Agents</h2>
      </div>

      <div v-if="loading" class="loading">Loading agents...</div>
      <div v-else-if="error" class="error">{{ error }}</div>
      
      <div v-else class="agents-grid">
        <div v-for="agent in agents" :key="agent.id" class="agent-card">
          <div class="agent-header">
            <h3>{{ agent.name }}</h3>
            <span class="status-badge" :class="agent.status">{{ agent.status }}</span>
          </div>
          <div class="agent-details">
            <p><strong>Role:</strong> <span class="badge" :class="agent.role">{{ agent.role }}</span></p>
            <p><strong>Team:</strong> {{ agent.team }}</p>
            <p><strong>Project ID:</strong> {{ agent.project_id }}</p>
            <p><strong>Last Seen:</strong> {{ formatDate(agent.last_seen) }}</p>
            <p><strong>Created:</strong> {{ formatDate(agent.created_at) }}</p>
          </div>
        </div>
      </div>

      <div v-if="!loading && agents.length === 0" class="empty-state">
        <h3>No agents registered</h3>
        <p>Agents will appear here when they are registered to projects</p>
      </div>
    </div>
  </div>
</template>

<script>
import { onMounted } from 'vue'
import { useAgentStore } from '../stores/agentStore'

export default {
  name: 'Agents',
  setup() {
    const agentStore = useAgentStore()

    onMounted(() => {
      agentStore.fetchAgents()
    })

    const formatDate = (dateString) => {
      return new Date(dateString).toLocaleString()
    }

    return {
      agents: agentStore.agents,
      loading: agentStore.loading,
      error: agentStore.error,
      formatDate
    }
  }
}
</script>
