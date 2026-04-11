import { DEFAULT_THEME, PRIMARY_COLORS, NEUTRAL_COLORS } from '~/config/theme'

/**
 * 客户端插件：应用半径、字体、颜色模式、自定义 CSS 等无法在 SSR 阶段处理的设置。
 * 颜色（primary/neutral）已由 theme-early.ts 和 theme-head.server.ts 处理，
 * 此处再次调用 applyAll 是幂等的（同值不会触发任何视觉变化）。
 *
 * 同时将当前主题配置同步到 window.__nuxtblog_theme 响应式对象，供插件读取。
 */
function findHex(colors: { name: string; hex: string }[], name: string): string {
  return colors.find(c => c.name === name)?.hex ?? ''
}

export default defineNuxtPlugin({
  name: 'blog-theme',
  setup() {
    const optionsStore = useOptionsStore()
    const colorMode = useColorMode()
    const { applyAll, applyHeadCode, applyBodyCode } = useTheme()

    const apply = () => {
      const primary  = optionsStore.get('theme_primary',    DEFAULT_THEME.primary)
      const neutral  = optionsStore.get('theme_neutral',    DEFAULT_THEME.neutral)
      const radius   = optionsStore.get('theme_radius',     DEFAULT_THEME.radius)
      const font     = optionsStore.get('theme_font',       DEFAULT_THEME.font)
      const fontSize = Number(optionsStore.get('theme_font_size', String(DEFAULT_THEME.fontSize))) || DEFAULT_THEME.fontSize

      applyAll({
        primary,
        neutral,
        radius,
        colorMode: optionsStore.get('theme_color_mode', DEFAULT_THEME.colorMode),
        font,
        fontSize,
        customCss: optionsStore.get('theme_custom_css', DEFAULT_THEME.customCss),
      })
      applyHeadCode(optionsStore.get('theme_head_code', ''))
      applyBodyCode(optionsStore.get('theme_body_code', ''))

      // Sync theme tokens for plugins
      const themeTokens = (window as any).__nuxtblog_theme
      if (themeTokens) {
        Object.assign(themeTokens, {
          primary,
          neutral,
          primaryHex: findHex(PRIMARY_COLORS, primary),
          neutralHex: findHex(NEUTRAL_COLORS, neutral),
          radius,
          font,
          fontSize,
        })
      }
    }

    // Keep resolved colorMode in sync (never 'system', always 'light' or 'dark')
    watch(() => colorMode.value, (mode) => {
      const themeTokens = (window as any).__nuxtblog_theme
      if (themeTokens) {
        themeTokens.colorMode = mode === 'dark' ? 'dark' : 'light'
      }
    }, { immediate: true })

    if (optionsStore.loaded) {
      apply()
    } else {
      optionsStore.load().then(apply)
    }
  },
})
