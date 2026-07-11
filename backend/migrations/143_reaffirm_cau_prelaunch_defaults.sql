-- Re-apply CAU API pre-launch gates for environments that already ran an earlier 142.
-- Keep this migration idempotent; legal document body defaults are provided by SettingService.

INSERT INTO settings (key, value, updated_at)
VALUES
  ('login_agreement_enabled', 'true', NOW()),
  ('login_agreement_mode', 'checkbox', NOW()),
  ('login_agreement_updated_at', '2026-05-25', NOW()),
  ('invitation_code_enabled', 'true', NOW())
ON CONFLICT (key) DO UPDATE SET
  value = EXCLUDED.value,
  updated_at = NOW();
