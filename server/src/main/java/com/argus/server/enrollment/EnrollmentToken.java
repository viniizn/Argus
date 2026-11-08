package com.argus.server.enrollment;

import jakarta.persistence.*;

import java.time.Duration;
import java.time.Instant;
import java.util.UUID;

@Entity
@Table(name = "enrollment_tokens")
public class EnrollmentToken {

    @Id
    private UUID id;

    @Column(name = "token_hash", nullable = false, updatable = false)
    private String tokenHash;

    @Column(name = "created_by", nullable = false, updatable = false)
    private String createdBy;

    @Column(name = "created_at", nullable = false, updatable = false)
    private Instant createdAt;

    @Column(name = "expires_at", nullable = false, updatable = false)
    private Instant expiresAt;

    @Column(name = "used_at")
    private Instant usedAt;

    @Column(name = "used_by_agent_id")
    private UUID usedByAgentId;

    protected EnrollmentToken() {}

    public static EnrollmentToken issue(String tokenHash, String createdBy, Duration ttl) {
        EnrollmentToken t = new EnrollmentToken();
        t.id = UUID.randomUUID();
        t.tokenHash = tokenHash;
        t.createdBy = createdBy;
        t.createdAt = Instant.now();
        t.expiresAt = t.createdAt.plus(ttl);
        return t;
    }

    public UUID getId() { return id; }
    public Instant getExpiresAt() { return expiresAt; }
}