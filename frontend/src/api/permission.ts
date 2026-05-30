import request from './request'
import type { PageData } from '@/types/api'
import type { PermissionInfo, CreatePermissionParams, UpdatePermissionParams } from '@/types/permission'

export function getPermissionListApi(params: { page: number; page_size: number }) {
  return request.get<any, { code: number; data: PageData<PermissionInfo>; message: string }>('/permissions', { params })
}

export function getAllPermissionsApi() {
  return request.get<any, { code: number; data: PermissionInfo[]; message: string }>('/permissions/all')
}

export function getPermissionApi(id: number) {
  return request.get<any, { code: number; data: PermissionInfo; message: string }>(`/permissions/${id}`)
}

export function createPermissionApi(data: CreatePermissionParams) {
  return request.post<any, { code: number; data: { id: number }; message: string }>('/permissions', data)
}

export function updatePermissionApi(id: number, data: UpdatePermissionParams) {
  return request.put<any, { code: number; data: null; message: string }>(`/permissions/${id}`, data)
}

export function deletePermissionApi(id: number) {
  return request.delete<any, { code: number; data: null; message: string }>(`/permissions/${id}`)
}
