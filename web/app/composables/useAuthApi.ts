export interface AuthUser {
  id: number
  username: string
  email: string
  display_name: string
  bio: string
  role: number
  status: number
  avatar_id?: number
  avatar?: string
  locale: string
  created_at?: string
  has_password: boolean
}

export interface AuthTokens {
  access_token: string
  refresh_token: string
  expires_in: number
  user: AuthUser
}

export const useAuthApi = () => {
  const { apiFetch } = useApiFetch()

  const login = (login: string, password: string) =>
    apiFetch<AuthTokens>('/auth/login', { method: 'POST', body: { login, password } })

  const register = (data: { username: string; email: string; password: string; display_name?: string; code?: string }) =>
    apiFetch<AuthTokens>('/auth/register', { method: 'POST', body: data })

  const me = () =>
    apiFetch<{ user: AuthUser }>('/auth/me')

  const logout = () =>
    apiFetch('/auth/logout', { method: 'POST' })

  const refresh = (refreshToken: string) =>
    apiFetch<{ access_token: string; expires_in: number }>('/auth/refresh', {
      method: 'POST',
      body: { refresh_token: refreshToken },
    })

  const getOAuthProviders = () =>
    apiFetch<{ providers: string[] }>('/auth/oauth/providers').then(r => r.providers)

  return { login, register, me, logout, refresh, getOAuthProviders }
}
