import { Test, TestingModule } from '@nestjs/testing';
import { PermisoService } from '@core/permisos/permisoService';
import PermisosPort  from '@core/permisos/permisoPort';
import { PermisosPortToken } from '@api/permisos/permiso-port.token';

describe('PermisoService', () => {
  let permisoService: PermisoService;
  let mockPermisoPort: jest.Mocked<PermisosPort>;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [
        PermisoService,
        {
            provide: PermisosPortToken,
          useValue: {
            crearPermisos: jest.fn(),
            obtenerPermisos: jest.fn(),
            obtenerPermisosXid: jest.fn(),
            delPermiso: jest.fn(),
            actualizaPermiso: jest.fn(),
          },
        },
      ],
    }).compile();

    permisoService = module.get<PermisoService>(PermisoService);
    mockPermisoPort = module.get<jest.Mocked<PermisosPort>>(PermisosPortToken);
  });

  it('Debe encontrar un permiso por id', async () => {
    const mockUser = { id: 1, nombres: "Ciro", email: "s@gmail.com" };

    // El método espera un objeto con `{ id }`
    mockPermisoPort.obtenerPermisosXid.mockResolvedValue(mockUser);

    // Llamar con el objeto correcto
    const permiso = await permisoService.obtenerPermisoXid({ id: "s@gmail.com" });

    expect(permiso).toEqual(mockUser);
    expect(mockPermisoPort.obtenerPermisosXid).toHaveBeenCalledWith({ id: "s@gmail.com" });
  });
});
