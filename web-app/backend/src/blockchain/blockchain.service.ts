import { Injectable } from '@nestjs/common';
import { PrismaClient, Subscription } from '@prisma/client';

@Injectable()
export class BlockchainService {
    constructor(private readonly prismaClient: PrismaClient) {}

    async getSubscriptionsByUserId(userId: string): Promise<Subscription[]> {
        return this.prismaClient.subscription.findMany({ where: { userId } });
    }

    async getSubscriptionsByUserIdAndCurreny(userId: string, currency: string): Promise<Subscription[]> {
        return this.prismaClient.subscription.findMany({ where: { userId, currency } });
    }

    async subscribe(
        userId: string,
        currency: string,
        address: string,
        blockHeight: number,
        startingBlockHeight: number,
    ): Promise<Subscription> {
        // TODO: use logged-in user instead
        await this.prismaClient.user.findUniqueOrThrow({ where: { id: userId } });

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
        // TODO: use logged-in user instead
        await this.prismaClient.user.findUniqueOrThrow({ where: { id: userId } });

        const subscription = await this.prismaClient.subscription.findUniqueOrThrow(
            { where: {
                userId_currency_account: { userId, account, currency }}
             },
        );

        await this.prismaClient.subscription.delete({
            where: { id: subscription.id, userId, account, currency },
        });
    }
}