export interface AIConfig {
  id: string
  name: string
  api_format: string   // "openai" | "claude"
  label?: string       // optional display label, e.g. "OpenAI", "My Proxy"
  api_key: string
  model: string
  base_url: string
  timeout_ms: number
  is_active: boolean
}

export interface AIConfigsRes {
  items: AIConfig[]
  active_id: string
}

export function useAiApi() {
  const { apiFetch } = useApiFetch()

  const listConfigs = () =>
    apiFetch<AIConfigsRes>('/admin/ai/configs')

  const createConfig = (data: Partial<AIConfig>) =>
    apiFetch<{ item: AIConfig }>('/admin/ai/configs', { method: 'POST', body: data })

  const updateConfig = (id: string, data: Partial<AIConfig>) =>
    apiFetch<{ item: AIConfig }>(`/admin/ai/configs/${id}`, { method: 'PUT', body: data })

  const deleteConfig = (id: string) =>
    apiFetch(`/admin/ai/configs/${id}`, { method: 'DELETE' })

  const activateConfig = (id: string) =>
    apiFetch(`/admin/ai/configs/${id}/activate`, { method: 'POST' })

  const testConfig = (id: string) =>
    apiFetch<{ ok: boolean; message: string }>('/admin/ai/test', { method: 'POST', body: { id } })

  return { listConfigs, createConfig, updateConfig, deleteConfig, activateConfig, testConfig }
}
