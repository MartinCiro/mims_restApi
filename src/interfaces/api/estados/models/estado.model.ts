export interface EstadoData {
  nombre: string;
  descripcion?: string;
}

export interface EstadoDataXid {
  id: number | string;
}

export type EstadoDataUpdate = Partial<Omit<EstadoData, 'id'>> & EstadoDataXid;