export interface AnnouncementItem {
  id: number
  title: string
  content: string
  type: 'info' | 'warning' | 'success' | 'danger'
  created_at: string
  unread: boolean
}

export interface AnnouncementListRes {
  list: AnnouncementItem[]
  total: number
  unread_count: number
}

export const useAnnouncementApi = () => {
  const { apiFetch } = useApiFetch()

  const getAnnouncements = (params: { page?: number; size?: number }) =>
    apiFetch<AnnouncementListRes>('/announcements', { method: 'GET', params })

  const markAnnouncementsRead = () =>
    apiFetch<void>('/announcements/read', { method: 'PUT' })

  const adminListAnnouncements = (params: { page?: number; size?: number }) =>
    apiFetch<AnnouncementListRes>('/admin/announcements', { method: 'GET', params })

  const adminCreateAnnouncement = (body: { title: string; content: string; type: string }) =>
    apiFetch<{ id: number }>('/admin/announcements', { method: 'POST', body })

  const adminUpdateAnnouncement = (id: number, body: { title: string; content: string; type: string }) =>
    apiFetch<void>(`/admin/announcements/${id}`, { method: 'PUT', body })

  const adminDeleteAnnouncement = (id: number) =>
    apiFetch<void>(`/admin/announcements/${id}`, { method: 'DELETE' })

  return {
    getAnnouncements,
    markAnnouncementsRead,
    adminListAnnouncements,
    adminCreateAnnouncement,
    adminUpdateAnnouncement,
    adminDeleteAnnouncement,
  }
}
