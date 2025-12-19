package gm.zona_fit.internal.infrastructure.persistence.adapters;

import gm.zona_fit.internal.core.auth.entities.Usuario;
import gm.zona_fit.internal.core.auth.ports.UserRepositoryPort;
import gm.zona_fit.internal.infrastructure.persistence.entities.UsuarioJpaEntity;
import gm.zona_fit.internal.infrastructure.persistence.mappers.UsuarioMapper;
import gm.zona_fit.internal.infrastructure.persistence.repositories.jpa.UsuarioJpaRepository;
import org.springframework.stereotype.Component;

@Component
public class UserRepositoryAdapter implements UserRepositoryPort {
    
    private final UsuarioJpaRepository usuarioJpaRepository;
    
    public UserRepositoryAdapter(UsuarioJpaRepository usuarioJpaRepository) {
        this.usuarioJpaRepository = usuarioJpaRepository;
    }
    
    @Override
    public Integer createUser(Usuario usuario, String email) {
        if (usuario == null || email == null || email.isBlank()) {
            throw new IllegalArgumentException("Usuario y email son requeridos");
        }
        
        // Verificar si el email ya existe
        if (usuarioJpaRepository.existsByEmail(email)) {
            throw new RuntimeException("El email ya está registrado");
        }
        
        // Verificar si el username ya existe
        if (usuario.getUsername() != null && usuarioJpaRepository.existsByUsername(usuario.getUsername())) {
            throw new RuntimeException("El nombre de usuario ya existe");
        }
        
        UsuarioJpaEntity entity = UsuarioMapper.toEntity(usuario, email);
        UsuarioJpaEntity savedEntity = usuarioJpaRepository.save(entity);
        
        return savedEntity.getId();
    }
}