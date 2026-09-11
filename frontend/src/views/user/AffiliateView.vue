<template>
  <AppLayout>
    <main class="affiliate-page">
      <section class="affiliate-hero">
        <div class="affiliate-hero__copy">
          <p class="affiliate-kicker">AFFILIATE REBATE</p>
          <h1>邀请返利</h1>
          <div class="affiliate-hero__actions">
            <button
              type="button"
              class="affiliate-button affiliate-button--primary"
              :disabled="!detail?.aff_code"
              @click="copyInviteLink"
            >
              <Icon name="copy" size="sm" />
              复制邀请链接
            </button>
            <button
              type="button"
              class="affiliate-button affiliate-button--secondary"
              :disabled="!detail?.aff_code"
              @click="copyCode"
            >
              复制邀请码
            </button>
          </div>
        </div>

        <div class="affiliate-hero__rate">
          <span>当前返利比例</span>
          <strong>{{ formattedRebateRate }}%</strong>
          <p>按被邀请用户实际产生的有效消费计算。</p>
        </div>
      </section>

      <div v-if="loading" class="affiliate-loading">
        <Icon name="refresh" size="lg" class="animate-spin" />
        <span>正在加载邀请返利数据</span>
      </div>

      <template v-else-if="detail">
        <section class="affiliate-stats" aria-label="返利概览">
          <article class="affiliate-stat affiliate-stat--gold">
            <span>可转返利</span>
            <strong>{{ formatCurrency(detail.aff_quota) }}</strong>
            <p>可立即转入余额</p>
          </article>
          <article class="affiliate-stat">
            <span>已邀请用户</span>
            <strong>{{ formatCount(detail.aff_count) }}</strong>
            <p>通过你的链接注册</p>
          </article>
          <article class="affiliate-stat">
            <span>历史返利</span>
            <strong>{{ formatCurrency(detail.aff_history_quota) }}</strong>
            <p>累计产生的返利额度</p>
          </article>
        </section>

        <section class="affiliate-grid">
          <article class="affiliate-panel affiliate-panel--share">
            <div class="affiliate-panel__head">
              <div>
                <p class="affiliate-kicker">SHARE</p>
                <h2>你的专属邀请信息</h2>
              </div>
            </div>

            <div class="affiliate-copy-card">
              <label>邀请码</label>
              <div>
                <code>{{ detail.aff_code }}</code>
                <button type="button" @click="copyCode">
                  <Icon name="copy" size="sm" />
                  复制
                </button>
              </div>
            </div>

            <div class="affiliate-copy-card">
              <label>邀请链接</label>
              <div>
                <code>{{ inviteLink }}</code>
                <button type="button" @click="copyInviteLink">
                  <Icon name="copy" size="sm" />
                  复制
                </button>
              </div>
            </div>

            <ol class="affiliate-rules">
              <li><span>01</span><p>对方通过你的邀请链接或邀请码注册后，会绑定为你的邀请用户。</p></li>
              <li><span>02</span><p>对方后续充值余额、购买订阅或兑换订阅卡密时，你按比例获得返利。</p></li>
              <li><span>03</span><p>返利先进入可转返利额度，点击转入余额后即可用于站内消费。</p></li>
            </ol>
          </article>

          <article class="affiliate-panel affiliate-panel--transfer">
            <div class="affiliate-panel__head">
              <div>
                <p class="affiliate-kicker">BALANCE</p>
                <h2>返利转余额</h2>
              </div>
            </div>
            <p class="affiliate-transfer__amount">{{ formatCurrency(detail.aff_quota) }}</p>
            <p class="affiliate-transfer__desc">
              转入后会和普通余额合并，可直接用于 API 调用或站内购买。
            </p>
            <button
              type="button"
              class="affiliate-button affiliate-button--primary affiliate-button--full"
              :disabled="transferring || detail.aff_quota <= 0"
              @click="transferQuota"
            >
              <Icon v-if="transferring" name="refresh" size="sm" class="animate-spin" />
              <Icon v-else name="dollar" size="sm" />
              {{ transferring ? '转入中' : '转入余额' }}
            </button>
            <p v-if="detail.aff_quota <= 0" class="affiliate-transfer__empty">当前没有可转入的返利额度。</p>
          </article>
        </section>

        <section class="affiliate-panel affiliate-panel--records">
          <div class="affiliate-panel__head">
            <div>
              <p class="affiliate-kicker">RECORDS</p>
              <h2>邀请记录</h2>
            </div>
            <span>{{ invitees.length }} 人</span>
          </div>

          <div v-if="invitees.length === 0" class="affiliate-empty">
            暂无邀请记录。复制邀请链接发给朋友，对方注册并消费后这里会显示返利明细。
          </div>

          <div v-else class="affiliate-table-wrap">
            <table class="affiliate-table">
              <thead>
                <tr>
                  <th>用户</th>
                  <th>邮箱</th>
                  <th class="align-right">累计返利</th>
                  <th>注册时间</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="item in invitees" :key="item.user_id">
                  <td>{{ item.username || '-' }}</td>
                  <td>{{ item.email || '-' }}</td>
                  <td class="align-right affiliate-table__money">{{ formatCurrency(item.total_rebate) }}</td>
                  <td>{{ formatDateTime(item.created_at) || '-' }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </section>
      </template>
    </main>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import userAPI from '@/api/user'
import type { AffiliateInvitee, UserAffiliateDetail } from '@/types'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'
import { useClipboard } from '@/composables/useClipboard'
import { formatCurrency, formatDateTime } from '@/utils/format'
import { extractApiErrorMessage } from '@/utils/apiError'

const appStore = useAppStore()
const authStore = useAuthStore()
const { copyToClipboard } = useClipboard()

const loading = ref(true)
const transferring = ref(false)
const detail = ref<UserAffiliateDetail | null>(null)

const invitees = computed<AffiliateInvitee[]>(() => detail.value?.invitees ?? [])

const inviteLink = computed(() => {
  const code = detail.value?.aff_code || ''
  if (!code) return ''
  if (typeof window === 'undefined') return `/register?aff=${encodeURIComponent(code)}`
  return `${window.location.origin}/register?aff=${encodeURIComponent(code)}`
})

const formattedRebateRate = computed(() => {
  const value = detail.value?.effective_rebate_rate_percent ?? 0
  const rounded = Math.round(value * 100) / 100
  return Number.isInteger(rounded) ? String(rounded) : rounded.toString()
})

function formatCount(value: number): string {
  return value.toLocaleString()
}

async function loadAffiliateDetail(silent = false): Promise<void> {
  if (!silent) loading.value = true
  try {
    detail.value = await userAPI.getAffiliateDetail()
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, '加载邀请返利数据失败'))
  } finally {
    if (!silent) loading.value = false
  }
}

async function copyCode(): Promise<void> {
  if (!detail.value?.aff_code) return
  await copyToClipboard(detail.value.aff_code, '邀请码已复制')
}

async function copyInviteLink(): Promise<void> {
  if (!inviteLink.value) return
  await copyToClipboard(inviteLink.value, '邀请链接已复制')
}

async function transferQuota(): Promise<void> {
  if (!detail.value || detail.value.aff_quota <= 0 || transferring.value) return
  transferring.value = true
  try {
    const resp = await userAPI.transferAffiliateQuota()
    appStore.showSuccess(`已转入余额：${formatCurrency(resp.transferred_quota)}`)
    await Promise.all([
      loadAffiliateDetail(true),
      authStore.refreshUser().catch(() => undefined),
    ])
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, '转入余额失败'))
  } finally {
    transferring.value = false
  }
}

onMounted(() => {
  void loadAffiliateDetail()
})
</script>

<style scoped>
.affiliate-page {
  display: grid;
  gap: 1rem;
}

.affiliate-hero,
.affiliate-panel,
.affiliate-stat {
  position: relative;
  overflow: hidden;
  border: 1px solid rgba(43, 132, 83, 0.38);
  border-radius: 8px;
  background:
    linear-gradient(145deg, rgba(5, 29, 24, 0.78), rgba(4, 24, 20, 0.88)),
    radial-gradient(circle at 92% 0%, rgba(34, 197, 94, 0.08), transparent 17rem);
  box-shadow:
    0 18px 54px rgba(0, 0, 0, 0.22),
    inset 0 1px 0 rgba(34, 197, 94, 0.08);
}

.affiliate-hero::before,
.affiliate-panel::before,
.affiliate-stat::before {
  content: '';
  position: absolute;
  inset: 0;
  pointer-events: none;
  background:
    radial-gradient(circle at 14% 22%, rgba(34, 197, 94, 0.08) 0 1px, transparent 1px 8px),
    radial-gradient(circle at 82% 72%, rgba(43, 132, 83, 0.08) 0 1px, transparent 1px 9px),
    linear-gradient(90deg, rgba(255, 255, 255, 0.035), transparent 22%, transparent 78%, rgba(34, 197, 94, 0.03));
  background-size: 31px 31px, 43px 43px, 100% 100%, 100% 100%;
  opacity: 0.7;
}

.affiliate-hero > *,
.affiliate-panel > *,
.affiliate-stat > * {
  position: relative;
  z-index: 1;
}

.affiliate-hero {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(17rem, 0.32fr);
  gap: 1rem;
  min-height: 17rem;
  padding: clamp(1.1rem, 2.4vw, 1.8rem);
}

.affiliate-hero__copy {
  display: flex;
  flex-direction: column;
  justify-content: center;
  max-width: 48rem;
}

.affiliate-kicker {
  margin: 0 0 0.65rem;
  color: #22c55e;
  font-size: 0.72rem;
  font-weight: 850;
  letter-spacing: 0.12em;
}

.affiliate-hero h1 {
  margin: 0;
  color: #edf7ee;
  font-size: clamp(2.2rem, 5vw, 4.5rem);
  font-weight: 950;
  line-height: 0.95;
  letter-spacing: -0.04em;
}

.affiliate-hero__actions {
  display: flex;
  flex-wrap: wrap;
  gap: 0.65rem;
  margin-top: 1.4rem;
}

.affiliate-hero__rate {
  display: flex;
  min-height: 13rem;
  flex-direction: column;
  justify-content: center;
  border: 1px solid rgba(43, 132, 83, 0.38);
  border-radius: 8px;
  background:
    radial-gradient(circle at 100% 0%, rgba(34, 197, 94, 0.12), transparent 12rem),
    linear-gradient(145deg, rgba(3, 22, 18, 0.64), rgba(5, 29, 24, 0.82));
  padding: 1.2rem;
}

.affiliate-hero__rate span,
.affiliate-stat span {
  color: #22c55e;
  font-size: 0.78rem;
  font-weight: 800;
}

.affiliate-hero__rate strong {
  margin-top: 0.4rem;
  color: #22c55e;
  font-size: clamp(3rem, 7vw, 5rem);
  font-weight: 950;
  line-height: 1;
}

.affiliate-hero__rate p,
.affiliate-stat p,
.affiliate-transfer__desc,
.affiliate-transfer__empty {
  color: rgba(220, 239, 224, 0.68);
  font-size: 0.86rem;
  line-height: 1.7;
}

.affiliate-button {
  display: inline-flex;
  min-height: 2.55rem;
  align-items: center;
  justify-content: center;
  gap: 0.45rem;
  border-radius: 8px;
  padding: 0 1rem;
  font-size: 0.88rem;
  font-weight: 750;
  text-decoration: none;
  transition:
    transform 0.16s ease,
    border-color 0.16s ease,
    background-color 0.16s ease;
}

.affiliate-button:not(:disabled):hover {
  transform: translateY(-1px);
}

.affiliate-button:disabled {
  cursor: not-allowed;
  opacity: 0.52;
}

.affiliate-button--primary {
  border: 1px solid #27643f;
  background: #1f6b3b;
  color: #ffffff;
  box-shadow: 0 12px 24px rgba(31, 107, 59, 0.15);
}

.affiliate-button--primary:not(:disabled):hover {
  border-color: #184b2d;
  background: #18542f;
}

.affiliate-button--secondary {
  border: 1px solid rgba(43, 132, 83, 0.38);
  background: rgba(3, 22, 18, 0.64);
  color: #edf7ee;
}

.affiliate-button--secondary:not(:disabled):hover {
  border-color: rgba(34, 197, 94, 0.52);
  background: rgba(34, 197, 94, 0.12);
}

.affiliate-button--full {
  width: 100%;
  margin-top: 1rem;
}

.affiliate-loading {
  display: flex;
  min-height: 14rem;
  align-items: center;
  justify-content: center;
  gap: 0.7rem;
  border: 1px solid rgba(43, 132, 83, 0.38);
  border-radius: 8px;
  background: rgba(5, 29, 24, 0.78);
  color: rgba(220, 239, 224, 0.68);
  font-weight: 700;
}

.affiliate-stats {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 1rem;
}

.affiliate-stat {
  min-height: 8.2rem;
  padding: 1rem;
}

.affiliate-stat strong {
  display: block;
  margin-top: 0.55rem;
  color: #edf7ee;
  font-size: clamp(1.55rem, 3vw, 2.2rem);
  font-weight: 900;
  line-height: 1.05;
}

.affiliate-stat--gold {
  border-color: rgba(34, 197, 94, 0.38);
  background:
    radial-gradient(circle at 88% 0%, rgba(34, 197, 94, 0.12), transparent 12rem),
    linear-gradient(145deg, rgba(5, 29, 24, 0.82), rgba(4, 24, 20, 0.9));
}

.affiliate-stat--warning strong {
  color: #7ee39b;
}

.affiliate-grid {
  display: grid;
  grid-template-columns: minmax(0, 1.45fr) minmax(18rem, 0.55fr);
  gap: 1rem;
}

.affiliate-panel {
  padding: clamp(1rem, 2vw, 1.35rem);
}

.affiliate-panel__head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 1rem;
  margin-bottom: 1rem;
}

.affiliate-panel__head h2 {
  margin: 0;
  color: #edf7ee;
  font-size: 1.15rem;
  font-weight: 850;
}

.affiliate-panel__head > span {
  border: 1px solid rgba(43, 132, 83, 0.38);
  border-radius: 999px;
  background: rgba(3, 22, 18, 0.64);
  padding: 0.35rem 0.65rem;
  color: rgba(220, 239, 224, 0.68);
  font-size: 0.78rem;
  font-weight: 750;
}

.affiliate-copy-card {
  display: grid;
  gap: 0.45rem;
  margin-top: 0.75rem;
}

.affiliate-copy-card label {
  color: rgba(220, 239, 224, 0.68);
  font-size: 0.82rem;
  font-weight: 750;
}

.affiliate-copy-card > div {
  display: flex;
  align-items: center;
  gap: 0.55rem;
  border: 1px solid rgba(43, 132, 83, 0.38);
  border-radius: 8px;
  background: rgba(3, 22, 18, 0.64);
  padding: 0.5rem 0.55rem 0.5rem 0.75rem;
}

.affiliate-copy-card code {
  min-width: 0;
  flex: 1;
  overflow: hidden;
  color: #edf7ee;
  font-size: 0.88rem;
  font-weight: 800;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.affiliate-copy-card button {
  display: inline-flex;
  min-height: 2rem;
  align-items: center;
  gap: 0.35rem;
  border: 1px solid rgba(43, 132, 83, 0.38);
  border-radius: 7px;
  background: rgba(3, 22, 18, 0.64);
  padding: 0 0.7rem;
  color: #edf7ee;
  font-size: 0.8rem;
  font-weight: 750;
}

.affiliate-copy-card button:hover {
  border-color: rgba(34, 197, 94, 0.52);
  background: rgba(34, 197, 94, 0.12);
}

.affiliate-rules {
  display: grid;
  gap: 0.65rem;
  margin: 1rem 0 0;
  padding: 0;
  list-style: none;
}

.affiliate-rules li {
  display: grid;
  grid-template-columns: 2.4rem minmax(0, 1fr);
  gap: 0.75rem;
  align-items: flex-start;
  border: 1px solid rgba(43, 132, 83, 0.32);
  border-radius: 8px;
  background: rgba(3, 22, 18, 0.64);
  padding: 0.75rem;
}

.affiliate-rules span {
  color: #22c55e;
  font-size: 0.72rem;
  font-weight: 900;
}

.affiliate-rules p {
  margin: 0;
  color: rgba(220, 239, 224, 0.68);
  font-size: 0.86rem;
  line-height: 1.65;
}

.affiliate-transfer__amount {
  margin: 0.2rem 0 0;
  color: #22c55e;
  font-size: clamp(2.1rem, 5vw, 3.4rem);
  font-weight: 950;
  letter-spacing: -0.03em;
}

.affiliate-transfer__desc {
  margin: 0.75rem 0 0;
}

.affiliate-transfer__empty {
  margin: 0.7rem 0 0;
  color: #7ee39b;
}

.affiliate-empty {
  border: 1px dashed rgba(43, 132, 83, 0.38);
  border-radius: 8px;
  background: rgba(3, 22, 18, 0.64);
  padding: 2.5rem 1rem;
  color: rgba(220, 239, 224, 0.68);
  text-align: center;
  font-size: 0.92rem;
}

.affiliate-table-wrap {
  overflow-x: auto;
}

.affiliate-table {
  width: 100%;
  min-width: 680px;
  border-collapse: collapse;
  font-size: 0.88rem;
}

.affiliate-table th,
.affiliate-table td {
  border-bottom: 1px solid rgba(43, 132, 83, 0.28);
  padding: 0.8rem 0.7rem;
  color: rgba(220, 239, 224, 0.68);
  text-align: left;
}

.affiliate-table th {
  color: rgba(220, 239, 224, 0.56);
  font-size: 0.78rem;
  font-weight: 800;
}

.affiliate-table tbody tr:last-child td {
  border-bottom: 0;
}

.affiliate-table tbody td:first-child {
  color: #edf7ee;
  font-weight: 750;
}

.affiliate-table .align-right {
  text-align: right;
}

.affiliate-table__money {
  color: #1f6b3b;
  font-weight: 850;
}

@media (max-width: 980px) {
  .affiliate-hero,
  .affiliate-grid {
    grid-template-columns: 1fr;
  }

  .affiliate-stats {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 640px) {
  .affiliate-stats {
    grid-template-columns: 1fr;
  }

  .affiliate-hero {
    min-height: auto;
  }

  .affiliate-hero__rate {
    min-height: 10rem;
  }

  .affiliate-button {
    width: 100%;
  }

  .affiliate-copy-card > div {
    align-items: stretch;
    flex-direction: column;
  }

  .affiliate-copy-card button {
    justify-content: center;
  }
}
</style>
