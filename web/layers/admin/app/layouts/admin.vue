<template>
  <div class="h-dvh flex flex-col overflow-hidden bg-default" data-admin-layout>
    <!-- Header -->
    <AdminHeader class="shrink-0 h-14 md:h-16 z-40 border-b border-default" />

    <div class="flex-1 flex overflow-hidden min-h-0">
      <!-- 移动端遮罩 -->
      <Transition name="fade">
        <div
          v-if="mobileOpen"
          class="fixed inset-0 bg-black/50 z-30 lg:hidden"
          @click="closeMobile"
        />
      </Transition>

      <!-- 侧边栏：移动端=固定抽屉，桌面端=内联 -->
      <div
        class="
          fixed top-14 md:top-16 bottom-0 left-0 z-40
          lg:static lg:top-auto lg:bottom-auto lg:z-auto
          transition-transform duration-300 ease-in-out
          lg:transition-none
        "
        :class="mobileOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'"
      >
        <AdminSidebar class="h-full" />
      </div>

      <!-- 主内容区 -->
      <main class="flex-1 min-w-0 h-full overflow-hidden flex flex-col">
        <slot />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
const { mobileOpen, closeMobile } = useAdminSidebar();

// Phase 2.4: install nuxtblogAdmin global + load plugin scripts
if (import.meta.client) {
  const { installNuxtblogAdmin } = await import('~/composables/useNuxtblogAdmin');
  const { usePluginLoader } = await import('~/composables/usePluginLoader');

  installNuxtblogAdmin();
  const { loadPlugins } = usePluginLoader();
  onMounted(() => { loadPlugins(); });
}
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.25s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
