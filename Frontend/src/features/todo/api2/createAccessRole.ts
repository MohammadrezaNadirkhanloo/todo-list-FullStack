import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useToast } from "@/hooks/useToast";
import { api } from "@/services/httpService";
import { accessRoleKeys } from "./queryKeys";
import type { AccessRoleFormData } from "../form/accessRole.schema";
import type { AccessRole } from "../types";

export const createAccessRole = async (data: AccessRoleFormData) => {
  const response = await api.post<AccessRole>("/acl/roles", data);
  return response.data;
};

export function useCreateAccessRole() {
  const queryClient = useQueryClient();
  const toast = useToast();

  return useMutation({
    mutationFn: createAccessRole,
    // با return کردنِ Promise، مقدار isPending تا پایان تازه‌سازی کش true می‌ماند
    // و مودال زودتر از به‌روز شدن جدول بسته نمی‌شود.
    onSuccess: (created) => {
      toast.success("toast.success.create");
      queryClient.setQueryData(accessRoleKeys.detail(created.id), created);
      return queryClient.invalidateQueries({ queryKey: accessRoleKeys.table() });
    },
    onError: () => {
      toast.error("toast.error.create");
    },
  });
}
