import { defineStore } from 'pinia'

/**
 * Phase 2.2 & 2.3: Plugin Contribution Points Registry
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
}

export interface PluginMenuItem {
  pluginId: string
  slot: string
  command: string
  when?: string
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
  }>
  menus?: Record<string, Array<{
    command: string
    when?: string
  }>>
  views?: Record<string, Array<{
    id: string
    title: string
    type?: string
    icon?: string
  }>>
}

export const usePluginContributionsStore = defineStore('plugin-contributions', () => {
  const commands = ref<PluginCommand[]>([])
  const navigation = ref<PluginNavItem[]>([])
  const menuItems = ref<PluginMenuItem[]>([])
  const viewItems = ref<PluginViewItem[]>([])

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
            title: view.title,
            type: view.type,
            icon: view.icon,
          })
        }
      }
    }
  }

  /** Unregister all contributions from a plugin. */
  function unregisterPlugin(pluginId: string) {
    commands.value = commands.value.filter(c => c.pluginId !== pluginId)
    navigation.value = navigation.value.filter(n => n.pluginId !== pluginId)
    menuItems.value = menuItems.value.filter(m => m.pluginId !== pluginId)
    viewItems.value = viewItems.value.filter(v => v.pluginId !== pluginId)
  }

  /** Get navigation items for a specific slot, sorted by order. */
  function getNavigation(slot: string) {
    return computed(() =>
      navigation.value
        .filter(n => n.slot === slot)
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
  }

  return {
    commands,
    navigation,
    menuItems,
    viewItems,
    registerPlugin,
    unregisterPlugin,
    getNavigation,
    getMenuItems,
    getViewItems,
    getCommand,
    clear,
  }
})
