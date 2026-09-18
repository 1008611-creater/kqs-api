# AGENTS.md

> 本文件是**唯一入口**。任何 AI 助手或新加入的工程师，从这里开始。
>
> 本文件只做导航与约束声明，**不复述** `docs/` 下的具体内容。要改流程细节，改对应文档，不要往这里复制。

---

## 1. 这个项目是什么

**矿泉水 API（kqs-api）** —— 一个基于 sub2api 二次开发的 AI API 中转与卡密销售站点。

| 项 | 值 |
| --- | --- |
| 线上域名 | `https://api.cauai.fun` |
| 本仓库 | `1008611-creater/kqs-api` |
| 技术栈 | Go 后端（Ent ORM + Gin）+ Vue3 前端（pnpm）+ PostgreSQL 16 + Redis |
| 部署形态 | Docker 镜像 + Cloudflare Worker 代理 + Cloudflare R2 备份 |
| 当前阶段 | **受控邀请制 Beta**，非公开收费推广（见 `docs/RELEASE_CANDIDATE_20260525_CN.md`） |

### 与上游的关系（重要，先读这段再动手）

本仓库**不是** git fork。它与上游 `Wei-Shaw/sub2api` **没有共同祖先**——历史上是一次导入式快照。

这意味着：

- **不能** `git merge upstream/main`，也不能 `git cherry-pick`（无共同祖先，必然全量冲突）
- 上游同步只能走**增量摘取**：逐个文件比对，手工移植有价值的改动
- 基线 ref：`kqs/upstream-base`（当前指向 `881f3202694c6bc932446931a30c27d9675178b9`，2026-09-15）
- 完整分析见 `docs/ITERATION_OS_CN.md` 与 `.workbuddy-ai/memory/`

**同步上游前必须做的判断**：这个改动是「上游通用修复」还是「上游为新功能重构」？后者往往牵动几十个文件，摘取成本远高于收益。宁可放弃，不要引入半截代码。

---

## 2. 开工前必读

按需读，不要一次全读：

| 你要做的事 | 先读 |
| --- | --- |
| 改任何代码 | 本文件第 3、4 节 + `DEV_GUIDE.md` |
| 理解项目全貌与迭代节奏 | `docs/ITERATION_OS_CN.md` |
| 发布 / 回滚 | `docs/RELEASE_RUNBOOK_CN.md` |
| 备份 / 恢复演练 | `docs/BACKUP_RESTORE_DRILL_CN.md`、`docs/R2_BACKUP_SETUP_CN.md` |
| 上线前自检 | `docs/PRELAUNCH_GATE_CN.md`、`docs/RELEASE_CANDIDATE_20260525_CN.md` |
| 团队上手 / 卡密套餐 | `docs/TEAM_QUICKSTART_CN.md` |
| 支付接入 API | `docs/ADMIN_PAYMENT_INTEGRATION_API.md` |
| 全部文档索引 | `docs/INDEX.md` |
| 硬性红线 | `docs/CONSTRAINTS.md` |

---

## 3. 硬约束（违反即打回，详见 `docs/CONSTRAINTS.md`）

以下几条是**每次改动都要自检**的，不是参考建议：

1. **包管理器是 pnpm，不是 npm。** 改 `frontend/package.json` 后必须同步提交 `pnpm-lock.yaml`，否则 CI 的 `--frozen-lockfile` 直接失败。
2. **改 `backend/ent/schema/*.go` 后必须 `go generate ./ent` 并提交生成代码。**
3. **给 Go interface 加方法后，所有测试 stub 必须补齐**，否则整包编译失败。
4. **数据库迁移的唯一标识是文件名，不是编号。** 编号重复是允许的；但**已执行过的迁移文件内容不能改**（checksum 校验会导致服务启动失败）。
5. **不得提交任何真实密钥。** 提交前跑 secret-scan。
6. **不改 `.git` 历史，不 force push 到 main。**

---

## 4. 变更风险分级（决定你需要做到哪一步）

改动前先给自己分级。**先分级，再动手**——分级决定了验证成本和是否需要人工闸门。

| 级别 | 典型改动 | 必须做的验证 | 是否需人工批准 |
| --- | --- | --- | --- |
| **L0 文档/注释** | 改 `.md`、代码注释、i18n 文案 | 无 | 否 |
| **L1 前端 UI** | 改 `.vue`、`.ts` 组件与类型 | `pnpm run typecheck` + `pnpm run lint:check` + 相关 vitest | 否 |
| **L2 后端业务** | 改 handler/service/repository，不动 DB schema | L1 全部 + `go test -tags=unit ./...` + `golangci-lint run ./...` | 否 |
| **L3 触及数据/部署** | 新增迁移、改 schema、改 Dockerfile、改 CI、改 `deploy/`、改用例配置 | L2 全部 + 迁移幂等性自检 + **回滚方案写在 PR 描述里** | **是** |

**L3 是闸门。** 涉及线上数据或不可逆动作时，先给方案、等批准、再动手。批准过一次不代表后续都批。

---

## 5. 验证：`scripts/verify`

一条命令跑完所有质量门：

```bash
bash scripts/verify.sh          # 默认：跑本机能跑的全部
bash scripts/verify.sh --quick  # 只跑毫秒级静态检查
```

**关于分层降级**：本机 Go 工具链当前处于损坏状态（`C:\Go\src` 下 std 核心实现文件已被系统性删除，且网络受限无法重装）。`verify.sh` 会自动检测，跑不了的层会**打印醒目警告并标记为「未验证」**，绝不假装通过。

脚本最终会输出一张汇总表，明确区分三类结果：

- ✅ **通过** —— 确实跑过且成功
- ⚠️ **未验证** —— 环境限制导致没跑成，需依赖 CI 或人工
- ❌ **失败** —— 跑了，且失败

**只看最后一行不够。** 有 ⚠️ 就说明这次改动还没被完整验证过。

CI 侧的对应关系见 `.github/workflows/backend-ci.yml`（push 与 PR 都触发）。

---

## 6. 文档体系

```
AGENTS.md                  ← 你在这里，唯一入口
docs/
  INDEX.md                 ← 全部文档索引
  CONSTRAINTS.md           ← 硬约束全集（每条都可被脚本/CI 检查）
  ITERATION_OS_CN.md       ← 项目地图、迭代节奏、发布闸门
  TEAM_QUICKSTART_CN.md    ← 团队上手、卡密套餐
  PRELAUNCH_GATE_CN.md     ← 上线前检查清单
  RELEASE_RUNBOOK_CN.md    ← 发布与回滚步骤
  RELEASE_CANDIDATE_*.md   ← 某次发布的验证证据
  BACKUP_RESTORE_DRILL_CN.md
  R2_BACKUP_SETUP_CN.md
  PAYMENT.md / PAYMENT_CN.md
  ADMIN_PAYMENT_INTEGRATION_API.md
DEV_GUIDE.md               ← 环境配置与 11 个常见坑点
openspec/                  ← 上游的变更记录写法（proposal/tasks/verification/spec）
```

### 写文档的规矩

- **本地产文档用 `_CN.md` 后缀**，英文文档不加后缀。这样与上游重叠的同名文件（`PAYMENT.md` 等）能一眼区分。
- **不要往 `AGENTS.md` 里堆内容。** 新流程写进 `docs/`，在这里加一行导航。
- **文档里不写死本机绝对路径。** 用相对路径或 `$REPO_ROOT`。
- **写完文档要更新 `docs/INDEX.md`。**

---

## 7. 提交规范

- Commit message 用 `type: 描述`，type 取值：`feat` / `fix` / `docs` / `chore` / `ci` / `refactor`
- 一个提交只做一件事。文档修复与代码改动**分开提交**。
- PR 模板见 `.github/PULL_REQUEST_TEMPLATE.md`。
- **不要自动 push。** 推到远端是外部动作，需要明确批准。

---

## 8. 环境陷阱速查

本机曾踩过的坑，动手前扫一眼：

| 陷阱 | 真相 |
| --- | --- |
| `make` 在 Windows 不可用 | 直接用 Makefile 里的原始命令，见 `DEV_GUIDE.md` 坑 8 |
| Go 工具链损坏 | `C:\Go\src` std 被裁剪，本地无法编译。用容器或依赖 CI |
| 网络受限 | `curl` 官方源报 SSL handshake failed，无法装依赖 |
| psql 中文路径 | 复制到纯英文路径再执行 |
| PowerShell 吞 bcrypt 的 `$` | 写 SQL 文件用 `psql -f`，不要 `psql -c` |
| 并行会话争抢 | 多个 AI 会话同时操作 git 会清 ref、报 EPERM。操作前先确认没有别的会话在跑 |
