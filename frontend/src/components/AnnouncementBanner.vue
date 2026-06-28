<template>
  <div v-if="visibleAnnouncements.length" class="announcement-banner">
    <div
      v-for="a in visibleAnnouncements"
      :key="a.id"
      class="announcement-item"
    >
      <router-link :to="`/article/${a.slug}`" class="announcement-link">
        <span class="announcement-icon">📢</span>
        <span class="announcement-title">{{ a.title }}</span>
        <span class="announcement-date">{{ formatDate(a.created_at) }}</span>
      </router-link>
      <button class="announcement-close" @click="dismiss(a.id)" title="关闭">&times;</button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { announcements as api } from '../api'

const announcements = ref([])

const DISMISS_KEY = 'zb_dismissed_announcements'

function getDismissed() {
  try {
    const raw = localStorage.getItem(DISMISS_KEY)
    if (!raw) return {}
    const data = JSON.parse(raw)
    if (data.expires && Date.now() > data.expires) {
      localStorage.removeItem(DISMISS_KEY)
      return {}
    }
    return data.ids || {}
  } catch {
    return {}
  }
}

function dismiss(id) {
  const dismissed = getDismissed()
  dismissed[id] = true
  localStorage.setItem(DISMISS_KEY, JSON.stringify({
    ids: dismissed,
    expires: Date.now() + 24 * 60 * 60 * 1000,
  }))
  announcements.value = announcements.value.filter(a => a.id !== id)
}

const visibleAnnouncements = ref([])

function formatDate(d) {
  if (!d) return ''
  return new Date(d).toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' })
}

onMounted(async () => {
  try {
    const res = await api.list()
    if (res.code === 200 && res.data) {
      const dismissed = getDismissed()
      announcements.value = res.data
      visibleAnnouncements.value = res.data.filter(a => !dismissed[a.id])
    }
  } catch {}
})
</script>

<style scoped>
.announcement-banner {
  background: var(--color-deep);
  padding: 0;
}

.announcement-item {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--space-2) var(--space-4);
  border-bottom: 1px solid rgba(255,255,255,0.08);
  position: relative;
}

.announcement-item:last-child {
  border-bottom: none;
}

.announcement-link {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  color: #fff;
  text-decoration: none;
  font-size: var(--text-sm);
  max-width: var(--max-width-wide);
  width: 100%;
  justify-content: center;
}

.announcement-link:hover .announcement-title {
  text-decoration: underline;
}

.announcement-icon {
  flex-shrink: 0;
}

.announcement-title {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.announcement-date {
  color: rgba(255,255,255,0.5);
  font-size: var(--text-xs);
  flex-shrink: 0;
  margin-left: var(--space-2);
}

.announcement-close {
  position: absolute;
  right: var(--space-3);
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  color: rgba(255,255,255,0.4);
  font-size: var(--text-lg);
  cursor: pointer;
  line-height: 1;
  padding: 0 var(--space-1);
}

.announcement-close:hover {
  color: rgba(255,255,255,0.8);
}
</style>
