import { Injectable, Inject } from '@nestjs/common';
import RolesPort from './rolPort';
import { RolData, RolDataUpdate, RolDataXid } from '@api/roles/models/rol.model';

@Injectable() 
export class RolService {
  constructor(
    @Inject('RolesPort') private rolPort: RolesPort
  ) {}

  async obtenerRoles() {
    return await this.rolPort.obtenerRoles();
  }

  async crearRol(rolData: RolData) {
    return await this.rolPort.crearRoles(rolData);
  }

  async obtenerRolXid(rolData: RolDataXid) {
    return await this.rolPort.obtenerRolesXid(rolData);
  }

  async upRol(rolData: RolDataUpdate) {
    return await this.rolPort.actualizaRol(rolData);
  }

  async delRol(rolData: RolDataXid) {
    return await this.rolPort.delRol(rolData);
  }
}
