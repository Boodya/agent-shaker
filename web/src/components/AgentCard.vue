<template>
  <div class="agent-card">
    <!-- Header with name and status -->
    <div class="agent-card-header">
      <h4 class="agent-card-title">{{ agent.name }}</h4>
      <span :class="[
        'status-badge',
        agent.status === 'active' ? 'status-active' : 
        agent.status === 'idle' ? 'status-idle' : 'status-offline'
      ]">
        {{ agent.status }}
      </span>
    </div>

    <!-- Agent details -->
    <div class="agent-card-body">
      <div class="agent-detail-item">
        <span class="detail-label">Role:</span>
        <span :class="[
          'role-badge',
          agent.role === 'frontend' ? 'role-frontend' : 'role-backend'
        ]">
          {{ agent.role }}
        </span>
      </div>
      <div class="agent-detail-item">
        <span class="detail-label">Team:</span>
        <span class="detail-value">{{ agent.team || 'Not assigned' }}</span>
      </div>
      <div class="agent-detail-item">
        <span class="detail-label">Last Seen:</span>
        <span class="detail-value">{{ formatDate(agent.last_seen) }}</span>
      </div>
    </div>

    <!-- Action buttons at bottom -->
    <div class="agent-card-footer">
      <button 
        @click="$emit('setup', agent)" 
        class="action-btn btn-setup" 
        title="Download MCP Setup"
      >
        <span class="btn-icon">⚙️</span>
        <span class="btn-text">Setup</span>
      </button>
      <button 
        @click="$emit('edit', agent)" 
        class="action-btn btn-edit" 
        title="Edit Agent"
      >
        <span class="btn-icon">✏️</span>
        <span class="btn-text">Edit</span>
      </button>
      <button 
        @click="$emit('delete', agent)" 
        class="action-btn btn-delete" 
        title="Delete Agent"
      >
        <span class="btn-icon">🗑️</span>
        <span class="btn-text">Delete</span>
      </button>
    </div>
  </div>
</template>

<script>
export default {
  name: 'AgentCard',
  props: {
    agent: {
      type: Object,
      required: true
    }
  },
  emits: ['setup', 'edit', 'delete'],
  methods: {
    formatDate(dateString) {
      if (!dateString) return 'Never'
      return new Date(dateString).toLocaleString()
    }
  }
}
</script>

<style scoped>
.agent-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  transition: all 0.3s ease;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.agent-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}

/* Header */
.agent-card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 1.25rem 1.5rem;
  border-bottom: 1px solid #e5e7eb;
  gap: 1rem;
}

.agent-card-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: #1f2937;
  margin: 0;
  line-height: 1.4;
  flex: 1;
}

.status-badge {
  padding: 0.375rem 0.875rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.025em;
  white-space: nowrap;
}

.status-active {
  background-color: #d1fae5;
  color: #065f46;
}

.status-idle {
  background-color: #fef3c7;
  color: #92400e;
}

.status-offline {
  background-color: #fee2e2;
  color: #991b1b;
}

/* Body */
.agent-card-body {
  padding: 1.5rem;
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.agent-detail-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.detail-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: #6b7280;
  min-width: 80px;
}

.detail-value {
  font-size: 0.875rem;
  color: #374151;
}

.role-badge {
  padding: 0.25rem 0.75rem;
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: capitalize;
}

.role-frontend {
  background-color: #dbeafe;
  color: #1e40af;
}

.role-backend {
  background-color: #fce7f3;
  color: #9f1239;
}

/* Footer */
.agent-card-footer {
  display: flex;
  gap: 0.5rem;
  padding: 1rem 1.5rem;
  background-color: #f9fafb;
  border-top: 1px solid #e5e7eb;
}

.action-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 0.625rem 1rem;
  border: 1px solid transparent;
  border-radius: 8px;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.action-btn:hover {
  transform: translateY(-1px);
}

.action-btn:active {
  transform: translateY(0);
}

.btn-icon {
  font-size: 1rem;
}

.btn-text {
  display: none;
}

/* Button variants */
.btn-setup {
  background-color: #f3e8ff;
  color: #6b21a8;
  border-color: #e9d5ff;
}

.btn-setup:hover {
  background-color: #e9d5ff;
  border-color: #d8b4fe;
}

.btn-edit {
  background-color: #dbeafe;
  color: #1e40af;
  border-color: #bfdbfe;
}

.btn-edit:hover {
  background-color: #bfdbfe;
  border-color: #93c5fd;
}

.btn-delete {
  background-color: #fee2e2;
  color: #991b1b;
  border-color: #fecaca;
}

.btn-delete:hover {
  background-color: #fecaca;
  border-color: #fca5a5;
}

/* Show text on larger screens */
@media (min-width: 640px) {
  .btn-text {
    display: inline;
  }
  
  .btn-icon {
    font-size: 0.875rem;
  }
}

/* Responsive adjustments */
@media (max-width: 640px) {
  .agent-card-header {
    padding: 1rem;
  }

  .agent-card-body {
    padding: 1rem;
  }

  .agent-card-footer {
    padding: 0.75rem 1rem;
  }

  .action-btn {
    padding: 0.5rem;
  }
}
</style>
