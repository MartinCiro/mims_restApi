package gm.zona_fit.internal.infrastructure.config;

import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class AuthConfig {

    @Value("${app.auth.session-ttl:7200}")
    private int sessionTTL;

    @Value("${app.auth.refresh-threshold:300}")
    private int refreshThreshold;

    @Value("${app.auth.jwt-expire-time:3600}")
    private int jwtExpireTime;

    @Bean
    public gm.zona_fit.internal.core.auth.service.AuthConfig coreAuthConfig() {
        return new gm.zona_fit.internal.core.auth.service.AuthConfig(
                sessionTTL, refreshThreshold, jwtExpireTime);
    }
}