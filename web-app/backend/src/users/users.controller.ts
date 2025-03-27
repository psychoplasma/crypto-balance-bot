import { Controller, Get } from '@nestjs/common';
import { UsersService } from './users.service';
import { mapToUserResponse, UserResponse } from './dto/user-response.dto';

@Controller('api/users')
export class UsersController {
  constructor(private readonly usersService: UsersService) {}

  // FIXME: Add guard for user role
  @Get()
  async findAll(): Promise<UserResponse[]> {
    const users = await this.usersService.findAll();
    return users.map(mapToUserResponse);
  }
}
