<template>
  <el-dialog
    v-model="visible"
    :title="config.title || '详情'"
    :width="config.width || '600px'"
    :close-on-click-modal="false"
    destroy-on-close
  >
    <div v-loading="loading" class="detail-body">
      <template v-if="detail">
        <div class="detail-header">
          <el-avatar
            v-if="config.header.avatar"
            :size="48"
            :src="typeof config.header.avatar === 'function' ? config.header.avatar(detail) : config.header.avatar"
            style="background: linear-gradient(135deg, var(--color-primary), var(--color-primary-light)); flex-shrink: 0"
          >
            {{ (typeof config.header.title === 'function' ? config.header.title(detail) : config.header.title).charAt(0) }}
          </el-avatar>
          <div class="detail-header-info">
            <div class="detail-title">
              {{ typeof config.header.title === 'function' ? config.header.title(detail) : config.header.title }}
            </div>
            <div class="detail-subtitle">
              {{ typeof config.header.subtitle === 'function' ? config.header.subtitle(detail) : config.header.subtitle }}
            </div>
          </div>
          <el-tag
            v-if="config.header.status"
            :type="config.header.status(detail).type"
            size="small"
            class="detail-status-tag"
          >
            {{ config.header.status(detail).label }}
          </el-tag>
        </div>

        <el-divider />

        <template v-for="(section, si) in config.sections" :key="si">
          <div v-if="section.title" class="section-title">{{ section.title }}</div>
          <el-descriptions :column="2" border>
            <template v-for="item in section.items" :key="item.label">
              <el-descriptions-item :label="item.label" :span="item.span || 1">
                <template v-if="item.slot">
                  <slot :name="item.slot" :detail="detail" />
                </template>
                <template v-else>
                  {{ typeof item.value === 'function' ? item.value(detail) : item.value }}
                </template>
              </el-descriptions-item>
            </template>
          </el-descriptions>
        </template>
      </template>
      <el-empty v-else-if="!loading" description="暂无数据" />
    </div>

    <template #footer>
      <el-button @click="visible = false">关闭</el-button>
      <template v-if="config.actions">
        <el-button
          v-for="action in config.actions"
          :key="action.label"
          :type="action.type || 'default'"
          :icon="action.icon"
          @click="action.onClick(detail)"
        >
          {{ action.label }}
        </el-button>
      </template>
    </template>
  </el-dialog>
</template>

<script setup lang="ts" generic="T extends Record<string, any>">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import type { Component } from 'vue'

interface DetailAction {
  label: string
  icon?: Component
  type?: 'primary' | 'default'
  onClick: (detail: any) => void
}

interface DetailSection {
  title?: string
  items: {
    label: string
    value: string | ((d: any) => string)
    slot?: string
    span?: number
  }[]
}

interface DetailDialogConfig<T = Record<string, any>> {
  title?: string
  fetchDetail: (id: number) => Promise<T>
  width?: string
  header: {
    avatar?: string | ((d: any) => string)
    title: string | ((d: any) => string)
    subtitle: string | ((d: any) => string)
    status?: (d: any) => { label: string; type: 'success' | 'danger' | 'warning' | 'info' }
  }
  sections: DetailSection[]
  actions?: DetailAction[]
}

const props = defineProps<{
  config: DetailDialogConfig<T>
}>()

defineOptions({ name: 'DetailDialog' })

const visible = ref(false)
const loading = ref(false)
const detail = ref<T | null>(null)

async function open(id: number) {
  visible.value = true
  loading.value = true
  detail.value = null
  try {
    detail.value = await props.config.fetchDetail(id)
  } catch {
    ElMessage.error('加载详情失败')
    visible.value = false
  } finally {
    loading.value = false
  }
}

defineExpose({ open })
</script>

<style scoped lang="scss">
.detail-body {
  min-height: 80px;
}
.detail-header {
  display: flex;
  align-items: center;
  gap: 14px;
}
.detail-header-info {
  flex: 1;
}
.detail-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--color-text);
}
.detail-subtitle {
  font-size: 13px;
  color: var(--color-text-muted);
  margin-top: 4px;
}
.detail-status-tag {
  flex-shrink: 0;
}
.section-title {
  font-size: 14px;
  font-weight: 600;
  color: var(--color-text);
  margin: 16px 0 10px;
}
</style>
