import { Module } from '@nestjs/common';
import { PermisoService } from '@core/permisos/permisoService';
import { PermisoController } from './permisoController';
import { RedisService } from '@shared/cache/redis.service';

import  PermisosAdapter  from '@db/permisoAdapter';
import { PermisosPortToken } from './permiso-port.token';

@Module({
  controllers: [PermisoController],
  providers: [
    PermisoService,
    {
      provide: PermisosPortToken, 
      useClass: PermisosAdapter,   
    },
    RedisService
  ],
  exports: [PermisoService],
})
export class PermisoModule {}
