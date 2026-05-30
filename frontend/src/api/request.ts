import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import axios from 'axios'
import type { AxiosInstance, AxiosError } from 'axios'
import { ElMessage } from 'element-plus'
import type { ApiResponse } from '@/types/api'

const request: AxiosInstance = axios.create({
  baseURL: '/api/v1',
  timeout: 15000,
})

request.interceptors.request.use((config) => {
  const token = useAuthStore().token
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

request.interceptors.response.use(
  (response) => {
    const res = response.data as ApiResponse
    if (res.code !== 0) {
      ElMessage.error(res.message || '请求失败')
      if (res.code === 401) {
        const authStore = useAuthStore()
        authStore.logout()
        const router = useRouter()
        router.push('/login')
      }
      return Promise.reject(new Error(res.message))
    }
    return response.data
  },
  (error: AxiosError<ApiResponse>) => {
    if (error.response?.status === 401) {
      const authStore = useAuthStore()
      authStore.logout()
      const router = useRouter()
      router.push('/login')
      return Promise.reject(error)
    }
    const message = error.response?.data?.message || error.message || '网络错误'
    ElMessage.error(message)
    return Promise.reject(error)
  },
)

export default request
