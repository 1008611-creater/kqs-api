-- Keep user_subscriptions.validity_days aligned with the actual term length.
-- Older create paths only set starts_at/expires_at, so multi-day cards could
-- be persisted with the schema default of 1 day.

UPDATE user_subscriptions
SET validity_days = GREATEST(
  1,
  LEAST(
    36500,
    CEIL(EXTRACT(EPOCH FROM (expires_at - starts_at)) / 86400.0)::integer
  )
)
WHERE starts_at IS NOT NULL
  AND expires_at IS NOT NULL
  AND validity_days IS DISTINCT FROM GREATEST(
    1,
    LEAST(
      36500,
      CEIL(EXTRACT(EPOCH FROM (expires_at - starts_at)) / 86400.0)::integer
    )
  );

CREATE OR REPLACE FUNCTION sync_user_subscription_validity_days()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
  IF NEW.starts_at IS NOT NULL AND NEW.expires_at IS NOT NULL THEN
    NEW.validity_days := GREATEST(
      1,
      LEAST(
        36500,
        CEIL(EXTRACT(EPOCH FROM (NEW.expires_at - NEW.starts_at)) / 86400.0)::integer
      )
    );
  END IF;
  RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS trg_sync_user_subscription_validity_days ON user_subscriptions;
CREATE TRIGGER trg_sync_user_subscription_validity_days
BEFORE INSERT OR UPDATE OF starts_at, expires_at
ON user_subscriptions
FOR EACH ROW
EXECUTE FUNCTION sync_user_subscription_validity_days();
