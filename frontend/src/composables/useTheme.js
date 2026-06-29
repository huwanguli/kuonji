import { ref, watchEffect } from 'vue'

const STORAGE_KEY = 'zb-theme'
const theme = ref(localStorage.getItem(STORAGE_KEY) || 'light')

watchEffect(() => {
  document.documentElement.setAttribute('data-theme', theme.value)
  localStorage.setItem(STORAGE_KEY, theme.value)
})

export function useTheme() {
  const toggle = () => {
    theme.value = theme.value === 'light' ? 'dark' : 'light'
  }

  const isDark = () => theme.value === 'dark'

  return { theme, toggle, isDark }
}
