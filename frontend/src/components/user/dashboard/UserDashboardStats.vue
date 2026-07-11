<template>
  <section class="dashboard-overview">
    <div class="dashboard-hero">
      <div class="dashboard-hero-copy min-w-0">
        <p class="dashboard-eyebrow">{{ t('dashboard.serviceDesk') }}</p>
        <h2 class="dashboard-hero-title mt-3 text-2xl font-black tracking-tight md:text-3xl">
          {{ t('dashboard.heroTitle') }}
        </h2>

        <div class="mt-6 flex flex-wrap gap-3">
          <router-link to="/keys" class="btn btn-primary min-h-11 rounded-xl px-4">
            <Icon name="key" size="md" class="mr-2" />
            {{ t('dashboard.createApiKey') }}
          </router-link>
          <router-link to="/usage" class="btn btn-secondary min-h-11 rounded-xl px-4">
            <Icon name="chart" size="md" class="mr-2" />
            {{ t('dashboard.viewUsage') }}
          </router-link>
        </div>
      </div>

      <div class="dashboard-hero-panel">
        <div class="flex items-center justify-between gap-3">
          <div>
            <p class="dashboard-hero-panel-kicker text-xs font-semibold uppercase tracking-[0.18em]">
              {{ t('dashboard.routeHealth') }}
            </p>
          </div>
          <div class="dashboard-status-dot" :class="{ 'is-ready': activeKeyCount > 0 }"></div>
        </div>

        <div class="dashboard-hero-stats-grid mt-3">
          <div v-if="!isSimple" class="dashboard-hero-stat">
            <span>{{ t('dashboard.balance') }}</span>
            <strong>${{ formatBalance(balance) }}</strong>
          </div>
          <div class="dashboard-hero-stat">
            <span>{{ t('dashboard.apiKeys') }}</span>
            <strong>{{ stats.total_api_keys || 0 }}</strong>
          </div>
          <div class="dashboard-hero-stat">
            <span>{{ t('common.active') }}</span>
            <strong>{{ activeKeyCount }}</strong>
          </div>
          <div class="dashboard-hero-stat">
            <span>{{ t('dashboard.todayCost') }}</span>
            <strong>${{ formatCost(stats.today_actual_cost || 0) }}</strong>
          </div>
        </div>
      </div>
    </div>

    <div class="dashboard-metric-grid">
      <article v-for="metric in primaryMetrics" :key="metric.key" class="dashboard-metric">
        <div class="dashboard-metric-icon">
          <Icon :name="metric.icon" size="md" />
        </div>
        <div class="min-w-0">
          <p class="dashboard-metric-label">{{ metric.label }}</p>
          <p class="dashboard-metric-value">{{ metric.value }}</p>
          <p class="dashboard-metric-detail">{{ metric.detail }}</p>
        </div>
      </article>
    </div>

  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import type { UserDashboardStats as UserStatsType } from '@/api/usage'

const props = defineProps<{
  stats: UserStatsType
  balance: number
  isSimple: boolean
}>()
const { t } = useI18n()

const activeKeyCount = computed(() => props.stats?.active_api_keys || 0)
const primaryMetrics = computed(() => [
  {
    key: 'todayRequests',
    icon: 'chart' as const,
    label: t('dashboard.todayRequests'),
    value: formatNumber(props.stats?.today_requests || 0),
    detail: `${t('common.total')}: ${formatNumber(props.stats?.total_requests || 0)}`
  },
  {
    key: 'todayTokens',
    icon: 'cube' as const,
    label: t('dashboard.todayTokens'),
    value: formatTokens(props.stats?.today_tokens || 0),
    detail: `${t('dashboard.input')}: ${formatTokens(props.stats?.today_input_tokens || 0)}`
  },
  {
    key: 'performance',
    icon: 'bolt' as const,
    label: t('dashboard.performance'),
    value: `${formatTokens(props.stats?.rpm || 0)} RPM`,
    detail: `${formatTokens(props.stats?.tpm || 0)} TPM`
  },
  {
    key: 'avgResponse',
    icon: 'clock' as const,
    label: t('dashboard.avgResponse'),
    value: formatDuration(props.stats?.average_duration_ms || 0),
    detail: t('dashboard.averageTime')
  }
])

const formatBalance = (b: number) =>
  new Intl.NumberFormat('en-US', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  }).format(b)

const formatNumber = (n: number) => n.toLocaleString()
const formatCost = (c: number) => c.toFixed(4)
const formatTokens = (value: number) => {
  if (value >= 1_000_000_000) return `${(value / 1_000_000_000).toFixed(2)}B`
  if (value >= 1_000_000) return `${(value / 1_000_000).toFixed(1)}M`
  if (value >= 1000) return `${(value / 1000).toFixed(1)}K`
  return value.toString()
}
const formatDuration = (ms: number) => ms >= 1000 ? `${(ms / 1000).toFixed(2)}s` : `${ms.toFixed(0)}ms`
</script>

<style scoped>
.dashboard-overview {
  display: grid;
  gap: 1.25rem;
}

.dashboard-hero {
  display: grid;
  grid-template-columns: minmax(16rem, 1fr) minmax(22rem, 0.78fr);
  align-items: center;
  gap: clamp(0.9rem, 1.7vw, 1.25rem);
  overflow: hidden;
  border: 1px solid rgba(2, 43, 18, 0.1);
  border-radius: 1.75rem;
  background:
    radial-gradient(circle at 8% 12%, rgba(0, 135, 60, 0.12), transparent 34%),
    linear-gradient(135deg, rgba(255, 255, 252, 0.92), rgba(241, 248, 232, 0.84));
  box-shadow: 0 24px 70px rgba(2, 43, 18, 0.09);
  padding: clamp(1rem, 1.55vw, 1.25rem);
}

:global(.dark) .dashboard-hero {
  border-color: rgba(0, 135, 60, 0.22);
  background:
    radial-gradient(circle at 8% 12%, rgba(0, 135, 60, 0.14), transparent 34%),
    radial-gradient(circle at 88% 12%, rgba(217, 162, 27, 0.1), transparent 30%),
    linear-gradient(135deg, rgba(242, 251, 246, 0.94), rgba(224, 240, 232, 0.9));
  box-shadow: 0 24px 70px rgba(0, 0, 0, 0.16);
}

.dashboard-eyebrow {
  font-size: 0.72rem;
  font-weight: 800;
  letter-spacing: 0.18em;
  text-transform: uppercase;
  color: rgba(0, 107, 49, 0.76);
}

:global(.dark) .dashboard-eyebrow {
  color: rgba(0, 107, 49, 0.78);
}

.dashboard-hero-title {
  color: #022b12;
  max-width: 13em;
  line-height: 1.08;
}

:global(.dark) .dashboard-hero-title {
  color: #043b1d;
}

.dashboard-hero-panel-kicker {
  color: rgba(2, 43, 18, 0.48);
}

:global(.dark) .dashboard-hero-panel-kicker {
  color: rgba(2, 43, 18, 0.5);
}

.dashboard-hero-panel-status {
  color: #022b12;
}

:global(.dark) .dashboard-hero-panel-status {
  color: #043b1d;
}

.dashboard-hero-panel,
.dashboard-platforms,
.dashboard-metric {
  border: 1px solid rgba(2, 43, 18, 0.08);
  background: rgba(255, 255, 255, 0.68);
  box-shadow: 0 14px 36px rgba(2, 43, 18, 0.06);
  backdrop-filter: blur(18px);
}

:global(.dark) .dashboard-hero-panel,
:global(.dark) .dashboard-platforms,
:global(.dark) .dashboard-metric {
  border-color: rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.055);
  box-shadow: 0 14px 36px rgba(0, 0, 0, 0.16);
}

.dashboard-hero-panel {
  display: flex;
  width: min(100%, 34rem);
  min-height: 9.5rem;
  flex-direction: column;
  justify-content: space-between;
  justify-self: end;
  border-radius: 1.35rem;
  padding: clamp(0.82rem, 1.15vw, 0.95rem);
}

.dashboard-hero-copy {
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  padding: 0.25rem 0;
}

.dashboard-hero-stats-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0.62rem;
}

:global(.dark) .dashboard-hero-panel {
  border-color: rgba(0, 135, 60, 0.14);
  background: rgba(255, 255, 255, 0.5);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.45), 0 16px 38px rgba(2, 43, 18, 0.08);
}

.dashboard-status-dot {
  height: 0.72rem;
  width: 0.72rem;
  border-radius: 999px;
  background: #d9a21b;
  box-shadow: 0 0 0 5px rgba(217, 162, 27, 0.12);
}

.dashboard-status-dot.is-ready {
  background: #00873c;
  box-shadow: 0 0 0 5px rgba(0, 135, 60, 0.12);
}

.dashboard-hero-stat {
  display: flex;
  min-height: 3.45rem;
  flex-direction: column;
  justify-content: center;
  border-radius: 0.9rem;
  background: rgba(2, 43, 18, 0.045);
  padding: 0.62rem 0.78rem;
}

:global(.dark) .dashboard-hero-stat {
  background: rgba(0, 107, 49, 0.06);
}

.dashboard-hero-stat span,
.dashboard-metric-label {
  display: block;
  font-size: 0.75rem;
  font-weight: 650;
  color: rgba(2, 43, 18, 0.55);
}

:global(.dark) .dashboard-hero-stat span,
:global(.dark) .dashboard-metric-label {
  color: rgba(255, 255, 255, 0.54);
}

:global(.dark) .dashboard-hero-stat span {
  color: rgba(2, 43, 18, 0.58);
}

.dashboard-hero-stat strong {
  display: block;
  margin-top: 0.25rem;
  color: #022b12;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  font-size: clamp(0.98rem, 1.25vw, 1.12rem);
  font-weight: 850;
}

:global(.dark) .dashboard-hero-stat strong {
  color: rgba(255, 255, 255, 0.92);
}

:global(.dark) .dashboard-hero-stat strong {
  color: #022b12;
}

.dashboard-metric-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 1rem;
}

.dashboard-metric {
  display: flex;
  gap: 0.9rem;
  min-width: 0;
  border-radius: 1.35rem;
  padding: 1rem;
}

.dashboard-metric-icon {
  display: flex;
  width: 2.6rem;
  height: 2.6rem;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  border-radius: 1rem;
  background: rgba(0, 135, 60, 0.1);
  color: #006b31;
}

:global(.dark) .dashboard-metric-icon {
  background: rgba(0, 135, 60, 0.16);
  color: #bde8cb;
}

.dashboard-metric-value {
  margin-top: 0.3rem;
  color: #022b12;
  font-size: 1.18rem;
  font-weight: 850;
  letter-spacing: -0.02em;
}

:global(.dark) .dashboard-metric-value {
  color: rgba(255, 255, 255, 0.92);
}

.dashboard-metric-detail {
  margin-top: 0.18rem;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: rgba(2, 43, 18, 0.48);
  font-size: 0.75rem;
}

:global(.dark) .dashboard-metric-detail {
  color: rgba(255, 255, 255, 0.42);
}

@media (max-width: 1024px) {
  .dashboard-hero {
    grid-template-columns: 1fr;
  }

  .dashboard-metric-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (min-width: 641px) and (max-width: 1024px) {
  .dashboard-hero {
    gap: 0.85rem;
  }

  .dashboard-hero-panel {
    min-height: 9.7rem;
  }

  .dashboard-hero-stats-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .dashboard-hero-stat {
    min-height: 3.75rem;
    padding-inline: 0.8rem;
  }
}

@media (max-width: 640px) {
  .dashboard-hero {
    border-radius: 1.5rem;
    padding: 1rem;
  }

  .dashboard-metric-grid {
    grid-template-columns: 1fr;
  }
}
</style>
