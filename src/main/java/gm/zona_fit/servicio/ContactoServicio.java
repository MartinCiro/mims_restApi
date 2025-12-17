package gm.zona_fit.servicio;

import java.util.List;

import gm.zona_fit.modelo.Contacto;

public interface ContactoServicio {
    List<Contacto> listarContactos();
    Contacto buscarContactoPorId(Long id);
    List<Contacto> buscarContactoPorNombre(String nombre);
    Contacto guardarContacto(Contacto contacto);
    Contacto actualizarContacto(Long id, Contacto contacto);
    void eliminarContacto(Long id);
}