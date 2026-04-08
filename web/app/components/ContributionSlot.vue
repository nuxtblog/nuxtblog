<script setup lang="ts">
/**
 * Phase 2.3: ContributionSlot — renders plugin UI extensions in named anchor points.
 *
 * Usage:
 *   <ContributionSlot name="post-editor:toolbar" />
 *   <ContributionSlot name="admin:sidebar-nav" />
 *   <ContributionSlot name="post-editor:context" :ctx="editorCtx" />
 */
import { usePluginContextStore } from '~/stores/plugin-context'
import { usePluginContributionsStore } from '~/stores/plugin-contributions'

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

/** Filter menu items by their `when` condition. */
const visibleMenuItems = computed(() =>
  menuItems.value.filter(item =>
    !item.when || contextStore.evaluateWhen(item.when),
  ),
)

function handleCommand(commandId: string) {
  emit('command', commandId, props.ctx)
}
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
</template>
