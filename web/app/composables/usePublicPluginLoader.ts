/**
 * Public Plugin Loader — fetches enabled plugins and loads their public scripts.
 *
 * Called once in public layouts (default, home, post). Handles:
 * 1. Fetching enabled plugins with public contributes info from the public endpoint
 * 2. Registering contribution points in the contributions store
 * 3. Loading public activation modules (official/local = inline, community = iframe sandbox)
 * 4. Creating per-plugin permission-scoped nuxtblogPublic API objects
 */
import type { PluginContributes } from '~/stores/plugin-contributions'
import { installNuxtblogPublic, eventBus } from '~/composables/useNuxtblogPublic'
import type { PublicPluginPermissions } from '~/composables/useNuxtblogPublic'
import { registerPluginVersion } from '~/composables/usePluginComponents'
import { registerPluginI18n, createPluginI18nApi, createPluginTComposable } from '~/composables/usePluginI18n'

interface PluginClientItem {
  id: string
  title: string
  icon: string
  version: string
  trust_level: string
  contributes?: string // raw JSON
  permissions?: string // raw JSON
  i18n?: string // raw JSON
}

interface PageDef {
  path?: string
  slot: string
  component?: string
  module?: string
  title?: string
  nav?: { slot?: string; icon?: string; order?: number }
}

interface ActivationEntry {
  scope: string // admin | public
  module: string
}

// Module-level flag — safe because this only runs on the client (guarded by import.meta.client)
let _loaded = false
// Cached createPluginApi — installed once in setup sync context, reused in async _doLoad()
let _createPluginApi: ((meta: PublicPluginPermissions) => Record<string, any>) | null = null
/** Reactive flag — true once public plugins have finished loading. */
export const publicPluginsLoaded = ref(false)
/** Resolves when public plugins have been loaded (or immediately if already done). */
export let pluginsLoadedPromise: Promise<void> | null = null

export function usePublicPluginLoader() {
  const plugins = ref<PluginClientItem[]>([])
  const contributionsStore = usePluginContributionsStore()
  const { locale } = useI18n()
  // Expose i18n factory for plugin modules (source rewriting injects the call)
  ;(window as any).__nuxtblog_i18n = (pluginId: string) =>
    createPluginTComposable(pluginId, locale)
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase as string
  // Extract backend origin for asset URLs (e.g. "http://localhost:9000" from "http://localhost:9000/api/v1")
  const backendOrigin = apiBase.replace(/\/api\/v\d+$/, '')

  // Install nuxtblogPublic in setup sync context (only once) — must happen before any await
  // so that Vue composables (useRoute, useColorMode, useToast, etc.) have access to the component instance.
  if (!_createPluginApi) {
    const { createPluginApi } = installNuxtblogPublic()
    _createPluginApi = createPluginApi
  }

  async function loadPlugins() {
    if (_loaded) return
    if (pluginsLoadedPromise) return pluginsLoadedPromise
    pluginsLoadedPromise = _doLoad()
    return pluginsLoadedPromise
  }

  async function _doLoad() {
    try {
      const res = await $fetch<{ code: number; data: { items: PluginClientItem[] } }>(`${apiBase}/plugins/client`)
      plugins.value = res?.data?.items || []

      if (plugins.value.length === 0) {
        _loaded = true
        return
      }

      for (const plugin of plugins.value) {
        // Register version for cache busting
        registerPluginVersion(plugin.id, plugin.version)

        // Parse and register i18n messages
        if (plugin.i18n) {
          try {
            const i18nMessages = JSON.parse(plugin.i18n)
            registerPluginI18n(plugin.id, i18nMessages)
          }
          catch { /* ignore */ }
        }

        // Parse contributes (now includes pages, activation, styles)
        let contributes: PluginContributes = {}
        if (plugin.contributes) {
          try {
            contributes = JSON.parse(plugin.contributes)
          }
          catch (e) {
            console.warn(`[public-plugin-loader] failed to parse contributes for ${plugin.id}:`, e)
          }
        }

        // Extract public pages from contributes.pages
        const allPages: PageDef[] = contributes.pages || []
        const parsedPages = allPages.filter(p => p.slot === 'public')
        const pageNavItems = parsedPages
          .filter(p => p.nav)
          .map(p => ({
            slot: p.nav!.slot || 'public:header-actions',
            title: p.title || p.component || '',
            icon: p.nav!.icon,
            route: p.path || `/p/${encodeURIComponent(plugin.id)}/${encodeURIComponent(p.component || '')}`,
            order: p.nav!.order ?? 100,
          }))
        if (pageNavItems.length > 0) {
          contributes.navigation = [...(contributes.navigation || []), ...pageNavItems]
        }

        // Single registerPlugin call (calls unregisterPlugin internally, which clears pluginPages)
        if (contributes.commands?.length || contributes.navigation?.length
          || contributes.menus && Object.keys(contributes.menus).length
          || contributes.views && Object.keys(contributes.views).length) {
          contributionsStore.registerPlugin(plugin.id, contributes)
        }

        // Register public pages AFTER registerPlugin to avoid being wiped by unregisterPlugin()
        for (const page of parsedPages) {
          if (!page.path || !page.component) continue
          contributionsStore.registerPluginPage({
            pluginId: plugin.id,
            component: page.component,
            title: page.title,
            path: page.path,
            moduleFile: page.module || 'public.mjs',
          })
        }

        // Parse permissions
        let permissions: string[] = []
        if (plugin.permissions) {
          try { permissions = JSON.parse(plugin.permissions) }
          catch { /* empty */ }
        }

        // Load public activation module
        const activations: ActivationEntry[] = contributes.activation || []
        const publicActivation = activations.find(a => a.scope === 'public')
        if (publicActivation) {
          const meta: PublicPluginPermissions = {
            pluginId: plugin.id,
            trustLevel: plugin.trust_level,
            permissions,
          }
          await loadPublicScript(plugin, publicActivation.module, meta, _createPluginApi!, backendOrigin, plugin.i18n)
        }
      }

      _loaded = true
      publicPluginsLoaded.value = true
    }
    catch (e) {
      console.warn('[public-plugin-loader] failed to load plugins:', e)
    }
  }

  return { loadPlugins, plugins }
}

/**
 * Load a plugin's public activation module. Loading method depends on trust_level:
 * - official/local: loaded as inline <script> with per-plugin nuxtblogPublic
 * - community: loaded in a sandboxed iframe with postMessage bridge
 */
async function loadPublicScript(
  plugin: PluginClientItem,
  moduleFile: string,
  meta: PublicPluginPermissions,
  createPluginApi: (meta: PublicPluginPermissions) => Record<string, any>,
  backendOrigin: string,
  i18nJSON?: string,
) {
  const scriptUrl = `${backendOrigin}/api/plugins/${encodeURIComponent(plugin.id)}/assets/${moduleFile}?v=${plugin.version}`

  if (plugin.trust_level === 'official' || plugin.trust_level === 'local') {
    try {
      // Set per-plugin API on window before loading the script
      const pluginApi = createPluginApi(meta)
      ;(window as any).nuxtblogPublic = pluginApi

      const script = document.createElement('script')
      script.type = 'module'
      script.src = scriptUrl
      script.dataset.pluginId = plugin.id
      document.head.appendChild(script)

      await new Promise<void>((resolve, reject) => {
        script.onload = () => resolve()
        script.onerror = () => reject(new Error(`Failed to load public module for ${plugin.id}`))
      })
    }
    catch (e) {
      console.warn(`[public-plugin-loader] failed to load public module for ${plugin.id}:`, e)
    }
  }
  else {
    loadInSandbox(plugin.id, scriptUrl, meta, backendOrigin, i18nJSON)
  }
}

/**
 * Load a community plugin's public script in a sandboxed iframe.
 * Only safe APIs (page, theme, notify, slots.render with text only) are available.
 */
function loadInSandbox(pluginId: string, scriptUrl: string, _meta: PublicPluginPermissions, _backendOrigin: string, i18nJSON?: string) {
  const contributionsStore = usePluginContributionsStore()
  const iframe = document.createElement('iframe')
  iframe.style.display = 'none'
  iframe.sandbox.add('allow-scripts')
  iframe.dataset.pluginId = pluginId

  const i18nMsgs = i18nJSON || '{}'
  const html = `
    <!DOCTYPE html>
    <html>
    <head><script>
    window.__i18n = { messages: ${i18nMsgs}, locale: navigator.language.startsWith('zh') ? 'zh' : 'en' };
    window.nuxtblogPublic = {
      page: {
        getRoute: function() { return window.__route || {}; },
        getPost: function() { return window.__post || null; },
        getType: function() { return window.__pageType || 'other'; }
      },
      theme: {
        getMode: function() { return window.__themeMode || 'light'; },
        onChange: function(cb) {
          window.__themeCallbacks = window.__themeCallbacks || [];
          window.__themeCallbacks.push(cb);
          return { dispose: function() { var i = window.__themeCallbacks.indexOf(cb); if (i >= 0) window.__themeCallbacks.splice(i, 1); } };
        }
      },
      notify: {
        success: function(msg) { parent.postMessage({type:'plugin:notify',pluginId:'${pluginId}',level:'success',message:msg},'*'); },
        error: function(msg) { parent.postMessage({type:'plugin:notify',pluginId:'${pluginId}',level:'error',message:msg},'*'); },
        info: function(msg) { parent.postMessage({type:'plugin:notify',pluginId:'${pluginId}',level:'info',message:msg},'*'); }
      },
      slots: {
        render: function(slotId, content, options) {
          parent.postMessage({type:'plugin:slotRender',pluginId:'${pluginId}',slotId:slotId,content:String(content),order:(options&&options.order)||100},'*');
        }
      },
      i18n: {
        t: function(key, fallback) {
          var m = window.__i18n.messages;
          var l = window.__i18n.locale;
          return (m[l] && m[l][key]) || (m['zh'] && m['zh'][key]) || fallback || key;
        },
        get locale() { return window.__i18n.locale; }
      },
      events: {
        emit: function(eventName) {
          if (typeof eventName !== 'string' || eventName.indexOf('${pluginId}:') !== 0) return;
          var args = Array.prototype.slice.call(arguments, 1);
          parent.postMessage({type:'plugin:eventEmit',pluginId:'${pluginId}',eventName:eventName,args:args},'*');
        },
        on: function(eventName, handler) {
          window.__eventHandlers = window.__eventHandlers || {};
          window.__eventHandlers[eventName] = window.__eventHandlers[eventName] || [];
          window.__eventHandlers[eventName].push(handler);
          parent.postMessage({type:'plugin:eventOn',pluginId:'${pluginId}',eventName:eventName},'*');
          return function() {
            var arr = window.__eventHandlers[eventName];
            if (arr) { var i = arr.indexOf(handler); if (i >= 0) arr.splice(i, 1); }
          };
        },
        once: function(eventName, handler) {
          var unsub;
          var wrapper = function() { if (unsub) unsub(); handler.apply(null, arguments); };
          unsub = window.nuxtblogPublic.events.on(eventName, wrapper);
          return unsub;
        }
      }
    };
    window.addEventListener('message', function(e) {
      if (!e.data) return;
      if (e.data.type === 'themeUpdate' && window.__themeCallbacks) {
        window.__themeCallbacks.forEach(function(cb) { try { cb(e.data.mode); } catch(ex) {} });
      }
      if (e.data.type === 'plugin:eventDispatch' && window.__eventHandlers) {
        var handlers = window.__eventHandlers[e.data.eventName];
        if (handlers) handlers.forEach(function(h) { try { h.apply(null, e.data.args || []); } catch(ex) {} });
      }
      if (e.data.type === 'localeChange') {
        window.__i18n.locale = e.data.locale;
      }
    });
    <\/script>
    <script type="module" src="${scriptUrl}"><\/script>
    </head><body></body></html>
  `

  iframe.srcdoc = html
  document.body.appendChild(iframe)

  // Listen for postMessage from sandboxed iframe
  window.addEventListener('message', (e) => {
    if (e.source !== iframe.contentWindow) return
    const data = e.data
    if (!data || !data.type || data.pluginId !== pluginId) return

    switch (data.type) {
      case 'plugin:notify': {
        const toast = useToast()
        toast.add({ title: data.message, color: data.level || 'info' })
        break
      }
      case 'plugin:slotRender': {
        contributionsStore.registerContentBlock({
          pluginId,
          slot: data.slotId,
          content: data.content,
          trustLevel: 'community',
          hasHtmlPermission: false, // community plugins never get html:inject
          order: data.order ?? 100,
        })
        break
      }
      case 'plugin:eventEmit': {
        // Community plugin emitting an event — prefix already validated in iframe
        if (data.eventName?.startsWith(`${pluginId}:`)) {
          eventBus.emit(data.eventName, ...(data.args || []))
        }
        break
      }
      case 'plugin:eventOn': {
        // Community plugin subscribing — forward matching events back to iframe
        if (data.eventName) {
          eventBus.on(data.eventName, (...args: any[]) => {
            iframe.contentWindow?.postMessage(
              { type: 'plugin:eventDispatch', eventName: data.eventName, args },
              '*',
            )
          })
        }
        break
      }
    }
  })
}
