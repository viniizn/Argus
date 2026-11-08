package com.argus.server.agent;

import jakarta.persistence.*;
import java.time.Instant;
import java.util.UUID;

@Entity
@Table(name = "agents")
public class Agent {

    @Id
    private UUID id;

    @Column(nullable = false)
    private String hostname;

    @Column(nullable = false)
    private String os;

    @Column(nullable = false)
    private String architecture;

    @Column(name = "ip_address", nullable = false)
    private String ipAddress;

    @Column(nullable = false)
    private String version;

    @Enumerated(EnumType.STRING)
    @Column(nullable = false)
    private AgentStatus status;

    @Column(name = "last_seen_at")
    private Instant lastSeenAt;

    @Column(name = "registered_at", nullable = false, updatable = false)
    private Instant registeredAt;

    protected Agent() {}

    public static Agent register(UUID id, String hostname, String os,
                                 String architecture, String ipAddress,
                                 String version) {
        Agent agent = new Agent();
        agent.id = id;
        agent.hostname = hostname;
        agent.os = os;
        agent.architecture = architecture;
        agent.ipAddress = ipAddress;
        agent.version = version;
        agent.status = AgentStatus.ONLINE;
        agent.registeredAt = Instant.now();
        agent.lastSeenAt = Instant.now();
        return agent;
    }

    public void recordHeartbeat(String hostname, String os,
                                String architecture, String ipAddress,
                                String version) {
        this.hostname = hostname;
        this.os = os;
        this.architecture = architecture;
        this.ipAddress = ipAddress;
        this.version = version;
        this.status = AgentStatus.ONLINE;
        this.lastSeenAt = Instant.now();
    }

    public UUID getId() { return id; }
    public String getHostname() { return hostname; }
    public AgentStatus getStatus() { return status; }
    public Instant getLastSeenAt() { return lastSeenAt; }
    public Instant getRegisteredAt() { return registeredAt; }
}