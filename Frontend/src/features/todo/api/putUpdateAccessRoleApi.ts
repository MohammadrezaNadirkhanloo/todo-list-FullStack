import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useToast } from "@/hooks/useToast";
import { api } from "@/services/httpService";
import type { AccessRoleFormData } from "../form/accessRole.schema";

export const putUpdateAccessRole = async ({
  id,
  data,
}: {
  id: number | string;
  data: Partial<AccessRoleFormData>;
}) => {
  const response = await api.put(`acl/roles/${id}`, data);
  return response.data;
};

export function useUpdateAccessRole() {
  const queryClient = useQueryClient();
  const toast = useToast();
  const { mutate, isPending } = useMutation({
    mutationFn: putUpdateAccessRole,
    onSuccess: () => {
      queryClient.refetchQueries({ queryKey: ["table-query", "/acl/roles"] });
      toast.success("toast.success.updated");
    },
  });

  return {
    UpdateAccessRole: mutate,
    isLoadingSend: isPending,
  };
}
