import { Controller, Post, Body } from '@nestjs/common';
import { AuthService } from './auth.service';
import { CreateUserDto } from '../users/dto/create-user.dto';
import { LoginUserDto } from './dto/login-user.dto';
import { mapToUserResponse, UserResponse } from '../users/dto/user-response.dto';

@Controller('api/auth')
export class AuthController {
  constructor(private readonly authService: AuthService) {}

  @Post('signup')
  async signUp(@Body() createUserDto: CreateUserDto): Promise<UserResponse> {
    const user = await this.authService.signUp(
      createUserDto.email,
      createUserDto.password,
      createUserDto.name,
    );
    return mapToUserResponse(user);
  }

  @Post('login')
  async login(@Body() loginUserDto: LoginUserDto): Promise<{ accessToken: string }> {
    return this.authService.login(loginUserDto.email, loginUserDto.password);
  }
}