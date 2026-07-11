# 生产发布流水线

## 入口

GitHub Actions 的 `Release Production` 为唯一正常发布入口。它必须使用受保护的 `production` Environment，并在执行前获得人工批准。

## 门禁

1. `pnpm run build`
2. Linux `go build -tags embed`
3. 二进制嵌入 HTML 检查
4. 上传后二进制 SHA-256 检查
5. 服务器镜像内可执行文件与嵌入 HTML 检查
6. 应用容器健康检查
7. `/health`、`/dashboard`、`/login`、`/keys`、`/monitor`、`/guide` 外部探测
8. 使用 `PROD_SMOKE_API_KEY` 的 `/v1/models` 探测

## 前端覆盖目录

当前生产环境存在 `/opt/sub2api/data/public` 覆盖目录，它优先于应用内嵌的前端资源。发布脚本会同步构建后的前端并为主入口生成 release-specific 文件名，避免 Cloudflare 对旧 immutable JS 的缓存覆盖新版页面。

覆盖目录不是独立源码来源。任何人工在服务器中修改它的行为都会在下一次正式发布中被覆盖。
