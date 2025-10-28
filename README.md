```bash
.
├── cmd/
│   └── api/
│       └── main.go                    # Punto de entrada de la aplicación
├── config/
│   └── config.go                      # Configuración y variables de entorno
├── internal/                          # Código interno de la aplicación
│   ├── app/
│   │   └── app.go                     # Inicialización y configuración principal de la app
│   ├── core/                          # Lógica de negocio (Hexagonal Core)
│   │   ├── auth/                      # Módulo de autenticación
│   │   │   ├── entities.go            # Entidades: User, AuthData
│   │   │   ├── ports.go               # Interfaces: AuthPort
│   │   │   └── service.go             # Lógica de negocio: AuthService
│   │   ├── estados/                   # Módulo de estados
│   │   │   ├── entities.go            # Entidades: Estado, EstadoData
│   │   │   ├── ports.go               # Interfaces: EstadosPort
│   │   │   └── service.go             # Lógica de negocio: EstadoService
│   │   └── login/                     # Caso de uso específico de login
│   │       └── service.go             # Lógica completa del flujo de login
│   ├── infrastructure/                # Adaptadores de infraestructura
│   │   ├── database/
│   │   │   ├── database.go            # Conexión y configuración de BD
│   │   │   └── models/
│   │   │       └── usuario.go         # Modelos GORM para BD
│   │   ├── jwt/
│   │   │   └── jwt_service.go         # Servicio JWT (generación/validación)
│   │   ├── redis/
│   │   │   ├── cache.go               # Cliente y operaciones Redis
│   │   │   └── initializer.go         # Inicialización de Redis
│   │   └── repositories/              # Implementaciones de puertos (adaptadores BD)
│   │       ├── auth_repository.go     # AuthAdapter implementa AuthPort
│   │       └── estado_repository.go   # EstadosAdapter implementa EstadosPort
│   └── interfaces/                    # Adaptadores de interfaces (HTTP)
│       └── api/
│           ├── common/
│           │   └── responses.go       # Estructuras estándar de respuesta HTTP
│           ├── handlers/              # Handlers HTTP (controladores)
│           │   ├── auth/
│           │   │   └── handler.go     # Handler para operaciones de auth
│           │   ├── common/
│           │   │   └── handler.go     # Health checks y rutas comunes
│           │   ├── estados/
│           │   │   └── handler.go     # Handler para CRUD de estados
│           │   ├── login/
│           │   │   └── handler.go     # Handler específico para login
│           │   └── usuarios/
│           │       └── handler.go     # Handler para usuarios (placeholder)
│           ├── middlewares/           # Middlewares HTTP
│           │   ├── auth.go            # Middleware de autenticación JWT
│           │   └── permissions.go     # Middleware de verificación de permisos
│           └── routes/
│               └── routes.go          # Configuración de rutas y enrutamiento
└── pkg/
    └── utils/                         # Utilidades compartidas
        ├── http_helpers.go            # Helpers para HTTP
        ├── password.go                # Servicio de encriptación de passwords
        └── validators.go              # Validadores y utilities de validación
```

## **Estructura Arquitectura Hexagonal:**

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