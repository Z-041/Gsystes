<template>
  <el-tag :type="tagType" size="small" class="role-tag">
    {{ name || '未分配' }}
  </el-tag>
</template>

<script setup lang="ts">
import { computed } from 'vue'

defineOptions({ name: 'RoleTag' })

const props = defineProps<{
  name?: string
}>()

const tagType = computed(() => {
  if (!props.name) return 'info'
  const map: Record<string, string> = {
    管理员: 'danger',
    admin: 'danger',
  }
  return map[props.name] || 'primary'
})
</script>
