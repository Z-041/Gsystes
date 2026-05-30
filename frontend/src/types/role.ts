export interface RoleInfo {
  id: number
  name: string
  code: string
  description: string
  status: number
  permissions?: PermissionSimple[]
  created_at: string
  updated_at: string
}

export interface RoleSimple {
  id: number
  name: string
  code: string
}

export interface CreateRoleParams {
  name: string
  code: string
  description?: string
}

export interface UpdateRoleParams {
  name: string
  code: string
  description?: string
  status?: number
}

export interface AssignPermissionsParams {
  permission_ids: number[]
}

export interface PermissionSimple {
  id: number
  name: string
  code: string
  type: number
  parent_id: number
  path: string
  method: string
  sort: number
}
