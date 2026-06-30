<template>
  <div class="editor container-wide">
    <div class="editor-header">
        <h2 class="editor-title">{{ isEdit ? '编辑文章' : '写新文章' }}</h2>
        <div class="editor-actions">
          <button @click="importMdTrigger" class="btn-import" :disabled="saving">导入 .md</button>
          <input
            ref="mdFileInput"
            type="file"
            accept=".md"
            class="file-hidden"
            @change="handleMdImport"
          />
          <button @click="publish(0)" class="btn-draft" :disabled="saving">存草稿</button>
          <button @click="publish(2)" class="btn-private" :disabled="saving">私密</button>
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
        <div class="cover-input-wrap">
          <input v-model="form.cover" type="text" placeholder="封面图 URL（选填）" class="input cover-input" />
          <button class="cover-upload-btn" @click="uploadCoverTrigger" title="上传封面图">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="8.5" cy="8.5" r="1.5"/><polyline points="21 15 16 10 5 21"/></svg>
          </button>
          <input
            ref="coverFileInput"
            type="file"
            accept=".jpg,.jpeg,.png,.gif,.webp"
            class="file-hidden"
            @change="handleCoverFileSelect"
          />
        </div>
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
              <span>设为公告（置顶卡片样式）</span>
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

      <section v-if="isEdit && articleId" class="editor-comments">
        <div class="ec-toolbar">
          <h3 class="ec-heading">评论管理 ({{ editorCommentTotal }})</h3>
          <label class="ec-select-all" v-if="editorComments.length">
            <input type="checkbox" :checked="editorAllSelected" @change="editorToggleAll" />
            全选
          </label>
          <div class="ec-batch-actions" v-if="editorSelectedIds.length">
            <button @click="editorBatchStatus(1)" class="ec-btn ec-approve">通过 ({{ editorSelectedIds.length }})</button>
            <button @click="editorBatchStatus(2)" class="ec-btn ec-delete">删除 ({{ editorSelectedIds.length }})</button>
          </div>
        </div>

        <div v-if="editorCommentsLoading" class="ec-loading">加载评论中...</div>
        <div v-else-if="!editorComments.length" class="ec-empty">暂无评论</div>
        <div v-else class="ec-list">
          <div v-for="c in editorComments" :key="c.id" class="ec-row">
            <input type="checkbox" :checked="editorSelectedIds.includes(c.id)" @change="editorToggleSelect(c.id)" class="ec-check" />
            <div class="ec-body">
              <div class="ec-meta">
                <strong>{{ c.author }}</strong>
                <span :class="['ec-status', 's-' + c.status]">{{ editorStatusLabel(c.status) }}</span>
                <time>{{ editorFormatDateTime(c.created_at) }}</time>
                <span v-if="c.parent_id" class="ec-parent">回复 #{{ c.parent_id }}</span>
              </div>
              <p class="ec-content">{{ c.content }}</p>
            </div>
          </div>
        </div>
      </section>
    </div>

    <input
      ref="batchImageInput"
      type="file"
      accept=".jpg,.jpeg,.png,.gif,.webp,.svg"
      multiple
      class="file-hidden"
      @change="handleBatchImageSelect"
    />

    <div v-if="imgImportDlg.visible" class="modal-overlay" @click.self="closeImageImportDlg">
      <div class="modal modal-import">
        <h3>检测到 {{ imgImportDlg.refs.length }} 个本地图片引用</h3>
        <ul class="import-img-list">
          <li v-for="r in imgImportDlg.results" :key="r.original" :class="'import-img-item ' + r.status">
            <span class="import-img-path" :title="r.original">{{ r.original }}</span>
            <span v-if="r.status === 'pending'" class="import-img-badge pending">待上传</span>
            <span v-else-if="r.status === 'uploading'" class="import-img-badge uploading">上传中...</span>
            <span v-else-if="r.status === 'uploaded'" class="import-img-badge uploaded">已上传</span>
            <span v-else-if="r.status === 'notfound'" class="import-img-badge notfound">未找到文件</span>
            <span v-else-if="r.status === 'failed'" class="import-img-badge failed">上传失败</span>
          </li>
        </ul>
        <div class="modal-actions">
          <button @click="skipImageImport" class="btn-cancel">跳过</button>
          <button @click="triggerBatchImageUpload" class="btn-upload">选择图片上传</button>
          <button v-if="imgImportDlg.results.some(r => r.status === 'notfound' || r.status === 'failed')"
            @click="closeImageImportDlg" class="btn-close-dlg">关闭</button>
          <button v-else-if="imgImportDlg.results.every(r => r.status === 'uploaded')"
            @click="closeImageImportDlg" class="btn-close-dlg">完成</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { articles as api, categories as catApi, tags as tagApi, upload as uploadApi, comments as commentsApi } from '../api'
import { marked } from 'marked'
import { addHeadingIds } from '../utils/html'

const route = useRoute()
const router = useRouter()
const isEdit = computed(() => !!route.params.id)
const articleId = computed(() => isEdit.value ? Number(route.params.id) : null)
const loggedIn = ref(!!localStorage.getItem('token'))
const form = ref({ title: '', slug: '', excerpt: '', cover: '', content_md: '', status: 1, category_id: null, tag_ids: [], series: '', series_order: 0, is_announcement: 0 })
const categories = ref([])
const allTags = ref([])
const saving = ref(false)
const saveMsg = ref('')
const textareaRef = ref(null)
const fileInput = ref(null)
const coverFileInput = ref(null)
const mdFileInput = ref(null)
const batchImageInput = ref(null)
const imgPicker = ref({ visible: false, style: {}, src: '' })
const imgImportDlg = ref({ visible: false, refs: [], results: [] })

const editorComments = ref([])
const editorCommentTotal = ref(0)
const editorCommentsLoading = ref(false)
const editorSelectedIds = ref([])

const editorAllSelected = computed(() => {
  return editorComments.value.length > 0 && editorSelectedIds.value.length === editorComments.value.length
})

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
  fetchEditorComments()
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
      if (status === 2) {
        saveMsg.value = '已设为私密！'
      } else {
        saveMsg.value = status === 1 ? '已发布！' : '草稿已保存！'
      }
      if (status === 1 && res.data?.slug) {
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

/* MD file import */
function importMdTrigger() {
  mdFileInput.value?.click()
}

function handleMdImport(e) {
  const file = e.target.files?.[0]
  if (!file) { e.target.value = ''; return }
  const reader = new FileReader()
  reader.onload = () => {
    form.value.content_md = reader.result || ''
    const name = file.name.replace(/\.md$/i, '')
    if (!isEdit.value && name) {
      form.value.title = name
    }
    saveMsg.value = '已导入：' + file.name
    setTimeout(() => { if (saveMsg.value?.startsWith('已导入')) saveMsg.value = '' }, 2000)
    checkLocalImages()
  }
  reader.readAsText(file)
  e.target.value = ''
}

function parseImageRefs(md) {
  const refs = []
  const seen = new Set()
  const mdPattern = /!\[[^\]]*\]\(([^)]+)\)/g
  const htmlPattern = /<img[^>]+src=["']([^"']+)["'][^>]*>/gi
  let m
  while ((m = mdPattern.exec(md)) !== null) {
    const path = m[1].trim()
    if (isLocalPath(path) && !seen.has(path)) {
      seen.add(path)
      refs.push({ original: path, filename: path.split('/').pop().split('?')[0].split('#')[0] })
    }
  }
  while ((m = htmlPattern.exec(md)) !== null) {
    const path = m[1].trim()
    if (isLocalPath(path) && !seen.has(path)) {
      seen.add(path)
      refs.push({ original: path, filename: path.split('/').pop().split('?')[0].split('#')[0] })
    }
  }
  return refs
}

function isLocalPath(p) {
  return !/^(https?:|data:|\/uploads\/)/i.test(p)
}

function checkLocalImages() {
  const refs = parseImageRefs(form.value.content_md)
  if (!refs.length) return
  imgImportDlg.value = { visible: true, refs, results: refs.map(r => ({ ...r, status: 'pending' })) }
}

function skipImageImport() {
  imgImportDlg.value.visible = false
}

function triggerBatchImageUpload() {
  batchImageInput.value?.click()
}

async function handleBatchImageSelect(e) {
  const files = Array.from(e.target.files || [])
  e.target.value = ''
  if (!files.length || !imgImportDlg.value.visible) return

  const refs = imgImportDlg.value.refs
  const results = imgImportDlg.value.results
  const fileByName = {}
  for (const f of files) fileByName[f.name] = f

  let uploaded = 0
  for (const r of refs) {
    const file = fileByName[r.filename]
    if (!file) {
      results.find(x => x.original === r.original).status = 'notfound'
      continue
    }
    try {
      results.find(x => x.original === r.original).status = 'uploading'
      const res = await uploadApi.image(file)
      if (res.code === 200 && res.data) {
        const url = res.data.url
        form.value.content_md = replaceImagePath(form.value.content_md, r.original, url)
        results.find(x => x.original === r.original).status = 'uploaded'
        results.find(x => x.original === r.original).url = url
        uploaded++
      } else {
        results.find(x => x.original === r.original).status = 'failed'
      }
    } catch {
      results.find(x => x.original === r.original).status = 'failed'
    }
  }

  saveMsg.value = `已上传 ${uploaded}/${refs.length} 张图片`
  setTimeout(() => { if (saveMsg.value?.startsWith('已上传')) saveMsg.value = '' }, 4000)
}

function replaceImagePath(md, oldPath, newUrl) {
  const escaped = escapeRegExp(oldPath)
  md = md.replace(new RegExp(`!\\[([^\\]]*)\\]\\(${escaped}\\)`, 'g'), `![$1](${newUrl})`)
  md = md.replace(new RegExp(`(<img[^>]*src=["'])${escaped}(["'][^>]*>)`, 'gi'), `$1${newUrl}$2`)
  return md
}

function closeImageImportDlg() {
  imgImportDlg.value.visible = false
}

/* Editor comment management */
async function fetchEditorComments() {
  if (!articleId.value) return
  editorCommentsLoading.value = true
  try {
    const res = await commentsApi.adminList({ article_id: articleId.value, page_size: 100 })
    if (res.code === 200 && res.data) {
      editorComments.value = res.data.list || []
      editorCommentTotal.value = res.data.total || editorComments.value.length
    }
  } catch {} finally {
    editorCommentsLoading.value = false
  }
}

function editorToggleSelect(id) {
  const idx = editorSelectedIds.value.indexOf(id)
  if (idx >= 0) {
    editorSelectedIds.value.splice(idx, 1)
  } else {
    editorSelectedIds.value.push(id)
  }
}

function editorToggleAll() {
  if (editorAllSelected.value) {
    editorSelectedIds.value = []
  } else {
    editorSelectedIds.value = editorComments.value.map(c => c.id)
  }
}

async function editorBatchStatus(status) {
  if (!editorSelectedIds.value.length) return
  for (const id of editorSelectedIds.value) {
    try { await commentsApi.updateStatus(id, status) } catch {}
  }
  editorSelectedIds.value = []
  fetchEditorComments()
}

function editorStatusLabel(s) {
  if (s === 0) return '待审核'
  if (s === 2) return '已删除'
  return '正常'
}

function editorFormatDateTime(s) {
  if (!s) return ''
  const d = new Date(s)
  return `${d.getMonth() + 1}/${d.getDate()} ${d.getHours()}:${String(d.getMinutes()).padStart(2, '0')}`
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

function uploadCoverTrigger() {
  coverFileInput.value?.click()
}

function handleFileSelect(e) {
  const file = e.target.files?.[0]
  if (file) uploadFile(file)
  e.target.value = ''
}

async function handleCoverFileSelect(e) {
  const file = e.target.files?.[0]
  if (!file) return
  saveMsg.value = '封面上传中...'
  try {
    const res = await uploadApi.image(file)
    if (res.code === 200 && res.data) {
      form.value.cover = res.data.url
      saveMsg.value = '封面已上传！'
      setTimeout(() => { if (saveMsg.value === '封面已上传！') saveMsg.value = '' }, 2000)
    } else {
      saveMsg.value = '上传失败'
    }
  } catch {
    saveMsg.value = '上传失败，请重试。'
  }
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

.cover-input-wrap {
  flex: 1;
  min-width: 200px;
  display: flex;
  gap: 2px;
}

.cover-input {
  flex: 1;
  min-width: 0;
}

.cover-upload-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 var(--space-3);
  background: var(--color-card);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  color: var(--color-muted);
  transition: all var(--duration) var(--ease);
}

.cover-upload-btn:hover {
  color: var(--color-vermilion);
  border-color: var(--color-vermilion);
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

.btn-import {
  font-family: var(--font-body);
  font-size: var(--text-sm);
  font-weight: 500;
  color: var(--color-muted);
  background: var(--color-card);
  border: 1px solid var(--color-border);
  padding: var(--space-2) var(--space-4);
  border-radius: var(--radius);
  cursor: pointer;
  transition: all var(--duration) var(--ease);
}

.btn-import:hover:not(:disabled) {
  color: var(--color-ink);
  border-color: var(--color-ink);
}

.btn-private {
  font-family: var(--font-body);
  font-size: var(--text-sm);
  font-weight: 500;
  color: var(--color-vermilion);
  background: var(--color-card);
  border: 1px solid var(--color-vermilion);
  padding: var(--space-2) var(--space-4);
  border-radius: var(--radius);
  cursor: pointer;
  transition: opacity var(--duration) var(--ease);
}

.btn-private:hover:not(:disabled) { opacity: 0.85; }

.btn-publish:disabled,
.btn-draft:disabled,
.btn-private:disabled,
.btn-import:disabled {
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

/* Image import dialog */
.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 300;
}

.modal-import {
  background: var(--color-paper);
  border-radius: var(--radius);
  padding: var(--space-6);
  max-width: 520px;
  width: 90%;
  box-shadow: 0 8px 32px rgba(0,0,0,0.12);
  max-height: 80vh;
  overflow-y: auto;
}

.modal-import h3 {
  font-family: var(--font-display);
  font-size: var(--text-lg);
  color: var(--color-deep);
  margin-bottom: var(--space-4);
}

.import-img-list {
  list-style: none;
  margin: 0 0 var(--space-4) 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.import-img-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-3);
  padding: var(--space-2) var(--space-3);
  border-radius: var(--radius-sm);
  font-size: var(--text-sm);
}

.import-img-item.pending { background: var(--color-card); }
.import-img-item.uploading { background: #eff6ff; }
.import-img-item.uploaded { background: #f0fdf4; }
.import-img-item.notfound { background: #fef2f2; }
.import-img-item.failed { background: #fff7ed; }

.import-img-path {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: var(--color-ink);
  font-family: var(--font-mono);
  font-size: var(--text-xs);
}

.import-img-badge {
  font-size: 0.65rem;
  font-weight: 600;
  padding: 1px 6px;
  border-radius: 3px;
  white-space: nowrap;
  flex-shrink: 0;
}

.import-img-badge.pending { background: #f1f5f9; color: #64748b; }
.import-img-badge.uploading { background: #dbeafe; color: #1e40af; }
.import-img-badge.uploaded { background: #d1fae5; color: #065f46; }
.import-img-badge.notfound { background: #fee2e2; color: #991b1b; }
.import-img-badge.failed { background: #ffedd5; color: #9a3412; }

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-3);
  margin-top: var(--space-2);
}

.btn-cancel {
  font-family: var(--font-body);
  font-size: var(--text-sm);
  color: var(--color-muted);
  background: none;
  border: 1px solid var(--color-border);
  padding: var(--space-2) var(--space-4);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: border-color var(--duration) var(--ease);
}

.btn-cancel:hover { border-color: var(--color-muted); }

.btn-upload {
  font-family: var(--font-body);
  font-size: var(--text-sm);
  font-weight: 600;
  color: white;
  background: var(--color-vermilion);
  border: none;
  padding: var(--space-2) var(--space-4);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: opacity var(--duration) var(--ease);
}

.btn-upload:hover { opacity: 0.85; }

.btn-close-dlg {
  font-family: var(--font-body);
  font-size: var(--text-sm);
  font-weight: 500;
  color: var(--color-ink);
  background: var(--color-card);
  border: 1px solid var(--color-border);
  padding: var(--space-2) var(--space-4);
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: border-color var(--duration) var(--ease);
}

.btn-close-dlg:hover { border-color: var(--color-ink); }

/* Editor comment management */
.editor-comments {
  margin-top: var(--space-6);
  border: 1px solid var(--color-border);
  border-radius: var(--radius);
}

.ec-toolbar {
  display: flex;
  align-items: center;
  gap: var(--space-4);
  padding: var(--space-3) var(--space-4);
  border-bottom: 1px solid var(--color-border);
  font-size: var(--text-xs);
}

.ec-heading {
  font-family: var(--font-body);
  font-size: var(--text-base);
  font-weight: 600;
  color: var(--color-deep);
  margin: 0;
}

.ec-select-all {
  display: flex;
  align-items: center;
  gap: var(--space-1);
  color: var(--color-muted);
  cursor: pointer;
}

.ec-batch-actions {
  display: flex;
  gap: var(--space-2);
  flex: 1;
  justify-content: flex-end;
}

.ec-btn {
  font-family: var(--font-body);
  font-size: var(--text-xs);
  font-weight: 500;
  padding: 2px 10px;
  border-radius: var(--radius-sm);
  cursor: pointer;
  border: 1px solid transparent;
  transition: all var(--duration) var(--ease);
}

.ec-approve { color: #065f46; background: #d1fae5; }
.ec-approve:hover { border-color: #065f46; }
.ec-delete { color: var(--color-vermilion); background: #fee2e2; }
.ec-delete:hover { border-color: var(--color-vermilion); }

.ec-loading,
.ec-empty {
  text-align: center;
  padding: var(--space-6);
  font-size: var(--text-sm);
  color: var(--color-muted);
}

.ec-list {
  max-height: 400px;
  overflow-y: auto;
}

.ec-row {
  display: flex;
  gap: var(--space-3);
  padding: var(--space-3) var(--space-4);
  border-bottom: 1px solid var(--color-border);
  align-items: flex-start;
}

.ec-row:last-child { border-bottom: none; }

.ec-check {
  margin-top: 4px;
  flex-shrink: 0;
  accent-color: var(--color-vermilion);
}

.ec-body {
  flex: 1;
  min-width: 0;
}

.ec-meta {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  flex-wrap: wrap;
  margin-bottom: var(--space-1);
  font-size: var(--text-xs);
  color: var(--color-muted);
}

.ec-meta strong {
  color: var(--color-deep);
  font-size: var(--text-sm);
}

.ec-status {
  font-weight: 600;
  padding: 0 6px;
  border-radius: 3px;
  font-size: 0.65rem;
}

.ec-parent { color: var(--color-muted); }

.ec-content {
  font-size: var(--text-sm);
  line-height: 1.4;
  color: var(--color-ink);
  word-break: break-word;
}
</style>
