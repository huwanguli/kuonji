<template>
  <div class="admin-page container-wide">

    <div class="page-header">
      <h2 class="page-title">文章管理</h2>
      <router-link to="/editor" class="btn-new">写新文章</router-link>
    </div>

    <div class="filter-row">
      <button :class="['filter-btn', { active: statusFilter === null }]" @click="setFilter(null)">全部</button>
      <button :class="['filter-btn', { active: statusFilter === 1 }]" @click="setFilter(1)">已发布</button>
      <button :class="['filter-btn', { active: statusFilter === 0 }]" @click="setFilter(0)">草稿</button>
    </div>

    <div v-if="loading" class="center">加载中...</div>
    <div v-else-if="!articles.length" class="center empty">暂无文章</div>
    <div v-else class="article-table">
      <div v-for="a in articles" :key="a.id" class="table-row">
        <div class="row-main">
          <span :class="['status-dot', a.status === 1 ? 'published' : 'draft']"></span>
          <router-link :to="`/editor/${a.id}`" class="row-title">{{ a.title }}</router-link>
          <span class="row-meta">{{ formatDate(a.created_at) }}</span>
          <span class="row-views" v-if="a.status === 1">{{ a.view_count }} 阅读</span>
        </div>
        <div class="row-actions">
          <router-link :to="`/editor/${a.id}`" class="row-action">编辑</router-link>
          <button @click="deleteArticle(a.id)" class="row-action danger">删除</button>
        </div>
      </div>

      <div class="pagination" v-if="total > pageSize">
        <button class="page-btn" :disabled="page <= 1" @click="goPage(page - 1)">&larr; 上页</button>
        <span class="page-info">{{ page }} / {{ totalPages }}</span>
        <button class="page-btn" :disabled="page >= totalPages" @click="goPage(page + 1)">下页 &rarr;</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { articles as api } from '../api'

const articles = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = 20
const loading = ref(true)
const statusFilter = ref(null)

const totalPages = computed(() => Math.ceil(total.value / pageSize))

async function fetch() {
  loading.value = true
  const params = { page: page.value, page_size: pageSize }
  if (statusFilter.value !== null) params.status = statusFilter.value
  try {
    const res = await api.adminList(params)
    if (res.code === 200) {
      articles.value = res.data.list || []
      total.value = res.data.total || 0
    }
  } catch {} finally {
    loading.value = false
  }
}

function setFilter(v) {
  statusFilter.value = v
  page.value = 1
  fetch()
}

function goPage(p) {
  page.value = p
  fetch()
}

async function deleteArticle(id) {
  if (!confirm('确定删除这篇文章？')) return
  try {
    const res = await api.delete(id)
    if (res.code === 200) fetch()
  } catch {}
}

function formatDate(s) {
  if (!s) return ''
  const d = new Date(s)
  return `${d.getFullYear()}/${d.getMonth() + 1}/${d.getDate()}`
}

onMounted(fetch)
</script>

<style scoped>
.admin-page {
  padding-bottom: var(--space-16);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-4);
}

.page-title {
  font-family: var(--font-display);
  font-size: var(--text-xl);
  color: var(--color-deep);
}

.btn-new {
  font-family: var(--font-body);
  font-size: var(--text-sm);
  font-weight: 600;
  color: white;
  background: var(--color-vermilion);
  text-decoration: none;
  padding: var(--space-2) var(--space-4);
  border-radius: var(--radius);
  transition: opacity var(--duration) var(--ease);
}

.btn-new:hover { opacity: 0.85; }

.filter-row {
  display: flex;
  gap: var(--space-2);
  margin-bottom: var(--space-4);
}

.filter-btn {
  font-family: var(--font-body);
  font-size: var(--text-sm);
  color: var(--color-muted);
  background: none;
  border: 1px solid var(--color-border);
  padding: var(--space-1) var(--space-4);
  border-radius: 100px;
  cursor: pointer;
  transition: all var(--duration) var(--ease);
}

.filter-btn:hover { color: var(--color-ink); }

.filter-btn.active {
  color: var(--color-paper);
  background: var(--color-ink);
  border-color: var(--color-ink);
}

.center { text-align: center; padding: var(--space-16) 0; color: var(--color-muted); }
.empty { font-style: italic; }

.article-table {
  border: 1px solid var(--color-border);
  border-radius: var(--radius);
  overflow: hidden;
}

.table-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-3) var(--space-4);
  border-bottom: 1px solid var(--color-border);
  gap: var(--space-4);
}

.table-row:last-child { border-bottom: none; }
.table-row:hover { background: var(--color-card); }

.row-main {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  flex: 1;
  min-width: 0;
  flex-wrap: wrap;
}

.status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.status-dot.published { background: #16a34a; }
.status-dot.draft { background: var(--color-muted); }

.row-title {
  font-size: var(--text-base);
  font-weight: 500;
  color: var(--color-ink);
  text-decoration: none;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.row-title:hover { color: var(--color-vermilion); }

.row-meta, .row-views {
  font-size: var(--text-xs);
  color: var(--color-muted);
  white-space: nowrap;
}

.row-actions {
  display: flex;
  gap: var(--space-4);
  flex-shrink: 0;
}

.row-action {
  font-size: var(--text-sm);
  color: var(--color-muted);
  background: none;
  border: none;
  cursor: pointer;
  font-family: var(--font-body);
  text-decoration: none;
  transition: color var(--duration) var(--ease);
}

.row-action:hover { color: var(--color-ink); }
.row-action.danger:hover { color: var(--color-vermilion); }

.pagination {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: var(--space-4);
  padding: var(--space-6) 0;
}

.page-btn {
  font-family: var(--font-body);
  font-size: var(--text-sm);
  color: var(--color-ink);
  background: none;
  border: 1px solid var(--color-border);
  padding: var(--space-2) var(--space-4);
  border-radius: var(--radius);
  cursor: pointer;
}

.page-btn:hover:not(:disabled) { border-color: var(--color-vermilion); }
.page-btn:disabled { opacity: 0.3; cursor: not-allowed; }
.page-info { font-size: var(--text-sm); color: var(--color-muted); }

@media (max-width: 640px) {
  .table-row { flex-direction: column; align-items: flex-start; gap: var(--space-2); }
  .row-actions { align-self: flex-end; }
}
</style>
