import type { Metadata } from "next";
import { LoginForm } from "./_components/login-form";

export const metadata: Metadata = { title: "Sign In" };

export default function LoginPage() {
  return (
    <div className="w-full max-w-sm space-y-6 rounded-xl border bg-card p-8 shadow-sm">
      <div className="space-y-1">
        <h1 className="font-semibold text-2xl tracking-tight">Sign in</h1>
        <p className="text-muted-foreground text-sm">
          Enter your credentials to continue
        </p>
      </div>
      <LoginForm />
    </div>
  );
}
