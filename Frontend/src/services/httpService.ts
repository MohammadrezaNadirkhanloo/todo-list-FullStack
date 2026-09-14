import { showToast } from "@/utils/toast";
import axios from "axios";

export const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL ??"localhost:5005",
  withCredentials: true,
  timeout: 30_000,
  headers: {
    "Content-Type": "application/json",
    Accept: "application/json",
  },
});

// api.interceptors.request.use((config) => {
//   const locale = document.documentElement.lang;
//   if (locale) config.headers.set("Accept-Language", locale);
//   return config;
// });


api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem("token");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  },
);

// Response Interceptor
api.interceptors.response.use(
  (response) => response,
  (error) => {
    // const lang = localStorage.getItem("Language") || "fa";
    // const isFa = lang === "fa";

    if (!error.response) {
      showToast.error({
        message: "Network Error",
        description: "Please check your internet connection",
      });
      return Promise.reject(error);
    }

    const status = error.response?.status;
    const serverMessage = error.response?.data?.message;

    switch (status) {
      case 400:
        showToast.error({
          message: "Bad Request",
          description: serverMessage || "The submitted data is invalid",
        });
        break;

      case 401:
        // useSessionStore.getState().clearSession();
        // showToast.error({
        //   message: isFa ? "نشست منقضی شده" : "Session expired",
        //   description: isFa ? "لطفاً دوباره وارد شوید" : "Please login again",
        // });
        // const currentPath = window.location.pathname + window.location.search;
        // window.location.href = `/auth?redirect=${encodeURIComponent(currentPath)}`;
        break;

      case 403:
        showToast.error({
          message: "Forbidden",
          description: "You are not allowed to perform this action",
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
        window.location.href = "/500";
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

    return Promise.reject(error);
  },
);
