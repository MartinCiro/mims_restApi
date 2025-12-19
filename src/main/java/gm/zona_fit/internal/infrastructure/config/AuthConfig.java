package gm.zona_fit.internal.infrastructure.config;

import gm.zona_fit.internal.core.auth.service.AuthConfig as CoreAuthConfig;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class AuthConfig {
    
    @Value("${app.auth.session-ttl:7200}") // 2 horas por defecto
    private int sessionTTL;
    
    @Value("${app.auth.refresh-threshold:300}") // 5 minutos por defecto
    private int refreshThreshold;
    
    @Value("${app.auth.jwt-expire-time:3600}") // 1 hora por defecto
    private int jwtExpireTime;
    
    @Bean
    public CoreAuthConfig coreAuthConfig() {
        return new CoreAuthConfig(sessionTTL, refreshThreshold, jwtExpireTime);
    }
}