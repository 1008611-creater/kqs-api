<template>
  <div class="auth-shell">
    <main class="auth-main">
      <section class="auth-panel">
        <router-link to="/home" class="auth-brand">
          <BrandLogo :src="siteLogo || defaultLogo" :alt="siteName" variant="auth" />
          <span>
            <strong>{{ siteName }}</strong>
          </span>
        </router-link>

        <div class="auth-content">
          <slot />
        </div>

        <div class="auth-footer">
          <slot name="footer" />
        </div>

<<<<<<< HEAD
    <!-- Content Container -->
    <div class="relative z-10 w-full max-w-md">
      <!-- Logo/Brand -->
      <div class="mb-8 text-center">
        <!-- Custom Logo or Default Logo -->
        <template v-if="settingsLoaded">
          <div
            class="mb-4 inline-flex h-16 w-16 items-center justify-center overflow-hidden rounded-2xl shadow-lg shadow-primary-500/30"
          >
            <img :src="siteLogo || '/logo.svg'" alt="Logo" class="h-full w-full object-contain" />
          </div>
          <h1 class="text-gradient mb-2 text-3xl font-bold">
            {{ siteName }}
          </h1>
          <p class="text-sm text-gray-500 dark:text-dark-400">
            {{ siteSubtitle }}
          </p>
        </template>
      </div>

      <!-- Card Container -->
      <div class="card-glass rounded-2xl p-8 shadow-glass">
        <slot />
      </div>

      <!-- Footer Links -->
      <div class="mt-6 text-center text-sm">
        <slot name="footer" />
      </div>

      <!-- Copyright -->
      <div class="mt-8 text-center text-xs text-gray-400 dark:text-dark-500">
        &copy; {{ currentYear }} {{ siteName }}. All rights reserved.
      </div>
    </div>
=======
        <div class="auth-copyright">
          &copy; {{ currentYear }} {{ siteName }}. All rights reserved.
        </div>
      </section>
    </main>
>>>>>>> 1157e4cfe6271a81c9dc96df84af6c6a3b22f831
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useAppStore } from '@/stores'
import BrandLogo from '@/components/common/BrandLogo.vue'
import { sanitizeUrl } from '@/utils/url'

const appStore = useAppStore()

const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || '矿泉水API')
const siteLogo = computed(() => sanitizeUrl(appStore.siteLogo || '', { allowRelative: true, allowDataUrl: true }))
const defaultLogo = '/kqs-water-logo.svg?v=20260612'
const currentYear = computed(() => new Date().getFullYear())

onMounted(() => {
  appStore.fetchPublicSettings()
})
</script>

<style scoped>
.auth-shell {
  position: relative;
  min-height: 100dvh;
  overflow: hidden;
  background:
    radial-gradient(ellipse at 18% 12%, rgba(34, 197, 94, 0.06), transparent 28rem),
    radial-gradient(ellipse at 84% 78%, rgba(27, 93, 70, 0.28), transparent 34rem),
    linear-gradient(115deg, rgba(255, 255, 255, 0.018), transparent 38%, rgba(255, 255, 255, 0.008) 70%, transparent),
    linear-gradient(135deg, #08231d 0%, #0b342a 44%, #082920 100%);
}

.auth-shell::before,
.auth-shell::after {
  content: '';
  position: absolute;
  inset: 0;
  pointer-events: none;
}

.auth-shell::before {
  background:
    linear-gradient(180deg, rgba(2, 16, 12, 0.36) 0%, rgba(3, 24, 18, 0.4) 46%, rgba(2, 15, 12, 0.52) 100%),
    url('/kqs-mountain-water-gilded-v2.webp?v=20260730-image2') center top / 100% auto no-repeat;
  opacity: 1;
  filter: saturate(0.82) contrast(1.05) brightness(0.86);
}

.auth-shell::after {
  background: radial-gradient(ellipse at 50% 42%, rgba(2, 12, 9, 0.04), rgba(2, 12, 9, 0.14) 70%, rgba(2, 12, 9, 0.3) 100%);
}

.auth-main {
  position: relative;
  z-index: 1;
  display: flex;
  min-height: 100dvh;
  align-items: center;
  justify-content: center;
  padding: clamp(1rem, 4vw, 2.5rem);
}

.auth-panel {
  width: min(100%, 28rem);
  color: #edf7ee;
}

.auth-brand {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.9rem;
  margin-bottom: 1.6rem;
  text-decoration: none;
}

.auth-brand span {
  min-width: 0;
}

.auth-brand strong {
  display: block;
}

.auth-brand strong {
  color: #edf7ee;
  font-size: 1.35rem;
  font-weight: 900;
  letter-spacing: 0;
  line-height: 1.1;
}

.auth-brand :deep(.brand-logo) {
  border-color: rgba(43, 132, 83, 0.38);
  border-radius: 8px;
  background: rgba(5, 29, 24, 0.78);
  box-shadow: 0 14px 34px rgba(0, 0, 0, 0.24);
}

.auth-brand :deep(.brand-logo::before),
.auth-brand :deep(.brand-logo::after) {
  display: none;
}

.auth-brand :deep(.brand-logo__plate) {
  border-radius: 6px;
  background: transparent;
}

.auth-content {
  color: #edf7ee;
  border-radius: 8px;
  padding: clamp(1.3rem, 4vw, 1.9rem);
  background: rgba(5, 29, 24, 0.78);
  box-shadow: 0 18px 54px rgba(0, 0, 0, 0.26);
}

.auth-content :deep(h2) {
  color: #edf7ee;
  letter-spacing: 0;
}

.auth-content :deep(p),
.auth-content :deep(.input-label),
.auth-content :deep(.input-hint) {
  color: rgba(220, 239, 224, 0.68);
}

.auth-content :deep(.input) {
  min-height: 2.85rem;
  border-radius: 8px;
  border-color: rgba(43, 132, 83, 0.38);
  background: rgba(2, 18, 15, 0.72);
  color: #edf7ee;
}

.auth-content :deep(.input::placeholder) {
  color: rgba(220, 239, 224, 0.42);
}

.auth-content :deep(.btn-primary) {
  border-radius: 8px;
}

.auth-footer {
  margin-top: 1rem;
  text-align: center;
  color: rgba(220, 239, 224, 0.66);
  font-size: 0.875rem;
}

.auth-copyright {
  margin-top: 1.3rem;
  text-align: center;
  color: rgba(220, 239, 224, 0.36);
  font-size: 0.75rem;
}

@media (max-width: 640px) {
  .auth-main {
    align-items: flex-start;
    padding: 1rem;
  }

  .auth-brand {
    margin-top: 0.5rem;
  }

  .auth-content {
    padding: 1.15rem;
  }
}
</style>
