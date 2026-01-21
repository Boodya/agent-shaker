import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../services/api'

export const useAgentStore = defineStore('agents', () => {
  const agents = ref([])
  const loading = ref(false)
  const error = ref(null)

  const fetchAgents = async () => {
    loading.value = true
    error.value = null
    try {
      agents.value = await api.getAgents()
    } catch (err) {
      error.value = err.message
    } finally {
      loading.value = false
    }
  }

  const fetchProjectAgents = async (projectId) => {
    loading.value = true
    error.value = null
    try {
      agents.value = await api.getProjectAgents(projectId)
    } catch (err) {
      error.value = err.message
    } finally {
      loading.value = false
    }
  }

  const createAgent = async (data) => {
    loading.value = true
    error.value = null
    try {
      const agent = await api.createAgent(data)
      agents.value.push(agent)
      return agent
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const updateAgentStatus = async (id, status) => {
    try {
      await api.updateAgentStatus(id, status)
      const agent = agents.value.find(a => a.id === id)
      if (agent) {
        agent.status = status
      }
    } catch (err) {
      error.value = err.message
      throw err
    }
  }

  return {
    agents,
    loading,
    error,
    fetchAgents,
    fetchProjectAgents,
    createAgent,
    updateAgentStatus
  }
})
