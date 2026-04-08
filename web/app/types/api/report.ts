/** Matches backend ReportItem */
export interface ReportItem {
  id: number
  reporter_id: number
  reporter_name: string
  target_type: string   // post|comment|user
  target_id: number
  target_name: string
  reason: string
  detail: string
  status: string        // pending|resolved|dismissed
  notes: string
  created_at: string
  resolved_at?: string
}

/** Matches backend ReportCreateReq */
export interface CreateReportRequest {
  target_type: 'post' | 'comment' | 'user'
  target_id: number
  reason: string
  detail?: string
}

/** Matches backend ReportListReq */
export interface ReportListRequest {
  status?: 'pending' | 'resolved' | 'dismissed' | 'all'
  page?: number
  size?: number
}

/** Matches backend ReportListRes */
export interface ReportListResponse {
  items: ReportItem[]
  total: number
}

/** Matches backend ReportHandleReq */
export interface ReportHandleRequest {
  status: 'resolved' | 'dismissed'
  notes?: string
}
