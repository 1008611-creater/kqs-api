<template>
  <AppLayout>
    <main class="subscription-shop">
      <section class="subscription-shop__hero">
        <div class="subscription-shop__intro">
          <p class="subscription-shop__eyebrow">订阅服务</p>
          <h1>订阅套餐</h1>
          <div class="subscription-shop__hero-actions">
            <a class="subscription-shop__primary" :href="LDXP_SHOP_URL" target="_blank" rel="noopener noreferrer" @click="trackOpenShop">
              <Icon name="externalLink" size="sm" />
              打开链动小铺
            </a>
            <RouterLink class="subscription-shop__secondary" to="/purchase">
              购买余额
            </RouterLink>
          </div>
        </div>

        <form class="subscription-shop__redeem" @submit.prevent="handleRedeem">
          <div>
            <h2>兑换订阅卡密</h2>
            <p>在链动小铺购买后，复制订单里的卡密，在这里兑换开通。</p>
          </div>
          <label class="subscription-shop__field">
            <span>卡密</span>
            <input
              v-model="redeemCode"
              :disabled="redeeming"
              autocomplete="off"
              placeholder="粘贴链动小铺卡密"
            >
          </label>
          <button class="subscription-shop__primary" type="submit" :disabled="!redeemCode.trim() || redeeming">
            <span v-if="redeeming">正在兑换</span>
            <span v-else>立即兑换</span>
          </button>
          <p v-if="redeemError" class="subscription-shop__message subscription-shop__message--error">
            {{ redeemError }}
          </p>
          <p v-else-if="redeemResult" class="subscription-shop__message subscription-shop__message--success">
            {{ redeemResult.message || '兑换成功，订阅已刷新' }}
          </p>
        </form>
      </section>

      <section class="subscription-shop__layout">
        <div class="subscription-shop__plans">
          <div class="subscription-shop__section-head">
            <div>
              <h2>选择订阅</h2>
              <p>日卡按当天额度使用，周卡每天刷新一次每日额度。</p>
            </div>
            <div class="subscription-shop__tabs" aria-label="套餐筛选">
              <button
                v-for="tab in planTabs"
                :key="tab.value"
                type="button"
                :class="{ active: planTab === tab.value }"
                @click="planTab = tab.value"
              >
                {{ tab.label }}
              </button>
            </div>
          </div>

          <div v-if="plansLoading" class="subscription-shop__loading">
            <LoadingSpinner size="md" />
          </div>
          <div v-else-if="visiblePlans.length" class="subscription-shop__plan-grid">
            <article v-for="plan in visiblePlans" :key="plan.key" class="subscription-plan">
              <div class="subscription-plan__top">
                <div>
                  <h3>{{ plan.name }}</h3>
                  <p>{{ plan.validity_days > 1 ? '周卡：每天刷新额度' : '日卡：当天额度' }}</p>
                </div>
                <span>{{ formatValidity(plan.validity_days) }}</span>
              </div>

              <div class="subscription-plan__price">
                <strong>¥{{ formatPrice(plan.price) }}</strong>
                <span>链动小铺卡密</span>
              </div>

              <dl class="subscription-plan__metrics">
                <div>
                  <dt>每日额度</dt>
                  <dd>${{ formatQuota(plan.daily_limit_usd) }}</dd>
                </div>
                <div v-if="plan.validity_days > 1">
                  <dt>周期额度</dt>
                  <dd>${{ formatQuota(plan.weekly_limit_usd) }}</dd>
                </div>
                <div>
                  <dt>订阅分组</dt>
                  <dd>{{ displayGroupName(plan) }}</dd>
                </div>
                <div>
                  <dt>计费倍率</dt>
                  <dd>{{ formatRate(plan.rate_multiplier) }}</dd>
                </div>
              </dl>

              <a
                class="subscription-shop__secondary subscription-shop__secondary--full"
                :href="planPurchaseUrl(plan)"
                target="_blank"
                rel="noopener noreferrer"
                @click="trackPlanPurchase(plan)"
              >
                去购买
                <Icon name="arrowRight" size="xs" />
              </a>
            </article>
          </div>
          <div v-else class="subscription-shop__empty">
            暂无可购买订阅套餐。
          </div>
        </div>

        <aside class="subscription-shop__side">
          <section class="subscription-panel">
            <div class="subscription-panel__head">
              <div>
                <h2>当前订阅</h2>
                <p>兑换成功后会自动显示在这里。</p>
              </div>
              <button type="button" aria-label="刷新订阅" @click="refreshAll">
                <Icon name="refresh" size="sm" />
              </button>
            </div>

            <div v-if="subscriptionStore.loading" class="subscription-shop__side-loading">
              <LoadingSpinner size="sm" />
            </div>
            <div v-else-if="activeSubscriptions.length" class="subscription-list">
              <article v-for="sub in activeSubscriptions" :key="sub.id" class="subscription-current">
                <div class="subscription-current__top">
                  <div>
                    <h3>{{ sub.group?.name || `分组 #${sub.group_id}` }}</h3>
                    <p>{{ formatExpires(sub.expires_at) }}</p>
                  </div>
                  <span v-if="sub.quota_bonus_multiplier && sub.quota_bonus_multiplier > 1">
                    {{ formatMultiplier(sub.quota_bonus_multiplier) }}
                  </span>
                </div>
                <div class="subscription-current__meter">
                  <div :style="{ width: `${subscriptionDailyPercent(sub)}%` }"></div>
                </div>
                <dl class="subscription-current__stats">
                  <div>
                    <dt>每日额度</dt>
                    <dd>${{ formatQuota(subscriptionDailyLimit(sub)) }}</dd>
                  </div>
                  <div>
                    <dt>今日已用</dt>
                    <dd>${{ formatQuota(sub.daily_usage_usd) }}</dd>
                  </div>
                  <div>
                    <dt>次日结转</dt>
                    <dd>${{ formatQuota(sub.daily_rollover_usd || 0) }}</dd>
                  </div>
                  <div>
                    <dt>今日剩余</dt>
                    <dd>${{ formatQuota(subscriptionDailyRemaining(sub)) }}</dd>
                  </div>
                </dl>
              </article>
            </div>
            <div v-else class="subscription-shop__empty subscription-shop__empty--compact">
              还没有可用订阅。购买卡密后在上方兑换即可开通。
            </div>
          </section>

          <section class="subscription-panel">
            <h2>拼卡奖励</h2>
            <div class="subscription-bonus">
              <div v-for="bonus in groupBuyBonuses" :key="bonus.target">
                <strong>{{ bonus.target }} 人</strong>
                <span>{{ bonus.label }}</span>
              </div>
            </div>
            <p class="subscription-panel__note">
              先兑换对应套餐，才能创建或加入同套餐拼卡房间。成团后订阅额度自动升级。
            </p>
          </section>
        </aside>
      </section>

      <section class="subscription-shop__hall">
        <div class="subscription-shop__section-head">
          <div>
            <h2>拼卡大厅</h2>
            <p>拼卡只提升订阅额度，不改变链动小铺售价。</p>
          </div>
          <div class="subscription-shop__room-actions">
            <select v-model.number="selectedCreatePlanID">
              <option :value="0">选择已激活套餐</option>
              <option v-for="plan in activePlans" :key="plan.key" :value="plan.id">{{ plan.name }}</option>
            </select>
            <select v-model.number="selectedTargetCount">
              <option :value="3">3 人</option>
              <option :value="5">5 人</option>
              <option :value="10">10 人</option>
            </select>
            <button type="button" :disabled="!canCreateRoom || groupBuySubmitting" @click="createRoom">
              发起拼卡
            </button>
          </div>
        </div>

        <div v-if="groupBuyLoading" class="subscription-shop__loading">
          <LoadingSpinner size="md" />
        </div>
        <div v-else-if="allRooms.length" class="subscription-shop__rooms">
          <article v-for="room in allRooms" :key="room.id" class="subscription-room">
            <div class="subscription-room__top">
              <div>
                <h3>{{ room.plan_name }}</h3>
                <p>{{ room.group_name }}</p>
              </div>
              <span :class="roomStatusClass(room.status)">
                {{ roomStatusLabel(room.status) }}
              </span>
            </div>
            <div class="subscription-room__progress">
              <div>
                <span>{{ room.member_count }} / {{ room.target_count }} 人</span>
                <span>{{ formatMultiplier(room.bonus_multiplier) }}</span>
              </div>
              <i><b :style="{ width: `${roomProgress(room)}%` }"></b></i>
            </div>
            <dl class="subscription-room__quota">
              <div>
                <dt>日额度</dt>
                <dd>${{ formatQuota(room.base_daily_limit_usd) }} -> ${{ formatQuota(room.upgraded_daily_limit_usd) }}</dd>
              </div>
              <div>
                <dt>周额度</dt>
                <dd>${{ formatQuota(room.base_weekly_limit_usd) }} -> ${{ formatQuota(room.upgraded_weekly_limit_usd) }}</dd>
              </div>
            </dl>
            <div class="subscription-room__foot">
              <span>{{ formatRoomTime(room) }}</span>
              <button type="button" :disabled="!room.can_join || groupBuySubmitting" @click="joinRoom(room)">
                {{ room.joined ? '已加入' : room.status === 'completed' ? '已成团' : '加入' }}
              </button>
            </div>
          </article>
        </div>
        <div v-else class="subscription-shop__empty">
          现在还没有拼卡房间。兑换订阅后可以先发起一个。
        </div>
      </section>
    </main>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import Icon from '@/components/icons/Icon.vue'
import { paymentAPI } from '@/api/payment'
import { redeemAPI } from '@/api/redeem'
import { groupBuyAPI, type GroupBuyRoom, type GroupBuyTargetCount } from '@/api/groupBuy'
import { useAppStore } from '@/stores/app'
import { useSubscriptionStore } from '@/stores/subscriptions'
import { LDXP_FALLBACK_SUBSCRIPTION_PLANS, LDXP_SHOP_URL, SUBSCRIPTION_DISPLAY_RATE_MULTIPLIER } from '@/constants/onboarding'
import { extractApiErrorMessage } from '@/utils/apiError'
import type { SubscriptionPlan } from '@/types/payment'
import type { UserSubscription } from '@/types'

type PlanTab = 'all' | 'day' | 'week'

interface ShopPlan {
  key: string
  id: number
  group_id: number
  group_name?: string
  rate_multiplier: number
  name: string
  price: number
  validity_days: number
  daily_limit_usd: number
  weekly_limit_usd: number
  sort_order: number
  is_fallback: boolean
}

const fallbackPlans: ShopPlan[] = LDXP_FALLBACK_SUBSCRIPTION_PLANS.map((plan, index) => ({
  key: `fallback-${index + 1}`,
  id: 0,
  group_id: 0,
  group_name: '订阅分组',
  rate_multiplier: SUBSCRIPTION_DISPLAY_RATE_MULTIPLIER,
  name: plan.name,
  price: plan.price,
  validity_days: plan.validity_days,
  daily_limit_usd: plan.daily_limit_usd,
  weekly_limit_usd: plan.weekly_limit_usd,
  sort_order: plan.sort_order,
  is_fallback: true
}))

const appStore = useAppStore()
const subscriptionStore = useSubscriptionStore()

const plans = ref<ShopPlan[]>([])
const rooms = ref<GroupBuyRoom[]>([])
const myRooms = ref<GroupBuyRoom[]>([])
const plansLoading = ref(false)
const groupBuyLoading = ref(false)
const groupBuySubmitting = ref(false)
const redeeming = ref(false)
const redeemCode = ref('')
const redeemError = ref('')
const redeemResult = ref<{ message?: string; type?: string } | null>(null)
const planTab = ref<PlanTab>('all')
const selectedCreatePlanID = ref(0)
const selectedTargetCount = ref<GroupBuyTargetCount>(3)

const planTabs = [
  { value: 'all' as const, label: '全部' },
  { value: 'day' as const, label: '日卡' },
  { value: 'week' as const, label: '周卡' }
]

const groupBuyBonuses = [
  { target: 3, label: '额度 +10%' },
  { target: 5, label: '额度 +17%' },
  { target: 10, label: '额度 +34%' }
]

const activeSubscriptions = computed(() => subscriptionStore.activeSubscriptions)

const activeGroupIDs = computed(() => {
  return new Set(activeSubscriptions.value.map((sub) => sub.group_id))
})

const activePlans = computed(() => {
  return plans.value.filter((plan) => plan.id > 0 && activeGroupIDs.value.has(plan.group_id))
})

const visiblePlans = computed(() => {
  if (planTab.value === 'day') return plans.value.filter((plan) => plan.validity_days <= 1)
  if (planTab.value === 'week') return plans.value.filter((plan) => plan.validity_days > 1)
  return plans.value
})

const allRooms = computed(() => {
  const seen = new Set<number>()
  return [...myRooms.value, ...rooms.value].filter((room) => {
    if (seen.has(room.id)) return false
    seen.add(room.id)
    return true
  })
})

const canCreateRoom = computed(() => selectedCreatePlanID.value > 0 && activePlans.value.some((plan) => plan.id === selectedCreatePlanID.value))

watch(activePlans, (items) => {
  if (!items.some((plan) => plan.id === selectedCreatePlanID.value)) {
    selectedCreatePlanID.value = items[0]?.id || 0
  }
})

async function loadPlans() {
  plansLoading.value = true
  try {
    let loaded: ShopPlan[] = []
    try {
      const checkout = await paymentAPI.getCheckoutInfo()
      loaded = normalizePlans(checkout.data?.plans || [])
    } catch {
      const response = await paymentAPI.getPlans()
      loaded = normalizePlans(response.data || [])
    }
    plans.value = loaded.length ? loaded : fallbackPlans
  } catch (error) {
    console.warn('Failed to load shop plans, using fallback plans:', error)
    plans.value = fallbackPlans
  } finally {
    plansLoading.value = false
  }
}

async function loadGroupBuys() {
  groupBuyLoading.value = true
  try {
    const hall = await groupBuyAPI.listGroupBuys()
    rooms.value = hall.rooms || []
    myRooms.value = hall.my_rooms || []
  } catch (error) {
    rooms.value = []
    myRooms.value = []
    console.warn('Failed to load group-buy hall:', error)
  } finally {
    groupBuyLoading.value = false
  }
}

async function refreshAll() {
  await Promise.all([
    subscriptionStore.fetchActiveSubscriptions(true).catch(() => []),
    loadGroupBuys()
  ])
}

async function handleRedeem() {
  const code = redeemCode.value.trim()
  if (!code) return

  redeeming.value = true
  redeemError.value = ''
  redeemResult.value = null
  try {
    const result = await redeemAPI.redeem(code)
    redeemResult.value = result
    redeemCode.value = ''
    appStore.showSuccess('兑换成功，订阅状态已刷新')
    await refreshAll()
  } catch (error) {
    redeemError.value = extractApiErrorMessage(error, '兑换失败，请检查卡密')
  } finally {
    redeeming.value = false
  }
}

async function createRoom() {
  if (!canCreateRoom.value) return
  groupBuySubmitting.value = true
  try {
    await groupBuyAPI.createGroupBuyRoom({
      plan_id: selectedCreatePlanID.value,
      target_count: selectedTargetCount.value
    })
    appStore.showSuccess('拼卡房间已创建')
    await loadGroupBuys()
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, '创建拼卡房间失败'))
  } finally {
    groupBuySubmitting.value = false
  }
}

async function joinRoom(room: GroupBuyRoom) {
  if (!room.can_join) return
  groupBuySubmitting.value = true
  try {
    await groupBuyAPI.joinGroupBuyRoom(room.id)
    appStore.showSuccess('已加入拼卡房间')
    await refreshAll()
  } catch (error) {
    appStore.showError(extractApiErrorMessage(error, '加入拼卡失败'))
  } finally {
    groupBuySubmitting.value = false
  }
}

function planPurchaseUrl(_plan: ShopPlan): string {
  return LDXP_SHOP_URL
}

function trackOpenShop() {
  appStore.showInfo('已打开链动小铺，付款后请回本站兑换卡密。')
}

function trackPlanPurchase(plan: ShopPlan) {
  appStore.showInfo(`已打开链动小铺，请选择：${plan.name}`)
}

function normalizePlans(rawPlans: SubscriptionPlan[]): ShopPlan[] {
  return rawPlans
    .filter((plan) => plan.for_sale !== false)
    .map((plan, index) => {
      const inferredDaily = inferDailyQuota(plan)
      const daily = toPositiveNumber(plan.daily_limit_usd) || inferredDaily
      const weekly = toPositiveNumber(plan.weekly_limit_usd) || daily * Number(plan.validity_days || 1)
      return {
        key: `plan-${plan.id}`,
        id: Number(plan.id),
        group_id: Number(plan.group_id),
        group_name: plan.group_name,
        rate_multiplier: SUBSCRIPTION_DISPLAY_RATE_MULTIPLIER,
        name: plan.name,
        price: toFiniteNumber(plan.price),
        validity_days: Number(plan.validity_days || 1),
        daily_limit_usd: daily,
        weekly_limit_usd: weekly,
        sort_order: Number(plan.sort_order ?? index),
        is_fallback: false
      }
    })
    .sort((a, b) => a.sort_order - b.sort_order || a.price - b.price)
}

function displayGroupName(plan: ShopPlan): string {
  if (!plan.group_name) return '订阅分组'
  return plan.group_name.replace(/1(?:\.0)?倍/g, `${SUBSCRIPTION_DISPLAY_RATE_MULTIPLIER}倍`)
}

function inferDailyQuota(plan: SubscriptionPlan): number {
  const matched = `${plan.name} ${plan.description || ''}`.match(/(\d+(?:\.\d+)?)\s*刀/)
  if (matched) return toFiniteNumber(matched[1])
  return 0
}

function toPositiveNumber(value: unknown): number {
  const n = toFiniteNumber(value)
  return n > 0 ? n : 0
}

function toFiniteNumber(value: unknown): number {
  const n = Number(value ?? 0)
  return Number.isFinite(n) ? n : 0
}

function formatPrice(value: number): string {
  const n = toFiniteNumber(value)
  return n % 1 === 0 ? n.toFixed(0) : n.toFixed(2).replace(/0$/, '')
}

function formatQuota(value: number | null | undefined): string {
  const n = toFiniteNumber(value)
  if (n <= 0) return '不限'
  return n % 1 === 0 ? n.toFixed(0) : n.toFixed(2)
}

function formatValidity(days: number): string {
  return days > 1 ? `${days} 天` : '1 天'
}

function formatRate(value: number | null | undefined): string {
  const n = toFiniteNumber(value)
  if (n <= 0) return `${SUBSCRIPTION_DISPLAY_RATE_MULTIPLIER} 倍`
  return `${n.toFixed(n % 1 === 0 ? 0 : 2)} 倍`
}

function formatMultiplier(value: number): string {
  const n = toFiniteNumber(value)
  if (n <= 1) return '基础额度'
  return `x${n.toFixed(2)}`
}

function formatExpires(value: string | null): string {
  if (!value) return '长期有效'
  const expires = new Date(value)
  if (Number.isNaN(expires.getTime())) return value
  const days = Math.max(0, Math.ceil((expires.getTime() - Date.now()) / 86_400_000))
  return `${days} 天后到期`
}

function formatRoomTime(room: GroupBuyRoom): string {
  if (room.status === 'completed') return '已成团'
  const expires = new Date(room.expires_at)
  if (Number.isNaN(expires.getTime())) return ''
  const minutes = Math.max(0, Math.ceil((expires.getTime() - Date.now()) / 60_000))
  if (minutes >= 60) return `剩余 ${Math.ceil(minutes / 60)} 小时`
  return `剩余 ${minutes} 分钟`
}

function subscriptionDailyLimit(sub: UserSubscription): number | null {
  return sub.daily_limit_override_usd ?? sub.group?.daily_limit_usd ?? null
}

function subscriptionDailyRemaining(sub: UserSubscription): number | null {
  if (sub.daily_remaining_usd !== undefined && sub.daily_remaining_usd !== null) {
    return sub.daily_remaining_usd
  }
  const available = sub.daily_available_usd ?? ((subscriptionDailyLimit(sub) || 0) + (sub.daily_rollover_usd || 0))
  if (!available) return null
  return Math.max(0, available - (sub.daily_usage_usd || 0))
}

function subscriptionDailyPercent(sub: UserSubscription): number {
  const available = sub.daily_available_usd ?? ((subscriptionDailyLimit(sub) || 0) + (sub.daily_rollover_usd || 0))
  if (!available) return 0
  return Math.min(100, Math.round(((sub.daily_usage_usd || 0) / available) * 100))
}

function roomProgress(room: GroupBuyRoom): number {
  if (room.target_count <= 0) return 0
  return Math.min(100, Math.round((room.member_count / room.target_count) * 100))
}

function roomStatusLabel(status: string): string {
  if (status === 'completed') return '已成团'
  if (status === 'expired') return '已过期'
  if (status === 'canceled') return '已取消'
  return '进行中'
}

function roomStatusClass(status: string): string {
  const base = 'subscription-room__status'
  if (status === 'completed') return `${base} subscription-room__status--done`
  if (status === 'expired' || status === 'canceled') return `${base} subscription-room__status--muted`
  return base
}

onMounted(async () => {
  await Promise.all([
    loadPlans(),
    subscriptionStore.fetchActiveSubscriptions().catch(() => []),
    loadGroupBuys()
  ])
})
</script>

<style scoped>
.subscription-shop {
  display: grid;
  gap: 1rem;
}

.subscription-shop__hero,
.subscription-shop__plans,
.subscription-shop__side > section,
.subscription-shop__hall,
.subscription-plan,
.subscription-room {
  border: 1px solid rgba(43, 132, 83, 0.38);
  border-radius: 8px;
  background: rgba(5, 29, 24, 0.78);
  box-shadow: 0 18px 54px rgba(0, 0, 0, 0.22);
}

.subscription-shop__hero {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(19rem, 0.38fr);
  gap: 1rem;
  padding: clamp(1rem, 2vw, 1.55rem);
}

.subscription-shop__intro {
  display: flex;
  min-height: 14rem;
  flex-direction: column;
  justify-content: center;
}

.subscription-shop__eyebrow {
  margin: 0 0 0.65rem;
  color: #22c55e;
  font-size: 0.72rem;
  font-weight: 850;
  letter-spacing: 0.12em;
}

.subscription-shop h1,
.subscription-shop h2,
.subscription-shop h3,
.subscription-shop p,
.subscription-shop dl {
  margin: 0;
}

.subscription-shop h1 {
  color: #edf7ee;
  font-size: clamp(2.1rem, 4vw, 3.8rem);
  font-weight: 900;
  letter-spacing: 0;
  line-height: 1;
}

.subscription-shop__hero-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 0.7rem;
  margin-top: 1.2rem;
}

.subscription-shop__redeem,
.subscription-panel {
  display: grid;
  gap: 0.85rem;
  border-radius: 8px;
  background: rgba(3, 22, 18, 0.72);
  padding: 1rem;
}

.subscription-shop__redeem h2,
.subscription-panel h2 {
  color: #edf7ee;
  font-size: 1rem;
  font-weight: 850;
}

.subscription-shop__redeem p,
.subscription-panel p {
  margin-top: 0.25rem;
  color: rgba(220, 239, 224, 0.62);
  font-size: 0.86rem;
  line-height: 1.55;
}

.subscription-shop__field {
  display: grid;
  gap: 0.38rem;
  color: rgba(220, 239, 224, 0.76);
  font-size: 0.8rem;
  font-weight: 750;
}

.subscription-shop__field input,
.subscription-shop__room-actions select {
  width: 100%;
  min-height: 2.5rem;
  border: 1px solid rgba(43, 132, 83, 0.38);
  border-radius: 8px;
  background: rgba(3, 22, 18, 0.64);
  padding: 0 0.75rem;
  color: #edf7ee;
  font-size: 0.9rem;
  outline: none;
}

.subscription-shop__field input:focus,
.subscription-shop__room-actions select:focus {
  border-color: #22c55e;
  box-shadow: 0 0 0 3px rgba(34, 197, 94, 0.14);
}

.subscription-shop__primary,
.subscription-shop__secondary,
.subscription-shop__room-actions button,
.subscription-room__foot button,
.subscription-panel__head button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.4rem;
  border-radius: 8px;
  font-weight: 800;
  text-decoration: none;
  transition: background-color 0.16s ease, border-color 0.16s ease, color 0.16s ease;
}

.subscription-shop__primary,
.subscription-shop__room-actions button {
  min-height: 2.5rem;
  border: 1px solid rgba(34, 197, 94, 0.32);
  background: #16a34a;
  padding: 0 1rem;
  color: #ffffff;
}

.subscription-shop__primary:hover,
.subscription-shop__room-actions button:hover {
  background: #0f8a3f;
}

.subscription-shop__secondary {
  min-height: 2.45rem;
  border: 1px solid rgba(43, 132, 83, 0.38);
  background: rgba(3, 22, 18, 0.64);
  padding: 0 1rem;
  color: #edf7ee;
  font-size: 0.88rem;
}

.subscription-shop__secondary:hover {
  border-color: rgba(34, 197, 94, 0.52);
  background: rgba(34, 197, 94, 0.12);
}

.subscription-shop__secondary--full {
  width: 100%;
  margin-top: auto;
}

.subscription-shop button:disabled,
.subscription-shop__primary:disabled,
.subscription-room__foot button:disabled {
  cursor: not-allowed;
  opacity: 0.5;
}

.subscription-shop__message {
  border-radius: 8px;
  padding: 0.65rem 0.75rem;
  font-size: 0.84rem;
  font-weight: 700;
}

.subscription-shop__message--error {
  background: #fff1f0;
  color: #b42318;
}

.subscription-shop__message--success {
  background: rgba(34, 197, 94, 0.14);
  color: #7ee39b;
}

.subscription-shop__layout {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(18rem, 22rem);
  gap: 1rem;
  align-items: start;
}

.subscription-shop__plans,
.subscription-shop__hall {
  padding: 1rem;
}

.subscription-shop__section-head,
.subscription-panel__head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 1rem;
}

.subscription-shop__section-head h2 {
  color: #edf7ee;
  font-size: 1.18rem;
  font-weight: 900;
}

.subscription-shop__section-head p {
  margin-top: 0.28rem;
  color: rgba(220, 239, 224, 0.62);
  font-size: 0.88rem;
}

.subscription-shop__tabs {
  display: inline-flex;
  gap: 0.18rem;
  border: 1px solid rgba(43, 132, 83, 0.38);
  border-radius: 8px;
  background: rgba(3, 22, 18, 0.72);
  padding: 0.18rem;
}

.subscription-shop__tabs button {
  min-height: 2rem;
  border: 0;
  border-radius: 6px;
  background: transparent;
  padding: 0 0.82rem;
  color: rgba(220, 239, 224, 0.62);
  font-size: 0.84rem;
  font-weight: 800;
}

.subscription-shop__tabs button.active {
  background: rgba(34, 197, 94, 0.14);
  color: #edf7ee;
  box-shadow: none;
}

.subscription-shop__loading,
.subscription-shop__side-loading {
  display: flex;
  min-height: 10rem;
  align-items: center;
  justify-content: center;
}

.subscription-shop__plan-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 0.75rem;
  margin-top: 1rem;
}

.subscription-plan {
  display: flex;
  min-width: 0;
  flex-direction: column;
  gap: 0.9rem;
  padding: 1rem;
}

.subscription-plan__top,
.subscription-current__top,
.subscription-room__top,
.subscription-room__foot {
  display: flex;
  justify-content: space-between;
  gap: 0.75rem;
}

.subscription-plan__top h3,
.subscription-current__top h3,
.subscription-room__top h3 {
  color: #edf7ee;
  font-size: 0.98rem;
  font-weight: 900;
}

.subscription-plan__top p,
.subscription-current__top p,
.subscription-room__top p,
.subscription-room__foot span {
  margin-top: 0.22rem;
  color: rgba(220, 239, 224, 0.62);
  font-size: 0.78rem;
}

.subscription-plan__top > span,
.subscription-current__top > span {
  align-self: flex-start;
  border-radius: 6px;
  background: rgba(34, 197, 94, 0.14);
  padding: 0.28rem 0.5rem;
  color: #7ee39b;
  font-size: 0.76rem;
  font-weight: 850;
  white-space: nowrap;
}

.subscription-plan__price {
  display: flex;
  align-items: end;
  gap: 0.5rem;
}

.subscription-plan__price strong {
  color: #edf7ee;
  font-size: 2rem;
  line-height: 1;
}

.subscription-plan__price span {
  color: rgba(220, 239, 224, 0.62);
  font-size: 0.82rem;
}

.subscription-plan__metrics,
.subscription-current__stats,
.subscription-room__quota {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0.5rem;
}

.subscription-plan__metrics div,
.subscription-current__stats div,
.subscription-room__quota div {
  min-width: 0;
  border: 1px solid rgba(43, 132, 83, 0.32);
  border-radius: 8px;
  background: rgba(3, 22, 18, 0.64);
  padding: 0.65rem;
}

.subscription-plan__metrics dt,
.subscription-current__stats dt,
.subscription-room__quota dt {
  color: rgba(220, 239, 224, 0.62);
  font-size: 0.72rem;
  font-weight: 750;
}

.subscription-plan__metrics dd,
.subscription-current__stats dd,
.subscription-room__quota dd {
  overflow: hidden;
  margin: 0.25rem 0 0;
  color: #edf7ee;
  font-size: 0.86rem;
  font-weight: 850;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.subscription-shop__side {
  display: grid;
  gap: 1rem;
}

.subscription-panel__head button {
  width: 2rem;
  height: 2rem;
  border: 1px solid rgba(43, 132, 83, 0.38);
  background: rgba(3, 22, 18, 0.64);
  color: #22c55e;
}

.subscription-list {
  display: grid;
  gap: 0.75rem;
}

.subscription-current {
  border: 1px solid rgba(43, 132, 83, 0.32);
  border-radius: 8px;
  background: rgba(3, 22, 18, 0.64);
  padding: 0.85rem;
}

.subscription-current__meter,
.subscription-room__progress i {
  display: block;
  overflow: hidden;
  height: 0.44rem;
  border-radius: 999px;
  background: rgba(43, 132, 83, 0.22);
}

.subscription-current__meter {
  margin: 0.75rem 0;
}

.subscription-current__meter div,
.subscription-room__progress b {
  display: block;
  height: 100%;
  border-radius: inherit;
  background: #22c55e;
}

.subscription-bonus {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 0.45rem;
}

.subscription-bonus div {
  border: 1px solid rgba(43, 132, 83, 0.32);
  border-radius: 8px;
  background: rgba(3, 22, 18, 0.64);
  padding: 0.65rem 0.4rem;
  text-align: center;
}

.subscription-bonus strong,
.subscription-bonus span {
  display: block;
}

.subscription-bonus strong {
  color: #edf7ee;
  font-size: 0.86rem;
}

.subscription-bonus span {
  margin-top: 0.16rem;
  color: #22c55e;
  font-size: 0.74rem;
  font-weight: 800;
}

.subscription-panel__note {
  margin-top: 0;
}

.subscription-shop__room-actions {
  display: grid;
  grid-template-columns: minmax(11rem, 1fr) 6rem auto;
  gap: 0.5rem;
  align-items: center;
}

.subscription-shop__room-actions button {
  white-space: nowrap;
}

.subscription-shop__rooms {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 0.75rem;
  margin-top: 1rem;
}

.subscription-room {
  display: grid;
  gap: 0.85rem;
  padding: 1rem;
}

.subscription-room__status {
  align-self: flex-start;
  border-radius: 6px;
  background: rgba(34, 197, 94, 0.14);
  padding: 0.28rem 0.5rem;
  color: #7ee39b;
  font-size: 0.74rem;
  font-weight: 850;
  white-space: nowrap;
}

.subscription-room__status--done {
  background: rgba(34, 197, 94, 0.14);
  color: #7ee39b;
}

.subscription-room__status--muted {
  background: rgba(34, 197, 94, 0.08);
  color: rgba(220, 239, 224, 0.62);
}

.subscription-room__progress div {
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.45rem;
  color: rgba(220, 239, 224, 0.62);
  font-size: 0.78rem;
  font-weight: 750;
}

.subscription-room__foot {
  align-items: center;
}

.subscription-room__foot button {
  min-height: 2rem;
  border: 1px solid rgba(43, 132, 83, 0.38);
  background: rgba(3, 22, 18, 0.64);
  padding: 0 0.8rem;
  color: #edf7ee;
  font-size: 0.82rem;
}

.subscription-shop__empty {
  margin-top: 1rem;
  border: 1px dashed rgba(43, 132, 83, 0.38);
  border-radius: 8px;
  background: rgba(3, 22, 18, 0.64);
  padding: 1.5rem;
  color: rgba(220, 239, 224, 0.62);
  font-size: 0.9rem;
  text-align: center;
}

.subscription-shop__empty--compact {
  margin-top: 0;
  padding: 1rem;
  text-align: left;
}

@media (max-width: 1100px) {
  .subscription-shop__layout,
  .subscription-shop__hero {
    grid-template-columns: 1fr;
  }

  .subscription-shop__plan-grid,
  .subscription-shop__rooms {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 720px) {
  .subscription-shop__section-head,
  .subscription-panel__head {
    display: grid;
  }

  .subscription-shop__tabs {
    width: 100%;
  }

  .subscription-shop__tabs button {
    flex: 1;
  }

  .subscription-shop__hero-actions,
  .subscription-shop__primary,
  .subscription-shop__secondary {
    width: 100%;
  }

  .subscription-shop__plan-grid,
  .subscription-shop__rooms,
  .subscription-plan__metrics,
  .subscription-current__stats,
  .subscription-room__quota {
    grid-template-columns: 1fr;
  }

  .subscription-shop__room-actions {
    grid-template-columns: 1fr;
    width: 100%;
  }
}
</style>
