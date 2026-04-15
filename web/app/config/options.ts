/**
 * Options Schema — single source of truth for ALL site options.
 *
 * Every option has a Zod schema (for runtime validation), a typed default,
 * autoload flag and a display label.
 *
 * TypeScript infers OptionValue<K> directly from the default field, so
 * useOption.ts callers get proper types without any manual casting.
 *
 * Usage:
 *   const { getOption } = useOption()
 *   const n = getOption('posts_per_page')           // number
 *   const layout = getOption('post_default_layout') // string
 *   const widgets = getOption('blog_sidebar_widgets') // WidgetConfig[]
 */

import { z } from 'zod'
import { RoleCapabilitiesSchema } from '~/config/permissions'

// ---------------------------------------------------------------------------
// JSON sub-schemas (canonical shapes for all json-typed options)
// ---------------------------------------------------------------------------

export const WidgetConfigSchema = z.object({
  id:          z.string(),
  label:       z.string(),
  enabled:     z.boolean(),
  title:       z.string().optional(),
  showRecent:  z.boolean().optional(),
  showHot:     z.boolean().optional(),
  maxCount:    z.number().int().positive().optional(),
  pluginSettings: z.record(z.unknown()).optional(),
})

export const SectionConfigSchema = z.object({
  id:                 z.string(),
  label:              z.string(),
  enabled:            z.boolean(),
  count:              z.number().int().positive(),
  title:              z.string().optional(),
  layout:             z.string().optional(),
  includeCategoryIds: z.array(z.number()).optional(),
  excludeCategoryIds: z.array(z.number()).optional(),
})

export const NavMenuItemSchema = z.object({
  local_id:        z.string(),
  label:           z.string(),
  url:             z.string(),
  object_type:     z.enum(['custom', 'page', 'category', 'archive', 'action', 'separator']),
  object_id:       z.number(),
  target:          z.string(),
  css_classes:     z.string(),
  parent_local_id: z.string(),
})

export const NavCustomMenuSchema = z.object({
  id:          z.string(),
  name:        z.string(),
  description: z.string(),
  items:       z.array(NavMenuItemSchema),
})

export const OAuthProviderConfigSchema = z.object({
  enabled:      z.boolean(),
  clientId:     z.string(),
  clientSecret: z.string(),
  callbackUrl:  z.string(),
})

export type OAuthProviderConfig = z.infer<typeof OAuthProviderConfigSchema>

// GenericProviderConfig — full config for a frontend-defined OAuth2 provider.
export const GenericOAuthFieldMapSchema = z.object({
  id:     z.string(),
  email:  z.string(),
  name:   z.string(),
  avatar: z.string(),
})

export const GenericOAuthProviderConfigSchema = OAuthProviderConfigSchema.extend({
  slug:        z.string(),
  label:       z.string(),
  icon:        z.string(),
  authUrl:     z.string(),
  tokenUrl:    z.string(),
  userInfoUrl: z.string(),
  scopes:      z.array(z.string()),
  fields:      GenericOAuthFieldMapSchema,
})

export type GenericOAuthProviderConfig = z.infer<typeof GenericOAuthProviderConfigSchema>

// Inferred types — re-exported so composables can import from here instead
// of defining their own interfaces.
export type WidgetConfig    = z.infer<typeof WidgetConfigSchema>
export type SectionConfig   = z.infer<typeof SectionConfigSchema>
export type NavMenuItem     = z.infer<typeof NavMenuItemSchema>
export type NavCustomMenu   = z.infer<typeof NavCustomMenuSchema>

// ---------------------------------------------------------------------------
// Option entry type
// ---------------------------------------------------------------------------

type OptionKind = 'string' | 'number' | 'boolean' | 'json'

interface OptionEntry<T = unknown> {
  type:     OptionKind
  schema:   z.ZodType<T>
  default:  T
  autoload: boolean
  label:    string
}

// ---------------------------------------------------------------------------
// Builder helpers — keep call sites concise
// ---------------------------------------------------------------------------

function optStr(def: string, label: string, autoload = true): OptionEntry<string> {
  return { type: 'string', schema: z.string().default(def), default: def, autoload, label }
}

function optNum(def: number, label: string, autoload = true): OptionEntry<number> {
  return { type: 'number', schema: z.number().default(def), default: def, autoload, label }
}

function optBool(def: boolean, label: string, autoload = true): OptionEntry<boolean> {
  return { type: 'boolean', schema: z.boolean().default(def), default: def, autoload, label }
}

function optJSON<T>(
  schema: z.ZodType<T>,
  def: T,
  label: string,
  autoload = true,
): OptionEntry<T> {
  return { type: 'json', schema, default: def, autoload, label }
}

// ---------------------------------------------------------------------------
// Schema definition
// ---------------------------------------------------------------------------

export const OPTIONS_SCHEMA = {
  // ── 站点基本信息 ──────────────────────────────────────────────────────────
  site_name:           optStr('个人博客',              '站点名称'),
  site_description:    optStr('',                      '站点描述'),
  site_url:            optStr('',                      '站点 URL'),
  admin_email:         optStr('',                      '管理员邮箱'),
  language:            optStr('zh-CN',                 '站点语言'),
  allow_registration:  optBool(false,                  '开放注册'),
  footer_text:         optStr('',                      '页脚文字'),
  icp_number:          optStr('',                      'ICP 备案号'),
  police_number:       optStr('',                      '公安备案号'),

  // ── 封面图 ────────────────────────────────────────────────────────────────
  default_post_cover:  optStr('/images/default-cover.svg', '默认封面图'),
  error_post_cover:    optStr('/images/default-cover.svg', '加载失败封面图'),

  // ── 写作 ──────────────────────────────────────────────────────────────────
  default_editor:      optStr('markdown', '默认编辑器'),
  auto_save:           optBool(true,      '自动保存'),
  auto_save_interval:  optNum(60,         '自动保存间隔（秒）'),

  // ── 文章列表 ──────────────────────────────────────────────────────────────
  posts_per_page:      optNum(10, '每页文章数'),

  // ── 文章页面 ──────────────────────────────────────────────────────────────
  post_default_layout:    optStr('hero',  '文章默认封面布局'),
  post_sidebar_enabled:   optBool(false,  '文章默认显示侧栏'),
  blog_sidebar_widgets:   optJSON(
    z.array(WidgetConfigSchema),
    [] as WidgetConfig[],
    '文章侧栏小部件',
  ),

  // ── 静态页面 ──────────────────────────────────────────────────────────────
  page_default_layout:    optStr('none',  '页面默认封面布局'),
  page_sidebar_enabled:   optBool(false,  '页面默认显示侧栏'),
  page_sidebar_widgets:   optJSON(
    z.array(WidgetConfigSchema),
    [] as WidgetConfig[],
    '页面侧栏小部件',
  ),

  // ── 首页 ──────────────────────────────────────────────────────────────────
  homepage_sidebar_enabled: optBool(false, '首页显示侧栏'),
  homepage_sidebar_widgets: optJSON(
    z.array(WidgetConfigSchema),
    [] as WidgetConfig[],
    '首页侧栏小部件',
  ),
  homepage_sections: optJSON(
    z.array(SectionConfigSchema),
    [] as SectionConfig[],
    '首页版块配置',
  ),

  // ── 导航菜单 ──────────────────────────────────────────────────────────────
  primary_menu:   optJSON(z.array(NavMenuItemSchema), [] as NavMenuItem[], '主导航'),
  secondary_menu: optJSON(z.array(NavMenuItemSchema), [] as NavMenuItem[], '副导航'),
  footer_menu:    optJSON(z.array(NavMenuItemSchema), [] as NavMenuItem[], '页脚导航'),
  social_menu:    optJSON(z.array(NavMenuItemSchema), [] as NavMenuItem[], '社交链接'),
  header_actions: optJSON(z.array(NavMenuItemSchema), [
    { local_id: 'action:lang_switcher', label: '语言切换', url: '', object_type: 'action' as const, object_id: 0, target: '', css_classes: '', parent_local_id: '' },
    { local_id: 'action:theme_toggle', label: '主题切换', url: '', object_type: 'action' as const, object_id: 0, target: '', css_classes: '', parent_local_id: '' },
    { local_id: 'action:messages', label: '消息', url: '', object_type: 'action' as const, object_id: 0, target: '', css_classes: '', parent_local_id: '' },
    { local_id: 'action:notifications', label: '通知', url: '', object_type: 'action' as const, object_id: 0, target: '', css_classes: '', parent_local_id: '' },
  ] as NavMenuItem[], '顶栏操作按钮'),

  user_menu: optJSON(z.array(NavMenuItemSchema), [
    { local_id: 'action:admin_dashboard', label: '管理后台', url: '', object_type: 'action' as const, object_id: 0, target: '', css_classes: '', parent_local_id: '' },
    { local_id: 'sep:1', label: '', url: '', object_type: 'separator' as const, object_id: 0, target: '', css_classes: '', parent_local_id: '' },
    { local_id: 'action:user_profile', label: '个人主页', url: '', object_type: 'action' as const, object_id: 0, target: '', css_classes: '', parent_local_id: '' },
    { local_id: 'action:user_posts', label: '我的文章', url: '', object_type: 'action' as const, object_id: 0, target: '', css_classes: '', parent_local_id: '' },
    { local_id: 'action:user_favorites', label: '我的收藏', url: '', object_type: 'action' as const, object_id: 0, target: '', css_classes: '', parent_local_id: '' },
    { local_id: 'sep:2', label: '', url: '', object_type: 'separator' as const, object_id: 0, target: '', css_classes: '', parent_local_id: '' },
    { local_id: 'action:user_settings', label: '设置', url: '', object_type: 'action' as const, object_id: 0, target: '', css_classes: '', parent_local_id: '' },
    { local_id: 'sep:3', label: '', url: '', object_type: 'separator' as const, object_id: 0, target: '', css_classes: '', parent_local_id: '' },
    { local_id: 'action:user_logout', label: '退出登录', url: '', object_type: 'action' as const, object_id: 0, target: '', css_classes: '', parent_local_id: '' },
  ] as NavMenuItem[], '用户菜单'),

  floating_toolbar: optJSON(z.array(NavMenuItemSchema), [
    { local_id: 'action:profile_login', label: '个人中心', url: '', object_type: 'action' as const, object_id: 0, target: '', css_classes: '', parent_local_id: '' },
    { local_id: 'action:checkin', label: '每日签到', url: '', object_type: 'action' as const, object_id: 0, target: '', css_classes: '', parent_local_id: '' },
    { local_id: 'action:ft_notifications', label: '通知', url: '', object_type: 'action' as const, object_id: 0, target: '', css_classes: '', parent_local_id: '' },
    { local_id: 'sep:1', label: '', url: '', object_type: 'separator' as const, object_id: 0, target: '', css_classes: '', parent_local_id: '' },
    { local_id: 'action:ft_theme_toggle', label: '主题切换', url: '', object_type: 'action' as const, object_id: 0, target: '', css_classes: '', parent_local_id: '' },
  ] as NavMenuItem[], '悬浮工具栏'),

  post_actions: optJSON(z.array(NavMenuItemSchema), [
    { local_id: 'action:post_like', label: '点赞', url: '', object_type: 'action' as const, object_id: 0, target: '', css_classes: '', parent_local_id: '' },
    { local_id: 'action:post_comment', label: '评论', url: '', object_type: 'action' as const, object_id: 0, target: '', css_classes: '', parent_local_id: '' },
    { local_id: 'action:post_bookmark', label: '收藏', url: '', object_type: 'action' as const, object_id: 0, target: '', css_classes: '', parent_local_id: '' },
    { local_id: 'sep:1', label: '', url: '', object_type: 'separator' as const, object_id: 0, target: '', css_classes: '', parent_local_id: '' },
    { local_id: 'action:post_share', label: '分享', url: '', object_type: 'action' as const, object_id: 0, target: '', css_classes: '', parent_local_id: '' },
  ] as NavMenuItem[], '文章操作栏'),
  nav_custom_menus: optJSON(
    z.array(NavCustomMenuSchema),
    [] as NavCustomMenu[],
    '自定义菜单',
    false,
  ),

  // ── 外观主题 ──────────────────────────────────────────────────────────────
  theme_primary:    optStr('violet',    '主色调'),
  theme_neutral:    optStr('zinc',      '中性色'),
  theme_radius:     optStr('0.375rem',  '圆角大小'),
  theme_color_mode: optStr('system',    '颜色模式（light/dark/system）'),
  theme_font:       optStr('system',    '字体'),
  theme_font_size:  optNum(16,          '字体大小（px）'),
  theme_custom_css: optStr('',          '自定义 CSS', false),

  // ── 媒体上传大小限制 ──────────────────────────────────────────────────────
  media_size_limits: optJSON(
    z.object({
      image:    z.number().positive(),
      video:    z.number().positive(),
      audio:    z.number().positive(),
      document: z.number().positive(),
      other:    z.number().positive(),
    }),
    { image: 10, video: 100, audio: 20, document: 20, other: 10 },
    '媒体上传大小限制（MB）',
  ),

  // ── 角色权限 ──────────────────────────────────────────────────────────────
  role_capabilities: optJSON(
    RoleCapabilitiesSchema,
    {} as Record<string, string[]>,
    '角色权限配置',
  ),

  // ── OAuth 三方登录 ─────────────────────────────────────────────────────────
  // oauth_providers: list of custom (frontend-added) provider slugs
  oauth_providers: optJSON(z.array(z.string()), [] as string[], '自定义 OAuth 提供商列表', false),
  oauth_github: optJSON(
    OAuthProviderConfigSchema,
    { enabled: false, clientId: '', clientSecret: '', callbackUrl: 'http://localhost:9000/api/v1/auth/oauth/github/callback' } as OAuthProviderConfig,
    'GitHub OAuth',
    false,
  ),
  oauth_google: optJSON(
    OAuthProviderConfigSchema,
    { enabled: false, clientId: '', clientSecret: '', callbackUrl: 'http://localhost:9000/api/v1/auth/oauth/google/callback' } as OAuthProviderConfig,
    'Google OAuth',
    false,
  ),
  oauth_qq: optJSON(
    OAuthProviderConfigSchema,
    { enabled: false, clientId: '', clientSecret: '', callbackUrl: 'http://localhost:9000/api/v1/auth/oauth/qq/callback' } as OAuthProviderConfig,
    'QQ OAuth',
    false,
  ),

  // ── 插件视图可见性 ──────────────────────────────────────────────────────────
  disabled_plugin_views: optJSON(z.array(z.string()), [] as string[], '已禁用的插件视图'),

  // ── 评论 ──────────────────────────────────────────────────────────────────
  default_allow_comments:      optBool(true,  '新文章默认开启评论'),
  comment_moderation:          optBool(false, '所有评论需审核'),
  comment_require_name_email:  optBool(true,  '评论需填写姓名和邮箱'),
  comment_max_links:           optNum(2,      '评论最多链接数'),
  comment_blacklist:           optStr('',     '评论黑名单'),
} satisfies Record<string, OptionEntry>

export type OptionKey   = keyof typeof OPTIONS_SCHEMA
export type OptionValue<K extends OptionKey> = typeof OPTIONS_SCHEMA[K]['default']
