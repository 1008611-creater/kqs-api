#!/usr/bin/env bash
set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)
deploy_script="$repo_root/ops/production/deploy-remote.sh"
tmp_dir=$(mktemp -d)
fake_bin="$tmp_dir/bin"
mkdir -p "$fake_bin"

cleanup() {
  rm -rf "$tmp_dir"
}
trap cleanup EXIT

cat >"$fake_bin/docker" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
case "${1:-}" in
  build)
    cat >/dev/null
    ;;
  run|logs|compose)
    ;;
  inspect)
    printf '%s\n' "${FAKE_DOCKER_HEALTH:-healthy}"
    ;;
  *)
    echo "unexpected fake docker command: $*" >&2
    exit 1
    ;;
esac
EOF

cat >"$fake_bin/curl" <<'EOF'
#!/usr/bin/env bash
exit "${FAKE_CURL_EXIT:-0}"
EOF

cat >"$fake_bin/sleep" <<'EOF'
#!/usr/bin/env bash
exit 0
EOF
chmod 700 "$fake_bin/docker" "$fake_bin/curl" "$fake_bin/sleep"

create_release() {
  local target=$1
  mkdir -p "$target/frontend/assets"
  printf '<html><script src="/assets/index-test.js"></script></html>\n' >"$target/frontend/index.html"
  printf 'console.log("release");\n' >"$target/frontend/assets/index-test.js"
  tar -C "$target/frontend" -czf "$target/frontend-dist.tar.gz" .
  printf '<html>embedded frontend marker</html>\n' >"$target/sub2api"
  sha256sum "$target/sub2api" | awk '{print $1}'
}

run_success_case() {
  local app="$tmp_dir/app-success"
  local release="$tmp_dir/release-success"
  local checksum
  mkdir -p "$app/data/public" "$app/backups/releases"
  printf 'SUB2API_IMAGE=sub2api:before-success\n' >"$app/.env"
  printf 'old public content\n' >"$app/data/public/old-marker"
  checksum=$(create_release "$release")

  PATH="$fake_bin:$PATH" APP_DIR="$app" APP_CONTAINER="fake-sub2api" \
    "$deploy_script" "$release" "sub2api:release-success" "$checksum"

  grep -qx 'SUB2API_IMAGE=sub2api:release-success' "$app/.env"
  test -f "$app/data/public/assets/index-test-release-success.js"
  test ! -e "$app/data/public/old-marker"
  test ! -d "$app/data/.public-stage-release-success"
  printf '%s\n' "success deployment path passed"
}

run_rollback_case() {
  local app="$tmp_dir/app-rollback"
  local release="$tmp_dir/release-rollback"
  local checksum
  mkdir -p "$app/data/public" "$app/backups/releases"
  printf 'SUB2API_IMAGE=sub2api:before-rollback\n' >"$app/.env"
  printf 'old public content\n' >"$app/data/public/old-marker"
  checksum=$(create_release "$release")

  if PATH="$fake_bin:$PATH" APP_DIR="$app" APP_CONTAINER="fake-sub2api" \
    FAKE_DOCKER_HEALTH="starting" "$deploy_script" "$release" "sub2api:release-rollback" "$checksum"; then
    echo "expected unhealthy deployment to fail" >&2
    exit 1
  fi

  grep -qx 'SUB2API_IMAGE=sub2api:before-rollback' "$app/.env"
  test -f "$app/data/public/old-marker"
  test ! -e "$app/data/public/assets/index-test-release-rollback.js"
  printf '%s\n' "unhealthy deployment rollback path passed"
}

run_success_case
run_rollback_case
