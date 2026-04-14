<script setup lang="ts">
import type { UiMenuItem, NavMenuSlotKey } from '~/types/api/navMenu'

const props = defineProps<{
  items: UiMenuItem[]
  slotKey: string
}>()

const HEADER_ACTION_ICONS: Record<string, string> = {
  'action:lang_switcher': 'i-tabler-language',
  'action:theme_toggle': 'i-tabler-sun-moon',
  'action:messages': 'i-tabler-message',
  'action:notifications': 'i-tabler-bell',
}

const USER_MENU_ICONS: Record<string, string> = {
  'action:admin_dashboard': 'i-tabler-layout-dashboard',
  'action:user_profile': 'i-tabler-user',
  'action:user_posts': 'i-tabler-file-text',
  'action:user_favorites': 'i-tabler-bookmark',
  'action:user_settings': 'i-tabler-settings',
  'action:user_logout': 'i-tabler-logout',
}

const FLOATING_TOOLBAR_ICONS: Record<string, string> = {
  'action:profile_login': 'i-tabler-user-circle',
  'action:checkin': 'i-tabler-calendar-event',
  'action:ft_notifications': 'i-tabler-bell',
  'action:ft_theme_toggle': 'i-tabler-sun-moon',
}

const POST_ACTION_ICONS: Record<string, string> = {
  'action:post_like': 'i-tabler-heart',
  'action:post_comment': 'i-tabler-message-circle',
  'action:post_bookmark': 'i-tabler-bookmark',
  'action:post_share': 'i-tabler-share-2',
}

const isHeaderActions = computed(() => props.slotKey === 'header_actions')
const isUserMenu = computed(() => props.slotKey === 'user_menu')
const isFloatingToolbar = computed(() => props.slotKey === 'floating_toolbar')
const isPostActions = computed(() => props.slotKey === 'post_actions')
const isVerticalSlot = computed(() => isUserMenu.value || isFloatingToolbar.value || isPostActions.value)

const rootItems = computed(() => props.items.filter(i => i.depth === 0))

function getChildren(parentId: string): UiMenuItem[] {
  return props.items.filter(i => i.parent_local_id === parentId)
}

function getIconForSlot(item: UiMenuItem): string {
  if (item.type === 'action') {
    if (isHeaderActions.value) return HEADER_ACTION_ICONS[item.local_id] || 'i-tabler-click'
    if (isUserMenu.value) return USER_MENU_ICONS[item.local_id] || 'i-tabler-click'
    if (isFloatingToolbar.value) return FLOATING_TOOLBAR_ICONS[item.local_id] || 'i-tabler-click'
    if (isPostActions.value) return POST_ACTION_ICONS[item.local_id] || 'i-tabler-click'
  }
  return item.cssClasses || 'i-tabler-link'
}
</script>

<template>
  <UCard>
    <template #header>
      <h3 class="font-semibold text-highlighted">{{ $t('admin.appearance.menus.menu_preview') }}</h3>
    </template>

    <!-- Header actions preview -->
    <nav v-if="isHeaderActions" class="flex flex-wrap items-center gap-2 min-h-8 bg-elevated/40 rounded-md p-4">
      <template v-for="item in rootItems" :key="item.local_id">
        <div v-if="item.type === 'separator'" class="h-6 w-px bg-border-default border border-default" />
        <UTooltip v-else :text="item.label">
          <div class="size-8 rounded-md flex items-center justify-center bg-default border border-default hover:bg-elevated transition-colors">
            <UIcon :name="getIconForSlot(item)" class="size-4 text-muted" />
          </div>
        </UTooltip>
      </template>
      <span v-if="!rootItems.length" class="text-sm text-muted">{{ $t('admin.appearance.menus.empty_preview') }}</span>
    </nav>

    <!-- Vertical slot preview (user_menu, floating_toolbar, post_actions) -->
    <div v-else-if="isVerticalSlot" class="min-h-8 bg-elevated/40 rounded-md p-4">
      <div class="inline-flex flex-col gap-0.5 rounded-md bg-default border border-default p-1.5 min-w-48">
        <template v-for="item in rootItems" :key="item.local_id">
          <div v-if="item.type === 'separator'" class="h-px bg-border-default border border-default my-1 mx-2" />
          <div v-else class="flex items-center gap-2.5 px-3 py-1.5 rounded hover:bg-elevated transition-colors cursor-default">
            <UIcon :name="getIconForSlot(item)" class="size-4 text-muted shrink-0" />
            <span class="text-sm text-highlighted">{{ item.label }}</span>
          </div>
        </template>
        <span v-if="!rootItems.length" class="text-sm text-muted px-3 py-2">{{ $t('admin.appearance.menus.empty_preview') }}</span>
      </div>
    </div>

    <!-- Standard menu preview -->
    <nav v-else class="flex flex-wrap gap-6 min-h-8 bg-elevated/40 rounded-md p-4">
      <template v-for="item in rootItems" :key="item.local_id">
        <div v-if="item.type === 'separator'" class="h-6 w-px bg-border-default border border-default self-center" />
        <div v-else class="relative group">
          <a
            href="#"
            class="text-sm font-medium text-highlighted hover:text-primary flex items-center gap-1"
            @click.prevent>
            {{ item.label }}
            <UIcon v-if="getChildren(item.local_id).length" name="i-tabler-chevron-down" class="size-3 text-muted" />
          </a>
          <div
            v-if="getChildren(item.local_id).length"
            class="absolute top-full left-0 mt-1 bg-default border border-default rounded-md shadow-lg py-1.5 min-w-36 opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all z-10">
            <template v-for="child in getChildren(item.local_id)" :key="child.local_id">
              <div v-if="child.type === 'separator'" class="h-px bg-border-default border border-default my-1 mx-2" />
              <a v-else href="#" class="block px-3 py-1.5 text-sm text-highlighted hover:bg-elevated" @click.prevent>{{ child.label }}</a>
            </template>
          </div>
        </div>
      </template>
      <span v-if="!rootItems.length" class="text-sm text-muted">{{ $t('admin.appearance.menus.empty_preview') }}</span>
    </nav>
  </UCard>
</template>
