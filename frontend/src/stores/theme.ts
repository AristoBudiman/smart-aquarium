import { defineStore } from 'pinia'
import { computed, ref, watch } from 'vue'

export type ThemeMode = 'light' | 'dark' | 'system'

const STORAGE_KEY = 'aquasense_theme'

function readStoredMode(): ThemeMode {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (raw === 'light' || raw === 'dark' || raw === 'system') {
      return raw
    }
  } catch {
    // Gunakan mode Ikuti sistem jika localStorage tidak dapat diakses.
  }
  return 'system'
}

function getSystemDark(): boolean {
  try {
    return typeof window !== 'undefined' && window.matchMedia('(prefers-color-scheme: dark)').matches
  } catch {
    return false
  }
}

export const useThemeStore = defineStore('theme', () => {
  const mode = ref<ThemeMode>(readStoredMode())
  const systemDark = ref<boolean>(getSystemDark())
  let listenerBound = false

  const isDark = computed(() =>
    mode.value === 'dark' ? true : mode.value === 'light' ? false : systemDark.value,
  )

  function syncDom() {
    if (typeof document === 'undefined') return
    const dark = isDark.value
    const root = document.documentElement
    root.classList.toggle('dark', dark)
    root.style.colorScheme = dark ? 'dark' : 'light'
    const meta = document.querySelector('meta[name="theme-color"]')
    if (meta) {
      meta.setAttribute('content', dark ? '#0b1220' : '#0d9488')
    }
  }

  function init() {
    if (typeof window !== 'undefined' && !listenerBound) {
      listenerBound = true
      try {
        const mql = window.matchMedia('(prefers-color-scheme: dark)')
        systemDark.value = mql.matches
        const onChange = (event: MediaQueryListEvent) => {
          systemDark.value = event.matches
        }
        mql.addEventListener('change', onChange)
      } catch {
        // Abaikan jika matchMedia tidak didukung.
      }
    }
    syncDom()
  }

  function setMode(next: ThemeMode) {
    mode.value = next
    try {
      localStorage.setItem(STORAGE_KEY, next)
    } catch {
      // Abaikan kegagalan tulis ke localStorage.
    }
    syncDom()
  }

  watch(isDark, () => {
    syncDom()
  })

  return {
    mode,
    isDark,
    init,
    setMode,
  }
})
