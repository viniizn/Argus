package com.argus.server.security;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class SecretGeneratorTest {

    private final SecretGenerator generator = new SecretGenerator();

    @Test
    void generatesSecretsOfExpectedShape() {
        GeneratedSecret secret = generator.generate();

        assertEquals(43, secret.raw().length());   // 32 bytes -> 43 base64url chars
        assertEquals(64, secret.hash().length());  // sha-256 -> 64 hex chars
        assertTrue(secret.raw().matches("[A-Za-z0-9_-]+"));
        assertTrue(secret.hash().matches("[0-9a-f]+"));
    }

    @Test
    void generatesDifferentSecretsEachTime() {
        assertNotEquals(generator.generate().raw(), generator.generate().raw());
    }

    @Test
    void hashIsDeterministic() {
        assertEquals(generator.hash("abc"), generator.hash("abc"));
    }

    @Test
    void hashMatchesKnownSha256Vector() {
        // Well-known SHA-256 test vector for "abc".
        assertEquals(
                "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad",
                generator.hash("abc"));
    }

    @Test
    void matchesAcceptsCorrectSecret() {
        GeneratedSecret secret = generator.generate();
        assertTrue(generator.matches(secret.raw(), secret.hash()));
    }

    @Test
    void matchesRejectsWrongSecret() {
        GeneratedSecret secret = generator.generate();
        assertFalse(generator.matches("wrong-secret", secret.hash()));
    }

    @Test
    void matchesRejectsNulls() {
        assertFalse(generator.matches(null, "x"));
        assertFalse(generator.matches("x", null));
    }

    @Test
    void toStringDoesNotLeakRawSecret() {
        GeneratedSecret secret = generator.generate();
        assertFalse(secret.toString().contains(secret.raw()));
    }
}