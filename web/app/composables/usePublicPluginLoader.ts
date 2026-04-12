/**
 * Public Plugin Loader — fetches enabled plugins and loads their public_js scripts.
 *
 * Called once in public layouts (default, home, post). Handles:
 * 1. Fetching enabled plugins with public_js info from the public endpoint
 * 2. Registering contribution points in the contributions store
 * 3. Loading public_js scripts (official/local = inline, community = iframe sandbox)
 * 4. Creating per-plugin permission-scoped nuxtblogPublic API objects
 */
import type { PluginContributes } from '~/stores/plugin-contributions'
import { installNuxtblogPublic, eventBus } from '~/composables/useNuxtblogPublic'
import type { PublicPluginPermissions } from '~/composables/useNuxtblogPublic'

interface PluginClientItem {
  id: string
  title: string
  icon: string
  version: string
  trust_level: string
  public_js?: string
  contributes?: string // raw JSON
  permissions?: string // raw JSON
  pages?: string // raw JSON
}

interface PageDef {
  path: string
  slot: string
  component: string
  title?: string
  nav?: { slot?: string; icon?: string; order?: number }
}

// Module-level flag — safe because this only runs on the client (guarded by import.meta.client)
let _loaded = false

export function usePublicPluginLoader() {
  const plugins = ref<PluginClientItem[]>([])
  const contributionsStore = usePluginContributionsStore()
  const config = useRuntimeConfig()
  const apiBase = config.public.apiBase as string
  // Extract backend origin for asset URLs (e.g. "http://localhost:9000" from "http://localhost:9000/api/v1")
  const backendOrigin = apiBase.replace(/\/api\/v\d+$/, '')

  async function loadPlugins() {
    if (_loaded) return

    try {
      const res = await $fetch<{ code: number; data: { items: PluginClientItem[] } }>(`${apiBase}/plugins/client`)
      plugins.value = res?.data?.items || []

      if (plugins.value.length === 0) {
        _loaded = true
        return
      }

      // Install the nuxtblogPublic global API
      const { createPluginApi } = installNuxtblogPublic()

      for (const plugin of plugins.value) {
        // Register contribution points
        if (plugin.contributes) {
          try {
            const contributes: PluginContributes = JSON.parse(plugin.contributes)
            contributionsStore.registerPlugin(plugin.id, contributes)
          }
          catch (e) {
            console.warn(`[public-plugin-loader] failed to parse contributes for ${plugin.id}:`, e)
          }
        }

        // Register public plugin pages
        if (plugin.pages) {
          try {
            const pages: PageDef[] = JSON.parse(plugin.pages)
            for (const page of pages) {
              if (page.slot !== 'public') continue
              contributionsStore.registerPluginPage({
                pluginId: plugin.id,
                component: page.component,
                title: page.title,
                moduleFile: 'public.mjs',
              })
              // Register nav item if page has nav config
              if (page.nav) {
                contributionsStore.registerPlugin(plugin.id, {
                  navigation: [{
                    slot: page.nav.slot || 'public:header-actions',
                    title: page.title || page.component,
                    icon: page.nav.icon,
                    route: `/p/${encodeURIComponent(plugin.id)}/${encodeURIComponent(page.component)}`,
                    order: page.nav.order ?? 100,
                  }],
                })
              }
            }
          }
          catch (e) {
            console.warn(`[public-plugin-loader] failed to parse pages for ${plugin.id}:`, e)
          }
        }

        // Parse permissions
        let permissions: string[] = []
        if (plugin.permissions) {
          try { permissions = JSON.parse(plugin.permissions) }
          catch { /* empty */ }
        }

        // Load public_js with per-plugin API
        if (plugin.public_js) {
          const meta: PublicPluginPermissions = {
            pluginId: plugin.id,
            trustLevel: plugin.trust_level,
            permissions,
          }
          await loadPublicScript(plugin, meta, createPluginApi, backendOrigin)
        }
      }

      _loaded = true
    }
    catch (e) {
      console.warn('[public-plugin-loader] failed to load plugins:', e)
    }
  }

  return { loadPlugins, plugins }
}

/**
 * Load a plugin's public_js script. Loading method depends on trust_level:
 * - official/local: loaded as inline <script> with per-plugin nuxtblogPublic
 * - community: loaded in a sandboxed iframe with postMessage bridge
 */
async function loadPublicScript(
  plugin: PluginClientItem,
  meta: PublicPluginPermissions,
  createPluginApi: (meta: PublicPluginPermissions) => Record<string, any>,
  backendOrigin: string,
) {
  const assetFilename = plugin.public_js!.split('/').pop()!
  const scriptUrl = `${backendOrigin}/api/plugins/${encodeURIComponent(plugin.id)}/assets/${assetFilename}?v=${plugin.version}`

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
        script.onerror = () => reject(new Error(`Failed to load public_js for ${plugin.id}`))
      })
    }
    catch (e) {
      console.warn(`[public-plugin-loader] failed to load public_js for ${plugin.id}:`, e)
    }
  }
  else {
    loadInSandbox(plugin.id, scriptUrl, meta, backendOrigin)
  }
}

/**
 * Load a community plugin's public_js in a sandboxed iframe.
 * Only safe APIs (page, theme, notify, slots.render with text only) are available.
 */
function loadInSandbox(pluginId: string, scriptUrl: string, _meta: PublicPluginPermissions, _backendOrigin: string) {
  const contributionsStore = usePluginContributionsStore()
  const iframe = document.createElement('iframe')
  iframe.style.display = 'none'
  iframe.sandbox.add('allow-scripts')
  iframe.dataset.pluginId = pluginId

  const html = `
    <!DOCTYPE html>
    <html>
    <head><script>
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
