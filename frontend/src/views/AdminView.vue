<template>
  <div class="admin-page container-wide">

    <div v-if="!isLoggedIn">
      <h2 class="admin-title">管理员登录</h2>
      <p v-if="expiredMsg" class="expired-hint">{{ expiredMsg }}</p>
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
        <router-link to="/admin/comments" class="quick-card">
          <span class="quick-icon">&#9993;</span>
          <span class="quick-label">评论管理</span>
          <span class="quick-stat">{{ stats.comments }} 条</span>
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

        <section class="admin-section admin-section-wide">
          <h3>系列</h3>
          <ul class="taxo-list" v-if="adminSeries.length">
            <li v-for="s in adminSeries" :key="s.id">
              <div class="taxo-item">
                <span>{{ s.name }}</span>
                <span class="taxo-count">{{ s.count }} 篇</span>
                <div class="taxo-actions">
                  <button @click="startEdit(s)" class="mini-btn">编辑</button>
                  <button @click="deleteSeries(s.id)" class="del-btn">&times;</button>
                </div>
              </div>
              <div v-if="editingSeries === s.id" class="series-edit-row">
                <input v-model="editSeries.name" type="text" placeholder="系列名称" class="input-sm" />
                <div class="cover-input-wrap">
                  <input v-model="editSeries.cover" type="text" placeholder="封面图 URL" class="input-sm" />
                  <button @click="clickEditCover" class="cover-upload-btn" title="上传封面">
                    <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="8.5" cy="8.5" r="1.5"/><polyline points="21 15 16 10 5 21"/></svg>
                  </button>
                </div>
                <input v-model="editSeries.desc" type="text" placeholder="简介" class="input-sm" />
                <button @click="saveSeriesEdit(s)" class="add-btn">保存</button>
              </div>
            </li>
          </ul>
          <p class="series-hint" v-else>暂无系列。</p>
          <div class="inline-form-series">
            <input v-model="newSeries.name" type="text" placeholder="系列名称" class="input-sm" @keyup.enter="addSeries" />
            <div class="cover-input-wrap">
              <input v-model="newSeries.cover" type="text" placeholder="封面图 URL（选填）" class="input-sm" />
              <button @click="uploadSeriesCover" class="cover-upload-btn" title="上传封面">
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="8.5" cy="8.5" r="1.5"/><polyline points="21 15 16 10 5 21"/></svg>
              </button>
            </div>
            <input
              ref="seriesCoverInput"
              type="file"
              accept=".jpg,.jpeg,.png,.gif,.webp"
              class="file-hidden"
              @change="handleSeriesCoverUpload"
            />
            <input v-model="newSeries.desc" type="text" placeholder="简介（选填）" class="input-sm" />
            <button @click="addSeries" class="add-btn">新建</button>
          </div>
          <input
            ref="editSeriesCoverInput"
            type="file"
            accept=".jpg,.jpeg,.png,.gif,.webp"
            class="file-hidden"
            @change="handleEditSeriesCoverUpload"
          />
        </section>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { auth, categories as catApi, tags as tagApi, series as seriesApi, upload as uploadApi, articles as articleApi, comments as commentsApi } from '../api'

const route = useRoute()
const username = ref('')
const password = ref('')
const submitting = ref(false)
const error = ref('')
const expiredMsg = ref(route.query.expired === '1' ? '登录已过期，请重新登录' : '')
const isLoggedIn = ref(false)

const adminCategories = ref([])
const adminTags = ref([])
const adminSeries = ref([])
const newCat = reactive({ name: '' })
const newTag = reactive({ name: '' })
const newSeries = reactive({ name: '', cover: '', desc: '' })
const editingSeries = ref(null)
const editSeries = reactive({ name: '', cover: '', desc: '' })
const seriesCoverInput = ref(null)
const editSeriesCoverInput = ref(null)
const stats = reactive({ articles: 0, comments: 0 })

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
    const [catRes, tagRes, seriesRes, articleRes, commentRes] = await Promise.all([
      catApi.list(),
      tagApi.list(),
      seriesApi.list(),
      articleApi.adminList({ page: 1, page_size: 1 }),
      commentsApi.adminList({ page: 1, page_size: 1 }),
    ])
    if (catRes.code === 200) adminCategories.value = catRes.data || []
    if (tagRes.code === 200) adminTags.value = tagRes.data || []
    if (seriesRes.code === 200) adminSeries.value = seriesRes.data || []
    if (articleRes.code === 200) stats.articles = articleRes.data.total || 0
    if (commentRes.code === 200) stats.comments = commentRes.data.total || 0
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

async function deleteSeries(id) {
  if (!confirm('确定删除这个系列？相关文章的系列归属将被清空，但文章不受影响。')) return
  try {
    const res = await seriesApi.delete(id)
    if (res.code === 200) loadDashboard()
  } catch { alert('删除失败，请重试') }
}

function startEdit(s) {
  editingSeries.value = s.id
  editSeries.name = s.name
  editSeries.cover = s.cover || ''
  editSeries.desc = s.description || ''
}

async function saveSeriesEdit(s) {
  if (!editSeries.name.trim()) return
  try {
    const res = await seriesApi.update(s.id, {
      name: editSeries.name.trim(),
      cover: editSeries.cover.trim(),
      description: editSeries.desc.trim(),
    })
    if (res.code === 200) {
      editingSeries.value = null
      loadDashboard()
    }
  } catch { alert('保存失败，请重试') }
}

async function addSeries() {
  if (!newSeries.name.trim()) return
  try {
    const res = await seriesApi.create({ name: newSeries.name.trim(), cover: newSeries.cover.trim(), description: newSeries.desc.trim() })
    if (res.code === 200) {
      newSeries.name = ''
      newSeries.cover = ''
      newSeries.desc = ''
      loadDashboard()
    }
  } catch {}
}

function uploadSeriesCover() {
  seriesCoverInput.value?.click()
}

async function handleSeriesCoverUpload(e) {
  const file = e.target.files?.[0]
  if (!file) return
  try {
    const res = await uploadApi.image(file)
    if (res.code === 200 && res.data) {
      newSeries.cover = res.data.url
    }
  } catch {}
  e.target.value = ''
}

function clickEditCover() {
  editSeriesCoverInput.value?.click()
}

async function handleEditSeriesCoverUpload(e) {
  const file = e.target.files?.[0]
  if (!file) return
  try {
    const res = await uploadApi.image(file)
    if (res.code === 200 && res.data) {
      editSeries.cover = res.data.url
    }
  } catch {}
  e.target.value = ''
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
.expired-hint { font-size: var(--text-sm); color: #b45309; margin-bottom: var(--space-4); padding: var(--space-2) var(--space-3); background: #fef3c7; border-radius: var(--radius-sm); }

/* Dashboard */
.quick-links {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--space-4);
  max-width: 600px;
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
  width: 100%;
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

.taxo-count {
  font-size: var(--text-xs);
  color: var(--color-muted);
}

.taxo-actions {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.mini-btn {
  font-size: var(--text-xs);
  color: var(--color-muted);
  background: none;
  border: 1px solid var(--color-border);
  padding: 0 var(--space-2);
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-family: var(--font-body);
  transition: all var(--duration) var(--ease);
}

.mini-btn:hover {
  color: var(--color-ink);
  border-color: var(--color-ink);
}

.series-hint {
  font-size: var(--text-sm);
  color: var(--color-muted);
}

.admin-section-wide {
  grid-column: 1 / -1;
}

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

.series-edit-row {
  display: flex;
  gap: var(--space-2);
  padding: var(--space-2) 0 var(--space-3);
  flex-wrap: wrap;
  align-items: center;
}

.series-edit-row .input-sm:first-child {
  flex: 0 0 140px;
}

.series-edit-row .input-sm {
  flex: 1;
  min-width: 100px;
}

.cover-input-wrap {
  display: flex;
  gap: 2px;
  flex: 1;
  min-width: 120px;
}

.cover-input-wrap .input-sm {
  flex: 1;
  min-width: 0;
}

.cover-upload-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 var(--space-2);
  background: var(--color-card);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  color: var(--color-muted);
  flex-shrink: 0;
}

.cover-upload-btn:hover {
  color: var(--color-vermilion);
  border-color: var(--color-vermilion);
}

.file-hidden { display: none; }

.inline-form-series {
  display: flex;
  gap: var(--space-2);
  flex-wrap: wrap;
}

.inline-form-series .input-sm:first-child {
  flex: 0 0 140px;
}

.inline-form-series .input-sm {
  flex: 1;
  min-width: 120px;
}
</style>
