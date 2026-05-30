export interface PermissionInfo {
  id: number
  name: string
  code: string
  type: number
  parent_id: number
  path: string
  method: string
  sort: number
  children?: PermissionInfo[]
  created_at?: string
  updated_at?: string
}

export interface CreatePermissionParams {
  name: string
  code: string
  type: number
  parent_id?: number
  path?: string
  method?: string
  sort?: number
}

export interface UpdatePermissionParams {
  name: string
  code: string
  type: number
  parent_id?: number
  path?: string
  method?: string
  sort?: number
}
