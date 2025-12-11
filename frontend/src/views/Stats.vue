<template>
    <div class="min-h-screen bg-gradient-to-br from-gray-50 via-white to-blue-50 dark:from-gray-900 dark:via-gray-900 dark:to-gray-800 text-gray-800 dark:text-gray-200 transition-colors duration-300">
        <div class="stats-content container mx-auto px-4 md:px-6 py-8 max-w-7xl">
            <!-- 页面标题 -->
            <div class="stats-header mb-10 text-center">
                <h1 class="text-4xl md:text-5xl font-black mb-3 bg-clip-text text-transparent bg-gradient-to-r from-blue-600 to-purple-600">系统统计</h1>
                <p class="text-gray-600 dark:text-gray-400 text-lg">全面了解您的图床使用情况</p>
            </div>
            
            <!-- 加载状态 -->
            <div v-if="loading" class="loading-container flex flex-col items-center justify-center py-32">
                <div class="relative">
                    <div class="w-20 h-20 border-4 border-gray-200 dark:border-gray-700 border-t-blue-500 rounded-full animate-spin"></div>
                    <div class="absolute inset-0 flex items-center justify-center">
                        <i class="ri-bar-chart-line text-2xl text-blue-500 animate-pulse"></i>
                    </div>
                </div>
                <p class="text-gray-500 dark:text-gray-400 mt-6 font-bold text-lg animate-pulse">加载统计数据中...</p>
            </div>
            
            <div v-else class="space-y-8">
                <!-- 核心统计卡片 -->
                <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                    <!-- 总图片数 -->
                    <div class="stat-card bg-gradient-to-br from-blue-500 to-blue-600 rounded-3xl shadow-xl p-6 text-white transform hover:scale-105 transition-all duration-300">
                        <div class="flex items-center justify-between mb-4">
                            <div class="w-14 h-14 bg-white/20 backdrop-blur-sm rounded-2xl flex items-center justify-center">
                                <i class="ri-image-line text-3xl"></i>
                            </div>
                            <div class="text-right">
                                <div class="text-4xl font-black">{{ formatNumber(stats.total_images) }}</div>
                                <div class="text-sm text-blue-100 font-medium mt-1">总图片数</div>
                            </div>
                        </div>
                    </div>

                    <!-- 总存储空间 -->
                    <div class="stat-card bg-gradient-to-br from-purple-500 to-purple-600 rounded-3xl shadow-xl p-6 text-white transform hover:scale-105 transition-all duration-300">
                        <div class="flex items-center justify-between mb-4">
                            <div class="w-14 h-14 bg-white/20 backdrop-blur-sm rounded-2xl flex items-center justify-center">
                                <i class="ri-database-2-line text-3xl"></i>
                            </div>
                            <div class="text-right">
                                <div class="text-4xl font-black">{{ formatFileSize(stats.total_size) }}</div>
                                <div class="text-sm text-purple-100 font-medium mt-1">总存储空间</div>
                            </div>
                        </div>
                    </div>

                    <!-- 今日上传 -->
                    <div class="stat-card bg-gradient-to-br from-green-500 to-green-600 rounded-3xl shadow-xl p-6 text-white transform hover:scale-105 transition-all duration-300">
                        <div class="flex items-center justify-between mb-4">
                            <div class="w-14 h-14 bg-white/20 backdrop-blur-sm rounded-2xl flex items-center justify-center">
                                <i class="ri-upload-cloud-2-line text-3xl"></i>
                            </div>
                            <div class="text-right">
                                <div class="text-4xl font-black">{{ formatNumber(stats.today_uploads) }}</div>
                                <div class="text-sm text-green-100 font-medium mt-1">今日上传</div>
                            </div>
                        </div>
                    </div>

                    <!-- 本月上传 -->
                    <div class="stat-card bg-gradient-to-br from-orange-500 to-orange-600 rounded-3xl shadow-xl p-6 text-white transform hover:scale-105 transition-all duration-300">
                        <div class="flex items-center justify-between mb-4">
                            <div class="w-14 h-14 bg-white/20 backdrop-blur-sm rounded-2xl flex items-center justify-center">
                                <i class="ri-calendar-check-line text-3xl"></i>
                            </div>
                            <div class="text-right">
                                <div class="text-4xl font-black">{{ formatNumber(stats.month_uploads) }}</div>
                                <div class="text-sm text-orange-100 font-medium mt-1">本月上传</div>
                            </div>
                        </div>
                    </div>
                </div>
                
                <!-- 访客统计卡片 -->
                <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                    <!-- 总访问量 -->
                    <div class="stat-card bg-gradient-to-br from-indigo-500 to-indigo-600 rounded-3xl shadow-xl p-6 text-white transform hover:scale-105 transition-all duration-300">
                        <div class="flex items-center justify-between mb-4">
                            <div class="w-14 h-14 bg-white/20 backdrop-blur-sm rounded-2xl flex items-center justify-center">
                                <i class="ri-eye-line text-3xl"></i>
                            </div>
                            <div class="text-right">
                                <div class="text-4xl font-black">{{ formatNumber(visitStats.total_visits) }}</div>
                                <div class="text-sm text-indigo-100 font-medium mt-1">总访问量</div>
                            </div>
                        </div>
                    </div>

                    <!-- 今日访问 -->
                    <div class="stat-card bg-gradient-to-br from-teal-500 to-teal-600 rounded-3xl shadow-xl p-6 text-white transform hover:scale-105 transition-all duration-300">
                        <div class="flex items-center justify-between mb-4">
                            <div class="w-14 h-14 bg-white/20 backdrop-blur-sm rounded-2xl flex items-center justify-center">
                                <i class="ri-calendar-event-line text-3xl"></i>
                            </div>
                            <div class="text-right">
                                <div class="text-4xl font-black">{{ formatNumber(visitStats.today_visits) }}</div>
                                <div class="text-sm text-teal-100 font-medium mt-1">今日访问</div>
                            </div>
                        </div>
                    </div>

                    <!-- 独立访客 -->
                    <div class="stat-card bg-gradient-to-br from-rose-500 to-rose-600 rounded-3xl shadow-xl p-6 text-white transform hover:scale-105 transition-all duration-300">
                        <div class="flex items-center justify-between mb-4">
                            <div class="w-14 h-14 bg-white/20 backdrop-blur-sm rounded-2xl flex items-center justify-center">
                                <i class="ri-user-3-line text-3xl"></i>
                            </div>
                            <div class="text-right">
                                <div class="text-4xl font-black">{{ formatNumber(visitStats.unique_visitors) }}</div>
                                <div class="text-sm text-rose-100 font-medium mt-1">独立访客</div>
                            </div>
                        </div>
                    </div>

                    <!-- 注册用户 -->
                    <div class="stat-card bg-gradient-to-br from-amber-500 to-amber-600 rounded-3xl shadow-xl p-6 text-white transform hover:scale-105 transition-all duration-300">
                        <div class="flex items-center justify-between mb-4">
                            <div class="w-14 h-14 bg-white/20 backdrop-blur-sm rounded-2xl flex items-center justify-center">
                                <i class="ri-user-star-line text-3xl"></i>
                            </div>
                            <div class="text-right">
                                <div class="text-4xl font-black">{{ formatNumber(visitStats.total_users) }}</div>
                                <div class="text-sm text-amber-100 font-medium mt-1">注册用户</div>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- 访问趋势图表 -->
                <div class="bg-white dark:bg-gray-800 rounded-3xl shadow-xl p-6">
                    <div class="flex items-center justify-between mb-6">
                        <h2 class="text-xl font-bold text-gray-800 dark:text-gray-200">
                            <i class="ri-line-chart-line mr-2 text-blue-500"></i>7日访问趋势
                        </h2>
                        <div class="text-sm text-gray-500">
                            <span class="inline-flex items-center px-3 py-1 rounded-full bg-blue-100 dark:bg-blue-900 text-blue-600 dark:text-blue-300">
                                <i class="ri-pulse-line mr-1"></i>
                                实时更新
                            </span>
                        </div>
                    </div>
                    <div class="h-64 flex items-end justify-between gap-2">
                        <div 
                            v-for="(item, index) in visitStats.visit_trend" 
                            :key="index"
                            class="flex-1 flex flex-col items-center"
                        >
                            <div class="w-full relative group">
                                <!-- 数值提示 -->
                                <div class="absolute -top-8 left-1/2 transform -translate-x-1/2 bg-gray-800 text-white px-2 py-1 rounded text-xs opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap">
                                    {{ formatNumber(item.count) }} 次访问
                                </div>
                                <!-- 柱状条 -->
                                <div 
                                    class="w-full bg-gradient-to-t from-blue-500 to-indigo-500 rounded-t-lg transition-all duration-500 hover:from-blue-600 hover:to-indigo-600"
                                    :style="{ height: calculateTrendHeight(item.count) + 'px', minHeight: '4px' }"
                                ></div>
                            </div>
                            <div class="mt-2 text-xs text-gray-500 dark:text-gray-400">{{ formatTrendDate(item.date) }}</div>
                        </div>
                    </div>
                </div>

                <!-- 访客来源 TOP 10 -->
                <div class="bg-gradient-to-br from-white via-white to-green-50 dark:from-gray-800 dark:via-gray-800 dark:to-green-900/20 rounded-3xl shadow-xl p-6 border border-green-100/50 dark:border-green-800/30">
                    <h2 class="text-xl font-bold text-gray-800 dark:text-gray-200 mb-6 flex items-center">
                        <span class="w-10 h-10 rounded-full bg-gradient-to-br from-green-500 to-teal-500 flex items-center justify-center mr-3 shadow-lg">
                            <i class="ri-global-line text-white text-lg"></i>
                        </span>
                        访客来源 TOP 10
                        <span v-if="visitStats.is_admin" class="ml-3 text-xs font-medium px-2 py-1 bg-green-100 dark:bg-green-900/50 text-green-600 dark:text-green-300 rounded-full">
                            <i class="ri-shield-user-line mr-1"></i>管理员视图
                        </span>
                    </h2>
                    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div 
                            v-for="(item, index) in visitStats.top_ips?.slice(0, 10)" 
                            :key="index"
                            class="flex items-center p-4 bg-gradient-to-r from-gray-50 to-green-50/30 dark:from-gray-700/50 dark:to-green-900/20 rounded-xl hover:from-green-50 hover:to-green-100/50 dark:hover:from-gray-700 dark:hover:to-green-900/30 transition-all duration-300 border border-transparent hover:border-green-200 dark:hover:border-green-700/50"
                        >
                            <!-- 排名 -->
                            <div 
                                class="w-10 h-10 rounded-full flex items-center justify-center text-sm font-bold mr-4 shadow-md"
                                :class="getRankClass(index)"
                            >
                                {{ index + 1 }}
                            </div>
                            <!-- IP和归属地 -->
                            <div class="flex-1 min-w-0">
                                <div class="font-mono text-sm text-gray-800 dark:text-gray-200 truncate font-semibold">
                                    {{ item.ip }}
                                </div>
                                <div class="text-xs text-gray-500 dark:text-gray-400 flex items-center mt-1">
                                    <span class="inline-flex items-center bg-gray-100 dark:bg-gray-800 px-2 py-0.5 rounded">
                                        <i class="ri-map-pin-line mr-1 text-green-500"></i>
                                        <span class="truncate">
                                            {{ item.country || '未知' }}
                                            <template v-if="item.region && item.region !== item.city"> · {{ item.region }}</template>
                                            <template v-if="item.city"> · {{ item.city }}</template>
                                        </span>
                                    </span>
                                </div>
                            </div>
                            <!-- 访问次数 -->
                            <div class="text-right ml-4">
                                <div class="text-2xl font-black text-transparent bg-clip-text bg-gradient-to-r from-green-600 to-teal-600 dark:from-green-400 dark:to-teal-400">{{ formatNumber(item.count) }}</div>
                                <div class="text-xs text-gray-500 dark:text-gray-400">次访问</div>
                            </div>
                        </div>
                    </div>
                    <div v-if="!visitStats.top_ips?.length" class="text-center py-10 text-gray-500">
                        <div class="w-16 h-16 mx-auto mb-3 rounded-full bg-gray-100 dark:bg-gray-700 flex items-center justify-center">
                            <i class="ri-inbox-line text-3xl text-gray-300 dark:text-gray-600"></i>
                        </div>
                        <p class="font-medium">暂无访问数据</p>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { Message } from '@arco-design/web-vue'

const loading = ref(false)
const stats = ref({
    total_images: 0,
    total_size: 0,
    today_uploads: 0,
    month_uploads: 0,
    recent_images: [],
    upload_trend: [],
    format_stats: [],
    size_distribution: []
})

const visitStats = ref({
    total_visits: 0,
    today_visits: 0,
    unique_visitors: 0,
    total_users: 0,
    active_users: 0,
    visit_trend: [],
    top_ips: [],
    image_call_rank: [],
    is_admin: false
})

const loadStats = async () => {
    loading.value = true
    try {
        // 并行加载所有统计数据
        const [dashboardRes, visitRes] = await Promise.all([
            fetch('/api/stats/dashboard', { credentials: 'include' }),
            fetch('/api/stats/visits', { credentials: 'include' })
        ])
        
        if (dashboardRes.ok) {
            const result = await dashboardRes.json()
            stats.value = result.data
        }

        if (visitRes.ok) {
            const result = await visitRes.json()
            visitStats.value = result.data
        }
    } catch (error) {
        console.error('加载统计数据失败:', error)
        Message.error('加载统计数据失败')
    } finally {
        loading.value = false
    }
}

// 计算趋势图高度
const calculateTrendHeight = (count) => {
    if (!visitStats.value.visit_trend || visitStats.value.visit_trend.length === 0) return 4
    const max = Math.max(...visitStats.value.visit_trend.map(item => item.count))
    if (max === 0) return 4
    return Math.max(4, (count / max) * 200)
}

// 格式化趋势日期
const formatTrendDate = (dateString) => {
    if (!dateString) return ''
    const date = new Date(dateString)
    return `${date.getMonth() + 1}/${date.getDate()}`
}

// 获取排名样式
const getRankClass = (index) => {
    const classes = [
        'bg-gradient-to-br from-yellow-400 to-yellow-600 text-white',  // 第1名 金色
        'bg-gradient-to-br from-gray-300 to-gray-500 text-white',      // 第2名 银色
        'bg-gradient-to-br from-amber-600 to-amber-800 text-white',    // 第3名 铜色
    ]
    return classes[index] || 'bg-gray-200 dark:bg-gray-600 text-gray-600 dark:text-gray-300'
}

// 格式化数字
const formatNumber = (num) => {
    if (!num) return '0'
    if (num >= 10000) {
        return (num / 10000).toFixed(1) + '万'
    }
    return num.toLocaleString('zh-CN')
}

// 格式化文件大小
const formatFileSize = (bytes) => {
    if (!bytes) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

onMounted(() => {
    loadStats()
})
</script>

<style scoped>
/* 自定义滚动条 */
.max-h-96::-webkit-scrollbar {
    width: 6px;
}

.max-h-96::-webkit-scrollbar-track {
    background: transparent;
}

.max-h-96::-webkit-scrollbar-thumb {
    background-color: rgba(156, 163, 175, 0.5);
    border-radius: 3px;
}

.max-h-96::-webkit-scrollbar-thumb:hover {
    background-color: rgba(156, 163, 175, 0.8);
}

/* 卡片动画 */
.stat-card {
    animation: fadeInUp 0.5s ease-out forwards;
    opacity: 0;
}

@keyframes fadeInUp {
    from {
        opacity: 0;
        transform: translateY(20px);
    }
    to {
        opacity: 1;
        transform: translateY(0);
    }
}

.stat-card:nth-child(1) { animation-delay: 0.1s; }
.stat-card:nth-child(2) { animation-delay: 0.2s; }
.stat-card:nth-child(3) { animation-delay: 0.3s; }
.stat-card:nth-child(4) { animation-delay: 0.4s;
}

/* 数字计数动画 */
@keyframes countUp {
    from {
        opacity: 0;
        transform: translateY(10px);
    }
    to {
        opacity: 1;
        transform: translateY(0);
    }
}

/* 渐变文字闪烁 */
@keyframes shimmer {
    0% {
        background-position: -200% center;
    }
    100% {
        background-position: 200% center;
    }
}

/* 加载动画增强 */
.loading-container {
    animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
    0%, 100% {
        opacity: 1;
    }
    50% {
        opacity: 0.7;
    }
}

/* 页面标题动画 */
.stats-header h1 {
    animation: slideInDown 0.8s ease-out forwards;
}

.stats-header p {
    animation: fadeIn 1s ease-out 0.3s forwards;
    opacity: 0;
}

@keyframes slideInDown {
    from {
        opacity: 0;
        transform: translateY(-30px);
    }
    to {
        opacity: 1;
        transform: translateY(0);
    }
}

@keyframes fadeIn {
    from {
        opacity: 0;
    }
    to {
        opacity: 1;
    }
}
</style>