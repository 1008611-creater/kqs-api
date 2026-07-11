-- Remove donation wording from the public support contact text.
--
-- The current public flow sells LDXP card codes and uses in-site redeem.
-- Support copy should describe how to contact the administrator without
-- implying a donation flow.

INSERT INTO settings (key, value)
VALUES ('contact_info', '微信客服：联系管理员时请备注注册邮箱或用户名')
ON CONFLICT (key) DO UPDATE
SET value = EXCLUDED.value
WHERE settings.value IS NULL
   OR settings.value = ''
   OR settings.value ~* '(捐赠|赞助|打赏|收款|捐钱|donat|donate|donation|sponsor|tip)';
