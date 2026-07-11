-- CAU motto launch gift.
-- Public code: JMSZDJYTXZYC, from "解民生之多艰，育天下之英才".
-- Store only IP hashes, never raw IP addresses.

CREATE TABLE IF NOT EXISTS cau_motto_gift_claims (
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

CREATE INDEX IF NOT EXISTS idx_cau_motto_gift_claims_user_id
    ON cau_motto_gift_claims(user_id);
CREATE INDEX IF NOT EXISTS idx_cau_motto_gift_claims_claimed_at
    ON cau_motto_gift_claims(claimed_at);

COMMENT ON TABLE cau_motto_gift_claims IS 'CAU motto launch gift claim records; IP is hashed for privacy.';
COMMENT ON COLUMN cau_motto_gift_claims.campaign_code IS 'Campaign code, e.g. JMSZDJYTXZYC.';
COMMENT ON COLUMN cau_motto_gift_claims.ip_hash IS 'SHA-256 hash of normalized client IP and campaign salt.';

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
VALUES (
    '农大API 校训福利开放领取',
    '### 30 美金额度，先到先得\n\n输入中国农业大学校训“解民生之多艰，育天下之英才”的拼音首字母：`JMSZDJYTXZYC`，即可免费领取 **30 美金** 额度。\n\n领取方式：进入 [兑换码页面](/redeem)，粘贴上面的口令并点击兑换。\n\n每个网络/IP 仅限领取一次，换账号也不能重复领取。后续会逐步改成农大邮箱后缀专属福利。',
    'active',
    'popup',
    '{}'::jsonb,
    NOW(),
    NOW(),
    NOW()
);
