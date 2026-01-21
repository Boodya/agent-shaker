<template>
  <div>
    <div class="mb-8">
      <h2 class="text-3xl font-bold text-gray-900">All Agents</h2>
    </div>

    <div v-if="loading" class="text-center py-12 text-gray-500">Loading agents...</div>
    <div v-else-if="error" class="p-4 bg-red-50 text-red-600 rounded-md mb-4">{{ error }}</div>
    
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div v-for="agent in agents" :key="agent.id" class="bg-white p-6 rounded-lg shadow-sm">
        <div class="flex justify-between items-start mb-4">
          <h3 class="text-xl font-semibold text-gray-900">{{ agent.name }}</h3>
          <span :class="['px-3 py-1 rounded-full text-sm font-semibold', agent.status === 'active' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800' ]">{{ agent.status }}</span>
        </div>
        <div class="space-y-2">
          <p class="text-gray-600"><strong class="font-medium text-gray-900">Role:</strong> <span :class="['inline-block px-2 py-1 rounded text-xs font-semibold', agent.role === 'frontend' ? 'bg-blue-100 text-blue-800' : 'bg-pink-100 text-pink-800' ]">{{ agent.role }}</span></p>
          <p class="text-gray-600"><strong class="font-medium text-gray-900">Team:</strong> {{ agent.team }}</p>
          <p class="text-gray-600"><strong class="font-medium text-gray-900">Project ID:</strong> {{ agent.project_id }}</p>
          <p class="text-gray-600"><strong class="font-medium text-gray-900">Last Seen:</strong> {{ formatDate(agent.last_seen) }}</p>
          <p class="text-gray-600"><strong class="font-medium text-gray-900">Created:</strong> {{ formatDate(agent.created_at) }}</p>
        </div>
      </div>
    </div>

    <div v-if="!loading && agents.length === 0" class="text-center py-12">
      <h3 class="text-xl font-semibold text-gray-900 mb-2">No agents registered</h3>
      <p class="text-gray-600">Agents will appear here when they are registered to projects</p>
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
