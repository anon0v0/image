<template>
  <div class="fixed inset-0 flex items-center justify-center overflow-hidden transition-colors duration-300 z-40">
    <!-- 动态渐变背景 -->
    <div class="absolute inset-0 bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900">
      <!-- 动态光效 -->
      <div class="absolute top-0 left-0 w-full h-full overflow-hidden">
        <div class="absolute -top-[40%] -left-[20%] w-[80%] h-[80%] rounded-full bg-gradient-to-r from-blue-500/30 to-purple-500/30 blur-[120px] animate-blob"></div>
        <div class="absolute top-[20%] -right-[20%] w-[60%] h-[60%] rounded-full bg-gradient-to-r from-pink-500/30 to-orange-500/30 blur-[100px] animate-blob animation-delay-2000"></div>
        <div class="absolute -bottom-[30%] left-[20%] w-[70%] h-[70%] rounded-full bg-gradient-to-r from-cyan-500/30 to-blue-500/30 blur-[100px] animate-blob animation-delay-4000"></div>
      </div>
      
      <!-- 星空粒子效果 -->
      <div class="stars-container absolute inset-0 overflow-hidden">
        <div v-for="i in 50" :key="i" 
          class="absolute w-1 h-1 bg-white rounded-full animate-twinkle"
          :style="{
            left: Math.random() * 100 + '%',
            top: Math.random() * 100 + '%',
            animationDelay: Math.random() * 5 + 's',
            opacity: Math.random() * 0.5 + 0.3
          }"
        ></div>
      </div>
      
      <!-- 网格背景 -->
      <div class="absolute inset-0 bg-[linear-gradient(rgba(255,255,255,0.03)_1px,transparent_1px),linear-gradient(90deg,rgba(255,255,255,0.03)_1px,transparent_1px)] bg-[size:50px_50px]"></div>
    </div>

    <div class="w-full max-w-md px-6 relative z-10">
      <!-- 主卡片 -->
      <div class="login-card relative bg-white/10 backdrop-blur-xl rounded-3xl shadow-2xl p-6 md:p-8 border border-white/20 overflow-hidden">
        <!-- 卡片光效边框 -->
        <div class="absolute inset-0 rounded-3xl bg-gradient-to-r from-blue-500/20 via-purple-500/20 to-pink-500/20 opacity-0 hover:opacity-100 transition-opacity duration-500"></div>
        
        <!-- 标题区 -->
        <div class="text-center mb-6 relative z-10">
          <router-link to="/" class="inline-block mb-4 group">
            <div class="relative">
              <!-- Logo光环 -->
              <div class="absolute inset-0 bg-gradient-to-r from-blue-500 to-purple-500 rounded-xl blur-lg opacity-50 group-hover:opacity-80 transition-opacity animate-pulse"></div>
              <!-- Logo主体 -->
              <div class="relative w-14 h-14 bg-gradient-to-br from-blue-500 via-purple-500 to-pink-500 rounded-xl flex items-center justify-center text-white text-2xl shadow-2xl group-hover:scale-110 transition-transform duration-300 ring-2 ring-white/30">
                <i class="ri-gallery-fill"></i>
              </div>
            </div>
          </router-link>
          
          <h1 class="text-2xl font-black text-white mb-2 tracking-tight">
            <span class="bg-gradient-to-r from-white via-blue-200 to-purple-200 bg-clip-text text-transparent">欢迎回来</span>
          </h1>
          <p class="text-white/60 text-xs font-medium">登录以管理您的精彩图库</p>
        </div>

        <!-- 登录方式 -->
        <div class="space-y-4 relative z-10">
          <!-- GitHub 登录 -->
          <button 
            @click="handleGithubLogin" 
            :disabled="isLoading"
            class="github-btn w-full flex items-center justify-center gap-2 bg-white/10 hover:bg-white/20 backdrop-blur-sm text-white py-3 rounded-xl font-bold transition-all border border-white/20 hover:border-white/40 shadow-lg hover:shadow-xl transform hover:-translate-y-0.5 disabled:opacity-70 disabled:cursor-not-allowed group"
          >
            <i class="ri-github-fill text-xl"></i>
            <span>{{ isLoading && loginType === 'github' ? '正在连接...' : '使用 GitHub 登录' }}</span>
          </button>

          <!-- 分隔线 -->
          <div class="relative flex py-3 items-center">
            <div class="flex-grow border-t border-white/20"></div>
            <div class="mx-4 px-4 py-1 bg-white/10 backdrop-blur-sm rounded-full">
              <span class="text-white/60 text-xs uppercase tracking-widest font-bold">管理员登录</span>
            </div>
            <div class="flex-grow border-t border-white/20"></div>
          </div>

          <!-- 管理员登录表单 -->
          <form @submit.prevent="handleAdminLogin" class="space-y-3">
            <div class="space-y-2">
              <label class="text-xs font-bold text-white/70 ml-1 uppercase tracking-wider">用户名</label>
              <div class="relative group">
                <input 
                  v-model="username" 
                  type="text" 
                  required
                  class="w-full pl-11 pr-4 py-3 rounded-xl border border-white/20 bg-white/5 backdrop-blur-sm focus:ring-2 focus:ring-purple-500/50 focus:border-purple-500/50 focus:bg-white/10 outline-none transition-all text-white placeholder-white/40 font-medium text-sm"
                  placeholder="请输入管理员账号"
                />
                <div class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 bg-gradient-to-br from-blue-500 to-purple-500 rounded flex items-center justify-center">
                  <i class="ri-user-line text-white text-sm"></i>
                </div>
              </div>
            </div>
            
            <div class="space-y-2">
              <label class="text-xs font-bold text-white/70 ml-1 uppercase tracking-wider">密码</label>
              <div class="relative group">
                <input 
                  v-model="password" 
                  type="password" 
                  required
                  class="w-full pl-11 pr-4 py-3 rounded-xl border border-white/20 bg-white/5 backdrop-blur-sm focus:ring-2 focus:ring-purple-500/50 focus:border-purple-500/50 focus:bg-white/10 outline-none transition-all text-white placeholder-white/40 font-medium text-sm"
                  placeholder="请输入密码"
                />
                <div class="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 bg-gradient-to-br from-purple-500 to-pink-500 rounded flex items-center justify-center">
                  <i class="ri-lock-line text-white text-sm"></i>
                </div>
              </div>
            </div>

            <button 
              type="submit" 
              :disabled="isLoading"
              class="login-btn w-full relative overflow-hidden bg-gradient-to-r from-blue-500 via-purple-500 to-pink-500 hover:from-blue-600 hover:via-purple-600 hover:to-pink-600 text-white py-3 rounded-xl font-bold transition-all shadow-xl shadow-purple-500/30 hover:shadow-purple-500/50 transform hover:-translate-y-0.5 disabled:opacity-70 disabled:cursor-not-allowed mt-1"
            >
              <!-- 按钮光效 -->
              <div class="absolute inset-0 bg-gradient-to-r from-transparent via-white/20 to-transparent -translate-x-full hover:translate-x-full transition-transform duration-1000"></div>
              <span v-if="isLoading && loginType === 'admin'" class="relative flex items-center justify-center gap-3">
                <div class="w-6 h-6 border-3 border-white/30 border-t-white rounded-full animate-spin"></div>
                登录中...
              </span>
              <span v-else class="relative flex items-center justify-center gap-2">
                <i class="ri-login-box-line text-xl"></i>
                立即登录
              </span>
            </button>
          </form>
        </div>

        <!-- 底部 -->
        <div class="mt-4 text-center relative z-10">
          <p class="text-xs text-white/40 font-medium flex items-center justify-center gap-1">
            <i class="ri-information-line"></i>
            未注册？目前仅支持 GitHub 快捷登录
          </p>
        </div>
        
        <!-- 装饰元素 -->
        <div class="absolute top-4 right-4 w-20 h-20 bg-gradient-to-br from-blue-500/20 to-purple-500/20 rounded-full blur-2xl"></div>
        <div class="absolute bottom-4 left-4 w-16 h-16 bg-gradient-to-br from-pink-500/20 to-orange-500/20 rounded-full blur-2xl"></div>
      </div>
      
      <!-- 版权信息 -->
      <div class="text-center mt-4">
        <p class="text-white/30 text-xs font-medium">© 2026 Anon图床 · 安全 · 快速 · 免费</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import Message from '@/utils/message.js'

const router = useRouter()
const username = ref('')
const password = ref('')
const isLoading = ref(false)
const loginType = ref('') // 'github' or 'admin'

const handleGithubLogin = () => {
    if (isLoading.value) return
    isLoading.value = true
    loginType.value = 'github'
    
    // 立即跳转，移除不必要的延迟
    window.location.href = '/api/auth/github/login'
}

const handleAdminLogin = async () => {
    if (isLoading.value) return
    if (!username.value || !password.value) {
        Message.warning('请输入用户名和密码')
        return
    }

    isLoading.value = true
    loginType.value = 'admin'

    try {
        const response = await fetch('/api/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                username: username.value,
                password: password.value
            })
        })

        const data = await response.json()

        if (response.ok && data.code === 200) {
            Message.success('登录成功')
            // 保存用户信息到本地
             if (typeof localStorage !== 'undefined') {
                localStorage.setItem('userInfo', JSON.stringify({
                    username: data.user?.username || username.value,
                    role: 'admin',
                    user_id: data.user?.id
                }))
            }
            // 触发 storage 事件更新状态
            window.dispatchEvent(new Event('storage'))
            
            router.push('/')
        } else {
            Message.error(data.message || '登录失败')
        }
    } catch (error) {
        console.error('Login error:', error)
        Message.error('网络错误，请稍后重试')
    } finally {
        isLoading.value = false
        loginType.value = ''
    }
}
</script>

<style scoped>
/* 动态背景动画 */
@keyframes blob {
    0%, 100% {
        transform: translate(0, 0) scale(1);
    }
    25% {
        transform: translate(20px, -30px) scale(1.1);
    }
    50% {
        transform: translate(-20px, 20px) scale(0.9);
    }
    75% {
        transform: translate(30px, 10px) scale(1.05);
    }
}

.animate-blob {
    animation: blob 15s infinite ease-in-out;
}

.animation-delay-2000 {
    animation-delay: 2s;
}

.animation-delay-4000 {
    animation-delay: 4s;
}

/* 星星闪烁动画 */
@keyframes twinkle {
    0%, 100% {
        opacity: 0.3;
        transform: scale(1);
    }
    50% {
        opacity: 1;
        transform: scale(1.2);
    }
}

.animate-twinkle {
    animation: twinkle 3s infinite ease-in-out;
}

/* 登录卡片入场动画 */
.login-card {
    animation: cardFloat 0.8s ease-out;
}

@keyframes cardFloat {
    0% {
        opacity: 0;
        transform: translateY(30px);
    }
    100% {
        opacity: 1;
        transform: translateY(0);
    }
}

/* GitHub 按钮悬浮效果 */
.github-btn::before {
    content: '';
    position: absolute;
    inset: 0;
    background: linear-gradient(45deg, transparent, rgba(255,255,255,0.1), transparent);
    transform: translateX(-100%);
    transition: transform 0.5s;
}

.github-btn:hover::before {
    transform: translateX(100%);
}

/* 登录按钮光效 */
.login-btn::after {
    content: '';
    position: absolute;
    top: -50%;
    left: -50%;
    width: 200%;
    height: 200%;
    background: linear-gradient(45deg, transparent, rgba(255,255,255,0.1), transparent);
    transform: rotate(45deg);
    animation: btnShine 3s infinite;
}

@keyframes btnShine {
    0% {
        transform: rotate(45deg) translateX(-100%);
    }
    100% {
        transform: rotate(45deg) translateX(100%);
    }
}

/* 输入框聚焦效果 */
input:focus {
    box-shadow: 0 0 0 4px rgba(139, 92, 246, 0.2);
}
</style>