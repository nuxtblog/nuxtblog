export type MomentVisibility = 1 | 2 | 3  // 1=public 2=private 3=followers

export interface MomentAuthorItem {
  id: number
  username: string
  nickname: string
  avatar?: string
}

export interface MomentMediaItem {
  id: number
  url: string
  mime_type: string
  width?: number
  height?: number
}

export interface MomentStatsItem {
  view_count: number
  like_count: number
  comment_count: number
}

export interface MomentItem {
  id: number
  author_id: number
  content: string
  visibility: MomentVisibility
  created_at: string
  updated_at: string
  author?: MomentAuthorItem
  media?: MomentMediaItem[]
  stats?: MomentStatsItem
}

export interface CreateMomentRequest {
  content: string
  visibility?: MomentVisibility
  media_ids?: number[]
}

export interface UpdateMomentRequest {
  content?: string
  visibility?: MomentVisibility
  media_ids?: number[]
}

export interface MomentQueryRequest {
  page?: number
  page_size?: number
  author_id?: number
  visibility?: number
  keyword?: string
}

export interface PaginatedMoments {
  data: MomentItem[]
  total: number
  page: number
  page_size: number
  total_pages: number
}
