# Kuonji Blog · 前端文档

## 技术栈

| 层 | 技术 |
|---|---|
| 框架 | Vue 3.5 (Composition API, `<script setup>`) |
| 构建 | Vite 8.1 |
| 路由 | Vue Router 4.6 (history mode) |
| HTTP | Axios 1.18 |
| Markdown 渲染 | marked 18.0（编辑器预览）+ github-markdown-css 5.9（全局样式） |
| 样式 | 手写 CSS 设计系统（CSS 变量 + 亮/暗主题） |
| 状态管理 | 无（组件内 `ref`/`reactive`） |

## 项目结构

```
frontend/src/
├── App.vue                    # 根组件：导航栏 + 路由视图 + 页脚
├── main.js                    # 入口：挂载 Vue + Router + 全局 CSS
│
├── api/
│   └── index.js               # Axios 实例 + 所有 API 方法
│
├── router/
│   └── index.js               # 路由表 + 滚动行为
│
├── assets/
│   └── main.css               # 全局样式：CSS 变量、Reset、暗色主题覆盖
│
├── composables/
│   └── useTheme.js            # 主题切换逻辑（亮/暗 + localStorage 持久化）
│
├── utils/
│   └── html.js                # addHeadingIds（为标题生成锚点 ID）
│
├── components/                # 可复用组件
│   ├── NavBar.vue             # 顶部导航栏（公开/管理两种模式）
│   ├── Footer.vue             # 页脚版权
│   ├── Sidebar.vue            # 侧边栏：分类/系列/标签筛选 + TOC 目录
│   ├── ArticleList.vue        # 文章列表容器（含分页）
│   ├── ArticleCard.vue        # 文章卡片（公告特殊样式）
│   ├── CommentList.vue        # 评论列表（嵌套回复树）
│   └── CommentForm.vue        # 评论提交表单
│
└── views/                     # 页面级组件
    ├── HomeView.vue           # 首页：文章列表 + 分类/标签/系列筛选
    ├── ArticleView.vue        # 文章详情：正文 + 系列导航 + 评论区
    ├── SeriesView.vue         # 系列详情：封面/简介 + 文章列表
    ├── EditorView.vue         # 文章编辑器：Markdown 编辑 + 实时预览
    ├── AdminView.vue          # 管理入口：登录 / 仪表盘 + 分类/标签/系列管理
    ├── AdminArticles.vue      # 文章管理：列表 + 状态筛选 + 评论管理
    ├── AdminArticleView.vue   # 文章预览（只读）：正文 + 评论区管理
    └── AdminComments.vue      # 全局评论管理
```

## 页面架构与路由

```
 ┌──────────────────────────────────────────────────────────────┐
 │                         NavBar（固定顶栏）                     │
 ├──────────────────────────────────────────────────────────────┤
 │                                                              │
 │  ┌──────────┐  ┌──────────────────────────────────────────┐  │
 │  │          │  │                                          │  │
 │  │ Sidebar  │  │         <router-view>                    │  │
 │  │          │  │                                          │  │
 │  │ 分类     │  │                                          │  │
 │  │ 系列     │  │                                          │  │
 │  │ 标签     │  │                                          │  │
 │  │ TOC      │  │                                          │  │
 │  │          │  │                                          │  │
 │  └──────────┘  └──────────────────────────────────────────┘  │
 │                                                              │
 ├──────────────────────────────────────────────────────────────┤
 │                         Footer                               │
 └──────────────────────────────────────────────────────────────┘
```

### 路由表

| 路径 | 组件 | 说明 |
|---|---|---|
| `/` | `HomeView` | 首页：文章列表 + 分类/标签/系列筛选 |
| `/article/:slug` | `ArticleView` | 文章详情：正文 + 系列导航 + 评论 |
| `/series/:name` | `SeriesView` | 系列详情：封面/简介 + 文章列表 |
| `/admin` | `AdminView` | 未登录→登录页；已登录→仪表盘 |
| `/admin/articles` | `AdminArticles` | 文章管理（状态筛选 + 行内评论管理） |
| `/admin/articles/view/:id` | `AdminArticleView` | 管理员文章预览（只读）+ 评论管理 |
| `/admin/comments` | `AdminComments` | 全局评论管理（状态筛选 + 批量操作） |
| `/editor` | `EditorView` | 新建文章 |
| `/editor/:id` | `EditorView` | 编辑文章 |
| `/:pathMatch(.*)*` | → `/` | 404 重定向到首页 |

### 页面跳转关系图

```
                          ┌─────────────────────┐
                          │      HomeView (/)    │
                          │  文章列表 + 筛选      │
                          └──────┬──────────────┘
                                 │ 点击标题
                                 ▼
                          ┌─────────────────────┐
                          │ ArticleView          │
                          │ /article/:slug       │
                          │ 正文 + 系列导航 + 评论  │
                          └─────────────────────┘
                                 │
              ┌──────────────────┼──────────────────┐
              ▼                  ▼                  ▼
      系列上一篇/下一篇    点击标签→HomeView?tag_id=X    无跳转（页内锚点/评论）

═══════════════════════════════════════════════════════════════

                          ┌─────────────────────┐
                          │  AdminView (/admin)  │
                          │  登录 → 仪表盘        │
                          └──────┬──────────────┘
                                 │
         ┌───────────┬───────────┼───────────┬───────────┐
         ▼           ▼           ▼           ▼           ▼
    AdminArticles  AdminComments  EditorView  EditorView   (分类/标签 CRUD
    /admin/articles /admin/comments /editor   /editor/:id   仪表盘内行操作)

         │ 点击标题          │ 编辑按钮
         ▼                   ▼
    AdminArticleView     EditorView
    /admin/articles/     /editor/:id
    view/:id
```

## 组件树

```
App.vue
├── NavBar.vue
│   └── useTheme
├── Footer.vue
└── <router-view>
    │
    ├── HomeView.vue
    │   ├── Sidebar.vue          (分类/系列/标签列表)
    │   └── ArticleList.vue
    │       └── ArticleCard.vue  (每篇文章卡片，公告特殊样式)
    │
    ├── ArticleView.vue
    │   ├── Sidebar.vue          (+ TOC 目录)
    │   ├── CommentForm.vue
    │   └── CommentList.vue      (嵌套回复树)
    │
    ├── SeriesView.vue         (不含子组件，全部内联)
    │
    ├── EditorView.vue           (不含子组件，全部内联)
    │
    ├── AdminView.vue            (不含子组件，分类/标签 CRUD 内联)
    │
    ├── AdminArticles.vue        (不含子组件，评论管理内联)
    │
    ├── AdminArticleView.vue     (不含子组件，评论管理内联)
    │
    └── AdminComments.vue        (不含子组件)
```

## API 层

### Axios 实例配置

- **Base URL:** `/api`（开发时 Vite proxy 到 `localhost:8080`）
- **超时:** 10s
- **响应拦截器:** 检测 `code === 401`，清除 token 并重定向到 `/admin?expired=1`

### API 方法一览

```
articles
├── list(params)              → GET  /api/articles            (公开)
├── getBySlug(slug)           → GET  /api/articles/:slug      (公开)
├── adminList(params)         → GET  /api/admin/articles      (需认证)
├── adminDetail(id)           → GET  /api/admin/articles/:id  (需认证)
├── create(data)              → POST /api/admin/articles      (需认证)
├── update(id, data)          → PUT  /api/admin/articles/:id  (需认证)
└── delete(id)                → DEL  /api/admin/articles/:id  (需认证)

categories
├── list()                    → GET  /api/categories          (公开)
├── adminList()               → GET  /api/admin/categories    (需认证)
├── create(data)              → POST /api/admin/categories    (需认证)
├── update(id, data)          → PUT  /api/admin/categories/:id(需认证)
└── delete(id)                → DEL  /api/admin/categories/:id(需认证)

tags
├── list()                    → GET  /api/tags                (公开)
├── create(data)              → POST /api/admin/tags          (需认证)
├── update(id, data)          → PUT  /api/admin/tags/:id      (需认证)
└── delete(id)                → DEL  /api/admin/tags/:id      (需认证)

comments
├── list(articleId, params)   → GET  /api/comments            (公开)
├── create(data)              → POST /api/comments            (公开)
├── adminList(params)         → GET  /api/admin/comments      (需认证)
├── updateStatus(id, status)  → PUT  /api/admin/comments/:id  (需认证)
└── delete(id)                → DEL  /api/admin/comments/:id  (需认证)

auth
├── login(username, password) → POST /api/admin/login         (公开)
└── profile()                 → GET  /api/admin/profile       (需认证)

upload
└── image(file)               → POST /api/admin/upload        (需认证, multipart)

series
├── list()                    → GET  /api/series              (公开)
├── detail(name)              → GET  /api/series/:name        (公开)
├── create(data)              → POST /api/admin/series        (需认证)
├── update(id, data)          → PUT  /api/admin/series/:id    (需认证)
└── delete(id)                → DEL  /api/admin/series/:id    (需认证)

announcements
└── list()                    → GET  /api/announcements       (公开)
```

## 认证流程

```
用户输入账号密码 → POST /api/admin/login
                       │
           ┌───────────┴───────────┐
           ▼                       ▼
       code=200               code≠200
    token→localStorage        显示错误提示
    跳转 AdminView 仪表盘
           │
           ▼
   每个需认证的 API 调用带上
   Authorization: Bearer <token>
           │
           ▼
   API 返回 code=401 时
   → 清除 token → 重定向 /admin?expired=1
```

无路由守卫。每个管理页面在 `onMounted` 时自行检查 `localStorage.getItem('token')` 判断是否已登录。

## 关键功能

### 1. 文章编辑器（EditorView.vue）

**功能清单：**
- 左右分栏：Markdown 源码 + 实时预览（marked 解析）
- 工具栏按钮：加粗 / 斜体 / 行内代码 / 引用（光标处插入语法）
- 图片上传：按钮选择 / 剪贴板粘贴 / 拖拽 / 封面上传
- 预览中点击图片弹出尺寸选择器（25% / 50% / 75% / 100% / 原始）
- .md 文件导入：自动填充内容+标题，检测本地图片引用并批量上传
- 分类下拉 + 系列输入 + 标签勾选 + 公告 checkbox
- 三态保存：存草稿 (status=0) / 私密 (status=2) / 发布 (status=1)
- 客户端校验：标题和正文为空时阻止提交并显示错误提示
- 编辑已有文章时底部展示评论管理面板（批量通过/删除）

### 2. 文章渲染（ArticleView / AdminArticleView / EditorPreview）

**渲染管道：**

```
后端 goldmark 预渲染 content_html → 存入数据库
                                      │
                                      ▼
前端 ArticleView / AdminArticleView   │
  addHeadingIds(content_html)        │
  v-html 渲染到 .markdown-body       │
  github-markdown-css 提供样式        │
                                      │
编辑器预览（EditorView）              │
  marked.parse(content_md)           │
  addHeadingIds(html)                │
  v-html 渲染到 .markdown-body        │
  github-markdown-css 提供样式         │
```

**样式来源：**
- `github-markdown-css/github-markdown-light.css`：全局 markdown 元素样式
- 各组件 scoped CSS：字体族覆盖 (`--font-display`)、图片居中、scroll-margin-top（锚点跳转）
- `main.css` 暗色覆盖：`html[data-theme="dark"] .markdown-body { ... }`

### 3. 评论系统

**数据模型：**
- 状态值：0=待审核（预留）/ 1=通过（默认）/ 2=已删除（软删除）
- 支持嵌套回复（`parent_id`），前端平铺渲染（`_depth` 缩进）
- 作者名缓存到 `localStorage`

**评论管理位置：**
- `AdminComments.vue`：全局评论管理（状态筛选 + 单条操作）
- `AdminArticles.vue`：文章列表行内评论面板（展开管理）
- `AdminArticleView.vue`：文章预览页底部评论管理
- `EditorView.vue`：编辑页底部评论管理

### 4. 主题切换

- 亮色/暗色两种主题
- 通过 `useTheme` composable 管理
- `html[data-theme="dark"]` 切换 CSS 变量值
- 持久化到 `localStorage` key `zb-theme`
- NavBar 右上角太阳/月亮图标切换

## 设计系统

### CSS 变量（亮色主题）

```css
--color-paper:    #fdfbf7    /* 页面背景 */
--color-ink:      #1a1a1a    /* 正文颜色 */
--color-vermilion:#e63946    /* 强调色（链接/按钮） */
--color-deep:     #2b2d42    /* 标题颜色 */
--color-card:     #f5f3ee    /* 卡片/区块背景 */
--color-border:   #d1cdc6    /* 边框/分隔线 */
--color-muted:    #8a8780    /* 次要文字 */
--color-navbar-bg:rgba(253,251,247,0.85)
```

### 字体

| 用途 | 字体 |
|---|---|
| 标题 | Playfair Display |
| 正文 | Inter |
| 代码 | JetBrains Mono |

### 布局

- 文章宽度：720px（`.container`）
- 管理页宽度：1200px（`.container-wide`）
- 响应式断点：900px（折叠侧边栏）、640px（单列布局）

### 暗色主题

`html[data-theme="dark"]` 覆盖所有颜色变量 + markdown-body 内部元素（表格、代码块、引用、分割线等）。

## 开发

```bash
cd frontend
npm install
npm run dev        # 启动开发服务器（proxy → localhost:8080）
npm run build      # 生产构建
npm run preview    # 预览生产构建
```
