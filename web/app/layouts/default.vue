<template>
  <div class="min-h-screen flex flex-col">
    <BlogHeader />
    <main class="grow mx-auto w-full px-4 sm:px-6 lg:px-8 pb-8 overflow-x-clip">
      <slot />
    </main>
    <BlogFooter />
    <ClientOnly><BlogFloatingToolbar /></ClientOnly>
    <BlogBackToTop />
  </div>
</template>

<script setup lang="ts">
if (import.meta.client) {
  const { loadPlugins } = usePublicPluginLoader();
  onMounted(() => loadPlugins());
}
</script>
