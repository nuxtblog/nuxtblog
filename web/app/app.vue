<script setup lang="ts">
const route = useRoute()
const pluginApi = usePluginApi()

const { data: stylesData } = await useAsyncData(
  'plugin-styles',
  () => pluginApi.getStyles().catch(() => ({ css: '' })),
)

// Derive page type from route path for CSS targeting
const pageType = computed(() => {
  const p = route.path
  if (p === '/') return 'home'
  if (p.startsWith('/admin')) return 'admin'
  if (p.match(/^\/posts\/.+/)) return 'post'
  if (p.match(/^\/docs\/[^/]+\/[^/]+/)) return 'doc'
  if (p.match(/^\/docs\/[^/]+/)) return 'collection'
  if (p === '/docs') return 'docs'
  if (p.match(/^\/pages\/.+/)) return 'page'
  if (p === '/moments') return 'moments'
  if (p === '/archive') return 'archive'
  if (p.startsWith('/categories')) return 'category'
  if (p.startsWith('/tags')) return 'tags'
  if (p.startsWith('/user')) return 'user'
  if (p.startsWith('/auth')) return 'auth'
  return 'default'
})

useHead(computed(() => ({
  style: stylesData.value?.css ? [{ children: stylesData.value.css }] : [],
  bodyAttrs: { 'data-page': pageType.value },
})))
</script>

<template>
  <UApp :toaster="{ duration: 8000 }">
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
  </UApp>
</template>
