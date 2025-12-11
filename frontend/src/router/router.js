import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: {
      title: '登录',
      public: true
    }
  },
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/Home.vue'),
    meta: {
      title: 'Anon图床',
      public: true
    }
  },
  {
    path: '/upload',
    name: 'Upload',
    component: () => import('@/views/Upload.vue'),
    meta: {
      title: '上传图片'
    }
  },
  {
    path: '/gallery',
    redirect: '/'
  },
  {
    path: '/random',
    name: 'Random',
    component: () => import('@/views/Random.vue'),
    meta: {
      title: '随机图片',
      public: true
    }
  },
  {
    path: '/stats',
    name: 'Stats',
    component: () => import('@/views/Stats.vue'),
    meta: {
      title: '系统统计'
    }
  },
  {
    path: '/settings',
    name: 'Settings',
    component: () => import('@/views/Settings.vue'),
    meta: {
      title: '系统设置'
    }
  },
  {
    path: '/ai-function',
    name: 'AIFunction',
    component: () => import('@/views/AIFunction.vue'),
    meta: {
      title: 'AI功能'
    }
  }
]

// 创建路由实例
const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

// 全局前置守卫 - 处理页面标题和登录验证
router.beforeEach(async (to, from, next) => {
  document.title = to.meta.title || 'Anon图床';

  // 如果是公开页面，直接放行
  if (to.meta.public) {
    return next();
  }

  // 尝试从 localStorage 获取登录状态
  const userInfoStr = localStorage.getItem('userInfo');
  if (userInfoStr) {
    return next();
  }

  // 如果没有本地缓存，且访问的是受保护页面，则必须验证
  try {
    const response = await fetch('/api/user/status', {
      method: 'GET'
    });
    const result = await response.json()
    if (result.code === 200 && result.data.logged_in == true) {
      localStorage.setItem('userInfo', JSON.stringify({
        username: result.data.username,
        role: result.data.role,
        user_id: result.data.user_id
      }));
      return next();
    }
    next('/login');
  } catch (error) {
    console.error('验证登录状态失败:', error);
    next('/login');
  }
});

export default router