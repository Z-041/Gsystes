<template>
  <div class="stat-card" :style="{ background: bg }">
    <div class="stat-icon">
      <slot name="icon" />
    </div>
    <div class="stat-info">
      <div class="stat-value">{{ value }}</div>
      <div class="stat-label">{{ label }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineOptions({ name: 'StatCard' })

withDefaults(defineProps<{
  label: string
  value: string | number
  bg?: string
}>(), {
  bg: 'linear-gradient(135deg, var(--color-primary), var(--color-primary-light))',
})
</script>

<style scoped lang="scss">
.stat-card {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 18px 20px;
  border-radius: var(--radius-md);
  color: #fff;
}
.stat-icon {
  width: 44px;
  height: 44px;
  border-radius: var(--radius-sm);
  background: rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.stat-info {
  flex: 1;
}
.stat-value {
  font-size: 24px;
  font-weight: 700;
  line-height: 1.2;
}
.stat-label {
  font-size: 12px;
  opacity: 0.85;
  margin-top: 2px;
}
</style>
