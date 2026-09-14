import { api } from "@/services/httpService";
import { useMutation } from "@tanstack/react-query";
import { useNavigate, useSearchParams } from "react-router";
import { useToast } from "@/hooks/useToast";


interface LoginResponse {
  token: string;
}

export const postLoginUsers = async (data: any) => {
  const response = await api.post<LoginResponse>("auth/login", data);
  return response.data;
};

export function useCreateUsers() {
  const toast = useToast();
  const navigate = useNavigate();


  const { mutate, isPending } = useMutation({
    mutationFn: postLoginUsers,
    onSuccess: async (data) => {
      localStorage.setItem("token", data.token);
      toast.success("toast.success.login");
      navigate("");
    },
  });

  return {
    LoginUsers: mutate,
    isLoadingSend: isPending,
  };
}