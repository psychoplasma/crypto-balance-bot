import { Injectable, NotFoundException } from '@nestjs/common';
import { PrismaClient, Subscription } from '@prisma/client';

@Injectable()
export class BlockchainService {
  constructor(private readonly prismaClient: PrismaClient) {}

  async getSubscriptionsByUserId(userId: string): Promise<Subscription[]> {
    return await this.prismaClient.subscription.findMany({ where: { userId } });
  }

  async getSubscriptionsByUserIdAndCurreny(userId: string, currency: string): Promise<Subscription[]> {
    return this.prismaClient.subscription.findMany({
      where: { userId, currency },
    });
  }

  async subscribe(
    userId: string,
    currency: string,
    address: string,
    blockHeight: number,
    startingBlockHeight: number,
  ): Promise<Subscription> {
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
      totalReceived: 0,
      totalSpent: 0,
      filters: '',
    }});
  }

  async unsubscribe(userId: string, currency: string, account: string): Promise<void> {
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
  }
}
