package gm.zona_fit.internal.interfaces.api.config;

import gm.zona_fit.internal.interfaces.api.middleware.AuthMiddleware;
import gm.zona_fit.internal.interfaces.api.middleware.PermissionMiddleware;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.servlet.config.annotation.CorsRegistry;
import org.springframework.web.servlet.config.annotation.InterceptorRegistry;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer;

@Configuration
public class WebMvcConfig implements WebMvcConfigurer {
    
    private final AuthMiddleware authMiddleware;
    private final PermissionMiddleware permissionMiddleware;
    
    public WebMvcConfig(AuthMiddleware authMiddleware, PermissionMiddleware permissionMiddleware) {
        this.authMiddleware = authMiddleware;
        this.permissionMiddleware = permissionMiddleware;
    }
    
    @Override
    public void addInterceptors(InterceptorRegistry registry) {
        // Middleware de permisos primero (depende de auth)
        registry.addInterceptor(permissionMiddleware)
                .addPathPatterns("/api/**")
                .excludePathPatterns("/api/auth/login", "/api/auth/register", "/api/public/**");
        
        // Middleware de auth
        registry.addInterceptor(authMiddleware)
                .addPathPatterns("/api/**")
                .excludePathPatterns("/api/auth/login", "/api/auth/register", "/api/public/**");
    }
    
    @Override
    public void addCorsMappings(CorsRegistry registry) {
        registry.addMapping("/api/**")
                .allowedOrigins("http://localhost:3000") // Ajustar según tu frontend
                .allowedMethods("GET", "POST", "PUT", "DELETE", "OPTIONS")
                .allowedHeaders("*")
                .allowCredentials(true)
                .maxAge(3600);
    }
}