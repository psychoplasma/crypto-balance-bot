import { Module } from '@nestjs/common';
import { PrismaClient } from '@prisma/client';
import { BlockchainService } from './blockchain.service';
import { BlockchainController } from './blockchain.controller';

@Module({
  imports: [],
  controllers: [BlockchainController],
  providers: [BlockchainService, PrismaClient],
})
export class BlockchainModule {}