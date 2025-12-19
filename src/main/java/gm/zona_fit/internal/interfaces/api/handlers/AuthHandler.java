package gm.zona_fit.internal.interfaces.api.handlers;

import gm.zona_fit.internal.core.auth.service.AuthService;
import gm.zona_fit.internal.core.auth.service.RefreshResult;
import gm.zona_fit.internal.interfaces.api.adapters.AuthAdapter;
import gm.zona_fit.internal.interfaces.api.dtos.*;
import jakarta.servlet.http.Cookie;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.time.LocalDateTime;
import java.util.Map;

@RestController
@RequestMapping("/api/auth")
public class AuthHandler {
    
    private static final Logger logger = LoggerFactory.getLogger(AuthHandler.class);
    
    private final AuthService authService;
    
    public AuthHandler(AuthService authService) {
        this.authService = authService;
    }
    
    @PostMapping("/login")
    public ResponseEntity<ApiResponseDto<LoginResponseDto>> login(
            @RequestBody LoginRequestDto requestDto,
            HttpServletResponse response) {
        
        try {
            logger.info("Intento de login para email: {}", requestDto.getEmail());
            
            var domainRequest = AuthAdapter.toDomain(requestDto);
            var authResponse = authService.loginUser(domainRequest);
            
            if (authResponse == null) {
                return ResponseEntity
                    .status(HttpStatus.UNAUTHORIZED)
                    .body(ApiResponseDto.error(401, "Credenciales inválidas"));
            }
            
            // Crear cookie HTTP
            Cookie cookie = new Cookie("auth_token", authResponse.getCookie());
            cookie.setHttpOnly(true);
            cookie.setSecure(true); // En producción usar true
            cookie.setPath("/");
            cookie.setMaxAge((int) authResponse.getExpiresAt().until(LocalDateTime.now()).getSeconds());
            
            response.addCookie(cookie);
            
            var responseDto = new LoginResponseDto(
                authResponse.getMessage(),
                authResponse.getCookie(),
                authResponse.getExpiresAt()
            );
            
            logger.info("Login exitoso para email: {}", requestDto.getEmail());
            
            return ResponseEntity.ok(ApiResponseDto.success(responseDto));
            
        } catch (Exception e) {
            logger.error("Error en login: {}", e.getMessage(), e);
            return ResponseEntity
                .status(HttpStatus.INTERNAL_SERVER_ERROR)
                .body(ApiResponseDto.error(500, "Error interno del servidor"));
        }
    }
    
    @GetMapping("/profile")
    public ResponseEntity<ApiResponseDto<ProfileResponseDto>> getProfile(
            @CookieValue(name = "auth_token", required = false) String authToken) {
        
        try {
            if (authToken == null || authToken.isBlank()) {
                return ResponseEntity
                    .status(HttpStatus.UNAUTHORIZED)
                    .body(ApiResponseDto.error(401, "No autenticado"));
            }
            
            var user = authService.validateCookie(authToken);
            if (user == null) {
                return ResponseEntity
                    .status(HttpStatus.UNAUTHORIZED)
                    .body(ApiResponseDto.error(401, "Token inválido o expirado"));
            }
            
            var profileResponse = authService.getUserProfile(user.getId());
            
            var responseDto = new ProfileResponseDto(
                profileResponse.getId(),
                profileResponse.getNombre(),
                profileResponse.getRol(),
                profileResponse.getPermisos()
            );
            
            return ResponseEntity.ok(ApiResponseDto.success(responseDto));
            
        } catch (Exception e) {
            logger.error("Error obteniendo perfil: {}", e.getMessage(), e);
            return ResponseEntity
                .status(HttpStatus.INTERNAL_SERVER_ERROR)
                .body(ApiResponseDto.error(500, "Error interno del servidor"));
        }
    }
    
    @PostMapping("/register")
    public ResponseEntity<ApiResponseDto<LoginResponseDto>> register(
            @RequestBody RegisterRequestDto requestDto,
            @CookieValue(name = "auth_token", required = false) String authToken,
            HttpServletResponse response) {
        
        try {
            logger.info("Intento de registro para usuario: {}", requestDto.getUsername());
            
            // Validar si el usuario actual tiene permisos (si está autenticado)
            var currentUser = authToken != null ? authService.validateCookie(authToken) : null;
            
            var domainRequest = AuthAdapter.toDomain(requestDto);
            var result = authService.registerUser(domainRequest, currentUser);
            
            // El resultado puede ser String (mensaje) o AuthResponse
            if (result instanceof String) {
                // Registro exitoso sin autologin
                return ResponseEntity.ok(ApiResponseDto.success(null, (String) result));
            } else if (result instanceof gm.zona_fit.internal.core.auth.service.dto.AuthResponse) {
                var authResponse = (gm.zona_fit.internal.core.auth.service.dto.AuthResponse) result;
                
                // Crear cookie HTTP
                Cookie cookie = new Cookie("auth_token", authResponse.getCookie());
                cookie.setHttpOnly(true);
                cookie.setSecure(true);
                cookie.setPath("/");
                cookie.setMaxAge((int) authResponse.getExpiresAt().until(LocalDateTime.now()).getSeconds());
                
                response.addCookie(cookie);
                
                var responseDto = new LoginResponseDto(
                    authResponse.getMessage(),
                    authResponse.getCookie(),
                    authResponse.getExpiresAt()
                );
                
                logger.info("Registro exitoso para usuario: {}", requestDto.getUsername());
                return ResponseEntity.ok(ApiResponseDto.success(responseDto));
            }
            
            return ResponseEntity
                .status(HttpStatus.BAD_REQUEST)
                .body(ApiResponseDto.error(400, "Solicitud inválida"));
            
        } catch (IllegalArgumentException e) {
            logger.warn("Error de validación en registro: {}", e.getMessage());
            return ResponseEntity
                .status(HttpStatus.BAD_REQUEST)
                .body(ApiResponseDto.error(400, e.getMessage()));
        } catch (Exception e) {
            logger.error("Error en registro: {}", e.getMessage(), e);
            return ResponseEntity
                .status(HttpStatus.INTERNAL_SERVER_ERROR)
                .body(ApiResponseDto.error(500, "Error interno del servidor"));
        }
    }
    
    @PostMapping("/logout")
    public ResponseEntity<ApiResponseDto<Void>> logout(
            @CookieValue(name = "auth_token", required = false) String authToken,
            HttpServletRequest request,
            HttpServletResponse response) {
        
        try {
            if (authToken == null || authToken.isBlank()) {
                return ResponseEntity.ok(ApiResponseDto.success(null, "No hay sesión activa"));
            }
            
            var user = authService.validateCookie(authToken);
            if (user != null) {
                try {
                    Integer userId = Integer.parseInt(user.getId());
                    authService.logout(userId);
                } catch (NumberFormatException e) {
                    logger.warn("User ID inválido en logout: {}", user.getId());
                }
            }
            
            // Eliminar cookie
            Cookie cookie = new Cookie("auth_token", null);
            cookie.setHttpOnly(true);
            cookie.setSecure(true);
            cookie.setPath("/");
            cookie.setMaxAge(0);
            
            response.addCookie(cookie);
            
            logger.info("Logout exitoso");
            return ResponseEntity.ok(ApiResponseDto.success(null, "Sesión cerrada exitosamente"));
            
        } catch (Exception e) {
            logger.error("Error en logout: {}", e.getMessage(), e);
            return ResponseEntity
                .status(HttpStatus.INTERNAL_SERVER_ERROR)
                .body(ApiResponseDto.error(500, "Error interno del servidor"));
        }
    }
    
    @PostMapping("/refresh")
    public ResponseEntity<ApiResponseDto<RefreshResponseDto>> refresh(
            @CookieValue(name = "auth_token", required = false) String authToken,
            HttpServletResponse response) {
        
        try {
            if (authToken == null || authToken.isBlank()) {
                return ResponseEntity
                    .status(HttpStatus.UNAUTHORIZED)
                    .body(ApiResponseDto.error(401, "No autenticado"));
            }
            
            var refreshResult = authService.checkAndRefreshSessionFromCookie(authToken);
            
            if (refreshResult.isRefreshed() && refreshResult.getNewCookie() != null) {
                // Actualizar cookie
                Cookie cookie = new Cookie("auth_token", refreshResult.getNewCookie());
                cookie.setHttpOnly(true);
                cookie.setSecure(true);
                cookie.setPath("/");
                cookie.setMaxAge((int) refreshResult.getNewExpiresAt().until(LocalDateTime.now()).getSeconds());
                
                response.addCookie(cookie);
                
                var responseDto = new RefreshResponseDto(
                    true,
                    refreshResult.getNewCookie(),
                    refreshResult.getNewExpiresAt(),
                    refreshResult.getTtlSeconds(),
                    refreshResult.getMessage()
                );
                
                return ResponseEntity.ok(ApiResponseDto.success(responseDto));
            }
            
            var responseDto = new RefreshResponseDto(
                false,
                null,
                null,
                refreshResult.getTtlSeconds(),
                refreshResult.getMessage()
            );
            
            return ResponseEntity.ok(ApiResponseDto.success(responseDto));
            
        } catch (Exception e) {
            logger.error("Error en refresh: {}", e.getMessage(), e);
            return ResponseEntity
                .status(HttpStatus.INTERNAL_SERVER_ERROR)
                .body(ApiResponseDto.error(500, "Error interno del servidor"));
        }
    }
    
    @GetMapping("/session-info")
    public ResponseEntity<ApiResponseDto<Map<String, Object>>> getSessionInfo(
            @CookieValue(name = "auth_token", required = false) String authToken) {
        
        try {
            if (authToken == null || authToken.isBlank()) {
                return ResponseEntity
                    .status(HttpStatus.UNAUTHORIZED)
                    .body(ApiResponseDto.error(401, "No autenticado"));
            }
            
            var user = authService.validateCookie(authToken);
            if (user == null) {
                return ResponseEntity
                    .status(HttpStatus.UNAUTHORIZED)
                    .body(ApiResponseDto.error(401, "Token inválido o expirado"));
            }
            
            Integer userId = Integer.parseInt(user.getId());
            var sessionInfo = authService.getSessionInfo(userId);
            
            return ResponseEntity.ok(ApiResponseDto.success(sessionInfo));
            
        } catch (Exception e) {
            logger.error("Error obteniendo info de sesión: {}", e.getMessage(), e);
            return ResponseEntity
                .status(HttpStatus.INTERNAL_SERVER_ERROR)
                .body(ApiResponseDto.error(500, "Error interno del servidor"));
        }
    }
}