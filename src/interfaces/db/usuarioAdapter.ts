import UsuariosPort from 'core/usuarios/usuarioPort';
import { Usuario } from 'core/auth/entities/Usuario';
import { Prisma, PrismaClient } from '@prisma/client';
import { validarExistente, validarNoExistente } from 'api/utils/validaciones';
import { Injectable } from '@nestjs/common';
const prisma = new PrismaClient();

@Injectable()
export default class UsuariosAdapter implements UsuariosPort {

  async crearUsuarios(usuarioData: {
    username: string;
    nombres: string;
    apellidos: string;
    pass: string;
    id_rol?: number | string;
  }) {
    try {
      const result = await prisma.$transaction(async (tx: PrismaClient) => {
        // 1. Estado activo y rol en paralelo
        const [estado, rolData] = await Promise.all([
          tx.estado.upsert({
            where: { nombre: 'activo' },
            update: {},
            create: { nombre: 'activo' },
            select: { id: true }
          }),
          (async () => {
            if (usuarioData.id_rol) {
              const rol = await tx.rol.findUnique({
                where: { id: Number(usuarioData.id_rol) },
                select: { id: true }
              });
              if (!rol) throw new Error('ROL_NO_EXISTE');
              return { rol, permisos: [] };
            } else {
              return this.asignarRolInvitado(tx);
            }
          })()
        ]);

        // 2. Crear usuario
        const user = new Usuario(usuarioData.username, usuarioData.pass);
        return await tx.usuario.create({
          data: {
            username: usuarioData.username,
            nombres: usuarioData.nombres,
            apellidos: usuarioData.apellidos,
            pass: user.getEncryptedPassword(),
            id_estado: estado.id,
            id_rol: rolData.rol.id
          },
          select: { id: true }
        });
      });

      return result;
    } catch (error: any) {
      // Manejo de errores
      if (error.message === 'ROL_NO_EXISTE') {
        throw {
          ok: false,
          status_cod: 400,
          data: "El rol especificado no existe"
        };
      }

      const validacion = validarExistente(error.code, usuarioData.username);
      if (!validacion.ok) {
        throw {
          ok: false,
          status_cod: 409,
          data: validacion.data
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
        ok: error.ok || false,
        status_cod: error.status_cod || 400,
        data: error.message || error.data || "Ocurrió un error consultando el usuario"
      };
    }
  }

  private async asignarRolInvitado(tx: Prisma.TransactionClient) {
    // Crear rol invitado y permisos en una sola transacción
    const [rolInvitado, permisoLee] = await Promise.all([
      tx.rol.upsert({
        where: { nombre: 'Invitado' },
        update: {},
        create: { nombre: 'Invitado' },
        select: { id: true }
      }),
      tx.permiso.upsert({
        where: { nombre: 'anuncios:Lee' },
        update: {},
        create: {
          nombre: 'anuncios:Lee',
          descripcion: 'Permiso de solo lectura en anuncios'
        },
        select: { id: true }
      })
    ]);

    await tx.rolXPermiso.create({
      data: {
        id_rol: rolInvitado.id,
        id_permiso: permisoLee.id
      },
      skipDuplicates: true
    });

    return { rol: rolInvitado, permisos: [permisoLee] };
  }

  async obtenerUsuarios() {
    try {
      const usuarios = await prisma.usuario.findMany({
        select: {
          id: true,
          nombres: true,
          apellidos: true,
          username: true,
          estado: {
            select: { nombre: true }
          },
          rol: {
            select: {
              nombre: true
            }
          }
        }
      });

      if (!usuarios.length) {
        throw {
          ok: false,
          status_cod: 404,
          data: "No se encontraron datos"
        };
      }

      return usuarios.map((usuario:
        {
          nombres: any; 
          apellidos: any;
          username: any; 
          rol: { nombre: any; };
          estado: { nombre: any; };
        }) => ({
          nombres: usuario.nombres,
          apellidos: usuario.apellidos,
          nom_user: usuario.username,
          rol: usuario.rol.nombre,
          estado: usuario.estado.nombre,
        }));
    } catch (error: any) {
      throw {
        ok: error.ok || false,
        status_cod: error.status_cod || 400,
        data: error.message || "Ocurrió un error consultando el usuario"
      };
    }
  }

  async obtenerUsuariosXid(usuarioData: { id: string | number; }) {
    try {
      const usuario = await prisma.usuario.findUnique({
        where: { id: Number(usuarioData.id) },
        select: {
          id: true,
          nombres: true,
          estado: {
            select: { nombre: true }
          },
          rol: {
            select: { nombre: true }
          }
        }
      });

      if (!usuario) {
        throw {
          ok: false,
          status_cod: 409,
          data: "El usuario solicitado no existe en la base de datos",
        };
      }

      return {
        id: usuario.id,
        nombres: usuario.nombres,
        estado: usuario.estado.nombre,
        rol: usuario.rol.nombre
      };
    } catch (error: any) {
      throw {
        ok: error.ok || false,
        status_cod: error.status_cod || 400,
        data: error.message || "Ocurrió un error consultando el usuario"
      };
    }
  }

  async delUsuario(usuarioData: { id: string }) {
    try {
      const usuario = await prisma.usuario.delete({
        where: { id: Number(usuarioData.id) },
      });

      return {
        ok: true,
        message: "Usuario eliminado correctamente",
        usuario: usuario.id,
      };
    } catch (error: any) {
      if (error.code === "P2025") {
        throw {
          ok: false,
          status_cod: 409,
          data: "El usuario solicitado no existe en la base de datos",
        };
      }
      throw {
        ok: error.ok || false,
        status_cod: error.status_cod || 400,
        data: error.message || "Ocurrió un error consultando el usuario"
      };
    }
  }

  async actualizaUsuario(usuarioData: {
    nombres?: string;
    apellidos?: string;
    id_estado?: number | string;
    id_rol?: number | string;
    id: number | string;
  }) {
    try {
      const { id, id_estado, id_rol, ...updates } = usuarioData;

      const usuarioActualizado = await prisma.usuario.update({
        where: { id: Number(id) },
        data: {
          ...updates,
          id_estado: id_estado !== undefined ? Number(id_estado) : undefined,
          id_rol: id_rol !== undefined ? Number(id_rol) : undefined,
        },
      });

      return {
        ok: true,
        message: "Usuario actualizado correctamente",
        usuario: usuarioActualizado,
      };
    } catch (error: any) {
      validarExistente(error.code, "El usuario solicitado");
      throw {
        ok: error.ok || false,
        status_cod: error.status_cod || 400,
        data: error.message || "Ocurrió un error consultando el usuario"
      };
    }
  }
}