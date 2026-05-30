<template>
  <header class="app-header">
    <div class="header-left">
      <el-icon class="menu-toggler" @click="appStore.toggleSidebar">
        <Fold v-if="!appStore.sidebarCollapsed" />
        <Expand v-else />
      </el-icon>
      <el-icon class="mobile-toggler" @click="appStore.toggleMobileMenu">
        <Expand />
      </el-icon>
      <el-breadcrumb separator="/">
        <el-breadcrumb-item :to="{ path: '/dashboard' }">首页</el-breadcrumb-item>
        <el-breadcrumb-item v-if="route.meta.title">{{ route.meta.title }}</el-breadcrumb-item>
      </el-breadcrumb>
    </div>

    <div class="header-right">
      <el-badge :value="notifStore.unreadCount" :hidden="notifStore.unreadCount === 0" class="header-badge">
        <el-popover placement="bottom-end" :width="360" trigger="click">
          <template #reference>
            <el-button link>
              <el-icon :size="20"><Bell /></el-icon>
            </el-button>
          </template>
          <div class="notif-popover">
            <div class="notif-header">
              <span>通知</span>
              <el-button link size="small" @click="notifStore.markAllRead">全部已读</el-button>
            </div>
            <el-scrollbar max-height="300px">
              <div v-if="notifStore.items.length === 0" class="notif-empty">暂无通知</div>
              <div
                v-for="item in notifStore.items"
                :key="item.id"
                class="notif-item"
                :class="{ unread: !item.read }"
              >
                <div class="notif-title">
                  <span class="notif-dot" v-if="!item.read" />
                  {{ item.title }}
                </div>
                <div class="notif-msg">{{ item.message }}</div>
                <div class="notif-time">{{ formatTime(item.timestamp) }}</div>
              </div>
            </el-scrollbar>
          </div>
        </el-popover>
      </el-badge>

      <el-button link @click="appStore.toggleTheme" class="theme-btn">
        <el-icon :size="18"><Sunny v-if="appStore.isDark" /><Moon v-else /></el-icon>
      </el-button>

      <el-dropdown trigger="click" @command="handleCommand">
        <div class="user-dropdown-trigger">
          <el-avatar :size="32" style="background: linear-gradient(135deg, var(--color-primary), var(--color-primary-light))">
            {{ avatarLetter }}
          </el-avatar>
          <span class="user-name">{{ authStore.userInfo?.nickname || authStore.userInfo?.username }}</span>
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="profile">
              <el-icon><UserFilled /></el-icon>个人中心
            </el-dropdown-item>
            <el-dropdown-item command="password">
              <el-icon><Key /></el-icon>修改密码
            </el-dropdown-item>
            <el-dropdown-item divided command="logout">
              <el-icon><SwitchButton /></el-icon>退出登录
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>

    <PasswordDialog ref="passwordRef" />
  </header>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useAppStore } from '@/stores/app'
import { useNotificationStore } from '@/stores/notification'
import { Fold, Expand, Bell, Sunny, Moon, UserFilled, Key, SwitchButton } from '@element-plus/icons-vue'
import PasswordDialog from '@/components/common/PasswordDialog.vue'

defineOptions({ name: 'AppHeader' })

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const appStore = useAppStore()
const notifStore = useNotificationStore()

const passwordRef = ref<InstanceType<typeof PasswordDialog>>()

const avatarLetter = computed(() => {
  const u = authStore.userInfo
  return (u?.nickname?.charAt(0) || u?.username?.charAt(0) || 'U')
})

function formatTime(ts: number): string {
  const d = new Date(ts)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`
}

function handleCommand(cmd: string) {
  if (cmd === 'profile') {
    router.push('/profile')
  } else if (cmd === 'password') {
    passwordRef.value?.open()
  } else if (cmd === 'logout') {
    authStore.logout()
    router.push('/login')
  }
}
</script>

<style scoped lang="scss">
.app-header {
  height: var(--header-height);
  background: var(--color-bg-header);
  border-bottom: 1px solid var(--color-border-light);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  flex-shrink: 0;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.menu-toggler {
  font-size: 20px;
  cursor: pointer;
  color: var(--color-text-secondary);
}
.mobile-toggler {
  font-size: 20px;
  cursor: pointer;
  color: var(--color-text-secondary);
  display: none;
}
@media (max-width: 768px) {
  .menu-toggler { display: none; }
  .mobile-toggler { display: block; }
}

.header-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.header-badge {
  :deep(.el-badge__content) { border: none; }
}

.theme-btn {
  font-size: 18px;
}

.user-dropdown-trigger {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 2px 4px;
  border-radius: var(--radius-sm);
  &:hover { background: var(--color-bg-hover); }
}
.user-name {
  font-size: 14px;
  color: var(--color-text);
}

.notif-popover {
  .notif-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-bottom: 10px;
    border-bottom: 1px solid var(--color-border-light);
    margin-bottom: 8px;
    font-weight: 600;
  }
  .notif-empty {
    text-align: center;
    color: var(--color-text-muted);
    padding: 24px 0;
  }
  .notif-item {
    padding: 10px 0;
    border-bottom: 1px solid var(--color-border-light);
    &.unread { background: var(--color-bg-active); margin: 0 -12px; padding: 10px 12px; }
  }
  .notif-title {
    font-weight: 500;
    display: flex;
    align-items: center;
    gap: 6px;
  }
  .notif-dot {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background: var(--color-primary);
    flex-shrink: 0;
  }
  .notif-msg {
    color: var(--color-text-secondary);
    font-size: 13px;
    margin-top: 4px;
  }
  .notif-time {
    color: var(--color-text-muted);
    font-size: 11px;
    margin-top: 4px;
  }
}
</style>
