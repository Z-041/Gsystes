import request from './request'
import type { PageData } from '@/types/api'
import type { OperationLogInfo, LogFilter } from '@/types/log'

export function getOperationLogListApi(params: LogFilter & { page: number; page_size: number }) {
  return request.get<any, { code: number; data: PageData<OperationLogInfo>; message: string }>('/logs', { params })
}
