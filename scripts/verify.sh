#!/usr/bin/env bash
# scripts/verify.sh —— 项目统一质量门
#
# 设计原则：分层降级 + 诚实汇总。
#   L0 秒级静态检查（永远跑）
#   L1 前端检查（需要 node_modules）
#   L2 后端检查（需要可用 Go 工具链）
#
# 跑不了的层会标记为「未验证」并打印醒目警告，绝不假装通过。
# 最终汇总表区分：通过 / 未验证 / 失败。
#
# 用法：
#   bash scripts/verify.sh           # 跑全部能跑的
#   bash scripts/verify.sh --quick   # 只跑 L0

set -uo pipefail

# ---------- 定位仓库根 ----------
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
cd "${REPO_ROOT}"

QUICK=0
[ "${1:-}" = "--quick" ] && QUICK=1

# ---------- 输出色彩 ----------
if [ -t 1 ]; then
  C_RED=$'\033[31m'; C_GRN=$'\033[32m'; C_YEL=$'\033[33m'
  C_CYA=$'\033[36m'; C_BLD=$'\033[1m'; C_RST=$'\033[0m'
else
  C_RED=''; C_GRN=''; C_YEL=''; C_CYA=''; C_BLD=''; C_RST=''
fi

# ---------- 结果收集 ----------
# 状态取值：PASS / SKIP / FAIL
RESULT_NAMES=()
RESULT_STATES=()
RESULT_NOTES=()

record() {  # record <名称> <状态> <备注>
  RESULT_NAMES+=("$1")
  RESULT_STATES+=("$2")
  RESULT_NOTES+=("$3")
}

section() {
  printf '\n%s%s %s%s\n' "${C_CYA}" "${C_BLD}" "$1" "${C_RST}"
}

# ---------- L0：静态检查 ----------

l0_gofmt() {
  section "L0-1 · gofmt 格式检查（仅本次改动的文件）"

  # 注意：仓库里有 13 个上游遗留的未格式化文件。把它们全部列为门禁失败
  # 会让这个检查永远红着，从而被无视——门禁失效比没有门禁更糟。
  # 因此只检查「本次工作区相对 HEAD 有改动的」Go 文件。
  if ! command -v gofmt >/dev/null 2>&1; then
    printf '  %s⚠ gofmt 不在 PATH%s\n' "${C_YEL}" "${C_RST}"
    record "gofmt" "SKIP" "gofmt 不可用"
    return
  fi

  if ! git rev-parse HEAD >/dev/null 2>&1; then
    printf '  %s⚠ 非 git 仓库，跳过%s\n' "${C_YEL}" "${C_RST}"
    record "gofmt" "SKIP" "非 git 仓库"
    return
  fi

  local changed
  changed="$(git status --porcelain -- '*.go' 2>/dev/null | awk '{print $NF}' | grep -E '\.go$' || true)"

  if [ -z "${changed}" ]; then
    printf '  %s✓ 本次无 Go 文件改动%s\n' "${C_GRN}" "${C_RST}"
    record "gofmt" "PASS" "无改动"
    return
  fi

  local bad="" f
  while IFS= read -r f; do
    [ -z "${f}" ] && continue
    [ -f "${f}" ] || continue
    if [ -n "$(gofmt -l "${f}" 2>/dev/null)" ]; then
      bad="${bad}${f}"$'\n'
    fi
  done <<< "${changed}"

  if [ -z "${bad}" ]; then
    printf '  %s✓ 本次改动的 %d 个 Go 文件格式正确%s\n' "${C_GRN}" "$(printf '%s' "${changed}" | grep -c .)" "${C_RST}"
    record "gofmt" "PASS" ""
  else
    printf '  %s✗ 本次改动中有文件格式不正确：%s\n' "${C_RED}" "${C_RST}"
    printf '%s' "${bad}" | sed '/^$/d; s/^/      /'
    printf '      修复：gofmt -w <文件>\n'
    record "gofmt" "FAIL" "有文件未格式化"
  fi
}

l0_migration_immutable() {
  section "L0-2 · 已提交迁移文件改动检查（C5）"

  # C5 的正确判据：已提交进本仓库历史的迁移文件，不得在本次工作区里被修改。
  # 不能用「对比上游」——本仓库与上游无共同祖先，文件名体系不同，
  # 那样会把两边同编号不同名的文件全部误报成「被改动」。
  #
  # 因此只看工作区状态：相对 HEAD 的 M（已跟踪且被修改）与 D（已跟踪且被删除）。
  if ! git rev-parse HEAD >/dev/null 2>&1; then
    printf '  %s⚠ 非 git 仓库，跳过%s\n' "${C_YEL}" "${C_RST}"
    record "已提交迁移改动" "SKIP" "非 git 仓库"
    return
  fi

  local status_out
  status_out="$(git status --porcelain -- backend/migrations/ 2>/dev/null || true)"

  # 只关心已跟踪文件的修改/删除（第一列是 M/D/R），忽略 ?? 未跟踪（=新增迁移，正常）
  local problems
  problems="$(printf '%s\n' "${status_out}" | grep -E '^[ ]?[MDR]' || true)"

  if [ -z "${problems}" ]; then
    local added
    added="$(printf '%s\n' "${status_out}" | grep -cE '^\?\?' || true)"
    printf '  %s✓ 未修改任何已提交的迁移文件（工作区新增 %s 个）%s\n' "${C_GRN}" "${added}" "${C_RST}"
    record "已提交迁移改动" "PASS" ""
  else
    printf '  %s✗ 检测到已提交的迁移文件被修改或删除（C5 违规）：%s\n' "${C_RED}" "${C_RST}"
    printf '%s\n' "${problems}" | sed 's/^/      /'
    printf '      %s提示：已执行的迁移内容不得修改。要修正请新增迁移文件。%s\n' "${C_YEL}" "${C_RST}"
    record "已提交迁移改动" "FAIL" "历史迁移被改动"
  fi
}

l0_doc_encoding() {
  section "L0-3 · 文档编码检查（C13）"
  local py=""
  for cand in python3 python; do
    command -v "${cand}" >/dev/null 2>&1 && { py="${cand}"; break; }
  done

  if [ -z "${py}" ]; then
    printf '  %s⚠ Python 不可用，跳过%s\n' "${C_YEL}" "${C_RST}"
    record "文档编码" "SKIP" "Python 不可用"
    return
  fi

  local out
  out="$("${py}" - <<'PYEOF' 2>/dev/null || echo "__ERR__"
import pathlib, sys
bad = []
for p in pathlib.Path("docs").rglob("*.md"):
    raw = p.read_bytes()
    if raw.startswith(b"\xef\xbb\xbf"):
        bad.append(f"{p} 含 BOM")
        continue
    try:
        t = raw.decode("utf-8")
    except UnicodeDecodeError:
        bad.append(f"{p} 非 UTF-8")
        continue
    if "\ufffd" in t:
        bad.append(f"{p} 含 {t.count(chr(0xfffd))} 个替换字符")
for p in pathlib.Path(".").glob("*.md"):
    raw = p.read_bytes()
    if raw.startswith(b"\xef\xbb\xbf"):
        bad.append(f"{p} 含 BOM")
print("\n".join(bad))
PYEOF
)"

  if [ "${out}" = "__ERR__" ]; then
    printf '  %s⚠ 检查脚本执行失败%s\n' "${C_YEL}" "${C_RST}"
    record "文档编码" "SKIP" "检查脚本异常"
  elif [ -z "${out}" ]; then
    printf '  %s✓ 所有 Markdown 文件为无 BOM 的 UTF-8，无乱码%s\n' "${C_GRN}" "${C_RST}"
    record "文档编码" "PASS" ""
  else
    printf '  %s✗ 发现问题：%s\n' "${C_RED}" "${C_RST}"
    printf '%s\n' "${out}" | sed 's/^/      /'
    record "文档编码" "FAIL" "编码或乱码问题"
  fi
}

l0_ascii_quotes() {
  section "L0-4 · 配置引号检查（C12）"
  local py=""
  for cand in python3 python; do
    command -v "${cand}" >/dev/null 2>&1 && { py="${cand}"; break; }
  done

  if [ -z "${py}" ]; then
    printf '  %s⚠ Python 不可用，跳过%s\n' "${C_YEL}" "${C_RST}"
    record "配置引号" "SKIP" "Python 不可用"
    return
  fi

  local out
  out="$("${py}" - <<'PYEOF' 2>/dev/null || echo "__ERR__"
import pathlib
FULLWIDTH = "\u201c\u201d\u2018\u2019"
bad = []
EXCLUDE = {".git", "node_modules", "dist", "vendor"}
for pat in ("**/*.json", "**/*.yaml", "**/*.yml", "**/*.toml"):
    for p in pathlib.Path(".").glob(pat):
        if any(part in EXCLUDE for part in p.parts):
            continue
        try:
            t = p.read_text(encoding="utf-8")
        except Exception:
            continue
        hits = sum(t.count(c) for c in FULLWIDTH)
        if hits:
            bad.append(f"{p} 含 {hits} 个全角引号")
print("\n".join(bad))
PYEOF
)"

  if [ "${out}" = "__ERR__" ]; then
    printf '  %s⚠ 检查脚本执行失败%s\n' "${C_YEL}" "${C_RST}"
    record "配置引号" "SKIP" "检查脚本异常"
  elif [ -z "${out}" ]; then
    printf '  %s✓ 配置文件中未发现全角引号%s\n' "${C_GRN}" "${C_RST}"
    record "配置引号" "PASS" ""
  else
    printf '  %s✗ 发现问题：%s\n' "${C_RED}" "${C_RST}"
    printf '%s\n' "${out}" | sed 's/^/      /'
    record "配置引号" "FAIL" "含全角引号"
  fi
}

l0_secret_scan() {
  section "L0-5 · 密钥扫描（C7）"
  if [ ! -f "tools/secret_scan.py" ]; then
    printf '  %s⚠ tools/secret_scan.py 不存在，跳过%s\n' "${C_YEL}" "${C_RST}"
    record "密钥扫描" "SKIP" "扫描脚本缺失"
    return
  fi

  local py=""
  for cand in python3 python; do
    command -v "${cand}" >/dev/null 2>&1 && { py="${cand}"; break; }
  done

  if [ -z "${py}" ]; then
    printf '  %s⚠ Python 不可用，跳过%s\n' "${C_YEL}" "${C_RST}"
    record "密钥扫描" "SKIP" "Python 不可用"
    return
  fi

  if "${py}" tools/secret_scan.py >/tmp/kqs_secret_scan.log 2>&1; then
    printf '  %s✓ 未发现密钥泄漏%s\n' "${C_GRN}" "${C_RST}"
    record "密钥扫描" "PASS" ""
  else
    printf '  %s✗ 疑似发现密钥，详见 /tmp/kqs_secret_scan.log%s\n' "${C_RED}" "${C_RST}"
    tail -20 /tmp/kqs_secret_scan.log | sed 's/^/      /'
    record "密钥扫描" "FAIL" "疑似泄漏"
  fi
}

l0_lock_sync() {
  section "L0-6 · pnpm lock 同步检查（C2）"
  if ! git rev-parse HEAD >/dev/null 2>&1; then
    record "lock 同步" "SKIP" "非 git 仓库"
    return
  fi

  # 工作区未提交改动
  local files
  files="$(git status --porcelain 2>/dev/null | awk '{print $NF}')"
  # 相对基线的改动
  local base_files=""
  if git rev-parse --verify kqs/upstream-base >/dev/null 2>&1; then
    base_files="$(git diff --name-only kqs/upstream-base -- frontend/ 2>/dev/null || true)"
  fi

  local touched_pkg touched_lock
  touched_pkg="$(printf '%s\n%s\n' "${files}" "${base_files}" | grep -E 'frontend/package\.json$' || true)"
  touched_lock="$(printf '%s\n%s\n' "${files}" "${base_files}" | grep -E 'frontend/pnpm-lock\.yaml$' || true)"

  if [ -n "${touched_pkg}" ] && [ -z "${touched_lock}" ]; then
    printf '  %s✗ package.json 有改动但 pnpm-lock.yaml 未同步（C2 违规）%s\n' "${C_RED}" "${C_RST}"
    printf '      修复：cd frontend && pnpm install\n'
    record "lock 同步" "FAIL" "lock 未同步"
  else
    printf '  %s✓ package.json 与 pnpm-lock.yaml 同步%s\n' "${C_GRN}" "${C_RST}"
    record "lock 同步" "PASS" ""
  fi
}

l0_pnpm_modules() {
  section "L0-7 · 包管理器检查（C1）"
  if [ -d "frontend/node_modules/.pnpm" ]; then
    printf '  %s✓ frontend/node_modules 由 pnpm 安装%s\n' "${C_GRN}" "${C_RST}"
    record "包管理器" "PASS" ""
  elif [ -d "frontend/node_modules" ]; then
    printf '  %s✗ frontend/node_modules 存在但非 pnpm 结构（C1 违规）%s\n' "${C_RED}" "${C_RST}"
    printf '      修复：cd frontend && rm -rf node_modules && pnpm install\n'
    record "包管理器" "FAIL" "非 pnpm 结构"
  else
    printf '  %s- node_modules 未安装（跳过，不算失败）%s\n' "${C_CYA}" "${C_RST}"
    record "包管理器" "PASS" "未安装依赖"
  fi
}

# ---------- L1：前端 ----------

detect_node() {
  # 优先用托管 node
  local managed="C:/Users/lsb/.workbuddy-ai/binaries/node/versions/22.22.2-2/node.exe"
  if [ -x "${managed}" ]; then
    echo "${managed}"
    return
  fi
  command -v node >/dev/null 2>&1 && echo "node" || echo ""
}

l1_frontend() {
  section "L1 · 前端检查"
  if [ ! -d "frontend/node_modules" ]; then
    printf '  %s⚠ frontend/node_modules 不存在，无法运行前端检查%s\n' "${C_YEL}" "${C_RST}"
    printf '      安装：cd frontend && pnpm install\n'
    record "前端 typecheck" "SKIP" "依赖未安装"
    record "前端 lint" "SKIP" "依赖未安装"
    return
  fi

  if ! command -v pnpm >/dev/null 2>&1; then
    printf '  %s⚠ pnpm 不在 PATH%s\n' "${C_YEL}" "${C_RST}"
    record "前端 typecheck" "SKIP" "pnpm 不可用"
    record "前端 lint" "SKIP" "pnpm 不可用"
    return
  fi

  printf '  → pnpm run typecheck\n'
  if (cd frontend && pnpm run typecheck >/tmp/kqs_tsc.log 2>&1); then
    printf '  %s✓ typecheck 通过%s\n' "${C_GRN}" "${C_RST}"
    record "前端 typecheck" "PASS" ""
  else
    printf '  %s✗ typecheck 失败%s\n' "${C_RED}" "${C_RST}"
    tail -20 /tmp/kqs_tsc.log | sed 's/^/      /'
    record "前端 typecheck" "FAIL" "见日志"
  fi

  printf '  → pnpm run lint:check\n'
  if (cd frontend && pnpm run lint:check >/tmp/kqs_lint.log 2>&1); then
    printf '  %s✓ lint 通过%s\n' "${C_GRN}" "${C_RST}"
    record "前端 lint" "PASS" ""
  else
    printf '  %s✗ lint 失败%s\n' "${C_RED}" "${C_RST}"
    tail -20 /tmp/kqs_lint.log | sed 's/^/      /'
    record "前端 lint" "FAIL" "见日志"
  fi
}

# ---------- L2：后端 ----------

go_toolchain_healthy() {
  command -v go >/dev/null 2>&1 || return 1

  # 用仓库内的临时目录，避免 Git Bash 下 mktemp 返回 MSYS 路径
  # 与 Windows 原生 go 的路径解析冲突（会拼成 e:\...\C:\Users\... 的畸形路径）。
  local tmp="${REPO_ROOT}/.verify-probe"
  rm -rf "${tmp}" 2>/dev/null
  mkdir -p "${tmp}" 2>/dev/null || return 1

  printf 'package main\n\nimport (\n\t"fmt"\n\t"unicode/utf8"\n)\n\nfunc main() { fmt.Println(utf8.RuneCountInString("x")) }\n' > "${tmp}/main.go"

  local rc=0
  (cd "${tmp}" && GO111MODULE=off go build -o /dev/null . >/dev/null 2>&1) || rc=$?

  rm -rf "${tmp}" 2>/dev/null
  return ${rc}
}

l2_backend() {
  section "L2 · 后端检查"

  if ! command -v go >/dev/null 2>&1; then
    printf '  %s⚠ go 不在 PATH%s\n' "${C_YEL}" "${C_RST}"
    record "后端编译" "SKIP" "go 不可用"
    record "后端单测" "SKIP" "go 不可用"
    record "golangci-lint" "SKIP" "go 不可用"
    return
  fi

  if ! go_toolchain_healthy; then
    printf '  %s%s' "${C_RED}" "${C_BLD}"
    printf '  ┌────────────────────────────────────────────────────────────┐\n'
    printf '  │  ⚠  Go 工具链损坏，后端检查全部跳过                        │\n'
    printf '  └────────────────────────────────────────────────────────────┘\n'
    printf '%s' "${C_RST}"
    printf '      症状：连最小程序都无法编译（std 源文件缺失）\n'
    printf '      %s本次改动在后端侧【完全未验证】，必须依赖 CI。%s\n' "${C_YEL}" "${C_RST}"
    printf '      CI: .github/workflows/backend-ci.yml\n'
    printf '      替代方案：docker run --rm -v "%s/backend:/app" -w /app golang:1.26.6-alpine go test ./...\n' "${REPO_ROOT}"
    record "后端编译" "SKIP" "Go 工具链损坏"
    record "后端单测" "SKIP" "Go 工具链损坏"
    record "golangci-lint" "SKIP" "Go 工具链损坏"
    return
  fi

  printf '  → go build ./...（backend）\n'
  if (cd backend && go build ./... >/tmp/kqs_gobuild.log 2>&1); then
    printf '  %s✓ 编译通过%s\n' "${C_GRN}" "${C_RST}"
    record "后端编译" "PASS" ""
  else
    printf '  %s✗ 编译失败%s\n' "${C_RED}" "${C_RST}"
    tail -20 /tmp/kqs_gobuild.log | sed 's/^/      /'
    record "后端编译" "FAIL" "见日志"
  fi

  printf '  → go test -tags=unit ./...（backend）\n'
  if (cd backend && go test -tags=unit ./... >/tmp/kqs_gotest.log 2>&1); then
    printf '  %s✓ 单元测试通过%s\n' "${C_GRN}" "${C_RST}"
    record "后端单测" "PASS" ""
  else
    printf '  %s✗ 单元测试失败%s\n' "${C_RED}" "${C_RST}"
    tail -20 /tmp/kqs_gotest.log | sed 's/^/      /'
    record "后端单测" "FAIL" "见日志"
  fi

  if command -v golangci-lint >/dev/null 2>&1; then
    printf '  → golangci-lint run ./...（backend）\n'
    if (cd backend && golangci-lint run ./... >/tmp/kqs_gclint.log 2>&1); then
      printf '  %s✓ lint 通过%s\n' "${C_GRN}" "${C_RST}"
      record "golangci-lint" "PASS" ""
    else
      printf '  %s✗ lint 失败%s\n' "${C_RED}" "${C_RST}"
      tail -20 /tmp/kqs_gclint.log | sed 's/^/      /'
      record "golangci-lint" "FAIL" "见日志"
    fi
  else
    printf '  %s- golangci-lint 未安装，跳过%s\n' "${C_CYA}" "${C_RST}"
    record "golangci-lint" "SKIP" "未安装"
  fi
}

# ---------- 主流程 ----------

printf '%s%s\n' "${C_BLD}" "════════════════════════════════════════════════════════════"
printf '  kqs-api 质量门  verify.sh\n'
printf '  仓库：%s\n' "${REPO_ROOT}"
printf '════════════════════════════════════════════════════════════%s\n' "${C_RST}"

l0_gofmt
l0_migration_immutable
l0_doc_encoding
l0_ascii_quotes
l0_secret_scan
l0_lock_sync
l0_pnpm_modules

if [ "${QUICK}" -eq 0 ]; then
  l1_frontend
  l2_backend
else
  printf '\n%s（--quick 模式，跳过 L1/L2）%s\n' "${C_CYA}" "${C_RST}"
  record "前端检查" "SKIP" "--quick"
  record "后端检查" "SKIP" "--quick"
fi

# ---------- 汇总 ----------

printf '\n%s%s\n' "${C_BLD}" "════════════════════════════════════════════════════════════"
printf '  汇总\n'
printf '════════════════════════════════════════════════════════════%s\n' "${C_RST}"

pass=0; skip=0; fail=0
for i in "${!RESULT_NAMES[@]}"; do
  name="${RESULT_NAMES[$i]}"
  state="${RESULT_STATES[$i]}"
  note="${RESULT_NOTES[$i]}"
  case "${state}" in
    PASS) printf '  %s✓ 通过%s    %-22s %s\n' "${C_GRN}" "${C_RST}" "${name}" "${note}"; pass=$((pass+1)) ;;
    SKIP) printf '  %s⚠ 未验证%s  %-22s %s\n' "${C_YEL}" "${C_RST}" "${name}" "${note}"; skip=$((skip+1)) ;;
    FAIL) printf '  %s✗ 失败%s    %-22s %s\n' "${C_RED}" "${C_RST}" "${name}" "${note}"; fail=$((fail+1)) ;;
  esac
done

printf '\n  通过 %d · 未验证 %d · 失败 %d\n\n' "${pass}" "${skip}" "${fail}"

if [ "${fail}" -gt 0 ]; then
  printf '%s%s  结论：存在失败项，不要提交。%s\n\n' "${C_BLD}" "${C_RED}" "${C_RST}"
  exit 1
fi

if [ "${skip}" -gt 0 ]; then
  printf '%s%s  结论：无失败项，但有 %d 项未验证。%s\n' "${C_BLD}" "${C_YEL}" "${skip}" "${C_RST}"
  printf '        未验证 ≠ 通过。涉及这些层的改动需依赖 CI 确认。\n\n'
  exit 0
fi

printf '%s%s  结论：全部通过。%s\n\n' "${C_BLD}" "${C_GRN}" "${C_RST}"
exit 0
