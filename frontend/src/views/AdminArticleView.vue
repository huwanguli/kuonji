<template>
  <div class="admin-page container-wide">
    <div class="page-header">
      <router-link to="/admin/articles" class="back-link">&larr; 返回文章管理</router-link>
      <router-link :to="`/editor/${articleId}`" class="btn-edit" v-if="article">编辑</router-link>
    </div>

    <div v-if="loading" class="center">加载中...</div>
    <div v-else-if="!article" class="center empty">文章不存在。</div>
    <template v-else>
      <div class="article-preview">
        <div class="preview-meta">
          <span :class="['status-badge', article.status === 1 ? 'published' : article.status === 2 ? 'private' : 'draft']">
            {{ article.status === 1 ? '已发布' : article.status === 2 ? '私密' : '草稿' }}
          </span>
          <time>{{ formatDate(article.created_at) }}</time>
          <span>{{ article.view_count }} 阅读</span>
        </div>
        <h1 class="preview-title">{{ article.title }}</h1>
        <div class="preview-content markdown-body" v-html="article.content_html"></div>
      </div>

      <section class="comments-section">
        <div class="cp-toolbar">
          <h3 class="cp-heading">评论管理 ({{ commentTotal }})</h3>
          <label class="cp-select-all" v-if="comments.length">
            <input type="checkbox" :checked="allSelected" @change="toggleAll" />
            全选
          </label>
          <div class="cp-batch-actions" v-if="selectedIds.length">
            <button @click="batchUpdateStatus(1)" class="cp-btn cp-approve">通过 ({{ selectedIds.length }})</button>
            <button @click="batchUpdateStatus(2)" class="cp-btn cp-delete">删除 ({{ selectedIds.length }})</button>
          </div>
        </div>

        <div v-if="commentLoading" class="cp-loading">加载评论中...</div>
        <div v-else-if="!comments.length" class="cp-empty">暂无评论</div>
        <div v-else class="cp-list">
          <div v-for="c in comments" :key="c.id" class="cp-row">
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
      </section>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { articles, comments as commentsApi } from '../api'
import { addHeadingIds } from '../utils/html'

const route = useRoute()
const articleId = computed(() => Number(route.params.id))
const article = ref(null)
const loading = ref(true)

const comments = ref([])
const commentTotal = ref(0)
const commentLoading = ref(false)
const selectedIds = ref([])

const allSelected = computed(() => {
  return comments.value.length > 0 && selectedIds.value.length === comments.value.length
})

onMounted(async () => {
  try {
    const res = await articles.adminDetail(articleId.value)
    if (res.code === 200) {
      article.value = res.data
      article.value.content_html = addHeadingIds(article.value.content_html)
    }
  } catch {} finally {
    loading.value = false
  }
  fetchComments()
})

async function fetchComments() {
  commentLoading.value = true
  try {
    const res = await commentsApi.adminList({ article_id: articleId.value, page_size: 100 })
    if (res.code === 200 && res.data) {
      comments.value = res.data.list || []
      commentTotal.value = res.data.total || comments.value.length
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
    selectedIds.value = comments.value.map(c => c.id)
  }
}

async function batchUpdateStatus(status) {
  if (!selectedIds.value.length) return
  for (const id of selectedIds.value) {
    try { await commentsApi.updateStatus(id, status) } catch {}
  }
  selectedIds.value = []
  fetchComments()
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
</script>

<style scoped>
.admin-page {
  padding-bottom: var(--space-16);
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-6);
}

.back-link {
  font-size: var(--text-sm);
  color: var(--color-muted);
  text-decoration: none;
  transition: color var(--duration) var(--ease);
}

.back-link:hover { color: var(--color-ink); }

.btn-edit {
  font-size: var(--text-sm);
  font-weight: 600;
  color: white;
  background: var(--color-vermilion);
  text-decoration: none;
  padding: var(--space-2) var(--space-4);
  border-radius: var(--radius);
  transition: opacity var(--duration) var(--ease);
}

.btn-edit:hover { opacity: 0.85; }

.center { text-align: center; padding: var(--space-16) 0; color: var(--color-muted); }
.empty { font-style: italic; }

.preview-meta {
  display: flex;
  align-items: center;
  gap: var(--space-4);
  font-size: var(--text-sm);
  color: var(--color-muted);
  margin-bottom: var(--space-4);
}

.status-badge {
  font-size: var(--text-xs);
  font-weight: 600;
  padding: 1px 8px;
  border-radius: 3px;
}

.status-badge.published { background: #d1fae5; color: #065f46; }
.status-badge.private { background: #fee2e2; color: #991b1b; }
.status-badge.draft { background: #f1f5f9; color: #64748b; }

.preview-title {
  font-family: var(--font-display);
  font-size: var(--text-3xl);
  color: var(--color-deep);
  margin-bottom: var(--space-6);
  line-height: 1.3;
}

.preview-content {
  font-size: var(--text-base);
  line-height: 1.85;
  color: var(--color-ink);
  margin-bottom: var(--space-10);
  background: transparent;
}

.preview-content :deep(h1),
.preview-content :deep(h2),
.preview-content :deep(h3),
.preview-content :deep(h4),
.preview-content :deep(h5),
.preview-content :deep(h6) {
  font-family: var(--font-display);
  scroll-margin-top: 100px;
}
.preview-content :deep(del) { color: var(--color-muted); }

/* Comments section */
.comments-section {
  border: 1px solid var(--color-border);
  border-radius: var(--radius);
  margin-top: var(--space-10);
}

.cp-toolbar {
  display: flex;
  align-items: center;
  gap: var(--space-4);
  padding: var(--space-3) var(--space-4);
  border-bottom: 1px solid var(--color-border);
  font-size: var(--text-xs);
}

.cp-heading {
  font-family: var(--font-body);
  font-size: var(--text-base);
  font-weight: 600;
  color: var(--color-deep);
  margin: 0;
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
  justify-content: flex-end;
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

.cp-loading,
.cp-empty {
  text-align: center;
  padding: var(--space-8);
  font-size: var(--text-sm);
  color: var(--color-muted);
}

.cp-list {
  max-height: 500px;
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
</style>
