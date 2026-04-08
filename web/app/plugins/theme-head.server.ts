/**
 * 服务端插件：在 SSR 渲染前设置 appConfig.ui.colors，
 * 使 NuxtUI 的 colors.js 插件生成正确颜色的 CSS 写入 HTML。
 *
 * 执行顺序：options.ts（字母序 o < t）先加载完选项，
 * 本插件再读取并修改 appConfig → SSR HTML 里 <style id="nuxt-ui-colors">
 * 直接包含用户保存的颜色，客户端首次绘制就是正确颜色。
 */
export default defineNuxtPlugin(() => {
  const opts = useOptionsStore()
  const appConfig = useAppConfig()

  const primary = opts.get('theme_primary', 'violet')
  const neutral = opts.get('theme_neutral', 'zinc')

  appConfig.ui.colors.primary = primary
  appConfig.ui.colors.neutral = neutral
})
