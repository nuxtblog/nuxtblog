/**
 * enforce: 'pre' 插件：在 NuxtUI 的 colors.js 模块插件初始化之前，
 * 从 SSR payload 恢复的 useState 里读取保存的颜色，设置到 appConfig。
 *
 * 执行顺序（Nuxt 4）：
 *   1. enforce:'pre' 用户插件（本文件）
 *   2. 模块插件（NuxtUI colors.js）← 此时 appConfig 已经是正确颜色
 *   3. 普通用户插件（options.ts、theme.client.ts 等）
 *
 * 为什么不用 useOptionsStore()：
 *   Pinia 是模块插件（第 2 步），enforce:'pre' 时可能尚未安装。
 *   直接访问 useState('options:raw') 是 Nuxt 核心功能，无需 Pinia。
 *   服务端 SSR 已经把 useState 序列化到 payload，客户端水合时立即可用。
 */
export default defineNuxtPlugin({
  name: 'blog-theme-early',
  enforce: 'pre',
  setup() {
    // 直接读 options store 的底层 useState（与 useOptionsStore 共享同一 key）
    const raw = useState<Record<string, string>>('options:raw', () => ({}))

    const get = (key: string, fallback: string): string => {
      const v = raw.value[key]
      if (!v) return fallback
      try {
        const parsed = JSON.parse(v)
        return typeof parsed === 'string' ? parsed : fallback
      } catch {
        return fallback
      }
    }

    const appConfig = useAppConfig()
    appConfig.ui.colors.primary = get('theme_primary', 'violet')
    appConfig.ui.colors.neutral = get('theme_neutral', 'zinc')
  },
})
