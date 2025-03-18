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
  validateSessionAction,
} from '../actions/session';

interface AuthContextType {
  isAuthenticated: boolean;
  authLogin: (userId: string, token: string) => Promise<void>;
  authLogout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType>({
  isAuthenticated: false,
  authLogin: async () => {},
  authLogout: async () => {},
});

export function AuthProvider({ children }: { children: ReactNode }) {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);

  useEffect(() => {
    async function validateSession() {
      const { success } = await validateSessionAction();

      console.log(`authenticating: ${success}`);
      setIsAuthenticated(success);
    };
    validateSession();
  }, []);

  const authLogin = async (userId: string, token: string) => {
    await createSessionAction(userId, token);
    setIsAuthenticated(true);
  };

  const authLogout = async () => {
    await clearSessionAction();
    setIsAuthenticated(false);
  };

  return (
    <AuthContext.Provider
      value={{
        isAuthenticated,
        authLogin,
        authLogout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export const useAuth = () => useContext(AuthContext);
