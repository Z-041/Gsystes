<template>
  <div class="log-page">
    <div class="page-header">
      <h2>操作日志</h2>
    </div>

    <CrudTable ref="tableRef" :config="tableConfig" :icon-map="iconMap">
      <template #expand="{ row }">
        <div class="expand-detail">
          <div class="detail-row">
            <span class="detail-label">请求地址：</span>
            <span>{{ (row as any).path }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">请求方法：</span>
            <el-tag size="small" :type="methodColor((row as any).method)">{{ (row as any).method }}</el-tag>
          </div>
          <div class="detail-row">
            <span class="detail-label">请求体：</span>
            <pre class="code-block">{{ (row as any).request_body || '-' }}</pre>
          </div>
          <div class="detail-row">
            <span class="detail-label">响应码：</span>
            <span>{{ (row as any).status_code }}</span>
          </div>
          <div class="detail-row">
            <span class="detail-label">耗时：</span>
            <span>{{ (row as any).duration }}ms</span>
          </div>
        </div>
      </template>
    </CrudTable>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import CrudTable from '@/components/common/CrudTable.vue'
import { getOperationLogListApi } from '@/api/log'

defineOptions({ name: 'OperationLog' })

const iconMap = {}
const tableRef = ref<any>(null)

function methodColor(method: string): string {
  const map: Record<string, string> = { GET: 'success', POST: 'primary', PUT: 'warning', DELETE: 'danger' }
  return map[method] || 'info'
}

const tableConfig = {
  fetchApi: (params: Record<string, any>) => getOperationLogListApi(params),
  columns: [
    { prop: 'username', label: '操作人', width: 120 },
    { prop: 'module', label: '模块', width: 100 },
    { prop: 'action', label: '操作类型', width: 100 },
    { prop: 'path', label: '请求路径', minWidth: 220, showOverflowTooltip: true },
    { prop: 'ip', label: 'IP地址', width: 140 },
    { prop: 'created_at', label: '操作时间', width: 180 },
  ],
  searchFields: [
    { prop: 'username', label: '操作人', type: 'input' as const, placeholder: '模糊搜索', width: '160px' },
    { prop: 'method', label: '方法', type: 'select' as const,
      options: [
        { label: 'GET', value: 'GET' },
        { label: 'POST', value: 'POST' },
        { label: 'PUT', value: 'PUT' },
        { label: 'DELETE', value: 'DELETE' },
      ],
      width: '120px' },
    { prop: 'date_range', label: '时间范围', type: 'date-range' as const, width: '260px' },
  ],
  showExpand: true,
  cacheKey: 'OperationLog',
  pageSize: 20,
}
</script>

<style scoped lang="scss">
.page-header {
  margin-bottom: 16px;
  h2 { font-size: var(--font-size-xl); font-weight: 700; color: var(--color-text); }
}

.expand-detail {
  padding: 12px 16px;
}

.detail-row {
  display: flex;
  align-items: flex-start;
  margin-bottom: 8px;
}

.detail-label {
  color: var(--color-text-secondary);
  font-size: 13px;
  min-width: 80px;
  flex-shrink: 0;
}

.code-block {
  background: var(--color-bg-hover);
  padding: 6px 10px;
  border-radius: var(--radius-sm);
  font-size: 12px;
  font-family: monospace;
  max-height: 120px;
  overflow-y: auto;
  color: var(--color-text);
  white-space: pre-wrap;
  word-break: break-all;
}
</style>
