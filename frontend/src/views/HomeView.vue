<template>
  <div class="page-layout container-wide">
    <Sidebar
      :categories="categories"
      :series-list="seriesList"
      :tags="allTags"
    />

    <main class="page-main">
      <header class="hero">
        <h1 class="hero-title" v-html="quote.title"></h1>
        <p class="hero-sub">{{ quote.sub }}</p>
      </header>

      <div class="active-filter" v-if="activeCategory || activeTag || activeSeries">
        <span v-if="activeCategory">
          分类：{{ categoryName }}
          <button class="clear-btn" @click="clearFilter">&times;</button>
        </span>
        <span v-if="activeTag">
          标签：{{ tagName }}
          <button class="clear-btn" @click="clearFilter">&times;</button>
        </span>
        <span v-if="activeSeries">
          系列：{{ activeSeries }}
          <button class="clear-btn" @click="clearFilter">&times;</button>
        </span>
      </div>

      <ArticleList
        :articles="articles"
        :total="total"
        :page="page"
        :page-size="pageSize"
        :loading="loading"
        @page-change="goToPage"
      />
    </main>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { articles as api, categories as catApi, tags as tagApi, series as seriesApi } from '../api'
import ArticleList from '../components/ArticleList.vue'
import Sidebar from '../components/Sidebar.vue'

const route = useRoute()
const router = useRouter()

const articles = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = 10
const loading = ref(true)

const categories = ref([])
const allTags = ref([])
const seriesList = ref([])

const activeCategory = ref(null)
const activeTag = ref(null)
const activeSeries = ref(null)

const quotes = [
  { title: '写代码的人，<br/>最终都在写自己。', sub: 'Notes on code, design & the craft of building.' },
  { title: '天下武功，<br/>唯快不破。', sub: 'Speed is the essence of computing.' },
  { title: '代码即思考。', sub: 'Code as externalized thought.' },
  { title: '把复杂变简单，<br/>就是最好的工程。', sub: 'Simplicity is the ultimate sophistication.' },
  { title: '想清楚再写，<br/>比写清楚再想快一百倍。', sub: 'First think, then type.' },
  { title: '每段代码<br/>都是一个决定。', sub: 'Every line of code is a decision.' },
  { title: '好的架构<br/>是对变化的免疫。', sub: 'Good architecture resists entropy.' },
]
const quote = ref(quotes[Math.floor(Math.random() * quotes.length)])

const categoryName = computed(() => {
  const c = categories.value.find(c => c.id === activeCategory.value)
  return c ? c.name : ''
})

const tagName = computed(() => {
  const t = allTags.value.find(t => t.id === activeTag.value)
  return t ? t.name : ''
})

async function fetchArticles(p = 1) {
  loading.value = true
  const params = { page: p, page_size: pageSize }
  if (activeCategory.value) params.category_id = activeCategory.value
  if (activeTag.value) params.tag_id = activeTag.value
  if (activeSeries.value) params.series = activeSeries.value
  try {
    const res = await api.list(params)
    if (res.code === 200) {
      articles.value = res.data.list || []
      total.value = res.data.total || 0
    }
  } finally {
    loading.value = false
  }
}

async function fetchCategories() {
  try {
    const res = await catApi.list()
    if (res.code === 200) categories.value = res.data || []
  } catch {}
}

async function fetchTags() {
  try {
    const res = await tagApi.list()
    if (res.code === 200) allTags.value = res.data || []
  } catch {}
}

async function fetchSeries() {
  try {
    const res = await seriesApi.list()
    if (res.code === 200) seriesList.value = res.data || []
  } catch {}
}

function goToPage(p) {
  page.value = p
  fetchArticles(p)
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

function filterByCategory(id) {
  activeCategory.value = id
  activeTag.value = null
  activeSeries.value = null
  page.value = 1
  updateQuery()
  fetchArticles()
}

function clearFilter() {
  activeCategory.value = null
  activeTag.value = null
  activeSeries.value = null
  page.value = 1
  updateQuery()
  fetchArticles()
}

function readQuery() {
  activeCategory.value = null
  activeTag.value = null
  activeSeries.value = null
  const cid = parseInt(route.query.category_id)
  const tid = parseInt(route.query.tag_id)
  const s = route.query.series || ''
  if (cid > 0) activeCategory.value = cid
  if (tid > 0) activeTag.value = tid
  if (s) activeSeries.value = s
}

function updateQuery() {
  const q = {}
  if (activeCategory.value) q.category_id = activeCategory.value
  if (activeTag.value) q.tag_id = activeTag.value
  if (activeSeries.value) q.series = activeSeries.value
  router.replace({ query: q })
}

watch(() => route.query, () => {
  readQuery()
  page.value = 1
  fetchArticles()
}, { immediate: true })

onMounted(() => {
  fetchCategories()
  fetchTags()
  fetchSeries()
})
</script>

<style scoped>
.page-layout {
  display: grid;
  grid-template-columns: 200px 1fr;
  gap: var(--space-10);
  padding-top: var(--space-8);
  padding-bottom: var(--space-16);
}

@media (max-width: 900px) {
  .page-layout {
    grid-template-columns: 1fr;
  }
}

.page-main {
  min-width: 0;
}

.hero {
  text-align: left;
  padding: var(--space-6) 0 var(--space-8);
}

@media (max-width: 900px) {
  .hero {
    text-align: center;
    padding: var(--space-8) 0 var(--space-6);
  }
}

.hero-title {
  font-family: var(--font-display);
  font-size: var(--text-4xl);
  font-weight: 700;
  color: var(--color-deep);
  line-height: 1.2;
  letter-spacing: -0.02em;
  max-width: 640px;
  margin: 0 0 var(--space-4);
}

.hero-sub {
  font-size: var(--text-lg);
  color: var(--color-muted);
  font-weight: 300;
}

/* Active filter indicator */
.active-filter {
  display: flex;
  gap: var(--space-4);
  flex-wrap: wrap;
  padding-bottom: var(--space-4);
  font-size: var(--text-sm);
  color: var(--color-muted);
}

.clear-btn {
  font-size: var(--text-lg);
  color: var(--color-muted);
  background: none;
  border: none;
  cursor: pointer;
  padding: 0 var(--space-1);
  line-height: 1;
}

.clear-btn:hover {
  color: var(--color-vermilion);
}
</style>
