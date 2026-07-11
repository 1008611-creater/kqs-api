UPDATE announcements
SET content = $$输入农大校训首字母，即可免费领取 30 美金额度，先到先得。$$,
    updated_at = NOW()
WHERE title = '农大API 校训福利开放领取';
