// Go int64 snowflake IDs (≈1.9×10¹⁸) exceed JS Number.MAX_SAFE_INTEGER (9×10¹⁵).
// JSON.parse silently rounds them, causing "not found" errors when IDs are sent
// back to the server in URL paths.
// Fix: before JSON.parse runs, replace bare integer values with 16+ digits with
// their quoted-string equivalents. 16 digits is safe: the largest normal field
// (e.g. file_size) never exceeds ~10¹³, while snowflake IDs are always ≥10¹⁸.
const parseSafeJSON = (text: string): unknown => {
  // Lookbehind: value must follow : , or [ (i.e. it's a JSON value, not inside a string).
  // Lookahead: value must be followed by , ] } or end.
  // This safely converts snowflake IDs (19 digits) without touching numbers inside strings.
  const fixed = text.replace(/(?<=[:,\[]\s*)(-?\d{16,})(?=\s*[,\]\}]|$)/g, '"$1"')
  return JSON.parse(fixed)
}

// Wrapper for GoFrame API responses: { code: 0, message: "", data: T }
export const useApiFetch = () => {
  const config = useRuntimeConfig()
  const baseURL = config.public.apiBase as string
  const authStore = useAuthStore()
  // Read from the same cookie @nuxtjs/i18n uses — avoids useI18n() composable restriction
  const localeCookie = useCookie<string>('i18n_locale', { default: () => 'zh' })

  let refreshPromise: Promise<boolean> | null = null

  const doFetch = async <T>(
    path: string,
    options?: Parameters<typeof $fetch>[1],
  ): Promise<{ code: number; message: string; data: T }> => {
    const url = path.startsWith('http') ? path : `${baseURL}${path}`
    const token = authStore.token
    const headers: Record<string, string> = {
      'Accept-Language': localeCookie.value,
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    }
    return $fetch<{ code: number; message: string; data: T }>(url, {
      ...options,
      headers: { ...headers, ...(options?.headers as Record<string, string> | undefined) },
      parseResponse: parseSafeJSON,
    })
  }

  const tryAutoRefresh = async (): Promise<boolean> => {
    if (!authStore.token) return false
    if (!refreshPromise) {
      refreshPromise = authStore.tryRefresh().finally(() => { refreshPromise = null })
    }
    return refreshPromise
  }

  const apiFetch = async <T>(
    path: string,
    options?: Parameters<typeof $fetch>[1],
  ): Promise<T> => {
    let response = await doFetch<T>(path, options)

    // Auto-refresh on 401: try once, then retry the original request
    if (response.code === 401 && authStore.token) {
      const ok = await tryAutoRefresh()
      if (ok) {
        response = await doFetch<T>(path, options)
      }
    }

    if (response.code !== 0) {
      throw new Error(response.message || 'API error')
    }
    return response.data
  }

  return { apiFetch }
}
