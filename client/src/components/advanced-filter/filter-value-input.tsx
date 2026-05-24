"use client";

import * as React from "react";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import type { FilterField } from "./types";

interface FilterValueInputProps {
  field: FilterField | undefined;
  operator: string;
  value: unknown;
  onChange: (value: unknown) => void;
}

const NO_VALUE_OPERATORS = new Set([
  "isEmpty",
  "isNotEmpty",
  "isTrue",
  "isFalse",
]);

export function FilterValueInput({
  field,
  operator,
  value,
  onChange,
}: FilterValueInputProps) {
  if (!field || NO_VALUE_OPERATORS.has(operator)) {
    return null;
  }

  const fieldType = field.fieldType;

  if (
    fieldType === "single_select" ||
    fieldType === "radio_select" ||
    fieldType === "creatable_single_select" ||
    fieldType === "single_relation"
  ) {
    return (
      <Select
        value={typeof value === "string" ? value : ""}
        onValueChange={onChange}
      >
        <SelectTrigger className="h-8 min-w-[120px] text-xs">
          <SelectValue placeholder="Select value" />
        </SelectTrigger>
        <SelectContent>
          {field.options?.map((opt) => (
            <SelectItem key={opt.value} value={opt.value} className="text-xs">
              {opt.label}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    );
  }

  if (
    fieldType === "multi_select" ||
    fieldType === "creatable_multi_select" ||
    fieldType === "multi_relation"
  ) {
    const current = Array.isArray(value) ? value : [];
    const firstVal = typeof current[0] === "string" ? current[0] : "";
    return (
      <Select value={firstVal} onValueChange={(v) => onChange([v])}>
        <SelectTrigger className="h-8 min-w-[120px] text-xs">
          <SelectValue placeholder="Select value" />
        </SelectTrigger>
        <SelectContent>
          {field.options?.map((opt) => (
            <SelectItem key={opt.value} value={opt.value} className="text-xs">
              {opt.label}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    );
  }

  if (
    fieldType === "integer_field" ||
    fieldType === "float_field" ||
    fieldType === "currency_field" ||
    fieldType === "rating_field"
  ) {
    return (
      <Input
        type="number"
        className="h-8 min-w-[120px] text-xs"
        placeholder="Value"
        value={typeof value === "number" ? value : ""}
        onChange={(e) =>
          onChange(e.target.value === "" ? "" : Number(e.target.value))
        }
      />
    );
  }

  if (fieldType === "date_field" || fieldType === "datetime_field") {
    return (
      <Input
        type="date"
        className="h-8 min-w-[140px] text-xs"
        value={typeof value === "string" ? value.slice(0, 10) : ""}
        onChange={(e) => onChange(e.target.value)}
      />
    );
  }

  return (
    <Input
      className="h-8 min-w-[120px] text-xs"
      placeholder="Value"
      value={typeof value === "string" ? value : ""}
      onChange={(e) => onChange(e.target.value)}
    />
  );
}
