import { Injectable, InternalServerErrorException } from '@nestjs/common';
import { PrismaClient, User } from '@prisma/client';

@Injectable()
export class UsersService {
  constructor(private readonly prismaClient: PrismaClient) {}

  async create(email: string, password: string, name?: string): Promise<User> {
    try {
      return this.prismaClient.user.create({
        data: {
          email,
          password,
          name,
        },
      });
    } catch (error) {
      console.error('Error while creating user:', error);
      throw new InternalServerErrorException();
    }
  }

  async findAll(): Promise<User[]> {
    try {
      return await this.prismaClient.user.findMany();
    } catch (error) {
      console.error('Error while fetching users:', error);
      throw new InternalServerErrorException();
    }
  }

  async findOne(id: string): Promise<User | null> {
    try {
      return await this.prismaClient.user.findUnique({ where: { id } });
    } catch (error) {
      console.error('Error while fetching user by id:', error);
      throw new InternalServerErrorException();
    }
  }

  async findByEmail(email: string): Promise<User | null> {
    try {
      return await this.prismaClient.user.findUnique({ where: { email } });
    } catch (error) {
      console.error('Error while fetching user by email:', error);
      throw new InternalServerErrorException();
    }
  }

  async update(id: string, password?: string, name?: string): Promise<User> {
    try {
      return await this.prismaClient.user.update({
        where: { id },
        data: {
          password,
          name,
        },
      });
    } catch (error) {
      console.error('Error while updating user:', error);
      throw new InternalServerErrorException();
    }
  }

  async remove(id: string): Promise<void> {
    try {
      await this.prismaClient.user.delete({ where: { id } });
    } catch (error) {
      console.error('Error while deleting user:', error);
      throw new InternalServerErrorException();
    }
  }
}
