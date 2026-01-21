import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../services/api'

export const useContextStore = defineStore('contexts', () => {
  const contexts = ref([])
  const currentContext = ref(null)
  const loading = ref(false)
  const error = ref(null)

  const fetchContexts = async () => {
    loading.value = true
    error.value = null
    try {
      contexts.value = await api.getContexts()
    } catch (err) {
      error.value = err.message
    } finally {
      loading.value = false
    }
  }

  const fetchProjectContexts = async (projectId) => {
    loading.value = true
    error.value = null
    try {
      contexts.value = await api.getProjectContexts(projectId)
    } catch (err) {
      error.value = err.message
    } finally {
      loading.value = false
    }
  }

  const fetchContext = async (id) => {
    loading.value = true
    error.value = null
    try {
      currentContext.value = await api.getContext(id)
    } catch (err) {
      error.value = err.message
    } finally {
      loading.value = false
    }
  }

  const createContext = async (data) => {
    loading.value = true
    error.value = null
    try {
      const context = await api.createContext(data)
      contexts.value.unshift(context)
      return context
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const updateContext = async (id, data) => {
    loading.value = true
    error.value = null
    try {
      const updatedContext = await api.updateContext(id, data)
      const index = contexts.value.findIndex(c => c.id === id)
      if (index !== -1) {
        contexts.value[index] = updatedContext
      }
      return updatedContext
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  const deleteContext = async (id) => {
    loading.value = true
    error.value = null
    try {
      await api.deleteContext(id)
      contexts.value = contexts.value.filter(c => c.id !== id)
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    contexts,
    currentContext,
    loading,
    error,
    fetchContexts,
    fetchProjectContexts,
    fetchContext,
    createContext,
    updateContext,
    deleteContext
  }
})
