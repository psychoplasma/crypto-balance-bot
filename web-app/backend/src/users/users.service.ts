import { Injectable } from '@nestjs/common';
import { PrismaClient, User } from '@prisma/client';

@Injectable()
export class UsersService {
    constructor(private readonly prismaClient: PrismaClient) {}

    async create(email: string, password: string, name?: string): Promise<User> {
        return this.prismaClient.user.create({
            data: {
                email,
                password,
                name,
            },
        });
    }

    async findAll(): Promise<User[]> {
        return await this.prismaClient.user.findMany();
    }

    async findOne(id: string): Promise<User | null> {
        return await this.prismaClient.user.findUnique({ where: { id } });
    }

    async findByEmail(email: string): Promise<User | null> {
        return await this.prismaClient.user.findUnique({ where: { email } });
    }

    async update(id: string, password?: string, name?: string): Promise<User> {
        return await this.prismaClient.user.update({
            where: { id },
            data: {
                password,
                name,
            },
        });
    }

    async remove(id: string): Promise<void> {
        await this.prismaClient.user.delete({ where: { id } });
    }
}