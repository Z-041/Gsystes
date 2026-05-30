<template>
  <div class="user-list-page">
    <div class="page-header">
      <h2>用户管理</h2>
    </div>

    <CrudTable ref="tableRef" :config="tableConfig" :icon-map="iconMap">
      <template #column-user="{ row }">
        <div class="user-cell">
          <el-avatar :size="36" style="background: linear-gradient(135deg, var(--color-primary), var(--color-primary-light)); flex-shrink: 0">
            {{ (row as UserInfo).nickname?.charAt(0) || (row as UserInfo).username?.charAt(0) }}
          </el-avatar>
          <div class="user-info">
            <div class="user-name">{{ (row as UserInfo).nickname || (row as UserInfo).username }}</div>
            <div class="user-desc">@{{ (row as UserInfo).username }}</div>
          </div>
        </div>
      </template>
      <template #column-role="{ row }">
        <RoleTag :name="(row as UserInfo).role?.name" />
      </template>
      <template #column-status="{ row }">
        <StatusTag :status="(row as UserInfo).status" />
      </template>
    </CrudTable>

    <FormDialog ref="userFormRef" :config="userFormConfig" />

    <DetailDialog ref="detailRef" :config="detailConfig">
      <template #role="{ detail }">
        <RoleTag :name="(detail as UserInfo).role?.name" />
      </template>
      <template #status="{ detail }">
        <StatusTag :status="(detail as UserInfo).status" />
      </template>
    </DetailDialog>

    <BatchRoleDialog ref="batchRef" :roles="roleList" @success="() => tableRef?.refresh()" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus, Upload, Download, View, Edit, Key, SwitchButton, Delete } from '@element-plus/icons-vue'
import type { UserInfo } from '@/types/user'
import type { RoleSimple } from '@/types/role'
import CrudTable from '@/components/common/CrudTable.vue'
import FormDialog from '@/components/common/FormDialog.vue'
import DetailDialog from '@/components/common/DetailDialog.vue'
import RoleTag from '@/components/common/RoleTag.vue'
import StatusTag from '@/components/common/StatusTag.vue'
import BatchRoleDialog from '@/components/business/BatchRoleDialog.vue'
import { getUserListApi, getUserApi, createUserApi, updateUserApi, deleteUserApi, updateUserStatusApi } from '@/api/user'
import { getAllRolesApi } from '@/api/role'
import { invalidateCache } from '@/composables/useCache'

defineOptions({ name: 'UserList' })

const iconMap = {
  Plus, Upload, Download, View, Edit, Key, SwitchButton, Delete,
}

const tableRef = ref<any>(null)
const userFormRef = ref<any>(null)
const detailRef = ref<any>(null)
const batchRef = ref<any>(null)

const roleOptions = ref<{ label: string; value: number }[]>([])
const roleList = ref<RoleSimple[]>([])

async function loadRoles() {
  try {
    const res = await getAllRolesApi()
    roleList.value = res.data
    roleOptions.value = res.data.map((r) => ({ value: r.id, label: r.name }))
  } catch { /* */ }
}

async function handleToggleStatus(row: UserInfo) {
  const newStatus = row.status === 1 ? 2 : 1
  try {
    await updateUserStatusApi(row.id, newStatus)
    ElMessage.success(newStatus === 1 ? '已启用' : '已禁用')
    tableRef.value?.refresh()
    invalidateCache('UserList')
  } catch { /* */ }
}

async function handleDelete(row: UserInfo) {
  try {
    await deleteUserApi(row.id)
    ElMessage.success('删除成功')
    tableRef.value?.refresh()
    invalidateCache('UserList')
  } catch { /* */ }
}

const tableConfig = {
  fetchApi: (params: Record<string, any>) => getUserListApi(params),
  columns: [
    { prop: 'user', label: '用户', minWidth: 200, slot: 'user' },
    { prop: 'email', label: '邮箱', minWidth: 180 },
    { prop: 'role', label: '角色', width: 120, slot: 'role' },
    { prop: 'status', label: '状态', width: 90, slot: 'status' },
    { prop: 'created_at', label: '创建时间', width: 180 },
  ],
  searchFields: [
    { prop: 'username', label: '用户名', type: 'input' as const, placeholder: '模糊搜索', width: '180px' },
    { prop: 'role_id', label: '角色', type: 'select' as const, options: roleOptions, width: '140px' },
    { prop: 'status', label: '状态', type: 'select' as const, options: [{ label: '正常', value: 1 }, { label: '禁用', value: 2 }], width: '120px' },
  ],
  toolbarActions: [
    { label: '新增用户', icon: 'Plus', type: 'primary' as const, permission: 'user:create',
      onClick: () => userFormRef.value?.open('create') },
    { label: '导出', icon: 'Download', permission: 'user:read' },
    { label: '批量分配角色', icon: 'Key', type: 'warning' as const, permission: 'user:update',
      onBatch: (selected: UserInfo[]) => batchRef.value?.open(selected.map((s: UserInfo) => s.id)) },
  ],
  rowActions: [
    { label: '详情', icon: 'View', type: 'primary' as const, permission: 'user:read',
      onClick: (row: UserInfo) => detailRef.value?.open(row.id) },
    { label: '编辑', icon: 'Edit', type: 'primary' as const, permission: 'user:update',
      onClick: (row: UserInfo) => userFormRef.value?.open('edit', row.id) },
    { label: (r: any) => (r as UserInfo).status === 1 ? '禁用' : '启用', icon: 'SwitchButton',
      type: (r: any) => (r as UserInfo).status === 1 ? 'warning' : 'success',
      confirm: (r: any) => (r as UserInfo).status === 1 ? '确认禁用该用户？' : '确认启用该用户？',
      onClick: (row: UserInfo) => handleToggleStatus(row) },
    { label: '删除', icon: 'Delete', type: 'danger' as const, permission: 'user:delete',
      confirm: '确认删除该用户？', onClick: (row: UserInfo) => handleDelete(row) },
  ],
  showSelection: true,
  cacheKey: 'UserList',
  actionColumnWidth: '280px',
}

const userFormConfig = {
  title: (mode: string) => mode === 'create' ? '新增用户' : '编辑用户',
  width: '560px',
  fields: [
    { prop: 'username', label: '用户名', type: 'input' as const,
      rules: [{ required: true, min: 3, max: 64, message: '3-64个字符', trigger: 'blur' as const }],
      disabledOnEdit: true },
    { prop: 'password', label: '密码', type: 'password' as const,
      rules: [{ min: 6, max: 128, message: '6-128个字符', trigger: 'blur' as const }],
      placeholder: '新建时必填，编辑时留空不修改' },
    { prop: 'nickname', label: '昵称', type: 'input' as const, placeholder: '选填' },
    { prop: 'email', label: '邮箱', type: 'input' as const,
      rules: [{ type: 'email' as const, message: '邮箱格式不正确', trigger: 'blur' as const }] },
    { prop: 'phone', label: '手机号', type: 'input' as const, placeholder: '选填' },
    { prop: 'role_id', label: '角色', type: 'select' as const, options: roleOptions,
      rules: [{ required: true, message: '请选择角色', trigger: 'change' as const }] },
  ],
  fetchDetail: (id: number) => getUserApi(id).then((r) => r.data),
  onSubmit: async (payload: any, id?: number) => {
    if (id) {
      await updateUserApi(id, payload)
    } else {
      await createUserApi(payload)
    }
  },
  onSuccess: () => {
    tableRef.value?.refresh()
    invalidateCache('UserList')
  },
}

const detailConfig = {
  title: '用户详情',
  fetchDetail: (id: number) => getUserApi(id).then((r) => r.data),
  width: '600px',
  header: {
    avatar: (u: any) => (u as UserInfo).avatar || '',
    title: (u: any) => (u as UserInfo).nickname || (u as UserInfo).username,
    subtitle: (u: any) => `@${(u as UserInfo).username}`,
    status: (u: any) => (u as UserInfo).status === 1 ? { label: '正常', type: 'success' as const } : { label: '禁用', type: 'danger' as const },
  },
  sections: [
    {
      items: [
        { label: '邮箱', value: (u: any) => (u as UserInfo).email || '-' },
        { label: '手机号', value: (u: any) => (u as UserInfo).phone || '-' },
        { label: '角色', value: '', slot: 'role' },
        { label: '状态', value: '', slot: 'status' },
        { label: '创建时间', value: (u: any) => (u as UserInfo).created_at },
        { label: '更新时间', value: (u: any) => (u as UserInfo).updated_at },
      ],
    },
  ],
}

onMounted(loadRoles)
</script>

<style scoped lang="scss">
.page-header {
  margin-bottom: 16px;
  h2 { font-size: var(--font-size-xl); font-weight: 700; color: var(--color-text); }
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-info {
  display: flex;
  flex-direction: column;
}

.user-name {
  font-weight: 500;
  color: var(--color-text);
}

.user-desc {
  font-size: 12px;
  color: var(--color-text-muted);
}
</style>
