export interface UsuarioData {
  nombres: string;
  id_rol?: number;
  apellido: string;
  numero_documento: string;
  estado_id?: number;
  nom_user: string;
  passwd: string;
}

export interface UsuarioDataXid {
  numero_documento: number | string;
}

export type UsuarioDataUpdate = Partial<Omit<UsuarioData, 'numero_documento'>> & UsuarioDataXid;