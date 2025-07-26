import { RolData, RolDataUpdate, RolDataXid } from '@api/roles/models/rol.model';
export default interface RolesPort {
    crearRoles(rolData: RolData): Promise<any>;
    obtenerRoles(): Promise<any>;
    obtenerRolesXid(rolData: RolDataXid): Promise<any>;
    delRol(rolData: RolDataXid): Promise<any>;
    actualizaRol(rolData: RolDataUpdate): Promise<any>;
}

