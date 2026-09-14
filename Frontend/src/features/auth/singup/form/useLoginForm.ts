// hooks/useLoginForm.ts
import { zodResolver } from "@hookform/resolvers/zod";
import { useMemo } from "react";
import { useForm } from "react-hook-form";
import { createLoginSchema, type LoginFormData } from "./auth.schema";

export function useLoginForm() {
  
  const schema = useMemo(() => createLoginSchema(), []);
  
  return useForm<LoginFormData>({
    resolver: zodResolver(schema),
    defaultValues: {
      username: "",
      password: "",
    },
  });
}