<template>
  <KqsUserShell v-if="!isAdminRoute">
    <slot />
  </KqsUserShell>
  <div v-else class="relative min-h-[100dvh] overflow-x-hidden bg-[#f5f8f0] text-gray-900 dark:bg-[#06140d] dark:text-gray-100">
    <div class="pointer-events-none fixed inset-0 bg-[radial-gradient(circle_at_18%_14%,rgba(0,135,60,0.07),transparent_24%),radial-gradient(circle_at_82%_10%,rgba(212,160,23,0.08),transparent_18%),linear-gradient(180deg,rgba(255,255,255,0.54),rgba(255,255,255,0.2))] dark:bg-[radial-gradient(circle_at_18%_14%,rgba(0,162,74,0.15),transparent_24%),radial-gradient(circle_at_82%_10%,rgba(212,160,23,0.08),transparent_18%),linear-gradient(180deg,rgba(8,26,18,0.95),rgba(6,20,13,0.98))]"></div>
    <div class="pointer-events-none fixed inset-0 bg-[linear-gradient(rgba(0,135,60,0.035)_1px,transparent_1px),linear-gradient(90deg,rgba(0,135,60,0.035)_1px,transparent_1px)] bg-[size:72px_72px] opacity-60 dark:bg-[linear-gradient(rgba(255,255,255,0.04)_1px,transparent_1px),linear-gradient(90deg,rgba(255,255,255,0.04)_1px,transparent_1px)]"></div>

    <!-- Sidebar -->
    <AppSidebar />

    <!-- Main Content Area -->
    <div
      class="relative min-h-screen transition-all duration-300"
      :class="[sidebarCollapsed ? 'lg:ml-[72px]' : 'lg:ml-64']"
    >
      <!-- Header -->
      <AppHeader />

      <!-- Main Content -->
      <main class="p-4 md:p-6 lg:p-8">
        <slot />
      </main>
    </div>

    <UserOnboardingModal v-if="showUserOnboardingModal" />
  </div>
</template>

<script setup lang="ts">
import '@/styles/onboarding.css'
import { computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAppStore } from '@/stores'
import { useAuthStore } from '@/stores/auth'
import { useOnboardingTour } from '@/composables/useOnboardingTour'
import { useOnboardingStore } from '@/stores/onboarding'
import AppSidebar from './AppSidebar.vue'
import AppHeader from './AppHeader.vue'
import KqsUserShell from './KqsUserShell.vue'
import UserOnboardingModal from '@/components/user/onboarding/UserOnboardingModal.vue'

const appStore = useAppStore()
const authStore = useAuthStore()
const route = useRoute()
const sidebarCollapsed = computed(() => appStore.sidebarCollapsed)
const isAdmin = computed(() => authStore.user?.role === 'admin')
const isAdminRoute = computed(() => route.path.startsWith('/admin'))
const showUserOnboardingModal = computed(() => authStore.isAuthenticated && !isAdminRoute.value)

const { replayTour } = useOnboardingTour({
  storageKey: isAdmin.value ? 'admin_guide' : 'user_guide',
  autoStart: true
})

const onboardingStore = useOnboardingStore()

onMounted(() => {
  onboardingStore.setReplayCallback(replayTour)
})

defineExpose({ replayTour })
</script>
