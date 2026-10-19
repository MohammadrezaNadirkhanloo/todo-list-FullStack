import { z } from "zod";

const optionSchema = z.object({
  value: z.string(),
  label: z.string(),
});

export const createTodoSchema = () => {
  return z.object({
    title: z
      .string()
      .trim()
      .min(1, { message: "Task title is required" })
      .min(3, { message: "Task title must be at least 3 characters" })
      .max(100, { message: "Task title must be at most 100 characters" }),

    dueDate: z
      .union([z.string(), z.date()])
      .nullable()
      .optional(),

    category: optionSchema
      .nullable()
      .optional(),


    priority: optionSchema
      .nullable()
       .optional(),


    tags: z
      .string()
      .trim()

  });
};


export type TodoFormInput = z.input<
  ReturnType<typeof createTodoSchema>
>;
export type TodoFormData = z.infer<
  ReturnType<typeof createTodoSchema>
>;

