import request from './request'
import type { LoginParams, LoginResult } from '@/types/user'
import type { PermissionInfo } from '@/types/permission'

export function loginApi(data: LoginParams) {
  return request.post<any, { code: number; data: LoginResult; message: string }>('/auth/login', data)
}

export function getCurrentMenusApi() {
  return request.get<any, { code: number; data: PermissionInfo[]; message: string }>('/auth/menus')
}

export function getCurrentPermissionsApi() {
  return request.get<any, { code: number; data: string[]; message: string }>('/auth/permissions')
}
