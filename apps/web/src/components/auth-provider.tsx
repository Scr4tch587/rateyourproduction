"use client";

import {
  createContext,
  useContext,
  useEffect,
  useCallback,
  useMemo,
  useState,
  type ReactNode,
} from "react";
import { apiFetch } from "@/lib/api";
import type {
  LoginRequest,
  SessionProfile,
  SignupRequest,
} from "@/lib/types";

interface AuthContextValue {
  user: SessionProfile | null;
  loading: boolean;
  login: (payload: LoginRequest) => Promise<SessionProfile>;
  signup: (payload: SignupRequest) => Promise<SessionProfile>;
  logout: () => Promise<void>;
  refresh: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<SessionProfile | null>(null);
  const [loading, setLoading] = useState(true);

  const refresh = useCallback(async () => {
    try {
      const me = await apiFetch<SessionProfile>("/api/v1/auth/me", {
        cache: "no-store",
      });
      setUser(me);
    } catch {
      setUser(null);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    refresh();
  }, [refresh]);

  const login = useCallback(async (payload: LoginRequest) => {
    const loggedIn = await apiFetch<SessionProfile>("/api/v1/auth/login", {
      method: "POST",
      body: JSON.stringify(payload),
    });
    setUser(loggedIn);
    return loggedIn;
  }, []);

  const signup = useCallback(async (payload: SignupRequest) => {
    const signedUp = await apiFetch<SessionProfile>("/api/v1/auth/signup", {
      method: "POST",
      body: JSON.stringify(payload),
    });
    setUser(signedUp);
    return signedUp;
  }, []);

  const logout = useCallback(async () => {
    await apiFetch<{ status: string }>("/api/v1/auth/logout", {
      method: "POST",
    });
    setUser(null);
  }, []);

  const value = useMemo(
    () => ({ user, loading, login, signup, logout, refresh }),
    [user, loading, login, signup, logout, refresh],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within AuthProvider");
  }
  return context;
}
