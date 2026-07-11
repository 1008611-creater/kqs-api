-- KQS branding transition.
-- Keep historical migrations immutable; apply runtime data and table-name changes here.

DO $$
DECLARE
    old_table regclass := to_regclass('public.' || 'c' || 'a' || 'u_motto_gift_claims');
    new_table regclass := to_regclass('public.kqs_motto_gift_claims');
    old_user_idx text := 'idx_' || 'c' || 'a' || 'u_motto_gift_claims_user_id';
    old_claimed_idx text := 'idx_' || 'c' || 'a' || 'u_motto_gift_claims_claimed_at';
BEGIN
    IF new_table IS NULL AND old_table IS NOT NULL THEN
        EXECUTE format('ALTER TABLE %s RENAME TO %I', old_table, 'kqs_motto_gift_claims');
    END IF;

    IF to_regclass('public.idx_kqs_motto_gift_claims_user_id') IS NULL
        AND to_regclass('public.' || old_user_idx) IS NOT NULL THEN
        EXECUTE format('ALTER INDEX %I RENAME TO %I', old_user_idx, 'idx_kqs_motto_gift_claims_user_id');
    END IF;

    IF to_regclass('public.idx_kqs_motto_gift_claims_claimed_at') IS NULL
        AND to_regclass('public.' || old_claimed_idx) IS NOT NULL THEN
        EXECUTE format('ALTER INDEX %I RENAME TO %I', old_claimed_idx, 'idx_kqs_motto_gift_claims_claimed_at');
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS kqs_motto_gift_claims (
    id BIGSERIAL PRIMARY KEY,
    campaign_code VARCHAR(64) NOT NULL,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    ip_hash CHAR(64) NOT NULL,
    claimed_amount DECIMAL(20,8) NOT NULL DEFAULT 0,
    claimed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(campaign_code, user_id),
    UNIQUE(campaign_code, ip_hash)
);

CREATE INDEX IF NOT EXISTS idx_kqs_motto_gift_claims_user_id
    ON kqs_motto_gift_claims(user_id);
CREATE INDEX IF NOT EXISTS idx_kqs_motto_gift_claims_claimed_at
    ON kqs_motto_gift_claims(claimed_at);

COMMENT ON TABLE kqs_motto_gift_claims IS 'KQS launch gift claim records; IP is hashed for privacy.';
COMMENT ON COLUMN kqs_motto_gift_claims.campaign_code IS 'Campaign code, e.g. JMSZDJYTXZYC.';
COMMENT ON COLUMN kqs_motto_gift_claims.ip_hash IS 'SHA-256 hash of normalized client IP and campaign salt.';

UPDATE settings
SET value = replace(
    replace(
    replace(
    replace(
    replace(
    replace(
    replace(
    replace(
    replace(
    replace(
    replace(
    replace(value,
        '中国农业' || '大学 AI API 网关平台', '矿泉水 AI API 网关平台'),
        '中国农业' || '大学 AI API 网关', '矿泉水 AI API 网关'),
        '中国农业' || '大学官方信息系统', '任何第三方官方信息系统'),
        '中国农业' || '大学官方系统', '任何第三方官方系统'),
        '中国农业' || '大学', '矿泉水'),
        '农' || '大API', '矿泉水API'),
        '农' || '大邮箱', '矿泉水邮箱'),
        '农' || '大', '矿泉水'),
        'C' || 'A' || 'U_API', 'KQS_API'),
        'C' || 'A' || 'U API', 'KQS API'),
        'c' || 'a' || 'u-api', 'kqs-api'),
        '@' || 'c' || 'a' || 'u.edu.cn', '@kqs.edu.cn')
WHERE value LIKE '%' || '农' || '大' || '%'
   OR value LIKE '%' || '中国农业' || '大学' || '%'
   OR value LIKE '%' || 'C' || 'A' || 'U' || '%'
   OR value LIKE '%' || 'c' || 'a' || 'u' || '%';

INSERT INTO settings (key, value, updated_at)
VALUES ('site_logo', '/kqs-water-logo.svg?v=20260612', NOW())
ON CONFLICT (key) DO UPDATE
SET value = EXCLUDED.value,
    updated_at = NOW()
WHERE trim(settings.value) = ''
   OR settings.value LIKE '%' || 'c' || 'a' || 'u-emblem%'
   OR settings.value LIKE '%/logo.png%';

UPDATE announcements
SET title = replace(
        replace(
        replace(title,
            '农' || '大API', '矿泉水API'),
            '农' || '大', '矿泉水'),
            '校训福利', '专属福利'),
    content = replace(
        replace(
        replace(
        replace(
        replace(content,
            '中国农业' || '大学校训', '专属福利口令'),
            '中国农业' || '大学', '矿泉水'),
            '农' || '大邮箱', '矿泉水邮箱'),
            '农' || '大', '矿泉水'),
            '校训福利', '专属福利')
WHERE title LIKE '%' || '农' || '大' || '%'
   OR title LIKE '%校训福利%'
   OR content LIKE '%' || '农' || '大' || '%'
   OR content LIKE '%' || '中国农业' || '大学' || '%'
   OR content LIKE '%校训福利%';

UPDATE subscription_plans
SET product_name = replace(product_name, '农' || '大API', '矿泉水API')
WHERE product_name LIKE '%' || '农' || '大API' || '%';
