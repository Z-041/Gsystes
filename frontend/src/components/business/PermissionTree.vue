<template>
  <div class="permission-tree">
    <div class="tree-toolbar">
      <el-button type="primary" :icon="Plus" v-permission="'perm:manage'" @click="openDialog()">
        新增权限
      </el-button>
    </div>

    <el-tree
      ref="treeRef"
      :data="treeData"
      :props="treeProps"
      node-key="id"
      default-expand-all
      highlight-current
      :expand-on-click-node="false"
    >
      <template #default="{ data }">
        <span class="tree-node">
          <span class="tree-node-label">
            <el-icon :size="16"><Folder v-if="data.type === 1" /><Document v-else /></el-icon>
            <span>{{ data.name }}</span>
            <el-tag v-if="data.type === 1" size="small" type="info" class="type-tag">菜单</el-tag>
            <el-tag v-else size="small" type="warning" class="type-tag">按钮</el-tag>
          </span>
          <span class="tree-node-actions">
            <el-button link type="primary" :icon="Plus" size="small" @click.stop="openDialog(data.id)" v-permission="'perm:manage'">子级</el-button>
            <el-button link type="primary" :icon="Edit" size="small" @click.stop="openDialog(data.id, data)" v-permission="'perm:manage'">编辑</el-button>
            <el-popconfirm title="确认删除？" @confirm="handleDelete(data.id)">
              <template #reference>
                <el-button link type="danger" size="small" v-permission="'perm:manage'">删除</el-button>
              </template>
            </el-popconfirm>
          </span>
        </span>
      </template>
    </el-tree>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus, Edit, Folder, Document } from '@element-plus/icons-vue'
import { getAllPermissionsApi, deletePermissionApi } from '@/api/permission'
import type { PermissionInfo } from '@/types/permission'

defineOptions({ name: 'PermissionTree' })

const emit = defineEmits<{
  'edit': [parentId?: number, data?: PermissionInfo]
  'delete': []
}>()

const treeData = ref<PermissionInfo[]>([])
const treeRef = ref()

const treeProps = {
  children: 'children',
  label: 'name',
}

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

async function loadTree() {
  try {
    const res = await getAllPermissionsApi()
    treeData.value = buildTree(res.data)
  } catch {
    // handled by interceptor
  }
}

function openDialog(parentId?: number, data?: PermissionInfo) {
  emit('edit', parentId, data)
}

async function handleDelete(id: number) {
  try {
    await deletePermissionApi(id)
    ElMessage.success('删除成功')
    emit('delete')
    loadTree()
  } catch {
    // handled by interceptor
  }
}

onMounted(loadTree)

defineExpose({ loadTree })
</script>

<style scoped lang="scss">
.permission-tree {
  .tree-toolbar {
    margin-bottom: 12px;
  }
}

.tree-node {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-right: 12px;
}

.tree-node-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  .type-tag {
    margin-left: 6px;
  }
}

.tree-node-actions {
  display: flex;
  gap: 2px;
  opacity: 0.5;
  transition: opacity 0.2s;
}

:deep(.el-tree-node__content:hover) .tree-node-actions {
  opacity: 1;
}
</style>
