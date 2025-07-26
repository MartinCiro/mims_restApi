import { IsString, IsOptional, IsArray, ArrayMinSize } from 'class-validator';

export class ActualizarPermisoDto {
  @IsOptional()
  @IsString({ message: 'El texto de descripción no es valido' })
  readonly descripcion?: string;

  @IsArray({ message: 'Debe ser una lista de permisos' })
  @ArrayMinSize(1, { message: 'Debe proporcionar al menos un permiso' })
  @IsString({ each: true, message: 'Cada permiso debe ser un texto válido' })
  readonly permisos!: string[];
}
