import { describe, expect, it } from 'vitest'
import {
  CODEX_ALL_MODELS_SCOPE,
  CODEX_DOMESTIC_MODEL_PROFILES,
  buildCodexModelCatalogJson,
  buildCodexConfigToml,
  buildCodexDomesticInstallScript,
  buildCodexDomesticProfilesToml,
  isCodexCombinedModelGroup,
  isCodexDomesticTextModel
} from '../codexAppProfiles'

describe('codexAppProfiles', () => {
  it('exports domestic text model profiles without GPT, Claude, image, vision, or video models', () => {
    expect(CODEX_DOMESTIC_MODEL_PROFILES.length).toBeGreaterThanOrEqual(24)

    for (const profile of CODEX_DOMESTIC_MODEL_PROFILES) {
      expect(isCodexDomesticTextModel(profile.model)).toBe(true)
      expect(profile.model).not.toMatch(/gpt|claude|image|vision|video|cogview|cogvideo/i)
    }
  })

  it('builds minimal Codex TOML for normal OpenAI groups', () => {
    const toml = buildCodexConfigToml('https://api.example.com')

    expect(toml).toContain('[model_providers.KQS_API]')
    expect(toml).toContain('base_url = "https://api.example.com"')
    expect(toml).toContain('model = "gpt-5.5"')
    expect(toml).not.toContain('review_model')
    expect(toml).not.toContain('model_reasoning_effort')
    expect(toml).not.toContain('model_context_window')
    expect(toml).not.toContain('model_auto_compact_token_limit')
    expect(toml).not.toContain('[profiles.gpt-5-5]')
    expect(toml).not.toContain('[profiles.qwen-plus]')
    expect(toml).not.toContain('[profiles.deepseek-reasoner]')
  })

  it('builds GPT plus domestic Codex TOML profiles for bundle groups', () => {
    const toml = buildCodexConfigToml('https://api.example.com', {
      modelSet: 'gpt-domestic',
      modelCatalogPath: '~/.codex/model-catalog-kqs-api.json'
    })

    expect(toml).toContain('[model_providers.KQS_API]')
    expect(toml).toContain('base_url = "https://api.example.com"')
    expect(toml).toContain('model_catalog_json = "~/.codex/model-catalog-kqs-api.json"')
    expect(toml).toContain('# GPT')
    expect(toml).toContain('[profiles.gpt-5-5]')
    expect(toml).toContain('# Qwen')
    expect(toml).toContain('[profiles.qwen-plus]')
    expect(toml).toContain('# DeepSeek')
    expect(toml).toContain('[profiles.deepseek-reasoner]')
    expect(toml).toContain('# Kimi')
    expect(toml).toContain('[profiles.kimi-latest]')
    expect(toml).toContain('model_provider = "KQS_API"')
    expect(toml).not.toContain('cogview')
  })

  it('builds an installer script that backs up existing Codex config', () => {
    const script = buildCodexDomesticInstallScript('https://api.example.com', 'sk-test', {
      shell: 'powershell'
    })

    expect(script).toContain('Copy-Item $configPath "$configPath.bak-$timestamp"')
    expect(script).toContain('Copy-Item $authPath "$authPath.bak-$timestamp"')
    expect(script).toContain('model-catalog-kqs-api.json')
    expect(script).toContain('Remove-KqsManagedBlocks')
    expect(script).toContain('[profiles.qwen-plus]')
    expect(script).toContain('"OPENAI_API_KEY": "sk-test"')
  })

  it('can generate just the domestic profiles block', () => {
    const toml = buildCodexDomesticProfilesToml('LOCAL_PROVIDER')

    expect(toml).toContain('model_provider = "LOCAL_PROVIDER"')
    expect(toml).toContain('[profiles.hunyuan-code]')
    expect(toml).toContain('[profiles.yi-large]')
  })

  it('builds a Codex App model catalog containing GPT and domestic text models', () => {
    const catalog = JSON.parse(buildCodexModelCatalogJson('gpt-domestic'))
    const slugs = catalog.models.map((model: { slug: string }) => model.slug)

    expect(slugs).toContain('gpt-5.5')
    expect(slugs).toContain('gpt-5.4')
    expect(slugs).toContain('qwen-plus')
    expect(slugs).toContain('deepseek-reasoner')
    expect(slugs).toContain('kimi-latest')
    expect(slugs).not.toContain('claude-sonnet-4-6')
    expect(slugs.some((slug: string) => /image|vision|video|cogview/i.test(slug))).toBe(false)
  })

  it('detects only the combined OpenAI Codex model group', () => {
    expect(isCodexCombinedModelGroup({
      platform: 'openai',
      supportedModelScopes: [CODEX_ALL_MODELS_SCOPE]
    })).toBe(true)
    expect(isCodexCombinedModelGroup({
      platform: 'openai',
      name: 'codex直连gpt5.5&国产模型合集'
    })).toBe(true)
    expect(isCodexCombinedModelGroup({
      platform: 'openai',
      name: '普通 GPT 分组'
    })).toBe(false)
    expect(isCodexCombinedModelGroup({
      platform: 'anthropic',
      name: 'GPT + 国产模型合集',
      supportedModelScopes: [CODEX_ALL_MODELS_SCOPE]
    })).toBe(false)
  })
})
