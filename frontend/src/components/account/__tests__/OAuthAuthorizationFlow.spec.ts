import { describe, expect, it, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { defineComponent } from 'vue'
import OAuthAuthorizationFlow from '../OAuthAuthorizationFlow.vue'

vi.mock('vue-i18n', async () => {
  const actual = await vi.importActual<typeof import('vue-i18n')>('vue-i18n')
  return {
    ...actual,
    useI18n: () => ({
      t: (key: string) => key
    })
  }
})

vi.mock('@/composables/useClipboard', () => ({
  useClipboard: () => ({
    copied: { value: false },
    copyToClipboard: vi.fn()
  })
}))

const IconStub = defineComponent({
  name: 'Icon',
  template: '<span data-testid="icon" />'
})

const mountFlow = () =>
  mount(OAuthAuthorizationFlow, {
    props: {
      addMethod: 'oauth',
      platform: 'openai',
      showCookieOption: false,
      showRefreshTokenOption: true,
      showMobileRefreshTokenOption: true,
      showCodexSessionImportOption: true,
      showHelp: false
    },
    global: {
      stubs: {
        Icon: IconStub
      }
    }
  })

describe('OAuthAuthorizationFlow', () => {
  it('defaults OpenAI imports to the Codex JSON/AT/RT batch path', async () => {
    const wrapper = mountFlow()

    expect(wrapper.vm.inputMethod).toBe('codex_session')
    expect(wrapper.text()).toContain('admin.accounts.oauth.openai.codexSessionInputLabel')
    expect(wrapper.text()).toContain('admin.accounts.oauth.openai.codexSessionAcceptRefreshToken')

    await wrapper.get('textarea').setValue('[{"refresh_token":"rt_test_1"},{"accessToken":"access_test_2"}]')
    await wrapper.get('button.btn-primary').trigger('click')

    expect(wrapper.emitted('import-codex-session')?.[0]).toEqual([
      '[{"refresh_token":"rt_test_1"},{"accessToken":"access_test_2"}]'
    ])
  })

  it('restores Codex import as the OpenAI default after reset', () => {
    const wrapper = mountFlow()

    wrapper.vm.inputMethod = 'manual'
    wrapper.vm.reset()

    expect(wrapper.vm.inputMethod).toBe('codex_session')
  })
})
