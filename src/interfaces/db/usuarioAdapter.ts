import UsuariosPort from '@core/usuarios/usuarioPort';
import { Usuario } from '@core/auth/entities/Usuario';
import { Prisma, PrismaClient } from '@prisma/client';
import { validarExistente, validarNoExistente } from '@utils/validaciones';
import { Injectable, ForbiddenException } from '@nestjs/common';
import { UsuarioData, UsuarioDataUpdate, UsuarioDataXid } from '@api/usuarios/models/usuario.model';

const prisma = new PrismaClient();

@Injectable()
export default class UsuariosAdapter implements UsuariosPort {

  async crearUsuarios(usuarioData: UsuarioData) {
    try {
      // Obtener o crear estado "Activo"
      const estado = await prisma.estado.upsert({
        where: { nombre_estado: 'Activo' },
        update: {},
        create: { nombre_estado: 'Activo', descripcion: 'Usuario activo' },
        select: { id: true }
      });

      // Determinar o crear el rol
      let idRol: number;
      if (usuarioData.id_rol) {
        const rol = await prisma.rol.findUnique({
          where: { id: Number(usuarioData.id_rol) },
          select: { id: true }
        });

        if (!rol) throw new ForbiddenException(`El rol con ID ${usuarioData.id_rol} no existe`);
        idRol = rol.id;
      } else {
        // Rol invitado
        const rolInvitado = await prisma.rol.upsert({
          where: { nombre_rol: 'Invitado' },
          update: {},
          create: {
            nombre_rol: 'Invitado',
            descripcion: 'Rol para usuarios invitados'
          },
          select: { id: true }
        });

        const permisoLee = await prisma.permiso.upsert({
          where: { nombre_permiso: 'anuncios:Lee' },
          update: {},
          create: {
            nombre_permiso: 'anuncios:Lee',
            descripcion: 'Permiso de solo lectura en anuncios'
          },
          select: { id: true }
        });

        // Verificar si la relación ya existe
        const existeRelacion = await prisma.rolXPermiso.findUnique({
          where: {
            id_rol_id_permiso: {
              id_rol: rolInvitado.id,
              id_permiso: permisoLee.id
            }
          }
        });

        if (!existeRelacion) {
          await prisma.rolXPermiso.create({
            data: {
              id_rol: rolInvitado.id,
              id_permiso: permisoLee.id
            }
          });
        }

        idRol = rolInvitado.id;
      }

      // Encriptar la contraseña
      const user = new Usuario(
        usuarioData.nom_user,
        usuarioData.passwd,
        idRol,
        estado.id
      );

      // Crear el usuario
      const nuevoUsuario = await prisma.usuario.create({
        data: {
          documento: usuarioData.numero_documento,
          nombres: usuarioData.nombres,
          apellidos: usuarioData.apellido,
          username: usuarioData.nom_user,
          pass: user.getEncryptedPassword(),
          id_estado: estado.id,
          id_rol: idRol
        },
        select: {
          id: true,
          documento: true,
          nombres: true,
          apellidos: true,
          username: true
        }
      });

      return nuevoUsuario;
    } catch (error: any) {
      const validacion = validarExistente(error.code, error.meta?.target);
      if (!validacion.ok) {
        throw {
          ok: validacion.ok,
          status_cod: 409,
          data: validacion.data,
        };
      }

      const resultado = error.meta?.target?.[0] || "valor";
      const valNoExistente = validarNoExistente(error.code, `El ${resultado} asignado`);

      if (!valNoExistente.ok) {
        throw {
          ok: false,
          status_cod: 409,
          data: valNoExistente.data
        };
      }

      throw {
        ok: false,
        status_cod: 400,
        data: error.data || "Ocurrió un error creando el usuario"
      };
    }
  }

  async obtenerUsuarios() {
    try {
      const usuarios = await prisma.usuario.findMany({
        select: {
          nombres: true,
          apellidos: true,
          documento: true,
          username: true,
          rol: {
            select: {
              nombre_rol: true, // Cambiado de 'nombre' a 'nombre_rol'
              rolXPermiso: {
                select: {
                  permiso: {
                    select: {
                      nombre_permiso: true // Cambiado de 'nombre' a 'nombre_permiso'
                    }
                  }
                }
              }
            }
          },
          estado: {
            select: { nombre_estado: true } // Cambiado de 'nombre' a 'nombre_estado'
          }
        }
      });

      if (usuarios.length === 0) {
        throw {
          ok: true,
          status_cod: 200,
          data: "No se han encontrado usuarios"
        };
      }

      return usuarios.map((usuario) => {
        const permisosPlano = usuario.rol.rolXPermiso.map((rp) => rp.permiso.nombre_permiso);
        const permisosAgrupados = permisosPlano.reduce((grupos: Record<string, string[]>, permiso: string) => {
          const [categoria, accion] = permiso.split(':');
          if (!categoria || !accion) return grupos; // Validación para evitar errores

          if (!grupos[categoria]) grupos[categoria] = [];
          grupos[categoria].push(accion);
          return grupos;
        }, {});

        return {
          nombres: usuario.nombres,
          apellidos: usuario.apellidos, // Cambiado de 'apellido' a 'apellidos'
          estado: usuario.estado.nombre_estado, // Cambiado a 'nombre_estado'
          rol: usuario.rol.nombre_rol, // Cambiado a 'nombre_rol'
          documento: usuario.documento,
          username: usuario.username, // Cambiado de 'nom_user' a 'username'
          permisos: permisosAgrupados
          // Nota: Los siguientes campos no existen en tu modelo Prisma actual:
          // email, info_perfil, num_contacto, fecha_nacimiento, fecha_registro
        };
      });
    } catch (error: any) {
      throw {
        ok: error.ok || false,
        status_cod: error.status_cod || 400,
        data: error.message || error.data || "Ocurrió un error consultando el usuario"
      };
    }
  }

  async obtenerUsuariosXid(usuarioData: UsuarioDataXid) {
    try {
      const usuario = await prisma.usuario.findUnique({
        where: { documento: usuarioData.numero_documento.toString() },
        select: {
          documento: true,
          nombres: true,
          estado: {
            select: { nombre_estado: true } // Cambiado de 'nombre' a 'nombre_estado'
          },
          rol: {
            select: { nombre_rol: true } // Cambiado de 'nombre' a 'nombre_rol'
          }
        }
      });

      if (!usuario) {
        throw {
          ok: true,
          status_cod: 200,
          data: "El usuario solicitado no existe en la base de datos",
        };
      }

      return {
        id: usuario.documento,
        nombres: usuario.nombres,
        estado_usuario: usuario.estado.nombre_estado, // Cambiado a 'nombre_estado'
        rol: usuario.rol.nombre_rol // Cambiado a 'nombre_rol'
      };
    } catch (error: any) {
      throw {
        ok: error.ok || false,
        status_cod: error.status_cod || 400,
        data: error.message || error.data || "Ocurrió un error consultando el usuario",
      };
    }
  }

  async delUsuario(usuarioData: UsuarioDataXid) {
    try {
      const usuario = await prisma.usuario.delete({
        where: { documento: usuarioData.numero_documento.toString() },
        select: {
          documento: true,
          nombres: true,
          apellidos: true
        }
      });

      return {
        ok: true,
        message: "Usuario eliminado correctamente",
        usuario: usuario.documento,
      };
    } catch (error: any) {
      if (error.code === "P2025") {
        throw {
          ok: false,
          status_cod: 404, // Cambiado a 404 (Not Found) que es más apropiado
          data: "El usuario solicitado no existe en la base de datos",
        };
      }

      // Manejar error de violación de clave foránea
      if (error.code === "P2003") {
        throw {
          ok: false,
          status_cod: 409,
          data: "No se puede eliminar el usuario porque tiene registros relacionados",
        };
      }

      throw {
        ok: false,
        status_cod: 400,
        data: error.message || "Ocurrió un error eliminando el usuario",
      };
    }
  }

  async actualizaUsuario(usuarioData: UsuarioDataUpdate) {
    try {
      const { numero_documento, ...updates } = usuarioData;

      // Preparar datos para actualización según tu modelo Prisma
      const dataToUpdate: any = {};

      // Mapear campos básicos
      if (updates.nombres !== undefined) dataToUpdate.nombres = updates.nombres;
      if (updates.apellido !== undefined) dataToUpdate.apellidos = updates.apellido;
      if (updates.nom_user !== undefined) dataToUpdate.username = updates.nom_user;
      if (updates.passwd !== undefined) {
        // Encriptar nueva contraseña si se proporciona
        const user = new Usuario(updates.nom_user || '', updates.passwd);
        dataToUpdate.pass = user.getEncryptedPassword();
      }
      if (updates.email !== undefined) {
        // Nota: El campo email no existe en tu modelo actual
        // dataToUpdate.email = updates.email;
        console.warn('Campo email no está disponible en el modelo actual');
      }

      // Manejar relaciones
      if (updates.estado_id !== undefined) {
        dataToUpdate.id_estado = Number(updates.estado_id);
      }

      if (updates.id_rol !== undefined) {
        dataToUpdate.id_rol = Number(updates.id_rol);
      }

      // Nota: Los siguientes campos no existen en tu modelo Prisma actual:
      // - fecha_nacimiento
      // - numero_contacto

      const usuarioActualizado = await prisma.usuario.update({
        where: { documento: numero_documento.toString() },
        data: dataToUpdate,
        select: {
          documento: true,
          nombres: true,
          apellidos: true,
          username: true,
          rol: {
            select: { nombre_rol: true }
          },
          estado: {
            select: { nombre_estado: true }
          }
        }
      });

      return {
        ok: true,
        message: "Usuario actualizado correctamente",
        usuario: usuarioActualizado,
      };
    } catch (error: any) {
      // Manejar usuario no encontrado
      if (error.code === "P2025") {
        throw {
          ok: false,
          status_cod: 404,
          data: "El usuario no existe en la base de datos",
        };
      }

      // Manejar violación de constraints únicos
      if (error.code === "P2002") {
        const campo = error.meta?.target?.[0] || 'campo';
        throw {
          ok: false,
          status_cod: 409,
          data: `El ${campo} ya está en uso por otro usuario`,
        };
      }

      // Manejar relaciones no existentes
      if (error.code === "P2003") {
        throw {
          ok: false,
          status_cod: 409,
          data: "El rol o estado asignado no existe",
        };
      }

      throw {
        ok: false,
        status_cod: 400,
        data: error.message || "Ocurrió un error actualizando el usuario",
      };
    }
  }

}

