"use client";

import { LogOut } from "lucide-react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { useAuthStore } from "@/store/auth";

export function LogoutButton() {
  const router = useRouter();
  const logout = useAuthStore((s) => s.logout);

  function handleLogout() {
    logout();
    router.push("/login");
  }

  return (
    <Button
      variant="ghost"
      size="icon"
      className="size-8"
      onClick={handleLogout}
    >
      <LogOut className="size-4" />
      <span className="sr-only">Sign out</span>
    </Button>
  );
}
