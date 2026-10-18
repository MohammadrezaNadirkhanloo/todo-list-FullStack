import { api } from "@/services/httpService";
import { useQuery } from "@tanstack/react-query";

export const getListTodos = async () => {
  const response = await api.get("/todos/asdasd");
  return response.data;
};

export function useGetListTodos() {
  const query = useQuery({
    queryKey: ["list-todo"],
    queryFn: () => getListTodos(),
    staleTime: 0,
  });

  return {
    ListTodos: query.data,
    isLoading: query.isLoading,
    isFetching: query.isFetching,
  };
}
