export interface ReportItem {
  id: number
  reporter_id: number
  reporter_name: string
  target_type: string
  target_id: number
  target_name: string
  reason: string
  detail: string
  status: string
  notes: string
  created_at: string
  resolved_at?: string
}

export const useReportApi = () => {
  const { apiFetch } = useApiFetch()

  const create = (targetType: string, targetId: number, reason: string, detail = '') =>
    apiFetch<{ id: number }>('/reports', {
      method: 'POST',
      body: { target_type: targetType, target_id: targetId, reason, detail },
    })

  const list = (status = 'pending', page = 1, size = 20) =>
    apiFetch<{ items: ReportItem[]; total: number }>(`/reports?status=${status}&page=${page}&size=${size}`)

  const handle = (id: number, status: 'resolved' | 'dismissed', notes = '') =>
    apiFetch(`/reports/${id}`, { method: 'PUT', body: { status, notes } })

  return { create, list, handle }
}
