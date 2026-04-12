/**
 * Phase 2.4: Plugin Loader — fetches enabled plugins and loads their admin_js scripts.
 *
 * Called once in the admin layout. Handles:
 * 1. Fetching enabled plugins with their contributes/admin_js info
 * 2. Registering contribution points in the contributions store
 * 3. Loading admin_js scripts (official/local = inline, community = iframe sandbox)
 */
import type { PluginContributes } from '~/stores/plugin-contributions'
import { usePluginContributionsStore } from '~/stores/plugin-contributions'
import { loadPluginModule, registerPluginVersion } from '~/composables/usePluginComponents'

interface PluginClientItem {
  id: string
  title: string
  icon: string
  version: string
  trust_level: string
  admin_js?: string
  public_js?: string
  contributes?: string // raw JSON
  pages?: string // raw JSON — Phase 4.2
}

interface PageDef {
  path: string
  slot: string
  component: string
  title?: string
  nav?: { group?: string; icon?: string; order?: number }
}

export function usePluginLoader() {
  const { apiFetch } = useApiFetch()
  const contributionsStore = usePluginContributionsStore()
  const loaded = ref(false)
  const plugins = ref<PluginClientItem[]>([])

  async function loadPlugins() {
    if (loaded.value) return

    try {
      const res = await apiFetch<{ items: PluginClientItem[] }>('/admin/plugins/client')
      plugins.value = res?.items || []

      // Register plugin versions for cache busting, then contribution points
      for (const plugin of plugins.value) {
        registerPluginVersion(plugin.id, plugin.version)
        // Merge contributes + page navigation into a single registerPlugin call
        // to avoid the second call wiping out the first via unregisterPlugin().
        let contributes: PluginContributes = {}

        if (plugin.contributes) {
          try {
            contributes = JSON.parse(plugin.contributes)
          }
          catch (e) {
            console.warn(`[plugin-loader] failed to parse contributes for ${plugin.id}:`, e)
          }
        }

        // Phase 4.2: merge page nav items into the same contributes object
        if (plugin.pages) {
          try {
            const pages: PageDef[] = JSON.parse(plugin.pages)
            const pageNavItems = pages
              .filter(p => p.slot === 'admin' && p.nav)
              .map(p => ({
                slot: 'admin:sidebar-nav',
                parent: p.nav?.group,
                title: p.title || p.component,
                icon: p.nav?.icon,
                route: `/admin/plugin-page/${encodeURIComponent(plugin.id)}/${encodeURIComponent(p.component)}`,
                order: p.nav?.order ?? 100,
              }))
            if (pageNavItems.length > 0) {
              contributes.navigation = [...(contributes.navigation || []), ...pageNavItems]
            }
          }
          catch (e) {
            console.warn(`[plugin-loader] failed to parse pages for ${plugin.id}:`, e)
          }
        }

        // Single registration: contributes + page nav merged
        if (contributes.commands?.length || contributes.navigation?.length
          || contributes.menus && Object.keys(contributes.menus).length
          || contributes.views && Object.keys(contributes.views).length) {
          contributionsStore.registerPlugin(plugin.id, contributes)
        }

        // Load admin_js
        if (plugin.admin_js) {
          await loadAdminScript(plugin)
        }
      }

      loaded.value = true
    }
    catch (e) {
      console.warn('[plugin-loader] failed to load plugins:', e)
    }
  }

  return { loadPlugins, loaded, plugins }
}

/**
 * Load a plugin's admin_js script. The loading method depends on trust_level:
 * - official/local: loaded via loadPluginModule (source rewriting + activate())
 * - community: loaded in a sandboxed iframe with postMessage bridge
 */
async function loadAdminScript(plugin: PluginClientItem) {
  // admin_js may include a directory prefix (e.g. "dist/admin.mjs"), but the
  // backend saves assets as flat filenames. Strip any directory prefix.
  const assetFilename = plugin.admin_js!.split('/').pop()!

  if (plugin.trust_level === 'official' || plugin.trust_level === 'local') {
    // Main context — load via loadPluginModule (source rewriting + Blob URL import)
    try {
      const mod = await loadPluginModule(plugin.id, assetFilename, plugin.version)
      if (typeof mod.activate === 'function') {
        await mod.activate((window as any).nuxtblogAdmin)
      }
    }
    catch (e) {
      console.warn(`[plugin-loader] failed to load admin_js for ${plugin.id}:`, e)
    }
  }
  else {
    // Community plugins: sandboxed iframe
    // The iframe only has access to a restricted postMessage-based API
    const scriptUrl = `/api/plugins/${encodeURIComponent(plugin.id)}/assets/${assetFilename}?v=${plugin.version}`
    loadInSandbox(plugin.id, scriptUrl)
  }
}

/**
 * Load a community plugin's admin_js in a sandboxed iframe.
 * Only suggest() is available (no set()), no cookie/DOM access.
 */
function loadInSandbox(pluginId: string, scriptUrl: string) {
  const iframe = document.createElement('iframe')
  iframe.style.display = 'none'
  iframe.sandbox.add('allow-scripts')
  iframe.dataset.pluginId = pluginId

  // Create a minimal HTML page that loads the script and bridges postMessage
  const html = `
    <!DOCTYPE html>
    <html>
    <head><script>
    // Restricted nuxtblogAdmin for sandboxed plugins
    window.nuxtblogAdmin = {
      watch: function(field, cb) {
        var id = Math.random().toString(36).slice(2);
        window.__watchers = window.__watchers || {};
        window.__watchers[id] = cb;
        parent.postMessage({type:'plugin:watch',pluginId:'${pluginId}',field:field,watchId:id},'*');
        return {dispose:function(){delete window.__watchers[id];parent.postMessage({type:'plugin:unwatch',watchId:id},'*')}};
      },
      suggest: function(field, value) {
        parent.postMessage({type:'plugin:suggest',pluginId:'${pluginId}',field:field,value:value},'*');
      },
      getPost: function() { return window.__currentPost || {}; },
      commands: {
        register: function(id, handler) {
          window.__commandHandlers = window.__commandHandlers || {};
          window.__commandHandlers[id] = handler;
          parent.postMessage({type:'plugin:registerCommand',pluginId:'${pluginId}',commandId:id},'*');
          return {dispose:function(){delete window.__commandHandlers[id]}};
        }
      },
      http: {
        get: function(path) { return fetch(path).then(function(r){return r.json()}).then(function(d){return {ok:true,data:d}}).catch(function(e){return {ok:false,error:String(e)}}); },
        post: function(path, body) { return fetch(path,{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)}).then(function(r){return r.json()}).then(function(d){return {ok:true,data:d}}).catch(function(e){return {ok:false,error:String(e)}}); }
      },
      notify: {
        success: function(msg) { parent.postMessage({type:'plugin:notify',level:'success',message:msg},'*'); },
        error: function(msg) { parent.postMessage({type:'plugin:notify',level:'error',message:msg},'*'); },
        info: function(msg) { parent.postMessage({type:'plugin:notify',level:'info',message:msg},'*'); }
      }
    };
    // Listen for messages from parent
    window.addEventListener('message', function(e) {
      if (e.data && e.data.type === 'fieldUpdate' && window.__watchers) {
        var cb = window.__watchers[e.data.watchId];
        if (cb) cb(e.data.value);
      }
      if (e.data && e.data.type === 'plugin:executeCommand' && window.__commandHandlers) {
        var h = window.__commandHandlers[e.data.commandId];
        if (h) h(e.data.ctx);
      }
    });
    <\/script>
    <script type="module" src="${scriptUrl}"><\/script>
    </head><body></body></html>
  `

  iframe.srcdoc = html
  document.body.appendChild(iframe)

  // Track field watchers so we can dispose them on unwatch
  const sandboxWatchers = new Map<string, { dispose: () => void }>()

  // Listen for postMessage from the sandboxed iframe
  window.addEventListener('message', (e) => {
    if (e.source !== iframe.contentWindow) return
    const data = e.data
    if (!data || !data.type) return

    switch (data.type) {
      case 'plugin:watch': {
        const disposable = (window as any).nuxtblogAdmin?.watch(data.field, (value: string) => {
          iframe.contentWindow?.postMessage(
            { type: 'fieldUpdate', watchId: data.watchId, value },
            '*',
          )
        })
        if (disposable) {
          sandboxWatchers.set(data.watchId, disposable)
        }
        break
      }
      case 'plugin:unwatch': {
        const disposable = sandboxWatchers.get(data.watchId)
        if (disposable) {
          disposable.dispose()
          sandboxWatchers.delete(data.watchId)
        }
        break
      }
      case 'plugin:suggest':
        (window as any).nuxtblogAdmin?.suggest(data.field, data.value)
        break
      case 'plugin:notify': {
        const toast = useToast()
        toast.add({ title: data.message, color: data.level || 'info' })
        break
      }
      case 'plugin:registerCommand': {
        // Bridge: register a handler in the host that forwards execution to the iframe
        ;(window as any).nuxtblogAdmin?.commands.register(data.commandId, async (ctx: any) => {
          // Strip function properties before sending to iframe
          const safe = Object.fromEntries(
            Object.entries(ctx).filter(([, v]) => typeof v !== 'function'),
          )
          iframe.contentWindow?.postMessage(
            { type: 'plugin:executeCommand', commandId: data.commandId, ctx: safe },
            '*',
          )
        })
        break
      }
    }
  })
}
