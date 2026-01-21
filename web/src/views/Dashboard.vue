<template>
  <div>
    <div class="mb-8">
      <h2 class="text-3xl font-bold text-gray-900">Dashboard</h2>
      <p class="text-gray-600 mt-2">Overview of your MCP Task Tracker</p>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
      <div class="bg-white p-6 rounded-lg shadow-sm hover:shadow-md transition-shadow flex items-center gap-4">
        <div class="text-4xl">📁</div>
        <div>
          <h3 class="text-2xl font-bold text-blue-600">{{ projects.length }}</h3>
          <p class="text-gray-600 text-sm">Projects</p>
        </div>
      </div>

      <div class="bg-white p-6 rounded-lg shadow-sm hover:shadow-md transition-shadow flex items-center gap-4">
        <div class="text-4xl">🤖</div>
        <div>
          <h3 class="text-2xl font-bold text-blue-600">{{ agents.length }}</h3>
          <p class="text-gray-600 text-sm">Agents</p>
        </div>
      </div>

      <div class="bg-white p-6 rounded-lg shadow-sm hover:shadow-md transition-shadow flex items-center gap-4">
        <div class="text-4xl">📋</div>
        <div>
          <h3 class="text-2xl font-bold text-blue-600">{{ tasks.length }}</h3>
          <p class="text-gray-600 text-sm">Total Tasks</p>
        </div>
      </div>

      <div class="bg-white p-6 rounded-lg shadow-sm hover:shadow-md transition-shadow flex items-center gap-4">
        <div class="text-4xl">✅</div>
        <div>
          <h3 class="text-2xl font-bold text-blue-600">{{ completedTasks }}</h3>
          <p class="text-gray-600 text-sm">Completed</p>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div class="bg-white p-6 rounded-lg shadow-sm">
        <h3 class="text-xl font-semibold text-gray-900 mb-4">Recent Projects</h3>
        <div class="space-y-3">
          <div v-for="project in recentProjects" :key="project.id" class="rounded-lg overflow-hidden">
            <router-link :to="`/projects/${project.id}`" class="block p-4 bg-gray-50 hover:bg-gray-100 transition-colors text-inherit no-underline">
              <h4 class="mb-1">{{ project.name }}</h4>
              <p class="text-gray-600 text-sm mb-2">{{ project.description }}</p>
              <span :class=" [
                'inline-block px-3 py-1 rounded-full text-xs font-semibold uppercase',
                project.status === 'active' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
              ]">{{ project.status }}</span>
            </router-link>
          </div>
          <p v-if="projects.length === 0" class="text-center py-12 text-gray-500">No projects yet</p>
        </div>
      </div>

      <div class="bg-white p-6 rounded-lg shadow-sm">
        <h3 class="text-xl font-semibold text-gray-900 mb-4">Active Agents</h3>
        <div class="space-y-3">
          <div v-for="agent in activeAgents" :key="agent.id" class="flex justify-between items-center p-3 bg-gray-50 rounded-lg">
            <div class="flex items-center gap-3">
              <span class="font-medium">{{ agent.name }}</span>
              <span :class=" [
                'px-2 py-1 rounded text-xs font-semibold',
                agent.role === 'frontend' ? 'bg-blue-100 text-blue-800' : 'bg-pink-100 text-pink-800'
              ]">{{ agent.role }}</span>
            </div>
            <span :class=" [
              'px-3 py-1 rounded-full text-sm font-semibold',
              agent.status === 'active' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
            ]">{{ agent.status }}</span>
          </div>
          <p v-if="activeAgents.length === 0" class="text-center py-12 text-gray-500">No active agents</p>
        </div>
      </div>

      <div class="col-span-full bg-white p-6 rounded-lg shadow-sm">
        <h3 class="text-xl font-semibold text-gray-900 mb-4">Recent Tasks</h3>
        <div class="space-y-3">
          <div v-for="task in recentTasks" :key="task.id" class="flex justify-between items-center p-3 bg-gray-50 rounded-lg">
            <div class="flex-1">
              <h4 class="font-medium">{{ task.title }}</h4>
              <p class="text-gray-600 text-sm">{{ task.description }}</p>
            </div>
            <div class="flex gap-2">
              <span :class=" [
                'px-2 py-1 rounded text-xs font-semibold',
                task.priority === 'high' ? 'bg-red-100 text-red-800' : 
                task.priority === 'medium' ? 'bg-yellow-100 text-yellow-800' : 'bg-blue-100 text-blue-800'
              ]">{{ task.priority }}</span>
              <span :class=" [
                'px-2 py-1 rounded text-xs font-semibold',
                task.status === 'done' ? 'bg-green-100 text-green-800' : 
                task.status === 'in_progress' ? 'bg-blue-100 text-blue-800' : 
                task.status === 'pending' ? 'bg-gray-100 text-gray-800' : 'bg-red-100 text-red-800'
              ]">{{ task.status }}</span>
            </div>
          </div>
          <p v-if="tasks.length === 0" class="text-center py-12 text-gray-500">No tasks yet</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { computed, onMounted } from 'vue'
import { useProjectStore } from '../stores/projectStore'
import { useAgentStore } from '../stores/agentStore'
import { useTaskStore } from '../stores/taskStore'

export default {
  name: 'Dashboard',
  setup() {
    const projectStore = useProjectStore()
    const agentStore = useAgentStore()
    const taskStore = useTaskStore()

    onMounted(() => {
      projectStore.fetchProjects()
      agentStore.fetchAgents()
      taskStore.fetchTasks()
    })

    const recentProjects = computed(() => projectStore.projects.slice(0, 5))
    const activeAgents = computed(() => 
      agentStore.agents.filter(a => a.status === 'active').slice(0, 5)
    )
    const recentTasks = computed(() => taskStore.tasks.slice(0, 10))
    const completedTasks = computed(() => 
      taskStore.tasks.filter(t => t.status === 'done').length
    )

    return {
      projects: projectStore.projects,
      agents: agentStore.agents,
      tasks: taskStore.tasks,
      recentProjects,
      activeAgents,
      recentTasks,
      completedTasks
    }
  }
}
</script>
