<template>
  <div class="projects-page">
    <div class="container">
      <div class="page-header">
        <h2>Projects</h2>
        <button @click="showCreateModal = true" class="btn btn-primary">
          + Create Project
        </button>
      </div>

      <div v-if="loading" class="loading">Loading projects...</div>
      <div v-else-if="error" class="error">{{ error }}</div>
      
      <div v-else class="projects-grid">
        <div v-for="project in projects" :key="project.id" class="project-card">
          <router-link :to="`/projects/${project.id}`" class="card-link">
            <div class="card-header">
              <h3>{{ project.name }}</h3>
              <span class="badge" :class="project.status">{{ project.status }}</span>
            </div>
            <p class="card-description">{{ project.description }}</p>
            <div class="card-footer">
              <span class="card-date">Created {{ formatDate(project.created_at) }}</span>
            </div>
          </router-link>
        </div>
      </div>

      <div v-if="!loading && projects.length === 0" class="empty-state">
        <h3>No projects yet</h3>
        <p>Create your first project to get started</p>
        <button @click="showCreateModal = true" class="btn btn-primary">
          Create Project
        </button>
      </div>
    </div>

    <!-- Create Project Modal -->
    <div v-if="showCreateModal" class="modal-overlay" @click="showCreateModal = false">
      <div class="modal" @click.stop>
        <h3>Create New Project</h3>
        <form @submit.prevent="handleCreateProject">
          <div class="form-group">
            <label>Project Name</label>
            <input 
              v-model="newProject.name" 
              type="text" 
              placeholder="Enter project name"
              required
            />
          </div>
          <div class="form-group">
            <label>Description</label>
            <textarea 
              v-model="newProject.description" 
              placeholder="Enter project description"
              rows="4"
            ></textarea>
          </div>
          <div class="modal-actions">
            <button type="button" @click="showCreateModal = false" class="btn btn-secondary">
              Cancel
            </button>
            <button type="submit" class="btn btn-primary">
              Create Project
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useProjectStore } from '../stores/projectStore'

export default {
  name: 'Projects',
  setup() {
    const projectStore = useProjectStore()
    const showCreateModal = ref(false)
    const newProject = ref({ name: '', description: '' })

    onMounted(() => {
      projectStore.fetchProjects()
    })

    const handleCreateProject = async () => {
      try {
        await projectStore.createProject(newProject.value)
        showCreateModal.value = false
        newProject.value = { name: '', description: '' }
      } catch (error) {
        console.error('Failed to create project:', error)
      }
    }

    const formatDate = (dateString) => {
      return new Date(dateString).toLocaleDateString()
    }

    return {
      projects: projectStore.projects,
      loading: projectStore.loading,
      error: projectStore.error,
      showCreateModal,
      newProject,
      handleCreateProject,
      formatDate
    }
  }
}
</script>
