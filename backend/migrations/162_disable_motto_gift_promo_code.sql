-- 福利口令（JMSZDJYTXZYC）在 promo_codes 里的那条种子记录必须停用。
-- 原因：注册接口的 ApplyPromoCode 只校验状态/次数/是否用过，不校验邮箱归属，
-- 导致任意邮箱带 promo_code 注册即可白拿 30 美金额度，绕过兑换页的邮箱 + 每用户 + 每 IP 三重校验。
-- 兑换页 /redeem 走的是 RedeemService 的专用分支（isKqsMottoGiftCode），不读 promo_codes，
-- 所以停用这条记录只会堵住注册绕过路径，不影响正常福利领取。
-- 注意：不要回改 144 迁移（checksum 校验会导致启动失败）。

UPDATE promo_codes
SET status = 'disabled',
    updated_at = NOW()
WHERE code = 'JMSZDJYTXZYC';
