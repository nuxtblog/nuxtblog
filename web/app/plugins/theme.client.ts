import { DEFAULT_THEME } from '~/config/theme'

/**
 * 客户端插件：应用半径、字体、颜色模式、自定义 CSS 等无法在 SSR 阶段处理的设置。
 * 颜色（primary/neutral）已由 theme-early.ts 和 theme-head.server.ts 处理，
 * 此处再次调用 applyAll 是幂等的（同值不会触发任何视觉变化）。
 */
export default defineNuxtPlugin({
  name: 'blog-theme',
  setup() {
    const optionsStore = useOptionsStore()
    const { applyAll, applyHeadCode, applyBodyCode } = useTheme()

    const apply = () => {
      applyAll({
        primary:   optionsStore.get('theme_primary',    DEFAULT_THEME.primary),
        neutral:   optionsStore.get('theme_neutral',    DEFAULT_THEME.neutral),
        radius:    optionsStore.get('theme_radius',     DEFAULT_THEME.radius),
        colorMode: optionsStore.get('theme_color_mode', DEFAULT_THEME.colorMode),
        font:      optionsStore.get('theme_font',       DEFAULT_THEME.font),
        fontSize:  Number(optionsStore.get('theme_font_size', String(DEFAULT_THEME.fontSize))) || DEFAULT_THEME.fontSize,
        customCss: optionsStore.get('theme_custom_css', DEFAULT_THEME.customCss),
      })
      applyHeadCode(optionsStore.get('theme_head_code', ''))
      applyBodyCode(optionsStore.get('theme_body_code', ''))
    }

    if (optionsStore.loaded) {
      apply()
    } else {
      optionsStore.load().then(apply)
    }
  },
})
