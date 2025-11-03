## 📁 **ESTRUCTURA COMPLETA DEL PROYECTO**

``` bash
template-nest/
├── 📄 Dockerfile
├── 📄 Makefile
├── 📄 README.md
├── 📄 go.mod
├── 📄 go.sum
├── 📄 docker-compose.yml
├── 📄 .env
├── 📄 .env.template
├── 📄 package-lock.json
├── 📁 cmd/
│   └── 📁 api/
│       └── 📄 main.go
├── 📁 config/
│   └── 📄 config.go
├── 📁 example/
├── 📁 infrastructure/
│   └── 📁 database/
│       └── 📁 models/
│           └── 📄 tables.go
├── 📁 internal/
│   ├── 📁 app/
│   │   └── 📄 app.go
│   ├── 📁 core/
│   │   ├── 📁 auth/
│   │   │   ├── 📄 entities.go
│   │   │   ├── 📄 ports.go
│   │   │   └── 📄 service.go
│   │   ├── 📁 common/
│   │   │   └── 📄 entities.go
│   │   ├── 📁 estados/
│   │   │   ├── 📄 entities.go
│   │   │   ├── 📄 ports.go
│   │   │   └── 📄 service.go
│   │   └── 📁 login/
│   │       └── 📄 service.go
│   ├── 📁 infrastructure/
│   │   ├── 📁 database/
│   │   │   ├── 📄 database.go
│   │   │   ├── 📄 db_manager.go
│   │   │   └── 📁 models/
│   │   │       └── 📄 usuario.go
│   │   ├── 📁 jwt/
│   │   │   └── 📄 jwt_service.go
│   │   ├── 📁 redis/
│   │   │   ├── 📄 cache.go
│   │   │   └── 📄 initializer.go
│   │   └── 📁 repositories/
│   │       ├── 📄 auth_repository.go
│   │       ├── 📄 estado_repository.go
│   │       ├── 📄 permiso_repository.go
│   │       ├── 📄 rol_repository.go
│   │       └── 📄 usuario_repository.go
│   └── 📁 interfaces/
│       └── 📁 api/
│           ├── 📁 common/
│           │   └── 📄 responses.go
│           ├── 📁 estados/
│           │   ├── 📄 adapters.go
│           │   └── 📄 dtos.go
│           ├── 📁 handlers/
│           │   ├── 📁 auth/
│           │   │   ├── 📄 handler.go
│           │   │   └── 📄 profile_handler.go
│           │   ├── 📁 common/
│           │   │   └── 📄 handler.go
│           │   ├── 📁 estados/
│           │   │   └── 📄 handler.go
│           │   ├── 📁 login/
│           │   │   └── 📄 handler.go
│           │   └── 📁 usuarios/
│           │       └── 📄 handler.go
│           ├── 📁 middlewares/
│           │   ├── 📄 auth.go
│           │   ├── 📄 permissions.go
│           │   └── 📄 permissions_constants.go
│           └── 📁 routes/
│               ├── 📄 middleware_helpers.go
│               └── 📄 routes.go
├── 📁 manifests/                          # 🆕 KUBERNETES MANIFESTS
│   ├── 📁 go-api/
│   │   ├── 📄 configmap-simple.yaml       # ConfigMap simplificado
│   │   ├── 📄 configmap.yaml              # ConfigMap completo
│   │   ├── 📄 deployment.yaml             # Deployment de la API Go
│   │   └── 📄 service.yaml                # Service de la API
│   ├── 📁 postgres/
│   │   ├── 📄 deployment.yaml             # Deployment de PostgreSQL
│   │   ├── 📄 secret.yaml                 # Secret de PostgreSQL
│   │   └── 📄 service.yaml                # Service de PostgreSQL
│   ├── 📁 redis/
│   │   ├── 📄 deployment.yaml             # Deployment de Redis
│   │   └── 📄 service.yaml                # Service de Redis
│   └── 📁 shared/
│       ├── 📄 network-policies.yaml       # Políticas de red
│       ├── 📄 secrets-verification.txt    # Verificación de secrets (NO COMMIT)
│       ├── 📄 secrets.yaml                # Secrets de la aplicación
│       └── 📄 volumes.yaml                # Volúmenes persistentes
├── 📁 pkg/
│   ├── 📁 logger/
│   │   └── 📄 logger.go
│   └── 📁 utils/
│       ├── 📄 http_helpers.go
│       ├── 📄 password.go
│       └── 📄 validators.go
└── 📁 scripts/
    └── 📄 generate-secrets.sh             # Script para generar secrets

40 directories, 65 files
```

### 🚀 **ARCHIVOS KUBERNETES**

#### **`manifests/go-api/`**

- **`deployment.yaml`** - 3 réplicas de tu API Go
- **`service.yaml`** - Service ClusterIP en puerto 80
- **`configmap.yaml`** - Configuración de la aplicación
- **`configmap-simple.yaml`** - Versión simplificada

#### **`manifests/postgres/`**

- **`deployment.yaml`** - PostgreSQL con volumen persistente
- **`service.yaml`** - Service en puerto 5432
- **`secret.yaml`** - Credenciales de base de datos

#### **`manifests/redis/`**

- **`deployment.yaml`** - Redis con autenticación
- **`service.yaml`** - Service en puerto 6379

#### **`manifests/shared/`**

- **`volumes.yaml`** - PVCs para PostgreSQL (5Gi) y Redis (1Gi)
- **`secrets.yaml`** - Secrets de aplicación (JWT, DB, Redis)
- **`network-policies.yaml`** - Seguridad de red
- **`secrets-verification.txt`** - Para desarrollo (NO committear)

#### 📋 **COMANDOS DE DESPLIEGUE**

```bash
# Flujo completo
make build                    # Construir imagen Docker
make generate-secrets         # Generar secrets desde .env
make deploy                   # Desplegar en Kubernetes
make status                   # Verificar estado

# Comandos útiles
make logs                     # Ver logs de la API
make port-forward             # Acceder localmente (localhost:8080)
make clean                    # Limpiar recursos
```


## **Estructura Arquitectura Hexagonal:**

### **Flujo de Datos:**

``` go
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

``` go
HTTP Request → Routes → Middlewares → Handlers → Services → Ports → Adapters → BD/Redis
HTTP Response ← Handlers ← Services ← Ports ← Adapters ← BD/Redis
```

## ⚒️ Utilidades

### 🔐 Generar secretos 
```bash
make generate-secrets
```

### 🗃️ Ver contenido del Secret llamado "app-secrets"
```bash
kubectl get secret app-secrets -o yaml
```

### ⚙️ Aplicar archivos de configuración
```bash
kubectl apply -f manifests/go-api/configmap.yaml
```

### 🔄 Reiniciar cluster
```bash
kubectl rollout restart deployment/go-api-deployment
```

### ✅ Ver estado pod

```bash
kubectl get pods -w
```

### ⏳ Validar pods pendientes

```bash
kubectl describe pod postgres-deployment-xxxxx
```

### 🗑️ Eliminar PVs

```bash
kubectl delete pv postgres-prod-pv redis-prod-pv
```

### 📄 Logs en tiempo real
```bash
kubectl logs -f go-api-deployment-xxxxx
```

### 🖥️ Monitorear y filtrar los pods que tengan la etiqueta app=go-api

```bash
kubectl get pods -l app=go-api -w
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
