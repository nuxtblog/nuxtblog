// composables/useNotificationApi.ts

export interface NotificationItem {
  id: number
  type: 'follow' | 'like' | 'comment' | 'reply' | 'mention' | 'system'
  sub_type?: string
  user_name?: string
  avatar?: string
  action?: string
  title?: string
  content: string
  related_title?: string
  related_link?: string
  read: boolean
  created_at: string
}

export interface NotificationListRes {
  list: NotificationItem[]
  total: number
  page: number
  size: number
  total_pages: number
  unread: number
}

export const useNotificationApi = () => {
  const { apiFetch } = useApiFetch()

  const getNotifications = async (params: {
    user_id: number
    filter?: string
    page?: number
    size?: number
  }) => {
    return await apiFetch<NotificationListRes>('/notifications', {
      method: 'GET',
      params,
    })
  }

  const getUnreadCount = async (userId: number) => {
    return await apiFetch<{ count: number }>('/notifications/unread-count', {
      method: 'GET',
      params: { user_id: userId },
    })
  }

  const markRead = async (id: number) => {
    return await apiFetch<void>(`/notifications/${id}/read`, {
      method: 'PUT',
    })
  }

  const markAllRead = async (userId: number) => {
    return await apiFetch<void>('/notifications/read-all', {
      method: 'PUT',
      params: { user_id: userId },
    })
  }

  const deleteNotification = async (id: number) => {
    return await apiFetch<void>(`/notifications/${id}`, {
      method: 'DELETE',
    })
  }

  const clearNotifications = async (userId: number, filter: string = 'all') => {
    return await apiFetch<void>('/notifications/clear', {
      method: 'DELETE',
      params: { user_id: userId, filter },
    })
  }

  return {
    getNotifications,
    getUnreadCount,
    markRead,
    markAllRead,
    deleteNotification,
    clearNotifications,
  }
}
