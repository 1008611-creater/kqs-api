ALTER TABLE usage_logs
	ADD COLUMN IF NOT EXISTS upstream_request_bytes BIGINT NOT NULL DEFAULT 0,
	ADD COLUMN IF NOT EXISTS upstream_response_bytes BIGINT NOT NULL DEFAULT 0;

CREATE INDEX IF NOT EXISTS idx_usage_logs_created_account_upstream_traffic
	ON usage_logs (created_at, account_id);
