export const useCheckinApi = () => {
  const { apiFetch } = useApiFetch()

  const doCheckin = () =>
    apiFetch<{ already_checked_in: boolean; streak: number }>('/checkin', { method: 'POST' })

  const getStatus = () =>
    apiFetch<{ checked_in_today: boolean; streak: number }>('/checkin/status')

  return { doCheckin, getStatus }
}
