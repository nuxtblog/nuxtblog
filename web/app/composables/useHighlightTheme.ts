import { HIGHLIGHT_THEMES } from '~/utils/highlightThemes'

const CDNJS_BASE = 'https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.11.1/styles'

export function useHighlightTheme() {
  if (!import.meta.client) return

  const optionsStore = useOptionsStore()

  const themeId = computed(() => optionsStore.get('code_highlight_theme', 'github'))

  const applyTheme = (id: string) => {
    const theme = HIGHLIGHT_THEMES.find(t => t.id === id) ?? HIGHLIGHT_THEMES[0]!

    // Remove existing hljs theme links
    document.querySelectorAll('link[data-hljs-theme]').forEach(el => el.remove())

    // Create light theme link
    const lightLink = document.createElement('link')
    lightLink.rel = 'stylesheet'
    lightLink.href = `${CDNJS_BASE}/${theme.light}.min.css`
    lightLink.media = '(prefers-color-scheme: light)'
    lightLink.dataset.hljsTheme = 'light'
    document.head.appendChild(lightLink)

    // Create dark theme link
    const darkLink = document.createElement('link')
    darkLink.rel = 'stylesheet'
    darkLink.href = `${CDNJS_BASE}/${theme.dark}.min.css`
    darkLink.media = '(prefers-color-scheme: dark)'
    darkLink.dataset.hljsTheme = 'dark'
    document.head.appendChild(darkLink)

    // Also handle class-based dark mode by fetching CSS and injecting as style
    applyClassBasedTheme(theme.light, theme.dark)
  }

  const applyClassBasedTheme = async (lightName: string, darkName: string) => {
    try {
      const [lightCss, darkCss] = await Promise.all([
        fetch(`${CDNJS_BASE}/${lightName}.min.css`).then(r => r.text()),
        fetch(`${CDNJS_BASE}/${darkName}.min.css`).then(r => r.text()),
      ])

      let styleEl = document.getElementById('hljs-theme-class') as HTMLStyleElement | null
      if (!styleEl) {
        styleEl = document.createElement('style')
        styleEl.id = 'hljs-theme-class'
        document.head.appendChild(styleEl)
      }

      styleEl.textContent = `
:root:not(.dark) { ${lightCss} }
:root.dark { ${darkCss} }
`
    } catch {
      // Fallback: media query links still work
    }
  }

  watch(themeId, (id) => applyTheme(id), { immediate: true })
}
