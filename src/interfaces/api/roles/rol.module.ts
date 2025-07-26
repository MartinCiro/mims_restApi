import { Module } from '@nestjs/common';
import { RolService } from '@core/roles/rolService';
import { RolController } from './rolController';
import { RedisService } from '@shared/cache/redis.service';

import  RolesAdapter  from '@db/rolAdapter';
import { RolesPortToken } from './rol-port.token';

@Module({
  controllers: [RolController],
  providers: [
    RolService,
    {
      provide: RolesPortToken, 
      useClass: RolesAdapter,   
    },
    RedisService
  ],
  exports: [RolService],
})
export class RolModule {}
