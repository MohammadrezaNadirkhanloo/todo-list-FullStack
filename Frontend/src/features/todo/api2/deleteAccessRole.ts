import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useToast } from "@/hooks/useToast";
import { api } from "@/services/httpService";
import { accessRoleKeys } from "./queryKeys";

export const deleteAccessRole = async (ids: number[]) => {
  // توجه: بعضی proxyها بدنه‌ی DELETE را حذف می‌کنند.
  // اگر بک‌اند پشت nginx/gateway است، POST /acl/roles/bulk-delete امن‌تر است.
  const response = await api.delete<void>("/acl/users/roles", {
    data: { roles: ids },
  });
  return response.data;
};

export function useDeleteAccessRole() {
  const queryClient = useQueryClient();
  const toast = useToast();

  return useMutation({
    mutationFn: deleteAccessRole,
    onSuccess: (_result, ids) => {
      toast.success("toast.success.delete");
      // کش رکوردهای حذف‌شده را پاک می‌کنیم تا داده‌ی مرده باقی نماند.
      ids.forEach((id) =>
        queryClient.removeQueries({ queryKey: accessRoleKeys.detail(id) }),
      );
      return queryClient.invalidateQueries({ queryKey: accessRoleKeys.table() });
    },
    onError: () => {
      toast.error("toast.error.delete");
    },
  });
}
