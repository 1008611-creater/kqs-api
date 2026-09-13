/**
 * User API endpoints
 * Handles user profile management and password changes
 */

import { apiClient } from './client'
import {
  resolveWeChatOAuthStartStrict,
  prepareOAuthBindAccessTokenCookie,
  type WeChatOAuthPublicSettings,
} from './auth'
import type {
  User,
  ChangePasswordRequest,
  NotifyEmailEntry,
  UserAuthProvider,
  UserAffiliateDetail,
  AffiliateTransferResponse,
  PlatformQuotasResponse,
} from '@/types'

export type UserOnboardingStep = 'purchase' | 'redeem' | 'key' | 'config'
export type UserOnboardingEvent =
  | 'redeem_success'
  | 'key_created'
  | 'config_copied'
  | 'first_request'

export interface UserOnboardingState {
  version: number
  current_step: UserOnboardingStep
  completed_steps: UserOnboardingStep[]
  last_event?: string
  started_at?: string | null
  redeem_completed_at?: string | null
  key_created_at?: string | null
  config_copied_at?: string | null
  first_request_at?: string | null
  updated_at: string
  completed: boolean
  progress_percent: number
}

export interface UpdateUserOnboardingStateRequest {
  event: UserOnboardingEvent
}

/**
 * Get current user profile
 * @returns User profile data
 */
export async function getProfile(): Promise<User> {
  const { data } = await apiClient.get<User>('/user/profile')
  return data
}

/**
 * Update current user profile
 * @param profile - Profile data to update
 * @returns Updated user profile data
 */
export async function updateProfile(profile: {
  username?: string
  avatar_url?: string | null
  balance_notify_enabled?: boolean
  balance_notify_threshold?: number | null
  balance_notify_extra_emails?: NotifyEmailEntry[]
}): Promise<User> {
  const { data } = await apiClient.put<User>('/user', profile)
  return data
}

/**
 * Change current user password
 * @param passwords - Old and new password
 * @returns Success message
 */
export async function changePassword(
  oldPassword: string,
  newPassword: string
): Promise<{ message: string }> {
  const payload: ChangePasswordRequest = {
    old_password: oldPassword,
    new_password: newPassword
  }

  const { data } = await apiClient.put<{ message: string }>('/user/password', payload)
  return data
}

/**
 * Send verification code for adding a notify email
 * @param email - Email address to verify
 */
export async function sendNotifyEmailCode(email: string): Promise<void> {
  await apiClient.post('/user/notify-email/send-code', { email })
}

/**
 * Verify and add a notify email
 * @param email - Email address to add
 * @param code - Verification code
 */
export async function verifyNotifyEmail(email: string, code: string): Promise<void> {
  await apiClient.post('/user/notify-email/verify', { email, code })
}

/**
 * Remove a notify email
 * @param email - Email address to remove
 */
export async function removeNotifyEmail(email: string): Promise<void> {
  await apiClient.delete('/user/notify-email', { data: { email } })
}

/**
 * Toggle a notify email's disabled state
 * @param email - Email address (empty string for primary email placeholder)
 * @param disabled - Whether to disable the email
 */
export async function toggleNotifyEmail(email: string, disabled: boolean): Promise<User> {
  const { data } = await apiClient.put<User>('/user/notify-email/toggle', { email, disabled })
  return data
}

export async function sendEmailBindingCode(email: string): Promise<void> {
  await apiClient.post('/user/account-bindings/email/send-code', { email })
}

export async function bindEmailIdentity(payload: {
  email: string
  verify_code: string
  password: string
}): Promise<User> {
  const { data } = await apiClient.post<User>('/user/account-bindings/email', payload)
  return data
}

export async function unbindAuthIdentity(provider: BindableOAuthProvider): Promise<User> {
  const { data } = await apiClient.delete<User>(`/user/account-bindings/${provider}`)
  return data
}

export type BindableOAuthProvider = Exclude<UserAuthProvider, 'email'>

interface BuildOAuthBindingStartURLOptions {
  redirectTo?: string
  wechatOAuthSettings?: WeChatOAuthPublicSettings | null
}

export function resolveWeChatOAuthMode(): 'open' | 'mp' {
  if (typeof navigator === 'undefined') {
    return 'open'
  }
  return /MicroMessenger/i.test(navigator.userAgent) ? 'mp' : 'open'
}

function resolveWeChatOAuthBindingMode(
  settings?: WeChatOAuthPublicSettings | null
): 'open' | 'mp' | null {
  if (settings) {
    return resolveWeChatOAuthStartStrict(settings).mode
  }
  return resolveWeChatOAuthMode()
}

export function buildOAuthBindingStartURL(
  provider: BindableOAuthProvider,
  options: BuildOAuthBindingStartURLOptions = {}
): string | null {
  const redirectTo = options.redirectTo?.trim() || '/profile'
  const apiBase = (import.meta.env.VITE_API_BASE_URL as string | undefined) || '/api/v1'
  const normalized = apiBase.replace(/\/$/, '')
  const params = new URLSearchParams({
    redirect: redirectTo,
    intent: 'bind_current_user'
  })

  if (provider === 'wechat') {
    const mode = resolveWeChatOAuthBindingMode(options.wechatOAuthSettings)
    if (!mode) {
      return null
    }
    params.set('mode', mode)
  }

  return `${normalized}/auth/oauth/${provider}/bind/start?${params.toString()}`
}

export async function startOAuthBinding(
  provider: BindableOAuthProvider,
  options: BuildOAuthBindingStartURLOptions = {}
): Promise<void> {
  if (typeof window === 'undefined') {
    return
  }
  const startURL = buildOAuthBindingStartURL(provider, options)
  if (!startURL) {
    return
  }
  await prepareOAuthBindAccessTokenCookie()
  window.location.href = startURL
}

export async function getAffiliateDetail(): Promise<UserAffiliateDetail> {
  const { data } = await apiClient.get<UserAffiliateDetail>('/user/aff')
  return data
}

export async function transferAffiliateQuota(): Promise<AffiliateTransferResponse> {
  const { data } = await apiClient.post<AffiliateTransferResponse>('/user/aff/transfer')
  return data
}

<<<<<<< HEAD
/**
 * 获取当前用户的平台限额 + 用量。
 */
export async function getMyPlatformQuotas(): Promise<PlatformQuotasResponse> {
  const { data } = await apiClient.get<PlatformQuotasResponse>('/user/platform-quotas')
=======
export async function getOnboardingState(): Promise<UserOnboardingState> {
  const { data } = await apiClient.get<UserOnboardingState>('/user/onboarding-state')
  return data
}

export async function updateOnboardingState(
  payload: UpdateUserOnboardingStateRequest
): Promise<UserOnboardingState> {
  const { data } = await apiClient.patch<UserOnboardingState>('/user/onboarding-state', payload)
>>>>>>> 1157e4cfe6271a81c9dc96df84af6c6a3b22f831
  return data
}

export const userAPI = {
  getProfile,
  updateProfile,
  changePassword,
  sendNotifyEmailCode,
  verifyNotifyEmail,
  removeNotifyEmail,
  toggleNotifyEmail,
  sendEmailBindingCode,
  bindEmailIdentity,
  unbindAuthIdentity,
  buildOAuthBindingStartURL,
  startOAuthBinding,
  getAffiliateDetail,
  transferAffiliateQuota,
<<<<<<< HEAD
  getMyPlatformQuotas,
=======
  getOnboardingState,
  updateOnboardingState
>>>>>>> 1157e4cfe6271a81c9dc96df84af6c6a3b22f831
}

export default userAPI
