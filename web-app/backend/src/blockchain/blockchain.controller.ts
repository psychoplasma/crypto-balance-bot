import { Controller, Get, Post, Body, Param, Delete, UseGuards, Request } from '@nestjs/common';
import { BlockchainService } from './blockchain.service';
import { CreateSubscriptionDto } from './dto/create-subscription.dto'
import { DeleteSubscriptionDto } from './dto/delete-subscription.dto';
import { AuthGuard } from '../auth/auth.guard';

@Controller('api/subscriptions')
export class BlockchainController {
  constructor(private readonly blockchainService: BlockchainService) {}

  @UseGuards(AuthGuard)
  @Post()
  async createSubscription(
    @Request() req: { user: { sub: string, username: string} },
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
  @Delete()
  async deleteSubscription(
    @Request() req: { user: { sub: string, username: string} },
    @Body() deleteSubscriptionDto: DeleteSubscriptionDto,
  ) {
    return await this.blockchainService.unsubscribe(
      req.user.sub,
      deleteSubscriptionDto.currency,
      deleteSubscriptionDto.address,
    );
  }
  @UseGuards(AuthGuard)
  @Get()
  async getUserSubscriptions(
    @Request() req: { user: { sub: string, username: string} },
  ): Promise<any[]> {
    return await this.blockchainService.getSubscriptionsByUserId(req.user.sub);
  }

  @UseGuards(AuthGuard)
  @Get(':currency')
  async getUserSubscriptionsByCurrency(
    @Request() req: { user: { sub: string, username: string} },
    @Param('currency') currency: string,
  ): Promise<any[]> {
    return await this.blockchainService
      .getSubscriptionsByUserIdAndCurreny(req.user.sub, currency);
  }
}
