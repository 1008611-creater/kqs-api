# Cloudflare R2 备份配置手册

当前发布闸门要求备份可用：S3/R2 配置存在、连接测试成功、计划备份开启、最近一次成功备份小于 30 小时。

## 当前状态（2026-05-25）

- Cloudflare R2 已启用
- 私有 bucket 已配置为 `sub2api-backups`
- 矿泉水API 后台 S3/R2 连接已测试成功
- 定时备份已开启：`30 2 * * *`，即北京时间每天 02:30
- 已完成一次手动备份，备份记录 ID 是 `0cf00f51`

## 数据边界与隐私

当前后台备份功能执行的是 PostgreSQL 全量导出，并将压缩后的数据库备份上传到 R2。它不会主动读取或上传用户电脑中的真实文件、图片或文档目录。

数据库备份仍可能包含账号标识、余额、卡密兑换记录、API Key 元数据、站点配置以及用户错误日志，因此仍属于敏感运营数据。

- bucket 必须保持私有，不配置公开访问域名
- R2 token 仅授权 `sub2api-backups` 的 `Object Read & Write` 权限
- 不把 Access Key、Secret、备份下载链接写入聊天、仓库或公告
- 恢复演练只在隔离临时环境进行，不把备份交给普通用户

## 1. 启用 R2

打开 Cloudflare Dashboard：

```text
https://dash.cloudflare.com/?to=/:account/r2/overview
```

进入 `R2 object storage`，按页面提示启用 R2。Cloudflare 官方说明：必须先购买/启用 R2，之后才能生成 R2 API token。

## 2. 创建 Bucket

启用后可以在 Dashboard 里创建 bucket，也可以用 Wrangler：

```powershell
npx wrangler r2 bucket create sub2api-backups
npx wrangler r2 bucket list
```

Bucket 名称建议固定为：

```text
sub2api-backups
```

不要公开这个 bucket。它只用于数据库备份文件。

## 3. 创建 R2 API Token

在 `R2 object storage` 页面找到 `Account Details`，点击 `API Tokens` 旁边的 `Manage`，创建 token。

推荐配置：

- Token 类型：`Create Account API token`，如果页面权限不足，再用 `Create User API token`
- 权限：`Object Read & Write`
- 范围：只绑定 `sub2api-backups` 这个 bucket

创建完成后，只复制一次：

- `Access Key ID`
- `Secret Access Key`

Secret 只在 Cloudflare 页面显示一次，不要发给用户、不要写进仓库、不要贴到聊天里。

## 4. 填入矿泉水API后台

进入管理员后台：`设置 -> 数据备份`，填写：

```text
Endpoint: https://3b78db8f07bf8738064ca1606fe89cb5.r2.cloudflarestorage.com
Region: auto
Bucket: sub2api-backups
Prefix: backups/
Access Key ID: Cloudflare 生成的 Access Key ID
Secret Access Key: Cloudflare 生成的 Secret Access Key
Force Path Style: 关闭
```

保存前先点 `测试连接`。测试通过后再保存。

## 5. 开启计划备份

推荐配置：

```text
启用: 是
Cron: 30 2 * * *
保留天数: 14
保留份数: 30
```

说明：`30 2 * * *` 表示每天 02:30。容器时区当前配置为 `Asia/Shanghai`。

## 6. 发布前验证

手动创建一次备份，等待状态变为 `completed`。然后运行：

```powershell
Set-Location E:\codex\kqsapi
.\scripts\prelaunch-readiness.ps1 -ApiBaseUrl https://api.cauai.fun
```

必须看到这些项目通过：

- `backup s3 config`
- `backup s3 connection`
- `backup schedule`
- `backup freshness`
- `backup failures`

如果 `backup freshness` 失败，不要正式公开公布。

## 参考

- Cloudflare R2 API token 文档：`https://developers.cloudflare.com/r2/api/tokens/`
- Cloudflare R2 创建 bucket 文档：`https://developers.cloudflare.com/r2/buckets/create-buckets/`
- Cloudflare R2 S3 API 兼容性：`https://developers.cloudflare.com/r2/api/s3/api/`
- Cloudflare R2 Go SDK 示例：`https://developers.cloudflare.com/r2/examples/aws/aws-sdk-go/`
