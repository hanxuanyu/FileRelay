<template>
  <AdminLayout>
    <!-- 页面标题和搜索 -->
    <div class="mb-8">
      <div class="md:flex md:items-center md:justify-between">
        <div class="flex-1 min-w-0">
          <h2 class="text-3xl font-bold text-gray-900 dark:text-gray-100">{{ t('admin.batches.title') }}</h2>
          <p class="mt-2 text-gray-600 dark:text-gray-400">{{ t('admin.batches.description') }}</p>
        </div>
        <div class="mt-4 flex md:mt-0 md:ml-4">
          <Button @click="refreshData" variant="outline" size="sm" class="mr-2">
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
            </svg>
            {{ t('common.refresh') }}
          </Button>
          <Button @click="showCleanConfirmDialog = true" variant="outline" size="sm" class="mr-2">
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
            </svg>
            {{ t('admin.batches.manualClean') }}
          </Button>
        </div>
      </div>
    </div>

        <!-- 筛选和搜索 -->
        <Card class="mb-6">
          <CardContent class="pt-6">
            <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
              <div>
                <Label for="status-filter">{{ t('admin.batches.filters.statusLabel') }}</Label>
                <Select v-model="filters.status">
                  <SelectTrigger class="mt-1.5">
                    <SelectValue :placeholder="t('admin.batches.filters.statusAll')" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">{{ t('admin.batches.filters.statusAll') }}</SelectItem>
                    <SelectItem value="active">{{ t('admin.batches.filters.statusActive') }}</SelectItem>
                    <SelectItem value="expired">{{ t('admin.batches.filters.statusExpired') }}</SelectItem>
                    <SelectItem value="deleted">{{ t('admin.batches.filters.statusDeleted') }}</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <div class="md:col-span-2">
                <Label for="pickup-code">{{ t('admin.batches.filters.pickupCodeSearch') }}</Label>
                <Input
                  id="pickup-code"
                  v-model="filters.pickupCode"
                  autocomplete="off"
                  :placeholder="t('admin.batches.filters.pickupCodePlaceholder')"
                />
              </div>
              <div class="flex items-end">
                <Button @click="searchBatches" class="mr-2">
                  <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
                  </svg>
                  {{ t('common.search') }}
                </Button>
                <Button variant="outline" @click="clearFilters">
                  {{ t('common.clear') }}
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>

        <!-- 批次列表 -->
        <Card>
          <CardHeader>
            <div class="flex items-center justify-between">
              <div>
                <CardTitle>{{ t('admin.batches.list.title') }}</CardTitle>
                <CardDescription>
                  {{ t('admin.batches.list.description', { total: pagination.total, page: pagination.page }) }}
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
            <div v-else-if="batches.length === 0" class="text-center py-12">
              <svg class="w-16 h-16 mx-auto text-gray-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
              </svg>
              <p class="text-gray-500 text-lg">{{ t('admin.batches.list.empty') }}</p>
              <p class="text-gray-400 text-sm">{{ t('admin.batches.list.emptyHint') }}</p>
            </div>

            <!-- 批次表格 -->
            <div v-else class="overflow-x-auto">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>{{ t('admin.batches.list.columns.pickupCode') }}</TableHead>
                    <TableHead>{{ t('admin.batches.list.columns.type') }}</TableHead>
                    <TableHead>{{ t('admin.batches.list.columns.status') }}</TableHead>
                    <TableHead>{{ t('admin.batches.list.columns.downloadCount') }}</TableHead>
                    <TableHead>{{ t('admin.batches.list.columns.createdAt') }}</TableHead>
                    <TableHead>{{ t('admin.batches.list.columns.expireAt') }}</TableHead>
                    <TableHead class="text-center">{{ t('admin.batches.list.columns.operation') }}</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  <TableRow v-for="batch in batches" :key="batch.id" class="dark:hover:bg-black/20">
                    <TableCell class="font-medium font-mono">
                      {{ batch.pickup_code }}
                    </TableCell>
                    <TableCell>
                      <Badge variant="outline">
                        {{ batch.type === 'text' ? t('admin.batches.types.text') : t('admin.batches.types.file') }}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      <Badge :variant="getStatusVariant(batch.status)">
                        {{ getStatusText(batch.status) }}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      {{ batch.download_count }}{{ batch.max_downloads ? ` / ${batch.max_downloads}` : '' }}
                    </TableCell>
                    <TableCell>{{ formatDate(batch.created_at) }}</TableCell>
                    <TableCell>
                      <span v-if="batch.expire_at">{{ formatDate(batch.expire_at) }}</span>
                      <span v-else class="text-gray-500">{{ t('admin.batches.detail.neverExpire') }}</span>
                    </TableCell>
                    <TableCell>
                      <div class="flex items-center space-x-2">
                        <Button variant="outline" size="sm" @click="showBatchDetail(batch)">
                          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
                          </svg>
                        </Button>
                        <Button variant="outline" size="sm" @click="editBatch(batch)">
                          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
                          </svg>
                        </Button>
                        <Button variant="destructive" size="sm" @click="deleteBatch(batch)" :disabled="deleting.has(batch.id)">
                          <svg v-if="deleting.has(batch.id)" class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 0 1 8-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 0 1 4 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                          </svg>
                          <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
                          </svg>
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                </TableBody>
              </Table>
            </div>

            <!-- 分页 -->
            <div v-if="pagination.total > pagination.pageSize" class="flex items-center justify-between mt-6">
              <div class="text-sm text-gray-700">
                {{ t('admin.batches.pagination.showing', {
                  from: (pagination.page - 1) * pagination.pageSize + 1,
                  to: Math.min(pagination.page * pagination.pageSize, pagination.total),
                  total: pagination.total
                }) }}
              </div>
              <div class="flex items-center space-x-2">
                <Button
                  variant="outline"
                  size="sm"
                  @click="changePage(pagination.page - 1)"
                  :disabled="pagination.page <= 1"
                >
                  {{ t('admin.batches.pagination.previous') }}
                </Button>
                <Button
                  variant="outline"
                  size="sm"
                  @click="changePage(pagination.page + 1)"
                  :disabled="pagination.page >= Math.ceil(pagination.total / pagination.pageSize)"
                >
                  {{ t('admin.batches.pagination.next') }}
                </Button>
              </div>
            </div>
          </CardContent>
        </Card>

    <!-- 批次编辑对话框 -->
    <Dialog :open="showEditDialog" @update:open="showEditDialog = $event">
      <DialogContent class="sm:max-w-[500px] max-h-[90vh] flex flex-col p-0">
        <DialogHeader class="px-6 pt-6 pb-4 flex-shrink-0">
          <DialogTitle>{{ t('admin.batches.edit.title') }}</DialogTitle>
          <DialogDescription>
            {{ t('admin.batches.edit.description') }}
          </DialogDescription>
        </DialogHeader>
        
        <div v-if="editingBatch" class="space-y-4 px-6 overflow-y-auto flex-1">
          <div>
            <Label>{{ t('admin.batches.detail.pickupCode') }}</Label>
            <Input :model-value="editingBatch.pickup_code" disabled class="mt-1.5 bg-gray-50" />
            <p class="text-xs text-gray-500 mt-1">{{ t('admin.batches.edit.pickupCodeHint', { length: configLimits.pickupCodeLength }) }}</p>
          </div>

          <div>
            <Label for="edit-type">{{ t('admin.batches.edit.typeLabel') }}</Label>
            <Select v-model="editForm.type">
              <SelectTrigger id="edit-type" class="mt-1.5">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="file">{{ t('admin.batches.types.file') }}</SelectItem>
                <SelectItem value="text">{{ t('admin.batches.types.text') }}</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div v-if="editForm.type === 'text'">
            <Label for="edit-content">{{ t('admin.batches.detail.textContent') }}</Label>
            <Textarea
              id="edit-content"
              v-model="editForm.content"
              :placeholder="t('admin.batches.edit.contentPlaceholder')"
              class="mt-1.5 resize-none"
              rows="4"
              :maxlength="1000000"
            />
            <p class="text-xs text-gray-500 mt-1">{{ t('admin.batches.edit.contentLength', { count: editForm.content.length }) }}</p>
          </div>

          <div>
            <Label for="edit-remark">{{ t('admin.batches.edit.remarkLabel') }}</Label>
            <Textarea
              id="edit-remark"
              v-model="editForm.remark"
              :placeholder="t('admin.batches.edit.remarkPlaceholder')"
              class="mt-1.5 resize-none"
              rows="3"
            />
          </div>

          <div>
            <Label for="edit-expire-type">{{ t('admin.batches.edit.expireTypeLabel') }}</Label>
            <Select v-model="editForm.expire_type">
              <SelectTrigger id="edit-expire-type" class="mt-1.5">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="permanent">{{ t('admin.batches.edit.expireTypePermanent') }}</SelectItem>
                <SelectItem value="time">{{ t('admin.batches.edit.expireTypeTime') }}</SelectItem>
                <SelectItem value="download">{{ t('admin.batches.edit.expireTypeDownload') }}</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div v-if="editForm.expire_type === 'time'">
            <Label>{{ t('admin.batches.edit.expireAtLabel') }}</Label>
            <div class="flex gap-2 mt-1.5">
              <div class="flex-1">
                <Popover v-model:open="datePickerOpen">
                  <PopoverTrigger as-child>
                    <Button
                      variant="outline"
                      class="w-full justify-start font-normal"
                    >
                      <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
                      </svg>
                      {{ expireDateValue ? expireDateValue.toDate(getLocalTimeZone()).toLocaleDateString('zh-CN') : t('admin.batches.edit.selectDate') }}
                    </Button>
                  </PopoverTrigger>
                  <PopoverContent class="w-auto overflow-hidden p-0" align="start">
                    <Calendar
                      :model-value="expireDateValue"
                      layout="month-and-year"
                      @update:model-value="(value) => {
                        if (value) {
                          expireDateValue = value
                          datePickerOpen = false
                        }
                      }"
                      locale="zh-CN"
                    />
                  </PopoverContent>
                </Popover>
              </div>
              <div class="w-32">
                <Input
                  v-model="expireTime"
                  type="text"
                  placeholder="23:59"
                  maxlength="5"
                  class="bg-background"
                  @blur="validateTimeFormat"
                />
              </div>
            </div>
            <p class="text-xs text-gray-500 mt-1">{{ t('admin.batches.edit.expireTimeHint', { days: configLimits.maxRetentionDays }) }}</p>
          </div>

          <div v-if="editForm.expire_type === 'download'">
            <Label for="edit-max-downloads">{{ t('admin.batches.detail.maxDownloads') }}</Label>
            <Input
              id="edit-max-downloads"
              v-model.number="editForm.max_downloads"
              type="number"
              min="1"
              max="9999"
              :placeholder="t('admin.batches.edit.maxDownloadsPlaceholder')"
              class="mt-1.5"
            />
            <p class="text-xs text-gray-500 mt-1">{{ t('admin.batches.edit.maxDownloadsHint') }}</p>
          </div>

          <div>
            <Label for="edit-status">{{ t('admin.batches.edit.statusLabel') }}</Label>
            <Select v-model="editForm.status">
              <SelectTrigger id="edit-status" class="mt-1.5">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="active">{{ t('admin.batches.statuses.active') }}</SelectItem>
                <SelectItem value="expired">{{ t('admin.batches.statuses.expired') }}</SelectItem>
                <SelectItem value="deleted">{{ t('admin.batches.statuses.deleted') }}</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div>
            <Label for="edit-download-count">{{ t('admin.batches.detail.downloadCount') }}</Label>
            <div class="flex items-center gap-2 mt-1.5">
              <Input
                id="edit-download-count"
                v-model.number="editForm.download_count"
                type="number"
                disabled
                class="bg-gray-50"
              />
              <Button
                variant="outline"
                size="sm"
                @click="editForm.download_count = 0"
                :disabled="editForm.download_count === 0"
              >
                {{ t('common.reset') }}
              </Button>
            </div>
            <p class="text-xs text-gray-500 mt-1">{{ t('admin.batches.edit.downloadCountHint') }}</p>
          </div>
        </div>

        <DialogFooter class="px-6 pb-6 pt-4 flex-shrink-0 border-t">
          <Button variant="outline" @click="showEditDialog = false">
            {{ t('common.cancel') }}
          </Button>
          <Button @click="saveEdit" :disabled="saving">
            <svg v-if="saving" class="animate-spin -ml-1 mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 0 1 8-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 0 1 4 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            {{ saving ? t('admin.batches.edit.saving') : t('common.save') }}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- 删除确认对话框 -->
    <AlertDialog :open="showDeleteDialog" @update:open="showDeleteDialog = $event">
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>{{ t('admin.batches.delete.title') }}</AlertDialogTitle>
          <AlertDialogDescription>
            {{ t('admin.batches.delete.description', { code: deletingBatch?.pickup_code }) }}
            <br />
            {{ t('admin.batches.delete.warning') }}
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>{{ t('common.cancel') }}</AlertDialogCancel>
          <AlertDialogAction @click="confirmDelete" class="bg-red-600 hover:bg-red-700">
            {{ t('common.delete') }}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>

    <!-- 清理确认对话框 -->
    <AlertDialog :open="showCleanConfirmDialog" @update:open="showCleanConfirmDialog = $event">
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>{{ t('admin.batches.clean.title') }}</AlertDialogTitle>
          <AlertDialogDescription>
            {{ t('admin.batches.clean.description') }}
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>{{ t('common.cancel') }}</AlertDialogCancel>
          <AlertDialogAction @click="confirmClean" :disabled="cleaning" class="bg-orange-600 hover:bg-orange-700">
            <svg v-if="cleaning" class="animate-spin -ml-1 mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 0 1 8-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 0 1 4 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            {{ cleaning ? t('admin.batches.clean.cleaning') : t('admin.batches.clean.confirmButton') }}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>

    <!-- 批次详情对话框 -->
    <Dialog :open="showDetailDialog" @update:open="showDetailDialog = $event">
      <DialogContent class="sm:max-w-[600px] max-h-[90vh] flex flex-col p-0">
        <DialogHeader class="px-6 pt-6 pb-4 flex-shrink-0">
          <DialogTitle>{{ t('admin.batches.detail.title') }}</DialogTitle>
          <DialogDescription>
            {{ t('admin.batches.detail.description') }}
          </DialogDescription>
        </DialogHeader>
        
        <div v-if="selectedBatch" class="px-6 space-y-4 overflow-y-auto flex-1">
          <div class="grid grid-cols-2 gap-4 text-sm">
            <div class="space-y-2">
              <p><span class="font-medium">{{ t('admin.batches.detail.pickupCode') }}:</span> {{ selectedBatch.pickup_code }}</p>
              <p><span class="font-medium">{{ t('admin.batches.detail.type') }}:</span> {{ selectedBatch.type === 'text' ? t('admin.batches.types.text') : t('admin.batches.types.file') }}</p>
              <p><span class="font-medium">{{ t('admin.batches.detail.status') }}:</span> 
                <Badge :variant="getStatusVariant(selectedBatch.status)" class="ml-2">
                  {{ getStatusText(selectedBatch.status) }}
                </Badge>
              </p>
              <p><span class="font-medium">{{ t('admin.batches.detail.downloadCount') }}:</span> {{ selectedBatch.download_count }}{{ selectedBatch.max_downloads ? ` / ${selectedBatch.max_downloads}` : '' }}</p>
            </div>
            <div class="space-y-2">
              <p><span class="font-medium">{{ t('admin.batches.detail.createdAt') }}:</span> {{ formatDate(selectedBatch.created_at) }}</p>
              <p><span class="font-medium">{{ t('admin.batches.detail.updatedAt') }}:</span> {{ formatDate(selectedBatch.updated_at) }}</p>
              <p><span class="font-medium">{{ t('admin.batches.detail.expireAt') }}:</span> 
                <span v-if="selectedBatch.expire_at">{{ formatDate(selectedBatch.expire_at) }}</span>
                <span v-else class="text-gray-500">{{ t('admin.batches.detail.neverExpire') }}</span>
              </p>
              <p><span class="font-medium">{{ t('admin.batches.detail.expireType') }}:</span> {{ getExpireTypeText(selectedBatch.expire_type) }}</p>
            </div>
          </div>
          
          <div v-if="selectedBatch.remark">
            <p class="font-medium mb-1">{{ t('admin.batches.detail.remark') }}:</p>
            <p class="text-gray-600 bg-gray-50 p-2 rounded">{{ selectedBatch.remark }}</p>
          </div>

          <div v-if="selectedBatch.type === 'text' && selectedBatch.content">
            <p class="font-medium mb-1">{{ t('admin.batches.detail.textContent') }}:</p>
            <div class="bg-gray-50 dark:bg-gray-800 p-3 rounded border border-gray-200 dark:border-gray-700 max-h-40 overflow-y-auto transition-colors">
              <pre class="text-sm whitespace-pre-wrap">{{ selectedBatch.content }}</pre>
            </div>
          </div>

          <div v-if="selectedBatch.file_items && selectedBatch.file_items.length > 0">
            <p class="font-medium mb-2">{{ t('admin.batches.detail.fileList') }} ({{ selectedBatch.file_items.length }} {{ t('admin.batches.detail.filesCount') }}):</p>
            <div class="space-y-2 max-h-40 overflow-y-auto">
              <div v-for="file in selectedBatch.file_items" :key="file.id" class="flex items-center justify-between p-2 bg-gray-50 rounded">
                <div>
                  <p class="font-medium">{{ file.original_name }}</p>
                  <p class="text-xs text-gray-500">{{ formatFileSize(file.size) }} • {{ file.mime_type }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>

        <DialogFooter class="px-6 pb-6 pt-4 flex-shrink-0 border-t">
          <Button @click="showDetailDialog = false">
            {{ t('common.close') }}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- Toast 组件 -->
    <Toaster />
  </AdminLayout>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { toast } from 'vue-sonner'
import { useI18n } from '@/composables/useI18n'
import { usePublicConfig } from '@/composables/usePublicConfig'
import type { DateValue } from '@internationalized/date'
import { getLocalTimeZone, fromDate } from '@internationalized/date'

// 组件导入
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Textarea } from '@/components/ui/textarea'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Skeleton } from '@/components/ui/skeleton'
import { Calendar } from '@/components/ui/calendar'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from '@/components/ui/alert-dialog'
import { Toaster } from '@/components/ui/sonner'
import AdminLayout from '@/layouts/AdminLayout.vue'

// API 和工具导入
import { adminApi, utils, type FileBatch } from '@/lib/api'

const { t } = useI18n()
const { config: publicConfig } = usePublicConfig()

// 响应式数据
const loading = ref(true)
const batches = ref<FileBatch[]>([])
const selectedBatch = ref<FileBatch | null>(null)
const showDetailDialog = ref(false)
const deleting = ref(new Set<string>())

// 编辑相关
const showEditDialog = ref(false)
const editingBatch = ref<FileBatch | null>(null)
const saving = ref(false)
const editForm = reactive({
  type: 'file' as 'file' | 'text',
  content: '',
  remark: '',
  expire_type: 'permanent' as 'time' | 'download' | 'permanent',
  expire_at: '',
  max_downloads: 0,
  download_count: 0,
  status: 'active' as 'active' | 'expired' | 'deleted'
})

// 日期时间选择器状态
const expireDateValue = ref<DateValue | undefined>()
const expireTime = ref('23:59')
const datePickerOpen = ref(false)

// 配置约束计算属性
const configLimits = computed(() => ({
  pickupCodeLength: publicConfig.value?.security?.pickup_code_length || 6,
  maxFileSize: publicConfig.value?.upload?.max_file_size_mb || 100,
  maxBatchFiles: publicConfig.value?.upload?.max_batch_files || 10,
  maxRetentionDays: publicConfig.value?.upload?.max_retention_days || 30
}))

// 删除确认相关
const showDeleteDialog = ref(false)
const deletingBatch = ref<FileBatch | null>(null)

// 清理确认相关
const showCleanConfirmDialog = ref(false)
const cleaning = ref(false)

const filters = reactive({
  status: 'all',
  pickupCode: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 方法
const fetchBatches = async () => {
  loading.value = true
  
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.pageSize
    }

    if (filters.status && filters.status !== 'all') {
      params.status = filters.status
    }

    if (filters.pickupCode.trim()) {
      params.pickup_code = filters.pickupCode.trim()
    }

    const response = await adminApi.getBatches(params)
    
    if (response.data.code === 200) {
      batches.value = response.data.data.data
      pagination.total = response.data.data.total
      pagination.page = response.data.data.page
      pagination.pageSize = response.data.data.page_size
    } else {
      throw new Error(response.data.msg || t('admin.batches.messages.fetchFailed'))
    }
  } catch (error: any) {
    console.error('获取批次列表失败:', error)
    toast.error(error.response?.data?.msg || t('admin.batches.messages.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const searchBatches = () => {
  pagination.page = 1
  fetchBatches()
}

const clearFilters = () => {
  filters.status = 'all'
  filters.pickupCode = ''
  pagination.page = 1
  fetchBatches()
}

const refreshData = () => {
  fetchBatches()
  toast.success(t('admin.batches.messages.refreshSuccess'))
}

const changePage = (page: number) => {
  pagination.page = page
  fetchBatches()
}

const showBatchDetail = async (batch: FileBatch) => {
  try {
    const response = await adminApi.getBatchDetail(batch.id)
    if (response.data.code === 200) {
      selectedBatch.value = response.data.data
      showDetailDialog.value = true
    }
  } catch (error: any) {
    toast.error(t('admin.batches.messages.fetchDetailFailed'))
  }
}

const editBatch = async (batch: FileBatch) => {
  try {
    // 获取批次详细信息以便编辑
    const response = await adminApi.getBatchDetail(batch.id)
    if (response.data.code === 200) {
      editingBatch.value = response.data.data
      
      // 填充编辑表单
      editForm.type = editingBatch.value.type
      editForm.content = editingBatch.value.content || ''
      editForm.remark = editingBatch.value.remark || ''
      editForm.expire_type = editingBatch.value.expire_type
      editForm.max_downloads = editingBatch.value.max_downloads || 0
      editForm.download_count = editingBatch.value.download_count || 0
      editForm.status = editingBatch.value.status
      
      // 处理过期时间格式
      if (editingBatch.value.expire_at) {
        const date = new Date(editingBatch.value.expire_at)
        // 设置日期选择器的值
        expireDateValue.value = fromDate(date, getLocalTimeZone())
        // 设置时间
        const hours = String(date.getHours()).padStart(2, '0')
        const minutes = String(date.getMinutes()).padStart(2, '0')
        expireTime.value = `${hours}:${minutes}`
        // 同时设置editForm.expire_at以兼容
        const year = date.getFullYear()
        const month = String(date.getMonth() + 1).padStart(2, '0')
        const day = String(date.getDate()).padStart(2, '0')
        editForm.expire_at = `${year}-${month}-${day}T${hours}:${minutes}`
      } else {
        expireDateValue.value = undefined
        expireTime.value = '23:59'
        editForm.expire_at = ''
      }
      
      showEditDialog.value = true
    }
  } catch (error: any) {
    toast.error(t('admin.batches.messages.fetchFailed'))
  }
}

const saveEdit = async () => {
  if (!editingBatch.value) return

  // 数据验证
  if (editForm.type === 'text' && !editForm.content.trim()) {
    toast.error(t('admin.batches.edit.errors.contentRequired'))
    return
  }

  if (editForm.expire_type === 'time' && editForm.expire_at) {
    const expireDate = new Date(editForm.expire_at)
    const maxDate = new Date()
    maxDate.setDate(maxDate.getDate() + configLimits.value.maxRetentionDays)
    
    if (expireDate > maxDate) {
      toast.error(t('admin.batches.edit.errors.expireTimeExceeds', { days: configLimits.value.maxRetentionDays }))
      return
    }

    if (expireDate < new Date()) {
      toast.error(t('admin.batches.edit.errors.expireTimePast'))
      return
    }
  }

  if (editForm.expire_type === 'download') {
    if (!editForm.max_downloads || editForm.max_downloads < 1) {
      toast.error(t('admin.batches.edit.errors.maxDownloadsMin'))
      return
    }
    if (editForm.max_downloads > 9999) {
      toast.error(t('admin.batches.edit.errors.maxDownloadsMax'))
      return
    }
  }

  if (editForm.download_count < 0) {
    toast.error(t('admin.batches.edit.errors.downloadCountNegative'))
    return
  }

  saving.value = true

  try {
    const updateData: any = {
      type: editForm.type,
      content: editForm.type === 'text' ? editForm.content : undefined,
      remark: editForm.remark || null,
      expire_type: editForm.expire_type,
      max_downloads: editForm.max_downloads || 0,
      download_count: editForm.download_count,
      status: editForm.status
    }

    // 处理过期时间
    if (editForm.expire_type === 'time' && expireDateValue.value) {
      // 组合日期和时间
      const dateObj = expireDateValue.value.toDate(getLocalTimeZone())
      const [hours = '0', minutes = '0'] = expireTime.value.split(':')
      dateObj.setHours(parseInt(hours), parseInt(minutes), 0, 0)
      updateData.expire_at = dateObj.toISOString()
    } else {
      updateData.expire_at = null
    }

    const response = await adminApi.updateBatch(editingBatch.value.id, updateData)
    
    if (response.data.code === 200) {
      toast.success(t('admin.batches.edit.saveSuccess'))
      showEditDialog.value = false
      fetchBatches()
    } else {
      throw new Error(response.data.msg || t('admin.batches.edit.saveFailed'))
    }
  } catch (error: any) {
    console.error('更新批次失败:', error)
    toast.error(error.response?.data?.msg || t('admin.batches.edit.saveFailed'))
  } finally {
    saving.value = false
  }
}

const deleteBatch = (batch: FileBatch) => {
  deletingBatch.value = batch
  showDeleteDialog.value = true
}

const confirmDelete = async () => {
  if (!deletingBatch.value) return

  const batchId = deletingBatch.value.id
  deleting.value.add(batchId)

  try {
    const response = await adminApi.deleteBatch(batchId)
    
    if (response.data.code === 200) {
      toast.success(t('admin.batches.delete.deleteSuccess'))
      showDeleteDialog.value = false
      deletingBatch.value = null
      fetchBatches()
    } else {
      throw new Error(response.data.msg || t('admin.batches.delete.deleteFailed'))
    }
  } catch (error: any) {
    console.error('删除批次失败:', error)
    toast.error(error.response?.data?.msg || t('admin.batches.delete.deleteFailed'))
  } finally {
    deleting.value.delete(batchId)
  }
}

const confirmClean = async () => {
  cleaning.value = true

  try {
    const response = await adminApi.cleanBatches()
    
    if (response.data.code === 200) {
      toast.success(t('admin.batches.clean.cleanSuccess'))
      showCleanConfirmDialog.value = false
      // 刷新列表以反映清理后的状态
      fetchBatches()
    } else {
      throw new Error(response.data.msg || t('admin.batches.clean.cleanFailed'))
    }
  } catch (error: any) {
    console.error('清理失败:', error)
    toast.error(error.response?.data?.msg || t('admin.batches.clean.cleanFailed'))
  } finally {
    cleaning.value = false
  }
}

const getStatusVariant = (status: string) => {
  switch (status) {
    case 'active':
      return 'default'
    case 'expired':
      return 'secondary'
    case 'deleted':
      return 'destructive'
    default:
      return 'outline'
  }
}

const getStatusText = (status: string) => {
  switch (status) {
    case 'active':
      return t('admin.batches.statuses.active')
    case 'expired':
      return t('admin.batches.statuses.expired')
    case 'deleted':
      return t('admin.batches.statuses.deleted')
    default:
      return t('admin.batches.statuses.unknown')
  }
}

const getExpireTypeText = (expireType: string) => {
  switch (expireType) {
    case 'time':
      return t('admin.batches.detail.expireTypeTime')
    case 'download':
      return t('admin.batches.detail.expireTypeDownload')
    case 'permanent':
      return t('admin.batches.detail.expireTypePermanent')
    default:
      return t('admin.batches.detail.expireTypeUnknown')
  }
}

const formatDate = (dateString: string): string => {
  return utils.formatDate(dateString)
}

const formatFileSize = (bytes: number): string => {
  return utils.formatFileSize(bytes)
}

// 验证和格式化时间输入
const validateTimeFormat = () => {
  const timePattern = /^([0-1]?[0-9]|2[0-3]):([0-5][0-9])$/
  
  if (!expireTime.value) {
    expireTime.value = '23:59'
    return
  }
  
  // 尝试自动补全格式
  let time = expireTime.value.replace(/[^\d:]/g, '')
  
  // 如果只输入了数字，尝试智能解析
  if (!time.includes(':')) {
    if (time.length === 3) {
      // 如 "959" -> "9:59"
      time = time[0] + ':' + time.slice(1)
    } else if (time.length === 4) {
      // 如 "2359" -> "23:59"
      time = time.slice(0, 2) + ':' + time.slice(2)
    }
  }
  
  // 验证格式
  if (timePattern.test(time)) {
    // 格式化为标准的 HH:MM 格式
    const [h = '0', m = '0'] = time.split(':')
    expireTime.value = `${h.padStart(2, '0')}:${m.padStart(2, '0')}`
  } else {
    toast.error(t('admin.batches.edit.errors.timeFormat'))
    expireTime.value = '23:59'
  }
}

// 组件挂载
onMounted(() => {
  fetchBatches()
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

/* 内容区域动画 */
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

/* 按钮悬停效果 */
.hover\:scale-105:hover {
  transform: scale(1.05);
  transition: transform 0.2s ease;
}
</style>