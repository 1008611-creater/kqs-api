# 矿泉水API 发布候选记录 - 2026-05-25

本记录描述当前已经验证过的事实，以及阻止正式公开发布的剩余闸门。当前阶段适合继续管理员验收和受控邀请制 Beta，不应直接作为无门槛公开发布声明。

## 已验证证据

| 项目 | 结果 | 证据 |
| --- | --- | --- |
| 固定镜像运行 | 通过 | `sub2api-gg` 运行 `sub2api:kqs-api-prelaunch-20260525-2058`，健康状态 `healthy` |
| 健康检查 | 通过 | `http://localhost:18080/health` 与 `https://api.cauai.fun/health` 返回 `{"status":"ok"}` |
| 后端测试 | 通过 | 使用 `golang:1.26.6-alpine` 容器执行 `go test ./...` 通过 |
| 前端构建 | 通过 | Docker 镜像构建过程中 `pnpm run build` 成功 |
| 迁移启动 | 通过 | `142_cau_prelaunch_defaults.sql` 已兼容历史 checksum，`143_reaffirm_cau_prelaunch_defaults.sql` 已成功应用 |
| Cloudflare Worker | 通过 | `sub2api-proxy` 已部署，版本 `768ca6f6-93c9-4d6c-8e8a-1ab4cf9fadf6` |
| 可信客户端 IP | 通过 | Worker 规范化访问 IP 头，后端启用 `SERVER_TRUSTED_PLATFORM=cloudflare`，启动日志不再提示 `server.trusted_proxies` 为空 |
| URL 白名单 | 通过 | 已启用 `security.url_allowlist.enabled`，显式放行 `api.openai.com`、`auth.openai.com`、`chatgpt.com` 等必要上游，启动日志不再提示白名单关闭 |
| Codex 上游烟测 | 部分通过 | 后台账号测试已到达 `chatgpt.com`，返回上游 `usage_limit_reached`；说明白名单未拦截，但当前 active OAuth 账号额度不足，尚无成功生成证据 |
| 邀请注册门槛 | 通过 | 干净访客打开 `/register` 可见邀请码输入框与协议确认入口 |
| 公开合规文档 | 通过 | `/legal/terms`、`/legal/privacy`、`/legal/refund`、`/legal/usage-policy` 均可读取正文 |
| R2 基础备份 | 通过 | bucket `sub2api-backups` 已连接，定时策略 `30 2 * * *` 已开启，已有手动备份 `0cf00f51` |

## 数据与隐私边界

- 当前 R2 备份来自 PostgreSQL 全量导出，不会主动采集用户电脑中的真实文件或文档目录。
- 备份仍包含业务敏感数据，bucket 保持私有，token 仅限指定 bucket 的对象读写。
- 不要求用户提交完整 Key、完整卡密、银行卡、验证码或 cookie 以处理客服问题。

## 尚未放行的闸门

1. 工作区仍包含大量未提交改动；尚未形成带 git SHA 的冻结发布提交和 tag。
2. 还未在隔离临时环境完成一次 R2 备份恢复演练。
3. 需要在至少丢失一个 OpenAI OAuth 上游账号额度恢复或补充新账号后，再用全新测试用户走完公网链路：邀请码注册、卡密兑换、创建 `codex` Key、复制配置、`/responses` 成功。
4. 需要触发并记录重复卡密、错误卡密、余额不足、无权限 Key 等错误场景的用户指引表现。
5. 正式宣布前应做一次旧镜像回滚演练并记录健康检查结果。
6. `CORS allowed_origins` 留空在当前同源 UI/API 部署下表示默认拒绝浏览器跨域请求，已记录为可接受策略；未来若拆分浏览器前端，再配置精确白名单。

## 放量结论

- 当前结论：基础发布闸门已通过，可以进行管理员自测和少量受邀请用户 Beta 验收。
- 当前结论：尚不建议直接面向学校范围公开发布与收费推广。
- 当前最小动作：补充或等待一个可用 OpenAI OAuth 上游账号额度，随后创建一个全新测试账号完成公网全流程与错误场景验收。
