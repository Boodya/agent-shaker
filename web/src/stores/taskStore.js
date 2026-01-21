import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../services/api'

export const useTaskStore = defineStore('tasks', () => {
  const tasks = ref([])
  const loading = ref(false)
  const error = ref(null)

  const fetchTasks = async () => {
    loading.value = true
    error.value = null
    try {
      tasks.value = await api.getTasks()
    } catch (err) {
      error.value = err.message
    } finally {
      loading.value = false
    }
  }

  const fetchProjectTasks = async (projectId) => {
    loading.value = true
    error.value = null
    try {
      tasks.value = await api.getProjectTasks(projectId)
    } catch (err) {
      error.value = err.message
    } finally {
      loading.value = false
    }
  }

  const fetchAgentTasks = async (agentId) => {
    loading.value = true
    error.value = null
    try {
      tasks.value = await api.getAgentTasks(agentId)
    } catch (err) {
      error.value = err.message
    } finally {
      loading.value = false
    }
  }

  const createTask = async (data) => {
    loading.value = true
    error.value = null
    try {
      const task = await api.createTask(data)
      tasks.value.push(task)
      return task
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const updateTaskStatus = async (id, status) => {
    try {
      await api.updateTaskStatus(id, status)
      const task = tasks.value.find(t => t.id === id)
      if (task) {
        task.status = status
      }
    } catch (err) {
      error.value = err.message
      throw err
    }
  }

  return {
    tasks,
    loading,
    error,
    fetchTasks,
    fetchProjectTasks,
    fetchAgentTasks,
    createTask,
    updateTaskStatus
  }
})
