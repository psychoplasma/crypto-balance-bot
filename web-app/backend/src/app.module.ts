import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { AuthModule } from './auth/auth.module';
import { UsersModule } from './users/users.module';
import { BlockchainModule } from './blockchain/blockchain.module';

@Module({
  imports: [
    AuthModule,
    BlockchainModule,
    ConfigModule.forRoot(),
    UsersModule,
  ],
})
export class AppModule {}