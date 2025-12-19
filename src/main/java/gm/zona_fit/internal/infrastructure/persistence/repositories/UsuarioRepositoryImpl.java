package gm.zona_fit.internal.infrastructure.persistence.repositories;

import gm.zona_fit.internal.core.auth.ports.UserRepositoryPort;
import gm.zona_fit.internal.core.auth.entities.Usuario;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public class UsuarioRepositoryImpl implements UserRepositoryPort {
    
    private final UsuarioJpaRepository jpaRepository;
    
    public UsuarioRepositoryImpl(UsuarioJpaRepository jpaRepository) {
        this.jpaRepository = jpaRepository;
    }
    
    @Override
    public Integer createUser(Usuario usuario, String email) {
        UsuarioJpaEntity entity = UsuarioMapper.toEntity(usuario);
        entity.setEmail(email);
        UsuarioJpaEntity saved = jpaRepository.save(entity);
        return saved.getId();
    }
}

// JPA Repository interface
interface UsuarioJpaRepository extends JpaRepository<UsuarioJpaEntity, Integer> {
    Optional<UsuarioJpaEntity> findByEmail(String email);
    Optional<UsuarioJpaEntity> findByUsername(String username);
}