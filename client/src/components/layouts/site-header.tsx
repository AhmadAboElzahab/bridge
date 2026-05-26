import { LayoutGrid } from "lucide-react";
import Link from "next/link";
import { ActiveLink } from "@/components/active-link";
import { ModeToggle } from "@/components/layouts/mode-toggle";
import { LogoutButton } from "@/components/logout-button";
import { buttonVariants } from "@/components/ui/button";
import { cn } from "@/lib/utils";

export function SiteHeader() {
  return (
    <header className="sticky top-0 z-50 w-full bg-muted">
      <div className="flex h-9 items-center">
        <Link
          href="/maids"
          className={cn(buttonVariants({ variant: "ghost", size: "icon" }), "size-8")}
        >
          <LayoutGrid />
        </Link>
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
