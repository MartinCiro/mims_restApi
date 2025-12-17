package gm.zona_fit.repositorio;

import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.atomic.AtomicLong;

import org.springframework.stereotype.Repository;

import gm.zona_fit.modelo.Contacto;

@Repository
public class ContactoRepositorio {
    // Simulación de base de datos en memoria
    private final List<Contacto> contactos = new ArrayList<>();
    private final AtomicLong contadorId = new AtomicLong(1);
    
    // Métodos CRUD básicos
    public List<Contacto> findAll() {
        return new ArrayList<>(contactos); // Retorna copia
    }
    
    public Contacto findById(Long id) {
        return contactos.stream()
                .filter(c -> c.getId().equals(id))
                .findFirst()
                .orElse(null);
    }
    
    public List<Contacto> findByNombre(String nombre) {
        return contactos.stream()
                .filter(c -> c.getNombre().toLowerCase().contains(nombre.toLowerCase()))
                .toList();
    }
    
    public Contacto save(Contacto contacto) {
        if (contacto.getId() == null) {
            // Nuevo contacto
            contacto.setId(contadorId.getAndIncrement());
            contactos.add(contacto);
        } else {
            // Actualizar contacto existente
            contactos.removeIf(c -> c.getId().equals(contacto.getId()));
            contactos.add(contacto);
        }
        return contacto;
    }
    
    public void deleteById(Long id) {
        contactos.removeIf(c -> c.getId().equals(id));
    }
    
    public boolean existsById(Long id) {
        return contactos.stream().anyMatch(c -> c.getId().equals(id));
    }
}