export type MenuItemType = 'custom' | 'page' | 'category' | 'archive'

/**
 * Built-in menu slots registered by the theme.
 * These cannot be deleted — only their items can be edited.
 * Each slot is stored as its own option key (same as the slot key).
 *
 * Frontend usage:
 *   const items = optionsStore.getJSON<NavMenuItem[]>('primary_menu', [])
 */
export const NAV_MENU_SLOTS = {
  primary_menu: '主导航',
} as const

export const NAV_MENU_SLOT_HINTS: Record<keyof typeof NAV_MENU_SLOTS, string> = {
  primary_menu: '显示在页面顶部导航栏，支持下拉子菜单',
}

export type NavMenuSlotKey = keyof typeof NAV_MENU_SLOTS

export interface NavMenuItem {
  local_id: string        // client-side ID, used for parent refs
  label: string
  url: string
  object_type: MenuItemType
  object_id: number
  target: string          // '' | '_blank'
  css_classes: string
  parent_local_id: string // '' = root item
}

/**
 * UI-only representation of a menu item (not persisted directly).
 * Adds depth, parent_local_id, expanded, and flattened type fields.
 */
export interface UiMenuItem {
  local_id: string
  label: string
  url: string
  type: MenuItemType
  object_id: number
  openInNewTab: boolean
  cssClasses: string
  depth: number
  parent_local_id: string
  expanded: boolean
}

/**
 * Custom (user-created) menu — stored in `nav_custom_menus` array.
 * These can be created and deleted freely.
 */
export interface NavCustomMenu {
  id: string
  name: string
  description: string
  items: NavMenuItem[]
}
