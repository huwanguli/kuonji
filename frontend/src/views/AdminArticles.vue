<template>
  <div class="admin-page container-wide">

    <div class="page-header">
      <h2 class="page-title">文章管理</h2>
      <router-link to="/editor" class="btn-new">写新文章</router-link>
    </div>

    <div class="filter-row">
      <button :class="['filter-btn', { active: statusFilter === null }]" @click="setFilter(null)">全部</button>
      <button :class="['filter-btn', { active: statusFilter === 1 }]" @click="setFilter(1)">已发布</button>
      <button :class="['filter-btn', { active: statusFilter === 2 }]" @click="setFilter(2)">私密</button>
      <button :class="['filter-btn', { active: statusFilter === 0 }]" @click="setFilter(0)">草稿</button>
    </div>

    <div v-if="loading" class="center">加载中...</div>
    <div v-else-if="!articles.length" class="center empty">暂无文章</div>
    <div v-else class="article-table">
      <div v-for="a in articles" :key="a.id" class="table-group">
        <div class="table-row">
          <div class="row-main">
            <span :class="['status-dot', a.status === 1 ? 'published' : a.status === 2 ? 'private' : 'draft']"></span>
            <router-link :to="`/admin/articles/view/${a.id}`" class="row-title">{{ a.title }}</router-link>
            <span class="row-meta">{{ formatDate(a.created_at) }}</span>
            <span class="row-views" v-if="a.status !== 0">{{ a.view_count }} 阅读</span>
          </div>
          <div class="row-actions">
            <button @click="toggleComments(a)" class="row-action">
              评论 {{ commentCounts[a.id] != null ? '(' + commentCounts[a.id] + ')' : '' }}
            </button>
            <router-link :to="`/editor/${a.id}`" class="row-action">编辑</router-link>
            <button @click="deleteArticle(a.id)" class="row-action danger">删除</button>
          </div>
        </div>

        <div v-if="openPanel === a.id" class="comment-panel">
          <div class="cp-toolbar">
            <label class="cp-select-all">
              <input type="checkbox" :checked="allSelected" @change="toggleAll" />
              全选
            </label>
            <div class="cp-batch-actions" v-if="selectedIds.length">
              <button @click="batchUpdateStatus(1)" class="cp-btn cp-approve">通过 ({{ selectedIds.length }})</button>
              <button @click="batchUpdateStatus(2)" class="cp-btn cp-delete">删除 ({{ selectedIds.length }})</button>
            </div>
            <button @click="openPanel = null" class="cp-close">&times;</button>
          </div>

          <div v-if="commentLoading" class="cp-loading">加载评论中...</div>
          <div v-else-if="!panelComments.length" class="cp-empty">暂无评论</div>
          <div v-else class="cp-list">
            <div v-for="c in panelComments" :key="c.id" class="cp-row">
              <input type="checkbox" :checked="selectedIds.includes(c.id)" @change="toggleSelect(c.id)" class="cp-check" />
              <div class="cp-body">
                <div class="cp-meta">
                  <strong>{{ c.author }}</strong>
                  <span :class="['cp-status', 's-' + c.status]">{{ statusLabel(c.status) }}</span>
                  <time>{{ formatDateTime(c.created_at) }}</time>
                  <span v-if="c.parent_id" class="cp-parent">回复 #{{ c.parent_id }}</span>
                </div>
                <p class="cp-content">{{ c.content }}</p>
              </div>
            </div>
          </div>
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
import { ref, computed, onMounted, nextTick } from 'vue'
import { articles as api, comments as commentsApi } from '../api'

const articles = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = 20
const loading = ref(true)
const statusFilter = ref(null)

const openPanel = ref(null)
const panelComments = ref([])
const commentCounts = ref({})
const commentLoading = ref(false)
const selectedIds = ref([])

const allSelected = computed(() => {
  return panelComments.value.length > 0 && selectedIds.value.length === panelComments.value.length
})

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

async function toggleComments(a) {
  if (openPanel.value === a.id) {
    openPanel.value = null
    return
  }
  openPanel.value = a.id
  selectedIds.value = []
  panelComments.value = []
  commentLoading.value = true
  try {
    const res = await commentsApi.adminList({ article_id: a.id, page_size: 100 })
    if (res.code === 200 && res.data) {
      panelComments.value = res.data.list || []
      commentCounts.value[a.id] = res.data.total || panelComments.value.length
    }
  } catch {} finally {
    commentLoading.value = false
  }
}

function toggleSelect(id) {
  const idx = selectedIds.value.indexOf(id)
  if (idx >= 0) {
    selectedIds.value.splice(idx, 1)
  } else {
    selectedIds.value.push(id)
  }
}

function toggleAll() {
  if (allSelected.value) {
    selectedIds.value = []
  } else {
    selectedIds.value = panelComments.value.map(c => c.id)
  }
}

async function batchUpdateStatus(status) {
  if (!selectedIds.value.length) return
  for (const id of selectedIds.value) {
    try { await commentsApi.updateStatus(id, status) } catch {}
  }
  selectedIds.value = []
  const a = articles.value.find(x => x.id === openPanel.value)
  if (a) toggleComments(a)
}

function statusLabel(s) {
  if (s === 0) return '待审核'
  if (s === 2) return '已删除'
  return '正常'
}

function formatDate(s) {
  if (!s) return ''
  const d = new Date(s)
  return `${d.getFullYear()}/${d.getMonth() + 1}/${d.getDate()}`
}

function formatDateTime(s) {
  if (!s) return ''
  const d = new Date(s)
  return `${d.getMonth() + 1}/${d.getDate()} ${d.getHours()}:${String(d.getMinutes()).padStart(2, '0')}`
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

.table-group {
  border-bottom: 1px solid var(--color-border);
}

.table-group:last-child { border-bottom: none; }

.table-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-3) var(--space-4);
  gap: var(--space-4);
}

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
.status-dot.private { background: var(--color-vermilion); }
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

/* Comment panel */
.comment-panel {
  background: var(--color-card);
  border-top: 1px solid var(--color-border);
  padding: 0;
}

.cp-toolbar {
  display: flex;
  align-items: center;
  gap: var(--space-4);
  padding: var(--space-2) var(--space-4);
  border-bottom: 1px solid var(--color-border);
  font-size: var(--text-xs);
}

.cp-select-all {
  display: flex;
  align-items: center;
  gap: var(--space-1);
  color: var(--color-muted);
  cursor: pointer;
}

.cp-batch-actions {
  display: flex;
  gap: var(--space-2);
  flex: 1;
}

.cp-btn {
  font-family: var(--font-body);
  font-size: var(--text-xs);
  font-weight: 500;
  padding: 2px 10px;
  border-radius: var(--radius-sm);
  cursor: pointer;
  border: 1px solid transparent;
  transition: all var(--duration) var(--ease);
}

.cp-approve {
  color: #065f46;
  background: #d1fae5;
}

.cp-approve:hover { border-color: #065f46; }

.cp-delete {
  color: var(--color-vermilion);
  background: #fee2e2;
}

.cp-delete:hover { border-color: var(--color-vermilion); }

.cp-close {
  font-size: var(--text-lg);
  color: var(--color-muted);
  background: none;
  border: none;
  cursor: pointer;
  padding: 0 var(--space-1);
  margin-left: auto;
}

.cp-close:hover { color: var(--color-ink); }

.cp-loading,
.cp-empty {
  text-align: center;
  padding: var(--space-8);
  font-size: var(--text-sm);
  color: var(--color-muted);
}

.cp-list {
  max-height: 400px;
  overflow-y: auto;
}

.cp-row {
  display: flex;
  gap: var(--space-3);
  padding: var(--space-3) var(--space-4);
  border-bottom: 1px solid var(--color-border);
  align-items: flex-start;
}

.cp-row:last-child { border-bottom: none; }

.cp-check {
  margin-top: 4px;
  flex-shrink: 0;
  accent-color: var(--color-vermilion);
}

.cp-body {
  flex: 1;
  min-width: 0;
}

.cp-meta {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  flex-wrap: wrap;
  margin-bottom: var(--space-1);
  font-size: var(--text-xs);
  color: var(--color-muted);
}

.cp-meta strong {
  color: var(--color-deep);
  font-size: var(--text-sm);
}

.cp-status {
  font-weight: 600;
  padding: 0 6px;
  border-radius: 3px;
  font-size: 0.65rem;
}

.s-0 { background: #fef3c7; color: #92400e; }
.s-1 { background: #d1fae5; color: #065f46; }
.s-2 { background: #f1f5f9; color: #64748b; }

.cp-parent {
  color: var(--color-muted);
}

.cp-content {
  font-size: var(--text-sm);
  line-height: 1.4;
  color: var(--color-ink);
  word-break: break-word;
}

@media (max-width: 640px) {
  .table-row { flex-direction: column; align-items: flex-start; gap: var(--space-2); }
  .row-actions { align-self: flex-end; }
  .cp-toolbar { flex-wrap: wrap; gap: var(--space-2); }
}
</style>
