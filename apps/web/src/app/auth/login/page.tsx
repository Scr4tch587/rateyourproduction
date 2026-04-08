import Link from "next/link";
import { AuthForm } from "@/components/auth-form";

export default function LoginPage() {
  return (
    <div className="max-w-sm mx-auto px-4 py-16">
      <h1 className="font-serif text-xl font-bold mb-2">Log in</h1>
      <p className="text-xs text-muted-foreground mb-6">
        Access your theatre journal and continue logging productions.
      </p>
      <AuthForm mode="login" />
      <div className="mt-4 text-xs">
        <span className="text-muted-foreground">No account? </span>
        <Link href="/auth/signup" className="text-burgundy hover:text-burgundy-light transition-colors">
          Sign up
        </Link>
      </div>
    </div>
  );
}
