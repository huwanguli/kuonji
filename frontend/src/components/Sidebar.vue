<template>
  <aside class="sidebar">
    <nav v-if="categories && categories.length" class="sidebar-section">
      <h3 class="sidebar-heading">分类</h3>
      <ul class="sidebar-nav">
        <li class="sidebar-item">
          <router-link
            to="/"
            :class="['sidebar-link', { active: !activeCategory && !activeSeries }]"
          >全部</router-link>
        </li>
        <li v-for="cat in categories" :key="cat.id" class="sidebar-item">
          <router-link
            :to="`/?category_id=${cat.id}`"
            :class="['sidebar-link', { active: activeCategory === cat.id }]"
          >{{ cat.name }}</router-link>
        </li>
      </ul>
    </nav>

    <nav v-if="seriesList && seriesList.length" class="sidebar-section">
      <h3 class="sidebar-heading">系列</h3>
      <ul class="sidebar-nav">
        <li v-for="s in seriesList" :key="s.name" class="sidebar-item">
          <router-link
            :to="`/?series=${s.name}`"
            :class="['sidebar-link', { active: activeSeries === s.name }]"
          >{{ s.name }}</router-link>
          <span class="sidebar-count">{{ s.count }}</span>
        </li>
      </ul>
    </nav>

    <nav v-if="tags && tags.length" class="sidebar-section">
      <h3 class="sidebar-heading">标签</h3>
      <div class="sidebar-tags">
        <router-link
          v-for="t in tags"
          :key="t.id"
          :to="`/?tag_id=${t.id}`"
          :class="['tag-pill', { active: activeTag === t.id }]"
        >#{{ t.name }}</router-link>
      </div>
    </nav>

    <nav v-if="toc.length" class="sidebar-section sidebar-toc">
      <h3 class="sidebar-heading">目录</h3>
      <ul class="sidebar-nav">
        <li
          v-for="item in toc"
          :key="item.id"
          :class="['sidebar-item', item.level === 2 ? 'toc-h2' : '', item.level === 3 ? 'toc-sub' : '']"
        >
          <a
            :href="`#${item.id}`"
            :class="['sidebar-link', { active: activeTocId === item.id }]"
            @click.prevent="scrollTo(item.id)"
          >{{ item.text }}</a>
        </li>
      </ul>
    </nav>
  </aside>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { useRoute } from 'vue-router'

const props = defineProps({
  categories: { type: Array, default: () => [] },
  seriesList: { type: Array, default: () => [] },
  tags: { type: Array, default: () => [] },
  contentHtml: { type: String, default: '' },
})

const route = useRoute()

const activeCategory = ref(null)
const activeTag = ref(null)
const activeSeries = ref(null)
const toc = ref([])
const activeTocId = ref('')

function buildToc() {
  if (!props.contentHtml) {
    toc.value = []
    return
  }
  const parser = new DOMParser()
  const doc = parser.parseFromString(props.contentHtml, 'text/html')
  const headings = doc.querySelectorAll('h1, h2, h3')
  toc.value = Array.from(headings).map(h => ({
    id: h.id,
    text: h.textContent,
    level: parseInt(h.tagName.replace('H', '')),
  })).filter(item => item.id)
}

function readQuery() {
  const cid = parseInt(route.query.category_id)
  activeCategory.value = cid > 0 ? cid : null
  activeTag.value = parseInt(route.query.tag_id) > 0 ? parseInt(route.query.tag_id) : null
  activeSeries.value = route.query.series || null
}

function scrollTo(id) {
  const el = document.getElementById(id)
  if (el) {
    el.scrollIntoView({ behavior: 'smooth', block: 'start' })
    activeTocId.value = id
  }
}

function onScroll() {
  if (!toc.value.length) return
  const headings = toc.value.map(item => document.getElementById(item.id)).filter(Boolean)
  let activeId = ''
  for (const h of headings) {
    const rect = h.getBoundingClientRect()
    if (rect.top <= 120) {
      activeId = h.id
    }
  }
  activeTocId.value = activeId || (headings[0]?.id || '')
}

onMounted(() => {
  readQuery()
  buildToc()
  if (toc.value.length) {
    window.addEventListener('scroll', onScroll, { passive: true })
  }
})

watch(() => props.contentHtml, () => {
  buildToc()
  if (toc.value.length) {
    window.addEventListener('scroll', onScroll, { passive: true })
  }
})

onUnmounted(() => {
  window.removeEventListener('scroll', onScroll)
})
</script>

<style scoped>
.sidebar {
  position: sticky;
  top: calc(var(--navbar-height, 64px) + var(--space-8));
  align-self: start;
  max-height: calc(100vh - var(--navbar-height, 64px) - var(--space-16));
  overflow-y: auto;
  padding-right: var(--space-4);
}

.sidebar-section {
  margin-bottom: var(--space-8);
}

.sidebar-heading {
  font-family: var(--font-display);
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-deep);
  margin-bottom: var(--space-3);
  text-transform: uppercase;
  letter-spacing: 0.06em;
}

.sidebar-nav {
  list-style: none;
}

.sidebar-item {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.toc-h2 {
  padding-left: var(--space-3);
}

.toc-sub {
  padding-left: var(--space-6);
}

.sidebar-link {
  display: block;
  font-size: var(--text-sm);
  color: var(--color-muted);
  text-decoration: none;
  padding: var(--space-1) 0;
  transition: color var(--duration) var(--ease);
  line-height: 1.6;
}

.sidebar-link:hover {
  color: var(--color-ink);
}

.sidebar-link.active {
  color: var(--color-vermilion);
  font-weight: 500;
}

.sidebar-count {
  font-size: var(--text-xs);
  color: var(--color-muted);
  background: var(--color-card);
  padding: 0 6px;
  border-radius: 10px;
  flex-shrink: 0;
}

.sidebar-toc .sidebar-link {
  font-size: var(--text-xs);
}

.sidebar-tags {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-2);
}

.tag-pill {
  font-size: var(--text-xs);
  color: var(--color-muted);
  background: var(--color-card);
  padding: 3px 10px;
  border-radius: 100px;
  text-decoration: none;
  font-weight: 500;
  transition: all var(--duration) var(--ease);
  border: 1px solid transparent;
  white-space: nowrap;
}

.tag-pill:hover {
  color: var(--color-ink);
  border-color: var(--color-border);
}

.tag-pill.active {
  color: var(--color-vermilion);
  border-color: var(--color-vermilion);
  background: transparent;
}

@media (max-width: 900px) {
  .sidebar {
    display: none;
  }
}
</style>
