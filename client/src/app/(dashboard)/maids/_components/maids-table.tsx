"use client";

import { ModelDataGrid } from "@/app/(dashboard)/_components/model-data-grid";
import type { FormField, UserTab } from "@/types/api";

interface MaidsTableProps {
  formFields: FormField[];
  tab: UserTab;
}

export function MaidsTable({ formFields, tab }: MaidsTableProps) {
  return <ModelDataGrid model="maids" formFields={formFields} tab={tab} />;
}
