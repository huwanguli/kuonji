<template>
  <div class="comment-section" v-if="flatList.length">
    <h3 class="comment-heading">评论 ({{ total }})</h3>

    <div
      v-for="c in flatList"
      :key="c.id"
      :class="['comment', { 'comment-indent': c._depth > 0 }]"
    >
      <div class="comment-header">
        <span class="comment-author">{{ c.author }}</span>
        <time :datetime="c.created_at" class="comment-time">{{ formatDate(c.created_at) }}</time>
      </div>
      <p class="comment-body">
        <span v-if="c.parent_author" class="reply-to-inline">回复 @{{ c.parent_author }}：</span>{{ c.content }}
      </p>
      <a class="comment-reply" href="#" @click.prevent="toggleReply(c.id, c.author)">回复</a>

      <form v-if="replyTarget === c.id" class="reply-inline" @submit.prevent="doReply">
        <span class="reply-label">回复 @{{ c.author }}：</span>
        <div class="reply-fields">
          <input v-model="replyAuthor" type="text" placeholder="昵称" required class="rfi" />
          <button type="submit" class="rfi-btn" :disabled="replying">{{ replying ? '提交中...' : '回复' }}</button>
          <button type="button" class="rfi-cancel" @click="cancelReply">取消</button>
        </div>
        <textarea v-model="replyContent" placeholder="写下你的回复..." required rows="2" class="rft"></textarea>
        <p v-if="replyError" class="rfi-err">回复提交失败，请稍后再试。</p>
        <p v-if="replyTarget === c.id && replyOk" class="rfi-ok">回复已发表！</p>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { comments as commentsApi } from '../api'

const props = defineProps({
  comments: { type: Array, default: () => [] },
  total: { type: Number, default: 0 },
  articleId: { type: Number, required: true },
})

const emit = defineEmits(['posted'])

function buildTree(list) {
  const map = {}
  const roots = []
  for (const c of list) {
    map[c.id] = { ...c, replies: [] }
  }
  for (const c of list) {
    if (c.parent_id && map[c.parent_id]) {
      map[c.parent_id].replies.push(map[c.id])
    } else {
      roots.push(map[c.id])
    }
  }
  return roots
}

function flattenTree(nodes, depth) {
  const out = []
  for (const n of nodes) {
    out.push({ ...n, _depth: depth })
    if (n.replies && n.replies.length) {
      out.push(...flattenTree(n.replies, depth + 1))
    }
  }
  return out
}

const flatList = computed(() => {
  const tree = buildTree(props.comments)
  return flattenTree(tree, 0)
})

const replyTarget = ref(null)
const replyAuthor = ref('')
const replyContent = ref('')
const replying = ref(false)
const replyError = ref('')
const replyOk = ref('')

function toggleReply(id, author) {
  if (replyTarget.value === id) {
    cancelReply()
    return
  }
  replyTarget.value = id
  replyAuthor.value = localStorage.getItem('zb_comment_author') || ''
  replyContent.value = ''
  replyError.value = ''
  replyOk.value = ''
}

function cancelReply() {
  replyTarget.value = null
  replyAuthor.value = ''
  replyContent.value = ''
  replyError.value = ''
  replyOk.value = ''
}

async function doReply() {
  if (!replyAuthor.value.trim() || !replyContent.value.trim()) return
  replying.value = true
  replyError.value = ''
  replyOk.value = ''
  try {
    await commentsApi.create({
      article_id: props.articleId,
      parent_id: replyTarget.value,
      author: replyAuthor.value.trim(),
      content: replyContent.value.trim(),
    })
    localStorage.setItem('zb_comment_author', replyAuthor.value.trim())
    replyOk.value = true
    setTimeout(() => {
      replyOk.value = false
      cancelReply()
    }, 800)
    emit('posted')
  } catch {
    replyError.value = '回复提交失败，请稍后再试。'
  } finally {
    replying.value = false
  }
}

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
.comment-section {
  margin-top: var(--space-12);
  padding-top: var(--space-8);
  border-top: 1px solid var(--color-border);
}

.comment-heading {
  font-family: var(--font-display);
  font-size: var(--text-xl);
  margin: 0 0 var(--space-6);
  color: var(--color-deep);
}

.comment {
  padding: var(--space-4) 0;
  border-bottom: 1px solid var(--color-border);
}

.comment-indent {
  padding-left: var(--space-6);
  border-left: 2px solid var(--color-card);
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
  margin: 0 0 var(--space-2);
}

.reply-to-inline {
  color: var(--color-muted);
}

.comment-reply {
  font-size: var(--text-xs);
  color: var(--color-muted);
  text-decoration: none;
}

.comment-reply:hover {
  color: var(--color-vermilion);
}

.reply-inline {
  margin-top: var(--space-2);
  padding: var(--space-3);
  background: var(--color-card);
  border-radius: var(--radius-sm);
}

.reply-label {
  font-size: var(--text-xs);
  color: var(--color-muted);
  display: block;
  margin-bottom: var(--space-2);
}

.reply-fields {
  display: flex;
  gap: var(--space-2);
  align-items: center;
  margin-bottom: var(--space-2);
}

.rfi {
  font-family: var(--font-body);
  font-size: var(--text-sm);
  padding: var(--space-2) var(--space-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  background: var(--color-paper);
  color: var(--color-ink);
  outline: none;
  flex: 1;
  max-width: 160px;
}

.rfi:focus {
  border-color: var(--color-vermilion);
}

.rft {
  font-family: var(--font-body);
  font-size: var(--text-sm);
  padding: var(--space-2) var(--space-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  background: var(--color-paper);
  color: var(--color-ink);
  outline: none;
  width: 100%;
  resize: vertical;
}

.rft:focus {
  border-color: var(--color-vermilion);
}

.rfi-btn {
  font-family: var(--font-body);
  font-size: var(--text-xs);
  font-weight: 600;
  color: #fff;
  background: var(--color-vermilion);
  border: none;
  padding: var(--space-1) var(--space-3);
  border-radius: var(--radius-sm);
  cursor: pointer;
  white-space: nowrap;
}

.rfi-btn:hover:not(:disabled) {
  opacity: 0.85;
}

.rfi-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.rfi-cancel {
  font-size: var(--text-xs);
  color: var(--color-muted);
  background: none;
  border: none;
  cursor: pointer;
  white-space: nowrap;
}

.rfi-cancel:hover {
  color: var(--color-ink);
}

.rfi-err {
  margin-top: var(--space-2);
  font-size: var(--text-xs);
  color: var(--color-vermilion);
}

.rfi-ok {
  margin-top: var(--space-2);
  font-size: var(--text-xs);
  color: #16a34a;
}
</style>
