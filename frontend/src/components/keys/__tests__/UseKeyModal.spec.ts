import { describe, expect, it, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import { nextTick } from 'vue'

vi.mock('vue-i18n', () => ({
  useI18n: () => ({
    t: (key: string) => key
  })
}))

vi.mock('@/composables/useClipboard', () => ({
  useClipboard: () => ({
    copyToClipboard: vi.fn().mockResolvedValue(true)
  })
}))

import UseKeyModal from '../UseKeyModal.vue'

describe('UseKeyModal', () => {
  it('renders a minimal OpenAI Codex config by default', () => {
    const wrapper = mount(UseKeyModal, {
      props: {
        show: true,
        apiKey: 'sk-test',
        baseUrl: 'https://example.com',
        platform: 'openai',
        groupName: '普通 GPT 分组',
        supportedModelScopes: []
      },
      global: {
        stubs: {
          BaseDialog: {
            template: '<div><slot /><slot name="footer" /></div>'
          },
          Icon: {
            template: '<span />'
          }
        }
      }
    })

    const text = wrapper.text()
    expect(text).toContain('keys.useKeyModal.manual.title')
    expect(text).not.toContain('keys.useKeyModal.codexApp.title')
    expect(text).toContain('~/.codex')
    expect(text).toContain('~/.codex/config.toml')
    expect(text).toContain('~/.codex/auth.json')
    expect(text).not.toContain('install-codex-openai.sh')

    const codeBlocks = wrapper.findAll('pre code').map((code) => code.text()).join('\n')
    expect(codeBlocks).toContain('[model_providers.KQS_API]')
    expect(codeBlocks).toContain('model = "gpt-5.5"')
    expect(codeBlocks).toContain('"OPENAI_API_KEY": "sk-test"')
    expect(codeBlocks).not.toContain('[profiles.gpt-5-5]')
    expect(codeBlocks).not.toContain('[profiles.qwen-plus]')
    expect(codeBlocks).not.toContain('model_context_window')
    expect(codeBlocks).not.toContain('model_auto_compact_token_limit')
    expect(codeBlocks).not.toContain('model_catalog_json')
  })

  it('keeps Codex config minimal even for previous GPT plus domestic bundle groups', () => {
    const wrapper = mount(UseKeyModal, {
      props: {
        show: true,
        apiKey: 'sk-test',
        baseUrl: 'https://example.com',
        platform: 'openai',
        groupName: 'codex直连gpt5.5&国产模型合集',
        supportedModelScopes: ['codex_all_models']
      },
      global: {
        stubs: {
          BaseDialog: {
            template: '<div><slot /><slot name="footer" /></div>'
          },
          Icon: {
            template: '<span />'
          }
        }
      }
    })

    expect(wrapper.text()).not.toContain('keys.useKeyModal.codexApp.title')
    const codeBlocks = wrapper.findAll('pre code').map((code) => code.text()).join('\n')
    expect(codeBlocks).toContain('[model_providers.KQS_API]')
    expect(codeBlocks).toContain('model = "gpt-5.5"')
    expect(codeBlocks).not.toContain('[profiles.gpt-5-5]')
    expect(codeBlocks).not.toContain('[profiles.codex-mini-latest]')
    expect(codeBlocks).not.toContain('[profiles.qwen-plus]')
    expect(codeBlocks).not.toContain('model = "qwen-plus"')
    expect(codeBlocks).not.toContain('[profiles.deepseek-reasoner]')
    expect(codeBlocks).not.toContain('model = "deepseek-reasoner"')
    expect(codeBlocks).not.toContain('[profiles.kimi-latest]')
    expect(codeBlocks).not.toContain('[profiles.glm-4-6]')
    expect(codeBlocks).not.toContain('[profiles.doubao-thinking]')
    expect(codeBlocks).not.toContain('[profiles.hunyuan-code]')
    expect(codeBlocks).not.toContain('codex -p qwen-plus')
    expect(codeBlocks).not.toContain('model-catalog-kqs-api.json')
    expect(codeBlocks).not.toContain('install-codex-openai.sh')
    expect(codeBlocks).not.toContain('"slug": "qwen-plus"')
    expect(codeBlocks).not.toContain('claude-sonnet')
    expect(codeBlocks).not.toContain('gpt-image')
    expect(codeBlocks).not.toContain('cogview')
  })

  it('renders GPT-5.4 mini entry in OpenCode config', async () => {
    const wrapper = mount(UseKeyModal, {
      props: {
        show: true,
        apiKey: 'sk-test',
        baseUrl: 'https://example.com/v1',
        platform: 'openai'
      },
      global: {
        stubs: {
          BaseDialog: {
            template: '<div><slot /><slot name="footer" /></div>'
          },
          Icon: {
            template: '<span />'
          }
        }
      }
    })

    const opencodeTab = wrapper.findAll('button').find((button) =>
      button.text().includes('keys.useKeyModal.cliTabs.opencode')
    )

    expect(opencodeTab).toBeDefined()
    await opencodeTab!.trigger('click')
    await nextTick()

    const codeBlock = wrapper.find('pre code')
    expect(codeBlock.exists()).toBe(true)
    expect(codeBlock.text()).toContain('"name": "GPT-5.4 Mini"')
    expect(codeBlock.text()).not.toContain('"name": "GPT-5.4 Nano"')
  })
})
