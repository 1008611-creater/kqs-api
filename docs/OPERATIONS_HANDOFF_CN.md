# 矿泉水 API 多电脑 Codex 交接

## 目标

任何获授权电脑上的 Codex 都可以从同一 Git 提交构建、测试和申请发布，不需要登录腾讯云控制台，也不需要复制生产数据库或 root 私钥。

## 架构

```text
开发电脑 / Codex -> 私有 Git 仓库 -> GitHub Actions production 环境 -> 首尔生产服务器
```

- 生产服务器：`43.133.254.221:22`
- 生产目录：`/opt/sub2api`
- 应用：`sub2api-gg`
- 数据库：`sub2api-gg-postgres`
- Redis：`sub2api-gg-redis`
- 域名：`https://api.cauai.fun`
- Cloudflare Tunnel 只由首尔服务器维护；禁止任何开发电脑连接生产 Tunnel。

## 新电脑首次接入

1. 使用自己的 GitHub 账号获得私有仓库权限。
2. 克隆仓库并运行：

   ```powershell
   ./scripts/handoff-preflight.ps1
   ```

3. 仅在本地使用独立开发配置和数据库；不得复制 `/opt/sub2api/data` 到本地后作为生产写入端。
4. 每个任务从最新 `main` 创建 `feature/<task>` 分支，完成测试后提交 PR。
5. 生产发布只能通过 `Release Production` 工作流，并由 production Environment 批准。

## 生产 Secrets

仅保存于 GitHub `production` Environment，绝不提交进仓库：

```text
PROD_SSH_HOST
PROD_SSH_PORT
PROD_SSH_USER
PROD_SSH_PRIVATE_KEY
PROD_SSH_KNOWN_HOSTS
PROD_SMOKE_API_KEY
```

Cloudflare DNS、Tunnel 或缓存变更使用独立 Environment 和最小权限 Token，不进入常规发布工作流。

## 发布与回滚

- 正常发布：Actions 构建前端、构建 `-tags embed` Linux 二进制、上传并调用 `ops/production/deploy-remote.sh`。
- 发布脚本会备份 `.env` 和 `data/public` 覆盖前端，只重建 `sub2api-gg`，不重启 Postgres/Redis。
- 应用健康检查失败时自动恢复此前镜像。
- 紧急回滚：恢复 `/opt/sub2api/backups/releases/.env.before-*` 中的 `SUB2API_IMAGE`，然后运行：

  ```bash
  cd /opt/sub2api
  docker compose up -d --no-deps --force-recreate sub2api
  ```

## 不可违反的生产规则

- 前端发布必须走 `go build -tags embed`。
- 不允许从开发电脑启动生产 Cloudflare Tunnel。
- 不允许重启 PostgreSQL 或 Redis 作为普通发布步骤。
- 不允许把 API Key、SSH 私钥、Cloudflare Token、生产 `.env`、数据库备份提交到 Git。
- 生产数据库只有首尔服务器是写入主节点。
