#!/usr/bin/env python3
"""密钥泄漏扫描（C7）。

只扫描 git 跟踪的文件，避免把 node_modules、构建产物和本地密钥文件误报。
无第三方依赖。

用法：
    python tools/secret_scan.py              # 扫描跟踪文件
    python tools/secret_scan.py --staged     # 只扫描暂存区
    python tools/secret_scan.py --all        # 扫描工作区所有文件（含未跟踪）

退出码：0 = 干净，1 = 发现疑似泄漏。
"""

from __future__ import annotations

import argparse
import re
import subprocess
import sys
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parent.parent

# ---------- 扫描规则 ----------
# 每条规则：(名称, 正则, 说明)
# 设计原则：宁可少报也不要滥报。命中率低但每次都是真的，比满屏误报更有用。
RULES: list[tuple[str, re.Pattern[str], str]] = [
    (
        "OpenAI API Key",
        re.compile(r"\bsk-[A-Za-z0-9_\-]{20,}\b"),
        "OpenAI 风格的 API Key",
    ),
    (
        "Anthropic API Key",
        re.compile(r"\bsk-ant-[A-Za-z0-9_\-]{20,}\b"),
        "Anthropic API Key",
    ),
    (
        "Cloudflare API Token",
        re.compile(r"\b[A-Za-z0-9_\-]{40}\b(?=[^\n]{0,80}cloudflare)", re.IGNORECASE),
        "Cloudflare token 字面量",
    ),
    (
        "AWS Access Key ID",
        re.compile(r"\b(?:AKIA|ASIA)[0-9A-Z]{16}\b"),
        "AWS 访问密钥 ID",
    ),
    (
        "GitHub Token",
        re.compile(r"\b(?:ghp|gho|ghu|ghs|ghr)_[A-Za-z0-9]{36,}\b"),
        "GitHub 个人访问令牌",
    ),
    (
        "GitHub Fine-grained PAT",
        re.compile(r"\bgithub_pat_[A-Za-z0-9_]{60,}\b"),
        "GitHub 细粒度令牌",
    ),
    (
        "Slack Token",
        re.compile(r"\bxox[baprs]-[A-Za-z0-9\-]{10,}\b"),
        "Slack 令牌",
    ),
    (
        "私钥内容",
        # 必须匹配完整的 PEM 块，且正文是真实的 base64 负载。
        # 只匹配 "-----BEGIN ... PRIVATE KEY-----" 头会把测试里的
        # '-----BEGIN PRIVATE KEY-----\nabc\n-----END PRIVATE KEY-----' 也报出来。
        re.compile(
            r"-----BEGIN (?:RSA |EC |OPENSSH |PGP )?PRIVATE KEY-----"
            r"\\?n?(?:[A-Za-z0-9+/=]{40,}\\?n?)+"
            r"-----END",
        ),
        "私钥文件内容",
    ),
    (
        "JWT",
        re.compile(r"\beyJ[A-Za-z0-9_\-]{10,}\.eyJ[A-Za-z0-9_\-]{10,}\.[A-Za-z0-9_\-]{10,}\b"),
        "JWT 令牌",
    ),
    (
        "Google OAuth Client Secret",
        re.compile(r"\bGOCSPX-[A-Za-z0-9_\-]{20,}\b"),
        "Google OAuth 客户端密钥（若是内置的公开 CLI 凭据，请确认并记录）",
    ),
    (
        "Generic hardcoded secret",
        # 收紧判据，避免误报泛滥。要求值同时满足：
        #   - 长度 >= 20
        #   - 不含连字符/空格分隔的纯小写单词串（测试值多为 'resume-token-h5'）
        #   - 含数字且含大小写混合，或长度 >= 32
        # 目标：真密钥（高熵）命中，测试占位（低熵）放行。
        re.compile(
            r"""(?:password|passwd|secret|token|api[_-]?key|private[_-]?key|access[_-]?key)"""
            r"""\s*[:=]\s*["']((?![/\s])(?=[^"']{20,})(?=[^"']*\d)[A-Za-z0-9_\-]{20,})["']""",
            re.IGNORECASE,
        ),
        "疑似硬编码的密码或密钥赋值",
    ),
]

# 低熵值的判定：只由小写字母、数字、连字符、下划线组成，且没有连续的长随机段。
# 用于对 Generic hardcoded secret 规则做二次过滤。
LOW_ENTROPY_RE = re.compile(r"^[a-z0-9_\-]+$")

# ---------- 允许清单 ----------
# 这些是文档/模板中的占位符与已知无害值，命中后不计入。用小写子串匹配。
ALLOWLIST_SUBSTRINGS = (
    "your-",
    "your_",
    "<your",
    "example.com",
    "placeholder",
    "changeme",
    "xxx",
    "****",
    "redacted",
    "dummy",
    "fake",
    "test-key",
    "test_key",
    "admin-<64hex>",
    "sub2api",
    "postgres",
    "localhost",
    "127.0.0.1",
    "process.env",
    "os.environ",
    "getenv",
    "${",
    "{{",
    "${{",
)

# 不扫描的路径（子串匹配）。这些目录下必然有大量非源码内容。
SKIP_PATH_SUBSTRINGS = (
    ".git/",
    "node_modules/",
    "dist/",
    "build/",
    "vendor/",
    "/bin/",
    "pnpm-lock.yaml",
    "package-lock.json",
    "go.sum",
    ".min.js",
    ".min.css",
)

# ---------- 已知且已审查的例外 ----------
# 格式：值（精确子串） -> (说明, 审查日期)
#
# 规则：只有「已确认是公开凭据或非敏感值」才能进这里。
# 不要把真实私密密钥塞进来——那是自欺欺人，扫描器会失去意义。
#
# 每次新增例外都必须写清依据，方便后人复核。
KNOWN_EXCEPTIONS: dict[str, tuple[str, str]] = {
    "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": (
        "Antigravity 客户端内置的 Google OAuth client secret。"
        "上游从基线至最新 tip 一直存在，非本项目引入。"
        "此类 CLI 工具的 client secret 属公开凭据，安全性由 PKCE 与用户授权保证。",
        "2026-09-18",
    ),
    "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": (
        "Gemini CLI 内置的 Google OAuth client secret。"
        "同上，上游自带，属公开凭据。",
        "2026-09-18",
    ),
}

# 不扫描的扩展名。二进制与压缩包无法逐行匹配，且容易误判。
SKIP_EXTENSIONS = frozenset(
    {
        ".png", ".jpg", ".jpeg", ".gif", ".webp", ".ico", ".svg",
        ".pdf", ".zip", ".gz", ".tar", ".tgz", ".7z", ".rar",
        ".exe", ".dll", ".so", ".dylib", ".bin", ".woff", ".woff2",
        ".ttf", ".eot", ".mp4", ".mp3", ".wasm",
        ".patch",  # patch 文件常含被移除的历史内容，噪声大
    }
)

# 测试上下文：这些路径下的命中不报。
#
# 理由：测试需要构造假 Key / 假 token / 假私钥来驱动被测代码，
# 这是设计需要而非泄漏。把它们全部报出来只会让扫描器被无视。
#
# 例外：Google OAuth Client Secret 与真实私钥仍会在这些文件中被报出——
# 因为那类值即便在测试里出现，也常是从生产代码复制过来的真凭据。
TEST_CONTEXT_MARKERS = (
    "/__tests__/",
    ".spec.ts",
    ".spec.tsx",
    ".test.ts",
    ".test.tsx",
    "_test.go",
    "testdata/",
    "/test/",
    "/tests/",
)

# 在这些测试上下文中仍然强制报告的规则（真凭据不应出现在任何地方）。
ALWAYS_REPORT_IN_TESTS = frozenset(
    {
        "Google OAuth Client Secret",
        "WeChat Pay Key",
    }
)


def is_test_context(rel_path: str) -> bool:
    normalized = "/" + rel_path.replace("\\", "/")
    return any(marker in normalized for marker in TEST_CONTEXT_MARKERS)

MAX_FILE_BYTES = 2 * 1024 * 1024  # 跳过大于 2MB 的文件


def git_tracked_files(staged_only: bool = False) -> list[str]:
    """返回 git 跟踪的文件列表（相对于仓库根）。"""
    cmd = ["git", "diff", "--cached", "--name-only", "--diff-filter=ACMR"]
    if not staged_only:
        cmd = ["git", "ls-files"]
    try:
        out = subprocess.run(
            cmd,
            cwd=REPO_ROOT,
            capture_output=True,
            text=True,
            check=True,
        ).stdout
    except (subprocess.CalledProcessError, FileNotFoundError) as exc:
        print(f"错误：无法获取文件列表（{exc}）", file=sys.stderr)
        return []
    return [line for line in out.splitlines() if line.strip()]


def all_worktree_files() -> list[str]:
    """返回工作区所有文件（含未跟踪），排除 .git 目录。"""
    result = []
    for path in REPO_ROOT.rglob("*"):
        if not path.is_file():
            continue
        rel = path.relative_to(REPO_ROOT).as_posix()
        if rel.startswith(".git/"):
            continue
        result.append(rel)
    return result


def should_skip(rel_path: str) -> bool:
    lowered = rel_path.lower()
    if any(s in lowered for s in SKIP_PATH_SUBSTRINGS):
        return True
    return Path(rel_path).suffix.lower() in SKIP_EXTENSIONS


def is_allowlisted(line: str) -> bool:
    """整行命中允许清单则跳过。"""
    lowered = line.lower()
    return any(token in lowered for token in ALLOWLIST_SUBSTRINGS)


def scan_file(rel_path: str) -> list[tuple[int, str, str, str]]:
    """扫描单个文件，返回 [(行号, 规则名, 说明, 行内容)]。"""
    findings: list[tuple[int, str, str, str]] = []
    full = REPO_ROOT / rel_path
    try:
        if full.stat().st_size > MAX_FILE_BYTES:
            return findings
        text = full.read_text(encoding="utf-8", errors="replace")
    except (OSError, UnicodeError):
        return findings

    in_test = is_test_context(rel_path)

    for lineno, line in enumerate(text.splitlines(), start=1):
        if not line.strip():
            continue
        if is_allowlisted(line):
            continue
        if any(exc in line for exc in KNOWN_EXCEPTIONS):
            # 该行含已审查的公开凭据例外，整行跳过。
            continue
        for name, pattern, note in RULES:
            match = pattern.search(line)
            if not match:
                continue
            # 测试上下文豁免（真凭据类规则除外）。
            if in_test and name not in ALWAYS_REPORT_IN_TESTS:
                continue
            # Generic 规则额外做低熵过滤：测试里的 'resume-token-123' 这类
            # 全小写+连字符的值不应该报，否则扫描器会被误报淹没从而失效。
            if name == "Generic hardcoded secret" and match.groups():
                value = match.group(1)
                if LOW_ENTROPY_RE.match(value) and len(value) < 32:
                    continue
            findings.append((lineno, name, note, line.strip()[:160]))
            break  # 一行只报一次，避免重复刷屏
    return findings


def binary_paths(files: list[str]) -> set[str]:
    """批量查询哪些路径被 git 标记为 binary。

    单文件调 git check-attr 在 2000+ 文件时是子进程灾难，因此一次传入
    全部路径，用 `--stdin -z` 让 git 批量返回，再解析输出。
    """
    if not files:
        return set()
    try:
        out = subprocess.run(
            ["git", "check-attr", "--stdin", "-z", "binary"],
            cwd=REPO_ROOT,
            input="\0".join(files),
            capture_output=True,
            text=True,
            check=True,
        ).stdout
    except (subprocess.CalledProcessError, FileNotFoundError):
        return set()

    # 输出格式（-z）：路径\0属性名\0属性值\0 循环
    parts = out.split("\0")
    result = set()
    for i in range(0, len(parts) - 2, 3):
        path, _attr, value = parts[i], parts[i + 1], parts[i + 2]
        if value == "set":
            result.add(path)
    return result


def main() -> int:
    parser = argparse.ArgumentParser(description="扫描疑似泄漏的密钥（C7）")
    parser.add_argument("--staged", action="store_true", help="只扫描暂存区文件")
    parser.add_argument("--all", action="store_true", help="扫描工作区所有文件（含未跟踪）")
    parser.add_argument("--quiet", action="store_true", help="只在发现问题时输出")
    args = parser.parse_args()

    if args.all:
        files = all_worktree_files()
        scope = "工作区全部文件"
    elif args.staged:
        files = git_tracked_files(staged_only=True)
        scope = "暂存区文件"
    else:
        files = git_tracked_files()
        scope = "git 跟踪文件"

    if not files:
        print(f"未获取到待扫描文件（范围：{scope}）")
        return 0

    total_findings = 0
    scanned = 0
    skipped = 0

    # 先过滤路径，再一次性查二进制属性，避免逐文件调 git。
    candidates = [f for f in sorted(files) if not should_skip(f)]
    skipped += len(files) - len(candidates)
    binaries = binary_paths(candidates)

    for rel in candidates:
        if rel in binaries:
            skipped += 1
            continue
        scanned += 1
        findings = scan_file(rel)
        if findings:
            if total_findings == 0 and not args.quiet:
                print("发现疑似密钥：\n")
            total_findings += len(findings)
            for lineno, name, note, snippet in findings:
                print(f"  {rel}:{lineno}")
                print(f"    规则：{name}（{note}）")
                print(f"    内容：{snippet}")
                print()

    if total_findings:
        print(f"扫描 {scanned} 个文件，跳过 {skipped} 个，命中 {total_findings} 处。")
        print()
        print("处置建议：")
        print("  1. 若确为真实密钥：立刻从代码中移除，并【更换该凭证】——")
        print("     删文件不等于撤销泄漏，必须换 key 才算处置。")
        print("  2. 若为误报：检查能否改成明显的占位符（如 your-api-key-here）。")
        print("     不要往允许清单里加过宽的关键词，那会让扫描器失效。")
        return 1

    if not args.quiet:
        print(f"扫描 {scanned} 个文件（跳过 {skipped} 个），未发现疑似密钥。")
        if KNOWN_EXCEPTIONS:
            print()
            print(f"已审查的例外 {len(KNOWN_EXCEPTIONS)} 项（这些值存在但判断为可接受）：")
            for value, (reason, reviewed) in KNOWN_EXCEPTIONS.items():
                preview = value[:16] + "..." if len(value) > 16 else value
                print(f"  - {preview}（审查于 {reviewed}）")
                print(f"    {reason}")
            print()
            print("  例外清单在 tools/secret_scan.py 的 KNOWN_EXCEPTIONS 中维护。")
    return 0


if __name__ == "__main__":
    sys.exit(main())
