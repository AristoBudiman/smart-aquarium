<script setup lang="ts">
defineProps<{
  eyebrow?: string
  title: string
  subtitle?: string
}>()
</script>

<template>
  <header class="page-header-wrap">
    <div v-if="$slots.breadcrumb" class="breadcrumb">
      <slot name="breadcrumb" />
    </div>
    <div class="page-heading" :class="{ 'heading-with-action': Boolean($slots.actions) }">
      <div>
        <p v-if="eyebrow" class="eyebrow">{{ eyebrow }}</p>
        <h1>{{ title }}</h1>
        <p v-if="subtitle" class="muted">{{ subtitle }}</p>
      </div>
      <div v-if="$slots.actions" class="page-heading-actions">
        <slot name="actions" />
      </div>
    </div>
  </header>
</template>
