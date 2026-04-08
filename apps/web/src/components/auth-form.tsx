"use client";

import { useState, type FormEvent } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/components/auth-provider";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

interface AuthFormProps {
  mode: "login" | "signup";
}

export function AuthForm({ mode }: AuthFormProps) {
  const router = useRouter();
  const { login, signup } = useAuth();
  const [displayName, setDisplayName] = useState("");
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [submitting, setSubmitting] = useState(false);

  async function handleSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    setSubmitting(true);
    setError("");
    try {
      if (mode === "signup") {
        const user = await signup({
          username,
          email,
          password,
          display_name: displayName || undefined,
        });
        router.push(`/profile/${user.username}`);
      } else {
        const user = await login({ email, password });
        router.push(`/profile/${user.username}`);
      }
      router.refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Request failed");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-3">
      {mode === "signup" && (
        <>
          <Input
            value={displayName}
            onChange={(e) => setDisplayName(e.target.value)}
            placeholder="Display name"
            className="bg-background text-sm"
          />
          <Input
            required
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            placeholder="Username"
            className="bg-background text-sm"
          />
        </>
      )}
      <Input
        required
        type="email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        placeholder="Email"
        className="bg-background text-sm"
      />
      <Input
        required
        type="password"
        minLength={8}
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        placeholder="Password"
        className="bg-background text-sm"
      />
      {error && <div className="text-xs text-destructive">{error}</div>}
      <Button type="submit" disabled={submitting} className="w-full">
        {submitting ? "Submitting..." : mode === "signup" ? "Create account" : "Log in"}
      </Button>
    </form>
  );
}
