import { createRouter, createWebHashHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login/LoginView.vue'),
    meta: { title: '登录', noAuth: true },
  },
  {
    path: '/',
    component: () => import('@/components/layout/DefaultLayout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/DashboardView.vue'),
        meta: { title: '仪表盘', icon: 'Monitor' },
      },
      {
        path: 'users',
        name: 'UserList',
        component: () => import('@/views/user/UserList.vue'),
        meta: { title: '用户管理', icon: 'User', permission: 'user:read', keepAlive: true },
      },
      {
        path: 'roles',
        name: 'RoleList',
        component: () => import('@/views/role/RoleList.vue'),
        meta: { title: '角色管理', icon: 'Key', permission: 'role:read', keepAlive: true },
      },
      {
        path: 'permissions',
        name: 'PermissionManage',
        component: () => import('@/views/permission/PermissionManage.vue'),
        meta: { title: '权限管理', icon: 'Lock', permission: 'perm:manage', keepAlive: true },
      },
      {
        path: 'logs',
        name: 'OperationLog',
        component: () => import('@/views/log/OperationLog.vue'),
        meta: { title: '操作日志', icon: 'Document', permission: 'log:read', keepAlive: true },
      },
      {
        path: 'profile',
        name: 'Profile',
        component: () => import('@/views/profile/ProfileInfo.vue'),
        meta: { title: '个人中心', icon: 'UserFilled', hidden: true },
      },
    ],
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/dashboard',
  },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

export default router
export { routes }
