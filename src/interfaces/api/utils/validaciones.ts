
import { HttpException } from '@nestjs/common/exceptions/http.exception';
import { HttpStatus } from '@nestjs/common/enums/http-status.enum';
import { ResponseBody } from '@api/models/ResponseBody';

export const validarBlank = (valor: any, nombre: string): void => {
  if (!valor) {
    throw new HttpException(`No se ha proporcionado ${nombre}`, HttpStatus.BAD_REQUEST);
  }
};

//metodos para validar si existe y otros
export const validarExistente = (
  code: string,
  target?: string[] | string
): { ok: boolean; data?: string } => {
  if (code === 'P2002') {
    const campo = Array.isArray(target) ? target[0] : target;
    return {
      ok: false,
      data: `${campo} ya existe en la base de datos`
    };
  }

  return { ok: true };
};

export const validarNoExistente = (valor: string, nombre: string | number) => {
  if (valor === 'P2025') return {
      ok: false,
      status_cod: 400,
      data: `${nombre} no existe en la base de datos`
    };
  return { ok: true };
};

export function handleException(error: any): never {
  if (typeof error === 'object' && error !== null && 'status_cod' in error && 'data' in error) {
    const statusCode = typeof error.status_cod === 'number' ? error.status_cod : HttpStatus.INTERNAL_SERVER_ERROR;
    const data = typeof error.data === 'string' ? error.data : 'Error desconocido';
    throw new HttpException(new ResponseBody(false, statusCode, data), statusCode);
  }
  console.error(error);
  throw new HttpException(
    new ResponseBody(false, HttpStatus.INTERNAL_SERVER_ERROR, 'Error interno del servidor'),
    HttpStatus.INTERNAL_SERVER_ERROR
  );
}

export function capitalize(text: string) {
  return text.charAt(0).toUpperCase() + text.slice(1).toLowerCase();
}
