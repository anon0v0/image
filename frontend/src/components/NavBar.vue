<template>
  <!-- 顶部导航栏 -->
  <header class="bg-light-100/80 dark:bg-dark-300/80 backdrop-blur-md border-b border-light-200 dark:border-dark-100 py-3 fixed top-0 left-0 right-0 z-40 transition-all duration-300">
    <div class="container mx-auto px-4">
      <div class="flex justify-between items-center">
        <!-- 左侧Logo -->
        <div class="flex items-center gap-3">
          <router-link to="/" class="flex items-center gap-2 font-semibold text-xl hover:opacity-80 transition-opacity">
            <img src="/logo.png" alt="图床Logo" class="h-10 object-contain">
            <span>Anon图床</span>
          </router-link>
        </div>
        
        <!-- 中间导航菜单 (桌面端) -->
        <nav class="hidden lg:flex items-center gap-1">
          <router-link
            v-for="item in navItems"
            :key="item.path"
            :to="item.path"
            :class="[
              'flex items-center px-4 py-2 rounded-md transition-all duration-200 font-medium',
              isRouteActive(item.path) ? 'bg-primary/10 text-primary' : 'text-secondary hover:bg-light-200 dark:hover:bg-dark-200 hover:text-primary'
            ]"
          >
            <i :class="`ri-${item.icon} mr-2`"></i>
            <span>{{ item.name }}</span>
          </router-link>
        </nav>

        <!-- 右侧操作区 -->
        <div class="flex items-center gap-3">
          <!-- 搜索框 -->
          <div class="hidden md:flex items-center relative">
            <input 
              v-model="searchQuery"
              @keyup.enter="handleSearch"
              type="text" 
              placeholder="搜索图片..." 
              class="w-48 lg:w-56 pl-9 pr-3 py-2 rounded-lg border border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-800 focus:ring-2 focus:ring-primary/50 focus:border-primary outline-none transition-all text-sm"
            />
            <i class="ri-search-line absolute left-3 text-gray-400"></i>
          </div>

          <button 
            @click="toggleTheme"
            class="w-10 h-10 rounded-md bg-light-200 dark:bg-dark-100 text-secondary hover:bg-light-300 dark:hover:bg-dark-200 hover:text-primary transition-all duration-200 flex items-center justify-center"
          >
            <i class="ri-moon-clear-line dark:hidden"></i>
            <i class="ri-sun-line dark:inline-block hidden"></i>
          </button>

          <!-- 用户信息显示（仅登录后显示） -->
          <div v-if="isLoggedIn" class="hidden md:flex items-center gap-2 h-10 px-3 rounded-md bg-light-200 dark:bg-dark-100">
            <div class="w-6 h-6 rounded-full bg-gradient-to-br from-primary to-primary-dark flex items-center justify-center text-white font-bold text-xs">
              {{ userInitial }}
            </div>
            <div class="flex flex-col justify-center leading-tight">
              <span class="text-xs font-medium text-gray-800 dark:text-gray-200">{{ username }}</span>
              <span class="text-[10px] text-gray-500 dark:text-gray-400">{{ isAdmin ? '管理员' : '用户' }}</span>
            </div>
          </div>

          <!-- 登录/退出按钮 - 分开显示 -->
          <button 
            v-if="isLoggedIn"
            @click="handleLogout"
            class="px-4 h-10 rounded-md bg-light-200 dark:bg-dark-100 text-secondary hover:bg-light-300 dark:hover:bg-dark-200 hover:text-primary transition-all duration-200 flex items-center justify-center gap-2"
            title="退出登录"
          >
            <i class="ri-logout-circle-r-line"></i>
            <span class="hidden sm:inline">退出</span>
          </button>
          <router-link 
            v-else
            to="/login"
            class="px-4 h-10 rounded-md bg-primary text-white hover:bg-primary-dark transition-all duration-200 flex items-center justify-center gap-2"
            title="登录"
          >
            <i class="ri-login-circle-line"></i>
            <span class="hidden sm:inline">登录</span>
          </router-link>
          
          <!-- 移动端菜单按钮 -->
          <button 
            class="lg:hidden w-10 h-10 rounded-md bg-light-200 dark:bg-dark-100 text-secondary hover:bg-light-300 dark:hover:bg-dark-200 transition-all duration-200 flex items-center justify-center"
            @click="toggleMobileMenu"
          >
            <i class="ri-menu-line"></i>
          </button>
        </div>
      </div>
    </div>
    
    <!-- 移动端菜单展开 -->
    <div v-if="mobileMenuOpen" class="lg:hidden border-t border-light-200 dark:border-dark-100 mt-3 p-4 bg-light-100 dark:bg-dark-300">
        <!-- 移动端搜索框 -->
        <div class="mb-4 relative">
          <input 
            v-model="searchQuery"
            @keyup.enter="handleSearch"
            type="text" 
            placeholder="搜索图片..." 
            class="w-full pl-9 pr-3 py-2 rounded-lg border border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-800 focus:ring-2 focus:ring-primary/50 focus:border-primary outline-none transition-all text-sm"
          />
          <i class="ri-search-line absolute left-3 top-1/2 -translate-y-1/2 text-gray-400"></i>
        </div>
        
        <nav class="flex flex-col gap-2">
          <router-link
            v-for="item in navItems"
            :key="item.path"
            :to="item.path"
            :class="[
              'flex items-center px-4 py-3 rounded-md transition-all duration-200 font-medium',
              isRouteActive(item.path) ? 'bg-primary/10 text-primary' : 'text-secondary hover:bg-light-200 dark:hover:bg-dark-200 hover:text-primary'
            ]"
            @click="mobileMenuOpen = false"
          >
            <i :class="`ri-${item.icon} mr-3`"></i>
            <span>{{ item.name }}</span>
          </router-link>
        </nav>
    </div>
  </header>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

const mobileMenuOpen = ref(false)
const isLoggedIn = ref(false)
const isAdmin = ref(false)
const username = ref('')
const searchQuery = ref('')

const userInitial = computed(() => {
  if (!username.value) return '游'
  return username.value.charAt(0).toUpperCase()
})

// 搜索处理
const handleSearch = () => {
  if (searchQuery.value.trim()) {
    // 跳转到首页并传递搜索参数
    router.push({ path: '/', query: { search: searchQuery.value.trim() } })
    mobileMenuOpen.value = false
  }
}

// 导航菜单数据
const allNavItems = [
  { path: '/', icon: 'home-4-line', name: '首页', roles: ['admin', 'user', 'guest'] },
  { path: '/random', icon: 'shuffle-line', name: '随机图片', roles: ['admin', 'user', 'guest'] },
  { path: '/upload', icon: 'upload-cloud-2-line', name: '上传', roles: ['admin'] },
  { path: '/ai-function', icon: 'robot-2-line', name: 'AI功能', roles: ['admin'] },
  { path: '/stats', icon: 'pie-chart-2-line', name: '统计', roles: ['admin', 'user'] },
  { path: '/settings', icon: 'settings-3-line', name: '设置', roles: ['admin'] },
]

const navItems = computed(() => {
  return allNavItems.filter(item => {
    // 游客可见
    if (item.roles.includes('guest')) return true
    
    // 未登录，只显示 guest 页面
    if (!isLoggedIn.value) return false
    
    // 管理员可见所有 admin 页面
    if (isAdmin.value && item.roles.includes('admin')) return true
    
    // 普通用户可见 user 页面
    if (!isAdmin.value && item.roles.includes('user')) return true
    
    return false
  })
})

const isRouteActive = (targetPath) => {
  if (targetPath === '/') {
    return route.path === '/'
  }
  return route.path.startsWith(targetPath)
}

const toggleMobileMenu = () => {
    mobileMenuOpen.value = !mobileMenuOpen.value
}

// 从 API 检查登录状态
const checkLoginStatus = async () => {
  try {
    const response = await fetch('/api/user/status', {
      credentials: 'include' // 包含 cookie
    })
    const data = await response.json()
    
    if (data.code === 200 && data.data.logged_in) {
      isLoggedIn.value = true
      isAdmin.value = data.data.role === 'admin'
      username.value = data.data.username || 'User'
      
      // 同步到 localStorage（用于其他页面）
      if (typeof localStorage !== 'undefined') {
        localStorage.setItem('userInfo', JSON.stringify({
          username: data.data.username || 'user',
          role: data.data.role || 'user',
          user_id: data.data.user_id
        }))
      }
    } else {
      isLoggedIn.value = false
      isAdmin.value = false
      username.value = ''
      
      // 清除 localStorage
      if (typeof localStorage !== 'undefined') {
        localStorage.removeItem('userInfo')
      }
    }
  } catch (error) {
    console.error('检查登录状态失败:', error)
    isLoggedIn.value = false
    isAdmin.value = false
    username.value = ''
  }
}

// 主题切换功能
const storageKey = 'theme-preference'

const detectUserThemePreference = () => {
  if (typeof localStorage !== 'undefined' && localStorage.getItem(storageKey)) {
    return localStorage.getItem(storageKey)
  }
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
}

const applyTheme = (theme) => {
  const htmlElement = document.documentElement
  if (theme === 'dark') {
    htmlElement.classList.add('dark')
  } else {
    htmlElement.classList.remove('dark')
  }
  if (typeof localStorage !== 'undefined') {
    localStorage.setItem(storageKey, theme)
  }
}

const toggleTheme = () => {
  const currentTheme = localStorage.getItem(storageKey) || 'light'
  const newTheme = currentTheme === 'dark' ? 'light' : 'dark'
  applyTheme(newTheme)
}

// 登出功能
const handleLogout = async () => {
  try {
    await fetch('/api/logout', {
      method: 'POST',
      credentials: 'include'
    })
    
    // 清除本地状态
    if (typeof localStorage !== 'undefined') {
      localStorage.removeItem('userInfo')
    }
    isLoggedIn.value = false
    isAdmin.value = false
    username.value = ''
    
    // 跳转到首页
    router.push('/').catch(() => {})
  } catch (error) {
    console.error('登出失败:', error)
  }
}

// 组件挂载时初始化
onMounted(() => {
  // 初始化主题
  const initialTheme = detectUserThemePreference()
  applyTheme(initialTheme)
  
  // 检查登录状态
  checkLoginStatus()
  
  // 监听 storage 变化（用于多标签页同步）
  window.addEventListener('storage', checkLoginStatus)
})
</script>