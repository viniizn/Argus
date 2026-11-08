package com.argus.server.enrollment;

import com.argus.server.security.GeneratedSecret;
import com.argus.server.security.SecretGenerator;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.Duration;

@Service
public class EnrollmentTokenService {

    private static final Logger log = LoggerFactory.getLogger(EnrollmentTokenService.class);
    private static final Duration TTL = Duration.ofMinutes(15);

    private final EnrollmentTokenRepository repository;
    private final SecretGenerator secrets;

    public EnrollmentTokenService(EnrollmentTokenRepository repository, SecretGenerator secrets) {
        this.repository = repository;
        this.secrets = secrets;
    }

    @Transactional
    public IssuedToken issue(String createdBy) {
        GeneratedSecret secret = secrets.generate();
        EnrollmentToken token = repository.save(
                EnrollmentToken.issue(secret.hash(), createdBy, TTL));
        log.info("Enrollment token issued tokenId={} createdBy={} expiresAt={}",
                token.getId(), createdBy, token.getExpiresAt());
        return new IssuedToken(secret.raw(), token.getExpiresAt());
    }
}