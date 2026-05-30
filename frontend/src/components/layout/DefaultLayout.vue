<template>
  <div class="default-layout">
    <AppHeader />
    <div class="layout-body">
      <AppSidebar />
      <main class="layout-content">
        <router-view v-slot="{ Component }">
          <keep-alive :include="cachedNames">
            <component :is="Component" />
          </keep-alive>
        </router-view>
      </main>
    </div>

    <transition name="fade">
      <div v-if="appStore.mobileMenuOpen" class="mobile-overlay" @click="appStore.closeMobileMenu" />
    </transition>
    <transition name="slide-left">
      <aside v-if="appStore.mobileMenuOpen" class="mobile-sidebar">
        <AppSidebar @item-click="appStore.closeMobileMenu" />
      </aside>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import AppHeader from './AppHeader.vue'
import AppSidebar from './AppSidebar.vue'
import { useAppStore } from '@/stores/app'
import { routes } from '@/router'

defineOptions({ name: 'DefaultLayout' })

const appStore = useAppStore()

const cachedNames = computed(() => {
  const names: string[] = []
  const layoutRoute = routes.find((r) => r.path === '/')
  if (layoutRoute?.children) {
    for (const child of layoutRoute.children) {
      if (child.meta?.keepAlive && child.name) {
        names.push(child.name as string)
      }
    }
  }
  return names
})
</script>

<style scoped lang="scss">
.default-layout {
  height: 100vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.layout-body {
  flex: 1;
  min-height: 0;
  display: flex;
  overflow: hidden;
}

.layout-content {
  flex: 1;
  min-width: 0;
  overflow-y: auto;
  padding: 20px;
  background: var(--color-bg);
}

.mobile-overlay {
  display: none;
}
.mobile-sidebar {
  display: none;
}

@media (max-width: 768px) {
  .layout-content {
    padding: 14px;
  }
  .mobile-overlay {
    display: block;
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.4);
    z-index: 200;
  }
  .mobile-sidebar {
    display: block;
    position: fixed;
    top: 0;
    left: 0;
    bottom: 0;
    width: 220px;
    height: 100vh;
    background: var(--color-bg-sidebar);
    z-index: 201;
    box-shadow: var(--shadow-lg);
    overflow-y: auto;
  }
}
</style>
