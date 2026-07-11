param(
  [switch]$RequireClean,
  [switch]$CheckProduction
)

$ErrorActionPreference = "Stop"
$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
$failures = [System.Collections.Generic.List[string]]::new()

function Assert-Command {
  param([string]$Name)
  if (-not (Get-Command $Name -ErrorAction SilentlyContinue)) {
    $script:failures.Add("Missing required command: $Name")
    return
  }
  Write-Host "[PASS] command: $Name"
}

function Assert-Path {
  param([string]$RelativePath)
  if (-not (Test-Path (Join-Path $repoRoot $RelativePath))) {
    $script:failures.Add("Missing required handoff file: $RelativePath")
    return
  }
  Write-Host "[PASS] file: $RelativePath"
}

Push-Location $repoRoot
try {
  "git", "node", "corepack", "go" | ForEach-Object { Assert-Command $_ }

  $requiredFiles = @(
    ".github/workflows/ci.yml",
    ".github/workflows/release-production.yml",
    "ops/production/deploy-remote.sh",
    "ops/production/verify-public.sh",
    "docs/OPERATIONS_HANDOFF_CN.md",
    "docs/RELEASE_PIPELINE_CN.md",
    "docs/SECRETS_INVENTORY_TEMPLATE.md"
  )
  $requiredFiles | ForEach-Object { Assert-Path $_ }

  $origin = git remote get-url origin 2>$null
  if ($LASTEXITCODE -ne 0 -or [string]::IsNullOrWhiteSpace($origin)) {
    $failures.Add("Git remote 'origin' is not configured.")
  } elseif ($origin -match "Wei-Shaw/sub2api") {
    $failures.Add("origin still points at upstream ($origin). Set origin to the private KQS repository before sharing this checkout.")
  } else {
    Write-Host "[PASS] private origin: $origin"
  }

  $status = @(git status --short)
  if ($RequireClean -and $status.Count -gt 0) {
    $failures.Add("Worktree is dirty ($($status.Count) entries). Commit or stash before a handoff/release.")
  } elseif ($status.Count -gt 0) {
    Write-Host "[WARN] worktree has $($status.Count) changed entries. This is expected only while actively developing."
  } else {
    Write-Host "[PASS] worktree is clean"
  }

  if ($CheckProduction) {
    if ([string]::IsNullOrWhiteSpace($env:KQS_PROD_SSH_HOST)) {
      $failures.Add("KQS_PROD_SSH_HOST is required with -CheckProduction. Use an SSH host alias; do not put a private key in this repository.")
    } else {
      ssh -o BatchMode=yes -o StrictHostKeyChecking=yes $env:KQS_PROD_SSH_HOST "docker inspect -f '{{.State.Health.Status}}' sub2api-gg; curl -fsS --max-time 8 http://127.0.0.1:18080/health >/dev/null"
      if ($LASTEXITCODE -ne 0) {
        $failures.Add("Production SSH/app health check failed for $($env:KQS_PROD_SSH_HOST).")
      } else {
        Write-Host "[PASS] production app container and local health endpoint"
      }
    }
  }
} finally {
  Pop-Location
}

if ($failures.Count -gt 0) {
  $failures | ForEach-Object { Write-Error "[FAIL] $_" }
  exit 1
}

Write-Host "Handoff preflight passed." -ForegroundColor Green
