import { keepPreviousData, useQuery } from "@tanstack/react-query";

import { api } from "@/services/httpService";

import type { ListResponse, Task, TodoListParams } from "../types";
import { todoKeys } from "./keys";

export const getListTodos = async (params: TodoListParams) => {
  const response = await api.get<ListResponse<Task>>("/todos", { params });
  return response.data;
};

export function useGetListTodos(params: TodoListParams = {}) {
  const query = useQuery({
    queryKey: todoKeys.list(params),
    queryFn: () => getListTodos(params),
    placeholderData: keepPreviousData,
  });

  return {
    ListTodos: query.data,
    isLoading: query.isLoading,
    isFetching: query.isFetching,
  };
}