import request from './request'
import type { PageData } from '@/types/api'
import type {
  UserInfo,
  CreateUserParams,
  UpdateUserParams,
  ChangePasswordParams,
  BatchAssignRoleParams,
  UpdateProfileParams,
} from '@/types/user'

export function getUserListApi(params: {
  page: number
  page_size: number
  username?: string
  status?: number
  role_id?: number
}) {
  return request.get<any, { code: number; data: PageData<UserInfo>; message: string }>('/users', { params })
}

export function getUserApi(id: number) {
  return request.get<any, { code: number; data: UserInfo; message: string }>(`/users/${id}`)
}

export function createUserApi(data: CreateUserParams) {
  return request.post<any, { code: number; data: { id: number }; message: string }>('/users', data)
}

export function updateUserApi(id: number, data: UpdateUserParams) {
  return request.put<any, { code: number; data: null; message: string }>(`/users/${id}`, data)
}

export function deleteUserApi(id: number) {
  return request.delete<any, { code: number; data: null; message: string }>(`/users/${id}`)
}

export function updateUserStatusApi(id: number, status: number) {
  return request.put<any, { code: number; data: null; message: string }>(`/users/${id}/status`, { status })
}

export function assignUserRoleApi(id: number, role_id: number) {
  return request.put<any, { code: number; data: null; message: string }>(`/users/${id}/role`, { role_id })
}

export function batchAssignRoleApi(data: BatchAssignRoleParams) {
  return request.post<any, { code: number; data: null; message: string }>('/users/batch/role', data)
}

export function getUsersByRoleApi(roleId: number) {
  return request.get<any, { code: number; data: UserInfo[]; message: string }>(`/users/by-role/${roleId}`)
}

export function getUserProfileApi() {
  return request.get<any, { code: number; data: UserInfo; message: string }>('/users/profile')
}

export function updateProfileApi(data: UpdateProfileParams) {
  return request.put<any, { code: number; data: null; message: string }>('/users/profile', data)
}

export function changePasswordApi(data: ChangePasswordParams) {
  return request.put<any, { code: number; data: null; message: string }>('/users/password', data)
}

export function uploadAvatarApi(file: File) {
  const form = new FormData()
  form.append('avatar', file)
  return request.post<any, { code: number; data: { url: string }; message: string }>('/users/avatar', form)
}

export function importUsersApi(file: File) {
  const form = new FormData()
  form.append('file', file)
  return request.post<any, { code: number; data: { count: number }; message: string }>('/users/import', form)
}

export function exportUsersApi() {
  return request.get('/users/export', { responseType: 'blob' })
}
