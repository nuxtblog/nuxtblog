import type {
  MomentItem, CreateMomentRequest, UpdateMomentRequest, MomentQueryRequest, PaginatedMoments,
} from '~/types/api/moment'

export const useMomentApi = () => {
  const { apiFetch } = useApiFetch()

  const getMoments = (params?: MomentQueryRequest) =>
    apiFetch<PaginatedMoments>('/moments', { params })

  const getMoment = (id: number) =>
    apiFetch<MomentItem>(`/moments/${id}`)

  const createMoment = (data: CreateMomentRequest) =>
    apiFetch<{ id: number }>('/moments', { method: 'POST', body: data })

  const updateMoment = (id: number, data: UpdateMomentRequest) =>
    apiFetch<void>(`/moments/${id}`, { method: 'PUT', body: data })

  const deleteMoment = (id: number) =>
    apiFetch<void>(`/moments/${id}`, { method: 'DELETE' })

  const incrementView = (id: number) =>
    apiFetch<void>(`/moments/${id}/view`, { method: 'POST' })

  return { getMoments, getMoment, createMoment, updateMoment, deleteMoment, incrementView }
}
