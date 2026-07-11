-- Normalize old frontend default API key limits to unlimited during creation.
-- 0 means unlimited. This trigger intentionally runs only on INSERT, so
-- manually edited limits on existing keys are preserved.

CREATE OR REPLACE FUNCTION public.api_keys_new_default_limits_to_unlimited()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
  IF NEW.quota IS NOT NULL AND abs(NEW.quota::numeric - 15.9::numeric) < 0.00000001 THEN
    NEW.quota := 0;
  END IF;

  IF NEW.rate_limit_5h IS NOT NULL
     AND NEW.rate_limit_1d IS NOT NULL
     AND NEW.rate_limit_7d IS NOT NULL
     AND abs(NEW.rate_limit_5h::numeric - 5::numeric) < 0.00000001
     AND abs(NEW.rate_limit_1d::numeric - 15.9::numeric) < 0.00000001
     AND abs(NEW.rate_limit_7d::numeric - 34.9::numeric) < 0.00000001 THEN
    NEW.rate_limit_5h := 0;
    NEW.rate_limit_1d := 0;
    NEW.rate_limit_7d := 0;
  END IF;

  RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS api_keys_new_default_limits_to_unlimited_trg ON public.api_keys;

CREATE TRIGGER api_keys_new_default_limits_to_unlimited_trg
BEFORE INSERT ON public.api_keys
FOR EACH ROW
EXECUTE FUNCTION public.api_keys_new_default_limits_to_unlimited();
