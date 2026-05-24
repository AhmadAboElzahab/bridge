"use client";

import Cookies from "js-cookie";
import { create } from "zustand";
import { persist } from "zustand/middleware";
import { TOKEN_KEY } from "@/lib/api";

interface AuthState {
  token: string | null;
  setToken: (token: string) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      token: null,
      setToken: (token) => {
        Cookies.set(TOKEN_KEY, token, { expires: 7, sameSite: "lax" });
        set({ token });
      },
      logout: () => {
        Cookies.remove(TOKEN_KEY);
        set({ token: null });
      },
    }),
    {
      name: "bridge-auth",
      partialize: (state) => ({ token: state.token }),
    },
  ),
);
