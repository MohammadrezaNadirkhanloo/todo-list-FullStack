import { api } from "@/services/httpService";
import { useQuery } from "@tanstack/react-query";
import type { AccessRole } from "../types";

export const getAccessRole = async (id: number | string) => {
  const response = await api.get<AccessRole>(`acl/roles/${id}`);
  return response.data;
};

export function useGetAccessRole(id?: number | string) {
  const query = useQuery({
    queryKey: ["AccessRole-detail", id],
    queryFn: () => getAccessRole(id as number | string),
    enabled: id !== undefined,
    staleTime: 0,
  });

  return {
    AccessRole: query.data,
    isLoading: query.isLoading,
    isFetching: query.isFetching,
  };
}
