import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface ToastItem {
  id: number
  kind: 'success' | 'error' | 'info'
  message: string
}

export const useToastStore = defineStore('toast', () => {
  const items = ref<ToastItem[]>([])
  let nextId = 1

  function show(message: string, kind: ToastItem['kind'] = 'success', durationMs = 3400) {
    const id = nextId
    nextId += 1
    items.value.push({ id, kind, message })
    window.setTimeout(() => {
      dismiss(id)
    }, durationMs)
  }

  function success(message: string) {
    show(message, 'success')
  }

  function error(message: string) {
    show(message, 'error', 4200)
  }

  function dismiss(id: number) {
    items.value = items.value.filter((item) => item.id !== id)
  }

  return { items, show, success, error, dismiss }
})
