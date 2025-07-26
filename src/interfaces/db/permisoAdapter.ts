import { PrismaClient, Permiso } from '@prisma/client';
import { Injectable } from '@nestjs/common';
import PermisosPort from '@core/permisos/permisoPort';
import { validarExistente, capitalize } from '@utils/validaciones';

const prisma = new PrismaClient();

@Injectable()
export default class PermisosAdapter implements PermisosPort {


async crearPermisos(permisoData: { permisos: string[]; descripcion?: string | null }) {
    const { permisos, descripcion } = permisoData;

    try {
      // Define explícitamente el tipo del array
      const permisosCreados: Permiso[] = [];  // <--- Aquí está el fix

      for (const nombre of permisos) {
        const permisoExistente = await prisma.permiso.findUnique({
          where: { nombre },
        });

        if (permisoExistente) continue;
        
        const nuevoPermiso = await prisma.permiso.create({
          data: {
            nombre,
            descripcion,
          },
        });

        permisosCreados.push(nuevoPermiso);  // <--- Ahora no habrá error
      }

      return {
        ok: true,
        message: "Permisos creados correctamente",
        permisos: permisosCreados,
      };
    } catch (error: any) {
      throw {
        ok: false,
        status_cod: error.status_cod || 400,
        data: error.data || "Error inesperado creando permisos",
      };
    }
}

  async obtenerPermisos() {
    try {
      const permisos = await prisma.permiso.findMany({
        select: {
          id: true,
          nombre: true,        // Ejemplo: "anuncios:crear"
          descripcion: true    // Ejemplo: "Permiso para anuncios"
        }
      });

      if (permisos.length === 0) {
        throw {
          ok: true,
          status_cod: 200,
          data: "No se han encontrado ningun permiso"
        };
      }

      const agrupados: Record<string, { descripcion: string, acciones: Set<string> }> = {};

      for (const permiso of permisos) {
        const [entidad, accionRaw] = permiso.nombre.split(':');
        const accion = capitalize(accionRaw);

        if (!agrupados[entidad]) {
          agrupados[entidad] = {
            descripcion: permiso.descripcion ?? "",
            acciones: new Set([accion])
          };
        } else {
          agrupados[entidad].acciones.add(accion);
        }
      }

      // Convertir a array con el formato deseado
      const resultado = Object.entries(agrupados).map(
        ([entidad, grupo]: [string, { descripcion: string; acciones: Set<string> }]) => ({
          descripcion: grupo.descripcion,
          [entidad]: Array.from(grupo.acciones)
        })
      );

      return resultado;

    } catch (error: any) {
      throw {
        ok: error.ok || false,
        status_cod: error.status_cod || 400,
        data: error.message || error.data || "Ocurrió un error consultando los permisos"
      };
    }
  }

  async actualizaPermiso(permisoData: { descripcion?: string; permisos: string[] }) {
    const { descripcion, permisos } = permisoData;

    try {
      // 1. Obtener el prefijo (nombre de la entidad, ej: "cedula")
      const entidad = permisos[0]?.split(":")[0];
      if (!entidad) throw new Error("Nombre de permiso no válido");

      // 2. Obtener todos los permisos actuales de esa entidad
      const permisosExistentes = await prisma.permiso.findMany({
        where: {
          nombre: {
            startsWith: `${entidad}:`
          }
        }
      });

      // 3. Extraer nombres actuales y los que deben eliminarse
      const nombresActuales = permisosExistentes.map(p => p.nombre);
      const permisosAEliminar = nombresActuales.filter(nombre => !permisos.some(p => p === nombre));

      // 4. Eliminar los permisos que ya no están
      await Promise.all(
        permisosAEliminar.map(nombre =>
          prisma.permiso.delete({ where: { nombre } })
        )
      );

      // 5. Actualizar los permisos que se mantienen (no crear nuevos)
      await Promise.all(
        permisos.map(async (nombre) => {
          return await prisma.permiso.update({
            where: { nombre },
            data: { descripcion }
          });
        })
      );

      return {
        ok: true,
        message: `Los permisos de ${entidad} han sido actualizados`,
      };
    } catch (error: any) {
      const validacion = validarExistente(error.code, error.meta?.target);
      if (!validacion.ok) {
        throw {
          ok: validacion.ok,
          status_cod: 409,
          data: validacion.data,
        };
      }
      throw {
        ok: false,
        status_cod: 400,
        data: error.message || "Ocurrió un error actualizando los permisos",
      };
    }
  }
}
