<template>
  <div class="admin-comments container-wide">
    <div class="page-header">
      <h2 class="page-title">评论管理</h2>
      <router-link to="/admin" class="back-link">&larr; 返回面板</router-link>
    </div>

    <div v-if="!loggedIn" class="not-logged">
      <p>请先 <router-link to="/admin">登录</router-link></p>
    </div>

    <div v-else>
      <div class="filter-bar">
        <label class="filter-label">
          <select v-model="statusFilter" @change="fetchList(1)" class="select">
            <option :value="null">全部</option>
            <option :value="0">待审核</option>
            <option :value="1">已通过</option>
            <option :value="2">已删除</option>
          </select>
        </label>
      </div>

      <div v-if="loading" class="loading">加载中...</div>

      <div v-else-if="!comments.length" class="empty">暂无评论。</div>

      <div v-else>
        <div class="comment-table">
          <div class="comment-row" v-for="c in comments" :key="c.id">
            <div class="row-main">
              <div class="row-meta">
                <strong>{{ c.author }}</strong>
                <span class="row-article">文章 #{{ c.article_id }}</span>
                <time>{{ formatDate(c.created_at) }}</time>
                <span :class="['status-badge', 'status-' + c.status]">{{ statusLabel(c.status) }}</span>
              </div>
              <p class="row-content">{{ c.content }}</p>
              <p v-if="c.parent_id" class="row-parent">回复 #{{ c.parent_id }}</p>
            </div>
            <div class="row-actions">
              <button
                v-if="c.status !== 1"
                class="action-btn action-approve"
                @click="updateStatus(c.id, 1)"
              >通过</button>
              <button
                v-if="c.status !== 2"
                class="action-btn action-delete"
                @click="confirmDelete(c)"
              >删除</button>
              <button
                v-if="c.status === 2"
                class="action-btn action-approve"
                @click="updateStatus(c.id, 1)"
              >恢复</button>
            </div>
          </div>
        </div>

        <div class="pagination" v-if="totalPages > 1">
          <button class="page-btn" :disabled="page <= 1" @click="goPage(page - 1)">&larr; 上一页</button>
          <span class="page-info">第 {{ page }} / {{ totalPages }} 页</span>
          <button class="page-btn" :disabled="page >= totalPages" @click="goPage(page + 1)">下一页 &rarr;</button>
        </div>
      </div>
    </div>

    <div v-if="deleteTarget" class="modal-overlay" @click.self="deleteTarget = null">
      <div class="modal">
        <p>确认删除这条评论？</p>
        <p class="delete-preview">"{{ deleteTarget.content }}"</p>
        <p class="modal-note">子评论不受影响，"回复 @{{ deleteTarget.author }}" 照常显示</p>
        <div class="modal-actions">
          <button class="btn-cancel" @click="deleteTarget = null">取消</button>
          <button class="btn-delete" @click="doDelete">确认删除</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { comments as commentsApi } from '../api'

const loggedIn = ref(!!localStorage.getItem('token'))
const comments = ref([])
const page = ref(1)
const totalPages = ref(1)
const loading = ref(true)
const statusFilter = ref(null)
const deleteTarget = ref(null)

async function fetchList(p = 1) {
  loading.value = true
  try {
    const params = { page: p, page_size: 20 }
    if (statusFilter.value !== null) params.status = statusFilter.value
    const res = await commentsApi.adminList(params)
    if (res.code === 200 && res.data) {
      comments.value = res.data.list || []
      page.value = res.data.page || 1
      totalPages.value = Math.ceil((res.data.total || 0) / 20) || 1
    }
  } catch {} finally {
    loading.value = false
  }
}

function goPage(p) {
  fetchList(p)
}

async function updateStatus(id, status) {
  try {
    await commentsApi.updateStatus(id, status)
    fetchList(page.value)
  } catch {}
}

function confirmDelete(c) {
  deleteTarget.value = c
}

async function doDelete() {
  if (!deleteTarget.value) return
  try {
    await commentsApi.delete(deleteTarget.value.id)
    deleteTarget.value = null
    fetchList(page.value)
  } catch {}
}

function statusLabel(s) {
  if (s === 0) return '待审核'
  if (s === 2) return '已删除'
  return '已通过'
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString('zh-CN', {
    year: 'numeric', month: '2-digit', day: '2-digit',
    hour: '2-digit', minute: '2-digit',
  })
}

onMounted(() => {
  if (loggedIn.value) fetchList()
})
</script>

<style scoped>
.admin-comments {
  padding: var(--space-8) 0 var(--space-16);
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--space-6);
}

.page-title {
  font-family: var(--font-display);
  font-size: var(--text-2xl);
  color: var(--color-deep);
}

.back-link {
  font-size: var(--text-sm);
  color: var(--color-muted);
  text-decoration: none;
}

.back-link:hover {
  color: var(--color-vermilion);
}

.not-logged {
  text-align: center;
  padding: var(--space-12) 0;
  color: var(--color-muted);
}

.not-logged a {
  color: var(--color-vermilion);
}

.loading,
.empty {
  text-align: center;
  padding: var(--space-12) 0;
  color: var(--color-muted);
  font-size: var(--text-sm);
}

.filter-bar {
  margin-bottom: var(--space-4);
}

.filter-label {
  font-size: var(--text-sm);
  color: var(--color-muted);
}

.select {
  font-family: var(--font-body);
  font-size: var(--text-sm);
  padding: var(--space-2) var(--space-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  background: var(--color-paper);
  color: var(--color-ink);
  outline: none;
  cursor: pointer;
}

.select:focus {
  border-color: var(--color-vermilion);
}

.comment-table {
  border: 1px solid var(--color-border);
  border-radius: var(--radius);
  overflow: hidden;
}

.comment-row {
  display: flex;
  align-items: flex-start;
  gap: var(--space-4);
  padding: var(--space-4);
  border-bottom: 1px solid var(--color-border);
  background: var(--color-paper);
}

.comment-row:last-child {
  border-bottom: none;
}

.row-main {
  flex: 1;
  min-width: 0;
}

.row-meta {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  flex-wrap: wrap;
  margin-bottom: var(--space-2);
  font-size: var(--text-xs);
  color: var(--color-muted);
}

.row-meta strong {
  color: var(--color-deep);
  font-size: var(--text-sm);
}

.row-article {
  color: var(--color-muted);
}

.status-badge {
  font-weight: 600;
  padding: 0 6px;
  border-radius: 3px;
  font-size: 0.7rem;
}

.status-0 {
  background: #fef3c7;
  color: #92400e;
}

.status-1 {
  background: #d1fae5;
  color: #065f46;
}

.status-2 {
  background: #f1f5f9;
  color: #64748b;
}

.row-content {
  font-size: var(--text-sm);
  line-height: 1.5;
  color: var(--color-ink);
  word-break: break-word;
}

.row-parent {
  font-size: var(--text-xs);
  color: var(--color-muted);
  margin-top: var(--space-1);
}

.row-actions {
  display: flex;
  gap: var(--space-2);
  flex-shrink: 0;
  padding-top: 2px;
}

.action-btn {
  font-family: var(--font-body);
  font-size: var(--text-xs);
  font-weight: 500;
  padding: var(--space-1) var(--space-3);
  border-radius: var(--radius-sm);
  cursor: pointer;
  border: 1px solid transparent;
  transition: all var(--duration) var(--ease);
}

.action-approve {
  color: #065f46;
  background: #d1fae5;
}

.action-approve:hover {
  border-color: #065f46;
}

.action-delete {
  color: var(--color-vermilion);
  background: #fee2e2;
}

.action-delete:hover {
  border-color: var(--color-vermilion);
}

.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-4);
  padding: var(--space-8) 0;
}

.page-btn {
  font-family: var(--font-body);
  font-size: var(--text-sm);
  font-weight: 500;
  color: var(--color-ink);
  background: none;
  border: 1px solid var(--color-border);
  padding: var(--space-2) var(--space-4);
  border-radius: var(--radius);
  cursor: pointer;
  transition: all var(--duration) var(--ease);
}

.page-btn:hover:not(:disabled) {
  border-color: var(--color-vermilion);
  color: var(--color-vermilion);
}

.page-btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.page-info {
  font-size: var(--text-sm);
  color: var(--color-muted);
}

.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 100;
}

.modal {
  background: var(--color-paper);
  border-radius: var(--radius);
  padding: var(--space-6);
  max-width: 400px;
  width: 90%;
  box-shadow: 0 8px 32px rgba(0,0,0,0.12);
}

.modal p {
  margin-bottom: var(--space-4);
  font-size: var(--text-base);
}

.delete-preview {
  color: var(--color-muted);
  font-style: italic;
  font-size: var(--text-sm);
}

.modal-note {
  font-size: var(--text-xs);
  color: var(--color-muted);
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-3);
}

.btn-cancel {
  font-family: var(--font-body);
  font-size: var(--text-sm);
  color: var(--color-muted);
  background: none;
  border: 1px solid var(--color-border);
  padding: var(--space-2) var(--space-4);
  border-radius: var(--radius-sm);
  cursor: pointer;
}

.btn-cancel:hover {
  border-color: var(--color-muted);
}

.btn-delete {
  font-family: var(--font-body);
  font-size: var(--text-sm);
  font-weight: 600;
  color: #fff;
  background: var(--color-vermilion);
  border: none;
  padding: var(--space-2) var(--space-4);
  border-radius: var(--radius-sm);
  cursor: pointer;
}

.btn-delete:hover {
  opacity: 0.85;
}
</style>
