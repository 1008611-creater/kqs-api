import { computed, ref, watch } from 'vue'
import { defineStore } from 'pinia'
import { userAPI } from '@/api'
import { useAuthStore } from './auth'
import type {
  UserOnboardingEvent,
  UserOnboardingState,
  UserOnboardingStep
} from '@/api/user'
import { USER_SUBSCRIPTIONS_VISIBLE } from '@/constants/subscriptions'

const STEP_ORDER: UserOnboardingStep[] = USER_SUBSCRIPTIONS_VISIBLE
  ? ['purchase', 'redeem', 'key', 'config']
  : ['redeem', 'key', 'config']
const STORAGE_PREFIX = 'sub2api:user-onboarding:v1'

const createDefaultState = (): UserOnboardingState => ({
  version: 1,
  current_step: STEP_ORDER[0],
  completed_steps: [],
  last_event: '',
  started_at: null,
  redeem_completed_at: null,
  key_created_at: null,
  config_copied_at: null,
  first_request_at: null,
  updated_at: new Date().toISOString(),
  completed: false,
  progress_percent: 0
})

const storageKeyForUser = (userId: number) => `${STORAGE_PREFIX}:${userId}`

const normalizeSteps = (steps: UserOnboardingStep[] = []) => {
  const set = new Set(steps)
  return STEP_ORDER.filter((step) => set.has(step))
}

const deriveNextStep = (completedSteps: UserOnboardingStep[]): UserOnboardingStep => {
  const set = new Set(completedSteps)
  return STEP_ORDER.find((step) => !set.has(step)) || 'config'
}

const withCompletedSteps = (
  state: UserOnboardingState,
  steps: UserOnboardingStep[]
): UserOnboardingState => {
  const completedSteps = normalizeSteps([...state.completed_steps, ...steps])
  return {
    ...state,
    completed_steps: completedSteps,
    current_step: deriveNextStep(completedSteps),
    completed: completedSteps.length === STEP_ORDER.length,
    progress_percent: Math.round((completedSteps.length / STEP_ORDER.length) * 100),
    updated_at: new Date().toISOString()
  }
}

const normalizeState = (input?: Partial<UserOnboardingState> | null): UserOnboardingState => {
  const fallback = createDefaultState()
  const completedSteps = normalizeSteps(input?.completed_steps || [])
  return {
    ...fallback,
    ...input,
    current_step:
      input?.current_step && STEP_ORDER.includes(input.current_step)
        ? input.current_step
        : deriveNextStep(completedSteps),
    completed_steps: completedSteps,
    completed: input?.completed ?? completedSteps.length === STEP_ORDER.length,
    progress_percent:
      typeof input?.progress_percent === 'number'
        ? input.progress_percent
        : Math.round((completedSteps.length / STEP_ORDER.length) * 100),
    updated_at: input?.updated_at || fallback.updated_at
  }
}

const applyEventOptimistically = (
  state: UserOnboardingState,
  event: UserOnboardingEvent
): UserOnboardingState => {
  const now = new Date().toISOString()
  const startedAt = state.started_at || now
  const base = { ...state, started_at: startedAt, last_event: event, updated_at: now }

  switch (event) {
    case 'redeem_success':
      return withCompletedSteps(
        { ...base, redeem_completed_at: base.redeem_completed_at || now },
        USER_SUBSCRIPTIONS_VISIBLE ? ['purchase', 'redeem'] : ['redeem']
      )
    case 'key_created':
      return withCompletedSteps(
        { ...base, key_created_at: base.key_created_at || now },
        USER_SUBSCRIPTIONS_VISIBLE ? ['purchase', 'redeem', 'key'] : ['redeem', 'key']
      )
    case 'config_copied':
      return withCompletedSteps(
        { ...base, config_copied_at: base.config_copied_at || now },
        STEP_ORDER
      )
    case 'first_request':
      return withCompletedSteps(
        { ...base, first_request_at: base.first_request_at || now },
        STEP_ORDER
      )
  }
}

const eventAlreadyApplied = (state: UserOnboardingState, event: UserOnboardingEvent) => {
  switch (event) {
    case 'redeem_success':
      return state.completed_steps.includes('redeem')
    case 'key_created':
      return state.completed_steps.includes('key')
    case 'config_copied':
      return state.config_copied_at != null || state.completed_steps.includes('config')
    case 'first_request':
      return state.first_request_at != null
  }
}

export const useUserOnboardingStore = defineStore('userOnboarding', () => {
  const authStore = useAuthStore()
  const state = ref<UserOnboardingState>(createDefaultState())
  const loading = ref(false)
  const syncing = ref(false)
  const loadedForUserId = ref<number | null>(null)
  const guideReplayNonce = ref(0)

  const currentStep = computed(() => state.value.current_step)
  const completedSteps = computed(() => state.value.completed_steps)
  const progressPercent = computed(() => state.value.progress_percent)
  const isComplete = computed(() => state.value.completed)

  const persistLocal = () => {
    const userId = authStore.user?.id
    if (!userId) return
    try {
      localStorage.setItem(storageKeyForUser(userId), JSON.stringify(state.value))
    } catch {
      // Local persistence is a convenience; backend state remains authoritative.
    }
  }

  const replaceState = (next: UserOnboardingState) => {
    state.value = normalizeState(next)
    persistLocal()
  }

  const loadLocal = (userId: number) => {
    try {
      const cached = localStorage.getItem(storageKeyForUser(userId))
      if (cached) {
        state.value = normalizeState(JSON.parse(cached))
      }
    } catch {
      state.value = createDefaultState()
    }
  }

  const bootstrap = async (force = false) => {
    const userId = authStore.user?.id
    if (!userId) {
      state.value = createDefaultState()
      loadedForUserId.value = null
      return
    }
    if (!force && loadedForUserId.value === userId) {
      return
    }

    loading.value = true
    loadedForUserId.value = userId
    loadLocal(userId)
    try {
      replaceState(await userAPI.getOnboardingState())
    } catch (error) {
      console.warn('Failed to load onboarding state:', error)
    } finally {
      loading.value = false
    }
  }

  const recordEvent = async (event: UserOnboardingEvent) => {
    const userId = authStore.user?.id
    if (!userId || eventAlreadyApplied(state.value, event)) {
      return state.value
    }

    state.value = applyEventOptimistically(state.value, event)
    persistLocal()

    syncing.value = true
    try {
      const updated = await userAPI.updateOnboardingState({ event })
      replaceState(updated)
      return state.value
    } catch (error) {
      console.warn('Failed to update onboarding state:', error)
      return state.value
    } finally {
      syncing.value = false
    }
  }

  const hasCompleted = (step: UserOnboardingStep) => completedSteps.value.includes(step)

  const requestGuideReplay = () => {
    guideReplayNonce.value += 1
  }

  watch(
    () => authStore.user?.id,
    (userId) => {
      if (userId) {
        bootstrap(true)
      } else {
        state.value = createDefaultState()
        loadedForUserId.value = null
      }
    },
    { immediate: true }
  )

  return {
    state,
    loading,
    syncing,
    currentStep,
    completedSteps,
    progressPercent,
    isComplete,
    guideReplayNonce,
    bootstrap,
    recordEvent,
    hasCompleted,
    requestGuideReplay
  }
})
