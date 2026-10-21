import { router } from "@/routes";
import { showToast } from "@/utils/toast";
import axios, { AxiosError } from "axios";

declare module "axios" {
  interface AxiosRequestConfig {
    _retry?: boolean;
    _skipRefresh?: boolean;
  }
}

type ApiError = {
  success?: boolean;
  message?: string;
  error?: { code?: string; message?: string };
};

export const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  withCredentials: true,
  timeout: 30_000,
  headers: {
    "Content-Type": "application/json",
    Accept: "application/json",
  },
});

const getCookie = (name: string): string | null => {
  const match = document.cookie.match(
    new RegExp(
      `(?:^|; )${name.replace(/([.$?*|{}()[\]\\/+^])/g, "\\$1")}=([^;]*)`,
    ),
  );
  return match ? decodeURIComponent(match[1]) : null;
};

// ───────────────────────── Request Interceptor ─────────────────────────

api.interceptors.request.use(
  (config) => {
    const csrfToken = getCookie("csrf_token");
    if (csrfToken) {
      config.headers["X-CSRF-Token"] = csrfToken;
    }
    return config;
  },
  (error) => Promise.reject(error),
);

// ───────────────────────── Refresh ─────────────────────────

const NO_REFRESH_PATHS = [
  "/auth/login",
  "/auth/register",
  "/auth/refresh",
  "/auth/logout",
];

const REFRESHABLE_CODES = ["TOKEN_REQUIRED", "TOKEN_EXPIRED", "TOKEN_INVALID"];

const shouldSkipRefresh = (url?: string) =>
  !!url && NO_REFRESH_PATHS.some((p) => url.includes(p));

let refreshPromise: Promise<void> | null = null;

const redirectToSignIn = () => {
  showToast.error({
    message: "Session expired",
    description: "Please login again",
  });
  if (router.state.location.pathname !== "/sign-in") {
    router.navigate("/sign-in");
  }
};

const refreshSession = (): Promise<void> => {
  refreshPromise ??= api
    .post("/auth/refresh", null, { _skipRefresh: true })
    .then(() => undefined)
    .catch((err: AxiosError) => {
      const status = err.response?.status;
      if (status === 401 || status === 403) {
        redirectToSignIn();
      }
      throw err;
    })
    .finally(() => {
      refreshPromise = null;
    });

  return refreshPromise;
};

// ───────────────────────── Error toasts ─────────────────────────

const handleError = (error: AxiosError<ApiError>) => {
  if (!error.response) {
    showToast.error({
      message: "Network Error",
      description: "Please check your internet connection",
    });
    return;
  }

  const status = error.response.status;
  const serverMessage =
    error.response.data?.message || error.response.data?.error?.message;

  switch (status) {
    case 400:
      showToast.error({
        message: "Bad Request",
        description: serverMessage || "The submitted data is invalid",
      });
      break;
    case 401:
      showToast.error({
        message: "Unauthorized",
        description: serverMessage || "Invalid credentials or session",
      });
      break;
    case 403:
      showToast.error({
        message: "Forbidden",
        description:
          serverMessage || "You are not allowed to perform this action",
      });
      break;
    case 404:
      showToast.error({
        message: "Not Found",
        description: serverMessage || "The requested resource was not found",
      });
      break;
    case 422:
      showToast.error({
        message: "Validation Error",
        description: serverMessage || "Please check the fields",
      });
      break;
    case 429:
      showToast.warning({
        message: "Too Many Requests",
        description: "Please wait a moment and try again",
      });
      break;
    case 500:
    case 502:
    case 503:
      showToast.error({
        message: "Server Error",
        description: "A server error occurred. Please try again later",
      });
      break;
    default:
      showToast.error({
        message: "Unknown Error",
        description: serverMessage || "Please try again",
      });
  }
};

// ───────────────────────── Response Interceptor ─────────────────────────

api.interceptors.response.use(
  (response) => response,
  async (error: AxiosError<ApiError>) => {
    const original = error.config;
    const code = error.response?.data?.error?.code;

    const canRefresh =
      error.response?.status === 401 &&
      (!code || REFRESHABLE_CODES.includes(code)) &&
      !!original &&
      !original._retry &&
      !original._skipRefresh &&
      !shouldSkipRefresh(original.url);

    if (canRefresh && original) {
      original._retry = true;

      try {
        await refreshSession();
      } catch (refreshError) {
        const s = (refreshError as AxiosError).response?.status;
        if (s !== 401 && s !== 403) {
          handleError(refreshError as AxiosError<ApiError>);
        }
        return Promise.reject(refreshError);
      }

      return api(original);
    }

    handleError(error);
    return Promise.reject(error);
  },
);
