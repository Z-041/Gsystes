import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface NotificationItem {
  id: string
  username: string
  title: string
  message: string
  timestamp: number
  read: boolean
}

let autoId = 0

export const useNotificationStore = defineStore('notification', () => {
  const items = ref<NotificationItem[]>([])
  const unreadCount = ref(0)

  function push(username: string, title: string, message: string) {
    const item: NotificationItem = {
      id: `notif-${Date.now()}-${++autoId}`,
      username,
      title,
      message,
      timestamp: Date.now(),
      read: false,
    }
    items.value.unshift(item)
    unreadCount.value++
    if (items.value.length > 50) {
      items.value.pop()
    }
  }

  function markAllRead() {
    items.value.forEach((i) => (i.read = true))
    unreadCount.value = 0
  }

  function clearAll() {
    items.value = []
    unreadCount.value = 0
  }

  return { items, unreadCount, push, markAllRead, clearAll }
})
