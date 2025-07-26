import { Module } from '@nestjs/common';
import { EstadoService } from '@core/estados/estadoService';
import { EstadoController } from './estadoController';
import { RedisService } from '@shared/cache/redis.service';

import  EstadosAdapter  from '@db/estadoAdapter';
import { EstadosPortToken } from './estado-port.token';

@Module({
  controllers: [EstadoController],
  providers: [
    EstadoService,
    {
      provide: EstadosPortToken, 
      useClass: EstadosAdapter,   
    },
    RedisService
  ],
  exports: [EstadoService],
})
export class EstadoModule {}
