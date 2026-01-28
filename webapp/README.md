# 文件中转站前端

一个基于 Vue 3 + TypeScript + Vite + Shadcn-Vue 构建的文件中转站前端应用。

## 项目简介

文件中转站是一个极简的文件暂存服务，用户可以像使用现实中的暂存柜一样，将文件、文本等内容临时存储在网站中，并获得一个取件码用于在其他设备上提取文件。

## 主要功能

### 用户功能
- **文件上传**: 支持批量上传多种格式的文件
- **文本分享**: 快速分享长文本内容
- **取件码提取**: 使用取件码获取文件或文本
- **过期策略**: 支持按时间或下载次数设置过期规则
- **快捷操作**: 一键复制、打包下载等便捷功能

### 管理功能
- **批次管理**: 查看、编辑、删除文件批次
- **API Token**: 创建和管理 API 访问凭证
- **数据统计**: 系统运行状态和使用统计
- **权限控制**: 管理员密码保护的后台系统

## 技术栈

- **框架**: Vue 3 + TypeScript
- **构建工具**: Vite
- **UI 组件**: Shadcn-Vue (基于 Tailwind CSS)
- **路由**: Vue Router
- **HTTP 客户端**: Axios
- **状态管理**: 组合式 API
- **代码规范**: ESLint + Prettier

## 快速开始

### 环境要求

- Node.js 18+
- npm 或 yarn 或 pnpm

### 安装依赖

```bash
npm install
```

### 开发环境

```bash
npm run dev
```

启动后访问 `http://localhost:5173`

### 构建生产版本

```bash
npm run build
```

### 预览构建结果

```bash
npm run preview
```

## 环境配置

创建环境变量文件：

### `.env` (通用配置)
```env
VITE_API_URL=http://localhost:8080
VITE_APP_TITLE=文件中转站
VITE_APP_DESCRIPTION=安全便捷的文件暂存服务
```

### `.env.development` (开发环境)
```env
VITE_API_URL=http://localhost:8080
```

### `.env.production` (生产环境)
```env
VITE_API_URL=/api
```

## 页面路由

### 用户页面
- `/` - 首页 (取件/存件切换)
- `/upload` - 文件上传页面  
- `/pickup` - 取件页面
- `/pickup/:code` - 直接取件 (URL 包含取件码)

### 管理页面 (隐藏入口)
- `/admin/login` - 管理员登录
- `/admin` - 管理概览
- `/admin/batches` - 批次管理
- `/admin/tokens` - API Token 管理

## 功能特点

### 用户体验
- **极简设计**: 首页默认取件，一键切换存件模式
- **智能识别**: 自动识别文件类型并显示对应图标
- **快捷操作**: 支持剪贴板粘贴取件码、一键复制等
- **进度反馈**: 上传进度显示和状态提示
- **响应式设计**: 完美适配桌面和移动设备

### 管理功能
- **数据统计**: 实时显示系统运行数据
- **批次管理**: 支持搜索、筛选、分页查看
- **权限控制**: Token 可设置权限范围和过期时间
- **操作日志**: 详细的操作记录和状态跟踪

## 许可证

MIT License