import type { FieldType } from "@/types/api";

export type Conjunction = "and" | "or";

export interface FilterItemNode {
  id: string;
  type: "ITEM";
  field: string;
  operator: string;
  value: unknown;
}

export interface FilterGroupNode {
  id: string;
  type: "GROUP";
  conjunction: Conjunction;
  children: (FilterItemNode | FilterGroupNode)[];
}

export type FilterNode = FilterItemNode | FilterGroupNode;

export interface FilterField {
  label: string;
  value: string;
  fieldType: FieldType;
  options?: { label: string; value: string }[];
}
