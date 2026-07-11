export const normalizeSupportedModelScopesForPlatform = (
  platform: string,
  scopes: string[] | undefined,
): string[] => {
  if (platform === "openai") {
    return (scopes ?? []).filter((scope) => scope === "codex_all_models");
  }
  if (platform !== "antigravity") return [];
  return scopes ?? [];
};
