'use server'

import { createSession, deleteSession, getSession, updateSession } from '@/lib/session';
import { login as loginApi } from '@/lib/api';

export async function login_(email: string, password: string): Promise<{ userId: string, token: string}> {
  const { id, token } = await loginApi(email, password);
  return { userId: id, token: token!! };
}

export async function login(userId: string, token: string) {
  try {
    await createSession(userId, token);
    return { success: true };
  } catch (error) {
    return { success: false, error: 'Failed to create session' };
  }
}

export async function logout() {
  try {
    await deleteSession();
    return { success: true };
  } catch (error) {
    return { success: false, error: 'Failed to delete session' };
  }
}

export async function refresh() {
  try {
    await updateSession();
    return { success: true };
  } catch (error) {
    return { success: false, error: 'Failed to update session' };
  }
}

export async function validate() {
  try {
    const { isAuth } = await getSession();
    return { success: isAuth };
  } catch (error) {
    return { success: false, error: 'Failed to update session' };
  }
}

export async function validateAndGet() {
  try {
    const { isAuth, token, userId } = await getSession();
    return { error: null, isAuth, token, userId };
  } catch (error) {
    return { error, isAuth: false };
  }
}
