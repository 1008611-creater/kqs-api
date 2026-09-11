import { describe, expect, it } from 'vitest'
import {
  OPENAI_CC_SWITCH_CODEX_MODEL,
  buildCcSwitchImportDeeplink
} from '@/utils/ccswitchImport'
import type { GroupPlatform } from '@/types'

function paramsFromDeeplink(deeplink: string): URLSearchParams {
  const query = deeplink.split('?')[1] || ''
  return new URLSearchParams(query)
}

describe('ccswitchImport utils', () => {
  const baseInput = {
    baseUrl: 'https://api.example.com',
    providerName: 'Sub2API',
    apiKey: 'sk-test',
    usageScript: 'return true'
  }

  it('adds the Codex model parameter for OpenAI imports', () => {
    const params = paramsFromDeeplink(
      buildCcSwitchImportDeeplink({
        ...baseInput,
        platform: 'openai',
        clientType: 'claude'
      })
    )

    expect(params.get('resource')).toBe('provider')
    expect(params.get('app')).toBe('codex')
    expect(params.get('endpoint')).toBe(baseInput.baseUrl)
    expect(params.get('model')).toBe(OPENAI_CC_SWITCH_CODEX_MODEL)
    expect(params.get('model')).toBe('gpt-5.6-terra')
    expect(params.get('configFormat')).toBe('json')
    const codexConfig = JSON.parse(new TextDecoder().decode(Uint8Array.from(atob(params.get('config') || ''), char => char.charCodeAt(0))))
    expect(codexConfig.config).toContain('wire_api = "responses"')
    expect(codexConfig.config).toContain(`base_url = "${baseInput.baseUrl}"`)
    expect(atob(params.get('usageScript') || '')).toBe(baseInput.usageScript)
  })

  it('encodes non-ASCII provider names without breaking the import link', () => {
    const deeplink = buildCcSwitchImportDeeplink({
      ...baseInput,
      providerName: '矿泉水 API',
      platform: 'openai',
      clientType: 'claude'
    })
    const params = paramsFromDeeplink(deeplink)
    const codexConfig = JSON.parse(new TextDecoder().decode(Uint8Array.from(atob(params.get('config') || ''), char => char.charCodeAt(0))))

    expect(params.get('name')).toBe('矿泉水 API')
    expect(codexConfig.config).toContain('矿泉水 API')
  })

  it.each([
    { platform: 'anthropic' as GroupPlatform, clientType: 'claude' as const, app: 'claude' },
    { platform: 'gemini' as GroupPlatform, clientType: 'gemini' as const, app: 'gemini' }
  ])('does not add a model parameter for $platform imports', ({ platform, clientType, app }) => {
    const params = paramsFromDeeplink(
      buildCcSwitchImportDeeplink({
        ...baseInput,
        platform,
        clientType
      })
    )

    expect(params.get('app')).toBe(app)
    expect(params.get('endpoint')).toBe(baseInput.baseUrl)
    expect(params.has('model')).toBe(false)
    expect(params.has('config')).toBe(false)
  })

  it('keeps Antigravity imports on the selected client endpoint without a model parameter', () => {
    const params = paramsFromDeeplink(
      buildCcSwitchImportDeeplink({
        ...baseInput,
        platform: 'antigravity',
        clientType: 'gemini'
      })
    )

    expect(params.get('app')).toBe('gemini')
    expect(params.get('endpoint')).toBe(`${baseInput.baseUrl}/antigravity`)
    expect(params.has('model')).toBe(false)
  })
})
