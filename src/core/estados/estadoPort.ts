import { EstadoData, EstadoDataUpdate, EstadoDataXid } from '@api/estados/models/estado.model';
export default interface EstadosPort {
    crearEstados(usuarioData: EstadoData): Promise<any>;
    obtenerEstados(): Promise<any>;
    obtenerEstadosXid(usuarioData: EstadoDataXid): Promise<any>;
    delEstado(usuarioData: EstadoDataXid): Promise<any>;
    actualizaEstado(usuarioData: EstadoDataUpdate): Promise<any>;
}

