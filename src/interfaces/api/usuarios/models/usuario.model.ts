export interface UsuarioData {
  nombres: string;
  id_rol?: number;
  apellido: string;
  numero_documento: string;
  email: string;
  estado_id?: number;
  info_perfil?: string;
  nom_user: string;
  numero_contacto?: string;
  fecha_nacimiento: string;
  passwd: string;
}

export interface UsuarioDataXid {
  numero_documento: number | string;
}

export type UsuarioDataUpdate = Partial<Omit<UsuarioData, 'numero_documento'>> & UsuarioDataXid;