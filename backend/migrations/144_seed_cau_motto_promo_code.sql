-- CAU-only campus promo: pinyin initials of the motto
-- "解民生之多艰，育天下之英才" => JMSZDJYTXZYC

INSERT INTO settings (key, value, updated_at)
VALUES ('promo_code_enabled', 'true', NOW())
ON CONFLICT (key) DO UPDATE SET
  value = EXCLUDED.value,
  updated_at = NOW();

INSERT INTO promo_codes (
  code,
  bonus_amount,
  max_uses,
  status,
  expires_at,
  notes,
  created_at,
  updated_at
)
VALUES (
  'JMSZDJYTXZYC',
  30,
  0,
  'active',
  NULL,
  'CAU motto initials launch promo; $30 signup balance; one use per account via promo_code_usages.',
  NOW(),
  NOW()
)
ON CONFLICT (code) DO UPDATE SET
  bonus_amount = EXCLUDED.bonus_amount,
  max_uses = EXCLUDED.max_uses,
  status = EXCLUDED.status,
  expires_at = EXCLUDED.expires_at,
  notes = EXCLUDED.notes,
  updated_at = NOW();
