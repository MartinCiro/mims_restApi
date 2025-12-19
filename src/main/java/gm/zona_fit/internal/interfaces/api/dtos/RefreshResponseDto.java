package gm.zona_fit.internal.interfaces.api.dtos;

import com.fasterxml.jackson.annotation.JsonFormat;
import com.fasterxml.jackson.annotation.JsonProperty;

import java.time.LocalDateTime;

public class RefreshResponseDto {
    
    @JsonProperty("refreshed")
    private boolean refreshed;
    
    @JsonProperty("new_cookie")
    private String newCookie;
    
    @JsonProperty("new_expires_at")
    @JsonFormat(shape = JsonFormat.Shape.STRING, pattern = "yyyy-MM-dd'T'HH:mm:ss")
    private LocalDateTime newExpiresAt;
    
    @JsonProperty("ttl_seconds")
    private long ttlSeconds;
    
    @JsonProperty("message")
    private String message;
    
    // Constructor
    public RefreshResponseDto(boolean refreshed, String newCookie, LocalDateTime newExpiresAt, 
                             long ttlSeconds, String message) {
        this.refreshed = refreshed;
        this.newCookie = newCookie;
        this.newExpiresAt = newExpiresAt;
        this.ttlSeconds = ttlSeconds;
        this.message = message;
    }
    
    // Getters
    public boolean isRefreshed() { return refreshed; }
    public String getNewCookie() { return newCookie; }
    public LocalDateTime getNewExpiresAt() { return newExpiresAt; }
    public long getTtlSeconds() { return ttlSeconds; }
    public String getMessage() { return message; }
}