import { Controller, Post, Body, HttpException, HttpStatus } from '@nestjs/common';
import AuthService  from '@core/auth/authService';
import { AuthDto } from './dtos/auth.dto';
import { ResponseBody } from '@api/models/ResponseBody';

@Controller('auth')
export class AuthController {
  constructor(private readonly authService: AuthService) {}
  
  @Post('login')
  async login(@Body() body: AuthDto): Promise<ResponseBody<any>> {
    const { username, enpass: password } = body;

    if (!username || !password) throw new HttpException(new ResponseBody(false, HttpStatus.BAD_REQUEST, "Usuario y contraseña son obligatorios"), HttpStatus.BAD_REQUEST);

    try {
      const auth = await this.authService.loginUser({ username, password });
      return new ResponseBody(auth.ok, auth.statusCode, auth.result);
    } catch (error: any) {
      if (typeof error === 'object' && error !== null && 'status_cod' in error && 'data' in error) {
        const err = error as { status_cod: unknown; data: unknown };
        const statusCode = typeof err.status_cod === 'number' ? err.status_cod : HttpStatus.INTERNAL_SERVER_ERROR;
        const data = typeof err.data === 'string' ? err.data : 'Error desconocido';
        throw new HttpException(new ResponseBody(false, statusCode, data), statusCode);
      }
      throw new HttpException(new ResponseBody(false, HttpStatus.INTERNAL_SERVER_ERROR, 'Error interno del servidor'), HttpStatus.INTERNAL_SERVER_ERROR);
    }
  }
}
