export interface User {
  email: string,
  name?: string,
  createdAt: Date,
}

export interface Subscription {
  id: string;
  userId: string;
  blockHeight: number;
  startingBlockHeight: number;
  currency: string;
  account: string;
  filters: string;
}