<template>
  <aside class="app-sidebar" :class="{ collapsed: appStore.sidebarCollapsed }">
    <div class="sidebar-logo">
      <span class="logo-icon">G</span>
      <transition name="fade">
        <span v-if="!appStore.sidebarCollapsed" class="logo-text">Gsystes</span>
      </transition>
    </div>

    <el-scrollbar>
      <el-menu
        :default-active="activeMenu"
        :collapse="appStore.sidebarCollapsed"
        :unique-opened="true"
        background-color="transparent"
        :text-color="'var(--color-text-secondary)'"
        :active-text-color="'var(--color-primary)'"
        router
        @select="onMenuSelect"
      >
        <template v-for="menu in filteredMenuItems" :key="menu.path">
          <el-sub-menu v-if="menu.children && menu.children.length > 0" :index="menu.path">
            <template #title>
              <el-icon v-if="menu.icon"><component :is="iconMap[menu.icon]" /></el-icon>
              <span>{{ menu.title }}</span>
            </template>
            <template v-for="child in menu.children" :key="child.path">
              <el-menu-item :index="child.path">
                <span>{{ child.title }}</span>
              </el-menu-item>
            </template>
          </el-sub-menu>
          <el-menu-item v-else :index="menu.path">
            <el-icon v-if="menu.icon"><component :is="iconMap[menu.icon]" /></el-icon>
            <template #title>{{ menu.title }}</template>
          </el-menu-item>
        </template>
      </el-menu>
    </el-scrollbar>
  </aside>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'
import { routes } from '@/router'
import type { Component } from 'vue'
import {
  Monitor, User, Key, Lock, Document, UserFilled,
} from '@element-plus/icons-vue'

defineOptions({ name: 'AppSidebar' })

const emit = defineEmits<{
  'item-click': []
}>()

const route = useRoute()
const appStore = useAppStore()
const authStore = useAuthStore()

const iconMap: Record<string, Component> = {
  Monitor, User, Key, Lock, Document, UserFilled,
}

interface MenuItem {
  path: string
  title: string
  icon?: string
  children?: MenuItem[]
  permission?: string
  hidden?: boolean
}

const activeMenu = computed(() => route.path)

const filteredMenuItems = computed(() => {
  const layoutRoute = routes.find((r) => r.path === '/')
  if (!layoutRoute?.children) return []

  return layoutRoute.children
    .filter((child) => {
      if (child.meta?.hidden) return false
      const perm = child.meta?.permission as string | undefined
      if (perm) return authStore.hasPermission(perm)
      return true
    })
    .map((child) => ({
      path: `/${child.path}`,
      title: (child.meta?.title as string) || child.path,
      icon: child.meta?.icon as string | undefined,
    }))
})

function onMenuSelect() {
  emit('item-click')
}
</script>

<style scoped lang="scss">
.app-sidebar {
  width: var(--sidebar-width);
  height: 100%;
  min-height: 0;
  background: var(--color-bg-sidebar);
  border-right: 1px solid var(--color-border-light);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
  transition: width var(--transition-normal);

  &.collapsed {
    width: 64px;
  }
}

.sidebar-logo {
  height: var(--header-height);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  flex-shrink: 0;

  .logo-icon {
    width: 36px;
    height: 36px;
    background: linear-gradient(135deg, var(--color-primary), var(--color-primary-light));
    border-radius: var(--radius-sm);
    display: flex;
    align-items: center;
    justify-content: center;
    color: #fff;
    font-weight: 700;
    font-size: 18px;
    flex-shrink: 0;
  }
  .logo-text {
    font-size: 17px;
    font-weight: 700;
    color: var(--color-text);
  }
}

:deep(.el-scrollbar) {
  flex: 1;
  min-height: 0;
}

.el-menu {
  border-right: none;
  padding: 8px 0;
}

:deep(.el-menu-item),
:deep(.el-sub-menu__title) {
  border-radius: 0;
  margin: 2px 8px;
  border-radius: var(--radius-sm);
}

:deep(.el-menu-item.is-active) {
  background: var(--color-bg-active) !important;
}

@media (max-width: 768px) {
  .app-sidebar {
    display: none;
  }
}
</style>
