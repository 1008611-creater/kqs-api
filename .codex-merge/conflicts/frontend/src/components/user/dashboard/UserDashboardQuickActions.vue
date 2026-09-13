<template>
  <section class="quick-panel h-full">
    <div class="quick-panel-header">
      <p class="text-xs font-bold uppercase tracking-[0.18em] text-primary-700 dark:text-primary-200">
        {{ t('dashboard.quickActions') }}
      </p>
      <h2 class="mt-2 text-lg font-bold text-primary-950 dark:text-white">
        {{ t('dashboard.createApiKey') }}
      </h2>
    </div>

<<<<<<< HEAD
      <button @click="router.push('/usage')" class="group flex w-full items-center gap-4 rounded-xl bg-gray-50 p-4 text-left transition-all duration-200 hover:bg-gray-100 dark:bg-dark-800/50 dark:hover:bg-dark-800">
        <div class="flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-xl bg-emerald-100 transition-transform group-hover:scale-105 dark:bg-emerald-900/30">
          <Icon name="chart" size="lg" class="text-emerald-600 dark:text-emerald-400" />
        </div>
        <div class="min-w-0 flex-1">
          <p class="text-sm font-medium text-gray-900 dark:text-white">{{ t('dashboard.viewUsage') }}</p>
          <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('dashboard.checkDetailedLogs') }}</p>
        </div>
        <Icon
          name="chevronRight"
          size="md"
          class="text-gray-400 transition-colors group-hover:text-emerald-500 dark:text-dark-500"
        />
      </button>

      <button v-if="canUseBatchImage" @click="router.push('/batch-image')" class="group flex w-full items-center gap-4 rounded-xl bg-gray-50 p-4 text-left transition-all duration-200 hover:bg-gray-100 dark:bg-dark-800/50 dark:hover:bg-dark-800">
        <div class="flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-xl bg-sky-100 transition-transform group-hover:scale-105 dark:bg-sky-900/30">
          <Icon name="sparkles" size="lg" class="text-sky-600 dark:text-sky-400" />
        </div>
        <div class="min-w-0 flex-1">
          <p class="text-sm font-medium text-gray-900 dark:text-white">{{ t('dashboard.batchImageAgent') }}</p>
          <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('dashboard.batchImageAgentDesc') }}</p>
        </div>
        <Icon
          name="chevronRight"
          size="md"
          class="text-gray-400 transition-colors group-hover:text-sky-500 dark:text-dark-500"
        />
      </button>

      <button @click="router.push('/redeem')" class="group flex w-full items-center gap-4 rounded-xl bg-gray-50 p-4 text-left transition-all duration-200 hover:bg-gray-100 dark:bg-dark-800/50 dark:hover:bg-dark-800">
        <div class="flex h-12 w-12 flex-shrink-0 items-center justify-center rounded-xl bg-amber-100 transition-transform group-hover:scale-105 dark:bg-amber-900/30">
          <Icon name="gift" size="lg" class="text-amber-600 dark:text-amber-400" />
        </div>
        <div class="min-w-0 flex-1">
          <p class="text-sm font-medium text-gray-900 dark:text-white">{{ t('dashboard.redeemCode') }}</p>
          <p class="text-xs text-gray-500 dark:text-dark-400">{{ t('dashboard.addBalanceWithCode') }}</p>
        </div>
        <Icon
          name="chevronRight"
          size="md"
          class="text-gray-400 transition-colors group-hover:text-amber-500 dark:text-dark-500"
        />
=======
    <div class="flex flex-1 flex-col space-y-2.5 p-3">
      <button
        v-for="action in actions"
        :key="action.key"
        @click="handleAction(action)"
        class="quick-action"
      >
        <span class="quick-action-icon">
          <Icon :name="action.icon" size="md" />
        </span>
        <span class="min-w-0 flex-1">
          <span class="block text-sm font-semibold text-primary-950 dark:text-white">
            {{ action.title }}
          </span>
          <span class="mt-0.5 block truncate text-xs text-primary-900/52 dark:text-white/46">
            {{ action.description }}
          </span>
        </span>
        <Icon name="chevronRight" size="sm" class="text-primary-900/35 dark:text-white/35" />
>>>>>>> 1157e4cfe6271a81c9dc96df84af6c6a3b22f831
      </button>
    </div>
  </section>
</template>

<script setup lang="ts">
<<<<<<< HEAD
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import { useBatchImageAccess } from '@/composables/useBatchImageAccess'
const router = useRouter()
const { t } = useI18n()
const { canUseBatchImage, refreshBatchImageAccess } = useBatchImageAccess()

onMounted(() => {
  void refreshBatchImageAccess()
})
=======
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'

const router = useRouter()
const { t } = useI18n()

type QuickActionIcon = 'key' | 'chart' | 'gift' | 'creditCard'

interface QuickAction {
  key: string
  icon: QuickActionIcon
  title: string
  description: string
  to?: string
  href?: string
}

const actions = computed<QuickAction[]>(() => {
  const items: QuickAction[] = []

  items.push({
    key: 'purchase',
    to: '/purchase',
    icon: 'creditCard',
    title: '购买余额',
    description: '跳转链动小铺购买卡密'
  })

  items.push(
  {
    key: 'redeem',
    to: '/redeem',
    icon: 'gift',
    title: t('dashboard.redeemCode'),
    description: t('dashboard.addBalanceWithCode')
  },
  {
    key: 'keys',
    to: '/keys',
    icon: 'key',
    title: t('dashboard.createApiKey'),
    description: t('dashboard.generateNewKey')
  },
  {
    key: 'usage',
    to: '/usage',
    icon: 'chart',
    title: t('dashboard.viewUsage'),
    description: t('dashboard.checkDetailedLogs')
  }
  )

  return items
})

const handleAction = (action: QuickAction) => {
  if (action.to) {
    router.push(action.to)
    return
  }
  if (action.href) {
    window.open(action.href, '_blank', 'noopener,noreferrer')
  }
}
>>>>>>> 1157e4cfe6271a81c9dc96df84af6c6a3b22f831
</script>

<style scoped>
.quick-panel {
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

:global(.dark) .quick-panel {
  border-color: rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.055);
  box-shadow: 0 16px 42px rgba(0, 0, 0, 0.16);
}

.quick-panel-header {
  border-bottom: 1px solid rgba(2, 43, 18, 0.08);
  padding: 1.1rem 1.25rem;
}

:global(.dark) .quick-panel-header {
  border-color: rgba(255, 255, 255, 0.1);
}

.quick-action {
  display: flex;
  width: 100%;
  align-items: center;
  gap: 0.9rem;
  border-radius: 1.1rem;
  padding: 0.95rem;
  text-align: left;
  transition: background 180ms ease, transform 180ms ease;
}

.quick-action:hover {
  background: rgba(0, 135, 60, 0.07);
  transform: translateY(-1px);
}

:global(.dark) .quick-action:hover {
  background: rgba(255, 255, 255, 0.07);
}

.quick-action:active {
  transform: translateY(0);
}

.quick-action-icon {
  display: flex;
  width: 2.75rem;
  height: 2.75rem;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  border-radius: 1rem;
  background: rgba(0, 135, 60, 0.1);
  color: #006b31;
}

:global(.dark) .quick-action-icon {
  background: rgba(0, 135, 60, 0.16);
  color: #bde8cb;
}
</style>
