import { Injectable, Inject } from '@nestjs/common';
import EstadosPort from './estadoPort';
import { EstadoData, EstadoDataUpdate, EstadoDataXid } from '@api/estados/models/estado.model';

@Injectable() 
export class EstadoService {
  constructor(
    @Inject('EstadosPort') private estadoPort: EstadosPort
  ) {}

  async obtenerEstados() {
    return await this.estadoPort.obtenerEstados();
  }

  async crearEstado(estadoData: EstadoData) {
    return await this.estadoPort.crearEstados(estadoData);
  }

  async obtenerEstadoXid(estadoData: EstadoDataXid) {
    return await this.estadoPort.obtenerEstadosXid(estadoData);
  }

  async upEstado(estadoData: EstadoDataUpdate) {
    return await this.estadoPort.actualizaEstado(estadoData);
  }

  async delEstado(estadoData: EstadoDataXid) {
    return await this.estadoPort.delEstado(estadoData);
  }
}
