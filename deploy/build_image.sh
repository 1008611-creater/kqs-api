#!/usr/bin/env bash
# 本地构建镜像的快速脚本。正式发布时请传入固定 tag，避免 latest 漂移。

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"

TAG="${1:-${SUB2API_IMAGE:-}}"
if [ -z "${TAG}" ]; then
    SHORT_SHA="$(git -C "${REPO_ROOT}" rev-parse --short HEAD 2>/dev/null || echo local)"
    TAG="sub2api:kqs-api-${SHORT_SHA}"
fi

docker build -t "${TAG}" \
    --build-arg GOPROXY=https://goproxy.cn,direct \
    --build-arg GOSUMDB=sum.golang.google.cn \
    --build-arg NPM_CONFIG_REGISTRY="${NPM_CONFIG_REGISTRY:-https://registry.npmmirror.com}" \
    --build-arg PNPM_VERSION="${PNPM_VERSION:-9.15.9}" \
    -f "${REPO_ROOT}/deploy/Dockerfile" \
    "${REPO_ROOT}"

echo "Built image: ${TAG}"
