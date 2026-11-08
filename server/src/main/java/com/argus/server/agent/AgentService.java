package com.argus.server.agent;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
public class AgentService {

    private static final Logger log = LoggerFactory.getLogger(AgentService.class);

    private final AgentRepository repository;

    public AgentService(AgentRepository repository) {
        this.repository = repository;
    }

    @Transactional
    public void processHeartbeat(HeartbeatRequest request) {
        Agent agent = repository.findById(request.agentId())
                .orElseThrow(() -> new IllegalStateException("Authenticated agent not found"));
        agent.recordHeartbeat(request.hostname(), request.os(), request.architecture(),
                request.ipAddress(), request.agentVersion());
        log.info("Heartbeat received agentId={} hostname={}", request.agentId(), request.hostname());
    }
}