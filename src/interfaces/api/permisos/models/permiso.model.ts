export interface PermisoData {
  permisos: string[];
  descripcion?: string;
}

export interface PermisoDataXid {
  id: number | string;
}

export type PermisoDataUpdate = Partial<Omit<PermisoData, 'id'>> & PermisoDataXid;