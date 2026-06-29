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
    </div>
  </nav>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

const isAdminPage = computed(() => {
  return route.path.startsWith('/admin') || route.path.startsWith('/editor')
})

function logout() {
  localStorage.removeItem('token')
  router.push('/admin')
}
</script>

<style scoped>
.navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
  background: rgba(253, 251, 247, 0.85);
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
</style>
