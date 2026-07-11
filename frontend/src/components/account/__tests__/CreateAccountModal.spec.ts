import { describe, expect, it, vi } from 'vitest'
import { mount } from '@vue/test-utils'

import CreateAccountModal from '../CreateAccountModal.vue'

const { getWebSearchEmulationConfig, listTLSProfiles } = vi.hoisted(() => ({
  getWebSearchEmulationConfig: vi.fn(),
  listTLSProfiles: vi.fn()
}))

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string) => key
    })
  }
})

vi.mock('@/api/admin', () => ({
  adminAPI: {
    settings: {
      getWebSearchEmulationConfig
    },
    tlsFingerprintProfiles: {
      list: listTLSProfiles
    },
    accounts: {
      checkMixedChannelRisk: vi.fn(),
      create: vi.fn(),
      importCodexSession: vi.fn(),
      exchangeCode: vi.fn()
    }
  }
}))

vi.mock('@/stores/app', () => ({
  useAppStore: () => ({
    showError: vi.fn(),
    showSuccess: vi.fn(),
    showWarning: vi.fn(),
    showInfo: vi.fn()
  })
}))

vi.mock('@/stores/auth', () => ({
  useAuthStore: () => ({
    token: 'test-token'
  })
}))

vi.mock('@/composables/useQuotaNotifyState', () => ({
  useQuotaNotifyState: () => ({
    globalEnabled: { value: false },
    state: { value: null },
    loadGlobalState: vi.fn(),
    writeToExtra: vi.fn()
  })
}))

const oauthState = () => ({
  authUrl: { value: '' },
  sessionId: { value: '' },
  loading: { value: false },
  error: { value: '' },
  exchangeCode: vi.fn(),
  validateRefreshToken: vi.fn(),
  reset: vi.fn()
})

vi.mock('@/composables/useAccountOAuth', async () => {
  const actual = await vi.importActual<typeof import('@/composables/useAccountOAuth')>('@/composables/useAccountOAuth')
  return {
    ...actual,
    useAccountOAuth: () => oauthState()
  }
})

vi.mock('@/composables/useOpenAIOAuth', () => ({
  useOpenAIOAuth: () => oauthState()
}))

vi.mock('@/composables/useGeminiOAuth', () => ({
  useGeminiOAuth: () => oauthState()
}))

vi.mock('@/composables/useAntigravityOAuth', () => ({
  useAntigravityOAuth: () => oauthState()
}))

const mountModal = () =>
  mount(CreateAccountModal, {
    props: {
      show: true,
      proxies: [],
      groups: []
    },
    global: {
      stubs: {
        BaseDialog: {
          template: '<div data-test="base-dialog"><slot /><slot name="footer" /></div>'
        },
        ConfirmDialog: true,
        Icon: true,
        Select: true,
        ProxySelector: true,
        ProxyAdBanner: true,
        GroupSelector: true,
        ModelWhitelistSelector: true,
        QuotaLimitCard: true,
        OAuthAuthorizationFlow: true
      }
    }
  })

describe('CreateAccountModal', () => {
  it('does not require account name for OpenAI OAuth paste import flow', () => {
    getWebSearchEmulationConfig.mockResolvedValue({ enabled: false, providers: [] })
    listTLSProfiles.mockResolvedValue({ items: [] })

    const wrapper = mountModal()
    const nameInput = wrapper.get('[data-tour="account-form-name"]')

    expect(nameInput.attributes('required')).toBeUndefined()
    expect(wrapper.text()).toContain('common.optional')
    expect(wrapper.text()).toContain('admin.accounts.createGuide.openaiOauth.nameOptionalHint')
  })

  it('low-cost supply preset fills priority, weight, and max body limit', async () => {
    getWebSearchEmulationConfig.mockResolvedValue({ enabled: false, providers: [] })
    listTLSProfiles.mockResolvedValue({ items: [] })

    const wrapper = mountModal()

    await wrapper.get('[data-testid="account-supply-preset-lowCostSmall"]').trigger('click')

    expect((wrapper.get('#create-account-priority').element as HTMLInputElement).value).toBe('1')
    expect((wrapper.get('#create-account-load-factor').element as HTMLInputElement).value).toBe('1')
    expect((wrapper.get('#create-account-max-body').element as HTMLInputElement).value).toBe('10')
  })
})
