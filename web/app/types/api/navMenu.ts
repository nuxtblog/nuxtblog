export type MenuItemType = 'custom' | 'page' | 'category' | 'archive' | 'action' | 'separator'

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
  header_actions: '顶栏操作按钮',
  user_menu: '用户菜单',
  floating_toolbar: '悬浮工具栏',
  post_actions: '文章操作栏',
} as const

export const NAV_MENU_SLOT_HINTS: Record<keyof typeof NAV_MENU_SLOTS, string> = {
  primary_menu: '显示在页面顶部导航栏，支持下拉子菜单',
  header_actions: '显示在顶栏右侧的操作按钮，如语言切换、主题切换等',
  user_menu: '登录用户的下拉菜单',
  floating_toolbar: '页面右下角的悬浮快捷工具栏',
  post_actions: '文章页面左侧的互动操作栏',
}

export const HEADER_BUILTIN_ACTIONS = [
  { id: 'action:lang_switcher', label: '语言切换', icon: 'i-tabler-language' },
  { id: 'action:theme_toggle', label: '主题切换', icon: 'i-tabler-sun-moon' },
  { id: 'action:messages', label: '消息', icon: 'i-tabler-message' },
  { id: 'action:notifications', label: '通知', icon: 'i-tabler-bell' },
] as const

export const USER_MENU_BUILTIN_ACTIONS = [
  { id: 'action:admin_dashboard', label: '管理后台', icon: 'i-tabler-layout-dashboard' },
  { id: 'action:user_profile', label: '个人主页', icon: 'i-tabler-user' },
  { id: 'action:user_posts', label: '我的文章', icon: 'i-tabler-file-text' },
  { id: 'action:user_favorites', label: '我的收藏', icon: 'i-tabler-bookmark' },
  { id: 'action:user_settings', label: '设置', icon: 'i-tabler-settings' },
  { id: 'action:user_logout', label: '退出登录', icon: 'i-tabler-logout' },
] as const

export const FLOATING_TOOLBAR_BUILTIN_ACTIONS = [
  { id: 'action:profile_login', label: '个人中心', icon: 'i-tabler-user-circle' },
  { id: 'action:ft_notifications', label: '通知', icon: 'i-tabler-bell' },
  { id: 'action:ft_theme_toggle', label: '主题切换', icon: 'i-tabler-sun-moon' },
] as const

export const POST_ACTIONS_BUILTIN_ACTIONS = [
  { id: 'action:post_like', label: '点赞', icon: 'i-tabler-heart' },
  { id: 'action:post_comment', label: '评论', icon: 'i-tabler-message-circle' },
  { id: 'action:post_bookmark', label: '收藏', icon: 'i-tabler-bookmark' },
  { id: 'action:post_share', label: '分享', icon: 'i-tabler-share-2' },
] as const

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

// ── Slot config: drives AddPanel / Structure / Preview per slot ────────────

export interface NavMenuSlotConfig {
  builtinActions?: readonly { id: string; label: string; icon: string }[]
  showPages?: boolean
  showCategories?: boolean
  showCustomLink?: boolean
  customLinkHasIcon?: boolean
  supportsNesting?: boolean
  contributionSlot?: string
  defaultItems?: NavMenuItem[]
}

function actionsToDefaultItems(actions: readonly { id: string; label: string }[]): NavMenuItem[] {
  return actions.map(a => ({
    local_id: a.id,
    label: a.label,
    url: '',
    object_type: 'action' as MenuItemType,
    object_id: 0,
    target: '',
    css_classes: '',
    parent_local_id: '',
  }))
}

export const NAV_MENU_SLOT_CONFIGS: Partial<Record<NavMenuSlotKey, NavMenuSlotConfig>> = {
  primary_menu: { showPages: true, showCategories: true, showCustomLink: true, supportsNesting: true },
  header_actions: { builtinActions: HEADER_BUILTIN_ACTIONS, showCustomLink: true, customLinkHasIcon: true, contributionSlot: 'public:header-actions', defaultItems: actionsToDefaultItems(HEADER_BUILTIN_ACTIONS) },
  user_menu: { builtinActions: USER_MENU_BUILTIN_ACTIONS, showCustomLink: true, contributionSlot: 'public:user-menu', defaultItems: actionsToDefaultItems(USER_MENU_BUILTIN_ACTIONS) },
  floating_toolbar: { builtinActions: FLOATING_TOOLBAR_BUILTIN_ACTIONS, contributionSlot: 'public:floating-toolbar', defaultItems: actionsToDefaultItems(FLOATING_TOOLBAR_BUILTIN_ACTIONS) },
  post_actions: { builtinActions: POST_ACTIONS_BUILTIN_ACTIONS, contributionSlot: 'public:post-actions', defaultItems: actionsToDefaultItems(POST_ACTIONS_BUILTIN_ACTIONS) },
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
