"use client";

import {
  DndContext,
  KeyboardSensor,
  PointerSensor,
  closestCenter,
  useSensor,
  useSensors,
  type DragEndEvent,
} from "@dnd-kit/core";
import {
  SortableContext,
  arrayMove,
  sortableKeyboardCoordinates,
  useSortable,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { useDirection } from "@radix-ui/react-direction";
import type { Table } from "@tanstack/react-table";
import { Eye, EyeOff, GripVertical, Settings2 } from "lucide-react";
import * as React from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { cn } from "@/lib/utils";

interface ColItem {
  id: string;
  label: string;
  visible: boolean;
}

interface SortableColItemProps {
  item: ColItem;
  isSearching: boolean;
  onToggle: (id: string) => void;
}

function SortableColItem({
  item,
  isSearching,
  onToggle,
}: SortableColItemProps) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: item.id });

  return (
    <div
      ref={setNodeRef}
      style={{
        transform: CSS.Transform.toString(transform),
        transition,
        opacity: isDragging ? 0.4 : 1,
      }}
      className="flex items-center gap-1.5 rounded-md px-1 py-1 text-sm hover:bg-accent"
    >
      {!isSearching && (
        <button
          type="button"
          className="cursor-grab touch-none text-muted-foreground/50 hover:text-muted-foreground"
          {...attributes}
          {...listeners}
        >
          <GripVertical className="size-3.5" />
        </button>
      )}
      {isSearching && <div className="w-3.5" />}
      <span
        className={cn(
          "flex-1 truncate",
          !item.visible && "text-muted-foreground",
        )}
      >
        {item.label}
      </span>
      <button
        type="button"
        aria-label={item.visible ? "Hide column" : "Show column"}
        onClick={() => onToggle(item.id)}
        className={cn(
          "rounded p-0.5 transition-colors",
          item.visible
            ? "text-foreground hover:text-muted-foreground"
            : "text-muted-foreground/40 hover:text-muted-foreground",
        )}
      >
        {item.visible ? (
          <Eye className="size-3.5" />
        ) : (
          <EyeOff className="size-3.5" />
        )}
      </button>
    </div>
  );
}

export interface ColumnSavePayload {
  id: string;
  visible: boolean;
  order: number;
}

interface DataGridViewMenuProps<TData> extends React.ComponentProps<
  typeof PopoverContent
> {
  table: Table<TData>;
  disabled?: boolean;
  onSaveColumns?: (columns: ColumnSavePayload[]) => void;
}

export function DataGridViewMenu<TData>({
  table,
  disabled,
  className,
  onSaveColumns,
  ...props
}: DataGridViewMenuProps<TData>) {
  const dir = useDirection();
  const [search, setSearch] = React.useState("");

  const [items, setItems] = React.useState<ColItem[]>(() =>
    table
      .getAllLeafColumns()
      .filter((c) => typeof c.accessorFn !== "undefined" && c.getCanHide())
      .map((c) => ({
        id: c.id,
        label: String(c.columnDef.meta?.label ?? c.id),
        visible: c.getIsVisible(),
      })),
  );

  const sensors = useSensors(
    useSensor(PointerSensor),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    }),
  );

  function applyToTable(nextItems: ColItem[], reordered: boolean) {
    table.setColumnVisibility(
      Object.fromEntries(nextItems.map((c) => [c.id, c.visible])),
    );
    if (reordered) {
      table.setColumnOrder(["select", ...nextItems.map((c) => c.id)]);
    }
    onSaveColumns?.(
      nextItems.map((c, i) => ({ id: c.id, visible: c.visible, order: i + 1 })),
    );
  }

  function handleToggle(id: string) {
    const next = items.map((c) =>
      c.id === id ? { ...c, visible: !c.visible } : c,
    );
    setItems(next);
    applyToTable(next, false);
  }

  function handleDragEnd(event: DragEndEvent) {
    const { active, over } = event;
    if (!over || active.id === over.id) return;
    const oldIdx = items.findIndex((c) => c.id === active.id);
    const newIdx = items.findIndex((c) => c.id === over.id);
    const next = arrayMove(items, oldIdx, newIdx);
    setItems(next);
    applyToTable(next, true);
  }

  const trimmed = search.trim().toLowerCase();
  const filtered = trimmed
    ? items.filter((c) => c.label.toLowerCase().includes(trimmed))
    : items;

  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button
          aria-label="Toggle columns"
          role="combobox"
          dir={dir}
          variant="ghost"
          size="sm"
          className="ms-auto hidden h-8 font-normal lg:flex"
          disabled={disabled}
        >
          <Settings2 className="text-muted-foreground" />
          View
        </Button>
      </PopoverTrigger>
      <PopoverContent
        dir={dir}
        className={cn("w-56 p-2", className)}
        {...props}
      >
        <p className="mb-2 px-1 text-muted-foreground text-xs font-medium">
          Columns
        </p>
        <Input
          placeholder="Search columns..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="mb-1.5 h-7 text-xs"
        />
        <DndContext
          sensors={sensors}
          collisionDetection={closestCenter}
          onDragEnd={handleDragEnd}
        >
          <SortableContext
            items={items.map((c) => c.id)}
            strategy={verticalListSortingStrategy}
          >
            <div className="max-h-64 overflow-y-auto">
              {filtered.length === 0 && (
                <p className="py-2 text-center text-muted-foreground text-xs">
                  No columns found.
                </p>
              )}
              {filtered.map((item) => (
                <SortableColItem
                  key={item.id}
                  item={item}
                  isSearching={!!trimmed}
                  onToggle={handleToggle}
                />
              ))}
            </div>
          </SortableContext>
        </DndContext>
      </PopoverContent>
    </Popover>
  );
}
