export interface PluginCapabilities {
  http?: { allow?: string[]; timeout_ms?: number }
  store?: { read?: boolean; write?: boolean }
  events?: { subscribe?: string[] }
  db?: {
    own?: boolean
    tables?: Array<{ table: string; ops: string[] }>
    raw?: boolean
  } | boolean
}

export interface PluginItem {
  id: string
  title: string
  description: string
  version: string
  author: string
  icon: string
  repo_url: string
  enabled: boolean
  installed_at: string
  capabilities: string // raw JSON string
  /** 'builtin' = Go native (cannot uninstall), 'external' = installed via zip/github */
  source: 'builtin' | 'external'
  /** Plugin runtime type: builtin (Go native), js (Goja), yaml (declarative), full (JS + frontend) */
  type: 'builtin' | 'js' | 'yaml' | 'full'
  /** True when the plugin requires the Go server to be rebuilt and restarted to take effect (builtin plugins only). */
  need_restart?: boolean
}

/** Returns human-readable capability badge labels for a plugin. */
export function parseCapabilityBadges(capsJSON: string, t: (k: string) => string): Array<{ label: string; color: string }> {
  let caps: PluginCapabilities = {}
  try { caps = JSON.parse(capsJSON || '{}') } catch { /* ignore */ }

  const isFullAccess = !caps.http && !caps.store && !caps.events
  if (isFullAccess) {
    return [{ label: t('admin.plugins.cap_full_access'), color: 'warning' }]
  }

  const badges: Array<{ label: string; color: string }> = []
  if (caps.http) badges.push({ label: t('admin.plugins.cap_http'), color: 'info' })
  if (caps.store?.write) badges.push({ label: t('admin.plugins.cap_store_rw'), color: 'primary' })
  else if (caps.store?.read) badges.push({ label: t('admin.plugins.cap_store_r'), color: 'neutral' })
  if (caps.events) badges.push({ label: t('admin.plugins.cap_events'), color: 'success' })
  return badges
}

export interface MarketplaceItem {
  name: string
  title: string
  description: string
  version: string
  author: string
  icon: string
  repo: string
  homepage: string
  tags: string[]
  type: string
  /** 'compiled' = requires server restart, 'interpreted' = hot reload */
  runtime: string
  is_official: boolean
  license: string
  sdk_version: string
  trust_level: string
  capabilities: string[]
  features: string[]
  published_at: string
  updated_at: string
}

export interface PluginSettingField {
  key: string
  label: string
  /** "string" | "password" | "number" | "boolean" | "select" | "textarea" */
  type: string
  required?: boolean
  default?: unknown
  placeholder?: string
  description?: string
  options?: string[]
  /** Group name for visual grouping in the settings form */
  group?: string
  /** Condition expression, e.g. "advanced_mode === true". Field is hidden when falsy. */
  showIf?: string
  /** When true, other plugins can read this setting via ctx.plugins.getSetting() */
  shared?: boolean
}

export interface PluginWindowBucket {
  minute: string
  invocations: number
  errors: number
}

export interface PluginStatsRes {
  plugin_id: string
  invocations: number
  errors: number
  avg_duration_ms: number
  last_error?: string
  last_error_at?: string
  history: PluginWindowBucket[]
}

export interface PluginErrorEntry {
  at: string
  event: string
  message: string
  input_diff?: string
  /** "filter" | "handler" | "pipeline" | "route" */
  phase?: string
  /** JS stack trace when available */
  stack?: string
  /** Execution time in milliseconds */
  duration_ms?: number
}

// Pipeline types (Phase 4)
export interface StepDef {
  type: 'js' | 'webhook' | 'condition'
  name: string
  fn?: string
  url?: string
  headers?: Record<string, string>
  if?: string
  then?: StepDef[]
  else?: StepDef[]
  timeout_ms?: number
  retry?: number
}

export interface PipelineDef {
  name: string
  trigger: string
  steps: StepDef[]
}

export interface PluginDependency {
  id: string
  version?: string
  optional?: boolean
}

export interface PluginManifest {
  depends?: PluginDependency[]
  pipelines?: PipelineDef[]
  webhooks?: Array<{ url: string; events: string[]; headers?: Record<string, string> }>
  [key: string]: unknown
}

export interface PluginPreviewInfo {
  name: string
  title: string
  description: string
  version: string
  author: string
  icon: string
  priority: number
  has_css: boolean
  capabilities: {
    http?: { allow?: string[]; timeout_ms?: number }
    store?: { read?: boolean; write?: boolean }
    events?: { subscribe?: string[] }
    db?: {
      own?: boolean
      tables?: Array<{ table: string; ops: string[] }>
      raw?: boolean
    } | boolean
  }
  depends?: PluginDependency[]
  settings: PluginSettingField[]
  webhooks: Array<{ url: string; events: string[] }>
  pipelines: Array<{ name: string; trigger: string; step_count: number }>
  permissions?: string[]
}

export const usePluginApi = () => {
  const { apiFetch } = useApiFetch()

  const list = () =>
    apiFetch<{ items: PluginItem[] }>('/admin/plugins')

  const install = (repoUrl: string, expectedVersion?: string) =>
    apiFetch<{ item: PluginItem }>('/admin/plugins', {
      method: 'POST',
      body: { repo_url: repoUrl, expected_version: expectedVersion },
    })

  const uploadZip = (file: File) => {
    const form = new FormData()
    form.append('file', file)
    return apiFetch<{ item: PluginItem }>('/admin/plugins/upload', {
      method: 'POST',
      body: form,
    })
  }

  const uninstall = (id: string) =>
    apiFetch<{ need_restart?: boolean }>(`/admin/plugins/${encodeURIComponent(id)}`, { method: 'DELETE' })

  const unloadImpact = (id: string) =>
    apiFetch<{ will_unload: string[] }>(`/admin/plugins/${encodeURIComponent(id)}/unload-impact`)

  const batchUninstall = (ids: string[]) =>
    apiFetch<{ succeeded: string[]; failed: string[]; need_restart?: boolean }>('/admin/plugins/batch-uninstall', {
      method: 'POST',
      body: { ids },
    })

  const toggle = (id: string, enabled: boolean) =>
    apiFetch(`/admin/plugins/${encodeURIComponent(id)}`, {
      method: 'PATCH',
      body: { enabled },
    })

  const getSettings = (id: string) =>
    apiFetch<{ schema: PluginSettingField[]; values: Record<string, unknown> }>(
      `/admin/plugins/${encodeURIComponent(id)}/settings`,
    )

  const updateSettings = (id: string, values: Record<string, unknown>) =>
    apiFetch(`/admin/plugins/${encodeURIComponent(id)}/settings`, {
      method: 'PUT',
      body: { values },
    })

  const update = (id: string) =>
    apiFetch<{ item: PluginItem }>(`/admin/plugins/${encodeURIComponent(id)}/update`, { method: 'POST' })

  const batchUpdate = (ids: string[]) =>
    apiFetch<{ succeeded: string[]; failed: string[] }>('/admin/plugins/batch-update', {
      method: 'POST',
      body: { ids },
    })

  const getStyles = () =>
    apiFetch<{ css: string }>('/plugins/styles')

  const marketplace = (keyword?: string, type?: string) =>
    apiFetch<{ items: MarketplaceItem[]; synced_at: string }>('/admin/plugins/marketplace', {
      params: { keyword, type },
    })

  const syncMarketplace = () =>
    apiFetch<{ count: number; synced_at: string }>('/admin/plugins/marketplace/sync', { method: 'POST' })

  const getStats = (id: string) =>
    apiFetch<PluginStatsRes>(`/admin/plugins/${encodeURIComponent(id)}/stats`)

  const getErrors = (id: string) =>
    apiFetch<{ items: PluginErrorEntry[] }>(`/admin/plugins/${encodeURIComponent(id)}/errors`)

  const getManifest = (id: string) =>
    apiFetch<{ manifest: string }>(`/admin/plugins/${encodeURIComponent(id)}/manifest`)

  const updateManifest = (id: string, manifest: string) =>
    apiFetch(`/admin/plugins/${encodeURIComponent(id)}/manifest`, {
      method: 'PUT',
      body: { manifest },
    })

  const preview = (repo: string) =>
    apiFetch<PluginPreviewInfo>('/admin/plugins/preview', { params: { repo } })

  /** Phase 2.4: get enabled plugins with client-side info (contributes, admin_js) */
  const clientList = () =>
    apiFetch<{ items: PluginClientItem[] }>('/admin/plugins/client')

  return { list, install, uploadZip, uninstall, unloadImpact, batchUninstall, update, batchUpdate, toggle, getSettings, updateSettings, getStyles, marketplace, syncMarketplace, getStats, getErrors, getManifest, updateManifest, preview, clientList }
}

/** Phase 2.4: Client-side plugin info for admin panel */
export interface PluginClientItem {
  id: string
  title: string
  icon: string
  trust_level: string
  admin_js?: string
  public_js?: string
  contributes?: string
}
