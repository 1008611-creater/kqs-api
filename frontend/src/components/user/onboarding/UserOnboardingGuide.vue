<template>
  <section class="onboarding-guide" :class="`onboarding-guide--${variant}`">
    <div class="onboarding-guide__glow onboarding-guide__glow--left"></div>
    <div class="onboarding-guide__glow onboarding-guide__glow--right"></div>

    <div class="onboarding-guide__header">
      <div class="min-w-0">
        <p class="onboarding-guide__eyebrow">{{ t('newUserGuide.eyebrow') }}</p>
        <h2 class="onboarding-guide__title">{{ t('newUserGuide.title') }}</h2>
        <p class="onboarding-guide__description">{{ t('newUserGuide.description') }}</p>
      </div>
      <div class="onboarding-guide__side">
        <div class="onboarding-guide__badge">
          <Icon name="sparkles" size="sm" />
          <span>{{ t(`newUserGuide.context.${context}`) }}</span>
        </div>
        <div class="onboarding-guide__progress">
          <span>{{ t('newUserGuide.progress', { done: completedCount, total: stepOrder.length }) }}</span>
          <div class="onboarding-guide__progress-track">
            <i :style="{ width: `${progressPercent}%` }"></i>
          </div>
        </div>
      </div>
    </div>

    <div class="onboarding-guide__steps" :aria-label="t('newUserGuide.stepsLabel')">
      <article
        v-for="(step, index) in steps"
        :key="step.id"
        class="onboarding-step"
        :class="`onboarding-step--${stepStatus(step.id)}`"
      >
        <div class="onboarding-step__index">{{ String(index + 1).padStart(2, '0') }}</div>
        <div class="onboarding-step__icon">
          <Icon :name="stepStatus(step.id) === 'done' ? 'check' : step.icon" size="md" />
        </div>
        <div class="min-w-0 flex-1">
          <h3 class="onboarding-step__title">{{ step.title }}</h3>
          <p class="onboarding-step__description">{{ step.description }}</p>
          <RouterLink v-if="step.route" :to="step.route" class="onboarding-step__action">
            {{ step.action }}
            <Icon name="arrowRight" size="xs" />
          </RouterLink>
          <a
            v-else
            :href="step.href"
            target="_blank"
            rel="noopener noreferrer"
            class="onboarding-step__action"
          >
            {{ step.action }}
            <Icon name="externalLink" size="xs" />
          </a>
        </div>
      </article>
    </div>

    <div v-if="showPackages && USER_SUBSCRIPTIONS_VISIBLE" class="onboarding-packages">
      <div class="onboarding-section-heading">
        <p>{{ t('newUserGuide.packages.title') }}</p>
        <span>{{ t('newUserGuide.packages.subtitle') }}</span>
      </div>
      <div class="onboarding-packages__grid">
        <a
          v-for="pack in packages"
          :key="pack.id"
          :href="pack.url"
          target="_blank"
          rel="noopener noreferrer"
          class="onboarding-package"
          :class="{ 'onboarding-package--recommended': pack.recommended }"
        >
          <div class="flex items-center justify-between gap-3">
            <span class="onboarding-package__name">
              {{ t(`newUserGuide.packages.names.${pack.id}`) }}
            </span>
            <span v-if="pack.recommended" class="onboarding-package__tag">
              {{ t('newUserGuide.packages.recommended') }}
            </span>
          </div>
          <div class="onboarding-package__price">
            <strong>{{ pack.priceCny }}</strong>
            <span>{{ t('newUserGuide.packages.priceUnit') }}</span>
          </div>
          <p class="onboarding-package__balance">
            {{ pack.quotaLabel }} {{ t('newUserGuide.packages.balanceLabel') }}
          </p>
          <span class="onboarding-package__button">
            {{ t('newUserGuide.actions.buyNow') }}
            <Icon name="externalLink" size="xs" />
          </span>
        </a>
      </div>
    </div>

    <div v-if="showConfig" class="onboarding-config">
      <div class="onboarding-section-heading">
        <p>{{ t('newUserGuide.config.title') }}</p>
        <span>{{ t('newUserGuide.config.subtitle') }}</span>
      </div>

      <div class="onboarding-config__grid">
        <div class="onboarding-code-card">
          <div class="onboarding-code-card__top">
            <span>{{ t('newUserGuide.config.configFilename') }}</span>
            <button type="button" class="onboarding-copy-button" @click="copyCodexConfig">
              <Icon name="clipboard" size="xs" />
              {{ t('newUserGuide.config.copyConfig') }}
            </button>
          </div>
          <pre><code>{{ codexConfigDisplay }}</code></pre>
        </div>

        <div class="onboarding-key-card">
          <div class="onboarding-key-card__icon">
            <Icon :name="apiKey ? 'checkCircle' : 'key'" size="lg" />
          </div>
          <h3>
            {{ apiKey ? t('newUserGuide.config.keyReady') : t('newUserGuide.config.noKeyTitle') }}
          </h3>
          <p>
            {{
              apiKey
                ? t('newUserGuide.config.maskedKeyLabel', { key: maskedApiKey })
                : t('newUserGuide.config.noKeyDescription')
            }}
          </p>
          <button
            type="button"
            class="onboarding-key-card__button"
            :disabled="!apiKey"
            @click="copyEnvCommand"
          >
            <Icon name="copy" size="sm" />
            {{ t('newUserGuide.config.copyEnv') }}
          </button>
        </div>
      </div>
    </div>

    <div class="onboarding-guidance">
      <div v-if="hasUngroupedKey" class="onboarding-guidance__item onboarding-guidance__item--warning">
        <Icon name="exclamationTriangle" size="sm" />
        <span>{{ t('newUserGuide.guidance.noGroup') }}</span>
      </div>
      <div class="onboarding-guidance__item">
        <Icon name="terminal" size="sm" />
        <span>{{ t('newUserGuide.guidance.endpoint', { endpoint: responsesEndpoint }) }}</span>
      </div>
      <div class="onboarding-guidance__item">
        <Icon name="dollar" size="sm" />
        <span>{{ t('newUserGuide.guidance.insufficientBalance') }}</span>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import { useClipboard } from '@/composables/useClipboard'
import { useUserOnboardingStore } from '@/stores/userOnboarding'
import { maskApiKey as maskApiKeyValue } from '@/utils/maskApiKey'
import {
  CODEX_API_BASE_URL,
  CODEX_DEFAULT_MODEL,
  CODEX_RESPONSES_ENDPOINT,
  LDXP_PACKAGES,
  LDXP_SHOP_URL
} from '@/constants/onboarding'
import { USER_SUBSCRIPTIONS_VISIBLE } from '@/constants/subscriptions'

type GuideStepId = 'purchase' | 'redeem' | 'key' | 'config'
type GuideContext = 'dashboard' | 'redeem' | 'keys'
type GuideVariant = 'full' | 'compact' | 'modal'
type LocalIcon =
  | 'check'
  | 'creditCard'
  | 'gift'
  | 'key'
  | 'terminal'

interface GuideStep {
  id: GuideStepId
  icon: LocalIcon
  title: string
  description: string
  action: string
  route?: string
  href?: string
}

const props = withDefaults(defineProps<{
  currentStep?: GuideStepId
  context?: GuideContext
  variant?: GuideVariant
  showPackages?: boolean
  showConfig?: boolean
  apiKey?: string
  hasUngroupedKey?: boolean
}>(), {
  currentStep: 'purchase',
  context: 'dashboard',
  variant: 'full',
  showPackages: true,
  showConfig: false,
  apiKey: '',
  hasUngroupedKey: false
})

const { t } = useI18n()
const { copyToClipboard } = useClipboard()
const userOnboardingStore = useUserOnboardingStore()

const packages = LDXP_PACKAGES
const responsesEndpoint = CODEX_RESPONSES_ENDPOINT
const stepOrder = computed<GuideStepId[]>(() =>
  USER_SUBSCRIPTIONS_VISIBLE ? ['purchase', 'redeem', 'key', 'config'] : ['redeem', 'key', 'config']
)

onMounted(() => {
  userOnboardingStore.bootstrap().catch((error) => {
    console.warn('Failed to bootstrap onboarding guide:', error)
  })
})

const steps = computed<GuideStep[]>(() => {
  const items: GuideStep[] = []

  if (USER_SUBSCRIPTIONS_VISIBLE) {
    items.push({
    id: 'purchase',
    icon: 'creditCard',
    title: t('newUserGuide.steps.purchase.title'),
    description: t('newUserGuide.steps.purchase.description'),
    action: t('newUserGuide.actions.buyNow'),
    href: LDXP_SHOP_URL
    })
  }

  items.push(
  {
    id: 'redeem',
    icon: 'gift',
    title: t('newUserGuide.steps.redeem.title'),
    description: t('newUserGuide.steps.redeem.description'),
    action: t('newUserGuide.actions.goRedeem'),
    route: '/redeem'
  },
  {
    id: 'key',
    icon: 'key',
    title: t('newUserGuide.steps.key.title'),
    description: t('newUserGuide.steps.key.description'),
    action: t('newUserGuide.actions.createKey'),
    route: '/keys'
  },
  {
    id: 'config',
    icon: 'terminal',
    title: t('newUserGuide.steps.config.title'),
    description: t('newUserGuide.steps.config.description'),
    action: t('newUserGuide.actions.openKeys'),
    route: '/keys'
  }
  )

  return items
})

const fallbackCompletedSteps = computed<GuideStepId[]>(() => {
  const currentIndex = stepOrder.value.indexOf(props.currentStep)
  if (currentIndex <= 0) return []
  return stepOrder.value.slice(0, currentIndex)
})

const completedSteps = computed<GuideStepId[]>(() =>
  userOnboardingStore.loading
    ? fallbackCompletedSteps.value
    : userOnboardingStore.completedSteps
)
const completedCount = computed(() => completedSteps.value.length)
const progressPercent = computed(() =>
  userOnboardingStore.loading
    ? Math.round((completedSteps.value.length / stepOrder.value.length) * 100)
    : userOnboardingStore.progressPercent
)
const currentStep = computed(() =>
  userOnboardingStore.loading ? props.currentStep : userOnboardingStore.currentStep || props.currentStep
)
const currentIndex = computed(() => Math.max(stepOrder.value.indexOf(currentStep.value), 0))

const stepStatus = (step: GuideStepId) => {
  if (completedSteps.value.includes(step)) return 'done'
  const index = stepOrder.value.indexOf(step)
  if (index === currentIndex.value) return 'current'
  return 'next'
}

const codexConfigDisplay = computed(() => [
  `model_provider = "OpenAI"`,
  `model = "${CODEX_DEFAULT_MODEL}"`,
  `review_model = "${CODEX_DEFAULT_MODEL}"`,
  `model_reasoning_effort = "xhigh"`,
  '',
  '[model_providers.OpenAI]',
  `name = "OpenAI"`,
  `base_url = "${CODEX_API_BASE_URL}"`,
  `wire_api = "responses"`,
  `requires_openai_auth = true`
].join('\n'))

const codexConfigToCopy = computed(() => [
  '# ~/.codex/config.toml',
  codexConfigDisplay.value,
  '',
  '# OPENAI_API_KEY 请使用本页按钮单独复制到终端环境变量，不建议写进 config.toml。'
].join('\n'))

const maskedApiKey = computed(() => props.apiKey ? maskApiKeyValue(props.apiKey) : '')

const envCommand = computed(() => {
  if (!props.apiKey) return ''
  return [
    '# Windows PowerShell',
    `$env:OPENAI_API_KEY = "${props.apiKey}"`,
    '',
    '# macOS / Linux',
    `export OPENAI_API_KEY="${props.apiKey}"`
  ].join('\n')
})

const copyCodexConfig = async () => {
  const success = await copyToClipboard(codexConfigToCopy.value, t('newUserGuide.config.copiedConfig'))
  if (success) {
    userOnboardingStore.recordEvent('config_copied').catch((error) => {
      console.warn('Failed to record config copy:', error)
    })
  }
}

const copyEnvCommand = async () => {
  if (!envCommand.value) return
  const success = await copyToClipboard(envCommand.value, t('newUserGuide.config.copiedEnv'))
  if (success) {
    userOnboardingStore.recordEvent('config_copied').catch((error) => {
      console.warn('Failed to record env copy:', error)
    })
  }
}
</script>

<style scoped>
.onboarding-guide {
  position: relative;
  overflow: hidden;
  border: 1px solid rgba(0, 137, 78, 0.14);
  border-radius: 1.5rem;
  background:
    linear-gradient(135deg, rgba(239, 250, 242, 0.94), rgba(255, 255, 255, 0.84)),
    radial-gradient(circle at 12% 12%, rgba(238, 195, 73, 0.24), transparent 32%);
  box-shadow: 0 18px 52px rgba(2, 43, 18, 0.08);
  padding: 1.15rem;
  backdrop-filter: blur(18px);
}

:global(.dark) .onboarding-guide {
  border-color: rgba(255, 255, 255, 0.1);
  background:
    linear-gradient(135deg, rgba(7, 35, 22, 0.9), rgba(14, 18, 16, 0.86)),
    radial-gradient(circle at 12% 12%, rgba(238, 195, 73, 0.12), transparent 34%);
  box-shadow: 0 18px 52px rgba(0, 0, 0, 0.18);
}

.onboarding-guide--compact {
  border-radius: 1.25rem;
  padding: 1rem;
}

.onboarding-guide--modal {
  border-color: rgba(0, 137, 78, 0.1);
  border-radius: 1.65rem;
  background:
    linear-gradient(135deg, rgba(247, 255, 249, 0.72), rgba(255, 255, 255, 0.5)),
    radial-gradient(circle at 12% 10%, rgba(238, 195, 73, 0.18), transparent 28%);
  box-shadow: none;
  padding: 1rem;
}

:global(.dark) .onboarding-guide--modal {
  border-color: rgba(255, 255, 255, 0.08);
  background:
    linear-gradient(135deg, rgba(13, 44, 28, 0.52), rgba(255, 255, 255, 0.035)),
    radial-gradient(circle at 12% 10%, rgba(238, 195, 73, 0.1), transparent 30%);
  box-shadow: none;
}

.onboarding-guide__glow {
  position: absolute;
  pointer-events: none;
  border-radius: 999px;
  filter: blur(1px);
}

.onboarding-guide__glow--left {
  left: -4rem;
  top: -5rem;
  width: 14rem;
  height: 14rem;
  background: rgba(238, 195, 73, 0.2);
}

.onboarding-guide__glow--right {
  right: -5rem;
  bottom: -6rem;
  width: 18rem;
  height: 18rem;
  background: rgba(0, 137, 78, 0.14);
}

.onboarding-guide__header,
.onboarding-guide__steps,
.onboarding-packages,
.onboarding-config,
.onboarding-guidance {
  position: relative;
  z-index: 1;
}

.onboarding-guide__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 1rem;
}

.onboarding-guide__eyebrow {
  font-size: 0.72rem;
  font-weight: 850;
  letter-spacing: 0.16em;
  text-transform: uppercase;
  color: #00894e;
}

:global(.dark) .onboarding-guide__eyebrow {
  color: #8de0ad;
}

.onboarding-guide__title {
  margin-top: 0.35rem;
  color: #032f1c;
  font-size: clamp(1.15rem, 1.9vw, 1.65rem);
  font-weight: 850;
  letter-spacing: -0.035em;
}

:global(.dark) .onboarding-guide__title {
  color: #f8fff9;
}

.onboarding-guide__description {
  margin-top: 0.35rem;
  max-width: 58rem;
  color: rgba(3, 47, 28, 0.68);
  line-height: 1.7;
}

:global(.dark) .onboarding-guide__description {
  color: rgba(255, 255, 255, 0.62);
}

.onboarding-guide__side {
  display: flex;
  flex-shrink: 0;
  flex-direction: column;
  align-items: flex-end;
  gap: 0.65rem;
  min-width: 12rem;
}

.onboarding-guide__badge {
  display: inline-flex;
  flex-shrink: 0;
  align-items: center;
  gap: 0.45rem;
  border: 1px solid rgba(0, 137, 78, 0.16);
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.64);
  color: #006b31;
  padding: 0.45rem 0.7rem;
  font-size: 0.78rem;
  font-weight: 800;
}

:global(.dark) .onboarding-guide__badge {
  border-color: rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.07);
  color: #c4f0d2;
}

.onboarding-guide__progress {
  width: min(13rem, 100%);
  color: rgba(3, 47, 28, 0.62);
  font-size: 0.72rem;
  font-weight: 800;
  text-align: right;
}

:global(.dark) .onboarding-guide__progress {
  color: rgba(255, 255, 255, 0.6);
}

.onboarding-guide__progress-track {
  margin-top: 0.4rem;
  height: 0.42rem;
  overflow: hidden;
  border-radius: 999px;
  background: rgba(0, 137, 78, 0.12);
}

.onboarding-guide__progress-track i {
  display: block;
  height: 100%;
  border-radius: inherit;
  background: linear-gradient(90deg, #00894e, #eec349);
  transition: width 360ms ease;
}

.onboarding-guide__steps {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 0.75rem;
  margin-top: 1rem;
}

.onboarding-step {
  position: relative;
  display: flex;
  min-height: 8.25rem;
  gap: 0.75rem;
  border: 1px solid rgba(2, 43, 18, 0.08);
  border-radius: 1.1rem;
  background: rgba(255, 255, 255, 0.72);
  padding: 0.85rem;
  transition: transform 180ms ease, border-color 180ms ease, box-shadow 180ms ease;
}

.onboarding-step:hover {
  border-color: rgba(0, 137, 78, 0.2);
  box-shadow: 0 14px 28px rgba(2, 43, 18, 0.07);
  transform: translateY(-1px);
}

:global(.dark) .onboarding-step {
  border-color: rgba(255, 255, 255, 0.08);
  background: rgba(255, 255, 255, 0.055);
}

.onboarding-step--current {
  border-color: rgba(0, 137, 78, 0.28);
  box-shadow: inset 0 0 0 1px rgba(0, 137, 78, 0.08), 0 16px 34px rgba(0, 88, 38, 0.08);
}

.onboarding-step--done {
  border-color: rgba(0, 137, 78, 0.18);
  background: rgba(238, 250, 242, 0.86);
}

:global(.dark) .onboarding-step--done {
  background: rgba(0, 137, 78, 0.1);
}

.onboarding-step__index {
  position: absolute;
  right: 0.75rem;
  top: 0.6rem;
  color: rgba(3, 47, 28, 0.15);
  font-size: 1.25rem;
  font-weight: 900;
  line-height: 1;
}

:global(.dark) .onboarding-step__index {
  color: rgba(255, 255, 255, 0.14);
}

.onboarding-step__icon {
  display: flex;
  width: 2.3rem;
  height: 2.3rem;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  border-radius: 0.85rem;
  background: rgba(0, 137, 78, 0.1);
  color: #006b31;
}

.onboarding-step--current .onboarding-step__icon {
  background: #00894e;
  color: white;
}

.onboarding-step--done .onboarding-step__icon {
  background: rgba(0, 137, 78, 0.16);
}

:global(.dark) .onboarding-step__icon {
  background: rgba(0, 137, 78, 0.18);
  color: #bde8cb;
}

.onboarding-step__title {
  padding-right: 2.4rem;
  color: #032f1c;
  font-size: 0.94rem;
  font-weight: 820;
}

:global(.dark) .onboarding-step__title {
  color: white;
}

.onboarding-step__description {
  margin-top: 0.35rem;
  color: rgba(3, 47, 28, 0.6);
  font-size: 0.82rem;
  line-height: 1.55;
}

:global(.dark) .onboarding-step__description {
  color: rgba(255, 255, 255, 0.58);
}

.onboarding-step__action {
  display: inline-flex;
  align-items: center;
  gap: 0.28rem;
  margin-top: 0.7rem;
  color: #00894e;
  font-size: 0.8rem;
  font-weight: 800;
}

.onboarding-step__action:hover {
  color: #006b31;
}

:global(.dark) .onboarding-step__action {
  color: #9ce8b7;
}

.onboarding-section-heading {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: 0.75rem;
  margin-bottom: 0.75rem;
}

.onboarding-section-heading p {
  color: #032f1c;
  font-weight: 850;
}

.onboarding-section-heading span {
  color: rgba(3, 47, 28, 0.56);
  font-size: 0.82rem;
}

:global(.dark) .onboarding-section-heading p {
  color: white;
}

:global(.dark) .onboarding-section-heading span {
  color: rgba(255, 255, 255, 0.55);
}

.onboarding-packages,
.onboarding-config {
  margin-top: 1rem;
}

.onboarding-packages__grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 0.75rem;
}

.onboarding-package {
  display: block;
  border: 1px solid rgba(2, 43, 18, 0.08);
  border-radius: 1.1rem;
  background: rgba(255, 255, 255, 0.76);
  padding: 0.95rem;
  transition: transform 180ms ease, border-color 180ms ease, box-shadow 180ms ease;
}

.onboarding-package:hover {
  border-color: rgba(0, 137, 78, 0.24);
  box-shadow: 0 16px 34px rgba(2, 43, 18, 0.08);
  transform: translateY(-2px);
}

:global(.dark) .onboarding-package {
  border-color: rgba(255, 255, 255, 0.08);
  background: rgba(255, 255, 255, 0.055);
}

.onboarding-package--recommended {
  border-color: rgba(238, 195, 73, 0.42);
  background: linear-gradient(145deg, rgba(255, 249, 229, 0.92), rgba(255, 255, 255, 0.76));
}

:global(.dark) .onboarding-package--recommended {
  background: linear-gradient(145deg, rgba(238, 195, 73, 0.13), rgba(255, 255, 255, 0.055));
}

.onboarding-package__name {
  color: #032f1c;
  font-size: 0.9rem;
  font-weight: 850;
}

:global(.dark) .onboarding-package__name {
  color: white;
}

.onboarding-package__tag {
  border-radius: 999px;
  background: rgba(238, 195, 73, 0.2);
  color: #8a6412;
  padding: 0.22rem 0.5rem;
  font-size: 0.68rem;
  font-weight: 850;
}

.onboarding-package__price {
  margin-top: 0.8rem;
  color: #032f1c;
}

.onboarding-package__price strong {
  font-size: 1.75rem;
  font-weight: 900;
  letter-spacing: -0.04em;
}

.onboarding-package__price span {
  margin-left: 0.2rem;
  color: rgba(3, 47, 28, 0.56);
  font-size: 0.86rem;
}

:global(.dark) .onboarding-package__price,
:global(.dark) .onboarding-package__price span {
  color: white;
}

.onboarding-package__balance {
  margin-top: 0.15rem;
  color: rgba(3, 47, 28, 0.6);
  font-size: 0.82rem;
}

:global(.dark) .onboarding-package__balance {
  color: rgba(255, 255, 255, 0.58);
}

.onboarding-package__button {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  margin-top: 0.85rem;
  color: #00894e;
  font-size: 0.8rem;
  font-weight: 850;
}

.onboarding-config__grid {
  display: grid;
  grid-template-columns: minmax(0, 1.6fr) minmax(16rem, 0.8fr);
  gap: 0.75rem;
}

.onboarding-code-card,
.onboarding-key-card {
  border: 1px solid rgba(2, 43, 18, 0.08);
  border-radius: 1.1rem;
  background: rgba(255, 255, 255, 0.78);
}

:global(.dark) .onboarding-code-card,
:global(.dark) .onboarding-key-card {
  border-color: rgba(255, 255, 255, 0.08);
  background: rgba(255, 255, 255, 0.055);
}

.onboarding-code-card {
  overflow: hidden;
}

.onboarding-code-card__top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.75rem;
  border-bottom: 1px solid rgba(2, 43, 18, 0.08);
  padding: 0.75rem 0.85rem;
  color: rgba(3, 47, 28, 0.62);
  font-size: 0.78rem;
  font-weight: 800;
}

:global(.dark) .onboarding-code-card__top {
  border-color: rgba(255, 255, 255, 0.08);
  color: rgba(255, 255, 255, 0.62);
}

.onboarding-code-card pre {
  overflow-x: auto;
  padding: 0.9rem;
  color: #d9fbe3;
  background: #072316;
  font-size: 0.78rem;
  line-height: 1.65;
}

.onboarding-copy-button {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  border-radius: 999px;
  background: rgba(0, 137, 78, 0.1);
  color: #006b31;
  padding: 0.35rem 0.6rem;
  font-size: 0.75rem;
  font-weight: 850;
  transition: background 180ms ease;
}

.onboarding-copy-button:hover {
  background: rgba(0, 137, 78, 0.16);
}

:global(.dark) .onboarding-copy-button {
  background: rgba(0, 137, 78, 0.18);
  color: #bde8cb;
}

.onboarding-key-card {
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 1rem;
}

.onboarding-key-card__icon {
  display: flex;
  width: 2.75rem;
  height: 2.75rem;
  align-items: center;
  justify-content: center;
  border-radius: 1rem;
  background: rgba(0, 137, 78, 0.1);
  color: #006b31;
}

:global(.dark) .onboarding-key-card__icon {
  background: rgba(0, 137, 78, 0.18);
  color: #bde8cb;
}

.onboarding-key-card h3 {
  margin-top: 0.8rem;
  color: #032f1c;
  font-weight: 850;
}

.onboarding-key-card p {
  margin-top: 0.35rem;
  color: rgba(3, 47, 28, 0.62);
  font-size: 0.84rem;
  line-height: 1.6;
}

:global(.dark) .onboarding-key-card h3 {
  color: white;
}

:global(.dark) .onboarding-key-card p {
  color: rgba(255, 255, 255, 0.6);
}

.onboarding-key-card__button {
  display: inline-flex;
  width: fit-content;
  align-items: center;
  gap: 0.45rem;
  margin-top: 0.9rem;
  border-radius: 999px;
  background: #00894e;
  color: white;
  padding: 0.55rem 0.8rem;
  font-size: 0.8rem;
  font-weight: 850;
  transition: transform 180ms ease, background 180ms ease, opacity 180ms ease;
}

.onboarding-key-card__button:hover:not(:disabled) {
  background: #006b31;
  transform: translateY(-1px);
}

.onboarding-key-card__button:disabled {
  cursor: not-allowed;
  opacity: 0.48;
}

.onboarding-guidance {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 0.6rem;
  margin-top: 1rem;
}

.onboarding-guidance__item {
  display: flex;
  align-items: flex-start;
  gap: 0.45rem;
  border: 1px solid rgba(2, 43, 18, 0.07);
  border-radius: 0.95rem;
  background: rgba(255, 255, 255, 0.58);
  color: rgba(3, 47, 28, 0.66);
  padding: 0.65rem 0.7rem;
  font-size: 0.78rem;
  line-height: 1.55;
}

.onboarding-guidance__item svg {
  flex-shrink: 0;
  margin-top: 0.12rem;
  color: #00894e;
}

.onboarding-guidance__item--warning svg {
  color: #b87900;
}

:global(.dark) .onboarding-guidance__item {
  border-color: rgba(255, 255, 255, 0.08);
  background: rgba(255, 255, 255, 0.05);
  color: rgba(255, 255, 255, 0.62);
}

@media (max-width: 1100px) {
  .onboarding-guide__steps {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .onboarding-config__grid,
  .onboarding-guidance {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 760px) {
  .onboarding-guide__header,
  .onboarding-section-heading {
    align-items: flex-start;
    flex-direction: column;
  }

  .onboarding-guide__side {
    align-items: flex-start;
    width: 100%;
  }

  .onboarding-guide__progress {
    width: 100%;
    text-align: left;
  }

  .onboarding-guide__steps,
  .onboarding-packages__grid {
    grid-template-columns: 1fr;
  }
}
</style>
