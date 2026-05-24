"use client";

import { useQuery } from "@tanstack/react-query";
import { Plus } from "lucide-react";
import { parseAsInteger, useQueryState } from "nuqs";
import { DataTableSkeleton } from "@/components/data-table/data-table-skeleton";
import { Button } from "@/components/ui/button";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { fetchTabs } from "@/services/tabs.service";
import { MaidsTable } from "./maids-table";

export function MaidsTableContainer() {
  const { data, isLoading, isError } = useQuery({
    queryKey: ["tabs", "Maid"],
    queryFn: () => fetchTabs("Maid"),
  });

  const defaultTabId =
    data?.tabs.find((t) => t.is_default)?.id ?? data?.tabs[0]?.id;

  const [activeTabId, setActiveTabId] = useQueryState(
    "tab",
    parseAsInteger.withDefault(defaultTabId ?? 0),
  );

  if (isLoading) {
    return <DataTableSkeleton columnCount={6} rowCount={10} />;
  }

  if (isError || !data) {
    return (
      <div className="flex h-40 items-center justify-center text-muted-foreground text-sm">
        Failed to load tabs. Make sure the backend is running.
      </div>
    );
  }

  const { tabs, form_fields } = data;
  const resolvedTabId = activeTabId || defaultTabId || 0;
  const activeTab = tabs.find((t) => t.id === resolvedTabId) ?? tabs[0];

  if (!activeTab) {
    return (
      <div className="flex h-40 items-center justify-center text-muted-foreground text-sm">
        No tabs found. Create a tab to get started.
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center gap-2">
        <Tabs
          value={String(activeTab.id)}
          onValueChange={(v: string) => void setActiveTabId(Number(v))}
        >
          <TabsList className="h-8">
            {tabs.map((tab) => (
              <TabsTrigger
                key={tab.id}
                value={String(tab.id)}
                className="h-6 px-3 text-xs"
              >
                {tab.tab_name}
              </TabsTrigger>
            ))}
          </TabsList>
        </Tabs>
        <Button
          variant="ghost"
          size="icon"
          className="size-7 shrink-0"
          disabled
        >
          <Plus className="size-3.5" />
          <span className="sr-only">Add tab</span>
        </Button>
      </div>

      <MaidsTable key={activeTab.id} formFields={form_fields} tab={activeTab} />
    </div>
  );
}
