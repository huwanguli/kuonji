import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
})

api.interceptors.response.use(
  (res) => res.data,
  (err) => Promise.reject(err)
)

export function authHeader() {
  const token = localStorage.getItem('token')
  return token ? { Authorization: `Bearer ${token}` } : {}
}

export const articles = {
  list(params = {}) {
    return api.get('/articles', { params: { page: 1, page_size: 10, ...params } })
  },
  getBySlug(slug) {
    return api.get(`/articles/${slug}`)
  },
  adminList(params = {}) {
    return api.get('/admin/articles', { params, headers: authHeader() })
  },
  create(data) {
    return api.post('/admin/articles', data, { headers: authHeader() })
  },
  update(id, data) {
    return api.put(`/admin/articles/${id}`, data, { headers: authHeader() })
  },
  delete(id) {
    return api.delete(`/admin/articles/${id}`, { headers: authHeader() })
  },
}

export const series = {
  list() {
    return api.get('/series')
  },
}

export const announcements = {
  list() {
    return api.get('/announcements')
  },
}

export const categories = {
  list() {
    return api.get('/categories')
  },
  adminList() {
    return api.get('/admin/categories', { headers: authHeader() })
  },
  create(data) {
    return api.post('/admin/categories', data, { headers: authHeader() })
  },
  update(id, data) {
    return api.put(`/admin/categories/${id}`, data, { headers: authHeader() })
  },
  delete(id) {
    return api.delete(`/admin/categories/${id}`, { headers: authHeader() })
  },
}

export const tags = {
  list() {
    return api.get('/tags')
  },
  create(data) {
    return api.post('/admin/tags', data, { headers: authHeader() })
  },
  update(id, data) {
    return api.put(`/admin/tags/${id}`, data, { headers: authHeader() })
  },
  delete(id) {
    return api.delete(`/admin/tags/${id}`, { headers: authHeader() })
  },
}

export const comments = {
  list(articleId, params = {}) {
    return api.get('/comments', { params: { article_id: articleId, ...params } })
  },
  create(data) {
    return api.post('/comments', data)
  },
  adminList(params = {}) {
    return api.get('/admin/comments', { params, headers: authHeader() })
  },
  updateStatus(id, status) {
    return api.put(`/admin/comments/${id}`, { status }, { headers: authHeader() })
  },
  delete(id) {
    return api.delete(`/admin/comments/${id}`, { headers: authHeader() })
  },
}

export const auth = {
  login(username, password) {
    return api.post('/admin/login', { username, password })
  },
  profile() {
    return api.get('/admin/profile', { headers: authHeader() })
  },
}

export const upload = {
  image(file) {
    const fd = new FormData()
    fd.append('file', file)
    return api.post('/admin/upload', fd, { headers: { ...authHeader(), 'Content-Type': 'multipart/form-data' } })
  },
}
