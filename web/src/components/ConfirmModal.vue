<template>
  <div v-if="show" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click.self="$emit('close')">
    <div class="bg-white p-6 rounded-lg max-w-sm w-full mx-4">
      <h3 class="text-xl font-semibold mb-4 text-red-600">⚠️ {{ title }}</h3>
      <p class="text-gray-600 mb-2">{{ message }}</p>
      <p v-if="warning" class="text-sm text-orange-600 mb-6">⚠️ {{ warning }}</p>
      <div class="flex justify-end gap-3">
        <button 
          @click="$emit('close')" 
          class="bg-gray-200 hover:bg-gray-300 text-gray-800 px-4 py-2 rounded-md font-medium transition-colors"
        >
          Cancel
        </button>
        <button 
          @click="$emit('confirm')" 
          class="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-md font-medium transition-colors"
        >
          {{ confirmText }}
        </button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ConfirmModal',
  props: {
    show: {
      type: Boolean,
      required: true
    },
    title: {
      type: String,
      default: 'Confirm Delete'
    },
    message: {
      type: String,
      required: true
    },
    warning: {
      type: String,
      default: ''
    },
    confirmText: {
      type: String,
      default: 'Delete'
    }
  },
  emits: ['close', 'confirm']
}
</script>
