<template>
  <div class="tasks-page">
    <div class="container">
      <div class="page-header">
        <h2>All Tasks</h2>
      </div>

      <div class="filters">
        <select v-model="statusFilter" class="filter-select">
          <option value="">All Status</option>
          <option value="pending">Pending</option>
          <option value="in_progress">In Progress</option>
          <option value="done">Done</option>
          <option value="blocked">Blocked</option>
        </select>
        <select v-model="priorityFilter" class="filter-select">
          <option value="">All Priorities</option>
          <option value="low">Low</option>
          <option value="medium">Medium</option>
          <option value="high">High</option>
        </select>
      </div>

      <div v-if="loading" class="loading">Loading tasks...</div>
      <div v-else-if="error" class="error">{{ error }}</div>
      
      <div v-else class="tasks-list">
        <div v-for="task in filteredTasks" :key="task.id" class="task-card">
          <div class="task-header">
            <h3>{{ task.title }}</h3>
            <div class="task-badges">
              <span class="priority" :class="task.priority">{{ task.priority }}</span>
              <span class="status" :class="task.status">{{ task.status }}</span>
            </div>
          </div>
          <p class="task-description">{{ task.description }}</p>
          <div class="task-meta">
            <span><strong>Project:</strong> {{ task.project_id }}</span>
            <span><strong>Agent:</strong> {{ task.agent_id }}</span>
            <span><strong>Created:</strong> {{ formatDate(task.created_at) }}</span>
          </div>
        </div>
      </div>

      <div v-if="!loading && filteredTasks.length === 0" class="empty-state">
        <h3>No tasks found</h3>
        <p>{{ statusFilter || priorityFilter ? 'Try adjusting your filters' : 'Tasks will appear here when they are created' }}</p>
      </div>
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
