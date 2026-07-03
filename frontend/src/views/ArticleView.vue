<template>
  <div class="page-layout container-wide">
    <Sidebar :content-html="article?.content_html || ''" />

    <main class="page-main">
      <div v-if="article" class="article-content-wrapper">
        <nav class="breadcrumb">
          <router-link to="/" class="back-link">首页</router-link>
          <template v-if="article.series">
            <span class="breadcrumb-sep">/</span>
            <router-link :to="`/series/${article.series}`" class="back-link">{{ article.series }}</router-link>
          </template>
        </nav>

        <div class="article-cover" v-if="article.cover">
          <img :src="article.cover" :alt="article.title" />
        </div>

        <header class="article-header">
          <div class="article-meta">
            <router-link
              v-if="article.category"
              :to="`/?category_id=${article.category.id}`"
              class="article-category"
            >{{ article.category.name }}</router-link>
            <time :datetime="article.created_at">{{ formatDate(article.created_at) }}</time>
            <span class="article-views">{{ article.view_count }} 次阅读</span>
          </div>
          <h1 class="article-title">{{ article.title }}</h1>
          <div class="article-tags" v-if="article.tags && article.tags.length">
            <router-link
              v-for="tag in article.tags"
              :key="tag.id"
              :to="`/?tag_id=${tag.id}`"
              class="chip"
            >{{ tag.name }}</router-link>
          </div>
        </header>

        <div class="article-content markdown-body" v-html="article.content_html"></div>

        <div v-if="article.series && (prevArticle || nextArticle || seriesArticles.length)" class="series-nav-wrap">
          <div class="series-nav-header">
            <router-link :to="`/series/${article.series}`" class="series-nav-name">{{ article.series }}</router-link>
            <span class="series-nav-pos" v-if="seriesTotal">(第 {{ article.series_order || '—' }} 篇 / 共 {{ seriesTotal }} 篇)</span>
          </div>

          <nav class="series-nav">
            <router-link v-if="prevArticle" :to="`/article/${prevArticle.slug}`" class="series-link series-prev">
              <span class="series-label">&larr; 上一篇</span>
              <span class="series-link-title">{{ prevArticle.title }}</span>
            </router-link>
            <span v-else class="series-link series-void"></span>
            <router-link v-if="nextArticle" :to="`/article/${nextArticle.slug}`" class="series-link series-next">
              <span class="series-label">下一篇 &rarr;</span>
              <span class="series-link-title">{{ nextArticle.title }}</span>
            </router-link>
            <span v-else class="series-link series-void"></span>
          </nav>

          <button v-if="seriesArticles.length > 1" class="series-toc-toggle" @click="seriesExpanded = !seriesExpanded">
            目录 {{ seriesExpanded ? '▾' : '▸' }}
          </button>
          <div v-if="seriesExpanded && seriesArticles.length" class="series-toc">
            <router-link
              v-for="(sa, i) in seriesArticles"
              :key="sa.slug"
              :to="`/article/${sa.slug}`"
              :class="['series-toc-item', { active: sa.slug === article.slug }]"
            >
              <span class="series-toc-order">{{ i + 1 }}</span>
              <span class="series-toc-title">{{ sa.title }}</span>
            </router-link>
          </div>
        </div>

        <CommentForm :article-id="article.id" @posted="fetchComments" />
        <CommentList :comments="comments" :total="commentTotal" :article-id="article.id" @posted="fetchComments" />
      </div>

      <div v-else-if="loading" class="loading">加载中...</div>
      <div v-else class="not-found">文章不存在。</div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { articles, comments as commentsApi } from '../api'
import { addHeadingIds } from '../utils/html'
import CommentList from '../components/CommentList.vue'
import CommentForm from '../components/CommentForm.vue'
import Sidebar from '../components/Sidebar.vue'

const route = useRoute()
const article = ref(null)
const prevArticle = ref(null)
const nextArticle = ref(null)
const seriesTotal = ref(0)
const seriesArticles = ref([])
const seriesExpanded = ref(false)
const loading = ref(true)
const comments = ref([])
const commentTotal = ref(0)

async function fetchArticle() {
  loading.value = true
  try {
    const res = await articles.getBySlug(route.params.slug)
    if (res.code === 200) {
      article.value = res.data.article
      article.value.content_html = addHeadingIds(article.value.content_html)
      prevArticle.value = res.data.prev_in_series || null
      nextArticle.value = res.data.next_in_series || null
      seriesTotal.value = res.data.series_total || 0
      seriesArticles.value = res.data.series_articles || []
      fetchComments()
    } else {
      article.value = null
    }
  } catch {
    article.value = null
  } finally {
    loading.value = false
  }
}

async function fetchComments() {
  if (!article.value) return
  try {
    const res = await commentsApi.list(article.value.id)
    if (res.code === 200 && res.data) {
      comments.value = res.data.list || []
      commentTotal.value = res.data.total || 0
    }
  } catch {}
}

function formatDate(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return `${d.getFullYear()}年${d.getMonth() + 1}月${d.getDate()}日`
}

onMounted(fetchArticle)
watch(() => route.params.slug, fetchArticle)
</script>

<style scoped>
.page-layout {
  display: grid;
  grid-template-columns: 200px 1fr;
  gap: var(--space-10);
  padding-top: var(--space-8);
  padding-bottom: var(--space-20);
}

@media (max-width: 900px) {
  .page-layout {
    grid-template-columns: 1fr;
  }
}

.loading,
.not-found {
  text-align: center;
  padding: var(--space-20) 0;
  color: var(--color-muted);
  font-size: var(--text-lg);
}

/* Breadcrumb */
.breadcrumb {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) 0 var(--space-6);
}

.back-link {
  font-size: var(--text-base);
  color: var(--color-muted);
  text-decoration: none;
  transition: color var(--duration) var(--ease);
  padding: var(--space-1) var(--space-2);
  border-radius: var(--radius-sm);
}

.back-link:hover {
  color: var(--color-vermilion);
  background: var(--color-card);
}

.breadcrumb-sep {
  font-size: var(--text-base);
  color: var(--color-border);
  margin: 0 var(--space-1);
}

/* Header */
.article-header {
  padding: var(--space-2) 0 var(--space-6);
}

.article-meta {
  display: flex;
  align-items: center;
  gap: var(--space-4);
  margin-bottom: var(--space-4);
  flex-wrap: wrap;
}

.article-category {
  font-size: var(--text-xs);
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: var(--color-vermilion);
  text-decoration: none;
  transition: opacity var(--duration) var(--ease);
}

.article-category:hover {
  opacity: 0.7;
}

time,
.article-views {
  font-size: var(--text-sm);
  color: var(--color-muted);
}

.article-title {
  font-family: var(--font-display);
  font-size: var(--text-5xl);
  font-weight: 700;
  color: var(--color-deep);
  line-height: 1.15;
  letter-spacing: -0.025em;
  margin-bottom: var(--space-4);
}

.article-tags {
  display: flex;
  gap: var(--space-2);
  flex-wrap: wrap;
}

.chip {
  font-size: var(--text-xs);
  color: var(--color-muted);
  background: var(--color-card);
  padding: 3px 12px;
  border-radius: 100px;
  font-weight: 500;
  text-decoration: none;
  transition: all var(--duration) var(--ease);
  border: 1px solid transparent;
}

.chip:hover {
  color: var(--color-ink);
  border-color: var(--color-border);
}

.article-content-wrapper {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius);
  padding: var(--space-8);
  box-shadow: var(--shadow-card);
}

/* Content */
.article-cover {
  margin-bottom: var(--space-8);
  border-radius: var(--radius);
  overflow: hidden;
}

.article-cover img {
  width: 100%;
  max-height: 480px;
  object-fit: cover;
  display: block;
  margin: 0;
}

.article-content {
  font-size: var(--text-lg);
  line-height: 1.85;
  color: var(--color-ink);
  padding: var(--space-8) 0;
  border-top: 1px solid var(--color-border);
  border-bottom: 1px solid var(--color-border);
  background: transparent;
}

.article-content :deep(h1),
.article-content :deep(h2),
.article-content :deep(h3),
.article-content :deep(h4),
.article-content :deep(h5),
.article-content :deep(h6) {
  font-family: var(--font-display);
  scroll-margin-top: 100px;
}

.article-content :deep(img) {
  display: block;
  margin-left: auto;
  margin-right: auto;
}

/* Series navigation */
.series-nav-wrap {
  margin-top: var(--space-8);
  padding-top: var(--space-6);
  border-top: 1px solid var(--color-border);
}

.series-nav-header {
  display: flex;
  align-items: baseline;
  gap: var(--space-3);
  margin-bottom: var(--space-4);
}

.series-nav-name {
  font-family: var(--font-display);
  font-size: var(--text-lg);
  font-weight: 600;
  color: var(--color-deep);
  text-decoration: none;
  transition: color var(--duration) var(--ease);
}

.series-nav-name:hover {
  color: var(--color-vermilion);
}

.series-nav-pos {
  font-size: var(--text-sm);
  color: var(--color-muted);
}

.series-nav {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-4);
  margin-bottom: var(--space-4);
}

.series-link {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
  padding: var(--space-4);
  border: 1px solid var(--color-border);
  border-radius: var(--radius);
  text-decoration: none;
  transition: all var(--duration) var(--ease);
}

.series-link:hover {
  border-color: var(--color-muted);
  box-shadow: 0 2px 8px rgba(0,0,0,0.04);
}

.series-next {
  text-align: right;
}

.series-label {
  font-size: var(--text-xs);
  color: var(--color-muted);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.series-link-title {
  font-size: var(--text-sm);
  color: var(--color-ink);
  font-weight: 500;
}

.series-void {
  visibility: hidden;
}

.series-toc-toggle {
  display: block;
  width: 100%;
  text-align: left;
  font-family: var(--font-body);
  font-size: var(--text-sm);
  color: var(--color-muted);
  background: none;
  border: none;
  padding: var(--space-2) 0;
  cursor: pointer;
  transition: color var(--duration) var(--ease);
}

.series-toc-toggle:hover {
  color: var(--color-ink);
}

.series-toc {
  border: 1px solid var(--color-border);
  border-radius: var(--radius);
  overflow: hidden;
  margin-bottom: var(--space-2);
}

.series-toc-item {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-2) var(--space-3);
  text-decoration: none;
  font-size: var(--text-sm);
  border-bottom: 1px solid var(--color-border);
  transition: background var(--duration) var(--ease);
}

.series-toc-item:last-child {
  border-bottom: none;
}

.series-toc-item:hover {
  background: var(--color-card);
}

.series-toc-item.active {
  background: var(--color-card);
}

.series-toc-item.active .series-toc-title {
  color: var(--color-vermilion);
  font-weight: 600;
}

.series-toc-order {
  font-family: var(--font-body);
  font-size: var(--text-xs);
  color: var(--color-muted);
  width: 20px;
  text-align: center;
  flex-shrink: 0;
}

.series-toc-title {
  color: var(--color-ink);
}

@media (max-width: 640px) {
  .article-title {
    font-size: var(--text-3xl);
  }

  .series-nav {
    grid-template-columns: 1fr;
  }

  .series-void {
    display: none;
  }
}
</style>
