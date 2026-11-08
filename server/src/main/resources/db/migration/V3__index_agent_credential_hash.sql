CREATE UNIQUE INDEX idx_agents_credential_hash
    ON agents(credential_hash)
    WHERE credential_hash IS NOT NULL;