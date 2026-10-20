import { useMutation, useQueryClient } from "@tanstack/react-query";

import { api } from "@/services/httpService";

import { todoKeys } from "./keys";

export const deleteTodo = async (id: number) => {
  await api.delete(`/todos/${id}`);
};

export function useDeleteTodo() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: deleteTodo,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: todoKeys.all }),
  });
}