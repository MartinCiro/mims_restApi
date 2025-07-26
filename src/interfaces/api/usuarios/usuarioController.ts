import {
  Controller, Post, Body, HttpException, HttpStatus, HttpCode,
  UsePipes, ValidationPipe, Get, Put, Delete, UseGuards, Req
} from '@nestjs/common';
import { AuthGuard } from '@core/auth/guards/auth.guard';
import { Permissions, Public } from '@core/auth/decorators/permissions.decorator';
import { ResponseBody } from '@api/models/ResponseBody';
import { UsuarioService } from '@core/usuarios/usuarioService';
import { handleException } from '@utils/validaciones';
import { CrearUsuarioDto } from './dtos/crearUsuario.dto';
import { PermissionsGuard } from '@core/auth/guards/permissions.guard';
import { ObtenerUsuariosDto } from './dtos/obtenerUsuario.dto';
import { EliminarUsuarioDto } from './dtos/eliminarUsuario.dto';
import { ActualizarUsuarioDto } from './dtos/actualizarUsuario.dto';

@Controller('usuarios')
@UseGuards(AuthGuard) // Todas las rutas requieren autenticación
export class UsuarioController {
  constructor(private readonly usuarioService: UsuarioService) { }

  @Post(['', 'register'])
  @HttpCode(HttpStatus.CREATED)
  @Public()
  @UseGuards(PermissionsGuard)
  async registerUsuario(@Body() body: CrearUsuarioDto): Promise<ResponseBody<string>> {
    try {
      await this.usuarioService.crearUsuario(body);
      return new ResponseBody<string>(true, 201, "Se ha creado el usuario exitosamente");
    } catch (error) {
      handleException(error);
    }
  }

  @Post()
  @HttpCode(HttpStatus.CREATED)
  @UseGuards(PermissionsGuard)
  @Permissions('usuarios:Crea')
  async crearUsuario(@Body() body: CrearUsuarioDto): Promise<ResponseBody<string>> {
    try {
      await this.usuarioService.crearUsuario(body);
      return new ResponseBody<string>(true, 201, "Se ha creado el usuario exitosamente");
    } catch (error) {
      handleException(error);
    }
  }
 
  @Get()
  @HttpCode(HttpStatus.OK)
  @UseGuards(PermissionsGuard)
  @Permissions('usuarios:Lee')
  @UsePipes(new ValidationPipe({
    whitelist: true, transform: true, exceptionFactory: (errors) => {
      const mensajes = errors.map(err => ({
        campo: err.property,
        mensaje: err.constraints ? Object.values(err.constraints).join(', ') : ''
      }));
      return new HttpException(new ResponseBody(false, HttpStatus.BAD_REQUEST, mensajes), HttpStatus.BAD_REQUEST);
    }
  }))
  async obtenerUsuarios(@Body() body: ObtenerUsuariosDto): Promise<ResponseBody<any>> {
    try {
      const usuarios = body.numero_documento
        ? await this.usuarioService.obtenerUsuarioXid({ numero_documento: body.numero_documento })
        : await this.usuarioService.obtenerUsuarios();

      return new ResponseBody<any>(true, 200, usuarios);
    } catch (error) {
      handleException(error);
    }
  }

  @Put()
  @HttpCode(HttpStatus.OK)
  @UseGuards(PermissionsGuard)
  @Permissions('usuarios:Actualiza')
  @UsePipes(new ValidationPipe({
    whitelist: true, transform: true, exceptionFactory: (errors) => {
      const mensajes = errors.map(err => ({
        campo: err.property,
        mensaje: err.constraints ? Object.values(err.constraints).join(', ') : ''
      }));
      return new HttpException(new ResponseBody(false, HttpStatus.BAD_REQUEST, mensajes), HttpStatus.BAD_REQUEST);
    }
  }))
  async actualizarUsuario(@Body() body: ActualizarUsuarioDto): Promise<ResponseBody<string>> {
    if (Object.keys(body).length === 0) {
      throw new HttpException(
        new ResponseBody(false, HttpStatus.BAD_REQUEST, "Debe proporcionar al menos un campo para actualizar."),
        HttpStatus.BAD_REQUEST,
      );
    }

    try {
      await this.usuarioService.upUsuario(body);
      return new ResponseBody(true, HttpStatus.OK, "Usuario actualizado exitosamente.");
    } catch (error) {
      handleException(error);
    }
  }

  @Delete()
  @UseGuards(PermissionsGuard)
  @Permissions('usuarios:Elimina')
  @UsePipes(new ValidationPipe({
    whitelist: true, transform: true, exceptionFactory: (errors) => {
      const mensajes = errors.map(err => ({
        campo: err.property,
        mensaje: err.constraints ? Object.values(err.constraints).join(', ') : ''
      }));
      return new HttpException(new ResponseBody(false, HttpStatus.BAD_REQUEST, mensajes), HttpStatus.BAD_REQUEST);
    }
  }))

  async delUsuario(@Body() eliminarUsuarioDto: EliminarUsuarioDto): Promise<ResponseBody<string>> {
    try {
      await this.usuarioService.delUsuario({ numero_documento: eliminarUsuarioDto.numero_documento});
      return new ResponseBody(true, 201, "Se ha eliminado el usuario exitosamente");
    } catch (error) {
      handleException(error);
    }
  }

  @UseGuards(PermissionsGuard)
  @Permissions('usuarios:LeeProfile')
  @UsePipes(new ValidationPipe({
    whitelist: true, transform: true, exceptionFactory: (errors) => {
      const mensajes = errors.map(err => ({
        campo: err.property,
        mensaje: err.constraints ? Object.values(err.constraints).join(', ') : ''
      }));
      return new HttpException(new ResponseBody(false, HttpStatus.BAD_REQUEST, mensajes), HttpStatus.BAD_REQUEST);
    }
  }))
  @Get('perfil')
  async getPerfil(@Req() request: any) {
    return new ResponseBody(true, 200, { mensaje: 'Usuario autenticado', usuario: request.user });
  }
}
