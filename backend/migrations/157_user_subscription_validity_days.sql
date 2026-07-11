DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = current_schema()
          AND table_name = 'user_subscriptions'
          AND column_name = 'validity_days'
    ) THEN
        ALTER TABLE user_subscriptions
            ADD COLUMN validity_days INTEGER NOT NULL DEFAULT 1;

        UPDATE user_subscriptions
        SET validity_days = GREATEST(
            1,
            CEIL(EXTRACT(EPOCH FROM (expires_at - starts_at)) / 86400.0)::integer
        )
        WHERE starts_at IS NOT NULL
          AND expires_at IS NOT NULL;
    END IF;
END $$;
