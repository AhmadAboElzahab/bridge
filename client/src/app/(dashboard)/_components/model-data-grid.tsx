"use client";

import { DirectionProvider } from "@radix-ui/react-direction";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import type { ColumnDef, SortingState } from "@tanstack/react-table";
import * as React from "react";
import { AdvancedFilter } from "@/components/advanced-filter";
import { fromApiFilter } from "@/components/advanced-filter/to-api-filter";
import type { FilterField } from "@/components/advanced-filter/types";
import { DataGrid } from "@/components/data-grid/data-grid";
import { DataGridKeyboardShortcuts } from "@/components/data-grid/data-grid-keyboard-shortcuts";
import { DataGridRowHeightMenu } from "@/components/data-grid/data-grid-row-height-menu";
import { getDataGridSelectColumn } from "@/components/data-grid/data-grid-select-column";
import { DataGridSortMenu } from "@/components/data-grid/data-grid-sort-menu";
import type { ColumnSavePayload } from "@/components/data-grid/data-grid-view-menu";
import { DataGridViewMenu } from "@/components/data-grid/data-grid-view-menu";
import { useDataGrid } from "@/hooks/use-data-grid";
import { useIsomorphicLayoutEffect } from "@/hooks/use-isomorphic-layout-effect";
import { useWindowSize } from "@/hooks/use-window-size";
import { getFilterFn } from "@/lib/data-grid-filters";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import {
  fetchModelIndex,
  PAGE_SIZE,
  updateModelRow,
} from "@/services/data.service";
import { updateTabColumns, updateTabSorting } from "@/services/tabs.service";
import type { FieldType, FilterGroup, FormField, UserTab } from "@/types/api";
import type { CellOpts } from "@/types/data-grid";

type Row = Record<string, unknown>;

// Matches ISO-8601 date/datetime strings, e.g. "1987-02-02T04:00:00+04:00"
const ISO_DATE_RE = /^\d{4}-\d{2}-\d{2}(T[\d:.+Z-].*)?$/;

function guessMediaType(url: string): string {
  const ext = (url.split("?")[0]?.split(".").pop() ?? "").toLowerCase();
  if (["jpg", "jpeg", "png", "gif", "webp", "svg", "avif"].includes(ext))
    return "image/jpeg";
  if (["mp4", "webm", "mov", "avi"].includes(ext)) return "video/mp4";
  if (["mp3", "wav", "ogg"].includes(ext)) return "audio/mpeg";
  if (ext === "pdf") return "application/pdf";
  return "application/octet-stream";
}

function urlBasename(url: string): string {
  return url.split("?")[0]?.split("/").pop() ?? url;
}

function toFileCellData(
  val: unknown,
): import("@/types/data-grid").FileCellData[] {
  if (!val) return [];

  const toItem = (
    v: unknown,
  ): import("@/types/data-grid").FileCellData | null => {
    if (typeof v === "string") {
      if (!v) return null;
      return {
        id: v,
        name: urlBasename(v),
        size: 0,
        type: guessMediaType(v),
        url: v,
      };
    }
    if (typeof v === "object" && v !== null) {
      const o = v as Record<string, unknown>;
      const url = String(o.url ?? o.file_url ?? o.path ?? "");
      if (!url) return null;
      return {
        id: String(o.id ?? url),
        name: String(o.name ?? o.filename ?? urlBasename(url)),
        size: Number(o.size ?? o.file_size ?? 0),
        type: String(
          o.type ?? o.mime_type ?? o.content_type ?? guessMediaType(url),
        ),
        url,
      };
    }
    return null;
  };

  if (Array.isArray(val)) {
    return val
      .map(toItem)
      .filter((x): x is import("@/types/data-grid").FileCellData => x !== null);
  }
  const item = toItem(val);
  return item ? [item] : [];
}

function toStr(val: unknown): string {
  if (val === null || val === undefined) return "";
  if (typeof val === "object") {
    const o = val as Record<string, unknown>;
    return String(o.label ?? o.name ?? o.title ?? o.url ?? "");
  }
  return String(val);
}

function toSelectId(val: unknown): string {
  if (val === null || val === undefined) return "";
  if (typeof val === "object") {
    const o = val as Record<string, unknown>;
    return String(o.value ?? o.id ?? "");
  }
  return String(val);
}

function normalizeValue(val: unknown, fieldType: FieldType): unknown {
  if (val === null || val === undefined) return null;

  switch (fieldType) {
    case "date_field":
    case "datetime_field":
      return typeof val === "string" ? val : null;

    case "integer_field":
    case "float_field":
    case "currency_field":
    case "rating_field":
      return typeof val === "number" ? val : Number(val) || null;

    case "boolean_field":
      return Boolean(val);

    case "single_select":
    case "radio_select":
    case "creatable_single_select":
    case "single_relation":
      return toSelectId(val);

    case "multi_select":
    case "creatable_multi_select":
    case "multi_relation":
      return Array.isArray(val) ? val.map(toSelectId) : [];

    case "url_field":
    case "social_url_field":
    case "link_field":
    case "image_field":
    case "file_field":
    case "video_field":
    case "attachments_field":
      return toFileCellData(val);

    default:
      if (typeof val === "object" && !Array.isArray(val)) {
        return toStr(val);
      }
      if (typeof val === "string" && ISO_DATE_RE.test(val)) {
        return val;
      }
      return val;
  }
}

function normalizeRow(row: Row, formFields: FormField[]): Row {
  const out: Row = { id: row.id };
  for (const field of formFields) {
    out[field.field_key] = normalizeValue(
      row[field.field_key],
      field.form_field_type,
    );
  }
  return out;
}

function fieldTypeToCellOpts(field: FormField): CellOpts {
  const options = (field.options ?? []).map((o) => ({
    label: String(o.label),
    value: String(o.value),
  }));

  switch (field.form_field_type) {
    case "date_field":
    case "datetime_field":
      return { variant: "date" };

    case "integer_field":
      return { variant: "number", step: 1 };
    case "float_field":
    case "currency_field":
      return { variant: "number", step: 0.01 };
    case "rating_field":
      return { variant: "number", min: 0, max: 5, step: 1 };

    case "boolean_field":
      return { variant: "checkbox" };

    case "single_select":
    case "radio_select":
    case "creatable_single_select":
    case "single_relation":
      return { variant: "select", options };

    case "multi_select":
    case "creatable_multi_select":
    case "multi_relation":
      return { variant: "multi-select", options };

    case "rich_field":
      return { variant: "long-text" };

    case "url_field":
    case "social_url_field":
    case "link_field":
      return { variant: "url" };

    case "image_field":
      return {
        variant: "file",
        accept: "image/*",
        maxFileSize: 10 * 1024 * 1024,
      };
    case "video_field":
      return {
        variant: "file",
        accept: "video/*",
        maxFileSize: 100 * 1024 * 1024,
      };
    case "file_field":
      return { variant: "file", maxFileSize: 10 * 1024 * 1024 };
    case "attachments_field":
      return {
        variant: "file",
        multiple: true,
        maxFiles: 10,
        maxFileSize: 10 * 1024 * 1024,
        accept: "image/*,video/*,audio/*,.pdf,.doc,.docx,.xls,.xlsx,.ppt,.pptx",
      };

    default:
      return { variant: "short-text" };
  }
}

interface ModelDataGridProps {
  model: string;
  formFields: FormField[];
  tab: UserTab;
}

export function ModelDataGrid({ model, formFields, tab }: ModelDataGridProps) {
  const queryClient = useQueryClient();
  const windowSize = useWindowSize({ defaultHeight: 760 });
  const [data, setData] = React.useState<Row[]>([]);
  const prevDataRef = React.useRef<Row[]>([]);
  const [page, setPage] = React.useState(1);
  const [apiFilters, setApiFilters] = React.useState<
    FilterGroup | Record<string, never>
  >(() => (tab.filters && "type" in tab.filters ? tab.filters : {}));

  // apiSorting drives the backend ORDER BY — it mirrors the TanStack sort state
  const [apiSorting, setApiSorting] = React.useState<SortingState>(
    () => tab.sorting ?? [],
  );

  // Persist sort state to the tab record (debounced) and re-fetch from page 1
  const saveSortingTimerRef = React.useRef<ReturnType<
    typeof setTimeout
  > | null>(null);
  const handleSortingChange = React.useCallback(
    (updater: SortingState | ((prev: SortingState) => SortingState)) => {
      const next =
        typeof updater === "function" ? updater(apiSorting) : updater;
      setApiSorting(next);
      setPage(1);
      setData([]);
      prevDataRef.current = [];
      if (saveSortingTimerRef.current)
        clearTimeout(saveSortingTimerRef.current);
      saveSortingTimerRef.current = setTimeout(() => {
        void updateTabSorting(tab.id, next);
      }, 600);
    },
    [tab.id, apiSorting],
  );

  // Reset page and filters when the active tab changes
  // (sorting + full state reset handled by key={tab.id} remount in the parent)
  const prevTabIdRef = React.useRef(tab.id);
  if (prevTabIdRef.current !== tab.id) {
    prevTabIdRef.current = tab.id;
    setPage(1);
    setApiFilters(tab.filters && "type" in tab.filters ? tab.filters : {});
  }

  // Reset to page 1 whenever filters change
  const handleFilterChange = React.useCallback(
    (filters: FilterGroup | Record<string, never>) => {
      setPage(1);
      setApiFilters(filters);
    },
    [],
  );

  const initialFilterNode = React.useMemo(
    () =>
      tab.filters && "type" in tab.filters
        ? fromApiFilter(tab.filters as FilterGroup)
        : undefined,
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [tab.id],
  );

  const filterFields = React.useMemo<FilterField[]>(
    () =>
      formFields.map((f) => ({
        label: f.label,
        value: f.field_key,
        fieldType: f.form_field_type,
        options: f.options?.map((o) => ({
          label: String(o.label),
          value: String(o.value),
        })),
      })),
    [formFields],
  );

  const { data: apiData } = useQuery({
    queryKey: [model, tab.id, page, apiFilters, apiSorting],
    queryFn: () =>
      fetchModelIndex(model, {
        tab_id: tab.id,
        filters: apiFilters,
        sorting: apiSorting.map((s) => ({ field_key: s.id, desc: s.desc })),
        search_term: tab.search_term ?? "",
        columns: tab.columns,
        page,
        size: PAGE_SIZE,
        search: "",
      }),
    placeholderData: (prev) => prev,
    staleTime: 30_000,
  });

  const totalRowCount = apiData?.meta?.totalRowCount ?? 0;
  const totalPages = Math.max(1, Math.ceil(totalRowCount / PAGE_SIZE));

  // Prefetch next page so pagination feels instant
  React.useEffect(() => {
    if (page >= totalPages) return;
    void queryClient.prefetchQuery({
      queryKey: [model, tab.id, page + 1, apiFilters, apiSorting],
      queryFn: () =>
        fetchModelIndex(model, {
          tab_id: tab.id,
          filters: apiFilters,
          sorting: apiSorting.map((s) => ({ field_key: s.id, desc: s.desc })),
          search_term: tab.search_term ?? "",
          columns: tab.columns,
          page: page + 1,
          size: PAGE_SIZE,
          search: "",
        }),
      staleTime: 30_000,
    });
  }, [model, tab, page, apiFilters, totalPages, queryClient]);

  useIsomorphicLayoutEffect(() => {
    if (apiData?.data) {
      const rows = (apiData.data as Row[]).map((r) =>
        normalizeRow(r, formFields),
      );
      setData(rows);
      prevDataRef.current = rows;
    }
  }, [apiData?.data, formFields]);

  const fieldTypeMap = React.useMemo(
    () => new Map(formFields.map((f) => [f.field_key, f.form_field_type])),
    [formFields],
  );

  const handleDataChange = React.useCallback(
    async (newData: Row[]) => {
      const prevData = prevDataRef.current;

      // Optimistic update: show new values immediately
      prevDataRef.current = newData;
      setData(newData);

      const patches: Array<{
        rowId: number | string;
        changed: Record<string, unknown>;
      }> = [];

      for (let i = 0; i < newData.length; i++) {
        const newRow = newData[i];
        const prevRow = prevData[i];
        if (!prevRow || !newRow) continue;

        const rowId = newRow.id as number | string;
        if (rowId == null) continue;

        const changed: Record<string, unknown> = {};
        for (const key of Object.keys(newRow)) {
          if (key === "id" || newRow[key] === prevRow[key]) continue;

          const fieldType = fieldTypeMap.get(key);

          // Many-to-many associations can't be patched via a simple column update
          if (
            fieldType === "multi_relation" ||
            fieldType === "multi_select" ||
            fieldType === "creatable_multi_select"
          ) {
            continue;
          }

          // single_relation stores a foreign key column named field_key + "_id"
          if (fieldType === "single_relation") {
            const rawId = newRow[key];
            changed[`${key}_id`] =
              rawId !== null && rawId !== "" ? Number(rawId) : null;
            continue;
          }

          changed[key] = newRow[key];
        }

        if (Object.keys(changed).length > 0) {
          patches.push({ rowId, changed });
        }
      }

      try {
        await Promise.all(
          patches.map(({ rowId, changed }) =>
            updateModelRow(model, rowId, changed),
          ),
        );
      } catch {
        // Rollback optimistic update on failure
        prevDataRef.current = prevData;
        setData(prevData);
        toast.error("Failed to save changes");
      }
    },
    [model, fieldTypeMap],
  );

  const filterFn = React.useMemo(() => getFilterFn<Row>(), []);

  const columnVisibility = React.useMemo(
    () =>
      tab.columns.reduce<Record<string, boolean>>((acc, col) => {
        acc[col.field_key] = col.visible;
        return acc;
      }, {}),
    [tab.columns],
  );

  const sortedFields = React.useMemo(() => {
    const colMap = new Map(tab.columns.map((c) => [c.field_key, c]));
    return [...formFields].sort((a, b) => {
      const aOrder = colMap.get(a.field_key)?.order ?? a.table_order;
      const bOrder = colMap.get(b.field_key)?.order ?? b.table_order;
      return aOrder - bOrder;
    });
  }, [formFields, tab.columns]);

  const columns = React.useMemo<ColumnDef<Row>[]>(() => {
    const colMap = new Map(tab.columns.map((c) => [c.field_key, c]));
    return [
      getDataGridSelectColumn<Row>({ enableRowMarkers: true }),
      ...sortedFields.map(
        (field): ColumnDef<Row> => ({
          id: field.field_key,
          accessorKey: field.field_key,
          header: field.label,
          minSize: colMap.get(field.field_key)?.width ?? 180,
          filterFn,
          meta: {
            label: field.label,
            cell: fieldTypeToCellOpts(field),
          },
        }),
      ),
    ];
  }, [sortedFields, tab.columns, filterFn]);

  const onFilesUpload = React.useCallback(
    async ({
      files,
    }: {
      files: File[];
      rowIndex: number;
      columnId: string;
    }) => {
      return files.map((file) => ({
        id: crypto.randomUUID(),
        name: file.name,
        size: file.size,
        type: file.type,
        url: URL.createObjectURL(file),
      }));
    },
    [],
  );

  const getRowId = React.useCallback((row: Row) => String(row.id ?? ""), []);

  const columnOrder = React.useMemo(
    () => ["select", ...sortedFields.map((f) => f.field_key)],
    [sortedFields],
  );

  const handleSaveColumns = React.useCallback(
    async (
      cols: import("@/components/data-grid/data-grid-view-menu").ColumnSavePayload[],
    ) => {
      try {
        await updateTabColumns(
          tab.id,
          cols.map((c) => ({
            field_key: c.id,
            visible: c.visible,
            order: c.order,
          })),
        );
        void queryClient.invalidateQueries({
          queryKey: ["tabs", tab.model_name],
        });
      } catch {
        toast.error("Failed to save column settings");
      }
    },
    [tab.id, tab.model_name, queryClient],
  );

  const { table, ...dataGridProps } = useDataGrid({
    data,
    columns,
    onDataChange: handleDataChange,
    onFilesUpload,
    getRowId,
    initialState: {
      columnPinning: { left: ["select"] },
      columnVisibility,
      columnOrder,
      sorting: tab.sorting ?? [],
    },
    manualSorting: true,
    onSortingChange: handleSortingChange,
    enableSearch: true,
    enablePaste: true,
  });

  const height = Math.max(400, windowSize.height - 140);

  return (
    <DirectionProvider dir="ltr">
      <div className="flex flex-col">
        <div
          role="toolbar"
          aria-orientation="horizontal"
          className="flex h-9 items-center gap-2 w-full self-end"
        >
          <AdvancedFilter
            fields={filterFields}
            value={initialFilterNode}
            onChange={handleFilterChange}
            align="end"
          />
          <DataGridSortMenu table={table} align="end" />
          <DataGridRowHeightMenu table={table} align="end" />
          <DataGridViewMenu
            table={table}
            align="end"
            onSaveColumns={handleSaveColumns}
          />
        </div>
        <DataGridKeyboardShortcuts enableSearch={!!dataGridProps.searchState} />
        <DataGrid {...dataGridProps} table={table} height={height} />
        <div className="flex items-center justify-between px-1 text-muted-foreground  text-sm">
          <span>
            {totalRowCount > 0
              ? `${(page - 1) * PAGE_SIZE + 1}–${Math.min(page * PAGE_SIZE, totalRowCount)} of ${totalRowCount} rows`
              : "0 rows"}
          </span>
          <div className="flex items-center  gap-2">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => setPage((p) => p - 1)}
              disabled={page <= 1}
            >
              Previous
            </Button>
            <span className="tabular-nums">
              Page {page} of {totalPages}
            </span>
            <Button
              variant="ghost"
              size="sm"
              onClick={() => setPage((p) => p + 1)}
              disabled={page >= totalPages}
            >
              Next
            </Button>
          </div>
        </div>
      </div>
    </DirectionProvider>
  );
}
