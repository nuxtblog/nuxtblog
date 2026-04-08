import type {
  CreateUserRequest,
  UpdateUserRequest,
  UserQueryRequest,
  UserResponse,
  UserListData,
  UserStatsResponse,
  ResetPasswordRequest,
} from '~/types/api/user'

export const useUserApi = () => {
  const { apiFetch } = useApiFetch()

  const getUsers = async (params?: UserQueryRequest) =>
    apiFetch<UserListData>('/users', { method: 'GET', params })

  const getUser = async (id: number) =>
    apiFetch<UserResponse>(`/users/${id}`, { method: 'GET' })

  const createUser = async (data: CreateUserRequest) =>
    apiFetch<{ id: number }>('/users', { method: 'POST', body: data })

  const updateUser = async (id: number, data: UpdateUserRequest) =>
    apiFetch<void>(`/users/${id}`, { method: 'PUT', body: data })

  const deleteUser = async (id: number) =>
    apiFetch<void>(`/users/${id}`, { method: 'DELETE' })

  const resetUserPassword = async (id: number, data: ResetPasswordRequest) =>
    apiFetch<void>(`/users/${id}/password`, { method: 'PUT', body: data })

  const getUserStats = async () =>
    apiFetch<UserStatsResponse>('/users/stats', { method: 'GET' })

  return { getUsers, getUser, createUser, updateUser, deleteUser, resetUserPassword, getUserStats }
}
