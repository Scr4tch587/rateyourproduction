"use client";

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from "react";
import { apiFetch } from "@/lib/api";
import { getSupabaseClient } from "@/lib/supabase";
import type { SessionProfile } from "@/lib/types";

interface SignupPayload {
  username: string;
  email: string;
  password: string;
  display_name?: string;
}

interface LoginPayload {
  email: string;
  password: string;
}

interface AuthContextValue {
  user: SessionProfile | null;
  accessToken: string | null;
  loading: boolean;
  login: (payload: LoginPayload) => Promise<SessionProfile>;
  signup: (payload: SignupPayload) => Promise<SessionProfile | null>;
  logout: () => Promise<void>;
  refresh: (tokenOverride?: string | null) => Promise<SessionProfile | null>;
}

const AuthContext = createContext<AuthContextValue | null>(null);

async function fetchProfile(accessToken: string): Promise<SessionProfile> {
  return apiFetch<SessionProfile>("/api/v1/auth/me", {
    cache: "no-store",
    headers: {
      Authorization: `Bearer ${accessToken}`,
    },
  });
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const supabase = useMemo(() => getSupabaseClient(), []);
  const [user, setUser] = useState<SessionProfile | null>(null);
  const [accessToken, setAccessToken] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  const refresh = useCallback(async (tokenOverride?: string | null) => {
    let token = tokenOverride ?? null;

    if (tokenOverride === undefined) {
      const {
        data: { session },
      } = await supabase.auth.getSession();
      token = session?.access_token ?? null;
    }

    setAccessToken(token);

    if (!token) {
      setUser(null);
      setLoading(false);
      return null;
    }

    try {
      const profile = await fetchProfile(token);
      setUser(profile);
      return profile;
    } catch {
      setUser(null);
      return null;
    } finally {
      setLoading(false);
    }
  }, [supabase]);

  useEffect(() => {
    void refresh();

    const {
      data: { subscription },
    } = supabase.auth.onAuthStateChange((_event, session) => {
      void refresh(session?.access_token ?? null);
    });

    return () => subscription.unsubscribe();
  }, [refresh, supabase]);

  const login = useCallback(async (payload: LoginPayload) => {
    const { data, error } = await supabase.auth.signInWithPassword(payload);
    if (error) {
      throw error;
    }
    const profile = await refresh(data.session?.access_token ?? null);
    if (!profile) {
      throw new Error("Authenticated, but profile could not be loaded.");
    }
    return profile;
  }, [refresh, supabase]);

  const signup = useCallback(async (payload: SignupPayload) => {
    const { data, error } = await supabase.auth.signUp({
      email: payload.email,
      password: payload.password,
      options: {
        data: {
          username: payload.username,
          display_name: payload.display_name,
        },
      },
    });
    if (error) {
      throw error;
    }

    if (!data.session?.access_token) {
      setLoading(false);
      return null;
    }

    return refresh(data.session.access_token);
  }, [refresh, supabase]);

  const logout = useCallback(async () => {
    const { error } = await supabase.auth.signOut();
    if (error) {
      throw error;
    }
    setAccessToken(null);
    setUser(null);
  }, [supabase]);

  const value = useMemo(
    () => ({ user, accessToken, loading, login, signup, logout, refresh }),
    [user, accessToken, loading, login, signup, logout, refresh],
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
