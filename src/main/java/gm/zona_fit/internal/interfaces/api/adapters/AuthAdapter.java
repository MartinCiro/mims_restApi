package gm.zona_fit.internal.interfaces.api.adapters;

import gm.zona_fit.internal.core.auth.service.dto.LoginRequest;
import gm.zona_fit.internal.core.auth.service.dto.RegisterRequest;
import gm.zona_fit.internal.interfaces.api.dtos.LoginRequestDto;
import gm.zona_fit.internal.interfaces.api.dtos.RegisterRequestDto;

public class AuthAdapter {
    
    private AuthAdapter() {
        // Clase de utilidad, no instanciable
    }
    
    public static LoginRequest toDomain(LoginRequestDto dto) {
        if (dto == null) {
            return null;
        }
        
        LoginRequest domain = new LoginRequest();
        domain.setEmail(dto.getEmail());
        domain.setPassword(dto.getPassword());
        
        return domain;
    }
    
    public static RegisterRequest toDomain(RegisterRequestDto dto) {
        if (dto == null) {
            return null;
        }
        
        RegisterRequest domain = new RegisterRequest();
        domain.setUsername(dto.getUsername());
        domain.setEmail(dto.getEmail());
        domain.setPassword(dto.getPassword());
        
        // Convertir rol_nombre a rolId si es necesario
        // Esto se hará en el servicio
        
        return domain;
    }
}