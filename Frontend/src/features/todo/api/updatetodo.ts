import { useMutation, useQueryClient } from "@tanstack/react-query";

import { api } from "@/services/httpService";

import type { Task, UpdateTaskInput } from "../types";
import { todoKeys } from "./keys";

export const updateTodo = async ({
  id,
  ...body
}: UpdateTaskInput & { id: number }) => {
  const response = await api.patch<{ result: Task }>(`/todos/${id}`, body);
  return response.data.result;
};

export function useUpdateTodo() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: updateTodo,
    onSuccess: () => queryClient.invalidateQueries({ queryKey: todoKeys.all }),
  });
}