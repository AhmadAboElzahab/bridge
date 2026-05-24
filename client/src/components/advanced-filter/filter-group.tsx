"use client";

import { GripVertical, Plus, PlusSquare } from "lucide-react";
import * as React from "react";
import { type ItemInterface, ReactSortable } from "react-sortablejs";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { MAX_NESTING_LEVEL } from "@/hooks/use-advanced-filter";
import { cn } from "@/lib/utils";
import { FilterItemRow } from "./filter-item-row";
import { getDefaultOperator } from "./operators";
import type {
  FilterField,
  FilterGroupNode,
  FilterItemNode,
  FilterNode,
} from "./types";

interface SortableNode extends ItemInterface {
  id: string;
  node: FilterNode;
}

interface FilterGroupProps {
  group: FilterGroupNode;
  fields: FilterField[];
  depth?: number;
  onAddItem: (groupId: string, field: string, operator: string) => void;
  onAddGroup: (parentId: string) => void;
  onRemoveNode: (id: string) => void;
  onUpdateItem: (
    id: string,
    updates: Partial<Omit<FilterItemNode, "id" | "type">>,
  ) => void;
  onUpdateConjunction: (groupId: string, conjunction: "and" | "or") => void;
  onReorderChildren: (groupId: string, newChildren: FilterNode[]) => void;
  showHandle?: boolean;
}

export function FilterGroup({
  group,
  fields,
  depth = 0,
  onAddItem,
  onAddGroup,
  onRemoveNode,
  onUpdateItem,
  onUpdateConjunction,
  onReorderChildren,
  showHandle = false,
}: FilterGroupProps) {
  const defaultField = fields[0];
  const isNested = depth > 0;

  const sortableItems: SortableNode[] = group.children.map((c) => ({
    id: c.id,
    node: c,
  }));

  function handleSetList(newItems: SortableNode[]) {
    onReorderChildren(
      group.id,
      newItems.map((item) => item.node),
    );
  }

  function handleAddItem() {
    if (!defaultField) return;
    onAddItem(
      group.id,
      defaultField.value,
      getDefaultOperator(defaultField.fieldType),
    );
  }

  return (
    <div
      className={cn(
        "flex flex-col gap-2",
        isNested && "rounded-md border border-border bg-muted/30 p-2",
      )}
    >
      {/* Group drag handle + header */}
      <div className="flex items-center gap-1.5">
        {showHandle && (
          <button
            type="button"
            className="drag-handle cursor-grab text-muted-foreground hover:text-foreground"
            aria-label="Drag group"
          >
            <GripVertical className="size-3.5" />
          </button>
        )}

        {/* AND / OR — always visible when there are children */}
        {group.children.length > 0 && (
          <div className="flex items-center gap-1.5">
            <span className="text-muted-foreground text-xs">Match</span>
            <Select
              value={group.conjunction}
              onValueChange={(v) =>
                onUpdateConjunction(group.id, v as "and" | "or")
              }
            >
              <SelectTrigger className="h-7 w-[70px] text-xs">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="and" className="text-xs">
                  All
                </SelectItem>
                <SelectItem value="or" className="text-xs">
                  Any
                </SelectItem>
              </SelectContent>
            </Select>
            <span className="text-muted-foreground text-xs">of the following</span>
          </div>
        )}

        {isNested && (
          <Button
            variant="ghost"
            size="sm"
            className="ml-auto h-6 px-2 text-muted-foreground text-xs hover:text-destructive"
            onClick={() => onRemoveNode(group.id)}
          >
            Remove group
          </Button>
        )}
      </div>

      {/* Children — items and nested groups, all draggable */}
      <ReactSortable
        list={sortableItems}
        setList={handleSetList}
        group={{ name: "advanced-filter", pull: true, put: true }}
        animation={150}
        handle=".drag-handle"
        ghostClass="opacity-40"
        dragClass="shadow-md"
        className="flex flex-col gap-1.5"
      >
        {group.children.map((child) => (
          <div key={child.id}>
            {child.type === "ITEM" ? (
              <FilterItemRow
                item={child}
                fields={fields}
                onUpdate={(updates) => onUpdateItem(child.id, updates)}
                onRemove={() => onRemoveNode(child.id)}
                showHandle
              />
            ) : (
              <FilterGroup
                group={child}
                fields={fields}
                depth={depth + 1}
                showHandle
                onAddItem={onAddItem}
                onAddGroup={onAddGroup}
                onRemoveNode={onRemoveNode}
                onUpdateItem={onUpdateItem}
                onUpdateConjunction={onUpdateConjunction}
                onReorderChildren={onReorderChildren}
              />
            )}
          </div>
        ))}
      </ReactSortable>

      {/* Add buttons */}
      <div className="flex items-center gap-1.5">
        <Button
          variant="ghost"
          size="sm"
          className="h-7 gap-1.5 px-2 text-xs"
          onClick={handleAddItem}
        >
          <Plus className="size-3.5" />
          Add filter
        </Button>
        {depth < MAX_NESTING_LEVEL && (
          <Button
            variant="ghost"
            size="sm"
            className="h-7 gap-1.5 px-2 text-xs"
            onClick={() => onAddGroup(group.id)}
          >
            <PlusSquare className="size-3.5" />
            Add group
          </Button>
        )}
      </div>
    </div>
  );
}
