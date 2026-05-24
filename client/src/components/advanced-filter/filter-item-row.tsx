"use client";

import { GripVertical, Trash2 } from "lucide-react";
import * as React from "react";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { FilterValueInput } from "./filter-value-input";
import { getDefaultOperator, getOperators } from "./operators";
import type { FilterField, FilterItemNode } from "./types";

interface FilterItemRowProps {
  item: FilterItemNode;
  fields: FilterField[];
  onUpdate: (updates: Partial<Omit<FilterItemNode, "id" | "type">>) => void;
  onRemove: () => void;
  showHandle?: boolean;
}

export function FilterItemRow({
  item,
  fields,
  onUpdate,
  onRemove,
  showHandle = true,
}: FilterItemRowProps) {
  const selectedField = fields.find((f) => f.value === item.field);
  const operators = selectedField ? getOperators(selectedField.fieldType) : [];

  function handleFieldChange(fieldValue: string) {
    const field = fields.find((f) => f.value === fieldValue);
    if (!field) return;
    onUpdate({
      field: fieldValue,
      operator: getDefaultOperator(field.fieldType),
      value: "",
    });
  }

  function handleOperatorChange(op: string) {
    onUpdate({ operator: op, value: "" });
  }

  return (
    <div className="flex items-center gap-1.5">
      {showHandle && (
        <Button
          variant="ghost"
          size="icon"
          className="drag-handle size-7 cursor-grab text-muted-foreground"
        >
          <GripVertical className="size-3.5" />
        </Button>
      )}

      <Select value={item.field} onValueChange={handleFieldChange}>
        <SelectTrigger className="h-8 min-w-[130px] text-xs">
          <SelectValue placeholder="Select field" />
        </SelectTrigger>
        <SelectContent>
          {fields.map((f) => (
            <SelectItem key={f.value} value={f.value} className="text-xs">
              {f.label}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>

      {operators.length > 0 && (
        <Select value={item.operator} onValueChange={handleOperatorChange}>
          <SelectTrigger className="h-8 min-w-[130px] text-xs">
            <SelectValue placeholder="Operator" />
          </SelectTrigger>
          <SelectContent>
            {operators.map((op) => (
              <SelectItem key={op.value} value={op.value} className="text-xs">
                {op.label}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      )}

      <FilterValueInput
        field={selectedField}
        operator={item.operator}
        value={item.value}
        onChange={(v) => onUpdate({ value: v })}
      />

      <Button
        variant="ghost"
        size="icon"
        className="size-7 text-muted-foreground hover:text-destructive"
        onClick={onRemove}
      >
        <Trash2 className="size-3.5" />
      </Button>
    </div>
  );
}
