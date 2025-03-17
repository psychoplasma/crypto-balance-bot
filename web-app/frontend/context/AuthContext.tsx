'use client';

import {
  createContext,
  ReactNode,
  useContext,
  useState,
  useEffect,
} from 'react';
import {
  login as authLogin,
  logout as authLogout,
  validate,
} from '../actions/actions';

interface AuthContextType {
  isAuthenticated: boolean;
  login: (userId: string, token: string) => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType>({
  isAuthenticated: false,
  login: async () => {},
  logout: async () => {},
});

export function AuthProvider({ children }: { children: ReactNode }) {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);

  useEffect(() => {
    async function validateSession() {
      const { success } = await validate();
      setIsAuthenticated(success);
    };
    validateSession();
  }, []);

  const login = async (userId: string, token: string) => {
    await authLogin(userId, token);
    setIsAuthenticated(true);
  };

  const logout = async () => {
    await authLogout();
    setIsAuthenticated(false);
  };

  return (
    <AuthContext.Provider
      value={{
        isAuthenticated,
        login,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export const useAuth = () => useContext(AuthContext);
