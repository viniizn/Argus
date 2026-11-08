package com.argus.server.enrollment;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Modifying;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;

import java.time.Instant;
import java.util.UUID;

public interface EnrollmentTokenRepository extends JpaRepository<EnrollmentToken, UUID> {

    @Modifying(clearAutomatically = true, flushAutomatically = true)
    @Query("""
            UPDATE EnrollmentToken t
               SET t.usedAt = :now, t.usedByAgentId = :agentId
             WHERE t.tokenHash = :hash
               AND t.usedAt IS NULL
               AND t.expiresAt > :now
            """)
    int consume(@Param("hash") String hash,
                @Param("agentId") UUID agentId,
                @Param("now") Instant now);
}