<template>
  <nav class="navbar" role="navigation" aria-label="主导航">
    <div class="navbar-inner">
      <div class="nav-left">
        <router-link to="/" class="logo" aria-label="Kuonji 首页">
          Kuonji
        </router-link>
        <span v-if="isAdminPage" class="logo-sep">/</span>
        <router-link v-if="isAdminPage" to="/admin" class="logo logo-admin">admin</router-link>
      </div>
      <div class="nav-right">
        <div class="nav-links">
        <template v-if="isAdminPage">
          <router-link to="/admin" class="nav-link" exact-active-class="active">面板</router-link>
          <router-link to="/admin/articles" class="nav-link" active-class="active">文章</router-link>
          <router-link to="/editor" class="nav-link" active-class="active">写文章</router-link>
          <button @click="logout" class="nav-link logout-link">退出</button>
        </template>
        <template v-else>
          <router-link to="/" class="nav-link" exact-active-class="active">首页</router-link>
          <a href="https://github.com/huwanguli" target="_blank" rel="noopener" class="nav-link gh-link">GitHub</a>
        </template>
      </div>
      <button class="theme-toggle" @click="toggle" :aria-label="isDark() ? '切换亮色模式' : '切换暗色模式'">
        <svg v-if="isDark()" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="5"/><line x1="12" y1="1" x2="12" y2="3"/><line x1="12" y1="21" x2="12" y2="23"/><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/><line x1="1" y1="12" x2="3" y2="12"/><line x1="21" y1="12" x2="23" y2="12"/><line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/></svg>
        <svg v-else width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/></svg>
      </button>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useTheme } from '../composables/useTheme'

const { toggle, isDark } = useTheme()

const route = useRoute()

const isAdminPage = computed(() => {
  return route.path.startsWith('/admin') || route.path.startsWith('/editor')
})

function logout() {
  localStorage.removeItem('token')
  window.location.href = '/admin'
}
</script>

<style scoped>
.navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
  background: var(--color-navbar-bg);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border-bottom: 1px solid var(--color-border);
}

.navbar-inner {
  max-width: var(--max-width-wide);
  margin: 0 auto;
  padding: 0 var(--space-6);
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.nav-left {
  display: flex;
  align-items: baseline;
  gap: var(--space-2);
}

.logo {
  font-family: var(--font-display);
  font-size: var(--text-2xl);
  font-weight: 700;
  color: var(--color-ink);
  text-decoration: none;
  letter-spacing: -0.01em;
  transition: color var(--duration) var(--ease);
}

.logo:hover {
  color: var(--color-vermilion);
}

.logo-sep {
  color: var(--color-border);
  font-size: var(--text-xl);
  font-weight: 300;
}

.logo-admin {
  font-size: var(--text-sm);
  font-weight: 400;
  color: var(--color-muted);
  text-transform: uppercase;
  letter-spacing: 0.06em;
  font-family: var(--font-body);
}

.logo-admin:hover {
  color: var(--color-ink);
}

.nav-links {
  display: flex;
  gap: var(--space-6);
  align-items: center;
}

.nav-right {
  display: flex;
  align-items: center;
  gap: var(--space-4);
}

.nav-link {
  font-size: var(--text-sm);
  font-weight: 500;
  color: var(--color-muted);
  text-decoration: none;
  letter-spacing: 0.02em;
  padding: var(--space-1) 0;
  border-bottom: 2px solid transparent;
  transition: all var(--duration) var(--ease);
  background: none;
  border-top: none;
  border-left: none;
  border-right: none;
  cursor: pointer;
  font-family: var(--font-body);
}

.nav-link:hover,
.nav-link.active {
  color: var(--color-ink);
  border-bottom-color: var(--color-vermilion);
}

.logout-link:hover {
  color: var(--color-vermilion) !important;
  border-bottom-color: var(--color-vermilion);
}

.theme-toggle {
  background: none;
  border: 1px solid var(--color-border);
  border-radius: 50%;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: var(--color-muted);
  padding: 0;
  margin-left: var(--space-2);
  transition: color var(--duration) var(--ease), border-color var(--duration) var(--ease);
}

.theme-toggle:hover {
  color: var(--color-ink);
  border-color: var(--color-ink);
}
</style>
