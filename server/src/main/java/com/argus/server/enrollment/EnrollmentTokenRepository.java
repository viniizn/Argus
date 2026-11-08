package com.argus.server.enrollment;

import org.springframework.data.jpa.repository.JpaRepository;

import java.util.UUID;

public interface EnrollmentTokenRepository extends JpaRepository<EnrollmentToken, UUID> {}