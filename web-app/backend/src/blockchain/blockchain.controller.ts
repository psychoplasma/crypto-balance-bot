import { Controller, Get, Post, Body, Param, Delete } from '@nestjs/common';
import { BlockchainService } from './blockchain.service';
import { CreateSubscriptionDto } from './dto/create-subscription.dto'
import { DeleteSubscriptionDto } from './dto/delete-subscription.dto';

@Controller('subscriptions')
export class BlockchainController {
  constructor(private readonly blockchainService: BlockchainService) {}

  @Post(':userId')
  async createSubscription(
    @Param('userId') userId: string,
    @Body() createSubscriptionDto: CreateSubscriptionDto,
  ) {
    return await this.blockchainService.subscribe(
      userId,
      createSubscriptionDto.currency,
      createSubscriptionDto.address,
      createSubscriptionDto.blockHeight,
      createSubscriptionDto.startingBlockheight,
    );
  }

  @Delete(':userId')
  async deleteSubscription(
    @Param('userId') userId: string,
    @Body() deleteSubscriptionDto: DeleteSubscriptionDto,
  ) {
    return await this.blockchainService.unsubscribe(
      userId,
      deleteSubscriptionDto.currency,
      deleteSubscriptionDto.address,
    );
  }

  @Get(':userId')
  async getUserSubscriptions(@Param('userId') userId: string): Promise<any[]> {
    return await this.blockchainService.getSubscriptionsByUserId(userId);
  }

  @Get(':userId/:currency')
  async getUserSubscriptionsByCurrency(
    @Param('userId') userId: string,
    @Param('currency') currency: string,
  ): Promise<any[]> {
    return await this.blockchainService
      .getSubscriptionsByUserIdAndCurreny(userId, currency);
  }
}
