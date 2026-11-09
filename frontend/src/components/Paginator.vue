<script setup lang="ts">
import { ChevronLeft, ChevronRight } from 'lucide-vue-next'

const props = defineProps<{
  page: number
  totalPages: number
  totalItems: number
}>()

const emit = defineEmits<{
  'update:page': [nextPage: number]
}>()

function prev() {
  if (props.page > 1) emit('update:page', props.page - 1)
}

function next() {
  if (props.page < props.totalPages) emit('update:page', props.page + 1)
}
</script>

<template>
  <nav v-if="totalPages > 1" class="paginator" aria-label="Navigasi halaman daftar">
    <button
      class="button button-quiet paginator-btn"
      type="button"
      :disabled="page <= 1"
      @click="prev"
    >
      <ChevronLeft :size="18" aria-hidden="true" />
      <span>Sebelumnya</span>
    </button>
    <span class="paginator-status">
      Halaman <strong>{{ page }}</strong> dari <strong>{{ totalPages }}</strong>
    </span>
    <button
      class="button button-quiet paginator-btn"
      type="button"
      :disabled="page >= totalPages"
      @click="next"
    >
      <span>Berikutnya</span>
      <ChevronRight :size="18" aria-hidden="true" />
    </button>
  </nav>
</template>
