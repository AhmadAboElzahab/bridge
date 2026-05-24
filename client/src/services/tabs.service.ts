import { api } from "@/lib/api";
import type { TabsResponse } from "@/types/api";

export async function fetchTabs(model: string): Promise<TabsResponse> {
  const res = await api.get<TabsResponse>("/tabs", { params: { model } });
  return res.data;
}

export async function createTab(modelName: string, tabName: string) {
  const res = await api.post("/tabs", {
    model_name: modelName,
    tab_name: tabName,
  });
  return res.data;
}

export async function renameTab(id: number, tabName: string) {
  const res = await api.put(`/tabs/${id}`, { tab_name: tabName });
  return res.data;
}

export async function deleteTab(id: number) {
  const res = await api.delete(`/tabs/${id}`);
  return res.data;
}
