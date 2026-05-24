"use client";

import * as TabsPrimitive from "@radix-ui/react-tabs";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { MoreHorizontal, Pencil, Plus, Trash2 } from "lucide-react";
import { parseAsInteger, useQueryState } from "nuqs";
import * as React from "react";
import { toast } from "sonner";
import { DataTableSkeleton } from "@/components/data-table/data-table-skeleton";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { cn } from "@/lib/utils";
import {
  createTab,
  deleteTab,
  fetchTabs,
  renameTab,
} from "@/services/tabs.service";
import type { UserTab } from "@/types/api";
import { UsersTable } from "./users-table";

export function UsersTableContainer() {
  const queryClient = useQueryClient();
  const { data, isLoading, isError } = useQuery({
    queryKey: ["tabs", "User"],
    queryFn: () => fetchTabs("User"),
  });

  const defaultTabId =
    data?.tabs.find((t) => t.is_default)?.id ?? data?.tabs[0]?.id;

  const [activeTabId, setActiveTabId] = useQueryState(
    "tab",
    parseAsInteger.withDefault(defaultTabId ?? 0),
  );

  const [addDialogOpen, setAddDialogOpen] = React.useState(false);
  const [renameDialogOpen, setRenameDialogOpen] = React.useState(false);
  const [tabName, setTabName] = React.useState("");
  const [editingTab, setEditingTab] = React.useState<UserTab | null>(null);

  const createMutation = useMutation({
    mutationFn: (name: string) => createTab("User", name),
    onSuccess: (result) => {
      void queryClient.invalidateQueries({ queryKey: ["tabs", "User"] });
      setAddDialogOpen(false);
      setTabName("");
      if (result?.tab_id) void setActiveTabId(result.tab_id as number);
      toast.success("Tab created");
    },
    onError: () => toast.error("Failed to create tab"),
  });

  const renameMutation = useMutation({
    mutationFn: ({ id, name }: { id: number; name: string }) =>
      renameTab(id, name),
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["tabs", "User"] });
      setRenameDialogOpen(false);
      setTabName("");
      setEditingTab(null);
      toast.success("Tab renamed");
    },
    onError: () => toast.error("Failed to rename tab"),
  });

  const deleteMutation = useMutation({
    mutationFn: (id: number) => deleteTab(id),
    onSuccess: (_, deletedId) => {
      void queryClient.invalidateQueries({ queryKey: ["tabs", "User"] });
      if (activeTabId === deletedId) {
        void setActiveTabId(defaultTabId ?? 0);
      }
      toast.success("Tab deleted");
    },
    onError: () => toast.error("Failed to delete tab"),
  });

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
        <TabsPrimitive.Root
          value={String(activeTab.id)}
          onValueChange={(v) => void setActiveTabId(Number(v))}
        >
          <TabsPrimitive.List className="inline-flex h-8 items-center rounded-lg bg-muted p-1 text-muted-foreground">
            {tabs.map((tab) => (
              <div key={tab.id} className="group flex items-center">
                <TabsPrimitive.Trigger
                  value={String(tab.id)}
                  className={cn(
                    "inline-flex h-6 items-center justify-center whitespace-nowrap rounded-md px-2 text-xs font-medium ring-offset-background transition-all",
                    "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2",
                    "disabled:pointer-events-none disabled:opacity-50",
                    "data-[state=active]:bg-background data-[state=active]:text-foreground data-[state=active]:shadow",
                  )}
                >
                  {tab.tab_name}
                </TabsPrimitive.Trigger>
                {!tab.is_default && (
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <button
                        type="button"
                        className="invisible ml-0.5 inline-flex size-4 items-center justify-center rounded opacity-60 transition-opacity hover:opacity-100 group-hover:visible"
                      >
                        <MoreHorizontal className="size-3" />
                      </button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="start" side="bottom">
                      <DropdownMenuItem
                        onClick={() => {
                          setEditingTab(tab);
                          setTabName(tab.tab_name);
                          setRenameDialogOpen(true);
                        }}
                      >
                        <Pencil className="mr-2 size-3.5" />
                        Rename
                      </DropdownMenuItem>
                      <DropdownMenuItem
                        className="text-destructive focus:text-destructive"
                        onClick={() => deleteMutation.mutate(tab.id)}
                      >
                        <Trash2 className="mr-2 size-3.5" />
                        Delete
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                )}
              </div>
            ))}
          </TabsPrimitive.List>
        </TabsPrimitive.Root>

        <Button
          variant="ghost"
          size="icon"
          className="size-7 shrink-0"
          onClick={() => {
            setTabName("");
            setAddDialogOpen(true);
          }}
        >
          <Plus className="size-3.5" />
          <span className="sr-only">Add tab</span>
        </Button>
      </div>

      <UsersTable formFields={form_fields} tab={activeTab} />

      <Dialog open={addDialogOpen} onOpenChange={setAddDialogOpen}>
        <DialogContent className="sm:max-w-sm">
          <DialogHeader>
            <DialogTitle>New tab</DialogTitle>
          </DialogHeader>
          <div className="grid gap-2">
            <Label htmlFor="add-tab-name">Tab name</Label>
            <Input
              id="add-tab-name"
              value={tabName}
              onChange={(e) => setTabName(e.target.value)}
              placeholder="e.g. All users"
              onKeyDown={(e) => {
                if (e.key === "Enter" && tabName.trim()) {
                  createMutation.mutate(tabName.trim());
                }
              }}
              autoFocus
            />
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setAddDialogOpen(false)}>
              Cancel
            </Button>
            <Button
              disabled={!tabName.trim() || createMutation.isPending}
              onClick={() => createMutation.mutate(tabName.trim())}
            >
              Create
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={renameDialogOpen} onOpenChange={setRenameDialogOpen}>
        <DialogContent className="sm:max-w-sm">
          <DialogHeader>
            <DialogTitle>Rename tab</DialogTitle>
          </DialogHeader>
          <div className="grid gap-2">
            <Label htmlFor="rename-tab-name">Tab name</Label>
            <Input
              id="rename-tab-name"
              value={tabName}
              onChange={(e) => setTabName(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === "Enter" && tabName.trim() && editingTab) {
                  renameMutation.mutate({
                    id: editingTab.id,
                    name: tabName.trim(),
                  });
                }
              }}
              autoFocus
            />
          </div>
          <DialogFooter>
            <Button
              variant="outline"
              onClick={() => {
                setRenameDialogOpen(false);
                setEditingTab(null);
              }}
            >
              Cancel
            </Button>
            <Button
              disabled={!tabName.trim() || renameMutation.isPending}
              onClick={() => {
                if (editingTab) {
                  renameMutation.mutate({
                    id: editingTab.id,
                    name: tabName.trim(),
                  });
                }
              }}
            >
              Rename
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
