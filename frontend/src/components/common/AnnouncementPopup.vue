<template>
  <Teleport to="body">
    <Transition name="announcement-popup">
      <div
        v-if="announcementStore.currentPopup"
        class="announcement-popup"
        role="dialog"
        aria-modal="true"
        :aria-label="announcementStore.currentPopup.title"
      >
        <div class="announcement-popup__backdrop"></div>

        <div class="announcement-popup__stage">
          <section class="announcement-popup__shell" @click.stop>
            <div class="announcement-popup__ambient announcement-popup__ambient--gold"></div>
            <div class="announcement-popup__ambient announcement-popup__ambient--green"></div>

            <header class="announcement-popup__header">
              <div class="announcement-popup__icon" aria-hidden="true">
                <svg viewBox="0 0 24 24" fill="none">
                  <path
                    d="M12 3.2 18.7 6v5.2c0 4.1-2.7 7.9-6.7 9.6-4-1.7-6.7-5.5-6.7-9.6V6L12 3.2Z"
                    stroke="currentColor"
                    stroke-width="1.8"
                    stroke-linejoin="round"
                  />
                  <path
                    d="M8.7 12.1h6.6M9.7 9.1h4.6M10.5 15.1h3"
                    stroke="currentColor"
                    stroke-width="1.8"
                    stroke-linecap="round"
                  />
                </svg>
              </div>

              <div class="min-w-0 flex-1">
                <div class="announcement-popup__eyebrow">
                  <span>{{ t('announcements.unread') }}</span>
                  <i></i>
                  <time>{{ formatRelativeWithDateTime(announcementStore.currentPopup.created_at) }}</time>
                </div>
                <h2>{{ announcementStore.currentPopup.title }}</h2>
              </div>
            </header>

            <main class="announcement-popup__body">
              <div class="announcement-popup__rail"></div>
              <div
                class="markdown-body announcement-popup__markdown"
                v-html="renderedContent"
              ></div>
            </main>

            <footer class="announcement-popup__footer">
              <div class="announcement-popup__actions">
                <RouterLink
                  v-if="showRedeemAction"
                  to="/redeem"
                  class="announcement-popup__primary"
                  @click="handleDismiss"
                >
                  去兑换福利
                  <svg viewBox="0 0 20 20" fill="none" aria-hidden="true">
                    <path d="M7 4.5 12.5 10 7 15.5" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round" />
                  </svg>
                </RouterLink>
                <button type="button" class="announcement-popup__secondary" @click="handleDismiss">
                  {{ t('announcements.markRead') }}
                </button>
              </div>
            </footer>
          </section>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { marked } from 'marked'
import DOMPurify from 'dompurify'
import { useAnnouncementStore } from '@/stores/announcements'
import { formatRelativeWithDateTime } from '@/utils/format'

const { t } = useI18n()
const announcementStore = useAnnouncementStore()

marked.setOptions({
  breaks: true,
  gfm: true
})

const renderedContent = computed(() => {
  const content = announcementStore.currentPopup?.content
  if (!content) return ''
  const html = marked.parse(content) as string
  return DOMPurify.sanitize(html)
})

const showRedeemAction = computed(() =>
  announcementStore.currentPopup?.title?.includes('专属福利') ?? false
)

function handleDismiss() {
  announcementStore.dismissPopup()
}

watch(
  () => announcementStore.currentPopup,
  (popup) => {
    document.body.style.overflow = popup ? 'hidden' : ''
  }
)

onBeforeUnmount(() => {
  document.body.style.overflow = ''
})
</script>

<style scoped>
.announcement-popup {
  position: fixed;
  inset: 0;
  z-index: 120;
}

.announcement-popup__backdrop {
  position: absolute;
  inset: 0;
  background:
    radial-gradient(circle at 15% 15%, rgba(238, 195, 73, 0.2), transparent 28%),
    radial-gradient(circle at 88% 76%, rgba(0, 137, 78, 0.2), transparent 30%),
    rgba(4, 20, 13, 0.52);
  backdrop-filter: blur(18px);
}

.announcement-popup__stage {
  position: relative;
  display: flex;
  min-height: 100dvh;
  align-items: center;
  justify-content: center;
  padding: 1.25rem;
}

.announcement-popup__shell {
  position: relative;
  width: min(42rem, 100%);
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.56);
  border-radius: 2rem;
  background:
    linear-gradient(135deg, rgba(248, 255, 249, 0.97), rgba(255, 252, 239, 0.93)),
    rgba(255, 255, 255, 0.92);
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.82),
    0 34px 90px rgba(2, 43, 18, 0.28);
}

:global(.dark) .announcement-popup__shell {
  border-color: rgba(255, 255, 255, 0.12);
  background:
    linear-gradient(135deg, rgba(8, 37, 23, 0.97), rgba(23, 27, 18, 0.95)),
    rgba(9, 24, 16, 0.96);
  box-shadow:
    inset 0 1px 0 rgba(255, 255, 255, 0.1),
    0 34px 90px rgba(0, 0, 0, 0.48);
}

.announcement-popup__ambient {
  position: absolute;
  pointer-events: none;
  border-radius: 999px;
  transform: translateZ(0);
}

.announcement-popup__ambient--gold {
  left: -6rem;
  top: -7rem;
  width: 15rem;
  height: 15rem;
  background: rgba(238, 195, 73, 0.24);
}

.announcement-popup__ambient--green {
  right: -7rem;
  bottom: -5rem;
  width: 19rem;
  height: 19rem;
  background: rgba(0, 137, 78, 0.15);
}

.announcement-popup__header,
.announcement-popup__body,
.announcement-popup__footer {
  position: relative;
  z-index: 1;
}

.announcement-popup__header {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1.35rem 1.35rem 1rem;
}

.announcement-popup__icon {
  display: flex;
  width: 3.6rem;
  height: 3.6rem;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  border: 1px solid rgba(0, 137, 78, 0.16);
  border-radius: 1.25rem;
  background:
    linear-gradient(135deg, rgba(255, 255, 255, 0.86), rgba(239, 250, 242, 0.72)),
    rgba(255, 255, 255, 0.72);
  color: #00894e;
  box-shadow: 0 14px 34px rgba(0, 88, 38, 0.12);
}

.announcement-popup__icon svg {
  width: 1.9rem;
  height: 1.9rem;
}

.announcement-popup__eyebrow {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 0.45rem;
  color: rgba(3, 47, 28, 0.58);
  font-size: 0.75rem;
  font-weight: 820;
}

.announcement-popup__eyebrow span {
  border: 1px solid rgba(0, 137, 78, 0.16);
  border-radius: 999px;
  background: rgba(0, 137, 78, 0.09);
  color: #006b31;
  padding: 0.22rem 0.5rem;
}

.announcement-popup__eyebrow i {
  width: 0.28rem;
  height: 0.28rem;
  border-radius: 999px;
  background: rgba(0, 137, 78, 0.35);
}

.announcement-popup__header h2 {
  margin-top: 0.4rem;
  color: #032f1c;
  font-size: clamp(1.4rem, 3vw, 2.1rem);
  font-weight: 900;
  letter-spacing: -0.04em;
  line-height: 1.08;
}

:global(.dark) .announcement-popup__eyebrow,
:global(.dark) .announcement-popup__header h2 {
  color: #f8fff9;
}

:global(.dark) .announcement-popup__eyebrow span {
  border-color: rgba(141, 224, 173, 0.22);
  background: rgba(0, 137, 78, 0.18);
  color: #c4f0d2;
}

.announcement-popup__body {
  display: grid;
  grid-template-columns: 0.25rem minmax(0, 1fr);
  gap: 1rem;
  max-height: min(52dvh, 26rem);
  overflow-y: auto;
  padding: 0.35rem 1.35rem 1.1rem;
}

.announcement-popup__rail {
  border-radius: 999px;
  background: linear-gradient(180deg, #00894e, #eec349);
}

.announcement-popup__markdown {
  color: rgba(3, 47, 28, 0.72);
  line-height: 1.75;
}

:global(.dark) .announcement-popup__markdown {
  color: rgba(255, 255, 255, 0.68);
}

.announcement-popup__markdown :deep(h3) {
  margin: 0 0 0.7rem;
  color: #032f1c;
  font-size: 1.05rem;
  font-weight: 860;
}

:global(.dark) .announcement-popup__markdown :deep(h3) {
  color: #f8fff9;
}

.announcement-popup__markdown :deep(p) {
  margin: 0.65rem 0;
}

.announcement-popup__markdown :deep(strong) {
  color: #006b31;
  font-weight: 900;
}

.announcement-popup__markdown :deep(code) {
  border: 1px solid rgba(0, 137, 78, 0.18);
  border-radius: 0.65rem;
  background: rgba(0, 137, 78, 0.09);
  color: #006b31;
  padding: 0.18rem 0.45rem;
  font-weight: 850;
}

.announcement-popup__markdown :deep(a) {
  color: #00894e;
  font-weight: 850;
  text-decoration: none;
}

.announcement-popup__footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  border-top: 1px solid rgba(0, 137, 78, 0.12);
  background: rgba(255, 255, 255, 0.62);
  padding: 0.95rem 1.15rem 1.05rem;
  backdrop-filter: blur(16px);
}

:global(.dark) .announcement-popup__footer {
  border-top-color: rgba(255, 255, 255, 0.08);
  background: rgba(8, 22, 14, 0.72);
}

.announcement-popup__footer p {
  max-width: 26rem;
  color: rgba(3, 47, 28, 0.58);
  font-size: 0.78rem;
  line-height: 1.55;
}

:global(.dark) .announcement-popup__footer p {
  color: rgba(255, 255, 255, 0.58);
}

.announcement-popup__actions {
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-end;
  gap: 0.7rem;
}

.announcement-popup__primary,
.announcement-popup__secondary {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.45rem;
  border-radius: 999px;
  padding: 0.72rem 1rem;
  font-size: 0.84rem;
  font-weight: 850;
  transition:
    transform 180ms cubic-bezier(0.16, 1, 0.3, 1),
    border-color 180ms ease,
    background 180ms ease,
    color 180ms ease,
    box-shadow 180ms ease;
}

.announcement-popup__primary {
  border: 1px solid #00894e;
  background: #00894e;
  color: white;
  box-shadow: 0 16px 32px rgba(0, 88, 38, 0.18);
}

.announcement-popup__primary svg {
  width: 1rem;
  height: 1rem;
}

.announcement-popup__primary:hover {
  background: #006b31;
}

.announcement-popup__secondary {
  border: 1px solid rgba(0, 137, 78, 0.18);
  background: rgba(255, 255, 255, 0.68);
  color: rgba(3, 47, 28, 0.74);
}

.announcement-popup__secondary:hover {
  border-color: rgba(0, 137, 78, 0.28);
  background: rgba(239, 250, 242, 0.92);
}

.announcement-popup__primary:active,
.announcement-popup__secondary:active {
  transform: scale(0.98) translateY(1px);
}

.announcement-popup-enter-active,
.announcement-popup-leave-active {
  transition: opacity 240ms ease;
}

.announcement-popup-enter-from,
.announcement-popup-leave-to {
  opacity: 0;
}

.announcement-popup-enter-active .announcement-popup__shell,
.announcement-popup-leave-active .announcement-popup__shell {
  transition:
    transform 300ms cubic-bezier(0.16, 1, 0.3, 1),
    opacity 240ms ease;
}

.announcement-popup-enter-from .announcement-popup__shell,
.announcement-popup-leave-to .announcement-popup__shell {
  opacity: 0;
  transform: translateY(1rem) scale(0.985);
}

.announcement-popup__body::-webkit-scrollbar {
  width: 8px;
}

.announcement-popup__body::-webkit-scrollbar-track {
  background: transparent;
}

.announcement-popup__body::-webkit-scrollbar-thumb {
  border-radius: 999px;
  background: rgba(0, 137, 78, 0.28);
}

@media (max-width: 760px) {
  .announcement-popup__stage {
    align-items: flex-end;
    padding: 0.75rem;
  }

  .announcement-popup__shell {
    border-radius: 1.5rem;
  }

  .announcement-popup__header,
  .announcement-popup__footer {
    align-items: stretch;
    flex-direction: column;
  }

  .announcement-popup__actions,
  .announcement-popup__primary,
  .announcement-popup__secondary {
    width: 100%;
  }
}
</style>
