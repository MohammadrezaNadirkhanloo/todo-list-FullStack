import { useToast } from "@/hooks/useToast";
import { api } from "@/services/httpService";
import { useMutation } from "@tanstack/react-query";
import { useNavigate } from "react-router";


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
    onSuccess: async () => {
      toast.success("toast.success.login");
      navigate("/todo-list");
    },
  });

  return {
    LoginUsers: mutate,
    isLoadingSend: isPending,
  };
}