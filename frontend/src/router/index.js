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
