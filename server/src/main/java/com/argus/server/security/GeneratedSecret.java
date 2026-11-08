package com.argus.server.security;

public record GeneratedSecret(String raw, String hash) {

    @Override
    public String toString() {
        return "GeneratedSecret[raw=<redacted>, hash=" + hash + "]";
    }
}
