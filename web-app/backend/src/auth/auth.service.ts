import {
  hash as hashPassword,
  compare as comparePasswords,
} from 'bcrypt';
import { Injectable, InternalServerErrorException, UnauthorizedException } from '@nestjs/common';
import { UsersService } from '../users/users.service';
import { JwtService } from '@nestjs/jwt';
import { User } from '@prisma/client';

const SALT_ROUNDS = 12;

@Injectable()
export class AuthService {
  constructor(
    private readonly usersService: UsersService,
    private readonly jwtService: JwtService,
  ) {}

  async signUp(email: string, password: string, name?: string): Promise<User> {
    const passwordHash = await hashPassword(password, SALT_ROUNDS);
    return await this.usersService.create(email, passwordHash, name);
  }

  async login(email: string, password: string): Promise<{ accessToken: string }> {
    const user = await this.usersService.findByEmail(email);

    try {
      if (user && await comparePasswords(password, user.password)) {
        const payload = { username: user.email, sub: user.id };
        const { password: _, ...userWithoutPassword } = user;
        return {
          ...userWithoutPassword,
          accessToken: this.jwtService.sign(payload),
        };
      }

      throw new UnauthorizedException('Invalid credentials');
    } catch (error) {
      console.error('Error while validating password:', error);
      throw new InternalServerErrorException(error);
    }
  }
}
