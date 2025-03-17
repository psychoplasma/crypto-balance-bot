export interface User {
  id: string;
  email: string;
  name?: string;
  createdAt: Date;
  token?: string;
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