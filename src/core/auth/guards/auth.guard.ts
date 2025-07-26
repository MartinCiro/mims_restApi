import { Injectable, CanActivate, ExecutionContext, UnauthorizedException, HttpException } from '@nestjs/common';
import { Reflector } from '@nestjs/core';
import { verifyJWT } from '@core/auth/service/jwtService';
import { ResponseBody } from '@api/models/ResponseBody'; 

// Caché en memoria para almacenar información de usuarios autenticados
const userCache = new Map<string, any>();

@Injectable()
export class AuthGuard implements CanActivate {
  constructor(private reflector: Reflector) {}

  async canActivate(context: ExecutionContext): Promise<boolean> {
    const request = context.switchToHttp().getRequest();
    const rawToken = request.headers['jwt'] || request.headers['authorization'];

    if (!rawToken) {
      const error = new UnauthorizedException('No se ha proporcionado token');
      throw new HttpException(
        new ResponseBody(false, 404, error.message),
        404
      );
    }

    // Eliminar prefijo "Bearer " si está presente
    const token = rawToken.startsWith('Bearer ') ? rawToken.slice(7) : rawToken;

    // Validar formato del token con una expresión regular
    const jwtRegex = /^[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+\.[A-Za-z0-9-_]+$/;
    if (!jwtRegex.test(token)) {
      const error = new UnauthorizedException('El token proporcionado no tiene un formato válido');
      throw new HttpException(
        new ResponseBody(false, 404, error.message),
        404
      );
    }

    try {
      const { userInfo, jwt } = await verifyJWT(token);
      
      // Almacenar usuario en caché
      if (!userInfo.userInfo?.doc) {
        const error = new UnauthorizedException('Token inválido');
        throw new HttpException(
          new ResponseBody(false, 404, error.message),
          404
        );
      }
      userCache.set(userInfo.userInfo.id_user.toString(), userInfo);

      // Adjuntar la información del usuario a la solicitud
      request.user = userInfo;
      if (jwt) request.newToken = jwt;

      return true;
    } catch (error: any) {
      console.error(error);
      const err = new UnauthorizedException(error.data || 'Token inválido o expirado');
      throw new HttpException(
        new ResponseBody(false, 404, err.message),
        404
      );
    }
  }

  /* // Método estático para obtener información del usuario desde la caché
  static getUserInfo(id_user: string) {
    return userCache.get(id_user);
  } */
}

