<template>
  <div id="app">
      <!-- 背景网格 -->
    <div class="fixed inset-0 bg-grid opacity-70 dark:opacity-50"></div>
    
    <!-- 装饰性背景元素 -->
    <div class="fixed top-20 -left-20 w-64 h-64 bg-primary/10 dark:bg-primary/20 rounded-full decorative-blur animate-pulse-slow"></div>
    <div class="fixed bottom-20 -right-20 w-80 h-80 bg-primary-dark/10 dark:bg-primary-dark/20 rounded-full decorative-blur animate-pulse-slow" style="animation-delay: 1s;"></div>
    
    <Navbar />
    
    <!-- 主内容区 -->
    <main class="flex-grow pt-24 pb-16 px-4 relative z-10">
        <router-view></router-view>
    </main>
  </div>
</template>

<script setup>
import Navbar from "@/components/NavBar.vue";
import { onMounted } from 'vue';

onMounted(async () => {
    try {
        // 每次加载应用时，检查后端会话状态
        const response = await fetch('/api/user/status');
        if (response.ok) {
            const result = await response.json();
            if (result.code === 200 && result.data.logged_in) {
                // 后端已登录，同步到前端 localStorage
                const userInfo = {
                    username: result.data.username || 'User',
                    role: result.data.role || 'user'
                };
                localStorage.setItem('userInfo', JSON.stringify(userInfo));
                // 触发 storage 事件以便其他组件更新
                window.dispatchEvent(new Event('storage'));
            } else {
                // 后端未登录，清除前端状态
                // localStorage.removeItem('userInfo');
                // window.dispatchEvent(new Event('storage'));
            }
        }
    } catch (error) {
        console.error('Session sync failed:', error);
    }
});
</script>