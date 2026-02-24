import { createContext, useContext, useState, useEffect, ReactNode } from "react";

interface AuthContextType {
  token: string | null;
  setToken: (token: string | null) => void;
  logout: () => void;
  isAdmin: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setTokenState] = useState<string | null>(null);
  const [isAdmin, setIsAdmin] = useState(false);

  useEffect(() => {
    const savedToken = localStorage.getItem("bastion_token");
    if (savedToken) {
      updateState(savedToken);
    }
  }, []);

  const updateState = (newToken: string | null) => {
    if (newToken) {
      try {
        const payload = JSON.parse(atob(newToken.split('.')[1]));
        setIsAdmin(payload.admin === true);
      } catch (e) {
        setIsAdmin(false);
      }
    } else {
      setIsAdmin(false);
    }
    setTokenState(newToken);
  };

  const setToken = (newToken: string | null) => {
    if (newToken) {
      localStorage.setItem("bastion_token", newToken);
    } else {
      localStorage.removeItem("bastion_token");
    }
    updateState(newToken);
  };

  const logout = () => setToken(null);

  return (
    <AuthContext.Provider value={{ token, setToken, logout, isAdmin }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
