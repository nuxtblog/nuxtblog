<script setup lang="ts">
import type { UiMenuItem } from '~/types/api/navMenu'

const props = defineProps<{
  items: UiMenuItem[]
  isHeaderActions: boolean
}>()

const HEADER_ACTION_ICONS: Record<string, string> = {
  'action:lang_switcher': 'i-tabler-language',
  'action:theme_toggle': 'i-tabler-sun-moon',
  'action:messages': 'i-tabler-message',
  'action:notifications': 'i-tabler-bell',
}

const rootItems = computed(() => props.items.filter(i => i.depth === 0))

function getChildren(parentId: string): UiMenuItem[] {
  return props.items.filter(i => i.parent_local_id === parentId)
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
        <UTooltip :text="item.label">
          <div class="size-8 rounded-md flex items-center justify-center bg-default border border-default hover:bg-elevated transition-colors">
            <UIcon
              :name="item.type === 'action' ? (HEADER_ACTION_ICONS[item.local_id] || 'i-tabler-click') : (item.cssClasses || 'i-tabler-link')"
              class="size-4 text-muted"
            />
          </div>
        </UTooltip>
      </template>
      <span v-if="!rootItems.length" class="text-sm text-muted">
        {{ $t('admin.appearance.menus.empty_preview') }}
      </span>
    </nav>
    <!-- Standard menu preview -->
    <nav v-else class="flex flex-wrap gap-6 min-h-8 bg-elevated/40 rounded-md p-4">
      <template v-for="item in rootItems" :key="item.local_id">
        <div class="relative group">
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
            <a
              v-for="child in getChildren(item.local_id)"
              :key="child.local_id"
              href="#"
              class="block px-3 py-1.5 text-sm text-highlighted hover:bg-elevated"
              @click.prevent>{{ child.label }}</a>
          </div>
        </div>
      </template>
      <span v-if="!rootItems.length" class="text-sm text-muted">
        {{ $t('admin.appearance.menus.empty_preview') }}
      </span>
    </nav>
  </UCard>
</template>
