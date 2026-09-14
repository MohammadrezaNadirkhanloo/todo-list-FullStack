import { redirect, type LoaderFunctionArgs } from "react-router";
import { useSessionStore } from "@/zustand/permission.store";
import { showToast } from "@/utils/toast";

export async function protectedLoader({ request }: LoaderFunctionArgs) {
  const token = localStorage.getItem("token");
  const lang = localStorage.getItem("Language") || "fa";
  const isFa = lang === "fa";
  
  // استخراج مسیر فعلی برای برگشت بعد از لاگین
  const url = new URL(request.url);
  const redirectTo = url.pathname + url.search;

  if (!token) {
    showToast.error({
      message: isFa ? "نشست منقضی شده" : "Session expired",
      description: isFa ? "لطفاً دوباره وارد شوید" : "Please login again",
    });
    return redirect(`/auth?redirect=${encodeURIComponent(redirectTo)}`);
  }

  await useSessionStore.getState().ensureSession();
  
  if (!useSessionStore.getState().isAuthenticated) {
    useSessionStore.getState().clearSession();
    showToast.error({
      message: isFa ? "توکن نامعتبر است" : "Invalid Token",
      description: isFa ? "لطفاً دوباره وارد شوید" : "Please login again",
    });
    return redirect(`/auth?redirect=${encodeURIComponent(redirectTo)}`);
  }

  return null;
}