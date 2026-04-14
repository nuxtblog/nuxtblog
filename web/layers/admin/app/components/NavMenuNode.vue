<template>
  <VueDraggable
    class="nav-menu-drag-area"
    tag="ul"
    :model-value="modelValue"
    :group="groupConfig"
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <li v-for="item in modelValue" :key="item.local_id">
      <!-- Node content -->
      <div
        class="group rounded-md border px-3 py-2 flex items-center gap-2 transition-all"
        :class="item.type === 'separator'
          ? 'border-dashed border-muted bg-elevated/20'
          : 'border-default bg-default hover:border-primary/40 hover:shadow-sm'"
      >
        <!-- Drag handle -->
        <UIcon
          name="i-tabler-grip-vertical"
          class="size-4 shrink-0 cursor-grab text-muted hover:text-primary active:cursor-grabbing"
        />

        <!-- Expand/collapse (nesting slots, non-separator items) -->
        <button
          v-if="supportsNesting && item.type !== 'separator'"
          class="flex h-6 w-6 items-center justify-center shrink-0 rounded hover:bg-elevated transition-colors"
          @click="toggleExpanded(item.local_id)"
        >
          <UIcon
            name="i-tabler-chevron-right"
            class="size-4 transition-transform duration-200"
            :class="[
              expandedSet.has(item.local_id) ? 'rotate-90' : '',
              item.children?.length ? 'text-muted' : 'text-transparent',
            ]"
          />
        </button>

        <!-- Type icon -->
        <UIcon :name="getMenuIconName(item.type)" class="size-4 shrink-0" :class="item.type === 'separator' ? 'text-muted' : 'text-primary'" />

        <!-- Label -->
        <span v-if="item.type === 'separator'" class="flex-1 text-sm text-muted italic">—— {{ $t('admin.appearance.menus.type_separator') }} ——</span>
        <span v-else class="flex-1 text-sm font-medium text-highlighted truncate">{{ item.label }}</span>

        <!-- Badges -->
        <UBadge v-if="item.type !== 'separator'" :label="typeLabel(item.type)" color="neutral" variant="soft" size="xs" class="hidden sm:flex shrink-0" />
        <UBadge v-if="item.local_id.startsWith('plugin:')" label="插件" color="info" variant="soft" size="xs" class="shrink-0" />

        <!-- Actions -->
        <div class="flex items-center gap-0.5 opacity-0 group-hover:opacity-100 transition-opacity shrink-0">
          <UButton
            v-if="item.type !== 'separator'"
            color="neutral" variant="ghost" icon="i-tabler-pencil" size="xs"
            @click="toggleEditing(item.local_id)" />
          <UButton
            color="error" variant="ghost" icon="i-tabler-trash" size="xs"
            @click="$emit('delete', item.local_id)" />
        </div>
      </div>

      <!-- Edit panel -->
      <div v-if="item.type !== 'separator' && editingSet.has(item.local_id)" class="ml-6 mt-1 mb-2 p-3 border border-default rounded-md bg-elevated/30 space-y-2">
        <div :class="item.type === 'action' ? '' : 'grid grid-cols-2 gap-2'">
          <UFormField :label="$t('admin.appearance.menus.nav_label')">
            <UInput :model-value="item.label" size="sm" class="w-full" @update:model-value="$emit('update-field', item.local_id, 'label', $event)" />
          </UFormField>
          <UFormField v-if="item.type !== 'action'" label="URL">
            <UInput :model-value="item.url" size="sm" class="w-full" @update:model-value="$emit('update-field', item.local_id, 'url', $event)" />
          </UFormField>
        </div>
        <UFormField :label="item.type === 'action' ? $t('admin.appearance.menus.action_icon_label') : $t('admin.appearance.menus.css_classes')">
          <UInput :model-value="item.cssClasses" :placeholder="item.type === 'action' ? 'i-tabler-brand-github' : 'my-class'" size="sm" class="w-full" @update:model-value="$emit('update-field', item.local_id, 'cssClasses', $event)" />
        </UFormField>
        <div class="flex items-center gap-2">
          <UCheckbox :model-value="item.openInNewTab" @update:model-value="$emit('update-field', item.local_id, 'openInNewTab', $event)" />
          <span
            class="text-sm text-highlighted cursor-pointer select-none"
            @click="$emit('update-field', item.local_id, 'openInNewTab', !item.openInNewTab)">
            {{ $t('admin.appearance.menus.open_new_tab') }}
          </span>
        </div>
      </div>

      <!-- Children (recursive) — always render drop zone for nesting-enabled non-separator items -->
      <NavMenuNode
        v-if="supportsNesting && item.type !== 'separator' && expandedSet.has(item.local_id)"
        :model-value="item.children ?? []"
        :supports-nesting="supportsNesting"
        :max-depth="maxDepth"
        :current-depth="currentDepth + 1"
        @update:model-value="$emit('update-children', item.local_id, $event)"
        @delete="(id: string) => $emit('delete', id)"
        @update-children="(id: string, children: UiMenuTreeItem[]) => $emit('update-children', id, children)"
        @update-field="(id: string, field: string, val: any) => $emit('update-field', id, field, val)"
      />
    </li>
  </VueDraggable>
</template>

<script setup lang="ts">
import { VueDraggable } from 'vue-draggable-plus'
import type { MenuItemType } from '~/types/api/navMenu'

export interface UiMenuTreeItem {
  local_id: string
  label: string
  url: string
  type: MenuItemType
  object_id: number
  openInNewTab: boolean
  cssClasses: string
  children: UiMenuTreeItem[]
}

const props = withDefaults(defineProps<{
  modelValue: UiMenuTreeItem[]
  supportsNesting: boolean
  maxDepth?: number
  currentDepth?: number
}>(), {
  maxDepth: 2,
  currentDepth: 0,
})

defineEmits<{
  'update:modelValue': [value: UiMenuTreeItem[]]
  'update-children': [localId: string, children: UiMenuTreeItem[]]
  'update-field': [localId: string, field: string, value: any]
  'delete': [localId: string]
}>()

const { t } = useI18n()

const groupConfig = computed(() => {
  if (!props.supportsNesting) return { name: 'nav-menu' }
  // Allow nesting up to maxDepth
  const canPut = props.currentDepth < props.maxDepth
  return { name: 'nav-menu', pull: true, put: canPut }
})

// Expanded nodes (show children)
const expandedSet = ref<Set<string>>(new Set())

const initExpanded = () => {
  const s = new Set<string>()
  const traverse = (items: UiMenuTreeItem[]) => {
    for (const item of items) {
      // Expand ALL items so their child drop zone is always rendered
      s.add(item.local_id)
      if (item.children?.length) traverse(item.children)
    }
  }
  traverse(props.modelValue)
  expandedSet.value = s
}

watch(() => props.modelValue, initExpanded, { immediate: true, deep: true })

function toggleExpanded(id: string) {
  if (expandedSet.value.has(id)) expandedSet.value.delete(id)
  else expandedSet.value.add(id)
}

// Editing nodes (show edit panel)
const editingSet = ref<Set<string>>(new Set())

function toggleEditing(id: string) {
  if (editingSet.value.has(id)) editingSet.value.delete(id)
  else editingSet.value.add(id)
}

function getMenuIconName(type: MenuItemType) {
  return ({
    page: 'i-tabler-file',
    category: 'i-tabler-folder',
    custom: 'i-tabler-link',
    archive: 'i-tabler-archive',
    action: 'i-tabler-click',
    separator: 'i-tabler-separator',
  }[type] ?? 'i-tabler-circle')
}

function typeLabel(type: MenuItemType) {
  return ({
    page: t('admin.appearance.menus.type_page'),
    category: t('admin.appearance.menus.type_category'),
    custom: t('admin.appearance.menus.type_custom'),
    archive: t('admin.appearance.menus.type_archive'),
    action: t('admin.appearance.menus.type_action'),
    separator: t('admin.appearance.menus.type_separator'),
  }[type] ?? type)
}
</script>

<style scoped>
.nav-menu-drag-area {
  list-style: none;
  margin: 0;
  padding: 0;
}

.nav-menu-drag-area > li {
  margin-bottom: 2px;
}

/* Nested levels get indentation with dashed border */
.nav-menu-drag-area .nav-menu-drag-area {
  border-left: 2px dashed var(--ui-border);
  margin-left: 16px;
  padding-left: 8px;
  margin-top: 2px;
  min-height: 4px;
}
</style>
