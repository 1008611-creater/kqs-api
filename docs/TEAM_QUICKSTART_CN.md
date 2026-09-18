# 矿泉水API 快速上手

这是一份给管理员和团队成员都能直接照着做的简版流程。

当前部署示例：

- 站点地址：`https://api.cauai.fun`
- 管理后台：`/admin`
- 用户密钥页：`/keys`
- 默认模型：`gpt-5.5`

## 1. 管理员先做什么

1. 打开管理后台：`https://api.cauai.fun/admin/accounts`
2. 导入或检查上游账号
3. 把可用账号放进同一个分组，比如 `codex`
4. 用"测试连接"确认账号可用
5. 准备邀请码，并让成员自行注册、兑换和创建各自的 API Key

## 2. 成员先做什么

1. 打开注册页：`https://api.cauai.fun/register`
2. 用管理员提供的邀请码注册自己的账号，并同意服务条款
3. 登录：`https://api.cauai.fun/login`
4. 在新手引导中选择链动小铺套餐购买卡密
5. 到 `/redeem` 输入卡密，确认余额增加
6. 到 `/keys` 自己创建 `codex` 分组的密钥，并复制 Codex 配置
7. 每个成员都用自己的账号和自己的 Key

当前卡密套餐：

- 购买入口统一使用链动小铺店铺页：`https://pay.ldxp.cn/shop/B59CCLX7`
- 店铺内有 `4.9 元 / 15.9 美金余额`、`9.9 元 / 34.9 美金余额`、`19.9 元 / 79.9 美金余额` 三档卡密

## 3. 上游账号怎么加

1. 在 `账号管理` 里点 `添加账号`
2. 选择对应授权方式
3. 导入完成后，选中账号
4. 点 `分组管理`
5. 选择 `codex`
6. 批量更新

## 4. 成员怎么创建可用 Key

1. 成员先注册、登录并完成卡密兑换
2. 成员打开自己的密钥页：`https://api.cauai.fun/keys`
3. 点 `创建密钥`
4. 给这个 Key 选择 `codex` 分组，并保留默认额度/速率保护
5. 创建后只由本人保存完整 Key
6. 不共享同一个 Key，不把完整 Key 发给管理员或同学

## 5. 成员电脑怎么配

编辑本机 `~/.codex/config.toml`，放成下面这样：

```toml
model = "gpt-5.5"
review_model = "gpt-5.5"

[model_providers.OpenAI]
name = "OpenAI"
base_url = "https://api.cauai.fun"
wire_api = "responses"
requires_openai_auth = true
```

然后在 Codex 里填自己拿到的 API Key。

如果成员用的是支持环境变量的客户端，也可以在 `/keys` 页面点 `使用密钥`，直接复制页面给出的配置。

## 6. 怎么测试是不是通了

最简单的测试方法：

1. 保存配置
2. 重启 Codex
3. 发一句短测试，比如 `只回复两个字：通了`

如果返回 `通了`，就说明链路正常。

## 7. 常见问题

- 如果报 `No available accounts`，先检查上游账号是否都在同一个分组里
- 如果报鉴权错误，先确认 API Key 是否填对
- 如果模型不对，确认 `model = "gpt-5.5"`
- 如果本机还在连 `localhost:18080`，把 `base_url` 改成公网域名
