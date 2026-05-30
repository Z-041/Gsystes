<template>
  <el-dialog
    v-model="visible"
    title="分配权限"
    width="520px"
    :close-on-click-modal="false"
    destroy-on-close
  >
    <div v-loading="loading">
      <el-tree
        v-if="treeData.length > 0"
        ref="treeRef"
        :data="treeData"
        :props="treeProps"
        node-key="id"
        show-checkbox
        default-expand-all
        check-strictly
      />
      <el-empty v-else-if="!loading" description="暂无权限数据" />
    </div>
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" :loading="submitting" @click="handleSave">保存</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import type { ElTree } from 'element-plus'
import { getAllPermissionsApi } from '@/api/permission'
import { getRolePermissionsApi, assignPermissionsApi } from '@/api/role'
import type { PermissionInfo } from '@/types/permission'

defineOptions({ name: 'AssignPermissionDialog' })

const emit = defineEmits<{
  success: []
}>()

const visible = ref(false)
const loading = ref(false)
const submitting = ref(false)
const treeRef = ref<InstanceType<typeof ElTree>>()
const treeData = ref<PermissionInfo[]>([])
const roleId = ref(0)

const treeProps = { children: 'children', label: 'name' }

function buildTree(list: PermissionInfo[]): PermissionInfo[] {
  const map = new Map<number, PermissionInfo>()
  const roots: PermissionInfo[] = []

  list.forEach((item) => {
    map.set(item.id, { ...item, children: [] })
  })

  list.forEach((item) => {
    const node = map.get(item.id)!
    if (item.parent_id && map.has(item.parent_id)) {
      map.get(item.parent_id)!.children!.push(node)
    } else {
      roots.push(node)
    }
  })

  return roots
}

function flattenIds(nodes: PermissionInfo[]): number[] {
  const ids: number[] = []
  for (const n of nodes) {
    ids.push(n.id)
    if (n.children) ids.push(...flattenIds(n.children))
  }
  return ids
}

async function open(id: number) {
  roleId.value = id
  visible.value = true
  loading.value = true
  treeData.value = []

  try {
    const [allRes, roleRes] = await Promise.all([
      getAllPermissionsApi(),
      getRolePermissionsApi(id),
    ])
    treeData.value = buildTree(allRes.data)
    await nextTick()
    const checkedIds = flattenIds(roleRes.data)
    if (checkedIds.length > 0 && treeRef.value) {
      treeRef.value.setCheckedKeys(checkedIds, false)
    }
  } catch {
    visible.value = false
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  if (!treeRef.value) return
  const checkedKeys = treeRef.value.getCheckedKeys(false) as number[]

  submitting.value = true
  try {
    await assignPermissionsApi(roleId.value, { permission_ids: checkedKeys })
    ElMessage.success('权限分配成功')
    visible.value = false
    emit('success')
  } finally {
    submitting.value = false
  }
}

defineExpose({ open })
</script>
