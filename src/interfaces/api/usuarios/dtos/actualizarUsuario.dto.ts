import { IsString, IsOptional, IsInt, IsNumber, IsNotEmpty } from 'class-validator';

export class ActualizarUsuarioDto {
  @IsOptional()
  @IsString({ message: 'El texto de apellido no es valido' })
  readonly apellidos!: string;

  @IsOptional()
  @IsNumber({ allowNaN: false, allowInfinity: false }, { message: 'Debe ser un número' })
  readonly id_estado!: number;

  @IsNotEmpty({ message: 'El nombre es obligatorio' })
  @IsString({ message: 'El texto de nombre no es valido' })
  readonly nombres!: string

  @IsNotEmpty({ message: 'El nombre de usuario es obligatorio' })
  @IsString({ message: 'El texto en nombre de usuario no es valido' })
  readonly username!: string

  @IsNumber({ allowNaN: false, allowInfinity: false }, { message: 'Debe ser un número' })
  readonly id_rol!: number

  @IsNotEmpty({ message: 'El identificador de usuario es obligatorio' })
  @IsString({ message: 'El texto en identificador de usuario no es valido' })
  readonly numero_documento!: string
}
