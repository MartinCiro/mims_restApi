export interface RolData {
  nombre: string;
  descripcion?: string;
  permisos: (string | number)[];
}

export interface RolDataXid {
  id: number | string;
}

export type RolDataUpdate = Partial<Omit<RolData, 'id'>> & RolDataXid;