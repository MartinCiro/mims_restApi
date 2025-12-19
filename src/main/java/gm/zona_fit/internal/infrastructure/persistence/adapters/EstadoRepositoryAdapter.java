package gm.zona_fit.internal.infrastructure.persistence.adapters;

import gm.zona_fit.internal.core.auth.entities.Estado;
import gm.zona_fit.internal.core.auth.ports.EstadoRepositoryPort;
import gm.zona_fit.internal.infrastructure.persistence.entities.EstadoJpaEntity;
import gm.zona_fit.internal.infrastructure.persistence.repositories.jpa.EstadoJpaRepository;
import org.springframework.stereotype.Component;

import java.util.Optional;

@Component
public class EstadoRepositoryAdapter implements EstadoRepositoryPort {
    
    private final EstadoJpaRepository estadoJpaRepository;
    
    public EstadoRepositoryAdapter(EstadoJpaRepository estadoJpaRepository) {
        this.estadoJpaRepository = estadoJpaRepository;
    }
    
    @Override
    public Integer findEstadoIdByName(String nombre) {
        if (nombre == null || nombre.isBlank()) {
            throw new IllegalArgumentException("Nombre de estado es requerido");
        }
        
        Optional<EstadoJpaEntity> estadoEntity = estadoJpaRepository.findByNombreEstado(nombre);
        return estadoEntity.map(EstadoJpaEntity::getId).orElse(null);
    }
    
    @Override
    public Estado findByExactName(String nombre) {
        if (nombre == null || nombre.isBlank()) {
            return null;
        }
        
        Optional<EstadoJpaEntity> estadoEntity = estadoJpaRepository.findByNombreEstado(nombre);
        return estadoEntity.map(entity -> new Estado(entity.getId(), entity.getNombreEstado()))
                          .orElse(null);
    }
    
    @Override
    public Integer findOrCreateActive() {
        // Buscar estado "activo"
        Optional<EstadoJpaEntity> activoEntity = estadoJpaRepository.findByNombreEstado("activo");
        
        if (activoEntity.isPresent()) {
            return activoEntity.get().getId();
        }
        
        // Crear estado "activo" si no existe
        EstadoJpaEntity nuevoEstado = new EstadoJpaEntity("activo", "Estado activo para usuarios");
        EstadoJpaEntity saved = estadoJpaRepository.save(nuevoEstado);
        
        return saved.getId();
    }
}