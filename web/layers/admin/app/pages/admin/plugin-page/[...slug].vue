<script setup lang="ts">
/**
 * Phase 4.2: Dynamic plugin page renderer.
 *
 * Catches all routes under /admin/plugin-page/* and renders the
 * corresponding plugin Vue component loaded from admin.mjs.
 *
 * Plugin pages are registered via contributes.navigation with routes
 * pointing to /admin/plugin-page/{pluginId}/{pageName}.
 */
import { createPluginAsyncComponent } from '~/composables/usePluginComponents'

definePageMeta({
  layout: 'admin',
})

const route = useRoute()
const slugParts = computed(() => {
  const slug = route.params.slug
  return Array.isArray(slug) ? slug : [slug]
})

const pluginId = computed(() => slugParts.value[0] || '')
const componentName = computed(() => slugParts.value[1] || 'default')

// Use shallowRef + watch to avoid creating a new async component on every computed evaluation
const PluginComponent = shallowRef<ReturnType<typeof createPluginAsyncComponent> | null>(null)
let _lastKey = ''
watch([pluginId, componentName], ([pid, cname]) => {
  const key = `${pid}:${cname}`
  if (key === _lastKey) return
  _lastKey = key
  PluginComponent.value = pid ? createPluginAsyncComponent(pid, cname) : null
}, { immediate: true })
</script>

<template>
  <div class="h-full overflow-y-auto">
    <ClientOnly>
      <div v-if="!PluginComponent" class="flex items-center justify-center h-64 text-muted">
        Plugin not found
      </div>
      <component :is="PluginComponent" v-else />
    </ClientOnly>
  </div>
</template>
