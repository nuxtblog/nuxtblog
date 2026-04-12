/**
 * Plugin Component Loader
 *
 * Loads plugin Vue components from pre-compiled admin.mjs files.
 * Plugins build their admin UI as ES modules with Vue externalized to
 * /_shared/vue.mjs. At load time we rewrite that import to use the
 * host app's Vue instance (window.__nuxtblog_vue) directly, ensuring
 * the plugin gets the exact same Vue runtime as the host Nuxt app.
 */
import { defineAsyncComponent, type Component } from 'vue'

const loadedComponents: Record<string, Component> = {}
const loadedModules: Record<string, Record<string, any>> = {}
// Version registry — populated by plugin loaders, consumed by component loaders
const pluginVersions: Record<string, string> = {}

/** Register a plugin's version so component loaders can use it for cache busting. */
export function registerPluginVersion(pluginId: string, version: string) {
  pluginVersions[pluginId] = version
}

/** Get a plugin's registered version. */
export function getPluginVersion(pluginId: string): string | undefined {
  return pluginVersions[pluginId]
}

/**
 * Load a plugin's admin.mjs, rewriting Vue imports to use the host's Vue.
 * Returns the full module namespace (cached).
 */
export async function loadPluginModule(pluginId: string, moduleFile = 'admin.mjs', version?: string): Promise<Record<string, any>> {
  const ver = version || pluginVersions[pluginId] || ''
  const cacheKey = `${pluginId}:${moduleFile}:${ver}`
  if (loadedModules[cacheKey]) return loadedModules[cacheKey]

  const qs = ver ? `?v=${ver}` : ''
  const url = `/api/plugins/${encodeURIComponent(pluginId)}/assets/${moduleFile}${qs}`
  const resp = await fetch(url)
  if (!resp.ok) throw new Error(`Failed to load plugin "${pluginId}": HTTP ${resp.status}`)

  let source = await resp.text()

  // Rewrite: import { ref, computed as r, ... } from "/_shared/vue.mjs"
  // To:      const { ref, computed: r, ... } = window.__nuxtblog_vue;
  // Note: import uses "as" for aliases, destructuring uses ":"
  source = source.replace(
    /import\s*\{([^}]+)\}\s*from\s*["'][^"']*vue\.mjs["']\s*;?/g,
    (_match, imports) => {
      const fixed = imports.replace(/([\w$]+)\s+as\s+([\w$]+)/g, '$1: $2')
      return `const {${fixed}} = window.__nuxtblog_vue;`
    },
  )
  // Also handle: import Vue from "/_shared/vue.mjs"
  source = source.replace(
    /import\s+(\w+)\s+from\s*["'][^"']*vue\.mjs["']\s*;?/g,
    (_match, name) => `const ${name} = window.__nuxtblog_vue;`,
  )

  // Rewrite: import { UButton, UCard as Card, ... } from "/_shared/ui.mjs"
  // To:      const { UButton, UCard: Card, ... } = window.__nuxtblog_ui;
  source = source.replace(
    /import\s*\{([^}]+)\}\s*from\s*["'][^"']*ui\.mjs["']\s*;?/g,
    (_match, imports) => {
      const fixed = imports.replace(/([\w$]+)\s+as\s+([\w$]+)/g, '$1: $2')
      return `const {${fixed}} = window.__nuxtblog_ui;`
    },
  )

  // Rewrite: import { AdminPageContainer, ... } from "/_shared/admin.mjs"
  // To:      const { AdminPageContainer, ... } = window.__nuxtblog_admin;
  source = source.replace(
    /import\s*\{([^}]+)\}\s*from\s*["'][^"']*admin\.mjs["']\s*;?/g,
    (_match, imports) => {
      const fixed = imports.replace(/([\w$]+)\s+as\s+([\w$]+)/g, '$1: $2')
      return `const {${fixed}} = window.__nuxtblog_admin;`
    },
  )

  // Inject authenticated fetch: override global fetch for plugin API calls
  // so plugins don't need to manually attach Authorization headers.
  // Token is stored in cookie "blog_token" by the auth store.
  const fetchShim = `
const __origFetch = globalThis.fetch;
const __pluginFetch = (url, opts = {}) => {
  if (typeof url === 'string' && (url.startsWith('/api/plugin/') || url.startsWith('/api/v1/'))) {
    const token = document.cookie.match(/(?:^|;)\\s*blog_token=([^;]*)/)?.[1];
    if (token) {
      opts.headers = Object.assign({}, opts.headers, { Authorization: 'Bearer ' + decodeURIComponent(token) });
    }
  }
  return __origFetch(url, opts);
};
`
  source = fetchShim + source.replace(/\bfetch\s*\(/g, '__pluginFetch(')

  const blob = new Blob([source], { type: 'application/javascript' })
  const blobUrl = URL.createObjectURL(blob)
  try {
    const mod = await import(/* @vite-ignore */ blobUrl)
    loadedModules[cacheKey] = mod
    return mod
  }
  finally {
    URL.revokeObjectURL(blobUrl)
  }
}

/**
 * Load a named export from a plugin's admin.mjs file as a Vue component.
 * Components are cached after first load.
 */
export async function loadPluginComponent(pluginId: string, componentName: string, moduleFile = 'admin.mjs'): Promise<Component> {
  const key = `${pluginId}:${moduleFile}:${componentName}`
  if (loadedComponents[key]) return loadedComponents[key]

  const mod = await loadPluginModule(pluginId, moduleFile)
  const component = mod[componentName] || mod.default
  if (!component) {
    throw new Error(`Component "${componentName}" not found in plugin "${pluginId}"`)
  }

  loadedComponents[key] = component
  return component
}

/**
 * Create a defineAsyncComponent wrapper for lazy-loading a plugin component.
 * Returns a component that can be used directly in templates.
 */
export function createPluginAsyncComponent(pluginId: string, componentName: string, moduleFile = 'admin.mjs'): Component {
  return defineAsyncComponent(() => loadPluginComponent(pluginId, componentName, moduleFile))
}

/**
 * Get all loaded plugin components.
 */
export function getLoadedPluginComponents(): Record<string, Component> {
  return { ...loadedComponents }
}
