import RolesPort from '@core/roles/rolPort';
import { PrismaClient } from '@prisma/client';
import { validarExistente, validarNoExistente } from '@api/utils/validaciones';
import { Injectable } from '@nestjs/common';
import { ForbiddenException } from '@nestjs/common';
import { ResponseBody } from '@api/models/ResponseBody';
import { RolData, RolDataUpdate, RolDataXid } from '@api/roles/models/rol.model';

const prisma = new PrismaClient();

@Injectable()
export default class RolesAdapter implements RolesPort {

  async crearRoles(rolData: RolData) {
    try {
      // Buscar los permisos en la BD que coincidan con los proporcionados
      const permisosEncontrados = await prisma.permiso.findMany({
        where: {
          OR: rolData.permisos.map((permiso) => ({
            OR: [{ id: typeof permiso === "number" ? permiso : undefined }, { nombre: typeof permiso === "string" ? permiso : undefined }]
          }))
        }
      });

      // Si no se encontraron permisos, buscar el permiso "Lectura"
      if (permisosEncontrados.length === 0) {
        const permisoLectura = await prisma.permiso.findFirst({
          where: { nombre: "Lectura" }
        });

        if (!permisoLectura) throw new ForbiddenException("No se encontraron permisos válidos y tampoco existe el permiso 'Lectura'.");

        permisosEncontrados.push(permisoLectura);
      }

      // Crear el nuevo rol con los permisos encontrados
      const nuevoRol = await prisma.rol.create({
        data: {
          nombre: rolData.nombre,
          descripcion: rolData.descripcion
        }
      });

      await prisma.rolXPermiso.createMany({
        data: permisosEncontrados.map((permiso: { id: any }) => ({
          id_rol: nuevoRol.id,
          id_permiso: permiso.id
        }))
      });

      return nuevoRol;
    } catch (error: any) {
      const validacion = validarExistente(error.code, rolData.nombre);
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
        ok: false,
        status_cod: 400,
        data: error.data || "Ocurrió un error consultando el rol"
      };
    }
  }

  async obtenerRoles() {
    try {
      const roles = await prisma.rol.findMany({
        select: {
          id: true,
          nombre: true,
          descripcion: true,
          rolXPermiso: {
            select: {
              permiso: {
                select: {
                  nombre: true
                }
              }
            }
          }
        }
      });

      if (!roles.length) throw new ForbiddenException("No se encontraron roles");


      return roles.map((rol: { id: any; nombre: any; descripcion: any; rolXPermiso: any[]; }) => ({
        id: rol.id,
        nombre: rol.nombre,
        descripcion: rol.descripcion,
        permisos: rol.rolXPermiso.map((rp: { permiso: { nombre: any; }; }) => rp.permiso.nombre) // Agrupar nombres de permisos en una lista
      }));
    } catch (error: any) {
      throw {
        ok: false,
        status_cod: 400,
        data: error.message || "Ocurrió un error consultando el rol"
      };
    }
  }

  async obtenerRolesXid(rolData: RolDataXid) {
    try {
      const rol = await prisma.rol.findUnique({
        where: { id: Number(rolData.id) },
        select: {
          id: true,
          nombre: true,
          descripcion: true,
          rolXPermiso: {
            select: {
              permiso: {
                select: {
                  nombre: true
                }
              }
            }
          }
        }
      });

      if (!rol) throw new ForbiddenException("El rol solicitado no existe en la base de datos");
      return {
        id: rol.id,
        nombre: rol.nombre,
        descripcion: rol.descripcion,
        permisos: rol.rolXPermiso.map((rp: { permiso: { nombre: any; }; }) => rp.permiso.nombre)
      };
    } catch (error: any) {
      throw {
        ok: false,
        status_cod: 400,
        data: error.message || "Ocurrió un error consultando el rol",
      };
    }
  }


  async delRol(rolData: RolDataXid) {
    try {
      const rol = await prisma.rol.delete({
        where: { id: Number(rolData.id) },
      });

      return {
        ok: true,
        message: "Rol eliminado correctamente",
        rol: rol.id,
      };
    } catch (error: any) {
      validarNoExistente(error.code, "El rol solicitado");
      throw {
        ok: false,
        status_cod: 400,
        data: error.message || "Ocurrió un error eliminando el rol",
      };
    }
  }

  async actualizaRol(rolData: RolDataUpdate) {
    try {
      const { id, permisos, ...updates } = rolData;

      // Verificar si el rol existe
      const rolExistente = await prisma.rol.findUnique({
        where: { id: Number(id) }
      });

      if (!rolExistente) {
        throw {
          ok: false,
          status_cod: 400,
          message: "El rol solicitado no existe en la base de datos",
        }
      }

      // Eliminar permisos asociados al rol en RolXPermiso
      await prisma.rolXPermiso.deleteMany({
        where: { id_rol: Number(id) }
      });

      // Actualizar el rol con los nuevos datos
      const rolActualizado = await prisma.rol.update({
        where: { id: Number(id) },
        data: updates
      });

      // Si hay permisos nuevos, insertarlos en RolXPermiso
      if (permisos && permisos.length > 0) {
        const permisosEncontrados = await prisma.permiso.findMany({
          where: {
            OR: permisos.map((permiso) => ({
              id: typeof permiso === "number" ? permiso : undefined,
              nombre: typeof permiso === "string" ? permiso : undefined
            }))
          }
        });

        if (permisosEncontrados.length > 0) {
          await prisma.rolXPermiso.createMany({
            data: permisosEncontrados.map((permiso: { id: any; }) => ({
              id_rol: Number(id),
              id_permiso: permiso.id
            }))
          });
        }
      }

      return {
        ok: true,
        message: "Rol actualizado correctamente",
        rol: rolActualizado
      };
    } catch (error: any) {
      validarExistente(error.code, "El rol solicitado");
      throw {
        ok: false,
        status_cod: 400,
        data: error.message || error.data || "Ocurrió un error actualizando el rol",
      };
    }
  }
}

