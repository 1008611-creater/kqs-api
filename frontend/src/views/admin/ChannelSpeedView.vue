<template>
  <AppLayout>
    <TablePageLayout>
      <template #filters>
        <div class="space-y-4">
          <div class="grid grid-cols-2 gap-3 lg:grid-cols-5">
            <div class="speed-stat">
              <span class="speed-stat-label">测试账号</span>
              <strong>{{ summary?.total ?? 0 }}</strong>
            </div>
            <div class="speed-stat">
              <span class="speed-stat-label">成功率</span>
              <strong>{{ formatPercent(summary?.success_rate ?? 0) }}</strong>
            </div>
            <div class="speed-stat">
              <span class="speed-stat-label">平均延迟</span>
              <strong>{{ formatMs(summary?.avg_latency_ms ?? 0) }}</strong>
            </div>
            <div class="speed-stat">
              <span class="speed-stat-label">P95 延迟</span>
              <strong>{{ formatMs(summary?.p95_latency_ms ?? 0) }}</strong>
            </div>
            <div class="speed-stat col-span-2 lg:col-span-1">
              <span class="speed-stat-label">最快渠道</span>
              <strong class="truncate text-sm">{{ fastestDisplay }}</strong>
            </div>
          </div>

          <div
            v-if="runMessage || loading"
            class="speed-feedback"
            :class="runStatusClass"
          >
            <div>
              <div class="font-semibold">
                {{ loading ? '测速任务执行中' : runMessage }}
              </div>
              <div class="mt-1 text-xs opacity-80">
                {{ loading ? loadingHint : lastRunAt ? `完成时间：${lastRunAt}` : '等待测速' }}
              </div>
            </div>
            <span class="speed-feedback-pill">{{ targetModeLabel }}</span>
          </div>

          <div class="speed-model-banner">
            <div>
              <div class="text-xs font-semibold uppercase tracking-wide text-blue-600 dark:text-blue-300">固定检测模型</div>
              <div class="mt-1 flex flex-wrap items-center gap-2">
                <span class="speed-model-name">{{ SPEED_TEST_MODEL }}</span>
                <span class="speed-model-note">只检测启用、可调度账号；停用分组不进入候选。</span>
              </div>
            </div>
            <button class="btn btn-primary h-10" :disabled="loading" @click="runSpeedTest">
              <Icon name="refresh" size="sm" :class="loading ? 'animate-spin' : ''" />
              <span>{{ loading ? '检测中' : '检测 gpt-5.5' }}</span>
            </button>
          </div>

          <div v-if="groupedResultCards.length > 0" class="space-y-5">
            <section v-for="group in groupedResultCards" :key="group.name" class="speed-monitor-section">
              <div class="mb-2">
                <h3>{{ group.name }}</h3>
                <p>共 {{ group.items.length }} 个检测项，成功率 {{ formatPercent(group.successRate) }}</p>
              </div>
              <div class="speed-monitor-grid">
                <article v-for="item in group.items" :key="item.account_id" class="speed-monitor-card">
                  <div class="speed-monitor-card-head">
                    <div class="flex min-w-0 items-center gap-3">
                      <span class="speed-monitor-icon">
                        <Icon name="chart" size="sm" />
                      </span>
                      <div class="min-w-0">
                        <h4>{{ item.name }}</h4>
                        <p>{{ SPEED_TEST_MODEL }} · #{{ item.account_id }}</p>
                      </div>
                    </div>
                    <span :class="['speed-status', item.success ? 'speed-status-ok' : 'speed-status-bad']">
                      {{ item.success ? '正常' : '失败' }}
                    </span>
                  </div>

                  <div class="speed-monitor-metrics">
                    <div>
                      <span>首包延迟</span>
                      <strong>{{ item.success ? formatMs(item.latency_ms) : '-' }}</strong>
                    </div>
                    <div>
                      <span>最近检测</span>
                      <strong>{{ item.finished_at ? formatUnixTime(item.finished_at) : '-' }}</strong>
                    </div>
                  </div>

                  <div class="speed-monitor-card-foot">
                    <div>
                      <span>本次可用率</span>
                      <strong :class="item.success ? 'text-emerald-600 dark:text-emerald-300' : 'text-red-600 dark:text-red-300'">
                        {{ item.success ? '100%' : '0%' }}
                      </strong>
                    </div>
                    <div>
                      <span>状态条</span>
                      <div class="speed-bars" aria-hidden="true">
                        <i
                          v-for="bar in statusBars(item)"
                          :key="bar.index"
                          :class="bar.className"
                        ></i>
                      </div>
                    </div>
                  </div>

                  <p v-if="item.error" class="speed-monitor-error" :title="item.error">{{ item.error }}</p>
                </article>
              </div>
            </section>
          </div>

          <div class="speed-account-panel">
            <div class="flex flex-wrap items-end gap-3">
              <label class="speed-field min-w-[260px] flex-1">
                <span>指定账号测速</span>
                <div class="relative">
                  <Icon name="search" size="sm" class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
                  <input
                    v-model.trim="accountSearch"
                    class="input pl-9"
                    type="text"
                    placeholder="输入账号名 / 备注 / ID 后查找"
                    @focus="openAccountPicker"
                    @keydown.enter.prevent="loadAccountOptions"
                  />

                  <div v-if="accountPickerOpen" class="speed-account-dropdown">
                    <div v-if="accountOptionsLoading" class="speed-account-empty">账号候选加载中...</div>
                    <div v-else-if="accountOptionsError" class="speed-account-empty text-red-600 dark:text-red-300">
                      {{ accountOptionsError }}
                    </div>
                    <div v-else-if="accountOptions.length === 0" class="speed-account-empty">
                      没有匹配账号
                    </div>
                    <template v-else>
                      <button
                        v-for="account in accountOptions"
                        :key="account.id"
                        type="button"
                        class="speed-account-option"
                        :disabled="isAccountSelected(account.id)"
                        @mousedown.prevent="selectAccount(account)"
                      >
                        <div class="min-w-0">
                          <div class="truncate font-semibold text-gray-900 dark:text-white">
                            {{ account.name }}
                          </div>
                          <div class="truncate text-xs text-gray-500 dark:text-gray-400">
                            #{{ account.id }} · {{ platformLabel(account.platform) }} · {{ account.type }} · {{ accountStatusLabel(account.status) }}
                          </div>
                        </div>
                        <span :class="['speed-mini-status', account.schedulable ? 'speed-mini-status-ok' : 'speed-mini-status-muted']">
                          {{ account.schedulable ? '可调度' : '不调度' }}
                        </span>
                      </button>
                    </template>
                  </div>
                </div>
              </label>

              <button class="btn btn-secondary h-10" :disabled="accountOptionsLoading" @click="loadAccountOptions">
                <Icon name="refresh" size="sm" :class="accountOptionsLoading ? 'animate-spin' : ''" />
                <span>查找账号</span>
              </button>

              <button
                v-if="accountPickerOpen"
                class="btn btn-secondary h-10"
                @click="accountPickerOpen = false"
              >
                收起
              </button>
            </div>

            <div v-if="selectedAccounts.length > 0" class="mt-3 flex flex-wrap items-center gap-2">
              <span
                v-for="account in selectedAccounts"
                :key="account.id"
                class="speed-account-chip"
              >
                {{ account.name }} #{{ account.id }}
                <button type="button" @click="removeSelectedAccount(account.id)">×</button>
              </span>
              <button class="speed-clear-button" type="button" @click="clearSelectedAccounts">清空</button>
              <span class="text-xs text-gray-500 dark:text-gray-400">
                已选 {{ selectedAccounts.length }} 个，开始测速时优先只测这些账号。
              </span>
            </div>
          </div>

          <div class="flex flex-wrap items-end gap-3">
            <label class="speed-field min-w-[220px] flex-1">
              <span>搜索账号</span>
              <div class="relative">
                <Icon name="search" size="sm" class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
                <input v-model.trim="filters.search" class="input pl-9" type="text" placeholder="账号名 / 备注" />
              </div>
            </label>

            <label class="speed-field w-full sm:w-36">
              <span>平台</span>
              <select v-model="filters.platform" class="input">
                <option value="">全部</option>
                <option value="anthropic">Claude</option>
                <option value="openai">OpenAI</option>
                <option value="gemini">Gemini</option>
                <option value="antigravity">Antigravity</option>
              </select>
            </label>

            <label class="speed-field w-full sm:w-44">
              <span>分组</span>
              <select v-model="filters.group_id" class="input">
                <option value="" disabled>请选择启用分组</option>
                <option v-for="group in activeGroups" :key="group.id" :value="String(group.id)">
                  {{ group.name }}
                </option>
              </select>
            </label>

            <label class="speed-field w-full sm:w-36">
              <span>状态</span>
              <select v-model="filters.status" class="input">
                <option value="active">启用</option>
                <option value="">全部</option>
                <option value="inactive">停用</option>
                <option value="error">异常</option>
              </select>
            </label>

            <label class="speed-field w-full sm:w-44">
              <span>测试模型</span>
              <div class="speed-fixed-model">{{ SPEED_TEST_MODEL }}</div>
            </label>

            <label class="speed-field w-28">
              <span>数量</span>
              <input v-model.number="filters.limit" class="input" min="1" max="80" type="number" />
            </label>

            <label class="speed-field w-28">
              <span>并发</span>
              <input v-model.number="filters.concurrency" class="input" min="1" max="6" type="number" />
            </label>

            <label class="flex h-10 items-center gap-2 rounded-lg border border-gray-200 px-3 text-sm text-gray-700 dark:border-dark-500 dark:text-gray-200">
              <input v-model="filters.only_schedulable" type="checkbox" class="rounded border-gray-300 text-primary-600 focus:ring-primary-500" />
              只测可调度
            </label>

            <button class="btn btn-primary h-10" :disabled="loading" @click="runSpeedTest">
              <Icon name="refresh" size="sm" :class="loading ? 'animate-spin' : ''" />
              <span>{{ loading ? '检测中' : '开始检测' }}</span>
            </button>
          </div>
        </div>
      </template>

      <template #table>
        <div class="table-wrapper">
          <table>
            <thead>
              <tr>
                <th class="cursor-pointer" @click="setSort('name')">渠道</th>
                <th class="cursor-pointer" @click="setSort('platform')">平台</th>
                <th>分组</th>
                <th>代理</th>
                <th class="cursor-pointer" @click="setSort('success')">状态</th>
                <th class="cursor-pointer" @click="setSort('latency_ms')">延迟</th>
                <th>速度条</th>
                <th>错误</th>
              </tr>
            </thead>
            <tbody>
              <tr v-if="loading">
                <td colspan="8" class="py-10 text-center text-gray-500 dark:text-gray-400">
                  正在按当前筛选测速...
                </td>
              </tr>
              <tr v-else-if="sortedResults.length === 0">
                <td colspan="8" class="py-10 text-center text-gray-500 dark:text-gray-400">
                  {{ emptyTableMessage }}
                </td>
              </tr>
              <template v-else>
                <tr v-for="item in sortedResults" :key="item.account_id">
                  <td>
                    <div class="font-semibold text-gray-900 dark:text-white">{{ item.name }}</div>
                    <div class="text-xs text-gray-500">#{{ item.account_id }} · {{ item.type }}</div>
                  </td>
                  <td>
                    <span class="speed-pill">{{ platformLabel(item.platform) }}</span>
                  </td>
                  <td>
                    <div class="max-w-[260px] truncate text-sm">
                      {{ groupLabel(item) }}
                    </div>
                  </td>
                  <td>
                    <span v-if="item.proxy_id">#{{ item.proxy_id }}</span>
                    <span v-else class="text-gray-400">直连</span>
                  </td>
                  <td>
                    <span :class="['speed-status', item.success ? 'speed-status-ok' : 'speed-status-bad']">
                      {{ item.success ? qualityLabel(item.latency_ms) : '失败' }}
                    </span>
                  </td>
                  <td class="font-mono">
                    {{ item.success ? formatMs(item.latency_ms) : '-' }}
                  </td>
                  <td class="min-w-[160px]">
                    <div class="speed-bar-track">
                      <div
                        class="speed-bar"
                        :class="speedBarClass(item)"
                        :style="{ width: speedBarWidth(item) }"
                      ></div>
                    </div>
                  </td>
                  <td>
                    <div class="max-w-[360px] truncate text-xs text-red-600 dark:text-red-300" :title="item.error || ''">
                      {{ item.error || '-' }}
                    </div>
                  </td>
                </tr>
              </template>
            </tbody>
          </table>
        </div>
      </template>
    </TablePageLayout>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import TablePageLayout from '@/components/layout/TablePageLayout.vue'
import { Icon } from '@/components/icons'
import { adminAPI } from '@/api/admin'
import type { Account, AdminGroup } from '@/types'
import type {
  AccountSpeedTestItem,
  AccountSpeedTestRequest,
  AccountSpeedTestSummary
} from '@/api/admin/accounts'
import { useAppStore } from '@/stores'

type SortKey = 'name' | 'platform' | 'success' | 'latency_ms'
type RunState = 'idle' | 'running' | 'success' | 'empty' | 'error'

const appStore = useAppStore()
const SPEED_TEST_MODEL = 'gpt-5.5'
const PREFERRED_SPEED_GROUP_NAME = 'codex-plus-2.5倍'

const loading = ref(false)
const groups = ref<AdminGroup[]>([])
const accountOptions = ref<Account[]>([])
const selectedAccounts = ref<Account[]>([])
const accountSearch = ref('')
const accountOptionsLoading = ref(false)
const accountOptionsLoaded = ref(false)
const accountOptionsError = ref('')
const accountPickerOpen = ref(false)
const results = ref<AccountSpeedTestItem[]>([])
const summary = ref<AccountSpeedTestSummary | null>(null)
const sortKey = ref<SortKey>('latency_ms')
const sortOrder = ref<'asc' | 'desc'>('asc')
const runState = ref<RunState>('idle')
const runMessage = ref('')
const lastRunAt = ref('')

const filters = ref({
  search: '',
  platform: '',
  group_id: '',
  status: 'active',
  model_id: SPEED_TEST_MODEL,
  limit: 30,
  concurrency: 3,
  only_schedulable: true
})

const selectedAccountIds = computed(() => selectedAccounts.value.map((account) => account.id))

const activeGroups = computed(() => groups.value.filter((group) => group.status === 'active'))

const currentSpeedGroupName = computed(() => {
  if (!filters.value.group_id) return ''
  const groupID = Number(filters.value.group_id)
  return activeGroups.value.find((group) => group.id === groupID)?.name || ''
})

const maxLatency = computed(() => {
  const values = results.value.filter((item) => item.success).map((item) => item.latency_ms)
  return Math.max(1, ...values)
})

const sortedResults = computed(() => {
  return [...results.value].sort((a, b) => {
    let av: string | number | boolean = a[sortKey.value]
    let bv: string | number | boolean = b[sortKey.value]
    if (sortKey.value === 'latency_ms') {
      av = a.success ? a.latency_ms : Number.MAX_SAFE_INTEGER
      bv = b.success ? b.latency_ms : Number.MAX_SAFE_INTEGER
    }
    if (typeof av === 'string' && typeof bv === 'string') {
      return sortOrder.value === 'asc' ? av.localeCompare(bv) : bv.localeCompare(av)
    }
    const delta = Number(av) - Number(bv)
    return sortOrder.value === 'asc' ? delta : -delta
  })
})

const groupedResultCards = computed(() => {
  const buckets = new Map<string, AccountSpeedTestItem[]>()
  for (const item of sortedResults.value) {
    const name = primaryGroupName(item)
    if (!buckets.has(name)) buckets.set(name, [])
    buckets.get(name)?.push(item)
  }
  return Array.from(buckets.entries()).map(([name, items]) => ({
    name,
    items,
    successRate: items.length > 0
      ? items.filter((item) => item.success).length / items.length
      : 0
  }))
})

const fastestDisplay = computed(() => {
  if (!summary.value?.fastest_name) return '-'
  return `${summary.value.fastest_name} · ${formatMs(summary.value.fastest_latency_ms ?? 0)}`
})

const targetModeLabel = computed(() => {
  if (selectedAccounts.value.length > 0) return `指定账号 ${selectedAccounts.value.length}`
  return '当前筛选'
})

const loadingHint = computed(() => {
  if (selectedAccounts.value.length > 0) return `正在测试 ${selectedAccounts.value.length} 个指定账号`
  return '正在按当前筛选测速'
})

const runStatusClass = computed(() => {
  if (loading.value || runState.value === 'running') return 'speed-feedback-running'
  if (runState.value === 'success') return 'speed-feedback-success'
  if (runState.value === 'empty') return 'speed-feedback-empty'
  if (runState.value === 'error') return 'speed-feedback-error'
  return 'speed-feedback-idle'
})

const emptyTableMessage = computed(() => {
  if (runState.value === 'empty' && runMessage.value) return runMessage.value
  if (runState.value === 'error' && runMessage.value) return runMessage.value
  return '还没有测速结果'
})

onMounted(async () => {
  try {
    groups.value = await adminAPI.groups.getAll()
    selectDefaultSpeedGroup()
  } catch (error) {
    console.error('Failed to load groups:', error)
  }
  await loadAccountOptions()
})

function selectDefaultSpeedGroup() {
  if (filters.value.group_id || activeGroups.value.length === 0) return
  const target = activeGroups.value.find((group) => group.name === PREFERRED_SPEED_GROUP_NAME)
    ?? activeGroups.value.find((group) => group.name.toLowerCase().includes('gptplus2.5'))
    ?? activeGroups.value.find((group) => group.name.toLowerCase().includes('plus') && group.name.includes('2.5'))
  if (target) {
    filters.value.group_id = String(target.id)
  }
}

function setSort(key: SortKey) {
  if (sortKey.value === key) {
    sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
    return
  }
  sortKey.value = key
  sortOrder.value = key === 'latency_ms' ? 'asc' : 'desc'
}

async function runSpeedTest() {
  const ids = selectedAccountIds.value
  if (ids.length > 80) {
    appStore.showError('一次最多选择 80 个账号测速')
    return
  }

  loading.value = true
  results.value = []
  summary.value = null
  runState.value = 'running'
  runMessage.value = ''
  lastRunAt.value = ''
  try {
    const normalizedLimit = normalizeNumber(filters.value.limit, 30, 1, 80)
    const payload: AccountSpeedTestRequest = {
      account_ids: ids.length > 0 ? ids : undefined,
      search: filters.value.search || undefined,
      platform: filters.value.platform || undefined,
      status: filters.value.status || undefined,
      group_id: filters.value.group_id ? Number(filters.value.group_id) : undefined,
      model_id: SPEED_TEST_MODEL,
      limit: ids.length > 0 ? Math.max(normalizedLimit, ids.length) : normalizedLimit,
      concurrency: normalizeNumber(filters.value.concurrency, 3, 1, 6),
      only_schedulable: ids.length > 0 ? false : filters.value.only_schedulable
    }
    const response = await adminAPI.accounts.speedTest(payload)
    const nextResults = Array.isArray(response?.results) ? response.results : []
    const nextSummary = normalizeSpeedTestSummary(response?.summary, nextResults)
    results.value = nextResults
    summary.value = nextSummary
    lastRunAt.value = formatDateTime(new Date())

    if (nextResults.length === 0) {
      runState.value = 'empty'
      runMessage.value = ids.length > 0
        ? '测速完成，但选中的账号没有返回结果。请确认账号是否仍存在。'
        : '测速完成，但当前筛选下没有符合条件的账号。'
      appStore.showWarning(runMessage.value)
      return
    }

    runState.value = 'success'
    runMessage.value = `测速完成：共 ${nextSummary.total} 个，成功 ${nextSummary.success}，失败 ${nextSummary.failed}`
    appStore.showSuccess(runMessage.value)
  } catch (error: any) {
    console.error('Speed test failed:', error)
    const message = extractErrorMessage(error, '测速失败')
    runState.value = 'error'
    runMessage.value = message
    lastRunAt.value = formatDateTime(new Date())
    appStore.showError(message)
  } finally {
    loading.value = false
  }
}

async function loadAccountOptions() {
  accountOptionsLoading.value = true
  accountOptionsError.value = ''
  try {
    const response = await adminAPI.accounts.list(1, 50, {
      platform: filters.value.platform || undefined,
      status: filters.value.status || undefined,
      group: filters.value.group_id || undefined,
      search: accountSearch.value || undefined,
      sort_by: 'name',
      sort_order: 'asc',
      lite: 'true'
    })

    const items = Array.isArray(response.items) ? [...response.items] : []
    const directID = parseAccountSearchID(accountSearch.value)
    if (directID) {
      try {
        const directAccount = await adminAPI.accounts.getById(directID)
        if (matchesAccountPickerFilters(directAccount)) {
          items.unshift(directAccount)
        }
      } catch {
        // Ignore direct ID misses; the normal search result still gives feedback.
      }
    }

    accountOptions.value = dedupeAccounts(items)
    accountOptionsLoaded.value = true
    accountPickerOpen.value = true
  } catch (error: any) {
    console.error('Failed to load speed test accounts:', error)
    accountOptions.value = []
    accountOptionsError.value = extractErrorMessage(error, '账号候选加载失败')
  } finally {
    accountOptionsLoading.value = false
  }
}

function openAccountPicker() {
  accountPickerOpen.value = true
  if (!accountOptionsLoaded.value) {
    void loadAccountOptions()
  }
}

function selectAccount(account: Account) {
  if (isAccountSelected(account.id)) return
  selectedAccounts.value = [...selectedAccounts.value, account]
  accountSearch.value = ''
  accountPickerOpen.value = false
}

function removeSelectedAccount(accountID: number) {
  selectedAccounts.value = selectedAccounts.value.filter((account) => account.id !== accountID)
}

function clearSelectedAccounts() {
  selectedAccounts.value = []
}

function isAccountSelected(accountID: number): boolean {
  return selectedAccountIds.value.includes(accountID)
}

function normalizeNumber(value: number, fallback: number, min: number, max: number): number {
  if (!Number.isFinite(value)) return fallback
  return Math.min(max, Math.max(min, Math.floor(value)))
}

function normalizeSpeedTestSummary(
  raw: AccountSpeedTestSummary | undefined,
  items: AccountSpeedTestItem[]
): AccountSpeedTestSummary {
  if (raw && Number.isFinite(raw.total)) return raw
  return summarizeSpeedTestItems(items)
}

function summarizeSpeedTestItems(items: AccountSpeedTestItem[]): AccountSpeedTestSummary {
  const successful = items.filter((item) => item.success)
  const latencies = successful.map((item) => item.latency_ms).filter((latency) => latency > 0).sort((a, b) => a - b)
  const fastest = successful.reduce<AccountSpeedTestItem | null>((best, item) => {
    if (!best) return item
    return item.latency_ms < best.latency_ms ? item : best
  }, null)
  const slowest = successful.reduce<AccountSpeedTestItem | null>((worst, item) => {
    if (!worst) return item
    return item.latency_ms > worst.latency_ms ? item : worst
  }, null)
  const avg = latencies.length > 0
    ? Math.round(latencies.reduce((sum, latency) => sum + latency, 0) / latencies.length)
    : 0
  const p95Index = latencies.length > 0 ? Math.min(latencies.length - 1, Math.ceil(latencies.length * 0.95) - 1) : -1

  return {
    total: items.length,
    success: successful.length,
    failed: items.length - successful.length,
    avg_latency_ms: avg,
    p95_latency_ms: p95Index >= 0 ? latencies[p95Index] : 0,
    fastest_id: fastest?.account_id,
    fastest_name: fastest?.name,
    fastest_latency_ms: fastest?.latency_ms,
    slowest_id: slowest?.account_id,
    slowest_name: slowest?.name,
    slowest_latency_ms: slowest?.latency_ms,
    success_rate: items.length > 0 ? successful.length / items.length : 0
  }
}

function parseAccountSearchID(value: string): number | null {
  const trimmed = value.trim().replace(/^#/, '')
  if (!/^\d+$/.test(trimmed)) return null
  const id = Number(trimmed)
  return Number.isSafeInteger(id) && id > 0 ? id : null
}

function matchesAccountPickerFilters(account: Account): boolean {
  if (filters.value.platform && account.platform !== filters.value.platform) return false
  if (filters.value.status && account.status !== filters.value.status) return false
  if (filters.value.group_id) {
    const groupID = Number(filters.value.group_id)
    const groupIDs = account.group_ids ?? account.groups?.map((group) => group.id) ?? []
    if (!groupIDs.includes(groupID)) return false
  }
  return true
}

function dedupeAccounts(accounts: Account[]): Account[] {
  const seen = new Set<number>()
  const out: Account[] = []
  for (const account of accounts) {
    if (!account || seen.has(account.id)) continue
    seen.add(account.id)
    out.push(account)
  }
  return out
}

function formatMs(value: number): string {
  if (!value || value <= 0) return '-'
  if (value >= 1000) return `${(value / 1000).toFixed(value >= 10000 ? 1 : 2)}s`
  return `${Math.round(value)}ms`
}

function formatDateTime(value: Date): string {
  return value.toLocaleString('zh-CN', { hour12: false })
}

function formatPercent(value: number): string {
  return `${Math.round(value * 100)}%`
}

function platformLabel(platform: string): string {
  const labels: Record<string, string> = {
    anthropic: 'Claude',
    openai: 'OpenAI',
    gemini: 'Gemini',
    antigravity: 'Antigravity'
  }
  return labels[platform] || platform
}

function accountStatusLabel(status: string): string {
  const labels: Record<string, string> = {
    active: '启用',
    inactive: '停用',
    error: '异常'
  }
  return labels[status] || status
}

function groupLabel(item: AccountSpeedTestItem): string {
  if (currentSpeedGroupName.value) return currentSpeedGroupName.value
  if (item.groups?.length) return item.groups.map((group) => group.name).join(' / ')
  if (item.group_ids?.length) return item.group_ids.map((id) => `#${id}`).join(' / ')
  return '未分组'
}

function primaryGroupName(item: AccountSpeedTestItem): string {
  if (currentSpeedGroupName.value) return currentSpeedGroupName.value
  if (item.groups?.length) return item.groups[0].name
  if (item.group_ids?.length) return `分组 #${item.group_ids[0]}`
  return '未分组'
}

function formatUnixTime(value: number): string {
  if (!value) return '-'
  return new Date(value * 1000).toLocaleString('zh-CN', { hour12: false })
}

function statusBars(item: AccountSpeedTestItem): { index: number; className: string }[] {
  const out: { index: number; className: string }[] = []
  const base = item.success ? 'speed-bars-ok' : 'speed-bars-bad'
  const latencyClass = item.success && item.latency_ms > 4000 ? 'speed-bars-warn' : base
  for (let i = 0; i < 36; i++) {
    const isCurrent = i === 35
    const className = isCurrent ? latencyClass : base
    out.push({ index: i, className })
  }
  return out
}

function qualityLabel(latency: number): string {
  if (latency <= 1500) return '很快'
  if (latency <= 4000) return '可用'
  if (latency <= 9000) return '偏慢'
  return '很慢'
}

function speedBarWidth(item: AccountSpeedTestItem): string {
  if (!item.success || item.latency_ms <= 0) return '6%'
  const pct = Math.max(8, Math.min(100, (item.latency_ms / maxLatency.value) * 100))
  return `${pct}%`
}

function speedBarClass(item: AccountSpeedTestItem): string {
  if (!item.success) return 'speed-bar-bad'
  if (item.latency_ms <= 1500) return 'speed-bar-fast'
  if (item.latency_ms <= 4000) return 'speed-bar-ok'
  return 'speed-bar-slow'
}

function extractErrorMessage(error: any, fallback: string): string {
  return String(
    error?.response?.data?.detail ||
    error?.response?.data?.message ||
    error?.message ||
    error?.reason ||
    error?.code ||
    fallback
  )
}
</script>

<style scoped>
.speed-feedback {
  @apply flex flex-wrap items-center justify-between gap-3 rounded-lg border px-4 py-3 text-sm;
}

.speed-feedback-idle,
.speed-feedback-running {
  @apply border-sky-200 bg-sky-50 text-sky-800 dark:border-sky-500/30 dark:bg-sky-500/10 dark:text-sky-200;
}

.speed-feedback-success {
  @apply border-emerald-200 bg-emerald-50 text-emerald-800 dark:border-emerald-500/30 dark:bg-emerald-500/10 dark:text-emerald-200;
}

.speed-feedback-empty {
  @apply border-amber-200 bg-amber-50 text-amber-800 dark:border-amber-500/30 dark:bg-amber-500/10 dark:text-amber-200;
}

.speed-feedback-error {
  @apply border-red-200 bg-red-50 text-red-800 dark:border-red-500/30 dark:bg-red-500/10 dark:text-red-200;
}

.speed-feedback-pill {
  @apply rounded-full bg-white/70 px-3 py-1 text-xs font-semibold dark:bg-dark-700/70;
}

.speed-model-banner {
  @apply flex flex-wrap items-center justify-between gap-3 rounded-lg border border-blue-200 bg-blue-50/80 px-4 py-3 dark:border-blue-500/30 dark:bg-blue-500/10;
}

.speed-model-name {
  @apply rounded-md bg-white px-2.5 py-1 font-mono text-sm font-semibold text-gray-950 shadow-sm dark:bg-dark-700 dark:text-white;
}

.speed-model-note {
  @apply text-xs text-gray-500 dark:text-gray-300;
}

.speed-fixed-model {
  @apply flex h-10 items-center rounded-lg border border-gray-200 bg-gray-50 px-3 font-mono text-sm font-semibold text-gray-900 dark:border-dark-500 dark:bg-dark-700 dark:text-white;
}

.speed-monitor-section h3 {
  @apply text-sm font-semibold text-gray-900 dark:text-white;
}

.speed-monitor-section p {
  @apply text-xs text-gray-500 dark:text-gray-400;
}

.speed-monitor-grid {
  @apply grid gap-3 lg:grid-cols-2;
}

.speed-monitor-card {
  @apply rounded-lg border border-gray-200 bg-white p-4 shadow-sm dark:border-dark-500 dark:bg-dark-700;
}

.speed-monitor-card-head {
  @apply flex items-start justify-between gap-3;
}

.speed-monitor-icon {
  @apply flex h-10 w-10 shrink-0 items-center justify-center rounded-lg border border-gray-200 bg-gray-50 text-gray-800 dark:border-dark-500 dark:bg-dark-800 dark:text-gray-100;
}

.speed-monitor-card h4 {
  @apply truncate text-base font-semibold text-gray-950 dark:text-white;
}

.speed-monitor-card-head p {
  @apply mt-0.5 truncate text-xs text-gray-500 dark:text-gray-400;
}

.speed-monitor-metrics {
  @apply mt-5 grid grid-cols-2 gap-3;
}

.speed-monitor-metrics > div {
  @apply rounded-lg border border-gray-100 bg-gray-50 p-3 dark:border-dark-600 dark:bg-dark-800;
}

.speed-monitor-metrics span,
.speed-monitor-card-foot span {
  @apply block text-xs text-gray-500 dark:text-gray-400;
}

.speed-monitor-metrics strong {
  @apply mt-1 block truncate font-mono text-lg font-semibold text-gray-950 dark:text-white;
}

.speed-monitor-card-foot {
  @apply mt-4 grid gap-3 border-t border-gray-100 pt-4 dark:border-dark-600 sm:grid-cols-[140px_1fr];
}

.speed-monitor-card-foot strong {
  @apply mt-1 block text-xl font-semibold;
}

.speed-bars {
  @apply mt-2 flex h-7 items-end gap-1 overflow-hidden;
}

.speed-bars i {
  @apply block h-7 w-1.5 shrink-0 rounded-full;
}

.speed-bars-ok {
  @apply bg-emerald-500;
}

.speed-bars-warn {
  @apply bg-amber-500;
}

.speed-bars-bad {
  @apply bg-red-500;
}

.speed-monitor-error {
  @apply mt-3 truncate rounded-md bg-red-50 px-3 py-2 text-xs text-red-700 dark:bg-red-500/10 dark:text-red-300;
}

.speed-account-panel {
  @apply rounded-lg border border-gray-200 bg-gray-50/70 p-3 dark:border-dark-500 dark:bg-dark-700/40;
}

.speed-account-dropdown {
  @apply absolute z-30 mt-2 max-h-80 w-full overflow-y-auto rounded-lg border border-gray-200 bg-white p-1 shadow-lg dark:border-dark-500 dark:bg-dark-700;
}

.speed-account-empty {
  @apply px-3 py-4 text-center text-sm text-gray-500 dark:text-gray-400;
}

.speed-account-option {
  @apply flex w-full items-center justify-between gap-3 rounded-md px-3 py-2 text-left transition hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-50 dark:hover:bg-dark-600;
}

.speed-mini-status {
  @apply shrink-0 rounded-full px-2 py-0.5 text-xs font-semibold;
}

.speed-mini-status-ok {
  @apply bg-emerald-100 text-emerald-700 dark:bg-emerald-500/15 dark:text-emerald-300;
}

.speed-mini-status-muted {
  @apply bg-gray-100 text-gray-500 dark:bg-dark-600 dark:text-gray-300;
}

.speed-account-chip {
  @apply inline-flex items-center gap-1.5 rounded-full bg-primary-100 px-3 py-1 text-xs font-semibold text-primary-700 dark:bg-primary-500/15 dark:text-primary-200;
}

.speed-account-chip button {
  @apply text-sm leading-none text-primary-500 hover:text-primary-700 dark:text-primary-200 dark:hover:text-white;
}

.speed-clear-button {
  @apply text-xs font-semibold text-gray-500 hover:text-gray-800 dark:text-gray-400 dark:hover:text-gray-100;
}

.speed-stat {
  @apply flex min-h-[74px] flex-col justify-between rounded-lg border border-gray-200 bg-white px-4 py-3 dark:border-dark-500 dark:bg-dark-700;
}

.speed-stat-label {
  @apply text-xs font-medium text-gray-500 dark:text-gray-400;
}

.speed-stat strong {
  @apply truncate text-xl font-semibold text-gray-900 dark:text-white;
}

.speed-field {
  @apply flex flex-col gap-1.5;
}

.speed-field span {
  @apply text-xs font-medium text-gray-500 dark:text-gray-400;
}

.speed-pill {
  @apply inline-flex rounded-md bg-gray-100 px-2 py-1 text-xs font-semibold text-gray-700 dark:bg-dark-600 dark:text-gray-200;
}

.speed-status {
  @apply inline-flex rounded-full px-2.5 py-1 text-xs font-semibold;
}

.speed-status-ok {
  @apply bg-emerald-100 text-emerald-700 dark:bg-emerald-500/15 dark:text-emerald-300;
}

.speed-status-bad {
  @apply bg-red-100 text-red-700 dark:bg-red-500/15 dark:text-red-300;
}

.speed-bar-track {
  @apply h-2.5 w-full overflow-hidden rounded-full bg-gray-100 dark:bg-dark-600;
}

.speed-bar {
  @apply h-full rounded-full transition-all duration-300;
}

.speed-bar-fast {
  @apply bg-emerald-500;
}

.speed-bar-ok {
  @apply bg-sky-500;
}

.speed-bar-slow {
  @apply bg-amber-500;
}

.speed-bar-bad {
  @apply bg-red-500;
}
</style>
