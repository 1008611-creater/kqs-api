-- Subscription group-buy hall and per-user quota upgrade overrides.

ALTER TABLE user_subscriptions
  ADD COLUMN IF NOT EXISTS daily_limit_override_usd DECIMAL(20,8),
  ADD COLUMN IF NOT EXISTS weekly_limit_override_usd DECIMAL(20,8),
  ADD COLUMN IF NOT EXISTS monthly_limit_override_usd DECIMAL(20,8),
  ADD COLUMN IF NOT EXISTS quota_bonus_multiplier DECIMAL(10,4) NOT NULL DEFAULT 1,
  ADD COLUMN IF NOT EXISTS quota_bonus_source TEXT;

CREATE INDEX IF NOT EXISTS idx_user_subscriptions_quota_bonus_source
  ON user_subscriptions (quota_bonus_source)
  WHERE deleted_at IS NULL AND quota_bonus_source IS NOT NULL;

CREATE TABLE IF NOT EXISTS group_buy_rooms (
  id BIGSERIAL PRIMARY KEY,
  plan_id BIGINT NOT NULL REFERENCES subscription_plans(id) ON DELETE CASCADE,
  group_id BIGINT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
  owner_user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  target_count INT NOT NULL,
  bonus_multiplier DECIMAL(10,4) NOT NULL DEFAULT 1,
  status VARCHAR(20) NOT NULL DEFAULT 'active',
  expires_at TIMESTAMPTZ NOT NULL,
  completed_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,
  CONSTRAINT group_buy_rooms_target_count_check CHECK (target_count IN (3, 5, 10)),
  CONSTRAINT group_buy_rooms_bonus_multiplier_check CHECK (bonus_multiplier >= 1),
  CONSTRAINT group_buy_rooms_status_check CHECK (status IN ('active', 'completed', 'expired', 'canceled'))
);

CREATE INDEX IF NOT EXISTS idx_group_buy_rooms_status_expires
  ON group_buy_rooms (status, expires_at)
  WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_group_buy_rooms_plan_id
  ON group_buy_rooms (plan_id)
  WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS group_buy_room_members (
  id BIGSERIAL PRIMARY KEY,
  room_id BIGINT NOT NULL REFERENCES group_buy_rooms(id) ON DELETE CASCADE,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  user_subscription_id BIGINT NOT NULL REFERENCES user_subscriptions(id) ON DELETE CASCADE,
  joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  daily_limit_before DECIMAL(20,8),
  daily_limit_after DECIMAL(20,8),
  weekly_limit_before DECIMAL(20,8),
  weekly_limit_after DECIMAL(20,8),
  quota_applied_at TIMESTAMPTZ,
  UNIQUE (room_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_group_buy_room_members_user_id
  ON group_buy_room_members (user_id);

CREATE INDEX IF NOT EXISTS idx_group_buy_room_members_subscription_id
  ON group_buy_room_members (user_subscription_id);
