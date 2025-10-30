package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"api_go/internal/infrastructure/jwt"
	"api_go/internal/infrastructure/redis"
	"api_go/internal/interfaces/api/common"
	"api_go/pkg/logger"
)

type PermissionsMiddleware struct {
	redisService *redis.Cache
}

type UserData struct {
	Permisos []string `json:"permisos"`
}

func NewPermissionsMiddleware(redisService *redis.Cache) *PermissionsMiddleware {
	return &PermissionsMiddleware{
		redisService: redisService,
	}
}

// Handler verifica los permisos del usuario
func (pm *PermissionsMiddleware) Handler(requiredPermissions []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Si no hay permisos requeridos, continuar
			if len(requiredPermissions) == 0 {
				next.ServeHTTP(w, r)
				return
			}

			// Obtener userInfo del contexto (seteado por AuthMiddleware)
			userInfo, ok := r.Context().Value("userInfo").(*jwt.VerifyResponse)
			if !ok || userInfo == nil {
				logger.Warn("intento de acceso sin información de usuario en contexto")
				response := common.NewErrorResponse(401, "No autenticado")
				common.WriteJSONResponse(w, response, 401)
				return
			}

			// ✅ CORREGIDO: Acceso correcto a la estructura anidada
			if userInfo.UserInfo == nil {
				logger.Warn("userInfo.UserInfo es nil")
				response := common.NewErrorResponse(403, "Información de usuario incompleta")
				common.WriteJSONResponse(w, response, 403)
				return
			}

			// Extraer ID del usuario
			userID := userInfo.UserInfo.IDUser
			if userID == 0 {
				logger.Warn("ID de usuario inválido en userInfo")
				response := common.NewErrorResponse(403, "ID de usuario no válido")
				common.WriteJSONResponse(w, response, 403)
				return
			}

			logger.Debug("verificando permisos",
				"user_id", userID,
				"required_permissions", requiredPermissions)

			// Obtener permisos desde Redis
			userData, err := pm.getUserPermissions(strconv.Itoa(userID))
			if err != nil {
				logger.Warn("error obteniendo permisos de usuario",
					"user_id", userID,
					"error", err)
				response := common.NewErrorResponse(403, "No posee permisos suficientes para realizar esta acción")
				common.WriteJSONResponse(w, response, 403)
				return
			}

			// Verificar si el usuario tiene al menos uno de los permisos requeridos
			hasPermission := pm.hasAnyPermission(userData.Permisos, requiredPermissions)
			if !hasPermission {
				logger.Warn("usuario sin permisos para endpoint",
					"user_id", userID,
					"required_permissions", requiredPermissions,
					"user_permissions", userData.Permisos)
				response := common.NewErrorResponse(403, "No posee permisos suficientes para realizar esta acción")
				common.WriteJSONResponse(w, response, 403)
				return
			}

			logger.Debug("permisos verificados exitosamente",
				"user_id", userID,
				"required_permissions", requiredPermissions)

			next.ServeHTTP(w, r)
		})
	}
}

// hasAnyPermission verifica si el usuario tiene al menos uno de los permisos requeridos
func (pm *PermissionsMiddleware) hasAnyPermission(userPermisos []string, requiredPermissions []string) bool {
	for _, requiredPerm := range requiredPermissions {
		for _, userPerm := range userPermisos {
			if strings.EqualFold(trimPermission(userPerm), trimPermission(requiredPerm)) {
				return true
			}
		}
	}
	return false
}

// hasAllPermission verifica si el usuario tiene TODOS los permisos requeridos
func (pm *PermissionsMiddleware) hasAllPermission(userPermisos []string, requiredPermissions []string) bool {
	for _, requiredPerm := range requiredPermissions {
		found := false
		for _, userPerm := range userPermisos {
			if strings.EqualFold(trimPermission(userPerm), trimPermission(requiredPerm)) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// trimPermission limpia espacios y convierte a minúsculas
func trimPermission(perm string) string {
	return strings.ToLower(strings.TrimSpace(perm))
}

func (pm *PermissionsMiddleware) getUserPermissions(userID string) (*UserData, error) {
	ctx := context.Background()

	cacheKey := fmt.Sprintf("user:%s", userID)
	userDataJSON, err := pm.redisService.Get(ctx, cacheKey)
	if err != nil {
		return nil, fmt.Errorf("error accediendo a Redis: %v", err)
	}

	if userDataJSON == "" {
		return nil, fmt.Errorf("datos de usuario no encontrados en cache para key: %s", cacheKey)
	}

	var userData UserData
	err = json.Unmarshal([]byte(userDataJSON), &userData)
	if err != nil {
		return nil, fmt.Errorf("error parseando datos de usuario: %v", err)
	}

	logger.Debug("permisos obtenidos de Redis",
		"user_id", userID,
		"permisos_count", len(userData.Permisos),
		"permisos", userData.Permisos)

	return &userData, nil
}
