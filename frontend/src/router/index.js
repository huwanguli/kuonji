import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'home',
    component: () => import('../views/HomeView.vue'),
  },
  {
    path: '/article/:slug',
    name: 'article',
    component: () => import('../views/ArticleView.vue'),
  },
  {
    path: '/admin',
    name: 'admin',
    component: () => import('../views/AdminView.vue'),
  },
  {
    path: '/admin/articles',
    name: 'admin-articles',
    component: () => import('../views/AdminArticles.vue'),
  },
  {
    path: '/admin/articles/view/:id',
    name: 'admin-article-view',
    component: () => import('../views/AdminArticleView.vue'),
  },
  {
    path: '/admin/comments',
    name: 'admin-comments',
    component: () => import('../views/AdminComments.vue'),
  },
  {
    path: '/editor',
    name: 'editor',
    component: () => import('../views/EditorView.vue'),
  },
  {
    path: '/editor/:id',
    name: 'editor-edit',
    component: () => import('../views/EditorView.vue'),
  },
  {
    path: '/series/:name',
    name: 'series',
    component: () => import('../views/SeriesView.vue'),
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/',
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) return savedPosition
    return { top: 0 }
  },
})

export default router
