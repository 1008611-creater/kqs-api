-- Seed CAU daily-refresh subscription plans.
-- These plans are intended for external shop/card-code sales first.
-- A week card grants the same daily quota every day for 7 days; it is not a
-- one-time weekly total.

WITH quota(daily_limit, sort_order) AS (
  VALUES
    (10::numeric, 10),
    (20::numeric, 20),
    (45::numeric, 45),
    (60::numeric, 60),
    (90::numeric, 90),
    (135::numeric, 135),
    (180::numeric, 180)
)
INSERT INTO groups (
  name,
  description,
  rate_multiplier,
  is_exclusive,
  status,
  platform,
  subscription_type,
  daily_limit_usd,
  weekly_limit_usd,
  monthly_limit_usd,
  default_validity_days,
  sort_order,
  created_at,
  updated_at,
  deleted_at
)
SELECT
  '每日' || daily_limit::text || '刀订阅',
  '订阅套餐专用分组：每天刷新 ' || daily_limit::text || ' 美元使用额度。',
  1.0,
  false,
  'active',
  'openai',
  'subscription',
  daily_limit,
  daily_limit * 7,
  NULL,
  1,
  2000 + sort_order,
  NOW(),
  NOW(),
  NULL
FROM quota
ON CONFLICT (name) WHERE deleted_at IS NULL DO UPDATE SET
  description = EXCLUDED.description,
  status = 'active',
  platform = 'openai',
  subscription_type = 'subscription',
  daily_limit_usd = EXCLUDED.daily_limit_usd,
  weekly_limit_usd = EXCLUDED.weekly_limit_usd,
  monthly_limit_usd = NULL,
  default_validity_days = 1,
  sort_order = EXCLUDED.sort_order,
  updated_at = NOW();

WITH subscription_groups AS (
  SELECT id
  FROM groups
  WHERE deleted_at IS NULL
    AND name IN (
      '每日10刀订阅',
      '每日20刀订阅',
      '每日45刀订阅',
      '每日60刀订阅',
      '每日90刀订阅',
      '每日135刀订阅',
      '每日180刀订阅'
    )
),
sched_accounts AS (
  SELECT DISTINCT a.id AS account_id
  FROM accounts a
  JOIN account_groups ag ON ag.account_id = a.id
  JOIN groups g ON g.id = ag.group_id
  WHERE g.name = 'codex'
    AND a.deleted_at IS NULL
    AND a.status = 'active'
    AND a.schedulable = true
)
INSERT INTO account_groups (account_id, group_id, priority)
SELECT sa.account_id, sg.id, 50
FROM sched_accounts sa
CROSS JOIN subscription_groups sg
ON CONFLICT (account_id, group_id) DO UPDATE SET priority = EXCLUDED.priority;

WITH plan_seed(name, description, price, original_price, validity_days, daily_limit, sort_order, product_name, features) AS (
  VALUES
    ('日卡10刀', '每天 10 美元额度，适合轻量试用。', 1.00::numeric, 1.00::numeric, 1, 10::numeric, 10, '农大API 日卡10刀',
      '每天$10额度' || E'\n' || '每日额度24小时窗口刷新' || E'\n' || '有效期1天' || E'\n' || '适用于 Codex / OpenAI 兼容模型'),
    ('日卡20刀', '每天 20 美元额度，适合日常轻度使用。', 1.90::numeric, 2.00::numeric, 1, 20::numeric, 20, '农大API 日卡20刀',
      '每天$20额度' || E'\n' || '每日额度24小时窗口刷新' || E'\n' || '有效期1天' || E'\n' || '适用于 Codex / OpenAI 兼容模型'),
    ('日卡45刀', '每天 45 美元额度，适合较高频使用。', 4.00::numeric, 4.50::numeric, 1, 45::numeric, 30, '农大API 日卡45刀',
      '每天$45额度' || E'\n' || '每日额度24小时窗口刷新' || E'\n' || '有效期1天' || E'\n' || '适用于 Codex / OpenAI 兼容模型'),
    ('日卡60刀', '每天 60 美元额度，适合集中使用。', 5.20::numeric, 6.00::numeric, 1, 60::numeric, 40, '农大API 日卡60刀',
      '每天$60额度' || E'\n' || '每日额度24小时窗口刷新' || E'\n' || '有效期1天' || E'\n' || '适用于 Codex / OpenAI 兼容模型'),
    ('日卡90刀', '每天 90 美元额度，适合高频使用。', 7.50::numeric, 9.00::numeric, 1, 90::numeric, 50, '农大API 日卡90刀',
      '每天$90额度' || E'\n' || '每日额度24小时窗口刷新' || E'\n' || '有效期1天' || E'\n' || '适用于 Codex / OpenAI 兼容模型'),
    ('周卡45刀', '连续 7 天，每天刷新 45 美元额度。', 25.00::numeric, 31.50::numeric, 7, 45::numeric, 110, '农大API 周卡45刀',
      '每天$45额度' || E'\n' || '每日额度24小时窗口刷新' || E'\n' || '有效期7天' || E'\n' || '可续费叠加有效期'),
    ('周卡90刀', '连续 7 天，每天刷新 90 美元额度。', 48.00::numeric, 63.00::numeric, 7, 90::numeric, 120, '农大API 周卡90刀',
      '每天$90额度' || E'\n' || '每日额度24小时窗口刷新' || E'\n' || '有效期7天' || E'\n' || '可续费叠加有效期'),
    ('周卡135刀', '连续 7 天，每天刷新 135 美元额度。', 70.00::numeric, 94.50::numeric, 7, 135::numeric, 130, '农大API 周卡135刀',
      '每天$135额度' || E'\n' || '每日额度24小时窗口刷新' || E'\n' || '有效期7天' || E'\n' || '可续费叠加有效期'),
    ('周卡180刀', '连续 7 天，每天刷新 180 美元额度。', 90.00::numeric, 126.00::numeric, 7, 180::numeric, 140, '农大API 周卡180刀',
      '每天$180额度' || E'\n' || '每日额度24小时窗口刷新' || E'\n' || '有效期7天' || E'\n' || '可续费叠加有效期')
),
resolved AS (
  SELECT ps.*, g.id AS group_id
  FROM plan_seed ps
  JOIN groups g ON g.deleted_at IS NULL AND g.name = '每日' || ps.daily_limit::text || '刀订阅'
),
updated AS (
  UPDATE subscription_plans sp
  SET group_id = r.group_id,
      description = r.description,
      price = r.price,
      original_price = r.original_price,
      validity_days = r.validity_days,
      validity_unit = 'day',
      features = r.features,
      product_name = r.product_name,
      for_sale = true,
      sort_order = r.sort_order,
      updated_at = NOW()
  FROM resolved r
  WHERE sp.name = r.name
  RETURNING sp.name
)
INSERT INTO subscription_plans (
  group_id,
  name,
  description,
  price,
  original_price,
  validity_days,
  validity_unit,
  features,
  product_name,
  for_sale,
  sort_order,
  created_at,
  updated_at
)
SELECT
  r.group_id,
  r.name,
  r.description,
  r.price,
  r.original_price,
  r.validity_days,
  'day',
  r.features,
  r.product_name,
  true,
  r.sort_order,
  NOW(),
  NOW()
FROM resolved r
WHERE NOT EXISTS (SELECT 1 FROM updated u WHERE u.name = r.name)
  AND NOT EXISTS (SELECT 1 FROM subscription_plans sp WHERE sp.name = r.name);
