<template>
  <div class="comment-list" v-if="comments.length">
    <h3 class="comment-heading">评论 ({{ total }})</h3>
    <div class="comment" v-for="c in comments" :key="c.id">
      <div class="comment-header">
        <strong class="comment-author">{{ c.author }}</strong>
        <time :datetime="c.created_at" class="comment-time">{{ formatDate(c.created_at) }}</time>
      </div>
      <p class="comment-body">{{ c.content }}</p>
      <div class="replies" v-if="c.replies && c.replies.length">
        <div class="reply" v-for="r in c.replies" :key="r.id">
          <div class="comment-header">
            <strong class="comment-author">{{ r.author }}</strong>
            <time :datetime="r.created_at" class="comment-time">{{ formatDate(r.created_at) }}</time>
          </div>
          <p class="comment-body">{{ r.content }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  comments: { type: Array, required: true },
  total: { type: Number, default: 0 },
})

function formatDate(dateStr) {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}
</script>

<style scoped>
.comment-list {
  margin-top: var(--space-12);
  padding-top: var(--space-8);
  border-top: 1px solid var(--color-border);
}

.comment-heading {
  font-family: var(--font-display);
  font-size: var(--text-xl);
  margin-bottom: var(--space-6);
  color: var(--color-deep);
}

.comment {
  padding: var(--space-4) 0;
  border-bottom: 1px solid var(--color-border);
}

.comment:last-child {
  border-bottom: none;
}

.comment-header {
  display: flex;
  align-items: baseline;
  gap: var(--space-3);
  margin-bottom: var(--space-2);
}

.comment-author {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-deep);
}

.comment-time {
  font-size: var(--text-xs);
  color: var(--color-muted);
}

.comment-body {
  font-size: var(--text-sm);
  line-height: 1.6;
  color: var(--color-ink);
  margin-bottom: 0;
}

.replies {
  margin-top: var(--space-3);
  padding-left: var(--space-6);
  border-left: 2px solid var(--color-card);
}

.reply {
  padding: var(--space-3) 0;
}
</style>
