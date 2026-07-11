param(
  [string]$HostName = "api.cauai.fun",
  [string]$HealthPath = "/health",
  [int]$Samples = 2,
  [string]$OutputDir = "",
  [string]$SshHost = "103.236.92.40",
  [int]$SshPort = 38961,
  [string]$SshUser = "root",
  [string]$SshKey = "C:\Users\lsb\.ssh\sub2api_mingkj_ed25519",
  [string]$KnownHosts = "D:\codex-work\gg\.ssh_known_hosts_sub2api",
  [switch]$SkipRemote
)

$ErrorActionPreference = "Stop"
$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
if ([string]::IsNullOrWhiteSpace($OutputDir)) {
  $OutputDir = Join-Path $repoRoot ".tools\line-speed"
}
New-Item -ItemType Directory -Force -Path $OutputDir | Out-Null

$stamp = Get-Date -Format "yyyyMMdd-HHmmss"
$edgeCsv = Join-Path $OutputDir "edge-candidates-$stamp.csv"
$mihomoCsv = Join-Path $OutputDir "mihomo-nodes-$stamp.csv"
$reportHtml = Join-Path $OutputDir "line-speed-report-$stamp.html"

function Join-Url {
  param([string]$BaseHost, [string]$Path)
  return "https://$BaseHost/$($Path.TrimStart('/'))"
}

function Convert-ToDoubleOrNull {
  param([string]$Value)
  $parsed = 0.0
  if ([double]::TryParse($Value, [Globalization.NumberStyles]::Float, [Globalization.CultureInfo]::InvariantCulture, [ref]$parsed)) {
    return $parsed
  }
  return $null
}

function Invoke-CurlTiming {
  param(
    [string]$Url,
    [string]$ResolveHost,
    [string]$Ip,
    [string]$Label,
    [int]$Sample
  )

  $resolveArgs = @()
  if (-not [string]::IsNullOrWhiteSpace($Ip)) {
    $resolveArgs = @("--resolve", "$ResolveHost`:443:$Ip")
  }

  $line = & curl.exe -k -sS --connect-timeout 5 --max-time 12 @resolveArgs -o NUL `
    -w "code=%{http_code} remote_ip=%{remote_ip} dns=%{time_namelookup} conn=%{time_connect} tls=%{time_appconnect} ttfb=%{time_starttransfer} total=%{time_total}`n" `
    $Url 2>$null

  $match = [regex]::Match($line, "code=(?<code>\d+) remote_ip=(?<remote>\S*) dns=(?<dns>[0-9.]+) conn=(?<conn>[0-9.]+) tls=(?<tls>[0-9.]+) ttfb=(?<ttfb>[0-9.]+) total=(?<total>[0-9.]+)")
  if (-not $match.Success) {
    return [pscustomobject]@{
      label = $Label
      ip = $Ip
      sample = $Sample
      code = 0
      remote_ip = ""
      dns_ms = $null
      connect_ms = $null
      tls_ms = $null
      ttfb_ms = $null
      total_ms = $null
      raw = $line
    }
  }

  return [pscustomobject]@{
    label = $Label
    ip = $Ip
    sample = $Sample
    code = [int]$match.Groups["code"].Value
    remote_ip = $match.Groups["remote"].Value
    dns_ms = [math]::Round((Convert-ToDoubleOrNull $match.Groups["dns"].Value) * 1000, 1)
    connect_ms = [math]::Round((Convert-ToDoubleOrNull $match.Groups["conn"].Value) * 1000, 1)
    tls_ms = [math]::Round((Convert-ToDoubleOrNull $match.Groups["tls"].Value) * 1000, 1)
    ttfb_ms = [math]::Round((Convert-ToDoubleOrNull $match.Groups["ttfb"].Value) * 1000, 1)
    total_ms = [math]::Round((Convert-ToDoubleOrNull $match.Groups["total"].Value) * 1000, 1)
    raw = $line.Trim()
  }
}

function Get-EdgeCandidates {
  $current = @()
  try {
    $current = Resolve-DnsName $HostName -Type A -ErrorAction Stop |
      Where-Object { $_.IPAddress } |
      ForEach-Object { $_.IPAddress }
  } catch {
    $current = @()
  }

  $knownFast = @(
    "104.19.0.1", "104.20.0.1", "104.24.0.1", "104.26.0.1", "104.27.0.1",
    "104.16.0.1", "104.17.0.1", "104.18.0.1",
    "162.159.129.1", "162.159.130.1",
    "172.67.147.202", "104.21.28.236"
  )

  $items = @()
  foreach ($ip in ($current + $knownFast | Select-Object -Unique)) {
    $label = if ($current -contains $ip) { "current-dns" } else { "candidate" }
    $items += [pscustomobject]@{ label = $label; ip = $ip }
  }
  return $items
}

function Write-HtmlReport {
  param(
    [object[]]$EdgeRows,
    [object[]]$EdgeSummary,
    [object[]]$MihomoRows,
    [string]$Path
  )

  $bestEdge = $EdgeSummary | Where-Object { $_.success -gt 0 } | Sort-Object avg_total_ms | Select-Object -First 5
  $bestMihomo = $MihomoRows | Where-Object { $_.ok -eq $true -and $_.delay_ms -ge 0 } | Sort-Object delay_ms | Select-Object -First 20

  $style = @"
body{font-family:Arial,Helvetica,sans-serif;margin:24px;color:#111827;background:#f8fafc}
h1,h2{margin:0 0 12px}
section{background:white;border:1px solid #e5e7eb;border-radius:8px;padding:16px;margin:0 0 16px}
table{border-collapse:collapse;width:100%;font-size:13px}
th,td{border-bottom:1px solid #e5e7eb;text-align:left;padding:8px}
th{background:#f3f4f6}
.ok{color:#047857;font-weight:700}.bad{color:#b91c1c;font-weight:700}
.muted{color:#6b7280}
"@

  $edgeHtml = ($bestEdge | ConvertTo-Html -Fragment -Property ip,label,success,avg_total_ms,min_total_ms,avg_ttfb_ms | Out-String)
  $mihomoHtml = ($bestMihomo | ConvertTo-Html -Fragment -Property name,type,group,now,delay_ms,url,error | Out-String)
  $rawEdgeHtml = ($EdgeRows | Sort-Object ip,sample | ConvertTo-Html -Fragment -Property label,ip,sample,code,remote_ip,dns_ms,connect_ms,tls_ms,ttfb_ms,total_ms | Out-String)

  $html = @"
<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <title>Sub2api line speed report $stamp</title>
  <style>$style</style>
</head>
<body>
  <h1>Sub2api line speed report</h1>
  <p class="muted">Generated: $stamp. Host: $HostName. Samples per edge IP: $Samples.</p>
  <section>
    <h2>Recommended Cloudflare edge IPs</h2>
    $edgeHtml
  </section>
  <section>
    <h2>Fastest mihomo proxy nodes</h2>
    $mihomoHtml
  </section>
  <section>
    <h2>Raw edge measurements</h2>
    $rawEdgeHtml
  </section>
</body>
</html>
"@
  Set-Content -Path $Path -Value $html -Encoding UTF8
}

$url = Join-Url $HostName $HealthPath
$edgeRows = New-Object System.Collections.Generic.List[object]
foreach ($candidate in Get-EdgeCandidates) {
  for ($i = 1; $i -le $Samples; $i++) {
    $row = Invoke-CurlTiming -Url $url -ResolveHost $HostName -Ip $candidate.ip -Label $candidate.label -Sample $i
    $edgeRows.Add($row)
    Write-Host ("edge {0,-14} {1,-11} sample={2} code={3} total={4}ms ttfb={5}ms" -f $candidate.ip, $candidate.label, $i, $row.code, $row.total_ms, $row.ttfb_ms)
  }
}

$edgeRows | Export-Csv -NoTypeInformation -Encoding UTF8 -Path $edgeCsv

$edgeSummary = $edgeRows |
  Group-Object ip |
  ForEach-Object {
    $ok = @($_.Group | Where-Object { $_.code -eq 200 -and $null -ne $_.total_ms })
    [pscustomobject]@{
      ip = $_.Name
      label = ($_.Group | Select-Object -First 1).label
      success = $ok.Count
      avg_total_ms = if ($ok.Count) { [math]::Round((($ok | Measure-Object total_ms -Average).Average), 1) } else { $null }
      min_total_ms = if ($ok.Count) { [math]::Round((($ok | Measure-Object total_ms -Minimum).Minimum), 1) } else { $null }
      avg_ttfb_ms = if ($ok.Count) { [math]::Round((($ok | Measure-Object ttfb_ms -Average).Average), 1) } else { $null }
    }
  } |
  Sort-Object @{ Expression = { if ($null -eq $_.avg_total_ms) { [double]::PositiveInfinity } else { $_.avg_total_ms } } }

$mihomoRows = @()
if (-not $SkipRemote) {
  $sshArgs = @(
    "-i", $SshKey,
    "-o", "UserKnownHostsFile=$KnownHosts",
    "-o", "StrictHostKeyChecking=yes",
    "-p", "$SshPort",
    "$SshUser@$SshHost"
  )

  $remotePython = @'
import json
import time
import urllib.parse
import urllib.request

controller = "http://127.0.0.1:9090"
test_url = "http://cp.cloudflare.com/generate_204"

def fetch_json(url, timeout=8):
    with urllib.request.urlopen(url, timeout=timeout) as resp:
        return json.loads(resp.read().decode("utf-8"))

proxies = fetch_json(controller + "/proxies")["proxies"]
wanted_groups = {"SSRDOG", "Auto", "GLOBAL", "🇭🇰 Hong Kong", "🇯🇵 Japan", "🇸🇬 Singapore", "🇺🇸 United States", "🇨🇳 Taiwan", "🇬🇧 Great Britain", "🇨🇦 Canada"}
region_tokens = ("Hong Kong", "Japan", "Singapore", "United States", "Taiwan", "Great Britain", "Canada", "Korea")

names = []
for name, p in proxies.items():
    ptype = p.get("type", "")
    if name in wanted_groups:
        names.append(name)
    elif any(token in name for token in region_tokens) and ptype not in ("Selector", "Fallback", "Direct", "Reject", "Compatible", "Pass"):
        names.append(name)

seen = set()
rows = []
for name in names[:90]:
    if name in seen:
        continue
    seen.add(name)
    p = proxies.get(name, {})
    encoded = urllib.parse.quote(name, safe="")
    started = time.time()
    ok = False
    delay = -1
    err = ""
    try:
        result = fetch_json(f"{controller}/proxies/{encoded}/delay?timeout=5000&url={urllib.parse.quote(test_url, safe='')}", timeout=8)
        delay = int(result.get("delay", -1))
        ok = delay >= 0
    except Exception as exc:
        err = str(exc)[:160]
    rows.append({
        "name": name,
        "type": p.get("type", ""),
        "group": p.get("provider-name", ""),
        "now": p.get("now", ""),
        "ok": ok,
        "delay_ms": delay,
        "elapsed_ms": int((time.time() - started) * 1000),
        "url": test_url,
        "error": err,
    })

print(json.dumps(rows, ensure_ascii=False))
'@

  $remotePayload = [Convert]::ToBase64String([Text.Encoding]::UTF8.GetBytes(($remotePython -replace "`r", "")))
  $remoteCommand = "python3 -c 'import base64; exec(base64.b64decode(""$remotePayload"").decode(""utf-8""))'"
  $remoteJson = & ssh.exe @sshArgs $remoteCommand
  $mihomoRows = @($remoteJson | ConvertFrom-Json)
  $mihomoRows | Export-Csv -NoTypeInformation -Encoding UTF8 -Path $mihomoCsv
}

Write-HtmlReport -EdgeRows @($edgeRows.ToArray()) -EdgeSummary @($edgeSummary) -MihomoRows @($mihomoRows) -Path $reportHtml

Write-Host ""
Write-Host "Top Cloudflare edge candidates:"
$edgeSummary | Select-Object -First 8 | Format-Table -AutoSize
if ($mihomoRows.Count -gt 0) {
  Write-Host ""
  Write-Host "Top mihomo nodes:"
  $mihomoRows | Where-Object { $_.ok -eq $true -and $_.delay_ms -ge 0 } | Sort-Object delay_ms | Select-Object -First 12 name,type,now,delay_ms,error | Format-Table -AutoSize
}
Write-Host ""
Write-Host "Wrote:"
Write-Host "  $edgeCsv"
if ($mihomoRows.Count -gt 0) {
  Write-Host "  $mihomoCsv"
}
Write-Host "  $reportHtml"
