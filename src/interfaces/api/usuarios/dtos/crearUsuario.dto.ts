import { IsNotEmpty, IsNumber, IsString, MinLength, IsOptional, IsDate } from 'class-validator';

export class CrearUsuarioDto {
  @IsNotEmpty({ message: 'El nombre es obligatorio' })
  @IsString({ message: 'El texto de nombre no es valido' })
  readonly nombres!: string

  @IsNotEmpty({ message: 'La contraseña es obligatoria' })
  @MinLength(6, { message: 'La contraseña debe tener al menos 6 caracteres' })
  readonly passwd!: string

  @IsOptional()
  @IsNumber({ allowNaN: false, allowInfinity: false }, { message: 'Debe ser un número' })
  readonly id_rol?: number;

  @IsNotEmpty({ message: 'El apellido es obligatorio' })
  @IsString({ message: 'El texto de apellido no es valido' })
  readonly apellido!: string;

  @IsNotEmpty({ message: 'El numero de documento es obligatorio' })
  @IsString({ message: 'El texto en numero de documento no es valido' })
  readonly numero_documento!: string

  @IsNotEmpty({ message: 'El correo es obligatorio' })
  @IsString({ message: 'El texto de correo no es valido' })
  readonly email!: string;

  @IsOptional()
  @IsNumber({ allowNaN: false, allowInfinity: false }, { message: 'Debe ser un número' })
  readonly estado_id?: number; 
  
  @IsOptional()
  @IsString({ message: 'El texto ingresado en la informacion del perfil no es valido' })
  readonly info_perfil?: string;

  @IsNotEmpty({ message: 'El nombre de usuario es obligatorio' })
  @IsString({ message: 'El texto en nombre de usuario no es valido' })
  readonly nom_user!: string;
  
  @IsOptional()
  @IsString({ message: 'El texto ingresado en el numero de contacto no es valido' })
  readonly numero_contacto?: string; 

  @IsNotEmpty({ message: 'La fecha es obligatoria' })
  @IsDate({ message: 'La fecha no es valida' })
  readonly fecha_nacimiento!: string;
}

