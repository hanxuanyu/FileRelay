<template>
  <div class="w-full h-full">
    <div class="container mx-auto px-3 sm:px-4 pb-6 sm:pb-12 flex items-center h-full">
        <div class="max-w-2xl mx-auto w-full">
        
        <!-- 主卡片 -->
        <Card class="shadow-2xl border-0 overflow-hidden">
          <CardContent class="pt-8 sm:pt-12 pb-8 sm:pb-10 px-4 sm:px-6">
            
            <!-- 顶部标题 -->
            <div class="text-center mb-8 sm:mb-10">
              <div class="inline-flex items-center justify-center w-16 h-16 sm:w-20 sm:h-20 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 mb-4 sm:mb-6 shadow-lg">
                <svg class="w-8 h-8 sm:w-10 sm:h-10 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"></path>
                </svg>
              </div>
              <h2 class="text-2xl sm:text-3xl font-bold text-gray-900 dark:text-gray-100 mb-2 sm:mb-3">{{ t('upload.title') }}</h2>
              <p class="text-sm sm:text-base text-gray-600 dark:text-gray-400 px-4 sm:px-0">{{ t('upload.description') }}</p>
            </div>

            <!-- 上传类型切换 -->
            <div class="max-w-md mx-auto mb-6 sm:mb-8 px-2 sm:px-0">
              <div ref="tabGroup" class="relative inline-flex w-full items-center rounded-full bg-gray-100 dark:bg-gray-800 p-1 transition-colors">
                <!-- 滑动指示器 -->
                <div 
                  :class="[
                    'absolute inset-y-1 bg-white dark:bg-gray-700 rounded-full shadow-sm',
                    tabIndicatorInitialized && 'transition-all duration-300 ease-in-out'
                  ]"
                  :style="{ 
                    left: tabIndicatorPosition,
                    width: tabIndicatorWidth
                  }"
                ></div>
                
                <button
                  ref="fileTabButton"
                  @click="switchTab('file')"
                  :class="[
                    'relative z-10 flex-1 flex items-center justify-center py-2 sm:py-2.5 px-2 sm:px-4 text-xs sm:text-sm font-semibold rounded-full transition-colors',
                    activeTab === 'file' 
                      ? 'text-blue-600 dark:text-blue-400' 
                      : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-200'
                  ]"
                >
                  <svg class="w-4 h-4 sm:mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
                  </svg>
                  <span class="hidden sm:inline ml-2">{{ t('upload.fileTab') }}</span>
                  <span class="sm:hidden ml-1">{{ t('upload.fileTabMobile') }}</span>
                </button>
                <button
                  ref="textTabButton"
                  @click="switchTab('text')"
                  :class="[
                    'relative z-10 flex-1 flex items-center justify-center py-2 sm:py-2.5 px-2 sm:px-4 text-xs sm:text-sm font-semibold rounded-full transition-colors',
                    activeTab === 'text'
                      ? 'text-purple-600 dark:text-purple-400'
                      : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-200'
                  ]"
                >
                  <svg class="w-4 h-4 sm:mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h8a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
                  </svg>
                  <span class="hidden sm:inline ml-2">{{ t('upload.textTab') }}</span>
                  <span class="sm:hidden ml-1">{{ t('upload.textTabMobile') }}</span>
                </button>
              </div>
            </div>

            <!-- 内容区域 -->
            <div class="space-y-4 sm:space-y-5">
              <!-- 文件上传界面 -->
              <div v-show="activeTab === 'file'" class="space-y-3 sm:space-y-4">
                <!-- 文件选择区域 -->
                <div
                  class="border-2 border-dashed rounded-xl p-4 sm:p-6 text-center transition-all cursor-pointer"
                  :class="{
                    'border-blue-400 bg-gradient-to-br from-blue-50 to-purple-50 dark:from-blue-950/30 dark:to-purple-950/30 shadow-inner': isDragging,
                    'border-gray-300 dark:border-gray-600 hover:border-blue-400 dark:hover:border-blue-500 hover:bg-gray-50 dark:hover:bg-gray-800': !isDragging
                  }"
                  @click="triggerFileInput"
                  @dragover.prevent="isDragging = true"
                  @dragleave.prevent="isDragging = false"
                  @drop.prevent="handleFileDrop"
                >
                  <input
                    ref="fileInput"
                    type="file"
                    multiple
                    class="hidden"
                    @change="handleFileSelect"
                  />
                  
                  <div class="inline-flex items-center justify-center w-10 h-10 sm:w-12 sm:h-12 rounded-full mb-2 sm:mb-3" :class="isDragging ? 'bg-gradient-to-br from-blue-500 to-purple-600' : 'bg-gray-200 dark:bg-gray-700'">
                    <svg class="h-5 w-5 sm:h-6 sm:w-6" :class="isDragging ? 'text-white' : 'text-gray-400 dark:text-gray-500'" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"/>
                    </svg>
                  </div>
                  
                  <p class="text-xs sm:text-sm font-semibold mb-1" :class="isDragging ? 'text-blue-600 dark:text-blue-400' : 'text-gray-900 dark:text-gray-100'">
                    <span class="hidden sm:inline">{{ isDragging ? t('upload.dragging') : t('upload.dragTip') }}</span>
                    <span class="sm:hidden">{{ isDragging ? t('upload.dragging') : t('upload.selectFile', '点击选择文件') }}</span>
                  </p>
                  <p class="text-xs text-gray-500 dark:text-gray-400">
                    <span class="hidden sm:inline">{{ t('upload.fileSizeLimit', { size: publicConfig.getFileSizeLimit() }) }}{{ t('upload.fileCountLimit', { count: config.upload?.max_batch_files || 10 }) }}</span>
                    <span class="sm:hidden">{{ t('upload.fileSizeLimit', { size: publicConfig.getFileSizeLimit() }) }}</span>
                  </p>
                </div>

                <!-- 已选择的文件列表 -->
                <div v-if="selectedFiles.length > 0" class="space-y-2">
                  <div class="flex items-center justify-between">
                    <span class="text-xs font-semibold text-gray-900 dark:text-gray-100">
                      {{ t('upload.selectedFiles', { count: selectedFiles.length }) }} ({{ totalSize }})
                    </span>
                    <Button
                      @click="clearFiles"
                      variant="ghost"
                      size="sm"
                      class="h-6 px-2 text-xs text-gray-600 hover:text-red-600"
                    >
                      <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                      </svg>
                      {{ t('upload.clearAll') }}
                    </Button>
                  </div>
                  
                  <div class="space-y-1.5 max-h-32 overflow-y-auto">
                    <div
                      v-for="(file, index) in selectedFiles"
                      :key="index"
                      class="flex items-center justify-between p-2 bg-gray-50 dark:bg-gray-800 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
                    >
                      <div class="flex items-center space-x-2 flex-1 min-w-0">
                        <div class="w-7 h-7 rounded bg-blue-100 dark:bg-blue-900/30 flex items-center justify-center flex-shrink-0">
                          <component :is="getFileIconComponent(file.name)" class="w-3.5 h-3.5 text-blue-600 dark:text-blue-400" />
                        </div>
                        <div class="flex-1 min-w-0">
                          <p class="text-xs font-medium text-gray-900 dark:text-gray-100 truncate">{{ file.name }}</p>
                          <p class="text-xs text-gray-500 dark:text-gray-400">{{ formatFileSize(file.size) }}</p>
                        </div>
                      </div>
                      
                      <button
                        @click="removeFile(index)"
                        class="p-0.5 text-gray-400 dark:text-gray-500 hover:text-red-600 dark:hover:text-red-500 transition-colors flex-shrink-0"
                      >
                        <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                        </svg>
                      </button>
                    </div>
                  </div>
                </div>
              </div>

              <!-- 文本保存界面 -->
              <div v-show="activeTab === 'text'">
                <Textarea
                  v-model="textContent"
                  :placeholder="t('upload.textPlaceholder')"
                  class="resize-none text-sm rounded-xl"
                  style="height: 140px;"
                />
                
                <div class="mt-2 flex justify-between items-center text-xs text-gray-500">
                  <span class="font-medium">{{ t('upload.textCount', { count: textContent.length }) }}</span>
                  <Button @click="pasteText" variant="ghost" size="sm" class="h-7 text-xs text-gray-600 hover:text-blue-600">
                    <svg class="w-3.5 h-3.5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"></path>
                    </svg>
                    {{ t('common.paste', '粘贴') }}
                  </Button>
                </div>
              </div>

              <!-- 配置区域 -->
              <div v-if="activeTab === 'file' ? selectedFiles.length > 0 : textContent.trim()" class="pt-4 border-t border-gray-200 dark:border-gray-700">
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <!-- 按时间过期 -->
                  <div>
                    <div class="text-xs text-gray-700 dark:text-gray-300 mb-2 font-semibold">{{ t('upload.settings.expireTime') }}</div>
                    <div class="flex flex-wrap gap-1.5">
                      <button
                        v-for="option in publicConfig.getExpireOptions().filter(opt => opt.value > 0)"
                        :key="'time-' + option.value"
                        @click="expireType = 'time'; expireDays = String(option.value); showCustomExpire = false"
                        :class="[
                          'px-2.5 py-1 text-xs font-semibold rounded-full transition-all',
                          expireType === 'time' && expireDays === String(option.value) && !showCustomExpire
                            ? 'bg-gradient-to-r from-blue-500 to-blue-600 text-white shadow-sm'
                            : 'bg-white dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-600 border border-gray-200 dark:border-gray-600'
                        ]"
                      >
                        {{ t(option.labelKey) }}
                      </button>
                      <button
                        @click="toggleCustomExpire"
                        :class="[
                          'px-2.5 py-1 text-xs font-semibold rounded-full transition-all flex items-center gap-1',
                          showCustomExpire
                            ? 'bg-gradient-to-r from-blue-500 to-blue-600 text-white shadow-sm'
                            : 'bg-white dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-600 border border-gray-200 dark:border-gray-600'
                        ]"
                      >
                        <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
                        </svg>
                        {{ t('upload.settings.custom') }}
                      </button>
                    </div>
                    <input
                      v-if="showCustomExpire"
                      v-model="customExpireDays"
                      type="number"
                      min="0.01"
                      :max="config.upload.max_retention_days"
                      step="0.01"
                      :placeholder="t('upload.settings.customExpirePlaceholder')"
                      class="mt-2 w-full px-3 py-1.5 text-xs border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      @input="handleCustomExpireInput"
                    />
                  </div>
                  
                  <!-- 按次数删除 -->
                  <div>
                    <div class="text-xs text-gray-700 dark:text-gray-300 mb-2 font-semibold">{{ t('upload.settings.maxDownloads') }}</div>
                    <div class="flex flex-wrap gap-1.5">
                      <button
                        v-for="option in publicConfig.getDownloadOptions().filter(opt => opt.value > 0)"
                        :key="'download-' + option.value"
                        @click="expireType = 'download'; maxDownloads = String(option.value); showCustomDownloads = false"
                        :class="[
                          'px-2.5 py-1 text-xs font-semibold rounded-full transition-all',
                          expireType === 'download' && maxDownloads === String(option.value) && !showCustomDownloads
                            ? 'bg-gradient-to-r from-purple-500 to-purple-600 text-white shadow-sm'
                            : 'bg-white dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-600 border border-gray-200 dark:border-gray-600'
                        ]"
                      >
                        {{ t(option.labelKey) }}
                      </button>
                      <button
                        @click="toggleCustomDownloads"
                        :class="[
                          'px-2.5 py-1 text-xs font-semibold rounded-full transition-all flex items-center gap-1',
                          showCustomDownloads
                            ? 'bg-gradient-to-r from-purple-500 to-purple-600 text-white shadow-sm'
                            : 'bg-white dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-600 border border-gray-200 dark:border-gray-600'
                        ]"
                      >
                        <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
                        </svg>
                        {{ t('upload.settings.custom') }}
                      </button>
                    </div>
                    <input
                      v-if="showCustomDownloads"
                      v-model="customMaxDownloads"
                      type="number"
                      min="1"
                      step="1"
                      :placeholder="t('upload.settings.customDownloadsPlaceholder')"
                      class="mt-2 w-full px-3 py-1.5 text-xs border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                      @input="handleCustomDownloadsInput"
                    />
                  </div>
                </div>

                <!-- 备注和按钮 -->
                <div class="space-y-3 mt-4">
                  <Textarea
                    id="remark"
                    v-model="remark"
                    :placeholder="t('upload.settings.remark', '备注（可选）')"
                    class="resize-none text-xs h-14 rounded-lg"
                    rows="2"
                  />

                  <div class="flex gap-2">
                    <Button
                      @click="activeTab === 'file' ? uploadFiles() : uploadText()"
                      :disabled="uploading || (activeTab === 'file' ? selectedFiles.length === 0 : !textContent.trim())"
                      class="flex-1 h-9 text-sm font-semibold" :class="activeTab === 'file' ? 'bg-gradient-to-r from-blue-500 to-blue-600 hover:from-blue-600 hover:to-blue-700' : 'bg-gradient-to-r from-purple-500 to-purple-600 hover:from-purple-600 hover:to-purple-700'"
                    >
                      <svg v-if="uploading" class="animate-spin -ml-1 mr-2 h-3.5 w-3.5" fill="none" viewBox="0 0 24 24">
                        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 0 1 8-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 0 1 4 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                      </svg>
                      {{ uploading ? `${uploadProgress}%` : (activeTab === 'file' ? t('upload.uploadFileButton', '寄存文件') : t('upload.uploadTextButton', '寄存文本')) }}
                    </Button>
                    
                    <Button
                      @click="router.push('/')"
                      variant="outline"
                      class="h-9 px-4 text-sm"
                    >
                      {{ t('common.back') }}
                    </Button>
                  </div>

                  <!-- 上传进度 -->
                  <div v-if="uploading" class="p-2 bg-gradient-to-br from-blue-50 to-purple-50 dark:from-blue-950/30 dark:to-purple-950/30 rounded-lg border transition-colors" :class="activeTab === 'file' ? 'border-blue-200 dark:border-blue-700' : 'border-purple-200 dark:border-purple-700'">
                    <Progress :value="uploadProgress" class="h-1.5" />
                  </div>
                </div>
              </div>
            </div>
          
          </CardContent>
        </Card>
        
        <!-- 成功对话框 -->
        <Dialog :open="showSuccessDialog" @update:open="handleSuccessDialogClose">
          <DialogContent class="sm:max-w-[500px]" @interact-outside="$event.preventDefault()" @escape-key-down="$event.preventDefault()">
            <DialogHeader>
              <div class="flex flex-col items-center mb-4">
                <div class="w-16 h-16 rounded-full bg-green-100 flex items-center justify-center mb-4">
                  <svg class="w-8 h-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
                  </svg>
                </div>
                <DialogTitle class="text-2xl text-gray-900 dark:text-gray-100 text-center">
                  {{ t('upload.uploadSuccess') }}
                </DialogTitle>
              </div>
              <DialogDescription class="text-center">
                {{ t('upload.successDescription', '您的内容已安全存入暂存柜，请妖善保管取件凭证') }}
              </DialogDescription>
            </DialogHeader>
            
            <div class="py-6">
              <div class="bg-gradient-to-br from-blue-50 to-purple-50 dark:from-blue-950/30 dark:to-purple-950/30 p-6 rounded-xl border-2 border-blue-200 dark:border-blue-800 shadow-inner transition-colors">
                <div class="text-center">
                  <p class="text-sm text-gray-600 dark:text-gray-300 mb-3 font-medium">{{ t('upload.pickupCodeGenerated') }}</p>
                  <p class="text-3xl font-mono font-bold text-transparent bg-clip-text bg-gradient-to-r from-blue-600 to-purple-600 dark:from-blue-400 dark:to-purple-400 tracking-[0.3em] mb-4">
                    {{ uploadResult?.pickup_code }}
                  </p>
                  <Button @click="copyPickupCode" class="w-full bg-gradient-to-r from-blue-500 to-purple-500 hover:from-blue-600 hover:to-purple-600 text-white shadow-lg">
                    <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"></path>
                    </svg>
                    {{ t('upload.copyPickupCode') }}
                  </Button>
                </div>
              </div>
              
              <div class="mt-6 space-y-2 bg-gray-50 dark:bg-gray-800 p-4 rounded-lg transition-colors">
                <div v-if="uploadResult?.expire_at" class="flex items-center text-sm text-gray-600 dark:text-gray-300">
                  <svg class="w-4 h-4 mr-2 text-gray-400 dark:text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                  </svg>
                  {{ t('upload.expiresAt', '暂存柜过期') }}: {{ formatDate(uploadResult.expire_at) }}
                </div>
                <div v-if="uploadResult?.max_downloads" class="flex items-center text-sm text-gray-600 dark:text-gray-300">
                  <svg class="w-4 h-4 mr-2 text-gray-400 dark:text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4"></path>
                  </svg>
                  {{ t('upload.maxDownloadsLimit', '最大取件次数') }}: {{ uploadResult.max_downloads }}
                </div>
              </div>
            </div>

            <DialogFooter class="sm:space-x-2 space-y-2 sm:space-y-0">
              <Button @click="continueUpload" variant="outline" class="w-full sm:w-auto">
                {{ t('upload.uploadAnother') }}
              </Button>
              <Button @click="router.push('/')" class="w-full sm:w-auto bg-gradient-to-r from-blue-500 to-purple-500 hover:from-blue-600 hover:to-purple-600 text-white">
                {{ t('nav.backToHome') }}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
        
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, h, onMounted, nextTick, watch } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import { useI18n } from '@/composables/useI18n'

// 组件导入
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { Card, CardContent } from '@/components/ui/card'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Progress } from '@/components/ui/progress'

// API 和工具导入
import { publicApi, utils, type UploadResponse } from '@/lib/api'
import { usePublicConfig } from '@/composables/usePublicConfig'

const router = useRouter()
const publicConfig = usePublicConfig()
const { config } = publicConfig
const { t } = useI18n()

// 响应式数据
const activeTab = ref('file')
const isDragging = ref(false)
const selectedFiles = ref<File[]>([])
const textContent = ref('')

// 共用配置
const expireType = ref('time')
const expireDays = ref('7')
const maxDownloads = ref('3')
const remark = ref('')

// 自定义输入相关
const showCustomExpire = ref(false)
const showCustomDownloads = ref(false)
const customExpireDays = ref('')
const customMaxDownloads = ref('')

// 上传状态
const uploading = ref(false)
const uploadProgress = ref(0)
const showSuccessDialog = ref(false)
const uploadResult = ref<UploadResponse | null>(null)

// DOM 引用
const fileInput = ref<HTMLInputElement | null>(null)
const tabGroup = ref<HTMLElement | null>(null)
const fileTabButton = ref<HTMLElement | null>(null)
const textTabButton = ref<HTMLElement | null>(null)

// 标签页指示器状态
const tabIndicatorWidth = ref('0px')
const tabIndicatorPosition = ref('0px')
const tabIndicatorInitialized = ref(false)

// 切换标签页
const switchTab = (tab: 'file' | 'text') => {
  activeTab.value = tab
}

// 更新标签页指示器位置和宽度
const updateTabIndicator = async () => {
  await nextTick()
  
  if (!tabGroup.value || !fileTabButton.value || !textTabButton.value) return
  
  const containerRect = tabGroup.value.getBoundingClientRect()
  let activeButton: HTMLElement | null = null
  
  if (activeTab.value === 'file') {
    activeButton = fileTabButton.value
  } else if (activeTab.value === 'text') {
    activeButton = textTabButton.value
  }
  
  if (activeButton) {
    const buttonRect = activeButton.getBoundingClientRect()
    tabIndicatorWidth.value = `${buttonRect.width}px`
    tabIndicatorPosition.value = `${buttonRect.left - containerRect.left}px`
  }
}

// 监听 activeTab 变化
watch(activeTab, async () => {
  await updateTabIndicator()
})

// 组件挂载后初始化
onMounted(async () => {
  // 加载公共配置
  await publicConfig.loadConfig()
  
  // 先设置位置，不显示动画
  await updateTabIndicator()
  // 使用 requestAnimationFrame 确保位置已经应用到 DOM
  requestAnimationFrame(() => {
    requestAnimationFrame(() => {
      // 启用动画效果
      tabIndicatorInitialized.value = true
    })
  })
  // 监听窗口大小变化，重新计算位置
  window.addEventListener('resize', updateTabIndicator)
})

// 计算属性
const totalSize = computed(() => {
  const total = selectedFiles.value.reduce((sum, file) => sum + file.size, 0)
  return utils.formatFileSize(total)
})

// 方法
const triggerFileInput = () => {
  fileInput.value?.click()
}

const handleFileSelect = (event: Event) => {
  const target = event.target as HTMLInputElement
  const files = target.files
  if (files) {
    addFiles(Array.from(files))
  }
}

const handleFileDrop = (event: DragEvent) => {
  isDragging.value = false
  const files = event.dataTransfer?.files
  if (files) {
    addFiles(Array.from(files))
  }
}

const addFiles = (files: File[]) => {
  // 过滤重复文件
  const newFiles = files.filter(file => {
    return !selectedFiles.value.some(existing => 
      existing.name === file.name && existing.size === file.size
    )
  })
  
  // 检查文件大小，过滤掉超过限制的文件
  const oversizedFiles = newFiles.filter(file => !publicConfig.validateFileSize(file))
  const validFiles = newFiles.filter(file => publicConfig.validateFileSize(file))
  
  if (oversizedFiles.length > 0) {
    const fileNames = oversizedFiles.map(f => f.name).join(', ')
    toast.error('文件超过大小限制', {
      duration: 4000,
      description: h('span', { class: 'text-gray-700 dark:text-gray-200' }, `以下文件超过 ${publicConfig.getFileSizeLimit()} 限制，已被移除：${fileNames}`)
    })
  }
  
  if (validFiles.length > 0) {
    selectedFiles.value.push(...validFiles)
    toast.success(`添加了 ${validFiles.length} 个文件`)
  }
  
  // 检查文件数量
  if (!publicConfig.validateFileCount(selectedFiles.value)) {
    toast.warning(`最多只能上传 ${config.value.upload?.max_batch_files || 10} 个文件`)
  }
}

const removeFile = (index: number) => {
  selectedFiles.value.splice(index, 1)
}

// 切换自定义过期时间输入
const toggleCustomExpire = () => {
  showCustomExpire.value = !showCustomExpire.value
  if (showCustomExpire.value) {
    expireType.value = 'time'
    showCustomDownloads.value = false
  }
}

// 切换自定义下载次数输入
const toggleCustomDownloads = () => {
  showCustomDownloads.value = !showCustomDownloads.value
  if (showCustomDownloads.value) {
    expireType.value = 'download'
    showCustomExpire.value = false
  }
}

// 处理自定义过期时间输入
const handleCustomExpireInput = () => {
  const value = parseFloat(customExpireDays.value)
  if (!isNaN(value) && value > 0) {
    const maxDays = config.value.upload.max_retention_days
    if (value > maxDays) {
      toast.warning(t('upload.errors.expireExceedsMax', { max: maxDays }))
      customExpireDays.value = String(maxDays)
      expireDays.value = String(maxDays)
    } else {
      expireDays.value = customExpireDays.value
    }
  }
}

// 处理自定义下载次数输入
const handleCustomDownloadsInput = () => {
  const value = parseInt(customMaxDownloads.value)
  if (!isNaN(value) && value > 0) {
    maxDownloads.value = customMaxDownloads.value
  }
}

const clearFiles = () => {
  selectedFiles.value = []
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

const uploadFiles = async () => {
  if (selectedFiles.value.length === 0) {
    toast.warning(t('upload.errors.noFiles'))
    return
  }

  // 再次校验文件大小
  const oversizedFiles = selectedFiles.value.filter(file => !publicConfig.validateFileSize(file))
  if (oversizedFiles.length > 0) {
    const fileNames = oversizedFiles.map(f => f.name).join(', ')
    toast.error(t('upload.errors.fileTooLarge'), {
      duration: 4000,
      description: h('span', { class: 'text-gray-700 dark:text-gray-200' }, `${t('upload.files', '以下文件')}:${fileNames}`)
    })
    return
  }

  // 再次校验文件数量
  if (!publicConfig.validateFileCount(selectedFiles.value)) {
    toast.error(t('upload.errors.tooManyFiles', { max: config.value.upload?.max_batch_files || 10 }))
    return
  }

  uploading.value = true
  uploadProgress.value = 0

  try {
    const formData = new FormData()
    
    // 添加文件
    selectedFiles.value.forEach(file => {
      formData.append('files', file)
    })

    // 添加配置
    const config = {
      expire_type: expireType.value,
      [expireType.value === 'time' ? 'expire_days' : 'max_downloads']: 
        expireType.value === 'time' ? parseInt(expireDays.value) : parseInt(maxDownloads.value),
      remark: remark.value.trim() || undefined
    }

    Object.entries(config).forEach(([key, value]) => {
      if (value !== undefined) {
        formData.append(key, String(value))
      }
    })

    const response = await publicApi.uploadFiles(formData, {
      onUploadProgress: (progressEvent: any) => {
        if (progressEvent.total) {
          uploadProgress.value = Math.round((progressEvent.loaded / progressEvent.total) * 100)
        }
      }
    })

    if (response.data.code === 200) {
      uploadResult.value = response.data.data
      showSuccessDialog.value = true
      toast.success(t('upload.uploadSuccess'))
    } else {
      throw new Error(response.data.msg || t('upload.uploadFailed'))
    }
  } catch (error: any) {
    handleUploadError(error, t('upload.uploadFailed'))
  } finally {
    uploading.value = false
    uploadProgress.value = 0
  }
}

const pasteText = async () => {
  try {
    const text = await navigator.clipboard.readText()
    if (text.trim()) {
      textContent.value = text
      toast.success(t('common.pasteSuccess'))
    } else {
      toast.warning(t('common.clipboardEmpty'))
    }
  } catch (error) {
    toast.error(t('common.clipboardFailed'))
  }
}

const uploadText = async () => {
  if (!textContent.value.trim()) {
    toast.warning(t('upload.errors.noText'))
    return
  }

  uploading.value = true
  uploadProgress.value = 0
  
  // 模拟文本上传进度
  const progressInterval = setInterval(() => {
    if (uploadProgress.value < 90) {
      uploadProgress.value += 10
    }
  }, 100)

  try {
    const config = {
      content: textContent.value.trim(),
      expire_type: expireType.value,
      [expireType.value === 'time' ? 'expire_days' : 'max_downloads']: 
        expireType.value === 'time' ? parseInt(expireDays.value) : parseInt(maxDownloads.value),
      remark: remark.value.trim() || undefined
    }

    const response = await publicApi.uploadText(config)

    if (response.data.code === 200) {
      uploadProgress.value = 100
      uploadResult.value = response.data.data
      showSuccessDialog.value = true
      toast.success('文本保存成功！')
    } else {
      throw new Error(response.data.msg || '保存失败')
    }
  } catch (error: any) {
    handleUploadError(error, '文本保存失败')
  } finally {
    clearInterval(progressInterval)
    uploading.value = false
    uploadProgress.value = 0
  }
}

const handleUploadError = (error: any, defaultMessage: string) => {
  // 不在控制台显示业务异常（如404等）
  if (error.response?.status >= 400 && error.response?.status < 500) {
    // 客户端错误，只显示toast，不打印到控制台
    toast.error(error.response?.data?.msg || defaultMessage)
  } else {
    // 服务器错误或网络错误，记录到控制台用于调试
    console.error(`${defaultMessage}:`, error)
    toast.error(error.response?.data?.msg || `${defaultMessage}，请重试`)
  }
}

const copyPickupCode = async () => {
  if (uploadResult.value?.pickup_code) {
    const success = await utils.copyToClipboard(uploadResult.value.pickup_code)
    if (success) {
      toast.success('取件码已复制到剪贴板')
    } else {
      toast.error('复制失败，请手动复制')
    }
  }
}

const continueUpload = () => {
  showSuccessDialog.value = false
  uploadResult.value = null
  
  // 清空表单内容
  if (activeTab.value === 'file') {
    clearFiles()
  } else {
    textContent.value = ''
  }
  
  // 重置共用配置为默认值
  expireType.value = 'time'
  expireDays.value = '7'
  maxDownloads.value = '3'
  remark.value = ''
}

// 处理成功对话框关闭事件（包括右上角X按钮）
const handleSuccessDialogClose = (open: boolean) => {
  if (!open) {
    // 关闭时执行与"继续上传"相同的逻辑
    continueUpload()
  }
}

const formatFileSize = (bytes: number): string => {
  return utils.formatFileSize(bytes)
}

const formatDate = (dateString: string): string => {
  return utils.formatDate(dateString)
}

// 图标组件映射
const getFileIconComponent = (filename: string) => {
  const icon = utils.getFileTypeIcon(filename)
  
  const iconMap: Record<string, any> = {
    image: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z' })
    ]),
    video: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z' })
    ]),
    audio: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M9 19V6l12-3v13M9 19c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zm12-3c0 1.105-1.343 2-3 2s-3-.895-3-2 1.343-2 3-2 3 .895 3 2zM9 10l12-3' })
    ]),
    archive: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M5 8l4 4 4-4m6-4v12a2 2 0 01-2 2H7a2 2 0 01-2-2V4a2 2 0 012-2h10a2 2 0 012 2z' })
    ]),
    file: () => h('svg', { fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', { 'stroke-linecap': 'round', 'stroke-linejoin': 'round', 'stroke-width': '2', d: 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z' })
    ])
  }

  return iconMap[icon] || iconMap.file
}
</script>

<style scoped>
/* 自定义样式已移至全局配置 */
</style>