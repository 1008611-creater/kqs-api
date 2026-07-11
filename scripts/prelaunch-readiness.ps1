param(
  [string]$ApiBaseUrl = "https://api.cauai.fun",
  [string]$AdminToken = "",
  [string]$AdminEmail = "",
  [string]$AdminPassword = "",
  [string]$AdminLoginBaseUrl = "",
  [string]$AdminApiBaseUrl = "",
  [string]$EnvFile = "",
  [int]$BackupFreshnessHours = 30,
  [switch]$AllowDirty,
  [switch]$AllowUnhardenedSecurity,
  [switch]$SkipNetwork
)

$ErrorActionPreference = "Stop"
$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
$results = New-Object System.Collections.Generic.List[object]

function Add-Check {
  param(
    [ValidateSet("PASS", "WARN", "FAIL")]
    [string]$Status,
    [string]$Name,
    [string]$Detail
  )
  $results.Add([pscustomobject]@{
    Status = $Status
    Name = $Name
    Detail = $Detail
  })
}

function Join-Url {
  param([string]$Base, [string]$Path)
  return ($Base.TrimEnd("/") + "/" + $Path.TrimStart("/"))
}

function Read-DotEnvFile {
  param([string]$Path)

  $values = @{}
  if (-not (Test-Path $Path)) {
    return $values
  }

  Get-Content $Path | ForEach-Object {
    $line = $_.Trim()
    if ([string]::IsNullOrWhiteSpace($line) -or $line.StartsWith("#") -or -not $line.Contains("=")) {
      return
    }

    $idx = $line.IndexOf("=")
    if ($idx -le 0) {
      return
    }

    $key = $line.Substring(0, $idx).Trim()
    $value = $line.Substring($idx + 1).Trim()
    if (($value.StartsWith('"') -and $value.EndsWith('"')) -or ($value.StartsWith("'") -and $value.EndsWith("'"))) {
      $value = $value.Substring(1, $value.Length - 2)
    }
    $values[$key] = $value
  }

  return $values
}

function Resolve-AdminToken {
  param(
    [string]$ProvidedToken,
    [string]$RepoRoot,
    [string]$EnvFilePath,
    [string]$Email,
    [string]$Password,
    [string]$LoginBaseUrl
  )

  if (-not [string]::IsNullOrWhiteSpace($ProvidedToken)) {
    Add-Check "PASS" "admin auth" "Using provided admin token for backup checks."
    return $ProvidedToken
  }

  $resolvedEnvFile = $EnvFilePath
  if ([string]::IsNullOrWhiteSpace($resolvedEnvFile)) {
    $resolvedEnvFile = Join-Path $RepoRoot "deploy/.env"
  }

  $envValues = Read-DotEnvFile $resolvedEnvFile
  if ([string]::IsNullOrWhiteSpace($Email) -and $envValues.ContainsKey("ADMIN_EMAIL")) {
    $Email = $envValues["ADMIN_EMAIL"]
  }
  if ([string]::IsNullOrWhiteSpace($Password) -and $envValues.ContainsKey("ADMIN_PASSWORD")) {
    $Password = $envValues["ADMIN_PASSWORD"]
  }
  if ([string]::IsNullOrWhiteSpace($LoginBaseUrl)) {
    $port = "8080"
    if ($envValues.ContainsKey("SERVER_PORT") -and -not [string]::IsNullOrWhiteSpace($envValues["SERVER_PORT"])) {
      $port = $envValues["SERVER_PORT"]
    }
    $LoginBaseUrl = "http://localhost:$port"
  }

  if ([string]::IsNullOrWhiteSpace($Email) -or [string]::IsNullOrWhiteSpace($Password)) {
    Add-Check "WARN" "admin auth" "Admin token was not provided and local admin credentials were not found; skipped authenticated backup checks."
    return ""
  }

  try {
    $loginUrl = Join-Url $LoginBaseUrl "/api/v1/auth/login"
    $body = @{ email = $Email; password = $Password } | ConvertTo-Json
    $login = Invoke-RestMethod -Uri $loginUrl -Method Post -ContentType "application/json" -Body $body -TimeoutSec 20
    $data = if ($login.data) { $login.data } else { $login }
    if ($data.requires_2fa) {
      Add-Check "FAIL" "admin auth" "Admin account requires 2FA; pass -AdminToken to run backup checks."
      return ""
    }
    if ([string]::IsNullOrWhiteSpace($data.access_token)) {
      Add-Check "FAIL" "admin auth" "Admin login succeeded but no access token was returned."
      return ""
    }
    Add-Check "PASS" "admin auth" "Obtained a temporary admin token from local credentials."
    return $data.access_token
  } catch {
    Add-Check "FAIL" "admin auth" "Admin login failed: $($_.Exception.Message)"
    return ""
  }
}

function Resolve-AdminApiBaseUrl {
  param(
    [string]$RepoRoot,
    [string]$EnvFilePath,
    [string]$LoginBaseUrl,
    [string]$ExplicitAdminApiBaseUrl,
    [string]$FallbackApiBaseUrl
  )

  if (-not [string]::IsNullOrWhiteSpace($ExplicitAdminApiBaseUrl)) {
    return $ExplicitAdminApiBaseUrl
  }
  if (-not [string]::IsNullOrWhiteSpace($LoginBaseUrl)) {
    return $LoginBaseUrl
  }

  $resolvedEnvFile = $EnvFilePath
  if ([string]::IsNullOrWhiteSpace($resolvedEnvFile)) {
    $resolvedEnvFile = Join-Path $RepoRoot "deploy/.env"
  }
  $envValues = Read-DotEnvFile $resolvedEnvFile
  if ($envValues.ContainsKey("SERVER_PORT") -and -not [string]::IsNullOrWhiteSpace($envValues["SERVER_PORT"])) {
    return "http://localhost:$($envValues["SERVER_PORT"])"
  }

  return $FallbackApiBaseUrl
}

Push-Location $repoRoot
try {
  $gitStatus = git status --short 2>$null
  if ($LASTEXITCODE -ne 0) {
    Add-Check "WARN" "git status" "Git is unavailable or this is not a git checkout."
  } elseif ($gitStatus) {
    $lineCount = ($gitStatus | Measure-Object).Count
    if ($AllowDirty) {
      Add-Check "WARN" "git status" "Worktree has $lineCount changed entries; allowed for local inspection."
    } else {
      Add-Check "FAIL" "git status" "Worktree has $lineCount changed entries. Freeze or separate changes before release."
    }
  } else {
    Add-Check "PASS" "git status" "Worktree is clean."
  }

  $requiredFiles = @(
    "docs/RELEASE_RUNBOOK_CN.md",
    "docs/PRELAUNCH_GATE_CN.md",
    "docs/BACKUP_RESTORE_DRILL_CN.md",
    "docs/R2_BACKUP_SETUP_CN.md",
    "frontend/src/constants/legalDocuments.ts",
    "frontend/src/views/user/GuideView.vue",
    "scripts/prelaunch-readiness.ps1"
  )
  foreach ($file in $requiredFiles) {
    if (Test-Path (Join-Path $repoRoot $file)) {
      Add-Check "PASS" "required file" $file
    } else {
      Add-Check "FAIL" "required file" "$file is missing."
    }
  }

  $compose = Get-Content (Join-Path $repoRoot "deploy/docker-compose.yml") -Raw
  if ($compose -match '\$\{SUB2API_IMAGE:-') {
    Add-Check "PASS" "fixed image config" "deploy/docker-compose.yml supports SUB2API_IMAGE."
  } else {
    Add-Check "FAIL" "fixed image config" "docker-compose still hardcodes a single image tag."
  }

  $buildScript = Get-Content (Join-Path $repoRoot "deploy/build_image.sh") -Raw
  if ($buildScript -match 'sub2api:kqs-api-\$\{SHORT_SHA\}' -or $buildScript -match 'Built image:') {
    Add-Check "PASS" "image build script" "deploy/build_image.sh can build a fixed tag."
  } else {
    Add-Check "WARN" "image build script" "Could not detect fixed-tag output in build script."
  }

  $docker = Get-Command docker -ErrorAction SilentlyContinue
  if ($docker) {
    $imageName = docker inspect sub2api-gg --format '{{.Config.Image}}' 2>$null
    if ($LASTEXITCODE -eq 0 -and $imageName) {
      if ($imageName -match ':latest$') {
        Add-Check "FAIL" "running image" "sub2api-gg is running $imageName. Release should use a fixed tag."
      } else {
        Add-Check "PASS" "running image" "sub2api-gg is running $imageName."
      }

      $containerLogs = @(docker logs sub2api-gg 2>&1)
      $urlAllowlistDisabled = @($containerLogs | Select-String -Pattern "security\.url_allowlist\.enabled=false").Count -gt 0
      $trustedProxiesEmpty = @($containerLogs | Select-String -Pattern "server\.trusted_proxies is empty").Count -gt 0
      $corsOriginsEmpty = @($containerLogs | Select-String -Pattern "CORS allowed_origins not configured").Count -gt 0

      if ($urlAllowlistDisabled) {
        if ($AllowUnhardenedSecurity) {
          Add-Check "WARN" "url allowlist" "URL allowlist is disabled; accepted only for controlled Beta while required upstream hosts are being confirmed."
        } else {
          Add-Check "FAIL" "url allowlist" "URL allowlist is disabled. Confirm required upstream hosts and enable it before public release."
        }
      } else {
        Add-Check "PASS" "url allowlist" "No disabled URL allowlist warning was observed in container startup logs."
      }

      if ($trustedProxiesEmpty) {
        if ($AllowUnhardenedSecurity) {
          Add-Check "WARN" "trusted client ip" "Trusted proxy chain is not configured; accepted only for controlled Beta until Cloudflare IP forwarding is verified."
        } else {
          Add-Check "FAIL" "trusted client ip" "Trusted proxy chain is not configured. Define and verify Cloudflare client IP forwarding before public release."
        }
      } else {
        Add-Check "PASS" "trusted client ip" "No empty trusted proxy warning was observed in container startup logs."
      }

      if ($corsOriginsEmpty) {
        Add-Check "PASS" "cors policy" "Cross-origin browser requests are denied by default; accepted for the current same-origin UI/API deployment."
      } else {
        Add-Check "PASS" "cors policy" "No empty allowed-origins warning was observed; verify any configured browser origins are intentional."
      }
    } else {
      Add-Check "WARN" "running image" "Container sub2api-gg was not found locally."
    }
  } else {
    Add-Check "WARN" "docker" "Docker CLI not found; skipped container image check."
  }

  if (-not $SkipNetwork) {
    try {
      $healthUrl = Join-Url $ApiBaseUrl "/health"
      $health = Invoke-RestMethod -Uri $healthUrl -Method Get -TimeoutSec 12
      $healthStatus = if ($health.PSObject.Properties.Name -contains "status") { $health.status } else { "ok" }
      Add-Check "PASS" "public health" "$healthUrl returned $healthStatus."
    } catch {
      Add-Check "FAIL" "public health" "Health check failed for ${ApiBaseUrl}: $($_.Exception.Message)"
    }

    $effectiveAdminToken = Resolve-AdminToken -ProvidedToken $AdminToken -RepoRoot $repoRoot -EnvFilePath $EnvFile -Email $AdminEmail -Password $AdminPassword -LoginBaseUrl $AdminLoginBaseUrl
    if ($effectiveAdminToken) {
      try {
        $headers = @{ Authorization = "Bearer $effectiveAdminToken" }
        $adminBaseUrl = Resolve-AdminApiBaseUrl -RepoRoot $repoRoot -EnvFilePath $EnvFile -LoginBaseUrl $AdminLoginBaseUrl -ExplicitAdminApiBaseUrl $AdminApiBaseUrl -FallbackApiBaseUrl $ApiBaseUrl
        Add-Check "PASS" "admin api" "Using $adminBaseUrl for authenticated backup checks."

        $s3Url = Join-Url $adminBaseUrl "/api/v1/admin/backups/s3-config"
        $s3Response = Invoke-RestMethod -Uri $s3Url -Method Get -Headers $headers -TimeoutSec 20
        $s3Config = if ($s3Response.data) { $s3Response.data } else { $s3Response }
        $s3Missing = @()
        if ([string]::IsNullOrWhiteSpace($s3Config.endpoint)) { $s3Missing += "endpoint" }
        if ([string]::IsNullOrWhiteSpace($s3Config.bucket)) { $s3Missing += "bucket" }
        if ([string]::IsNullOrWhiteSpace($s3Config.access_key_id)) { $s3Missing += "access_key_id" }
        if ($s3Missing.Count -gt 0) {
          Add-Check "FAIL" "backup s3 config" "Missing required S3/R2 fields: $($s3Missing -join ', ')."
        } else {
          Add-Check "PASS" "backup s3 config" "S3/R2 config is present for bucket $($s3Config.bucket)."
          try {
            $testUrl = Join-Url $adminBaseUrl "/api/v1/admin/backups/s3-config/test"
            $testBody = $s3Config | ConvertTo-Json -Depth 8
            $testResponse = Invoke-RestMethod -Uri $testUrl -Method Post -Headers $headers -ContentType "application/json" -Body $testBody -TimeoutSec 30
            $testData = if ($testResponse.data) { $testResponse.data } else { $testResponse }
            if ($testData.ok -eq $true) {
              Add-Check "PASS" "backup s3 connection" "S3/R2 connection test passed."
            } else {
              Add-Check "FAIL" "backup s3 connection" "S3/R2 connection test failed: $($testData.message)"
            }
          } catch {
            Add-Check "FAIL" "backup s3 connection" "S3/R2 connection test failed: $($_.Exception.Message)"
          }
        }

        $scheduleUrl = Join-Url $adminBaseUrl "/api/v1/admin/backups/schedule"
        $scheduleResponse = Invoke-RestMethod -Uri $scheduleUrl -Method Get -Headers $headers -TimeoutSec 20
        $schedule = if ($scheduleResponse.data) { $scheduleResponse.data } else { $scheduleResponse }
        if ($schedule.enabled -eq $true -and -not [string]::IsNullOrWhiteSpace($schedule.cron_expr)) {
          Add-Check "PASS" "backup schedule" "Scheduled backup is enabled: $($schedule.cron_expr)."
          if ($schedule.cron_expr -ne "30 2 * * *") {
            Add-Check "WARN" "backup schedule time" "Recommended Beijing 02:30 cron is '30 2 * * *'; current is '$($schedule.cron_expr)'."
          }
          if ([int]$schedule.retain_days -lt 14 -and [int]$schedule.retain_count -lt 30) {
            Add-Check "WARN" "backup retention" "Recommended retention is at least 14 days or 30 copies."
          } else {
            Add-Check "PASS" "backup retention" "Retention is retain_days=$($schedule.retain_days), retain_count=$($schedule.retain_count)."
          }
        } else {
          Add-Check "FAIL" "backup schedule" "Scheduled backup is disabled or missing cron expression."
        }

        $backupUrl = Join-Url $adminBaseUrl "/api/v1/admin/backups"
        $backupResponse = Invoke-RestMethod -Uri $backupUrl -Method Get -Headers $headers -TimeoutSec 20
        $items = @()
        if ($backupResponse.data -and $backupResponse.data.items) {
          $items = @($backupResponse.data.items)
        } elseif ($backupResponse.items) {
          $items = @($backupResponse.items)
        }

        $completed = @($items | Where-Object { $_.status -eq "completed" } | Sort-Object {
          if ($_.finished_at) { [datetime]$_.finished_at } else { [datetime]$_.started_at }
        } -Descending)

        if ($completed.Count -eq 0) {
          Add-Check "FAIL" "backup freshness" "No completed backup was found."
        } else {
          $latest = $completed[0]
          $finishedAt = if ($latest.finished_at) { [datetime]$latest.finished_at } else { [datetime]$latest.started_at }
          $ageHours = ((Get-Date).ToUniversalTime() - $finishedAt.ToUniversalTime()).TotalHours
          if ($ageHours -le $BackupFreshnessHours) {
            Add-Check "PASS" "backup freshness" "Latest completed backup $($latest.id) is $([math]::Round($ageHours, 1)) hours old."
          } else {
            Add-Check "FAIL" "backup freshness" "Latest completed backup $($latest.id) is $([math]::Round($ageHours, 1)) hours old."
          }
        }

        $failed = @($items | Where-Object { $_.status -eq "failed" })
        if ($failed.Count -gt 0) {
          Add-Check "WARN" "backup failures" "$($failed.Count) failed backup records exist. Review before public launch."
        } else {
          Add-Check "PASS" "backup failures" "No failed backup records returned by API."
        }
      } catch {
        Add-Check "FAIL" "backup api" "Backup API check failed: $($_.Exception.Message)"
      }
    } else {
      Add-Check "WARN" "backup checks" "Skipped backup API checks because no admin token was available."
    }
  } else {
    Add-Check "WARN" "network checks" "Skipped by -SkipNetwork."
  }
} finally {
  Pop-Location
}

$results | Format-Table -AutoSize

$failures = @($results | Where-Object { $_.Status -eq "FAIL" })
if ($failures.Count -gt 0) {
  Write-Host ""
  Write-Host "Prelaunch readiness failed with $($failures.Count) blocking issue(s)." -ForegroundColor Red
  exit 1
}

Write-Host ""
Write-Host "Prelaunch readiness passed with no blocking failures." -ForegroundColor Green
