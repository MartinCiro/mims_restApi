import { Genero } from '@prisma/client'
export interface ProgramacionData {
  rama: Genero;
  fecha: string;
  lugar: number;
  competencia: string;
  equipoLocal: number;
  equipoVisitante: number;
  torneo: number;
  categoria?: number;
  hora: string;
}

/* 
  "rama": "M",
  "fecha": "10/06/2025",
  "lugar": 1,
  "competencia": "fecha 2", 
  "equipoLocal": "1",
  "equipoVisitante": "10",
  "torneo": "1",
  "categoria": "2011 M",
  "hora": "10:22 p.m.",
  "dia": "martes",
  "id": 5,
*/

export interface ProgramacionDataXid {
  id: number | string;
}

export type ProgramacionDataUpdate = Partial<Omit<ProgramacionData, 'id'>> & ProgramacionDataXid;