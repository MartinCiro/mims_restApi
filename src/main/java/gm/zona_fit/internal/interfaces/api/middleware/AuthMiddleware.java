package gm.zona_fit.internal.interfaces.api.middleware;

import gm.zona_fit.internal.core.auth.service.AuthService;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Component;
import org.springframework.web.servlet.HandlerInterceptor;

@Component
public class AuthMiddleware implements HandlerInterceptor {
    
    private static final Logger logger = LoggerFactory.getLogger(AuthMiddleware.class);
    
    private final AuthService authService;
    
    public AuthMiddleware(AuthService authService) {
        this.authService = authService;
    }
    
    @Override
    public boolean preHandle(HttpServletRequest request, HttpServletResponse response, Object handler) 
            throws Exception {
        
        // Excluir endpoints públicos
        String path = request.getRequestURI();
        if (path.startsWith("/api/auth/login") || 
            path.startsWith("/api/auth/register") ||
            path.startsWith("/api/public/")) {
            return true;
        }
        
        // Obtener token de cookie
        String authToken = null;
        var cookies = request.getCookies();
        if (cookies != null) {
            for (var cookie : cookies) {
                if ("auth_token".equals(cookie.getName())) {
                    authToken = cookie.getValue();
                    break;
                }
            }
        }
        
        // Si no hay cookie, verificar header Authorization
        if (authToken == null) {
            authToken = request.getHeader("Authorization");
            if (authToken != null && authToken.startsWith("Bearer ")) {
                authToken = authToken.substring(7);
            }
        }
        
        if (authToken == null || authToken.isBlank()) {
            response.setStatus(HttpServletResponse.SC_UNAUTHORIZED);
            response.setContentType("application/json");
            response.getWriter().write("{\"error\": \"No autenticado\"}");
            return false;
        }
        
        try {
            var user = authService.validateCookie(authToken);
            if (user == null) {
                response.setStatus(HttpServletResponse.SC_UNAUTHORIZED);
                response.setContentType("application/json");
                response.getWriter().write("{\"error\": \"Token inválido o expirado\"}");
                return false;
            }
            
            // Agregar usuario al request para uso en controllers
            request.setAttribute("currentUser", user);
            
            // Actualizar último acceso de sesión
            try {
                Integer userId = Integer.parseInt(user.getId());
                authService.updateUserSessionAccess(userId);
            } catch (NumberFormatException e) {
                logger.warn("User ID inválido para actualizar acceso: {}", user.getId());
            }
            
            return true;
            
        } catch (Exception e) {
            logger.error("Error en middleware de autenticación: {}", e.getMessage(), e);
            response.setStatus(HttpServletResponse.SC_INTERNAL_SERVER_ERROR);
            response.setContentType("application/json");
            response.getWriter().write("{\"error\": \"Error interno del servidor\"}");
            return false;
        }
    }
}