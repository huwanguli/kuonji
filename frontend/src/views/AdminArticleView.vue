<template>
  <div class="admin-page container-wide">
    <div class="page-header">
      <router-link to="/admin/articles" class="back-link">&larr; 返回文章管理</router-link>
      <router-link :to="`/editor/${articleId}`" class="btn-edit" v-if="article">编辑</router-link>
    </div>

    <div v-if="loading" class="center">加载中...</div>
    <div v-else-if="!article" class="center empty">文章不存在。</div>
    <div v-else class="article-preview">
      <div class="preview-meta">
        <span :class="['status-badge', article.status === 1 ? 'published' : article.status === 2 ? 'private' : 'draft']">
          {{ article.status === 1 ? '已发布' : article.status === 2 ? '私密' : '草稿' }}
        </span>
        <time>{{ formatDate(article.created_at) }}</time>
        <span>{{ article.view_count }} 阅读</span>
      </div>
      <h1 class="preview-title">{{ article.title }}</h1>
      <div class="preview-content" v-html="article.content_html"></div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { articles } from '../api'
import { addHeadingIds } from '../utils/html'

const route = useRoute()
const articleId = computed(() => Number(route.params.id))
const article = ref(null)
const loading = ref(true)

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
})

function formatDate(s) {
  if (!s) return ''
  const d = new Date(s)
  return `${d.getFullYear()}/${d.getMonth() + 1}/${d.getDate()}`
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
}

.preview-content :deep(h1) { font-family: var(--font-display); font-size: var(--text-3xl); margin-bottom: var(--space-4); margin-top: var(--space-8); }
.preview-content :deep(h2) { font-family: var(--font-display); font-size: var(--text-2xl); margin-bottom: var(--space-3); margin-top: var(--space-6); }
.preview-content :deep(h3) { font-family: var(--font-display); font-size: var(--text-xl); margin-bottom: var(--space-3); margin-top: var(--space-6); }
.preview-content :deep(p) { margin-bottom: var(--space-4); }
.preview-content :deep(pre) { background: var(--color-deep); color: #e2e8f0; padding: var(--space-4); border-radius: var(--radius); overflow-x: auto; margin-bottom: var(--space-4); }
.preview-content :deep(code) { font-family: var(--font-mono); font-size: 0.9em; background: var(--color-card); padding: 0.1em 0.3em; border-radius: 2px; }
.preview-content :deep(pre code) { background: none; padding: 0; }
.preview-content :deep(img) { max-width: 100%; border-radius: var(--radius); }
.preview-content :deep(blockquote) { border-left: 3px solid var(--color-vermilion); padding-left: var(--space-4); color: var(--color-muted); margin: var(--space-4) 0; }
</style>
