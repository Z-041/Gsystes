<template>
  <div class="role-list-page">
    <div class="page-header">
      <h2>角色管理</h2>
    </div>

    <CrudTable ref="tableRef" :config="tableConfig" :icon-map="iconMap">
      <template #column-status="{ row }">
        <StatusTag :status="(row as RoleInfo).status" />
      </template>
    </CrudTable>

    <FormDialog ref="roleFormRef" :config="roleFormConfig" />

    <AssignPermissionDialog ref="permRef" @success="() => tableRef?.refresh()" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus, Key, Edit, Delete } from '@element-plus/icons-vue'
import type { RoleInfo } from '@/types/role'
import CrudTable from '@/components/common/CrudTable.vue'
import FormDialog from '@/components/common/FormDialog.vue'
import StatusTag from '@/components/common/StatusTag.vue'
import AssignPermissionDialog from '@/components/business/AssignPermissionDialog.vue'
import { getRoleListApi, getRoleApi, createRoleApi, updateRoleApi, deleteRoleApi } from '@/api/role'
import { invalidateCache } from '@/composables/useCache'

defineOptions({ name: 'RoleList' })

const iconMap = { Plus, Key, Edit, Delete }

const tableRef = ref<any>(null)
const roleFormRef = ref<any>(null)
const permRef = ref<any>(null)

async function handleDelete(row: RoleInfo) {
  try {
    await deleteRoleApi(row.id)
    ElMessage.success('删除成功')
    tableRef.value?.refresh()
    invalidateCache('RoleList')
  } catch { /* */ }
}

const tableConfig = {
  fetchApi: (params: Record<string, any>) => getRoleListApi(params),
  columns: [
    { prop: 'name', label: '角色名称', width: 160 },
    { prop: 'code', label: '角色编码', width: 160 },
    { prop: 'description', label: '描述', minWidth: 200 },
    { prop: 'status', label: '状态', width: 90, slot: 'status' },
    { prop: 'created_at', label: '创建时间', width: 180 },
  ],
  toolbarActions: [
    { label: '新增角色', icon: 'Plus', type: 'primary' as const, permission: 'role:create',
      onClick: () => roleFormRef.value?.open('create') },
  ],
  rowActions: [
    { label: '分配权限', icon: 'Key', type: 'primary' as const, permission: 'perm:assign',
      onClick: (row: RoleInfo) => permRef.value?.open(row.id) },
    { label: '编辑', icon: 'Edit', type: 'primary' as const, permission: 'role:update',
      onClick: (row: RoleInfo) => roleFormRef.value?.open('edit', row.id) },
    { label: '删除', icon: 'Delete', type: 'danger' as const, permission: 'role:delete',
      confirm: '确认删除该角色？', onClick: (row: RoleInfo) => handleDelete(row) },
  ],
  cacheKey: 'RoleList',
  actionColumnWidth: '280px',
}

const roleFormConfig = {
  title: (mode: string) => mode === 'create' ? '新增角色' : '编辑角色',
  width: '520px',
  fields: [
    { prop: 'name', label: '角色名称', type: 'input' as const,
      rules: [{ required: true, message: '请输入角色名称', trigger: 'blur' as const }],
      placeholder: '如：系统管理员' },
    { prop: 'code', label: '角色编码', type: 'input' as const,
      rules: [{ required: true, message: '请输入角色编码', trigger: 'blur' as const }],
      placeholder: '如：admin' },
    { prop: 'description', label: '描述', type: 'textarea' as const, placeholder: '选填' },
  ],
  fetchDetail: (id: number) => getRoleApi(id).then((r) => r.data),
  onSubmit: async (payload: any, id?: number) => {
    if (id) {
      await updateRoleApi(id, payload)
    } else {
      await createRoleApi(payload)
    }
  },
  onSuccess: () => {
    tableRef.value?.refresh()
    invalidateCache('RoleList')
  },
}
</script>

<style scoped lang="scss">
.page-header {
  margin-bottom: 16px;
  h2 { font-size: var(--font-size-xl); font-weight: 700; color: var(--color-text); }
}
</style>
