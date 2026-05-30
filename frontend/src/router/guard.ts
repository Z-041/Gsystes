import type { Router } from 'vue-router'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'

NProgress.configure({ showSpinner: false, trickleSpeed: 100 })

const WHITE_LIST = ['/login']

export function setupRouterGuard(router: Router) {
  router.beforeEach(async (to, _from, next) => {
    NProgress.start()

    const authStore = useAuthStore()

    if (WHITE_LIST.includes(to.path)) {
      if (authStore.token && to.path === '/login') {
        next('/dashboard')
      } else {
        next()
      }
      return
    }

    if (!authStore.token) {
      next(`/login?redirect=${encodeURIComponent(to.fullPath)}`)
      return
    }

    if (!authStore.userInfo) {
      try {
        await authStore.fetchUserInfo()
      } catch {
        authStore.logout()
        NProgress.done()
        next(`/login?redirect=${encodeURIComponent(to.fullPath)}`)
        return
      }
    }

    const requiredPermission = to.meta.permission as string | undefined
    if (requiredPermission && !authStore.hasPermission(requiredPermission)) {
      ElMessage.warning('没有访问权限')
      next('/dashboard')
      return
    }

    next()
  })

  router.afterEach((to) => {
    document.title = `${to.meta.title || 'Gsystes'} - Gsystes 管理系统`
    NProgress.done()
  })
}
