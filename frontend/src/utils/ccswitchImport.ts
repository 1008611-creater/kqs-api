import type { GroupPlatform } from '@/types'

// Terra is the default authority model for the current GPT-5.6 group.
// Sol remains available from the same imported provider when selected in Codex.
export const OPENAI_CC_SWITCH_CODEX_MODEL = 'gpt-5.6-terra'

export type CcSwitchClientType = 'claude' | 'gemini'

export interface CcSwitchImportConfig {
  app: string
  endpoint: string
  model?: string
  config?: string
}

export interface CcSwitchImportDeeplinkInput {
  baseUrl: string
  platform?: GroupPlatform | null
  clientType: CcSwitchClientType
  providerName: string
  apiKey: string
  usageScript: string
}

function trimTrailingSlashes(value: string): string {
  return value.trim().replace(/\/+$/, '')
}

function encodeBase64(value: string): string {
  const bytes = new TextEncoder().encode(value)
  let binary = ''
  for (let index = 0; index < bytes.length; index += 1) {
    binary += String.fromCharCode(bytes[index])
  }
  return btoa(binary)
}

function buildCodexConfig(providerName: string, endpoint: string, model: string): string {
  const quoteToml = (value: string) => JSON.stringify(value)
  const config = [
    'model_provider = "custom"',
    `model = ${quoteToml(model)}`,
    'model_reasoning_effort = "high"',
    'disable_response_storage = true',
    '[model_providers.custom]',
    `name = ${quoteToml(providerName)}`,
    `base_url = ${quoteToml(endpoint)}`,
    'wire_api = "responses"',
    'requires_openai_auth = true'
  ].join('\n')

  // CC Switch expects the inline config itself to be a Base64-encoded JSON
  // object whose `config` field contains Codex's TOML configuration.
  return encodeBase64(JSON.stringify({ config }))
}

export function resolveCcSwitchImportConfig(
  platform: GroupPlatform | undefined | null,
  clientType: CcSwitchClientType,
  baseUrl: string
): CcSwitchImportConfig {
  const normalizedBaseUrl = trimTrailingSlashes(baseUrl)
  switch (platform || 'anthropic') {
    case 'antigravity':
      return {
        app: clientType === 'gemini' ? 'gemini' : 'claude',
        endpoint: `${normalizedBaseUrl}/antigravity`
      }
    case 'openai':
      return {
        app: 'codex',
        endpoint: normalizedBaseUrl,
        model: OPENAI_CC_SWITCH_CODEX_MODEL
      }
    case 'gemini':
      return {
        app: 'gemini',
        endpoint: normalizedBaseUrl
      }
    default:
      return {
        app: 'claude',
        endpoint: normalizedBaseUrl
      }
  }
}

export function buildCcSwitchImportDeeplink(input: CcSwitchImportDeeplinkInput): string {
  const resolvedConfig = resolveCcSwitchImportConfig(input.platform, input.clientType, input.baseUrl)
  const config = resolvedConfig.app === 'codex' && resolvedConfig.model
    ? {
        ...resolvedConfig,
        config: buildCodexConfig(input.providerName, resolvedConfig.endpoint, resolvedConfig.model)
      }
    : resolvedConfig
  const baseUrl = trimTrailingSlashes(input.baseUrl)
  const entries: [string, string][] = [
    ['resource', 'provider'],
    ['app', config.app],
    ['name', input.providerName],
    ['homepage', baseUrl],
    ['endpoint', config.endpoint],
    ['apiKey', input.apiKey],
    ['usageEnabled', 'true'],
    ['usageScript', encodeBase64(input.usageScript)],
    ['usageAutoInterval', '30']
  ]

  if (config.model) {
    entries.splice(2, 0, ['model', config.model])
  }

  if (config.config) {
    entries.splice(7, 0, ['configFormat', 'json'], ['config', config.config])
  }

  return `ccswitch://v1/import?${new URLSearchParams(entries).toString()}`
}
