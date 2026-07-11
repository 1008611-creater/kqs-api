UPDATE announcements
SET content = $$近期网站注册入口遭遇批量恶意注册攻击。为保护真实同学的试用额度和平台稳定性，校训福利领取规则已调整：

- 30 美金额度仅限使用 `@cau.edu.cn` 农大邮箱注册的账号领取。
- 非农大邮箱账号将无法领取校训福利。
- 已发现的异常批量注册账号会被禁用并进入人工核验。
- 正常购买卡密、兑换余额和使用 API 不受影响。

请使用农大邮箱注册后，在 [兑换码页面](/redeem) 输入校训口令领取。$$,
    status = 'active',
    notify_mode = 'popup',
    targeting = '{}'::jsonb,
    starts_at = COALESCE(starts_at, NOW()),
    updated_at = NOW()
WHERE title = '农大API 校训福利开放领取';

INSERT INTO announcements (
    title,
    content,
    status,
    notify_mode,
    targeting,
    starts_at,
    created_at,
    updated_at
)
SELECT
    '农大API 校训福利领取规则调整',
    $$近期网站注册入口遭遇批量恶意注册攻击。为保护真实同学的试用额度和平台稳定性，校训福利领取规则已调整：

- 30 美金额度仅限使用 `@cau.edu.cn` 农大邮箱注册的账号领取。
- 非农大邮箱账号将无法领取校训福利。
- 已发现的异常批量注册账号会被禁用并进入人工核验。
- 正常购买卡密、兑换余额和使用 API 不受影响。

请使用农大邮箱注册后，在 [兑换码页面](/redeem) 输入校训口令领取。$$,
    'active',
    'popup',
    '{}'::jsonb,
    NOW(),
    NOW(),
    NOW()
WHERE NOT EXISTS (
    SELECT 1
    FROM announcements
    WHERE title = '农大API 校训福利领取规则调整'
);

UPDATE announcements
SET content = $$近期网站注册入口遭遇批量恶意注册攻击。为保护真实同学的试用额度和平台稳定性，校训福利领取规则已调整：

- 30 美金额度仅限使用 `@cau.edu.cn` 农大邮箱注册的账号领取。
- 非农大邮箱账号将无法领取校训福利。
- 已发现的异常批量注册账号会被禁用并进入人工核验。
- 正常购买卡密、兑换余额和使用 API 不受影响。

请使用农大邮箱注册后，在 [兑换码页面](/redeem) 输入校训口令领取。$$,
    status = 'active',
    notify_mode = 'popup',
    targeting = '{}'::jsonb,
    starts_at = COALESCE(starts_at, NOW()),
    updated_at = NOW()
WHERE title = '农大API 校训福利领取规则调整';
