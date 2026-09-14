import { useSessionStore } from "@/zustand/permission.store";
import { safeRedirectPath } from "@/lib/utils";
import { redirect, type LoaderFunctionArgs } from "react-router";

export async function authLoader({ request }: LoaderFunctionArgs) {
  const token = localStorage.getItem("token");
  if (!token) return null;
  await useSessionStore.getState().ensureSession();
  const { isAuthenticated, homeRoute } = useSessionStore.getState();

  if (isAuthenticated) {
    // اگر با ?redirect=... به /auth آمده و توکن معتبر دارد، به همان مسیر برگردد
    const redirectTo = safeRedirectPath(
      new URL(request.url).searchParams.get("redirect"),
    );
    throw redirect(redirectTo ?? homeRoute ?? "/dashboard");
  }
  return null;
}
