<template>
  <el-tag :type="status === 1 ? 'success' : 'danger'" size="small">
    {{ status === 1 ? '正常' : '禁用' }}
  </el-tag>
</template>

<script setup lang="ts">
defineOptions({ name: 'StatusTag' })

defineProps<{
  status: number
}>()
</script>
