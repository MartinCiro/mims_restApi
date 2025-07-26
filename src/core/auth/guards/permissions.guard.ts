import { Injectable, CanActivate, ExecutionContext, ForbiddenException } from '@nestjs/common';
import { Reflector } from '@nestjs/core';
import { RedisService } from '@shared/cache/redis.service';
import { ResponseBody } from 'src/interfaces/api/models/ResponseBody';
import { HttpException } from '@nestjs/common';

@Injectable()
export class PermissionsGuard implements CanActivate {
  constructor(private reflector: Reflector, private redisService: RedisService) {}

  async canActivate(context: ExecutionContext): Promise<boolean> {
    const requiredPermissions = this.reflector.get<string[]>('permissions', context.getHandler()) || [];
    if (!requiredPermissions.length) return true; 

    const request = context.switchToHttp().getRequest();
    const userId = request.user?.userInfo?.id_user;
    
    if (!userId) throw new ForbiddenException('No autenticado');
    
    let userData = await this.redisService.get(`user:${userId}`);
    if(typeof userData === 'string') userData = JSON.parse(userData);

    if (!userData || !userData.permisos) throw new HttpException(new ResponseBody(false, 401, 'No posee permisos suficientes para realizar esta accion' ), 401);

    // Verificar si el usuario tiene al menos uno de los permisos requeridos
    const hasPermission = requiredPermissions.some((permiso: String) => userData.permisos.includes(permiso));

    if (!hasPermission) throw new ForbiddenException('No posee permisos suficientes para realizar esta accion');

    return true;
  }
}
