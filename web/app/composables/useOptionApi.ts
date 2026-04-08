export const useOptionApi = () => {
  const { apiFetch } = useApiFetch()

  /** Get all autoloaded options (called on app init) */
  const getAutoload = () =>
    apiFetch<{ options: Record<string, string> }>('/options/autoload')

  /** Get a single option's raw JSON value */
  const getOption = async (key: string): Promise<unknown> => {
    const res = await apiFetch<{ key: string; value: string; autoload: number }>(`/options/${key}`)
    try {
      return JSON.parse(res.value)
    } catch {
      return res.value
    }
  }

  /** Set an option value (value is JSON-serialised automatically) */
  const setOption = (key: string, value: unknown) =>
    apiFetch<void>(`/options/${key}`, {
      method: 'PUT',
      body: { value: JSON.stringify(value) },
    })

  /** Delete an option */
  const deleteOption = (key: string) =>
    apiFetch<void>(`/options/${key}`, { method: 'DELETE' })

  // Legacy alias
  const get = getOption

  return { getAutoload, getOption, setOption, deleteOption, get }
}
