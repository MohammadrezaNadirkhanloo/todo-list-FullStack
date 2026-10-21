import React, { useMemo, useState } from "react";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

import { useDeleteTodo } from "./api/deletetodo";
import { useGetListTodos } from "./api/gettodos";
import { useUpdateTodo } from "./api/updatetodo";
import { EditTaskDialog } from "./components/EditTaskDialog";
import { TaskItem } from "./components/TaskItem";
import { TodoFilters } from "./components/TodoFilters";
import { DEFAULT_FILTERS } from "./constants";
import { TodoForm } from "./form/FormTodo";
import { useDebounce } from "./hooks/useDebounce";
import type {
  Task,
  TodoFilters as Filters,
  TodoListParams,
  UpdateTaskInput,
} from "./types";

const EmptyState = ({ message }: { message: string }) => (
  <Card>
    <CardContent className="pt-6">
      <p className="text-center text-muted-foreground italic">{message}</p>
    </CardContent>
  </Card>
);

const TodoList: React.FC = () => {
  const [filters, setFilters] = useState<Filters>(DEFAULT_FILTERS);
  const [editingTask, setEditingTask] = useState<Task | null>(null);

  const debouncedSearch = useDebounce(filters.search);

  const params = useMemo<TodoListParams>(
    () => ({
      search: debouncedSearch.trim() || undefined,
      category: filters.category || undefined,
      priority: filters.priority || undefined,
    }),
    [debouncedSearch, filters.category, filters.priority],
  );

  const hasFilters = Boolean(
    filters.search || filters.category || filters.priority,
  );

  const { ListTodos, isLoading, isFetching } = useGetListTodos(params);
  const updateTodo = useUpdateTodo();
  const deleteTodo = useDeleteTodo();

  const tasks = useMemo(() => ListTodos?.data ?? [], [ListTodos?.data]);

  const handleToggle = (task: Task) =>
    updateTodo.mutate({ id: task.id, done: !task.done });

  const handleDelete = (id: number) => deleteTodo.mutate(id);

  const handleSave = (id: number, input: UpdateTaskInput) =>
    updateTodo.mutate(
      { id, ...input },
      { onSuccess: () => setEditingTask(null) },
    );

  const renderList = () => {
    if (isLoading) return <EmptyState message="Loading tasks..." />;

    if (tasks.length === 0) {
      return (
        <EmptyState
          message={
            hasFilters
              ? "No tasks match your search criteria."
              : "No tasks yet. Add a new task to get started!"
          }
        />
      );
    }

    return (
      <div
        className={`space-y-3 transition-opacity ${
          isFetching ? "opacity-70" : ""
        }`}
      >
        {tasks.map((task) => (
          <TaskItem
            key={task.id}
            task={task}
            onToggle={handleToggle}
            onEdit={setEditingTask}
            onDelete={handleDelete}
          />
        ))}
      </div>
    );
  };

  return (
    <div className="p-6">
      <Card className="mb-6">
        <CardHeader>
          <CardTitle className="text-2xl font-bold text-center">
            Todo List
          </CardTitle>
        </CardHeader>
        <CardContent>
          <TodoForm />
        </CardContent>
      </Card>

      <TodoFilters value={filters} onChange={setFilters} />

      <div className="space-y-3">{renderList()}</div>

      <EditTaskDialog
        task={editingTask}
        isSaving={updateTodo.isPending}
        onClose={() => setEditingTask(null)}
        onSave={handleSave}
      />
    </div>
  );
};

export default TodoList;