<template>
  <div>
    <div class="mb-8">
      <h2 class="text-3xl font-bold text-gray-900">All Tasks</h2>
    </div>

    <div class="flex gap-4 mb-6">
      <select v-model="statusFilter" class="px-4 py-2 border border-gray-300 rounded-md bg-white focus:outline-none focus:ring-2 focus:ring-blue-500">
        <option value="">All Status</option>
        <option value="pending">Pending</option>
        <option value="in_progress">In Progress</option>
        <option value="done">Done</option>
        <option value="blocked">Blocked</option>
      </select>
      <select v-model="priorityFilter" class="px-4 py-2 border border-gray-300 rounded-md bg-white focus:outline-none focus:ring-2 focus:ring-blue-500">
        <option value="">All Priorities</option>
        <option value="low">Low</option>
        <option value="medium">Medium</option>
        <option value="high">High</option>
      </select>
    </div>

    <div v-if="loading" class="text-center py-12 text-gray-500">Loading tasks...</div>
    <div v-else-if="error" class="p-4 bg-red-50 text-red-600 rounded-md mb-4">{{ error }}</div>
    
    <div v-else class="space-y-4">
      <div v-for="task in filteredTasks" :key="task.id" class="bg-white p-6 rounded-lg shadow-sm">
        <div class="flex justify-between items-start mb-3">
          <h3 class="text-xl font-semibold text-gray-900">{{ task.title }}</h3>
          <div class="flex gap-2">
            <span :class="['px-2 py-1 rounded text-xs font-semibold', task.priority === 'high' ? 'bg-red-100 text-red-800' : task.priority === 'medium' ? 'bg-yellow-100 text-yellow-800' : 'bg-blue-100 text-blue-800']">{{ task.priority }}</span>
            <span :class="['px-2 py-1 rounded text-xs font-semibold', task.status === 'done' ? 'bg-green-100 text-green-800' : task.status === 'in_progress' ? 'bg-blue-100 text-blue-800' : task.status === 'pending' ? 'bg-gray-100 text-gray-800' : 'bg-red-100 text-red-800']">{{ task.status }}</span>
          </div>
        </div>
        <p class="text-gray-600 mb-4">{{ task.description }}</p>
        <div class="flex flex-wrap gap-4 text-sm text-gray-500">
          <span><strong class="font-medium text-gray-900">Project:</strong> {{ task.project_id }}</span>
          <span><strong class="font-medium text-gray-900">Agent:</strong> {{ task.agent_id }}</span>
          <span><strong class="font-medium text-gray-900">Created:</strong> {{ formatDate(task.created_at) }}</span>
        </div>
      </div>
    </div>

    <div v-if="!loading && filteredTasks.length === 0" class="text-center py-12">
      <h3 class="text-xl font-semibold text-gray-900 mb-2">No tasks found</h3>
      <p class="text-gray-600">{{ statusFilter || priorityFilter ? 'Try adjusting your filters' : 'Tasks will appear here when they are created' }}</p>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useTaskStore } from '../stores/taskStore'

export default {
  name: 'Tasks',
  setup() {
    const taskStore = useTaskStore()
    const statusFilter = ref('')
    const priorityFilter = ref('')

    onMounted(() => {
      taskStore.fetchTasks()
    })

    const filteredTasks = computed(() => {
      let filtered = taskStore.tasks

      if (statusFilter.value) {
        filtered = filtered.filter(t => t.status === statusFilter.value)
      }

      if (priorityFilter.value) {
        filtered = filtered.filter(t => t.priority === priorityFilter.value)
      }

      return filtered
    })

    const formatDate = (dateString) => {
      return new Date(dateString).toLocaleString()
    }

    return {
      tasks: taskStore.tasks,
      loading: taskStore.loading,
      error: taskStore.error,
      statusFilter,
      priorityFilter,
      filteredTasks,
      formatDate
    }
  }
}
</script>
