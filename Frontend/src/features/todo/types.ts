export type Priority = "low" | "medium" | "high";

export type Task = {
  id: number;
  title: string;
  done: boolean;
  dueDate: string | null;
  priority: Priority;
  category: string | null;
  tags: string[];
  createdAt: string;
};

export type ListResponse<T> = {
  data: T[];
  total: number;
  totalPages: number;
  page: number;
  pageSize: number;
};

export type TodoFilters = {
  search: string;
  category: string;
  priority: Priority | "";
};

export type TodoListParams = {
  search?: string;
  category?: string;
  priority?: Priority;
};

export type UpdateTaskInput = {
  title?: string;
  done?: boolean;
  dueDate?: string | null;
  category?: string | null;
  priority?: Priority;
  tags?: string[];
};

export type CreateTaskInput = {
  title: string;
  priority?: Priority;
  category?: string | null;
  dueDate?: string | null;
  tags?: string[];
};