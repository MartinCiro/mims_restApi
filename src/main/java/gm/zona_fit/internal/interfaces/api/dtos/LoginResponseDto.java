package gm.zona_fit.internal.interfaces.api.dtos;

import com.fasterxml.jackson.annotation.JsonFormat;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.time.LocalDateTime;

public class LoginResponseDto {
    
    @JsonProperty("message")
    private String message;
    
    @JsonProperty("cookie")
    private String cookie;
    
    @JsonProperty("expires_at")
    @JsonFormat(shape = JsonFormat.Shape.STRING, pattern = "yyyy-MM-dd'T'HH:mm:ss")
    private LocalDateTime expiresAt;
    
    // Constructor
    public LoginResponseDto(String message, String cookie, LocalDateTime expiresAt) {
        this.message = message;
        this.cookie = cookie;
        this.expiresAt = expiresAt;
    }
    
    // Getters
    public String getMessage() { return message; }
    public String getCookie() { return cookie; }
    public LocalDateTime getExpiresAt() { return expiresAt; }
}