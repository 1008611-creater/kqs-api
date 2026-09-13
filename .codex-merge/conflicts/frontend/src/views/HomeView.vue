<template>
  <!-- Custom Home Content: Full Page Mode -->
<<<<<<< HEAD
  <div v-if="hasHomeContent" class="min-h-screen">
=======
  <div v-if="homeContent" class="min-h-[100dvh]">
>>>>>>> 1157e4cfe6271a81c9dc96df84af6c6a3b22f831
    <!-- iframe mode -->
    <iframe
      v-if="isHomeContentUrl"
      :src="safeHomeContentUrl"
      class="min-h-[100dvh] w-full border-0"
      allowfullscreen
    ></iframe>
    <div v-else v-html="safeHomeContentHtml"></div>
  </div>

<<<<<<< HEAD
  <!-- Compact Home Page -->
  <div
    v-else-if="compactHomeEnabled"
    data-testid="compact-home"
    class="flex min-h-screen flex-col bg-gray-50 text-gray-900 dark:bg-dark-950 dark:text-white"
  >
    <header class="border-b border-gray-200 px-4 py-4 sm:px-6 dark:border-dark-800">
      <nav class="mx-auto flex max-w-5xl flex-wrap items-center justify-between gap-3 sm:gap-4">
        <div class="flex min-w-0 flex-1 items-center gap-3">
          <img
            :src="siteLogo || '/logo.svg'"
            alt="Logo"
            class="h-9 w-9 shrink-0 rounded-lg object-contain"
          />
          <span class="min-w-0 truncate text-base font-semibold">{{ siteName }}</span>
        </div>
        <div class="flex max-w-full shrink-0 flex-wrap items-center justify-end gap-2">
          <LocaleSwitcher />
          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg text-gray-500 hover:bg-gray-100 dark:text-dark-400 dark:hover:bg-dark-800"
            :title="t('home.viewDocs')"
          >
            <Icon name="book" size="md" />
          </a>
          <router-link
            v-if="showModelPlazaEntry"
            to="/model-plaza"
            class="flex h-10 shrink-0 items-center gap-1.5 rounded-lg px-2.5 text-sm font-medium text-gray-500 hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
            :title="t('nav.modelPlaza')"
          >
            <Icon name="grid" size="md" />
            <span class="hidden sm:inline">{{ t('nav.modelPlaza') }}</span>
          </router-link>
          <button
            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg text-gray-500 hover:bg-gray-100 dark:text-dark-400 dark:hover:bg-dark-800"
            :title="isDark ? t('home.switchToLight') : t('home.switchToDark')"
            @click="toggleTheme"
          >
            <Icon v-if="isDark" name="sun" size="md" />
            <Icon v-else name="moon" size="md" />
          </button>
          <router-link
            :to="isAuthenticated ? dashboardPath : '/login'"
            class="inline-flex min-h-10 shrink-0 items-center justify-center rounded-lg bg-gray-900 px-4 py-2 text-sm font-medium text-white hover:bg-gray-800 dark:bg-white dark:text-gray-900 dark:hover:bg-gray-200"
          >
            {{ isAuthenticated ? t('home.dashboard') : t('home.login') }}
          </router-link>
        </div>
      </nav>
    </header>

    <main class="flex min-w-0 flex-1 items-center justify-center px-4 py-16 sm:px-6">
      <div class="min-w-0 max-w-2xl text-center">
        <img
          :src="siteLogo || '/logo.svg'"
          alt="Logo"
          class="mx-auto mb-6 h-20 w-20 rounded-2xl object-contain"
        />
        <h1 class="[overflow-wrap:anywhere] text-3xl font-bold md:text-4xl">{{ siteName }}</h1>
        <p class="mt-4 whitespace-pre-wrap [overflow-wrap:anywhere] text-base text-gray-600 dark:text-dark-300">{{ siteSubtitle }}</p>
        <router-link
          :to="isAuthenticated ? dashboardPath : '/login'"
          class="mt-8 inline-flex min-h-10 items-center justify-center rounded-lg bg-primary-600 px-5 py-2.5 text-sm font-medium text-white hover:bg-primary-700"
        >
          {{ isAuthenticated ? t('home.goToDashboard') : t('home.login') }}
        </router-link>
      </div>
    </main>

    <footer class="min-w-0 border-t border-gray-200 px-4 py-5 text-center text-sm text-gray-500 [overflow-wrap:anywhere] sm:px-6 dark:border-dark-800 dark:text-dark-400">
      &copy; {{ currentYear }} {{ siteName }}
    </footer>
  </div>
=======
  <LoginView v-else-if="!isAuthenticated" />
>>>>>>> 1157e4cfe6271a81c9dc96df84af6c6a3b22f831

  <!-- Default Home Page -->
  <div
    v-else
    class="home-shell relative flex min-h-[100dvh] flex-col overflow-hidden"
  >
    <!-- Header -->
    <header class="relative z-20 px-6 py-4">
      <nav class="mx-auto flex max-w-6xl items-center justify-between">
        <!-- Logo -->
<<<<<<< HEAD
        <div class="flex items-center">
          <div class="h-10 w-10 overflow-hidden rounded-xl shadow-md">
            <img :src="siteLogo || '/logo.svg'" alt="Logo" class="h-full w-full object-contain" />
=======
        <div class="flex items-center gap-4">
          <BrandLogo :src="siteLogo || defaultLogo" alt="Logo" variant="home" />
          <div class="hidden leading-tight sm:block">
            <div class="text-xl font-black tracking-tight text-primary-950 dark:text-white">{{ siteName }}</div>
            <div class="text-[11px] font-semibold uppercase tracking-[0.2em] text-primary-700 dark:text-primary-100/80">
              KQS AI Gateway
            </div>
>>>>>>> 1157e4cfe6271a81c9dc96df84af6c6a3b22f831
          </div>
        </div>

        <!-- Nav Actions -->
        <div class="flex items-center gap-3">
          <!-- Language Switcher -->
          <LocaleSwitcher />

          <!-- Doc Link -->
          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="min-h-11 min-w-11 rounded-lg p-2 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-primary-100/70 dark:hover:bg-white/10 dark:hover:text-white"
            :title="t('home.viewDocs')"
          >
            <Icon name="book" size="md" />
          </a>

<<<<<<< HEAD
          <!-- Model Plaza Link -->
          <router-link
            v-if="showModelPlazaEntry"
            to="/model-plaza"
            class="inline-flex items-center gap-1.5 rounded-lg p-2 text-sm text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
            :title="t('nav.modelPlaza')"
          >
            <Icon name="grid" size="md" />
            <span class="hidden sm:inline">{{ t('nav.modelPlaza') }}</span>
          </router-link>

          <!-- Theme Toggle -->
          <button
            @click="toggleTheme"
            class="rounded-lg p-2 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
            :title="isDark ? t('home.switchToLight') : t('home.switchToDark')"
          >
            <Icon v-if="isDark" name="sun" size="md" />
            <Icon v-else name="moon" size="md" />
          </button>

=======
>>>>>>> 1157e4cfe6271a81c9dc96df84af6c6a3b22f831
          <!-- Login / Dashboard Button -->
          <router-link
            v-if="isAuthenticated"
            :to="dashboardPath"
            class="inline-flex items-center gap-1.5 rounded-full bg-gray-900 py-1 pl-1 pr-2.5 transition-colors hover:bg-gray-800 dark:bg-gray-800 dark:hover:bg-gray-700"
          >
            <span
              class="flex h-5 w-5 items-center justify-center rounded-full bg-gradient-to-br from-primary-400 to-primary-600 text-[10px] font-semibold text-white"
            >
              {{ userInitial }}
            </span>
            <span class="text-xs font-medium text-white">{{ t('home.dashboard') }}</span>
            <svg
              class="h-3 w-3 text-gray-400"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              stroke-width="2"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M4.5 19.5l15-15m0 0H8.25m11.25 0v11.25"
              />
            </svg>
          </router-link>
          <router-link
            v-else
            to="/login"
            class="inline-flex min-h-11 items-center rounded-full bg-primary-700 px-5 py-2 text-sm font-medium text-white transition-colors hover:bg-primary-800 dark:bg-primary-600 dark:hover:bg-primary-500"
          >
            {{ t('home.login') }}
          </router-link>
        </div>
      </nav>
    </header>

    <!-- Main Content -->
    <main class="relative z-10 flex-1 px-6 pb-16 pt-10 sm:pt-14 lg:pt-16">
      <div class="mx-auto max-w-6xl">
        <!-- Hero Section - Left/Right Layout -->
        <div class="mb-16 flex flex-col items-center justify-between gap-12 lg:flex-row lg:gap-16">
          <!-- Left: Text Content -->
          <div class="home-reveal flex-1 text-center lg:text-left">
            <div
              class="mb-5 inline-flex items-center gap-2 rounded-full border border-primary-700/20 bg-white/70 px-4 py-2 text-xs font-semibold text-primary-800 shadow-sm shadow-primary-900/5 backdrop-blur"
            >
              <span class="h-2 w-2 rounded-full bg-accent-400 shadow-[0_0_0_4px_rgba(251,191,36,0.18)]"></span>
              {{ t('home.eyebrow') }}
            </div>
            <h1
              class="mb-4 text-5xl font-black tracking-[-0.055em] text-primary-950 dark:text-[#f1f7f0] md:text-6xl lg:text-7xl"
            >
              {{ siteName }}
            </h1>
            <!-- CTA Button -->
            <div>
              <router-link
                :to="isAuthenticated ? dashboardPath : '/login'"
                class="btn btn-primary min-h-12 px-8 py-3 text-base shadow-lg shadow-primary-900/20 ring-1 ring-primary-300/30"
              >
                {{ isAuthenticated ? t('home.goToDashboard') : t('home.getStarted') }}
                <Icon name="arrowRight" size="md" class="ml-2" :stroke-width="2" />
              </router-link>
            </div>
          </div>

          <!-- Right: Terminal Animation -->
          <div class="home-reveal home-reveal-delay-1 flex flex-1 justify-center lg:justify-end">
            <div class="terminal-container">
              <div class="terminal-window">
                <!-- Window header -->
                <div class="terminal-header">
                  <div class="terminal-buttons">
                    <span class="btn-close"></span>
                    <span class="btn-minimize"></span>
                    <span class="btn-maximize"></span>
                  </div>
                  <span class="terminal-title">terminal</span>
                </div>
                <!-- Terminal content -->
                <div class="terminal-body">
                  <div class="code-line line-1">
                    <span class="code-prompt">$</span>
                    <span class="code-cmd">curl</span>
                    <span class="code-flag">-X POST</span>
                    <span class="code-url">/v1/messages</span>
                  </div>
                  <div class="code-line line-2">
                    <span class="code-comment"># {{ t('home.terminal.routing') }}</span>
                  </div>
                  <div class="code-line line-3">
                    <span class="code-success">200 OK</span>
                    <span class="code-response">{ "content": "Hello!" }</span>
                  </div>
                  <div class="code-line line-4">
                    <span class="code-prompt">$</span>
                    <span class="cursor"></span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Access flow and safeguards -->
        <div class="home-reveal home-reveal-delay-2 grid gap-4 lg:grid-cols-[1.12fr_0.88fr]">
          <section class="relative overflow-hidden rounded-[2rem] border border-primary-700/15 bg-white/70 p-7 shadow-[0_18px_50px_-32px_rgba(2,43,18,0.4)] dark:border-white/10 dark:bg-white/[0.045] sm:p-9">
            <div class="absolute right-0 top-0 h-32 w-32 translate-x-10 -translate-y-12 rounded-full bg-primary-500/10"></div>
            <p class="mb-3 text-xs font-semibold tracking-[0.16em] text-primary-700 dark:text-primary-200">
              {{ t('home.flow.label') }}
            </p>
            <h2 class="max-w-md text-2xl font-bold tracking-tight text-primary-950 dark:text-[#f1f7f0] sm:text-3xl">
              {{ t('home.flow.title') }}
            </h2>
            <div class="mt-8 grid gap-3 sm:grid-cols-3">
              <div class="access-step">
                <span>01</span>
                <strong>{{ t('home.tags.subscriptionToApi') }}</strong>
              </div>
              <div class="access-step">
                <span>02</span>
                <strong>{{ t('home.tags.stickySession') }}</strong>
              </div>
              <div class="access-step">
                <span>03</span>
                <strong>{{ t('home.tags.realtimeBilling') }}</strong>
              </div>
            </div>
          </section>

          <div class="grid gap-4">
            <article class="assurance-panel">
              <Icon name="shield" size="md" class="text-primary-600 dark:text-primary-300" />
              <div>
                <h3>{{ t('home.features.multiAccount') }}</h3>
                <p>{{ t('home.features.multiAccountDesc') }}</p>
              </div>
            </article>
            <article class="assurance-panel assurance-panel-accent">
              <Icon name="chart" size="md" class="text-accent-600 dark:text-accent-300" />
              <div>
                <h3>{{ t('home.features.balanceQuota') }}</h3>
                <p>{{ t('home.features.balanceQuotaDesc') }}</p>
              </div>
            </article>
          </div>
        </div>

      </div>
    </main>

    <!-- Footer -->
    <footer class="relative z-10 border-t border-gray-200/50 px-6 py-8 dark:border-dark-800/50">
      <div
        class="mx-auto flex max-w-6xl flex-col items-center justify-center gap-4 text-center sm:flex-row sm:text-left"
      >
        <p class="text-sm text-gray-500 dark:text-dark-400">
          &copy; {{ currentYear }} {{ siteName }}. {{ t('home.footer.allRightsReserved') }}
        </p>
        <div class="flex items-center gap-4">
          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="text-sm text-gray-500 transition-colors hover:text-gray-700 dark:text-dark-400 dark:hover:text-white"
          >
            {{ t('home.docs') }}
          </a>
          <a
            :href="githubUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="text-sm text-gray-500 transition-colors hover:text-gray-700 dark:text-dark-400 dark:hover:text-white"
          >
            GitHub
          </a>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore, useAppStore } from '@/stores'
import BrandLogo from '@/components/common/BrandLogo.vue'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import Icon from '@/components/icons/Icon.vue'
<<<<<<< HEAD
import { sanitizeUrl } from '@/utils/url'
import { FeatureFlags, isFeatureFlagEnabled } from '@/utils/featureFlags'
=======
import LoginView from '@/views/auth/LoginView.vue'
import { sanitizeHtml } from '@/utils/sanitize'
import { sanitizeUrl } from '@/utils/url'
>>>>>>> 1157e4cfe6271a81c9dc96df84af6c6a3b22f831

const { t } = useI18n()

const authStore = useAuthStore()
const appStore = useAppStore()

// Site settings - directly from appStore (already initialized from injected config)
<<<<<<< HEAD
const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'Sub2API')
const siteLogo = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '', { allowRelative: true, allowDataUrl: true }))
const siteSubtitle = computed(() => appStore.cachedPublicSettings?.site_subtitle || 'AI API Gateway Platform')
const docUrl = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.doc_url || appStore.docUrl || ''))
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')
const hasHomeContent = computed(() => homeContent.value.trim().length > 0)
const compactHomeEnabled = computed(() => appStore.cachedPublicSettings?.compact_home_enabled === true)
const modelPlazaEnabled = computed(() => isFeatureFlagEnabled(FeatureFlags.modelPlaza))
=======
const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || '矿泉水API')
const siteLogo = computed(() => appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '')
const docUrl = computed(() => appStore.cachedPublicSettings?.doc_url || appStore.docUrl || '')
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')
const safeHomeContentUrl = computed(() => sanitizeUrl(homeContent.value))
const safeHomeContentHtml = computed(() => sanitizeHtml(homeContent.value))
const defaultLogo = '/kqs-water-logo.svg?v=20260612'
>>>>>>> 1157e4cfe6271a81c9dc96df84af6c6a3b22f831

// Check if homeContent is a URL (for iframe display)
const isHomeContentUrl = computed(() => {
  return safeHomeContentUrl.value !== ''
})

// GitHub URL
const githubUrl = 'https://github.com/Wei-Shaw/sub2api'

// Auth state
const isAuthenticated = computed(() => authStore.isAuthenticated)
const modelPlazaRequiresAuth = computed(
  () => appStore.cachedPublicSettings?.model_plaza_require_auth === true,
)
const showModelPlazaEntry = computed(
  () => modelPlazaEnabled.value && (isAuthenticated.value || !modelPlazaRequiresAuth.value),
)
const isAdmin = computed(() => authStore.isAdmin)
const dashboardPath = computed(() => isAdmin.value ? '/admin/dashboard' : '/dashboard')
const userInitial = computed(() => {
  const user = authStore.user
  if (!user || !user.email) return ''
  return user.email.charAt(0).toUpperCase()
})

// Current year for footer
const currentYear = computed(() => new Date().getFullYear())

onMounted(() => {
  // Check auth state
  authStore.checkAuth()

  // Ensure public settings are loaded (will use cache if already loaded from injected config)
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
})
</script>

<style scoped>
.home-reveal {
  animation: home-reveal 0.72s cubic-bezier(0.16, 1, 0.3, 1) both;
}

.home-shell {
  background:
    radial-gradient(ellipse at 14% 10%, rgba(34, 197, 94, 0.06), transparent 28rem),
    radial-gradient(ellipse at 78% 74%, rgba(27, 93, 70, 0.28), transparent 34rem),
    linear-gradient(115deg, rgba(255, 255, 255, 0.022), transparent 38%, rgba(255, 255, 255, 0.01) 70%, transparent),
    linear-gradient(135deg, #08231d 0%, #0b342a 48%, #082920 100%);
}

.home-shell::before,
.home-shell::after {
  content: '';
  position: absolute;
  inset: 0;
  pointer-events: none;
}

.home-shell::before {
  z-index: 1;
  background:
    linear-gradient(180deg, rgba(2, 16, 12, 0.3) 0%, rgba(3, 24, 18, 0.32) 44%, rgba(2, 15, 12, 0.46) 100%),
    url('/kqs-mountain-water-gilded-v2.webp?v=20260730-image2') center top / 100% auto no-repeat;
  opacity: 1;
  mix-blend-mode: normal;
  filter: saturate(0.82) contrast(1.05) brightness(0.88);
}

.home-shell::after {
  z-index: 1;
  background:
    radial-gradient(ellipse at 50% 34%, transparent, rgba(2, 12, 9, 0.1) 68%, rgba(2, 12, 9, 0.24) 100%),
    linear-gradient(180deg, rgba(2, 12, 9, 0.03), transparent 44%, rgba(2, 12, 9, 0.12));
}

.home-shell header nav > div:first-child,
.home-shell header nav > div:last-child > a,
.home-shell header nav > div:last-child > button {
  border-color: rgba(43, 132, 83, 0.38);
  background-color: rgba(6, 41, 33, 0.44);
  color: #edf7ee;
}

.home-shell h1,
.home-shell h2,
.home-shell h3 {
  color: #edf7ee;
}

.home-shell header nav,
.home-shell header nav .text-primary-950,
.home-shell header nav .text-primary-700 {
  color: #edf7ee;
}

.home-shell header nav .text-primary-700 {
  color: rgba(220, 239, 224, 0.74);
}

.home-shell p {
  color: rgba(220, 239, 224, 0.7);
}

.home-shell .home-reveal > div:first-child,
.home-shell .access-step,
.home-shell .assurance-panel,
.home-shell section {
  border-radius: 8px;
  border-color: rgba(43, 132, 83, 0.38);
  background: rgba(5, 29, 24, 0.78);
  box-shadow: 0 18px 48px rgba(0, 0, 0, 0.22);
}

.home-shell :deep(.brand-logo) {
  border-color: rgba(43, 132, 83, 0.38);
  border-radius: 8px;
  background: rgba(5, 29, 24, 0.78);
  box-shadow:
    0 18px 48px rgba(0, 0, 0, 0.28),
    inset 0 1px 0 rgba(34, 197, 94, 0.1);
}

.home-shell :deep(.brand-logo::before),
.home-shell :deep(.brand-logo::after) {
  display: none;
}

.home-shell :deep(.brand-logo__plate) {
  border-radius: 6px;
  background: rgba(3, 22, 18, 0.64);
  box-shadow: inset 0 0 0 1px rgba(43, 132, 83, 0.22);
}

.home-reveal-delay-1 {
  animation-delay: 0.1s;
}

.home-reveal-delay-2 {
  animation-delay: 0.2s;
}

.access-step {
  display: flex;
  flex-direction: column;
  gap: 0.55rem;
  min-height: 5.5rem;
  padding: 0.875rem 1rem;
  border: 1px solid rgba(43, 132, 83, 0.38);
  border-radius: 1rem;
  background: rgba(3, 22, 18, 0.64);
  box-shadow: inset 0 1px 0 rgba(34, 197, 94, 0.08);
}

.access-step span {
  color: #22c55e;
  font-size: 0.6875rem;
  font-weight: 700;
  letter-spacing: 0.16em;
}

.access-step strong {
  color: #edf7ee;
  font-size: 0.875rem;
  font-weight: 600;
}

.assurance-panel {
  display: flex;
  align-items: flex-start;
  gap: 1.125rem;
  min-height: 9rem;
  padding: 1.5rem;
  border: 1px solid rgba(43, 132, 83, 0.38);
  border-radius: 8px;
  background: rgba(3, 22, 18, 0.64);
  transition:
    transform 0.28s cubic-bezier(0.16, 1, 0.3, 1),
    border-color 0.28s cubic-bezier(0.16, 1, 0.3, 1);
}

.assurance-panel:hover {
  transform: translateY(-2px);
  border-color: rgba(34, 197, 94, 0.52);
}

.assurance-panel h3 {
  margin-bottom: 0.45rem;
  color: #edf7ee;
  font-size: 1rem;
  font-weight: 650;
}

.assurance-panel p {
  color: rgba(220, 239, 224, 0.68);
  font-size: 0.875rem;
  line-height: 1.7;
}

.assurance-panel-accent {
  border-color: rgba(43, 132, 83, 0.42);
  background: rgba(5, 29, 24, 0.78);
}

:global(.dark) .access-step {
  border-color: rgba(142, 217, 164, 0.13);
  background: rgba(142, 217, 164, 0.055);
}

:global(.dark) .access-step span {
  color: #8ed9a4;
}

:global(.dark) .access-step strong,
:global(.dark) .assurance-panel h3 {
  color: #f1f7f0;
}

:global(.dark) .assurance-panel {
  border-color: rgba(142, 217, 164, 0.12);
  background: rgba(7, 38, 31, 0.72);
}

:global(.dark) .assurance-panel-accent {
  border-color: rgba(252, 211, 77, 0.14);
  background: rgba(252, 211, 77, 0.04);
}

:global(.dark) .assurance-panel p {
  color: rgba(220, 244, 228, 0.66);
}

@keyframes home-reveal {
  from {
    opacity: 0;
    transform: translateY(14px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Terminal Container */
.terminal-container {
  position: relative;
  display: inline-block;
}

.terminal-container::before {
  content: '';
  position: absolute;
  inset: -24px;
  z-index: -1;
  border-radius: 28px;
  background:
    linear-gradient(135deg, rgba(34, 197, 94, 0.08), transparent 38%),
    linear-gradient(315deg, rgba(10, 76, 57, 0.2), transparent 42%);
  filter: blur(8px);
}

/* Terminal Window */
.terminal-window {
  width: min(420px, calc(100vw - 3rem));
  background: #082b23;
  border-radius: 14px;
  box-shadow:
    0 25px 50px -12px rgba(0, 0, 0, 0.4),
    0 0 0 1px rgba(43, 132, 83, 0.38),
    inset 0 1px 0 rgba(255, 255, 255, 0.08);
  overflow: hidden;
  transform: perspective(1000px) rotateX(2deg) rotateY(-2deg);
  transition: transform 0.3s ease;
}

.terminal-window:hover {
  transform: perspective(1000px) rotateX(0deg) rotateY(0deg) translateY(-4px);
}

/* Terminal Header */
.terminal-header {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  background: rgba(8, 23, 14, 0.78);
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.terminal-buttons {
  display: flex;
  gap: 8px;
}

.terminal-buttons span {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.btn-close {
  background: #ef4444;
}
.btn-minimize {
  background: #eab308;
}
.btn-maximize {
  background: #22c55e;
}

.terminal-title {
  flex: 1;
  text-align: center;
  font-size: 12px;
  font-family: ui-monospace, monospace;
  color: #86efac;
  margin-right: 52px;
}

/* Terminal Body */
.terminal-body {
  padding: 20px 24px;
  font-family: ui-monospace, 'Fira Code', monospace;
  font-size: 14px;
  line-height: 2;
}

.code-line {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  opacity: 0;
  animation: line-appear 0.5s ease forwards;
}

.line-1 {
  animation-delay: 0.3s;
}
.line-2 {
  animation-delay: 1s;
}
.line-3 {
  animation-delay: 1.8s;
}
.line-4 {
  animation-delay: 2.5s;
}

@keyframes line-appear {
  from {
    opacity: 0;
    transform: translateY(5px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.code-prompt {
  color: #22c55e;
  font-weight: bold;
}
.code-cmd {
  color: #b7ebcb;
}
.code-flag {
  color: #fcd34d;
}
.code-url {
  color: #34d399;
}
.code-comment {
  color: #9ca3af;
  font-style: italic;
}
.code-success {
  color: #86efac;
  background: rgba(0, 135, 60, 0.22);
  padding: 2px 8px;
  border-radius: 4px;
  font-weight: 600;
}
.code-response {
  color: #fcd34d;
}

/* Blinking Cursor */
.cursor {
  display: inline-block;
  width: 8px;
  height: 16px;
  background: #2aa85b;
  animation: blink 1s step-end infinite;
}

@keyframes blink {
  0%,
  50% {
    opacity: 1;
  }
  51%,
  100% {
    opacity: 0;
  }
}

/* Dark mode adjustments */
:global(.dark) .terminal-window {
  box-shadow:
    0 25px 50px -12px rgba(0, 0, 0, 0.6),
    0 0 0 1px rgba(0, 135, 60, 0.2),
    0 0 40px rgba(0, 135, 60, 0.12),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
}

@media (prefers-reduced-motion: reduce) {
  .home-reveal,
  .code-line,
  .cursor {
    animation: none;
    opacity: 1;
    transform: none;
  }

  .terminal-window,
  .assurance-panel {
    transition: none;
  }
}
</style>
