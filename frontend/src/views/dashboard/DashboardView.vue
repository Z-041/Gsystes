<template>
  <div class="dashboard-page">
    <div class="page-header">
      <h2>仪表盘</h2>
    </div>

    <el-row :gutter="16">
      <el-col :xs="12" :sm="12" :md="6" v-for="card in statCards" :key="card.label">
        <StatCard :label="card.label" :value="card.value" :bg="card.bg">
          <template #icon>
            <el-icon :size="20"><component :is="card.icon" /></el-icon>
          </template>
        </StatCard>
      </el-col>
    </el-row>

    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :span="24">
        <el-card shadow="never" class="recent-card">
          <template #header>
            <div class="card-header-row">
              <span class="card-title">近期操作日志</span>
              <span class="ws-status" :class="{ live: wsConnected }">
                <span class="ws-dot" />
                {{ wsConnected ? '实时' : '离线' }}
              </span>
            </div>
          </template>
          <el-table :data="recentLogs" v-loading="logLoading" stripe>
            <el-table-column prop="username" label="操作人" width="120" />
            <el-table-column prop="module" label="模块" width="100" />
            <el-table-column prop="action" label="操作" width="100" />
            <el-table-column prop="path" label="请求路径" min-width="200" show-overflow-tooltip />
            <el-table-column prop="ip" label="IP" width="140" />
            <el-table-column prop="created_at" label="时间" width="180" />
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, shallowRef, onMounted, onUnmounted } from 'vue'
import { User, Key, TrendCharts, Document } from '@element-plus/icons-vue'
import { getDashboardStatsApi } from '@/api/dashboard'
import { getOperationLogListApi } from '@/api/log'
import type { OperationLogInfo } from '@/types/log'
import { useWebSocket } from '@/composables/useWebSocket'
import type { WsStatUpdate, WsLogEntry } from '@/composables/useWebSocket'
import StatCard from '@/components/common/StatCard.vue'
import type { Component } from 'vue'

defineOptions({ name: 'DashboardView' })

const statsLoading = ref(false)
const logLoading = ref(false)

const statCards = shallowRef<
  { label: string; value: string; icon: Component; bg: string }[]
>([
  { label: '用户总数', value: '--', icon: User, bg: 'linear-gradient(135deg, #6366f1, #a78bfa)' },
  { label: '角色数', value: '--', icon: Key, bg: 'linear-gradient(135deg, #22c55e, #4ade80)' },
  { label: '今日操作', value: '--', icon: TrendCharts, bg: 'linear-gradient(135deg, #f97316, #fb923c)' },
  { label: '操作日志', value: '--', icon: Document, bg: 'linear-gradient(135deg, #06b6d4, #22d3ee)' },
])

const recentLogs = ref<OperationLogInfo[]>([])
const MAX_LOGS = 50

const { connected: wsConnected, subscribe } = useWebSocket()

let unsubStat: (() => void) | null = null
let unsubLog: (() => void) | null = null

onMounted(async () => {
  statsLoading.value = true
  try {
    const res = await getDashboardStatsApi()
    const s = res.data
    statCards.value = statCards.value.map((c) => {
      if (c.label === '用户总数') return { ...c, value: String(s.user_count) }
      if (c.label === '角色数') return { ...c, value: String(s.role_count) }
      if (c.label === '今日操作' || c.label === '操作日志')
        return { ...c, value: String(s.today_log_count) }
      return c
    })
  } catch {
    // keep '--'
  } finally {
    statsLoading.value = false
  }

  logLoading.value = true
  try {
    const res = await getOperationLogListApi({ page: 1, page_size: 8 })
    recentLogs.value = res.data.list
  } catch {
    recentLogs.value = []
  } finally {
    logLoading.value = false
  }

  unsubStat = subscribe('stat_update', (payload: WsStatUpdate) => {
    statCards.value = statCards.value.map((c) => {
      if (c.label === '用户总数') return { ...c, value: String(payload.user_count) }
      if (c.label === '角色数') return { ...c, value: String(payload.role_count) }
      if (c.label === '今日操作' || c.label === '操作日志')
        return { ...c, value: String(payload.today_log_count) }
      return c
    })
  })

  unsubLog = subscribe('log_entry', (payload: WsLogEntry) => {
    const entry: OperationLogInfo = {
      id: payload.id,
      user_id: 0,
      username: payload.username,
      module: payload.module,
      action: payload.action,
      method: payload.method,
      path: payload.path,
      ip: payload.ip,
      duration: payload.duration,
      request_body: '',
      status_code: payload.status_code,
      created_at: payload.created_at,
    }
    recentLogs.value.unshift(entry)
    if (recentLogs.value.length > MAX_LOGS) {
      recentLogs.value.pop()
    }
  })
})

onUnmounted(() => {
  unsubStat?.()
  unsubLog?.()
})
</script>

<style scoped lang="scss">
.page-header {
  margin-bottom: 16px;
  h2 {
    font-size: var(--font-size-xl);
    font-weight: 700;
    color: var(--color-text);
  }
}

.recent-card {
  border-radius: var(--radius-md);
  border: 1px solid var(--color-border-light);

  .card-header-row {
    display: flex;
    align-items: center;
    gap: 10px;
  }
  .card-title {
    font-weight: 600;
  }
}

.ws-status {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 11px;
  color: var(--color-text-muted);
  padding: 2px 8px;
  border-radius: 20px;
  background: #f1f5f9;

  &.live {
    color: #16a34a;
    background: #dcfce7;
  }
}

.ws-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #cbd5e1;
}

.ws-status.live .ws-dot {
  background: #22c55e;
}
</style>
