<template>
  <div class="pt-4 px-4 md:px-6 lg:px-8 xl:container xl:mx-auto">
    <div class="max-w-3xl mx-auto">
      <h2 class="text-2xl font-bold mb-6 text-secondary">AI 功能配置</h2>
      
      <div class="bg-white dark:bg-dark-200 rounded-xl shadow-md p-6 mb-6">
        <h3 class="text-lg font-semibold mb-4">模型设置</h3>
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-secondary mb-1">API URL</label>
            <input 
              v-model="aiConfig.api_url" 
              type="text" 
              placeholder="例如: https://api.openai.com"
              class="w-full px-4 py-2 rounded-lg border border-light-300 dark:border-dark-100 bg-light-50 dark:bg-dark-300 focus:outline-none focus:ring-2 focus:ring-primary"
            />
          </div>
          
          <div>
            <label class="block text-sm font-medium text-secondary mb-1">API Key</label>
            <input 
              v-model="aiConfig.api_key" 
              type="password" 
              placeholder="sk-..."
              class="w-full px-4 py-2 rounded-lg border border-light-300 dark:border-dark-100 bg-light-50 dark:bg-dark-300 focus:outline-none focus:ring-2 focus:ring-primary"
            />
          </div>
          
          <div>
            <label class="block text-sm font-medium text-secondary mb-1">模型名称</label>
            <div class="flex gap-2">
              <div class="relative flex-1">
                <input 
                  v-model="aiConfig.model" 
                  type="text" 
                  placeholder="例如: gpt-3.5-turbo"
                  class="w-full px-4 py-2 rounded-lg border border-light-300 dark:border-dark-100 bg-light-50 dark:bg-dark-300 focus:outline-none focus:ring-2 focus:ring-primary"
                  list="model-list"
                />
                <datalist id="model-list">
                  <option v-for="m in availableModels" :key="m.id" :value="m.id">{{ m.id }}</option>
                </datalist>
              </div>
              <button 
                @click="fetchModels" 
                class="px-3 py-2 bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 rounded-lg transition-colors text-secondary flex items-center justify-center min-w-[42px]"
                title="获取可用模型"
                :disabled="fetchingModels"
              >
                <i class="ri-refresh-line text-lg" :class="{ 'animate-spin': fetchingModels }"></i>
              </button>
            </div>
          </div>
          
          <div class="pt-4">
            <button 
              @click="saveConfig" 
              class="px-6 py-2 bg-primary hover:bg-primary-dark text-white rounded-lg transition-colors duration-200"
              :disabled="saving"
            >
              {{ saving ? '保存中...' : '保存配置' }}
            </button>
          </div>
        </div>
      </div>

      <div class="bg-white dark:bg-dark-200 rounded-xl shadow-md p-6">
        <h3 class="text-lg font-semibold mb-4">批量操作</h3>
        <div class="flex items-center justify-between mb-4">
          <div>
            <p class="font-medium">AI 智能打标签</p>
            <p class="text-sm text-secondary mt-1">使用配置的 AI 模型对现有图片进行分析，生成中文标签并分类（不会修改文件名）</p>
          </div>
          <button 
            @click="startBatchTag"
            class="px-4 py-2 bg-secondary hover:bg-secondary-dark text-white rounded-lg transition-colors duration-200"
            :disabled="tagging || progress.is_running"
          >
            {{ tagging || progress.is_running ? '处理中...' : '开始打标签' }}
          </button>
        </div>

        <!-- 进度条 -->
        <div v-if="progress.is_running || progress.total > 0" class="mt-4 p-4 bg-light-50 dark:bg-dark-300 rounded-lg">
            <div class="flex justify-between text-sm mb-2">
                <span>处理进度: {{ progress.current }} / {{ progress.total }}</span>
                <span>{{ percentage }}%</span>
            </div>
            <div class="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2.5">
                <div class="bg-primary h-2.5 rounded-full transition-all duration-300" :style="{ width: percentage + '%' }"></div>
            </div>
            <p v-if="!progress.is_running && progress.total > 0" class="text-green-500 text-sm mt-2">
                任务已完成！
            </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, onUnmounted } from 'vue'

const aiConfig = ref({
  api_url: '',
  api_key: '',
  model: ''
})

const saving = ref(false)
const tagging = ref(false)
const fetchingModels = ref(false)
const availableModels = ref([])
const progress = ref({
    total: 0,
    current: 0,
    is_running: false
})
let pollTimer = null

const percentage = computed(() => {
    if (progress.value.total === 0) return 0
    return Math.round((progress.value.current / progress.value.total) * 100)
})

const fetchModels = async () => {
    if (!aiConfig.value.api_url || !aiConfig.value.api_key) {
        alert('请先填写 API URL 和 API Key，并保存配置')
        return
    }
    
    fetchingModels.value = true
    try {
        const response = await fetch('/api/ai/models', {
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('authToken')}`
            }
        })
        
        if (response.ok) {
            const result = await response.json()
            if (result.code === 200 && Array.isArray(result.data)) {
                availableModels.value = result.data
                alert(`成功获取 ${result.data.length} 个模型`)
            } else {
                alert(result.msg || '获取模型失败')
            }
        } else {
            const result = await response.json()
            alert(result.msg || '获取模型失败')
        }
    } catch (error) {
        console.error('获取模型失败', error)
        alert('网络请求失败')
    } finally {
        fetchingModels.value = false
    }
}

const loadConfig = async () => {
  try {
    const response = await fetch('/api/settings/ai', {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('authToken')}`
      }
    })
    if (response.ok) {
      const result = await response.json()
      if (result.code === 200) {
        aiConfig.value = result.data
      }
    }
  } catch (error) {
    console.error('加载配置失败', error)
  }
}

const saveConfig = async () => {
  saving.value = true
  try {
    const response = await fetch('/api/settings/ai', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('authToken')}`
      },
      body: JSON.stringify(aiConfig.value)
    })
    if (response.ok) {
      alert('配置已保存')
    } else {
      alert('保存失败')
    }
  } catch (error) {
    console.error('保存配置失败', error)
    alert('保存失败')
  } finally {
    saving.value = false
  }
}

const startBatchTag = async () => {
  if (!confirm('确定要开始批量打标签吗？这可能需要很长时间。')) return
  
  tagging.value = true
  try {
    const response = await fetch('/api/batch-tag', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('authToken')}`
      }
    })
    if (response.ok) {
      // alert('AI 打标签任务已在后台启动')
      startPolling()
    } else {
      alert('启动失败')
    }
  } catch (error) {
    console.error('启动任务失败', error)
    alert('启动失败')
  } finally {
    tagging.value = false
  }
}

const fetchProgress = async () => {
    try {
        const response = await fetch('/api/ai/progress', {
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('authToken')}`
            }
        })
        if (response.ok) {
            const result = await response.json()
            if (result.code === 200) {
                progress.value = result.data
                if (!result.data.is_running && result.data.total > 0 && result.data.current === result.data.total) {
                    stopPolling()
                }
            }
        }
    } catch (error) {
        console.error('获取进度失败', error)
    }
}

const startPolling = () => {
    if (pollTimer) clearInterval(pollTimer)
    fetchProgress() // 立即执行一次
    pollTimer = setInterval(fetchProgress, 2000)
}

const stopPolling = () => {
    if (pollTimer) clearInterval(pollTimer)
    pollTimer = null
}

onMounted(() => {
  loadConfig()
  // 页面加载时检查是否有正在运行的任务
  fetchProgress()
  startPolling()
})

onUnmounted(() => {
    stopPolling()
})
</script>
