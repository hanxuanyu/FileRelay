import { ref } from 'vue'
import { publicApi, type PublicConfig } from '@/lib/api'

// 默认配置
const defaultConfig: PublicConfig = {
  site: {
    name: '文件中转站',
    description: '安全、便捷的文件暂存服务',
    logo: '/favicon.png'
  },
  upload: {
    max_file_size_mb: 100,
    max_batch_files: 10,
    max_retention_days: 30,
    require_token: false
  },
  security: {
    pickup_code_length: 6
  },
  storage: {
    type: 'local'
  },
  api_token: {
    enabled: false
  }
}

// 全局配置状态
const config = ref<PublicConfig>(defaultConfig)
const loading = ref(false)
const loaded = ref(false) // 添加已加载标记

export function usePublicConfig() {
  // 加载公共配置
  const loadConfig = async (force = false) => {
    // 如果已经加载过且不是强制刷新，直接返回
    if (loaded.value && !force) {
      return
    }
    
    try {
      loading.value = true
      const response = await publicApi.getPublicConfig()
      config.value = response.data.data
      loaded.value = true // 标记为已加载
      // 更新页面标题和图标
      updateSiteInfo()
    } catch (error) {
      console.error('加载公共配置失败:', error)
      // 使用默认配置
      config.value = defaultConfig
    } finally {
      loading.value = false
    }
  }

  // 更新页面标题和图标
  const updateSiteInfo = () => {
    // 更新标题
    if (config.value.site.name) {
      document.title = config.value.site.name
    }

    // 更新 favicon
    if (config.value.site.logo) {
      const link = document.querySelector("link[rel~='icon']") as HTMLLinkElement
      if (link) {
        link.href = config.value.site.logo
      } else {
        const newLink = document.createElement('link')
        newLink.rel = 'icon'
        newLink.type = 'image/png'
        newLink.href = config.value.site.logo
        document.head.appendChild(newLink)
      }
    }
  }

  // 获取上传配置
  const getUploadLimits = () => {
    return {
      maxFileSizeMB: config.value.upload.max_file_size_mb,
      maxFileSize: config.value.upload.max_file_size_mb * 1024 * 1024, // 字节
      maxBatchFiles: config.value.upload.max_batch_files,
      maxRetentionDays: config.value.upload.max_retention_days
    }
  }

  // 生成失效时间选项
  const getExpireOptions = () => {
    const maxDays = config.value.upload.max_retention_days
    const options = [
      { labelKey: 'upload.settings.expireTimeOptions.1h', value: 1/24 },
      { labelKey: 'upload.settings.expireTimeOptions.6h', value: 6/24 },
      { labelKey: 'upload.settings.expireTimeOptions.24h', value: 1 },
      { labelKey: 'upload.settings.expireTimeOptions.72h', value: 3 },
      { labelKey: 'upload.settings.expireTimeOptions.168h', value: 7 }
    ]

    // 只显示不超过最大保存天数的选项
    const validOptions = options.filter(option => option.value <= maxDays)
    
    // 如果最大天数大于7，添加更多选项
    if (maxDays >= 14) {
      validOptions.push({ labelKey: 'upload.settings.expireTimeOptions.14d', value: 14 })
    }
    if (maxDays >= 30) {
      validOptions.push({ labelKey: 'upload.settings.expireTimeOptions.30d', value: 30 })
    }
    if (maxDays >= 60) {
      validOptions.push({ labelKey: 'upload.settings.expireTimeOptions.60d', value: 60 })
    }
    if (maxDays >= 90) {
      validOptions.push({ labelKey: 'upload.settings.expireTimeOptions.90d', value: 90 })
    }

    // 如果允许，添加永久选项
    validOptions.push({ labelKey: 'upload.settings.expireTimeOptions.never', value: 0 })

    return validOptions
  }

  // 生成下载次数选项
  const getDownloadOptions = () => {
    return [
      { labelKey: 'upload.settings.maxDownloadsOptions.1', value: 1 },
      { labelKey: 'upload.settings.maxDownloadsOptions.3', value: 3 },
      { labelKey: 'upload.settings.maxDownloadsOptions.5', value: 5 },
      { labelKey: 'upload.settings.maxDownloadsOptions.10', value: 10 },
      { labelKey: 'upload.settings.maxDownloadsOptions.20', value: 20 },
      { labelKey: 'upload.settings.maxDownloadsOptions.50', value: 50 },
      { labelKey: 'upload.settings.maxDownloadsOptions.unlimited', value: 0 }
    ]
  }

  // 验证文件大小
  const validateFileSize = (file: File): boolean => {
    const limits = getUploadLimits()
    return file.size <= limits.maxFileSize
  }

  // 验证文件数量
  const validateFileCount = (files: FileList | File[]): boolean => {
    const limits = getUploadLimits()
    return files.length <= limits.maxBatchFiles
  }

  // 格式化文件大小限制文本
  const getFileSizeLimit = (): string => {
    const sizeMB = config.value.upload.max_file_size_mb
    if (sizeMB >= 1024) {
      return `${(sizeMB / 1024).toFixed(1)}GB`
    }
    return `${sizeMB}MB`
  }

  return {
    config,
    loading,
    loadConfig,
    getUploadLimits,
    getExpireOptions,
    getDownloadOptions,
    validateFileSize,
    validateFileCount,
    getFileSizeLimit
  }
}

// 初始化配置（应用启动时调用）
export async function initPublicConfig() {
  const { loadConfig } = usePublicConfig()
  await loadConfig()
}