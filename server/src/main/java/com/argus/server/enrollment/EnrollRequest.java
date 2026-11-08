package com.argus.server.enrollment;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Size;

import java.util.UUID;

public record EnrollRequest(
        @NotBlank @Size(max = 128) String token,
        @NotNull UUID agentId,
        @NotBlank @Size(max = 255) String hostname,
        @NotBlank @Size(max = 50) String os,
        @NotBlank @Size(max = 50) String architecture,
        @NotBlank @Size(max = 45) String ipAddress,
        @NotBlank @Size(max = 20) String agentVersion
) {
    @Override
    public String toString() {
        return "EnrollRequest[agentId=" + agentId + ", hostname=" + hostname + ", token=<redacted>]";
    }
}