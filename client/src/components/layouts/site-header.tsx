import { LayoutGrid } from "lucide-react";
import Link from "next/link";
import { ActiveLink } from "@/components/active-link";
import { ModeToggle } from "@/components/layouts/mode-toggle";
import { LogoutButton } from "@/components/logout-button";
import { Button } from "@/components/ui/button";

export function SiteHeader() {
  return (
    <header className="sticky top-0 z-50 w-full bg-muted">
      <div className="flex h-9 items-center">
        <Button variant="ghost" size="icon" className="size-8" asChild>
          <Link href="/maids" viewTransition>
            <LayoutGrid />
          </Link>
        </Button>
        <nav className="flex w-full items-center text-sm">
          <ActiveLink href="/maids">Maids</ActiveLink>
          <ActiveLink href="/users">Users</ActiveLink>
        </nav>
        <nav className="flex flex-1 items-center gap-1 md:justify-end">
          <ModeToggle />
          <LogoutButton />
        </nav>
      </div>
    </header>
  );
}
