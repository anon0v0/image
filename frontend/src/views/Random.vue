<template>
    <div class="random-page fixed inset-0 overflow-hidden bg-gradient-to-br from-gray-900 via-gray-800 to-black text-white flex flex-col items-center justify-center pt-20 pb-4 px-4 z-40">
        <!-- 背景动画 -->
        <div class="absolute inset-0 overflow-hidden opacity-20">
            <div class="absolute w-96 h-96 bg-blue-500 rounded-full blur-3xl -top-48 -left-48 animate-pulse"></div>
            <div class="absolute w-96 h-96 bg-purple-500 rounded-full blur-3xl -bottom-48 -right-48 animate-pulse" style="animation-delay: 1s"></div>
        </div>

        <!-- 主内容区 -->
        <div class="relative z-10 w-full max-w-5xl flex-1 flex flex-col justify-center">

            <!-- 图片容器 -->
            <div class="image-container mb-4 rounded-3xl overflow-hidden shadow-2xl bg-white/5 backdrop-blur-sm border border-white/10 relative" style="max-height: calc(100vh - 220px);">
                <!-- 加载动画 -->
                <div v-if="loading || !imageLoaded" class="absolute inset-0 flex flex-col items-center justify-center">
                    <div class="loading-animation text-center">
                        <div class="relative inline-block mb-6">
                            <!-- 旋转圆环 -->
                            <div class="w-24 h-24 border-4 border-white/20 border-t-blue-500 rounded-full animate-spin"></div>
                            <!-- 中心图标 -->
                            <div class="absolute inset-0 flex items-center justify-center">
                                <i class="ri-image-line text-4xl text-blue-500 animate-pulse"></i>
                            </div>
                        </div>
                        <p class="text-gray-400 font-medium animate-pulse">正在加载随机图片...</p>
                    </div>
                </div>

                <!-- 图片显示 -->
                <div v-if="currentImage" class="relative group flex items-center justify-center" :class="{ 'opacity-0': !imageLoaded }">
                    <img 
                        :src="getFullUrl(currentImage.url)" 
                        :alt="currentImage.filename"
                        class="w-full h-auto object-contain transition-opacity duration-500"
                        style="max-height: calc(100vh - 240px);"
                        :class="{ 'opacity-100': imageLoaded }"
                        @load="onImageLoad"
                    />
                    
                    <!-- 悬浮信息 -->
                    <div class="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/80 via-black/40 to-transparent p-6 opacity-0 group-hover:opacity-100 transition-opacity duration-300">
                        <div class="flex flex-wrap gap-2 mb-3" v-if="currentImage.tags">
                            <span v-for="tag in currentImage.tags.split(',').slice(0, 5)" :key="tag" class="px-3 py-1 bg-blue-500/80 backdrop-blur-sm rounded-full text-xs font-bold">
                                #{{ tag }}
                            </span>
                        </div>
                        <!-- Filename removed as per instruction "隐藏随机图片文件名" -->
                        <div class="flex items-center gap-4 text-sm text-gray-300">
                            <span>{{ currentImage.width }} × {{ currentImage.height }}</span>
                            <span>{{ formatFileSize(currentImage.file_size) }}</span>
                            <span>{{ formatDate(currentImage.created_at) }}</span>
                        </div>
                    </div>
                </div>

                <!-- 空状态 -->
                <div v-else class="aspect-video flex items-center justify-center">
                    <div class="text-center">
                        <i class="ri-image-2-line text-6xl text-gray-600 mb-4"></i>
                        <p class="text-gray-500">暂无图片</p>
                    </div>
                </div>
            </div>

            <!-- 控制按钮 -->
            <div class="flex flex-wrap justify-center gap-3 mt-4">
                <button 
                    @click="loadRandomImage"
                    :disabled="loading"
                    class="px-6 py-3 bg-gradient-to-r from-blue-500 to-purple-600 hover:from-blue-600 hover:to-purple-700 disabled:from-gray-600 disabled:to-gray-700 rounded-xl font-bold shadow-lg shadow-blue-500/30 hover:shadow-xl transform hover:-translate-y-0.5 transition-all disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none flex items-center gap-2"
                >
                    <i class="ri-shuffle-line text-xl"></i>
                    <span>{{ loading ? '加载中...' : '随机一张' }}</span>
                </button>

                <!-- 复制按钮（仅登录后显示） -->
                <button 
                    v-if="isLoggedIn && currentImage"
                    @click="copyLink"
                    class="px-6 py-3 bg-gradient-to-r from-green-500 to-emerald-600 hover:from-green-600 hover:to-emerald-700 rounded-xl font-bold shadow-lg shadow-green-500/30 hover:shadow-xl transform hover:-translate-y-0.5 transition-all flex items-center gap-2"
                >
                    <i class="ri-file-copy-line text-xl"></i>
                    <span>复制链接</span>
                </button>

                <!-- 下载按钮（仅登录后显示） -->
                <button 
                    v-if="isLoggedIn && currentImage"
                    @click="downloadImage"
                    class="px-6 py-3 bg-gradient-to-r from-purple-500 to-pink-600 hover:from-purple-600 hover:to-pink-700 rounded-xl font-bold shadow-lg shadow-purple-500/30 hover:shadow-xl transform hover:-translate-y-0.5 transition-all flex items-center gap-2"
                >
                    <i class="ri-download-cloud-line text-xl"></i>
                    <span>下载图片</span>
                </button>

                <router-link 
                    to="/"
                    class="px-6 py-3 bg-white/10 hover:bg-white/20 backdrop-blur-sm rounded-xl font-bold shadow-lg transform hover:-translate-y-0.5 transition-all flex items-center gap-2 border border-white/20"
                >
                    <i class="ri-home-line text-xl"></i>
                    <span>返回首页</span>
                </router-link>
            </div>

            <!-- 键盘提示和登录提示 -->
            <div class="text-center mt-4 text-gray-500 text-xs space-y-1">
                <p>快捷键：<kbd class="px-2 py-0.5 bg-white/10 rounded text-xs font-mono mx-1">空格</kbd> 刷新</p>
                <p v-if="!isLoggedIn" class="text-yellow-400/80">
                    💡 登录后开启复制链接和下载图片功能
                </p>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const loading = ref(false)
const imageLoaded = ref(false)
const currentImage = ref(null)
const isLoggedIn = ref(false)

// 检查登录状态
const checkLoginStatus = async () => {
    try {
        const userInfoStr = localStorage.getItem('userInfo')
        if (userInfoStr) {
            isLoggedIn.value = true
        }
    } catch (error) {
        console.error('检查登录状态失败:', error)
    }
}

// 获取完整 URL
const getFullUrl = (path) => {
    if (!path) return ''
    if (typeof window !== 'undefined') {
        return window.location.origin + path
    }
    return path
}

// 加载随机图片
const loadRandomImage = async () => {
    if (loading.value) return
    
    loading.value = true
    imageLoaded.value = false // 重置图片加载状态
    try {
        const response = await fetch('/api/images?limit=1&sort_by=random', {
            credentials: 'include'
        })
        
        if (response.ok) {
            const result = await response.json()
            if (result.data && result.data.images && result.data.images.length > 0) {
                currentImage.value = result.data.images[0]
            }
        }
    } catch (error) {
        console.error('加载随机图片失败:', error)
    } finally {
        loading.value = false
    }
}

// 图片加载完成
const onImageLoad = () => {
    imageLoaded.value = true
}

// 格式化文件大小
const formatFileSize = (bytes) => {
    if (!bytes) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 格式化日期
const formatDate = (dateString) => {
    if (!dateString) return ''
    return new Date(dateString).toLocaleDateString('zh-CN')
}

// 键盘事件处理
const handleKeyPress = (event) => {
    if (event.code === 'Space') {
        event.preventDefault()
        loadRandomImage()
    }
}

// 复制链接
const copyLink = async () => {
    if (!currentImage.value) return
    const url = getFullUrl(currentImage.value.url)
    try {
        await navigator.clipboard.writeText(url)
        alert('链接已复制到剪贴板')
    } catch (error) {
        console.error('复制失败:', error)
    }
}

// 下载图片
const downloadImage = () => {
    if (!currentImage.value) return
    const link = document.createElement('a')
    link.href = getFullUrl(currentImage.value.url)
    link.download = currentImage.value.filename || 'image'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
}

onMounted(() => {
    checkLoginStatus()
    loadRandomImage()
    window.addEventListener('keydown', handleKeyPress)
})

onUnmounted(() => {
    window.removeEventListener('keydown', handleKeyPress)
})
</script>

<style scoped>
@keyframes fade-in {
    from {
        opacity: 0;
        transform: translateY(20px);
    }
    to {
        opacity: 1;
        transform: translateY(0);
    }
}

@keyframes scale-in {
    from {
        opacity: 0;
        transform: scale(0.95);
    }
    to {
        opacity: 1;
        transform: scale(1);
    }
}

.animate-fade-in {
    animation: fade-in 0.6s ease-out;
}

.animate-scale-in {
    animation: scale-in 0.5s ease-out;
}

kbd {
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}
</style>
