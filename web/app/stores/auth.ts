import { defineStore } from 'pinia'
import type { AuthUser } from '~/composables/useAuthApi'

export const useAuthStore = defineStore('auth', () => {
  const config = useRuntimeConfig()
  const baseURL = config.public.apiBase as string

  // Persist token in cookie (SSR-compatible) and mirror to ref
  const tokenCookie = useCookie<string | null>('blog_token', {
    maxAge: 60 * 60 * 24 * 30, // 30 days
    httpOnly: false,
    sameSite: 'lax',
  })
  const refreshTokenCookie = useCookie<string | null>('blog_refresh_token', {
    maxAge: 60 * 60 * 24 * 30,
    httpOnly: false,
    sameSite: 'lax',
  })

  const token = computed(() => tokenCookie.value ?? null)
  const user = ref<AuthUser | null>(null)
  const isLoggedIn = computed(() => !!token.value && !!user.value)

  const localeCookie = useCookie<string>('i18n_locale', { default: () => 'zh' })

  const authApi = useAuthApi()

  const setTokens = (access: string, refresh: string) => {
    tokenCookie.value = access
    refreshTokenCookie.value = refresh
  }

  const setUser = (u: AuthUser) => {
    user.value = u
  }

  const login = async (loginStr: string, password: string) => {
    const res = await authApi.login(loginStr, password)
    setTokens(res.access_token, res.refresh_token)
    setUser(res.user)
    return res
  }

  const register = async (data: { username: string; email: string; password: string; display_name?: string; code?: string }) => {
    const res = await authApi.register(data)
    setTokens(res.access_token, res.refresh_token)
    setUser(res.user)
    return res
  }

  const logout = async () => {
    if (token.value) {
      await authApi.logout().catch(() => {})
    }
    tokenCookie.value = null
    refreshTokenCookie.value = null
    user.value = null
  }

  const fetchMe = async () => {
    if (!token.value) return
    try {
      const res = await authApi.me()
      user.value = res.user
    } catch {
      // access token may have expired — try refresh before giving up
      const refreshed = await tryRefresh()
      if (refreshed) {
        try {
          const res = await authApi.me()
          user.value = res.user
          return
        } catch { /* fall through to clear */ }
      }
      tokenCookie.value = null
      refreshTokenCookie.value = null
      user.value = null
    }
  }

  const syncLocale = async () => {
    if (!user.value || !token.value) return
    const cookieLocale = localeCookie.value
    if (cookieLocale && cookieLocale !== user.value.locale) {
      await $fetch<{ code: number; message: string; data: unknown }>(
        `${baseURL}/users/${user.value.id}`,
        {
          method: 'PUT',
          headers: { Authorization: `Bearer ${token.value}`, 'Accept-Language': cookieLocale },
          body: { locale: cookieLocale },
        }
      ).catch(() => {})
      user.value = { ...user.value, locale: cookieLocale }
    }
  }

  const tryRefresh = async () => {
    if (!refreshTokenCookie.value) return false
    try {
      const res = await authApi.refresh(refreshTokenCookie.value)
      tokenCookie.value = res.access_token
      return true
    } catch {
      tokenCookie.value = null
      refreshTokenCookie.value = null
      user.value = null
      return false
    }
  }

  return { token, user, isLoggedIn, login, register, logout, fetchMe, tryRefresh, setUser, setTokens, syncLocale }
})
