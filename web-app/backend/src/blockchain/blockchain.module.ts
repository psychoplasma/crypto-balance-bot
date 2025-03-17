import { Module } from '@nestjs/common';
import { PrismaClient } from '@prisma/client';
import { BlockchainService } from './blockchain.service';
import { BlockchainController } from './blockchain.controller';
import { AuthModule } from '../auth/auth.module';
import { ConfigModule } from '@nestjs/config';

@Module({
  imports: [AuthModule, ConfigModule],
  controllers: [BlockchainController],
  providers: [BlockchainService, PrismaClient],
})
export class BlockchainModule {};
