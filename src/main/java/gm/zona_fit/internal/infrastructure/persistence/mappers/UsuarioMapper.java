package gm.zona_fit.internal.infrastructure.persistence.mappers;

public class UsuarioMapper {
    
    public static UsuarioJpaEntity toEntity(Usuario usuario, String email) {
        if (usuario == null) {
            return null;
        }
        
        UsuarioJpaEntity entity = new UsuarioJpaEntity();
        entity.setId(usuario.getId());
        entity.setUsername(usuario.getUsername());
        entity.setEmail(email);
        entity.setPasswordHash(usuario.getEncryptedPassword());
        entity.setIdRol(usuario.getIdRol());
        entity.setIdEstado(usuario.getIdEstado());
        
        return entity;
    }
    
    public static Usuario toDomain(UsuarioJpaEntity entity) {
        if (entity == null) {
            return null;
        }
        
        return Usuario.crearDesdePasswordEncriptado(
            entity.getUsername(),
            entity.getPasswordHash(),
            entity.getIdRol(),
            entity.getIdEstado()
        );
    }
}