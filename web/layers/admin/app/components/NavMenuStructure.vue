<script setup lang="ts">
import type { UiMenuItem, NavMenuSlotKey } from '~/types/api/navMenu'
import { NAV_MENU_SLOT_CONFIGS } from '~/types/api/navMenu'
import type { UiMenuTreeItem } from './NavMenuNode.vue'

const props = defineProps<{
  items: UiMenuItem[]
  slotKey: string
}>()

const emit = defineEmits<{
  'update:items': [items: UiMenuItem[]]
}>()

const slotConfig = computed(() => NAV_MENU_SLOT_CONFIGS[props.slotKey as NavMenuSlotKey])
const supportsNesting = computed(() => slotConfig.value?.supportsNesting ?? false)

// ── Flat ↔ Tree conversion ──────────────────────────────────────────────────

function flatToTree(items: UiMenuItem[]): UiMenuTreeItem[] {
  const map = new Map<string, UiMenuTreeItem>()
  const roots: UiMenuTreeItem[] = []

  for (const item of items) {
    map.set(item.local_id, {
      local_id: item.local_id,
      label: item.label,
      url: item.url,
      type: item.type,
      object_id: item.object_id,
      openInNewTab: item.openInNewTab,
      cssClasses: item.cssClasses,
      children: [],
    })
  }

  for (const item of items) {
    const node = map.get(item.local_id)!
    if (item.parent_local_id && map.has(item.parent_local_id)) {
      map.get(item.parent_local_id)!.children.push(node)
    } else {
      roots.push(node)
    }
  }

  return roots
}

function treeToFlat(tree: UiMenuTreeItem[], depth = 0, parentId = ''): UiMenuItem[] {
  const result: UiMenuItem[] = []
  for (const item of tree) {
    result.push({
      local_id: item.local_id,
      label: item.label,
      url: item.url,
      type: item.type,
      object_id: item.object_id,
      openInNewTab: item.openInNewTab,
      cssClasses: item.cssClasses,
      depth,
      parent_local_id: parentId,
      expanded: false,
    })
    if (item.children.length) {
      result.push(...treeToFlat(item.children, depth + 1, item.local_id))
    }
  }
  return result
}

// ── Tree state ──────────────────────────────────────────────────────────────

const tree = ref<UiMenuTreeItem[]>([])
let selfEmitting = false

watch(() => props.items, (items) => {
  if (selfEmitting) { selfEmitting = false; return }
  tree.value = flatToTree(items)
}, { immediate: true })

function emitFlat() {
  selfEmitting = true
  emit('update:items', treeToFlat(tree.value))
}

function onTreeUpdate(newTree: UiMenuTreeItem[]) {
  tree.value = newTree
  emitFlat()
}

function onUpdateChildren(parentId: string, children: UiMenuTreeItem[]) {
  const update = (nodes: UiMenuTreeItem[]): boolean => {
    for (const node of nodes) {
      if (node.local_id === parentId) {
        node.children = children
        return true
      }
      if (update(node.children)) return true
    }
    return false
  }
  update(tree.value)
  emitFlat()
}

function onUpdateField(localId: string, field: string, value: any) {
  const update = (nodes: UiMenuTreeItem[]): boolean => {
    for (const node of nodes) {
      if (node.local_id === localId) {
        ;(node as any)[field] = value
        return true
      }
      if (update(node.children)) return true
    }
    return false
  }
  update(tree.value)
  emitFlat()
}

function onDelete(localId: string) {
  const remove = (nodes: UiMenuTreeItem[]): UiMenuTreeItem[] =>
    nodes.filter(n => n.local_id !== localId).map(n => ({
      ...n,
      children: remove(n.children),
    }))
  tree.value = remove(tree.value)
  emitFlat()
}
</script>

<template>
  <UCard>
    <template #header>
      <div class="flex items-center justify-between">
        <h3 class="font-semibold text-highlighted">
          {{ $t('admin.appearance.menus.menu_structure') }}
        </h3>
        <span v-if="supportsNesting" class="text-xs text-muted">{{ $t('admin.appearance.menus.drag_nest_hint') }}</span>
      </div>
    </template>

    <div v-if="!tree.length" class="text-center py-10 border-2 border-dashed border-default rounded-md">
      <UIcon name="i-tabler-menu-2" class="size-10 text-muted mb-2 mx-auto block" />
      <p class="text-sm text-muted">{{ $t('admin.appearance.menus.menu_empty') }}</p>
    </div>

    <NavMenuNode
      v-else
      :model-value="tree"
      :supports-nesting="supportsNesting"
      @update:model-value="onTreeUpdate"
      @update-children="onUpdateChildren"
      @update-field="onUpdateField"
      @delete="onDelete"
    />
  </UCard>
</template>
