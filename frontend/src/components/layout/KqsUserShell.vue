<template>
  <div class="kqs-shell">
    <header class="kqs-topbar">
      <router-link to="/dashboard" class="kqs-brand" aria-label="返回仪表盘">
        <img :src="brandLogo" alt="" class="kqs-brand__logo">
        <span class="kqs-brand__name">{{ siteName }}</span>
      </router-link>

      <nav class="kqs-nav" aria-label="用户导航">
        <template v-for="item in navItems" :key="item.path">
          <a
            v-if="item.hard"
            :href="item.path"
            class="kqs-nav__link"
            :class="{ 'router-link-active': route.path === item.path }"
          >
            {{ item.label }}
          </a>
          <router-link
            v-else
            :to="item.path"
            class="kqs-nav__link"
          >
            {{ item.label }}
          </router-link>
        </template>
      </nav>

      <div class="kqs-actions">
        <router-link v-if="isAdmin" to="/admin/dashboard" class="kqs-admin-button">
          管理后台
        </router-link>
        <span v-if="user" class="kqs-balance">余额 ¥{{ balanceText }}</span>
        <button v-if="user" type="button" class="kqs-login-button" @click="handleLogout">
          退出
        </button>
        <router-link v-else to="/login" class="kqs-login-button">
          登录 / 注册
        </router-link>
      </div>
    </header>

    <main class="kqs-main">
      <slot />
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAppStore, useAuthStore } from '@/stores'
import { USER_SUBSCRIPTIONS_VISIBLE } from '@/constants/subscriptions'

const router = useRouter()
const route = useRoute()
const appStore = useAppStore()
const authStore = useAuthStore()

const user = computed(() => authStore.user)
const isAdmin = computed(() => authStore.isAdmin)
const siteName = computed(() => appStore.siteName || '矿泉水API')
const brandLogo = computed(() => appStore.siteLogo || '/kqs-water-logo.svg?v=20260613')
const balanceText = computed(() => (user.value?.balance ?? 0).toFixed(2))

const navItems = [
  { path: '/dashboard', label: '概览' },
  { path: '/keys', label: 'API Key' },
  { path: '/purchase', label: '购买余额' },
  ...(USER_SUBSCRIPTIONS_VISIBLE ? [{ path: '/shop', label: '订阅' }] : []),
  { path: '/affiliate', label: '邀请返利' },
  { path: '/usage', label: '用量' },
  { path: '/monitor', label: '监测', hard: true },
  { path: '/guide', label: '教程' },
]

async function handleLogout() {
  await authStore.logout()
  await router.push('/login')
}
</script>

<style scoped>
.kqs-shell {
  --kqs-bg: #08231d;
  --kqs-card: rgba(5, 29, 24, 0.78);
  --kqs-card-strong: rgba(4, 24, 20, 0.9);
  --kqs-card-soft: rgba(34, 197, 94, 0.14);
  --kqs-border: rgba(43, 132, 83, 0.38);
  --kqs-text: #edf7ee;
  --kqs-muted: rgba(220, 239, 224, 0.68);
  --kqs-action: #16a34a;
  --kqs-action-hover: #0f8a3f;
  min-height: 100dvh;
  position: relative;
  isolation: isolate;
  overflow-x: hidden;
  background:
    radial-gradient(ellipse at 11% 7%, rgba(34, 197, 94, 0.06), transparent 28rem),
    radial-gradient(ellipse at 84% 82%, rgba(27, 93, 70, 0.28), transparent 34rem),
    linear-gradient(115deg, rgba(255, 255, 255, 0.022), transparent 38%, rgba(255, 255, 255, 0.01) 70%, transparent),
    linear-gradient(135deg, #08231d 0%, #0b342a 48%, #082920 100%);
  color: var(--kqs-text);
}

.kqs-shell::before,
.kqs-shell::after {
  content: '';
  position: fixed;
  inset: 0;
  z-index: -1;
  pointer-events: none;
}

.kqs-shell::before {
  background:
    url('/kqs-mountain-water-texture.fast.jpg?v=20260616-fastbg') center top / cover no-repeat,
    linear-gradient(180deg, rgba(34, 197, 94, 0.03), transparent 42%, rgba(0, 0, 0, 0.12));
  opacity: 0.86;
  mix-blend-mode: normal;
  mask-image: linear-gradient(180deg, #000 0%, rgba(0, 0, 0, 0.72) 62%, rgba(0, 0, 0, 0.18) 100%);
}

.kqs-shell::after {
  background: radial-gradient(ellipse at 50% 52%, rgba(0, 0, 0, 0.18), transparent 36rem);
  opacity: 0.72;
}

.kqs-topbar {
  position: sticky;
  top: 0;
  z-index: 30;
  display: grid;
  grid-template-columns: minmax(12rem, 1fr) auto minmax(12rem, 1fr);
  gap: 1rem;
  align-items: center;
  width: min(100% - 1.25rem, 1380px);
  min-height: 3.25rem;
  margin: 0 auto;
  padding: 0.55rem 1rem;
  border: 1px solid var(--kqs-border);
  border-top: 0;
  border-radius: 0 0 8px 8px;
  background: var(--kqs-card-strong);
  box-shadow:
    0 22px 58px rgba(0, 0, 0, 0.3),
    inset 0 1px 0 rgba(34, 197, 94, 0.08);
}

.kqs-brand {
  display: inline-flex;
  min-width: 0;
  align-items: center;
  gap: 0.65rem;
  color: inherit;
  font-weight: 700;
  text-decoration: none;
}

.kqs-brand__logo {
  width: 1.65rem;
  height: 1.65rem;
  object-fit: contain;
}

.kqs-brand__name {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.kqs-nav {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.2rem;
  border: 1px solid var(--kqs-border);
  border-radius: 8px;
  background: rgba(3, 22, 18, 0.72);
  padding: 0.18rem;
}

.kqs-nav__link {
  display: inline-flex;
  align-items: center;
  min-height: 2rem;
  border-radius: 6px;
  padding: 0 0.72rem;
  color: var(--kqs-muted);
  font-size: 0.84rem;
  font-weight: 650;
  text-decoration: none;
  transition: background-color 0.16s ease, color 0.16s ease;
}

.kqs-nav__link:hover,
.kqs-nav__link.router-link-active {
  background: var(--kqs-card-soft);
  color: #edf7ee;
}

.kqs-actions {
  display: flex;
  min-width: 0;
  align-items: center;
  justify-content: flex-end;
  gap: 0.75rem;
}

.kqs-balance {
  overflow: hidden;
  color: #22c55e;
  font-size: 0.84rem;
  font-weight: 700;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.kqs-login-button {
  display: inline-flex;
  min-height: 2rem;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--kqs-border);
  border-radius: 8px;
  background: rgba(3, 22, 18, 0.64);
  padding: 0 0.82rem;
  color: #edf7ee;
  font-size: 0.84rem;
  font-weight: 700;
  text-decoration: none;
  transition: background-color 0.16s ease, border-color 0.16s ease;
}

.kqs-login-button:hover {
  border-color: rgba(43, 132, 83, 0.64);
  background: rgba(34, 197, 94, 0.12);
}

.kqs-admin-button {
  display: inline-flex;
  min-height: 2rem;
  align-items: center;
  justify-content: center;
  border: 1px solid rgba(34, 197, 94, 0.32);
  border-radius: 8px;
  background: var(--kqs-action);
  padding: 0 0.82rem;
  color: #ffffff;
  font-size: 0.84rem;
  font-weight: 750;
  text-decoration: none;
  transition: background-color 0.16s ease, border-color 0.16s ease;
}

.kqs-admin-button:hover {
  border-color: rgba(34, 197, 94, 0.52);
  background: var(--kqs-action-hover);
}

.kqs-main {
  position: relative;
  z-index: 1;
  width: min(100% - 1.25rem, 1380px);
  margin: 0 auto;
  padding: 0.85rem 0 2.25rem;
}

.kqs-main::before {
  content: '';
  position: absolute;
  inset: 0 -0.75rem auto;
  height: 13rem;
  z-index: -1;
  border-radius: 0 0 24px 24px;
  background:
    radial-gradient(circle at 18% 16%, rgba(34, 197, 94, 0.06), transparent 18rem),
    radial-gradient(circle at 72% 12%, rgba(0, 137, 78, 0.06), transparent 16rem),
    linear-gradient(180deg, rgba(34, 197, 94, 0.03), transparent);
  filter: blur(0.2px);
}

.kqs-shell :deep(.kqs-panel),
.kqs-shell :deep(.kqs-tutorial),
.kqs-shell :deep(.purchase-market__hero),
.kqs-shell :deep(.purchase-market__panel),
.kqs-shell :deep(.purchase-market__notice),
.kqs-shell :deep(.purchase-product),
.kqs-shell :deep(.subscription-shop__hero),
.kqs-shell :deep(.subscription-shop__plans),
.kqs-shell :deep(.subscription-shop__side > section),
.kqs-shell :deep(.subscription-shop__hall),
.kqs-shell :deep(.card),
.kqs-shell :deep(.table-container) {
  position: relative;
  overflow: hidden;
  background: var(--kqs-card);
  box-shadow:
    0 24px 68px rgba(0, 0, 0, 0.26),
    inset 0 1px 0 rgba(34, 197, 94, 0.08);
  border-color: var(--kqs-border);
}

.kqs-shell :deep(.kqs-panel),
.kqs-shell :deep(.kqs-tutorial),
.kqs-shell :deep(.purchase-market__hero),
.kqs-shell :deep(.purchase-market__panel),
.kqs-shell :deep(.purchase-market__notice),
.kqs-shell :deep(.purchase-product),
.kqs-shell :deep(.subscription-shop__hero),
.kqs-shell :deep(.subscription-shop__plans),
.kqs-shell :deep(.subscription-shop__side > section),
.kqs-shell :deep(.subscription-shop__hall),
.kqs-shell :deep(.card),
.kqs-shell :deep(.table-container) {
  color: #edf7ee;
}

.kqs-shell :deep(.kqs-panel > *),
.kqs-shell :deep(.kqs-tutorial > *),
.kqs-shell :deep(.purchase-market__hero > *),
.kqs-shell :deep(.purchase-market__panel > *),
.kqs-shell :deep(.purchase-market__notice > *),
.kqs-shell :deep(.purchase-product > *),
.kqs-shell :deep(.subscription-shop__hero > *),
.kqs-shell :deep(.subscription-shop__plans > *),
.kqs-shell :deep(.subscription-shop__side > section > *),
.kqs-shell :deep(.subscription-shop__hall > *),
.kqs-shell :deep(.card > *) {
  position: relative;
  z-index: 1;
}

@media (max-width: 860px) {
  .kqs-topbar {
    grid-template-columns: 1fr auto;
    width: min(100% - 1rem, 1380px);
    border-radius: 0 0 8px 8px;
  }

  .kqs-nav {
    grid-column: 1 / -1;
    order: 3;
    justify-content: flex-start;
    overflow-x: auto;
  }

  .kqs-actions {
    gap: 0.45rem;
  }

  .kqs-balance {
    max-width: 7rem;
  }

  .kqs-main {
    width: min(100% - 1rem, 1380px);
  }
}

@media (max-width: 520px) {
  .kqs-topbar {
    padding: 0.5rem;
  }

  .kqs-brand__name {
    max-width: 9rem;
  }

  .kqs-balance {
    display: none;
  }

  .kqs-admin-button {
    padding: 0 0.65rem;
  }
}
</style>
