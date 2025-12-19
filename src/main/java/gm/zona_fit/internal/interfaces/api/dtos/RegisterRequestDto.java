package gm.zona_fit.internal.interfaces.api.dtos;

import com.fasterxml.jackson.annotation.JsonProperty;

public class RegisterRequestDto {
    
    @JsonProperty("username")
    private String username;
    
    @JsonProperty("email")
    private String email;
    
    @JsonProperty("password")
    private String password;
    
    @JsonProperty("rol_nombre")
    private String rolNombre;
    
    // Constructor vacío para Jackson
    public RegisterRequestDto() {}
    
    public RegisterRequestDto(String username, String email, String password, String rolNombre) {
        this.username = username;
        this.email = email;
        this.password = password;
        this.rolNombre = rolNombre;
    }
    
    // Getters y Setters
    public String getUsername() { return username; }
    public void setUsername(String username) { this.username = username; }
    
    public String getEmail() { return email; }
    public void setEmail(String email) { this.email = email; }
    
    public String getPassword() { return password; }
    public void setPassword(String password) { this.password = password; }
    
    public String getRolNombre() { return rolNombre; }
    public void setRolNombre(String rolNombre) { this.rolNombre = rolNombre; }
}