"use client";

import { ModelDataGrid } from "@/app/(dashboard)/_components/model-data-grid";
import type { FormField, UserTab } from "@/types/api";

interface UsersTableProps {
  formFields: FormField[];
  tab: UserTab;
}

export function UsersTable({ formFields, tab }: UsersTableProps) {
  return <ModelDataGrid model="users" formFields={formFields} tab={tab} />;
}
