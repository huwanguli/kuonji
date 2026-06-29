<template>
  <div class="page-layout container-wide">
    <Sidebar :content-html="article?.content_html || ''" />

    <main class="page-main">
      <div v-if="article" class="article-content-wrapper">
        <nav class="breadcrumb">
          <router-link to="/" class="back-link">&larr; 返回首页</router-link>
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

        <div class="article-content" v-html="article.content_html"></div>

        <nav v-if="prevArticle || nextArticle" class="series-nav">
          <router-link v-if="prevArticle" :to="`/article/${prevArticle.slug}`" class="series-link series-prev">
            <span class="series-label">&larr; 上一篇</span>
            <span class="series-title">{{ prevArticle.title }}</span>
          </router-link>
          <span v-else class="series-link series-void"></span>
          <router-link v-if="nextArticle" :to="`/article/${nextArticle.slug}`" class="series-link series-next">
            <span class="series-label">下一篇 &rarr;</span>
            <span class="series-title">{{ nextArticle.title }}</span>
          </router-link>
          <span v-else class="series-link series-void"></span>
        </nav>

        <nav class="article-footer-nav">
          <router-link to="/" class="back-link">&larr; 返回首页</router-link>
        </nav>

        <CommentList :comments="comments" :total="commentTotal" />
        <CommentForm :article-id="article.id" @posted="fetchComments" />

        <section v-if="relatedArticles.length" class="related-section">
          <h3 class="related-title">相关文章</h3>
          <div class="related-grid">
            <router-link
              v-for="a in relatedArticles"
              :key="a.id"
              :to="`/article/${a.slug}`"
              class="related-card"
            >
              <span class="related-card-title">{{ a.title }}</span>
              <span class="related-card-date">{{ formatDate(a.created_at) }}</span>
            </router-link>
          </div>
        </section>
      </div>

      <div v-else-if="loading" class="loading">加载中...</div>
      <div v-else class="not-found">文章不存在。</div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
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
const loading = ref(true)
const comments = ref([])
const commentTotal = ref(0)
const relatedArticles = ref([])

async function fetchArticle() {
  loading.value = true
  try {
    const res = await articles.getBySlug(route.params.slug)
    if (res.code === 200) {
      article.value = res.data.article
      article.value.content_html = addHeadingIds(article.value.content_html)
      prevArticle.value = res.data.prev_in_series || null
      nextArticle.value = res.data.next_in_series || null
      fetchComments()
      fetchRelated()
    } else {
      article.value = null
    }
  } catch {
    article.value = null
  } finally {
    loading.value = false
  }
}

async function fetchRelated() {
  if (!article.value) return
  try {
    const params = { page_size: 3 }
    if (article.value.category_id) params.category_id = article.value.category_id
    const res = await articles.list(params)
    if (res.code === 200) {
      relatedArticles.value = (res.data.list || []).filter(a => a.id !== article.value.id).slice(0, 3)
    }
  } catch {}
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
  padding: var(--space-4) 0 var(--space-2);
}

.back-link {
  font-size: var(--text-sm);
  color: var(--color-muted);
  text-decoration: none;
  transition: color var(--duration) var(--ease);
}

.back-link:hover {
  color: var(--color-vermilion);
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
}

.article-content :deep(h1),
.article-content :deep(h2),
.article-content :deep(h3) {
  font-family: var(--font-display);
  color: var(--color-deep);
  margin-top: var(--space-8);
  margin-bottom: var(--space-4);
  scroll-margin-top: 100px;
}

.article-content :deep(h1) { font-size: var(--text-3xl); }
.article-content :deep(h2) { font-size: var(--text-2xl); }
.article-content :deep(h3) { font-size: var(--text-xl); }

.article-content :deep(p) {
  margin-bottom: var(--space-6);
}

.article-content :deep(ul),
.article-content :deep(ol) {
  margin-left: var(--space-6);
  margin-bottom: var(--space-6);
}

.article-content :deep(li) {
  margin-bottom: var(--space-2);
}

.article-content :deep(a) {
  color: var(--color-vermilion);
  text-decoration: underline;
  text-underline-offset: 2px;
}

.article-content :deep(img) {
  display: block;
  max-width: 100%;
  height: auto;
  margin: var(--space-6) auto;
  border-radius: var(--radius);
}

.article-content :deep(hr) {
  border: none;
  border-top: 1px solid var(--color-border);
  margin: var(--space-8) 0;
}

.article-content :deep(pre) {
  margin: var(--space-6) 0;
}

/* Series navigation */
.series-nav {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-4);
  padding: var(--space-6) 0;
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

.series-title {
  font-size: var(--text-sm);
  color: var(--color-ink);
  font-weight: 500;
}

.series-void {
  visibility: hidden;
}

.article-footer-nav {
  display: flex;
  justify-content: center;
  padding: var(--space-6) 0;
}

/* Related articles */
.related-section {
  margin-top: var(--space-12);
  padding-top: var(--space-8);
  border-top: 1px solid var(--color-border);
}

.related-title {
  font-family: var(--font-display);
  font-size: var(--text-xl);
  color: var(--color-deep);
  margin-bottom: var(--space-4);
}

.related-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: var(--space-4);
}

.related-card {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
  padding: var(--space-4);
  border: 1px solid var(--color-border);
  border-radius: var(--radius);
  text-decoration: none;
  transition: all var(--duration) var(--ease);
}

.related-card:hover {
  border-color: var(--color-muted);
  box-shadow: 0 2px 8px rgba(0,0,0,0.04);
}

.related-card-title {
  font-weight: 500;
  font-size: var(--text-sm);
  color: var(--color-ink);
  line-height: 1.4;
}

.related-card-date {
  font-size: var(--text-xs);
  color: var(--color-muted);
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
