'use server'

import {
  createSession,
  deleteSession,
  updateSession,
  validateSession,
} from '@/lib/session';
import { login as loginApi } from '@/lib/api';

export async function login_(email: string, password: string): Promise<{ userId: string, token: string}> {
  const { id, accessToken } = await loginApi(email, password);
  return { userId: id, token: accessToken!! };
}

export async function createSessionAction(userId: string, token: string) {
  try {
    await createSession(userId, token);
    return { success: true };
  } catch (error) {
    return { success: false, error: 'Failed to create session' };
  }
}

export async function clearSessionAction() {
  try {
    await deleteSession();
    return { success: true };
  } catch (error) {
    return { success: false, error: 'Failed to delete session' };
  }
}

export async function refreshSessionAction() {
  try {
    await updateSession();
    return { success: true };
  } catch (error) {
    return { success: false, error: 'Failed to update session' };
  }
}

export async function validateSessionAction() {
  try {
    const { isAuth } = await validateSession();
    return { success: isAuth };
  } catch (error) {
    return { success: false, error: 'Failed to update session' };
  }
}

export async function validateAndGet() {
  try {
    const { isAuth, token, userId } = await validateSession();
    return { success: true, isAuth, token, userId };
  } catch (error) {
    return { success: false, error: 'Failed to update session' };
  }
}
