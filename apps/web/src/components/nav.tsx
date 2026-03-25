import Link from "next/link";

export function Nav() {
  return (
    <header className="border-b border-border bg-card">
      <div className="max-w-7xl mx-auto px-4 h-10 flex items-center justify-between">
        <div className="flex items-center gap-6">
          <Link href="/" className="font-bold text-sm tracking-tight">
            RateYourProduction
          </Link>
          <nav className="flex items-center gap-4 text-xs text-muted-foreground">
            <Link href="/discover" className="hover:text-foreground transition-colors">
              Discover
            </Link>
          </nav>
        </div>
        <div className="flex items-center gap-4 text-xs text-muted-foreground">
          <Link href="/auth/login" className="hover:text-foreground transition-colors">
            Log in
          </Link>
          <Link
            href="/auth/signup"
            className="hover:text-foreground transition-colors"
          >
            Sign up
          </Link>
        </div>
      </div>
    </header>
  );
}
