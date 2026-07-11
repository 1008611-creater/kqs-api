<template>
  <div class="table-page-layout" :class="{ 'mobile-mode': isMobile }">
    <!-- 固定区域：操作按钮 -->
    <div v-if="$slots.actions" class="layout-section-fixed">
      <slot name="actions" />
    </div>

    <!-- 固定区域：搜索和过滤器 -->
    <div v-if="$slots.filters" class="layout-section-fixed">
      <slot name="filters" />
    </div>

    <!-- 滚动区域：表格 -->
    <div class="layout-section-scrollable">
      <div class="card table-scroll-container">
        <slot name="table" />
      </div>
    </div>

    <!-- 固定区域：分页器 -->
    <div v-if="$slots.pagination" class="layout-section-fixed">
      <slot name="pagination" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

const isMobile = ref(false)

const checkMobile = () => {
  isMobile.value = window.innerWidth < 1024
}

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})
</script>

<style scoped>
/* 桌面端：Flexbox 布局 */
.table-page-layout {
  @apply flex flex-col gap-6;
  height: calc(100vh - 64px - 4rem); /* 减去 header + lg:p-8 的上下padding */
}

.layout-section-fixed {
  @apply flex-shrink-0;
  border: 1px solid rgba(2, 43, 18, 0.08);
  border-radius: 1.35rem;
  background: rgba(255, 255, 255, 0.74);
  box-shadow: 0 14px 38px rgba(2, 43, 18, 0.055);
  padding: 0.9rem;
  backdrop-filter: blur(18px);
}

:global(.dark) .layout-section-fixed {
  border-color: rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.055);
  box-shadow: 0 14px 38px rgba(0, 0, 0, 0.14);
}

.layout-section-scrollable {
  @apply flex-1 min-h-0 flex flex-col;
}

/* 表格滚动容器 - 增强版表体滚动方案 */
.table-scroll-container {
  @apply flex flex-col overflow-hidden h-full rounded-2xl;
  border: 1px solid rgba(2, 43, 18, 0.08);
  background: rgba(255, 255, 255, 0.78);
  box-shadow: 0 18px 50px rgba(2, 43, 18, 0.065);
  backdrop-filter: blur(18px);
}

:global(.dark) .table-scroll-container {
  border-color: rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.055);
  box-shadow: 0 18px 50px rgba(0, 0, 0, 0.16);
}

.table-scroll-container :deep(.table-wrapper) {
  @apply flex-1 overflow-x-auto overflow-y-auto;
  /* 确保横向滚动条显示在最底部 */
  scrollbar-gutter: stable;
}

.table-scroll-container :deep(table) {
  @apply w-full;
  min-width: max-content; /* 关键：确保表格宽度根据内容撑开，从而触发横向滚动 */
  display: table; /* 使用标准 table 布局以支持 sticky 列 */
}

.table-scroll-container :deep(thead) {
  background: rgba(239, 250, 242, 0.84);
  backdrop-filter: blur(12px);
}

:global(.dark) .table-scroll-container :deep(thead) {
  background: rgba(255, 255, 255, 0.055);
}

.table-scroll-container :deep(tbody) {
  /* 保持默认 table-row-group 显示，不使用 block */
}

.table-scroll-container :deep(th) {
  @apply px-5 py-4 text-left text-sm font-semibold;
  border-bottom: 1px solid rgba(2, 43, 18, 0.08);
  color: rgba(2, 43, 18, 0.66);
}

:global(.dark) .table-scroll-container :deep(th) {
  border-color: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.62);
}

.table-scroll-container :deep(td) {
  @apply px-5 py-4 text-sm;
  border-bottom: 1px solid rgba(2, 43, 18, 0.055);
  color: rgba(2, 43, 18, 0.74);
}

:global(.dark) .table-scroll-container :deep(td) {
  border-color: rgba(255, 255, 255, 0.055);
  color: rgba(255, 255, 255, 0.72);
}

/* 移动端：恢复正常滚动 */
.table-page-layout.mobile-mode .table-scroll-container {
  @apply h-auto overflow-visible border-none shadow-none bg-transparent;
  backdrop-filter: none;
}

.table-page-layout.mobile-mode .layout-section-scrollable {
  @apply flex-none min-h-fit;
}

.table-page-layout.mobile-mode .table-scroll-container :deep(.table-wrapper) {
  @apply overflow-visible;
}

.table-page-layout.mobile-mode .table-scroll-container :deep(table) {
  @apply flex-none;
  display: table;
  min-width: 100%;
}
</style>
