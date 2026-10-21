import { useToast } from "@/hooks/useToast";
import { api } from "@/services/httpService";
import { useMutation } from "@tanstack/react-query";
import { useNavigate } from "react-router";

interface LoginResponse {
  token: string;
}

export const postCreatUsers = async (data: any) => {
  const response = await api.post<LoginResponse>("auth/register", data);
  return response.data;
};

export function useCreateUsers() {
  const toast = useToast();
  const navigate = useNavigate();

  const { mutate, isPending } = useMutation({
    mutationFn: postCreatUsers,
    onSuccess: async () => {
      navigate("/sign-in");
      toast.success("Creat Success");
    },
  });

  return {
    CreatUsers: mutate,
    isLoadingSend: isPending,
  };
}
