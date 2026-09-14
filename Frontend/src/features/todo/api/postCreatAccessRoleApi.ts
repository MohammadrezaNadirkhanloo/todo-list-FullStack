import { useToast } from "@/hooks/useToast";
import { api } from "@/services/httpService";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import type { AccessRoleFormData } from "../form/accessRole.schema";

export const postCreatAccessRole = async (data: AccessRoleFormData) => {
  const response = await api.post("acl/roles", data);
  return response.data;
};

export function useCreateAccessRole() {
  const toast = useToast();
  const queryClient = useQueryClient();
  const { mutate, isPending } = useMutation({
    mutationFn: postCreatAccessRole,
    onSuccess: () => {
      queryClient.refetchQueries({ queryKey: ["table-query", "/acl/roles"] });
      toast.success("toast.success.creat");
    },
  });

  return {
    CreatAccessRole: mutate,
    isLoadingSend: isPending,
  };
}
