# Zine API 文档

Base URL: `http://localhost:8080/api`

所有接口返回统一格式：

```json
{ "code": 200, "message": "success", "data": {} }
```

---

## 公开接口（无需认证）

### 文章

#### 文章列表

```
GET /api/articles
```

Query 参数：

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码，默认 1 |
| page_size | int | 否 | 每页条数，默认 10，最大 100 |
| category_id | int | 否 | 按分类筛选 |
| tag_id | int | 否 | 按标签筛选 |
| keyword | string | 否 | 标题模糊搜索 |
| series | string | 否 | 按系列筛选 |

> 公开接口只返回 `status=1`（已发布）的文章。

响应：

```json
{
  "code": 200,
  "data": {
    "list": [
      {
        "id": 1,
        "title": "Hello World",
        "slug": "hello-world",
        "content_md": "# Hello\n\n...",
        "content_html": "<h1>Hello</h1>\n...",
        "excerpt": "",
        "cover": "",
        "status": 1,
        "view_count": 42,
        "is_top": 0,
        "category_id": 1,
        "category": { "id": 1, "name": "Tech", "slug": "tech" },
        "tags": [
          { "id": 1, "name": "Go", "slug": "go" }
        ],
        "series": "",
        "series_order": 0,
        "created_at": "2025-06-27T10:00:00Z",
        "updated_at": "2025-06-27T10:00:00Z"
      }
    ],
    "total": 25,
    "page": 1,
    "page_size": 10
  }
}
```

#### 文章详情

```
GET /api/articles/:slug
```

响应：

```json
{
  "code": 200,
  "data": {
    "article": {
      "id": 1,
      "title": "Hello World",
      "slug": "hello-world",
      "content_md": "# Hello\n\n**Bold**",
      "content_html": "<h1>Hello</h1>\n<p><strong>Bold</strong></p>",
      "excerpt": "",
      "cover": "",
      "status": 1,
      "view_count": 43,
      "is_top": 0,
      "category": { "id": 1, "name": "Tech", "slug": "tech" },
      "tags": [...],
      "series": "building-a-compiler",
      "series_order": 2,
      "created_at": "2025-06-27T10:00:00Z",
      "updated_at": "2025-06-27T10:00:00Z"
    },
    "prev_in_series": { "slug": "part-1", "title": "Part 1" },
    "next_in_series": { "slug": "part-3", "title": "Part 3" }
  }
}
```

> - 每次访问详情会自增 `view_count`。
> - `prev_in_series` / `next_in_series` 仅当文章属于某个系列（`series` 非空且 `series_order > 0`）且存在前后篇时返回。

---

### 分类

```
GET /api/categories
```

```json
{
  "code": 200,
  "data": [
    { "id": 1, "name": "Tech", "slug": "tech", "description": "技术文章" },
    { "id": 2, "name": "Life", "slug": "life", "description": "" }
  ]
}
```

---

### 标签

```
GET /api/tags
```

```json
{
  "code": 200,
  "data": [
    { "id": 1, "name": "Go", "slug": "go" },
    { "id": 2, "name": "Gin", "slug": "gin" }
  ]
}
```

---

### 系列

```
GET /api/series
```

```json
{
  "code": 200,
  "data": [
    {
      "name": "building-a-compiler",
      "count": 5,
      "latest_slug": "building-a-compiler-part-5"
    }
  ]
}
```

> 只返回已发布文章数 ≥1 的系列，按文章数降序排列。

---

### 评论

#### 获取文章评论

```
GET /api/comments?article_id=1
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| article_id | int | 是 | 文章 ID |
| page | int | 否 | 页码，默认 1 |
| page_size | int | 否 | 每页条数，默认 20 |

```json
{
  "code": 200,
  "data": {
    "list": [
      {
        "id": 1,
        "article_id": 1,
        "parent_id": null,
        "author": "小王",
        "content": "写得真好",
        "status": 1,
        "created_at": "2025-06-27T12:00:00Z",
        "replies": [
          {
            "id": 2,
            "parent_id": 1,
            "author": "站长",
            "content": "谢谢！",
            "created_at": "2025-06-27T12:30:00Z"
          }
        ]
      }
    ],
    "total": 10,
    "page": 1,
    "page_size": 20
  }
}
```

#### 发表评论

```
POST /api/comments
```

Body (JSON):

```json
{
  "article_id": 1,
  "parent_id": null,
  "author": "小王",
  "email": "user@example.com",
  "content": "写得真好"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| article_id | int | 是 | 文章 ID |
| parent_id | int | 否 | 回复某条评论的 ID，空表示顶级评论 |
| author | string | 是 | 昵称 |
| email | string | 否 | 邮箱 |
| content | string | 是 | 评论内容 |

---

## 管理接口（需要 JWT）

前缀：`/api/admin`

请求头：`Authorization: Bearer <token>`

### 登录

```
POST /api/admin/login
```

Body:

```json
{ "username": "admin", "password": "admin123" }
```

响应：

```json
{
  "code": 200,
  "data": { "token": "eyJhbGciOi..." }
}
```

---

### 文章管理

#### 列表（含草稿）

```
GET /api/admin/articles
```

参数同公开接口的 `/api/articles`，但不过滤 `status`，可查看草稿。

#### 创建文章

```
POST /api/admin/articles
```

```json
{
  "title": "新文章",
  "slug": "",
  "content_md": "# Hello",
  "excerpt": "",
  "cover": "",
  "status": 1,
  "is_top": 0,
  "category_id": 1,
  "tag_ids": [1, 2],
  "series": "",
  "series_order": 0
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| title | string | 是 | 文章标题 |
| slug | string | 否 | 留空则自动从标题生成 |
| content_md | string | 是 | Markdown 原文 |
| excerpt | string | 否 | 摘要 |
| cover | string | 否 | 封面图 URL |
| status | int | 否 | 0=草稿, 1=发布 |
| is_top | int | 否 | 0=普通, 1=置顶 |
| category_id | int | 否 | 分类 ID |
| tag_ids | []int | 否 | 标签 ID 数组 |
| series | string | 否 | 系列标识（如 `building-a-compiler`） |
| series_order | int | 否 | 在系列中的序号（1-based） |

> `content_html` 由服务端自动从 `content_md` 渲染生成。

#### 更新文章

```
PUT /api/admin/articles/:id
```

Body 格式同创建。

#### 删除文章

```
DELETE /api/admin/articles/:id
```

---

### 分类管理

```json
POST   /api/admin/categories    { "name": "Tech", "slug": "", "description": "" }
PUT    /api/admin/categories/:id
DELETE /api/admin/categories/:id
```

---

### 标签管理

```json
POST   /api/admin/tags           { "name": "Go", "slug": "" }
PUT    /api/admin/tags/:id
DELETE /api/admin/tags/:id
```

---

### 评论管理

#### 列表

```
GET /api/admin/comments?page=1&page_size=20&status=0
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 默认 1 |
| page_size | int | 否 | 默认 20 |
| status | int | 否 | 筛选状态：0=待审核, 1=通过, 2=垃圾；不传则全部 |

#### 审核评论

```
PUT /api/admin/comments/:id
```

```json
{ "status": 1 }
```

| status | 含义 |
|--------|------|
| 0 | 待审核 |
| 1 | 通过 |
| 2 | 垃圾 |

#### 删除评论

```
DELETE /api/admin/comments/:id
```

---

### 文件上传

```
POST /api/admin/upload
```

Form-Data: `file` 字段

限制：`.jpg` `.jpeg` `.png` `.gif` `.webp` `.svg`，最大 10MB

响应：

```json
{
  "code": 200,
  "data": { "url": "/uploads/2025-06-27/a1b2c3.jpg" }
}
```

文件存入 `uploads/` 目录下按日期的子目录中，可通过 `/uploads/2025-06-27/a1b2c3.jpg` 直接访问。

---

### 个人资料

```
GET /api/admin/profile
```

```json
{
  "code": 200,
  "data": {
    "id": 1,
    "username": "admin",
    "nickname": "Admin",
    "created_at": "2025-01-01T00:00:00Z"
  }
}
```

---

## 错误码

| code | 说明 |
|------|------|
| 200 | 成功 |
| 400 | 参数错误 |
| 401 | 未认证 / Token 过期 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

---

## 静态文件

`/uploads/` 目录下的文件可直接通过 URL 访问，Gin 内置静态文件路由。
