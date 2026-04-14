<script setup lang="ts">
/**
 * Phase 2.3: ContributionSlot — renders plugin UI extensions in named anchor points.
 *
 * Usage:
 *   <ContributionSlot name="post-editor:toolbar" />
 *   <ContributionSlot name="admin:sidebar-nav" />
 *   <ContributionSlot name="post-editor:context" :ctx="editorCtx" />
 */
import DOMPurify from 'dompurify'
import type { ContributionSlotName } from '~/config/contribution-slots'
import { usePluginContextStore } from '~/stores/plugin-context'
import { usePluginContributionsStore, type PluginContentBlock, type PluginMenuItem } from '~/stores/plugin-contributions'

export interface MenuGroup {
  group: string
  items: PluginMenuItem[]
}

export interface GroupedEntry {
  /** Present for ungrouped single items */
  item?: PluginMenuItem
  /** Present for grouped items */
  group?: MenuGroup
}

defineOptions({ inheritAttrs: false })

const props = defineProps<{
  /** Slot name matching manifest menus/views keys */
  name: ContributionSlotName
  /** Optional context passed to command handlers */
  ctx?: Record<string, unknown>
}>()

const emit = defineEmits<{
  command: [commandId: string, ctx?: Record<string, unknown>]
}>()

const contextStore = usePluginContextStore()
const contributionsStore = usePluginContributionsStore()

const menuItems = contributionsStore.getMenuItems(props.name)
const viewItems = contributionsStore.getViewItems(props.name)
const navItems = contributionsStore.getNavigation(props.name)
const contentBlocks = contributionsStore.getContentBlocks(props.name)

/** Sanitize HTML content for safe rendering via DOMPurify. */
function sanitizeHtml(html: string): string {
  if (import.meta.server) return ''
  return DOMPurify.sanitize(html)
}

/** Filter menu items by their `when` condition. */
const visibleMenuItems = computed(() =>
  menuItems.value.filter(item =>
    !item.when || contextStore.evaluateWhen(item.when),
  ),
)

/** Group visible menu items: ungrouped items stay solo, grouped items merge into MenuGroup entries. */
const groupedMenuEntries = computed<GroupedEntry[]>(() => {
  const singles: GroupedEntry[] = []
  const groupMap = new Map<string, PluginMenuItem[]>()
  const groupOrder: string[] = []

  for (const item of visibleMenuItems.value) {
    if (!item.group) {
      singles.push({ item })
    } else {
      if (!groupMap.has(item.group)) {
        groupMap.set(item.group, [])
        groupOrder.push(item.group)
      }
      groupMap.get(item.group)!.push(item)
    }
  }

  // Insert group entries at the position of first occurrence
  const result: GroupedEntry[] = []
  let groupIdx = 0
  for (const entry of singles) {
    // Insert any groups that appeared before this ungrouped item
    while (groupIdx < groupOrder.length) {
      const gName = groupOrder[groupIdx]!
      const gItems = groupMap.get(gName)!
      const gFirstIdx = visibleMenuItems.value.indexOf(gItems[0]!)
      const curIdx = visibleMenuItems.value.indexOf(entry.item!)
      if (gFirstIdx < curIdx) {
        result.push({ group: { group: gName, items: gItems } })
        groupIdx++
      } else {
        break
      }
    }
    result.push(entry)
  }
  // Append remaining groups
  while (groupIdx < groupOrder.length) {
    const gName = groupOrder[groupIdx]!
    result.push({ group: { group: gName, items: groupMap.get(gName)! } })
    groupIdx++
  }
  return result
})

function handleCommand(commandId: string) {
  emit('command', commandId, props.ctx)
}

// ── renderFn support: mount DOM render functions when blocks appear ──────
const renderFnRefs = ref<Map<string, HTMLElement>>(new Map())
const cleanupFns = ref<Map<string, () => void>>(new Map())

function setRenderFnRef(block: PluginContentBlock, el: HTMLElement | null) {
  const key = `${block.pluginId}-${block.slot}`
  if (el && block.renderFn) {
    renderFnRefs.value.set(key, el)
    // Execute render function
    nextTick(() => {
      const cleanup = block.renderFn!(el)
      if (typeof cleanup === 'function') {
        cleanupFns.value.set(key, cleanup)
      }
    })
  }
}

onBeforeUnmount(() => {
  for (const cleanup of cleanupFns.value.values()) {
    try { cleanup() } catch {}
  }
  cleanupFns.value.clear()
})
</script>

<template>
  <!-- Navigation items (sidebar, topbar) -->
  <template v-for="nav in navItems" :key="`nav-${nav.pluginId}-${nav.route}`">
    <slot name="nav" :item="nav">
      <NuxtLink :to="nav.route" class="flex items-center gap-2 px-3 py-2 text-sm rounded-md hover:bg-default/50">
        <UIcon v-if="nav.icon" :name="nav.icon" class="size-4" />
        <span>{{ nav.title }}</span>
      </NuxtLink>
    </slot>
  </template>

  <!-- Menu items (toolbar buttons, context menu entries) -->
  <template v-for="(entry, idx) in groupedMenuEntries" :key="entry.item ? `menu-${entry.item.command}` : `group-${entry.group!.group}-${idx}`">
    <!-- Ungrouped: single button -->
    <template v-if="entry.item">
      <slot name="menu" :item="entry.item" :execute="() => handleCommand(entry.item!.command)">
        <UButton
          variant="ghost"
          size="xs"
          :icon="entry.item.icon"
          :label="entry.item.title"
          @click="handleCommand(entry.item.command)" />
      </slot>
    </template>

    <!-- Grouped: split button with dropdown -->
    <template v-else-if="entry.group">
      <slot name="menu-group" :group="entry.group" :execute="(cmd: string) => handleCommand(cmd)">
        <div class="inline-flex items-center">
          <!-- Primary button: first command in group -->
          <UTooltip :text="entry.group!.items[0]!.title ?? ''">
            <UButton
              variant="ghost"
              color="neutral"
              size="xs"
              :icon="entry.group!.items[0]!.icon"
              class="rounded-r-none"
              @click="handleCommand(entry.group!.items[0]!.command)" />
          </UTooltip>
          <!-- Dropdown arrow for all commands -->
          <UPopover :ui="{ content: 'p-1' }">
            <UButton
              variant="ghost"
              color="neutral"
              size="xs"
              icon="i-tabler-chevron-down"
              class="rounded-l-none px-0.5"
              :ui="{ leadingIcon: 'size-3' }" />
            <template #content>
              <div class="flex flex-col min-w-40">
                <button
                  v-for="gItem in entry.group!.items"
                  :key="gItem.command"
                  class="flex items-center gap-2 px-3 py-1.5 text-sm text-left rounded hover:bg-elevated transition-colors"
                  @click="handleCommand(gItem.command)">
                  <UIcon v-if="gItem.icon" :name="gItem.icon" class="size-4 text-muted" />
                  <span>{{ gItem.title }}</span>
                </button>
              </div>
            </template>
          </UPopover>
        </div>
      </slot>
    </template>
  </template>

  <!-- View panels -->
  <template v-for="view in viewItems" :key="`view-${view.pluginId}-${view.id}`">
    <slot name="view" :item="view">
      <div class="plugin-view-panel">
        <div class="flex items-center gap-2 mb-2 text-sm font-medium">
          <UIcon v-if="view.icon" :name="view.icon" class="size-4" />
          <span>{{ view.title }}</span>
        </div>
        <div :id="`plugin-view-${view.id}`" class="plugin-view-content" />
      </div>
    </slot>
  </template>

  <!-- Content blocks (injected by plugin scripts via slots.render) -->
  <template v-for="block in contentBlocks" :key="`block-${block.pluginId}-${block.slot}`">
    <!-- renderFn block: plugin provides a DOM render callback -->
    <div
      v-if="block.renderFn"
      :ref="(el: any) => setRenderFnRef(block, el)"
      class="plugin-content-block plugin-render-fn" />
    <!-- HTML block: trusted plugin with html:inject permission -->
    <div
      v-else-if="block.hasHtmlPermission && (block.trustLevel === 'official' || block.trustLevel === 'local')"
      v-html="sanitizeHtml(block.content)"
      class="plugin-content-block" />
    <!-- Text block: fallback for untrusted or no html:inject -->
    <div v-else class="plugin-content-block">{{ block.content }}</div>
  </template>
</template>
