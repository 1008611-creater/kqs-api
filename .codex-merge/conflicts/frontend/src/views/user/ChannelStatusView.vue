<template>
<<<<<<< HEAD
  <ChannelStatusV1View v-if="isV1" />
  <ChannelStatusV2View v-else />
=======
  <AppLayout>
    <main class="channel-status-page">
      <section class="channel-hero">
        <div class="channel-hero__copy">
          <p class="channel-kicker">CHANNEL MONITOR</p>
          <h1>渠道监测</h1>
        </div>

        <div class="channel-hero__status">
          <span class="channel-status-chip" :class="statusChipClass">
            <span :class="statusDotClass"></span>
            {{ statusLabel }}
          </span>
          <button type="button" class="channel-refresh" :disabled="loading || pro4Loading || gpt56Loading" title="刷新" @click="manualReload">
            <Icon name="refresh" size="md" :class="loading || pro4Loading || gpt56Loading ? 'animate-spin' : ''" />
          </button>
        </div>
      </section>

      <section class="channel-monitor-card" :class="gpt56SummaryToneClass">
        <div class="channel-monitor-card__head">
          <div>
            <p class="channel-kicker">置顶监测对象</p>
            <h2>{{ gpt56GroupName }}</h2>
          </div>
          <div class="channel-monitor-card__status">
            <span class="channel-status-chip" :class="gpt56StatusChipClass">
              <span :class="gpt56StatusDotClass"></span>
              {{ gpt56StatusLabel }}
            </span>
            <div class="channel-refresh-meta">
              <span>自动刷新 {{ countdownSeconds }}s</span>
              <span>最近检测 {{ formatDateTime(gpt56Status?.last_checked_at) }}</span>
            </div>
          </div>
        </div>

        <div class="channel-monitor-card__main">
          <div class="channel-monitor-card__summary">
            <p class="channel-kicker">当前结论</p>
            <h3>{{ gpt56SummaryTitle }}</h3>
            <p>{{ gpt56SummaryDescription }}</p>
            <div class="channel-monitor-card__model">
              <span>检测模型</span>
              <strong>{{ gpt56ModelName }}</strong>
            </div>
          </div>

          <div class="channel-monitor-card__value">
            <span>额度价值</span>
            <strong>{{ gpt56RateMetric.impact }}</strong>
            <em>{{ gpt56RateMetric.value }}</em>
          </div>
        </div>

        <div class="channel-monitor-card__metrics">
          <MetricBox label="1小时可用率" :value="formatPercent(gpt56Status?.success_rate_1h)" hint="短时间稳定性" :tone="metricTone(gpt56Status?.success_rate_1h)" />
          <MetricBox label="24小时可用率" :value="formatPercent(gpt56Status?.success_rate_24h)" hint="全天整体表现" :tone="metricTone(gpt56Status?.success_rate_24h)" />
          <MetricBox label="平均延迟" :value="formatLatency(gpt56Status?.average_latency_ms)" hint="常规请求体感" :tone="latencyTone(gpt56Status?.average_latency_ms)" />
          <MetricBox label="P95延迟" :value="formatLatency(gpt56Status?.p95_latency_ms)" hint="高峰慢请求参考" :tone="latencyTone(gpt56Status?.p95_latency_ms)" />
        </div>

        <p class="channel-monitor-card__note">{{ gpt56RateMetric.note }}</p>
      </section>

      <section class="channel-timeline-panel">
        <div class="channel-section-title">
          <div>
            <p>最近两小时</p>
            <h2>GPT-5.6 自动检测结果</h2>
          </div>
          <span>{{ gpt56TimelineSummary }}</span>
        </div>

        <div class="channel-timeline-frame">
          <div class="channel-timeline-axis">
            <span>2小时前</span>
            <span>现在</span>
          </div>
          <div class="channel-timeline" aria-label="GPT-5.6 最近两小时自动检测状态">
            <div
              v-for="(point, index) in gpt56Timeline"
              :key="`gpt56-${point.checked_at}-${index}`"
              class="channel-timeline__bar"
              :class="timelineClass(point.status)"
              :style="{ height: timelineHeight(point.status) }"
              :title="timelineTooltip(point)"
              :data-tooltip="timelineTooltip(point)"
            ></div>
            <template v-if="gpt56Timeline.length === 0">
              <div v-for="index in 24" :key="`gpt56-empty-${index}`" class="channel-timeline__bar channel-timeline__bar--unknown" style="height: 24%"></div>
            </template>
          </div>
        </div>

        <div class="channel-legend">
          <LegendDot status="operational" label="正常" />
          <LegendDot status="degraded" label="波动" />
          <LegendDot status="failed" label="不可用" />
          <LegendDot status="unknown" label="暂无数据" />
        </div>
      </section>

      <section class="channel-monitor-card" :class="summaryToneClass">
        <div class="channel-monitor-card__head">
          <div>
            <p class="channel-kicker">监测对象</p>
            <h2>{{ groupName }}</h2>
          </div>
          <div class="channel-monitor-card__status">
            <span class="channel-status-chip" :class="statusChipClass">
              <span :class="statusDotClass"></span>
              {{ statusLabel }}
            </span>
            <div class="channel-refresh-meta">
              <span>自动刷新 {{ countdownSeconds }}s</span>
              <span>最近检测 {{ formatDateTime(status?.last_checked_at) }}</span>
            </div>
          </div>
        </div>

        <div class="channel-monitor-card__main">
          <div class="channel-monitor-card__summary">
            <p class="channel-kicker">当前结论</p>
            <h3>{{ summaryTitle }}</h3>
            <p>{{ summaryDescription }}</p>
            <div class="channel-monitor-card__model">
              <span>检测模型</span>
              <strong>{{ modelName }}</strong>
            </div>
          </div>

          <div class="channel-monitor-card__value">
            <span>额度价值</span>
            <strong>{{ primaryRateMetric.impact }}</strong>
            <em>{{ primaryRateMetric.value }}</em>
          </div>
        </div>

        <div class="channel-monitor-card__metrics">
          <MetricBox label="1小时可用率" :value="formatPercent(status?.success_rate_1h)" hint="短时间稳定性" :tone="metricTone(status?.success_rate_1h)" />
          <MetricBox label="24小时可用率" :value="formatPercent(status?.success_rate_24h)" hint="全天整体表现" :tone="metricTone(status?.success_rate_24h)" />
          <MetricBox label="平均延迟" :value="formatLatency(status?.average_latency_ms)" hint="常规请求体感" :tone="latencyTone(status?.average_latency_ms)" />
          <MetricBox label="P95延迟" :value="formatLatency(status?.p95_latency_ms)" hint="高峰慢请求参考" :tone="latencyTone(status?.p95_latency_ms)" />
        </div>

        <p class="channel-monitor-card__note">{{ primaryRateMetric.note }}</p>
      </section>

      <section class="channel-timeline-panel">
        <div class="channel-section-title">
          <div>
            <p>最近两小时</p>
            <h2>自动检测结果</h2>
          </div>
          <span>{{ timelineSummary }}</span>
        </div>

        <div class="channel-timeline-frame">
          <div class="channel-timeline-axis">
            <span>2小时前</span>
            <span>现在</span>
          </div>
          <div class="channel-timeline" aria-label="最近两小时自动检测状态">
            <div
              v-for="(point, index) in timeline"
              :key="`${point.checked_at}-${index}`"
              class="channel-timeline__bar"
              :class="timelineClass(point.status)"
              :style="{ height: timelineHeight(point.status) }"
              :title="timelineTooltip(point)"
              :data-tooltip="timelineTooltip(point)"
            ></div>
            <template v-if="timeline.length === 0">
              <div v-for="index in 24" :key="`empty-${index}`" class="channel-timeline__bar channel-timeline__bar--unknown" style="height: 24%"></div>
            </template>
          </div>
        </div>

        <div class="channel-legend">
          <LegendDot status="operational" label="正常" />
          <LegendDot status="degraded" label="波动" />
          <LegendDot status="failed" label="不可用" />
          <LegendDot status="unknown" label="暂无数据" />
        </div>
      </section>

      <section class="channel-monitor-card" :class="pro4SummaryToneClass">
        <div class="channel-monitor-card__head">
          <div>
            <p class="channel-kicker">监测对象</p>
            <h2>{{ pro4GroupName }}</h2>
          </div>
          <div class="channel-monitor-card__status">
            <span class="channel-status-chip" :class="pro4StatusChipClass">
              <span :class="pro4StatusDotClass"></span>
              {{ pro4StatusLabel }}
            </span>
            <div class="channel-refresh-meta">
              <span>自动刷新 {{ countdownSeconds }}s</span>
              <span>最近检测 {{ formatDateTime(pro4Status?.last_checked_at) }}</span>
            </div>
          </div>
        </div>

        <div class="channel-monitor-card__main">
          <div class="channel-monitor-card__summary">
            <p class="channel-kicker">当前结论</p>
            <h3>{{ pro4SummaryTitle }}</h3>
            <p>{{ pro4SummaryDescription }}</p>
            <div class="channel-monitor-card__model">
              <span>检测模型</span>
              <strong>{{ pro4ModelName }}</strong>
            </div>
          </div>

          <div class="channel-monitor-card__value">
            <span>额度价值</span>
            <strong>{{ pro4RateMetric.impact }}</strong>
            <em>{{ pro4RateMetric.value }}</em>
          </div>
        </div>

        <div class="channel-monitor-card__metrics">
          <MetricBox label="1小时可用率" :value="formatPercent(pro4Status?.success_rate_1h)" hint="短时间稳定性" :tone="metricTone(pro4Status?.success_rate_1h)" />
          <MetricBox label="24小时可用率" :value="formatPercent(pro4Status?.success_rate_24h)" hint="全天整体表现" :tone="metricTone(pro4Status?.success_rate_24h)" />
          <MetricBox label="平均延迟" :value="formatLatency(pro4Status?.average_latency_ms)" hint="常规请求体感" :tone="latencyTone(pro4Status?.average_latency_ms)" />
          <MetricBox label="P95延迟" :value="formatLatency(pro4Status?.p95_latency_ms)" hint="高峰慢请求参考" :tone="latencyTone(pro4Status?.p95_latency_ms)" />
        </div>

        <p class="channel-monitor-card__note">{{ pro4RateMetric.note }}</p>
      </section>

      <section class="channel-timeline-panel">
        <div class="channel-section-title">
          <div>
            <p>最近两小时</p>
            <h2>pro4 倍自动检测结果</h2>
          </div>
          <span>{{ pro4TimelineSummary }}</span>
        </div>

        <div class="channel-timeline-frame">
          <div class="channel-timeline-axis">
            <span>2小时前</span>
            <span>现在</span>
          </div>
          <div class="channel-timeline" aria-label="pro4 倍最近两小时自动检测状态">
            <div
              v-for="(point, index) in pro4Timeline"
              :key="`pro4-${point.checked_at}-${index}`"
              class="channel-timeline__bar"
              :class="timelineClass(point.status)"
              :style="{ height: timelineHeight(point.status) }"
              :title="timelineTooltip(point)"
              :data-tooltip="timelineTooltip(point)"
            ></div>
            <template v-if="pro4Timeline.length === 0">
              <div v-for="index in 24" :key="`pro4-empty-${index}`" class="channel-timeline__bar channel-timeline__bar--unknown" style="height: 24%"></div>
            </template>
          </div>
        </div>

        <div class="channel-legend">
          <LegendDot status="operational" label="正常" />
          <LegendDot status="degraded" label="波动" />
          <LegendDot status="failed" label="不可用" />
          <LegendDot status="unknown" label="暂无数据" />
        </div>
      </section>

      <section class="channel-guidance">
        <article>
          <span class="channel-guidance__icon"><Icon name="checkCircle" size="sm" /></span>
          <div>
            <h3>正常</h3>
            <p>可以按平时方式调用。偶发慢请求通常会自动切换或重试。</p>
          </div>
        </article>
        <article>
          <span class="channel-guidance__icon"><Icon name="exclamationTriangle" size="sm" /></span>
          <div>
            <h3>部分波动</h3>
            <p>短时间可能出现响应变慢或个别请求失败，建议稍后重试。</p>
          </div>
        </article>
        <article>
          <span class="channel-guidance__icon"><Icon name="clock" size="sm" /></span>
          <div>
            <h3>不可用</h3>
            <p>如果状态连续异常，系统会继续检测，管理员会处理上游稳定性。</p>
          </div>
        </article>
      </section>
    </main>
  </AppLayout>
>>>>>>> 1157e4cfe6271a81c9dc96df84af6c6a3b22f831
</template>
<script setup lang="ts">
<<<<<<< HEAD
import { computed } from 'vue'
import { isChannelMonitorV1Mode } from '@/utils/featureFlags'
import ChannelStatusV1View from './ChannelStatusV1View.vue'
import ChannelStatusV2View from './ChannelStatusV2View.vue'

const isV1 = computed(() => isChannelMonitorV1Mode())
=======
import { computed, defineComponent, h, onBeforeUnmount, onMounted, ref } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import { publicStatus, type PublicChannelStatus, type PublicChannelStatusPoint } from '@/api/channelMonitor'
import { useAppStore } from '@/stores/app'
import { extractApiErrorMessage } from '@/utils/apiError'

const REFRESH_SECONDS = 60
const LOWEST_RECHARGE_RATE_CNY_PER_USD = 0.071
const USD_CNY_RATE_ESTIMATE = 6.76
const PRIMARY_GROUP_MULTIPLIER = 2
const PRO4_GROUP_MULTIPLIER = 4
const GPT56_GROUP_MULTIPLIER = 3

const appStore = useAppStore()
const status = ref<PublicChannelStatus | null>(null)
const pro4Status = ref<PublicChannelStatus | null>(null)
const gpt56Status = ref<PublicChannelStatus | null>(null)
const loading = ref(false)
const pro4Loading = ref(false)
const gpt56Loading = ref(false)
const countdownSeconds = ref(REFRESH_SECONDS)
let abortController: AbortController | null = null
let pro4AbortController: AbortController | null = null
let gpt56AbortController: AbortController | null = null
let reloadTimer: number | undefined
let countdownTimer: number | undefined

function createEmptyStatus(groupName: string): PublicChannelStatus {
  return {
    group_name: groupName,
    model_id: 'gpt-5.5',
    status: 'unknown',
    active_channels: 0,
    checked_channels: 0,
    success_rate_1h: null,
    success_rate_24h: null,
    average_latency_ms: null,
    p95_latency_ms: null,
    last_checked_at: null,
    generated_at: new Date().toISOString(),
    timeline: Array.from({ length: 24 }, (_, index) => ({
      status: 'unknown',
      checked_at: new Date(Date.now() - (23 - index) * 5 * 60 * 1000).toISOString(),
    })),
  }
}

const MetricBox = defineComponent({
  props: {
    label: { type: String, required: true },
    value: { type: String, required: true },
    hint: { type: String, required: true },
    tone: { type: String, default: 'neutral' },
    tooltip: { type: String, default: '' },
  },
  setup(props) {
    return () =>
      h('article', {
        class: ['channel-metric', `channel-metric--${props.tone}`, props.tooltip ? 'channel-metric--has-tooltip' : ''],
        title: props.tooltip || undefined,
        'data-tooltip': props.tooltip || undefined,
      }, [
        h('span', { class: 'channel-metric__label' }, props.label),
        h('strong', props.value),
        h('span', { class: 'channel-metric__hint' }, [
          props.hint,
          props.tooltip ? h('span', { class: 'channel-metric__info', 'aria-hidden': 'true' }, '?') : null,
        ]),
      ])
  },
})

const LegendDot = defineComponent({
  props: {
    status: { type: String, required: true },
    label: { type: String, required: true },
  },
  setup(props) {
    return () =>
      h('span', { class: 'channel-legend__item' }, [
        h('span', { class: ['channel-legend__dot', legendDotClass(props.status)] }),
        h('span', props.label),
      ])
  },
})

const timeline = computed<PublicChannelStatusPoint[]>(() => status.value?.timeline || [])
const pro4Timeline = computed<PublicChannelStatusPoint[]>(() => pro4Status.value?.timeline || [])
const gpt56Timeline = computed<PublicChannelStatusPoint[]>(() => gpt56Status.value?.timeline || [])
const groupName = computed(() => status.value?.group_name || 'GPT-plus')
const pro4GroupName = computed(() => pro4Status.value?.group_name || 'GPT-pro')
const gpt56GroupName = computed(() => gpt56Status.value?.group_name || 'GPT-5.6')
const modelName = computed(() => status.value?.model_id || 'gpt-5.5')
const pro4ModelName = computed(() => pro4Status.value?.model_id || 'gpt-5.5')
const gpt56ModelName = computed(() => gpt56Status.value?.model_id || 'gpt-5.6')
const primaryRateMetric = computed(() => buildRateMetric(PRIMARY_GROUP_MULTIPLIER))
const pro4RateMetric = computed(() => buildRateMetric(PRO4_GROUP_MULTIPLIER))
const gpt56RateMetric = computed(() => buildRateMetric(GPT56_GROUP_MULTIPLIER))

type NormalizedMonitorStatus = 'operational' | 'degraded' | 'failed' | 'unknown'

function normalizeMonitorStatus(value: string | null | undefined): NormalizedMonitorStatus {
  switch ((value || '').toLowerCase()) {
    case 'operational':
    case 'success':
    case 'ok':
    case 'healthy':
      return 'operational'
    case 'degraded':
    case 'warning':
    case 'partial':
    case 'partial_success':
      return 'degraded'
    case 'failed':
    case 'error':
    case 'timeout':
    case 'unavailable':
      return 'failed'
    default:
      return 'unknown'
  }
}

const statusLabel = computed(() => {
  return statusLabelFor(status.value)
})

const pro4StatusLabel = computed(() => {
  return statusLabelFor(pro4Status.value)
})

const gpt56StatusLabel = computed(() => {
  return statusLabelFor(gpt56Status.value)
})

function statusLabelFor(snapshot: PublicChannelStatus | null) {
  switch (normalizeMonitorStatus(snapshot?.status)) {
    case 'operational':
      return '运行正常'
    case 'degraded':
      return '部分波动'
    case 'failed':
      return '当前异常'
    default:
      return '暂无数据'
  }
}

const summaryTitle = computed(() => {
  return summaryTitleFor(status.value)
})

const pro4SummaryTitle = computed(() => {
  return summaryTitleFor(pro4Status.value)
})

const gpt56SummaryTitle = computed(() => {
  return summaryTitleFor(gpt56Status.value)
})

function summaryTitleFor(snapshot: PublicChannelStatus | null) {
  switch (normalizeMonitorStatus(snapshot?.status)) {
    case 'operational':
      return '当前可正常使用'
    case 'degraded':
      return '当前有波动，建议保留重试'
    case 'failed':
      return '当前检测异常'
    default:
      return '暂无检测数据'
  }
}

const summaryDescription = computed(() => {
  return summaryDescriptionFor(status.value)
})

const pro4SummaryDescription = computed(() => {
  return summaryDescriptionFor(pro4Status.value)
})

const gpt56SummaryDescription = computed(() => {
  return summaryDescriptionFor(gpt56Status.value)
})

function summaryDescriptionFor(snapshot: PublicChannelStatus | null) {
  switch (normalizeMonitorStatus(snapshot?.status)) {
    case 'operational':
      return '最近检测整体通过，延迟和可用率处于可接受范围。'
    case 'degraded':
      return '近期检测出现部分失败或延迟偏高，调用时可能出现慢响应。'
    case 'failed':
      return '最近检测未通过，短时间内可能无法稳定调用。'
    default:
      return '系统会自动刷新检测结果，当前还没有可展示的有效检测点。'
  }
}

const timelineSummary = computed(() => {
  return timelineSummaryFor(timeline.value)
})

const pro4TimelineSummary = computed(() => {
  return timelineSummaryFor(pro4Timeline.value)
})

const gpt56TimelineSummary = computed(() => {
  return timelineSummaryFor(gpt56Timeline.value)
})

function timelineSummaryFor(points: PublicChannelStatusPoint[]) {
  const counts = points.reduce(
    (acc, point) => {
      const normalized = normalizeMonitorStatus(point.status)
      acc[normalized] = (acc[normalized] || 0) + 1
      return acc
    },
    {} as Record<string, number>,
  )
  const total = points.length
  if (!total) return '暂无检测点'
  if ((counts.unknown || 0) === total) return '暂无检测点'
  const ok = counts.operational || 0
  const warn = counts.degraded || 0
  const bad = counts.failed || 0
  return `${total} 次检测 · 正常 ${ok} · 波动 ${warn} · 异常 ${bad}`
}

const statusChipClass = computed(() => {
  return statusChipClassFor(status.value)
})

const pro4StatusChipClass = computed(() => {
  return statusChipClassFor(pro4Status.value)
})

const gpt56StatusChipClass = computed(() => {
  return statusChipClassFor(gpt56Status.value)
})

function statusChipClassFor(snapshot: PublicChannelStatus | null) {
  switch (normalizeMonitorStatus(snapshot?.status)) {
    case 'operational':
      return 'channel-status-chip--ok'
    case 'degraded':
      return 'channel-status-chip--warn'
    case 'failed':
      return 'channel-status-chip--bad'
    default:
      return 'channel-status-chip--unknown'
  }
}

const summaryToneClass = computed(() => {
  return summaryToneClassFor(status.value)
})

const pro4SummaryToneClass = computed(() => {
  return summaryToneClassFor(pro4Status.value)
})

const gpt56SummaryToneClass = computed(() => {
  return summaryToneClassFor(gpt56Status.value)
})

function summaryToneClassFor(snapshot: PublicChannelStatus | null) {
  switch (normalizeMonitorStatus(snapshot?.status)) {
    case 'operational':
      return 'channel-summary-card--ok'
    case 'degraded':
      return 'channel-summary-card--warn'
    case 'failed':
      return 'channel-summary-card--bad'
    default:
      return 'channel-summary-card--unknown'
  }
}

const statusDotClass = computed(() => ['channel-status-chip__dot', legendDotClass(status.value?.status || 'unknown')])
const pro4StatusDotClass = computed(() => ['channel-status-chip__dot', legendDotClass(pro4Status.value?.status || 'unknown')])
const gpt56StatusDotClass = computed(() => ['channel-status-chip__dot', legendDotClass(gpt56Status.value?.status || 'unknown')])

function isAbortError(err: unknown) {
  const e = err as { name?: string; code?: string }
  return e?.name === 'AbortError' || e?.code === 'ERR_CANCELED'
}

function isExpectedMonitorEmpty(err: unknown) {
  const e = err as { status?: number; code?: string | number }
  return e?.status === 404 || e?.code === 404 || e?.code === 'NOT_FOUND'
}

async function reload(silent = false) {
  if (abortController) abortController.abort()
  if (pro4AbortController) pro4AbortController.abort()
  if (gpt56AbortController) gpt56AbortController.abort()
  const ctrl = new AbortController()
  const pro4Ctrl = new AbortController()
  const gpt56Ctrl = new AbortController()
  abortController = ctrl
  pro4AbortController = pro4Ctrl
  gpt56AbortController = gpt56Ctrl
  if (!silent) {
    loading.value = true
    pro4Loading.value = true
    gpt56Loading.value = true
  }
  try {
    const [primaryRes, pro4Res, gpt56Res] = await Promise.allSettled([
      publicStatus({ signal: ctrl.signal }),
      publicStatus({ profile: 'pro4', signal: pro4Ctrl.signal }),
      publicStatus({ profile: 'gpt56', signal: gpt56Ctrl.signal }),
    ])
    if (
      ctrl.signal.aborted ||
      pro4Ctrl.signal.aborted ||
      gpt56Ctrl.signal.aborted ||
      abortController !== ctrl ||
      pro4AbortController !== pro4Ctrl ||
      gpt56AbortController !== gpt56Ctrl
    ) return

    const failures: unknown[] = []
    if (primaryRes.status === 'fulfilled') {
      status.value = primaryRes.value
    } else if (!isAbortError(primaryRes.reason)) {
      if (!status.value) status.value = createEmptyStatus('gptplus2.5倍')
      if (!isExpectedMonitorEmpty(primaryRes.reason)) failures.push(primaryRes.reason)
    }

    if (pro4Res.status === 'fulfilled') {
      pro4Status.value = pro4Res.value
    } else if (!isAbortError(pro4Res.reason)) {
      if (!pro4Status.value) pro4Status.value = createEmptyStatus('pro4倍')
      if (!isExpectedMonitorEmpty(pro4Res.reason)) failures.push(pro4Res.reason)
    }

    if (gpt56Res.status === 'fulfilled') {
      gpt56Status.value = gpt56Res.value
    } else if (!isAbortError(gpt56Res.reason)) {
      if (!gpt56Status.value) gpt56Status.value = createEmptyStatus('GPT-5.6')
      if (!isExpectedMonitorEmpty(gpt56Res.reason)) failures.push(gpt56Res.reason)
    }

    countdownSeconds.value = REFRESH_SECONDS
    if (!silent && failures.length === 3) {
      appStore.showError(extractApiErrorMessage(failures[0], '加载服务状态失败'))
    }
  } finally {
    if (abortController === ctrl) {
      loading.value = false
      abortController = null
    }
    if (pro4AbortController === pro4Ctrl) {
      pro4Loading.value = false
      pro4AbortController = null
    }
    if (gpt56AbortController === gpt56Ctrl) {
      gpt56Loading.value = false
      gpt56AbortController = null
    }
  }
}

function manualReload() {
  countdownSeconds.value = REFRESH_SECONDS
  void reload(false)
}

function formatPercent(value: number | null | undefined) {
  if (typeof value !== 'number') return '-'
  return `${value.toFixed(1)}%`
}

function formatLatency(value: number | null | undefined) {
  if (typeof value !== 'number' || value <= 0) return '-'
  if (value >= 1000) return `${(value / 1000).toFixed(1)}s`
  return `${value}ms`
}

function buildRateMetric(multiplier: number) {
  const actualPrice = multiplier * LOWEST_RECHARGE_RATE_CNY_PER_USD
  const quotaRatio = USD_CNY_RATE_ESTIMATE / actualPrice
  const actualPriceText = actualPrice.toFixed(3)
  const quotaRatioText = quotaRatio.toFixed(1)
  return {
    value: `${multiplier}x`,
    hint: `1元 ≈ ${quotaRatioText}元官方额度起`,
    impact: `1元 ≈ ${quotaRatioText}元官方额度起`,
    note: `按 1 美元≈${USD_CNY_RATE_ESTIMATE.toFixed(2)} 元、最低充值倍率 ${LOWEST_RECHARGE_RATE_CNY_PER_USD} 估算，实际比例会随充值档位浮动。`,
    tooltip: [
      '按当前最低充值倍率和汇率估算。',
      '公式：1元可用官方人民币额度 = 美元人民币汇率 ÷（分组倍率 × 充值倍率）。',
      `当前按 1 美元 ≈ ${USD_CNY_RATE_ESTIMATE.toFixed(2)} 元、最低充值倍率 ${LOWEST_RECHARGE_RATE_CNY_PER_USD} 计算：`,
      `${USD_CNY_RATE_ESTIMATE.toFixed(2)} ÷ (${multiplier} × ${LOWEST_RECHARGE_RATE_CNY_PER_USD}) ≈ ${quotaRatioText}。`,
      `该分组最低约 ${actualPriceText} 元 / 官方 $1，充值档位不同，实际比例会浮动。`,
    ].join('\n'),
  }
}

function formatDateTime(value: string | null | undefined) {
  if (!value) return '等待首次检测'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return '等待首次检测'
  return date.toLocaleString()
}

function metricTone(value: number | null | undefined) {
  if (typeof value !== 'number') return 'neutral'
  if (value >= 97) return 'ok'
  if (value >= 90) return 'warn'
  return 'bad'
}

function latencyTone(value: number | null | undefined) {
  if (typeof value !== 'number' || value <= 0) return 'neutral'
  if (value <= 6000) return 'ok'
  if (value <= 15000) return 'warn'
  return 'bad'
}

function timelineClass(value: string) {
  switch (normalizeMonitorStatus(value)) {
    case 'operational':
      return 'channel-timeline__bar--ok'
    case 'degraded':
      return 'channel-timeline__bar--warn'
    case 'failed':
      return 'channel-timeline__bar--bad'
    default:
      return 'channel-timeline__bar--unknown'
  }
}

function legendDotClass(value: string) {
  switch (normalizeMonitorStatus(value)) {
    case 'operational':
      return 'channel-dot--ok'
    case 'degraded':
      return 'channel-dot--warn'
    case 'failed':
      return 'channel-dot--bad'
    default:
      return 'channel-dot--unknown'
  }
}

function timelineHeight(value: string) {
  switch (normalizeMonitorStatus(value)) {
    case 'operational':
      return '100%'
    case 'degraded':
      return '68%'
    case 'failed':
      return '42%'
    default:
      return '24%'
  }
}

function timelineLabel(value: string) {
  switch (normalizeMonitorStatus(value)) {
    case 'operational':
      return '正常'
    case 'degraded':
      return '波动'
    case 'failed':
      return '不可用'
    default:
      return '暂无数据'
  }
}

function timelineTooltip(point: PublicChannelStatusPoint) {
  return `检测时间：${formatDateTime(point.checked_at)}\n状态：${timelineLabel(point.status)}`
}

onMounted(() => {
  void reload(false)
  reloadTimer = window.setInterval(() => {
    void reload(true)
  }, REFRESH_SECONDS * 1000)
  countdownTimer = window.setInterval(() => {
    countdownSeconds.value = countdownSeconds.value <= 1 ? REFRESH_SECONDS : countdownSeconds.value - 1
  }, 1000)
})

onBeforeUnmount(() => {
  if (abortController) abortController.abort()
  if (pro4AbortController) pro4AbortController.abort()
  if (gpt56AbortController) gpt56AbortController.abort()
  if (reloadTimer) window.clearInterval(reloadTimer)
  if (countdownTimer) window.clearInterval(countdownTimer)
})
>>>>>>> 1157e4cfe6271a81c9dc96df84af6c6a3b22f831
</script>

<style scoped>
.channel-status-page {
  display: grid;
  gap: 0.75rem;
}

.channel-hero,
.channel-group-divider,
.channel-monitor-card,
.channel-summary-card,
.channel-overview,
.channel-timeline-panel,
.channel-guidance article {
  border: 1px solid rgba(43, 132, 83, 0.38);
  border-radius: 8px;
  background: rgba(5, 29, 24, 0.78);
}

.channel-monitor-card {
  position: relative;
  display: grid;
  gap: 1rem;
  overflow: hidden;
  padding: 1rem;
  background:
    radial-gradient(circle at 82% 26%, rgba(46, 148, 91, 0.18), transparent 32%),
    linear-gradient(135deg, rgba(34, 197, 94, 0.08), rgba(5, 29, 24, 0.92) 48%),
    rgba(5, 29, 24, 0.78);
}

.channel-monitor-card::before {
  position: absolute;
  inset: 0;
  border-top: 1px solid rgba(148, 226, 170, 0.18);
  background:
    linear-gradient(90deg, rgba(126, 227, 155, 0.08), transparent 28%),
    linear-gradient(180deg, rgba(255, 255, 255, 0.04), transparent 44%);
  content: "";
  pointer-events: none;
}

.channel-monitor-card > * {
  position: relative;
  z-index: 1;
}

.channel-monitor-card__head,
.channel-monitor-card__status,
.channel-monitor-card__main,
.channel-monitor-card__metrics {
  display: flex;
  gap: 1rem;
}

.channel-monitor-card__head {
  align-items: flex-start;
  justify-content: space-between;
}

.channel-monitor-card__head h2 {
  margin: 0.22rem 0 0;
  color: #edf7ee;
  font-size: 1.32rem;
  font-weight: 900;
}

.channel-monitor-card__status {
  flex-wrap: wrap;
  align-items: center;
  justify-content: flex-end;
}

.channel-monitor-card__main {
  display: grid;
  grid-template-columns: minmax(0, 0.9fr) minmax(22rem, 1.25fr);
  align-items: stretch;
}

.channel-monitor-card__summary {
  display: grid;
  align-content: start;
  gap: 0.55rem;
  min-width: 0;
}

.channel-monitor-card__summary h3 {
  margin: 0;
  color: #edf7ee;
  font-size: 1.45rem;
  font-weight: 900;
  line-height: 1.18;
}

.channel-monitor-card__summary p:not(.channel-kicker) {
  margin: 0;
  max-width: 42rem;
  color: rgba(220, 239, 224, 0.66);
  font-size: 0.9rem;
  line-height: 1.65;
}

.channel-monitor-card__model {
  display: inline-grid;
  width: fit-content;
  min-width: 8.5rem;
  gap: 0.28rem;
  margin-top: 0.2rem;
  border: 1px solid rgba(43, 132, 83, 0.32);
  border-radius: 8px;
  background: rgba(3, 22, 18, 0.58);
  padding: 0.62rem 0.74rem;
}

.channel-monitor-card__model span,
.channel-monitor-card__value span {
  color: rgba(220, 239, 224, 0.62);
  font-size: 0.72rem;
  font-weight: 850;
}

.channel-monitor-card__model strong {
  color: #edf7ee;
  font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
  font-size: 0.86rem;
}

.channel-monitor-card__value {
  position: relative;
  display: grid;
  align-content: center;
  min-height: 10.5rem;
  overflow: hidden;
  border: 1px solid rgba(90, 201, 148, 0.28);
  border-radius: 8px;
  background:
    radial-gradient(circle at 88% 18%, rgba(234, 179, 8, 0.2), transparent 24%),
    linear-gradient(135deg, rgba(27, 122, 78, 0.34), rgba(3, 22, 18, 0.62) 58%);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.06), 0 18px 48px rgba(0, 0, 0, 0.18);
  padding: 1.2rem 1.25rem;
}

.channel-monitor-card__value::after {
  position: absolute;
  right: -2.5rem;
  bottom: -3.8rem;
  width: 12rem;
  height: 12rem;
  border: 1px solid rgba(234, 179, 8, 0.18);
  border-radius: 999px;
  background: rgba(126, 227, 155, 0.06);
  content: "";
}

.channel-monitor-card__value strong {
  max-width: 14ch;
  margin-top: 0.52rem;
  color: #f4fff5;
  font-size: 3.65rem;
  font-weight: 950;
  line-height: 0.98;
  text-shadow: 0 0 28px rgba(126, 227, 155, 0.18);
}

.channel-monitor-card__value em {
  position: absolute;
  top: 1rem;
  right: 1rem;
  border: 1px solid rgba(234, 179, 8, 0.32);
  border-radius: 999px;
  background: rgba(234, 179, 8, 0.1);
  color: #f3d56b;
  font-style: normal;
  font-size: 0.86rem;
  font-weight: 900;
  padding: 0.3rem 0.58rem;
}

.channel-monitor-card__metrics {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 0.65rem;
}

.channel-monitor-card__note {
  justify-self: end;
  max-width: 42rem;
  margin: -0.2rem 0 0;
  color: rgba(220, 239, 224, 0.45);
  font-size: 0.68rem;
  font-weight: 650;
  line-height: 1.45;
  text-align: right;
}

.channel-hero {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 1rem;
  padding: clamp(0.95rem, 1.8vw, 1.45rem);
}

.channel-group-divider {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  margin-top: 0.35rem;
  padding: 1rem;
  background:
    linear-gradient(135deg, rgba(34, 197, 94, 0.08), rgba(5, 29, 24, 0.9) 48%),
    rgba(5, 29, 24, 0.78);
}

.channel-group-divider h2 {
  margin: 0.22rem 0 0;
  color: #edf7ee;
  font-size: 1.18rem;
  font-weight: 880;
}

.channel-group-divider p:not(.channel-kicker) {
  margin: 0.36rem 0 0;
  color: rgba(220, 239, 224, 0.68);
  font-size: 0.86rem;
  line-height: 1.55;
}

.channel-hero__copy {
  max-width: 52rem;
}

.channel-kicker,
.channel-overview__meta p,
.channel-section-title p {
  margin: 0;
  color: rgba(220, 239, 224, 0.68);
  font-size: 0.72rem;
  font-weight: 850;
  letter-spacing: 0;
  text-transform: uppercase;
}

.channel-hero h1 {
  margin: 0.32rem 0 0;
  color: #edf7ee;
  font-size: clamp(1.85rem, 3.25vw, 2.85rem);
  font-weight: 900;
  line-height: 0.96;
}

.channel-hero p:not(.channel-kicker),
.channel-summary-card p {
  margin: 0.8rem 0 0;
  max-width: 47rem;
  color: rgba(220, 239, 224, 0.68);
  font-size: 0.95rem;
  line-height: 1.75;
}

.channel-hero__status {
  display: flex;
  flex-shrink: 0;
  align-items: center;
  gap: 0.65rem;
}

.channel-status-chip {
  display: inline-flex;
  align-items: center;
  gap: 0.48rem;
  min-height: 2.15rem;
  border: 1px solid transparent;
  border-radius: 999px;
  padding: 0 0.85rem;
  font-size: 0.84rem;
  font-weight: 850;
}

.channel-status-chip__dot,
:deep(.channel-legend__dot) {
  display: inline-block;
  width: 0.5rem;
  height: 0.5rem;
  border-radius: 999px;
}

.channel-status-chip--ok {
  border-color: rgba(34, 197, 94, 0.32);
  background: rgba(34, 197, 94, 0.12);
  color: #7ee39b;
}

.channel-status-chip--warn {
  border-color: rgba(34, 197, 94, 0.22);
  background: rgba(34, 197, 94, 0.08);
  color: #edf7ee;
}

.channel-status-chip--bad {
  border-color: rgba(34, 197, 94, 0.22);
  background: rgba(34, 197, 94, 0.08);
  color: rgba(220, 239, 224, 0.72);
}

.channel-status-chip--unknown {
  border-color: rgba(43, 132, 83, 0.38);
  background: rgba(34, 197, 94, 0.12);
  color: rgba(220, 239, 224, 0.62);
}

.channel-refresh {
  display: inline-flex;
  width: 2.15rem;
  height: 2.15rem;
  align-items: center;
  justify-content: center;
  border: 1px solid rgba(43, 132, 83, 0.38);
  border-radius: 8px;
  background: rgba(3, 22, 18, 0.64);
  color: #edf7ee;
  transition: background-color 0.16s ease, border-color 0.16s ease, transform 0.16s ease;
}

.channel-refresh:hover:not(:disabled) {
  border-color: rgba(34, 197, 94, 0.52);
  background: rgba(34, 197, 94, 0.12);
  transform: translateY(-1px);
}

.channel-refresh:disabled {
  opacity: 0.58;
}

.channel-summary-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  padding: 0.9rem 1rem;
}

.channel-summary-card h2 {
  margin: 0.25rem 0 0;
  color: #edf7ee;
  font-size: clamp(1.2rem, 2vw, 1.65rem);
  font-weight: 880;
}

.channel-summary-card--ok {
  background: linear-gradient(135deg, rgba(34, 197, 94, 0.1), rgba(5, 29, 24, 0.92) 48%);
}

.channel-summary-card--warn {
  background: linear-gradient(135deg, rgba(34, 197, 94, 0.08), rgba(5, 29, 24, 0.92) 48%);
}

.channel-summary-card--bad {
  background: linear-gradient(135deg, rgba(34, 197, 94, 0.08), rgba(5, 29, 24, 0.92) 48%);
}

.channel-summary-card__meta {
  display: grid;
  flex: 0 0 auto;
  gap: 0.3rem;
  min-width: 8.5rem;
  border: 1px solid rgba(43, 132, 83, 0.38);
  border-radius: 8px;
  background: rgba(3, 22, 18, 0.64);
  padding: 0.75rem 0.85rem;
}

.channel-summary-card__meta span {
  color: rgba(220, 239, 224, 0.68);
  font-size: 0.75rem;
  font-weight: 750;
}

.channel-summary-card__meta strong {
  color: #edf7ee;
  font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
  font-size: 0.88rem;
}

.channel-overview {
  display: grid;
  gap: 0.8rem;
  padding: 0.9rem 1rem;
}

.channel-overview__meta,
.channel-section-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
}

.channel-overview__meta h2,
.channel-section-title h2 {
  margin: 0.22rem 0 0;
  color: #edf7ee;
  font-size: 1.15rem;
  font-weight: 850;
}

.channel-refresh-meta {
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-end;
  gap: 0.45rem;
}

.channel-refresh-meta span,
.channel-section-title > span {
  border: 1px solid rgba(43, 132, 83, 0.28);
  border-radius: 999px;
  background: rgba(3, 22, 18, 0.64);
  padding: 0.35rem 0.6rem;
  color: rgba(220, 239, 224, 0.68);
  font-size: 0.78rem;
  font-weight: 700;
}

.channel-metrics {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(10rem, 1fr));
  gap: 0.65rem;
}

:deep(.channel-metric) {
  position: relative;
  min-width: 0;
  overflow: hidden;
  border: 1px solid rgba(43, 132, 83, 0.28);
  border-radius: 8px;
  background: rgba(3, 22, 18, 0.64);
  padding: 0.78rem 0.85rem;
}

:deep(.channel-metric::before) {
  position: absolute;
  inset: 0 auto 0 0;
  width: 3px;
  background: rgba(43, 132, 83, 0.32);
  content: "";
}

:deep(.channel-metric--ok::before) {
  background: #2f9d57;
}

:deep(.channel-metric--warn::before) {
  background: #d4a017;
}

:deep(.channel-metric--bad::before) {
  background: #d94a3a;
}

:deep(.channel-metric--rate) {
  border-color: rgba(34, 197, 94, 0.34);
  background:
    linear-gradient(135deg, rgba(34, 197, 94, 0.12), rgba(3, 22, 18, 0.64) 52%),
    rgba(3, 22, 18, 0.64);
}

:deep(.channel-metric--rate::before) {
  background: #11a36a;
}

:deep(.channel-metric--has-tooltip) {
  overflow: visible;
}

:deep(.channel-metric--has-tooltip:hover) {
  z-index: 10;
}

:deep(.channel-metric--has-tooltip:hover::after) {
  position: absolute;
  left: 0.85rem;
  bottom: calc(100% + 0.55rem);
  z-index: 25;
  width: min(22rem, calc(100vw - 3rem));
  border: 1px solid rgba(90, 201, 148, 0.24);
  border-radius: 7px;
  background: rgba(17, 27, 22, 0.96);
  box-shadow: 0 14px 32px rgba(0, 0, 0, 0.28);
  color: #edf7ee;
  content: attr(data-tooltip);
  font-size: 0.74rem;
  font-weight: 650;
  line-height: 1.55;
  padding: 0.58rem 0.68rem;
  pointer-events: none;
  white-space: pre-line;
}

:deep(.channel-metric--has-tooltip:hover::before) {
  z-index: 1;
}

:deep(.channel-metric__label) {
  display: block;
  color: rgba(220, 239, 224, 0.68);
  font-size: 0.76rem;
  font-weight: 800;
}

:deep(.channel-metric strong) {
  display: block;
  margin-top: 0.38rem;
  overflow: hidden;
  color: #edf7ee;
  font-size: clamp(1.08rem, 1.8vw, 1.55rem);
  font-weight: 900;
  line-height: 1.05;
  text-overflow: ellipsis;
  white-space: nowrap;
}

:deep(.channel-metric__hint) {
  display: flex;
  align-items: center;
  gap: 0.35rem;
  margin-top: 0.38rem;
  color: rgba(220, 239, 224, 0.58);
  font-size: 0.74rem;
  font-weight: 650;
}

:deep(.channel-metric__info) {
  display: inline-flex;
  width: 1rem;
  height: 1rem;
  flex: 0 0 auto;
  align-items: center;
  justify-content: center;
  border: 1px solid rgba(90, 201, 148, 0.54);
  border-radius: 999px;
  color: #7ee39b;
  font-size: 0.68rem;
  font-weight: 900;
  line-height: 1;
}

.channel-timeline-panel {
  padding: 0.9rem 1rem;
}

.channel-timeline-frame {
  margin-top: 0.8rem;
  border: 1px solid rgba(43, 132, 83, 0.28);
  border-radius: 8px;
  background: rgba(3, 22, 18, 0.64);
  padding: 0.7rem;
}

.channel-timeline-axis {
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.5rem;
  color: rgba(220, 239, 224, 0.62);
  font-size: 0.72rem;
  font-weight: 750;
}

.channel-timeline {
  display: flex;
  align-items: flex-end;
  gap: 0.35rem;
  height: 4.8rem;
  border-radius: 6px;
  background:
    linear-gradient(rgba(34, 197, 94, 0.055) 1px, transparent 1px),
    rgba(3, 22, 18, 0.64);
  background-size: 100% 1.33rem;
  padding: 0.65rem;
}

.channel-timeline__bar {
  position: relative;
  min-width: 0.24rem;
  flex: 1;
  border-radius: 3px 3px 1px 1px;
  transition: height 0.2s ease, opacity 0.2s ease, transform 0.16s ease;
}

.channel-timeline__bar:hover {
  z-index: 5;
  transform: translateY(-2px);
}

.channel-timeline__bar:hover::after {
  position: absolute;
  left: 50%;
  bottom: calc(100% + 0.55rem);
  z-index: 20;
  width: max-content;
  max-width: 13rem;
  transform: translateX(-50%);
  border: 1px solid rgba(20, 24, 20, 0.16);
  border-radius: 7px;
  background: rgba(17, 22, 17, 0.94);
  box-shadow: 0 10px 22px rgba(17, 22, 17, 0.18);
  color: #edf7ee;
  content: attr(data-tooltip);
  font-size: 0.72rem;
  font-weight: 750;
  line-height: 1.45;
  padding: 0.42rem 0.55rem;
  pointer-events: none;
  text-align: left;
  white-space: pre-line;
}

.channel-timeline__bar:hover::before {
  position: absolute;
  left: 50%;
  bottom: calc(100% + 0.25rem);
  z-index: 19;
  width: 0.52rem;
  height: 0.52rem;
  transform: translateX(-50%) rotate(45deg);
  background: rgba(17, 22, 17, 0.94);
  content: "";
  pointer-events: none;
}

.channel-timeline__bar--ok,
.channel-dot--ok {
  background: #2f9d57;
}

.channel-timeline__bar--warn,
.channel-dot--warn {
  background: #d4a017;
}

.channel-timeline__bar--bad,
.channel-dot--bad {
  background: #d94a3a;
}

.channel-timeline__bar--unknown,
.channel-dot--unknown {
  background: rgba(43, 132, 83, 0.32);
}

.channel-legend {
  display: flex;
  flex-wrap: wrap;
  gap: 0.8rem 1rem;
  margin-top: 0.9rem;
  color: rgba(220, 239, 224, 0.62);
  font-size: 0.8rem;
  font-weight: 700;
}

:deep(.channel-legend__item) {
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
}

.channel-guidance {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 0.75rem;
}

.channel-guidance article {
  display: flex;
  min-width: 0;
  gap: 0.75rem;
  padding: 0.95rem;
}

.channel-guidance__icon {
  display: inline-flex;
  flex: 0 0 auto;
  width: 2rem;
  height: 2rem;
  align-items: center;
  justify-content: center;
  border: 1px solid rgba(43, 132, 83, 0.38);
  border-radius: 8px;
  background: rgba(3, 22, 18, 0.64);
  color: #22c55e;
}

.channel-guidance h3 {
  margin: 0;
  color: #edf7ee;
  font-size: 0.9rem;
  font-weight: 850;
}

.channel-guidance p {
  margin: 0.28rem 0 0;
  color: rgba(220, 239, 224, 0.62);
  font-size: 0.84rem;
  line-height: 1.55;
}

@media (max-width: 980px) {
  .channel-monitor-card__main {
    grid-template-columns: 1fr;
  }

  .channel-monitor-card__metrics {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .channel-metrics {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .channel-guidance {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 640px) {
  .channel-hero,
  .channel-group-divider,
  .channel-summary-card,
  .channel-monitor-card__head,
  .channel-monitor-card__status,
  .channel-overview__meta,
  .channel-section-title {
    align-items: flex-start;
    flex-direction: column;
  }

  .channel-hero__status,
  .channel-group-divider > .channel-status-chip,
  .channel-summary-card__meta {
    width: 100%;
  }

  .channel-hero__status {
    justify-content: space-between;
  }

  .channel-refresh-meta {
    justify-content: flex-start;
  }

  .channel-monitor-card__value {
    min-height: 8.6rem;
    padding: 1rem;
  }

  .channel-monitor-card__value strong {
    max-width: 12ch;
    font-size: 2.15rem;
  }

  .channel-monitor-card__metrics {
    grid-template-columns: 1fr;
  }

  .channel-monitor-card__note {
    justify-self: start;
    text-align: left;
  }

  .channel-metrics {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .channel-timeline {
    gap: 0.2rem;
    padding: 0.5rem;
  }
}
</style>
