/**
 * Dynamic plugin route registration (client-side only).
 *
 * Watches the plugin contributions store and registers routes via
 * router.addRoute(). These exact routes take priority over catch-all
 * pages during SPA navigation, providing a direct match without the
 * 404-check overhead in the catch-all fallback pages.
 *
 * For cold-start (direct URL / SSR), the catch-all pages
 * [...plugin-page].vue handle rendering.
 */
export default defineNuxtPlugin(() => {
  const router = useRouter()
  const store = usePluginContributionsStore()

  const registered = new Set<string>()

  watch(() => store.pluginPages, (pages) => {
    for (const page of pages) {
      if (!page.path) continue
      const name = `plugin:${page.pluginId}:${page.component}`
      if (registered.has(name)) continue
      registered.add(name)

      const meta = {
        pluginId: page.pluginId,
        pluginComponent: page.component,
        pluginModuleFile: page.moduleFile,
        pluginPageTitle: page.title,
      }
      const component = () => import('~/components/PluginPageRenderer.vue')

      // Exact path (e.g. /shop)
      router.addRoute({ name, path: page.path, component, meta })
      // Sub-paths (e.g. /shop/1, /shop/anything) — plugin component handles internal routing
      router.addRoute({
        name: `${name}:children`,
        path: `${page.path}/:pathMatch(.*)*`,
        component,
        meta,
      })
    }
  }, { deep: true })
})
