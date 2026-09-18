# 矿泉水API发布与回滚手册

本手册用于把一次改动从"本地能跑"推进到"可收费试运营"。核心原则：固定版本、先备份、可回滚、再公布。

## 当前候选版本状态（2026-05-25）

- 本机验证镜像：`sub2api:kqs-api-prelaunch-20260525-2058`
- 本机容器 `sub2api-gg` 已运行该固定 tag，`/health` 返回 `200`
- Cloudflare Worker `sub2api-proxy` 已部署版本 `768ca6f6-93c9-4d6c-8e8a-1ab4cf9fadf6`
- 当前候选已启用 `SERVER_TRUSTED_PLATFORM=cloudflare`，公网入口的可信客户端 IP 链路已纳入后端 `ClientIP()`
- 当前候选已启用 `SECURITY_URL_ALLOWLIST_ENABLED=true`，并显式放行 OpenAI/Codex 必需上游
- 当前代码工作区仍有未提交改动，因此该镜像属于发布候选，不等同于已冻结 tag
- 正式标记 `kqs-api-v0.2.0-beta` 前，必须先把发布范围整理成干净提交，并记录对应 git SHA
- 本次候选验证详情见 `docs/RELEASE_CANDIDATE_20260525_CN.md`

## 1. 发布前冻结

1. 确认当前目标版本，例如 `kqs-api-v0.2.0-beta`
2. 查看工作区：

```powershell
git status --short
git diff --name-only
```

3. 只保留本次发布需要的改动。无关改动要么单独提交，要么明确记录为不发布
4. 记录当前线上版本：

```powershell
docker inspect sub2api-gg --format '{{.Config.Image}}'
docker inspect sub2api-gg --format '{{.Image}}'
```

把结果写入发布记录，作为回滚目标。

## 2. 发布前备份

在后台进入 `管理后台 -> 数据备份`：

1. 确认 S3/R2 连接测试通过
2. 确认计划备份开启，建议北京时间每天 02:30，保留 14 天或 30 份
3. 点击手动备份，等待状态变为 `completed`
4. 记录备份 ID、文件名、完成时间和大小

禁止用生产库做恢复演练。恢复演练必须在临时环境执行。

## 3. 本地验证

后端：

```powershell
Set-Location E:\codex\kqsapi
$backend = (Resolve-Path .\backend).Path
docker run --rm -v "${backend}:/app/backend" -w /app/backend `
  -e GOPROXY=https://goproxy.cn,direct `
  -e GOSUMDB=sum.golang.google.cn `
  golang:1.26.3-alpine sh -lc `
  'export PATH=/usr/local/go/bin:$PATH; apk add --no-cache git ca-certificates tzdata >/dev/null; go test ./...'
```

说明：当 Windows 宿主机没有完整 Go toolchain 时，发布验证统一使用固定 Go 容器执行同一份源码的全量测试。

> 补充：本机 `C:\Go` 的 std 源文件已被裁剪，无法直接编译。请优先使用上面的容器方式，或依赖 GitHub Actions 的 `backend-ci.yml`。

前端：

```powershell
Set-Location E:\codex\kqsapi\frontend
pnpm run typecheck
pnpm run build
```

发布闸门：

```powershell
Set-Location E:\codex\kqsapi
.\scripts\prelaunch-readiness.ps1 -ApiBaseUrl https://api.cauai.fun
```

脚本会优先使用 `-AdminToken`，否则会尝试读取本机 `deploy/.env` 的管理员账号，在本机登录后只把临时 token 用于备份健康检查，不会打印密码或 token。

如果在 CI 或远程机器上没有 `deploy/.env`，显式传入管理员 token：

```powershell
.\scripts\prelaunch-readiness.ps1 -ApiBaseUrl https://api.cauai.fun -AdminToken $env:SUB2API_ADMIN_TOKEN
```

如果公网域名只用于 `/health`，但管理接口希望走本机，可显式指定：

```powershell
.\scripts\prelaunch-readiness.ps1 -ApiBaseUrl https://api.cauai.fun -AdminApiBaseUrl http://localhost:18080
```

备份闸门必须同时通过：S3/R2 配置存在、连接测试成功、计划备份开启、最近一次 completed 备份小于 30 小时。

## 4. 构建固定镜像

不要把正式发布绑定到 `latest`。发布镜像必须包含版本和 git 短 SHA。

PowerShell 示例：

```powershell
$sha = git rev-parse --short HEAD
$image = "sub2api:kqs-api-v0.2.0-beta-$sha"
bash .\deploy\build_image.sh $image
```

如果 Docker Hub 拉取基础镜像失败，先解决基础镜像拉取或配置镜像源；不要把临时 `docker cp` 二进制覆盖当作正式发布。

## 5. 上线

1. 修改 `deploy/.env`：

```env
SUB2API_IMAGE=sub2api:kqs-api-v0.2.0-beta-<sha>
```

2. 启动固定版本：

```powershell
docker compose --env-file .\deploy\.env -f .\deploy\docker-compose.yml up -d
```

3. 验证容器：

```powershell
docker ps --filter "name=sub2api-gg"
docker logs --tail 120 sub2api-gg
Invoke-RestMethod https://api.cauai.fun/health
```

4. 浏览器烟测：

- 无登录状态访问 `/register`，能看到邀请码输入和四份协议入口
- 无登录状态访问 `/legal/terms`、`/legal/privacy`、`/legal/refund`、`/legal/usage-policy`，正文均可读取
- `/dashboard` 能看到新手引导和支持入口
- `/redeem` 能看到卡密说明和支持入口
- `/keys` 新建 Key 默认有额度和速率保护，且可复制 Codex 配置

## 6. Cloudflare 缓存

前端静态资源或 logo 更新后，到 Cloudflare 执行缓存刷新：

1. 进入 `cauai.fun`
2. 打开 `Caching -> Configuration -> Purge Cache`
3. 优先 Purge URL：`https://api.cauai.fun/`、`https://api.cauai.fun/assets/*`
4. 若页面仍旧，用 Purge Everything

刷新后用无痕窗口访问 `https://api.cauai.fun/home` 和 `/login`。

## 7. 回滚

回滚只切镜像，不改数据库，除非本次发布包含不可兼容迁移。

```powershell
# 把 deploy/.env 里 SUB2API_IMAGE 改回发布前记录的镜像
docker compose --env-file .\deploy\.env -f .\deploy\docker-compose.yml up -d
Invoke-RestMethod https://api.cauai.fun/health
```

如果数据库迁移造成不可逆问题，使用发布前手动备份在临时环境先验证，再决定是否恢复生产。不要在没有验证的情况下直接覆盖生产库。

## 8. 正式公布条件

只有全部满足才发布给客户：

- 工作区已冻结成可追溯提交，发布 tag 与镜像 tag 均记录 git SHA
- 固定镜像 tag 已上线，线上不是未记录的 `latest` 漂移版本
- 发布前手动备份完成，并且最近自动备份小于 30 小时
- 登录条款、退款规则、隐私说明、可接受使用政策可访问
- 新用户流程完成：注册、兑换、创建 Key、复制配置、`/responses` 成功
- 错误提示覆盖：重复卡密、错误卡密、余额不足、Key 未分组、无可用 Key
- 已记录回滚镜像和回滚命令
