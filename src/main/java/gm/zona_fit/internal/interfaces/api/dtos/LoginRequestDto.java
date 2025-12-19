package gm.zona_fit.internal.interfaces.api.dtos;

import com.fasterxml.jackson.annotation.JsonProperty;

public class LoginRequestDto {
    
    @JsonProperty("email")
    private String email;
    
    @JsonProperty("password")
    private String password;
    
    // Constructor vacío para Jackson
    public LoginRequestDto() {}
    
    public LoginRequestDto(String email, String password) {
        this.email = email;
        this.password = password;
    }
    
    // Getters y Setters
    public String getEmail() { return email; }
    public void setEmail(String email) { this.email = email; }
    
    public String getPassword() { return password; }
    public void setPassword(String password) { this.password = password; }
}