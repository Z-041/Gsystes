import request from './request'
import type { PageData } from '@/types/api'
import type { RoleInfo, RoleSimple, CreateRoleParams, UpdateRoleParams, AssignPermissionsParams } from '@/types/role'
import type { PermissionInfo } from '@/types/permission'

export function getRoleListApi(params: { page: number; page_size: number }) {
  return request.get<any, { code: number; data: PageData<RoleInfo>; message: string }>('/roles', { params })
}

export function getAllRolesApi() {
  return request.get<any, { code: number; data: RoleSimple[]; message: string }>('/roles/all')
}

export function getRoleApi(id: number) {
  return request.get<any, { code: number; data: RoleInfo; message: string }>(`/roles/${id}`)
}

export function createRoleApi(data: CreateRoleParams) {
  return request.post<any, { code: number; data: { id: number }; message: string }>('/roles', data)
}

export function updateRoleApi(id: number, data: UpdateRoleParams) {
  return request.put<any, { code: number; data: null; message: string }>(`/roles/${id}`, data)
}

export function deleteRoleApi(id: number) {
  return request.delete<any, { code: number; data: null; message: string }>(`/roles/${id}`)
}

export function assignPermissionsApi(id: number, data: AssignPermissionsParams) {
  return request.post<any, { code: number; data: null; message: string }>(`/roles/${id}/permissions`, data)
}

export function getRolePermissionsApi(id: number) {
  return request.get<any, { code: number; data: PermissionInfo[]; message: string }>(`/roles/${id}/permissions`)
}
