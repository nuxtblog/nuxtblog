export interface ConversationItem {
  id: number
  other_user_id: number
  other_name: string
  other_avatar: string
  last_msg: string
  last_msg_at: string
  unread_count: number
}

export interface MessageItem {
  id: number
  sender_id: number
  content: string
  is_read: boolean
  created_at: string
}

export const useMessageApi = () => {
  const { apiFetch } = useApiFetch()

  const listConversations = (page = 1) =>
    apiFetch<{ items: ConversationItem[]; total: number; total_unread: number }>(`/messages?page=${page}`)

  const send = (toUserId: number, content: string) =>
    apiFetch<{ id: number; conversation_id: number }>(`/messages/${toUserId}`, {
      method: 'POST',
      body: { content },
    })

  const listMessages = (toUserId: number, beforeId = 0, size = 30) =>
    apiFetch<{ items: MessageItem[]; conversation_id: number; has_more: boolean }>(
      `/messages/${toUserId}/history?before_id=${beforeId}&size=${size}`
    )

  const unreadCount = () =>
    apiFetch<{ count: number }>('/messages/unread')

  return { listConversations, send, listMessages, unreadCount }
}
