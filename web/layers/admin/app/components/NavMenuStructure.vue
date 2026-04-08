<script setup lang="ts">
import type { UiMenuItem, MenuItemType } from '~/types/api/navMenu'

const props = defineProps<{
  items: UiMenuItem[]
}>()

const emit = defineEmits<{
  'update:items': [items: UiMenuItem[]]
}>()

const { t } = useI18n()

function getMenuIconName(type: MenuItemType) {
  return (
    {
      page: 'i-tabler-file',
      category: 'i-tabler-folder',
      custom: 'i-tabler-link',
      archive: 'i-tabler-archive',
    }[type] ?? 'i-tabler-circle'
  )
}

function typeLabel(type: MenuItemType) {
  return (
    {
      page: t('admin.appearance.menus.type_page'),
      category: t('admin.appearance.menus.type_category'),
      custom: t('admin.appearance.menus.type_custom'),
      archive: t('admin.appearance.menus.type_archive'),
    }[type] ?? type
  )
}

function indentItem(index: number) {
  const items = [...props.items]
  const item = items[index]
  const prev = items[index - 1]
  if (!item || !prev || item.depth >= 2 || item.depth > prev.depth) return
  items[index] = { ...item, depth: item.depth + 1, parent_local_id: prev.local_id }
  emit('update:items', items)
}

function outdentItem(index: number) {
  const items = [...props.items]
  const item = items[index]
  if (!item || item.depth <= 0) return
  items[index] = {
    ...item,
    depth: item.depth - 1,
    parent_local_id: items[index - 1]?.parent_local_id ?? '',
  }
  emit('update:items', items)
}

function deleteItem(index: number) {
  const item = props.items[index]
  if (!item) return
  emit('update:items', props.items.filter(
    (it, i) => i !== index && it.parent_local_id !== item.local_id,
  ))
}

function toggleExpanded(index: number) {
  const items = props.items.map((it, i) =>
    i === index ? { ...it, expanded: !it.expanded } : it,
  )
  emit('update:items', items)
}

function updateLabel(index: number, value: string) {
  emit('update:items', props.items.map((it, i) => i === index ? { ...it, label: value } : it))
}

function updateUrl(index: number, value: string) {
  emit('update:items', props.items.map((it, i) => i === index ? { ...it, url: value } : it))
}

function updateCssClasses(index: number, value: string) {
  emit('update:items', props.items.map((it, i) => i === index ? { ...it, cssClasses: value } : it))
}

function updateOpenInNewTab(index: number, value: boolean) {
  emit('update:items', props.items.map((it, i) => i === index ? { ...it, openInNewTab: value } : it))
}
</script>

<template>
  <UCard>
    <template #header>
      <div class="flex items-center justify-between">
        <h3 class="font-semibold text-highlighted">
          {{ $t('admin.appearance.menus.menu_structure') }}
        </h3>
        <span class="text-xs text-muted">{{ $t('admin.appearance.menus.indent_hint') }}</span>
      </div>
    </template>

    <div v-if="!items.length" class="text-center py-10 border-2 border-dashed border-default rounded-md">
      <UIcon name="i-tabler-menu-2" class="size-10 text-muted mb-2 mx-auto block" />
      <p class="text-sm text-muted">{{ $t('admin.appearance.menus.menu_empty') }}</p>
    </div>

    <div v-else class="space-y-2">
      <div
        v-for="(item, index) in items"
        :key="item.local_id"
        class="border border-default rounded-md bg-default"
        :style="{ marginLeft: item.depth * 28 + 'px' }">
        <div class="flex items-center gap-2 p-3">
          <UIcon name="i-tabler-grip-vertical" class="size-4 text-muted cursor-move shrink-0" />
          <UIcon :name="getMenuIconName(item.type)" class="size-4 text-primary shrink-0" />
          <span class="flex-1 text-sm font-medium text-highlighted truncate">{{ item.label }}</span>
          <UBadge :label="typeLabel(item.type)" color="neutral" variant="soft" size="xs" class="hidden sm:flex shrink-0" />
          <div class="flex items-center gap-0.5">
            <UButton
              v-if="item.depth < 2 && index > 0"
              color="neutral" variant="ghost" icon="i-tabler-arrow-bar-right" size="xs"
              :title="$t('admin.appearance.menus.indent')"
              @click="indentItem(index)" />
            <UButton
              v-if="item.depth > 0"
              color="neutral" variant="ghost" icon="i-tabler-arrow-bar-left" size="xs"
              :title="$t('admin.appearance.menus.outdent')"
              @click="outdentItem(index)" />
            <UButton
              color="neutral" variant="ghost" icon="i-tabler-pencil" size="xs"
              @click="toggleExpanded(index)" />
            <UButton
              color="error" variant="ghost" icon="i-tabler-trash" size="xs"
              @click="deleteItem(index)" />
          </div>
        </div>

        <div v-show="item.expanded" class="p-3 border-t border-default bg-elevated/30 space-y-2">
          <div class="grid grid-cols-2 gap-2">
            <UFormField :label="$t('admin.appearance.menus.nav_label')">
              <UInput :model-value="item.label" size="sm" class="w-full" @update:model-value="updateLabel(index, $event)" />
            </UFormField>
            <UFormField label="URL">
              <UInput :model-value="item.url" size="sm" class="w-full" @update:model-value="updateUrl(index, $event)" />
            </UFormField>
          </div>
          <UFormField :label="$t('admin.appearance.menus.css_classes')">
            <UInput :model-value="item.cssClasses" placeholder="my-class" size="sm" class="w-full" @update:model-value="updateCssClasses(index, $event)" />
          </UFormField>
          <div class="flex items-center gap-2">
            <UCheckbox :model-value="item.openInNewTab" @update:model-value="updateOpenInNewTab(index, $event)" />
            <span
              class="text-sm text-highlighted cursor-pointer select-none"
              @click="updateOpenInNewTab(index, !item.openInNewTab)">
              {{ $t('admin.appearance.menus.open_new_tab') }}
            </span>
          </div>
        </div>
      </div>
    </div>
  </UCard>
</template>
