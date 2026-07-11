<template>
  <AppLayout>
    <div class="kqs-dashboard">
      <section class="kqs-grid" aria-label="概况核心信息">
        <div class="kqs-panel kqs-usage">
          <div class="kqs-panel__header">
            <h2>使用量统计</h2>
            <button type="button" class="kqs-icon-button" title="刷新" :disabled="refreshing" @click="refreshAll">
              <Icon name="refresh" size="sm" />
            </button>
          </div>

          <div class="kqs-stat-grid">
            <article v-for="item in statsCards" :key="item.label" class="kqs-stat">
              <span>{{ item.label }}</span>
              <strong>{{ item.value }}</strong>
            </article>
          </div>

          <div class="kqs-chart" aria-label="近7日趋势图">
            <div class="kqs-chart__bars">
              <div
                v-for="bar in trendBars"
                :key="bar.date"
                class="kqs-chart__bar-item"
                :title="`${bar.date}: ${bar.requests} 次请求`"
              >
                <span class="kqs-chart__bar-value">{{ bar.requests }}</span>
                <span
                  class="kqs-chart__bar"
                  :class="{ 'kqs-chart__bar--empty': bar.requests === 0 }"
                  :style="{ height: `${bar.height}%` }"
                />
                <time>{{ bar.label }}</time>
              </div>
            </div>
          </div>
        </div>

        <div class="kqs-panel kqs-api">
          <div class="kqs-panel__header">
            <h2>我的 API</h2>
            <router-link to="/keys" class="kqs-text-link">管理全部</router-link>
          </div>

          <div class="kqs-api-row">
            <code>{{ apiBaseUrl }}</code>
            <button type="button" class="kqs-small-button" @click="copyText(apiBaseUrl, '接口地址')">
              <Icon name="copy" size="sm" />
              复制
            </button>
          </div>

          <div class="kqs-api-row">
            <code>{{ primaryKeyText }}</code>
            <button
              type="button"
              class="kqs-small-button"
              :disabled="!primaryKey"
              @click="showFullKey = !showFullKey"
            >
              <Icon name="eye" size="sm" />
              显示
            </button>
          </div>

          <div class="kqs-api-actions">
            <button type="button" class="kqs-block-button" :disabled="creatingKey" @click="createNewKey">
              <Icon name="plus" size="sm" />
              创建新 Key
            </button>
            <router-link to="/guide" class="kqs-block-button">
              <Icon name="document" size="sm" />
              对接文档
            </router-link>
            <router-link to="/keys" class="kqs-block-button">
              <Icon name="chevronRight" size="sm" />
              查看更多
            </router-link>
          </div>
        </div>
      </section>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import { useAppStore, useAuthStore } from '@/stores'
import { keysAPI } from '@/api/keys'
import { usageAPI, type UserDashboardStats } from '@/api/usage'
import { userGroupsAPI } from '@/api/groups'
import type { ApiKey, Group, TrendDataPoint } from '@/types'

const router = useRouter()
const appStore = useAppStore()
const authStore = useAuthStore()

const refreshing = ref(false)
const creatingKey = ref(false)
const stats = ref<UserDashboardStats | null>(null)
const apiKeys = ref<ApiKey[]>([])
const groups = ref<Group[]>([])
const trendData = ref<TrendDataPoint[]>([])
const showFullKey = ref(false)

const today = new Date()
const sevenDaysAgo = new Date(Date.now() - 6 * 24 * 60 * 60 * 1000)

function toDateInput(date: Date): string {
  return date.toISOString().slice(0, 10)
}

const apiBaseUrl = computed(() => {
  const configured = appStore.apiBaseUrl?.trim()
  if (configured) return configured
  return `${window.location.origin}/v1`
})

const primaryKey = computed(() => apiKeys.value.find((item) => item.status === 'active') ?? apiKeys.value[0] ?? null)
const primaryKeyText = computed(() => {
  if (!primaryKey.value) return '还没有 API Key'
  if (showFullKey.value) return primaryKey.value.key
  return maskKey(primaryKey.value.key)
})

const statsCards = computed(() => [
  { label: '今日请求', value: formatInteger(stats.value?.today_requests ?? 0) },
  { label: '今日消耗', value: `¥${formatMoney(stats.value?.today_actual_cost ?? 0)}` },
  { label: '总Token数', value: formatCompact(stats.value?.total_tokens ?? 0) },
])

const trendBars = computed(() => {
  const requestsByDate = new Map(
    trendData.value.map((point) => [point.date.slice(0, 10), point.requests]),
  )
  const days = Array.from({ length: 7 }, (_, index) => {
    const date = new Date(sevenDaysAgo)
    date.setDate(sevenDaysAgo.getDate() + index)
    const key = toDateInput(date)
    return {
      date: key,
      label: `${date.getMonth() + 1}/${date.getDate()}`,
      requests: requestsByDate.get(key) ?? 0,
    }
  })
  const maxRequests = Math.max(...days.map((point) => point.requests), 0)

  return days.map((point) => ({
    ...point,
    height: point.requests > 0 && maxRequests > 0
      ? Math.max(18, Math.round((point.requests / maxRequests) * 86))
      : 10,
  }))
})

function formatInteger(value: number): string {
  return new Intl.NumberFormat('en-US', { maximumFractionDigits: 0 }).format(value)
}

function formatMoney(value: number): string {
  return value.toFixed(2)
}

function formatCompact(value: number): string {
  if (value >= 1_000_000) return `${(value / 1_000_000).toFixed(value >= 10_000_000 ? 0 : 1)}M`
  if (value >= 1_000) return `${(value / 1_000).toFixed(value >= 10_000 ? 0 : 1)}K`
  return formatInteger(value)
}

function maskKey(key: string): string {
  if (!key) return ''
  if (key.length <= 12) return `${key.slice(0, 3)}••••${key.slice(-3)}`
  return `${key.slice(0, 3)}••••••••••••••••${key.slice(-4)}`
}

async function copyText(text: string, label: string) {
  try {
    await navigator.clipboard.writeText(text)
    appStore.showSuccess(`${label}已复制`)
  } catch {
    appStore.showError('复制失败，请手动复制')
  }
}

async function loadStats() {
  await authStore.refreshUser().catch(() => undefined)
  stats.value = await usageAPI.getDashboardStats()
}

async function loadTrend() {
  const response = await usageAPI.getDashboardTrend({
    start_date: toDateInput(sevenDaysAgo),
    end_date: toDateInput(today),
    granularity: 'day',
  })
  trendData.value = response.trend ?? []
}

async function loadKeys() {
  const response = await keysAPI.list(1, 8, { sort_by: 'created_at', sort_order: 'desc' })
  apiKeys.value = response.items ?? []
}

async function loadGroups() {
  groups.value = await userGroupsAPI.getAvailable()
}

async function createNewKey() {
  if (creatingKey.value) return
  const groupId = groups.value.find((group) => group.status === 'active')?.id
  if (!groupId) {
    appStore.showError('当前没有可用分组，请先到 API Key 页面选择分组')
    await router.push('/keys')
    return
  }

  creatingKey.value = true
  try {
    await keysAPI.create('矿泉水API', groupId)
    appStore.showSuccess('API Key 已创建')
    await loadKeys()
  } catch (error) {
    appStore.showError((error as { message?: string })?.message || '创建失败')
  } finally {
    creatingKey.value = false
  }
}

async function refreshAll() {
  if (refreshing.value) return
  refreshing.value = true
  const results = await Promise.allSettled([loadStats(), loadTrend(), loadKeys(), loadGroups()])
  if (results.some((result) => result.status === 'rejected')) {
    console.warn('Dashboard partially failed to load:', results)
  }
  refreshing.value = false
}

onMounted(() => {
  refreshAll()
})
</script>

<style scoped>
.kqs-dashboard {
  --panel: rgba(5, 29, 24, 0.78);
  --soft: rgba(3, 22, 18, 0.72);
  --line: rgba(43, 132, 83, 0.38);
  --text: #edf7ee;
  --muted: rgba(220, 239, 224, 0.64);
  --accent: #22c55e;
  display: grid;
  grid-template-columns: 1fr;
  gap: 0.85rem;
  min-height: calc(100dvh - 7rem);
  align-content: start;
}

.kqs-loading,
.kqs-panel {
  border: 1px solid var(--line);
  border-radius: 8px;
  background: var(--panel);
}

.kqs-loading {
  display: flex;
  min-height: 28rem;
  align-items: center;
  justify-content: center;
}

.kqs-grid {
  display: grid;
  grid-column: 1 / -1;
  grid-template-columns: 1fr;
  grid-template-rows: minmax(0, 1.18fr) minmax(0, 0.82fr);
  gap: 0.82rem;
  align-items: stretch;
  min-height: calc(100dvh - 7rem);
}

.kqs-panel {
  padding: 0.92rem;
  box-shadow: 0 18px 54px rgba(0, 0, 0, 0.22);
}

.kqs-usage,
.kqs-api {
  min-width: 0;
  min-height: 0;
  height: 100%;
}

.kqs-usage {
  display: flex;
  flex-direction: column;
}

.kqs-api {
  display: flex;
  flex-direction: column;
}

.kqs-panel__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.75rem;
  margin-bottom: 0.68rem;
}

.kqs-panel__header h2 {
  margin: 0;
  color: var(--text);
  font-size: 1rem;
  font-weight: 750;
}

.kqs-icon-button,
.kqs-small-button,
.kqs-block-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.35rem;
  border: 1px solid var(--line);
  border-radius: 8px;
  background: rgba(3, 22, 18, 0.64);
  color: var(--text);
  font-size: 0.84rem;
  font-weight: 680;
  text-decoration: none;
  transition: background-color 0.16s ease, border-color 0.16s ease, transform 0.16s ease;
}

.kqs-icon-button {
  width: 2rem;
  height: 2rem;
}

.kqs-small-button {
  min-height: 2rem;
  padding: 0 0.72rem;
}

.kqs-block-button {
  min-height: 2.25rem;
  padding: 0 0.9rem;
}

.kqs-icon-button:hover,
.kqs-small-button:hover,
.kqs-block-button:hover {
  border-color: rgba(34, 197, 94, 0.52);
  background: rgba(34, 197, 94, 0.12);
  transform: translateY(-1px);
}

.kqs-icon-button:disabled,
.kqs-small-button:disabled,
.kqs-block-button:disabled {
  cursor: not-allowed;
  opacity: 0.58;
  transform: none;
}

.kqs-stat-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 0.62rem;
}

.kqs-stat {
  min-height: 4rem;
  border-radius: 8px;
  background: var(--soft);
  padding: 0.74rem 0.86rem;
}

.kqs-stat span {
  display: block;
  color: var(--muted);
  font-size: 0.78rem;
  font-weight: 650;
}

.kqs-stat strong {
  display: block;
  margin-top: 0.3rem;
  color: var(--text);
  font-size: 1.28rem;
  font-weight: 760;
}

.kqs-chart {
  display: flex;
  flex: 1;
  min-height: 12rem;
  align-items: flex-end;
  justify-content: center;
  margin-top: 0.62rem;
  border-radius: 8px;
  background:
    linear-gradient(to top, rgba(43, 132, 83, 0.22) 1px, transparent 1px),
    var(--soft);
  background-size: 100% 25%;
  padding: 0.72rem 0.88rem 0.6rem;
}

.kqs-chart__empty {
  display: flex;
  align-items: center;
  gap: 0.45rem;
  color: var(--muted);
  font-size: 0.86rem;
  font-weight: 650;
}

.kqs-chart__bars {
  display: grid;
  grid-template-columns: repeat(7, minmax(0, 1fr));
  width: 100%;
  height: 100%;
  min-height: 6.65rem;
  align-items: end;
  gap: 0.54rem;
}

.kqs-chart__bar-item {
  display: grid;
  height: 100%;
  min-width: 0;
  grid-template-rows: 1rem 1fr 1rem;
  align-items: end;
  justify-items: center;
  gap: 0.28rem;
}

.kqs-chart__bar-value,
.kqs-chart__bar-item time {
  color: var(--muted);
  font-size: 0.68rem;
  font-weight: 700;
  line-height: 1;
}

.kqs-chart__bar {
  width: min(100%, 3.2rem);
  min-height: 0.7rem;
  border-radius: 6px 6px 2px 2px;
  background: linear-gradient(180deg, #22c55e, #0f7a3a);
  box-shadow: inset 0 1px 0 rgba(34, 197, 94, 0.18);
}

.kqs-chart__bar--empty {
  background: rgba(43, 132, 83, 0.2);
}

.kqs-api-row {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 0.55rem;
  margin-bottom: 0.65rem;
}

.kqs-api-row code {
  display: flex;
  min-width: 0;
  min-height: 2.25rem;
  align-items: center;
  overflow: hidden;
  border: 1px solid var(--line);
  border-radius: 7px;
  background: rgba(3, 22, 18, 0.64);
  padding: 0 0.78rem;
  color: #edf7ee;
  font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
  font-size: 0.78rem;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.kqs-api-actions {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 0.5rem;
  margin-top: auto;
}

.kqs-text-link {
  color: var(--accent);
  font-size: 0.8rem;
  font-weight: 700;
  text-decoration: none;
}

.kqs-text-link:hover {
  text-decoration: underline;
}

@media (max-width: 920px) {
  .kqs-dashboard {
    grid-template-columns: 1fr;
  }

  .kqs-grid {
    grid-template-columns: 1fr;
    grid-template-rows: auto auto;
    min-height: auto;
  }

  .kqs-chart {
    min-height: 10rem;
  }
}

@media (max-width: 640px) {
  .kqs-panel {
    padding: 0.85rem;
  }

  .kqs-stat-grid,
  .kqs-api-actions {
    grid-template-columns: 1fr;
  }

  .kqs-api-row {
    grid-template-columns: 1fr;
  }

}
</style>
