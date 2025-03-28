import { Injectable, InternalServerErrorException, NotFoundException } from '@nestjs/common';
import { PrismaClient, Subscription } from '@prisma/client';

@Injectable()
export class BlockchainService {
  constructor(private readonly prismaClient: PrismaClient) {}

  async getSubscriptionsByUserId(userId: string): Promise<Subscription[]> {
    try {
      return await this.prismaClient.subscription.findMany({ where: { userId } });
    } catch (error) {
      console.error('Error while fetching subscriptions by user id:', error);
      throw new InternalServerErrorException();
    }
  }

  async getSubscriptionsByUserIdAndCurreny(userId: string, currency: string): Promise<Subscription[]> {
    try {
      return this.prismaClient.subscription.findMany({
        where: { userId, currency },
      });
    } catch (error) {
      console.error('Error while fetching subscriptions by user id and currency:', error);
      throw new InternalServerErrorException();
    }
  }

  async subscribe(
    userId: string,
    currency: string,
    address: string,
    blockHeight: number,
    startingBlockHeight: number,
  ): Promise<Subscription> {
    try {
      const user = await this.prismaClient.user.findUnique({ where: { id: userId } });

      if (!user) {
        throw new NotFoundException(`User with ID ${userId} not found`);
      }

      return await this.prismaClient.subscription.create({ data: {
        userId,
        currency,
        account: address,
        blockHeight,
        startingBlockHeight,
        filters: [],
      }});
    } catch (error) {
      console.error('Error while creating subscription:', error);
      throw new InternalServerErrorException();
    }
  }

  async unsubscribe(userId: string, currency: string, account: string): Promise<void> {
    try {
      const user = await this.prismaClient.user.findUnique({ where: { id: userId } });

      if (!user) {
        throw new NotFoundException(`User with ID ${userId} not found`);
      }

      const subscription = await this.prismaClient.subscription.findUnique({
        where: { userId_currency_account: { userId, account, currency } },
      });

      if (!subscription) {
        throw new NotFoundException(`Subscription not found`);
      }

      await this.prismaClient.subscription.delete({
        where: { id: subscription.id, userId, account, currency },
      });
    } catch (error) {
      console.error('Error while deleting subscription:', error);
      throw new InternalServerErrorException();
    }
  }
}
