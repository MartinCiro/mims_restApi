import { UsuarioData, UsuarioDataUpdate, UsuarioDataXid } from '@api/usuarios/models/usuario.model';
export default interface UsuariosPort {
    obtenerUsuarios(): Promise<any>;
    obtenerUsuariosXid(usuarioData: UsuarioDataXid): Promise<any>;
    delUsuario(usuarioData: UsuarioDataXid): Promise<any>;
    crearUsuarios(usuarioData: UsuarioData): Promise<any>;
    actualizaUsuario(usuarioData: UsuarioDataUpdate): Promise<any>;
}