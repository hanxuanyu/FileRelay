import path from 'node:path'
import { defineConfig } from 'vite'
import tailwindcss from '@tailwindcss/vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue(), tailwindcss()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  build: {
    // 输出目录设置为上级目录的web文件夹
    outDir: path.resolve(__dirname, '../web'),
    // 允许清空项目外的输出目录
    emptyOutDir: true,
    // 启用代码压缩
    minify: 'terser',
    terserOptions: {
      compress: {
        // 移除console和debugger
        drop_console: true,
        drop_debugger: true,
        // 移除未使用的代码
        pure_funcs: ['console.log', 'console.info', 'console.debug'],
      },
      mangle: {
        // 启用变量名混淆
        safari10: true,
      },
      format: {
        // 移除注释
        comments: false,
      },
    },
    // 代码分割配置 - 合并为尽量少的文件
    rollupOptions: {
      output: {
        // 禁用代码分割，合并为单个JS文件
        inlineDynamicImports: true,
        // 配置输出文件名
        entryFileNames: 'assets/js/[name]-[hash].js',
        chunkFileNames: 'assets/js/[name]-[hash].js',
        assetFileNames: (assetInfo) => {
          // CSS文件
          if (assetInfo.name?.endsWith('.css')) {
            return 'assets/css/[name]-[hash][extname]'
          }
          // 图片和其他资源
          return 'assets/[ext]/[name]-[hash][extname]'
        },
      },
    },
    // 启用CSS代码分割
    cssCodeSplit: false,
    // 设置chunk大小警告限制
    chunkSizeWarningLimit: 2000,
    // 禁用source map以减小文件大小
    sourcemap: false,
    // 启用gzip压缩提示
    reportCompressedSize: true,
  },
})
