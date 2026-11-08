package com.argus.server.agent;

import jakarta.validation.Valid;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.Authentication;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/v1/agents")
public class AgentController {

    private final AgentService service;

    public AgentController(AgentService service) {
        this.service = service;
    }

    @PostMapping("/heartbeat")
    public ResponseEntity<Void> heartbeat(@Valid @RequestBody HeartbeatRequest request,
                                          Authentication authentication) {
        if (!request.agentId().equals(authentication.getPrincipal())) {
            return ResponseEntity.status(HttpStatus.FORBIDDEN).build();
        }
        service.processHeartbeat(request);
        return ResponseEntity.ok().build();
    }
}