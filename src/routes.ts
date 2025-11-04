import { join } from 'path';
import { Module } from '@nestjs/common';
import { AppController } from '@src/app.controller';
import { ServeStaticModule } from '@nestjs/serve-static';

import { RolModule } from '@api/roles/rol.module';
import { AuthModule } from '@api/auth/auth.module';
import { EstadoModule } from '@api/estados/estado.module';
import { PermisoModule } from '@api/permisos/permiso.module';

@Module({
  imports: [
    ServeStaticModule.forRoot({
      rootPath: join(__dirname, '..', 'public'),
      serveRoot: '/api-docs'
    }),
    PermisoModule,
    RolModule,
    EstadoModule,
    AuthModule,
  ],
  controllers: [AppController],
})
export class AppModule {}
