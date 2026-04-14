<script setup lang="ts">
/**
 * ContributionSlot — renders plugin UI extensions in named anchor points.
 *
 * Usage:
 *   <ContributionSlot name="post-editor:toolbar" />
 *   <ContributionSlot name="admin:sidebar-nav" />
 *   <ContributionSlot name="post-editor:context" :ctx="editorCtx" />
 */
import DOMPurify from 'dompurify'
import type { ContributionSlotName } from '~/config/contribution-slots'
import { usePluginContextStore } from '~/stores/plugin-context'
import { usePluginContributionsStore, type PluginContentBlock, type PluginMenuItem, type PluginViewItem } from '~/stores/plugin-contributions'
import { createPluginAsyncComponent } from '~/composables/usePluginComponents'

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
  /** When set, only render the single view/menu item matching this id */
  filterId?: string
}>()

const emit = defineEmits<{
  command: [commandId: string, ctx?: Record<string, unknown>]
}>()

const contextStore = usePluginContextStore()
const contributionsStore = usePluginContributionsStore()
const { getOption } = useOption()

const menuItems = contributionsStore.getMenuItems(props.name)
const allViewItems = contributionsStore.getViewItems(props.name)
const disabledViews = computed(() => getOption('disabled_plugin_views'))
const viewItems = computed(() => {
  const items = allViewItems.value.filter(v => !disabledViews.value.includes(v.id))
  return props.filterId ? items.filter(v => v.id === props.filterId) : items
})
const navItems = contributionsStore.getNavigation(props.name)
const contentBlocks = contributionsStore.getContentBlocks(props.name)

/** Sanitize HTML content for safe rendering via DOMPurify. */
function sanitizeHtml(html: string): string {
  if (import.meta.server) return ''
  return DOMPurify.sanitize(html)
}

/** Filter menu items by their `when` condition and optional filterId. */
const visibleMenuItems = computed(() => {
  let items = menuItems.value.filter(item =>
    !item.when || contextStore.evaluateWhen(item.when),
  )
  if (props.filterId) items = items.filter(m => m.command === props.filterId)
  return items
})

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

// ── Rich component views: resolve async components for views with component+module ──
const viewComponentCache = new Map<string, ReturnType<typeof createPluginAsyncComponent>>()
function getViewComponent(view: PluginViewItem) {
  if (!view.component || !view.module) return null
  const key = `${view.pluginId}:${view.component}:${view.module}`
  if (!viewComponentCache.has(key)) {
    viewComponentCache.set(key, createPluginAsyncComponent(view.pluginId, view.component, view.module))
  }
  return viewComponentCache.get(key)!
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
      <!-- Rich component view: render via createPluginAsyncComponent -->
      <template v-if="view.component && view.module">
        <ClientOnly>
          <component :is="getViewComponent(view)" v-if="getViewComponent(view)" />
        </ClientOnly>
      </template>
      <!-- Legacy declarative view panel -->
      <div v-else class="plugin-view-panel">
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
