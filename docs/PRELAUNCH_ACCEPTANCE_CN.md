# 上线前验收记录

> 本文件是 `docs/PRELAUNCH_GATE_CN.md` 的**执行记录**：记录每一项闸门实际跑到了什么程度、
> 拿到的证据是什么、还有什么没做完。只写可验证的事实，不写"应该没问题"。

- 仓库：`1008611-creater/kqs-api`
- 站点：`https://api.cauai.fun`
- 记录时间：2026-09-20
- 记录时的 HEAD：`a8949fc25`

---

## 0. 结论摘要

| 闸门项 | 状态 | 说明 |
| --- | --- | --- |
| 代码门禁（CI） | ✅ 通过 | 三个工作流全绿，见第 1 节 |
| 单测稳定性（flaky） | ✅ 已修复 | 定位到根因并修掉，见 2.1 |
| 凭据与生产数据保护 | ✅ 已加护栏 | 未泄露，但此前无防护，见 2.2 |
| 条款页线上可用 | ✅ 通过 | 四个页面均 200，见第 3 节 |
| 隔离备份恢复演练 | ❌ **未执行** | 阻塞：本机 Docker daemon 未启动 |
| `prelaunch-readiness.ps1` | ❌ **未执行** | 阻塞：需管理员凭据 |
| xlsx 安全例外 | ⏳ 待清理 | 2026-09-30 到期，需 pnpm 环境 |

**一句话**：代码层面的阻塞已清零；剩下的三项都是"必须有人在场"的动作，
不是靠改代码能推进的。具体需要做什么见第 4 节。

---

## 1. 代码门禁

最新提交 `a8949fc25` 的三个工作流：

| 工作流 | 关键 job | 结论 | 证据 |
| --- | --- | --- | --- |
| CI（verify） | — | success | run `35500722xxx` |
| CI（backend） | test / frontend / golangci-lint | success | 见下 |
| Security Scan | backend-security / frontend-security | success | 见下 |

（run id 以 GitHub Actions 页面为准；结论通过匿名 API 拉取 `conclusion` 字段核对。）

### 历史上最近一次失败已定位

提交 `25fd41580` 曾有一次 `test` job 失败。查 annotations 得到的真实内容是：

```
--- FAIL: TestAPIContracts (0.03s)
FAIL github.com/Wei-Shaw/sub2api/internal/server 0.066s
```

即契约快照落后于 DTO 新增字段，**不是**随机 flaky。已由 `9c51d838c` 修掉
（补 `supported_model_scopes`、`daily_rollover_usd`、`quota_bonus_multiplier`）。

---

## 2. 本轮修复项

### 2.1 单测 flaky：`TestOpenAIGatewayService_PrewarmReadHonorsParentContext`

**现象**：本机连续跑 `internal/service` 全量套件，8 轮中第 2 轮失败。

**现场证据**：

```
--- FAIL: TestOpenAIGatewayService_PrewarmReadHonorsParentContext (0.18s)
    openai_ws_forwarder_success_test.go:1035
    Error:      "184.5652ms" is not less than "180ms"
    Messages:   预热读取应受父 context 取消控制，不应阻塞到 read_timeout
```

**根因**：这不是数据竞争，是**时间断言的判别窗口太窄**。

- 用例让伪造连接阻塞 `readDelay = 200ms`，父 context 40ms 超时，
  然后断言 `elapsed < 180ms`。
- 判别逻辑本身是对的：遵守父 ctx 就 ~40ms 返回；不遵守就会等满 200ms。
- 但它只留了 **20ms 容错**。日志里 `cause=context deadline exceeded events=0`
  证明父 context **确实被遵守了**，行为完全正确，只是墙钟漂到了 184ms。

**修复**（`1e34a3fa9`）：不动断言意图，只把判别窗口拉宽。

- `readDelay`：`200ms` → `2s`
- 上限：`180ms` → `1s`

修复后：遵守父 ctx 仍然 ~40ms 返回（离上限 20 倍余量）；
一旦退化成只读 `read_timeout`，至少 2s 才返回，仍会**断言失败**。
也就是说**判别力没有下降**。

**验证**：该用例单独跑 40 次全过（2.257s，均值 ~56ms/次）；
全量套件在修复后连续 6 轮（run 3–8）全部 exit=0。

### 2.2 生产数据导出此前不受 gitignore 保护

**发现**：`git status` 显示 `ops/production/backup-20260912/` 未跟踪。查其内容：

| 文件 | 内容 | 敏感字段 |
| --- | --- | --- |
| `accounts.json` / `live-accounts.json` | 11 个上游账号 | **`credentials`**（上游 API Key / OAuth token） |
| `live-users.json` | 150 个用户 | email、余额 |
| `live-subscriptions.json` | 12 条订阅 | — |

**风险评估**：

- `git log --all -- ops/production/backup-20260912` 为空 → **从未进入 git 历史**
- 因此：不需要清历史、不需要轮换凭据（按固定规则，未泄露就不做这两件事）
- 但 `git check-ignore` 显示**没有任何规则覆盖它** →
  纯靠"没人手滑 `git add -A`"兜底，一次误操作就会把上游凭据推上 GitHub

**处置**（`a8949fc25`）：在 `.gitignore` 末尾追加两条护栏。
放在 `!ops/production/` 反忽略**之后**（gitignore 后置规则优先）：

```
/ops/production/backup-*/
/ops/production/sub2api-*
```

第二条顺带忽略 116MB 的发布二进制——原先的 `/sub2api-*` 锚定仓库根目录，
匹配不到 `ops/production/` 下的文件，所以它是漏网的。

**验证**：`git check-ignore -v` 确认三个敏感文件被忽略；
同时确认 `ops/production/*.sh` 部署脚本**未被误伤**（反忽略仍生效）。

**后续建议（非阻塞）**：gitignore 只是护栏。这些导出应当**移出工作区**，
放到仓库外的 secrets/备份目录，而不是留在工作区里靠忽略规则保护。

---

## 3. 条款合规（线上实测）

| 路径 | HTTP 状态 |
| --- | --- |
| `/legal/terms` | 200 |
| `/legal/privacy` | 200 |
| `/legal/refund` | 200 |
| `/legal/usage-policy` | 200 |

路由来源：`frontend/src/router/index.ts:171`（`/legal/:documentId`），
且已列入 `:737` 的 `BACKEND_MODE_ALLOWED_PATHS`。

> 备注：首轮探测时 `/legal/privacy` 返回 `000`（curl 失败），复查为 200。
> 判定为瞬时网络问题，非页面缺失。

---

## 4. 未完成项与解除阻塞的条件

### 4.1 隔离备份恢复演练 —— 未执行

手册：`docs/BACKUP_RESTORE_DRILL_CN.md`

阻塞原因：本机 Docker Desktop 的 Linux engine 未启动
（`dockerDesktopLinuxEngine` 连不上）。

**需要你做**：启动 Docker Desktop，然后按手册执行：
生产后台手动备份 → 临时 PostgreSQL + Redis 环境 → 配置同一 R2 bucket →
执行恢复 → 重启应用容器。

验收清单（四项全过才算完成）：
管理员能登录 / 用户列表存在 / 卡密兑换记录存在 / 用户余额正确。

### 4.2 `prelaunch-readiness.ps1` —— 未执行

阻塞原因：需要管理员凭据（`-AdminEmail` / `-AdminPassword` 或 `-AdminToken`）。

> 之前记录的"PowerShell 里 git 不可用"**已解除**：git 在 bash 里正常
> （`2.55.0`），只是那个 PowerShell 会话的 PATH 没带 Git。
> 用带 PATH 的终端跑即可：
> `$env:Path += ";C:\Program Files\Git\cmd"`

**需要你做**：

```powershell
$env:Path += ";C:\Program Files\Git\cmd"
.\scripts\prelaunch-readiness.ps1 -AdminEmail <邮箱> -AdminPassword <密码>
```

闸门要求（第 37 行）：备份 S3/R2 配置、连接测试、计划任务、备份新鲜度四项全部通过。

### 4.3 xlsx 安全例外 —— 待清理，2026-09-30 到期

`.github/audit-exceptions.yml` 有两条 xlsx 例外
（advisory `GHSA-4r6h-8v6p-xvw6`、`GHSA-5pgg-2g8v-p4x9`）。

前端唯一使用点：`frontend/src/views/admin/UsageView.vue:471`
（仅导出路径动态 `import('xlsx')`）。

阻塞原因：本机无 pnpm、corepack 损坏、无 `node_modules`，
改 lockfile 需要先装好依赖。

---

## 5. 观察项（有风险但暂无证据，本轮不改）

1. **其他紧时间断言**：`internal/handler/failover_loop_test.go` 里有
   50ms / 100ms / 200ms 一串上限，`openai_ws_pool_test.go:673` 有 80ms。
   本轮对其中三个嫌疑用例做了 40 次压力测试，**全部通过**，
   所以只列入观察，不盲改（改了反而可能削弱断言）。
2. **CI annotations 噪音**：annotations 会把测试里**故意触发的 panic 日志**
   （`openai.responses_panic_recovered`、`handler_dependencies_missing`）
   标成 failure 级别。不影响构建结论，但看失败明细时容易误判，建议后续过滤。
