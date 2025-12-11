<template>
    <div class="min-h-screen bg-gradient-to-br from-gray-50 via-white to-blue-50 dark:from-gray-900 dark:via-gray-900 dark:to-gray-800 text-gray-800 dark:text-gray-200 transition-colors duration-300">
        <!-- 轮播图区域 -->
        <div class="container mx-auto px-4 md:px-6 max-w-7xl mb-6 pt-4">
            <div class="carousel-section relative h-[320px] md:h-[480px] bg-gradient-to-br from-slate-900 via-purple-900/50 to-slate-900 overflow-hidden group rounded-3xl shadow-2xl border border-purple-500/20">
                <!-- 装饰性背景元素 -->
                <div class="absolute inset-0 overflow-hidden">
                    <div class="absolute -top-1/2 -left-1/4 w-[80%] h-[80%] bg-gradient-conic from-blue-500/20 via-purple-500/20 to-pink-500/20 rounded-full blur-3xl animate-slow-spin"></div>
                    <div class="absolute -bottom-1/2 -right-1/4 w-[60%] h-[60%] bg-gradient-conic from-cyan-500/20 via-blue-500/20 to-purple-500/20 rounded-full blur-3xl animate-slow-spin-reverse"></div>
                </div>
                
                <div 
                    v-for="(image, index) in carouselImages" 
                    :key="image.id"
                    class="absolute inset-0 transition-all duration-1000 ease-in-out"
                    :class="{ 'opacity-100 scale-100': currentCarouselIndex === index, 'opacity-0 scale-105': currentCarouselIndex !== index }"
                >
                    <!-- 背景模糊层 -->
                    <div class="absolute inset-0 overflow-hidden opacity-30">
                        <img :src="getThumbUrl(image.url)" class="w-full h-full object-cover blur-3xl scale-125 transform" />
                    </div>
                    
                    <!-- 主图容器 -->
                    <div class="absolute inset-0 flex items-center justify-center p-8 md:p-12">
                        <div class="relative max-w-full max-h-full">
                            <!-- 图片光晕效果 -->
                            <div class="absolute -inset-4 bg-gradient-to-r from-blue-500/30 via-purple-500/30 to-pink-500/30 rounded-3xl blur-2xl opacity-60 group-hover:opacity-80 transition-opacity"></div>
                            <img 
                                :src="getPreviewUrl(image.url)" 
                                class="relative max-w-full max-h-[260px] md:max-h-[380px] object-contain drop-shadow-2xl rounded-2xl ring-2 ring-white/20 group-hover:ring-purple-500/50 transition-all duration-500"
                                alt="Carousel Image"
                                loading="lazy"
                                @error="handlePreviewError($event, image.url)"
                            />
                        </div>
                    </div>
                    
                    <!-- 信息栏 -->
                    <div class="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black via-black/70 to-transparent pt-12 pb-5 px-6 md:px-10 z-10 text-white pointer-events-none">
                        <div class="container mx-auto max-w-7xl">
                            <div class="flex flex-col items-start space-y-2">
                                <div class="flex flex-wrap gap-2" v-if="image.tags">
                                    <span v-for="tag in image.tags.split(',').slice(0, 3)" :key="tag" class="px-3 py-1 bg-gradient-to-r from-blue-500/80 to-purple-500/80 backdrop-blur-md rounded-full text-xs font-bold tracking-wide shadow-lg border border-white/10">
                                        #{{ tag }}
                                    </span>
                                </div>
                                <h2 class="text-xl md:text-3xl font-black drop-shadow-2xl tracking-tight leading-tight text-white max-w-full truncate">
                                    {{ image.tags ? image.tags.split(',')[0] : image.filename }}
                                </h2>
                                <p class="text-xs text-gray-300/80 font-medium flex items-center gap-2">
                                    <span class="flex items-center gap-1"><i class="ri-aspect-ratio-line"></i>{{ image.width }}×{{ image.height }}</span>
                                    <span>•</span>
                                    <span class="flex items-center gap-1"><i class="ri-file-3-line"></i>{{ formatFileSize(image.file_size) }}</span>
                                </p>
                            </div>
                        </div>
                    </div>
                </div>
                
                <!-- 导航按钮 -->
                <button @click="prevSlide" class="absolute left-4 md:left-8 top-1/2 -translate-y-1/2 w-12 h-12 bg-white/10 hover:bg-white/25 backdrop-blur-xl text-white rounded-full transition-all z-30 opacity-0 group-hover:opacity-100 transform -translate-x-4 group-hover:translate-x-0 shadow-xl border border-white/30 flex items-center justify-center hover:scale-110">
                    <i class="ri-arrow-left-s-line text-2xl"></i>
                </button>
                <button @click="nextSlide" class="absolute right-4 md:right-8 top-1/2 -translate-y-1/2 w-12 h-12 bg-white/10 hover:bg-white/25 backdrop-blur-xl text-white rounded-full transition-all z-30 opacity-0 group-hover:opacity-100 transform translate-x-4 group-hover:translate-x-0 shadow-xl border border-white/30 flex items-center justify-center hover:scale-110">
                    <i class="ri-arrow-right-s-line text-2xl"></i>
                </button>
                
                <!-- 指示器 -->
                <div class="absolute bottom-5 left-1/2 -translate-x-1/2 z-30 flex gap-2 bg-black/30 backdrop-blur-md px-3 py-2 rounded-full border border-white/10">
                    <button 
                        v-for="(img, idx) in carouselImages.slice(0, 8)" 
                        :key="img.id" 
                        @click="currentCarouselIndex = idx"
                        class="h-2 rounded-full transition-all duration-500 hover:scale-125"
                        :class="currentCarouselIndex === idx ? 'w-8 bg-gradient-to-r from-blue-400 to-purple-500' : 'w-2 bg-white/40 hover:bg-white/70'"
                    ></button>
                </div>
            </div>
        </div>

        <!-- 画廊内容 -->
        <div class="gallery-content container mx-auto px-4 md:px-6 pb-20 max-w-7xl">
            <!-- 分类标签栏 -->
            <div class="filter-bar mb-6 bg-white/90 dark:bg-gray-800/90 backdrop-blur-xl py-3 px-2 rounded-2xl shadow-lg border border-gray-200/50 dark:border-gray-700/50 sticky top-20 z-40 transition-all duration-300">
                <div class="categories flex gap-2 overflow-x-auto scrollbar-hide px-2" style="padding-right: 20px;">
                    <button 
                        v-for="cat in categories" 
                        :key="cat"
                        @click="selectCategory(cat)"
                        class="category-btn flex-shrink-0 px-4 py-2 rounded-xl text-sm font-bold transition-all duration-300 relative overflow-hidden group flex items-center gap-1.5 whitespace-nowrap"
                        :class="selectedCategory === cat 
                            ? 'bg-gradient-to-r from-blue-500 via-purple-500 to-pink-500 text-white shadow-lg shadow-purple-500/30 scale-105' 
                            : 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600 hover:scale-105'"
                    >
                        <i :class="getCategoryIcon(cat)" class="text-sm"></i>
                        <span>{{ cat }}</span>
                    </button>
                    <!-- 末尾占位元素确保最后的标签可以完全滚动显示 -->
                    <div class="flex-shrink-0 w-4"></div>
                </div>
            </div>

            <!-- 加载状态 -->
            <div v-if="loading && images.length === 0" class="loading-container flex flex-col items-center justify-center py-40">
                <div class="relative">
                    <div class="w-20 h-20 border-4 border-gray-200 dark:border-gray-700 border-t-blue-500 rounded-full animate-spin"></div>
                    <div class="absolute inset-0 flex items-center justify-center">
                        <i class="ri-image-line text-2xl text-blue-500 animate-pulse"></i>
                    </div>
                </div>
                <p class="text-gray-500 dark:text-gray-400 mt-6 font-bold text-lg animate-pulse">加载精彩内容中...</p>
            </div>
            
            <!-- 图片网格 -->
            <div v-else-if="images.length > 0" class="images-container">
                <div class="images-grid grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4">
                    <div 
                        v-for="image in images" 
                        :key="image.id"
                        class="image-card group relative rounded-3xl overflow-hidden shadow-lg hover:shadow-2xl transition-all duration-400 bg-white dark:bg-gray-800 cursor-pointer transform hover:-translate-y-1 border-2 border-transparent hover:border-blue-500/30"
                        @click="openImageModal(image)"
                    >
                        <div class="image-wrapper relative aspect-square overflow-hidden bg-gray-200 dark:bg-gray-800">
                            <img 
                                :src="getThumbUrl(image.url)" 
                                :data-original="image.url"
                                :alt="image.filename"
                                loading="lazy"
                                class="w-full h-full object-cover transition-transform duration-400 ease-in-out group-hover:scale-110"
                                @error="handleThumbError"
                            />
                            
                            <!-- 悬浮信息遮罩 -->
                            <div class="absolute inset-0 bg-gradient-to-t from-black/95 via-black/60 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-400 ease-in-out flex flex-col justify-end p-3 pointer-events-none">
                                <div class="flex flex-wrap gap-1 mb-2" v-if="image.tags">
                                    <span 
                                        v-for="tag in image.tags.split(',').slice(0, 2)" 
                                        :key="tag" 
                                        class="px-2 py-0.5 bg-gradient-to-r from-blue-500/90 to-purple-500/90 backdrop-blur-sm text-white rounded-full text-[10px] font-bold"
                                    >
                                        #{{ tag }}
                                    </span>
                                </div>
                                <div class="flex justify-between items-center text-xs text-white font-bold">
                                    <span class="bg-black/50 px-2 py-1 rounded-lg backdrop-blur-sm">{{ image.width }}×{{ image.height }}</span>
                                    <span class="bg-black/50 px-2 py-1 rounded-lg backdrop-blur-sm">{{ formatFileSize(image.file_size) }}</span>
                                </div>
                            </div>

                            <!-- 查看图标 -->
                            <div class="absolute top-3 right-3 w-10 h-10 bg-white/20 backdrop-blur-md rounded-full flex items-center justify-center text-white opacity-0 group-hover:opacity-100 transition-all duration-400 transform scale-0 group-hover:scale-100">
                                <i class="ri-eye-line text-lg"></i>
                            </div>
                        </div>
                    </div>
                </div>
                
                <!-- 无限滚动触发器 - 增大高度便于触发 -->
                <div ref="loadMoreTrigger" class="h-96 flex justify-center items-center my-8">
                    <div v-if="loading" class="flex items-center gap-3 text-blue-500">
                        <div class="w-6 h-6 border-3 border-blue-300 border-t-blue-600 rounded-full animate-spin"></div>
                        <span class="font-bold">加载更多精彩...</span>
                    </div>
                    <div v-else-if="currentPage >= totalPages" class="flex flex-col items-center gap-3">
                        <div class="w-16 h-16 bg-gradient-to-br from-blue-100 to-purple-100 dark:from-gray-700 dark:to-gray-800 rounded-full flex items-center justify-center">
                            <i class="ri-check-double-line text-3xl text-blue-500"></i>
                        </div>
                        <p class="text-gray-400 text-sm font-bold">已经到底啦，没有更多内容了</p>
                    </div>
                </div>
            </div>
            
            <!-- 空状态 -->
            <div v-else class="empty-state flex flex-col items-center justify-center py-40 text-center">
                <div class="w-32 h-32 bg-gradient-to-br from-blue-100 to-purple-100 dark:from-gray-800 dark:to-gray-700 rounded-full flex items-center justify-center mb-8 shadow-xl">
                    <i class="ri-image-2-line text-6xl text-blue-500"></i>
                </div>
                <h3 class="text-2xl font-black mb-3 text-gray-800 dark:text-gray-200">暂无图片</h3>
                <p class="text-gray-500 dark:text-gray-400 mb-10 max-w-md">这里还是一片空白，快来上传你的第一张图片吧！</p>
                <router-link to="/upload" class="px-10 py-4 bg-gradient-to-r from-blue-500 to-purple-600 hover:from-blue-600 hover:to-purple-700 text-white rounded-2xl font-bold transition-all shadow-xl shadow-blue-500/30 hover:shadow-2xl hover:shadow-blue-500/50 transform hover:-translate-y-1 hover:scale-105">
                    <i class="ri-upload-cloud-2-line mr-2"></i>开始上传
                </router-link>
            </div>
        </div>

        <!-- 图片详情弹窗 -->
        <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center bg-black/95 backdrop-blur-xl p-4 md:p-8 animate-fade-in" @click.self="closeModal">
            <div class="bg-white dark:bg-gray-900 rounded-3xl overflow-hidden max-w-6xl w-full max-h-[90vh] flex flex-col md:flex-row shadow-2xl border border-gray-200 dark:border-gray-700 animate-scale-in">
                <!-- 图片预览区 -->
                <div class="flex-1 bg-gradient-to-br from-gray-900 to-gray-800 flex items-center justify-center p-8 relative group">
                    <!-- 加载动画 -->
                    <div v-if="imageLoading" class="absolute inset-0 flex items-center justify-center">
                        <div class="loading-animation text-center">
                            <div class="relative inline-block">
                                <!-- 旋转圆环 -->
                                <div class="w-20 h-20 border-4 border-white/20 border-t-blue-500 rounded-full animate-spin"></div>
                                <!-- 中心图标 -->
                                <div class="absolute inset-0 flex items-center justify-center">
                                    <i class="ri-image-line text-3xl text-blue-500 animate-pulse"></i>
                                </div>
                            </div>
                            <p class="mt-6 text-white/80 font-medium animate-pulse">图片加载中...</p>
                        </div>
                    </div>
                    
                    <!-- 图片 -->
                    <img 
                        :src="getFullUrl(currentImage.url)" 
                        class="max-w-full max-h-[80vh] object-contain rounded-2xl shadow-2xl ring-1 ring-white/10 transition-opacity duration-300" 
                        :class="{ 'opacity-0': imageLoading, 'opacity-100': !imageLoading }"
                        @load="imageLoading = false"
                        @error="handlePreviewError($event, currentImage.url)"
                    />
                </div>
                
                <!-- 信息侧栏 -->
                <div class="w-full md:w-[420px] p-8 flex flex-col bg-white dark:bg-gray-900 overflow-y-auto border-l border-gray-100 dark:border-gray-800">
                    <div class="flex justify-between items-start mb-8">
                        <h3 class="text-3xl font-black text-transparent bg-clip-text bg-gradient-to-r from-blue-600 to-purple-600">图片详情</h3>
                        <button @click="closeModal" class="w-10 h-10 rounded-full bg-gray-100 dark:bg-gray-800 flex items-center justify-center text-gray-500 hover:text-gray-900 dark:hover:text-white hover:bg-gray-200 dark:hover:bg-gray-700 transition-all transform hover:rotate-90">
                            <i class="ri-close-line text-2xl"></i>
                        </button>
                    </div>

                    <!-- 详细信息 -->
                    <div class="space-y-6 mb-8 flex-1">
                        <div class="p-5 bg-gradient-to-br from-blue-50 to-purple-50 dark:from-gray-800 dark:to-gray-700 rounded-2xl shadow-inner">
                            <span class="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider font-bold mb-3 block flex items-center gap-2">
                                <i class="ri-information-line"></i>基本信息
                            </span>
                            <div class="grid grid-cols-2 gap-4">
                                <div class="bg-white dark:bg-gray-800 p-3 rounded-xl shadow-sm">
                                    <span class="text-xs text-gray-400 block mb-1">尺寸</span>
                                    <span class="font-mono text-base font-bold text-blue-600 dark:text-blue-400">{{ currentImage.width }} × {{ currentImage.height }}</span>
                                </div>
                                <div class="bg-white dark:bg-gray-800 p-3 rounded-xl shadow-sm">
                                    <span class="text-xs text-gray-400 block mb-1">大小</span>
                                    <span class="font-mono text-base font-bold text-purple-600 dark:text-purple-400">{{ formatFileSize(currentImage.file_size) }}</span>
                                </div>
                                <div class="col-span-2 bg-white dark:bg-gray-800 p-3 rounded-xl shadow-sm">
                                    <span class="text-xs text-gray-400 block mb-1">上传时间</span>
                                    <span class="font-mono text-sm font-bold text-gray-700 dark:text-gray-300">{{ formatDate(currentImage.created_at) }}</span>
                                </div>
                            </div>
                        </div>

                        <div v-if="currentImage.tags" class="p-5 bg-gradient-to-br from-blue-50 to-purple-50 dark:from-gray-800 dark:to-gray-700 rounded-2xl shadow-inner">
                            <span class="text-xs text-gray-500 dark:text-gray-400 uppercase tracking-wider font-bold mb-3 block flex items-center gap-2">
                                <i class="ri-price-tag-3-line"></i>AI 标签
                            </span>
                            <div class="flex flex-wrap gap-2">
                                <span v-for="tag in currentImage.tags.split(',')" :key="tag" class="text-sm bg-gradient-to-r from-blue-500 to-purple-600 text-white px-4 py-2 rounded-xl font-bold hover:shadow-lg hover:scale-110 transition-all cursor-default">
                                    #{{ tag }}
                                </span>
                            </div>
                        </div>
                    </div>

                    <!-- 操作按钮 -->
                    <div class="actions space-y-3 mt-auto">
                        <template v-if="isLoggedIn">
                            <div class="grid grid-cols-2 gap-3">
                                <div class="relative group">
                                    <button class="w-full flex items-center justify-center gap-2 px-5 py-3.5 bg-gradient-to-r from-blue-500 to-blue-600 hover:from-blue-600 hover:to-blue-700 text-white rounded-xl text-sm font-bold transition-all shadow-lg hover:shadow-xl transform hover:-translate-y-0.5">
                                        <i class="ri-link-m text-lg"></i> 复制链接
                                    </button>
                                    <div class="absolute bottom-full left-0 w-full mb-2 bg-white dark:bg-gray-800 rounded-xl shadow-2xl border border-gray-100 dark:border-gray-700 overflow-hidden opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-200 z-10">
                                        <button @click="copyLink(getFullUrl(currentImage.url))" class="w-full text-left px-4 py-3 hover:bg-blue-50 dark:hover:bg-gray-700 text-sm font-medium transition-colors">原图链接</button>
                                        <button @click="copyLink(getPreviewUrl(currentImage.url), true)" class="w-full text-left px-4 py-3 hover:bg-blue-50 dark:hover:bg-gray-700 text-sm font-medium transition-colors">压缩图链接</button>
                                    </div>
                                </div>
                                
                                <div class="relative group">
                                    <button class="w-full flex items-center justify-center gap-2 px-5 py-3.5 bg-gradient-to-r from-purple-500 to-purple-600 hover:from-purple-600 hover:to-purple-700 text-white rounded-xl text-sm font-bold transition-all shadow-lg hover:shadow-xl transform hover:-translate-y-0.5">
                                        <i class="ri-code-s-slash-line text-lg"></i> 复制HTML
                                    </button>
                                    <div class="absolute bottom-full left-0 w-full mb-2 bg-white dark:bg-gray-800 rounded-xl shadow-2xl border border-gray-100 dark:border-gray-700 overflow-hidden opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-200 z-10">
                                        <button @click="copyHtml(getFullUrl(currentImage.url))" class="w-full text-left px-4 py-3 hover:bg-purple-50 dark:hover:bg-gray-700 text-sm font-medium transition-colors">原图 HTML</button>
                                        <button @click="copyHtml(getPreviewUrl(currentImage.url), true)" class="w-full text-left px-4 py-3 hover:bg-purple-50 dark:hover:bg-gray-700 text-sm font-medium transition-colors">压缩图 HTML</button>
                                    </div>
                                </div>

                                <div class="relative group col-span-2">
                                    <button class="w-full flex items-center justify-center gap-2 px-5 py-3.5 bg-gradient-to-r from-pink-500 to-pink-600 hover:from-pink-600 hover:to-pink-700 text-white rounded-xl text-sm font-bold transition-all shadow-lg hover:shadow-xl transform hover:-translate-y-0.5">
                                        <i class="ri-markdown-line text-lg"></i> 复制 Markdown
                                    </button>
                                    <div class="absolute bottom-full left-0 w-full mb-2 bg-white dark:bg-gray-800 rounded-xl shadow-2xl border border-gray-100 dark:border-gray-700 overflow-hidden opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-200 z-10">
                                        <button @click="copyMarkdown(getFullUrl(currentImage.url))" class="w-full text-left px-4 py-3 hover:bg-pink-50 dark:hover:bg-gray-700 text-sm font-medium transition-colors">原图 Markdown</button>
                                        <button @click="copyMarkdown(getPreviewUrl(currentImage.url), true)" class="w-full text-left px-4 py-3 hover:bg-pink-50 dark:hover:bg-gray-700 text-sm font-medium transition-colors">压缩图 Markdown</button>
                                    </div>
                                </div>
                            </div>
                            
                            <button @click="downloadImage(currentImage)" class="w-full flex items-center justify-center gap-2 px-5 py-4 bg-gradient-to-r from-green-500 to-emerald-600 hover:from-green-600 hover:to-emerald-700 text-white rounded-xl text-base font-black transition-all shadow-xl shadow-green-500/30 hover:shadow-2xl hover:shadow-green-500/50 transform hover:-translate-y-1">
                                <i class="ri-download-cloud-2-line text-xl"></i> 下载原图
                            </button>
                            <button v-if="isAdmin" @click="deleteImage(currentImage)" class="w-full flex items-center justify-center gap-2 px-5 py-3.5 border-2 border-red-200 dark:border-red-900/30 text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-xl text-sm font-bold transition-all">
                                <i class="ri-delete-bin-line text-lg"></i> 删除图片
                            </button>
                        </template>
                        <template v-else>
                            <div class="text-center p-8 bg-gradient-to-br from-blue-50 to-purple-50 dark:from-gray-800 dark:to-gray-700 rounded-2xl border-2 border-dashed border-blue-200 dark:border-gray-600">
                                <i class="ri-lock-line text-5xl text-blue-500 mb-4"></i>
                                <p class="text-sm text-gray-600 dark:text-gray-400 mb-6 font-medium">登录后可进行复制链接、下载原图等操作</p>
                                <router-link to="/login" class="inline-block w-full px-8 py-4 bg-gradient-to-r from-blue-500 to-purple-600 text-white rounded-xl text-base font-black hover:shadow-xl transition-all transform hover:-translate-y-1">
                                    <i class="ri-login-box-line mr-2"></i>立即登录
                                </router-link>
                            </div>
                        </template>
                    </div>
                </div>
            </div>
        </div>

        <!-- 固定在右下角的翻页控件 -->
        <div v-if="totalPages > 1" class="fixed bottom-6 right-6 z-50">
            <div class="bg-white/95 dark:bg-gray-800/95 backdrop-blur-xl rounded-2xl shadow-2xl border border-gray-200/50 dark:border-gray-700 p-3 flex flex-col gap-2">
                <!-- 页码显示 -->
                <div class="flex items-center justify-center gap-1 text-sm font-bold">
                    <span class="text-gray-500 dark:text-gray-400">第</span>
                    <span class="text-lg text-transparent bg-clip-text bg-gradient-to-r from-blue-500 to-purple-500">{{ currentPage }}</span>
                    <span class="text-gray-400">/</span>
                    <span class="text-gray-600 dark:text-gray-300">{{ totalPages }}</span>
                    <span class="text-gray-500 dark:text-gray-400">页</span>
                </div>
                
                <!-- 跳转控件 -->
                <div class="flex items-center gap-1">
                    <input 
                        v-model="jumpPageInput" 
                        @keyup.enter="handleJumpPage"
                        type="number" 
                        min="1" 
                        :max="totalPages"
                        class="w-14 px-2 py-1.5 rounded-lg border border-gray-300 dark:border-gray-600 bg-gray-50 dark:bg-gray-900 text-center focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none font-bold text-sm"
                        placeholder="页码"
                    />
                    <button @click="handleJumpPage" class="px-3 py-1.5 bg-gradient-to-r from-blue-500 to-purple-600 text-white rounded-lg text-sm font-bold hover:shadow-lg hover:scale-105 transition-all flex items-center gap-1">
                        <i class="ri-arrow-right-line"></i>
                        跳转
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted, computed, onUnmounted, watch, nextTick } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import Message from '@/utils/message.js'

// 获取完整 URL 辅助函数
const getFullUrl = (path) => {
  if (!path) return ''
  if (typeof window !== 'undefined') {
    return window.location.origin + path
  }
  return path
}

// 获取预览图 URL (1920p)
const getPreviewUrl = (url) => {
    if (!url) return ''
    const lastDotIndex = url.lastIndexOf('.')
    if (lastDotIndex === -1) return url
    return getFullUrl(url.substring(0, lastDotIndex) + '_preview.webp')
}

// 预览图加载失败回退到原图
const handlePreviewError = (event, originalUrl) => {
    const img = event.target
    const fullOriginal = getFullUrl(originalUrl)
    if (img.src !== fullOriginal) {
        img.src = fullOriginal
    }
}

// 响应式数据
const images = ref([])
const loading = ref(false)
const currentPage = ref(1)
const totalPages = ref(1)
const pageSize = ref(20)
const jumpPageInput = ref('')
const router = useRouter()
const route = useRoute()
const isLoggedIn = ref(false)
const isAdmin = ref(false)
const loadMoreTrigger = ref(null)
let observer = null
const hasScrolledPages = ref(false) // 标记是否已经通过滚动加载了多页

// 检查登录状态
const checkLoginStatus = async () => {
    if (typeof localStorage !== 'undefined') {
        const userInfoStr = localStorage.getItem('userInfo')
        if (userInfoStr) {
            try {
                const userInfo = JSON.parse(userInfoStr)
                isLoggedIn.value = true
                isAdmin.value = userInfo.role === 'admin'
            } catch (e) {}
        }
    }

    try {
        const response = await fetch('/api/user/status', { credentials: 'include' })
        const data = await response.json()
        if (data.code === 200 && data.data.logged_in) {
            isLoggedIn.value = true
            isAdmin.value = data.data.role === 'admin'
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
            if (typeof localStorage !== 'undefined') {
                localStorage.removeItem('userInfo')
            }
        }
    } catch (error) {
        console.error('检查登录状态失败:', error)
    }
}

checkLoginStatus()

// 轮播图相关
const carouselImages = ref([])
const currentCarouselIndex = ref(0)
let carouselInterval = null

// 搜索与分类
const searchQuery = ref('')
const selectedCategory = ref('全部')
const categories = ['全部', '动漫', '人物', '风景', '影视', '游戏', '美食', '动物', '艺术', '宇宙', '科技', '简约', '机车', '其他']

// 分类图标映射
const getCategoryIcon = (cat) => {
    const iconMap = {
        '全部': 'ri-apps-line',
        '动漫': 'ri-ghost-smile-line',
        '人物': 'ri-user-smile-line',
        '风景': 'ri-landscape-line',
        '影视': 'ri-movie-line',
        '游戏': 'ri-gamepad-line',
        '美食': 'ri-restaurant-line',
        '动物': 'ri-bear-smile-line',
        '艺术': 'ri-palette-line',
        '宇宙': 'ri-planet-line',
        '科技': 'ri-cpu-line',
        '简约': 'ri-contrast-2-line',
        '机车': 'ri-motorbike-line',
        '其他': 'ri-more-line'
    }
    return iconMap[cat] || 'ri-folder-line'
}

// 弹窗相关
const showModal = ref(false)
const currentImage = ref({})
const imageLoading = ref(false)

// 计算可见页码
const visiblePages = computed(() => {
    const count = 5
    const pages = []
    let start = currentPage.value - Math.floor(count / 2)
    start = Math.max(1, start)
    let end = start + count - 1
    end = Math.min(totalPages.value, end)
    if (end - start + 1 < count) {
        start = Math.max(1, end - count + 1)
    }
    for (let i = start; i <= end; i++) {
        pages.push(i)
    }
    return pages
})

// 生成缩略图 URL
const getThumbUrl = (url) => {
    if (!url) return ''
    const lastDotIndex = url.lastIndexOf('.')
    if (lastDotIndex === -1) return url
    return url.substring(0, lastDotIndex) + '_thumb' + url.substring(lastDotIndex)
}

const handleThumbError = (event) => {
    const img = event.target
    const originalUrl = img.dataset.original
    if (img.src.includes(originalUrl)) {
        img.src = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjAwIiBoZWlnaHQ9IjIwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iMTAwJSIgaGVpZ2h0PSIxMDAlIiBmaWxsPSIjZGRkIi8+PC9zdmc+'
    } else {
        img.src = originalUrl
    }
}

// 加载图片
const loadImages = async (isLoadMore = false) => {
    if (loading.value) return
    loading.value = true
    try {
        const params = new URLSearchParams({
            page: currentPage.value,
            limit: pageSize.value,
            sort_by: 'created_at',
            sort_order: 'desc',
            search: searchQuery.value,
            category: selectedCategory.value === '全部' ? '' : selectedCategory.value
        })
        const response = await fetch(`/api/images?${params}`, {
            credentials: 'include'
        })
        
        if (response.ok) {
            const result = await response.json()
            const newImages = result.data.images || []
            if (isLoadMore) {
                images.value = [...images.value, ...newImages]
            } else {
                images.value = newImages
            }
            totalPages.value = result.data.total_pages || 1
            jumpPageInput.value = ''
        }
    } catch (error) {
        console.error('加载图片错误:', error)
    } finally {
        loading.value = false
    }
}

// 加载轮播图 (随机)
const loadCarouselImages = async () => {
    try {
        const response = await fetch('/api/images?limit=10&sort_by=random', {
             credentials: 'include'
        })
        if (response.ok) {
            const result = await response.json()
            carouselImages.value = result.data.images || []
            if (carouselImages.value.length > 0) {
                startCarousel()
            }
        }
    } catch (error) {
        console.error('加载轮播图失败', error)
    }
}

const startCarousel = () => {
    if (carouselInterval) clearInterval(carouselInterval)
    carouselInterval = setInterval(() => {
        nextSlide()
    }, 5000)
}

const nextSlide = () => {
    currentCarouselIndex.value = (currentCarouselIndex.value + 1) % carouselImages.value.length
}

const prevSlide = () => {
    currentCarouselIndex.value = (currentCarouselIndex.value - 1 + carouselImages.value.length) % carouselImages.value.length
}

// 分页与筛选
const changePage = (page) => {
    if (page >= 1 && page <= totalPages.value) {
        currentPage.value = page
        images.value = []
        hasScrolledPages.value = false // 重置滚动标记
        loadImages(false)
    }
}

const handleJumpPage = () => {
    const page = parseInt(jumpPageInput.value)
    if (page && page >= 1 && page <= totalPages.value) {
        changePage(page)
    } else {
        Message.error('请输入有效的页码')
    }
}

const handleSearch = () => {
    currentPage.value = 1
    images.value = []
    loadImages(false)
}

const selectCategory = (cat) => {
    selectedCategory.value = cat
    currentPage.value = 1
    images.value = []
    hasScrolledPages.value = false // 重置滚动标记
    loadImages(false)
    // 滚动到顶部
    window.scrollTo({ top: 0, behavior: 'smooth' })
}

// 弹窗操作
const openImageModal = (image) => {
    currentImage.value = image
    imageLoading.value = true
    showModal.value = true
    document.body.style.overflow = 'hidden'
}

const closeModal = () => {
    showModal.value = false
    document.body.style.overflow = ''
}

const checkUrlExists = async (url) => {
    try {
        const response = await fetch(url, { method: 'HEAD' })
        return response.ok
    } catch (error) {
        return false
    }
}

const copyLink = async (url, isPreview = false) => {
    let targetUrl = url
    if (isPreview) {
        const exists = await checkUrlExists(url)
        if (!exists) {
            targetUrl = getFullUrl(currentImage.value.url)
            Message.warning('压缩图不存在，已复制原图链接')
        }
    }
    
    navigator.clipboard.writeText(targetUrl).then(() => {
        Message.success('链接已复制')
    })
}

const copyHtml = async (url, isPreview = false) => {
    let targetUrl = url
    if (isPreview) {
        const exists = await checkUrlExists(url)
        if (!exists) {
            targetUrl = getFullUrl(currentImage.value.url)
            Message.warning('压缩图不存在，已复制原图代码')
        }
    }

    const html = `<img src="${targetUrl}" alt="image" />`
    navigator.clipboard.writeText(html).then(() => {
        Message.success('HTML代码已复制')
    })
}

const copyMarkdown = async (url, isPreview = false) => {
    let targetUrl = url
    if (isPreview) {
        const exists = await checkUrlExists(url)
        if (!exists) {
            targetUrl = getFullUrl(currentImage.value.url)
            Message.warning('压缩图不存在，已复制原图代码')
        }
    }

    const md = `![](${targetUrl})`
    navigator.clipboard.writeText(md).then(() => {
        Message.success('Markdown代码已复制')
    })
}

const downloadImage = async (image) => {
    const imageUrl = getFullUrl(image.url)
    
    if (imageUrl.toLowerCase().endsWith('.webp')) {
        try {
            const response = await fetch(imageUrl)
            const blob = await response.blob()
            const bitmap = await createImageBitmap(blob)
            
            const canvas = document.createElement('canvas')
            canvas.width = bitmap.width
            canvas.height = bitmap.height
            
            const ctx = canvas.getContext('2d')
            ctx.drawImage(bitmap, 0, 0)
            
            canvas.toBlob((pngBlob) => {
                const url = URL.createObjectURL(pngBlob)
                const link = document.createElement('a')
                link.href = url
                link.download = image.filename.replace(/\.webp$/i, '.png')
                document.body.appendChild(link)
                link.click()
                document.body.removeChild(link)
                URL.revokeObjectURL(url)
            }, 'image/png')
            return
        } catch (e) {
            console.error('WebP转换PNG失败，降级为直接下载', e)
        }
    }

    const link = document.createElement('a')
    link.href = imageUrl
    link.download = image.filename
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
}

const deleteImage = async (image) => {
    if (!confirm('确定要删除这张图片吗?')) return
    
    try {
        const response = await fetch(`/api/images/${image.id}`, {
            method: 'DELETE',
            credentials: 'include'
        })
        
        if (response.ok) {
            Message.success('删除成功')
            closeModal()
            loadImages(false)
        } else {
            Message.error('删除失败')
        }
    } catch (error) {
        Message.error('网络错误')
    }
}

const formatFileSize = (bytes) => {
    if (!bytes) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const formatDate = (dateString) => {
    if (!dateString) return ''
    return new Date(dateString).toLocaleDateString('zh-CN')
}

// 无限滚动观察器
const setupIntersectionObserver = () => {
    // 先断开旧的observer
    if (observer) {
        observer.disconnect()
        observer = null
    }
    
    // 确保触发器元素存在
    if (!loadMoreTrigger.value) {
        return
    }
    
    observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting && !loading.value && currentPage.value < totalPages.value) {
                console.log('触发无限滚动，当前页:', currentPage.value)
                hasScrolledPages.value = true // 标记已经滚动加载
                currentPage.value++
                loadImages(true)
            }
        })
    }, { 
        root: null,
        rootMargin: '600px',
        threshold: 0
    })

    observer.observe(loadMoreTrigger.value)
    console.log('Observer已设置')
}

onMounted(() => {
    loadImages(false)
    loadCarouselImages()
})

// 监听images变化，每次加载完都重新设置observer
watch(() => images.value.length, () => {
    nextTick(() => {
        setupIntersectionObserver()
    })
}, { immediate: false })

// 监听路由查询参数变化（导航栏搜索）
watch(() => route.query.search, (newSearch) => {
    if (newSearch !== undefined) {
        searchQuery.value = newSearch || ''
        currentPage.value = 1
        images.value = []
        loadImages(false)
    }
}, { immediate: true })

onUnmounted(() => {
    if (carouselInterval) clearInterval(carouselInterval)
    if (observer) observer.disconnect()
})
</script>

<style scoped>
@keyframes fade-in {
    from {
        opacity: 0;
    }
    to {
        opacity: 1;
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
    animation: fade-in 0.3s ease-out;
}

.animate-scale-in {
    animation: scale-in 0.3s ease-out;
}

/* 隐藏数字输入框的上下箭头 */
input[type="number"]::-webkit-inner-spin-button,
input[type="number"]::-webkit-outer-spin-button {
    -webkit-appearance: none;
    margin: 0;
}

input[type="number"] {
    appearance: textfield;
    -moz-appearance: textfield;
}

/* 分类标签动画 */
@keyframes shimmer {
    0% {
        transform: translateX(-100%);
    }
    100% {
        transform: translateX(100%);
    }
}

.animate-shimmer {
    animation: shimmer 2s infinite;
}

.category-btn {
    backdrop-filter: blur(8px);
}

.category-btn::before {
    content: '';
    position: absolute;
    inset: 0;
    border-radius: inherit;
    opacity: 0;
    transition: opacity 0.3s;
    background: radial-gradient(circle at var(--x, 50%) var(--y, 50%), rgba(255,255,255,0.15), transparent 50%);
}

.category-btn:hover::before {
    opacity: 1;
}

/* 轮播图旋转动画 */
@keyframes slow-spin {
    from {
        transform: rotate(0deg);
    }
    to {
        transform: rotate(360deg);
    }
}

@keyframes slow-spin-reverse {
    from {
        transform: rotate(360deg);
    }
    to {
        transform: rotate(0deg);
    }
}

.animate-slow-spin {
    animation: slow-spin 30s linear infinite;
}

.animate-slow-spin-reverse {
    animation: slow-spin-reverse 25s linear infinite;
}

/* 隐藏滚动条 */
.scrollbar-hide {
    -ms-overflow-style: none;
    scrollbar-width: none;
}

.scrollbar-hide::-webkit-scrollbar {
    display: none;
}
</style>
