import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { PermissionInfo } from '@/types/permission'

export const useAppStore = defineStore('app', () => {
  const sidebarCollapsed = ref(false)
  const mobileMenuOpen = ref(false)
  const menuTree = ref<PermissionInfo[]>([])

  const isDark = ref(localStorage.getItem('theme') === 'dark')
  if (isDark.value) {
    document.documentElement.classList.add('dark')
  }

  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  function toggleMobileMenu() {
    mobileMenuOpen.value = !mobileMenuOpen.value
  }

  function closeMobileMenu() {
    mobileMenuOpen.value = false
  }

  function toggleTheme() {
    isDark.value = !isDark.value
    document.documentElement.classList.toggle('dark', isDark.value)
    localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
  }

  return {
    sidebarCollapsed,
    mobileMenuOpen,
    menuTree,
    isDark,
    toggleSidebar,
    toggleMobileMenu,
    closeMobileMenu,
    toggleTheme,
  }
})
