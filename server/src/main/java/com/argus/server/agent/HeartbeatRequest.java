package com.argus.server.agent;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import java.time.Instant;
import java.util.UUID;

public record HeartbeatRequest(
        @NotNull UUID agentId,
        @NotBlank String hostname,
        @NotBlank String os,
        @NotBlank String architecture,
        @NotBlank String ipAddress,
        @NotBlank String agentVersion,
        @NotNull Instant timestamp
) {}