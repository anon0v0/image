<template>
  <div class="pt-8 px-4 md:px-6 lg:px-8 max-w-7xl mx-auto min-h-screen">
    <!-- 上传区域 -->
    <section class="upload-section mb-12 animate-fade-in-up">
      <div class="bg-white dark:bg-gray-800 rounded-3xl shadow-xl p-8 md:p-12 border border-gray-100 dark:border-gray-700 transition-all duration-300 hover:shadow-2xl">
        <div class="text-center mb-8">
            <h2 class="text-3xl font-bold text-gray-900 dark:text-white mb-2">上传图片</h2>
            <p class="text-gray-500 dark:text-gray-400">支持拖拽、点击或粘贴图片</p>
        </div>

        <!-- 拖拽上传区域 -->
        <div 
          class="upload-area relative rounded-2xl border-2 border-dashed transition-all duration-300 cursor-pointer overflow-hidden group"
          :class="{ 
            'border-primary bg-primary/5': isDragOver,
            'border-gray-300 dark:border-gray-600 hover:border-primary dark:hover:border-primary hover:bg-gray-50 dark:hover:bg-gray-700/50': !isDragOver && !isUploading,
            'border-primary/50 bg-primary/5': isUploading
          }"
          @drop="handleDrop"
          @dragover="handleDragOver"
          @dragenter="handleDragEnter"
          @dragleave="handleDragLeave"
          @click="triggerFileInput"
        >
          <!-- 未上传状态 -->
          <div v-if="!isUploading" class="upload-content py-20 px-4 text-center">
            <div class="w-20 h-20 bg-primary/10 rounded-full flex items-center justify-center mx-auto mb-6 group-hover:scale-110 transition-transform duration-300">
                <i class="ri-upload-cloud-2-line text-4xl text-primary"></i>
            </div>
            <h3 class="text-xl font-bold text-gray-800 dark:text-gray-200 mb-2">点击或拖拽上传</h3>
            <p class="text-gray-500 dark:text-gray-400 text-sm mb-8 max-w-md mx-auto">支持 JPG, PNG, GIF, WebP 等常见格式，单张最大 30MB</p>
            
            <button class="px-8 py-3 bg-primary hover:bg-primary-dark text-white rounded-xl font-bold shadow-lg shadow-primary/30 hover:shadow-primary/50 transition-all transform hover:-translate-y-0.5">
              选择文件
            </button>
            <p class="mt-6 text-xs text-gray-400 flex items-center justify-center gap-1">
              <i class="ri-keyboard-line"></i> 支持 Ctrl+V 粘贴
            </p>
          </div>

          <!-- 上传进度状态 -->
          <div v-else class="upload-progress py-20 px-4 text-center">
            <div class="relative w-24 h-24 mx-auto mb-6">
                <svg class="w-full h-full transform -rotate-90">
                    <circle cx="48" cy="48" r="40" stroke="currentColor" stroke-width="8" fill="transparent" class="text-gray-200 dark:text-gray-700" />
                    <circle cx="48" cy="48" r="40" stroke="currentColor" stroke-width="8" fill="transparent" :stroke-dasharray="251.2" :stroke-dashoffset="251.2 - (251.2 * uploadProgress) / 100" class="text-primary transition-all duration-300 ease-out" />
                </svg>
                <div class="absolute inset-0 flex items-center justify-center text-lg font-bold text-primary">
                    {{ Math.round(uploadProgress) }}%
                </div>
            </div>
            <p class="text-gray-600 dark:text-gray-300 font-medium mb-2">正在上传 {{ uploadingCount }} 张图片...</p>
            <p class="text-sm text-gray-400">请稍候，精彩即将呈现</p>
          </div>
        </div>

        <!-- 隐藏的文件输入 -->
        <input 
          ref="fileInput"
          type="file"
          multiple
          accept="image/*"
          @change="handleFileSelect"
          class="hidden"
        />
      </div>
    </section>

    <!-- 最近上传的图片 -->
    <section class="recent-section pb-12 animate-fade-in-up" style="animation-delay: 0.1s">
      <div class="flex justify-between items-center mb-6 px-2">
        <h2 class="text-xl font-bold text-gray-900 dark:text-white flex items-center gap-2">
          <i class="ri-history-line text-primary"></i>
          最近上传
        </h2>
        <span class="px-3 py-1 bg-gray-100 dark:bg-gray-800 rounded-full text-xs font-medium text-gray-500 dark:text-gray-400">{{ recentImages.length }} 张图片</span>
      </div>

      <!-- 图片网格 -->
      <div v-if="recentImages.length > 0" class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-4">
        <div 
          v-for="image in recentImages" 
          :key="image.id"
          class="group relative aspect-square rounded-2xl overflow-hidden bg-gray-100 dark:bg-gray-800 shadow-sm hover:shadow-lg transition-all duration-300 cursor-pointer"
          @click="previewImage(image)"
        >
          <img 
            :src="getFullUrl(image.thumbnailUrl || image.url)"
            :alt="image.filename" 
            class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-110"
            loading="lazy"
          />
          
          <!-- 悬停操作栏 -->
          <div class="absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity duration-300 flex items-center justify-center gap-3 backdrop-blur-[2px]">
            <button 
                class="w-10 h-10 rounded-full bg-white/20 hover:bg-white text-white hover:text-primary flex items-center justify-center transition-all duration-200 backdrop-blur-md"
                title="复制链接"
                @click.stop="copyImageLink(image)"
              >
                <i class="ri-file-copy-line text-lg"></i>
            </button>
            <button 
              @click.stop="downloadImage(image)"
              class="w-10 h-10 rounded-full bg-white/20 hover:bg-white text-white hover:text-primary flex items-center justify-center transition-all duration-200 backdrop-blur-md"
              title="下载图片"
            >
              <i class="ri-download-cloud-2-line text-lg"></i>
            </button>
          </div>
          
          <!-- 文件名标签 -->
          <div class="absolute bottom-0 left-0 right-0 p-2 bg-gradient-to-t from-black/80 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300 pointer-events-none">
              <p class="text-white text-[10px] truncate text-center">{{ image.filename }}</p>
          </div>
        </div>
      </div>

      <!-- 无图片状态 -->
      <div v-else class="bg-white dark:bg-gray-800 rounded-3xl shadow-sm p-12 text-center border border-dashed border-gray-200 dark:border-gray-700">
        <div class="text-6xl text-gray-200 dark:text-gray-700 mb-4">
          <i class="ri-image-add-line"></i>
        </div>
        <p class="text-gray-500 dark:text-gray-400">暂无上传记录</p>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import Message from '@/utils/message.js'

// 核心：获取完整URL的函数
const getFullUrl = (path) => {
  if (!path) return ''
  if (typeof window !== 'undefined') {
    return window.location.origin + path
  }
  return path
}

// 响应式数据
const isDragOver = ref(false)
const isUploading = ref(false)
const uploadingCount = ref(0)
const uploadProgress = ref(0)
const recentImages = ref([])
const fileInput = ref(null)

// 拖拽处理
const handleDragOver = (e) => {
  e.preventDefault()
  isDragOver.value = true
}

const handleDragEnter = (e) => {
  e.preventDefault()
  isDragOver.value = true
}

const handleDragLeave = (e) => {
  e.preventDefault()
  if (!e.currentTarget.contains(e.relatedTarget)) {
    isDragOver.value = false
  }
}

const handleDrop = (e) => {
  e.preventDefault()
  isDragOver.value = false
  
  const files = Array.from(e.dataTransfer.files)
  const imageFiles = files.filter(file => file.type.startsWith('image/'))
  
  if (imageFiles.length > 0) {
    uploadFiles(imageFiles)
  } else {
    Message.warning('请拖拽图片文件')
  }
}

// 文件选择处理
const triggerFileInput = () => {
  if (!isUploading.value && fileInput.value) {
    fileInput.value.click()
  }
}

const handleFileSelect = (e) => {
  const files = Array.from(e.target.files)
  if (files.length > 0) {
    uploadFiles(files)
  }
  e.target.value = ''
}

// 剪贴板粘贴处理
const handlePaste = async (e) => {
  const items = e.clipboardData.items
  const imageFiles = []
  
  for (let item of items) {
    if (item.type.startsWith('image/')) {
      const file = item.getAsFile()
      if (file) {
        const timestamp = new Date().getTime()
        const extension = item.type.split('/')[1] || 'png'
        const newFile = new File([file], `paste-${timestamp}.${extension}`, {
          type: item.type
        })
        imageFiles.push(newFile)
      }
    }
  }
  
  if (imageFiles.length > 0) {
    e.preventDefault()
    uploadFiles(imageFiles)
  }
}

// 文件上传
const uploadFiles = async (files) => {
  if (isUploading.value) return
  
  isUploading.value = true
  uploadingCount.value = files.length
  uploadProgress.value = 0
  
  const formData = new FormData()
  files.forEach(file => {
    formData.append('images[]', file)
  })
  
  try {
    const progressInterval = setInterval(() => {
      if (uploadProgress.value < 95) {
        uploadProgress.value += Math.random() * 5
      }
    }, 150)
    
    const response = await fetch('/api/upload/images', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('authToken')}`
      },
      body: formData
    })
    
    clearInterval(progressInterval)
    uploadProgress.value = 100
    
    const result = await response.json()
    
    if (response.ok) {
      await loadRecentImages()
      Message.success('上传成功')
    } else {
      throw new Error(result.message || '上传失败')
    }
  } catch (error) {
    console.error('上传错误:', error)
    Message.error(`上传失败: ${error.message}`)
  } finally {
    isUploading.value = false
    uploadingCount.value = 0
    uploadProgress.value = 0
  }
}

// 加载最近上传的图片
const loadRecentImages = async () => {
  try {
    const response = await fetch('/api/images?limit=12', {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('authToken')}`
      }
    })
    
    if (response.ok) {
      const result = await response.json()
      recentImages.value = Array.isArray(result.data?.images) ? result.data.images : []
    }
  } catch (error) {
    console.error('加载图片失败:', error)
    recentImages.value = []
  }
}

const copyImageLink = async (image) => {
  if (!image) return
  const fullUrl = getFullUrl(image.url)
  try {
    await navigator.clipboard.writeText(fullUrl)
    Message.success('已复制链接')
  } catch (error) {
    console.error('复制失败', error)
    Message.error('复制失败')
  }
}

const downloadImage = async (image) => {
  if (!image || !image.url) return
  const fullUrl = getFullUrl(image.url)
  
  // 如果是 WebP，尝试转换为 PNG
  if (fullUrl.toLowerCase().endsWith('.webp')) {
      try {
          const response = await fetch(fullUrl)
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
              // 将文件名后缀改为 .png
              link.download = (image.filename || `image-${Date.now()}.webp`).replace(/\.webp$/i, '.png')
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
  link.href = fullUrl
  link.download = image.filename || `image-${Date.now()}.png`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

const previewImage = (image) => {
    window.open(getFullUrl(image.url), '_blank')
}

// 生命周期
onMounted(() => {
  document.addEventListener('paste', handlePaste)
  loadRecentImages()
})

onUnmounted(() => {
  document.removeEventListener('paste', handlePaste)
})
</script>
