import 'server-only';
import { Subscription, User } from './types';
import { getSession } from '@/lib/session';

const BACKEND_URL = process.env.NEXT_PUBLIC_BACKEND_URL || 'http://localhost:3000';

const getAuthHeaders = async () => {
  const { isAuth, token } = await getSession();
  return {
    'Content-Type': 'application/json',
    ...(isAuth ? { Authorization: `Bearer ${token}` } : {}),
  };
};

// TODO: Resolve response data for all api methods
export async function login(email: string, password: string): Promise<User> {
  const res = await fetch('/api/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ email, password }),
  });

  const ret = await res.json();
  return ret as User;
};

export const authenticatedRequest = async (url: string, options: RequestInit = {}) => {
  const authHeaders = await getAuthHeaders();
  return fetch(url, {
    ...options,
    headers: {
      ...options.headers,
      ...authHeaders,
    },
  });
};

export async function signup(email: string, password: string, name?: string): Promise<Response> {
  return await fetch(`${BACKEND_URL}/api/auth/signup`, {
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
  const res = await authenticatedRequest(`${BACKEND_URL}/api/subscriptions/${userId}`, {
    method: 'POST',
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
  const res = await authenticatedRequest(`${BACKEND_URL}/api/subscriptions/${userId}`, {
    method: 'GET',
  });
  const text = await res.text();
  return JSON.parse(text) as Subscription[];
}

export async function getSubscriptionsByCurrency(userId: string, currency: string): Promise<Subscription[]> {
  const res = await authenticatedRequest(`${BACKEND_URL}/api/subscriptions/${userId}/${currency}`, {
    method: 'GET',
  });
  const text = await res.text();
  return JSON.parse(text) as Subscription[];
}

export async function deleteSubscription(userId: string, currency: string, address: string): Promise<Response> {
  return await authenticatedRequest(`${BACKEND_URL}/api/subscriptions/${userId}`, {
    method: 'DELETE',
    body: JSON.stringify({ currency, address }),
  });
}

