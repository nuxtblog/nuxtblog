export const useHistoryApi = () => {
  const { apiFetch } = useApiFetch()

  const list = (page = 1, size = 20) =>
    apiFetch<{ items: HistoryItem[]; total: number }>(`/history?page=${page}&size=${size}`)

  const clear = () =>
    apiFetch('/history', { method: 'DELETE' })

  return { list, clear }
}

export interface HistoryItem {
  post_id: number
  post_title: string
  post_slug: string
  post_cover: string
  viewed_at: string
}
