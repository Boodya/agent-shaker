<template>
  <div class="stat-card">
    <div class="stat-header">
      <div class="stat-icon" :style="{ backgroundColor: iconBgColor }">
        <span>{{ icon }}</span>
      </div>
      <div class="stat-info">
        <h3 class="stat-title">{{ title }}</h3>
        <p class="stat-value">{{ value }}</p>
      </div>
    </div>
    <div v-if="breakdown && breakdown.length > 0" class="stat-breakdown">
      <div 
        v-for="item in breakdown" 
        :key="item.label"
        class="breakdown-item"
      >
        <span class="breakdown-label">{{ item.label }}</span>
        <span class="breakdown-value" :style="{ color: item.color }">{{ item.value }}</span>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'StatCard',
  props: {
    title: {
      type: String,
      required: true
    },
    value: {
      type: [Number, String],
      required: true
    },
    icon: {
      type: String,
      default: '📊'
    },
    iconBgColor: {
      type: String,
      default: '#3b82f6'
    },
    breakdown: {
      type: Array,
      default: () => []
    }
  }
}
</script>

<style scoped>
.stat-card {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  border: 1px solid #e5e7eb;
}

.stat-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transform: translateY(-2px);
}

.stat-header {
  display: flex;
  align-items: flex-start;
  gap: 1rem;
  margin-bottom: 1rem;
}

.stat-icon {
  width: 3rem;
  height: 3rem;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  flex-shrink: 0;
}

.stat-info {
  flex: 1;
}

.stat-title {
  font-size: 0.875rem;
  font-weight: 500;
  color: #6b7280;
  margin: 0 0 0.5rem 0;
}

.stat-value {
  font-size: 2rem;
  font-weight: 700;
  color: #111827;
  margin: 0;
  line-height: 1;
}

.stat-breakdown {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  padding-top: 1rem;
  border-top: 1px solid #f3f4f6;
}

.breakdown-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.875rem;
}

.breakdown-label {
  color: #6b7280;
  font-weight: 500;
}

.breakdown-value {
  font-weight: 600;
}
</style>
