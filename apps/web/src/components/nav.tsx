"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState, type FormEvent } from "react";
import { useAuth } from "@/components/auth-provider";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

export function Nav() {
  const router = useRouter();
  const { user, loading, logout } = useAuth();
  const [query, setQuery] = useState("");

  async function handleLogout() {
    await logout();
    router.push("/");
    router.refresh();
  }

  function handleSearch(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    router.push(`/discover?q=${encodeURIComponent(query.trim())}`);
  }

  return (
    <header className="border-b border-border bg-card">
      <div className="max-w-7xl mx-auto px-4 py-2 flex flex-col gap-2 md:flex-row md:items-center md:justify-between">
        <div className="flex items-center gap-6">
          <Link href="/" className="font-serif font-bold text-sm tracking-tight text-burgundy">
            RateYourProduction
          </Link>
          <nav className="flex items-center gap-5 text-xs text-muted-foreground">
            <Link href="/discover" className="hover:text-foreground transition-colors">
              Discover
            </Link>
            {user?.is_admin && (
              <Link href="/admin" className="hover:text-foreground transition-colors">
                Admin
              </Link>
            )}
          </nav>
        </div>

        <div className="flex flex-col gap-2 md:flex-row md:items-center md:gap-4">
          <form onSubmit={handleSearch} className="flex items-center gap-2">
            <Input
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              placeholder="Search works"
              className="h-7 min-w-[180px] text-xs bg-background"
            />
            <Button type="submit" size="xs">Search</Button>
          </form>

          <div className="flex items-center gap-4 text-xs text-muted-foreground">
            {loading ? null : user ? (
              <>
                <Link href={`/profile/${user.username}`} className="hover:text-foreground transition-colors">
                  @{user.username}
                </Link>
                <button
                  type="button"
                  onClick={handleLogout}
                  className="hover:text-foreground transition-colors cursor-pointer"
                >
                  Log out
                </button>
              </>
            ) : (
              <>
                <Link href="/auth/login" className="hover:text-foreground transition-colors">
                  Log in
                </Link>
                <Link
                  href="/auth/signup"
                  className="text-burgundy hover:text-burgundy-light transition-colors"
                >
                  Sign up
                </Link>
              </>
            )}
          </div>
        </div>
      </div>
    </header>
  );
}
