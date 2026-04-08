/**
 * Nuxt server route: GET /api/theme.css
 *
 * Generates a small CSS file that pre-applies the saved theme (primary color,
 * neutral color, border radius, font, font-size) BEFORE any JavaScript runs.
 * This eliminates the flash of unstyled content (FOUC) that would otherwise
 * occur when theme.client.ts applies the saved settings after hydration.
 *
 * The CSS mirrors exactly what useTheme composable does at runtime, but is
 * served as a render-blocking stylesheet linked in <head>.
 */

const SHADES = [50, 100, 200, 300, 400, 500, 600, 700, 800, 900, 950] as const

const FONT_STACKS: Record<string, string> = {
  system: 'system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif',
  sans:   '"Inter", "Noto Sans SC", sans-serif',
  serif:  '"Noto Serif SC", "Source Han Serif", Georgia, serif',
  mono:   '"JetBrains Mono", "Fira Code", Consolas, monospace',
}

const DEFAULTS = {
  primary:   'violet',
  neutral:   'zinc',
  radius:    '0.375rem',
  font:      'system',
  fontSize:  16,
  customCss: '',
}

function colorAliases(semantic: string, color: string): string {
  return SHADES.map((s) => `  --color-${semantic}-${s}: var(--color-${color}-${s});`).join('\n')
}

function parseSaved(raw: string | undefined, fallback: string): string {
  if (!raw) return fallback
  try {
    const v = JSON.parse(raw)
    return typeof v === 'string' ? v : fallback
  } catch {
    return fallback
  }
}

export default defineEventHandler(async (event) => {
  setResponseHeader(event, 'Content-Type', 'text/css; charset=utf-8')
  setResponseHeader(event, 'Cache-Control', 'no-store')

  const apiBase = process.env.NUXT_PUBLIC_API_BASE || 'http://localhost:9000/api/v1'

  let primary   = DEFAULTS.primary
  let neutral   = DEFAULTS.neutral
  let radius    = DEFAULTS.radius
  let font      = DEFAULTS.font
  let fontSize  = DEFAULTS.fontSize
  let customCss = DEFAULTS.customCss

  try {
    const res = await $fetch<{ code: number; data: { options: Record<string, string> } }>(
      `${apiBase}/options/autoload`,
    )
    if (res.code === 0 && res.data?.options) {
      const opts = res.data.options
      primary   = parseSaved(opts.theme_primary,   DEFAULTS.primary)
      neutral   = parseSaved(opts.theme_neutral,   DEFAULTS.neutral)
      radius    = parseSaved(opts.theme_radius,    DEFAULTS.radius)
      font      = parseSaved(opts.theme_font,      DEFAULTS.font)
      const fs  = parseSaved(opts.theme_font_size, String(DEFAULTS.fontSize))
      fontSize  = Number(fs) || DEFAULTS.fontSize
      customCss = parseSaved(opts.theme_custom_css, DEFAULTS.customCss)
    }
  } catch {
    // Backend unreachable — fall back to defaults silently
  }

  const fontStack = FONT_STACKS[font] ?? FONT_STACKS.system

  const parts = [
    ':root {',
    colorAliases('primary', primary),
    colorAliases('neutral', neutral),
    `  --ui-radius: ${radius};`,
    '}',
    'html {',
    `  font-size: ${fontSize}px;`,
    `  --blog-font-family: ${fontStack};`,
    '}',
  ]

  if (customCss?.trim()) {
    parts.push(customCss)
  }

  return parts.join('\n')
})
