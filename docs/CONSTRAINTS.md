# CONSTRAINTS.md — 硬约束全集

> 本文件列出本项目**不可妥协**的约束。
>
> 与普通「最佳实践」的区别：每条约束都标注了**检查手段**。如果一个约束无法被脚本、CI 或明确的 checklist 检查，它就不该出现在这里。
>
> 标注为「人工」的条目，说明当前技术手段无法自动检查，需要在 PR 描述里显式声明已确认。

---

## C1 · 包管理器必须是 pnpm

**约束**：前端依赖安装与构建一律使用 pnpm。禁止使用 npm 或 yarn 安装依赖。

**为什么**：CI 使用 `pnpm install --frozen-lockfile`。用 npm 装过一次会导致 `node_modules` 结构不兼容，后续 `pnpm install` 报 EPERM。

**检查手段**：`scripts/verify.sh` 检查 `frontend/node_modules/.pnpm` 是否存在（存在=pnpm 装的）。

---

## C2 · `pnpm-lock.yaml` 必须与 `package.json` 同步提交

**约束**：任何对 `frontend/package.json` 的依赖改动，必须在同一个提交里包含更新后的 `frontend/pnpm-lock.yaml`。

**为什么**：`--frozen-lockfile` 会在 lock 与 package.json 不一致时直接失败，CI 红。

**检查手段**：`scripts/verify.sh` 用 `git diff --name-only` 检查：若本次改动含 `package.json` 但含 `pnpm-lock.yaml` 则通过，否则失败。

---

## C3 · Ent schema 改动必须重新生成并提交

**约束**：改 `backend/ent/schema/*.go` 后必须执行 `go generate ./ent`，并把生成的代码一起提交。

**为什么**：生成代码不入库会导致别人的构建与你本机不一致；CI 编译用的就是仓库里的生成代码。

**检查手段**：`scripts/verify.sh` 检查 `backend/ent/schema/` 的修改时间是否晚于 `backend/ent/` 下生成文件（启发式，可能误报）；最终判据是 CI 编译通过。**部分人工。**

---

## C4 · Go interface 新增方法后必须补齐所有 stub

**约束**：给任何 interface 增加方法后，仓库内所有实现该 interface 的 struct（含测试 mock/stub）必须同步补齐该方法。

**为什么**：漏一个就整包编译失败，且报错信息指向 stub 文件而非你改的 interface，容易找不到原因。

**检查手段**：`go build ./...` + `go test -tags=unit ./...` 编译通过即为满足。本机跑不了时**必须依赖 CI**。

---

## C5 · 数据库迁移：文件名为唯一标识，已执行内容不可变更

**约束**：
1. 迁移的**唯一标识是文件名**（`schema_migrations.filename` 是主键），不是编号。编号重复是允许的，上游自己也在重复用号。
2. **已经执行过的迁移文件，内容不得修改。** 包括"只是改个注释"。

**为什么**：`migrations_runner.go` 会对每个已执行迁移做 checksum 校验。内容一变，服务启动时直接失败——而且是在生产启动路径上失败。

**要改一个已上线的迁移怎么办**：新增一个迁移文件来做修正，不要改原文件。

**检查手段**：**必须人工。** 迁移文件属于 L3 变更，PR 描述里必须声明「本次新增迁移 X 个，未修改任何已执行的迁移文件」。可用 `git diff --name-status kqs/upstream-base -- backend/migrations/` 辅助确认历史文件是否被改。

---

## C6 · 迁移必须幂等

**约束**：新增迁移脚本必须可重复执行不报错。使用 `IF NOT EXISTS` / `IF EXISTS` / `WHERE` 条件式 / `ON CONFLICT DO NOTHING`。

**为什么**：迁移可能在部分失败后被重跑；不具备幂等性会让恢复变得困难。

**检查手段**：**人工审查 SQL**。项目内既有迁移（如 `160_kqs_branding_transition.sql`、`162_disable_motto_gift_promo_code.sql`）可作为写法参照。

---

## C7 · 不得提交任何真实密钥

**约束**：API Key、OAuth token、数据库密码、R2/S3 凭据、卡密明文，一律不得进入仓库。示例值用明显的占位符。

**为什么**：仓库是公开的，且历史一旦推上去，删文件不等于删除记录——**必须换凭证才算处置**。

**检查手段**：`python tools/secret_scan.py`。CI 侧见 `.github/workflows/security-scan.yml`（govulncheck + gosec + pnpm audit）。

**注意**：项目内已有迁移用字符串拼接（`'c' || 'a' || 'u_motto_gift_claims'`）规避扫描器误报的先例。这是**为规避误报**，不是可以藏真密钥的借口。

---

## C8 · 不改写 git 历史，不 force push 到 main

**约束**：不 rebase 已推送的提交、不 `push --force` 到 `main`、不 `git reset --hard` 已推送的分支。

**为什么**：本仓库有并行 AI 会话与人工操作同时在跑。改写历史会导致其他人工作区损坏，且可能丢改动。

**检查手段**：**人工。** PR 模板中有确认项。

---

## C9 · 不可逆动作必须先给闸门

**约束**：涉及线上数据、发布、回滚、删除的操作，必须先说明方案与影响范围，等明确批准后再执行。

**范围包括**：`docker` 部署、数据库迁移执行、Cloudflare 配置变更、R2 备份/恢复、任何 `rm -rf`。

**为什么**：这些动作不可逆。批准过一次不代表后续都批。

**检查手段**：**人工。** L3 变更在 PR 描述里必须包含回滚方案。

---

## C10 · 本机 Go 环境不可信

**约束**：不得声称「本地编译/测试通过」除非确实跑成功。本机 `C:\Go\src` 下 std 核心实现文件已被系统性删除，连 hello world 都无法编译。

**为什么**：一次"假装验证过"会导致 Go 代码零验证地进入仓库。

**做法**：
- 本地验证优先用容器：`docker run --rm -v "$PWD/backend:/app" -w /app golang:1.26.3-alpine go test ./...`
- 跑不了就明确写「未验证」，依赖 CI
- `scripts/verify.sh` 会自动检测并标注

**检查手段**：`scripts/verify.sh` 的汇总表。

---

## C11 · 上游同步走增量摘取，不做整体合并

**约束**：与 `Wei-Shaw/sub2api` 同步时，逐个文件比对手工移植。**禁止** `git merge` / `git cherry-pick` / `git rebase` 上游分支。

**为什么**：两边无共同祖先，合并必然是全量冲突，且会引入上游针对其自身架构的重构（如把 i18n 从单文件 14767 行重构为 `locales/{en,zh}/` 34 个文件）。半截引入比不引入更危险。

**检查手段**：**人工。** 同步动作需在 PR 描述里列出「摘取了哪些上游提交、为什么、影响了哪些文件」。

---

## C12 · 文档写入用 ASCII 引号

**约束**：写代码、配置（JSON/YAML/TOML）、shell 命令时，字符串定界符一律用 ASCII 直引号（`"` / `'`）。中文自然语言内容不受此限。

**为什么**：中文全角引号会导致 JSON 解析失败、shell 命令错乱。

**检查手段**：`scripts/verify.sh` 检查 `.json` / `.yaml` / `.toml` 文件是否包含全角引号。

---

## C13 · 文档文件编码必须是 UTF-8（无 BOM）

**约束**：所有 `.md` 文件必须是无 BOM 的 UTF-8。

**为什么**：本仓库历史上发生过二次编码损坏（UTF-8 → 误按 GB18030 解码 → 再存为 UTF-8），导致 6 份文档正文出现 `\ufffd` 乱码，且**丢失的字节数不固定，无法可靠逆向还原**——只能人工重写。

**检查手段**：`scripts/verify.sh` 检查所有 `.md` 是否含 BOM 或 `\ufffd`。

---

## 约束总览

| 编号 | 约束 | 检查手段 |
| --- | --- | --- |
| C1 | 用 pnpm | 脚本 |
| C2 | lock 文件同步 | 脚本 |
| C3 | Ent 生成代码提交 | 脚本（启发式）+ CI |
| C4 | interface stub 补齐 | CI 编译 |
| C5 | 迁移内容不可变 | **人工** |
| C6 | 迁移幂等 | **人工** |
| C7 | 不提交密钥 | 脚本 + CI |
| C8 | 不改写历史 | **人工** |
| C9 | 不可逆动作先给闸门 | **人工** |
| C10 | 不假称本地验证通过 | 脚本 |
| C11 | 上游增量摘取 | **人工** |
| C12 | 配置文件用 ASCII 引号 | 脚本 |
| C13 | 文档 UTF-8 无 BOM | 脚本 |
