import 'server-only';
import { Subscription, User } from './types';
import { validateSession } from '@/lib/session';
import { redirect } from 'next/navigation';

const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:3000';

export async function login(email: string, password: string): Promise<User> {
  const res = await fetch(`${BACKEND_URL}/api/auth/login`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ email, password }),
  });

  if (!res.ok) {
    throw new Error(await res.text());
  }

  return (await res.json()) as User;
};

export async function signup(email: string, password: string, name?: string): Promise<void> {
  const res = await fetch(`${BACKEND_URL}/api/auth/signup`, {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json',
    },
    body: JSON.stringify({ email, password, name }),
  });

  if (!res.ok) {
    throw new Error(await res.text());
  }
}

export async function createSubscription(
  currency: string,
  address: string,
  blockHeight?: number,
  startingBlockheight?: number,
): Promise<Subscription> {
  const res = await authenticatedRequest(`${BACKEND_URL}/api/subscriptions`, {
    method: 'POST',
    body: JSON.stringify({
      currency,
      address,
      blockHeight,
      startingBlockheight,
    }),
  });

  if (!res.ok) {
    throw new Error(await res.text());
  }

  return (await res.json()) as Subscription;
}

export async function getSubscriptions(): Promise<Subscription[]> {
  const res = await authenticatedRequest(`${BACKEND_URL}/api/subscriptions`, {
    method: 'GET',
  });

  if (!res.ok) {
    throw new Error(await res.text());
  }
  return (await res.json()) as Subscription[];
}

export async function getSubscriptionsByCurrency(currency: string): Promise<Subscription[]> {
  const res = await authenticatedRequest(`${BACKEND_URL}/api/subscriptions/${currency}`, {
    method: 'GET',
  });

  if (!res.ok) {
    throw new Error(await res.text());
  }

  return (await res.json()) as Subscription[];
}

export async function deleteSubscription(currency: string, address: string): Promise<void> {
  const res = await authenticatedRequest(`${BACKEND_URL}/api/subscriptions`, {
    method: 'DELETE',
    body: JSON.stringify({ currency, address }),
  });

  if (!res.ok) {
    throw new Error(await res.text());
  }
}

async function authenticatedRequest(url: string, options: RequestInit = {}) {
  const authHeaders = await getAuthHeaders();
  return fetch(url, {
    ...options,
    headers: {
      ...options.headers,
      ...authHeaders,
    },
  });
};

async function getAuthHeaders() {
  const { isAuth, token } = await validateSession();

  if (!isAuth) {
    console.error('user is not autheticated. redirecting login page');
    redirect('/login');
  }

  return {
    'Content-Type': 'application/json',
    ...({ Authorization: `Bearer ${token}` }),
  };
};
