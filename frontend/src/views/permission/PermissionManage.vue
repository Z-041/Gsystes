<template>
  <div class="permission-page">
    <div class="page-header">
      <h2>权限管理</h2>
    </div>

    <PermissionTree ref="treeRef" @edit="handleEdit" @delete="triggerRefresh" />

    <FormDialog ref="permFormRef" :config="permFormConfig" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { PermissionInfo } from '@/types/permission'
import PermissionTree from '@/components/business/PermissionTree.vue'
import FormDialog from '@/components/common/FormDialog.vue'
import { getPermissionApi, createPermissionApi, updatePermissionApi } from '@/api/permission'
import { invalidateCache } from '@/composables/useCache'

defineOptions({ name: 'PermissionManage' })

const treeRef = ref<any>(null)
const permFormRef = ref<any>(null)
const editingParentId = ref<number>(0)

function triggerRefresh() {
  treeRef.value?.loadTree()
  invalidateCache('PermissionList')
}

function handleEdit(parentId?: number, data?: PermissionInfo) {
  editingParentId.value = parentId ?? 0
  if (data) {
    permFormRef.value?.open('edit', data.id)
  } else {
    permFormRef.value?.open('create')
  }
}

const permFormConfig = {
  title: (mode: string) => mode === 'create' ? '新增权限' : '编辑权限',
  width: '480px',
  fields: [
    { prop: 'name', label: '名称', type: 'input' as const,
      rules: [{ required: true, message: '请输入权限名称', trigger: 'blur' as const }],
      placeholder: '权限名称' },
    { prop: 'code', label: '编码', type: 'input' as const,
      rules: [{ required: true, message: '请输入权限编码', trigger: 'blur' as const }],
      placeholder: '如：user:read' },
    { prop: 'type', label: '类型', type: 'radio' as const,
      rules: [{ required: true, message: '请选择类型', trigger: 'change' as const }],
      options: [{ label: '菜单', value: 1 }, { label: '按钮', value: 2 }] },
    { prop: 'sort', label: '排序', type: 'number' as const, props: { min: 0, max: 999 } },
  ],
  fetchDetail: (id: number) => getPermissionApi(id).then((r) => r.data),
  onSubmit: async (payload: any, id?: number) => {
    const data = { ...payload, parent_id: editingParentId.value || 0 }
    if (id) {
      await updatePermissionApi(id, data)
    } else {
      await createPermissionApi(data)
    }
  },
  onSuccess: () => {
    treeRef.value?.loadTree()
    invalidateCache('PermissionList')
  },
}
</script>

<style scoped lang="scss">
.page-header {
  margin-bottom: 16px;
  h2 { font-size: var(--font-size-xl); font-weight: 700; color: var(--color-text); }
}
</style>
