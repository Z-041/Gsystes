import { ref, reactive } from 'vue'
import type { PageParams } from '@/types/api'

export function usePagination(fetchFn: (params: PageParams) => Promise<void>) {
  const total = ref(0)
  const loading = ref(false)
  const pageParam = reactive<PageParams>({
    page: 1,
    page_size: 20,
  })

  async function loadData() {
    loading.value = true
    try {
      await fetchFn({ ...pageParam })
    } finally {
      loading.value = false
    }
  }

  function handlePageChange(params: { page: number; page_size: number }) {
    pageParam.page = params.page
    pageParam.page_size = params.page_size
    loadData()
  }

  function resetPage() {
    pageParam.page = 1
    loadData()
  }

  return {
    total,
    loading,
    pageParam,
    loadData,
    handlePageChange,
    resetPage,
  }
}
