/**
 * nuxtblogPublic — global object exposed to plugin public_js scripts.
 *
 * Provides: page info, theme, notifications, slot rendering.
 * Sensitive APIs (user, http) require explicit permission declaration in plugin.yaml.
 */
import type { PluginContentBlock } from '~/stores/plugin-contributions'

export interface PublicPluginPermissions {
  pluginId: string
  trustLevel: string
  permissions: string[]
}

/**
 * Install the nuxtblogPublic global object on window.
 * Returns a factory to create per-plugin API objects with permission-scoped access.
 */
// ── Module-level event bus (shared across all plugins) ──────────────
const eventMap = new Map<string, Set<(...args: any[]) => void>>()

function eventOn(eventName: string, handler: (...args: any[]) => void): () => void {
  let handlers = eventMap.get(eventName)
  if (!handlers) {
    handlers = new Set()
    eventMap.set(eventName, handlers)
  }
  handlers.add(handler)
  return () => { handlers!.delete(handler) }
}

function eventEmit(eventName: string, ...args: any[]) {
  const handlers = eventMap.get(eventName)
  if (!handlers) return
  for (const handler of handlers) {
    try { handler(...args) }
    catch (e) { console.warn('[nuxtblogPublic] event handler error:', e) }
  }
}

/** Exposed for sandbox bridge in usePublicPluginLoader */
export const eventBus = { on: eventOn, emit: eventEmit }

export function installNuxtblogPublic() {
  const route = useRoute()
  const colorMode = useColorMode()
  const toast = useToast()
  const contributionsStore = usePluginContributionsStore()

  const themeCallbacks = new Set<(mode: string) => void>()

  // Watch for theme changes
  watch(() => colorMode.value, (mode) => {
    for (const cb of themeCallbacks) {
      try { cb(mode) }
      catch (e) { console.warn('[nuxtblogPublic] theme callback error:', e) }
    }
  })

  /** Create a permission-scoped API object for a specific plugin. */
  function createPluginApi(meta: PublicPluginPermissions) {
    const hasPermission = (perm: string) => meta.permissions.includes(perm)

    const api: Record<string, any> = {
      // ── Always available (no permission needed) ──────────────
      page: {
        getRoute() {
          return {
            path: route.path,
            params: { ...route.params },
            query: { ...route.query },
            name: route.name,
          }
        },
        getPost() {
          // Read from body data attributes set by the post page
          const body = document.body
          const id = body.dataset.pageId
          const slug = body.dataset.pageSlug
          if (!id) return null
          return { id, slug }
        },
        getType() {
          const path = route.path
          if (path === '/') return 'home'
          if (path.startsWith('/posts/')) return 'post'
          if (path.startsWith('/archive')) return 'archive'
          if (path.startsWith('/page/')) return 'page'
          if (path.startsWith('/docs')) return 'doc'
          return 'other'
        },
      },

      theme: {
        getMode() { return colorMode.value },
        onChange(cb: (mode: string) => void) {
          themeCallbacks.add(cb)
          return { dispose: () => { themeCallbacks.delete(cb) } }
        },
      },

      notify: {
        success(msg: string) { toast.add({ title: msg, color: 'success' }) },
        error(msg: string) { toast.add({ title: msg, color: 'error' }) },
        info(msg: string) { toast.add({ title: msg, color: 'info' }) },
      },

      events: {
        emit(eventName: string, ...args: any[]) {
          if (!eventName.startsWith(`${meta.pluginId}:`)) {
            console.warn(`[nuxtblogPublic] plugin "${meta.pluginId}" can only emit events prefixed with "${meta.pluginId}:"`)
            return
          }
          eventEmit(eventName, ...args)
        },
        on(eventName: string, handler: (...args: any[]) => void): () => void {
          return eventOn(eventName, handler)
        },
        once(eventName: string, handler: (...args: any[]) => void): () => void {
          const unsub = eventOn(eventName, (...args: any[]) => {
            unsub()
            handler(...args)
          })
          return unsub
        },
      },

      slots: {
        /**
         * Render content into a named slot.
         * @param content - string (text/HTML) or function (DOM render callback)
         *   When a function is passed, it receives the container HTMLElement and
         *   may return a cleanup function. Only official/local plugins can use functions.
         */
        render(slotId: string, content: string | ((container: HTMLElement) => void | (() => void)), options?: { order?: number }) {
          const isRenderFn = typeof content === 'function'
          const block: PluginContentBlock = {
            pluginId: meta.pluginId,
            slot: slotId,
            content: isRenderFn ? '' : content,
            trustLevel: meta.trustLevel,
            hasHtmlPermission: hasPermission('html:inject'),
            order: options?.order ?? 100,
            renderFn: isRenderFn && (meta.trustLevel === 'official' || meta.trustLevel === 'local')
              ? content as (container: HTMLElement) => void | (() => void)
              : undefined,
          }
          contributionsStore.registerContentBlock(block)
        },
      },
    }

    // ── Permission-gated APIs ──────────────────────────────────

    if (hasPermission('user:read')) {
      const authStore = useAuthStore()
      api.user = {
        isLoggedIn() { return authStore.isLoggedIn },
        getInfo() {
          if (!authStore.isLoggedIn || !authStore.user) return null
          return {
            id: authStore.user.id,
            displayName: authStore.user.display_name || authStore.user.username,
            avatar: authStore.user.avatar,
          }
        },
      }
    }

    if (hasPermission('http:auth')) {
      const { apiFetch } = useApiFetch()
      api.http = {
        async get(path: string) {
          try {
            const data = await apiFetch<unknown>(path)
            return { ok: true, data }
          }
          catch (e: unknown) {
            return { ok: false, data: null, error: String(e) }
          }
        },
        async post(path: string, body: object) {
          try {
            const data = await apiFetch<unknown>(path, { method: 'POST', body })
            return { ok: true, data }
          }
          catch (e: unknown) {
            return { ok: false, data: null, error: String(e) }
          }
        },
      }
    }

    return api
  }

  // Install a base (no-permission) version on window for simple scripts
  const baseApi = createPluginApi({
    pluginId: '__default__',
    trustLevel: 'community',
    permissions: [],
  })
  ;(window as any).nuxtblogPublic = baseApi

  return { createPluginApi }
}
