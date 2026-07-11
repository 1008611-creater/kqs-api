-- One-time cleanup for API keys created with the old frontend defaults.
-- Only exact old defaults are cleared. Non-default/manual limits are preserved.

CREATE TABLE IF NOT EXISTS public.api_key_default_limit_cleanup_20260604 (
  api_key_id bigint PRIMARY KEY,
  old_quota decimal(20,8) NOT NULL,
  old_rate_limit_5h decimal(20,8) NOT NULL,
  old_rate_limit_1d decimal(20,8) NOT NULL,
  old_rate_limit_7d decimal(20,8) NOT NULL,
  backed_up_at timestamptz NOT NULL DEFAULT now()
);

INSERT INTO public.api_key_default_limit_cleanup_20260604 (
  api_key_id,
  old_quota,
  old_rate_limit_5h,
  old_rate_limit_1d,
  old_rate_limit_7d
)
SELECT
  id,
  quota,
  rate_limit_5h,
  rate_limit_1d,
  rate_limit_7d
FROM public.api_keys
WHERE deleted_at IS NULL
  AND (
    abs(quota::numeric - 15.9::numeric) < 0.00000001
    OR (
      abs(rate_limit_5h::numeric - 5::numeric) < 0.00000001
      AND abs(rate_limit_1d::numeric - 15.9::numeric) < 0.00000001
      AND abs(rate_limit_7d::numeric - 34.9::numeric) < 0.00000001
    )
  )
ON CONFLICT (api_key_id) DO NOTHING;

UPDATE public.api_keys
SET
  quota = CASE
    WHEN abs(quota::numeric - 15.9::numeric) < 0.00000001 THEN 0
    ELSE quota
  END,
  rate_limit_5h = CASE
    WHEN abs(rate_limit_5h::numeric - 5::numeric) < 0.00000001
     AND abs(rate_limit_1d::numeric - 15.9::numeric) < 0.00000001
     AND abs(rate_limit_7d::numeric - 34.9::numeric) < 0.00000001 THEN 0
    ELSE rate_limit_5h
  END,
  rate_limit_1d = CASE
    WHEN abs(rate_limit_5h::numeric - 5::numeric) < 0.00000001
     AND abs(rate_limit_1d::numeric - 15.9::numeric) < 0.00000001
     AND abs(rate_limit_7d::numeric - 34.9::numeric) < 0.00000001 THEN 0
    ELSE rate_limit_1d
  END,
  rate_limit_7d = CASE
    WHEN abs(rate_limit_5h::numeric - 5::numeric) < 0.00000001
     AND abs(rate_limit_1d::numeric - 15.9::numeric) < 0.00000001
     AND abs(rate_limit_7d::numeric - 34.9::numeric) < 0.00000001 THEN 0
    ELSE rate_limit_7d
  END,
  updated_at = now()
WHERE deleted_at IS NULL
  AND (
    abs(quota::numeric - 15.9::numeric) < 0.00000001
    OR (
      abs(rate_limit_5h::numeric - 5::numeric) < 0.00000001
      AND abs(rate_limit_1d::numeric - 15.9::numeric) < 0.00000001
      AND abs(rate_limit_7d::numeric - 34.9::numeric) < 0.00000001
    )
  );
