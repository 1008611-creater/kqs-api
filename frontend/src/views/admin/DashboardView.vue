<template>
  <AppLayout>
    <div class="ops-overview">
      <header class="ops-overview__head">
        <div>
          <p class="ops-overview__eyebrow">运营控制台</p>
          <h1 class="ops-overview__title">经营概览</h1>
        </div>
        <button class="ops-overview__refresh" :disabled="loading" @click="refreshDashboard">刷新数据</button>
      </header>

      <div v-if="loading && !stats" class="ops-overview__loading">
        <LoadingSpinner />
      </div>

      <template v-else-if="stats">
        <section class="ops-kpis" aria-label="核心运营指标">
          <button class="ops-kpi ops-kpi--supply" @click="goTo('/admin/accounts')">
            <p class="ops-kpi__label">渠道供给</p>
            <p class="ops-kpi__value">{{ stats.normal_accounts }} <span>/ {{ stats.total_accounts }}</span></p>
            <p class="ops-kpi__meta" :class="stats.error_accounts > 0 ? 'ops-kpi__meta--attention' : 'ops-kpi__meta--healthy'">
              {{ stats.error_accounts > 0 ? `${stats.error_accounts} 个渠道需要处理` : '当前无异常渠道' }}
            </p>
          </button>

          <button class="ops-kpi" @click="goTo('/admin/ops')">
            <p class="ops-kpi__label">实时服务</p>
            <p class="ops-kpi__value">{{ formatNumber(stats.today_requests) }}</p>
            <p class="ops-kpi__meta">今日请求 · {{ stats.active_users }} 位活跃用户</p>
          </button>

          <button class="ops-kpi" @click="goTo('/admin/subscriptions')">
            <p class="ops-kpi__label">用户与权益</p>
            <p class="ops-kpi__value">{{ formatNumber(stats.total_users) }}</p>
            <p class="ops-kpi__meta">总用户 · 订阅、补发与重置</p>
          </button>

          <button class="ops-kpi" @click="goTo('/admin/usage')">
            <p class="ops-kpi__label">今日实际扣除</p>
            <p class="ops-kpi__value">${{ formatCost(stats.today_actual_cost) }}</p>
            <p class="ops-kpi__meta">{{ formatTokens(stats.today_tokens) }} Token · 仅今日口径</p>
          </button>
        </section>

        <section class="ops-traffic">
          <div class="ops-traffic__head">
            <div>
              <p class="ops-traffic__eyebrow">近 24 小时</p>
              <h2 class="ops-traffic__title">上游代理流量</h2>
            </div>
            <button class="ops-traffic__action" @click="goTo('/admin/ops')">进入运维处置</button>
          </div>

          <div v-if="upstreamTrafficLoading" class="ops-traffic__loading"><LoadingSpinner /></div>
          <template v-else-if="upstreamTraffic">
            <div class="ops-traffic__stats">
              <div class="ops-traffic__stat"><p>总流量</p><strong>{{ formatBytes(upstreamTraffic.total_bytes) }}</strong></div>
              <div class="ops-traffic__stat"><p>上行请求</p><strong>{{ formatBytes(upstreamTraffic.request_bytes) }}</strong></div>
              <div class="ops-traffic__stat"><p>下行响应</p><strong>{{ formatBytes(upstreamTraffic.response_bytes) }}</strong></div>
              <div class="ops-traffic__stat"><p>转发请求</p><strong>{{ formatNumber(upstreamTraffic.request_count) }}</strong></div>
            </div>
            <div class="ops-traffic__table-wrap">
              <table class="ops-traffic__table">
                <thead><tr><th>主要消耗渠道</th><th>流量</th><th>请求</th><th>平均耗时</th></tr></thead>
                <tbody>
                  <tr v-for="item in upstreamTraffic.top_accounts.slice(0, 5)" :key="item.account_id" @click="goTo('/admin/accounts')">
                    <td>{{ item.account_name }}</td>
                    <td>{{ formatBytes(item.total_bytes) }}</td>
                    <td>{{ formatNumber(item.request_count) }}</td>
                    <td>{{ formatDuration(item.avg_duration_ms) }}</td>
                  </tr>
                  <tr v-if="!upstreamTraffic.top_accounts.length"><td colspan="4" class="ops-traffic__empty">近 24 小时暂无上游流量</td></tr>
                </tbody>
              </table>
            </div>
          </template>
          <div v-else class="ops-traffic__empty">暂时无法读取上游流量，服务本身不受影响。</div>
        </section>

        <section class="ops-actions" aria-label="快捷操作">
          <button class="ops-action" @click="goTo('/admin/accounts')"><span><strong>渠道供给</strong><small>测速、摘除故障渠道、调整优先级和权重。</small></span><b aria-hidden="true">→</b></button>
          <button class="ops-action" @click="goTo('/admin/groups')"><span><strong>分组与模型</strong><small>配置用户可用模型、倍率和分组供给。</small></span><b aria-hidden="true">→</b></button>
          <button class="ops-action" @click="goTo('/admin/subscriptions')"><span><strong>订阅与兑换</strong><small>处理订阅发放、额度重置、补偿和兑换。</small></span><b aria-hidden="true">→</b></button>
        </section>
      </template>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { adminAPI } from '@/api/admin'
import { opsAPI, type OpsUpstreamTrafficStatsResponse } from '@/api/admin/ops'
import type { DashboardStats } from '@/types'
import AppLayout from '@/components/layout/AppLayout.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'

const router = useRouter()
const stats = ref<DashboardStats | null>(null)
const upstreamTraffic = ref<OpsUpstreamTrafficStatsResponse | null>(null)
const loading = ref(false)
const upstreamTrafficLoading = ref(false)

const formatNumber = (value: number | undefined) => Number(value || 0).toLocaleString()
const formatTokens = (value: number | undefined) => {
  const numeric = Number(value || 0)
  if (numeric >= 1_000_000_000) return `${(numeric / 1_000_000_000).toFixed(2)}B`
  if (numeric >= 1_000_000) return `${(numeric / 1_000_000).toFixed(2)}M`
  if (numeric >= 1_000) return `${(numeric / 1_000).toFixed(2)}K`
  return formatNumber(numeric)
}
const formatCost = (value: number | undefined) => {
  const numeric = Number(value || 0)
  return numeric >= 1000 ? `${(numeric / 1000).toFixed(2)}K` : numeric.toFixed(numeric >= 1 ? 2 : 3)
}
const formatBytes = (value: number | undefined) => {
  let bytes = Number(value || 0)
  if (bytes < 1024) return `${Math.round(bytes)} B`
  for (const unit of ['KB', 'MB', 'GB', 'TB']) { bytes /= 1024; if (bytes < 1024) return `${bytes.toFixed(bytes >= 100 ? 0 : 1)} ${unit}` }
  return `${bytes.toFixed(1)} PB`
}
const formatDuration = (value: number | null | undefined) => {
  const numeric = Number(value || 0)
  return numeric >= 1000 ? `${(numeric / 1000).toFixed(2)}s` : `${Math.round(numeric)}ms`
}
const goTo = (path: string) => { void router.push(path) }

const loadDashboard = async () => {
  loading.value = true
  upstreamTrafficLoading.value = true
  const now = new Date()
  const start = new Date(now.getTime() - 24 * 60 * 60 * 1000)
  try {
    const [nextStats, nextTraffic] = await Promise.all([
      adminAPI.dashboard.getStats(),
      opsAPI.getUpstreamTrafficStats({ start_time: start.toISOString(), end_time: now.toISOString() })
    ])
    stats.value = nextStats
    upstreamTraffic.value = nextTraffic
  } catch (error) {
    console.error('Failed to load operations overview:', error)
    if (!stats.value) window.setTimeout(() => undefined, 0)
  } finally {
    loading.value = false
    upstreamTrafficLoading.value = false
  }
}

const refreshDashboard = async () => {
  await loadDashboard()
}

onMounted(() => { void loadDashboard() })
</script>

<style scoped>
.ops-overview { display: grid; gap: 20px; max-width: 1440px; margin: 0 auto; color: #17201f; }
.ops-overview__head { display: flex; align-items: flex-end; justify-content: space-between; gap: 16px; padding: 4px 0 2px; }
.ops-overview__eyebrow, .ops-traffic__eyebrow { margin: 0 0 5px; color: #5d6b68; font-size: 11px; font-weight: 700; letter-spacing: 0; text-transform: uppercase; }
.ops-overview__title { margin: 0; font-size: 24px; font-weight: 650; line-height: 1.15; letter-spacing: 0; }
.ops-overview__refresh, .ops-traffic__action { min-height: 34px; border: 1px solid #d7e0dd; border-radius: 8px; background: #fff; color: #30413d; font-size: 13px; font-weight: 600; padding: 0 12px; transition: border-color .16s ease, background-color .16s ease; }
.ops-overview__refresh:hover, .ops-traffic__action:hover { border-color: #93b8af; background: #f5faf8; }
.ops-overview__refresh:focus-visible, .ops-traffic__action:focus-visible, .ops-kpi:focus-visible, .ops-action:focus-visible { outline: 2px solid #147d6c; outline-offset: 2px; }
.ops-overview__loading, .ops-traffic__loading { display: flex; align-items: center; justify-content: center; min-height: 260px; }
.ops-kpis { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); border: 1px solid #dce5e2; border-radius: 8px; background: #fff; overflow: hidden; }
.ops-kpi { min-width: 0; padding: 18px 20px; border-right: 1px solid #e5ecea; background: transparent; color: inherit; text-align: left; transition: background-color .16s ease; }
.ops-kpi:last-child { border-right: 0; }.ops-kpi:hover { background: #f5faf8; }.ops-kpi--supply { box-shadow: inset 3px 0 0 #147d6c; }
.ops-kpi__label { margin: 0; color: #687875; font-size: 12px; font-weight: 600; }.ops-kpi__value { margin: 9px 0 6px; color: #17201f; font-size: 25px; font-weight: 650; line-height: 1; }.ops-kpi__value span { color: #91a09d; font-size: 15px; font-weight: 600; }.ops-kpi__meta { margin: 0; color: #687875; font-size: 12px; line-height: 1.4; }.ops-kpi__meta--healthy { color: #147d6c; }.ops-kpi__meta--attention { color: #b45309; }
.ops-traffic { border: 1px solid #dce5e2; border-radius: 8px; background: #fff; overflow: hidden; }.ops-traffic__head { display: flex; justify-content: space-between; align-items: flex-start; gap: 18px; padding: 18px 20px; border-bottom: 1px solid #e5ecea; }.ops-traffic__title { margin: 0; font-size: 16px; font-weight: 650; line-height: 1.3; }.ops-traffic__stats { display: grid; grid-template-columns: repeat(4, minmax(0, 1fr)); border-bottom: 1px solid #e5ecea; }.ops-traffic__stat { min-width: 0; padding: 16px 20px; border-right: 1px solid #e5ecea; }.ops-traffic__stat:last-child { border-right: 0; }.ops-traffic__stat p { margin: 0; color: #687875; font-size: 12px; }.ops-traffic__stat strong { display: block; overflow: hidden; margin-top: 7px; color: #17201f; font-size: 19px; font-weight: 650; text-overflow: ellipsis; white-space: nowrap; }.ops-traffic__table-wrap { overflow-x: auto; }.ops-traffic__table { width: 100%; min-width: 620px; border-collapse: collapse; font-size: 13px; }.ops-traffic__table th { padding: 11px 20px; border-bottom: 1px solid #e5ecea; color: #778582; font-size: 11px; font-weight: 650; text-align: right; text-transform: uppercase; }.ops-traffic__table th:first-child, .ops-traffic__table td:first-child { text-align: left; }.ops-traffic__table td { padding: 13px 20px; border-bottom: 1px solid #edf1ef; color: #55635f; text-align: right; }.ops-traffic__table td:first-child { color: #24322f; font-weight: 600; }.ops-traffic__table tbody tr:not(:only-child) { cursor: pointer; }.ops-traffic__table tbody tr:not(:only-child):hover { background: #f5faf8; }.ops-traffic__table tbody tr:last-child td { border-bottom: 0; }.ops-traffic__empty { padding: 36px 20px !important; color: #778582 !important; font-weight: 400 !important; text-align: center !important; }
.ops-actions { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); border-top: 1px solid #dce5e2; }.ops-action { display: flex; min-width: 0; align-items: center; justify-content: space-between; gap: 16px; padding: 16px 4px 16px 0; color: #22312e; text-align: left; }.ops-action + .ops-action { padding-left: 20px; border-left: 1px solid #dce5e2; }.ops-action span { min-width: 0; }.ops-action strong, .ops-action small { display: block; }.ops-action strong { font-size: 13px; font-weight: 650; }.ops-action small { overflow: hidden; margin-top: 5px; color: #687875; font-size: 12px; line-height: 1.4; text-overflow: ellipsis; white-space: nowrap; }.ops-action b { color: #147d6c; font-size: 18px; font-weight: 500; transition: transform .16s ease; }.ops-action:hover b { transform: translateX(3px); }
:global(.dark) .ops-overview { color: #e3ece9; }:global(.dark) .ops-overview__title, :global(.dark) .ops-kpi__value, :global(.dark) .ops-traffic__title, :global(.dark) .ops-traffic__stat strong, :global(.dark) .ops-traffic__table td:first-child, :global(.dark) .ops-action { color: #edf5f2; }:global(.dark) .ops-kpi__label, :global(.dark) .ops-kpi__meta, :global(.dark) .ops-traffic__stat p, :global(.dark) .ops-traffic__table td, :global(.dark) .ops-action small { color: #a7b7b2; }:global(.dark) .ops-overview__refresh, :global(.dark) .ops-traffic__action, :global(.dark) .ops-kpis, :global(.dark) .ops-traffic { border-color: #31423e; background: #17211f; }:global(.dark) .ops-overview__refresh, :global(.dark) .ops-traffic__action { color: #d9e7e2; }:global(.dark) .ops-overview__refresh:hover, :global(.dark) .ops-traffic__action:hover, :global(.dark) .ops-kpi:hover, :global(.dark) .ops-traffic__table tbody tr:not(:only-child):hover { background: #20302c; }:global(.dark) .ops-kpi, :global(.dark) .ops-traffic__stat, :global(.dark) .ops-action + .ops-action { border-color: #31423e; }:global(.dark) .ops-traffic__head, :global(.dark) .ops-traffic__stats, :global(.dark) .ops-traffic__table th, :global(.dark) .ops-traffic__table td { border-color: #2a3a36; }
@media (max-width: 1023px) { .ops-kpis { grid-template-columns: repeat(2, minmax(0, 1fr)); }.ops-kpi:nth-child(2) { border-right: 0; }.ops-kpi:nth-child(-n+2) { border-bottom: 1px solid #e5ecea; }.ops-traffic__stats { grid-template-columns: repeat(2, minmax(0, 1fr)); }.ops-traffic__stat:nth-child(2) { border-right: 0; }.ops-traffic__stat:nth-child(-n+2) { border-bottom: 1px solid #e5ecea; }.ops-actions { grid-template-columns: 1fr; }.ops-action, .ops-action + .ops-action { padding: 14px 0; border-left: 0; }.ops-action + .ops-action { border-top: 1px solid #dce5e2; } }
@media (max-width: 640px) { .ops-overview__head, .ops-traffic__head { align-items: stretch; flex-direction: column; }.ops-overview__refresh, .ops-traffic__action { width: 100%; }.ops-kpis, .ops-traffic__stats { grid-template-columns: 1fr; }.ops-kpi, .ops-kpi:nth-child(2), .ops-traffic__stat, .ops-traffic__stat:nth-child(2) { border-right: 0; border-bottom: 1px solid #e5ecea; }.ops-kpi:last-child, .ops-traffic__stat:last-child { border-bottom: 0; }.ops-kpi:nth-child(-n+2), .ops-traffic__stat:nth-child(-n+2) { border-bottom: 1px solid #e5ecea; }.ops-overview__title { font-size: 22px; } }
@media (prefers-reduced-motion: reduce) { .ops-overview *, .ops-overview *::before, .ops-overview *::after { scroll-behavior: auto !important; transition-duration: .01ms !important; } }
</style>
