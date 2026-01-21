<template>
  <div class="documentation">
    <div class="docs-header">
      <h1>📚 Documentation</h1>
      <p>Browse project documentation and guides</p>
    </div>

    <div class="docs-layout">
      <!-- Sidebar with document list -->
      <aside class="docs-sidebar">
        <div class="docs-nav">
          <h3>📖 Guides</h3>
          <ul class="docs-list">
            <li 
              v-for="doc in guidesDocs" 
              :key="doc.path"
              :class="{ active: currentDoc === doc.path }"
              @click="loadDoc(doc.path)"
            >
              <span class="doc-icon">{{ doc.icon }}</span>
              <span class="doc-name">{{ doc.name }}</span>
            </li>
          </ul>

          <h3>🏗️ Architecture</h3>
          <ul class="docs-list">
            <li 
              v-for="doc in archDocs" 
              :key="doc.path"
              :class="{ active: currentDoc === doc.path }"
              @click="loadDoc(doc.path)"
            >
              <span class="doc-icon">{{ doc.icon }}</span>
              <span class="doc-name">{{ doc.name }}</span>
            </li>
          </ul>

          <h3>📡 API & Integration</h3>
          <ul class="docs-list">
            <li 
              v-for="doc in apiDocs" 
              :key="doc.path"
              :class="{ active: currentDoc === doc.path }"
              @click="loadDoc(doc.path)"
            >
              <span class="doc-icon">{{ doc.icon }}</span>
              <span class="doc-name">{{ doc.name }}</span>
            </li>
          </ul>

          <h3>ℹ️ About</h3>
          <ul class="docs-list">
            <li 
              v-for="doc in aboutDocs" 
              :key="doc.path"
              :class="{ active: currentDoc === doc.path }"
              @click="loadDoc(doc.path)"
            >
              <span class="doc-icon">{{ doc.icon }}</span>
              <span class="doc-name">{{ doc.name }}</span>
            </li>
          </ul>
        </div>
      </aside>

      <!-- Main content area -->
      <main class="docs-content">
        <div v-if="loading" class="loading">
          <div class="spinner"></div>
          <p>Loading documentation...</p>
        </div>

        <div v-else-if="error" class="error-state">
          <span class="error-icon">⚠️</span>
          <h2>Failed to Load Document</h2>
          <p>{{ error }}</p>
          <button @click="loadDoc(currentDoc)" class="btn btn-primary">
            Try Again
          </button>
        </div>

        <div v-else-if="content" class="markdown-content" v-html="renderedContent"></div>

        <div v-else class="empty-state">
          <span class="empty-icon">📄</span>
          <h2>Select a Document</h2>
          <p>Choose a document from the sidebar to view its contents</p>
        </div>
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import api from '@/services/api'

const currentDoc = ref('')
const content = ref('')
const loading = ref(false)
const error = ref(null)

// Document categories
const guidesDocs = [
  { name: 'Quick Start', path: 'QUICKSTART.md', icon: '🚀' },
  { name: 'Agent Setup Guide', path: 'AGENT_SETUP_GUIDE.md', icon: '🤖' },
  { name: 'Quick Start - Agent', path: 'QUICKSTART_AGENT.md', icon: '⚡' },
  { name: 'Contributing', path: 'CONTRIBUTING.md', icon: '🤝' },
]

const archDocs = [
  { name: 'Architecture', path: 'ARCHITECTURE.md', icon: '🏗️' },
  { name: 'Implementation Summary', path: 'IMPLEMENTATION_SUMMARY.md', icon: '📋' },
  { name: 'Docker Deployment', path: 'DOCKER_DEPLOYMENT.md', icon: '🐳' },
]

const apiDocs = [
  { name: 'API Documentation', path: 'docs/API.md', icon: '📡' },
  { name: 'Copilot Integration', path: 'docs/COPILOT_INTEGRATION.md', icon: '🔌' },
]

const aboutDocs = [
  { name: 'README', path: 'README.md', icon: 'ℹ️' },
  { name: 'Improvements', path: 'IMPROVEMENTS.md', icon: '✨' },
  { name: 'Vue Setup', path: 'VUE_SETUP_COMPLETE.md', icon: '🎨' },
]

// Configure marked options
marked.setOptions({
  breaks: true,
  gfm: true,
  headerIds: true,
  mangle: false,
})

// Render markdown to HTML
const renderedContent = computed(() => {
  if (!content.value) return ''
  
  try {
    const html = marked.parse(content.value)
    return DOMPurify.sanitize(html, {
      ALLOWED_TAGS: ['h1', 'h2', 'h3', 'h4', 'h5', 'h6', 'p', 'a', 'ul', 'ol', 'li', 
                     'blockquote', 'code', 'pre', 'strong', 'em', 'br', 'hr', 'table', 
                     'thead', 'tbody', 'tr', 'th', 'td', 'img', 'div', 'span'],
      ALLOWED_ATTR: ['href', 'src', 'alt', 'title', 'class', 'id']
    })
  } catch (err) {
    console.error('Failed to render markdown:', err)
    return '<p class="error">Failed to render markdown content</p>'
  }
})

// Load document content
async function loadDoc(path) {
  if (!path) return
  
  currentDoc.value = path
  loading.value = true
  error.value = null
  content.value = ''
  
  try {
    const response = await api.get(`/api/docs/${path}`)
    content.value = response.data
  } catch (err) {
    console.error('Failed to load document:', err)
    error.value = err.response?.data?.error || 'Failed to load document'
  } finally {
    loading.value = false
  }
}

// Load first document on mount
onMounted(() => {
  loadDoc('README.md')
})
</script>

<style scoped>
.documentation {
  padding: 2rem;
  max-width: 100%;
  height: calc(100vh - 120px);
}

.docs-header {
  margin-bottom: 2rem;
}

.docs-header h1 {
  margin: 0 0 0.5rem 0;
  color: var(--text-color);
}

.docs-header p {
  margin: 0;
  color: var(--text-muted);
}

.docs-layout {
  display: flex;
  gap: 2rem;
  height: calc(100% - 100px);
}

/* Sidebar */
.docs-sidebar {
  width: 280px;
  flex-shrink: 0;
  background: white;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 1.5rem;
  overflow-y: auto;
  height: fit-content;
  max-height: 100%;
}

.docs-nav h3 {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin: 1.5rem 0 0.75rem 0;
}

.docs-nav h3:first-child {
  margin-top: 0;
}

.docs-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.docs-list li {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 1rem;
  margin-bottom: 0.25rem;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s;
}

.docs-list li:hover {
  background: var(--background-color);
}

.docs-list li.active {
  background: var(--primary-color);
  color: white;
}

.doc-icon {
  font-size: 1.25rem;
  flex-shrink: 0;
}

.doc-name {
  font-size: 0.9375rem;
  font-weight: 500;
}

/* Content area */
.docs-content {
  flex: 1;
  background: white;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 2rem;
  overflow-y: auto;
  max-height: 100%;
}

/* Loading state */
.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem 2rem;
  color: var(--text-muted);
}

.spinner {
  width: 48px;
  height: 48px;
  border: 4px solid var(--border-color);
  border-top-color: var(--primary-color);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 1rem;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Error state */
.error-state {
  text-align: center;
  padding: 4rem 2rem;
}

.error-icon {
  font-size: 4rem;
  display: block;
  margin-bottom: 1rem;
}

.error-state h2 {
  color: var(--danger-color);
  margin-bottom: 0.5rem;
}

.error-state p {
  color: var(--text-muted);
  margin-bottom: 2rem;
}

/* Empty state */
.empty-state {
  text-align: center;
  padding: 4rem 2rem;
  color: var(--text-muted);
}

.empty-icon {
  font-size: 4rem;
  display: block;
  margin-bottom: 1rem;
  opacity: 0.5;
}

.empty-state h2 {
  margin-bottom: 0.5rem;
}

/* Markdown content styles */
.markdown-content {
  line-height: 1.7;
  color: var(--text-color);
}

.markdown-content :deep(h1) {
  font-size: 2rem;
  font-weight: 700;
  margin: 2rem 0 1rem 0;
  padding-bottom: 0.5rem;
  border-bottom: 2px solid var(--border-color);
}

.markdown-content :deep(h1:first-child) {
  margin-top: 0;
}

.markdown-content :deep(h2) {
  font-size: 1.5rem;
  font-weight: 600;
  margin: 1.75rem 0 0.75rem 0;
  color: var(--text-color);
}

.markdown-content :deep(h3) {
  font-size: 1.25rem;
  font-weight: 600;
  margin: 1.5rem 0 0.75rem 0;
}

.markdown-content :deep(p) {
  margin: 0 0 1rem 0;
}

.markdown-content :deep(a) {
  color: var(--primary-color);
  text-decoration: none;
}

.markdown-content :deep(a:hover) {
  text-decoration: underline;
}

.markdown-content :deep(ul),
.markdown-content :deep(ol) {
  margin: 0 0 1rem 0;
  padding-left: 2rem;
}

.markdown-content :deep(li) {
  margin: 0.5rem 0;
}

.markdown-content :deep(code) {
  background: var(--background-color);
  padding: 0.2rem 0.4rem;
  border-radius: 3px;
  font-family: 'Courier New', monospace;
  font-size: 0.875rem;
}

.markdown-content :deep(pre) {
  background: #f6f8fa;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  padding: 1rem;
  overflow-x: auto;
  margin: 1rem 0;
}

.markdown-content :deep(pre code) {
  background: none;
  padding: 0;
}

.markdown-content :deep(blockquote) {
  border-left: 4px solid var(--primary-color);
  padding-left: 1rem;
  margin: 1rem 0;
  color: var(--text-muted);
  font-style: italic;
}

.markdown-content :deep(table) {
  width: 100%;
  border-collapse: collapse;
  margin: 1rem 0;
}

.markdown-content :deep(th),
.markdown-content :deep(td) {
  border: 1px solid var(--border-color);
  padding: 0.75rem;
  text-align: left;
}

.markdown-content :deep(th) {
  background: var(--background-color);
  font-weight: 600;
}

.markdown-content :deep(img) {
  max-width: 100%;
  height: auto;
  border-radius: 6px;
  margin: 1rem 0;
}

.markdown-content :deep(hr) {
  border: none;
  border-top: 1px solid var(--border-color);
  margin: 2rem 0;
}

/* Responsive */
@media (max-width: 768px) {
  .documentation {
    padding: 1rem;
  }

  .docs-layout {
    flex-direction: column;
    height: auto;
  }

  .docs-sidebar {
    width: 100%;
    max-height: 400px;
  }

  .docs-content {
    min-height: 500px;
  }
}
</style>
