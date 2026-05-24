import type { Metadata } from "next";
import { Suspense } from "react";
import { DataTableSkeleton } from "@/components/data-table/data-table-skeleton";
import { UsersTableContainer } from "./_components/users-table-container";

export const metadata: Metadata = { title: "Users" };

export default function UsersPage() {
  return (
    <div className="container py-6">
      <div className="mb-6 space-y-1">
        <h1 className="font-semibold text-2xl tracking-tight">Users</h1>
        <p className="text-muted-foreground text-sm">
          Manage and view user accounts
        </p>
      </div>
      <Suspense fallback={<DataTableSkeleton columnCount={6} rowCount={10} />}>
        <UsersTableContainer />
      </Suspense>
    </div>
  );
}
