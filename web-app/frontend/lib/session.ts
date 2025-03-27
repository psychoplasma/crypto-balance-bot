import 'server-only';
import { cache } from 'react';
import { JWTPayload, SignJWT, jwtVerify } from 'jose';
import { cookies } from 'next/headers';
import { NextRequest, NextResponse } from "next/server";

const SESSION_SECRET = new TextEncoder().encode(process.env.SESSION_SECRET_KEY);
const COOKIE_NAME = 'user_session';
const EXPIRY = 7 * 24 * 60 * 60 * 1000; // 7 days

interface SessionPayload extends JWTPayload {
  userId: string;
  token: string;
  expiresAt: Date;
}

interface SessionOptions {
  path?: string;
  secure?: boolean;
  sameSite?: 'strict' | 'lax' | 'none';
  maxAge?: number;
}
 
export async function createSession(userId: string, token: string) {
  if (!token) {
    throw new Error('empty jwt token');
  }

  const expiresAt = new Date(Date.now() + EXPIRY);
  const session = await encrypt({ userId, token, expiresAt });
  const cookieStore = await cookies();
 
  cookieStore.set(COOKIE_NAME, session, {
    httpOnly: true,
    secure: process.env.NODE_ENV === 'production',
    expires: expiresAt,
    sameSite: 'lax',
    path: '/',
  });
}

export async function updateSession() {
  const cookieStore = await cookies();
  const session = cookieStore.get(COOKIE_NAME)?.value;
  const payload = await decrypt(session);
 
  if (!session || !payload) {
    return null;
  }
 
  const expires = new Date(Date.now() + EXPIRY);

  cookieStore.set(COOKIE_NAME, session, {
    httpOnly: true,
    secure: true,
    expires: expires,
    sameSite: 'lax',
    path: '/',
  });
}

export async function deleteSession() {
  const cookieStore = await cookies();
  cookieStore.delete(COOKIE_NAME);
}

export async function updateSession_(request: NextRequest) {
  const session = request.cookies.get(COOKIE_NAME)?.value;
  const payload = await decrypt(session);
 
  if (!session || !payload) {
    return null;
  }

  // Refresh the session so it doesn't expire
  payload.expiresAt = new Date(Date.now() + EXPIRY);

  const res = NextResponse.next();
  res.cookies.set({
    name: COOKIE_NAME,
    value: await encrypt(payload),
    httpOnly: true,
    expires: payload.expiresAt,
  });

  return res;
}

export async function validateSession() {
  const payload = await getSession();

  if (!payload
    || isPast(new Date(payload.expiresAt))
    || !payload.userId
    || !payload.token
  ) {
    return {
      isAuth: false,
    }
  }

  return {
    isAuth: true,
    userId: payload?.userId,
    token: payload?.token,
  }
}

const getSession = cache(async (): Promise<SessionPayload | null> => {
  const cookieStore = await cookies();
  const cookie = cookieStore.get(COOKIE_NAME)?.value;

  if (!cookie) {
    return null;
  }

  const payload = await decrypt(cookie);

  if (!payload) {
    return null;
  }
 
  return payload;
});

function isPast(time: Date): boolean {
  return time.getTime() < Date.now();
}

async function encrypt(payload: SessionPayload) {
  return await new SignJWT(payload)
    .setProtectedHeader({ alg: 'HS256' })
    .setIssuedAt()
    .setExpirationTime('24h')
    .sign(SESSION_SECRET);
}

async function decrypt(input: string | undefined = ''): Promise<SessionPayload> {
  const { payload } = await jwtVerify(input, SESSION_SECRET, {
    algorithms: ["HS256"],
  });

  return payload as SessionPayload;
}
