export default function AuthLayout({ children }: React.PropsWithChildren) {
  return (
    <div className="flex min-h-screen items-center justify-center bg-muted/40">
      {children}
    </div>
  );
}
