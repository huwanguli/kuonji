<template>
  <section>
    <div v-if="loading" class="loading">加载中...</div>
    <div v-else-if="!articles.length" class="empty">还没有文章，稍后再来看看。</div>
    <div v-else>
      <ArticleCard v-for="article in articles" :key="article.id" :article="article" />
      <div class="pagination" v-if="total > pageSize">
        <button
          class="page-btn"
          :disabled="page <= 1"
          @click="$emit('page-change', page - 1)"
          aria-label="Previous page"
        >&larr; 上一页</button>
        <span class="page-info">第 {{ page }} / {{ totalPages }} 页</span>
        <button
          class="page-btn"
          :disabled="page >= totalPages"
          @click="$emit('page-change', page + 1)"
          aria-label="下一页"
        >下一页 &rarr;</button>
      </div>
    </div>
  </section>
</template>

<script setup>
import { computed } from 'vue'
import ArticleCard from './ArticleCard.vue'

const props = defineProps({
  articles: { type: Array, required: true },
  total: { type: Number, default: 0 },
  page: { type: Number, default: 1 },
  pageSize: { type: Number, default: 10 },
  loading: { type: Boolean, default: false },
})

defineEmits(['page-change'])

const totalPages = computed(() => Math.ceil(props.total / props.pageSize))
</script>

<style scoped>
.loading,
.empty {
  text-align: center;
  padding: var(--space-16) 0;
  color: var(--color-muted);
  font-size: var(--text-lg);
}

.empty {
  font-style: italic;
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
</style>
