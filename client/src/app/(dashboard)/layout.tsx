import { SiteHeader } from "@/components/layouts/site-header";

export default function DashboardLayout({ children }: React.PropsWithChildren) {
  return (
    <div className="flex min-h-screen flex-col">
      <SiteHeader />
      <main className="flex-1">{children}</main>
    </div>
  );
}
