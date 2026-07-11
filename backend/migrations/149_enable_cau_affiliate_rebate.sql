-- CAU API invite rebate rollout.
-- Keep this idempotent so fresh environments and restored backups converge on
-- the same beta default: invitation links enabled, inviter rebate at 10%.

INSERT INTO settings (key, value, updated_at)
VALUES
  ('invitation_code_enabled', 'true', NOW()),
  ('affiliate_enabled', 'true', NOW()),
  ('affiliate_rebate_rate', '10.00000000', NOW()),
  ('affiliate_rebate_freeze_hours', '0', NOW()),
  ('affiliate_rebate_duration_days', '0', NOW()),
  ('affiliate_rebate_per_invitee_cap', '0.00000000', NOW())
ON CONFLICT (key) DO UPDATE SET
  value = EXCLUDED.value,
  updated_at = NOW();
