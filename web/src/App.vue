<template>
  <div id="app" class="min-h-screen bg-gray-50">
    <nav class="sticky top-0 z-50 bg-white border-b border-gray-200 shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex flex-col sm:flex-row justify-between items-center py-4 gap-3">
          <div class="flex items-center gap-3 sm:gap-4">
            <h1 class="text-lg sm:text-xl font-bold text-blue-600">🚀 MCP Task Tracker</h1>
            <span :class="['px-2 sm:px-3 py-1 rounded-full text-xs sm:text-sm font-semibold', connectionStatus === 'online' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800' ]">
              {{ connectionStatus === 'online' ? '● Connected' : '○ Disconnected' }}
            </span>
          </div>
          <div class="flex flex-wrap justify-center gap-1 sm:gap-2">
            <router-link to="/" class="text-gray-700 hover:text-blue-600 hover:bg-blue-50 px-3 sm:px-4 py-2 rounded-md transition-colors text-sm font-medium" active-class="bg-blue-600 text-white hover:bg-blue-700 hover:text-white">Dashboard</router-link>
            <router-link to="/projects" class="text-gray-700 hover:text-blue-600 hover:bg-blue-50 px-3 sm:px-4 py-2 rounded-md transition-colors text-sm font-medium" active-class="bg-blue-600 text-white hover:bg-blue-700 hover:text-white">Projects</router-link>
            <router-link to="/agents" class="text-gray-700 hover:text-blue-600 hover:bg-blue-50 px-3 sm:px-4 py-2 rounded-md transition-colors text-sm font-medium" active-class="bg-blue-600 text-white hover:bg-blue-700 hover:text-white">Agents</router-link>
            <router-link to="/tasks" class="text-gray-700 hover:text-blue-600 hover:bg-blue-50 px-3 sm:px-4 py-2 rounded-md transition-colors text-sm font-medium" active-class="bg-blue-600 text-white hover:bg-blue-700 hover:text-white">Tasks</router-link>
            <router-link to="/docs" class="text-gray-700 hover:text-blue-600 hover:bg-blue-50 px-3 sm:px-4 py-2 rounded-md transition-colors text-sm font-medium" active-class="bg-blue-600 text-white hover:bg-blue-700 hover:text-white">📚 Docs</router-link>
          </div>
        </div>
      </div>
    </nav>

    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 sm:py-8">
      <router-view />
    </main>

    <footer class="bg-white border-t border-gray-200 py-6 sm:py-8 mt-12 sm:mt-16">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center text-gray-500 text-sm">
        <p>&copy; 2026 MCP Task Tracker - AI Agent Coordination System</p>
      </div>
    </footer>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted } from 'vue'
import { useWebSocket } from './composables/useWebSocket'

export default {
  name: 'App',
  setup() {
    const connectionStatus = ref('offline')
    const { connect, disconnect, isConnected } = useWebSocket()

    onMounted(() => {
      connect()
      const checkConnection = setInterval(() => {
        connectionStatus.value = isConnected.value ? 'online' : 'offline'
      }, 1000)
      
      onUnmounted(() => {
        clearInterval(checkConnection)
        disconnect()
      })
    })

    return {
      connectionStatus
    }
  }
}
</script>
