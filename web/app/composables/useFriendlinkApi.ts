export interface FriendlinkItem {
  id: number
  name: string
  url: string
  logo: string
  description: string
  sort_order: number
  status: number
  created_at: string
}

export interface FriendlinkListRes {
  list: FriendlinkItem[]
  total: number
}

export const useFriendlinkApi = () => {
  const { apiFetch } = useApiFetch()

  const getFriendlinks = () =>
    apiFetch<{ list: FriendlinkItem[] }>('/friendlinks', { method: 'GET' })

  const adminListFriendlinks = (params: { page?: number; size?: number }) =>
    apiFetch<FriendlinkListRes>('/admin/friendlinks', { method: 'GET', params })

  const adminCreateFriendlink = (body: { name: string; url: string; logo?: string; description?: string; sort_order?: number; status?: number }) =>
    apiFetch<{ id: number }>('/admin/friendlinks', { method: 'POST', body })

  const adminUpdateFriendlink = (id: number, body: { name: string; url: string; logo?: string; description?: string; sort_order?: number; status?: number }) =>
    apiFetch<void>(`/admin/friendlinks/${id}`, { method: 'PUT', body })

  const adminDeleteFriendlink = (id: number) =>
    apiFetch<void>(`/admin/friendlinks/${id}`, { method: 'DELETE' })

  return {
    getFriendlinks,
    adminListFriendlinks,
    adminCreateFriendlink,
    adminUpdateFriendlink,
    adminDeleteFriendlink,
  }
}
