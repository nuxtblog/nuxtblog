// Comment review status: 1=pending 2=approved 3=spam 4=trash
export type CommentStatus = 1 | 2 | 3 | 4

// ── Create ─────────────────────────────────────────────────────────────────────

/** Matches backend CommentCreateReq */
export interface CreateCommentRequest {
  object_id: number
  object_type: 'post' | 'product' | 'video'
  parent_id?: number
  author_name: string
  author_email: string
  content: string
}

export interface CreateCommentResponse {
  id: number
  status: CommentStatus
}

// ── Update ─────────────────────────────────────────────────────────────────────

export interface UpdateCommentStatusRequest {
  status: CommentStatus
}

export interface UpdateCommentContentRequest {
  content: string
}

// ── Public query ───────────────────────────────────────────────────────────────

/** Matches backend CommentGetListReq */
export interface CommentQueryRequest {
  object_id: number
  object_type: 'post' | 'product' | 'video'
  page?: number
  size?: number
}

// ── Public response item ───────────────────────────────────────────────────────

/** Matches backend CommentItem (public endpoint) */
export interface CommentItem {
  id: number
  object_id: number
  object_type: string
  parent_id?: number
  user_id?: number
  author_name: string
  author_email?: string
  content: string
  status: CommentStatus
  created_at?: string
  replies?: CommentItem[]
}

export interface CommentListResponse {
  list: CommentItem[]
  total: number
  page: number
  size: number
}

// ── Admin response item ────────────────────────────────────────────────────────

export interface CommentAdminAuthor {
  id?: number
  name: string
  email?: string
  avatar?: string
}

/** Matches backend CommentAdminItem (admin endpoint) */
export interface CommentAdminItem {
  comment_id: number
  post_id: number
  content: string
  parent_id?: number
  status: string
  ip_address?: string
  user_agent?: string
  created_at: string
  updated_at: string
  author?: CommentAdminAuthor
  children?: CommentAdminItem[]
}

// ── Admin query ────────────────────────────────────────────────────────────────

/** Matches backend CommentAdminGetListReq */
export interface CommentAdminQueryRequest {
  object_type?: 'post' | 'product' | 'video'
  status?: string
  keyword?: string
  page?: number
  size?: number
}

export interface CommentAdminListResponse {
  list: CommentAdminItem[]
  total: number
  page: number
  size: number
}

// ── Stats ──────────────────────────────────────────────────────────────────────

export interface CommentStatsResponse {
  total: number
  approved: number
  pending: number
  spam: number
  trash: number
}

// ── Batch ──────────────────────────────────────────────────────────────────────

export interface BatchCommentUpdateStatusRequest {
  ids: number[]
  status: CommentStatus
}
