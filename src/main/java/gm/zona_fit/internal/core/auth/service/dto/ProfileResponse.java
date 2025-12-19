package gm.zona_fit.internal.core.auth.service.dto;

import java.util.List;

public class ProfileResponse {
    private String id;
    private String nombre;
    private String rol;
    private List<String> permisos;

    public ProfileResponse(String id, String nombre, String rol, List<String> permisos) {
        this.id = id;
        this.nombre = nombre;
        this.rol = rol;
        this.permisos = permisos;
    }

    // Getters
    public String getId() {
        return id;
    }

    public String getNombre() {
        return nombre;
    }

    public String getRol() {
        return rol;
    }

    public List<String> getPermisos() {
        return permisos;
    }
}