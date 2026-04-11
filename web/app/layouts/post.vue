<template>
  <div class="min-h-screen flex flex-col">
    <BlogHeader />
    <main class="grow max-w-7xl mx-auto w-full px-4 sm:px-6 lg:px-8 py-8 grid grid-cols-1 lg:grid-cols-12 gap-8">
      <div class="lg:col-span-9">
        <slot />
      </div>
      <BlogSidebar />
    </main>
    <BlogFooter />
  </div>
</template>

<script setup lang="ts">
if (import.meta.client) {
  const { loadPlugins } = usePublicPluginLoader()
  onMounted(() => loadPlugins())
}
</script>
