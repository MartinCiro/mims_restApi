package middlewares

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"api_go/internal/infrastructure/redis"
	"api_go/internal/interfaces/api/common"
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

func (pm *PermissionsMiddleware) Handler(requiredPermissions []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Si no hay permisos requeridos, continuar
			if len(requiredPermissions) == 0 {
				next.ServeHTTP(w, r)
				return
			}

			// Obtener usuario del contexto (seteado por AuthMiddleware)
			user := r.Context().Value("user")
			if user == nil {
				response := common.NewErrorResponse(403, "No autenticado")
				common.WriteJSONResponse(w, response, 403)
				return
			}

			// Extraer ID del usuario
			userInfo, ok := user.(map[string]interface{})
			if !ok {
				response := common.NewErrorResponse(403, "Error en información de usuario")
				common.WriteJSONResponse(w, response, 403)
				return
			}

			userInfoData, ok := userInfo["userInfo"].(map[string]interface{})
			if !ok {
				response := common.NewErrorResponse(403, "Error en información de usuario")
				common.WriteJSONResponse(w, response, 403)
				return
			}

			userId, ok := userInfoData["id_user"].(string)
			if !ok {
				response := common.NewErrorResponse(403, "ID de usuario no válido")
				common.WriteJSONResponse(w, response, 403)
				return
			}

			// Obtener permisos desde Redis
			userData, err := pm.getUserPermissions(userId)
			if err != nil {
				response := common.NewErrorResponse(401, "No posee permisos suficientes para realizar esta acción")
				common.WriteJSONResponse(w, response, 401)
				return
			}

			// Verificar si el usuario tiene al menos uno de los permisos requeridos
			hasPermission := false
			for _, requiredPerm := range requiredPermissions {
				for _, userPerm := range userData.Permisos {
					if strings.EqualFold(userPerm, requiredPerm) {
						hasPermission = true
						break
					}
				}
				if hasPermission {
					break
				}
			}

			if !hasPermission {
				response := common.NewErrorResponse(403, "No posee permisos suficientes para realizar esta acción")
				common.WriteJSONResponse(w, response, 403)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (pm *PermissionsMiddleware) getUserPermissions(userID string) (*UserData, error) {
	ctx := context.Background()

	userDataJSON, err := pm.redisService.Get(ctx, fmt.Sprintf("user:%s", userID))
	if err != nil || userDataJSON == "" {
		return nil, fmt.Errorf("datos de usuario no encontrados")
	}

	var userData UserData
	err = json.Unmarshal([]byte(userDataJSON), &userData)
	if err != nil {
		return nil, fmt.Errorf("error parseando datos de usuario")
	}

	return &userData, nil
}
