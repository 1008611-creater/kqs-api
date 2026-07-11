-- 151: Persist registration-related attempts for abuse investigation.

CREATE TABLE IF NOT EXISTS registration_audit_events (
    id                BIGSERIAL PRIMARY KEY,
    action            VARCHAR(40) NOT NULL,
    outcome           VARCHAR(40) NOT NULL,
    failure_reason    VARCHAR(120) NOT NULL DEFAULT '',
    user_id           BIGINT REFERENCES users(id) ON DELETE SET NULL,
    email_domain      VARCHAR(255) NOT NULL DEFAULT '',
    email_sha256      CHAR(64) NOT NULL DEFAULT '',
    client_ip         VARCHAR(80) NOT NULL DEFAULT '',
    user_agent        TEXT NOT NULL DEFAULT '',
    request_path      TEXT NOT NULL DEFAULT '',
    turnstile_present BOOLEAN NOT NULL DEFAULT FALSE,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_registration_audit_events_created_at
    ON registration_audit_events(created_at DESC);

CREATE INDEX IF NOT EXISTS idx_registration_audit_events_client_ip
    ON registration_audit_events(client_ip);

CREATE INDEX IF NOT EXISTS idx_registration_audit_events_email_domain
    ON registration_audit_events(email_domain);

CREATE INDEX IF NOT EXISTS idx_registration_audit_events_action_outcome
    ON registration_audit_events(action, outcome);
