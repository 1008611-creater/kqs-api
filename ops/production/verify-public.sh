#!/usr/bin/env bash
set -euo pipefail

BASE_URL=${1:-https://api.cauai.fun}
SMOKE_API_KEY=${SMOKE_API_KEY:-}
REQUIRE_SMOKE_API_KEY=${REQUIRE_SMOKE_API_KEY:-0}

for path in /health /dashboard /login /keys /monitor /guide; do
  code=$(curl -sS -o /dev/null -w '%{http_code}' --max-time 20 "$BASE_URL$path")
  [[ "$code" == "200" ]] || { echo "$path returned $code" >&2; exit 1; }
  echo "public route ok: $path"
done

if [[ -n "$SMOKE_API_KEY" ]]; then
  curl -fsS --max-time 30 -H "Authorization: Bearer $SMOKE_API_KEY" "$BASE_URL/v1/models" >/dev/null
  echo "authenticated models smoke ok"
elif [[ "$REQUIRE_SMOKE_API_KEY" == "1" ]]; then
  echo "authenticated API smoke is required but SMOKE_API_KEY is not set" >&2
  exit 1
else
  echo "authenticated API smoke skipped: SMOKE_API_KEY is not set"
fi
