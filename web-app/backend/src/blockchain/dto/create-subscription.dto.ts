export interface CreateSubscriptionDto {
  address: string;
  currency: string;
  blockHeight: number;
  startingBlockheight: number;
}