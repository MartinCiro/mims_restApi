package gm.zona_fit.internal.infrastructure.persistence.adapters;

import gm.zona_fit.internal.core.auth.entities.AuthData;
import gm.zona_fit.internal.core.auth.entities.User;
import gm.zona_fit.internal.core.auth.ports.AuthPort;
import gm.zona_fit.internal.infrastructure.persistence.entities.UsuarioJpaEntity;
import gm.zona_fit.internal.infrastructure.persistence.repositories.jpa.RolJpaRepository;
import gm.zona_fit.internal.infrastructure.persistence.repositories.jpa.UsuarioJpaRepository;
import org.springframework.stereotype.Component;

import java.util.Collections;
import java.util.List;
import java.util.Optional;

@Component
public class AuthAdapter implements AuthPort {
    
    private final UsuarioJpaRepository usuarioJpaRepository;
    private final RolJpaRepository rolJpaRepository;
    
    public AuthAdapter(UsuarioJpaRepository usuarioJpaRepository, 
                      RolJpaRepository rolJpaRepository) {
        this.usuarioJpaRepository = usuarioJpaRepository;
        this.rolJpaRepository = rolJpaRepository;
    }
    
    @Override
    public User retrieveUser(AuthData authData) {
        if (authData == null) {
            return null;
        }
        
        Optional<UsuarioJpaEntity> usuarioEntity = usuarioJpaRepository
            .findByEmailOrUsername(authData.getEmail(), authData.getUsername());
        
        if (usuarioEntity.isEmpty()) {
            return null;
        }
        
        UsuarioJpaEntity entity = usuarioEntity.get();
        
        // Obtener permisos del rol
        List<String> permisos = entity.getIdRol() != null 
            ? rolJpaRepository.findPermisosByRoleId(entity.getIdRol())
            : Collections.emptyList();
        
        return new User(
            entity.getId().toString(),
            entity.getUsername(),
            entity.getEmail(),
            entity.getPasswordHash(),
            entity.getIdRol(),
            entity.getIdEstado(),
            permisos,
            "rol_nombre" // Este vendría de otra consulta
        );
    }
    
    @Override
    public User retrieveUserById(Integer userId) {
        if (userId == null) {
            return null;
        }
        
        Optional<UsuarioJpaEntity> usuarioEntity = usuarioJpaRepository.findById(userId);
        
        if (usuarioEntity.isEmpty()) {
            return null;
        }
        
        UsuarioJpaEntity entity = usuarioEntity.get();
        
        // Obtener permisos del rol
        List<String> permisos = entity.getIdRol() != null 
            ? rolJpaRepository.findPermisosByRoleId(entity.getIdRol())
            : Collections.emptyList();
        
        return new User(
            entity.getId().toString(),
            entity.getUsername(),
            entity.getEmail(),
            entity.getPasswordHash(),
            entity.getIdRol(),
            entity.getIdEstado(),
            permisos,
            "rol_nombre" // Este vendría de otra consulta
        );
    }
    
    @Override
    public List<String> getPermissionsByRoleId(Integer roleId) {
        if (roleId == null) {
            return Collections.emptyList();
        }
        
        return rolJpaRepository.findPermisosByRoleId(roleId);
    }
}