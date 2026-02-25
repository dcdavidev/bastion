import type { ReactNode } from 'react';
import { createContext, use, useMemo, useState } from 'react';

interface AuthContextType {
  token: string | null;
  setToken: (token: string | null) => void;
  logout: () => void;
  isAdmin: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setToken] = useState<string | null>(() => {
    if (globalThis.window !== undefined) {
      return localStorage.getItem('bastion_token');
    }
    return null;
  });

  const isAdmin = useMemo(() => {
    if (!token) return false;
    try {
      const payload = JSON.parse(atob(token.split('.')[1]));
      return payload.admin === true;
    } catch {
      return false;
    }
  }, [token]);

  const handleSetToken = (newToken: string | null) => {
    if (globalThis.window !== undefined) {
      if (newToken) {
        localStorage.setItem('bastion_token', newToken);
      } else {
        localStorage.removeItem('bastion_token');
      }
    }
    setToken(newToken);
  };

  const logout = () => handleSetToken(null);

  return (
    <AuthContext value={{ token, setToken: handleSetToken, logout, isAdmin }}>
      {children}
    </AuthContext>
  );
}

export function useAuth() {
  const context = use(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
