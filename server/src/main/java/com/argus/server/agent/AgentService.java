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
        repository.findById(request.agentId()).ifPresentOrElse(
                agent -> {
                    agent.recordHeartbeat(
                            request.hostname(),
                            request.os(),
                            request.architecture(),
                            request.ipAddress(),
                            request.agentVersion()
                    );
                    log.info("Heartbeat received agentId={} hostname={} status=ONLINE",
                            request.agentId(), request.hostname());
                },
                () -> {
                    Agent agent = Agent.register(
                            request.agentId(),
                            request.hostname(),
                            request.os(),
                            request.architecture(),
                            request.ipAddress(),
                            request.agentVersion()
                    );
                    repository.save(agent);
                    log.info("Agent registered agentId={} hostname={}",
                            request.agentId(), request.hostname());
                }
        );
    }
}