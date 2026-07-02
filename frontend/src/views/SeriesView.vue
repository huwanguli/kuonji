<template>
  <div class="page-layout container-wide">
    <Sidebar :show-series="false" />

    <main class="page-main">
      <div v-if="loading" class="loading">加载中...</div>
      <div v-else-if="!series" class="not-found">系列不存在。</div>
      <template v-else>
        <nav class="breadcrumb">
          <router-link to="/" class="back-link">&larr; 返回首页</router-link>
        </nav>

        <header class="series-header">
          <div class="series-cover" v-if="series.cover">
            <img :src="series.cover" :alt="series.name" />
          </div>
          <h1 class="series-title">{{ series.name }}</h1>
          <p class="series-desc" v-if="series.description">{{ series.description }}</p>
          <p class="series-count">共 {{ articles.length }} 篇文章</p>
        </header>

        <div class="series-list">
          <router-link
            v-for="a in articles"
            :key="a.id"
            :to="`/article/${a.slug}`"
            class="series-row"
          >
            <span class="series-order">{{ a.series_order || '—' }}</span>
            <span class="series-article-title">{{ a.title }}</span>
            <span class="series-date">{{ formatDate(a.created_at) }}</span>
          </router-link>
        </div>
      </template>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import Sidebar from '../components/Sidebar.vue'
import { series as seriesApi } from '../api'

const route = useRoute()
const series = ref(null)
const articles = ref([])
const loading = ref(true)

function formatDate(s) {
  if (!s) return ''
  const d = new Date(s)
  return `${d.getFullYear()}/${d.getMonth() + 1}/${d.getDate()}`
}

onMounted(async () => {
  try {
    const res = await seriesApi.detail(route.params.name)
    if (res.code === 200 && res.data) {
      series.value = res.data.series || null
      articles.value = res.data.articles || []
    }
  } catch {} finally {
    loading.value = false
  }
})
</script>

<style scoped>
.page-layout {
  display: grid;
  grid-template-columns: 200px 1fr;
  gap: var(--space-8);
  padding-top: var(--space-8);
}

@media (max-width: 900px) {
  .page-layout {
    grid-template-columns: 1fr;
  }
}

.page-main {
  min-width: 0;
}

.loading,
.not-found {
  text-align: center;
  padding: var(--space-16) 0;
  color: var(--color-muted);
}

.breadcrumb {
  margin-bottom: var(--space-6);
}

.back-link {
  font-size: var(--text-sm);
  color: var(--color-muted);
  text-decoration: none;
  transition: color var(--duration) var(--ease);
}

.back-link:hover {
  color: var(--color-ink);
}

.series-header {
  margin-bottom: var(--space-8);
}

.series-title {
  font-family: var(--font-display);
  font-size: var(--text-3xl);
  color: var(--color-deep);
  margin-bottom: var(--space-2);
}

.series-desc {
  font-size: var(--text-base);
  color: var(--color-muted);
  margin-bottom: var(--space-2);
}

.series-count {
  font-size: var(--text-sm);
  color: var(--color-muted);
}

.series-cover {
  margin-bottom: var(--space-4);
}

.series-cover img {
  max-width: 100%;
  max-height: 300px;
  border-radius: var(--radius);
  object-fit: cover;
}

.series-list {
  display: flex;
  flex-direction: column;
  border: 1px solid var(--color-border);
  border-radius: var(--radius);
  overflow: hidden;
}

.series-row {
  display: flex;
  align-items: center;
  gap: var(--space-4);
  padding: var(--space-3) var(--space-4);
  text-decoration: none;
  border-bottom: 1px solid var(--color-border);
  transition: background var(--duration) var(--ease);
}

.series-row:last-child {
  border-bottom: none;
}

.series-row:hover {
  background: var(--color-card);
}

.series-order {
  font-family: var(--font-mono);
  font-size: var(--text-sm);
  color: var(--color-muted);
  width: 24px;
  text-align: center;
  flex-shrink: 0;
}

.series-article-title {
  flex: 1;
  font-size: var(--text-base);
  color: var(--color-ink);
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.series-date {
  font-size: var(--text-xs);
  color: var(--color-muted);
  white-space: nowrap;
}
</style>
