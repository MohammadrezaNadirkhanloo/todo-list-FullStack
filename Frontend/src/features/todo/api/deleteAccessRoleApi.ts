import { useToast } from "@/hooks/useToast";
import { api } from "@/services/httpService";
import { useMutation, useQueryClient } from "@tanstack/react-query";

export const deleteAccessRole = async (ids: number[]) => {
  const response = await api.delete("/acl/users/roles", {
    data: {
      roles: ids,
    },
  });

  return response.data;
};

export function useDeleteAccessRole() {
  const queryClient = useQueryClient();
  const toast = useToast();
  const { mutate, isPending } = useMutation({
    mutationFn: deleteAccessRole,
    onSuccess: () => {
      queryClient.refetchQueries({ queryKey: ["table-query", "/acl/roles"] }); //invalidateQueries با جزیات بالا ادرس بدم
      toast.success("toast.success.delete");
    },
  });

  return {
    DeleteAccessRole: mutate,
    isLoadingSend: isPending,
  };
}
