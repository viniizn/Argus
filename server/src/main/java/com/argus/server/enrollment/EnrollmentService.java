package com.argus.server.enrollment;

import com.argus.server.agent.Agent;
import com.argus.server.agent.AgentRepository;
import com.argus.server.security.GeneratedSecret;
import com.argus.server.security.SecretGenerator;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.Instant;
import java.util.Optional;

@Service
public class EnrollmentService {

    private static final Logger log = LoggerFactory.getLogger(EnrollmentService.class);

    public enum Outcome { ENROLLED, INVALID_TOKEN, ALREADY_ENROLLED }

    public record Result(Outcome outcome, EnrollResponse response) {}

    private final EnrollmentTokenRepository tokens;
    private final AgentRepository agents;
    private final SecretGenerator secrets;

    public EnrollmentService(EnrollmentTokenRepository tokens, AgentRepository agents,
                             SecretGenerator secrets) {
        this.tokens = tokens;
        this.agents = agents;
        this.secrets = secrets;
    }

    @Transactional
    public Result enroll(EnrollRequest request) {
        // Check before consuming so a conflict does not burn a valid token.
        Optional<Agent> existing = agents.findById(request.agentId());
        if (existing.isPresent() && existing.get().hasCredential()) {
            log.warn("Enroll rejected: agent already enrolled agentId={}", request.agentId());
            return new Result(Outcome.ALREADY_ENROLLED, null);
        }

        int consumed = tokens.consume(secrets.hash(request.token()), request.agentId(), Instant.now());
        if (consumed == 0) {
            log.warn("Enroll rejected: invalid, expired or used token agentId={}", request.agentId());
            return new Result(Outcome.INVALID_TOKEN, null);
        }

        Agent agent = existing.orElseGet(() -> Agent.register(
                request.agentId(), request.hostname(), request.os(),
                request.architecture(), request.ipAddress(), request.agentVersion()));

        GeneratedSecret credential = secrets.generate();
        agent.issueCredential(credential.hash());
        agents.save(agent);

        log.info("Agent enrolled agentId={} hostname={}", request.agentId(), request.hostname());
        return new Result(Outcome.ENROLLED, new EnrollResponse(request.agentId(), credential.raw()));
    }
}