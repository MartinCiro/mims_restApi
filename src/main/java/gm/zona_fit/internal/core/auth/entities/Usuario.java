package gm.zona_fit.internal.core.auth.entities;

public class Usuario {
    private Integer id;
    private String username;
    private String email;
    private String encryptedPassword;
    
    public boolean validarPassword(String plainPassword, PasswordEncoder encoder) {
        return encoder.matches(plainPassword, this.encryptedPassword);
    }
}