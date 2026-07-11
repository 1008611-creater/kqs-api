-- Allow users to hold multiple independent entitlements for the same
-- subscription group. This keeps duplicate purchases from extending time.
DROP INDEX CONCURRENTLY IF EXISTS user_subscriptions_user_group_unique_active;

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_user_subscriptions_user_group_status_expires_active
    ON user_subscriptions(user_id, group_id, status, expires_at)
    WHERE deleted_at IS NULL;
