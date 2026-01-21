import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '../services/api'

export const useProjectStore = defineStore('projects', () => {
  const projects = ref([])
  const currentProject = ref(null)
  const loading = ref(false)
  const error = ref(null)

  const fetchProjects = async () => {
    loading.value = true
    error.value = null
    try {
      projects.value = await api.getProjects()
    } catch (err) {
      error.value = err.message
    } finally {
      loading.value = false
    }
  }

  const fetchProject = async (id) => {
    loading.value = true
    error.value = null
    try {
      currentProject.value = await api.getProject(id)
    } catch (err) {
      error.value = err.message
    } finally {
      loading.value = false
    }
  }

  const createProject = async (data) => {
    loading.value = true
    error.value = null
    try {
      const project = await api.createProject(data)
      projects.value.push(project)
      return project
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    projects,
    currentProject,
    loading,
    error,
    fetchProjects,
    fetchProject,
    createProject
  }
})
