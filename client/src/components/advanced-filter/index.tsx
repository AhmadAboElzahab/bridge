"use client";

import { ListFilter, X } from "lucide-react";
import * as React from "react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Separator } from "@/components/ui/separator";
import { useAdvancedFilter } from "@/hooks/use-advanced-filter";
import { FilterGroup } from "./filter-group";
import { toApiFilter } from "./to-api-filter";
import type { FilterField, FilterGroupNode } from "./types";
import { countActiveFilters } from "./utils";

export type { FilterField, FilterGroupNode };

interface AdvancedFilterProps {
  fields: FilterField[];
  value?: FilterGroupNode;
  onChange: (apiFilter: ReturnType<typeof toApiFilter>) => void;
  align?: "start" | "center" | "end";
}

export function AdvancedFilter({
  fields,
  value,
  onChange,
  align = "end",
}: AdvancedFilterProps) {
  const [open, setOpen] = React.useState(false);

  const {
    ruleset,
    addItem,
    addGroup,
    removeNode,
    updateItem,
    updateConjunction,
    reorderChildren,
    clearAll,
  } = useAdvancedFilter(value, (next) => {
    onChange(toApiFilter(next));
  });

  const activeCount = countActiveFilters(ruleset);

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant="ghost"
          size="sm"
          className="h-8 gap-1.5 text-sm"
          aria-label="Filters"
        >
          <ListFilter className="size-3.5" />
          Filters
          {activeCount > 0 && (
            <Badge
              variant="secondary"
              className="size-4 rounded-full p-0 text-[10px] leading-none"
            >
              {activeCount}
            </Badge>
          )}
        </Button>
      </PopoverTrigger>

      <PopoverContent
        align={align}
        className="w-auto min-w-[480px] max-w-[700px] p-0"
        onInteractOutside={(e) => {
          if ((e.target as HTMLElement)?.closest("[role='dialog']")) return;
          setOpen(false);
        }}
      >
        <div className="flex items-center justify-between border-b px-3 py-2">
          <span className="font-medium text-sm">Filters</span>
          {activeCount > 0 && (
            <Button
              variant="ghost"
              size="sm"
              className="h-7 gap-1 px-2 text-muted-foreground text-xs"
              onClick={clearAll}
            >
              <X className="size-3" />
              Clear all
            </Button>
          )}
        </div>

        <div className="p-3">
          {ruleset.children.length === 0 ? (
            <p className="py-2 text-center text-muted-foreground text-xs">
              No filters applied. Add a filter to get started.
            </p>
          ) : (
            <FilterGroup
              group={ruleset}
              fields={fields}
              onAddItem={addItem}
              onAddGroup={addGroup}
              onRemoveNode={removeNode}
              onUpdateItem={updateItem}
              onUpdateConjunction={updateConjunction}
              onReorderChildren={reorderChildren}
            />
          )}
        </div>

        {ruleset.children.length === 0 && (
          <>
            <Separator />
            <div className="p-2">
              <Button
                variant="ghost"
                size="sm"
                className="w-full justify-start gap-1.5 text-xs"
                onClick={() => {
                  const first = fields[0];
                  if (first) addItem(ruleset.id, first.value, "contains");
                }}
                disabled={fields.length === 0}
              >
                <ListFilter className="size-3.5" />
                Add first filter
              </Button>
            </div>
          </>
        )}
      </PopoverContent>
    </Popover>
  );
}
