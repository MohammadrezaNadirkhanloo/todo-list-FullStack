import { toast as sonnerToast } from "sonner";
import type { ExternalToast } from "sonner";

type ToastType = "success" | "error" | "warning" | "info";

interface ShowToastOptions extends ExternalToast {
  message: string;
  description?: string;
}

const getPosition = (): "top-right" | "top-left" => {
  const lang = localStorage.getItem("Language");
  return lang === "fa" ? "top-left" : "top-right";
};


const toastStyles: Record<ToastType, React.CSSProperties> = {
  success: {
     fontFamily: "Vazir",
    "--normal-bg":
      "color-mix(in oklab, light-dark(var(--color-green-600), var(--color-green-400)) 10%, var(--background))",
    "--normal-text": "light-dark(var(--color-green-600), var(--color-green-400))",
    "--normal-border": "light-dark(var(--color-green-600), var(--color-green-400))",
  } as React.CSSProperties,

  error: {
     fontFamily: "Vazir",
    "--normal-bg": "color-mix(in oklab, var(--destructive) 10%, var(--background))",
    "--normal-text": "var(--destructive)",
    "--normal-border": "var(--destructive)",
  } as React.CSSProperties,

  warning: {
     fontFamily: "Vazir",
    "--normal-bg":
      "color-mix(in oklab, light-dark(var(--color-amber-600), var(--color-amber-400)) 10%, var(--background))",
    "--normal-text": "light-dark(var(--color-amber-600), var(--color-amber-400))",
    "--normal-border": "light-dark(var(--color-amber-600), var(--color-amber-400))",
  } as React.CSSProperties,

  info: {
     fontFamily: "Vazir",
    "--normal-bg":
      "color-mix(in oklab, light-dark(var(--color-sky-600), var(--color-sky-400)) 10%, var(--background))",
    "--normal-text": "light-dark(var(--color-sky-600), var(--color-sky-400))",
    "--normal-border": "light-dark(var(--color-sky-600), var(--color-sky-400))",
  } as React.CSSProperties,
};

export const showToast = {
  success: ({ message, description, ...options }: ShowToastOptions) => {
    return sonnerToast.success(message, {
      description,
      style: toastStyles.success,
      position: getPosition(), 
      ...options,              
    });
  },

  error: ({ message, description, ...options }: ShowToastOptions) => {
    return sonnerToast.error(message, {
      description,
      style: toastStyles.error,
      position: getPosition(),
      ...options,
    });
  },

  warning: ({ message, description, ...options }: ShowToastOptions) => {
    return sonnerToast.warning(message, {
      description,
      style: toastStyles.warning,
      position: getPosition(),
      ...options,
    });
  },

  info: ({ message, description, ...options }: ShowToastOptions) => {
    return sonnerToast.info(message, {
      description,
      style: toastStyles.info,
      position: getPosition(),
      ...options,
    });
  },

  default: ({ message, description, ...options }: ShowToastOptions) => {
    return sonnerToast(message, {
      description,
      position: getPosition(),
      ...options,
    });
  },
};