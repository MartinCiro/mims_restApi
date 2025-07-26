import { PasswordService } from "@core/auth/service/PasswordService";

export class Usuario {
    private encryptedPassword: string;

    constructor(
        public readonly username: string,
        password: string,
        public readonly id_rol?: number,
        public readonly id_estado?: number,
        isEncrypted: boolean = false 
    ) {
        this.encryptedPassword = isEncrypted ? password : PasswordService.encodePassword(password);
    }

    getEncryptedPassword(): string {
        return this.encryptedPassword;
    }

    comparePassword(plainPassword: string): boolean {
        return PasswordService.comparePasswords(plainPassword, this.encryptedPassword);
    }
}
