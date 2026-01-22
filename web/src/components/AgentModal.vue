<template>
  <div v-if="show" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="$emit('close')">
    <div class="bg-white p-6 rounded-lg max-w-md w-full mx-4">
      <h3 class="text-xl font-semibold mb-6">{{ isEdit ? 'Edit Agent' : 'Add Agent to Project' }}</h3>
      <form @submit.prevent="handleSubmit">
        <div class="mb-4">
          <label class="block text-sm font-medium text-gray-700 mb-2">Agent Name</label>
          <input 
            v-model="formData.name" 
            type="text" 
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500" 
            required 
          />
        </div>
        <div class="mb-4">
          <label class="block text-sm font-medium text-gray-700 mb-2">Role</label>
          <select 
            v-model="formData.role" 
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500" 
            required
          >
            <optgroup label="Development Roles">
              <option value="frontend">Frontend Developer</option>
              <option value="backend">Backend Developer</option>
              <option value="fullstack">Full Stack Developer</option>
              <option value="mobile">Mobile Developer</option>
              <option value="devops">DevOps Engineer</option>
              <option value="qa">QA Engineer</option>
              <option value="security">Security Engineer</option>
            </optgroup>
            <optgroup label="Agile Roles">
              <option value="product-owner">Product Owner</option>
              <option value="scrum-master">Scrum Master</option>
              <option value="agile-coach">Agile Coach</option>
            </optgroup>
            <optgroup label="R&D Roles">
              <option value="architect">Solution Architect</option>
              <option value="tech-lead">Tech Lead</option>
              <option value="researcher">Research Engineer</option>
              <option value="data-scientist">Data Scientist</option>
              <option value="ml-engineer">ML Engineer</option>
            </optgroup>
            <optgroup label="Design & UX">
              <option value="ux-designer">UX Designer</option>
              <option value="ui-designer">UI Designer</option>
              <option value="ux-researcher">UX Researcher</option>
            </optgroup>
          </select>
        </div>
        <div class="mb-4">
          <label class="block text-sm font-medium text-gray-700 mb-2">Team</label>
          <input 
            v-model="formData.team" 
            type="text" 
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500" 
          />
        </div>
        <div v-if="isEdit" class="mb-4">
          <label class="block text-sm font-medium text-gray-700 mb-2">Status</label>
          <select 
            v-model="formData.status" 
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="active">Active</option>
            <option value="idle">Idle</option>
            <option value="offline">Offline</option>
          </select>
        </div>
        <div class="flex justify-end gap-3 mt-6">
          <button 
            type="button" 
            @click="$emit('close')" 
            class="bg-gray-200 hover:bg-gray-300 text-gray-800 px-4 py-2 rounded-md font-medium transition-colors"
          >
            Cancel
          </button>
          <button 
            type="submit" 
            class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md font-medium transition-colors"
          >
            {{ isEdit ? 'Update' : 'Add' }} Agent
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import { ref, watch } from 'vue'

export default {
  name: 'AgentModal',
  props: {
    show: {
      type: Boolean,
      required: true
    },
    agent: {
      type: Object,
      default: null
    }
  },
  emits: ['close', 'save'],
  setup(props, { emit }) {
    const isEdit = ref(false)
    const formData = ref({
      name: '',
      role: 'frontend',
      team: '',
      status: 'active'
    })

    watch(() => props.agent, (newAgent) => {
      if (newAgent) {
        isEdit.value = true
        formData.value = {
          name: newAgent.name,
          role: newAgent.role,
          team: newAgent.team || '',
          status: newAgent.status || 'active'
        }
      } else {
        isEdit.value = false
        formData.value = {
          name: '',
          role: 'frontend',
          team: '',
          status: 'active'
        }
      }
    }, { immediate: true })

    const handleSubmit = () => {
      emit('save', { ...formData.value })
    }

    return {
      isEdit,
      formData,
      handleSubmit
    }
  }
}
</script>
