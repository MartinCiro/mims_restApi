import { Usuario } from './entities/Usuario';
import { generateJWT } from './service/jwtService';
import { Injectable, Inject } from '@nestjs/common';
import { ResponseBody } from '@api/models/ResponseBody';
import { RedisService } from '@shared/cache/redis.service';
import AuthPort from './authPort';

@Injectable()
export default class AuthService {
    constructor(
        @Inject('AuthPort') private authPort: AuthPort,
        private readonly redisService: RedisService
    ) {}

    async loginUser({ username, password }: { username: string; password: string }): Promise<ResponseBody<any>> {
        try {
            const usuarioRetrieved = await this.authPort.retrieveUser({ username });

            const isPasswordValid = new Usuario(
                usuarioRetrieved.id_user,
                usuarioRetrieved.password,
                usuarioRetrieved.id_rol,
                usuarioRetrieved.estado,
                true
            ).comparePassword(password);

            if (!isPasswordValid) throw new Error('Usuario o contraseña inválida');

            const userCacheKey = `user:${usuarioRetrieved.id_user}`;
            const eventKey = `event:usuario.logeado:${usuarioRetrieved.id_user}`;

            // Revisar si el usuario ya está en cache
            const cachedUser = await this.redisService.get(userCacheKey);
            let userData;

            if (cachedUser) {
                userData = JSON.parse(cachedUser);
            } else {
                userData = {
                    id_user: usuarioRetrieved.id_user,
                    nombre: usuarioRetrieved.usuario,
                    id_rol: usuarioRetrieved.id_rol
                };
            
                const userDataWithPermissions = {
                    ...userData,
                    permisos: usuarioRetrieved.permisos
                };
            
                await this.redisService.set(userCacheKey, JSON.stringify(userDataWithPermissions));
            }

            // Generar siempre un nuevo JWT para el usuario
            const token = generateJWT(userData);

            // Verificar si el evento ya se publicó en Redis
            const eventExists = await this.redisService.get(eventKey);

            if (!eventExists) await this.redisService.set(JSON.stringify(eventKey), 'true');
            
            return {
                ok: true,
                statusCode: 200,
                result: {
                    token,
                    usuario: {
                        id: userData.id_user,
                        nombre: userData.nombre,
                        rol: userData.id_rol
                    }
                }
            };
        } catch (error: any) {
            return {
                ok: false,
                statusCode: 401,
                result: error.data || 'Usuario o contraseña inválida'
            };
        }
    }
}
