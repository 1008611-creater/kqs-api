# Secrets 清单模板

此文件只记录 Secret 的用途、存储位置、轮换责任人和最后验证时间，不填写任何实际值。

| 名称 | 存储位置 | 用途 | 轮换后验证 |
| --- | --- | --- | --- |
| `PROD_SSH_PRIVATE_KEY` | GitHub production Environment | CI 发布到首尔服务器 | SSH 连通、部署 dry run |
| `PROD_SSH_KNOWN_HOSTS` | GitHub production Environment | 防止 SSH 主机伪造 | `ssh -o StrictHostKeyChecking=yes` |
| `PROD_SMOKE_API_KEY` | GitHub production Environment | 发布后的低成本 API 探测 | `/v1/models` 返回 200 |
| Cloudflare 最小权限 Token | 独立 Cloudflare Environment | DNS/Tunnel/缓存操作 | 对应 API 只读或变更验证 |
| 生产 `.env` | 首尔服务器 `/opt/sub2api/.env` | 应用运行时密钥 | 容器健康、登录、网关调用 |
