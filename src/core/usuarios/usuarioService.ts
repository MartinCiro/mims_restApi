import { Injectable, Inject } from '@nestjs/common';
import UsuariosPort from './usuarioPort';
import { UsuarioData, UsuarioDataUpdate, UsuarioDataXid } from '@api/usuarios/models/usuario.model';

@Injectable() 
export class UsuarioService {
  constructor(
    @Inject('UsuariosPort') private usuarioPort: UsuariosPort
  ) {}

  async obtenerUsuarios() {
    return await this.usuarioPort.obtenerUsuarios();
  }

  async crearUsuario(usuarioData: UsuarioData) {
    return await this.usuarioPort.crearUsuarios(usuarioData);
  }

  async obtenerUsuarioXid(usuarioData: UsuarioDataXid) {
    return await this.usuarioPort.obtenerUsuariosXid(usuarioData);
  }

  async upUsuario(usuarioData: UsuarioDataUpdate) {
    return await this.usuarioPort.actualizaUsuario(usuarioData);
  }

  async delUsuario(usuarioData: UsuarioDataXid) {
    return await this.usuarioPort.delUsuario(usuarioData);
  }
}
