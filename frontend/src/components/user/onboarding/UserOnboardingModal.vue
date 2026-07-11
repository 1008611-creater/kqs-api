<template>
  <Teleport to="body">
    <Transition name="onboarding-modal">
      <div
        v-if="visible"
        class="onboarding-modal"
        role="dialog"
        aria-modal="true"
        :aria-label="t('newUserGuide.modal.ariaLabel')"
      >
        <div class="onboarding-modal__backdrop"></div>
        <div class="onboarding-modal__stage">
          <div class="onboarding-modal__shell">
            <div class="onboarding-modal__ambient onboarding-modal__ambient--gold"></div>
            <div class="onboarding-modal__ambient onboarding-modal__ambient--green"></div>

            <div class="onboarding-modal__content">
              <UserOnboardingGuide
                :current-step="currentStep"
                :context="context"
                variant="modal"
                :show-packages="true"
                :show-config="true"
              />
            </div>

            <div class="onboarding-modal__footer">
              <div class="onboarding-modal__footer-copy">
                <p>{{ t('newUserGuide.modal.footerTitle') }}</p>
                <span>{{ t('newUserGuide.modal.footerDescription') }}</span>
              </div>
              <div class="onboarding-modal__actions">
                <button
                  type="button"
                  class="onboarding-modal__secondary"
                  :class="{ 'onboarding-modal__secondary--selected': dontRemindAgain }"
                  :aria-pressed="dontRemindAgain"
                  @click="toggleDontRemind"
                >
                  <Icon :name="dontRemindAgain ? 'checkCircle' : 'bell'" size="sm" />
                  {{
                    dontRemindAgain
                      ? t('newUserGuide.modal.dontRemindSelected')
                      : t('newUserGuide.modal.dontRemind')
                  }}
                </button>
                <button type="button" class="onboarding-modal__primary" @click="acknowledge">
                  {{ t('newUserGuide.modal.acknowledge') }}
                  <Icon name="arrowRight" size="sm" />
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import UserOnboardingGuide from './UserOnboardingGuide.vue'
import { useAnnouncementStore } from '@/stores/announcements'
import { useUserOnboardingStore } from '@/stores/userOnboarding'
import { USER_SUBSCRIPTIONS_VISIBLE } from '@/constants/subscriptions'

type GuideStepId = 'purchase' | 'redeem' | 'key' | 'config'
type GuideContext = 'dashboard' | 'redeem' | 'keys'

const MODAL_VERSION = 'v2'
const AUTO_SEEN_KEY = `kqs-api:new-user-guide-modal:${MODAL_VERSION}:seen`
const LEGACY_DONT_REMIND_KEY = `kqs-api:new-user-guide-modal:${MODAL_VERSION}:dont-remind`

const route = useRoute()
const { t } = useI18n()
const onboardingStore = useUserOnboardingStore()
const announcementStore = useAnnouncementStore()

const visible = ref(false)
const dontRemindAgain = ref(false)
const pendingAutoOpen = ref(false)

const routePath = computed(() => route.path)

const context = computed<GuideContext>(() => {
  if (routePath.value.startsWith('/redeem')) return 'redeem'
  if (routePath.value.startsWith('/keys')) return 'keys'
  return 'dashboard'
})

const currentStep = computed<GuideStepId>(() => {
  if (routePath.value.startsWith('/redeem')) return 'redeem'
  if (routePath.value.startsWith('/keys')) return 'config'
  return USER_SUBSCRIPTIONS_VISIBLE ? 'purchase' : 'redeem'
})

const canShowAfterAnnouncements = computed(() =>
  !announcementStore.loading && !announcementStore.hasPopupPending
)

const readFlag = (storage: Storage, key: string) => storage.getItem(key) === '1'
const writeFlag = (storage: Storage, key: string) => storage.setItem(key, '1')

const openGuide = async (force = false) => {
  try {
    await onboardingStore.bootstrap()
  } catch (error) {
    console.warn('Failed to bootstrap onboarding modal:', error)
  }

  if (!force) {
    try {
      await announcementStore.fetchAnnouncements(true)
    } catch (error) {
      console.warn('Failed to fetch announcements before onboarding modal:', error)
    }
  }

  if (!force) {
    try {
      if (readFlag(localStorage, AUTO_SEEN_KEY) || readFlag(localStorage, LEGACY_DONT_REMIND_KEY)) {
        return
      }
    } catch {
      // Storage can be unavailable in strict privacy modes; showing the guide is safer.
    }

    if (onboardingStore.isComplete) {
      return
    }

    if (!canShowAfterAnnouncements.value) {
      pendingAutoOpen.value = true
      return
    }
  }

  window.setTimeout(() => {
    if (force || canShowAfterAnnouncements.value) {
      visible.value = true
    } else {
      pendingAutoOpen.value = true
    }
  }, force ? 0 : 320)
}

onMounted(() => {
  openGuide()
})

watch(
  () => onboardingStore.guideReplayNonce,
  () => {
    openGuide(true)
  }
)

watch(canShowAfterAnnouncements, (canShow) => {
  if (!canShow || !pendingAutoOpen.value) return
  pendingAutoOpen.value = false
  window.setTimeout(() => {
    if (!visible.value && canShowAfterAnnouncements.value) {
      visible.value = true
    }
  }, 180)
})

const toggleDontRemind = () => {
  dontRemindAgain.value = !dontRemindAgain.value
}

const acknowledge = () => {
  try {
    writeFlag(localStorage, AUTO_SEEN_KEY)
  } catch {
    // Closing should still work even if storage writes fail.
  }
  visible.value = false
}
</script>

<style scoped>
.onboarding-modal {
  position: fixed;
  inset: 0;
  z-index: 70;
}

.onboarding-modal__backdrop {
  position: absolute;
  inset: 0;
  background:
    radial-gradient(circle at 12% 12%, rgba(238, 195, 73, 0.18), transparent 28%),
    radial-gradient(circle at 86% 74%, rgba(0, 137, 78, 0.18), transparent 30%),
    rgba(4, 20, 13, 0.46);
  backdrop-filter: blur(18px);
}

.onboarding-modal__stage {
  position: relative;
  display: flex;
  max-height: 100dvh;
  min-height: 100dvh;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  padding: 1.25rem;
}

.onboarding-modal__shell {
  position: relative;
  display: flex;
  flex-direction: column;
  width: min(72rem, 100%);
  max-height: min(calc(100dvh - 2.5rem), 58rem);
  min-height: 0;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.54);
  border-radius: 2rem;
  background:
    linear-gradient(135deg, rgba(248, 255, 249, 0.96), rgba(255, 252, 239, 0.92)),
    rgba(255, 255, 255, 0.92);
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.78),
    0 34px 90px rgba(2, 43, 18, 0.26);
}

:global(.dark) .onboarding-modal__shell {
  border-color: rgba(255, 255, 255, 0.12);
  background:
    linear-gradient(135deg, rgba(8, 37, 23, 0.97), rgba(23, 27, 18, 0.95)),
    rgba(9, 24, 16, 0.96);
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.1),
    0 34px 90px rgba(0, 0, 0, 0.48);
}

.onboarding-modal__ambient {
  position: absolute;
  pointer-events: none;
  border-radius: 999px;
  opacity: 0.74;
  transform: translateZ(0);
}

.onboarding-modal__ambient--gold {
  left: -8rem;
  top: -8rem;
  width: 18rem;
  height: 18rem;
  background: rgba(238, 195, 73, 0.24);
}

.onboarding-modal__ambient--green {
  right: -7rem;
  bottom: 3rem;
  width: 22rem;
  height: 22rem;
  background: rgba(0, 137, 78, 0.13);
}

.onboarding-modal__content {
  position: relative;
  flex: 1 1 auto;
  min-height: 0;
  overflow-y: auto;
  overscroll-behavior: contain;
  padding: 1rem 1rem 0.25rem;
  -webkit-overflow-scrolling: touch;
}

.onboarding-modal__footer {
  position: relative;
  display: flex;
  flex: 0 0 auto;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  border-top: 1px solid rgba(0, 137, 78, 0.12);
  background: rgba(255, 255, 255, 0.62);
  padding: 0.95rem 1rem 1rem;
  backdrop-filter: blur(16px);
}

:global(.dark) .onboarding-modal__footer {
  border-top-color: rgba(255, 255, 255, 0.08);
  background: rgba(8, 22, 14, 0.72);
}

.onboarding-modal__footer-copy p {
  color: #032f1c;
  font-size: 0.9rem;
  font-weight: 850;
}

.onboarding-modal__footer-copy span {
  margin-top: 0.2rem;
  display: block;
  color: rgba(3, 47, 28, 0.58);
  font-size: 0.78rem;
}

:global(.dark) .onboarding-modal__footer-copy p {
  color: #f7fff8;
}

:global(.dark) .onboarding-modal__footer-copy span {
  color: rgba(255, 255, 255, 0.58);
}

.onboarding-modal__actions {
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-end;
  gap: 0.7rem;
}

.onboarding-modal__secondary,
.onboarding-modal__primary {
  display: inline-flex;
  min-width: 0;
  align-items: center;
  justify-content: center;
  gap: 0.45rem;
  border-radius: 999px;
  padding: 0.72rem 1rem;
  font-size: 0.84rem;
  font-weight: 850;
  line-height: 1.2;
  text-align: center;
  transition:
    transform 180ms cubic-bezier(0.16, 1, 0.3, 1),
    border-color 180ms ease,
    background 180ms ease,
    color 180ms ease,
    box-shadow 180ms ease;
}

.onboarding-modal__secondary {
  border: 1px solid rgba(0, 137, 78, 0.18);
  background: rgba(255, 255, 255, 0.68);
  color: rgba(3, 47, 28, 0.74);
}

.onboarding-modal__secondary:hover {
  border-color: rgba(0, 137, 78, 0.28);
  background: rgba(239, 250, 242, 0.92);
}

.onboarding-modal__secondary--selected {
  border-color: rgba(0, 137, 78, 0.42);
  background: #00894e;
  color: white;
  box-shadow: 0 14px 28px rgba(0, 88, 38, 0.16);
}

.onboarding-modal__primary {
  border: 1px solid #00894e;
  background: #00894e;
  color: white;
  box-shadow: 0 16px 32px rgba(0, 88, 38, 0.18);
}

.onboarding-modal__primary:hover {
  background: #006b31;
}

.onboarding-modal__secondary:active,
.onboarding-modal__primary:active {
  transform: scale(0.98) translateY(1px);
}

:global(.dark) .onboarding-modal__secondary {
  border-color: rgba(255, 255, 255, 0.12);
  background: rgba(255, 255, 255, 0.08);
  color: rgba(255, 255, 255, 0.72);
}

:global(.dark) .onboarding-modal__secondary--selected {
  border-color: rgba(141, 224, 173, 0.52);
  background: #00894e;
  color: white;
}

.onboarding-modal-enter-active,
.onboarding-modal-leave-active {
  transition: opacity 240ms ease;
}

.onboarding-modal-enter-from,
.onboarding-modal-leave-to {
  opacity: 0;
}

.onboarding-modal-enter-active .onboarding-modal__shell,
.onboarding-modal-leave-active .onboarding-modal__shell {
  transition:
    transform 300ms cubic-bezier(0.16, 1, 0.3, 1),
    opacity 240ms ease;
}

.onboarding-modal-enter-from .onboarding-modal__shell,
.onboarding-modal-leave-to .onboarding-modal__shell {
  opacity: 0;
  transform: translateY(1rem) scale(0.985);
}

@media (max-width: 760px) {
  .onboarding-modal__stage {
    align-items: flex-end;
    padding:
      max(0.55rem, env(safe-area-inset-top))
      max(0.55rem, env(safe-area-inset-right))
      max(0.55rem, env(safe-area-inset-bottom))
      max(0.55rem, env(safe-area-inset-left));
  }

  .onboarding-modal__shell {
    width: 100%;
    max-height: calc(100svh - 1.1rem - env(safe-area-inset-top) - env(safe-area-inset-bottom));
    max-height: calc(100dvh - 1.1rem - env(safe-area-inset-top) - env(safe-area-inset-bottom));
    border-radius: 1.35rem;
  }

  .onboarding-modal__content {
    padding: 0.6rem 0.6rem 0.15rem;
  }

  .onboarding-modal__footer {
    align-items: stretch;
    flex-direction: column;
    gap: 0.7rem;
    padding: 0.75rem 0.7rem max(0.75rem, env(safe-area-inset-bottom));
  }

  .onboarding-modal__footer-copy p {
    font-size: 0.86rem;
    line-height: 1.35;
  }

  .onboarding-modal__footer-copy span {
    font-size: 0.72rem;
    line-height: 1.55;
  }

  .onboarding-modal__actions {
    display: grid;
    grid-template-columns: repeat(2, minmax(0, 1fr));
    gap: 0.55rem;
  }

  .onboarding-modal__secondary,
  .onboarding-modal__primary {
    min-height: 2.65rem;
    padding: 0.64rem 0.72rem;
    font-size: 0.8rem;
  }
}

@media (max-width: 360px) {
  .onboarding-modal__actions {
    grid-template-columns: 1fr;
  }
}
</style>
