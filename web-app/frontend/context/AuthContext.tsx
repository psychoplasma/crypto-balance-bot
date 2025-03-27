'use client';

import {
  createContext,
  ReactNode,
  useContext,
  useState,
  useEffect,
} from 'react';
import {
  createSessionAction,
  clearSessionAction,
  isAuthenticated,
} from '@/actions/session';

interface AuthContextType {
  isAuthenticated: boolean;
  isLoading: boolean;
  authLogin: (userId: string, token: string) => Promise<void>;
  authLogout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType>({
  isAuthenticated: false,
  isLoading: true,
  authLogin: async () => {},
  authLogout: async () => {},
});

export function AuthProvider({ children }: { children: ReactNode }) {
  const [authenticated, setAuthenticated] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    async function checkAuth() {
      const { isAuth } = await isAuthenticated();
      setAuthenticated(isAuth);
      setIsLoading(false);
    };

    checkAuth();
  }, []);

  const authLogin = async (userId: string, token: string) => {
    await createSessionAction(userId, token);
    setAuthenticated(true);
  };

  const authLogout = async () => {
    await clearSessionAction();
    setAuthenticated(false);
  };

  return (
    <AuthContext.Provider
      value={{
        isAuthenticated: authenticated,
        isLoading,
        authLogin,
        authLogout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export const useAuth = () => useContext(AuthContext);
