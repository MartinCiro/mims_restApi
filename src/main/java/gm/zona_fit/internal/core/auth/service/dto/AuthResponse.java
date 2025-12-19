package gm.zona_fit.internal.core.auth.service.dto;

import java.time.LocalDateTime;

public class AuthResponse {
    private String message;
    private LocalDateTime expiresAt;

    public AuthResponse(String message, LocalDateTime expiresAt) {
        this.message = message;
        this.expiresAt = expiresAt;
    }

    // Getters
    public String getMessage() {
        return message;
    }

    public LocalDateTime getExpiresAt() {
        return expiresAt;
    }
}