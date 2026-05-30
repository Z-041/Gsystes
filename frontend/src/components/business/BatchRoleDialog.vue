<template>
  <el-dialog
    v-model="visible"
    title="批量分配角色"
    width="420px"
    :close-on-click-modal="false"
    destroy-on-close
  >
    <el-form label-width="80px">
      <el-form-item label="目标角色">
        <el-select v-model="roleId" placeholder="请选择角色" style="width: 100%">
          <el-option
            v-for="r in roles"
            :key="r.id"
            :label="r.name"
            :value="r.id"
          />
        </el-select>
      </el-form-item>
      <el-form-item label="选中用户">
        <span class="count-text">{{ userIds.length }}</span>
        <span class="count-suffix">人</span>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="visible = false">取消</el-button>
      <el-button type="primary" :loading="submitting" @click="handleSubmit">确认分配</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import type { RoleSimple } from '@/types/role'
import { batchAssignRoleApi } from '@/api/user'

defineOptions({ name: 'BatchRoleDialog' })

defineProps<{
  roles: RoleSimple[]
}>()

const emit = defineEmits<{
  success: []
}>()

const visible = ref(false)
const submitting = ref(false)
const roleId = ref<number>()
const userIds = ref<number[]>([])

function open(ids: number[]) {
  userIds.value = ids
  roleId.value = undefined
  visible.value = true
}

async function handleSubmit() {
  if (!roleId.value) {
    ElMessage.warning('请选择角色')
    return
  }

  submitting.value = true
  try {
    await batchAssignRoleApi({
      user_ids: userIds.value,
      role_id: roleId.value,
    })
    ElMessage.success('分配成功')
    visible.value = false
    emit('success')
  } finally {
    submitting.value = false
  }
}

defineExpose({ open })
</script>

<style scoped lang="scss">
.count-text {
  color: var(--color-primary);
  font-weight: 600;
  font-size: 16px;
}
.count-suffix {
  color: var(--color-text-muted);
  margin-left: 4px;
  font-size: 13px;
}
</style>
