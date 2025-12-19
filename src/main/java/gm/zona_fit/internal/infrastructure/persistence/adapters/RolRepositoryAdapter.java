package gm.zona_fit.internal.infrastructure.persistence.adapters;

import gm.zona_fit.internal.core.auth.entities.Rol;
import gm.zona_fit.internal.core.auth.ports.RolRepositoryPort;
import gm.zona_fit.internal.infrastructure.persistence.entities.RolJpaEntity;
import gm.zona_fit.internal.infrastructure.persistence.repositories.jpa.RolJpaRepository;
import org.springframework.stereotype.Component;

import java.util.Optional;

@Component
public class RolRepositoryAdapter implements RolRepositoryPort {
    
    private final RolJpaRepository rolJpaRepository;
    
    public RolRepositoryAdapter(RolJpaRepository rolJpaRepository) {
        this.rolJpaRepository = rolJpaRepository;
    }
    
    @Override
    public Integer findRolIdByName(String nombre) {
        if (nombre == null || nombre.isBlank()) {
            throw new IllegalArgumentException("Nombre de rol es requerido");
        }
        
        Optional<RolJpaEntity> rolEntity = rolJpaRepository.findByNombreRolIgnoreCase(nombre);
        return rolEntity.map(RolJpaEntity::getId).orElse(null);
    }
    
    @Override
    public Rol findById(Integer id) {
        if (id == null) {
            return null;
        }
        
        Optional<RolJpaEntity> rolEntity = rolJpaRepository.findById(id);
        return rolEntity.map(entity -> new Rol(entity.getId(), entity.getNombreRol()))
                       .orElse(null);
    }
    
    @Override
    public Rol findByExactName(String nombre) {
        if (nombre == null || nombre.isBlank()) {
            return null;
        }
        
        Optional<RolJpaEntity> rolEntity = rolJpaRepository.findByNombreRol(nombre);
        return rolEntity.map(entity -> new Rol(entity.getId(), entity.getNombreRol()))
                       .orElse(null);
    }
    
    @Override
    public Integer findOrCreateGuest() {
        // Buscar rol "invitado"
        Optional<RolJpaEntity> invitadoEntity = rolJpaRepository.findByNombreRolIgnoreCase("invitado");
        
        if (invitadoEntity.isPresent()) {
            return invitadoEntity.get().getId();
        }
        
        // Crear rol "invitado" si no existe
        RolJpaEntity nuevoRol = new RolJpaEntity("invitado", "Rol para usuarios invitados");
        RolJpaEntity saved = rolJpaRepository.save(nuevoRol);
        
        return saved.getId();
    }
}