package com.argus.server.security;

import org.springframework.stereotype.Component;

import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.security.SecureRandom;
import java.util.Base64;
import java.util.HexFormat;

@Component
public class SecretGenerator {

    private static final int SECRET_BYTES = 32;

    private final SecureRandom random = new SecureRandom();

    public GeneratedSecret generate() {
        byte[] bytes = new byte[SECRET_BYTES];
        random.nextBytes(bytes);
        String raw = Base64.getUrlEncoder().withoutPadding().encodeToString(bytes);
        return new GeneratedSecret(raw, hash(raw));
    }

    public String hash(String raw) {
        try {
            MessageDigest digest = MessageDigest.getInstance("SHA-256");
            byte[] hashed = digest.digest(raw.getBytes(StandardCharsets.UTF_8));
            return HexFormat.of().formatHex(hashed);
        } catch (NoSuchAlgorithmException e) {
            // SHA-256 is required by the Java spec; this cannot happen.
            throw new IllegalStateException("SHA-256 not available", e);
        }
    }

    public boolean matches(String raw, String storedHash) {
        if (raw == null || storedHash == null) {
            return false;
        }
        byte[] computed = hash(raw).getBytes(StandardCharsets.UTF_8);
        byte[] stored = storedHash.getBytes(StandardCharsets.UTF_8);
        return MessageDigest.isEqual(computed, stored);
    }
}