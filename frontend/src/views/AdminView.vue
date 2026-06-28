<template>
  <div class="admin-page container-wide">

    <div v-if="!isLoggedIn">
      <h2 class="admin-title">管理员登录</h2>
      <form @submit.prevent="login" class="login-form">
        <input v-model="username" type="text" placeholder="用户名" required class="input" />
        <input v-model="password" type="password" placeholder="密码" required class="input" />
        <button type="submit" class="login-btn" :disabled="submitting">
          {{ submitting ? '登录中...' : '登 录' }}
        </button>
        <p v-if="error" class="error">{{ error }}</p>
      </form>
    </div>

    <div v-else class="dashboard">
      <h2 class="admin-title">管理面板</h2>

      <div class="quick-links">
        <router-link to="/editor" class="quick-card">
          <span class="quick-icon">&#9998;</span>
          <span class="quick-label">写新文章</span>
        </router-link>
        <router-link to="/admin/articles" class="quick-card">
          <span class="quick-icon">&#9776;</span>
          <span class="quick-label">文章管理</span>
          <span class="quick-stat">{{ stats.articles }} 篇</span>
        </router-link>
      </div>

      <div class="admin-sections">
        <section class="admin-section">
          <h3>分类</h3>
          <ul class="taxo-list" v-if="adminCategories.length">
            <li v-for="c in adminCategories" :key="c.id" class="taxo-item">
              <span>{{ c.name }}</span>
              <button @click="deleteCategory(c.id)" class="del-btn" title="删除">&times;</button>
            </li>
          </ul>
          <div class="inline-form">
            <input v-model="newCat.name" type="text" placeholder="新分类名称" class="input-sm" @keyup.enter="addCategory" />
            <button @click="addCategory" class="add-btn">添加</button>
          </div>
        </section>

        <section class="admin-section">
          <h3>标签</h3>
          <div class="taxo-chips" v-if="adminTags.length">
            <span v-for="t in adminTags" :key="t.id" class="admin-chip">
              {{ t.name }}
              <button @click="deleteTag(t.id)" class="del-chip">&times;</button>
            </span>
          </div>
          <div class="inline-form">
            <input v-model="newTag.name" type="text" placeholder="新标签名称" class="input-sm" @keyup.enter="addTag" />
            <button @click="addTag" class="add-btn">添加</button>
          </div>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { auth, categories as catApi, tags as tagApi, articles as articleApi } from '../api'

const username = ref('')
const password = ref('')
const submitting = ref(false)
const error = ref('')
const isLoggedIn = ref(false)

const adminCategories = ref([])
const adminTags = ref([])
const newCat = reactive({ name: '' })
const newTag = reactive({ name: '' })
const stats = reactive({ articles: 0 })

async function login() {
  submitting.value = true
  error.value = ''
  try {
    const res = await auth.login(username.value, password.value)
    if (res.code === 200) {
      localStorage.setItem('token', res.data.token)
      isLoggedIn.value = true
      loadDashboard()
    } else {
      error.value = res.message || '登录失败'
    }
  } catch {
    error.value = '登录失败，请重试。'
  } finally {
    submitting.value = false
  }
}

function logout() {
  localStorage.removeItem('token')
  isLoggedIn.value = false
}

async function loadDashboard() {
  try {
    const [catRes, tagRes, articleRes] = await Promise.all([
      catApi.list(),
      tagApi.list(),
      articleApi.adminList({ page: 1, page_size: 1 }),
    ])
    if (catRes.code === 200) adminCategories.value = catRes.data || []
    if (tagRes.code === 200) adminTags.value = tagRes.data || []
    if (articleRes.code === 200) stats.articles = articleRes.data.total || 0
  } catch {}
}

async function addCategory() {
  if (!newCat.name.trim()) return
  try {
    const res = await catApi.create({ name: newCat.name.trim() })
    if (res.code === 200) {
      newCat.name = ''
      loadDashboard()
    }
  } catch {}
}

async function deleteCategory(id) {
  try {
    const res = await catApi.delete(id)
    if (res.code === 200) loadDashboard()
  } catch {}
}

async function addTag() {
  if (!newTag.name.trim()) return
  try {
    const res = await tagApi.create({ name: newTag.name.trim() })
    if (res.code === 200) {
      newTag.name = ''
      loadDashboard()
    }
  } catch {}
}

async function deleteTag(id) {
  try {
    const res = await tagApi.delete(id)
    if (res.code === 200) loadDashboard()
  } catch {}
}

onMounted(() => {
  isLoggedIn.value = !!localStorage.getItem('token')
  if (isLoggedIn.value) loadDashboard()
})
</script>

<style scoped>
.admin-page {
  padding-bottom: var(--space-16);
}

.admin-title {
  font-family: var(--font-display);
  font-size: var(--text-2xl);
  color: var(--color-deep);
  margin-bottom: var(--space-6);
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  max-width: 400px;
}

.input {
  font-family: var(--font-body);
  font-size: var(--text-base);
  padding: var(--space-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  background: var(--color-paper);
  color: var(--color-ink);
}

.input:focus {
  outline: none;
  border-color: var(--color-vermilion);
}

.login-btn {
  font-family: var(--font-body);
  font-size: var(--text-sm);
  font-weight: 600;
  color: white;
  background: var(--color-vermilion);
  border: none;
  padding: var(--space-3);
  border-radius: var(--radius);
  cursor: pointer;
  transition: opacity var(--duration) var(--ease);
  letter-spacing: 0.5em;
}

.login-btn:hover:not(:disabled) { opacity: 0.85; }
.login-btn:disabled { opacity: 0.5; cursor: not-allowed; }
.error { font-size: var(--text-sm); color: var(--color-vermilion); }

/* Dashboard */
.quick-links {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-4);
  max-width: 500px;
  margin-bottom: var(--space-8);
}

@media (max-width: 400px) {
  .quick-links { grid-template-columns: 1fr; }
}

.quick-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-6);
  border: 1px solid var(--color-border);
  border-radius: var(--radius);
  text-decoration: none;
  transition: all var(--duration) var(--ease);
  color: var(--color-ink);
}

.quick-card:hover {
  border-color: var(--color-vermilion);
  box-shadow: 0 2px 8px rgba(0,0,0,0.04);
}

.quick-icon {
  font-size: var(--text-3xl);
  line-height: 1;
}

.quick-label {
  font-weight: 600;
  font-size: var(--text-sm);
}

.quick-stat {
  font-size: var(--text-xs);
  color: var(--color-muted);
}

/* Taxonomy sections */
.admin-sections {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-6);
  max-width: 600px;
}

@media (max-width: 500px) {
  .admin-sections { grid-template-columns: 1fr; }
}

.admin-section h3 {
  font-family: var(--font-display);
  font-size: var(--text-lg);
  color: var(--color-deep);
  margin-bottom: var(--space-3);
}

.taxo-list { list-style: none; margin-bottom: var(--space-3); }

.taxo-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--space-1) 0;
  font-size: var(--text-sm);
  color: var(--color-ink);
}

.del-btn {
  font-size: var(--text-lg);
  color: var(--color-muted);
  background: none;
  border: none;
  cursor: pointer;
  padding: 0 var(--space-1);
  line-height: 1;
  transition: color var(--duration) var(--ease);
}

.del-btn:hover { color: var(--color-vermilion); }

.taxo-chips { display: flex; gap: var(--space-2); flex-wrap: wrap; margin-bottom: var(--space-3); }

.admin-chip {
  font-size: var(--text-xs);
  color: var(--color-ink);
  background: var(--color-card);
  padding: 2px 8px;
  border-radius: 100px;
  font-weight: 500;
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.del-chip {
  font-size: 12px;
  color: var(--color-muted);
  background: none;
  border: none;
  cursor: pointer;
  padding: 0;
  line-height: 1;
}

.del-chip:hover { color: var(--color-vermilion); }

.inline-form { display: flex; gap: var(--space-2); }

.input-sm {
  font-family: var(--font-body);
  font-size: var(--text-sm);
  padding: var(--space-1) var(--space-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  background: var(--color-paper);
  color: var(--color-ink);
  flex: 1;
}

.input-sm:focus { outline: none; border-color: var(--color-vermilion); }

.add-btn {
  font-family: var(--font-body);
  font-size: var(--text-xs);
  font-weight: 600;
  color: white;
  background: var(--color-vermilion);
  border: none;
  padding: var(--space-1) var(--space-4);
  border-radius: var(--radius-sm);
  cursor: pointer;
  white-space: nowrap;
}

.add-btn:hover { opacity: 0.85; }
</style>
