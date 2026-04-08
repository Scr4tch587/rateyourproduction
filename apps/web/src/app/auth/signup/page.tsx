import Link from "next/link";
import { AuthForm } from "@/components/auth-form";

export default function SignupPage() {
  return (
    <div className="max-w-sm mx-auto px-4 py-16">
      <h1 className="font-serif text-xl font-bold mb-2">Sign up</h1>
      <p className="text-xs text-muted-foreground mb-6">
        Create an account to log productions, review performances, and submit missing entries.
      </p>
      <AuthForm mode="signup" />
      <div className="mt-4 text-xs">
        <span className="text-muted-foreground">Already have an account? </span>
        <Link href="/auth/login" className="text-burgundy hover:text-burgundy-light transition-colors">
          Log in
        </Link>
      </div>
    </div>
  );
}
