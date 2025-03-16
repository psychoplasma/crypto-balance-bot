import { Injectable, UnauthorizedException } from '@nestjs/common';
import { UsersService } from '../users/users.service';
import { JwtService } from '@nestjs/jwt';
import { User } from '@prisma/client';

@Injectable()
export class AuthService {
  constructor(
    private readonly usersService: UsersService,
    private readonly jwtService: JwtService,
  ) {}

  async signUp(email: string, password: string, name?: string): Promise<User> {
    // FIXME: Do not save plain password, instead use bcrypt and save password hash
    return await this.usersService.create(email, password, name);
  }

  async login(email: string, password: string): Promise<{ accessToken: string }> {
    const user = await this.usersService.findByEmail(email);
    if (user && user.password === password) {
      const payload = { username: user.email, sub: user.id };
      return {
        accessToken: this.jwtService.sign(payload),
      };
    }

    throw new UnauthorizedException('Invalid credentials');
  }
}