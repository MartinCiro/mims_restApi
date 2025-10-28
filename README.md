```bash
.
├── cmd/
│   └── api/
│       └── main.go                          # Punto de entrada de la aplicación
├── config/
│   └── config.go                            # Configuración y variables de entorno
├── internal/
│   ├── app/
│   │   └── app.go                           # Inicialización y configuración principal
│   ├── core/                                # Lógica de negocio (Hexagonal Core)
│   │   ├── auth/
│   │   │   ├── entities.go                  # Entidades: User, AuthData, LoginCredentials
│   │   │   ├── ports.go                     # Interfaces: AuthPort
│   │   │   └── service.go                   # Lógica: AuthService
│   │   ├── common/
│   │   │   └── entities.go                  # Entidades comunes: UserContext
│   │   ├── estados/
│   │   │   ├── entities.go                  # Entidades: Estado, EstadoData, EstadoDataUpdate, EstadoDataXid
│   │   │   ├── ports.go                     # Interfaces: EstadosPort
│   │   │   └── service.go                   # Lógica: EstadoService
│   │   └── login/
│   │       └── service.go                   # Caso de uso: LoginService, LoginResult
│   ├── infrastructure/                      # Adaptadores de infraestructura
│   │   ├── database/
│   │   │   ├── database.go                  # Conexión y configuración BD
│   │   │   ├── db_manager.go                # DBManager (simulador Prisma)
│   │   │   └── models/
│   │   │       └── usuario.go               # Modelos GORM (opcional)
│   │   ├── jwt/
│   │   │   └── jwt_service.go               # Servicio JWT: JWTService, UserInfo, JwtPayload
│   │   ├── redis/
│   │   │   ├── cache.go                     # Cliente Redis: Cache
│   │   │   └── initializer.go               # Inicialización Redis
│   │   └── repositories/                    # Implementaciones de puertos
│   │       ├── auth_repository.go           # AuthAdapter (coordinador)
│   │       ├── estado_repository.go         # EstadosAdapter
│   │       ├── permiso_repository.go        # PermisoRepository
│   │       ├── rol_repository.go            # RolRepository
│   │       └── usuario_repository.go        # UsuarioRepository
│   └── interfaces/                          # Adaptadores de interfaces (HTTP)
│       └── api/
│           ├── common/
│           │   └── responses.go             # ResponseBody[T], helpers de respuesta
│           ├── estados/
│           │   ├── adapters.go              # Adaptadores: ToEstadoData, FromEstado, etc.
│           │   └── dtos.go                  # DTOs: CreateEstadoRequest, UpdateEstadoRequest, etc.
│           ├── handlers/
│           │   ├── auth/
│           │   │   └── handler.go           # AuthHandler
│           │   ├── common/
│           │   │   └── handler.go           # HealthHandler, ReadyHandler
│           │   ├── estados/
│           │   │   └── handler.go           # EstadosHandler + DTOs de validación
│           │   ├── login/
│           │   │   └── handler.go           # LoginHandler + LoginRequestDTO
│           │   └── usuarios/
│           │       └── handler.go           # UsuariosHandler (placeholder)
│           ├── middlewares/
│           │   ├── auth.go                  # AuthMiddleware
│           │   └── permissions.go           # PermissionsMiddleware
│           └── routes/
│               └── routes.go                # Configuración de rutas y middlewares globales
└── pkg/
    └── utils/                               # Utilidades compartidas
        ├── http_helpers.go                  # Helpers HTTP: ValidateRequired, etc.
        ├── password.go                      # PasswordService
        ├── validation_errors.go             # ValidationError, ValidationErrors, HandleException
        └── validators.go                    # Validadores: ValidarBlank, ValidarEmail, etc.
```


## **Estructura Arquitectura Hexagonal:**

### **Flujo de Datos:**

```
HTTP Request 
    → Routes 
    → Middlewares (Auth, Permissions)
    → Handlers (adaptan HTTP → DTOs)
    → Services Core (lógica de negocio pura)
    → Ports (interfaces)
    → Repositories (adaptadores BD coordinados)
    → DBManager (simulador Prisma)
    → PostgreSQL

HTTP Response 
    ← Handlers (adaptan Domain → ResponseBody)  
    ← Services (objetos del dominio)
    ← Repositories (entidades del dominio)
```

### **Core (Dominio)**
- **`internal/core/`** - Lógica de negocio pura
- **`entities/`** - Entidades del dominio
- **`ports/`** - Interfaces que el dominio espera
- **`services/`** - Casos de uso y lógica de negocio

### **Infrastructure (Adaptadores)**
- **`internal/infrastructure/`** - Implementaciones concretas
- **`repositories/`** - Adaptadores de persistencia
- **`jwt/`, `redis/`** - Adaptadores de servicios externos

### **Interfaces (Controladores)**
- **`internal/interfaces/`** - Adaptadores de entrada (HTTP)
- **`handlers/`** - Controladores HTTP
- **`middlewares/`** - Middlewares de la API

### **Shared**
- **`pkg/utils/`** - Utilidades compartidas
- **`config/`** - Configuración

## **Flujo de Datos:**

```
HTTP Request → Routes → Middlewares → Handlers → Services → Ports → Adapters → BD/Redis
HTTP Response ← Handlers ← Services ← Ports ← Adapters ← BD/Redis
```


## 📋 Cambiar el Método de Autenticación (Ejemplo Email)

### **Archivos a Modificar**

#### **Core (Lógica de Negocio)**

```go
package auth
//internal/core/auth/ports.go
type AuthData struct {
    Email string `json:"email"` // ← CAMBIAR
}

// internal/core/auth/entities.go
type User struct {
    ID           int      `json:"id,omitempty"`
    Documento    string   `json:"documento,omitempty"`
    Username     string   `json:"username"`
    Email        string   `json:"email"` // Hacer obligatorio
    PasswordHash string   `json:"password_hash,omitempty"`
    IDRol        *int     `json:"id_rol,omitempty"`
    // ... otros campos
}
```

```go
// internal/core/login/service.go
package login

type LoginCredentials struct {
    Email    string `json:"email"`    // ← CAMBIAR
    Password string `json:"password"`
}

// En el método Execute, cambiar:
user, err := s.authPort.RetrieveUser(ctx, auth.AuthData{
    Email: credentials.Email, // ← CAMBIAR
})
```

```go
// internal/infrastructure/repositories/usuario_repository.go
package repositories

func (r *UsuarioRepository) FindByUsername(ctx context.Context, email string) (*auth.User, error) {
    ...
    err := r.dbManager.FindUnique(ctx, "usuario", &usuarioDB, map[string]interface{}{
        "email": email,
    })
```

#### **Infrastructure (Acceso a Datos)**

```go
// internal/infrastructure/repositories/usuario_repository.go
func (r *UsuarioRepository) FindByEmail(ctx context.Context, email string) (*auth.User, error) {
    // Buscar por email en BD
    err := r.dbManager.FindUnique(ctx, "usuario", &usuarioDB, map[string]interface{}{
        "email": email, // ← CAMBIAR
    })
}
```

```go
// internal/infrastructure/repositories/auth_repository.go
func (a *AuthAdapter) RetrieveUser(ctx context.Context, authData auth.AuthData) (*auth.User, error) {
    usuario, err := a.usuarioRepo.FindByEmail(ctx, authData.Email) // ← CAMBIAR
}
```

#### **Interfaces (HTTP)**

```go
// internal/interfaces/api/handlers/login/handler.go
type LoginRequestDTO struct {
    Email    string `json:"Email"`    // ← CAMBIAR
    Password string `json:"Password"`
}

func (dto *LoginRequestDTO) Validate() error {
    if strings.TrimSpace(dto.Email) == "" {
        return utils.NewValidationError("Email", "El correo electrónico es obligatorio")
    }
}
```