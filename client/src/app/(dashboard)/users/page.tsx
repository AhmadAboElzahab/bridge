import type { Metadata } from "next";
import { Suspense } from "react";
import { DataTableSkeleton } from "@/components/data-table/data-table-skeleton";
import { TableContainer } from "@/app/(dashboard)/_components/table-container";

export const metadata: Metadata = { title: "Users" };

export default function UsersPage() {
  return (
    <div className="container ">
      <Suspense fallback={<DataTableSkeleton columnCount={6} rowCount={10} />}>
        <TableContainer entityType="User" model="users" />
      </Suspense>
    </div>
  );
}
