import { skipToken, useQuery } from "@tanstack/react-query";
import { api } from "@/services/httpService";
import { accessRoleKeys } from "./queryKeys";
import type { AccessRole } from "../types";

export const getAccessRole = async (id: number | string, signal?: AbortSignal) => {
  const response = await api.get<AccessRole>(`/acl/roles/${id}`, { signal });
  return response.data;
};

export function useGetAccessRole(id?: number | string) {
  return useQuery({
    queryKey: accessRoleKeys.detail(id!),
    // skipToken جایگزین مدرنِ enabled است: هم کوئری را متوقف می‌کند،
    // هم تایپ id داخل تابع narrow می‌شود، پس دیگر به `id as number | string` نیاز نیست.
    queryFn: id == null ? skipToken : ({ signal }) => getAccessRole(id, signal),
  });
}
