export interface UserInfo {
  id: number
  username: string
  nickname: string
  email: string
  phone: string
  avatar: string
  status: number
  role_id: number
  role?: { id: number; name: string; code: string }
  created_at: string
  updated_at: string
}

export interface CreateUserParams {
  username: string
  password: string
  nickname?: string
  email?: string
  phone?: string
  role_id: number
}

export interface UpdateUserParams {
  nickname?: string
  email?: string
  phone?: string
  role_id?: number
  status?: number
}

export interface LoginParams {
  username: string
  password: string
}

export interface LoginResult {
  token: string
  user: {
    id: number
    username: string
    nickname: string
    avatar: string
    role_id: number
  }
}

export interface ChangePasswordParams {
  old_password: string
  new_password: string
}

export interface BatchAssignRoleParams {
  user_ids: number[]
  role_id: number
}

export interface UpdateProfileParams {
  nickname?: string
  email?: string
  phone?: string
}
