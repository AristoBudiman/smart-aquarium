<script setup lang="ts">
import { onMounted } from 'vue'
import { RouterView } from 'vue-router'
import AppToast from './components/AppToast.vue'
import { useThemeStore } from './stores/theme'

const theme = useThemeStore()

onMounted(() => {
  theme.init()
})
</script>

<template>
  <RouterView />
  <AppToast />
</template>
