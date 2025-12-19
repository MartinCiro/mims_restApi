package gm.zona_fit.internal.core.auth.service;

public class AuthConfig {
    private final int sessionTTL;
    private final int refreshThreshold;
    private final int jwtExpireTime;

    public AuthConfig(int sessionTTL, int refreshThreshold, int jwtExpireTime) {
        this.sessionTTL = sessionTTL;
        this.refreshThreshold = refreshThreshold;
        this.jwtExpireTime = jwtExpireTime;
    }

    public int getSessionTTL() {
        return sessionTTL;
    }

    public int getRefreshThreshold() {
        return refreshThreshold;
    }

    public int getJwtExpireTime() {
        return jwtExpireTime;
    }
}