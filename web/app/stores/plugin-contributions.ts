import { defineStore } from 'pinia'
import type { PluginSettingField } from '~/composables/usePluginApi'

/**
 * Plugin Contribution Points Registry
 *
 * Collects and manages all plugin contributions (commands, navigation,
 * menus, views) from enabled plugins. The admin panel queries this store
 * to render plugin UI extensions via <ContributionSlot>.
 */

export interface PluginCommand {
  pluginId: string
  id: string
  title: string
  titleEn?: string
  icon?: string
  shortcut?: string
}

export interface PluginNavItem {
  pluginId: string
  slot: string
  title: string
  icon?: string
  route: string
  order: number
  parent?: string
  /** Stable name for parent matching (auto-generated for plugin groups) */
  groupKey?: string
}

export interface PluginMenuItem {
  pluginId: string
  slot: string
  command: string
  when?: string
  group?: string
  // Resolved from commands registry:
  title?: string
  titleEn?: string
  icon?: string
  shortcut?: string
}

export interface PluginViewItem {
  pluginId: string
  slot: string
  id: string
  title: string
  type?: string
  icon?: string
  component?: string
  module?: string
  settings?: PluginSettingField[]
}

export interface PluginContentBlock {
  pluginId: string
  slot: string
  content: string
  trustLevel: string
  hasHtmlPermission: boolean
  order: number
  /** Optional DOM render function for dynamic UI injection (official/local only). */
  renderFn?: (container: HTMLElement) => void | (() => void)
}

export interface PluginContributes {
  commands?: Array<{
    id: string
    title: string
    title_en?: string
    icon?: string
    shortcut?: string
  }>
  navigation?: Array<{
    slot: string
    title: string
    icon?: string
    route: string
    order?: number
    parent?: string
    groupKey?: string
  }>
  menus?: Record<string, Array<{
    command: string
    when?: string
    group?: string
  }>>
  views?: Record<string, Array<{
    id: string
    title?: string
    type?: string
    icon?: string
    component?: string
    module?: string
    settings?: PluginSettingField[]
  }>>
  pages?: Array<{
    path?: string
    slot: string
    component?: string
    module?: string
    title?: string
    nav?: { type?: string; group?: string; icon?: string; order?: number; slot?: string }
  }>
  styles?: Array<{
    scope: string
    file: string
  }>
  activation?: Array<{
    scope: string
    module: string
  }>
}

export interface PluginPageDef {
  pluginId: string
  component: string
  title?: string
  /** Custom route path declared in plugin.yaml (e.g. '/shop', '/admin/shop/products') */
  path?: string
  /** 'admin.mjs' or 'public.mjs' — determines which module to load */
  moduleFile: string
}

export const usePluginContributionsStore = defineStore('plugin-contributions', () => {
  const commands = ref<PluginCommand[]>([])
  const navigation = ref<PluginNavItem[]>([])
  const menuItems = ref<PluginMenuItem[]>([])
  const viewItems = ref<PluginViewItem[]>([])
  const contentBlocks = ref<PluginContentBlock[]>([])
  const pluginPages = ref<PluginPageDef[]>([])

  /** Register contributions from a plugin manifest. */
  function registerPlugin(pluginId: string, contributes: PluginContributes) {
    // Remove any existing contributions from this plugin to prevent duplicates
    unregisterPlugin(pluginId)

    // Commands
    if (contributes.commands) {
      for (const cmd of contributes.commands) {
        commands.value.push({
          pluginId,
          id: cmd.id,
          title: cmd.title,
          titleEn: cmd.title_en,
          icon: cmd.icon,
          shortcut: cmd.shortcut,
        })
      }
    }

    // Navigation
    if (contributes.navigation) {
      for (const nav of contributes.navigation) {
        navigation.value.push({
          pluginId,
          slot: nav.slot,
          title: nav.title,
          icon: nav.icon,
          route: nav.route,
          order: nav.order ?? 100,
          parent: nav.parent,
          groupKey: nav.groupKey,
        })
      }
    }

    // Menus
    if (contributes.menus) {
      for (const [slot, entries] of Object.entries(contributes.menus)) {
        for (const entry of entries) {
          // Resolve command details
          const cmd = contributes.commands?.find(c => c.id === entry.command)
          menuItems.value.push({
            pluginId,
            slot,
            command: entry.command,
            when: entry.when,
            group: entry.group,
            title: cmd?.title,
            titleEn: cmd?.title_en,
            icon: cmd?.icon,
            shortcut: cmd?.shortcut,
          })
        }
      }
    }

    // Views
    if (contributes.views) {
      for (const [slot, views] of Object.entries(contributes.views)) {
        for (const view of views) {
          viewItems.value.push({
            pluginId,
            slot,
            id: view.id,
            title: view.title || '',
            type: view.type,
            icon: view.icon,
            component: view.component,
            module: view.module,
            settings: view.settings,
          })
        }
      }
    }
  }

  /** Register a content block from a plugin script (via slots.render). */
  function registerContentBlock(block: PluginContentBlock) {
    // Remove existing block from same plugin + slot to prevent duplicates
    contentBlocks.value = contentBlocks.value.filter(
      b => !(b.pluginId === block.pluginId && b.slot === block.slot),
    )
    contentBlocks.value.push(block)
  }

  /** Get content blocks for a specific slot, sorted by order. */
  function getContentBlocks(slot: string) {
    return computed(() =>
      contentBlocks.value
        .filter(b => b.slot === slot)
        .sort((a, b) => a.order - b.order),
    )
  }

  /** Register a public plugin page. */
  function registerPluginPage(page: PluginPageDef) {
    // Deduplicate
    pluginPages.value = pluginPages.value.filter(
      p => !(p.pluginId === page.pluginId && p.component === page.component),
    )
    pluginPages.value.push(page)
  }

  /** Get a plugin page definition by pluginId and component name. */
  function getPluginPage(pluginId: string, component: string) {
    return pluginPages.value.find(
      p => p.pluginId === pluginId && p.component === component,
    )
  }

  /** Unregister all contributions from a plugin. */
  function unregisterPlugin(pluginId: string) {
    commands.value = commands.value.filter(c => c.pluginId !== pluginId)
    navigation.value = navigation.value.filter(n => n.pluginId !== pluginId)
    menuItems.value = menuItems.value.filter(m => m.pluginId !== pluginId)
    viewItems.value = viewItems.value.filter(v => v.pluginId !== pluginId)
    contentBlocks.value = contentBlocks.value.filter(b => b.pluginId !== pluginId)
    pluginPages.value = pluginPages.value.filter(p => p.pluginId !== pluginId)
  }

  /** Get navigation items for a specific slot, sorted by order. */
  function getNavigation(slot: string) {
    return computed(() =>
      navigation.value
        .filter(n => n.slot === slot)
        .sort((a, b) => a.order - b.order),
    )
  }

  /** Get navigation items by parent menu name, sorted by order. */
  function getChildNavigation(parentName: string) {
    return computed(() =>
      navigation.value
        .filter(n => (n.parent || '') === parentName)
        .sort((a, b) => a.order - b.order),
    )
  }

  /** Get menu items for a specific slot. */
  function getMenuItems(slot: string) {
    return computed(() =>
      menuItems.value.filter(m => m.slot === slot),
    )
  }

  /** Get view items for a specific slot. */
  function getViewItems(slot: string) {
    return computed(() =>
      viewItems.value.filter(v => v.slot === slot),
    )
  }

  /** Get a command by its full ID. */
  function getCommand(commandId: string) {
    return commands.value.find(c => c.id === commandId)
  }

  /** Clear all contributions (e.g. on full reload). */
  function clear() {
    commands.value = []
    navigation.value = []
    menuItems.value = []
    viewItems.value = []
    contentBlocks.value = []
    pluginPages.value = []
  }

  return {
    commands,
    navigation,
    menuItems,
    viewItems,
    contentBlocks,
    pluginPages,
    registerPlugin,
    registerContentBlock,
    registerPluginPage,
    getPluginPage,
    unregisterPlugin,
    getNavigation,
    getChildNavigation,
    getMenuItems,
    getViewItems,
    getContentBlocks,
    getCommand,
    clear,
  }
})
