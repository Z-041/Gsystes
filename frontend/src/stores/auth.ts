import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { UserInfo, LoginParams } from '@/types/user'
import type { PermissionInfo } from '@/types/permission'
import { loginApi, getCurrentPermissionsApi, getCurrentMenusApi } from '@/api/auth'
import { getUserProfileApi } from '@/api/user'
import { useAppStore } from '@/stores/app'

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const userInfo = ref<UserInfo | null>(null)
  const permissions = ref<string[]>([])
  const menuTree = ref<PermissionInfo[]>([])

  const isLoggedIn = computed(() => !!token.value)

  let wsConnectFn: (() => void) | null = null
  let wsDisconnectFn: (() => void) | null = null

  function registerWs(connect: () => void, disconnect: () => void) {
    wsConnectFn = connect
    wsDisconnectFn = disconnect
  }

  function hasPermission(code: string): boolean {
    return permissions.value.includes(code)
  }

  async function login(params: LoginParams) {
    const res = await loginApi(params)
    const data = res.data
    setToken(data.token)
    userInfo.value = {
      id: data.user.id,
      username: data.user.username,
      nickname: data.user.nickname,
      email: '',
      phone: '',
      avatar: data.user.avatar,
      status: 1,
      role_id: data.user.role_id,
      created_at: '',
      updated_at: '',
    }
    const [permRes, menuRes] = await Promise.all([
      getCurrentPermissionsApi(),
      getCurrentMenusApi(),
    ])
    permissions.value = permRes.data
    const appStore = useAppStore()
    appStore.menuTree = menuRes.data
    wsConnectFn?.()
    return data
  }

  async function fetchUserInfo() {
    const [profileRes, permRes] = await Promise.all([
      getUserProfileApi(),
      getCurrentPermissionsApi(),
    ])
    userInfo.value = profileRes.data
    permissions.value = permRes.data
  }

  function logout() {
    wsDisconnectFn?.()
    setToken('')
    userInfo.value = null
    permissions.value = []
    menuTree.value = []
  }

  function setToken(value: string) {
    token.value = value
    if (value) {
      localStorage.setItem('token', value)
    } else {
      localStorage.removeItem('token')
    }
  }

  return {
    token,
    userInfo,
    permissions,
    menuTree,
    isLoggedIn,
    registerWs,
    hasPermission,
    login,
    fetchUserInfo,
    logout,
  }
})
