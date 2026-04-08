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

interface PluginClientItem {
  id: string
  title: string
  icon: string
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

      // Register contribution points
      for (const plugin of plugins.value) {
        if (plugin.contributes) {
          try {
            const contributes: PluginContributes = JSON.parse(plugin.contributes)
            contributionsStore.registerPlugin(plugin.id, contributes)
          }
          catch (e) {
            console.warn(`[plugin-loader] failed to parse contributes for ${plugin.id}:`, e)
          }
        }

        // Phase 4.2: register page nav items from manifest pages[]
        if (plugin.pages) {
          try {
            const pages: PageDef[] = JSON.parse(plugin.pages)
            const pageNavItems = pages
              .filter(p => p.slot === 'admin' && p.nav)
              .map(p => ({
                slot: 'admin:sidebar-nav',
                title: p.title || p.component,
                icon: p.nav?.icon,
                route: `/admin/plugin-page/${encodeURIComponent(plugin.id)}/${encodeURIComponent(p.component)}`,
                order: p.nav?.order ?? 100,
              }))
            if (pageNavItems.length > 0) {
              contributionsStore.registerPlugin(plugin.id, { navigation: pageNavItems })
            }
          }
          catch (e) {
            console.warn(`[plugin-loader] failed to parse pages for ${plugin.id}:`, e)
          }
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
 * - official/local: loaded as inline <script> in main page context
 * - community: loaded in a sandboxed iframe with postMessage bridge
 */
async function loadAdminScript(plugin: PluginClientItem) {
  // admin_js may include a directory prefix (e.g. "dist/admin.mjs"), but the
  // backend saves assets as flat filenames. Strip any directory prefix.
  const assetFilename = plugin.admin_js!.split('/').pop()!
  const scriptUrl = `/api/plugins/${encodeURIComponent(plugin.id)}/assets/${assetFilename}`

  if (plugin.trust_level === 'official' || plugin.trust_level === 'local') {
    // Main context — full nuxtblogAdmin access
    try {
      const script = document.createElement('script')
      script.type = 'module'
      script.src = scriptUrl
      script.dataset.pluginId = plugin.id
      document.head.appendChild(script)

      // Wait for script to load, then call activate() if exported
      await new Promise<void>((resolve, reject) => {
        script.onload = () => {
          // activate() is called by the script itself using nuxtblogAdmin
          resolve()
        }
        script.onerror = () => reject(new Error(`Failed to load admin_js for ${plugin.id}`))
      })
    }
    catch (e) {
      console.warn(`[plugin-loader] failed to load admin_js for ${plugin.id}:`, e)
    }
  }
  else {
    // Community plugins: sandboxed iframe
    // The iframe only has access to a restricted postMessage-based API
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
    // Listen for field updates from parent
    window.addEventListener('message', function(e) {
      if (e.data && e.data.type === 'fieldUpdate' && window.__watchers) {
        var cb = window.__watchers[e.data.watchId];
        if (cb) cb(e.data.value);
      }
    });
    <\/script>
    <script type="module" src="${scriptUrl}"><\/script>
    </head><body></body></html>
  `

  iframe.srcdoc = html
  document.body.appendChild(iframe)

  // Listen for postMessage from the sandboxed iframe
  window.addEventListener('message', (e) => {
    if (e.source !== iframe.contentWindow) return
    const data = e.data
    if (!data || !data.type) return

    switch (data.type) {
      case 'plugin:suggest':
        (window as any).nuxtblogAdmin?.suggest(data.field, data.value)
        break
      case 'plugin:notify': {
        const toast = useToast()
        toast.add({ title: data.message, color: data.level || 'info' })
        break
      }
    }
  })
}
