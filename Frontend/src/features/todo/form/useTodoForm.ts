import { zodResolver } from "@hookform/resolvers/zod";
import { useMemo } from "react";
import { useForm } from "react-hook-form";
import {
  createTodoSchema,
  type TodoFormData,
  type TodoFormInput,
} from "./todo.schema";

export function useTodoForm() {
  const schema = useMemo(() => createTodoSchema(), []);

  return useForm<TodoFormInput, any, TodoFormData>({
    resolver: zodResolver(schema),
    defaultValues: {
      title: "",
      dueDate: null,
      category: null,
      priority: null,
      tags: "",
    },
  });
}
