<template>
  <div v-if="show" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="$emit('close')">
    <div class="bg-white p-6 rounded-lg max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto">
      <h3 class="text-xl font-semibold mb-6">{{ isEdit ? 'Edit Context' : 'Add Context / Documentation' }}</h3>
      <form @submit.prevent="handleSubmit">
        <div class="mb-4">
          <label class="block text-sm font-medium text-gray-700 mb-2">Title</label>
          <input 
            v-model="formData.title" 
            type="text" 
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500" 
            required 
          />
        </div>
        <div class="mb-4">
          <label class="block text-sm font-medium text-gray-700 mb-2">Agent</label>
          <select 
            v-model="formData.agent_id" 
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500" 
            required
          >
            <option value="">Select an agent</option>
            <option v-for="agent in agents" :key="agent.id" :value="agent.id">
              {{ agent.name }}
            </option>
          </select>
        </div>
        <div class="mb-4">
          <label class="block text-sm font-medium text-gray-700 mb-2">Related Task (Optional)</label>
          <select 
            v-model="formData.task_id" 
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="">None</option>
            <option v-for="task in tasks" :key="task.id" :value="task.id">
              {{ task.title }}
            </option>
          </select>
        </div>
        <div class="mb-4">
          <label class="block text-sm font-medium text-gray-700 mb-2">Content (Markdown)</label>
          <textarea 
            v-model="formData.content" 
            rows="12" 
            placeholder="Write your documentation in Markdown format..."
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          ></textarea>
          <small class="text-gray-500 text-sm">Supports Markdown: **bold**, *italic*, [link](url), etc.</small>
        </div>
        <div class="mb-4">
          <label class="block text-sm font-medium text-gray-700 mb-2">Tags (comma-separated)</label>
          <input 
            v-model="formData.tagsString" 
            type="text" 
            placeholder="api, documentation, backend"
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
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
            {{ isEdit ? 'Update' : 'Create' }} Context
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import { ref, watch } from 'vue'

export default {
  name: 'ContextModal',
  props: {
    show: {
      type: Boolean,
      required: true
    },
    context: {
      type: Object,
      default: null
    },
    agents: {
      type: Array,
      required: true
    },
    tasks: {
      type: Array,
      required: true
    }
  },
  emits: ['close', 'save'],
  setup(props, { emit }) {
    const isEdit = ref(false)
    const formData = ref({
      title: '',
      agent_id: '',
      task_id: '',
      content: '',
      tagsString: ''
    })

    watch(() => props.context, (newContext) => {
      if (newContext) {
        isEdit.value = true
        formData.value = {
          title: newContext.title,
          agent_id: newContext.agent_id,
          task_id: newContext.task_id || '',
          content: newContext.content,
          tagsString: newContext.tags ? newContext.tags.join(', ') : ''
        }
      } else {
        isEdit.value = false
        formData.value = {
          title: '',
          agent_id: '',
          task_id: '',
          content: '',
          tagsString: ''
        }
      }
    }, { immediate: true })

    const handleSubmit = () => {
      const tags = formData.value.tagsString
        .split(',')
        .map(tag => tag.trim())
        .filter(tag => tag.length > 0)

      emit('save', {
        title: formData.value.title,
        agent_id: formData.value.agent_id,
        task_id: formData.value.task_id || null,
        content: formData.value.content,
        tags
      })
    }

    return {
      isEdit,
      formData,
      handleSubmit
    }
  }
}
</script>
