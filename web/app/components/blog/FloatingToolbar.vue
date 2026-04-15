<script setup lang="ts">
import type { NavMenuItem } from '~/types/api/navMenu'
import type { Component } from 'vue'
import { executePublicCommand } from '~/composables/useNuxtblogPublic'
import { createPluginAsyncComponent } from '~/composables/usePluginComponents'

const { t } = useI18n();
const authStore = useAuthStore();
const checkinApi = useCheckinApi();
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

// ── 签到状态 ──────────────────────────────────────────────────────────────────
const hasCheckedIn = ref(false);
const streak = ref(0);
const checkinAnim = ref(false);

// 登录后拉取今日签到状态
watch(
  () => authStore.isLoggedIn,
  async (loggedIn) => {
    if (!loggedIn) {
      hasCheckedIn.value = false;
      streak.value = 0;
      return;
    }
    try {
      const res = await checkinApi.getStatus();
      hasCheckedIn.value = res.checked_in_today;
      streak.value = res.streak;
    } catch {}
  },
  { immediate: true },
);

const doCheckIn = async () => {
  if (!authStore.isLoggedIn) {
    navigateTo("/auth/login");
    return;
  }
  if (hasCheckedIn.value) return;
  try {
    const res = await checkinApi.doCheckin();
    hasCheckedIn.value = true;
    streak.value = res.streak;
    checkinAnim.value = true;
    setTimeout(() => (checkinAnim.value = false), 2500);
  } catch {}
};

// Read menu config (array order = display order)
const toolbarItems = computed(() => getOption('floating_toolbar') as NavMenuItem[])

// Set of visible action IDs for quick lookup
const visibleActions = computed(() => new Set(toolbarItems.value.map(i => i.local_id)))
</script>

<template>
  <div
    class="fixed right-5 md:right-7 bottom-2 -translate-y-1/2 z-50 flex flex-col items-center">
    <!-- 签到成功气泡 -->
    <Transition
      enter-active-class="transition duration-300 ease-out"
      enter-from-class="opacity-0 translate-x-3 scale-95"
      enter-to-class="opacity-100 translate-x-0 scale-100"
      leave-active-class="transition duration-200 ease-in"
      leave-from-class="opacity-100 translate-x-0"
      leave-to-class="opacity-0 translate-x-3 scale-95">
      <div
        v-if="checkinAnim"
        class="absolute right-14 w-40 rounded-md bg-default shadow-xl ring-1 ring-default px-4 py-3 text-center pointer-events-none">
        <UIcon
          name="i-tabler-confetti"
          class="size-6 text-yellow-500 mx-auto mb-1" />
        <p class="text-sm font-bold text-highlighted">
          {{ $t("site.floating.checkin_success") }}
        </p>
        <p class="text-xs text-muted mt-0.5">
          {{ $t("site.floating.checkin_streak", { n: streak }) }}
        </p>
        <div
          class="absolute right-[-6px] top-1/2 -translate-y-1/2 size-3 rotate-45 bg-default ring-1 ring-default clip-arrow" />
      </div>
    </Transition>

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

        <!-- 今日签到 -->
        <UTooltip
          v-else-if="item.local_id === 'action:checkin'"
          :text="
            !authStore.isLoggedIn
              ? t('site.floating.checkin_login')
              : hasCheckedIn
                ? t('site.floating.checkin_done', { n: streak })
                : t('site.floating.checkin_todo')
          "
          side="left"
          :delay-duration="100">
          <button
            class="group relative flex items-center justify-center size-10 rounded-md transition-all duration-200"
            :class="
              hasCheckedIn
                ? 'bg-success/10 cursor-default'
                : 'hover:bg-amber-500/10 cursor-pointer'
            "
            @click="doCheckIn">
            <UIcon
              :name="
                hasCheckedIn
                  ? 'i-tabler-calendar-check'
                  : 'i-tabler-calendar-event'
              "
              class="size-5 transition-colors"
              :class="
                hasCheckedIn
                  ? 'text-success'
                  : 'text-muted group-hover:text-amber-500'
              " />
            <!-- pulse ring when not checked in and logged in -->
            <span
              v-if="authStore.isLoggedIn && !hasCheckedIn"
              class="absolute inset-0 rounded-md animate-ping bg-amber-400/20 pointer-events-none" />
            <!-- streak badge -->
            <span
              v-if="hasCheckedIn && streak > 1"
              class="absolute -top-0.5 -right-0.5 min-w-4 h-4 rounded-full bg-success text-white text-[10px] font-bold flex items-center justify-center px-0.5">
              {{ streak }}
            </span>
          </button>
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
