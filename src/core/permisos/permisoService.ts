import { Injectable, Inject } from '@nestjs/common';
import PermisosPort from './permisoPort';

import { PermisoData, PermisoDataUpdate, PermisoDataXid } from '@api/permisos/models/permiso.model';


@Injectable() 
export class PermisoService {
  constructor(
    @Inject('PermisosPort') private permisoPort: PermisosPort
  ) {}

  async obtenerPermisos() {
    return await this.permisoPort.obtenerPermisos();
  }

  async crearPermiso(permisoData: PermisoData) {
    return await this.permisoPort.crearPermisos(permisoData);
  }

  async upPermiso(permisoData: PermisoData) {
    return await this.permisoPort.actualizaPermiso(permisoData);
  }
}
