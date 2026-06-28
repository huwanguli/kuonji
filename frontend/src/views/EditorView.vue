<template>
  <div class="editor container-wide">
    <div class="editor-header">
        <h2 class="editor-title">{{ isEdit ? '编辑文章' : '写新文章' }}</h2>
        <div class="editor-actions">
          <button @click="publish(0)" class="btn-draft" :disabled="saving">存草稿</button>
          <button @click="publish(1)" class="btn-publish" :disabled="saving">发布</button>
        </div>
      </div>

    <div v-if="!loggedIn" class="not-logged-in">
      <p>请先 <router-link to="/admin">登录</router-link>。</p>
    </div>

    <div v-else class="editor-body">
      <input
        v-model="form.title"
        type="text"
        placeholder="文章标题"
        class="input title-input"
        @input="autoSlug"
      />

      <div class="meta-row">
        <input v-model="form.slug" type="text" placeholder="slug（留空自动生成）" class="input" />
        <input v-model="form.excerpt" type="text" placeholder="摘要（选填）" class="input" />
        <input v-model="form.cover" type="text" placeholder="封面图 URL（选填）" class="input" />
      </div>

      <div class="taxonomy-row">
        <div class="taxonomy-group">
          <label class="taxonomy-label">分类</label>
          <select v-model="form.category_id" class="select">
            <option :value="null">未分类</option>
            <option v-for="c in categories" :key="c.id" :value="c.id">{{ c.name }}</option>
          </select>
        </div>
        <div class="taxonomy-group">
          <label class="taxonomy-label">系列</label>
          <div class="series-inputs">
            <input v-model="form.series" type="text" placeholder="如 building-a-compiler" class="input" />
            <input v-model.number="form.series_order" type="number" placeholder="序号" class="input" min="0" />
          </div>
        </div>
      </div>

      <div class="taxonomy-row taxonomy-tags">
        <div class="taxonomy-group">
          <label class="taxonomy-label">选项</label>
          <div class="option-checkboxes">
            <label class="tag-checkbox">
              <input type="checkbox" v-model="form.is_announcement" :true-value="1" :false-value="0" />
              <span>设为公告（全站横幅 + 文章列表置顶）</span>
            </label>
          </div>
        </div>
      </div>

      <div class="taxonomy-row taxonomy-tags">
        <div class="taxonomy-group">
          <label class="taxonomy-label">标签</label>
          <div class="tag-checkboxes">
            <label v-for="t in allTags" :key="t.id" class="tag-checkbox">
              <input type="checkbox" :value="t.id" v-model="form.tag_ids" />
              <span>{{ t.name }}</span>
            </label>
          </div>
          <p v-if="!allTags.length" class="tag-hint">暂无标签，请到<a href="/admin">管理面板</a>添加。</p>
        </div>
      </div>

      <div class="editor-panes">
        <div class="pane pane-edit">
          <div class="pane-header">
            <span class="pane-label">Markdown</span>
            <div class="pane-tools">
              <button class="tool-btn" @click="insertSyntax('**', '**')" title="加粗">B</button>
              <button class="tool-btn" @click="insertSyntax('*', '*')" title="斜体"><em>I</em></button>
              <button class="tool-btn" @click="insertSyntax('`', '`')" title="代码">&lt;/&gt;</button>
              <button class="tool-btn" @click="insertBlock('> ')" title="引用">&ldquo;</button>
              <button class="tool-btn" @click="uploadTrigger" title="插入图片">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="8.5" cy="8.5" r="1.5"/><polyline points="21 15 16 10 5 21"/></svg>
              </button>
              <input
                ref="fileInput"
                type="file"
                accept=".jpg,.jpeg,.png,.gif,.webp,.svg"
                class="file-hidden"
                @change="handleFileSelect"
              />
            </div>
          </div>
          <textarea
            ref="textareaRef"
            v-model="form.content_md"
            placeholder="用 Markdown 书写..."
            class="textarea"
            @paste="handlePaste"
            @drop.prevent="handleDrop"
            @dragover.prevent
          ></textarea>
        </div>

        <div class="pane pane-preview">
          <div class="pane-header">
            <span class="pane-label">预览</span>
          </div>
          <div class="preview-content" v-html="previewHtml" @click="onPreviewClick"></div>
        </div>
      </div>

      <div v-if="imgPicker.visible" class="img-resize-popover" :style="imgPicker.style">
        <button v-for="w in ['25%','50%','75%','100%','原始']" :key="w"
          class="size-btn" @click="applyImgWidth(w)">
          {{ w }}
        </button>
      </div>

      <p v-if="saveMsg" :class="['save-msg', saveMsg.includes('失败') ? 'save-error' : '' ]">{{ saveMsg }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { articles as api, categories as catApi, tags as tagApi, upload as uploadApi } from '../api'
import { marked } from 'marked'
import { addHeadingIds } from '../utils/html'

const route = useRoute()
const router = useRouter()
const isEdit = computed(() => !!route.params.id)
const loggedIn = ref(!!localStorage.getItem('token'))
const form = ref({ title: '', slug: '', excerpt: '', cover: '', content_md: '', status: 1, category_id: null, tag_ids: [], series: '', series_order: 0, is_announcement: 0 })
const categories = ref([])
const allTags = ref([])
const saving = ref(false)
const saveMsg = ref('')
const textareaRef = ref(null)
const fileInput = ref(null)
const imgPicker = ref({ visible: false, style: {}, src: '' })

function onPreviewClick(e) {
  const img = e.target.closest('img')
  if (!img) { imgPicker.value.visible = false; return }
  const src = extractSrc(img)
  imgPicker.value.src = src
  imgPicker.value.style = {
    left: Math.min(e.clientX, window.innerWidth - 180) + 'px',
    top: (e.clientY - 40) + 'px',
  }
  imgPicker.value.visible = true
}

function extractSrc(img) {
  const u = new URL(img.src, window.location.origin)
  return u.pathname
}

function applyImgWidth(w) {
  const src = imgPicker.value.src
  imgPicker.value.visible = false
  if (!src) return
  const md = form.value.content_md
  let newMd
  if (w === '原始') {
    newMd = md.replace(
      new RegExp(`<img[^>]*src=["']?${escapeRegExp(src)}["']?[^>]*>`, 'g'),
      `![image](${src})`
    )
  } else {
    const idx = md.indexOf(src)
    if (idx === -1) return
    const lines = md.split('\n')
    for (let i = 0; i < lines.length; i++) {
      const li = lines[i].indexOf(src)
      if (li !== -1) {
        if (lines[i].includes('![') && lines[i].includes('](' + src + ')')) {
          lines[i] = `<img src="${src}" width="${w}">`
        } else {
          lines[i] = lines[i].replace(
            new RegExp(`<img[^>]*src=["']?${escapeRegExp(src)}["']?[^>]*>`, 'g'),
            `<img src="${src}" width="${w}">`
          )
        }
        break
      }
    }
    newMd = lines.join('\n')
  }
  form.value.content_md = newMd
}

function escapeRegExp(s) {
  return s.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
}

const previewHtml = computed(() => {
  if (!form.value.content_md) return '<p class="preview-empty">输入内容后实时预览...</p>'
  try {
    return addHeadingIds(marked.parse(form.value.content_md))
  } catch {
    return form.value.content_md
  }
})

function autoSlug() {
  if (isEdit.value) return
}

async function loadArticle() {
  if (!isEdit.value) return
  try {
    const res = await api.adminList()
    if (res.code === 200) {
      const list = res.data.list || []
      const found = list.find(a => a.id === Number(route.params.id))
      if (found) {
        form.value = {
          title: found.title,
          slug: found.slug,
          excerpt: found.excerpt || '',
          cover: found.cover || '',
          content_md: found.content_md,
          status: found.status,
          category_id: found.category_id || null,
          tag_ids: (found.tags || []).map(t => t.id),
          series: found.series || '',
          series_order: found.series_order || 0,
          is_announcement: found.is_announcement || 0,
        }
      }
    }
  } catch {}
}

async function publish(status) {
  saving.value = true
  saveMsg.value = ''
  try {
    form.value.status = status
    let res
    if (isEdit.value) {
      res = await api.update(Number(route.params.id), form.value)
    } else {
      res = await api.create(form.value)
    }
    if (res.code === 200) {
      saveMsg.value = status === 1 ? '已发布！' : '草稿已保存！'
      if (!isEdit.value) {
        setTimeout(() => router.push(`/article/${res.data.slug}`), 800)
      }
    } else {
      saveMsg.value = res.message || '保存失败'
    }
  } catch {
    saveMsg.value = '保存失败，请重试。'
  } finally {
    saving.value = false
  }
}

/* Markdown toolbar helpers */
function getTextarea() {
  return textareaRef.value
}

function insertAtCursor(before, after) {
  const ta = getTextarea()
  if (!ta) return
  const start = ta.selectionStart
  const end = ta.selectionEnd
  const text = form.value.content_md
  const selected = text.substring(start, end)
  form.value.content_md = text.substring(0, start) + before + selected + after + text.substring(end)
  nextTick(() => {
    ta.focus()
    ta.setSelectionRange(start + before.length, start + before.length + selected.length)
  })
}

function insertSyntax(before, after) {
  insertAtCursor(before, after || before)
}

function insertBlock(prefix) {
  const ta = getTextarea()
  if (!ta) return
  const start = ta.selectionStart
  const text = form.value.content_md
  const lineStart = text.lastIndexOf('\n', start - 1) + 1
  form.value.content_md = text.substring(0, lineStart) + prefix + text.substring(lineStart, start) + text.substring(start)
  nextTick(() => {
    ta.focus()
    ta.setSelectionRange(start + prefix.length, start + prefix.length)
  })
}

function insertImageMarkdown(url, alt) {
  const md = `![${alt || 'image'}](${url})`
  const ta = getTextarea()
  if (!ta) return
  const pos = ta.selectionStart
  const text = form.value.content_md
  form.value.content_md = text.substring(0, pos) + md + '\n' + text.substring(pos)
  nextTick(() => ta.focus())
}

/* Image upload */
async function uploadFile(file) {
  saveMsg.value = '图片上传中...'
  try {
    const res = await uploadApi.image(file)
    if (res.code === 200 && res.data) {
      const url = res.data.url
      insertImageMarkdown(url, file.name)
      saveMsg.value = '图片已插入！'
      setTimeout(() => { if (saveMsg.value === '图片已插入！') saveMsg.value = '' }, 2000)
    } else {
      saveMsg.value = '上传失败'
    }
  } catch {
    saveMsg.value = '上传失败，请重试。'
  }
}

function uploadTrigger() {
  fileInput.value?.click()
}

function handleFileSelect(e) {
  const file = e.target.files?.[0]
  if (file) uploadFile(file)
  e.target.value = ''
}

function handlePaste(e) {
  const items = e.clipboardData?.items
  if (!items) return
  for (const item of items) {
    if (item.type.startsWith('image/')) {
      e.preventDefault()
      const file = item.getAsFile()
      if (file) uploadFile(file)
      break
    }
  }
}

function handleDrop(e) {
  const file = e.dataTransfer?.files?.[0]
  if (file && file.type.startsWith('image/')) {
    uploadFile(file)
  }
}

async function fetchTaxonomies() {
  try {
    const [catRes, tagRes] = await Promise.all([catApi.list(), tagApi.list()])
    if (catRes.code === 200) categories.value = catRes.data || []
    if (tagRes.code === 200) allTags.value = tagRes.data || []
  } catch {}
}

onMounted(() => {
  loadArticle()
  fetchTaxonomies()
})
</script>

<style scoped>
.editor {
  padding: var(--space-8) 0 var(--space-20);
}

.editor-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-6);
}

.editor-title {
  font-family: var(--font-display);
  font-size: var(--text-2xl);
  color: var(--color-deep);
}

.editor-actions {
  display: flex;
  gap: var(--space-3);
}

.editor-body {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.input {
  font-family: var(--font-body);
  font-size: var(--text-base);
  padding: var(--space-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  background: var(--color-paper);
  color: var(--color-ink);
}

.input:focus,
.textarea:focus {
  outline: none;
  border-color: var(--color-vermilion);
}

.title-input {
  font-family: var(--font-display);
  font-size: var(--text-2xl);
  font-weight: 600;
}

.meta-row {
  display: flex;
  gap: var(--space-3);
  flex-wrap: wrap;
}

.meta-row .input {
  flex: 1;
  min-width: 200px;
}

/* Taxonomy */
.taxonomy-row {
  display: flex;
  gap: var(--space-4);
  flex-wrap: wrap;
}

.taxonomy-group {
  flex: 1;
  min-width: 200px;
}

.taxonomy-label {
  display: block;
  font-size: var(--text-xs);
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: var(--color-muted);
  margin-bottom: var(--space-2);
}

.select {
  font-family: var(--font-body);
  font-size: var(--text-base);
  padding: var(--space-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  background: var(--color-paper);
  color: var(--color-ink);
  width: 100%;
  cursor: pointer;
}

.select:focus {
  outline: none;
  border-color: var(--color-vermilion);
}

.series-inputs {
  display: flex;
  gap: var(--space-2);
}

.series-inputs .input:first-child {
  flex: 1;
}

.series-inputs .input:last-child {
  width: 72px;
}

.tag-checkboxes {
  display: flex;
  gap: var(--space-3);
  flex-wrap: wrap;
}

.tag-checkbox {
  display: flex;
  align-items: center;
  gap: var(--space-1);
  font-size: var(--text-sm);
  color: var(--color-muted);
  cursor: pointer;
}

.tag-checkbox input[type="checkbox"] {
  accent-color: var(--color-vermilion);
  cursor: pointer;
}

.tag-hint {
  font-size: var(--text-xs);
  color: var(--color-muted);
  margin-top: var(--space-2);
}

.tag-hint a {
  color: var(--color-vermilion);
}

/* Split panes */
.editor-panes {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-4);
  min-height: 60vh;
}

@media (max-width: 900px) {
  .editor-panes {
    grid-template-columns: 1fr;
  }
}

.pane {
  display: flex;
  flex-direction: column;
  border: 1px solid var(--color-border);
  border-radius: var(--radius);
  overflow: hidden;
  background: var(--color-paper);
}

.pane-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--space-2) var(--space-3);
  background: var(--color-card);
  border-bottom: 1px solid var(--color-border);
}

.pane-label {
  font-size: var(--text-xs);
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.06em;
  color: var(--color-muted);
}

.pane-tools {
  display: flex;
  gap: 2px;
}

.tool-btn {
  font-family: var(--font-body);
  font-size: var(--text-xs);
  font-weight: 600;
  color: var(--color-muted);
  background: none;
  border: 1px solid transparent;
  padding: 2px 8px;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: all var(--duration) var(--ease);
  line-height: 1.5;
}

.tool-btn:hover {
  color: var(--color-ink);
  background: var(--color-paper);
  border-color: var(--color-border);
}

.file-hidden {
  display: none;
}

.textarea {
  flex: 1;
  font-family: var(--font-mono);
  font-size: var(--text-sm);
  padding: var(--space-4);
  border: none;
  background: var(--color-paper);
  color: var(--color-ink);
  line-height: 1.7;
  resize: none;
}

.textarea:focus {
  outline: none;
}

.preview-content {
  flex: 1;
  padding: var(--space-4);
  overflow-y: auto;
  font-size: var(--text-base);
  line-height: 1.85;
  color: var(--color-ink);
}

.preview-empty {
  color: var(--color-muted);
  font-style: italic;
  font-size: var(--text-sm);
}

.preview-content :deep(h1) { font-family: var(--font-display); font-size: var(--text-3xl); margin-bottom: var(--space-4); }
.preview-content :deep(h2) { font-family: var(--font-display); font-size: var(--text-2xl); margin-bottom: var(--space-3); }
.preview-content :deep(h3) { font-family: var(--font-display); font-size: var(--text-xl); margin-bottom: var(--space-3); }
.preview-content :deep(p) { margin-bottom: var(--space-4); }
.preview-content :deep(pre) { background: var(--color-deep); color: #e2e8f0; padding: var(--space-4); border-radius: var(--radius); overflow-x: auto; }
.preview-content :deep(code) { font-family: var(--font-mono); font-size: 0.9em; background: var(--color-card); padding: 0.1em 0.3em; border-radius: 2px; }
.preview-content :deep(pre code) { background: none; padding: 0; }
.preview-content :deep(img) { max-width: 100%; border-radius: var(--radius); }
.preview-content :deep(blockquote) { border-left: 3px solid var(--color-vermilion); padding-left: var(--space-4); color: var(--color-muted); margin: var(--space-4) 0; }

.btns {
  display: flex;
  gap: var(--space-3);
}

.btn-publish {
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

.btn-publish:hover:not(:disabled) { opacity: 0.85; }

.btn-draft {
  font-family: var(--font-body);
  font-size: var(--text-sm);
  font-weight: 500;
  color: var(--color-ink);
  background: var(--color-card);
  border: 1px solid var(--color-border);
  padding: var(--space-2) var(--space-4);
  border-radius: var(--radius);
  cursor: pointer;
  transition: border-color var(--duration) var(--ease);
}

.btn-draft:hover:not(:disabled) {
  border-color: var(--color-vermilion);
}

.btn-publish:disabled,
.btn-draft:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.save-msg {
  text-align: center;
  font-size: var(--text-sm);
  color: #16a34a;
}

.save-error {
  color: var(--color-vermilion);
}

.not-logged-in {
  text-align: center;
  padding: var(--space-12) 0;
  color: var(--color-muted);
}

.not-logged-in a {
  color: var(--color-vermilion);
}

.img-resize-popover {
  position: fixed;
  z-index: 200;
  display: flex;
  gap: 2px;
  background: var(--color-deep);
  border-radius: var(--radius);
  padding: 3px;
  box-shadow: 0 4px 16px rgba(0,0,0,0.3);
}

.size-btn {
  background: transparent;
  border: none;
  color: #ccc;
  font-size: var(--text-xs);
  padding: 4px 10px;
  border-radius: 4px;
  cursor: pointer;
  font-family: var(--font-body);
  white-space: nowrap;
  transition: background 0.15s, color 0.15s;
}

.size-btn:hover {
  background: rgba(255,255,255,0.15);
  color: #fff;
}

.preview-content :deep(img) {
  cursor: pointer;
}
</style>
