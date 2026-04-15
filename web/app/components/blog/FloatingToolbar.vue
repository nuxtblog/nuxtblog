<script setup lang="ts">
import type { NavMenuItem } from '~/types/api/navMenu'
import type { Component } from 'vue'
import { executePublicCommand } from '~/composables/useNuxtblogPublic'
import { createPluginAsyncComponent } from '~/composables/usePluginComponents'

const authStore = useAuthStore();
const { getOption } = useOption();
const contributionsStore = usePluginContributionsStore();

// ── 插件 Popover 逻辑 ─────────────────────────────────────────────────────────
const pluginViews = contributionsStore.getViewItems('public:floating-toolbar')

function getPluginView(localId: string) {
  const pluginId = localId.replace('plugin:', '')
  return pluginViews.value.find(v => v.id === pluginId && v.component && v.module)
}

const viewComponentCache = new Map<string, Component>()

function resolvePluginComponent(localId: string) {
  const view = getPluginView(localId)
  if (!view) return null
  const key = `${view.pluginId}:${view.component}:${view.module}`
  if (!viewComponentCache.has(key)) {
    viewComponentCache.set(key, createPluginAsyncComponent(view.pluginId, view.component!, view.module!))
  }
  return viewComponentCache.get(key)!
}

// Hover-controlled popover state
const openPopoverId = ref<string | null>(null)
let hoverTimer: ReturnType<typeof setTimeout> | null = null

function onPluginEnter(id: string) {
  if (hoverTimer) clearTimeout(hoverTimer)
  openPopoverId.value = id
}

function onPluginLeave() {
  hoverTimer = setTimeout(() => { openPopoverId.value = null }, 200)
}

// Read menu config (array order = display order)
const toolbarItems = computed(() => getOption('floating_toolbar') as NavMenuItem[])

// Set of visible action IDs for quick lookup
const visibleActions = computed(() => new Set(toolbarItems.value.map(i => i.local_id)))
</script>

<template>
  <div
    class="fixed right-5 md:right-7 bottom-2 -translate-y-1/2 z-50 flex flex-col items-center">
    <!-- 工具栏 -->
    <div
      class="flex flex-col items-center gap-0.5 rounded-md bg-default/90 backdrop-blur-xl shadow-2xl shadow-black/10 ring-1 ring-default p-1.5">
      <template v-for="item in toolbarItems" :key="item.local_id">
        <!-- 分隔符 -->
        <div v-if="item.object_type === 'separator'" class="w-6 h-px bg-border-default border border-default my-1 rounded-full" />

        <!-- 个人中心 -->
        <UTooltip
          v-else-if="item.local_id === 'action:profile_login'"
          :text="$t('site.floating.profile_center')"
          side="left"
          :delay-duration="100">
          <NuxtLink
            :to="authStore.isLoggedIn ? '/user/profile' : '/auth/login'"
            class="group flex items-center justify-center size-10 rounded-md hover:bg-primary/10 transition-all duration-200">
            <BaseAvatar
              v-if="authStore.isLoggedIn"
              :src="authStore.user?.avatar"
              :alt="authStore.user?.display_name || authStore.user?.username"
              size="xs"
              class="ring-1 ring-default group-hover:ring-primary transition-all" />
            <UIcon
              v-else
              name="i-tabler-user-circle"
              class="size-5 text-muted group-hover:text-primary transition-colors" />
          </NuxtLink>
        </UTooltip>

        <!-- 消息通知 -->
        <UTooltip
          v-else-if="item.local_id === 'action:ft_notifications'"
          :text="$t('site.floating.notifications')"
          side="left"
          :delay-duration="100">
          <NuxtLink
            to="/notifications"
            class="group relative flex items-center justify-center size-10 rounded-md hover:bg-primary/10 transition-all duration-200">
            <UIcon
              name="i-tabler-bell"
              class="size-5 text-muted group-hover:text-primary transition-colors" />
          </NuxtLink>
        </UTooltip>

        <!-- 主题切换 -->
        <UTooltip
          v-else-if="item.local_id === 'action:ft_theme_toggle'"
          :text="$t('site.floating.theme_toggle')"
          side="left"
          :delay-duration="100">
          <UColorModeButton
            class="flex items-center justify-center size-10 rounded-md hover:bg-primary/10 transition-all duration-200 text-muted hover:text-primary" />
        </UTooltip>

        <!-- 插件项：Popover 模式（有 component 的插件） -->
        <UPopover
          v-else-if="item.local_id.startsWith('plugin:') && getPluginView(item.local_id)"
          :open="openPopoverId === item.local_id"
          side="left"
          :ui="{ content: 'p-0 w-72' }"
          @mouseenter="onPluginEnter(item.local_id)"
          @mouseleave="onPluginLeave"
        >
          <button class="group flex items-center justify-center size-10 rounded-md hover:bg-primary/10 transition-all duration-200">
            <UIcon :name="item.css_classes || 'i-tabler-puzzle'" class="size-5 text-muted group-hover:text-primary transition-colors" />
          </button>
          <template #content>
            <div @mouseenter="onPluginEnter(item.local_id)" @mouseleave="onPluginLeave">
              <ClientOnly>
                <component :is="resolvePluginComponent(item.local_id)" />
              </ClientOnly>
            </div>
          </template>
        </UPopover>

        <!-- 插件项：按钮模式（无 component 的命令型插件） -->
        <UTooltip
          v-else-if="item.local_id.startsWith('plugin:')"
          :text="item.label"
          side="left"
          :delay-duration="100">
          <button
            class="group flex items-center justify-center size-10 rounded-md hover:bg-primary/10 transition-all duration-200"
            @click="executePublicCommand(item.local_id.replace('plugin:', ''))">
            <UIcon
              :name="item.css_classes || 'i-tabler-puzzle'"
              class="size-5 text-muted group-hover:text-primary transition-colors" />
          </button>
        </UTooltip>

        <!-- 自定义链接 -->
        <UTooltip
          v-else-if="item.object_type === 'custom'"
          :text="item.label"
          side="left"
          :delay-duration="100">
          <NuxtLink
            :to="item.url"
            :target="item.target || undefined"
            class="group flex items-center justify-center size-10 rounded-md hover:bg-primary/10 transition-all duration-200">
            <UIcon
              :name="item.css_classes || 'i-tabler-link'"
              class="size-5 text-muted group-hover:text-primary transition-colors" />
          </NuxtLink>
        </UTooltip>
      </template>
    </div>
  </div>
</template>
