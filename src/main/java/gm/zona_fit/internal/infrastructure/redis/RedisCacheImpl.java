package gm.zona_fit.internal.infrastructure.redis;

import gm.zona_fit.internal.core.auth.service.RedisCache;
import org.springframework.data.redis.core.RedisTemplate;
import org.springframework.stereotype.Component;

import java.time.Duration;
import java.util.HashMap;
import java.util.Map;
import java.util.concurrent.TimeUnit;

@Component
public class RedisCacheImpl implements RedisCache {
    
    private final RedisTemplate<String, Object> redisTemplate;
    
    public RedisCacheImpl(RedisTemplate<String, Object> redisTemplate) {
        this.redisTemplate = redisTemplate;
    }
    
    @Override
    public Duration getTTL(String key) {
        Long ttlSeconds = redisTemplate.getExpire(key, TimeUnit.SECONDS);
        if (ttlSeconds == null || ttlSeconds < 0) {
            return Duration.ZERO;
        }
        return Duration.ofSeconds(ttlSeconds);
    }
    
    @Override
    public void refreshTTL(String key, Duration ttl) {
        redisTemplate.expire(key, ttl.getSeconds(), TimeUnit.SECONDS);
    }
    
    @Override
    public boolean exists(String key) {
        Boolean exists = redisTemplate.hasKey(key);
        return exists != null && exists;
    }
    
    @Override
    public void set(String key, Object value, Duration ttl) {
        if (ttl != null && !ttl.isZero() && !ttl.isNegative()) {
            redisTemplate.opsForValue().set(key, value, ttl.getSeconds(), TimeUnit.SECONDS);
        } else {
            redisTemplate.opsForValue().set(key, value);
        }
    }
    
    @Override
    public Object get(String key) {
        return redisTemplate.opsForValue().get(key);
    }
    
    @Override
    public void delete(String key) {
        redisTemplate.delete(key);
    }
    
    @Override
    public Map<String, Duration> keysWithTTL(String pattern) {
        Map<String, Duration> result = new HashMap<>();
        
        var keys = redisTemplate.keys(pattern);
        if (keys != null) {
            for (String key : keys) {
                Duration ttl = getTTL(key);
                if (!ttl.isZero() && !ttl.isNegative()) {
                    result.put(key, ttl);
                }
            }
        }
        
        return result;
    }
}