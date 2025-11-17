# Paso a Paso para Crear un Endpoint (Arquitectura Hexagonal)

## ORDEN:

### 1. **CORE - Entities** 
`internal/core/{entidad}/entities.go`
```go
// Definir las entidades de dominio
```

### 2. **CORE - Ports** 
`internal/core/{entidad}/ports.go`
```go
// Definir interfaces (contratos)
```

### 3. **CORE - Service** 
`internal/core/{entidad}/service.go`
```go
// Lógica de negocio que usa los ports
```

### 4. **INFRASTRUCTURE - Adapter** 
`internal/infrastructure/adapters/{entidad}_adapter.go`
```go
// Implementación concreta de los ports (BD, APIs externas)
```

### 5. **INTERFACES - DTOs** 
`internal/interfaces/api/{entidad}/dtos.go`
```go
// Estructuras para request/response HTTP
```

### 6. **INTERFACES - Handler** 
`internal/interfaces/api/{entidad}/handler.go`
```go
// Manejo de requests HTTP
```

### 7. **INTERFACES - Permisos** 
`internal/interfaces/api/middlewares/permissions_constants.go`
```go
// Constantes de permisos para el endpoint
```

### 8. **ROUTES** 
`internal/interfaces/api/routes/routes.go`
```go
// Registrar rutas en el router
```

---

## EJEMPLO COMPLETO: Endpoint "Categorías"

### Paso 1: Entities
`internal/core/categorias/entities.go`
```go
package categorias

import "time"

type Categoria struct {
    ID        int       `json:"id"`
    Nombre    string    `json:"nombre"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type CategoriaData struct {
    Nombre string `json:"nombre"`
}
```

### Paso 2: Ports
`internal/core/categorias/ports.go`
```go
package categorias

import "context"

type CategoriaPort interface {
    CrearCategoria(ctx context.Context, data CategoriaData) (*Categoria, error)
}
```

### Paso 3: Service
`internal/core/categorias/service.go`
```go
package categorias

import "context"

type CategoriaService struct {
    port CategoriaPort
}

func NewCategoriaService(port CategoriaPort) *CategoriaService {
    return &CategoriaService{port: port}
}

func (s *CategoriaService) ObtenerCategorias(ctx context.Context) ([]Categoria, error) {
    return s.port.ObtenerCategorias(ctx)
}
```

### Paso 4: Adapter
`internal/infrastructure/adapters/categorias_adapter.go`
```go
package adapters

import (
    "context"
    "api_go/internal/core/categorias"
    "api_go/internal/infrastructure/database/models"
    "gorm.io/gorm"
)

type CategoriasAdapter struct {
    db *gorm.DB
}

func NewCategoriasAdapter(db *gorm.DB) *CategoriasAdapter {
    return &CategoriasAdapter{db: db}
}

func (a *CategoriasAdapter) ObtenerCategorias(ctx context.Context) ([]categorias.Categoria, error) {
    var categoriasDB []models.Categoria
    err := a.db.WithContext(ctx).Find(&categoriasDB).Error
    if err != nil {
        return nil, err
    }

    result := make([]categorias.Categoria, len(categoriasDB))
    for i, cat := range categoriasDB {
        result[i] = categorias.Categoria{
            ID:        cat.ID,
            Nombre:    cat.Nombre,
            CreatedAt: cat.CreatedAt,
            UpdatedAt: cat.UpdatedAt,
        }
    }
    return result, nil
}

// Implementar otros métodos...
```

### Paso 5: DTOs
`internal/interfaces/api/categorias/dtos.go`
```go
package categorias

import "api_go/internal/core/categorias"

type CreateCategoriaRequest struct {
    Nombre string `json:"nombre" validate:"required,min=3"`
}

type CategoriaResponse struct {
    ID        int       `json:"id"`
    Nombre    string    `json:"nombre"`
    CreatedAt time.Time `json:"created_at"`
}

func ToCategoriaData(req CreateCategoriaRequest) categorias.CategoriaData {
    return categorias.CategoriaData{
        Nombre: req.Nombre,
    }
}

func FromCategoria(cat *categorias.Categoria) *CategoriaResponse {
    if cat == nil {
        return nil
    }
    return &CategoriaResponse{
        ID:        cat.ID,
        Nombre:    cat.Nombre,
        CreatedAt: cat.CreatedAt,
    }
}
```

### Paso 6: Handler
`internal/interfaces/api/categorias/handler.go`
```go
package categorias

import (
    "encoding/json"
    "net/http"
    "strconv"
    
    "api_go/internal/core/categorias"
    "api_go/internal/interfaces/api/common"
)

type CategoriaHandler struct {
    service *categorias.CategoriaService
}

func NewCategoriaHandler(service *categorias.CategoriaService) *CategoriaHandler {
    return &CategoriaHandler{service: service}
}

func (h *CategoriaHandler) ObtenerCategorias(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    categorias, err := h.service.ObtenerCategorias(ctx)
    if err != nil {
        common.WriteSimpleError(w, err.Error(), 500)
        return
    }

    response := common.NewSuccessResponse(categorias)
    common.WriteJSONResponse(w, response, 200)
}
```

### Paso 7: Permisos
`internal/interfaces/api/middlewares/permissions_constants.go`
```go
package middlewares

// Categorías
var PermissionCategoriasLeer = []string{"categorias:leer"}
```

### Paso 8: Routes
`internal/interfaces/api/routes/routes.go`
```go
// En SetupRoutes():
categoriaHandler := categorias_handler.NewCategoriaHandler(categoriaService)

// Rutas protegidas
protected.Handle("POST /api/categorias", 
    authMiddleware.RequireAuthAndPermission(permsMiddleware, middlewares.PermissionCategoriasLeer)(
        http.HandlerFunc(categoriaHandler.ObtenerCategorias),
    ))
```

---

**RESUMEN DEL ORDEN:**
1. 🎯 **Core** (Entities → Ports → Service)
2. 🏗️ **Infrastructure** (Adapter) 
3. 🌐 **Interfaces** (DTOs → Handler → Permisos → Routes)
