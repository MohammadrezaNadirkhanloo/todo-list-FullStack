import { z } from "zod";

export const createLoginSchema = () => {
  return z.object({
    username: z
      .string()
      .min(1, { message: "validation.username.required" })
      .min(3, { message: "validation.username.minLength" })
      .max(50, { message: "validation.username.maxLength" }),
    password: z
      .string()
      .min(1, { message: "validation.password.required" })
      .min(8, { message: "validation.password.minLength" }),
  });
};

export type LoginFormData = z.infer<ReturnType<typeof createLoginSchema>>;
