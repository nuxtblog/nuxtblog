export interface FollowStatus {
  following: boolean
  follower_count: number
  following_count: number
}

export interface FollowUserItem {
  id: number
  username: string
  display_name: string
  avatar: string
  bio: string
  followed_at: string
  article_count: number
  follower_count: number
  is_following_back: boolean
}

export interface FollowListResponse {
  list: FollowUserItem[]
  total: number
  page: number
  size: number
}

export const useFollowApi = () => {
  const { apiFetch } = useApiFetch()

  const follow = (userId: number) =>
    apiFetch<{ following: boolean }>(`/users/${userId}/follow`, { method: 'POST' })

  const unfollow = (userId: number) =>
    apiFetch<{ following: boolean }>(`/users/${userId}/follow`, { method: 'DELETE' })

  const removeFollower = (userId: number) =>
    apiFetch<{ removed: boolean }>(`/users/${userId}/follower`, { method: 'DELETE' })

  const getStatus = (userId: number) =>
    apiFetch<FollowStatus>(`/users/${userId}/follow-status`)

  const getFollowers = (userId: number, page = 1, size = 20) =>
    apiFetch<FollowListResponse>(`/users/${userId}/followers?page=${page}&size=${size}`)

  const getFollowing = (userId: number, page = 1, size = 20) =>
    apiFetch<FollowListResponse>(`/users/${userId}/following?page=${page}&size=${size}`)

  return { follow, unfollow, removeFollower, getStatus, getFollowers, getFollowing }
}
