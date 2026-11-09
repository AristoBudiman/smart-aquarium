<script setup lang="ts">
import type { Component } from 'vue'
import { Waves } from 'lucide-vue-next'

defineProps<{
  title: string
  description: string
  icon?: Component
}>()
</script>

<template>
  <div class="empty-state">
    <div class="empty-icon-wrap" aria-hidden="true">
      <component :is="icon ?? Waves" :size="30" :stroke-width="2" />
    </div>
    <h2>{{ title }}</h2>
    <p>{{ description }}</p>
    <div v-if="$slots.action" class="empty-state-action">
      <slot name="action" />
    </div>
  </div>
</template>
