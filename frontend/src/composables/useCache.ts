import { ref, onActivated } from 'vue'

const cacheVersions = new Map<string, number>()

export function useCache() {
  const _version = ref(0)

  function setup(key: string, refreshFn: () => Promise<void> | void) {
    if (!cacheVersions.has(key)) {
      cacheVersions.set(key, 0)
    }
    _version.value = cacheVersions.get(key)!

    onActivated(async () => {
      const current = cacheVersions.get(key)!
      if (current !== _version.value) {
        _version.value = current
        await refreshFn()
      }
    })
  }

  return { setup }
}

export function invalidateCache(key: string) {
  const current = cacheVersions.get(key) || 0
  cacheVersions.set(key, current + 1)
}
