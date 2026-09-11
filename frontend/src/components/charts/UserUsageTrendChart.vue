<template>
  <div class="card p-4">
    <div class="mb-4 flex items-center justify-between gap-3">
      <h3 class="text-sm font-semibold text-gray-900 dark:text-white">
        {{ t('admin.dashboard.userUsageTrend') }}
      </h3>
      <span class="text-xs font-medium text-gray-500 dark:text-gray-400">
        {{ t('admin.dashboard.actual') }}
      </span>
    </div>
    <div v-if="loading" class="flex h-56 items-center justify-center">
      <LoadingSpinner />
    </div>
    <div v-else-if="userTrendData.length > 0 && chartData" class="h-56">
      <Line :data="chartData" :options="lineOptions" />
    </div>
    <div
      v-else
      class="flex h-56 items-center justify-center text-sm text-gray-500 dark:text-gray-400"
    >
      {{ t('admin.dashboard.noDataAvailable') }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  CategoryScale,
  Chart as ChartJS,
  Filler,
  Legend,
  LinearScale,
  LineElement,
  PointElement,
  Tooltip
} from 'chart.js'
import { Line } from 'vue-chartjs'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import type { UserUsageTrendPoint } from '@/types'

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Tooltip, Legend, Filler)

const { t } = useI18n()

const props = defineProps<{
  userTrendData: UserUsageTrendPoint[]
  loading?: boolean
}>()

const palette = [
  '#0f766e', '#2563eb', '#c2410c', '#7c3aed', '#be123c', '#15803d',
  '#0369a1', '#a16207', '#9333ea', '#0f766e', '#dc2626', '#475569'
]

const isDarkMode = computed(() => document.documentElement.classList.contains('dark'))

const chartColors = computed(() => ({
  text: isDarkMode.value ? '#e5e7eb' : '#374151',
  grid: isDarkMode.value ? '#374151' : '#e5e7eb'
}))

const labels = computed(() => Array.from(new Set(props.userTrendData.map((point) => point.date))))

const users = computed(() => {
  const result = new Map<number, UserUsageTrendPoint>()
  for (const point of props.userTrendData) {
    if (!result.has(point.user_id)) result.set(point.user_id, point)
  }
  return Array.from(result.values())
})

const userLabel = (point: UserUsageTrendPoint): string => {
  if (point.username?.trim()) return `${point.username.trim()} (#${point.user_id})`
  if (point.email?.trim()) return `${point.email.trim()} (#${point.user_id})`
  return `#${point.user_id}`
}

const chartData = computed(() => {
  if (!labels.value.length || !users.value.length) return null

  return {
    labels: labels.value,
    datasets: users.value.map((user, index) => {
      const color = palette[index % palette.length]
      return {
        label: userLabel(user),
        data: labels.value.map((date) => {
          const point = props.userTrendData.find(
            (candidate) => candidate.user_id === user.user_id && candidate.date === date
          )
          return point?.actual_cost ?? 0
        }),
        borderColor: color,
        backgroundColor: `${color}18`,
        borderWidth: 2,
        pointRadius: 2,
        pointHoverRadius: 4,
        fill: false,
        tension: 0.25,
        spanGaps: true
      }
    })
  }
})

const lineOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  interaction: { intersect: false, mode: 'index' as const },
  plugins: {
    legend: {
      position: 'top' as const,
      labels: {
        color: chartColors.value.text,
        usePointStyle: true,
        pointStyle: 'circle',
        padding: 12,
        font: { size: 10 }
      }
    },
    tooltip: {
      callbacks: {
        label: (context: any) => `${context.dataset.label}: $${formatCost(Number(context.raw))}`
      }
    }
  },
  scales: {
    x: {
      grid: { color: chartColors.value.grid },
      ticks: { color: chartColors.value.text, font: { size: 10 }, maxTicksLimit: 18 }
    },
    y: {
      beginAtZero: true,
      grid: { color: chartColors.value.grid },
      ticks: {
        color: chartColors.value.text,
        font: { size: 10 },
        callback: (value: string | number) => `$${formatCost(Number(value))}`
      }
    }
  }
}))

const formatCost = (value: number): string => {
  if (value >= 1000) return `${(value / 1000).toFixed(2)}K`
  if (value >= 1) return value.toFixed(2)
  if (value >= 0.01) return value.toFixed(3)
  return value.toFixed(4)
}
</script>
