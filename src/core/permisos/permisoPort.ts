import { PermisoData } from '@api/permisos/models/permiso.model';
export default interface PermisosPort {
    crearPermisos(permisoData: PermisoData): Promise<any>;
    obtenerPermisos(): Promise<any>;
    actualizaPermiso(permisoData: PermisoData): Promise<any>;
}

