<template>
    <div class="text-gray-800 dark:text-gray-200">
        <!-- 页面头部 -->
        <div class="settings-header container mx-auto px-4 py-8 mb-6">
            <h1 class="page-title flex items-center text-2xl md:text-3xl font-bold">
                设置
            </h1>
            <p class="page-description text-gray-600 dark:text-gray-400 mt-2">管理您的账户</p>
        </div>

        <!-- 主要内容 -->
        <div class="settings-content container mx-auto px-4 pb-16">
            <div class="settings-panel max-w-2xl mx-auto bg-white dark:bg-gray-800 rounded-xl shadow-md overflow-hidden">
                <div class="panel-content p-6 md:p-8">
                    <h2 class="panel-title flex items-center text-xl font-semibold mb-8">
                        <span class="panel-icon mr-2 text-2xl">
                            <i class="ri-shield-user-line"></i>
                        </span>
                        账户设置
                    </h2>
                    
                    <!-- 账户修改表单 -->
                    <form @submit.prevent="updateAccount" class="account-form space-y-6">
                        <!-- 新用户名 -->
                        <div class="setting-group">
                            <label 
                                class="setting-label block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" 
                                for="newUsername"
                            >
                                新用户名（留空则不修改）
                            </label>
                            <input 
                                id="newUsername"
                                v-model="accountForm.newUsername"
                                type="text" 
                                class="setting-input w-full px-4 py-2.5 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 focus:ring-2 focus:ring-primary focus:border-primary dark:focus:ring-primary/70 dark:focus:border-primary/70 transition-colors outline-none"
                                placeholder="留空则不修改用户名"
                                minlength="3"
                                maxlength="20"
                            />
                        </div>
                        
                        <!-- 当前密码 -->
                        <div class="setting-group">
                            <label 
                                class="setting-label block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" 
                                for="currentPassword"
                            >
                                当前密码 <span class="text-red-500">*</span>
                            </label>
                            <input 
                                id="currentPassword"
                                v-model="accountForm.currentPassword"
                                type="password" 
                                class="setting-input w-full px-4 py-2.5 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 focus:ring-2 focus:ring-primary focus:border-primary dark:focus:ring-primary/70 dark:focus:border-primary/70 transition-colors outline-none"
                                placeholder="请输入当前密码以确认修改"
                                required
                            />
                        </div>
                        
                        <!-- 新密码 -->
                        <div class="setting-group">
                            <label 
                                class="setting-label block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" 
                                for="newPassword"
                            >
                                新密码（留空则不修改）
                            </label>
                            <input 
                                id="newPassword"
                                v-model="accountForm.newPassword"
                                type="password" 
                                class="setting-input w-full px-4 py-2.5 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 focus:ring-2 focus:ring-primary focus:border-primary dark:focus:ring-primary/70 dark:focus:border-primary/70 transition-colors outline-none"
                                placeholder="留空则不修改密码（至少6位）"
                                minlength="6"
                            />
                        </div>
                        
                        <!-- 确认新密码 -->
                        <div class="setting-group">
                            <label 
                                class="setting-label block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1" 
                                for="confirmPassword"
                            >
                                确认新密码
                            </label>
                            <input 
                                id="confirmPassword"
                                v-model="accountForm.confirmPassword"
                                type="password" 
                                class="setting-input w-full px-4 py-2.5 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 focus:ring-2 focus:ring-primary focus:border-primary dark:focus:ring-primary/70 dark:focus:border-primary/70 transition-colors outline-none"
                                placeholder="请再次输入新密码"
                            />
                        </div>
                        
                        <!-- 提交按钮 -->
                        <div class="setting-group pt-2">
                            <button 
                                type="submit" 
                                :disabled="isUpdatingAccount"
                                class="setting-btn accent w-full py-3 px-6 bg-primary hover:bg-primary/90 text-white font-medium rounded-lg transition-colors flex items-center justify-center gap-2 focus:ring-2 focus:ring-primary/50 focus:outline-none"
                            >
                                <span v-if="isUpdatingAccount" class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></span>
                                <span>保存修改</span>
                            </button>
                        </div>
                    </form>
                </div>
            </div>

            <!-- 系统维护面板 -->
            <div class="settings-panel max-w-2xl mx-auto bg-white dark:bg-gray-800 rounded-xl shadow-md overflow-hidden mt-6">
                <div class="panel-content p-6 md:p-8">
                    <h2 class="panel-title flex items-center text-xl font-semibold mb-8">
                        <span class="panel-icon mr-2 text-2xl">
                            <i class="ri-tools-line"></i>
                        </span>
                        系统维护
                    </h2>
                    
                    <div class="setting-group flex items-center justify-between">
                        <div>
                            <h3 class="font-medium text-gray-800 dark:text-gray-200">图片去重</h3>
                            <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">扫描并删除重复上传的图片（保留最早的一张）</p>
                        </div>
                        <button 
                            @click="handleDeduplicate"
                            :disabled="isDeduplicating"
                            class="px-4 py-2 bg-secondary hover:bg-secondary-dark text-white rounded-lg transition-colors duration-200 disabled:opacity-50"
                        >
                            {{ isDeduplicating ? '处理中...' : '开始去重' }}
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import message from '@/utils/message.js'

const router = useRouter()

// 表单数据
const accountForm = ref({
    newUsername: '',
    currentPassword: '',
    newPassword: '',
    confirmPassword: ''
})

// 加载状态
const isUpdatingAccount = ref(false)
const isDeduplicating = ref(false)

// 图片去重
const handleDeduplicate = async () => {
    if (!confirm('确定要开始去重吗？这将删除重复的图片文件。')) return
    
    isDeduplicating.value = true
    try {
        const response = await fetch('/api/deduplicate', {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${localStorage.getItem('authToken')}`
            }
        })
        const result = await response.json()
        if (response.ok) {
            message.success(result.msg || '去重完成')
        } else {
            message.error(result.msg || '去重失败')
        }
    } catch (error) {
        console.error('去重错误:', error)
        message.error('去重请求失败')
    } finally {
        isDeduplicating.value = false
    }
}

// 更新账户信息
const updateAccount = async () => {
    const { newUsername, currentPassword, newPassword, confirmPassword } = accountForm.value
    
    // 检查是否有任何修改
    const hasUsernameChange = newUsername && newUsername.trim() !== ''
    const hasPasswordChange = newPassword && newPassword.trim() !== ''
    
    if (!hasUsernameChange && !hasPasswordChange) {
        message.error('请输入要修改的用户名或密码')
        return
    }
    
    // 验证用户名（如果要修改）
    if (hasUsernameChange) {
        if (newUsername.length < 3) {
            message.error('用户名长度至少为3位')
            return
        }
        
        if (newUsername.length > 20) {
            message.error('用户名长度不能超过20位')
            return
        }
    }
    
    // 验证密码（如果要修改）
    if (hasPasswordChange) {
        if (newPassword.length < 6) {
            message.error('新密码长度至少为6位')
            return
        }
        
        if (newPassword !== confirmPassword) {
            message.error('两次输入的新密码不一致')
            return
        }
    }
    
    // 验证当前密码
    if (!currentPassword) {
        message.error('请输入当前密码以确认修改')
        return
    }
    
    try {
        isUpdatingAccount.value = true
        
        const response = await fetch('/api/account/change', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('authToken')}`
            },
            body: JSON.stringify({
                new_username: newUsername,
                current_password: currentPassword,
                new_password: newPassword
            })
        })
        
        const result = await response.json()
        
        if (!response.ok || !result.success) {
            // 未授权处理
            if (response.status === 401) {
                localStorage.removeItem('authToken')
                router.push('/login')
                return message.error('登录已过期，请重新登录')
            }
            throw new Error(result.message || '修改失败')
        }
        
        message.success('修改成功')

        // 清空表单
        accountForm.value = {
            newUsername: '',
            currentPassword: '',
            newPassword: '',
            confirmPassword: ''
        }
        
        // 刷新页面
        setTimeout(() => {
            window.location.reload();
        }, 1000)

    } catch (error) {
        message.error(error.message || '更新失败')
    } finally {
        isUpdatingAccount.value = false
    }
}
</script>