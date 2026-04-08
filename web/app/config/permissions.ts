/**
 * Capability / Permission System
 *
 * Defines all available capabilities and their default assignment per role.
 * Admins can override the defaults via the role_capabilities option (stored in
 * the options table and autoloaded on startup).
 *
 * Usage:
 *   const { can } = usePermissions()
 *   if (can('manage_options')) { ... }
 */

import { z } from 'zod'

// ── Capability definitions ─────────────────────────────────────────────────────

export const CAPABILITY_GROUPS = {
  admin: {
    label: '后台访问',
    caps: {
      access_admin: '访问后台管理',
    },
  },
  posts: {
    label: '文章管理',
    caps: {
      read_private_posts:   '阅读私密文章',
      edit_posts:           '编辑自己的文章',
      publish_posts:        '发布文章',
      edit_others_posts:    '编辑他人文章',
      delete_posts:         '删除自己的文章',
      delete_others_posts:  '删除他人文章',
    },
  },
  pages: {
    label: '页面管理',
    caps: {
      edit_pages:           '编辑页面',
      publish_pages:        '发布页面',
      edit_others_pages:    '编辑他人页面',
      delete_pages:         '删除页面',
    },
  },
  taxonomy: {
    label: '分类 & 标签',
    caps: {
      manage_categories:    '管理分类和标签',
    },
  },
  media: {
    label: '媒体管理',
    caps: {
      upload_files:         '上传文件',
      manage_media:         '管理所有媒体',
      delete_others_media:  '删除他人媒体',
    },
  },
  comments: {
    label: '评论管理',
    caps: {
      moderate_comments:    '审核评论',
      manage_comments:      '管理所有评论（删除/标垃圾）',
    },
  },
  users: {
    label: '用户管理',
    caps: {
      list_users:           '查看用户列表',
      create_users:         '创建用户',
      edit_users:           '编辑用户信息',
      delete_users:         '删除用户',
      promote_users:        '更改用户角色',
    },
  },
  settings: {
    label: '系统设置',
    caps: {
      manage_options:       '管理站点设置',
      manage_appearance:    '管理外观（主题/菜单）',
      manage_roles:         '管理角色权限',
    },
  },
} as const

// Flatten to a simple map { cap: label }
export const CAPABILITIES = Object.fromEntries(
  Object.values(CAPABILITY_GROUPS).flatMap((g) => Object.entries(g.caps)),
) as Record<Capability, string>

// ── Capability Zod schema ──────────────────────────────────────────────────────

// Derive the union type the same way as before, then build z.enum from it.
export type Capability = keyof {
  [G in keyof typeof CAPABILITY_GROUPS]:
    keyof (typeof CAPABILITY_GROUPS)[G]['caps']
}[keyof typeof CAPABILITY_GROUPS]

const ALL_CAP_KEYS = Object.keys(CAPABILITIES) as [Capability, ...Capability[]]

/** Validates that a value is a known capability string */
export const CapabilitySchema = z.enum(ALL_CAP_KEYS)

/** Validates an array of capability strings */
export const CapabilityArraySchema = z.array(CapabilitySchema)

// ── Role labels ────────────────────────────────────────────────────────────────

export const ROLE_LABELS: Record<number, string> = {
  1: '订阅者',
  2: '编辑',
  3: '管理员',
  4: '超级管理员',
}

/** Validates the saved role_capabilities object from the options store */
export const RoleCapabilitiesSchema = z.record(
  z.string().regex(/^\d+$/),   // key is role id as numeric string
  z.array(z.string()),         // values validated to Capability[] at call site
)

// ── Default capabilities per role ──────────────────────────────────────────────
// These are the factory defaults. Admins can override via role_capabilities option.

const ALL_CAPS = Object.keys(CAPABILITIES) as Capability[]

export const DEFAULT_ROLE_CAPABILITIES: Record<number, Capability[]> = z.record(
  z.coerce.number().int().positive(),
  CapabilityArraySchema,
).parse({
  /** Subscriber — read-only, can manage own profile */
  1: [],

  /** Editor — content management (own resources), no cross-user or system ops */
  2: [
    'access_admin',
    'read_private_posts',
    'edit_posts',
    'publish_posts',
    'delete_posts',
    'edit_pages',
    'publish_pages',
    'manage_categories',
    'upload_files',
    'manage_media',
    'moderate_comments',
    'manage_comments',
    'list_users',
  ],

  /** Admin — everything by default, customizable */
  3: ALL_CAPS,

  /** Super Admin — everything, always, immutable */
  4: ALL_CAPS,
})
