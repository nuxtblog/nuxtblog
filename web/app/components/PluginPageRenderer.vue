<script setup lang="ts">
/**
 * Universal plugin page renderer — three-state rendering:
 * 1. Plugins loading → empty (no flash of "not found")
 * 2. Loaded + matched → render plugin component
 * 3. Loaded + no match → 404 error
 *
 * Data sources (in priority order):
 * 1. route.meta — set by addRoute (SPA navigation after plugins loaded)
 * 2. path lookup — matches route.path against registered pluginPages (catch-all fallback)
 */
import { createPluginAsyncComponent } from '~/composables/usePluginComponents'
import { publicPluginsLoaded } from '~/composables/usePublicPluginLoader'
import { adminPluginsLoaded } from '~/composables/usePluginLoader'

const route = useRoute()
const store = usePluginContributionsStore()

// Resolve plugin page definition: prefer route.meta, fall back to path lookup
const pageDef = computed(() => {
  if (route.meta.pluginId) {
    return {
      pluginId: route.meta.pluginId as string,
      component: route.meta.pluginComponent as string,
      moduleFile: (route.meta.pluginModuleFile as string) || 'public.mjs',
      title: (route.meta.pluginPageTitle as string) || '',
    }
  }
  // Prefer the longest matching path to avoid parent routes shadowing child routes
  // e.g. /admin/shop/products/form should match before /admin/shop/products
  const matches = store.pluginPages.filter(p => p.path && (route.path === p.path || route.path.startsWith(p.path + '/')))
  if (matches.length === 0) return null
  return matches.reduce((best, p) => (p.path!.length > best.path!.length ? p : best))
})

const pluginId = computed(() => pageDef.value?.pluginId || '')
const componentName = computed(() => pageDef.value?.component || '')
const moduleFile = computed(() => pageDef.value?.moduleFile || 'public.mjs')
const pageTitle = computed(() => pageDef.value?.title || '')

useHead(() => ({ title: pageTitle.value || undefined }))

const isAdmin = computed(() => route.path.startsWith('/admin'))
const pluginsReady = computed(() =>
  isAdmin.value ? adminPluginsLoaded.value : publicPluginsLoaded.value,
)

// Once plugins are loaded but no page matched → 404
// Use immediate + nextTick: immediate handles SPA navigation when plugins are already loaded,
// nextTick ensures all reactive updates (pluginPages registration) have flushed before checking.
if (import.meta.client) {
  watch(pluginsReady, async (ready) => {
    if (!ready) return
    await nextTick()
    if (!pageDef.value) {
      showError(createError({ statusCode: 404, statusMessage: 'Page not found' }))
    }
  }, { immediate: true })
}

const PluginComponent = shallowRef<ReturnType<typeof createPluginAsyncComponent> | null>(null)
let _lastKey = ''
watch([pluginId, componentName, moduleFile], ([pid, cname, mfile]) => {
  const key = `${pid}:${cname}:${mfile}`
  if (key === _lastKey) return
  _lastKey = key
  PluginComponent.value = pid ? createPluginAsyncComponent(pid, cname, mfile) : null
}, { immediate: true })
</script>

<template>
  <div class="min-h-screen">
    <ClientOnly>
      <component :is="PluginComponent" v-if="PluginComponent" />
      <!-- While plugins are loading, render nothing to avoid "not found" flash -->
    </ClientOnly>
  </div>
</template>
