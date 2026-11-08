package com.argus.server.enrollment;

import java.time.Instant;

public record IssuedToken(String token, Instant expiresAt) {

    @Override
    public String toString() {
        return "IssuedToken[token=<redacted>, expiresAt=" + expiresAt + "]";
    }
}