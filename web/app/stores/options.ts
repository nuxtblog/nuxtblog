/**
 * Stores all autoload options (loaded once at app startup).
 * Values in the DB are JSON-encoded strings, e.g. `'"foo"'` or `'true'`.
 * Use `get(key, fallback)` to retrieve a parsed value.
 */
export const useOptionsStore = defineStore('options', () => {
  const _raw = useState<Record<string, string>>('options:raw', () => ({}))
  const loaded = useState<boolean>('options:loaded', () => false)

  const load = async () => {
    if (loaded.value) return
    const optionApi = useOptionApi()
    try {
      const data = await optionApi.getAutoload()
      _raw.value = data.options ?? {}
      loaded.value = true
    } catch (e) {
      console.error('[options] failed to load autoload options:', e)
    }
  }

  // Force reload from API (e.g. after admin saves settings)
  const reload = async () => {
    const optionApi = useOptionApi()
    try {
      const data = await optionApi.getAutoload()
      _raw.value = data.options ?? {}
    } catch (e) {
      console.error('[options] failed to reload autoload options:', e)
    }
  }

  /**
   * Returns the parsed value for a key.
   * Falls back to `defaultValue` when the key is missing or JSON.parse fails.
   */
  const get = (key: string, defaultValue = ''): string => {
    const raw = _raw.value[key]
    if (raw === undefined) return defaultValue
    try {
      const parsed = JSON.parse(raw)
      return typeof parsed === 'string' ? parsed : String(parsed)
    } catch {
      return raw
    }
  }

  /**
   * Returns the raw JSON-parsed value for a key.
   * Falls back to `fallback` when the key is missing or JSON.parse fails.
   */
  const getJSON = <T>(key: string, fallback: T): T => {
    const raw = _raw.value[key]
    if (raw === undefined) return fallback
    try {
      return JSON.parse(raw) as T
    } catch {
      return fallback
    }
  }

  return { loaded, load, reload, get, getJSON }
})
