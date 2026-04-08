export interface TokenItem {
  id: number
  name: string
  prefix: string        // "yblog_xxxxxxxx" — first 14 chars shown in list
  expires_at?: string
  last_used_at?: string
  created_at: string
}

export interface TokenCreateResult extends TokenItem {
  token: string         // full token, shown only once
}

export const useTokenApi = () => {
  const { apiFetch } = useApiFetch()

  const list = () =>
    apiFetch<TokenItem[]>('/users/tokens', { method: 'GET' })

  const create = (name: string, expiresInDays?: number) =>
    apiFetch<TokenCreateResult>('/users/tokens', {
      method: 'POST',
      body: { name, ...(expiresInDays ? { expires_in_days: expiresInDays } : {}) },
    })

  const revoke = (id: number) =>
    apiFetch<void>(`/users/tokens/${id}`, { method: 'DELETE' })

  return { list, create, revoke }
}
