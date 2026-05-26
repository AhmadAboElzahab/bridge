import { api } from "@/lib/api";

export interface UploadResult {
  url: string;
  blurhash?: string;
}

export async function uploadFile(
  file: File,
  folder?: string,
): Promise<UploadResult> {
  const form = new FormData();
  form.append("file", file);

  const params = folder ? `?folder=${encodeURIComponent(folder)}` : "";
  const res = await api.post<UploadResult>(`/upload${params}`, form, {
    headers: { "Content-Type": "multipart/form-data" },
  });

  return res.data;
}
