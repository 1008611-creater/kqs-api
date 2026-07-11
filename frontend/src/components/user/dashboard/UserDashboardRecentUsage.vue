<template>
  <section class="usage-panel h-full">
    <div class="usage-panel-header">
      <div>
        <p class="text-xs font-bold uppercase tracking-[0.18em] text-primary-700 dark:text-primary-200">
          {{ t('dashboard.last7Days') }}
        </p>
        <h2 class="mt-2 text-lg font-bold text-primary-950 dark:text-white">
          {{ t('dashboard.recentUsage') }}
        </h2>
      </div>
      <router-link
        to="/usage"
        class="hidden items-center gap-2 rounded-xl px-3 py-2 text-sm font-semibold text-primary-700 transition-colors hover:bg-primary-50 dark:text-primary-200 dark:hover:bg-white/8 sm:flex"
      >
        {{ t('dashboard.viewAllUsage') }}
        <Icon name="arrowRight" size="sm" />
      </router-link>
    </div>

    <div class="flex flex-1 flex-col p-4 sm:p-5">
      <div v-if="loading" class="space-y-3 py-2">
        <div v-for="item in 3" :key="item" class="usage-skeleton"></div>
      </div>
      <div v-else-if="data.length === 0" class="py-6">
        <EmptyState :title="t('dashboard.noUsageRecords')" :description="t('dashboard.startUsingApi')" />
      </div>
      <div v-else class="space-y-2.5">
        <div v-for="log in data" :key="log.id" class="usage-row">
          <div class="flex min-w-0 items-center gap-3">
            <div class="usage-model-mark">
              <Icon name="beaker" size="sm" />
            </div>
            <div class="min-w-0">
              <p class="truncate text-sm font-semibold text-primary-950 dark:text-white">
                {{ log.model }}
              </p>
              <p class="text-xs text-primary-900/48 dark:text-white/42">
                {{ formatDateTime(log.created_at) }}
              </p>
            </div>
          </div>
          <div class="text-right">
            <p class="text-sm font-semibold text-primary-800 dark:text-primary-200">
              ${{ formatCost(log.actual_cost) }}
              <span class="font-normal text-primary-900/35 dark:text-white/32">
                / ${{ formatCost(log.total_cost) }}
              </span>
            </p>
            <p class="text-xs text-primary-900/48 dark:text-white/42">
              {{ (log.input_tokens + log.output_tokens).toLocaleString() }} tokens
            </p>
          </div>
        </div>

        <router-link
          to="/usage"
          class="flex items-center justify-center gap-2 py-3 text-sm font-semibold text-primary-700 transition-colors hover:text-primary-900 dark:text-primary-200 dark:hover:text-white sm:hidden"
        >
          {{ t('dashboard.viewAllUsage') }}
          <Icon name="arrowRight" size="sm" />
        </router-link>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import EmptyState from '@/components/common/EmptyState.vue'
import Icon from '@/components/icons/Icon.vue'
import { formatDateTime } from '@/utils/format'
import type { UsageLog } from '@/types'

defineProps<{
  data: UsageLog[]
  loading: boolean
}>()
const { t } = useI18n()
const formatCost = (c: number) => c.toFixed(4)
</script>

<style scoped>
.usage-panel {
  display: flex;
  min-height: 100%;
  flex-direction: column;
  overflow: hidden;
  border: 1px solid rgba(2, 43, 18, 0.08);
  border-radius: 1.5rem;
  background: rgba(255, 255, 255, 0.72);
  box-shadow: 0 16px 42px rgba(2, 43, 18, 0.07);
  backdrop-filter: blur(18px);
}

:global(.dark) .usage-panel {
  border-color: rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.055);
  box-shadow: 0 16px 42px rgba(0, 0, 0, 0.16);
}

.usage-panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  border-bottom: 1px solid rgba(2, 43, 18, 0.08);
  padding: 1.1rem 1.25rem;
}

:global(.dark) .usage-panel-header {
  border-color: rgba(255, 255, 255, 0.1);
}

.usage-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  border-radius: 1.15rem;
  background: rgba(2, 43, 18, 0.035);
  padding: 0.95rem;
  transition: background 180ms ease, transform 180ms ease;
}

.usage-row:hover {
  background: rgba(0, 135, 60, 0.07);
  transform: translateY(-1px);
}

:global(.dark) .usage-row {
  background: rgba(255, 255, 255, 0.04);
}

:global(.dark) .usage-row:hover {
  background: rgba(255, 255, 255, 0.07);
}

.usage-model-mark {
  display: flex;
  width: 2.5rem;
  height: 2.5rem;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  border-radius: 0.95rem;
  background: rgba(0, 135, 60, 0.1);
  color: #006b31;
}

:global(.dark) .usage-model-mark {
  background: rgba(0, 135, 60, 0.16);
  color: #bde8cb;
}

.usage-skeleton {
  height: 4.5rem;
  border-radius: 1.1rem;
  background: linear-gradient(
    90deg,
    rgba(2, 43, 18, 0.04),
    rgba(0, 135, 60, 0.08),
    rgba(2, 43, 18, 0.04)
  );
  background-size: 200% 100%;
  animation: usage-skeleton 1.2s ease-in-out infinite;
}

:global(.dark) .usage-skeleton {
  background: linear-gradient(
    90deg,
    rgba(255, 255, 255, 0.04),
    rgba(255, 255, 255, 0.08),
    rgba(255, 255, 255, 0.04)
  );
  background-size: 200% 100%;
}

@keyframes usage-skeleton {
  0% {
    background-position: 100% 0;
  }
  100% {
    background-position: -100% 0;
  }
}
</style>
