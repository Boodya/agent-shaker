<template>
  <div>
    <div class="flex justify-between items-center mb-8">
      <h2 class="text-3xl font-bold text-gray-900">Projects</h2>
      <button @click="showCreateModal = true" class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md font-medium transition-colors">
        + Create Project
      </button>
    </div>

    <div v-if="loading" class="text-center py-12 text-gray-500">Loading projects...</div>
    <div v-else-if="error" class="p-4 bg-red-50 text-red-600 rounded-md mb-4">{{ error }}</div>
    
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div v-for="project in projects" :key="project.id" class="bg-white rounded-lg shadow-sm hover:shadow-md transition-shadow overflow-hidden">
        <router-link :to="`/projects/${project.id}`" class="block p-6 text-inherit no-underline">
          <div class="flex justify-between items-start mb-4">
            <h3 class="text-xl font-semibold text-gray-900">{{ project.name }}</h3>
            <span :class="['px-3 py-1 rounded-full text-xs font-semibold uppercase', project.status === 'active' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800']">{{ project.status }}</span>
          </div>
          <p class="text-gray-600 mb-4 leading-relaxed">{{ project.description }}</p>
          <div class="text-gray-500 text-sm">
            <span>Created {{ formatDate(project.created_at) }}</span>
          </div>
        </router-link>
      </div>
    </div>

    <div v-if="!loading && projects.length === 0" class="text-center py-12">
      <h3 class="text-xl font-semibold text-gray-900 mb-2">No projects yet</h3>
      <p class="text-gray-600 mb-4">Create your first project to get started</p>
      <button @click="showCreateModal = true" class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md font-medium transition-colors">
        Create Project
      </button>
    </div>

    <!-- Create Project Modal -->
    <div v-if="showCreateModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click="showCreateModal = false">
      <div class="bg-white p-6 rounded-lg max-w-md w-full mx-4" @click.stop>
        <h3 class="text-xl font-semibold mb-6">Create New Project</h3>
        <form @submit.prevent="handleCreateProject">
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">Project Name</label>
            <input 
              v-model="newProject.name" 
              type="text" 
              placeholder="Enter project name"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            />
          </div>
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">Description</label>
            <textarea 
              v-model="newProject.description" 
              placeholder="Enter project description"
              rows="4"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            ></textarea>
          </div>
          <div class="flex justify-end gap-3 mt-6">
            <button type="button" @click="showCreateModal = false" class="bg-gray-200 hover:bg-gray-300 text-gray-800 px-4 py-2 rounded-md font-medium transition-colors">
              Cancel
            </button>
            <button type="submit" class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md font-medium transition-colors">
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
