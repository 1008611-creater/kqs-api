export type CodexReasoningEffort = 'none' | 'low' | 'medium' | 'high' | 'xhigh'
export type CodexModelSet = 'gpt' | 'gpt-domestic'

export interface CodexDomesticProfile {
  brand: string
  profile: string
  label: string
  model: string
  effort?: CodexReasoningEffort
  contextWindow?: number
  description?: string
}

export const CODEX_PROVIDER_NAME = 'KQS_API'
export const CODEX_DEFAULT_MODEL = 'gpt-5.5'
export const CODEX_ALL_MODELS_SCOPE = 'codex_all_models'
export const CODEX_MODEL_CATALOG_FILENAME = 'model-catalog-kqs-api.json'
export const CODEX_TOP_LEVEL_MARKER_START = '# >>> KQS_API_CODEX_TOP_LEVEL'
export const CODEX_TOP_LEVEL_MARKER_END = '# <<< KQS_API_CODEX_TOP_LEVEL'
export const CODEX_MODELS_MARKER_START = '# >>> KQS_API_CODEX_MODELS'
export const CODEX_MODELS_MARKER_END = '# <<< KQS_API_CODEX_MODELS'
const LEGACY_CODEX_PROVIDER_NAME = String.fromCharCode(67, 65, 85) + '_API'
const LEGACY_CODEX_TOP_LEVEL_MARKER_START = '# >>> ' + LEGACY_CODEX_PROVIDER_NAME + '_CODEX_TOP_LEVEL'
const LEGACY_CODEX_TOP_LEVEL_MARKER_END = '# <<< ' + LEGACY_CODEX_PROVIDER_NAME + '_CODEX_TOP_LEVEL'
const LEGACY_CODEX_MODELS_MARKER_START = '# >>> ' + LEGACY_CODEX_PROVIDER_NAME + '_CODEX_MODELS'
const LEGACY_CODEX_MODELS_MARKER_END = '# <<< ' + LEGACY_CODEX_PROVIDER_NAME + '_CODEX_MODELS'

export const CODEX_GPT_MODEL_PROFILES: CodexDomesticProfile[] = [
  { brand: 'GPT', profile: 'gpt-5-5', label: 'GPT-5.5', model: 'gpt-5.5', effort: 'xhigh', contextWindow: 1050000 },
  { brand: 'GPT', profile: 'gpt-5-4', label: 'GPT-5.4', model: 'gpt-5.4', effort: 'xhigh', contextWindow: 1050000 },
  { brand: 'GPT', profile: 'gpt-5-3-codex', label: 'GPT-5.3 Codex', model: 'gpt-5.3-codex', effort: 'high', contextWindow: 400000 },
  { brand: 'GPT', profile: 'gpt-5-3-codex-spark', label: 'GPT-5.3 Codex Spark', model: 'gpt-5.3-codex-spark', effort: 'high', contextWindow: 128000 },
  { brand: 'GPT', profile: 'gpt-5-2', label: 'GPT-5.2', model: 'gpt-5.2', effort: 'high', contextWindow: 400000 },
  { brand: 'GPT', profile: 'codex-mini-latest', label: 'Codex Mini', model: 'codex-mini-latest', effort: 'medium', contextWindow: 200000 }
]

export const CODEX_DOMESTIC_MODEL_PROFILES: CodexDomesticProfile[] = [
  { brand: 'Qwen', profile: 'qwen-plus', label: 'Qwen Plus', model: 'qwen-plus', effort: 'high', contextWindow: 128000, description: 'Tongyi Qianwen Plus' },
  { brand: 'Qwen', profile: 'qwen-max', label: 'Qwen Max', model: 'qwen-max', effort: 'high', contextWindow: 128000, description: 'Tongyi Qianwen Max' },
  { brand: 'Qwen', profile: 'qwen3-235b', label: 'Qwen3 235B', model: 'qwen3-235b-a22b', effort: 'high', contextWindow: 128000 },
  { brand: 'Qwen', profile: 'qwen-coder', label: 'Qwen Coder', model: 'qwen2.5-coder-32b-instruct', effort: 'high', contextWindow: 128000 },
  { brand: 'Qwen', profile: 'qwen-qwq', label: 'QwQ', model: 'qwq-32b', effort: 'xhigh', contextWindow: 128000 },
  { brand: 'DeepSeek', profile: 'deepseek-chat', label: 'DeepSeek Chat', model: 'deepseek-chat', effort: 'high', contextWindow: 64000 },
  { brand: 'DeepSeek', profile: 'deepseek-reasoner', label: 'DeepSeek Reasoner', model: 'deepseek-reasoner', effort: 'xhigh', contextWindow: 64000 },
  { brand: 'DeepSeek', profile: 'deepseek-coder', label: 'DeepSeek Coder', model: 'deepseek-coder', effort: 'high', contextWindow: 64000 },
  { brand: 'DeepSeek', profile: 'deepseek-v3', label: 'DeepSeek V3', model: 'deepseek-v3', effort: 'high', contextWindow: 64000 },
  { brand: 'DeepSeek', profile: 'deepseek-r1', label: 'DeepSeek R1', model: 'deepseek-r1', effort: 'xhigh', contextWindow: 64000 },
  { brand: 'Kimi', profile: 'kimi-latest', label: 'Kimi Latest', model: 'kimi-latest', effort: 'high', contextWindow: 128000 },
  { brand: 'Kimi', profile: 'kimi-128k', label: 'Kimi 128K', model: 'moonshot-v1-128k', effort: 'high', contextWindow: 128000 },
  { brand: 'Kimi', profile: 'kimi-32k', label: 'Kimi 32K', model: 'moonshot-v1-32k', effort: 'medium', contextWindow: 32000 },
  { brand: 'GLM', profile: 'glm-4-6', label: 'GLM 4.6', model: 'glm-4.6', effort: 'high', contextWindow: 128000 },
  { brand: 'GLM', profile: 'glm-4-5', label: 'GLM 4.5', model: 'glm-4.5', effort: 'high', contextWindow: 128000 },
  { brand: 'GLM', profile: 'glm-4-plus', label: 'GLM 4 Plus', model: 'glm-4-plus', effort: 'high', contextWindow: 128000 },
  { brand: 'GLM', profile: 'glm-4-air', label: 'GLM 4 Air', model: 'glm-4-air', effort: 'medium', contextWindow: 128000 },
  { brand: 'Doubao', profile: 'doubao-pro-256k', label: 'Doubao Pro 256K', model: 'doubao-1.5-pro-256k', effort: 'high', contextWindow: 256000 },
  { brand: 'Doubao', profile: 'doubao-thinking', label: 'Doubao Thinking', model: 'doubao-1.5-thinking-pro', effort: 'xhigh', contextWindow: 128000 },
  { brand: 'Doubao', profile: 'doubao-pro-128k', label: 'Doubao Pro 128K', model: 'doubao-pro-128k', effort: 'high', contextWindow: 128000 },
  { brand: 'Doubao', profile: 'doubao-lite', label: 'Doubao Lite', model: 'doubao-lite-32k', effort: 'medium', contextWindow: 32000 },
  { brand: 'MiniMax', profile: 'minimax-pro', label: 'MiniMax Pro', model: 'abab6.5s-chat-pro', effort: 'high', contextWindow: 128000 },
  { brand: 'MiniMax', profile: 'minimax-chat', label: 'MiniMax Chat', model: 'abab6.5-chat', effort: 'medium', contextWindow: 32000 },
  { brand: 'Baidu', profile: 'ernie-4-latest', label: 'ERNIE 4 Latest', model: 'ernie-4.0-8k-latest', effort: 'high', contextWindow: 8000 },
  { brand: 'Baidu', profile: 'ernie-4-turbo', label: 'ERNIE 4 Turbo', model: 'ernie-4.0-turbo-8k', effort: 'high', contextWindow: 8000 },
  { brand: 'Baidu', profile: 'ernie-35-128k', label: 'ERNIE 3.5 128K', model: 'ernie-3.5-128k', effort: 'medium', contextWindow: 128000 },
  { brand: 'Spark', profile: 'spark-ultra', label: 'Spark Ultra', model: 'spark-ultra', effort: 'high', contextWindow: 128000 },
  { brand: 'Spark', profile: 'spark-max', label: 'Spark Max', model: 'spark-max', effort: 'high', contextWindow: 128000 },
  { brand: 'Spark', profile: 'spark-pro', label: 'Spark Pro', model: 'spark-pro', effort: 'medium', contextWindow: 32000 },
  { brand: 'Hunyuan', profile: 'hunyuan-turbo', label: 'Hunyuan Turbo', model: 'hunyuan-turbo', effort: 'high', contextWindow: 128000 },
  { brand: 'Hunyuan', profile: 'hunyuan-large', label: 'Hunyuan Large', model: 'hunyuan-large', effort: 'high', contextWindow: 128000 },
  { brand: 'Hunyuan', profile: 'hunyuan-pro', label: 'Hunyuan Pro', model: 'hunyuan-pro', effort: 'high', contextWindow: 128000 },
  { brand: 'Hunyuan', profile: 'hunyuan-256k', label: 'Hunyuan 256K', model: 'hunyuan-standard-256k', effort: 'medium', contextWindow: 256000 },
  { brand: 'Hunyuan', profile: 'hunyuan-code', label: 'Hunyuan Code', model: 'hunyuan-code', effort: 'high', contextWindow: 128000 },
  { brand: 'Yi', profile: 'yi-large', label: 'Yi Large', model: 'yi-large', effort: 'high', contextWindow: 32000 },
  { brand: 'Yi', profile: 'yi-200k', label: 'Yi 200K', model: 'yi-medium-200k', effort: 'medium', contextWindow: 200000 },
  { brand: 'Yi', profile: 'yi-34b', label: 'Yi 34B', model: 'yi-1.5-34b-chat', effort: 'medium', contextWindow: 32000 }
]

const blockedDomesticModelPattern = /(gpt|claude|image|vision|video|cogview|cogvideo|audio|tts)/i

export function isCodexDomesticTextModel(model: string): boolean {
  return !blockedDomesticModelPattern.test(model)
}

const activeDomesticProfiles = () =>
  CODEX_DOMESTIC_MODEL_PROFILES.filter((profile) => isCodexDomesticTextModel(profile.model))

const profileGroupsFromProfiles = (profiles: CodexDomesticProfile[]) => Array.from(
  profiles.reduce((groups, item) => {
    const entries = groups.get(item.brand) || []
    entries.push(item)
    groups.set(item.brand, entries)
    return groups
  }, new Map<string, CodexDomesticProfile[]>())
).map(([brand, profiles]) => ({ brand, profiles }))

export const CODEX_DOMESTIC_PROFILE_GROUPS = profileGroupsFromProfiles(activeDomesticProfiles())
export const CODEX_ALL_MODEL_PROFILES = [...CODEX_GPT_MODEL_PROFILES, ...activeDomesticProfiles()]
export const CODEX_ALL_MODEL_PROFILE_GROUPS = profileGroupsFromProfiles(CODEX_ALL_MODEL_PROFILES)

export interface CodexCombinedModelGroupInput {
  platform?: string | null
  name?: string | null
  description?: string | null
  supportedModelScopes?: string[] | null
}

export function isCodexCombinedModelGroup(input: CodexCombinedModelGroupInput): boolean {
  if (input.platform !== 'openai') return false
  const scopes = input.supportedModelScopes || []
  if (scopes.includes(CODEX_ALL_MODELS_SCOPE)) return true

  const text = `${input.name || ''} ${input.description || ''}`.toLowerCase()
  return (
    text.includes('gpt+国产') ||
    text.includes('gpt 加国产') ||
    text.includes('国产模型合集') ||
    text.includes('模型合集') ||
    (text.includes('gpt') && text.includes('国产')) ||
    (text.includes('codex') && text.includes('合集'))
  )
}

function tomlString(value: string): string {
  return `"${value.replace(/\\/g, '\\\\').replace(/"/g, '\\"')}"`
}

function shellSingleQuote(value: string): string {
  return `'${value.replace(/'/g, `'\\''`)}'`
}

function powershellHereString(value: string): string {
  return `@'\n${value.replace(/\r\n/g, '\n')}\n'@`
}

function regexEscape(value: string): string {
  return value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
}

function profilesForModelSet(modelSet: CodexModelSet): CodexDomesticProfile[] {
  return modelSet === 'gpt-domestic' ? CODEX_ALL_MODEL_PROFILES : CODEX_GPT_MODEL_PROFILES
}

function profileGroupsForModelSet(modelSet: CodexModelSet) {
  return modelSet === 'gpt-domestic'
    ? CODEX_ALL_MODEL_PROFILE_GROUPS
    : profileGroupsFromProfiles(CODEX_GPT_MODEL_PROFILES)
}

export function buildCodexProfilesToml(providerName = CODEX_PROVIDER_NAME, modelSet: CodexModelSet = 'gpt'): string {
  return profileGroupsForModelSet(modelSet).map(({ brand, profiles }) => {
    const blocks = profiles.map((profile) => [
      `[profiles.${profile.profile}]`,
      `model_provider = ${tomlString(providerName)}`,
      `model = ${tomlString(profile.model)}`,
      `model_reasoning_effort = ${tomlString(profile.effort || 'high')}`
    ].join('\n'))

    return [`# ${brand}`, ...blocks].join('\n\n')
  }).join('\n\n')
}

export function buildCodexDomesticProfilesToml(providerName = CODEX_PROVIDER_NAME): string {
  return buildCodexProfilesToml(providerName, 'gpt-domestic')
}

export interface BuildCodexConfigOptions {
  supportsWebsockets?: boolean
  modelSet?: CodexModelSet
  modelCatalogPath?: string
  includeProfiles?: boolean
}

export function buildCodexTopLevelToml(options: BuildCodexConfigOptions = {}): string {
  const catalogLine = options.modelCatalogPath
    ? `\nmodel_catalog_json = ${tomlString(options.modelCatalogPath)}`
    : ''

  return `${CODEX_TOP_LEVEL_MARKER_START}
model_provider = ${tomlString(CODEX_PROVIDER_NAME)}
model = ${tomlString(CODEX_DEFAULT_MODEL)}${catalogLine}
${CODEX_TOP_LEVEL_MARKER_END}`
}

export function buildCodexProviderAndProfilesToml(baseUrl: string, options: BuildCodexConfigOptions = {}): string {
  const websocketLine = options.supportsWebsockets ? '\nsupports_websockets = true' : ''
  const featuresBlock = options.supportsWebsockets
    ? `\n\n[features]\nresponses_websockets_v2 = true`
    : ''
  const modelSet = options.modelSet || 'gpt'
  const includeProfiles = options.includeProfiles ?? modelSet === 'gpt-domestic'
  const profilesBlock = includeProfiles
    ? `\n\n# Optional Codex CLI profiles:
${profilesForModelSet(modelSet)
  .slice(0, 3)
  .map((profile) => `# codex -p ${profile.profile}`)
  .join('\n')}

${buildCodexProfilesToml(CODEX_PROVIDER_NAME, modelSet)}`
    : ''

  return `${CODEX_MODELS_MARKER_START}
[model_providers.${CODEX_PROVIDER_NAME}]
name = ${tomlString(CODEX_PROVIDER_NAME)}
base_url = ${tomlString(baseUrl)}
wire_api = "responses"${websocketLine}
requires_openai_auth = true${featuresBlock}${profilesBlock}
${CODEX_MODELS_MARKER_END}`
}

export function buildCodexConfigToml(baseUrl: string, options: BuildCodexConfigOptions = {}): string {
  return [
    buildCodexTopLevelToml(options),
    buildCodexProviderAndProfilesToml(baseUrl, options)
  ].join('\n\n')
}

export function buildCodexAuthJson(apiKey: string): string {
  return JSON.stringify({ OPENAI_API_KEY: apiKey }, null, 2)
}

const reasoningLevels = [
  { effort: 'none', description: 'No extra reasoning' },
  { effort: 'low', description: 'Fast responses with lighter reasoning' },
  { effort: 'medium', description: 'Balanced speed and reasoning' },
  { effort: 'high', description: 'Deeper reasoning for complex tasks' },
  { effort: 'xhigh', description: 'Extra high reasoning depth' }
]

export function buildCodexModelCatalogJson(modelSet: CodexModelSet = 'gpt-domestic'): string {
  const profiles = profilesForModelSet(modelSet)
  const models = profiles.map((profile, index) => ({
    slug: profile.model,
    display_name: profile.label,
    description: profile.description || `${profile.brand} text model via KQS API`,
    context_window: profile.contextWindow || 128000,
    visibility: 'list',
    supported_in_api: true,
    minimal_client_version: '0.98.0',
    input_modalities: ['text'],
    default_reasoning_level: profile.effort || 'medium',
    supported_reasoning_levels: reasoningLevels,
    shell_type: 'shell_command',
    priority: index
  }))

  const uniqueModels = Array.from(
    new Map(models.map((model) => [model.slug, model])).values()
  )

  return JSON.stringify({ models: uniqueModels }, null, 2)
}

function ownedProfilesRegex(): string {
  return CODEX_ALL_MODEL_PROFILES.map((profile) => regexEscape(profile.profile)).join('|')
}

export function buildCodexInstallScript(
  baseUrl: string,
  apiKey: string,
  options: BuildCodexConfigOptions & { shell: 'powershell' | 'bash' } = { shell: 'bash' }
): string {
  const modelSet = options.modelSet || 'gpt'
  const includeCatalog = modelSet === 'gpt-domestic'
  const configOptions = {
    ...options,
    modelSet,
    modelCatalogPath: includeCatalog ? '__CODEX_MODEL_CATALOG_PATH__' : undefined
  }
  const topLevelContent = buildCodexTopLevelToml(configOptions)
  const sectionsContent = buildCodexProviderAndProfilesToml(baseUrl, configOptions)
  const authContent = buildCodexAuthJson(apiKey)
  const catalogContent = includeCatalog ? buildCodexModelCatalogJson(modelSet) : ''
  const profileSummary = profilesForModelSet(modelSet)
    .slice(0, 8)
    .map((profile) => profile.profile)
    .join(', ')
  const profileRegex = ownedProfilesRegex()
  const providerRegex = `${regexEscape(CODEX_PROVIDER_NAME)}|${regexEscape(LEGACY_CODEX_PROVIDER_NAME)}`
  const psSectionRegex = `^\\s*\\[(model_providers\\.(?:${providerRegex})|profiles\\.(?:${profileRegex}))\\]\\s*$`
  const awkSectionRegex = `^[[:space:]]*\\[(model_providers\\.(${providerRegex})|profiles\\.(${profileRegex}))\\][[:space:]]*$`

  if (options.shell === 'powershell') {
    const catalogWriteBlock = includeCatalog
      ? `
$catalogPath = Join-Path $codexDir "${CODEX_MODEL_CATALOG_FILENAME}"
${powershellHereString(catalogContent)} | Set-Content -Path $catalogPath -Encoding UTF8
$catalogForToml = $catalogPath -replace '\\\\', '\\\\\\\\'
$topBlock = $topBlock.Replace('__CODEX_MODEL_CATALOG_PATH__', $catalogForToml)
`
      : ''

    return `$ErrorActionPreference = "Stop"
$codexDir = Join-Path $HOME ".codex"
New-Item -ItemType Directory -Force -Path $codexDir | Out-Null
$timestamp = Get-Date -Format "yyyyMMddHHmmss"
$configPath = Join-Path $codexDir "config.toml"
$authPath = Join-Path $codexDir "auth.json"
if (Test-Path $configPath) { Copy-Item $configPath "$configPath.bak-$timestamp" -Force }
if (Test-Path $authPath) { Copy-Item $authPath "$authPath.bak-$timestamp" -Force }

$topBlock = ${powershellHereString(topLevelContent)}
$sectionsBlock = ${powershellHereString(sectionsContent)}
${catalogWriteBlock}
function Remove-KqsManagedBlocks([string]$Text) {
  $lines = $Text -split '\r?\n'
  $out = New-Object System.Collections.Generic.List[string]
  $skip = $false
  foreach ($line in $lines) {
    $trimmed = $line.Trim()
    if ($trimmed -eq '${CODEX_TOP_LEVEL_MARKER_START}' -or $trimmed -eq '${CODEX_MODELS_MARKER_START}' -or $trimmed -eq '${LEGACY_CODEX_TOP_LEVEL_MARKER_START}' -or $trimmed -eq '${LEGACY_CODEX_MODELS_MARKER_START}') { $skip = $true; continue }
    if ($trimmed -eq '${CODEX_TOP_LEVEL_MARKER_END}' -or $trimmed -eq '${CODEX_MODELS_MARKER_END}' -or $trimmed -eq '${LEGACY_CODEX_TOP_LEVEL_MARKER_END}' -or $trimmed -eq '${LEGACY_CODEX_MODELS_MARKER_END}') { $skip = $false; continue }
    if (-not $skip) { $out.Add($line) }
  }
  return ($out -join [Environment]::NewLine)
}

function Remove-KqsOwnedSections([string]$Text) {
  $lines = $Text -split '\r?\n'
  $out = New-Object System.Collections.Generic.List[string]
  $skip = $false
  foreach ($line in $lines) {
    if ($line -match '${psSectionRegex}') { $skip = $true; continue }
    if ($skip -and $line -match '^\s*\[') { $skip = $false }
    if (-not $skip) { $out.Add($line) }
  }
  return ($out -join [Environment]::NewLine)
}

function Remove-CodexTopLevelKeys([string]$Text) {
  $lines = $Text -split '\r?\n'
  $out = New-Object System.Collections.Generic.List[string]
  $inTopLevel = $true
  foreach ($line in $lines) {
    if ($line -match '^\s*\[') { $inTopLevel = $false }
    if ($inTopLevel -and $line -match '^\s*(model_provider|model|review_model|model_reasoning_effort|disable_response_storage|network_access|windows_wsl_setup_acknowledged|model_context_window|model_auto_compact_token_limit|model_catalog_json)\s*=') { continue }
    $out.Add($line)
  }
  return ($out -join [Environment]::NewLine)
}

$existing = if (Test-Path $configPath) { Get-Content -Raw -Path $configPath } else { "" }
$clean = Remove-CodexTopLevelKeys (Remove-KqsOwnedSections (Remove-KqsManagedBlocks $existing))
$parts = @($topBlock.Trim(), $clean.Trim(), $sectionsBlock.Trim()) | Where-Object { $_ -ne "" }
($parts -join ([Environment]::NewLine + [Environment]::NewLine)) + [Environment]::NewLine | Set-Content -Path $configPath -Encoding UTF8
${powershellHereString(authContent)} | Set-Content -Path $authPath -Encoding UTF8

Write-Host "Codex App / CLI config installed. Restart Codex App before testing."
Write-Host "Example profiles: ${profileSummary}"`
  }

  const catalogWriteBlock = includeCatalog
    ? `
catalog_path="$codex_dir/${CODEX_MODEL_CATALOG_FILENAME}"
cat > "$catalog_path" <<'CODEX_MODEL_CATALOG'
${catalogContent}
CODEX_MODEL_CATALOG
top_block="$(printf '%s' "$top_block" | sed "s#__CODEX_MODEL_CATALOG_PATH__#$catalog_path#g")"
`
    : ''

  return `#!/usr/bin/env bash
set -euo pipefail

codex_dir="$HOME/.codex"
mkdir -p "$codex_dir"
timestamp="$(date +%Y%m%d%H%M%S)"
config_path="$codex_dir/config.toml"
auth_path="$codex_dir/auth.json"
clean_path="$codex_dir/.config.clean-$timestamp"

if [ -f "$config_path" ]; then
  cp "$config_path" "$config_path.bak-$timestamp"
fi
if [ -f "$auth_path" ]; then
  cp "$auth_path" "$auth_path.bak-$timestamp"
fi

top_block="$(cat <<'CODEX_TOP_LEVEL'
${topLevelContent}
CODEX_TOP_LEVEL
)"
sections_block="$(cat <<'CODEX_SECTIONS'
${sectionsContent}
CODEX_SECTIONS
)"
${catalogWriteBlock}
remove_managed_blocks() {
  awk '
    /^# >>> KQS_API_CODEX_TOP_LEVEL$/ { skip=1; next }
    /^# <<< KQS_API_CODEX_TOP_LEVEL$/ { skip=0; next }
    /^# >>> KQS_API_CODEX_MODELS$/ { skip=1; next }
    /^# <<< KQS_API_CODEX_MODELS$/ { skip=0; next }
    $0 == "${LEGACY_CODEX_TOP_LEVEL_MARKER_START}" { skip=1; next }
    $0 == "${LEGACY_CODEX_TOP_LEVEL_MARKER_END}" { skip=0; next }
    $0 == "${LEGACY_CODEX_MODELS_MARKER_START}" { skip=1; next }
    $0 == "${LEGACY_CODEX_MODELS_MARKER_END}" { skip=0; next }
    !skip { print }
  '
}

remove_owned_sections() {
  awk -v section_re='${awkSectionRegex}' '
    $0 ~ section_re { skip_section=1; next }
    skip_section && $0 ~ /^[[:space:]]*\[/ { skip_section=0 }
    !skip_section { print }
  '
}

remove_top_level_keys() {
  awk '
    /^[[:space:]]*\[/ { in_top=0 }
    BEGIN { in_top=1 }
    in_top && $0 ~ /^[[:space:]]*(model_provider|model|review_model|model_reasoning_effort|disable_response_storage|network_access|windows_wsl_setup_acknowledged|model_context_window|model_auto_compact_token_limit|model_catalog_json)[[:space:]]*=/ { next }
    { print }
  '
}

if [ -f "$config_path" ]; then
  remove_managed_blocks < "$config_path" | remove_owned_sections | remove_top_level_keys > "$clean_path"
else
  : > "$clean_path"
fi

{
  printf '%s\n\n' "$top_block"
  if [ -s "$clean_path" ]; then
    cat "$clean_path"
    printf '\n\n'
  fi
  printf '%s\n' "$sections_block"
} > "$config_path"
rm -f "$clean_path"

cat > "$auth_path" <<'CODEX_AUTH'
${authContent}
CODEX_AUTH

echo ${shellSingleQuote('Codex App / CLI config installed. Restart Codex App before testing.')}
echo ${shellSingleQuote(`Example profiles: ${profileSummary}`)}`
}

export function buildCodexDomesticInstallScript(
  baseUrl: string,
  apiKey: string,
  options: BuildCodexConfigOptions & { shell: 'powershell' | 'bash' }
): string {
  return buildCodexInstallScript(baseUrl, apiKey, {
    ...options,
    modelSet: 'gpt-domestic'
  })
}
