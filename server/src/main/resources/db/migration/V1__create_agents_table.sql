CREATE TABLE agents (
    id          UUID PRIMARY KEY,
    hostname    VARCHAR(255) NOT NULL,
    os          VARCHAR(50)  NOT NULL,
    architecture VARCHAR(50) NOT NULL,
    ip_address  VARCHAR(45)  NOT NULL,
    version     VARCHAR(20)  NOT NULL,
    status      VARCHAR(20)  NOT NULL DEFAULT 'OFFLINE',
    last_seen_at TIMESTAMPTZ,
    registered_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_agents_status ON agents(status);