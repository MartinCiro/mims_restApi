package gm.zona_fit.controlador;

import gm.zona_fit.modelo.Contacto;
import gm.zona_fit.servicio.ContactoServicio;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/contactos")
@CrossOrigin(origins = "*") // Para permitir peticiones desde cualquier origen
public class ContactoControlador {
    
    @Autowired
    private ContactoServicio contactoServicio;
    
    // GET: Obtener todos los contactos
    @GetMapping
    public ResponseEntity<List<Contacto>> listarContactos() {
        List<Contacto> contactos = contactoServicio.listarContactos();
        return ResponseEntity.ok(contactos);
    }
    
    // GET: Obtener contacto por ID
    @GetMapping("/{id}")
    public ResponseEntity<Contacto> obtenerContacto(@PathVariable Long id) {
        Contacto contacto = contactoServicio.buscarContactoPorId(id);
        if (contacto != null) {
            return ResponseEntity.ok(contacto);
        }
        return ResponseEntity.status(HttpStatus.NOT_FOUND).build();
    }
    
    // GET: Buscar contactos por nombre
    @GetMapping("/buscar")
    public ResponseEntity<List<Contacto>> buscarPorNombre(@RequestParam String nombre) {
        List<Contacto> contactos = contactoServicio.buscarContactoPorNombre(nombre);
        return ResponseEntity.ok(contactos);
    }
    
    // POST: Crear nuevo contacto
    @PostMapping
    public ResponseEntity<Contacto> crearContacto(@RequestBody Contacto contacto) {
        Contacto nuevoContacto = contactoServicio.guardarContacto(contacto);
        return ResponseEntity.status(HttpStatus.CREATED).body(nuevoContacto);
    }
    
    // PUT: Actualizar contacto existente
    @PutMapping("/{id}")
    public ResponseEntity<Contacto> actualizarContacto(
            @PathVariable Long id, 
            @RequestBody Contacto contacto) {
        
        Contacto contactoActualizado = contactoServicio.actualizarContacto(id, contacto);
        if (contactoActualizado != null) {
            return ResponseEntity.ok(contactoActualizado);
        }
        return ResponseEntity.status(HttpStatus.NOT_FOUND).build();
    }
    
    // DELETE: Eliminar contacto
    @DeleteMapping("/{id}")
    public ResponseEntity<Void> eliminarContacto(@PathVariable Long id) {
        Contacto contacto = contactoServicio.buscarContactoPorId(id);
        if (contacto != null) {
            contactoServicio.eliminarContacto(id);
            return ResponseEntity.noContent().build();
        }
        return ResponseEntity.status(HttpStatus.NOT_FOUND).build();
    }
    
    // Endpoint adicional: Obtener contacto como matriz (como lo pediste originalmente)
    @GetMapping("/{id}/matriz")
    public ResponseEntity<String[][]> obtenerContactoComoMatriz(@PathVariable Long id) {
        Contacto contacto = contactoServicio.buscarContactoPorId(id);
        if (contacto != null) {
            String[][] matriz = new String[2][4];
            
            // Fila 0: Nombres de campos
            matriz[0][0] = "nombre";
            matriz[0][1] = "apellido";
            matriz[0][2] = "email";
            matriz[0][3] = "celular";
            
            // Fila 1: Valores
            matriz[1][0] = contacto.getNombre();
            matriz[1][1] = contacto.getApellido();
            matriz[1][2] = contacto.getEmail();
            matriz[1][3] = contacto.getCelular();
            
            return ResponseEntity.ok(matriz);
        }
        return ResponseEntity.status(HttpStatus.NOT_FOUND).build();
    }
}