import { api } from "@/lib/api";
import type { IndexPayload, IndexResponse } from "@/types/api";

export async function updateModelRow(
  model: string,
  id: number | string,
  patch: Record<string, unknown>,
): Promise<void> {
  await api.patch(`/${model}/${id}`, patch);
}

export const PAGE_SIZE = 50;

export async function fetchModelIndex<T = Record<string, unknown>>(
  model: string,
  payload: Omit<IndexPayload, "model">,
): Promise<IndexResponse<T>> {
  const res = await api.post<IndexResponse<T>>(`/${model}/index`, payload);
  return res.data;
}
