/**
 * Theme Configuration
 *
 * Defines all available color palettes, radius options, and fonts.
 * All shapes are Zod-validated; TypeScript types are inferred from schemas.
 */

import { z } from 'zod'

// ── Color ──────────────────────────────────────────────────────────────────────

export const ThemeColorSchema = z.object({
  /** Tailwind color name, e.g. 'violet' */
  name:  z.string(),
  /** Display label */
  label: z.string(),
  /** Tailwind 500-shade hex — used for swatch display only */
  hex:   z.string().regex(/^#[0-9a-fA-F]{6}$/),
})

export type ThemeColor = z.infer<typeof ThemeColorSchema>

export const PRIMARY_COLORS: ThemeColor[] = z.array(ThemeColorSchema).parse([
  { name: 'red',     label: '红色',   hex: '#ef4444' },
  { name: 'orange',  label: '橙色',   hex: '#f97316' },
  { name: 'amber',   label: '琥珀',   hex: '#f59e0b' },
  { name: 'yellow',  label: '黄色',   hex: '#eab308' },
  { name: 'lime',    label: '青柠',   hex: '#84cc16' },
  { name: 'green',   label: '绿色',   hex: '#22c55e' },
  { name: 'emerald', label: '翠绿',   hex: '#10b981' },
  { name: 'teal',    label: '蓝绿',   hex: '#14b8a6' },
  { name: 'cyan',    label: '青色',   hex: '#06b6d4' },
  { name: 'sky',     label: '天蓝',   hex: '#0ea5e9' },
  { name: 'blue',    label: '蓝色',   hex: '#3b82f6' },
  { name: 'indigo',  label: '靛青',   hex: '#6366f1' },
  { name: 'violet',  label: '紫罗兰', hex: '#8b5cf6' },
  { name: 'purple',  label: '紫色',   hex: '#a855f7' },
  { name: 'fuchsia', label: '品红',   hex: '#d946ef' },
  { name: 'pink',    label: '粉红',   hex: '#ec4899' },
  { name: 'rose',    label: '玫瑰',   hex: '#f43f5e' },
])

export const NEUTRAL_COLORS: ThemeColor[] = z.array(ThemeColorSchema).parse([
  { name: 'slate',   label: '石板灰', hex: '#64748b' },
  { name: 'gray',    label: '灰色',   hex: '#6b7280' },
  { name: 'zinc',    label: '锌灰',   hex: '#71717a' },
  { name: 'neutral', label: '中性灰', hex: '#737373' },
  { name: 'stone',   label: '石色',   hex: '#78716c' },
])

// ── Border radius ──────────────────────────────────────────────────────────────

export const RadiusOptionSchema = z.object({
  /** CSS value set on --ui-radius */
  value:   z.string(),
  label:   z.string(),
  /** border-radius for preview box */
  preview: z.string(),
})

export type RadiusOption = z.infer<typeof RadiusOptionSchema>

export const RADIUS_OPTIONS: RadiusOption[] = z.array(RadiusOptionSchema).parse([
  { value: '0',        label: '0',     preview: '0' },
  { value: '0.125rem', label: '0.125', preview: '2px' },
  { value: '0.25rem',  label: '0.25',  preview: '4px' },
  { value: '0.375rem', label: '0.375', preview: '6px' },
  { value: '0.5rem',   label: '0.5',   preview: '8px' },
])

// ── Font stacks ────────────────────────────────────────────────────────────────

export const FontOptionSchema = z.object({
  value: z.string(),
  label: z.string(),
  /** Actual CSS font-family value */
  stack: z.string(),
})

export type FontOption = z.infer<typeof FontOptionSchema>

export const FONT_OPTIONS: FontOption[] = z.array(FontOptionSchema).parse([
  {
    value: 'system',
    label: '系统默认',
    stack: 'system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif',
  },
  {
    value: 'sans',
    label: '无衬线',
    stack: '"Inter", "Noto Sans SC", sans-serif',
  },
  {
    value: 'serif',
    label: '衬线体',
    stack: '"Noto Serif SC", "Source Han Serif", Georgia, serif',
  },
  {
    value: 'mono',
    label: '等宽',
    stack: '"JetBrains Mono", "Fira Code", Consolas, monospace',
  },
])

// ── Theme settings shape ───────────────────────────────────────────────────────

export const ThemeSettingsSchema = z.object({
  primary:   z.string(),               // color name, e.g. 'violet'
  neutral:   z.string(),               // color name, e.g. 'zinc'
  radius:    z.string(),               // CSS value, e.g. '0.375rem'
  colorMode: z.enum(['light', 'dark', 'system']),
  font:      z.string(),               // font option value
  fontSize:  z.number().int().min(10).max(32),
  customCss: z.string(),
})

export type ThemeSettings = z.infer<typeof ThemeSettingsSchema>

export const DEFAULT_THEME: ThemeSettings = ThemeSettingsSchema.parse({
  primary:   'violet',
  neutral:   'zinc',
  radius:    '0.375rem',
  colorMode: 'system',
  font:      'system',
  fontSize:  16,
  customCss: '',
})
