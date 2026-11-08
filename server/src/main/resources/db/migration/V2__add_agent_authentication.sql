ALTER TABLE agents
    ADD COLUMN credential_hash VARCHAR(64),
    ADD COLUMN credential_issued_at TIMESTAMPTZ;

CREATE TABLE enrollment_tokens (
    id                  UUID PRIMARY KEY,
    token_hash          VARCHAR(64)  NOT NULL UNIQUE,
    created_by          VARCHAR(255) NOT NULL,
    created_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    expires_at          TIMESTAMPTZ  NOT NULL,
    used_at             TIMESTAMPTZ,
    used_by_agent_id    UUID,
    CONSTRAINT chk_enrollment_expiry CHECK (expires_at > created_at)
);

CREATE INDEX idx_enrollment_tokens_expires_at ON enrollment_tokens(expires_at);