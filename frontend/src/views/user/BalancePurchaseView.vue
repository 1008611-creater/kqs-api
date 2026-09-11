<template>
  <AppLayout>
    <main class="purchase-market">
      <section class="purchase-market__hero">
        <div class="purchase-market__intro">
          <p class="purchase-market__eyebrow">余额购买</p>
          <h1>购买余额</h1>
          <div class="purchase-market__actions">
            <a class="purchase-market__primary" :href="LDXP_SHOP_URL" target="_blank" rel="noopener noreferrer" @click="trackOpenShop">
              <Icon name="externalLink" size="sm" />
              打开链动小铺
            </a>
            <RouterLink class="purchase-market__secondary" to="/redeem">
              <Icon name="gift" size="sm" />
              已有卡密，去兑换
            </RouterLink>
            <RouterLink class="purchase-market__secondary" to="/shop">
              查看订阅套餐
            </RouterLink>
          </div>
        </div>

        <div class="purchase-market__redeem-card">
          <h2>购买后怎么到账</h2>
          <ol>
            <li><span>01</span><p>在链动小铺选择余额商品并付款。</p></li>
            <li><span>02</span><p>复制订单详情里的卡密。</p></li>
            <li><span>03</span><p>回到本站兑换，余额到账后即可使用。</p></li>
          </ol>
        </div>
      </section>

      <section class="purchase-market__panel">
        <div class="purchase-market__section-head">
          <div>
            <h2>余额商品</h2>
            <p>适合按量使用。买多少充多少，不绑定订阅周期。</p>
          </div>
        </div>

        <div class="purchase-market__grid">
          <article v-for="item in balanceProducts" :key="item.id" class="purchase-product" :class="{ 'purchase-product--recommended': item.recommended }">
            <div class="purchase-product__top">
              <div>
                <h3>{{ item.quotaLabel }}</h3>
                <p>{{ item.note || '一次性余额卡密' }}</p>
              </div>
              <span v-if="item.recommended">推荐</span>
            </div>

            <div class="purchase-product__price">
              <strong>¥{{ item.priceCny }}</strong>
              <span>链动小铺卡密</span>
            </div>

            <dl class="purchase-product__metrics">
              <div>
                <dt>到账额度</dt>
                <dd>${{ formatQuota(item.quotaUsd) }}</dd>
              </div>
              <div>
                <dt>商品类型</dt>
                <dd>余额卡</dd>
              </div>
            </dl>

            <a class="purchase-market__secondary purchase-market__secondary--full" :href="item.url" target="_blank" rel="noopener noreferrer" @click="trackProduct(item)">
              去购买
              <Icon name="arrowRight" size="xs" />
            </a>
          </article>
        </div>
      </section>

      <section class="purchase-market__notice">
        <div>
          <h2>余额和订阅的区别</h2>
          <p>余额是一次性到账，适合灵活使用；订阅是日卡/周卡额度，适合固定周期使用。两类商品分开购买、分开兑换。</p>
        </div>
        <RouterLink to="/shop">
          去订阅页
          <Icon name="arrowRight" size="xs" />
        </RouterLink>
      </section>
    </main>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import AppLayout from '@/components/layout/AppLayout.vue'
import Icon from '@/components/icons/Icon.vue'
import { LDXP_BALANCE_PACKAGES, LDXP_SHOP_URL, type LdxpPackage } from '@/constants/onboarding'
import { useAppStore } from '@/stores'

const appStore = useAppStore()
const balanceProducts = computed(() => (
  [...LDXP_BALANCE_PACKAGES].sort((left, right) => {
    const leftQuota = left.quotaUsd ?? Number.POSITIVE_INFINITY
    const rightQuota = right.quotaUsd ?? Number.POSITIVE_INFINITY
    if (leftQuota !== rightQuota) return leftQuota - rightQuota
    return Number.parseFloat(left.priceCny) - Number.parseFloat(right.priceCny)
  })
))

function formatQuota(value: number | undefined): string {
  if (!value || value <= 0) return '按卡密到账'
  return value % 1 === 0 ? value.toFixed(0) : value.toFixed(1)
}

function trackOpenShop() {
  appStore.showInfo('已打开链动小铺，付款后请回本站兑换卡密。')
}

function trackProduct(item: LdxpPackage) {
  appStore.showInfo(`已打开链动小铺，请选择：${item.quotaLabel}`)
}
</script>

<style scoped>
.purchase-market {
  display: grid;
  gap: 1rem;
}

.purchase-market__hero,
.purchase-market__panel,
.purchase-market__notice,
.purchase-product {
  border: 1px solid rgba(43, 132, 83, 0.38);
  border-radius: 8px;
  background: rgba(5, 29, 24, 0.78);
  box-shadow: 0 18px 54px rgba(0, 0, 0, 0.22);
}

.purchase-market__hero {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(18rem, 0.38fr);
  gap: 1rem;
  padding: clamp(1rem, 2vw, 1.55rem);
}

.purchase-market__intro {
  display: flex;
  min-height: 14rem;
  flex-direction: column;
  justify-content: center;
}

.purchase-market__eyebrow {
  margin: 0 0 0.65rem;
  color: #22c55e;
  font-size: 0.72rem;
  font-weight: 850;
  letter-spacing: 0.12em;
}

.purchase-market h1,
.purchase-market h2,
.purchase-market h3,
.purchase-market p,
.purchase-market dl {
  margin: 0;
}

.purchase-market h1 {
  color: #edf7ee;
  font-size: clamp(2.1rem, 4vw, 3.8rem);
  font-weight: 900;
  letter-spacing: 0;
  line-height: 1;
}

.purchase-market__actions {
  display: flex;
  flex-wrap: wrap;
  gap: 0.7rem;
  margin-top: 1.2rem;
}

.purchase-market__primary,
.purchase-market__secondary,
.purchase-market__notice a {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.4rem;
  min-height: 2.45rem;
  border-radius: 8px;
  padding: 0 1rem;
  font-size: 0.88rem;
  font-weight: 800;
  text-decoration: none;
  transition: background-color 0.16s ease, border-color 0.16s ease, color 0.16s ease;
}

.purchase-market__primary {
  border: 1px solid rgba(34, 197, 94, 0.32);
  background: #16a34a;
  color: #ffffff;
}

.purchase-market__primary:hover {
  background: #0f8a3f;
}

.purchase-market__secondary,
.purchase-market__notice a {
  border: 1px solid rgba(43, 132, 83, 0.38);
  background: rgba(3, 22, 18, 0.64);
  color: #edf7ee;
}

.purchase-market__secondary:hover,
.purchase-market__notice a:hover {
  border-color: rgba(34, 197, 94, 0.52);
  background: rgba(34, 197, 94, 0.12);
}

.purchase-market__secondary--full {
  width: 100%;
  margin-top: auto;
}

.purchase-market__redeem-card {
  display: grid;
  gap: 0.85rem;
  align-self: stretch;
  border-radius: 8px;
  background: rgba(3, 22, 18, 0.72);
  padding: 1rem;
}

.purchase-market__redeem-card h2,
.purchase-market__section-head h2,
.purchase-market__notice h2 {
  color: #edf7ee;
  font-size: 1.05rem;
  font-weight: 900;
}

.purchase-market__redeem-card ol {
  display: grid;
  gap: 0.65rem;
  margin: 0;
  padding: 0;
  list-style: none;
}

.purchase-market__redeem-card li {
  display: grid;
  grid-template-columns: 2rem minmax(0, 1fr);
  gap: 0.65rem;
  border: 1px solid rgba(43, 132, 83, 0.32);
  border-radius: 8px;
  background: rgba(5, 29, 24, 0.72);
  padding: 0.7rem;
}

.purchase-market__redeem-card li span {
  color: #22c55e;
  font-size: 0.78rem;
  font-weight: 900;
}

.purchase-market__redeem-card li p {
  color: rgba(220, 239, 224, 0.76);
  font-size: 0.86rem;
  line-height: 1.55;
}

.purchase-market__panel {
  padding: 1rem;
}

.purchase-market__section-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 1rem;
}

.purchase-market__section-head p,
.purchase-market__notice p {
  margin-top: 0.28rem;
  color: rgba(220, 239, 224, 0.62);
  font-size: 0.88rem;
  line-height: 1.55;
}

.purchase-market__grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 0.75rem;
  margin-top: 1rem;
}

.purchase-product {
  display: flex;
  min-width: 0;
  flex-direction: column;
  gap: 0.9rem;
  padding: 1rem;
}

.purchase-product--recommended {
  border-color: rgba(34, 197, 94, 0.52);
  box-shadow: 0 18px 54px rgba(0, 0, 0, 0.24);
}

.purchase-product__top {
  display: flex;
  justify-content: space-between;
  gap: 0.75rem;
}

.purchase-product__top h3 {
  color: #edf7ee;
  font-size: 0.98rem;
  font-weight: 900;
}

.purchase-product__top p {
  margin-top: 0.22rem;
  color: rgba(220, 239, 224, 0.62);
  font-size: 0.78rem;
}

.purchase-product__top > span {
  align-self: flex-start;
  border-radius: 6px;
  background: rgba(34, 197, 94, 0.14);
  padding: 0.28rem 0.5rem;
  color: #7ee39b;
  font-size: 0.76rem;
  font-weight: 850;
  white-space: nowrap;
}

.purchase-product__price {
  display: flex;
  align-items: end;
  gap: 0.5rem;
}

.purchase-product__price strong {
  color: #edf7ee;
  font-size: 2rem;
  line-height: 1;
}

.purchase-product__price span {
  color: rgba(220, 239, 224, 0.62);
  font-size: 0.82rem;
}

.purchase-product__metrics {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 0.5rem;
}

.purchase-product__metrics div {
  min-width: 0;
  border: 1px solid rgba(43, 132, 83, 0.32);
  border-radius: 8px;
  background: rgba(3, 22, 18, 0.64);
  padding: 0.65rem;
}

.purchase-product__metrics dt {
  color: rgba(220, 239, 224, 0.62);
  font-size: 0.72rem;
  font-weight: 750;
}

.purchase-product__metrics dd {
  overflow: hidden;
  margin: 0.25rem 0 0;
  color: #edf7ee;
  font-size: 0.86rem;
  font-weight: 850;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.purchase-market__notice {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  padding: 1rem;
}

.purchase-market__notice a {
  flex: 0 0 auto;
}

@media (max-width: 1100px) {
  .purchase-market__hero {
    grid-template-columns: 1fr;
  }

  .purchase-market__grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 720px) {
  .purchase-market__section-head,
  .purchase-market__notice {
    display: grid;
  }

  .purchase-market__grid,
  .purchase-product__metrics {
    grid-template-columns: 1fr;
  }

  .purchase-market__actions,
  .purchase-market__primary,
  .purchase-market__secondary,
  .purchase-market__notice a {
    width: 100%;
  }
}
</style>
