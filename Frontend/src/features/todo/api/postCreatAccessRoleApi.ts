import { useToast } from "@/hooks/useToast";
import { api } from "@/services/httpService";
import { useMutation, useQueryClient } from "@tanstack/react-query";

export const postCreatTodo = async (data: any) => {
  const response = await api.post("todos", data);
  return response.data;
};

export function useCreateTodo() {
  const toast = useToast();
  const queryClient = useQueryClient();
  const { mutate, isPending } = useMutation({
    mutationFn: postCreatTodo,
    onSuccess: () => {
      queryClient.refetchQueries({ queryKey: ["list-todo"] });
      toast.success("toast.success.creat");
    },
  });

  return {
    CreatTodo: mutate,
    isLoadingSend: isPending,
  };
}
