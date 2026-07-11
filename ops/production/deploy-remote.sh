#!/usr/bin/env bash
set -euo pipefail

# Called by the protected production workflow after it has uploaded a Linux
# binary and the matching frontend dist archive into RELEASE_DIR.
RELEASE_DIR=${1:?release directory is required}
RELEASE_TAG=${2:?release image tag is required}
EXPECTED_SHA256=${3:?expected binary sha256 is required}
APP_DIR=${APP_DIR:-/opt/sub2api}
APP_SERVICE=${APP_SERVICE:-sub2api}
APP_CONTAINER=${APP_CONTAINER:-sub2api-gg}
HEALTH_URL=${HEALTH_URL:-http://127.0.0.1:18080/health}

binary_path="$RELEASE_DIR/sub2api"
frontend_archive="$RELEASE_DIR/frontend-dist.tar.gz"
env_path="$APP_DIR/.env"
backup_dir="$APP_DIR/backups/releases"
release_id=$(basename "$RELEASE_DIR")
previous_image=""
released=0
env_changed=0
public_swapped=0
public_dir="$APP_DIR/data/public"
public_stage="$APP_DIR/data/.public-stage-$release_id"
public_previous_dir=""

require_file() {
  [[ -s "$1" ]] || { echo "required file is missing or empty: $1" >&2; exit 1; }
}

rollback() {
  if [[ "$released" -eq 1 || -z "$previous_image" ]]; then
    return
  fi

  echo "release failed; restoring previous release state" >&2
  if [[ "$public_swapped" -eq 1 && -d "$public_previous_dir" ]]; then
    rm -rf "$public_dir"
    mv "$public_previous_dir" "$public_dir"
    public_swapped=0
  fi
  rm -rf "$public_stage"

  if [[ "$env_changed" -ne 1 ]]; then
    return
  fi
  echo "release failed; restoring $previous_image" >&2
  sed -i "s|^SUB2API_IMAGE=.*|SUB2API_IMAGE=$previous_image|" "$env_path"
  (cd "$APP_DIR" && docker compose up -d --no-deps --force-recreate "$APP_SERVICE") || true
}

trap rollback ERR

require_file "$binary_path"
require_file "$frontend_archive"
[[ -d "$RELEASE_DIR" ]] || { echo "release directory does not exist: $RELEASE_DIR" >&2; exit 1; }
[[ "$RELEASE_TAG" =~ ^[A-Za-z0-9][A-Za-z0-9._/:@-]*$ ]] || {
  echo "invalid release image tag" >&2
  exit 1
}
[[ "$(sha256sum "$binary_path" | awk '{print $1}')" == "$EXPECTED_SHA256" ]] || {
  echo "uploaded binary sha256 mismatch" >&2
  exit 1
}

grep -q '^SUB2API_IMAGE=' "$env_path" || {
  echo "SUB2API_IMAGE is missing from $env_path" >&2
  exit 1
}
previous_image=$(sed -n 's/^SUB2API_IMAGE=//p' "$env_path" | tail -n 1)
[[ -n "$previous_image" ]] || { echo "previous image is empty" >&2; exit 1; }

mkdir -p "$backup_dir" "$public_dir"
timestamp=$(date +%Y%m%d_%H%M%S)
cp "$env_path" "$backup_dir/.env.before-${release_id}-${timestamp}"
tar -C "$public_dir" -czf "$backup_dir/frontend-public.before-${release_id}-${timestamp}.tar.gz" .

# The active image is the base. This avoids pulling a Docker build toolchain on
# the production server and preserves the known runtime dependencies.
printf '%s\n' \
  "FROM $previous_image" \
  'USER root' \
  'COPY sub2api /app/sub2api' \
  'RUN chmod 0755 /app/sub2api' \
  | docker build -t "$RELEASE_TAG" -f - "$RELEASE_DIR"

docker run --rm --entrypoint /bin/sh "$RELEASE_TAG" -c "test -x /app/sub2api && grep -a -q '<html' /app/sub2api"

# Production currently uses data/public as an intentional override. Build the
# replacement in a staging directory, then swap it only after all checks pass.
# This avoids exposing a partially extracted frontend during a release.
rm -rf "$public_stage"
mkdir -p "$public_stage"
tar -C "$public_stage" -xzf "$frontend_archive"
main_asset=$(grep -o '/assets/index-[A-Za-z0-9_-]*\.js' "$public_stage/index.html" | head -n 1)
[[ -n "$main_asset" ]] || { echo "frontend index has no main asset" >&2; exit 1; }
main_file="$public_stage${main_asset}"
[[ -s "$main_file" ]] || { echo "frontend main asset is missing: $main_file" >&2; exit 1; }
cache_busted_asset="${main_asset%.js}-${release_id}.js"
cp -f "$main_file" "$public_stage${cache_busted_asset}"
sed -i "s|${main_asset}|${cache_busted_asset}|g" "$public_stage/index.html"

public_previous_dir="$APP_DIR/data/.public-before-$release_id-$timestamp"
mv "$public_dir" "$public_previous_dir"
mv "$public_stage" "$public_dir"
public_swapped=1

sed -i "s|^SUB2API_IMAGE=.*|SUB2API_IMAGE=$RELEASE_TAG|" "$env_path"
env_changed=1
(cd "$APP_DIR" && docker compose up -d --no-deps --force-recreate "$APP_SERVICE")

for _ in $(seq 1 45); do
  status=$(docker inspect -f '{{if .State.Health}}{{.State.Health.Status}}{{else}}{{.State.Status}}{{end}}' "$APP_CONTAINER")
  if [[ "$status" == "healthy" ]] && curl -fsS --max-time 8 "$HEALTH_URL" >/dev/null; then
    released=1
    rm -rf "$public_previous_dir"
    echo "release complete: image=$RELEASE_TAG"
    exit 0
  fi
  sleep 2
done

docker logs --tail 120 "$APP_CONTAINER" >&2 || true
echo "application did not become healthy" >&2
rollback
trap - ERR
exit 1
