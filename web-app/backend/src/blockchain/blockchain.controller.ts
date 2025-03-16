import { Controller, Get, Post, Body, Param, Delete, UseGuards, Request } from '@nestjs/common';
import { BlockchainService } from './blockchain.service';
import { CreateSubscriptionDto } from './dto/create-subscription.dto'
import { DeleteSubscriptionDto } from './dto/delete-subscription.dto';
import { AuthGuard } from '../auth/auth.guard';

@Controller('api/subscriptions')
export class BlockchainController {
  constructor(private readonly blockchainService: BlockchainService) {}

  @UseGuards(AuthGuard)
  @Post(':userId')
  async createSubscription(
    @Request() req: { user: { sub: string, username: string} },
    @Param('userId') userId: string,
    @Body() createSubscriptionDto: CreateSubscriptionDto,
  ) {
    return await this.blockchainService.subscribe(
      req.user.sub,
      createSubscriptionDto.currency,
      createSubscriptionDto.address,
      createSubscriptionDto.blockHeight,
      createSubscriptionDto.startingBlockheight,
    );
  }

  @UseGuards(AuthGuard)
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
  @UseGuards(AuthGuard)
  @Get(':userId')
  async getUserSubscriptions(@Param('userId') userId: string): Promise<any[]> {
    return await this.blockchainService.getSubscriptionsByUserId(userId);
  }

  @UseGuards(AuthGuard)
  @Get(':userId/:currency')
  async getUserSubscriptionsByCurrency(
    @Param('userId') userId: string,
    @Param('currency') currency: string,
  ): Promise<any[]> {
    return await this.blockchainService
      .getSubscriptionsByUserIdAndCurreny(userId, currency);
  }
}
