package gm.zona_fit.internal.infrastructure.security;

import gm.zona_fit.internal.core.auth.service.CookieSigner;
import gm.zona_fit.internal.core.auth.service.SignedCookieData;
import io.jsonwebtoken.Claims;
import io.jsonwebtoken.Jwts;
import io.jsonwebtoken.security.Keys;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import javax.crypto.SecretKey;
import java.nio.charset.StandardCharsets;
import java.time.LocalDateTime;
import java.time.ZoneId;
import java.util.Date;

@Component
public class CookieSignerImpl implements CookieSigner {
    
    private final SecretKey secretKey;
    private final long expirationSeconds;
    
    public CookieSignerImpl(
            @Value("${app.jwt.secret}") String jwtSecret,
            @Value("${app.jwt.expiration:3600}") long expirationSeconds) {
        
        this.secretKey = Keys.hmacShaKeyFor(jwtSecret.getBytes(StandardCharsets.UTF_8));
        this.expirationSeconds = expirationSeconds;
    }
    
    @Override
    public String sign(SignedCookieData data) {
        Date expirationDate = Date.from(
            data.getExpiresAt().atZone(ZoneId.systemDefault()).toInstant()
        );
        
        Date issuedAtDate = Date.from(
            data.getIssuedAt().atZone(ZoneId.systemDefault()).toInstant()
        );
        
        return Jwts.builder()
                .subject(data.getUserId())
                .claim("username", data.getUsername())
                .claim("role", data.getRole())
                .claim("roleId", data.getRoleId())
                .issuedAt(issuedAtDate)
                .expiration(expirationDate)
                .signWith(secretKey)
                .compact();
    }
    
    @Override
    public SignedCookieData verify(String signedCookie) {
        try {
            Claims claims = Jwts.parser()
                    .verifyWith(secretKey)
                    .build()
                    .parseSignedClaims(signedCookie)
                    .getPayload();
            
            LocalDateTime expiresAt = claims.getExpiration()
                    .toInstant()
                    .atZone(ZoneId.systemDefault())
                    .toLocalDateTime();
            
            LocalDateTime issuedAt = claims.getIssuedAt()
                    .toInstant()
                    .atZone(ZoneId.systemDefault())
                    .toLocalDateTime();
            
            return new SignedCookieData(
                claims.getSubject(),
                claims.get("username", String.class),
                claims.get("role", String.class),
                claims.get("roleId", Integer.class),
                expiresAt,
                issuedAt
            );
        } catch (Exception e) {
            throw new IllegalArgumentException("Cookie inválida: " + e.getMessage(), e);
        }
    }
}