<template>
  <form @submit.prevent="submit" class="comment-form">
    <h4 class="form-title">发表评论</h4>
    <div class="form-row">
      <input
        v-model="author"
        type="text"
        placeholder="昵称"
        required
        class="input"
        aria-label="你的昵称"
      />
      <input
        v-model="email"
        type="email"
        placeholder="Email（选填）"
        class="input"
        aria-label="你的邮箱"
      />
    </div>
    <textarea
      v-model="content"
      placeholder="写下你的想法..."
      required
      rows="3"
      class="textarea"
      aria-label="评论内容"
    ></textarea>
    <div class="form-actions">
      <button type="submit" class="submit-btn" :disabled="submitting">
        {{ submitting ? '提交中...' : '提交评论' }}
      </button>
    </div>
    <p v-if="error" class="error">{{ error }}</p>
    <p v-if="success" class="success">评论已发表！</p>
  </form>
</template>

<script setup>
import { ref } from 'vue'
import { comments } from '../api'

const props = defineProps({
  articleId: { type: Number, required: true },
})

const emit = defineEmits(['posted'])

const author = ref('')
const email = ref('')
const content = ref('')
const submitting = ref(false)
const error = ref('')
const success = ref(false)

async function submit() {
  if (!author.value.trim() || !content.value.trim()) return
  submitting.value = true
  error.value = ''
  success.value = false
  try {
    await comments.create({
      article_id: props.articleId,
      author: author.value.trim(),
      email: email.value.trim(),
      content: content.value.trim(),
    })
    author.value = ''
    email.value = ''
    content.value = ''
    success.value = true
    emit('posted')
  } catch (e) {
    error.value = '评论提交失败，请稍后再试。'
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.comment-form {
  margin-top: var(--space-12);
  padding: var(--space-6);
  background: var(--color-card);
  border-radius: var(--radius);
}

.form-title {
  font-family: var(--font-display);
  font-size: var(--text-lg);
  margin-bottom: var(--space-4);
  color: var(--color-deep);
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-3);
  margin-bottom: var(--space-3);
}

@media (max-width: 500px) {
  .form-row {
    grid-template-columns: 1fr;
  }
}

.input,
.textarea {
  font-family: var(--font-body);
  font-size: var(--text-sm);
  padding: var(--space-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  background: var(--color-paper);
  color: var(--color-ink);
  outline: none;
  transition: border-color var(--duration) var(--ease);
}

.input:focus,
.textarea:focus {
  border-color: var(--color-vermilion);
}

.textarea {
  width: 100%;
  resize: vertical;
  margin-bottom: var(--space-3);
}

.form-actions {
  display: flex;
  justify-content: flex-end;
}

.submit-btn {
  font-family: var(--font-body);
  font-size: var(--text-sm);
  font-weight: 600;
  color: white;
  background: var(--color-vermilion);
  border: none;
  padding: var(--space-2) var(--space-6);
  border-radius: var(--radius);
  cursor: pointer;
  transition: opacity var(--duration) var(--ease);
}

.submit-btn:hover:not(:disabled) {
  opacity: 0.85;
}

.submit-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.error {
  margin-top: var(--space-3);
  font-size: var(--text-sm);
  color: var(--color-vermilion);
}

.success {
  margin-top: var(--space-3);
  font-size: var(--text-sm);
  color: #16a34a;
}
</style>
