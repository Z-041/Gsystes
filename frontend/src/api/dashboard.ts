import request from './request'

export interface DashboardStatsResponse {
  user_count: number
  role_count: number
  today_log_count: number
}

export function getDashboardStatsApi() {
  return request.get<any, { code: number; data: DashboardStatsResponse; message: string }>('/dashboard/stats')
}
