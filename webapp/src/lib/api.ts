import axios from 'axios'
import type { AxiosResponse } from 'axios'

/**
 * API 基础配置
 * 
 * 标准实践：
 * - baseURL 包含 API 前缀（如 /api）
 * - 实际请求时使用相对路径（如 /batches）
 * - axios 会自动拼接成完整路径（/api/batches）
 * 
 * 环境配置：
 * - 开发环境：http://localhost:8080/api
 * - 生产环境：/api（相对路径，自动使用当前域名）
 */
const API_BASE_URL = import.meta.env.VITE_API_URL || '/api'

// 创建 axios 实例
const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
})

// 导出 API_BASE_URL 供其他模块使用（如 SSE 连接）
export { API_BASE_URL }

// 请求拦截器 - 添加认证头
api.interceptors.request.use((config) => {
  // 优先使用管理员 token（登录后自动保存）
  const adminToken = localStorage.getItem('admin_token')
  if (adminToken) {
    config.headers.Authorization = `Bearer ${adminToken}`
  } else {
    // 如果没有管理员 token，检查 API token
    const apiToken = localStorage.getItem('api_token')
    if (apiToken) {
      config.headers.Authorization = `Bearer ${apiToken}`
    }
  }
  return config
})

// 响应拦截器 - 统一错误处理
api.interceptors.response.use(
  (response) => response,
  (error) => {
    // 处理 401 未授权错误
    if (error.response?.status === 401 && window.location.pathname.includes('/admin')) {
      localStorage.removeItem('admin_token')
      window.location.href = '/admin/login'
    }
    
    // 抑制 4xx 业务异常的控制台输出
    // 4xx 错误属于预期的业务异常，不需要在控制台显示
    if (error.response?.status >= 400 && error.response?.status < 500) {
      // 创建一个新的错误对象，但不让它显示在控制台
      // 通过设置 error.config.suppressError 标记来标识这是业务异常
      error.config = error.config || {}
      error.config.suppressError = true
    }
    
    return Promise.reject(error)
  }
)

// API 响应类型定义
interface ApiResponse<T = any> {
  code: number
  msg: string
  data: T
}

interface FileItem {
  id: string
  batch_id: string
  original_name: string
  size: number
  mime_type: string
  storage_path: string
  download_url?: string
  created_at: string
}

interface FileBatch {
  id: string
  pickup_code: string
  type: 'file' | 'text'
  content?: string
  remark: string
  expire_type: 'time' | 'download' | 'permanent'
  expire_at?: string
  max_downloads?: number
  download_count: number
  status: 'active' | 'expired' | 'deleted'
  file_items: FileItem[]
  created_at: string
  updated_at: string
}

interface UploadResponse {
  pickup_code: string
  batch_id: string
  expire_at?: string
  max_downloads?: number
}

interface PickupResponse {
  type: 'file' | 'text'
  content?: string
  remark: string
  expire_type: string
  expire_at?: string
  max_downloads?: number
  download_count: number
  files: FileItem[]
}

interface DownloadCountResponse {
  download_count: number
  max_downloads?: number
}

interface BatchListResponse {
  data: FileBatch[]
  page: number
  page_size: number
  total: number
}

interface APIToken {
  id: number
  name: string
  scope: string
  created_at: string
  expire_at?: string
  last_used_at?: string
  revoked: boolean
}

// 配置相关类型定义
interface SiteConfig {
  name: string
  description: string
  logo?: string
  base_url?: string
  port?: number
}

interface UploadConfig {
  max_file_size_mb: number
  max_batch_files: number
  max_retention_days: number
  require_token: boolean
}

interface StorageConfig {
  type: string
  local?: {
    path: string
  }
  s3?: {
    endpoint: string
    region: string
    access_key: string
    secret_key: string
    bucket: string
    use_ssl: boolean
  }
  webdav?: {
    url: string
    username: string
    password: string
    root: string
  }
}

interface SecurityConfig {
  pickup_code_length: number
  pickup_fail_limit: number
  admin_password_hash: string
  jwt_secret: string
}

interface DatabaseConfig {
  type: string  // sqlite, mysql, postgres
  path?: string  // SQLite only
  host?: string  // MySQL/PostgreSQL
  port?: number  // MySQL/PostgreSQL
  user?: string  // MySQL/PostgreSQL
  password?: string  // MySQL/PostgreSQL
  dbname?: string  // MySQL/PostgreSQL
  config?: string  // Extra DSN config
}

interface APITokenConfig {
  enabled: boolean
  max_tokens: number
  allow_admin_api: boolean
}

interface WebConfig {
  path: string
}

interface LogConfig {
  level: string
  file_path: string
}

interface SystemConfig {
  site: SiteConfig
  upload: UploadConfig
  storage: StorageConfig
  security: SecurityConfig
  database: DatabaseConfig
  api_token: APITokenConfig
  web: WebConfig
  log: LogConfig
}

interface PublicConfig {
  site: SiteConfig
  upload: UploadConfig
  security: {
    pickup_code_length: number
  }
  storage: {
    type: string
  }
  api_token: {
    enabled: boolean
  }
}

// 公共 API
export const publicApi = {
  // 获取公共配置
  getPublicConfig: (): Promise<AxiosResponse<ApiResponse<PublicConfig>>> => {
    return api.get('/config')
  },

  // 上传文件
  uploadFiles: (formData: FormData, config?: any): Promise<AxiosResponse<ApiResponse<UploadResponse>>> => {
    return api.post('/batches', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      ...config,
    })
  },

  // 上传文本
  uploadText: (data: {
    content: string
    remark?: string
    expire_type?: string
    expire_days?: number
    max_downloads?: number
  }): Promise<AxiosResponse<ApiResponse<UploadResponse>>> => {
    return api.post('/batches/text', data)
  },

  // 获取批次信息
  getBatch: (pickupCode: string): Promise<AxiosResponse<ApiResponse<PickupResponse>>> => {
    return api.get(`/batches/${pickupCode}`)
  },

  // 查询下载次数
  getDownloadCount: (pickupCode: string): Promise<AxiosResponse<ApiResponse<DownloadCountResponse>>> => {
    return api.get(`/batches/${pickupCode}/count`)
  },

  // 下载单个文件
  downloadFile: (fileId: string): Promise<AxiosResponse<Blob>> => {
    return api.get(`/files/${fileId}/download`, {
      responseType: 'blob',
    })
  },

  // 下载单个文件（带文件名）
  downloadFileWithName: (fileId: string, filename: string): Promise<AxiosResponse<Blob>> => {
    return api.get(`/files/${fileId}/${encodeURIComponent(filename)}`, {
      responseType: 'blob',
    })
  },

  // 批量下载文件
  downloadBatch: (pickupCode: string): Promise<AxiosResponse<Blob>> => {
    return api.get(`/batches/${pickupCode}/download`, {
      responseType: 'blob',
    })
  },
}

// 管理 API
export const adminApi = {
  // 管理员登录
  login: (password: string): Promise<AxiosResponse<ApiResponse<{ token: string }>>> => {
    return api.post('/admin/login', { password })
  },

  // 获取批次列表
  getBatches: (params?: {
    page?: number
    page_size?: number
    status?: string
    pickup_code?: string
  }): Promise<AxiosResponse<ApiResponse<BatchListResponse>>> => {
    return api.get('/admin/batches', { params })
  },

  // 获取批次详情
  getBatchDetail: (batchId: string): Promise<AxiosResponse<ApiResponse<FileBatch>>> => {
    return api.get(`/admin/batches/${batchId}`)
  },

  // 更新批次信息
  updateBatch: (batchId: string, data: {
    remark?: string
    expire_type?: string
    expire_at?: string
    max_downloads?: number
    status?: string
  }): Promise<AxiosResponse<ApiResponse<FileBatch>>> => {
    return api.put(`/admin/batches/${batchId}`, data)
  },

  // 删除批次
  deleteBatch: (batchId: string): Promise<AxiosResponse<ApiResponse<void>>> => {
    return api.delete(`/admin/batches/${batchId}`)
  },

  // 手动触发清理
  cleanBatches: (): Promise<AxiosResponse<ApiResponse<void>>> => {
    return api.post('/admin/batches/clean')
  },

  // 获取 API Token 列表
  getTokens: (): Promise<AxiosResponse<ApiResponse<APIToken[]>>> => {
    return api.get('/admin/api-tokens')
  },

  // 创建 API Token
  createToken: (data: {
    name: string
    scope?: string
    expire_at?: string
  }): Promise<AxiosResponse<ApiResponse<{ token: string; data: APIToken }>>> => {
    return api.post('/admin/api-tokens', data)
  },

  // 撤销 API Token
  revokeToken: (tokenId: number): Promise<AxiosResponse<ApiResponse<void>>> => {
    return api.post(`/admin/api-tokens/${tokenId}/revoke`)
  },

  // 恢复 API Token
  recoverToken: (tokenId: number): Promise<AxiosResponse<ApiResponse<void>>> => {
    return api.post(`/admin/api-tokens/${tokenId}/recover`)
  },

  // 删除 API Token
  deleteToken: (tokenId: number): Promise<AxiosResponse<ApiResponse<void>>> => {
    return api.delete(`/admin/api-tokens/${tokenId}`)
  },

  // 获取系统配置
  getConfig: (): Promise<AxiosResponse<ApiResponse<SystemConfig>>> => {
    return api.get('/admin/config')
  },

  // 更新系统配置
  updateConfig: (config: SystemConfig): Promise<AxiosResponse<ApiResponse<void>>> => {
    return api.put('/admin/config', config)
  },
}

// 工具函数
export const utils = {
  // 格式化文件大小
  formatFileSize: (bytes: number): string => {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
  },

  // 格式化日期
  formatDate: (dateString: string): string => {
    const date = new Date(dateString)
    return date.toLocaleString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
    })
  },

  // 复制到剪贴板
  copyToClipboard: async (text: string): Promise<boolean> => {
    try {
      await navigator.clipboard.writeText(text)
      return true
    } catch (err) {
      // 兼容旧浏览器
      try {
        const textArea = document.createElement('textarea')
        textArea.value = text
        document.body.appendChild(textArea)
        textArea.select()
        document.execCommand('copy')
        document.body.removeChild(textArea)
        return true
      } catch (fallbackErr) {
        console.error('复制失败:', fallbackErr)
        return false
      }
    }
  },

  // 下载文件
  downloadBlob: (blob: Blob, filename: string) => {
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = filename
    document.body.appendChild(a)
    a.click()
    window.URL.revokeObjectURL(url)
    document.body.removeChild(a)
  },

  // 获取文件图标类型
  getFileTypeIcon: (filename: string): string => {
    const ext = filename.split('.').pop()?.toLowerCase()
    
    if (['jpg', 'jpeg', 'png', 'gif', 'bmp', 'webp'].includes(ext || '')) {
      return 'image'
    } else if (['mp4', 'avi', 'mkv', 'mov', 'wmv', 'flv'].includes(ext || '')) {
      return 'video'
    } else if (['mp3', 'wav', 'flac', 'aac', 'ogg'].includes(ext || '')) {
      return 'audio'
    } else if (['pdf'].includes(ext || '')) {
      return 'pdf'
    } else if (['doc', 'docx'].includes(ext || '')) {
      return 'word'
    } else if (['xls', 'xlsx'].includes(ext || '')) {
      return 'excel'
    } else if (['ppt', 'pptx'].includes(ext || '')) {
      return 'powerpoint'
    } else if (['zip', 'rar', '7z', 'tar', 'gz'].includes(ext || '')) {
      return 'archive'
    } else if (['txt', 'md', 'json', 'xml', 'csv'].includes(ext || '')) {
      return 'text'
    } else if (['js', 'ts', 'html', 'css', 'php', 'py', 'java', 'cpp', 'c'].includes(ext || '')) {
      return 'code'
    }
    
    return 'file'
  }
}

export default api

// 类型定义导出
export type {
  ApiResponse,
  FileItem,
  FileBatch,
  UploadResponse,
  PickupResponse,
  BatchListResponse,
  APIToken,
  SystemConfig,
  PublicConfig,
  SiteConfig,
  UploadConfig,
  StorageConfig,
  SecurityConfig,
  DatabaseConfig,
  APITokenConfig,
}