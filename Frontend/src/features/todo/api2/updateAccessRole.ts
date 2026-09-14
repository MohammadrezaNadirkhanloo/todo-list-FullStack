import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useToast } from "@/hooks/useToast";
import { api } from "@/services/httpService";
import { accessRoleKeys } from "./queryKeys";
import type { AccessRoleFormData } from "../form/accessRole.schema";
import type { AccessRole } from "../types";

export interface UpdateAccessRoleVars {
  id: number | string;
  data: Partial<AccessRoleFormData>;
}

export const updateAccessRole = async ({ id, data }: UpdateAccessRoleVars) => {
  const response = await api.put<AccessRole>(`/acl/roles/${id}`, data);
  return response.data;
};

export function useUpdateAccessRole() {
  const queryClient = useQueryClient();
  const toast = useToast();

  return useMutation({
    mutationFn: updateAccessRole,
    onSuccess: (updated, { id }) => {
      toast.success("toast.success.updated");
      queryClient.setQueryData(accessRoleKeys.detail(id), updated);
      // باگ اصلی کد قبلی: فقط جدول باطل می‌شد و کوئری جزئیات دست‌نخورده
      // می‌ماند، پس صفحه‌ی جزئیات داده‌ی قدیمی نشان می‌داد.
      return Promise.all([
        queryClient.invalidateQueries({ queryKey: accessRoleKeys.table() }),
        queryClient.invalidateQueries({ queryKey: accessRoleKeys.detail(id) }),
      ]);
    },
    onError: () => {
      toast.error("toast.error.updated");
    },
  });
}
