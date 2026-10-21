import type { TodoListParams } from "../types";

export const todoKeys = {
  all: ["todos"] as const,
  list: (params: TodoListParams) => ["todos", "list", params] as const,
};