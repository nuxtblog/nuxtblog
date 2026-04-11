<script setup lang="ts">
/**
 * P3-2: Public plugin page renderer.
 *
 * Catches routes under /p/{pluginId}/{componentName} and renders
 * the corresponding plugin Vue component loaded from public.mjs.
 */
import { createPluginAsyncComponent } from '~/composables/usePluginComponents'

const route = useRoute()
const contributionsStore = usePluginContributionsStore()

const slugParts = computed(() => {
  const slug = route.params.slug
  return Array.isArray(slug) ? slug : [slug]
})

const pluginId = computed(() => slugParts.value[0] || '')
const componentName = computed(() => slugParts.value[1] || 'default')

// Look up the registered page to determine module file
const pageDef = computed(() =>
  contributionsStore.getPluginPage(pluginId.value, componentName.value),
)

// Fallback to public.mjs for public plugin pages
const moduleFile = computed(() => pageDef.value?.moduleFile || 'public.mjs')
const pageTitle = computed(() => pageDef.value?.title || '')

// Set page title reactively
useHead(() => ({
  title: pageTitle.value || undefined,
}))

// Use shallowRef + watch to avoid creating a new async component on every computed evaluation
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
      <div v-if="!PluginComponent" class="flex items-center justify-center h-64 text-gray-400">
        Plugin page not found
      </div>
      <component :is="PluginComponent" v-else />
    </ClientOnly>
  </div>
</template>
