import { PermisoData, PermisoDataUpdate, PermisoDataXid } from 'api/permisos/models/permiso.model';
export default interface PermisosPort {
    crearPermisos(permisoData: PermisoData): Promise<any>;
    obtenerPermisos(): Promise<any>;
    obtenerPermisosXid(permisoData: PermisoDataXid): Promise<any>;
    delPermiso(permisoData: PermisoDataXid): Promise<any>;
    actualizaPermiso(permisoData: PermisoDataUpdate): Promise<any>;
}

