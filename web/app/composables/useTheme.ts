import type { ThemeSettings } from '~/config/theme'
import { FONT_OPTIONS, DEFAULT_THEME } from '~/config/theme'

/**
 * useTheme — applies NuxtUI theme settings at runtime.
 *
 * Color switching:   via useAppConfig() — NuxtUI watches this and regenerates CSS vars
 * Border radius:     via CSS custom property --ui-radius on <html>
 * Color mode:        via useColorMode() from @nuxtjs/color-mode
 * Font / font-size:  via CSS custom properties on <html>
 * Custom CSS:        via an injected <style id="blog-custom-css"> tag
 */
export const useTheme = () => {
  const appConfig  = useAppConfig()
  const colorMode  = useColorMode()

  const applyPrimary = (color: string) => {
    appConfig.ui.colors.primary = color
  }

  const applyNeutral = (color: string) => {
    appConfig.ui.colors.neutral = color
  }

  const applyRadius = (radius: string) => {
    if (import.meta.client) {
      document.documentElement.style.setProperty('--ui-radius', radius)
    }
  }

  const applyColorMode = (mode: string) => {
    colorMode.preference = mode
  }

  const applyFont = (fontValue: string, fontSize: number) => {
    if (!import.meta.client) return
    const opt = FONT_OPTIONS.find((f) => f.value === fontValue)
    if (opt) {
      document.documentElement.style.setProperty('--blog-font-family', opt.stack)
    }
    // Set font-size on <html> so rem-based Tailwind utilities scale correctly
    document.documentElement.style.fontSize = `${fontSize}px`
  }

  const applyCustomCss = (css: string) => {
    if (!import.meta.client) return
    let el = document.getElementById('blog-custom-css') as HTMLStyleElement | null
    if (!el) {
      el = document.createElement('style')
      el.id = 'blog-custom-css'
      document.head.appendChild(el)
    }
    el.textContent = css
  }

  /**
   * Inject arbitrary HTML (script/style/link/meta) into a target element.
   * Handles <script> tags properly — innerHTML does NOT execute scripts.
   */
  const injectCode = (html: string, id: string, target: 'head' | 'body') => {
    if (!import.meta.client || !html.trim()) return
    // Remove previous injections
    document.querySelectorAll(`[data-blog-inject="${id}"]`).forEach((el) => el.remove())

    const targetEl = target === 'head' ? document.head : document.body
    // Wrap in body so DOMParser doesn't add its own <meta charset> to our node list
    const doc = new DOMParser().parseFromString(`<body>${html}</body>`, 'text/html')

    for (const node of [...doc.body.childNodes]) {
      let el: Node
      if ((node as Element).tagName === 'SCRIPT') {
        // Scripts cloned via innerHTML/cloneNode are NOT executed — must recreate
        const parsed = node as HTMLScriptElement
        const script = document.createElement('script')
        for (const attr of parsed.attributes) {
          script.setAttribute(attr.name, attr.value)
        }
        script.textContent = parsed.textContent
        el = script
      } else {
        el = node.cloneNode(true)
      }
      if ((el as Element).setAttribute) {
        (el as Element).setAttribute('data-blog-inject', id)
      }
      targetEl.appendChild(el)
    }
  }

  const applyHeadCode = (html: string) => injectCode(html, 'blog-head-code', 'head')
  const applyBodyCode = (html: string) => injectCode(html, 'blog-body-code', 'body')

  /** Apply all settings at once (used on startup and on save). */
  const applyAll = (settings: Partial<ThemeSettings>) => {
    const s = { ...DEFAULT_THEME, ...settings }
    applyPrimary(s.primary)
    applyNeutral(s.neutral)
    applyRadius(s.radius)
    applyColorMode(s.colorMode)
    applyFont(s.font, s.fontSize)
    applyCustomCss(s.customCss)
  }

  return { applyPrimary, applyNeutral, applyRadius, applyColorMode, applyFont, applyCustomCss, applyHeadCode, applyBodyCode, applyAll }
}
