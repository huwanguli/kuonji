<template>
  <article class="card">
    <div class="card-meta">
      <router-link
        v-if="article.category"
        :to="`/?category_id=${article.category.id}`"
        class="card-category"
      >{{ article.category.name }}</router-link>
      <time :datetime="article.created_at" class="card-date">{{ formatDate(article.created_at) }}</time>
    </div>
    <router-link :to="`/article/${article.slug}`" class="card-link">
      <h2 class="card-title">{{ article.title }}</h2>
      <p class="card-excerpt" v-if="article.excerpt">{{ article.excerpt }}</p>
      <p class="card-excerpt" v-else>{{ stripHtml(article.content_html) }}</p>
    </router-link>
    <div class="card-footer">
      <div class="card-tags" v-if="article.tags && article.tags.length">
        <router-link
          v-for="tag in article.tags"
          :key="tag.id"
          :to="`/?tag_id=${tag.id}`"
          class="chip"
        >{{ tag.name }}</router-link>
      </div>
      <router-link :to="`/article/${article.slug}`" class="card-read">
        阅读全文 &rarr;
      </router-link>
    </div>
  </article>
</template>

<script setup>
defineProps({
  article: { type: Object, required: true },
})

function formatDate(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return `${d.getFullYear()}年${d.getMonth() + 1}月${d.getDate()}日`
}

function stripHtml(html) {
  if (!html) return ''
  return html.replace(/<[^>]*>/g, '').substring(0, 200) + '...'
}
</script>

<style scoped>
.card {
  background: var(--color-paper);
  border: 1px solid var(--color-border);
  border-radius: var(--radius);
  padding: var(--space-6);
  margin-bottom: var(--space-6);
  transition: border-color var(--duration) var(--ease), box-shadow var(--duration) var(--ease);
  position: relative;
}

.card:hover {
  border-color: var(--color-muted);
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
}

.card-meta {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  margin-bottom: var(--space-3);
  flex-wrap: wrap;
}

.card-category {
  font-size: var(--text-xs);
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: var(--color-vermilion);
  text-decoration: none;
}

.card-category:hover {
  opacity: 0.7;
}

.card-date {
  font-size: var(--text-sm);
  color: var(--color-muted);
}

.card-link {
  display: block;
  text-decoration: none;
  color: inherit;
}

.card-link:hover .card-title {
  color: var(--color-vermilion);
}

.card-title {
  font-family: var(--font-display);
  font-size: var(--text-2xl);
  font-weight: 600;
  color: var(--color-deep);
  margin-bottom: var(--space-2);
  transition: color var(--duration) var(--ease);
  line-height: 1.25;
}

.card-excerpt {
  font-size: var(--text-base);
  color: var(--color-muted);
  line-height: 1.6;
  margin-bottom: 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  margin-top: var(--space-4);
  gap: var(--space-4);
}

.card-tags {
  display: flex;
  gap: var(--space-2);
  flex-wrap: wrap;
  flex: 1;
}

.chip {
  font-size: var(--text-xs);
  color: var(--color-muted);
  background: var(--color-card);
  padding: 2px 10px;
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

.card-read {
  font-size: var(--text-sm);
  font-weight: 500;
  color: var(--color-vermilion);
  text-decoration: none;
  white-space: nowrap;
}

.card-read:hover {
  opacity: 0.7;
}

@media (max-width: 500px) {
  .card {
    padding: var(--space-4);
    margin-bottom: var(--space-4);
  }

  .card-footer {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--space-2);
  }
}
</style>
