import 'server-only';
import { cookies } from 'next/headers';
import { cache } from 'react';
import { JWTPayload, SignJWT, jwtVerify } from 'jose';
 
const secretKey = process.env.SESSION_SECRET;
const encodedKey = new TextEncoder().encode(secretKey);

interface SessionPayload extends JWTPayload {
  userId: string;
  token: string;
}
 
export async function encrypt(payload: SessionPayload) {
  return new SignJWT(payload)
    .setProtectedHeader({ alg: 'HS256' })
    .setIssuedAt()
    .setExpirationTime('7d')
    .sign(encodedKey);
}
 
export async function decrypt(session: string | undefined = '') {
  try {
    const { payload } = await jwtVerify(session, encodedKey, {
      algorithms: ['HS256'],
    });
    return payload;
  } catch (error) {
    console.log('Failed to verify session');
  }
}
 
export async function createSession(userId: string, token: string) {
  const expiresAt = new Date(Date.now() + 7 * 24 * 60 * 60 * 1000);
  const session = await encrypt({ userId, token, expiresAt });
  const cookieStore = await cookies();
 
  cookieStore.set('session', session, {
    httpOnly: true,
    secure: true,
    expires: expiresAt,
    sameSite: 'lax',
    path: '/',
  });
}

export async function updateSession() {
  const session = (await cookies()).get('session')?.value;
  const payload = await decrypt(session);
 
  if (!session || !payload) {
    return null;
  }
 
  const expires = new Date(Date.now() + 7 * 24 * 60 * 60 * 1000);
 
  const cookieStore = await cookies();
  cookieStore.set('session', session, {
    httpOnly: true,
    secure: true,
    expires: expires,
    sameSite: 'lax',
    path: '/',
  });
}

export const getSession = cache(async () => {
  const cookie = (await cookies()).get('session')?.value;
  const session = await decrypt(cookie);
 
  return {
    isAuth: !!session?.userId || !!session?.token,
    userId: session?.userId,
    token: session?.token,
  };
});

export async function deleteSession() {
  const cookieStore = await cookies();
  cookieStore.delete('session');
}

interface SessionOptions {
  path?: string;
  secure?: boolean;
  sameSite?: 'strict' | 'lax' | 'none';
  maxAge?: number;
}
