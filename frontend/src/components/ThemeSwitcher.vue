<script setup lang="ts">
import { Monitor, Moon, Sun } from 'lucide-vue-next'
import { useThemeStore, type ThemeMode } from '../stores/theme'

const theme = useThemeStore()

const options: { value: ThemeMode; label: string; icon: typeof Sun }[] = [
  { value: 'light', label: 'Terang', icon: Sun },
  { value: 'dark', label: 'Gelap', icon: Moon },
  { value: 'system', label: 'Ikuti sistem', icon: Monitor },
]
</script>

<template>
  <div class="theme-switcher" role="radiogroup" aria-label="Pilih tema tampilan">
    <button
      v-for="item in options"
      :key="item.value"
      class="theme-switcher-btn"
      :class="{ 'is-active': theme.mode === item.value }"
      type="button"
      role="radio"
      :aria-checked="theme.mode === item.value"
      :aria-label="`Tema ${item.label}`"
      :title="`Tema: ${item.label}`"
      @click="theme.setMode(item.value)"
    >
      <component :is="item.icon" :size="16" :stroke-width="2.1" aria-hidden="true" />
      <span class="theme-switcher-label">{{ item.label }}</span>
    </button>
  </div>
</template>
