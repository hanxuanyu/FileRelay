<template>
  <AdminLayout>
    <div class="space-y-8">
      <!-- 页面标题和操作按钮 -->
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-3xl font-bold tracking-tight">{{ t('admin.config.title') }}</h1>
          <p class="text-muted-foreground mt-2">{{ t('admin.config.description') }}</p>
        </div>
        <div class="flex space-x-3">
          <Button variant="outline" @click="resetForm" :disabled="loading" size="sm" class="min-w-20">
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            {{ t('admin.config.resetButton') }}
          </Button>
          <Button @click="saveConfig" :disabled="loading" size="sm" class="min-w-24">
            <span v-if="loading" class="flex items-center">
              <svg class="animate-spin -ml-1 mr-3 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              {{ t('admin.config.saving') }}
            </span>
            <span v-else class="flex items-center">
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              {{ t('admin.config.saveButton') }}
            </span>
          </Button>
        </div>
      </div>

      <!-- 配置面板 -->
      <form @submit.prevent="saveConfig">
      <div class="grid grid-cols-1 xl:grid-cols-2 gap-6">
        <!-- 站点配置 -->
        <Card class="shadow-sm border-l-4 border-l-blue-500">
          <CardHeader>
            <CardTitle class="flex items-center">
              <svg class="w-5 h-5 mr-2 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9v-9m0-9v9" />
              </svg>
              {{ t('admin.config.site.title') }}
            </CardTitle>
            <CardDescription>{{ t('admin.config.site.description') }}</CardDescription>
          </CardHeader>
          <CardContent class="space-y-5">
            <div>
              <Label for="siteName" class="text-sm font-medium">{{ t('admin.config.site.nameLabel') }}</Label>
              <Input 
                id="siteName"
                v-model="config.site.name" 
                autocomplete="off"
                :placeholder="t('admin.config.site.namePlaceholder')"
                class="mt-2"
              />
            </div>
            <div>
              <Label for="siteDescription" class="text-sm font-medium">{{ t('admin.config.site.descriptionLabel') }}</Label>
              <Textarea 
                id="siteDescription"
                v-model="config.site.description" 
                :placeholder="t('admin.config.site.descriptionPlaceholder')"
                class="mt-2 resize-none"
                rows="3"
              />
            </div>
            <div>
              <Label for="siteLogo" class="text-sm font-medium">{{ t('admin.config.site.logoLabel') }}</Label>
              <Input 
                id="siteLogo"
                v-model="config.site.logo" 
                autocomplete="off"
                :placeholder="t('admin.config.site.logoPlaceholder')"
                class="mt-2"
              />
              <p class="text-xs text-muted-foreground mt-1">{{ t('admin.config.site.logoHint') }}</p>
            </div>
            <div>
              <Label for="siteBaseUrl" class="text-sm font-medium">{{ t('admin.config.site.baseUrlLabel') }}</Label>
              <Input 
                id="siteBaseUrl"
                v-model="config.site.base_url" 
                autocomplete="off"
                :placeholder="t('admin.config.site.baseUrlPlaceholder')"
                class="mt-2"
              />
              <p class="text-xs text-muted-foreground mt-1">{{ t('admin.config.site.baseUrlHint') }}</p>
            </div>
            <div>
              <Label for="sitePort" class="text-sm font-medium">{{ t('admin.config.site.portLabel') }}</Label>
              <Input 
                id="sitePort"
                v-model.number="config.site.port" 
                type="number"
                min="1"
                max="65535"
                autocomplete="off"
                :placeholder="t('admin.config.site.portPlaceholder')"
                class="mt-2"
              />
              <p class="text-xs text-muted-foreground mt-1">{{ t('admin.config.site.portHint') }}</p>
            </div>
          </CardContent>
        </Card>

        <!-- 上传配置 -->
        <Card class="shadow-sm border-l-4 border-l-green-500">
          <CardHeader>
            <CardTitle class="flex items-center">
              <svg class="w-5 h-5 mr-2 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
              </svg>
              {{ t('admin.config.upload.title') }}
            </CardTitle>
            <CardDescription>{{ t('admin.config.upload.description') }}</CardDescription>
          </CardHeader>
          <CardContent class="space-y-5">
            <div>
              <Label for="maxFileSizeMB" class="text-sm font-medium">{{ t('admin.config.upload.maxSizeLabel') }}</Label>
              <Input 
                id="maxFileSize"
                v-model.number="config.upload.max_file_size_mb" 
                type="number"
                min="1"
                autocomplete="off"
                :placeholder="t('admin.config.upload.maxSizePlaceholder')"
                class="mt-2"
              />
              <p class="text-xs text-muted-foreground mt-1">{{ t('admin.config.upload.maxSizeHint') }}</p>
            </div>
            <div>
              <Label for="maxBatchFiles" class="text-sm font-medium">{{ t('admin.config.upload.maxFilesLabel') }}</Label>
              <Input 
                id="maxBatchFiles"
                v-model.number="config.upload.max_batch_files" 
                type="number"
                min="1"
                autocomplete="off"
                :placeholder="t('admin.config.upload.maxFilesPlaceholder')"
                class="mt-2"
              />
              <p class="text-xs text-muted-foreground mt-1">{{ t('admin.config.upload.maxFilesHint') }}</p>
            </div>
            <div>
              <Label for="maxRetentionDays" class="text-sm font-medium">{{ t('admin.config.storage.maxExpireLabel') }}</Label>
              <Input 
                id="maxRetentionDays"
                v-model.number="config.upload.max_retention_days" 
                type="number"
                min="1"
                autocomplete="off"
                :placeholder="t('admin.config.storage.maxExpirePlaceholder')"
                class="mt-2"
              />
              <p class="text-xs text-muted-foreground mt-1">{{ t('admin.config.storage.maxExpireHint') }}</p>
            </div>
            <div class="flex items-center justify-between p-4 rounded-lg bg-muted/50">
              <div>
                <Label for="requireToken" class="text-base font-medium">{{ t('admin.config.upload.enableUploadLabel') }}</Label>
                <p class="text-sm text-muted-foreground mt-1">{{ t('admin.config.upload.enableUploadDesc') }}</p>
              </div>
              <Switch 
                id="requireToken"
                v-model="config.upload.require_token"
              />
            </div>
          </CardContent>
        </Card>

        <!-- 安全配置 -->
        <Card class="shadow-sm border-l-4 border-l-red-500">
          <CardHeader>
            <CardTitle class="flex items-center">
              <svg class="w-5 h-5 mr-2 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
              </svg>
              {{ t('admin.config.security.title') }}
            </CardTitle>
            <CardDescription>{{ t('admin.config.security.description') }}</CardDescription>
          </CardHeader>
          <CardContent class="space-y-5">
            <div>
              <Label for="pickupCodeLength" class="text-sm font-medium">{{ t('admin.config.security.pickupCodeLengthLabel') }}</Label>
              <Input 
                id="pickupCodeLength"
                v-model.number="config.security.pickup_code_length" 
                type="number"
                min="4"
                max="20"
                autocomplete="off"
                :placeholder="t('admin.config.security.pickupCodeLengthPlaceholder')"
                class="mt-2"
              />
              <p class="text-xs text-muted-foreground mt-1">{{ t('admin.config.security.pickupCodeLengthHint') }}</p>
            </div>
            <div>
              <Label for="pickupFailLimit" class="text-sm font-medium">{{ t('admin.config.security.pickupFailLimitLabel') }}</Label>
              <Input 
                id="pickupFailLimit"
                v-model.number="config.security.pickup_fail_limit" 
                type="number"
                min="0"
                autocomplete="off"
                :placeholder="t('admin.config.security.pickupFailLimitPlaceholder')"
                class="mt-2"
              />
              <p class="text-xs text-muted-foreground mt-1">{{ t('admin.config.security.pickupFailLimitHint') }}</p>
            </div>
            <div>
              <Label for="adminPassword" class="text-sm font-medium">{{ t('admin.config.security.adminPasswordLabel') }}</Label>
              <Input 
                id="adminPassword"
                v-model="adminPassword" 
                type="password"
                autocomplete="new-password"
                :placeholder="t('admin.config.security.adminPasswordPlaceholder')"
                class="mt-2"
              />
              <p class="text-xs text-muted-foreground mt-1">
                {{ t('admin.config.security.adminPasswordHint') }}
              </p>
            </div>
            <div>
              <Label for="jwtSecret" class="text-sm font-medium">{{ t('admin.config.security.jwtSecretLabel') }}</Label>
              <Input 
                id="jwtSecret"
                v-model="config.security.jwt_secret" 
                type="password"
                autocomplete="off"
                :placeholder="t('admin.config.security.jwtSecretPlaceholder')"
                class="mt-2"
              />
              <p class="text-xs text-muted-foreground mt-1">
                {{ t('admin.config.security.jwtSecretHint') }}
              </p>
            </div>
          </CardContent>
        </Card>

        <!-- API Token 配置 -->
        <Card class="shadow-sm border-l-4 border-l-purple-500">
          <CardHeader>
            <CardTitle class="flex items-center">
              <svg class="w-5 h-5 mr-2 text-purple-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
              </svg>
              {{ t('admin.config.apiToken.title') }}
            </CardTitle>
            <CardDescription>{{ t('admin.config.apiToken.description') }}</CardDescription>
          </CardHeader>
          <CardContent class="space-y-6">
            <div class="flex items-center justify-between p-4 rounded-lg bg-muted/50">
              <div>
                <Label for="apiTokenEnabled" class="text-base font-medium">{{ t('admin.config.apiToken.enabledLabel') }}</Label>
                <p class="text-sm text-muted-foreground mt-1">{{ t('admin.config.apiToken.enabledDesc') }}</p>
              </div>
              <Switch 
                id="apiTokenEnabled"
                v-model="config.api_token.enabled"
              />
            </div>
            
            <div class="space-y-5" :class="{ 'opacity-50': !config.api_token.enabled }">
              <div>
                <Label for="maxTokens" class="text-sm font-medium">{{ t('admin.config.apiToken.maxTokensLabel') }}</Label>
                <Input 
                  id="maxTokens"
                  v-model.number="config.api_token.max_tokens" 
                  type="number"
                  min="1"
                  autocomplete="off"
                  :placeholder="t('admin.config.apiToken.maxTokensPlaceholder')"
                  class="mt-2"
                  :disabled="!config.api_token.enabled"
                />
                <p class="text-xs text-muted-foreground mt-1">
                  {{ t('admin.config.apiToken.maxTokensHint') }}
                </p>
              </div>
              
              <div class="flex items-center justify-between p-4 rounded-lg bg-muted/50">
                <div>
                  <Label for="allowAdminAPI" class="text-base font-medium">{{ t('admin.config.apiToken.allowAdminLabel') }}</Label>
                  <p class="text-sm text-muted-foreground mt-1">{{ t('admin.config.apiToken.allowAdminDesc') }}</p>
                </div>
                <Switch 
                  id="allowAdminAPI"
                  v-model="config.api_token.allow_admin_api"
                  :disabled="!config.api_token.enabled"
                />
              </div>
            </div>
          </CardContent>
        </Card>

      <!-- 存储配置 -->
      <Card class="xl:col-span-2 shadow-sm border-l-4 border-l-orange-500">
        <CardHeader>
          <CardTitle class="flex items-center">
            <svg class="w-5 h-5 mr-2 text-orange-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 19a2 2 0 01-2-2V7a2 2 0 012-2h4l2 2h4a2 2 0 012 2v1M5 19h14a2 2 0 002-2v-5a2 2 0 00-2-2H9a2 2 0 00-2 2v5a2 2 0 01-2 2z" />
            </svg>
            {{ t('admin.config.storage.title') }}
          </CardTitle>
          <CardDescription>{{ t('admin.config.storage.description') }}</CardDescription>
        </CardHeader>
        <CardContent class="space-y-6">
          <div>
            <Label class="text-sm font-medium">{{ t('admin.config.storage.typeLabel') }}</Label>
            <Select v-model="config.storage.type" class="mt-2">
              <SelectTrigger>
                <SelectValue :placeholder="t('admin.config.storage.typePlaceholder')" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="local">{{ t('admin.config.storage.typeLocal') }}</SelectItem>
                <SelectItem value="s3">{{ t('admin.config.storage.typeS3') }}</SelectItem>
                <SelectItem value="webdav">{{ t('admin.config.storage.typeWebdav') }}</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <!-- 本地存储配置 -->
          <div v-if="config.storage.type === 'local'" class="space-y-4 border rounded-lg p-6 bg-gradient-to-r from-blue-50 to-indigo-50 dark:from-blue-950/20 dark:to-indigo-950/20">
            <h4 class="font-semibold text-lg flex items-center">
              <svg class="w-5 h-5 mr-2 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              {{ t('admin.config.storage.localTitle') }}
            </h4>
            <div>
              <Label for="localPath" class="text-sm font-medium">{{ t('admin.config.storage.pathLabel') }}</Label>
              <Input 
                id="localPath"
                v-model="config.storage.local!.path" 
                :placeholder="t('admin.config.storage.pathPlaceholder')"
                class="mt-2"
              />
              <p class="text-xs text-muted-foreground mt-1">{{ t('admin.config.storage.pathHint') }}</p>
            </div>
          </div>

          <!-- S3 存储配置 -->
          <div v-if="config.storage.type === 's3'" class="space-y-6 border rounded-lg p-6 bg-gradient-to-r from-green-50 to-emerald-50 dark:from-green-950/20 dark:to-emerald-950/20">
            <h4 class="font-semibold text-lg flex items-center">
              <svg class="w-5 h-5 mr-2 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M9 12l2 2 4-4" />
              </svg>
              {{ t('admin.config.storage.s3Title') }}
            </h4>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <Label for="s3Endpoint" class="text-sm font-medium">{{ t('admin.config.storage.s3EndpointLabel') }}</Label>
                <Input 
                  id="s3Endpoint"
                  v-model="config.storage.s3!.endpoint" 
                  autocomplete="off"
                  :placeholder="t('admin.config.storage.s3EndpointPlaceholder')"
                  class="mt-2"
                />
              </div>
              <div>
                <Label for="s3Region" class="text-sm font-medium">{{ t('admin.config.storage.s3RegionLabel') }}</Label>
                <Input 
                  id="s3Region"
                  v-model="config.storage.s3!.region" 
                  autocomplete="off"
                  :placeholder="t('admin.config.storage.s3RegionPlaceholder')"
                  class="mt-2"
                />
              </div>
              <div>
                <Label for="s3Bucket" class="text-sm font-medium">{{ t('admin.config.storage.s3BucketLabel') }}</Label>
                <Input 
                  id="s3Bucket"
                  v-model="config.storage.s3!.bucket" 
                  autocomplete="off"
                  :placeholder="t('admin.config.storage.s3BucketPlaceholder')"
                  class="mt-2"
                />
              </div>
              <div>
                <Label for="s3AccessKey" class="text-sm font-medium">{{ t('admin.config.storage.s3AccessKeyLabel') }}</Label>
                <Input 
                  id="s3AccessKey"
                  v-model="config.storage.s3!.access_key" 
                  autocomplete="off"
                  :placeholder="t('admin.config.storage.s3AccessKeyPlaceholder')"
                  class="mt-2"
                />
              </div>
              <div class="md:col-span-2">
                <Label for="s3SecretKey" class="text-sm font-medium">{{ t('admin.config.storage.s3SecretKeyLabel') }}</Label>
                <Input 
                  id="s3SecretKey"
                  v-model="config.storage.s3!.secret_key" 
                  type="password"
                  autocomplete="off"
                  :placeholder="t('admin.config.storage.s3SecretKeyPlaceholder')"
                  class="mt-2"
                />
              </div>
              <div class="md:col-span-2">
                <div class="flex items-center justify-between p-4 rounded-lg bg-white/80 dark:bg-gray-900/50">
                  <div>
                    <Label for="s3UseSSL" class="text-base font-medium">{{ t('admin.config.storage.s3UseSslLabel') }}</Label>
                    <p class="text-sm text-muted-foreground mt-1">{{ t('admin.config.storage.s3UseSslDesc') }}</p>
                  </div>
                  <Switch 
                    id="s3UseSSL"
                    v-model="config.storage.s3!.use_ssl"
                  />
                </div>
              </div>
            </div>
          </div>

          <!-- WebDAV 存储配置 -->
          <div v-if="config.storage.type === 'webdav'" class="space-y-6 border rounded-lg p-6 bg-gradient-to-r from-purple-50 to-violet-50 dark:from-purple-950/20 dark:to-violet-950/20">
            <h4 class="font-semibold text-lg flex items-center">
              <svg class="w-5 h-5 mr-2 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5.636 18.364a9 9 0 010-12.728m12.728 0a9 9 0 010 12.728m-9.9-2.829a5 5 0 010-7.07m7.072 0a5 5 0 010 7.07M13 12a1 1 0 11-2 0 1 1 0 012 0z" />
              </svg>
              {{ t('admin.config.storage.webdavTitle') }}
            </h4>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <Label for="webdavUrl" class="text-sm font-medium">{{ t('admin.config.storage.webdavUrlLabel') }}</Label>
                <Input 
                  id="webdavUrl"
                  v-model="config.storage.webdav!.url" 
                  autocomplete="off"
                  :placeholder="t('admin.config.storage.webdavUrlPlaceholder')"
                  class="mt-2"
                />
              </div>
              <div>
                <Label for="webdavRoot" class="text-sm font-medium">{{ t('admin.config.storage.webdavRootLabel') }}</Label>
                <Input 
                  id="webdavRoot"
                  v-model="config.storage.webdav!.root" 
                  autocomplete="off"
                  :placeholder="t('admin.config.storage.webdavRootPlaceholder')"
                  class="mt-2"
                />
              </div>
              <div>
                <Label for="webdavUsername" class="text-sm font-medium">{{ t('admin.config.storage.webdavUsernameLabel') }}</Label>
                <Input 
                  id="webdavUsername"
                  v-model="config.storage.webdav!.username" 
                  autocomplete="off"
                  :placeholder="t('admin.config.storage.webdavUsernamePlaceholder')"
                  class="mt-2"
                />
              </div>
              <div>
                <Label for="webdavPassword" class="text-sm font-medium">{{ t('admin.config.storage.webdavPasswordLabel') }}</Label>
                <Input 
                  id="webdavPassword"
                  v-model="config.storage.webdav!.password" 
                  type="password"
                  autocomplete="off"
                  :placeholder="t('admin.config.storage.webdavPasswordPlaceholder')"
                  class="mt-2"
                />
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- 数据库配置 -->
      <Card class="xl:col-span-2 shadow-sm border-l-4 border-l-indigo-500">
        <CardHeader>
          <CardTitle class="flex items-center">
            <svg class="w-5 h-5 mr-2 text-indigo-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4m0 5c0 2.21-3.582 4-8 4s-8-1.79-8-4" />
            </svg>
            {{ t('admin.config.database.title') }}
          </CardTitle>
          <CardDescription>{{ t('admin.config.database.description') }}</CardDescription>
        </CardHeader>
        <CardContent class="space-y-6">
          <div>
            <Label class="text-sm font-medium">{{ t('admin.config.database.typeLabel') }}</Label>
            <Select v-model="config.database.type" class="mt-2">
              <SelectTrigger>
                <SelectValue :placeholder="t('admin.config.database.typePlaceholder')" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="sqlite">SQLite</SelectItem>
                <SelectItem value="mysql">MySQL</SelectItem>
                <SelectItem value="postgres">PostgreSQL</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <!-- SQLite 配置 -->
          <div v-if="config.database.type === 'sqlite'" class="space-y-4 border rounded-lg p-6 bg-gradient-to-r from-indigo-50 to-blue-50 dark:from-indigo-950/20 dark:to-blue-950/20">
            <h4 class="font-semibold text-lg flex items-center">
              <svg class="w-5 h-5 mr-2 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              {{ t('admin.config.database.sqlite') }}
            </h4>
            <div>
              <Label for="databasePath" class="text-sm font-medium">{{ t('admin.config.database.pathLabel') }}</Label>
              <Input 
                id="databasePath"
                v-model="config.database.path" 
                autocomplete="off"
                :placeholder="t('admin.config.database.pathPlaceholder')"
                class="mt-2"
              />
              <p class="text-xs text-muted-foreground mt-1">{{ t('admin.config.database.pathHint') }}</p>
            </div>
          </div>

          <!-- MySQL/PostgreSQL 配置 -->
          <div v-if="config.database.type === 'mysql' || config.database.type === 'postgres'" class="space-y-6 border rounded-lg p-6 bg-gradient-to-r from-indigo-50 to-purple-50 dark:from-indigo-950/20 dark:to-purple-950/20">
            <h4 class="font-semibold text-lg flex items-center">
              <svg class="w-5 h-5 mr-2 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M5 12a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v4a2 2 0 01-2 2M5 12a2 2 0 00-2 2v4a2 2 0 002 2h14a2 2 0 002-2v-4a2 2 0 00-2-2m-2-4h.01M17 16h.01" />
              </svg>
              {{ config.database.type === 'mysql' ? t('admin.config.database.mysql') : t('admin.config.database.postgres') }}
            </h4>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <Label for="dbHost" class="text-sm font-medium">{{ t('admin.config.database.hostLabel') }}</Label>
                <Input 
                  id="dbHost"
                  v-model="config.database.host" 
                  autocomplete="off"
                  :placeholder="t('admin.config.database.hostPlaceholder')"
                  class="mt-2"
                />
              </div>
              <div>
                <Label for="dbPort" class="text-sm font-medium">{{ t('admin.config.database.portLabel') }}</Label>
                <Input 
                  id="dbPort"
                  v-model.number="config.database.port" 
                  type="number"
                  min="1"
                  max="65535"
                  autocomplete="off"
                  :placeholder="t('admin.config.database.portPlaceholder')"
                  class="mt-2"
                />
              </div>
              <div>
                <Label for="dbUser" class="text-sm font-medium">{{ t('admin.config.database.userLabel') }}</Label>
                <Input 
                  id="dbUser"
                  v-model="config.database.user" 
                  autocomplete="off"
                  :placeholder="t('admin.config.database.userPlaceholder')"
                  class="mt-2"
                />
              </div>
              <div>
                <Label for="dbPassword" class="text-sm font-medium">{{ t('admin.config.database.passwordLabel') }}</Label>
                <Input 
                  id="dbPassword"
                  v-model="config.database.password" 
                  type="password"
                  autocomplete="off"
                  :placeholder="t('admin.config.database.passwordPlaceholder')"
                  class="mt-2"
                />
              </div>
              <div>
                <Label for="dbName" class="text-sm font-medium">{{ t('admin.config.database.dbnameLabel') }}</Label>
                <Input 
                  id="dbName"
                  v-model="config.database.dbname" 
                  autocomplete="off"
                  :placeholder="t('admin.config.database.dbnamePlaceholder')"
                  class="mt-2"
                />
              </div>
              <div>
                <Label for="dbConfig" class="text-sm font-medium">{{ t('admin.config.database.configLabel') }}</Label>
                <Input 
                  id="dbConfig"
                  v-model="config.database.config" 
                  autocomplete="off"
                  :placeholder="t('admin.config.database.configPlaceholder')"
                  class="mt-2"
                />
                <p class="text-xs text-muted-foreground mt-1">{{ t('admin.config.database.configHint') }}</p>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- Web 配置 -->
      <Card class="shadow-sm border-l-4 border-l-cyan-500">
        <CardHeader>
          <CardTitle class="flex items-center">
            <svg class="w-5 h-5 mr-2 text-cyan-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9" />
            </svg>
            {{ t('admin.config.web.title') }}
          </CardTitle>
          <CardDescription>{{ t('admin.config.web.description') }}</CardDescription>
        </CardHeader>
        <CardContent>
          <div>
            <Label for="webPath" class="text-sm font-medium">{{ t('admin.config.web.pathLabel') }}</Label>
            <Input 
              id="webPath"
              v-model="config.web.path" 
              autocomplete="off"
              :placeholder="t('admin.config.web.pathPlaceholder')"
              class="mt-2"
            />
            <p class="text-xs text-muted-foreground mt-1">{{ t('admin.config.web.pathHint') }}</p>
          </div>
        </CardContent>
      </Card>

      <!-- 日志配置 -->
      <Card class="shadow-sm border-l-4 border-l-yellow-500">
        <CardHeader>
          <CardTitle class="flex items-center">
            <svg class="w-5 h-5 mr-2 text-yellow-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            {{ t('admin.config.log.title') }}
          </CardTitle>
          <CardDescription>{{ t('admin.config.log.description') }}</CardDescription>
        </CardHeader>
        <CardContent class="space-y-5">
          <div>
            <Label for="logLevel" class="text-sm font-medium">{{ t('admin.config.log.levelLabel') }}</Label>
            <Select v-model="config.log.level" class="mt-2">
              <SelectTrigger id="logLevel">
                <SelectValue :placeholder="t('admin.config.log.levelPlaceholder')" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="debug">{{ t('admin.config.log.levelDebug') }}</SelectItem>
                <SelectItem value="info">{{ t('admin.config.log.levelInfo') }}</SelectItem>
                <SelectItem value="warn">{{ t('admin.config.log.levelWarn') }}</SelectItem>
                <SelectItem value="error">{{ t('admin.config.log.levelError') }}</SelectItem>
              </SelectContent>
            </Select>
            <p class="text-xs text-muted-foreground mt-1">{{ t('admin.config.log.levelHint') }}</p>
          </div>
          <div>
            <Label for="logFilePath" class="text-sm font-medium">{{ t('admin.config.log.filePathLabel') }}</Label>
            <Input 
              id="logFilePath"
              v-model="config.log.file_path" 
              autocomplete="off"
              :placeholder="t('admin.config.log.filePathPlaceholder')"
              class="mt-2"
            />
            <p class="text-xs text-muted-foreground mt-1">{{ t('admin.config.log.filePathHint') }}</p>
          </div>
        </CardContent>
      </Card>
    </div>
    </form>

    <!-- Toast 组件 -->
    <Toaster />
    </div>
  </AdminLayout>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { toast } from 'vue-sonner'
import { adminApi, type SystemConfig } from '@/lib/api'
import { usePublicConfig } from '@/composables/usePublicConfig'
import { useI18n } from '@/composables/useI18n'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Toaster } from '@/components/ui/sonner'
import AdminLayout from '@/layouts/AdminLayout.vue'

const { t } = useI18n()

// 使用 toast
const showToast = toast

const loading = ref(false)
const adminPassword = ref('')
const originalConfig = ref<SystemConfig>()

const config = reactive<SystemConfig>({
  site: {
    name: '',
    description: '',
    logo: '',
    base_url: '',
    port: 8080
  },
  upload: {
    max_file_size_mb: 100,
    max_batch_files: 20,
    max_retention_days: 30,
    require_token: false
  },
  storage: {
    type: 'local',
    local: {
      path: 'storage_data'
    },
    s3: {
      endpoint: 's3.amazonaws.com',
      region: 'us-east-1',
      bucket: 'file-relay-bucket',
      access_key: 'your-access-key',
      secret_key: 'your-secret-key',
      use_ssl: false
    },
    webdav: {
      url: 'https://dav.example.com',
      username: 'user',
      password: 'pass',
      root: '/file-relay'
    }
  },
  security: {
    pickup_code_length: 6,
    pickup_fail_limit: 5,
    admin_password_hash: '',
    jwt_secret: ''
  },
  database: {
    type: 'sqlite',
    path: 'file_relay.db',
    host: '127.0.0.1',
    port: 3306,
    user: 'root',
    password: '',
    dbname: 'file_relay',
    config: ''
  },
  api_token: {
    enabled: true,
    max_tokens: 20,
    allow_admin_api: false
  },
  web: {
    path: 'web'
  },
  log: {
    level: 'info',
    file_path: ''
  }
})

const loadConfig = async () => {
  try {
    loading.value = true
    const response = await adminApi.getConfig()
    const configData = response.data.data  // 从ApiResponse中提取data字段
    
    // 直接赋值，保持 snake_case 格式
    Object.assign(config, configData)
    
    // 保存原始配置副本
    originalConfig.value = JSON.parse(JSON.stringify(config))
    
    console.log('配置加载成功:', config)
  } catch (error) {
    console.error('加载配置失败:', error)
    showToast.error(t('admin.config.loadFailed'))
  } finally {
    loading.value = false
  }
}

const resetForm = () => {
  if (originalConfig.value) {
    Object.assign(config, JSON.parse(JSON.stringify(originalConfig.value)))
    adminPassword.value = ''
    showToast.success(t('admin.config.resetSuccess'))
  }
}

const saveConfig = async () => {
  try {
    loading.value = true
    
    // 构建要提交的配置
    const configToSave = JSON.parse(JSON.stringify(config))
    
    // 如果设置了新密码，添加到 security 配置中
    if (adminPassword.value) {
      configToSave.security.admin_password = adminPassword.value
    }
    
    await adminApi.updateConfig(configToSave)
    
    // 更新原始配置
    originalConfig.value = JSON.parse(JSON.stringify(config))
    adminPassword.value = ''
    
    // 强制刷新公共配置，确保前端显示最新配置
    const { loadConfig: refreshPublicConfig } = usePublicConfig()
    await refreshPublicConfig(true)
    
    showToast.success(t('admin.config.saveSuccess'))
  } catch (error: any) {
    console.error('保存配置失败:', error)
    showToast.error(error.response?.data?.msg || t('admin.config.saveFailed'))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadConfig()
})
</script>