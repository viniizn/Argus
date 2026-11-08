package com.argus.server.enrollment;

import java.util.UUID;

public record EnrollResponse(UUID agentId, String credential) {

    @Override
    public String toString() {
        return "EnrollResponse[agentId=" + agentId + ", credential=<redacted>]";
    }
}