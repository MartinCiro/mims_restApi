import PermisosPort from 'core/permisos/permisoPort';
import { PrismaClient } from '@prisma/client';
import { validarExistente, validarNoExistente } from 'api/utils/validaciones';
import { Injectable } from '@nestjs/common';
import { ForbiddenException } from '@nestjs/common';
import { PermisoData, PermisoDataUpdate, PermisoDataXid } from 'api/permisos/models/permiso.model';

const prisma = new PrismaClient();

@Injectable()
export default class PermisosAdapter implements PermisosPort {

  async crearPermisos(permisoData: PermisoData) {
    const { permisos, descripcion } = permisoData;

    try {
      const permisosCreados = [];

      for (const nombre of permisos) {
        const permiso = await prisma.permiso.upsert({
          where: { nombre },
          update: {},
          create: {
            nombre,
            descripcion,
          },
        });

        permisosCreados.push(permiso);
      }

      return {
        ok: true,
        message: "Permisos creados correctamente",
        permisos: permisosCreados,
      };
    } catch (error: any) {
      throw {
        ok: error.ok || false,
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
          nombre: true,
          descripcion: true,
        }
      });
      if (!permisos.length) throw new ForbiddenException("No se encontró el permiso");


      return permisos.map((permiso: { id: any; nombre: any; descripcion: any; }) => ({
        id: permiso.id,
        nombre: permiso.nombre,
        descripcion: permiso.descripcion
      }));
    } catch (error: any) {
      throw {
        ok: false,
        status_cod: 400,
        data: error.message || "Ocurrió un error consultando el permiso"
      };
    }
  }

  async obtenerPermisosXid(permisoData: PermisoDataXid) {
    try {
      const permiso = await prisma.permiso.findUnique({
        where: { id: Number(permisoData.id) },
        select: {
          id: true,
          nombre: true,
          descripcion: true,
        }
      });

      if (!permiso) throw new ForbiddenException("No se encontró el permiso");

      return {
        id: permiso.id,
        nombre: permiso.nombre,
        descripcion: permiso.descripcion
      };
    } catch (error: any) {
      throw {
        ok: false,
        status_cod: 400,
        data: error.message || "Ocurrió un error consultando el permiso",
      };
    }
  }

  async delPermiso(permisoData: PermisoDataXid) {
    try {
      const permiso = await prisma.permiso.delete({
        where: { id: Number(permisoData.id) },
      });

      return {
        ok: true,
        message: "Permiso eliminado correctamente",
        permiso: permiso.id,
      };
    } catch (error: any) {
      validarNoExistente(error.code, "El permiso solicitado");
      throw {
        ok: false,
        status_cod: 400,
        data: error.message || "Ocurrió un error eliminando el permiso",
      };
    }
  }

  async actualizaPermiso(permisoData: PermisoDataUpdate) {
    try {
      const { id, ...updates } = permisoData;

      const permisoActualizado = await prisma.permiso.update({
        where: { id: Number(id) },
        data: {
          ...updates,
        },
      });

      return {
        ok: true,
        message: "Permiso actualizado correctamente",
        permiso: permisoActualizado,
      };
    } catch (error: any) {
      validarExistente(error.code, "El permiso solicitado");
      throw {
        ok: false,
        status_cod: 400,
        data: error.message || "Ocurrió un error actualizando el permiso",
      };
    }
  }
}

