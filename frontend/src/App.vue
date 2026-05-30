<template>
  <router-view />
</template>

<script setup lang="ts">
import { useWebSocket } from '@/composables/useWebSocket'

const { registerGlobal } = useWebSocket()
registerGlobal()
</script>
