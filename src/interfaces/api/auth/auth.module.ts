import { Module } from '@nestjs/common';
import AuthService from '@core/auth/authService';
import { AuthController } from '@api/auth/authController';

import  AuthAdapter  from '@db/authAdapter';
import { AuthPort } from '@api/auth/auth-port.token';
import { CacheModule } from '@shared/cache/cache.module';

@Module({
  controllers: [AuthController],
  imports: [CacheModule],
  providers: [
    AuthService,
    {
      provide: AuthPort, 
      useClass: AuthAdapter,   
    },
  ],
  exports: [AuthService],
})
export class AuthModule {}
