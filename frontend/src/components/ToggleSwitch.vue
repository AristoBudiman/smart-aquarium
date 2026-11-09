<script setup lang="ts">
defineProps<{
  modelValue: boolean
  disabled?: boolean
  label: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  toggle: []
}>()

function onClick() {
  emit('update:modelValue', true)
  emit('toggle')
}
</script>

<template>
  <button
    class="toggle-switch"
    :class="{ 'toggle-switch-on': modelValue }"
    role="switch"
    :aria-checked="modelValue"
    :aria-label="label"
    :disabled="disabled"
    type="button"
    @click="onClick"
  >
    <span class="toggle-switch-thumb"></span>
  </button>
</template>
