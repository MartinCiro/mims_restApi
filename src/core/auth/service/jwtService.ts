import jwt from 'jsonwebtoken';
import config from 'src/config';

interface JwtPayload {
    userInfo: any;
    id_user: number;
    username: string;
    id_rol: number;
    exp?: number;
}

export const generateJWT = (userInfo: any): string => {
    if (!config.JWT_SECRETO) throw new Error("JWT_SECRETO no está definido en la configuración.");
    return jwt.sign({ userInfo }, config.JWT_SECRETO, { expiresIn: 3600 });
};

export const verifyJWT = async (token: string): Promise<{ userInfo: JwtPayload; jwt?: string }> => {
    let response: { userInfo: JwtPayload; jwt?: string } = { userInfo: {} as JwtPayload };

    // Decodificar sin verificar la firma
    const decoded = jwt.decode(token) as JwtPayload | null;
    if (!decoded || !decoded.userInfo?.id_user) throw { ok: false, status_cod: 401, data: "El JWT es inválido" };
    // Si está en entorno de desarrollo, retornar sin verificar
    if (config.env === 'Dev') return { userInfo: decoded };

    try {
        // Verificar firma
        if (!config.JWT_SECRETO) throw new Error("JWT_SECRETO no está definido en la configuración.");

        // Verificar el token
        const verified = jwt.verify(token, config.JWT_SECRETO) as JwtPayload;
        // Calcular tiempo restante hasta la expiración
        if (verified.exp) {
            const expireDate = new Date(verified.exp * 1000);
            const now = new Date();
            const diffMins = Math.round((expireDate.getTime() - now.getTime()) / 60000);

            // Regenerar token si le quedan menos de 10 minutos
            if (diffMins < 10) response.jwt = generateJWT(verified);
        }
        response.userInfo = verified;
        return response;

    } catch (error: any) {
        throw error.name === 'TokenExpiredError' ? { ok: false, status_cod: 401, data: 'JWT expirado. Por favor inicie sesión nuevamente' } : { ok: false, status_cod: 401, data: "El JWT es inválido" };
    }
};