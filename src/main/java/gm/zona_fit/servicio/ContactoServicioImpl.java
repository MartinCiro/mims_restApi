package gm.zona_fit.servicio;

import java.util.List;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import gm.zona_fit.modelo.Contacto;
import gm.zona_fit.repositorio.ContactoRepositorio;

@Service
public class ContactoServicioImpl implements ContactoServicio {
    
    @Autowired
    private ContactoRepositorio contactoRepositorio;
    
    @Override
    public List<Contacto> listarContactos() {
        return contactoRepositorio.findAll();
    }
    
    @Override
    public Contacto buscarContactoPorId(Long id) {
        return contactoRepositorio.findById(id);
    }
    
    @Override
    public List<Contacto> buscarContactoPorNombre(String nombre) {
        return contactoRepositorio.findByNombre(nombre);
    }
    
    @Override
    public Contacto guardarContacto(Contacto contacto) {
        return contactoRepositorio.save(contacto);
    }
    
    @Override
    public Contacto actualizarContacto(Long id, Contacto contactoActualizado) {
        Contacto contactoExistente = contactoRepositorio.findById(id);
        if (contactoExistente != null) {
            contactoActualizado.setId(id);
            return contactoRepositorio.save(contactoActualizado);
        }
        return null;
    }
    
    @Override
    public void eliminarContacto(Long id) {
        contactoRepositorio.deleteById(id);
    }
}