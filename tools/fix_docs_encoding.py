#!/usr/bin/env python3
"""修复 docs/ 下二次编码损坏的 Markdown 文档。

损坏形态：UTF-8 字节被按 GB18030 误解码后，又以 UTF-8 保存。
修复方式：逆向还原 —— utf-8 解码 -> gb18030 编码 -> utf-8 解码。

不修改任何非损坏文件。输出修复质量报告，不掩盖损坏率。
"""
import codecs
import sys
from pathlib import Path

TARGETS = [
    "docs/ITERATION_OS_CN.md",
    "docs/TEAM_QUICKSTART_CN.md",
    "docs/BACKUP_RESTORE_DRILL_CN.md",
    "docs/R2_BACKUP_SETUP_CN.md",
    "docs/RELEASE_RUNBOOK_CN.md",
    "docs/RELEASE_CANDIDATE_20260525_CN.md",
    "docs/ADMIN_PAYMENT_INTEGRATION_API.md",
]

REPO = Path(__file__).resolve().parent.parent


def looks_like_mojibake(text: str) -> bool:
    """乱码特征判定。

    核心观察（实测数据）：UTF-8 被按 GB18030 误解码后，中文全角标点会整体消失。

        ITERATION_OS_CN.md          cjkMark=  1   乱码
        TEAM_QUICKSTART_CN.md       cjkMark=  1   乱码
        RELEASE_CANDIDATE_...md     cjkMark=  2   乱码
        R2_BACKUP_SETUP_CN.md       cjkMark=  3   乱码
        BACKUP_RESTORE_DRILL_CN.md  cjkMark= 37   正常
        RELEASE_RUNBOOK_CN.md       cjkMark= 59   正常
        PRELAUNCH_GATE_CN.md        cjkMark=111   正常
        PAYMENT_CN.md               cjkMark=155   正常

    分界清晰：乱码文件的中文标点数 < 5，正常文件 >= 29。
    取阈值 10 作为保守分界。不使用 ASCII 标点计数（代码块会干扰）。
    """
    if not text:
        return False
    sample = text[:8000]
    cjk = sum(1 for ch in sample if "\u4e00" <= ch <= "\u9fff")
    if cjk < 50:
        return False
    common_marks = sum(sample.count(m) for m in "，。、；：（）")
    return common_marks < 10


def restore(raw: bytes) -> tuple[str, float, int]:
    """返回 (修复后文本, 损坏率百分比, 损坏字符数)。"""
    body = raw[3:] if raw.startswith(codecs.BOM_UTF8) else raw
    mojibake = body.decode("utf-8")
    restored_bytes = mojibake.encode("gb18030", errors="replace")
    repaired = restored_bytes.decode("utf-8", errors="replace")
    bad = repaired.count("\ufffd")
    total = max(len(repaired), 1)
    return repaired, bad / total * 100, bad


def main() -> int:
    failures = []
    report = []

    for rel in TARGETS:
        path = REPO / rel
        if not path.is_file():
            failures.append((rel, "文件不存在"))
            continue

        raw = path.read_bytes()
        try:
            decoded = (raw[3:] if raw.startswith(codecs.BOM_UTF8) else raw).decode("utf-8")
        except UnicodeDecodeError as exc:
            failures.append((rel, f"非 UTF-8，跳过：{exc}"))
            continue

        if not looks_like_mojibake(decoded):
            report.append((rel, "已是正常文本，跳过", 0.0, 0))
            continue

        repaired, ratio, bad = restore(raw)

        if ratio > 15.0:
            failures.append((rel, f"损坏率过高 {ratio:.1f}%，拒绝写回"))
            continue

        path.write_bytes(codecs.BOM_UTF8 + repaired.encode("utf-8"))
        report.append((rel, "已修复", ratio, bad))

    print("=== 修复报告 ===")
    for rel, status, ratio, bad in report:
        print(f"{status:16} 损坏率={ratio:5.2f}%  坏字符={bad:4}  {rel}")

    if failures:
        print("\n=== 未处理 ===")
        for rel, why in failures:
            print(f"  {rel}: {why}")
        return 1

    print(f"\n完成：{sum(1 for _, s, _, _ in report if s == '已修复')} 个文件已修复")
    return 0


if __name__ == "__main__":
    sys.exit(main())
