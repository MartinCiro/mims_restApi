import { Test, TestingModule } from '@nestjs/testing';
import AuthService from '@core/auth/authService';
import AuthPortt from '@core/auth/authPort';
import { AuthPort } from '@api/auth/auth-port.token';
import bcrypt from 'bcryptjs';
import config from 'src/config';

describe('AuthService', () => {
  let authService: AuthService;
  let mockAuthPort: jest.Mocked<AuthPortt>;
  let hashedPassword: string;

  beforeEach(async () => {
    hashedPassword = await bcrypt.hash("123456", config.SALT?.toString() ? parseInt(config.SALT) : 10);
    const module: TestingModule = await Test.createTestingModule({
      providers: [
        AuthService,
        {
          provide: AuthPort, 
          useValue: {
            retrieveUser: jest.fn(), 
          },
        },
      ],
    }).compile();

    authService = module.get<AuthService>(AuthService);
    mockAuthPort = module.get<jest.Mocked<AuthPortt>>(AuthPort);
  });

  it('Debe encontrar un usuario por id', async () => {
    const mockUser = { 
      id_user: "s@gmail.com", 
      usuario: "Ciro",
      password: hashedPassword,  
      id_rol: 2, 
      status: "activo", 
      permisos: ["read", "write"] 
    };

    mockAuthPort.retrieveUser.mockResolvedValue(mockUser);

    const auth = await authService.loginUser({ id: "s@gmail.com", password: "123456" });
    expect(auth.ok).toBe(true);
    expect(auth.result.usuario.id).toBe(mockUser.id_user); 
    expect(mockAuthPort.retrieveUser).toHaveBeenCalledWith({ id: "s@gmail.com" });
  });
});
