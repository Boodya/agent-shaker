<template>
  <div id="app">
    <nav class="navbar">
      <div class="container">
        <div class="navbar-brand">
          <h1>🚀 MCP Task Tracker</h1>
          <span class="status-badge" :class="connectionStatus">
            {{ connectionStatus === 'online' ? '● Connected' : '○ Disconnected' }}
          </span>
        </div>
        <div class="navbar-menu">
          <router-link to="/" class="nav-item">Dashboard</router-link>
          <router-link to="/projects" class="nav-item">Projects</router-link>
          <router-link to="/agents" class="nav-item">Agents</router-link>
          <router-link to="/tasks" class="nav-item">Tasks</router-link>
          <router-link to="/docs" class="nav-item">📚 Docs</router-link>
        </div>
      </div>
    </nav>

    <main class="main-content">
      <router-view />
    </main>

    <footer class="footer">
      <div class="container">
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
