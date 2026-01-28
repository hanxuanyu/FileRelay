import { createRouter, createWebHistory } from 'vue-router'

// 主页面组件
const MainPage = () => import('@/views/MainPage.vue')

// 管理页面组件
const AdminLogin = () => import('@/views/admin/AdminLogin.vue')
const AdminDashboard = () => import('@/views/admin/AdminDashboard.vue')
const BatchManagement = () => import('@/views/admin/BatchManagement.vue')
const TokenManagement = () => import('@/views/admin/TokenManagement.vue')
const ConfigManagement = () => import('@/views/admin/ConfigManagement.vue')

// 检查管理员权限
function requireAuth(_to: any, _from: any, next: any) {
  const token = localStorage.getItem('admin_token')
  if (token) {
    next()
  } else {
    next('/admin/login')
  }
}

const routes = [
  {
    path: '/',
    name: 'HomePage',
    component: MainPage
  },
  {
    path: '/upload',
    name: 'UploadPage', 
    component: MainPage
  },
  {
    path: '/:code',
    name: 'HomePageWithCode',
    component: MainPage,
    props: true
  },
  // 管理页面路由
  {
    path: '/admin/login',
    name: 'AdminLogin',
    component: AdminLogin
  },
  {
    path: '/admin',
    name: 'AdminDashboard',
    component: AdminDashboard,
    beforeEnter: requireAuth
  },
  {
    path: '/admin/batches',
    name: 'BatchManagement', 
    component: BatchManagement,
    beforeEnter: requireAuth
  },
  {
    path: '/admin/tokens',
    name: 'TokenManagement',
    component: TokenManagement,
    beforeEnter: requireAuth
  },
  {
    path: '/admin/config',
    name: 'ConfigManagement',
    component: ConfigManagement,
    beforeEnter: requireAuth
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router