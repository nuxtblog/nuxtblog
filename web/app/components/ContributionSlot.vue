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
import { usePluginContextStore } from '~/stores/plugin-context'
import { usePluginContributionsStore, type PluginContentBlock } from '~/stores/plugin-contributions'

const props = defineProps<{
  /** Slot name matching manifest menus/views keys */
  name: string
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
  <template v-for="item in visibleMenuItems" :key="`menu-${item.pluginId}-${item.command}`">
    <slot name="menu" :item="item" :execute="() => handleCommand(item.command)">
      <UButton
        variant="ghost"
        size="xs"
        :icon="item.icon"
        :label="item.title"
        @click="handleCommand(item.command)"
      />
    </slot>
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
