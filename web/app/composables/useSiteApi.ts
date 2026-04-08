export const useSiteApi = () => {
  const { apiFetch } = useApiFetch()

  const getLanguages = () =>
    apiFetch<{ list: Array<{ code: string; name: string; label: string }> }>('/languages')

  return { getLanguages }
}
