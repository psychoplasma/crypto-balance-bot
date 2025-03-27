'use server'

import {
  createSession,
  deleteSession,
  updateSession,
  validateSession,
} from '@/lib/session';

export interface AuthResult {
  isAuth: boolean;
  userId?: string;
  token?: string;
}

export async function createSessionAction(userId: string, token: string) {
  try {
    await createSession(userId, token);
    return { success: true };
  } catch (error) {
    console.error('Error while creating session:', error);
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

export async function isAuthenticated(): Promise<AuthResult> {
  try {
    return validateSession();
  } catch (error) {
    console.error('Error while validating session:', error);
    return { isAuth: false };
  }
}
