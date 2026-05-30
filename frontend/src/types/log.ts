export interface OperationLogInfo {
  id: number
  user_id: number
  username: string
  module: string
  action: string
  method: string
  path: string
  ip: string
  duration: number
  request_body: string
  status_code: number
  created_at: string
}

export interface LogFilter {
  username?: string
  method?: string
  path?: string
  start_time?: string
  end_time?: string
}
