/** Matches backend AnnouncementItem */
export interface AnnouncementItem {
  id: number
  title: string
  content: string
  type: string   // info|warning|success|danger
  created_at: string
  unread: boolean
}

/** Matches backend AnnouncementListRes (public) */
export interface AnnouncementListResponse {
  list: AnnouncementItem[]
  total: number
  unread_count: number
}

/** Matches backend AnnouncementListAdminRes */
export interface AnnouncementAdminListResponse {
  list: AnnouncementItem[]
  total: number
}

/** Matches backend AnnouncementCreateReq */
export interface CreateAnnouncementRequest {
  title: string    // 1-200 chars
  content: string  // 1-5000 chars
  type?: 'info' | 'warning' | 'success' | 'danger'
}
