/** Matches backend NotificationItem */
export interface NotificationItem {
  id: number
  type: string        // follow|like|comment|reply|mention|system
  sub_type?: string   // system sub type
  user_name?: string  // actor name
  avatar?: string     // actor avatar URL
  action?: string     // e.g. "评论了你的文章"
  title?: string      // system notification title
  content: string
  related_title?: string
  related_link?: string
  read: boolean
  created_at: string
}

export interface NotificationListRequest {
  user_id: number
  filter?: 'all' | 'unread' | 'interaction' | 'system'
  page?: number
  size?: number
}

export interface NotificationListResponse {
  list: NotificationItem[]
  total: number
  page: number
  size: number
  total_pages: number
  unread: number
}

export interface NotificationUnreadCountResponse {
  count: number
}
