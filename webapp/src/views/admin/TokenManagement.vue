<template>
  <AdminLayout>
    <!-- 页面标题和操作 -->
    <div class="mb-8">
      <div class="md:flex md:items-center md:justify-between">
        <div class="flex-1 min-w-0">
          <h2 class="text-3xl font-bold text-gray-900 dark:text-gray-100">{{ t('admin.tokens.title') }}</h2>
          <p class="mt-2 text-gray-600 dark:text-gray-400">{{ t('admin.tokens.description') }}</p>
        </div>
        <div class="mt-4 flex md:mt-0 md:ml-4">
              <Button @click="refreshData" variant="outline" size="sm" class="mr-2">
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
                </svg>
                {{ t('admin.tokens.refresh') }}
              </Button>
              <Button @click="showCreateDialog = true" size="sm">
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
                </svg>
                {{ t('admin.tokens.create') }}
              </Button>
            </div>
          </div>
        </div>

        <!-- API Token 说明 -->
        <Card class="mb-6 border-blue-200 dark:border-blue-800 bg-blue-50 dark:bg-blue-950/30 transition-colors">
          <CardContent class="pt-6">
            <div class="flex">
              <div class="flex-shrink-0">
                <svg class="h-5 w-5 text-blue-400 dark:text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                </svg>
              </div>
              <div class="ml-3">
                <h3 class="text-sm font-medium text-blue-800 dark:text-blue-300">{{ t('admin.tokens.usage.title') }}</h3>
                <div class="mt-2 text-sm text-blue-700 dark:text-blue-400">
                  <ul class="list-disc list-inside space-y-1">
                    <li>{{ t('admin.tokens.usage.point1') }}</li>
                    <li>{{ t('admin.tokens.usage.point2') }}</li>
                    <li>{{ t('admin.tokens.usage.point3') }}</li>
                    <li>{{ t('admin.tokens.usage.point4') }}</li>
                  </ul>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>

        <!-- Token 列表 -->
        <Card>
          <CardHeader>
            <div class="flex items-center justify-between">
              <div>
                <CardTitle>{{ t('admin.tokens.list.title') }}</CardTitle>
                <CardDescription>
                  {{ t('admin.tokens.list.description') }}
                </CardDescription>
              </div>
            </div>
          </CardHeader>
          <CardContent>
            <!-- Loading 状态 -->
            <div v-if="loading" class="space-y-4">
              <Skeleton class="h-4 w-full" />
              <Skeleton class="h-4 w-3/4" />
              <Skeleton class="h-4 w-1/2" />
            </div>

            <!-- 空状态 -->
            <div v-else-if="tokens.length === 0" class="text-center py-12">
              <svg class="w-16 h-16 mx-auto text-gray-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"></path>
              </svg>
              <p class="text-gray-500 text-lg">{{ t('admin.tokens.list.empty') }}</p>
              <p class="text-gray-400 text-sm">{{ t('admin.tokens.list.emptyHint') }}</p>
            </div>

            <!-- Token 表格 -->
            <div v-else class="overflow-x-auto">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>{{ t('admin.tokens.list.columns.name') }}</TableHead>
                    <TableHead>{{ t('admin.tokens.list.columns.scope') }}</TableHead>
                    <TableHead>{{ t('admin.tokens.list.columns.status') }}</TableHead>
                    <TableHead>{{ t('admin.tokens.list.columns.createdAt') }}</TableHead>
                    <TableHead>{{ t('admin.tokens.list.columns.lastUsedAt') }}</TableHead>
                    <TableHead>{{ t('admin.tokens.list.columns.expireAt') }}</TableHead>
                    <TableHead class="text-center">{{ t('admin.tokens.list.columns.operation') }}</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  <TableRow v-for="token in tokens" :key="token.id" class="hover:bg-gray-50 dark:hover:bg-gray-800 transition-colors">
                    <TableCell class="font-medium">
                      {{ token.name }}
                    </TableCell>
                    <TableCell>
                      <Badge variant="outline">
                        {{ token.scope || 'ALL' }}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      <Badge :variant="token.revoked ? 'destructive' : 'default'">
                        {{ token.revoked ? t('admin.tokens.statuses.revoked') : t('admin.tokens.statuses.active') }}
                      </Badge>
                    </TableCell>
                    <TableCell>{{ formatDate(token.created_at) }}</TableCell>
                    <TableCell>
                      <span v-if="token.last_used_at">{{ formatDate(token.last_used_at) }}</span>
                      <span v-else class="text-gray-500">{{ t('admin.tokens.list.neverUsed') }}</span>
                    </TableCell>
                    <TableCell>
                      <span v-if="token.expire_at">{{ formatDate(token.expire_at) }}</span>
                      <span v-else class="text-gray-500">{{ t('admin.tokens.list.neverExpire') }}</span>
                    </TableCell>
                    <TableCell>
                      <div class="flex items-center space-x-2">
                        <Button 
                          v-if="!token.revoked"
                          variant="outline" 
                          size="sm" 
                          @click="openRevokeDialog(token)"
                          :disabled="revoking.has(token.id)"
                        >
                          <svg v-if="revoking.has(token.id)" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 0 1 8-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 0 1 4 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                          </svg>
                          <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728L5.636 5.636m12.728 12.728L18.364 5.636M5.636 18.364l12.728-12.728"></path>
                          </svg>
                          {{ t('admin.tokens.actions.revoke') }}
                        </Button>
                        <Button 
                          v-else
                          variant="outline" 
                          size="sm" 
                          @click="openRecoverDialog(token)"
                          :disabled="recovering.has(token.id)"
                          class="text-green-600 hover:text-green-700 dark:text-green-500 dark:hover:text-green-400 border-green-600 hover:border-green-700 dark:border-green-500 dark:hover:border-green-400"
                        >
                          <svg v-if="recovering.has(token.id)" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 0 1 8-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 0 1 4 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                          </svg>
                          <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
                          </svg>
                          {{ t('admin.tokens.actions.recover') }}
                        </Button>
                        <Button 
                          variant="destructive" 
                          size="sm" 
                          @click="openDeleteDialog(token)"
                          :disabled="deleting.has(token.id)"
                        >
                          <svg v-if="deleting.has(token.id)" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 0 1 8-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 0 1 4 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                          </svg>
                          <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
                          </svg>
                          {{ t('admin.tokens.actions.delete') }}
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                </TableBody>
              </Table>
            </div>
          </CardContent>
        </Card>

    <!-- 创建 Token 对话框 -->
    <Dialog :open="showCreateDialog" @update:open="showCreateDialog = $event">
      <DialogContent class="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>{{ t('admin.tokens.createDialog.title') }}</DialogTitle>
          <DialogDescription>
            {{ t('admin.tokens.createDialog.description') }}
          </DialogDescription>
        </DialogHeader>
        
        <form @submit.prevent="createToken" class="py-4 space-y-4">
          <div>
            <Label for="token-name">{{ t('admin.tokens.createDialog.nameLabel') }} *</Label>
            <Input
              id="token-name"
              v-model="createForm.name"
              autocomplete="off"
              :placeholder="t('admin.tokens.createDialog.namePlaceholder')"
              required
            />
          </div>

          <div>
            <Label for="token-scope" class="mb-2 block">{{ t('admin.tokens.createDialog.scopeLabel') }} *</Label>
            <Select v-model="createForm.scopes" multiple>
              <SelectTrigger class="w-full">
                <SelectValue :placeholder="t('admin.tokens.createDialog.scopePlaceholder')">
                  <span v-if="createForm.scopes.length === 0" class="text-gray-500">
                    {{ t('admin.tokens.createDialog.scopePlaceholder') }}
                  </span>
                  <span v-else class="flex gap-1 flex-wrap">
                    <Badge v-for="scope in createForm.scopes" :key="scope" variant="secondary" class="text-xs">
                      {{ t(`admin.tokens.scopes.${scope}`) }}
                    </Badge>
                  </span>
                </SelectValue>
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  <SelectLabel>{{ t('admin.tokens.createDialog.scopeOptions') }}</SelectLabel>
                  <SelectItem value="upload">
                    <div class="flex flex-col">
                      <span class="font-medium">{{ t('admin.tokens.createDialog.scopeUpload') }}</span>
                      <span class="text-xs text-gray-500">{{ t('admin.tokens.createDialog.scopeUploadDesc') }}</span>
                    </div>
                  </SelectItem>
                  <SelectItem value="pickup">
                    <div class="flex flex-col">
                      <span class="font-medium">{{ t('admin.tokens.createDialog.scopePickup') }}</span>
                      <span class="text-xs text-gray-500">{{ t('admin.tokens.createDialog.scopePickupDesc') }}</span>
                    </div>
                  </SelectItem>
                  <SelectItem value="admin">
                    <div class="flex flex-col">
                      <span class="font-medium">{{ t('admin.tokens.createDialog.scopeAdmin') }}</span>
                      <span class="text-xs text-gray-500">{{ t('admin.tokens.createDialog.scopeAdminDesc') }}</span>
                    </div>
                  </SelectItem>
                </SelectGroup>
              </SelectContent>
            </Select>
            <p class="text-xs text-gray-500 mt-2">
              {{ t('admin.tokens.createDialog.scopeHint', { count: createForm.scopes.length }) }}
            </p>
          </div>

          <div>
            <Label for="token-expire">{{ t('admin.tokens.createDialog.expireLabel') }}</Label>
            <Select v-model="createForm.expireType">
              <SelectTrigger class="mt-1.5">
                <SelectValue :placeholder="t('admin.tokens.createDialog.expirePlaceholder')" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="never">{{ t('admin.tokens.createDialog.expireNever') }}</SelectItem>
                <SelectItem value="7">{{ t('admin.tokens.createDialog.expire7Days') }}</SelectItem>
                <SelectItem value="30">{{ t('admin.tokens.createDialog.expire30Days') }}</SelectItem>
                <SelectItem value="90">{{ t('admin.tokens.createDialog.expire90Days') }}</SelectItem>
                <SelectItem value="365">{{ t('admin.tokens.createDialog.expire1Year') }}</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </form>

        <DialogFooter>
          <Button variant="outline" @click="showCreateDialog = false">
            {{ t('admin.tokens.createDialog.cancelButton') }}
          </Button>
          <Button @click="createToken" :disabled="creating || !createForm.name.trim() || createForm.scopes.length === 0">
            <svg v-if="creating" class="animate-spin -ml-1 mr-3 h-4 w-4" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 0 1 8-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 0 1 4 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            {{ creating ? t('admin.tokens.createDialog.creating') : t('admin.tokens.createDialog.createButton') }}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- Token 创建成功对话框 -->
    <Dialog :open="showTokenDialog" @update:open="showTokenDialog = $event">
      <DialogContent class="sm:max-w-[600px]">
        <DialogHeader>
          <DialogTitle class="text-green-600">{{ t('admin.tokens.successDialog.title') }}</DialogTitle>
          <DialogDescription>
            {{ t('admin.tokens.successDialog.description') }}
          </DialogDescription>
        </DialogHeader>
        
        <div v-if="newToken" class="py-4">
          <div class="bg-gray-50 p-4 rounded-lg border border-gray-200">
            <div class="flex items-center justify-between mb-3">
              <div>
                <p class="text-sm text-gray-600">{{ t('admin.tokens.successDialog.tokenName') }}</p>
                <p class="font-medium">{{ newToken.data.name }}</p>
              </div>
              <Badge variant="outline">{{ newToken.data.scope || 'ALL' }}</Badge>
            </div>
            
            <div class="space-y-2">
              <p class="text-sm text-gray-600">{{ t('admin.tokens.successDialog.tokenLabel') }}</p>
              <div class="flex items-center space-x-2">
                <code class="flex-1 p-2 bg-white border rounded text-sm font-mono break-all">{{ newToken.token }}</code>
                <Button @click="copyToken" variant="outline" size="sm">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"></path>
                  </svg>
                </Button>
              </div>
            </div>

            <div v-if="newToken.data.expire_at" class="mt-3 text-sm text-gray-600">
              {{ t('admin.tokens.createDialog.expireLabel') }}: {{ formatDate(newToken.data.expire_at) }}
            </div>
          </div>

          <div class="mt-4 p-4 bg-yellow-50 border border-yellow-200 rounded-lg">
            <div class="flex">
              <div class="flex-shrink-0">
                <svg class="h-5 w-5 text-yellow-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L4.082 16.5c-.77.833.192 2.5 1.732 2.5z"></path>
                </svg>
              </div>
              <div class="ml-3">
                <p class="text-sm font-medium text-yellow-800">{{ t('admin.tokens.successDialog.warning') }}</p>
                <p class="text-sm text-yellow-700 mt-1">
                  {{ t('admin.tokens.successDialog.warningDesc') }}
                </p>
              </div>
            </div>
          </div>
        </div>

        <DialogFooter>
          <Button @click="closeTokenDialog">
            {{ t('admin.tokens.successDialog.closeButton') }}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- 撤销确认对话框 -->
    <AlertDialog :open="showRevokeDialog" @update:open="showRevokeDialog = $event">
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>{{ t('admin.tokens.revoke.title') }}</AlertDialogTitle>
          <AlertDialogDescription>
            {{ t('admin.tokens.revoke.description', { name: pendingToken?.name }) }}
            <br />
            <span class="text-orange-600 dark:text-orange-400 font-medium">{{ t('admin.tokens.revoke.warning') }}</span>
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>{{ t('admin.tokens.revoke.cancelButton') }}</AlertDialogCancel>
          <AlertDialogAction @click="confirmRevoke" class="bg-orange-600 hover:bg-orange-700">
            {{ t('admin.tokens.revoke.confirmButton') }}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>

    <!-- 恢复确认对话框 -->
    <AlertDialog :open="showRecoverDialog" @update:open="showRecoverDialog = $event">
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>{{ t('admin.tokens.recover.title') }}</AlertDialogTitle>
          <AlertDialogDescription>
            {{ t('admin.tokens.recover.description', { name: pendingToken?.name }) }}
            <br />
            <span class="text-green-600 dark:text-green-400 font-medium">{{ t('admin.tokens.recover.warning') }}</span>
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>{{ t('admin.tokens.recover.cancelButton') }}</AlertDialogCancel>
          <AlertDialogAction @click="confirmRecover" class="bg-green-600 hover:bg-green-700">
            {{ t('admin.tokens.recover.confirmButton') }}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>

    <!-- 删除确认对话框 -->
    <AlertDialog :open="showDeleteDialog" @update:open="showDeleteDialog = $event">
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>{{ t('admin.tokens.delete.title') }}</AlertDialogTitle>
          <AlertDialogDescription>
            {{ t('admin.tokens.delete.description', { name: pendingToken?.name }) }}
            <br />
            <span class="text-red-600 dark:text-red-400 font-medium">{{ t('admin.tokens.delete.warning') }}</span>
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>{{ t('admin.tokens.delete.cancelButton') }}</AlertDialogCancel>
          <AlertDialogAction @click="confirmDelete" class="bg-red-600 hover:bg-red-700">
            {{ t('admin.tokens.delete.confirmButton') }}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>

    <!-- Toast 组件 -->
    <Toaster />
  </AdminLayout>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { toast } from 'vue-sonner'
import { useI18n } from '@/composables/useI18n'

// 组件导入
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Select, SelectContent, SelectGroup, SelectItem, SelectLabel, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Skeleton } from '@/components/ui/skeleton'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from '@/components/ui/alert-dialog'
import { Toaster } from '@/components/ui/sonner'
import AdminLayout from '@/layouts/AdminLayout.vue'

// API 和工具导入
import { adminApi, utils, type APIToken } from '@/lib/api'

// i18n
const { t } = useI18n()

// 响应式数据
const loading = ref(true)
const tokens = ref<APIToken[]>([])
const showCreateDialog = ref(false)
const showTokenDialog = ref(false)
const showRevokeDialog = ref(false)
const showRecoverDialog = ref(false)
const showDeleteDialog = ref(false)
const pendingToken = ref<APIToken | null>(null)
const creating = ref(false)
const revoking = ref(new Set<number>())
const recovering = ref(new Set<number>())
const deleting = ref(new Set<number>())

const createForm = reactive({
  name: '',
  scopes: [] as string[],
  expireType: 'never'
})

const newToken = ref<{ token: string; data: APIToken } | null>(null)

// 方法
const fetchTokens = async () => {
  loading.value = true
  
  try {
    const response = await adminApi.getTokens()
    
    if (response.data.code === 200) {
      tokens.value = response.data.data
    } else {
      throw new Error(response.data.msg || t('admin.tokens.messages.fetchFailed'))
    }
  } catch (error: any) {
    console.error('获取 Token 列表失败:', error)
    toast.error(error.response?.data?.msg || t('admin.tokens.messages.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const createToken = async () => {
  if (!createForm.name.trim()) {
    toast.warning(t('admin.tokens.createDialog.nameRequired'))
    return
  }

  creating.value = true

  try {
    const data: any = {
      name: createForm.name.trim()
    }

    // 将选中的权限数组转换为逗号分隔的字符串
    if (createForm.scopes.length > 0) {
      data.scope = createForm.scopes.join(',')
    }

    if (createForm.expireType && createForm.expireType !== 'never') {
      const days = parseInt(createForm.expireType)
      const expireAt = new Date()
      expireAt.setDate(expireAt.getDate() + days)
      data.expire_at = expireAt.toISOString()
    }

    const response = await adminApi.createToken(data)

    if (response.data.code === 200 || response.data.code === 201) {
      newToken.value = response.data.data
      showCreateDialog.value = false
      showTokenDialog.value = true
      
      // 重置表单
      createForm.name = ''
      createForm.scopes = []
      createForm.expireType = 'never'
      
      // 刷新列表
      fetchTokens()
      
      toast.success(t('admin.tokens.messages.createSuccess'))
    } else {
      throw new Error(response.data.msg || t('admin.tokens.messages.createFailed'))
    }
  } catch (error: any) {
    console.error('创建 Token 失败:', error)
    toast.error(error.response?.data?.msg || t('admin.tokens.messages.createFailed'))
  } finally {
    creating.value = false
  }
}

const openRevokeDialog = (token: APIToken) => {
  pendingToken.value = token
  showRevokeDialog.value = true
}

const confirmRevoke = async () => {
  if (!pendingToken.value) return
  
  const token = pendingToken.value
  showRevokeDialog.value = false
  revoking.value.add(token.id)

  try {
    const response = await adminApi.revokeToken(token.id)
    
    if (response.data.code === 200) {
      toast.success(t('admin.tokens.messages.revokeSuccess'))
      fetchTokens()
    } else {
      throw new Error(response.data.msg || t('admin.tokens.messages.revokeFailed'))
    }
  } catch (error: any) {
    console.error('撤销 Token 失败:', error)
    toast.error(error.response?.data?.msg || t('admin.tokens.messages.revokeFailed'))
  } finally {
    revoking.value.delete(token.id)
    pendingToken.value = null
  }
}

const openRecoverDialog = (token: APIToken) => {
  pendingToken.value = token
  showRecoverDialog.value = true
}

const confirmRecover = async () => {
  if (!pendingToken.value) return
  
  const token = pendingToken.value
  showRecoverDialog.value = false
  recovering.value.add(token.id)

  try {
    const response = await adminApi.recoverToken(token.id)
    
    if (response.data.code === 200) {
      toast.success(t('admin.tokens.messages.recoverSuccess'))
      fetchTokens()
    } else {
      throw new Error(response.data.msg || t('admin.tokens.messages.recoverFailed'))
    }
  } catch (error: any) {
    console.error('恢复 Token 失败:', error)
    toast.error(error.response?.data?.msg || t('admin.tokens.messages.recoverFailed'))
  } finally {
    recovering.value.delete(token.id)
    pendingToken.value = null
  }
}

const openDeleteDialog = (token: APIToken) => {
  pendingToken.value = token
  showDeleteDialog.value = true
}

const confirmDelete = async () => {
  if (!pendingToken.value) return
  
  const token = pendingToken.value
  showDeleteDialog.value = false
  deleting.value.add(token.id)

  try {
    const response = await adminApi.deleteToken(token.id)
    
    if (response.data.code === 200) {
      toast.success(t('admin.tokens.messages.deleteSuccess'))
      fetchTokens()
    } else {
      throw new Error(response.data.msg || t('admin.tokens.messages.deleteFailed'))
    }
  } catch (error: any) {
    console.error('删除 Token 失败:', error)
    toast.error(error.response?.data?.msg || t('admin.tokens.messages.deleteFailed'))
  } finally {
    deleting.value.delete(token.id)
    pendingToken.value = null
  }
}

const copyToken = async () => {
  if (newToken.value?.token) {
    const success = await utils.copyToClipboard(newToken.value.token)
    if (success) {
      toast.success(t('admin.tokens.messages.copySuccess'))
    } else {
      toast.error(t('admin.tokens.messages.copyFailed'))
    }
  }
}

const closeTokenDialog = () => {
  showTokenDialog.value = false
  newToken.value = null
}

const refreshData = () => {
  fetchTokens()
  toast.success(t('admin.tokens.messages.refreshSuccess'))
}

const formatDate = (dateString: string): string => {
  return utils.formatDate(dateString)
}

// 组件挂载
onMounted(() => {
  fetchTokens()
})
</script>

<style scoped>
/* 导航样式 */
.border-b-2 {
  border-bottom-width: 2px;
}

.hover\:border-gray-300:hover {
  border-color: #d1d5db;
}

.hover\:text-gray-700:hover {
  color: #374151;
}

/* 表格行悬停效果 */
.hover\:bg-gray-50:hover {
  background-color: #f9fafb;
  transition: background-color 0.2s ease;
}

/* 代码块样式 */
code {
  word-break: break-all;
  white-space: pre-wrap;
}

/* 动画效果 */
.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

/* 内容动画 */
.space-y-4 > * {
  animation: fadeInUp 0.3s ease forwards;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 警告框样式 */
.bg-yellow-50 {
  animation: slideInDown 0.3s ease-out;
}

@keyframes slideInDown {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 成功提示样式 */
.bg-blue-50 {
  animation: fadeIn 0.3s ease-in-out;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}
</style>