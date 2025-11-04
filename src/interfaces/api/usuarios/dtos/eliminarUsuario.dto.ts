import { IsNotEmpty, IsString} from 'class-validator';

export class EliminarUsuarioDto {
  @IsNotEmpty({ message: 'El identificador de usuario es obligatorio' })
  @IsString({ message: 'El texto en identificador de usuario no es valido' })
  numero_documento!: string;
}
