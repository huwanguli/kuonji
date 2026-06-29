# Kuonji

个人博客系统，Go + Gin 后端，Vue 3 前端。

[kuonji.xyz](https://kuonji.xyz)

## 技术栈

| 层 | 技术 |
|------|------|
| 后端框架 | Go 1.21+ / Gin |
| ORM | GORM v2 |
| 数据库 | MySQL 8.0+ |
| 配置 | Viper (yaml + 环境变量) |
| 日志 | logrus + lumberjack (按日切割) |
| 认证 | JWT (golang-jwt) |
| Markdown | goldmark |
| 前端框架 | Vue 3 + Vite |
| 路由 | Vue Router 4 |
| HTTP 客户端 | Axios |
| Markdown (前端) | marked |

## 快速启动

### 前置条件

- Go 1.21+
- MySQL 8.0+
- Node.js 18+

### 1. 数据库

创建数据库：

```sql
CREATE DATABASE zblog CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 2. 后端

```bash
cd backend
# 复制配置模板
cp config/config.example.yaml config/config.yaml
# 修改 config/config.yaml 中的数据库连接信息
go run main.go
```

首次启动自动建表并创建默认管理员 `admin` / `admin123`。

数据库密码等敏感配置可通过环境变量覆盖，见[生产部署](#生产部署)。

### 3. 前端

```bash
cd frontend
npm install
npm run dev
```

开发模式下前端 `/api` 和 `/uploads` 请求自动代理到 `localhost:8080`。

### 4. 访问

- 博客首页：`http://localhost:5173`
- 管理后台：`http://localhost:5173/admin`

## 项目结构

```
kuonji/
├── backend/
│   ├── main.go                    # 入口
│   ├── config/
│   │   ├── config.example.yaml    # 配置模板
│   │   └── config.yaml            # 配置（已 gitignore）
│   ├── Dockerfile                 # 后端镜像
│   ├── API.md                     # API 文档
│   └── internal/
│       ├── config/config.go       # Viper 配置加载
│       ├── logger/logger.go       # 日志初始化
│       ├── model/                 # GORM 模型 (User/Article/Category/Tag/Comment)
│       ├── repository/            # 数据访问层
│       ├── service/               # 业务逻辑层
│       ├── handler/               # HTTP 处理层
│       ├── middleware/             # JWT/CORS/Logger/Recovery
│       ├── router/router.go       # 路由注册
│       ├── dto/                   # 请求/响应结构体
│       └── utils/                 # 工具 (slug/MD/JWT)
├── frontend/
│   ├── Dockerfile                 # 前端镜像（Caddy 自动 HTTPS）
│   ├── Caddyfile                  # Caddy 配置（SPA + 反代 + 自动证书）
│   └── src/
│       ├── main.js                # 入口
│       ├── App.vue                # 根组件
│       ├── assets/main.css        # Design tokens + 全局样式
│       ├── api/index.js           # API 封装
│       ├── router/index.js        # 路由配置
│       ├── utils/html.js          # TOC 工具
│       ├── components/            # 组件
│       │   ├── NavBar.vue         # 导航栏（访客/管理员自适应）
│       │   ├── Footer.vue         # 页脚
│       │   ├── Sidebar.vue        # 侧边栏（分类/系列/标签/目录）
│       │   ├── ArticleCard.vue    # 文章卡片
│       │   ├── ArticleList.vue    # 文章列表 + 分页
│       │   ├── CommentList.vue    # 评论列表
│       │   └── CommentForm.vue    # 评论表单
│       └── views/                 # 页面
│           ├── HomeView.vue       # 首页
│           ├── ArticleView.vue    # 文章详情
│           ├── AdminView.vue      # 管理面板
│           ├── AdminArticles.vue  # 文章管理列表
│           ├── AdminComments.vue  # 评论管理
│           └── EditorView.vue     # Markdown 编辑器
├── docker-compose.yml             # 容器编排
└── .env.example                   # 环境变量模板
```

## 架构

```
Handler (参数校验、响应格式化)
  ↓
Service (业务逻辑、MD 渲染、事务)
  ↓
Repository (GORM 数据访问、分页)
  ↓
Model (GORM 模型定义)
```

## 路由一览

### 公开 API (`/api`)

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/articles` | 文章列表（分页/分类/标签/系列筛选） |
| GET | `/api/articles/:slug` | 文章详情（阅读量 +1，含系列前后篇） |
| GET | `/api/announcements` | 公告列表 |
| GET | `/api/series` | 系列列表 |
| GET | `/api/categories` | 分类列表 |
| GET | `/api/tags` | 标签列表 |
| GET | `/api/comments` | 文章评论 |
| POST | `/api/comments` | 发表评论 |

### 管理 API (`/api/admin`)，需 JWT

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/admin/login` | 登录 |
| GET | `/api/admin/profile` | 当前用户 |
| GET/POST/PUT/DELETE | `/api/admin/articles` | 文章 CRUD |
| POST/PUT/DELETE | `/api/admin/categories` | 分类管理 |
| POST/PUT/DELETE | `/api/admin/tags` | 标签管理 |
| GET/PUT/DELETE | `/api/admin/comments` | 评论审核 |
| POST | `/api/admin/upload` | 图片上传 |

详细请求/响应格式见 `backend/API.md`。

## 配置

`backend/config/config.yaml`（从 `config.example.yaml` 复制）：

```yaml
server:
  port: 8080
  mode: debug          # debug | release

database:
  host: 127.0.0.1
  port: 3306
  user: root
  password: ""
  dbname: zblog

jwt:
  secret: ""
  expire: 24h

log:
  level: debug
  file: /logs/zblog.log
  max_age: 30
  max_size: 100

upload:
  path: uploads/
  max_size: 10
  allowed_exts: [".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg"]
```

敏感字段（password、secret）留空，通过环境变量注入：

| 环境变量 | 对应字段 |
|----------|----------|
| `ZBLOG_DATABASE_HOST` | database.host |
| `ZBLOG_DATABASE_PASSWORD` | database.password |
| `ZBLOG_JWT_SECRET` | jwt.secret |
| `ZBLOG_SERVER_MODE` | server.mode |
| `ZBLOG_LOG_LEVEL` | log.level |

## 前端页面

| 路径 | 名称 | 说明 |
|------|------|------|
| `/` | 首页 | 随机标语 + 文章列表（支持分类/标签/系列筛选）+ 侧边栏 |
| `/article/:slug` | 文章详情 | 封面图 + 正文 + 目录 + 评论 |
| `/admin` | 管理面板 | 登录 + 分类/标签管理 |
| `/admin/articles` | 文章管理 | 列表/筛选/编辑/删除 + 文章评论管理 |
| `/admin/comments` | 评论管理 | 全部评论审核/删除/恢复 |
| `/editor` | 写文章 | Markdown 预览编辑器 + 封面/图片上传 |
| `/editor/:id` | 编辑文章 | 同上 |

导航栏在 `/admin*` 和 `/editor*` 路径下自动切换为管理导航，其余路径显示访客导航。

## 如何添加新功能

### 后端

1. `internal/model/` 添加 GORM 模型
2. `internal/repository/` 添加数据访问接口和实现
3. `internal/service/` 添加业务逻辑
4. `internal/handler/` 添加 HTTP 处理
5. `internal/router/router.go` 注册路由
6. `main.go` 中组装依赖

### 前端

1. `src/views/` 添加新页面
2. `src/router/index.js` 注册路由
3. `src/api/index.js` 添加 API 调用函数
4. 如需新组件：`src/components/` 添加

## 运行测试

```bash
# 后端（使用 SQLite 内存数据库，无需 MySQL）
cd backend && go test ./...

# 前端
cd frontend && npm run build
```

## 设计 Token

| 令牌 | 值 |
|------|-----|
| 背景色 | `#fdfbf7`（纸张白） |
| 正文色 | `#1a1a1a`（墨黑） |
| 强调色 | `#e63946`（朱红） |
| 标题色 | `#2b2d42`（深靛） |
| 标题字体 | Playfair Display |
| 正文字体 | Inter |
| 代码字体 | JetBrains Mono |

## 生产部署

### Docker Compose（推荐）

```bash
# 1. 配置环境变量
cp .env.example .env
# 编辑 .env 填入数据库密码和 JWT 密钥

# 2. 构建并启动
docker compose up -d --build
```

架构：`Caddy :80/:443` → 自动 Let's Encrypt HTTPS + 静态前端 + 反代 `/api` `/uploads` → `backend:8080` → MySQL。

### 手动部署

```bash
# 构建前端
cd frontend && npm run build     # 输出到 dist/

# 构建后端
cd backend && go build -o kuonji .

# 部署：Caddy 将 / 指向 frontend/dist/，/api /uploads 反代到后端 8080，自动获取 Let's Encrypt 证书
```
