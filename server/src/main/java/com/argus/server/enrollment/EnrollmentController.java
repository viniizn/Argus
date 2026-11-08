package com.argus.server.enrollment;

import jakarta.validation.Valid;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api/v1/agents/enroll")
public class EnrollmentController {

    private final EnrollmentService service;

    public EnrollmentController(EnrollmentService service) {
        this.service = service;
    }

    @PostMapping
    public ResponseEntity<EnrollResponse> enroll(@Valid @RequestBody EnrollRequest request) {
        EnrollmentService.Result result = service.enroll(request);
        return switch (result.outcome()) {
            case ENROLLED -> ResponseEntity.status(HttpStatus.CREATED).body(result.response());
            case INVALID_TOKEN -> ResponseEntity.status(HttpStatus.UNAUTHORIZED).build();
            case ALREADY_ENROLLED -> ResponseEntity.status(HttpStatus.CONFLICT).build();
        };
    }
}