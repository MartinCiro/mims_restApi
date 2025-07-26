import { PrismaClient } from '@prisma/client';
import AuthPort from '@core/auth/authPort';
import { validarExistente } from '@api/utils/validaciones';

const prisma = new PrismaClient();

class AuthAdapter implements AuthPort {

  async retrieveUser(authData: { username: string }) {
    try {
      const usuario = await prisma.usuario.findUnique({
        where: { username: authData.username },
        select: {
          id: true,
          nombres: true,
          apellidos: true,
          username: true,
          pass: true,
          id_rol: true,
          estado: {  
            select: {
              nombre: true
            }
          },
          rol: { 
            select: {
              nombre: true,
              rolXPermiso: {
                select: {
                  permiso: { select: { nombre: true } }
                }
              }
            }
          }
        },
      });
      if (!usuario) return null;

      return {
        id_user: usuario.id,
        usuario: usuario.nombres + " " + usuario.apellidos,
        password: usuario.pass,
        username: usuario.username,
        id_rol: usuario.id_rol,
        estado: usuario.estado.nombre,
        rol: usuario.rol.nombre,
        permisos: usuario.rol?.rolXPermiso.map((rp: { permiso: { nombre: any; }; }) => rp.permiso.nombre) || []
      };

    } catch (error: any) {
      console.error("Error en auth del usuario:", error);
      const validacion = validarExistente(error.code, authData.username);
      if (!validacion.ok) {
        throw {
          ok: validacion.ok,
          status_cod: 409,
          data: validacion.data,
        };
      }
      throw {
        ok: error.ok || false,
        status_cod: error.status_cod || 400,
        data: error.message || "Ocurrió un error consultando el usuario"
      };
    }
  }
}

export default AuthAdapter;
