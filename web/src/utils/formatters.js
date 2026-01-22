/**
 * Format date to locale string
 * @param {string} dateString - ISO date string
 * @returns {string} Formatted date string
 */
export const formatDate = (dateString) => {
  if (!dateString) return ''
  return new Date(dateString).toLocaleString()
}

/**
 * Convert tags string to array
 * @param {string} tagsString - Comma-separated tags
 * @returns {string[]} Array of trimmed tags
 */
export const parseTags = (tagsString) => {
  if (!tagsString) return []
  return tagsString
    .split(',')
    .map(tag => tag.trim())
    .filter(tag => tag.length > 0)
}

/**
 * Convert tags array to string
 * @param {string[]} tags - Array of tags
 * @returns {string} Comma-separated string
 */
export const tagsToString = (tags) => {
  if (!Array.isArray(tags)) return ''
  return tags.join(', ')
}

/**
 * Get unique tags from contexts
 * @param {Array} contexts - Array of context objects
 * @returns {string[]} Sorted array of unique tags
 */
export const getUniqueTags = (contexts) => {
  const tags = new Set()
  contexts.forEach(context => {
    if (context.tags && Array.isArray(context.tags)) {
      context.tags.forEach(tag => tags.add(tag))
    }
  })
  return Array.from(tags).sort()
}
