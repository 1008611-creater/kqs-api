# docs/ 索引

> 本文件是本目录所有文档的导航。**新增文档后必须在这里登记。**
>
> 顶层入口是仓库根目录的 `AGENTS.md`。

---

## 一、先读这几份

| 文档 | 讲什么 | 什么时候读 |
| --- | --- | --- |
| [`../AGENTS.md`](../AGENTS.md) | 项目定位、上游关系、变更分级、验证方式 | **任何时候，从这里开始** |
| [`CONSTRAINTS.md`](CONSTRAINTS.md) | 13 条硬约束，每条标注检查手段 | 动手前 |
| [`ITERATION_OS_CN.md`](ITERATION_OS_CN.md) | 项目地图（6 条链路）、迭代节奏、发布闸门 | 想理解项目全貌时 |
| [`../DEV_GUIDE.md`](../DEV_GUIDE.md) | 环境配置 + 11 个常见坑点 | 本地环境出问题时 |

---

## 二、发布与上线

按实际操作顺序排列：

| 文档 | 讲什么 | 状态 |
| --- | --- | --- |
| [`PRELAUNCH_GATE_CN.md`](PRELAUNCH_GATE_CN.md) | 上线前检查清单 | 有效 |
| [`RELEASE_PIPELINE_CN.md`](RELEASE_PIPELINE_CN.md) | 发布流水线概览 | 有效 |
| [`RELEASE_RUNBOOK_CN.md`](RELEASE_RUNBOOK_CN.md) | **发布与回滚的完整步骤**（含容器内执行测试的命令） | 有效 |
| [`RELEASE_CANDIDATE_20260525_CN.md`](RELEASE_CANDIDATE_20260525_CN.md) | 2026-05-25 候选版本的验证证据与未放行闸门 | 历史快照 |
| [`SECRETS_INVENTORY_TEMPLATE.md`](SECRETS_INVENTORY_TEMPLATE.md) | 密钥清单模板（填写后**不要提交**） | 模板 |

---

## 三、备份与恢复

| 文档 | 讲什么 |
| --- | --- |
| [`R2_BACKUP_SETUP_CN.md`](R2_BACKUP_SETUP_CN.md) | Cloudflare R2 备份的启用与配置步骤 |
| [`BACKUP_RESTORE_DRILL_CN.md`](BACKUP_RESTORE_DRILL_CN.md) | 恢复演练的环境、步骤与验收清单 |

**当前状态**：R2 bucket `sub2api-backups` 已连接，定时策略 `30 2 * * *` 已开启。**恢复演练尚未完成**——这是发布闸门中未放行的一项。

---

## 四、业务与运营

| 文档 | 讲什么 |
| --- | --- |
| [`TEAM_QUICKSTART_CN.md`](TEAM_QUICKSTART_CN.md) | 团队上手步骤、卡密套餐、Codex 客户端配置 |
| [`OPERATIONS_HANDOFF_CN.md`](OPERATIONS_HANDOFF_CN.md) | 运维交接说明 |
| [`PAYMENT_CN.md`](PAYMENT_CN.md) / [`PAYMENT.md`](PAYMENT.md) | 支付系统说明（中/英） |
| [`ADMIN_PAYMENT_INTEGRATION_API.md`](ADMIN_PAYMENT_INTEGRATION_API.md) | 管理员支付接入 API |

---

## 五、变更记录

上游采用 `openspec/changes/<name>/` 的结构记录每次变更，本项目沿用：

```
openspec/changes/<变更名>/
  proposal.md      ← Why / What Changes / Capabilities / Impact
  design.md        ← 技术方案
  tasks.md         ← 可勾选的任务清单
  verification.md  ← 验证证据与验证边界
  specs/<能力>/spec.md  ← 要求与场景（WHEN/THEN）
```

已有样例可参照：`openspec/changes/codex-manifest-pinned-accounts/`。

**关键点**：`verification.md` 必须写「**验证边界**」——明确说明哪些没验证。这不是减分项，是防止误判的必要信息。

---

## 六、上游对齐参考

| 路径 | 说明 |
| --- | --- |
| 基线 ref | `kqs/upstream-base` → `881f3202694c6bc932446931a30c27d9675178b9`（2026-09-15） |
| 上游仓库 | `Wei-Shaw/sub2api` |
| `openspec/` | 从上游引入的变更记录体系 |

**注意**：本仓库与上游**无共同祖先**，同步只能增量摘取。详见 `AGENTS.md` 第 1 节。

---

## 七、文档命名约定

| 后缀 | 含义 |
| --- | --- |
| `*_CN.md` | 本项目独有文档。中文为主。 |
| `*.md`（无后缀） | 与上游同名或对齐上游的文档。改动时需注意同步冲突。 |
| `*_TEMPLATE.md` | 模板，填写后另存，不要直接改模板。 |

**上游 `docs/` 与本地 `docs/` 的差异**：上游有 `ASYNC_IMAGE_TASKS.md`、`BATCH_IMAGE_MVP.md`、`COMPOSITE_GROUPS.md`、`PLUGIN_DEVELOPMENT.md`、`channel-monitor-v2-safe-defaults.md`、`legal/admin-compliance.{en,zh}.md` —— 本项目**均未引入**。
