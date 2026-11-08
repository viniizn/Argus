package com.argus.server.enrollment;

import org.springframework.http.HttpStatus;
import org.springframework.security.core.Authentication;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api/v1/admin/enrollment-tokens")
public class EnrollmentTokenAdminController {

    private final EnrollmentTokenService service;

    public EnrollmentTokenAdminController(EnrollmentTokenService service) {
        this.service = service;
    }

    @PostMapping
    @ResponseStatus(HttpStatus.CREATED)
    public IssuedToken create(Authentication authentication) {
        return service.issue(authentication.getName());
    }
}