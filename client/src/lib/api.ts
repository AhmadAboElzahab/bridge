import axios from "axios";
import Cookies from "js-cookie";
import { TOKEN_KEY } from "@/lib/constants";

export { TOKEN_KEY };

export const api = axios.create({
  baseURL: "/api",
  headers: { "Content-Type": "application/json" },
});

api.interceptors.request.use((config) => {
  const token = Cookies.get(TOKEN_KEY);
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (res) => res,
  (error) => {
    if (error.response?.status === 401) {
      Cookies.remove(TOKEN_KEY);
      window.location.href = "/login";
    }
    return Promise.reject(error);
  },
);
