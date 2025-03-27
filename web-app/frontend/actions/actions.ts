'use server';

import { Subscription, User } from "@/lib/types";
import {
  login,
  signup,
  createSubscription,
  getSubscriptions,
  getSubscriptionsByCurrency,
  deleteSubscription,
} from "@/lib/api";


export async function loginAction(email: string, password: string): Promise<User> {
  return await login(email, password);
}

export async function signupAction(email: string, password: string, name?: string): Promise<void> {
  await signup(email, password, name);
}

export async function createSubscriptionAction(
  currency: string,
  address: string,
  blockHeight?: number,
  startingBlockheight?: number,
): Promise<Subscription> {
  return await createSubscription(
    currency,
    address,
    blockHeight,
    startingBlockheight,
  );
}

export async function getSubscriptionsAction(): Promise<Subscription[]> {
  return await getSubscriptions();
}

export async function getSubscriptionsByCurrencyAction(
  currency: string,
): Promise<Subscription[]> {
  return await getSubscriptionsByCurrency(currency);
}

export async function deleteSubscriptionAction(
  currency: string,
  address: string,
): Promise<void> {
  return await deleteSubscription(currency, address);
}
