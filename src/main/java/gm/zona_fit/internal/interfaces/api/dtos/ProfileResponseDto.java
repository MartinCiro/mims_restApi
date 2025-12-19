package gm.zona_fit.internal.interfaces.api.dtos;

import com.fasterxml.jackson.annotation.JsonProperty;

import java.util.List;

public class ProfileResponseDto {
    
    @JsonProperty("id")
    private String id;
    
    @JsonProperty("nombre")
    private String nombre;
    
    @JsonProperty("rol")
    private String rol;
    
    @JsonProperty("permisos")
    private List<String> permisos;
    
    // Constructor
    public ProfileResponseDto(String id, String nombre, String rol, List<String> permisos) {
        this.id = id;
        this.nombre = nombre;
        this.rol = rol;
        this.permisos = permisos;
    }
    
    // Getters
    public String getId() { return id; }
    public String getNombre() { return nombre; }
    public String getRol() { return rol; }
    public List<String> getPermisos() { return permisos; }
}