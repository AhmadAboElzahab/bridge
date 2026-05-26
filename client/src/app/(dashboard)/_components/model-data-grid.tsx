"use client";

import { DirectionProvider } from "@radix-ui/react-direction";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import type { ColumnDef, SortingState } from "@tanstack/react-table";
import { MoveDiagonal, Search, X } from "lucide-react";
import { parseAsString, useQueryState } from "nuqs";
import * as React from "react";
import { AdvancedFilter } from "@/components/advanced-filter";
import { fromApiFilter } from "@/components/advanced-filter/to-api-filter";
import type { FilterField } from "@/components/advanced-filter/types";
import { DataGrid } from "@/components/data-grid/data-grid";
import { DataGridKeyboardShortcuts } from "@/components/data-grid/data-grid-keyboard-shortcuts";
import { DataGridRowHeightMenu } from "@/components/data-grid/data-grid-row-height-menu";
import { getDataGridSelectColumn } from "@/components/data-grid/data-grid-select-column";
import { DataGridSortMenu } from "@/components/data-grid/data-grid-sort-menu";
import { DataGridViewMenu } from "@/components/data-grid/data-grid-view-menu";
import { useDataGrid } from "@/hooks/use-data-grid";
import { useIsomorphicLayoutEffect } from "@/hooks/use-isomorphic-layout-effect";
import { useWindowSize } from "@/hooks/use-window-size";
import { getFilterFn } from "@/lib/data-grid-filters";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";
import {
  fetchModelIndex,
  fetchModelRow,
  PAGE_SIZE,
  updateModelRelation,
  updateModelRow,
} from "@/services/data.service";
import { uploadFile } from "@/services/upload.service";
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
      if (Array.isArray(val)) return val.map(toSelectId);
      if (typeof val === "string" && val.startsWith("[")) {
        try {
          const parsed = JSON.parse(val) as unknown[];
          return parsed.map(toSelectId);
        } catch {
          return [];
        }
      }
      return [];

    case "url_field":
    case "social_url_field":
    case "link_field":
      // "url" cell variant expects a plain string, not FileCellData[]
      if (typeof val === "string") return val;
      if (typeof val === "object" && val !== null) {
        const o = val as Record<string, unknown>;
        return String(o.url ?? o.value ?? "");
      }
      return "";

    case "image_field":
    case "file_field":
    case "video_field":
    case "attachments_field":
      return toFileCellData(val);

    case "array_field":
      return Array.isArray(val) ? val.map(toSelectId) : [];

    case "password_field":
      return null;

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

    case "string_field":
    case "phone_field":
    case "email_field":
    case "password_field":
    case "formula_field":
      return { variant: "short-text" };

    case "array_field":
      return { variant: "multi-select", options };

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

function ProfileFieldEditor({
  field,
  value,
  onChange,
}: {
  field: FormField;
  value: unknown;
  onChange: (val: unknown) => void;
}) {
  const type = field.form_field_type;

  if (type === "formula_field") {
    return (
      <div className="flex flex-col gap-1.5">
        <Label className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
          {field.label}
        </Label>
        <span className="text-sm text-muted-foreground">
          {value !== null && value !== undefined && value !== "" ? String(value) : "—"}
        </span>
      </div>
    );
  }

  if (type === "password_field") return null;

  if (type === "boolean_field") {
    return (
      <div className="flex flex-col gap-1.5">
        <Label className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
          {field.label}
        </Label>
        <div className="flex items-center gap-2">
          <Checkbox
            checked={!!value}
            onCheckedChange={(checked) => onChange(!!checked)}
          />
          <span className="text-sm">{value ? "Yes" : "No"}</span>
        </div>
      </div>
    );
  }

  if (
    type === "single_select" ||
    type === "radio_select" ||
    type === "creatable_single_select" ||
    type === "single_relation"
  ) {
    const options = field.options ?? [];
    return (
      <div className="flex flex-col gap-1.5">
        <Label className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
          {field.label}
        </Label>
        <Select value={String(value ?? "")} onValueChange={onChange}>
          <SelectTrigger className="w-full">
            <SelectValue placeholder="Select…" />
          </SelectTrigger>
          <SelectContent>
            {options.map((opt) => (
              <SelectItem key={String(opt.value)} value={String(opt.value)}>
                {String(opt.label)}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>
    );
  }

  if (
    type === "multi_select" ||
    type === "creatable_multi_select" ||
    type === "multi_relation" ||
    type === "array_field"
  ) {
    const arr = Array.isArray(value) ? (value as string[]) : [];
    const options = field.options ?? [];
    const available = options.filter((o) => !arr.includes(String(o.value)));
    return (
      <div className="flex flex-col gap-1.5">
        <Label className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
          {field.label}
        </Label>
        {arr.length > 0 && (
          <div className="flex flex-wrap gap-1">
            {arr.map((v, i) => {
              const opt = options.find((o) => String(o.value) === String(v));
              return (
                <span
                  key={i}
                  className="inline-flex items-center gap-1 rounded-full bg-muted px-2 py-0.5 text-xs"
                >
                  {opt ? String(opt.label) : String(v)}
                  <button
                    type="button"
                    onClick={() => onChange(arr.filter((_, j) => j !== i))}
                    className="hover:text-destructive leading-none"
                  >
                    ×
                  </button>
                </span>
              );
            })}
          </div>
        )}
        {available.length > 0 && (
          <Select
            value=""
            onValueChange={(v) => {
              if (v && !arr.includes(v)) onChange([...arr, v]);
            }}
          >
            <SelectTrigger className="w-full">
              <SelectValue placeholder="Add…" />
            </SelectTrigger>
            <SelectContent>
              {available.map((opt) => (
                <SelectItem key={String(opt.value)} value={String(opt.value)}>
                  {String(opt.label)}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        )}
      </div>
    );
  }

  if (type === "rich_field") {
    return (
      <div className="flex flex-col gap-1.5">
        <Label className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
          {field.label}
        </Label>
        <Textarea
          value={String(value ?? "")}
          onChange={(e) => onChange(e.target.value)}
          rows={5}
        />
      </div>
    );
  }

  if (
    type === "integer_field" ||
    type === "float_field" ||
    type === "currency_field" ||
    type === "rating_field"
  ) {
    return (
      <div className="flex flex-col gap-1.5">
        <Label className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
          {field.label}
        </Label>
        <Input
          type="number"
          value={value !== null && value !== undefined ? String(value) : ""}
          onChange={(e) =>
            onChange(e.target.value === "" ? null : Number(e.target.value))
          }
        />
      </div>
    );
  }

  if (type === "date_field") {
    const dateStr = value ? String(value).split("T")[0] : "";
    return (
      <div className="flex flex-col gap-1.5">
        <Label className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
          {field.label}
        </Label>
        <Input
          type="date"
          value={dateStr}
          onChange={(e) => onChange(e.target.value || null)}
        />
      </div>
    );
  }

  if (type === "datetime_field") {
    const raw = value ? String(value) : "";
    const dtStr = raw.split("+")[0]?.split("Z")[0]?.split(".")[0] ?? "";
    return (
      <div className="flex flex-col gap-1.5">
        <Label className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
          {field.label}
        </Label>
        <Input
          type="datetime-local"
          value={dtStr}
          onChange={(e) => onChange(e.target.value || null)}
        />
      </div>
    );
  }

  if (
    type === "image_field" ||
    type === "video_field" ||
    type === "file_field" ||
    type === "attachments_field"
  ) {
    const files = Array.isArray(value)
      ? (value as { url?: string; name?: string; type?: string }[])
      : [];
    const isMulti = type === "attachments_field";
    return (
      <div className="flex flex-col gap-1.5">
        <Label className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
          {field.label}
        </Label>
        {files.length > 0 && (
          <div className="flex flex-wrap gap-2">
            {files.map((file, i) =>
              file.type?.startsWith("image/") ? (
                <div key={i} className="relative group/img">
                  <img
                    src={file.url}
                    alt={file.name ?? ""}
                    className="h-20 w-20 rounded-md object-cover border"
                  />
                  <button
                    type="button"
                    onClick={() => onChange(files.filter((_, j) => j !== i))}
                    className="absolute -top-1 -right-1 size-4 rounded-full bg-destructive text-white text-xs flex items-center justify-center opacity-0 group-hover/img:opacity-100 transition-opacity"
                  >
                    ×
                  </button>
                </div>
              ) : (
                <div key={i} className="flex items-center gap-1">
                  <a
                    href={file.url}
                    target="_blank"
                    rel="noreferrer"
                    className="text-sm text-blue-600 hover:underline"
                  >
                    {file.name}
                  </a>
                  <button
                    type="button"
                    onClick={() => onChange(files.filter((_, j) => j !== i))}
                    className="text-xs text-muted-foreground hover:text-destructive"
                  >
                    ×
                  </button>
                </div>
              ),
            )}
          </div>
        )}
        <Input
          type="url"
          placeholder="Paste file URL…"
          onBlur={(e) => {
            const url = e.target.value.trim();
            if (!url) return;
            const newFile = {
              id: url,
              name: url.split("/").pop() ?? url,
              size: 0,
              type: guessMediaType(url),
              url,
            };
            onChange(isMulti ? [...files, newFile] : [newFile]);
            e.target.value = "";
          }}
        />
      </div>
    );
  }

  const inputType =
    type === "email_field"
      ? "email"
      : type === "phone_field"
        ? "tel"
        : type === "url_field" ||
            type === "social_url_field" ||
            type === "link_field"
          ? "url"
          : "text";

  return (
    <div className="flex flex-col gap-1.5">
      <Label className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
        {field.label}
      </Label>
      <Input
        type={inputType}
        value={String(value ?? "")}
        onChange={(e) => onChange(e.target.value)}
      />
    </div>
  );
}

function RowProfileModal({
  profileId,
  model,
  formFields,
  fieldTypeMap,
  onClose,
  onSaved,
  onRowUpdated,
}: {
  profileId: string | null;
  model: string;
  formFields: FormField[];
  fieldTypeMap: Map<string, FieldType>;
  onClose: () => void;
  onSaved: () => void;
  onRowUpdated: (id: string, data: Row) => void;
}) {
  const { data: rawRow, isLoading } = useQuery({
    queryKey: ["profile", model, profileId],
    queryFn: () => fetchModelRow(model, profileId!),
    enabled: !!profileId,
    staleTime: 0,
  });

  const [localData, setLocalData] = React.useState<Row | null>(null);
  const [isSaving, setIsSaving] = React.useState(false);

  // Reset when switching to a different row
  React.useEffect(() => {
    setLocalData(null);
  }, [profileId]);

  // Populate once the fetch resolves
  React.useEffect(() => {
    if (rawRow) setLocalData(normalizeRow(rawRow as Row, formFields));
  }, [rawRow, formFields]);

  const nameField = formFields.find((f) =>
    ["name", "full_name", "first_name", "title"].includes(f.field_key),
  );
  const title =
    localData && nameField
      ? String(localData[nameField.field_key] ?? "") || "Profile"
      : "Profile";

  const handleSave = async () => {
    if (!localData || !profileId) return;
    setIsSaving(true);
    try {
      const patch: Record<string, unknown> = {};
      const relations: { fieldKey: string; ids: number[] }[] = [];

      for (const key of Object.keys(localData)) {
        if (key === "id") continue;
        const val = localData[key];
        const fieldType = fieldTypeMap.get(key);

        if (fieldType === "formula_field" || fieldType === "password_field")
          continue;

        if (fieldType === "multi_relation") {
          const ids = (Array.isArray(val) ? val : [])
            .map((v) => Number(v))
            .filter((n) => Number.isFinite(n) && n > 0);
          relations.push({ fieldKey: key, ids });
          continue;
        }

        if (
          fieldType === "multi_select" ||
          fieldType === "creatable_multi_select"
        ) {
          patch[key] = JSON.stringify(Array.isArray(val) ? val : []);
          continue;
        }

        if (fieldType === "single_relation") {
          patch[`${key}_id`] =
            val !== null && val !== "" ? Number(val) : null;
          continue;
        }

        if (
          fieldType === "image_field" ||
          fieldType === "video_field" ||
          fieldType === "file_field"
        ) {
          const files = Array.isArray(val) ? (val as { url?: string }[]) : [];
          patch[key] = files[0]?.url ?? null;
          continue;
        }

        if (fieldType === "attachments_field") {
          const files = Array.isArray(val) ? (val as { url?: string }[]) : [];
          patch[key] = JSON.stringify(files.map((f) => f.url).filter(Boolean));
          continue;
        }

        patch[key] = val;
      }

      await Promise.all([
        updateModelRow(model, profileId, patch),
        ...relations.map(({ fieldKey, ids }) =>
          updateModelRelation(model, profileId, fieldKey, ids),
        ),
      ]);

      toast.success("Saved");
      onRowUpdated(profileId, localData);
      onSaved();
    } catch {
      toast.error("Failed to save");
    } finally {
      setIsSaving(false);
    }
  };

  const visibleFields = formFields.filter(
    (f) => f.form_field_type !== "password_field",
  );

  return (
    <Dialog open={!!profileId} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="max-w-4xl w-full max-h-[90vh] flex flex-col gap-0 p-0 overflow-hidden">
        <DialogHeader className="px-6 py-4 border-b shrink-0">
          <DialogTitle>{isLoading ? "Loading…" : title}</DialogTitle>
        </DialogHeader>
        <div className="overflow-y-auto flex-1 px-6 py-6">
          {isLoading || !localData ? (
            <div className="flex items-center justify-center py-12 text-sm text-muted-foreground">
              Loading…
            </div>
          ) : (
            <div key={profileId ?? ""} className="flex flex-col gap-5">
              {visibleFields.map((field) => (
                <ProfileFieldEditor
                  key={field.field_key}
                  field={field}
                  value={localData[field.field_key]}
                  onChange={(val) =>
                    setLocalData((prev) =>
                      prev ? { ...prev, [field.field_key]: val } : null,
                    )
                  }
                />
              ))}
            </div>
          )}
        </div>
        <DialogFooter className="px-6 py-4 border-t shrink-0">
          <Button variant="outline" onClick={onClose}>
            Cancel
          </Button>
          <Button onClick={handleSave} disabled={isSaving || isLoading || !localData}>
            {isSaving ? "Saving…" : "Save changes"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
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

  const [profileId, setProfileId] = useQueryState("profile", parseAsString);
  const [localProfileId, setLocalProfileId] = React.useState<string | null>(
    () => profileId,
  );

  // Sync URL → local state (handles browser back/forward)
  const prevProfileIdRef = React.useRef(profileId);
  if (prevProfileIdRef.current !== profileId) {
    prevProfileIdRef.current = profileId;
    setLocalProfileId(profileId);
  }

  const [searchInputValue, setSearchInputValue] = React.useState(
    tab.search_term ?? "",
  );
  const [serverSearch, setServerSearch] = React.useState(
    tab.search_term ?? "",
  );
  const serverSearchTimerRef = React.useRef<ReturnType<typeof setTimeout> | null>(null);

  const handleSearchInputChange = React.useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const value = e.target.value;
      setSearchInputValue(value);
      if (serverSearchTimerRef.current) clearTimeout(serverSearchTimerRef.current);
      serverSearchTimerRef.current = setTimeout(() => {
        setServerSearch(value.trim());
      }, 400);
    },
    [],
  );

  const handleSearchClear = React.useCallback(() => {
    setSearchInputValue("");
    setServerSearch("");
  }, []);

  const [page, setPage] = React.useState(1);

  const [apiFilters, setApiFilters] = React.useState<
    FilterGroup | Record<string, never>
  >(() => (tab.filters && "type" in tab.filters ? tab.filters : {}));

  // apiSorting drives the backend ORDER BY — it mirrors the TanStack sort state
  const [apiSorting, setApiSorting] = React.useState<SortingState>(
    () => tab.sorting ?? [],
  );

  // Persist sort state to the tab record (debounced); data stays visible during re-fetch
  const saveSortingTimerRef = React.useRef<ReturnType<typeof setTimeout> | null>(null);
  const handleSortingChange = React.useCallback(
    (updater: SortingState | ((prev: SortingState) => SortingState)) => {
      const next =
        typeof updater === "function" ? updater(apiSorting) : updater;
      setApiSorting(next);
      setPage(1);
      if (saveSortingTimerRef.current) clearTimeout(saveSortingTimerRef.current);
      saveSortingTimerRef.current = setTimeout(() => {
        void updateTabSorting(tab.id, next);
      }, 600);
    },
    [tab.id, apiSorting],
  );

  // Reset filters when the active tab changes
  const prevTabIdRef = React.useRef(tab.id);
  if (prevTabIdRef.current !== tab.id) {
    prevTabIdRef.current = tab.id;
    setPage(1);
    setApiFilters(tab.filters && "type" in tab.filters ? tab.filters : {});
  }

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
    queryKey: [model, tab.id, page, apiFilters, apiSorting, serverSearch],
    queryFn: () =>
      fetchModelIndex(model, {
        tab_id: tab.id,
        filters: apiFilters,
        sorting: apiSorting.map((s) => ({ field_key: s.id, desc: s.desc })),
        search_term: serverSearch,
        columns: tab.columns,
        page,
        size: PAGE_SIZE,
        search: serverSearch,
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
      queryKey: [model, tab.id, page + 1, apiFilters, apiSorting, serverSearch],
      queryFn: () =>
        fetchModelIndex(model, {
          tab_id: tab.id,
          filters: apiFilters,
          sorting: apiSorting.map((s) => ({ field_key: s.id, desc: s.desc })),
          search_term: serverSearch,
          columns: tab.columns,
          page: page + 1,
          size: PAGE_SIZE,
          search: serverSearch,
        }),
      staleTime: 30_000,
    });
  }, [model, tab, page, apiFilters, apiSorting, serverSearch, totalPages, queryClient]);

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

  // Pending server writes: rowId → accumulated patch + relations
  const pendingWritesRef = React.useRef<
    Map<string | number, { patch: Record<string, unknown>; relations: { fieldKey: string; ids: number[] }[] }>
  >(new Map());
  const writeFlushTimerRef = React.useRef<ReturnType<typeof setTimeout> | null>(null);

  const flushWrites = React.useCallback(async () => {
    const pending = pendingWritesRef.current;
    if (pending.size === 0) return;
    pendingWritesRef.current = new Map();

    const calls: Promise<void>[] = [];
    for (const [rowId, { patch, relations }] of pending) {
      if (Object.keys(patch).length > 0) {
        calls.push(updateModelRow(model, rowId, patch));
      }
      for (const { fieldKey, ids } of relations) {
        calls.push(updateModelRelation(model, rowId, fieldKey, ids));
      }
    }

    try {
      await Promise.all(calls);
    } catch {
      toast.error("Failed to save changes");
      void queryClient.invalidateQueries({ queryKey: [model] });
    }
  }, [model, queryClient]);

  // Flush on unmount so no pending writes are lost
  React.useEffect(() => {
    return () => {
      if (writeFlushTimerRef.current) clearTimeout(writeFlushTimerRef.current);
      void flushWrites();
    };
  }, [flushWrites]);

  const handleDataChange = React.useCallback(
    (newData: Row[]) => {
      const prevData = prevDataRef.current;

      // Optimistic update: show new values immediately
      prevDataRef.current = newData;
      setData(newData);

      const patches: Array<{
        rowId: number | string;
        changed: Record<string, unknown>;
      }> = [];
      const relationPatches: Array<{
        rowId: number | string;
        fieldKey: string;
        ids: number[];
      }> = [];

      for (let i = 0; i < newData.length; i++) {
        const newRow = newData[i];
        const prevRow = prevData[i];
        if (!prevRow || !newRow) continue;

        const rowId = newRow.id as number | string;
        if (rowId == null) continue;

        const changed: Record<string, unknown> = {};
        for (const key of Object.keys(newRow)) {
          if (key === "id") continue;

          const newVal = newRow[key];
          const prevVal = prevRow[key];

          // Array-aware equality: same reference OR same JSON for arrays
          const unchanged =
            newVal === prevVal ||
            (Array.isArray(newVal) &&
              Array.isArray(prevVal) &&
              JSON.stringify(newVal) === JSON.stringify(prevVal));
          if (unchanged) continue;

          const fieldType = fieldTypeMap.get(key);

          // computed fields are server-side only — never write back
          if (fieldType === "formula_field") continue;

          // many2many: queue a separate relation-replace call
          if (fieldType === "multi_relation") {
            const ids = (Array.isArray(newVal) ? newVal : [])
              .map((v) => Number(v))
              .filter((n) => Number.isFinite(n) && n > 0);
            relationPatches.push({ rowId, fieldKey: key, ids });
            continue;
          }

          // single file fields — store the URL of the uploaded file
          if (
            fieldType === "image_field" ||
            fieldType === "video_field" ||
            fieldType === "file_field"
          ) {
            const files = Array.isArray(newVal)
              ? (newVal as { url?: string }[])
              : [];
            changed[key] = files[0]?.url ?? null;
            continue;
          }

          // multi-file field — store JSON array of URLs
          if (fieldType === "attachments_field") {
            const files = Array.isArray(newVal)
              ? (newVal as { url?: string }[])
              : [];
            changed[key] = JSON.stringify(
              files.map((f) => f.url).filter(Boolean),
            );
            continue;
          }

          // multi_select / creatable_multi_select stored as a JSON string column
          if (
            fieldType === "multi_select" ||
            fieldType === "creatable_multi_select"
          ) {
            changed[key] = JSON.stringify(
              Array.isArray(newVal) ? newVal : [],
            );
            continue;
          }

          // single_relation maps to a foreign key column named field_key + "_id"
          if (fieldType === "single_relation") {
            const numId = Number(newVal);
            // 0 is Go's uint zero value — not a valid FK; treat as null (clear relation)
            changed[`${key}_id`] =
              Number.isFinite(numId) && numId > 0 ? numId : null;
            continue;
          }

          // url/link stored as a plain string
          if (
            fieldType === "url_field" ||
            fieldType === "social_url_field" ||
            fieldType === "link_field"
          ) {
            changed[key] =
              typeof newVal === "string" ? newVal : null;
            continue;
          }

          changed[key] = newVal;
        }

        if (Object.keys(changed).length > 0) {
          patches.push({ rowId, changed });
        }
      }

      // Queue writes — merge into pending map so rapid edits to the same row collapse
      for (const { rowId, changed } of patches) {
        const existing = pendingWritesRef.current.get(rowId) ?? { patch: {}, relations: [] };
        pendingWritesRef.current.set(rowId, {
          patch: { ...existing.patch, ...changed },
          relations: existing.relations,
        });
      }
      for (const { rowId, fieldKey, ids } of relationPatches) {
        const existing = pendingWritesRef.current.get(rowId) ?? { patch: {}, relations: [] };
        const otherRelations = existing.relations.filter((r) => r.fieldKey !== fieldKey);
        pendingWritesRef.current.set(rowId, {
          patch: existing.patch,
          relations: [...otherRelations, { fieldKey, ids }],
        });
      }

      if (writeFlushTimerRef.current) clearTimeout(writeFlushTimerRef.current);
      writeFlushTimerRef.current = setTimeout(() => { void flushWrites(); }, 500);
    },
    [fieldTypeMap, flushWrites],
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

  const expandColumn = React.useMemo<ColumnDef<Row>>(
    () => ({
      id: "expand",
      header: () => null,
      cell: ({ row }) => (
        <div className="group flex items-center justify-center h-full">
          <button
            type="button"
            onClick={(e) => {
              e.stopPropagation();
              const id = String(row.original.id);
              setLocalProfileId(id);
              void setProfileId(id);
            }}
            className="rounded p-0.5 text-muted-foreground opacity-0 group-hover:opacity-100 hover:text-foreground hover:bg-accent transition-all"
            aria-label="Open profile"
          >
            <MoveDiagonal className="size-3.5" />
          </button>
        </div>
      ),
      size: 36,
      enableHiding: false,
      enableResizing: false,
      enableSorting: false,
    }),
    [setLocalProfileId, setProfileId],
  );

  const columns = React.useMemo<ColumnDef<Row>[]>(() => {
    const colMap = new Map(tab.columns.map((c) => [c.field_key, c]));
    return [
      getDataGridSelectColumn<Row>({ enableRowMarkers: true }),
      expandColumn,
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
  }, [expandColumn, sortedFields, tab.columns, filterFn]);

  const onFilesUpload = React.useCallback(
    async ({
      files,
      columnId,
    }: {
      files: File[];
      rowIndex: number;
      columnId: string;
    }) => {
      const folder = `${model}/${columnId}`;
      try {
        const results = await Promise.all(
          files.map(async (file) => {
            const { url } = await uploadFile(file, folder);
            return {
              id: url,
              name: file.name,
              size: file.size,
              type: file.type,
              url,
            };
          }),
        );
        return results;
      } catch {
        toast.error("File upload failed");
        return [];
      }
    },
    [model],
  );

  const getRowId = React.useCallback((row: Row) => String(row.id ?? ""), []);

  const columnOrder = React.useMemo(
    () => ["select", "expand", ...sortedFields.map((f) => f.field_key)],
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
      columnPinning: { left: ["select", "expand"] },
      columnVisibility,
      columnOrder,
      sorting: tab.sorting ?? [],
    },
    manualSorting: true,
    onSortingChange: handleSortingChange,
    enableSearch: true,
    enablePaste: true,
  });

  const originalOnSearchRef = React.useRef(dataGridProps.searchState?.onSearch);
  React.useEffect(() => {
    originalOnSearchRef.current = dataGridProps.searchState?.onSearch;
  }, [dataGridProps.searchState?.onSearch]);

  const handleSearch = React.useCallback((query: string) => {
    originalOnSearchRef.current?.(query);
    if (serverSearchTimerRef.current) clearTimeout(serverSearchTimerRef.current);
    serverSearchTimerRef.current = setTimeout(() => {
      setServerSearch(query.trim());
    }, 400);
  }, []);

  const patchedSearchState = React.useMemo(() => {
    if (!dataGridProps.searchState) return undefined;
    return { ...dataGridProps.searchState, onSearch: handleSearch };
  }, [dataGridProps.searchState, handleSearch]);

  // Start invisible so the very first paint is opacity:0 — useEffect fires after
  // paint and triggers the CSS transition, avoiding any flash of content.
  const [visible, setVisible] = React.useState(false);
  React.useEffect(() => { setVisible(true); }, []);

  const height = Math.max(400, windowSize.height - 140);

  return (
    <>
      <DirectionProvider dir="ltr">
      <div
        className="flex flex-col"
        style={{ opacity: visible ? 1 : 0, transition: "opacity 0.2s ease" }}
      >
        <div
          role="toolbar"
          aria-orientation="horizontal"
          className="flex h-9 items-center gap-2 w-full self-end"
        >
          <div className="relative flex-1 max-w-xs">
            <Search className="absolute left-2.5 top-1/2 -translate-y-1/2 size-3.5 text-muted-foreground pointer-events-none" />
            <Input
              placeholder="Search..."
              value={searchInputValue}
              onChange={handleSearchInputChange}
              className="h-7 pl-8 pr-7 text-sm"
            />
            {searchInputValue && (
              <button
                type="button"
                onClick={handleSearchClear}
                className="absolute right-2 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
              >
                <X className="size-3.5" />
              </button>
            )}
          </div>
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
        <div style={{ height }}>
          <DataGrid {...dataGridProps} searchState={patchedSearchState} table={table} height={height} />
        </div>
        <div className="flex items-center justify-between px-1 text-muted-foreground text-sm">
          <span>
            {totalRowCount > 0
              ? `${(page - 1) * PAGE_SIZE + 1}–${Math.min(page * PAGE_SIZE, totalRowCount)} of ${totalRowCount} rows`
              : "0 rows"}
          </span>
          <div className="flex items-center gap-2">
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
    <RowProfileModal
      profileId={localProfileId}
      model={model}
      formFields={formFields}
      fieldTypeMap={fieldTypeMap}
      onClose={() => {
        setLocalProfileId(null);
        void setProfileId(null);
      }}
      onSaved={() => {
        void queryClient.invalidateQueries({ queryKey: [model] });
      }}
      onRowUpdated={(id, updated) => {
        setData((prev) =>
          prev.map((row) => (String(row.id) === id ? { ...row, ...updated } : row)),
        );
      }}
    />
    </>
  );
}
