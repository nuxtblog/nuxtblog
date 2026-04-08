export const useReactionApi = () => {
  const { apiFetch } = useApiFetch()
  const authStore = useAuthStore()

  const headers = computed(() =>
    authStore.token ? { Authorization: `Bearer ${authStore.token}` } : {}
  )

  const likePost = (id: number) =>
    apiFetch(`/posts/${id}/like`, { method: 'POST', headers: headers.value })

  const unlikePost = (id: number) =>
    apiFetch(`/posts/${id}/like`, { method: 'DELETE', headers: headers.value })

  const bookmarkPost = (id: number) =>
    apiFetch(`/posts/${id}/bookmark`, { method: 'POST', headers: headers.value })

  const unbookmarkPost = (id: number) =>
    apiFetch(`/posts/${id}/bookmark`, { method: 'DELETE', headers: headers.value })

  const getReaction = (id: number) =>
    apiFetch<{ liked: boolean; bookmarked: boolean }>(`/posts/${id}/reaction`, {
      headers: headers.value,
    })

  const getBookmarks = (page = 1, size = 20) =>
    apiFetch<{
      list: Array<{ id: number; title: string; slug: string; excerpt: string; created_at: string }>
      total: number
      page: number
      size: number
    }>(`/users/me/bookmarks?page=${page}&size=${size}`, { headers: headers.value })

  return { likePost, unlikePost, bookmarkPost, unbookmarkPost, getReaction, getBookmarks }
}
