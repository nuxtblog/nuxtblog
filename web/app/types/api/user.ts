// Backend enum types
export type UserRole   = 1 | 2 | 3 | 4  // 1=subscriber 2=editor 3=admin 4=super_admin
export type UserStatus = 1 | 2 | 3       // 1=active 2=banned 3=pending

// ── Auth ───────────────────────────────────────────────────────────────────────

export interface RegisterRequest {
  username: string
  password: string
  email: string
  nickname?: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user?: UserResponse
}

// ── Create / Update ────────────────────────────────────────────────────────────

/** Matches backend UserCreateReq */
export interface CreateUserRequest {
  username: string      // 3-30 chars
  password: string      // 8-64 chars
  email: string
  display_name: string  // 1-100 chars
  role: UserRole
  locale?: string
}

/** Matches backend UserUpdateReq */
export interface UpdateUserRequest {
  display_name?: string
  bio?: string
  avatar_id?: number
  locale?: string
  status?: UserStatus   // 1=active 2=banned 3=pending
  role?: UserRole       // 1=subscriber 2=editor 3=admin 4=super_admin
  location?: string
  website?: string
  github?: string
  twitter?: string
  instagram?: string
  linkedin?: string
  youtube?: string
  cover?: string
  cover_id?: number
}

export interface UpdatePasswordRequest {
  old_password: string
  new_password: string
}

export interface ResetPasswordRequest {
  new_password: string
}

// ── Query ──────────────────────────────────────────────────────────────────────

export interface UserQueryRequest {
  page?: number
  size?: number
  role?: UserRole
  status?: UserStatus
  keyword?: string
  order_by?: 'created_at' | 'id' | 'username'
  order?: 'asc' | 'desc'
}

// ── Responses ──────────────────────────────────────────────────────────────────

/** Matches backend UserItem JSON output */
export interface UserResponse {
  id: number
  username: string
  email: string
  display_name: string
  avatar_id?: number
  avatar?: string
  cover?: string
  bio: string
  role: UserRole
  status: UserStatus
  email_verified: number
  locale: string
  metas?: Record<string, string>
  last_login_at?: string
  created_at: string
  updated_at: string
}

/** Backend returns full UserItem in list too */
export type UserListResponse = UserResponse

/** List data wrapper */
export interface UserListData {
  list: UserListResponse[]
  total: number
  page: number
  size: number
}

export interface UserStatsResponse {
  total: number
  active: number
  pending: number
  banned: number
}

export interface BatchUpdateStatusRequest {
  ids: number[]
  status: UserStatus
}
