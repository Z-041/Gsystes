<template>
  <div class="table-pagination">
    <el-pagination
      v-model:current-page="innerPage"
      v-model:page-size="innerPageSize"
      :page-sizes="[10, 20, 50, 100]"
      :total="total"
      layout="total, sizes, prev, pager, next, jumper"
      background
      @size-change="onPageSizeChange"
      @current-change="onPageChange"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

defineOptions({ name: 'TablePagination' })

const props = withDefaults(defineProps<{
  page: number
  pageSize: number
  total: number
}>(), {
  page: 1,
  pageSize: 20,
  total: 0,
})

const emit = defineEmits<{
  change: [params: { page: number; page_size: number }]
}>()

const innerPage = ref(props.page)
const innerPageSize = ref(props.pageSize)

watch(() => props.page, (v) => { innerPage.value = v })
watch(() => props.pageSize, (v) => { innerPageSize.value = v })

function emitChange() {
  emit('change', { page: innerPage.value, page_size: innerPageSize.value })
}

function onPageChange() {
  emitChange()
}

function onPageSizeChange() {
  innerPage.value = 1
  emitChange()
}
</script>

<style scoped lang="scss">
.table-pagination {
  display: flex;
  justify-content: flex-end;
  padding: 16px 0 0;
}
</style>
