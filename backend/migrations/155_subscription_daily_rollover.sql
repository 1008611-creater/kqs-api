-- Add daily quota rollover for subscription plans.
-- At the next daily reset, unused daily quota is carried into this column.
ALTER TABLE user_subscriptions
    ADD COLUMN IF NOT EXISTS daily_rollover_usd DECIMAL(20, 10) NOT NULL DEFAULT 0;
