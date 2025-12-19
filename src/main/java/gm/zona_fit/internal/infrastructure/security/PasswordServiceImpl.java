package gm.zona_fit.internal.infrastructure.security;

import gm.zona_fit.internal.core.auth.entities.PasswordService;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Component;

@Component
public class PasswordServiceImpl implements PasswordService {
    
    private final PasswordEncoder passwordEncoder;
    
    public PasswordServiceImpl() {
        this.passwordEncoder = new BCryptPasswordEncoder();
    }
    
    @Override
    public boolean compararPasswords(String plainPassword, String encryptedPassword) {
        return passwordEncoder.matches(plainPassword, encryptedPassword);
    }
    
    // Método adicional para encriptar (no en la interfaz original)
    public String encriptarPassword(String plainPassword) {
        return passwordEncoder.encode(plainPassword);
    }
}