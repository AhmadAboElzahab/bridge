import type { Metadata } from "next";
import { Suspense } from "react";
import { DataTableSkeleton } from "@/components/data-table/data-table-skeleton";
import { MaidsTableContainer } from "./_components/maids-table-container";

export const metadata: Metadata = { title: "Maids" };

export default function MaidsPage() {
  return (
    <div className="container py-6">
      <div className="mb-6 space-y-1">
        <h1 className="font-semibold text-2xl tracking-tight">Maids</h1>
        <p className="text-muted-foreground text-sm">
          Manage and view maid profiles
        </p>
      </div>
      <Suspense fallback={<DataTableSkeleton columnCount={6} rowCount={10} />}>
        <MaidsTableContainer />
      </Suspense>
    </div>
  );
}
