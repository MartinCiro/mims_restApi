package gm.zona_fit.internal.core.auth.entities;

import java.util.Objects;

public class Usuario {
    private Integer id;
    private String username;
    private String email;
    private Integer idRol;
    private Integer idEstado;
    private String encryptedPassword;

    // Constructor principal
    private Usuario(String username, String encryptedPassword, Integer idRol, Integer idEstado) {
        this.username = username;
        this.encryptedPassword = encryptedPassword;
        this.idRol = idRol;
        this.idEstado = idEstado;
    }

    // Constructor adicional
    public Usuario(Integer id, String username, String encryptedPassword,
            Integer idRol, Integer idEstado, String email) {
        this.id = id;
        this.username = username;
        this.encryptedPassword = encryptedPassword;
        this.idRol = idRol;
        this.idEstado = idEstado;
        this.email = email;
    }

    // Factory method: NewUsuario (crea con password encriptado)
    public static Usuario crearNuevoUsuario(String username, String plainPassword, Integer idRol, Integer idEstado) {
        String encryptedPassword = encriptarPassword(plainPassword);
        return new Usuario(username, encryptedPassword, idRol, idEstado);
    }

    // Factory method: NewUsuarioFromEncrypted
    public static Usuario crearDesdePasswordEncriptado(String username, String encryptedPassword, Integer idRol,
            Integer idEstado) {
        return new Usuario(null, username, encryptedPassword, idRol, idEstado, null);
    }

    // GetEncryptedPassword
    public String getEncryptedPassword() {
        return encryptedPassword;
    }

    // ComparePassword (delegará a un servicio externo)
    public boolean compararPassword(String plainPassword, PasswordService passwordService) {
        return passwordService.compararPasswords(plainPassword, this.encryptedPassword);
    }

    // encodePassword (método privado estático)
    private static String encriptarPassword(String plainPassword) {
        // La implementación real irá en infrastructure
        // Por ahora solo placeholder
        return "encrypted:" + plainPassword;
    }

    // Getters y Setters
    public Integer getId() {
        return id;
    }

    public void setId(Integer id) {
        this.id = id;
    }

    public String getUsername() {
        return username;
    }

    public String getEmail() {
        return email;
    }

    public void setEmail(String email) {
        this.email = email;
    }

    public Integer getIdRol() {
        return idRol;
    }

    public Integer getIdEstado() {
        return idEstado;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o)
            return true;
        if (o == null || getClass() != o.getClass())
            return false;
        Usuario usuario = (Usuario) o;
        return Objects.equals(id, usuario.id) &&
                Objects.equals(username, usuario.username) &&
                Objects.equals(email, usuario.email);
    }

    @Override
    public int hashCode() {
        return Objects.hash(id, username, email);
    }

    @Override
    public String toString() {
        return "Usuario{" +
                "id=" + id +
                ", username='" + username + '\'' +
                ", email='" + email + '\'' +
                '}';
    }
}