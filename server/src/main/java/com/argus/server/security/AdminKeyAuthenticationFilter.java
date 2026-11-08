package com.argus.server.security;

import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.web.filter.OncePerRequestFilter;

import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.util.List;

/**
 * TEMPORARY: authenticates admin requests with a static key from the environment.
 * To be replaced by real user authentication + RBAC.
 */
public class AdminKeyAuthenticationFilter extends OncePerRequestFilter {

    private static final String HEADER = "X-Admin-Key";
    private static final String ADMIN_PATH = "/api/v1/admin/";

    private final byte[] expectedDigest; // null => admin API disabled

    public AdminKeyAuthenticationFilter(String bootstrapKey) {
        this.expectedDigest = (bootstrapKey == null || bootstrapKey.isBlank())
                ? null
                : sha256(bootstrapKey);
    }

    @Override
    protected boolean shouldNotFilter(HttpServletRequest request) {
        return !request.getRequestURI().startsWith(ADMIN_PATH);
    }

    @Override
    protected void doFilterInternal(HttpServletRequest request,
                                    HttpServletResponse response,
                                    FilterChain chain) throws ServletException, IOException {
        String provided = request.getHeader(HEADER);
        if (expectedDigest != null && provided != null
                && MessageDigest.isEqual(sha256(provided), expectedDigest)) {
            SecurityContextHolder.getContext().setAuthentication(
                    UsernamePasswordAuthenticationToken.authenticated(
                            "bootstrap-admin", null,
                            List.of(new SimpleGrantedAuthority("ROLE_ADMIN"))));
        }
        chain.doFilter(request, response);
    }

    private static byte[] sha256(String value) {
        try {
            return MessageDigest.getInstance("SHA-256")
                    .digest(value.getBytes(StandardCharsets.UTF_8));
        } catch (NoSuchAlgorithmException e) {
            throw new IllegalStateException("SHA-256 not available", e);
        }
    }
}