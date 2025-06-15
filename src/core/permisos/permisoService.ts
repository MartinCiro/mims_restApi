import { Injectable, Inject } from '@nestjs/common';
import PermisosPort from './permisoPort';

import { PermisoData, PermisoDataUpdate, PermisoDataXid } from 'api/permisos/models/permiso.model';

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

  async obtenerPermisoXid(permisoData: PermisoDataXid) {
    return await this.permisoPort.obtenerPermisosXid(permisoData);
  }

  async upPermiso(permisoData: PermisoDataUpdate) {
    return await this.permisoPort.actualizaPermiso(permisoData);
  }

  async delPermiso(permisoData: PermisoDataXid) {
    return await this.permisoPort.delPermiso(permisoData);
  }
}
