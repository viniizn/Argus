package com.argus.server.security;

import com.argus.server.agent.Agent;
import com.argus.server.agent.AgentRepository;
import com.argus.server.agent.AgentStatus;
import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.web.filter.OncePerRequestFilter;

import java.io.IOException;
import java.util.List;
import java.util.Optional;
import java.util.UUID;

/** Authenticates agent requests via "Authorization: Bearer <credential>". */
public class AgentCredentialFilter extends OncePerRequestFilter {

    private static final String PREFIX = "Bearer ";
    private static final String PROTECTED_PATH = "/api/v1/agents/heartbeat";
    private static final int MAX_CREDENTIAL_LENGTH = 128;

    private final AgentRepository agents;
    private final SecretGenerator secrets;

    public AgentCredentialFilter(AgentRepository agents, SecretGenerator secrets) {
        this.agents = agents;
        this.secrets = secrets;
    }

    @Override
    protected boolean shouldNotFilter(HttpServletRequest request) {
        return !request.getRequestURI().equals(PROTECTED_PATH);
    }

    @Override
    protected void doFilterInternal(HttpServletRequest request,
                                    HttpServletResponse response,
                                    FilterChain chain) throws ServletException, IOException {
        String header = request.getHeader("Authorization");
        if (header != null && header.startsWith(PREFIX)) {
            String credential = header.substring(PREFIX.length()).trim();
            if (!credential.isEmpty() && credential.length() <= MAX_CREDENTIAL_LENGTH) {
                Optional<Agent> agent = agents.findByCredentialHash(secrets.hash(credential));
                if (agent.isPresent() && agent.get().getStatus() != AgentStatus.REVOKED) {
                    UUID agentId = agent.get().getId();
                    SecurityContextHolder.getContext().setAuthentication(
                            UsernamePasswordAuthenticationToken.authenticated(
                                    agentId, null,
                                    List.of(new SimpleGrantedAuthority("ROLE_AGENT"))));
                }
            }
        }
        chain.doFilter(request, response);
    }
}