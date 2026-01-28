<script setup lang="ts">
  import 'vue-sonner/style.css'
</script>

<template>
  <div id="app">
    <router-view />
  </div>
</template>

<style>
#app {
  font-family: ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, "Noto Sans", sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

html, body {
  margin: 0;
  padding: 0;
  height: 100%;
}

#app {
  min-height: 100%;
}

/* 全局滚动条样式 - 默认隐藏，滚动时显示 */
* {
  scrollbar-width: thin;
  scrollbar-color: transparent transparent;
}

*:hover {
  scrollbar-color: rgba(156, 163, 175, 0.3) transparent;
}

.dark *:hover {
  scrollbar-color: rgba(156, 163, 175, 0.5) transparent;
}

/* Webkit 滚动条样式 */
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

::-webkit-scrollbar-track {
  background: transparent;
}

::-webkit-scrollbar-thumb {
  background: transparent;
  border-radius: 3px;
  transition: background 0.2s;
}

*:hover::-webkit-scrollbar-thumb {
  background: rgba(156, 163, 175, 0.3);
}

.dark *:hover::-webkit-scrollbar-thumb {
  background: rgba(156, 163, 175, 0.5);
}

::-webkit-scrollbar-thumb:hover {
  background: rgba(156, 163, 175, 0.5) !important;
}

.dark ::-webkit-scrollbar-thumb:hover {
  background: rgba(156, 163, 175, 0.7) !important;
}
</style>
