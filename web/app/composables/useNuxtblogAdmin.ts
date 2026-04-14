/**
 * nuxtblogAdmin — global object exposed to plugin admin_js scripts.
 *
 * Provides: field watching, suggest/set, getPost, commands, views, http, notify.
 * This composable sets up the global `nuxtblogAdmin` on window and manages
 * command handlers registered by plugins.
 */

// ── Command context types (by call site) ────────────────────────────────
export interface PostListCommandContext {
  source: 'post-list'
  postId: number
  postTitle: string
}

export interface EditorCommandContext {
  source: 'editor'
  post: { title: string; slug: string; content: string; excerpt: string; status: string }
  selection: string | null
  replace: (text: string) => void
  insert: (text: string) => void
  setContent: (html: string) => void
  setSlug: (text: string) => void
  setExcerpt: (text: string) => void
  addTags: (tags: Array<string | { name: string; slug?: string }>) => Promise<void>
}

export interface GenericCommandContext {
  source: string
  [key: string]: unknown
}

export type CommandContext = PostListCommandContext | EditorCommandContext | GenericCommandContext


export interface Disposable {
  dispose: () => void
}

type FieldWatcher = (value: string) => void
type CommandHandler = (ctx: CommandContext) => void | Promise<void>

const fieldWatchers = new Map<string, Set<FieldWatcher>>()
const commandHandlers = new Map<string, CommandHandler>()

// Reactive field values — the post editor writes to these
const fieldValues = reactive<Record<string, string>>({
  'post.title': '',
  'post.slug': '',
  'post.content': '',
  'post.excerpt': '',
})

// Suggestion queue — plugins write to these, editor may accept
const fieldSuggestions = reactive<Record<string, string | null>>({})

// Force-set queue — plugins write, editor must accept
const fieldSets = reactive<Record<string, string | null>>({})

/** Called by the post editor to update field values for watchers. */
export function updatePluginField(field: string, value: string) {
  fieldValues[field] = value
  const watchers = fieldWatchers.get(field)
  if (watchers) {
    for (const cb of watchers) {
      try { cb(value) }
      catch (e) { console.warn(`[plugin] field watcher error for ${field}:`, e) }
    }
  }
}

/** Called by the post editor to consume a field suggestion. */
export function consumeSuggestion(field: string): string | null {
  const val = fieldSuggestions[field]
  if (val != null) {
    fieldSuggestions[field] = null
  }
  return val ?? null
}

/** Called by the post editor to consume a force-set value. */
export function consumeFieldSet(field: string): string | null {
  const val = fieldSets[field]
  if (val != null) {
    fieldSets[field] = null
  }
  return val ?? null
}

/** Dispatch a command by ID. Returns false if no handler is registered. */
export async function dispatchCommand(commandId: string, ctx: CommandContext): Promise<boolean> {
  const handler = commandHandlers.get(commandId)
  if (!handler) {
    console.warn(`[plugin] no handler for command: ${commandId}`)
    return false
  }
  try {
    await handler(ctx)
    return true
  }
  catch (e) {
    console.error(`[plugin] command ${commandId} error:`, e)
    return false
  }
}


/** Get all field suggestions (reactive). */
export function useFieldSuggestions() {
  return fieldSuggestions
}

/** Get all field force-sets (reactive). */
export function useFieldSets() {
  return fieldSets
}

/**
 * Install the nuxtblogAdmin global object on window.
 * Called once in the admin layout.
 */
export function installNuxtblogAdmin() {
  const toast = useToast()
  const { apiFetch } = useApiFetch()

  const nuxtblogAdmin = {
    // ── Field watching ─────────────────────────────────
    watch(field: string, cb: FieldWatcher): Disposable {
      if (!fieldWatchers.has(field)) {
        fieldWatchers.set(field, new Set())
      }
      fieldWatchers.get(field)!.add(cb)
      return {
        dispose: () => { fieldWatchers.get(field)?.delete(cb) },
      }
    },

    // ── Field writing ──────────────────────────────────
    suggest(field: string, value: string) {
      fieldSuggestions[field] = value
    },
    set(field: string, value: string) {
      fieldSets[field] = value
    },

    // ── Read draft ─────────────────────────────────────
    getPost() {
      return {
        title: fieldValues['post.title'] || '',
        slug: fieldValues['post.slug'] || '',
        content: fieldValues['post.content'] || '',
        excerpt: fieldValues['post.excerpt'] || '',
        status: fieldValues['post.status'] || 'draft',
      }
    },

    // ── Commands ───────────────────────────────────────
    commands: {
      register(id: string, handler: CommandHandler): Disposable {
        commandHandlers.set(id, handler)
        return {
          dispose: () => { commandHandlers.delete(id) },
        }
      },
      async execute(id: string, ctx: CommandContext) {
        await dispatchCommand(id, ctx)
      },
    },

    // ── Views ──────────────────────────────────────────
    views: {
      register(_id: string, _provider: (webview: unknown) => void): Disposable {
        // Webview support — stub for now
        return { dispose: () => {} }
      },
    },

    // ── HTTP (calls plugin's own backend routes) ───────
    http: {
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
      async put(path: string, body: object) {
        try {
          const data = await apiFetch<unknown>(path, { method: 'PUT', body })
          return { ok: true, data }
        }
        catch (e: unknown) {
          return { ok: false, data: null, error: String(e) }
        }
      },
      async delete(path: string, body?: object) {
        try {
          const data = await apiFetch<unknown>(path, { method: 'DELETE', body })
          return { ok: true, data }
        }
        catch (e: unknown) {
          return { ok: false, data: null, error: String(e) }
        }
      },
    },

    // ── Notifications ──────────────────────────────────
    notify: {
      success(msg: string) { toast.add({ title: msg, color: 'success' }) },
      error(msg: string) { toast.add({ title: msg, color: 'error' }) },
      info(msg: string) { toast.add({ title: msg, color: 'info' }) },
    },

    // ── Editor extensions ─────────────────────────────────
    editor: {
      registerExtension(extension: unknown): { dispose: () => void } {
        const unregister = registerPluginExtension(extension as any)
        return { dispose: unregister }
      },
    },

    // ── Theme tokens ────────────────────────────────────
    theme: {
      /** Get a snapshot of current theme tokens. */
      getTokens() {
        const t = (window as any).__nuxtblog_theme
        return t ? { ...t } : null
      },
      /** Reactive theme tokens object. Changes trigger Vue watchers. */
      get tokens() {
        return (window as any).__nuxtblog_theme
      },
    },
  }

  // Install on window
  ;(window as any).nuxtblogAdmin = nuxtblogAdmin

  return nuxtblogAdmin
}
