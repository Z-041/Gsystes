<template>
  <div class="crud-table">
    <div v-if="searchFields && searchFields.length > 0" class="search-bar">
      <el-form :model="searchForm" inline>
        <el-form-item
          v-for="field in searchFields"
          :key="field.prop"
          :label="field.label"
        >
          <el-input
            v-if="field.type === 'input'"
            v-model="searchForm[field.prop]"
            :placeholder="field.placeholder"
            clearable
            :style="{ width: field.width || '180px' }"
            @clear="handleSearch"
            @keyup.enter="handleSearch"
          />
          <el-select
            v-else-if="field.type === 'select'"
            v-model="searchForm[field.prop]"
            :placeholder="field.placeholder || '全部'"
            clearable
            :style="{ width: field.width || '140px' }"
            @change="handleSearch"
          >
            <el-option
              v-for="opt in field.options"
              :key="opt.value"
              :label="opt.label"
              :value="opt.value"
            />
          </el-select>
          <el-date-picker
            v-else-if="field.type === 'date-range'"
            v-model="searchForm[field.prop]"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            :style="{ width: field.width || '260px' }"
            value-format="YYYY-MM-DD"
            @change="handleSearch"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="toolbar" v-if="toolbarActions.length > 0 || batchActions.length > 0">
      <div class="toolbar-left">
        <template v-for="action in toolbarActions" :key="action.label">
          <el-button
            v-if="!action.onBatch"
            :type="action.type || 'default'"
            :icon="action.icon ? toIcon(action.icon) : undefined"
            v-permission="action.permission || ''"
            @click="action.onClick?.()"
          >
            {{ action.label }}
          </el-button>
          <el-button
            v-else
            :type="action.type || 'default'"
            :icon="action.icon ? toIcon(action.icon) : undefined"
            :disabled="selectedRows.length === 0"
            v-permission="action.permission || ''"
            @click="action.onBatch(selectedRows)"
          >
            {{ action.label }}
            <span v-if="selectedRows.length > 0" class="batch-count">({{ selectedRows.length }})</span>
          </el-button>
        </template>
      </div>
    </div>

    <el-table
      ref="tableRef"
      :data="tableData"
      v-loading="loading"
      stripe
      :row-key="rowKey"
      @selection-change="onSelectionChange"
    >
      <el-table-column v-if="showSelection" type="selection" width="50" />
      <el-table-column v-if="showIndex" type="index" width="60" label="序号" />
      <el-table-column v-if="showExpand" type="expand" width="40">
        <template #default="{ row }">
          <slot name="expand" :row="row" />
        </template>
      </el-table-column>

      <template v-for="col in columns" :key="col.prop">
        <el-table-column
          v-if="!col.slot"
          :prop="col.prop"
          :label="col.label"
          :width="col.width"
          :min-width="col.minWidth"
          :fixed="col.fixed"
          :show-overflow-tooltip="col.showOverflowTooltip"
          :align="col.align"
        />
        <el-table-column
          v-else
          :label="col.label"
          :width="col.width"
          :min-width="col.minWidth"
          :fixed="col.fixed"
          :align="col.align"
        >
          <template #default="{ row }">
            <slot :name="`column-${col.slot}`" :row="row" />
          </template>
        </el-table-column>
      </template>

      <el-table-column
        v-if="rowActions.length > 0"
        label="操作"
        :width="actionColumnWidth"
        fixed="right"
      >
        <template #default="{ row }">
          <template v-for="action in rowActions" :key="(typeof action.label === 'function' ? action.label(row) : action.label)">
            <template v-if="!action.visible || action.visible(row)">
              <template v-if="action.confirm">
                <el-popconfirm
                  :title="typeof action.confirm === 'function' ? action.confirm(row) : action.confirm"
                  @confirm="action.onClick(row)"
                >
                  <template #reference>
                    <el-button
                      link
                      :type="getActionType(action, row)"
                      :icon="action.icon ? toIcon(action.icon) : undefined"
                      v-permission="action.permission || ''"
                    >
                      {{ typeof action.label === 'function' ? action.label(row) : action.label }}
                    </el-button>
                  </template>
                </el-popconfirm>
              </template>
              <el-button
                v-else
                link
                :type="getActionType(action, row)"
                :icon="action.icon ? toIcon(action.icon) : undefined"
                v-permission="action.permission || ''"
                @click="action.onClick(row)"
              >
                {{ typeof action.label === 'function' ? action.label(row) : action.label }}
              </el-button>
            </template>
          </template>
        </template>
      </el-table-column>
    </el-table>

    <TablePagination
      v-if="total > 0"
      :page="currentPage"
      :page-size="currentPageSize"
      :total="total"
      @change="onPageChange"
    />
  </div>
</template>

<script setup lang="ts" generic="T extends Record<string, any>">
import { ref, reactive, computed, unref, onMounted, onActivated } from 'vue'
import type { PageData } from '@/types/api'
import type { Component } from 'vue'
import { useCache } from '@/composables/useCache'
import TablePagination from './TablePagination.vue'

interface ColumnConfig {
  prop: string
  label: string
  width?: string | number
  minWidth?: string | number
  fixed?: 'left' | 'right'
  showOverflowTooltip?: boolean
  slot?: string
  align?: 'left' | 'center' | 'right'
}

interface SearchField {
  prop: string
  label: string
  type: 'input' | 'select' | 'date-range'
  placeholder?: string
  options?: { label: string; value: any }[]
  width?: string
}

interface ToolbarAction {
  label: string
  icon?: string
  type?: 'primary' | 'success' | 'warning' | 'danger' | 'default'
  permission?: string
  onClick?: () => void
  onBatch?: (selected: T[]) => void
}

interface RowAction {
  label: string | ((row: T) => string)
  icon?: string
  type?: string | ((row: T) => string)
  permission?: string
  confirm?: string | ((row: T) => string)
  visible?: (row: T) => boolean
  onClick: (row: T) => void
}

interface CrudTableConfig<T> {
  fetchApi: (params: Record<string, any>) => Promise<{ code: number; data: PageData<T>; message: string }>
  columns: ColumnConfig[]
  rowKey?: string
  pageSize?: number

  searchFields?: SearchField[]
  toolbarActions?: ToolbarAction[]
  batchActions?: ToolbarAction[]
  rowActions?: RowAction[]

  showSelection?: boolean
  showExpand?: boolean
  showIndex?: boolean
  actionColumnWidth?: string

  cacheKey?: string
  extraParams?: () => Record<string, any>
}

const props = withDefaults(defineProps<{
  config: CrudTableConfig<T>
  iconMap?: Record<string, Component>
}>(), {
  iconMap: () => ({}),
})

defineOptions({ name: 'CrudTable' })

const searchFields = computed(() =>
  (props.config.searchFields || []).map((f) => ({
    ...f,
    options: f.options ? unref(f.options) : f.options,
  })),
)
const toolbarActions = computed(() => props.config.toolbarActions || [])
const batchActions = computed(() => props.config.batchActions || [])
const rowActions = computed(() => props.config.rowActions || [])
const columns = computed(() => props.config.columns)
const rowKey = computed(() => props.config.rowKey || 'id')
const showSelection = computed(() => props.config.showSelection || false)
const showExpand = computed(() => props.config.showExpand || false)
const showIndex = computed(() => props.config.showIndex || false)
const actionColumnWidth = computed(() => props.config.actionColumnWidth || '280px')

const tableData = ref<T[]>([])
const total = ref(0)
const loading = ref(false)
const selectedRows = ref<T[]>([])
const currentPage = ref(1)
const currentPageSize = ref(props.config.pageSize || 20)

const searchForm = reactive<Record<string, any>>({})
if (searchFields.value.length > 0) {
  searchFields.value.forEach((f) => {
    searchForm[f.prop] = undefined
  })
}

const { setup: setupCache } = useCache()
if (props.config.cacheKey) {
  setupCache(props.config.cacheKey, fetchData)
}

function toIcon(name: string): Component | undefined {
  return props.iconMap?.[name]
}

function getActionType(action: RowAction, row: T): string {
  const t = action.type
  if (typeof t === 'function') return t(row) || 'primary'
  return t || 'primary'
}

function onSelectionChange(rows: T[]) {
  selectedRows.value = rows
}

async function fetchData() {
  loading.value = true
  try {
    const params: Record<string, any> = {
      page: currentPage.value,
      page_size: currentPageSize.value,
    }

    searchFields.value.forEach((f) => {
      if (f.type === 'date-range' && searchForm[f.prop]) {
        const range = searchForm[f.prop] as string[]
        if (range && range.length === 2) {
          params.start_time = range[0]
          params.end_time = range[1]
        }
      } else if (searchForm[f.prop] !== undefined && searchForm[f.prop] !== '') {
        params[f.prop] = searchForm[f.prop]
      }
    })

    if (props.config.extraParams) {
      Object.assign(params, props.config.extraParams())
    }

    const res = await props.config.fetchApi(params)
    const d = res.data
    tableData.value = d.list
    total.value = d.total
  } catch {
    // handled by axios interceptor
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  currentPage.value = 1
  fetchData()
}

function handleReset() {
  searchFields.value.forEach((f) => {
    searchForm[f.prop] = undefined
  })
  currentPage.value = 1
  fetchData()
}

function onPageChange(params: { page: number; page_size: number }) {
  currentPage.value = params.page
  currentPageSize.value = params.page_size
  fetchData()
}

defineExpose({
  refresh: fetchData,
  reset: handleReset,
  getSelected: () => selectedRows.value,
})

onMounted(() => {
  fetchData()
})

onActivated(() => {
  // Cache refresh handled by useCache
})
</script>

<style scoped lang="scss">
.crud-table {
  .search-bar {
    padding: 16px;
    background: var(--color-bg-card);
    border-radius: var(--radius-md);
    margin-bottom: 16px;
    border: 1px solid var(--color-border-light);
    :deep(.el-form-item) {
      margin-bottom: 0;
    }
  }
  .toolbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
  }
  .toolbar-left {
    display: flex;
    gap: 8px;
    flex-wrap: wrap;
  }
  .batch-count {
    color: var(--color-primary);
    font-weight: 600;
  }
  :deep(.el-table) {
    border-radius: var(--radius-md);
    overflow: hidden;
  }
}
</style>
