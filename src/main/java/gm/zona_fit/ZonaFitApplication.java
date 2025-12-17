package gm.zona_fit;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.CommandLineRunner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

import gm.zona_fit.modelo.Contacto;
import gm.zona_fit.servicio.ContactoServicio;

@SpringBootApplication
public class ZonaFitApplication implements CommandLineRunner {
    
    private static final Logger logger = LoggerFactory.getLogger(ZonaFitApplication.class);
    
    @Autowired
    private ContactoServicio contactoServicio;
    
    @Override
    public void run(String... args) throws Exception {
        // Datos de prueba iniciales
        cargarDatosPrueba();
        logger.info("=== API REST ZonaFit iniciada ===");
        logger.info("Disponible en: http://localhost:8080");
        logger.info("Endpoints disponibles:");
        logger.info("  GET    /api/contactos");
        logger.info("  GET    /api/contactos/{id}");
        logger.info("  POST   /api/contactos");
        logger.info("  PUT    /api/contactos/{id}");
        logger.info("  DELETE /api/contactos/{id}");
        logger.info("  GET    /api/contactos/buscar?nombre=valor");
    }
    
    private void cargarDatosPrueba() {
        // Agregar algunos contactos de ejemplo
        Contacto contacto1 = new Contacto("Juan", "Pérez", "juan@email.com", "123456789");
        Contacto contacto2 = new Contacto("María", "González", "maria@email.com", "987654321");
        Contacto contacto3 = new Contacto("Carlos", "López", "carlos@email.com", "555123456");
        
        contactoServicio.guardarContacto(contacto1);
        contactoServicio.guardarContacto(contacto2);
        contactoServicio.guardarContacto(contacto3);
        
        logger.info("Datos de prueba cargados: {} contactos", contactoServicio.listarContactos().size());
    }
    
    public static void main(String[] args) {
        SpringApplication.run(ZonaFitApplication.class, args);
    }
}