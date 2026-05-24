import { nanoid } from "nanoid";

export function generateId({ length = 8 }: { length?: number } = {}): string {
  return nanoid(length);
}
