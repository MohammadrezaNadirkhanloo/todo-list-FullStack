import { showToast } from "@/utils/toast";
import type { ExternalToast } from "sonner";

interface ToastOptions extends Omit<ExternalToast, "description"> {
  descriptionKey?: string;
  descriptionValues?: Record<string, any>;
}

export function useToast() {
  return {
    success: (messageKey: string, options?: ToastOptions) => {
      const { ...rest } = options || {};
      return showToast.success({
        message: messageKey,
        // description: descriptionKey
        //   ?  descriptionValues
        //   : undefined,
        ...rest,
      });
    },

    error: (messageKey: string, options?: ToastOptions) => {
      const { ...rest } = options || {};
      return showToast.error({
        message: messageKey,

        ...rest,
      });
    },

    warning: (messageKey: string, options?: ToastOptions) => {
      const { ...rest } = options || {};
      return showToast.warning({
        message: messageKey,

        ...rest,
      });
    },

    info: (messageKey: string, options?: ToastOptions) => {
      const { ...rest } = options || {};
      return showToast.info({
        message: messageKey,

        ...rest,
      });
    },
  };
}
