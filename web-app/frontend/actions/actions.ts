'use server';

import { deleteSession, createSession, getSession } from "@/lib/session";
import { Subscription } from "@/lib/types";
import {
  login,
  signup,
  createSubscription,
  getSubscriptions,
  getSubscriptionsByCurrency,
  deleteSubscription,
} from "../lib/api";

interface AuthResult {
  isAuth: boolean;
  userId?: string;
  token?: string;
}

export async function isAuthenticated(): Promise<AuthResult> {
  const session = await getSession();

  if (!session!.token || !session?.userId) {
    return { isAuth: false };
  }

  return {
    isAuth: true,
    userId: session.userId,
    token: session.token,
  };
}

export async function loginAction(email: string, password: string): Promise<{ userId?: string, error?: string}> {
  try {
    const res = await login(email, password);

    if (!res.accessToken) {
      return { error: 'empty jwt token' };
    }

    await createSession(res.id, res.accessToken);

    return { userId: res.id };
  } catch (e) {
    return { error: e as string };
  }
}

export async function logoutAction() {
  await deleteSession();
}

export async function signupAction(email: string, password: string, name?: string): Promise<void> {
  await signup(email, password, name);
}

export async function createSubscriptionAction(
  userId: string,
  currency: string,
  address: string,
  blockHeight?: number,
  startingBlockheight?: number,
): Promise<Subscription> {
  return await createSubscription(
    userId,
    currency,
    address,
    blockHeight,
    startingBlockheight,
  );
}

export async function getSubscriptionsAction(userId: string): Promise<Subscription[]> {
  return await getSubscriptions(userId);
}

export async function getSubscriptionsByCurrencyAction(
  userId: string,
  currency: string,
): Promise<Subscription[]> {
  return await getSubscriptionsByCurrency(userId, currency);
}

export async function deleteSubscriptionAction(userId: string, currency: string, address: string): Promise<void> {
  return await deleteSubscription(userId, currency, address);
}
