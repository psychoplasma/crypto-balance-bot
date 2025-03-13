import { Subscription } from './types';

const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:3000';

export async function login(email: string, password: string, name?: string): Promise<Response> {
  return await fetch(`${BACKEND_URL}/auth/login`, {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json',
    },
    body: JSON.stringify({ email, password, name }),
  });
}

export async function signup(email: string, password: string, name?: string): Promise<Response> {
  return await fetch(`${BACKEND_URL}/auth/signup`, {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json',
    },
    body: JSON.stringify({ email, password, name }),
  });
}

export async function createSubscription(
  userId: string,
  currency: string,
  address: string,
  blockHeight: number,
  startingBlockheight: number,
): Promise<Subscription> {
  const res = await fetch(`${BACKEND_URL}/subscriptions/${userId}`, {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      currency,
      address,
      blockHeight,
      startingBlockheight,
    }),
  });

  const text = await res.text();
  return JSON.parse(text) as Subscription;
}

export async function getSubscriptions(userId: string): Promise<Subscription[]> {
  const res = await fetch(`${BACKEND_URL}/subscriptions/${userId}`, {
    method: 'GET',
    headers: {
        'Content-Type': 'application/json',
    },
  });
  const text = await res.text();
  return JSON.parse(text) as Subscription[];
}

export async function getSubscriptionsByCurrency(userId: string, currency: string): Promise<Subscription[]> {
  const res = await fetch(`${BACKEND_URL}/subscriptions/${userId}/${currency}`, {
    method: 'GET',
    headers: {
        'Content-Type': 'application/json',
    },
  });
  const text = await res.text();
  return JSON.parse(text) as Subscription[];
}

export async function deleteSubscription(userId: string, currency: string, address: string): Promise<Response> {
  return await fetch(`${BACKEND_URL}/subscriptions/${userId}`, {
    method: 'DELETE',
    headers: {
        'Content-Type': 'application/json',
    },
    body: JSON.stringify({ currency, address }),
  });
}

