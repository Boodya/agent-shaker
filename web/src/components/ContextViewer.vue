<template>
  <div v-if="show" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="$emit('close')">
    <div class="bg-white p-6 rounded-lg max-w-4xl w-full mx-4 max-h-[90vh] overflow-y-auto">
      <div class="flex justify-between items-start mb-6">
        <h3 class="text-xl font-semibold text-gray-900">{{ context?.title }}</h3>
        <button @click="$emit('close')" class="text-gray-400 hover:text-gray-600 text-2xl">×</button>
      </div>
      <div>
        <div class="flex justify-between items-center mb-6">
          <div class="flex flex-wrap gap-2">
            <span v-for="tag in context?.tags" :key="tag" class="px-2 py-1 bg-gray-100 text-gray-800 rounded text-xs">{{ tag }}</span>
          </div>
          <div class="flex gap-4 text-sm text-gray-500">
            <span>Agent: {{ agentName }}</span>
            <span v-if="taskName">Task: {{ taskName }}</span>
            <span>{{ formattedDate }}</span>
          </div>
        </div>
        <div class="prose max-w-none" v-html="renderedContent"></div>
      </div>
    </div>
  </div>
</template>

<script>
import { computed } from 'vue'
import { marked } from 'marked'
import DOMPurify from 'dompurify'

export default {
  name: 'ContextViewer',
  props: {
    show: {
      type: Boolean,
      required: true
    },
    context: {
      type: Object,
      default: null
    },
    agentName: {
      type: String,
      default: 'Unknown'
    },
    taskName: {
      type: String,
      default: ''
    }
  },
  emits: ['close'],
  setup(props) {
    const renderedContent = computed(() => {
      if (!props.context?.content) return ''
      const html = marked(props.context.content)
      return DOMPurify.sanitize(html)
    })

    const formattedDate = computed(() => {
      if (!props.context?.created_at) return ''
      return new Date(props.context.created_at).toLocaleString()
    })

    return {
      renderedContent,
      formattedDate
    }
  }
}
</script>
