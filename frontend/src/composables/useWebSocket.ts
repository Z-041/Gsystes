import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useNotificationStore } from '@/stores/notification'

export interface WsLogEntry {
  id: number
  username: string
  module: string
  action: string
  method: string
  path: string
  ip: string
  duration: number
  status_code: number
  created_at: string
}

export interface WsStatUpdate {
  user_count: number
  role_count: number
  today_log_count: number
}

export interface WsNotification {
  username: string
  title: string
  message: string
}

type WsMessageHandler = (payload: any) => void

const INITIAL_DELAY = 1000
const MAX_DELAY = 30000
const BACKOFF_RATE = 2

let ws: WebSocket | null = null
let reconnectTimer: ReturnType<typeof setTimeout> | null = null
let reconnectDelay = INITIAL_DELAY
let destroyed = false

const handlers = new Map<string, Set<WsMessageHandler>>()

const connected = ref(false)

function buildWsUrl(token: string): string {
  if (!token) return ''
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  return `${protocol}//${window.location.host}/api/v1/../ws?token=${token}`
}

function scheduleReconnect() {
  if (destroyed) return
  if (reconnectTimer) clearTimeout(reconnectTimer)
  reconnectTimer = setTimeout(() => {
    reconnectDelay = Math.min(reconnectDelay * BACKOFF_RATE, MAX_DELAY)
    connect()
  }, reconnectDelay)
}

function connect() {
  destroyed = false
  const authStore = useAuthStore()
  if (!authStore.token) {
    scheduleReconnect()
    return
  }

  const url = buildWsUrl(authStore.token)
  if (!url) {
    scheduleReconnect()
    return
  }

  ws = new WebSocket(url)

  ws.onopen = () => {
    connected.value = true
    reconnectDelay = INITIAL_DELAY
  }

  ws.onmessage = (event: MessageEvent) => {
    try {
      const msg = JSON.parse(event.data)
      const typeHandlers = handlers.get(msg.type)
      if (typeHandlers) {
        typeHandlers.forEach((fn) => fn(msg.payload))
      }
      if (msg.type === 'notification') {
        const payload = msg.payload as WsNotification
        const notifStore = useNotificationStore()
        notifStore.push(payload.username, payload.title, payload.message)
      }
    } catch {
      // ignore malformed
    }
  }

  ws.onclose = () => {
    connected.value = false
    ws = null
    scheduleReconnect()
  }

  ws.onerror = () => {
    ws?.close()
  }
}

function disconnect() {
  destroyed = true
  handlers.clear()
  if (reconnectTimer) {
    clearTimeout(reconnectTimer)
    reconnectTimer = null
  }
  if (ws) {
    ws.onclose = null
    ws.close()
    ws = null
  }
  connected.value = false
}

function subscribe(type: string, handler: WsMessageHandler): () => void {
  if (!handlers.has(type)) {
    handlers.set(type, new Set())
  }
  handlers.get(type)!.add(handler)
  return () => {
    handlers.get(type)?.delete(handler)
  }
}

export function useWebSocket() {
  return {
    connected,
    subscribe,
    connect,
    disconnect,
    registerGlobal: () => {
      const authStore = useAuthStore()
      authStore.registerWs(connect, disconnect)
      if (authStore.token) {
        connect()
      }
    },
  }
}
